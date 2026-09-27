# 2026-09-20 — Feature 008 implementada: preguntar Tipo con botones

## Contexto

Con 001-007 completas y validadas, el usuario pidió una feature nueva:
que cada vez que se mande una foto sin caption, el bot pregunte si es
un ingreso o un egreso — pero con respuestas ya asignadas (botones),
sin que el usuario tenga que escribir nada. Reemplaza la suposición
fija `Tipo: Egreso` que 007 usaba para el camino de foto sin caption.

Spec corta escrita antes de codear (mismo patrón del resto de la
sesión), como `specs/features/008-preguntar-tipo-con-botones/`.

## Qué se implementó

- `internal/infrastructure/telegram/client.go`:
  - `Update.CallbackQuery *CallbackQuery` — el update que Telegram
    manda al tocar un botón inline.
  - `Message` ganó `MessageID` (hacía falta para poder editar el
    mensaje de la pregunta después).
  - `Bot.SendKeyboard` — manda un mensaje con botones inline
    (`reply_markup.inline_keyboard`), devuelve el `message_id`.
  - `Bot.AnswerCallbackQuery` — obligatorio llamarlo o el botón queda
    "cargando" en el cliente del usuario indefinidamente.
  - `Bot.EditMessageText` — para actualizar la pregunta una vez
    contestada.
- `internal/usecase/register_transaction.go`: `RegisterFromPhotoAI`
  gana el parámetro `txType domain.TransactionType` — ya no fija
  `Egreso` a mano, usa el que le pasan (ahora viene del botón, no de
  una suposición).
- `cmd/bot/main.go`:
  - `pendingPhotos map[int64][]byte` en `main`, pasado por todo el
    loop de polling — guarda la foto de un chat mientras espera la
    respuesta del botón. Sin mutex: el loop procesa updates
    secuencialmente, no hay concurrencia real.
  - `handlePhotoMessage`: sin caption, ya no llama a la IA directo —
    guarda la foto en `pendingPhotos` (pisa cualquier pregunta vieja
    sin contestar de ese chat) y manda los botones.
  - `handleCallbackQuery` nuevo: responde el callback siempre
    (incluso en caminos de error), busca la foto pendiente del chat
    (si no hay, edita el mensaje avisando que ya no aplica), parsea el
    `callback_data` (`tipoCallbackData`/`parseTipoCallback`), edita el
    mensaje con el tipo elegido, y recién ahí llama
    `RegisterFromPhotoAI` con ese tipo.
  - Dispatcher del loop: nueva rama para `update.CallbackQuery != nil`,
    antes de la rama de `Message`.
- Tests: `tipoCallbackData`/`parseTipoCallback` (round-trip + dato
  desconocido) en `cmd/bot/main_test.go`. Los 5 tests de
  `RegisterFromPhotoAI` en `internal/usecase` se ajustaron a la nueva
  firma (parámetro `txType`), sin cambiar lo que verificaban.

## Validación real (con el bot corriendo)

Una sola instancia confirmada. Foto real sin caption → aparecieron los
botones "Ingreso"/"Egreso" en Telegram (confirmado por el usuario).
Tocó "Ingreso" → sin errores en el log → confirmado en Notion vía
query real: `Name: Recarga Nequi, Amount: 120000, Type: Ingreso`. Bot
detenido limpio (`bot detenido`), sin instancias corriendo al cerrar.

## Estado al cierre de sesión

- `go build ./...`, `go vet ./...`, `go test ./...` limpios,
  `gofmt` corrido.
- Feature 008 marcada **implementado ✅**, movida a "Hecho" en
  `specs/constitution/roadmap.md`.
- **Las 8 features especificadas (001-008) están implementadas y
  validadas en vivo.**
