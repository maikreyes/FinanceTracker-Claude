# 002 · Transporte Telegram (long-polling) — Tareas

_Checklist accionable derivada del `plan.md`. Tareas pequeñas y concretas; marca `[x]` al completarlas._

- [x] `internal/infrastructure/telegram/client.go`: tipos `Update`,
      `Message` mínimos + `Bot.GetUpdates` + `Bot.SendMessage`.
- [x] `cmd/bot/main.go`: wiring completo (notion client → usecase 001 →
      telegram bot) + loop de polling + offset en memoria.
- [x] Backoff ante error de `GetUpdates` (no termina el proceso) —
      probado en vivo: un 409 Conflict real (dos instancias corriendo)
      disparó el backoff y reintentó sin morirse.
- [x] Apagado limpio con `signal.NotifyContext` — probado con `SIGINT`
      real sobre el proceso: log `bot detenido` tras cancelar el
      contexto en medio de un long-poll.
- [x] Ignorar updates sin `message.text` (fotos, stickers, etc.) sin
      crashear.
- [x] Prueba manual real: `go run ./cmd/bot`, dos mensajes reales desde
      Telegram — uno corto (formato inválido, respondió error sin
      crashear) y uno completo (`... 6000 efectivo comida`, registrado
      en Notion con Category/Type/Payment Method correctos).
- [x] `go build ./...`, `go vet ./...` en verde.
- [x] Validar contra los criterios de aceptación de `spec.md`.
- [x] Mover la feature a "Hecho" en `../../constitution/roadmap.md`.

## Hallazgo durante la prueba manual

Corrí dos instancias del bot en paralelo por accidente (un `pkill`
anterior no mató el proceso viejo) — Telegram devuelve `409 Conflict:
terminated by other getUpdates request` cuando hay más de un consumidor
de `getUpdates` para el mismo token. Confirma que el bot **no** soporta
correr más de una instancia a la vez (limitación de la API de Telegram,
no de este código) — anotar como "no hagas" operativo, no como bug.
