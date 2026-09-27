# 012 · Tareas

_Orden sugerido en `plan.md` — cada paso deja `go build/vet/test`
en verde antes de pasar al siguiente._

- [x] `internal/bot/state.go`: `ChatState`, `PendingConversation`
      (movido), `ConversationStep` (movido), `Store` interface.
- [x] `internal/bot/dispatch.go`: `Deps`, `Dispatch`, todos los
      `handle*`/`*Reply` movidos de `cmd/bot/main.go` sin cambiar
      comportamiento (solo cambiar mapas por `Store`). `pendingPhoto`
      ya guarda `file_id` (no bytes) desde este paso, no en uno
      posterior — la descarga se difiere a `handleTipoCallback`.
- [x] Mover los tests correspondientes de `cmd/bot/main_test.go` y
      `cmd/bot/reply_test.go` a `internal/bot/` (`reply_test.go`), más
      tests de integración nuevos de `Dispatch` en `dispatch_test.go`
      (foto cancela conversación pendiente, agregar categoría crea la
      categoría, flujo guiado completo end-to-end, chat no autorizado
      no hace nada) — no estaban en el plan original, se agregaron
      porque el riesgo real del refactor es la unificación de los tres
      mapas en `ChatState`, no las funciones puras ya cubiertas.
- [x] `cmd/bot/main.go` reducido: wiring + `memorystore.New()` +
      `runPollingLoop` llamando a `internalbot.Dispatch`.
- [x] Prueba manual real con `cmd/bot` (long-polling): mensaje de una
      línea, `resumen`, foto sin caption (botones + descarga vía
      `file_id` diferida), conversación `/egreso` completa — todos
      confirmados por el usuario. La interpretación por IA de esa foto
      falló por un `429 rate_limit_exceeded` de Groq (tier gratuito,
      límite de tokens de salida por minuto) — no es una regresión del
      refactor: ocurre después de la descarga vía `file_id` (que sí
      funcionó) y el bot mostró el mensaje de error esperado en vez de
      colgarse, confirmado por el usuario.
- [x] `internal/bot/memorystore/memorystore.go` + test.
- [x] Desviación del plan: `internal/infrastructure/telegram/client.go`
      pasó de `apiBaseURL` (var de paquete) a un campo `baseURL` por
      instancia de `Bot`, con `NewBotWithBaseURL(token, baseURL)`
      exportado — necesario porque los tests de `internal/bot`
      (paquete distinto) no pueden pisar una var no exportada de
      `telegram`. Mismo patrón que ya usa `notion`, adaptado para ser
      cross-package.
- [x] Reorganización al patrón de `maikreyes/Triggo/WebHook`
      (2026-09-26): `internal/` pasó a `pkg/<feature>/{handler,model,
      services}` + `pkg/ports`, con `api/` y `cmd/` como únicos puntos
      de entrada. Las rutas de arriba (`internal/bot`, `cmd/bot`, etc.)
      son las de antes de esa reorganización; ver `AGENT.md`.
- [x] `pkg/store/services/upstash.go` (Get/Set/Delete/Take con TTL,
      `MarkUpdateSeen` con `SET NX`) + test contra un Redis falso sobre
      `httptest.Server`.
- [x] `pkg/telegram/services/webhook.go`: `SetWebhook`/`DeleteWebhook`
      + test.
- [x] `api/webhook.go` + `Handler.WebhookHandler` (validación de secret
      en tiempo constante, decode, dedupe por `update_id`, `Dispatch`,
      siempre 200 salvo secret inválido, método o JSON inválidos) + tests.
- [x] `cmd/setwebhook/main.go` (`-url`, `-delete`, `-drop-pending`).
- [x] `.env.example`: `TELEGRAM_WEBHOOK_SECRET`,
      `UPSTASH_REDIS_REST_URL`, `UPSTASH_REDIS_REST_TOKEN`.
- [x] `.claude/rules/telegram-rules.md`: documentar que
      `getUpdates`/webhook son mutuamente excluyentes (mismo riesgo
      que los `409 Conflict` ya conocidos, pero entre modos).
- [x] Documentación de despliegue: `DEPLOY.md` (Upstash, secreto,
      Vercel, `cmd/setwebhook`, volver a local) y `vercel.json`.
- [x] `go build ./...` (incluye `api/`), `go vet ./...`,
      `go test -race ./...` en verde, `gofmt -l .` vacío.
- [ ] Validar en vivo: desplegar en Vercel con cuenta propia, correr
      `cmd/setwebhook` y probar mensaje, foto con botones y `/egreso`.
      Acción del usuario: nunca se validó contra Vercel ni Upstash
      reales.
- [ ] Mover la feature a "Hecho" en `../../constitution/roadmap.md`
      cuando se valide en vivo (ver la tarea de arriba).
- [x] Entrada de bitácora + actualizar `CLAUDE.md`.
