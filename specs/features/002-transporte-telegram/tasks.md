# 002 · Transporte Telegram (long-polling) — Tareas

_Checklist accionable derivada del `plan.md`. Tareas pequeñas y concretas; marca `[x]` al completarlas._

- [ ] `internal/infrastructure/telegram/client.go`: tipos `Update`,
      `Message` mínimos + `Bot.GetUpdates` + `Bot.SendMessage`.
- [ ] `cmd/bot/main.go`: wiring completo (notion client → usecase 001 →
      telegram bot) + loop de polling + offset en memoria.
- [ ] Backoff ante error de `GetUpdates` (no termina el proceso).
- [ ] Apagado limpio con `signal.NotifyContext`.
- [ ] Ignorar updates sin `message.text` (fotos, stickers, etc.) sin
      crashear.
- [ ] Prueba manual real: `go run ./cmd/bot`, mandar un mensaje de la
      feature 001 desde Telegram, confirmar que llega a Notion.
- [ ] `go build ./...`, `go vet ./...` en verde.
- [ ] Validar contra los criterios de aceptación de `spec.md`.
- [ ] Mover la feature a "Hecho" en `../../constitution/roadmap.md`.
