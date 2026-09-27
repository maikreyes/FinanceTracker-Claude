# 012 · Migración de arquitectura para desplegar en Vercel

**Estado:** código listo, falta validar en vivo 🚧 — reorganizado al
patrón de Triggo (`pkg/<feature>/{handler,model,services}`), con
store de Upstash, `api/webhook.go`, `cmd/setwebhook` y `DEPLOY.md`.
Pendiente: desplegar en Vercel con cuenta propia y probar de punta a
punta (acción del usuario). Las rutas del diseño de abajo son las de
antes de la reorganización; ver `AGENT.md` para las actuales.

## Qué hace

Cambia el bot de un proceso persistente (long-polling contra
`getUpdates`, estado de conversación en memoria) a un diseño
compatible con funciones serverless: Telegram manda cada update por
**webhook** a una función HTTP de Vercel, y el estado de conversación
pendiente (foto esperando botón, wizard de `/egreso`/`/ingreso`,
espera de nombre de categoría nueva) se guarda en un store externo
(Redis vía Upstash) en vez de en mapas de Go en memoria.

El resto del bot (parseo de mensajes, escritura en Notion, IA de
recibos, todos los comandos de 001-011) **no cambia de comportamiento**
— cambia solo cómo llega el update y dónde vive el estado entre
updates.

## Por qué

Confirmado con el usuario (turno anterior): el diseño actual
(long-polling + estado en memoria de proceso) es incompatible con
Vercel. Un despliegue serverless no mantiene un proceso corriendo
indefinidamente ni memoria compartida entre invocaciones — cada
request puede atenderla una instancia distinta, o una instancia fría
que perdió cualquier estado en memoria de la anterior.

## Diseño

### 1. Webhook en vez de long-polling

- `api/webhook.go` — función serverless de Vercel (convención del
  runtime de Go: `func Handler(w http.ResponseWriter, r *http.Request)`
  en un archivo bajo `api/`). Recibe el `POST` que manda Telegram por
  cada update.
- Valida el header `X-Telegram-Bot-Api-Secret-Token` contra
  `TELEGRAM_WEBHOOK_SECRET` — sin esto, cualquiera que adivine la URL
  podría mandar updates falsos al bot.
- Decodifica el body a `telegram.Update` (mismo tipo que ya existe) y
  llama a la lógica de despacho compartida (ver punto 3).
- `cmd/bot/main.go` (long-polling) se mantiene para desarrollo local —
  no todos los entornos de prueba tienen una URL pública HTTPS para
  recibir webhooks. Ambos entry points comparten la misma lógica de
  negocio.

### 2. Estado de conversación externalizado

Hoy son tres mapas en memoria (`pendingPhotos`,
`pendingConversations`, `pendingCategoryAdd`), mutuamente excluyentes
por diseño (activar uno cancela los otros dos). Se unifican en un solo
estado por chat:

```go
type ChatState struct {
    PendingPhotoFileID string
    Conversation       *PendingConversation
    PendingCategoryAdd bool
}
```

Y una interfaz `Store` (`Get`/`Set`/`Delete` por `chatID`) con dos
implementaciones:

- **En memoria** (`internal/bot/memorystore`) — la que ya existe hoy,
  reempaquetada. La sigue usando `cmd/bot/main.go` en desarrollo
  local.
- **Redis vía Upstash** (`internal/infrastructure/upstash`) — cliente
  HTTP propio (mismo patrón que Notion/Telegram/Groq: sin SDK externo,
  `net/http` + JSON) contra la REST API de Upstash
  (`UPSTASH_REDIS_REST_URL` + `UPSTASH_REDIS_REST_TOKEN`). El estado
  se serializa a JSON por `chat_id`, con TTL (ej. 1 hora) para que una
  conversación abandonada no quede colgada para siempre.

### 3. Lógica de despacho compartida

Se extrae el `switch` que hoy vive inline en `runPollingLoop` (más
todos los `handle*`) a un paquete nuevo `internal/bot`, con una
función `Dispatch(ctx, deps, update)` que no sabe si la llamó el
long-poller o el webhook. `cmd/bot/main.go` y `api/webhook.go` arman
sus propias dependencias (con su propio `Store`) y llaman a la misma
función.

