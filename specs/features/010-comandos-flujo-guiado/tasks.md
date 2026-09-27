# 010 · Comandos /ingreso y /egreso — Tareas

_Checklist accionable derivada del `plan.md`. Tareas pequeñas y concretas; marca `[x]` al completarlas._

- [x] `internal/usecase/parse_transaction.go`: `ParsePaymentMethod`.
- [x] `internal/usecase/register_transaction.go`: `ConversationInput` +
      `Registrar.RegisterFromFields`.
- [x] `internal/infrastructure/telegram/client.go`: `Bot.SetMyCommands`.
- [x] `cmd/bot/main.go`: `pendingConversations`, `parseStartCommand`/
      `isStartCommand`, `handleStartConversation`,
      `handleConversationReply`, dispatcher actualizado con la
      prioridad de `plan.md`.
- [x] `main()`: llama `SetMyCommands` al arrancar (no bloqueante).
- [x] Tests unitarios: `ParsePaymentMethod` (match exacto vs. texto con
      sobrante), `RegisterFromFields` con fakes (completo, sin
      concepto usa default, categoría no reconocida),
      `parseStartCommand`/`isStartCommand` (con `@BotUsername`,
      mayúsculas, sin "/").
- [x] Prueba manual real: `/egreso` (y luego probado con `/ingreso`)
      completo paso a paso — monto, concepto, medio de pago,
      categoría omitida con "-" — registrado en Notion (`Name: Pago
      préstamo, Amount: 120000, Type: Ingreso, Payment Method: Bank`,
      categoría vacía como se pidió).
- [x] `go build ./...`, `go vet ./...`, `go test ./...` en verde.
- [x] Validar contra los criterios de aceptación de `spec.md`.
- [x] Mover la feature a "Hecho" en `../../constitution/roadmap.md`.
