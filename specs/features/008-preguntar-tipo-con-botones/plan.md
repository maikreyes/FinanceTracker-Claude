# 008 · Preguntar Ingreso/Egreso con botones — Plan

_Cómo se implementa lo descrito en `spec.md`. Debe respetar la `constitution/`._

## Enfoque

Botones inline de Telegram (`inline_keyboard`) + `callback_query` (el
update que Telegram manda cuando se toca un botón). La foto descargada
se guarda en memoria (`map[int64][]byte`, por `chat_id`) mientras se
espera la respuesta — el bot hoy procesa updates secuencialmente (sin
goroutines concurrentes), así que no hace falta mutex.

## Implementación

1. `internal/infrastructure/telegram/client.go`:
   - `Update` gana `CallbackQuery *CallbackQuery`.
   - `CallbackQuery{ID string, Data string, Message *Message}`.
   - `Bot.SendKeyboard(ctx, chatID int64, text string, buttons
     []InlineButton) (messageID int, error)` — `sendMessage` con
     `reply_markup.inline_keyboard`. `InlineButton{Text, CallbackData
     string}`.
   - `Bot.AnswerCallbackQuery(ctx, callbackQueryID string) error` —
     Telegram exige responder o el botón queda "cargando" en el
     cliente del usuario.
   - `Bot.EditMessageText(ctx, chatID int64, messageID int, text
     string) error` — para actualizar el mensaje de la pregunta una
     vez contestada.
2. `cmd/bot/main.go`:
   - `pendingPhotos map[int64][]byte` en `main`, pasado a
     `runPollingLoop`.
   - `handlePhotoMessage`: si `Caption` vacío → guarda `photo` en
     `pendingPhotos[chatID]` (sobreescribe si ya había una), manda los
     botones con `SendKeyboard`, devuelve sin llamar a IA todavía. Con
     caption, sin cambios (006).
   - Dispatcher del loop: `update.CallbackQuery != nil` →
     `handleCallbackQuery`.
   - `handleCallbackQuery`: `AnswerCallbackQuery` primero (siempre,
     para no dejar el botón "cargando" ni en el camino de error);
     busca `pendingPhotos[chatID]` — si no hay, `EditMessageText` con
     "Esta pregunta ya no aplica, mandá la foto de nuevo." y termina;
     si hay, la saca del mapa (`delete`), parsea `Data` (`"tipo:egreso"`
     / `"tipo:ingreso"`) a `domain.TransactionType`, llama
     `Registrar.RegisterFromPhotoAI` (con el tipo elegido, ver abajo),
     `EditMessageText` con qué se eligió, y manda la confirmación
     (`successReply`/`errorReply`) como antes.
3. `internal/usecase/register_transaction.go`:
   `RegisterFromPhotoAI` gana un parámetro `txType domain.TransactionType`
   — ya no fija `Egreso` a mano, usa el que le pasan.

## Decisiones

- **Guardar la foto en memoria por `chat_id`, no un token/UUID por
  pregunta** — más simple; alcanza porque el bot procesa un chat
  autorizado a la vez (whitelist de 1) y el spec ya acepta que una
  pregunta vieja quede obsoleta si llega una foto nueva.
- **Interpretar con IA después de la respuesta, no antes** — evita
  gastar una llamada a Groq si el usuario nunca contesta el botón.
- **`AnswerCallbackQuery` siempre, incluso en el camino de error** —
  si no se llama, el botón del usuario queda con el ícono de "cargando"
  indefinidamente en su Telegram.
- **`RegisterFromPhotoAI` recibe el tipo como parámetro** en vez de
  seguir fijo en `Egreso` — cambio de firma controlado, la única
  responsable de llamarlo hoy es `handleCallbackQuery`.

## Riesgos

- **Bot se reinicia con una pregunta sin responder** — la foto en
  memoria se pierde, el botón que toque el usuario después no
  encuentra nada pendiente (camino ya cubierto: `EditMessageText`
  avisando que ya no aplica). Aceptado, no se persiste a disco.
- **Foto pesada guardada en memoria mientras se espera respuesta** —
  para un bot personal de uso esporádico no es un problema real; si
  se vuelve uno, se podría descargar la foto recién al responder en
  vez de al recibirla (cambiaría el diseño, no se hace ahora).