### 4. pendingPhoto guarda `file_id`, no los bytes

Hoy `handlePhotoMessage` descarga la foto completa antes de guardarla
en `pendingPhotos`. Con un store externo, guardar bytes de imagen en
Redis es peso innecesario — se guarda solo el `file_id` de Telegram, y
la descarga (`GetFile` + `DownloadFile`) se hace recién cuando llega
la respuesta del botón Ingreso/Egreso. Riesgo aceptado: si el usuario
tarda demasiado en tocar el botón, el `file_id` podría dejar de ser
válido (comportamiento ya implícito hoy, ahora explícito).

### 5. Configurar el webhook

`telegram.Bot` gana `SetWebhook(ctx, url, secretToken string) error` y
`DeleteWebhook(ctx) error`. Se agrega `cmd/setwebhook/main.go`, un
comando de un solo uso (no un servicio) para apuntar el webhook a la
URL de Vercel después de cada despliegue, o borrarlo para volver a
desarrollo local con long-polling.

## Riesgo operativo importante

Telegram no permite `getUpdates` (long-polling) mientras haya un
webhook activo — devuelve error. Una vez que el webhook de producción
esté seteado, correr `go run ./cmd/bot` en local **falla** hasta
borrar el webhook (`cmd/setwebhook -delete` o equivalente). Esto ya
mordió a esta sesión como `409 Conflict` con dos instancias de
long-polling — ahora es la misma clase de conflicto pero entre modos,
no entre instancias. Documentar en `.claude/rules/telegram-rules.md`.

## Criterios de aceptación

- [ ] `internal/bot.Dispatch` reproduce exactamente el comportamiento
      actual del dispatcher de `runPollingLoop` (mismos criterios de
      prioridad y cancelación de estados pendientes ya validados en
      001-011) — verificado con tests, no reimplementado desde cero.
- [ ] `cmd/bot/main.go` (long-polling + store en memoria) sigue
      funcionando igual que hoy para desarrollo local — validado en
      vivo con el bot real, igual que las features anteriores.
- [ ] `api/webhook.go` responde `403` (o similar) si el header del
      secret token no coincide, sin tocar ningún estado ni llamar a
      Notion/Telegram.
- [ ] `api/webhook.go` con el secret correcto decodifica el update y
      produce el mismo resultado que `internal/bot.Dispatch` llamado
      directo — validado con tests de integración locales (no
      requiere Vercel desplegado).
- [ ] El store de Upstash implementa la interfaz `Store` con TTL —
      validado con tests contra un `httptest.Server` que simula la
      REST API de Upstash (no se necesita una cuenta real de Upstash
      para los tests).
- [ ] `pendingPhoto` en cualquiera de los dos stores guarda un
      `file_id` (string corto), nunca los bytes de la imagen.
- [ ] `.env.example` documenta las variables nuevas:
      `TELEGRAM_WEBHOOK_SECRET`, `UPSTASH_REDIS_REST_URL`,
      `UPSTASH_REDIS_REST_TOKEN`.
- [ ] Documentación de despliegue (dónde: por decidir — `AGENT.md`,
      un `DEPLOY.md`, o una regla nueva) con los pasos: crear proyecto
      en Vercel, setear env vars, desplegar, correr
      `cmd/setwebhook` para apuntar el webhook a la URL desplegada.

## Fuera de alcance

- Desplegar realmente a Vercel en esta sesión — requiere que el
  usuario tenga/cree una cuenta de Vercel y (si se usa Upstash) una
  cuenta de Upstash, y configure las env vars ahí. Esta feature deja
  el código listo y documentado; el despliegue en sí es una acción del
  usuario, fuera del alcance de lo que se puede validar
  automáticamente en esta sesión.
- Migrar el estado a cualquier otro backend que no sea Redis/Upstash
  (ej. Vercel Postgres, un KV distinto) — se elige Upstash por ser la
  opción REST-friendly estándar para serverless, documentada por
  Vercel.
- Autenticación/autorización más allá del secret token del webhook y
  la whitelist de `chat_id` ya existente (003).
- Cambiar el modelo de `TELEGRAM_ALLOWED_CHAT_IDS` o cualquier otra
  regla de negocio ya implementada — esta feature es puramente de
  infraestructura/transporte.
