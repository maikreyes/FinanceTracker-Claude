# 2026-09-20 — Feature 012 (paso 1): internal/bot + memorystore

## Contexto

Con 001-011 completas, se especificó la feature 012 (migración a una
arquitectura compatible con Vercel: webhook + estado externalizado en
Redis). El usuario pidió arrancar por el primer paso del plan:
extraer la lógica de despacho a `internal/bot`, sin tocar todavía
Upstash/webhook.

## Qué se implementó

- `internal/bot/state.go`: `ChatState` (unifica los tres mapas que
  antes vivían en `cmd/bot` — `pendingPhotos`, `pendingConversations`,
  `pendingCategoryAdd` — en un solo struct por `chat_id`, mutuamente
  excluyentes por diseño del dispatcher, igual que antes),
  `PendingConversation`/`ConversationStep` (movidos), `Store`
  interface (`Get`/`Set`/`Delete`).
- `internal/bot/dispatch.go`: `Deps` + `Dispatch(ctx, deps, update)` —
  toda la lógica que antes vivía inline en `runPollingLoop` y en los
  `handle*`/`*Reply` de `cmd/bot`, movida tal cual (mismo orden de
  prioridad, mismas cancelaciones de estado pendiente). Cambio de
  diseño intencional incluido en este mismo paso (no diferido): la
  foto sin caption ahora guarda solo el `file_id` de Telegram en
  `ChatState` — la descarga real (`GetFile`+`DownloadFile`) se difiere
  a `handleTipoCallback`, para no meter bytes de imagen en un store
  externo más adelante.
- `internal/bot/memorystore/memorystore.go`: `Store` en memoria de
  proceso (mismo patrón sin-mutex de siempre), la que sigue usando
  `cmd/bot/main.go` en desarrollo local.
- `cmd/bot/main.go` reducido a wiring + `runPollingLoop`, que ahora
  solo trae updates y llama a `internalbot.Dispatch` — cero lógica de
  negocio propia.
- `internal/infrastructure/telegram/client.go`: desviación del plan
  original — `apiBaseURL` (var de paquete) pasó a ser un campo
  `baseURL` por instancia de `Bot`, con `NewBotWithBaseURL(token,
  baseURL)` exportado. Necesario porque los tests de `internal/bot`
  viven en otro paquete y no pueden pisar una var no exportada de
  `telegram` como si hacía `notion` con sus propios tests internos.
- Tests: funciones puras migradas tal cual
  (`isAllowed`/`matchesCommand`/`parseStartCommand`/
  `tipoCallbackData`), `successReply`/`errorReply`/`resumenReply`/
  `ultimoReply` migrados, más tests de integración nuevos de
  `Dispatch` con fakes de `usecase` y un servidor `httptest` que
  simula la API de Telegram: foto cancela conversación pendiente,
  agregar categoría crea la categoría, flujo guiado `/egreso` completo
  end-to-end, chat no autorizado no hace nada. `memorystore` tiene su
  propio test de `Get`/`Set`/`Delete`.

## Validación real (con el bot corriendo)

Una sola instancia confirmada por `ps aux`. Con `cmd/bot`
(long-polling, `memorystore`): mensaje de una línea, `resumen`, foto
sin caption (aparecieron los botones, la descarga diferida por
`file_id` funcionó) y conversación `/egreso` completa — los cuatro
confirmados por el usuario sin problemas. La interpretación por IA de
esa foto falló con `429 rate_limit_exceeded` de Groq (tier gratuito,
límite de tokens de salida por minuto — no relacionado con este
refactor: ocurre después de la descarga exitosa vía `file_id`); el bot
respondió con el mensaje de error esperado en vez de colgarse,
confirmado por el usuario. Bot detenido limpio, cero instancias
verificado.

## Estado al cierre de este paso

- `go build ./...`, `go vet ./...`, `go test ./...` limpios, `gofmt`
  corrido.
- Feature 012 marcada **en progreso 🚧** en `spec.md`/`tasks.md`/
  `roadmap.md` — paso 1 del plan completo, faltan: store de Upstash
  Redis, `SetWebhook`/`DeleteWebhook` en `telegram.Bot`, `api/webhook.go`,
  `cmd/setwebhook`, `.env.example`, documentación de despliegue.
- Sin instancias del bot corriendo.
