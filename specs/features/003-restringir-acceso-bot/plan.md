# 003 · Restringir acceso al bot — Plan

_Cómo se implementa lo descrito en `spec.md`. Debe respetar la `constitution/`._

## Enfoque

Un chequeo simple en el punto de entrada del loop de `cmd/bot/main.go`
(feature 002), antes de invocar el usecase — no en el usecase mismo,
porque "quién puede hablarle al bot" es una decisión de transporte, no
de negocio.

## Implementación

1. `internal/infrastructure/telegram/` (o `cmd/bot/main.go`
   directamente, ver decisiones): función `isAllowed(chatID int64,
   allowed []int64) bool`.
2. `cmd/bot/main.go`: parsea `TELEGRAM_ALLOWED_CHAT_IDS` (split por
   coma, `strconv.ParseInt` cada uno) una vez al arrancar. Si queda
   vacía, loggea advertencia y el loop igual arranca pero
   `isAllowed` siempre devuelve `false` (fail closed).
3. En el loop de polling (002): antes de despachar el mensaje al
   usecase, chequear `isAllowed(update.Message.Chat.ID, allowed)`. Si
   `false`, `continue` sin responder.

## Decisiones

- **Fail closed sin whitelist configurada** — más seguro que fail
  open; un despliegue mal configurado no debe quedar abierto a
  cualquiera por default.
- **No responder a chats no autorizados** — evita que un desconocido
  confirme que el bot existe y está vivo tocando datos reales.
- **Whitelist estática desde env, no editable en caliente** — un bot
  personal no necesita más que esto; menos superficie de bugs.
- **Chequeo en el transporte (`cmd/bot`), no en el usecase** — el
  usecase de 001 no debería saber nada de Telegram/chat_ids (regla de
  dependencia `infrastructure -> usecase -> domain` de
  `AGENT.md`/`.claude/rules/codigo.md`).

## Riesgos

- **Typo en `TELEGRAM_ALLOWED_CHAT_IDS`** (ID mal copiado) deja al
  usuario real bloqueado — mitigación: el log de advertencia al
  arrancar debe imprimir cuántos IDs cargó (no los valores, para no
  loguear IDs de terceros si algún día hay más de uno) para poder
  diagnosticar rápido.
