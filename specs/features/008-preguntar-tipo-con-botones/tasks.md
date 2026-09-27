# 008 · Preguntar Ingreso/Egreso con botones — Tareas

_Checklist accionable derivada del `plan.md`. Tareas pequeñas y concretas; marca `[x]` al completarlas._

- [x] `internal/infrastructure/telegram/client.go`: `CallbackQuery` en
      `Update`, `Bot.SendKeyboard`, `Bot.AnswerCallbackQuery`,
      `Bot.EditMessageText`. `Message` ganó `MessageID`.
- [x] `internal/usecase/register_transaction.go`:
      `RegisterFromPhotoAI` recibe `txType domain.TransactionType`
      como parámetro (ya no fija `Egreso`).
- [x] `cmd/bot/main.go`: `pendingPhotos map[int64][]byte` (creado en
      `main`, pasado por todo el loop), `handlePhotoMessage` bifurca
      sin caption → guarda foto + manda botones (ya no llama a IA
      directo), `handleCallbackQuery` nuevo.
- [x] Dispatcher del loop de polling: rama para
      `update.CallbackQuery != nil` (antes de la rama de `Message`).
- [x] Tests unitarios: `tipoCallbackData`/`parseTipoCallback`
      (round-trip + dato desconocido) en `cmd/bot/main_test.go`.
- [x] Prueba manual real: foto sin caption → aparecieron los botones
      "Ingreso"/"Egreso"; tocar "Ingreso" → interpretado y registrado
      en Notion (`Name: Recarga Nequi, Amount: 120000, Type:
      Ingreso`), sin errores en el log.
- [x] `go build ./...`, `go vet ./...`, `go test ./...` en verde.
- [x] Validar contra los criterios de aceptación de `spec.md`.
- [x] Mover la feature a "Hecho" en `../../constitution/roadmap.md`.
