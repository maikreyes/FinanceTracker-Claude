# 010 · Comandos /ingreso y /egreso — Plan

_Cómo se implementa lo descrito en `spec.md`. Debe respetar la `constitution/`._

## Enfoque

Estado de conversación en memoria por `chat_id`
(`map[int64]*pendingConversation`), mismo patrón que `pendingPhotos`
de la feature 008 — el loop de polling procesa updates
secuencialmente, no hace falta mutex. El registro final reusa
`resolveCategory` (ya existe en `Registrar`) a través de un método
nuevo, `RegisterFromFields`, para no reparsear un string reconstruido
(evita reintroducir la ambigüedad que el flujo guiado justamente
evita).

## Implementación

1. `internal/usecase/parse_transaction.go`: `ParsePaymentMethod(text
   string) (domain.PaymentMethod, bool)` — envoltorio público de
   `matchPaymentMethod` que exige que **todo** el texto matchee (no
   solo un prefijo), para el paso de medio de pago del wizard.
2. `internal/usecase/ports.go` o `register_transaction.go`:
   `ConversationInput{Type domain.TransactionType, Concept string,
   Amount float64, PaymentMethod domain.PaymentMethod, CategoryName
   string}`.
3. `internal/usecase/register_transaction.go`:
   `Registrar.RegisterFromFields(input ConversationInput)
   (RegisterResult, error)` — arma `domain.Transaction` directo de los
   campos (sin pasar por `ParseMessage`), resuelve categoría con
   `resolveCategory` si vino, `CreateExpense`.
4. `internal/infrastructure/telegram/client.go`: `Bot.SetMyCommands`
   — registra `/egreso`/`/ingreso` en el menú de comandos de Telegram
   (`POST setMyCommands`).
5. `cmd/bot/main.go`:
   - `main()`: llama `bot.SetMyCommands` una vez al arrancar (no
     bloqueante si falla — solo loguea, el bot igual funciona sin el
     menú).
   - `pendingConversations map[int64]*pendingConversation`, mismo
     manejo que `pendingPhotos`.
   - `parseStartCommand(text string) (domain.TransactionType, bool)` —
     reconoce `/egreso`/`/ingreso`, tolera el sufijo `@BotUsername`
     que Telegram agrega en algunos contextos.
   - Dispatcher del loop: orden de prioridad — foto (cancela
     conversación pendiente) → `/egreso`/`/ingreso` (arranca/reinicia
     conversación) → hay conversación pendiente para este chat →
     respuesta a esa conversación → `resumen` (cancela conversación
     pendiente) → texto normal (cancela conversación pendiente,
     mismo camino de 001) → nada.
   - `handleStartConversation`, `handleConversationReply` (switch por
     paso: monto → concepto → medio de pago → categoría).

## Decisiones

- **`RegisterFromFields`, no reconstruir un string y llamar
  `ParseMessage`** — los campos ya vienen separados y validados uno a
  uno; reconstruir un string de una línea reintroduciría exactamente
  la ambigüedad posicional que llevó a este flujo guiado (ej. un
  concepto que empieza con un número, o que contiene una palabra clave
  de medio de pago).
- **Categoría es el único paso que permite reintentar sin perder lo
  ya contestado** — monto y medio de pago tienen un formato
  cerrado/validable de antemano; una categoría inválida es la única
  falla que solo se puede detectar contra Notion en el momento de
  registrar, después de ya haber armado toda la transacción.
- **Cualquier acción "grande" (foto, `resumen`, mensaje de una línea)
  cancela una conversación pendiente** — evita el estado confuso de
  "¿esto es la respuesta a la pregunta vieja o un mensaje nuevo?".
- **`SetMyCommands` no bloquea el arranque si falla** — es una mejora
  de descubribilidad (aparecer en el menú), no algo de lo que dependa
  la funcionalidad; si Telegram lo rechaza, el bot sigue funcionando,
  los comandos igual se reconocen por texto.

## Riesgos

- **Usuario manda un comando desconocido** (ej. `/foo`) mientras no
  hay conversación pendiente — cae al camino de texto normal (001),
  que lo va a rechazar como `ErrUnrecognizedType` con un mensaje claro.
  Aceptado, no hace falta un caso especial para comandos no
  reconocidos.
- **Conversación abandonada a medias** (el usuario nunca contesta el
  último paso) — queda en memoria hasta que el bot se reinicia o el
  usuario manda otra cosa que la cancele. Sin timeout automático,
  aceptado para un bot personal de uso esporádico (mismo criterio que
  `pendingPhotos` en 008).
