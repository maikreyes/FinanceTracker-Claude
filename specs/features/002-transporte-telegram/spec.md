# 002 · Transporte Telegram (long-polling)

**Estado:** propuesta

## Qué hace

Conecta el bot con Telegram de verdad: recibe mensajes de texto por
long-polling, los pasa al usecase de la feature 001 (`ParseMessage` +
orquestación con Notion), y devuelve la respuesta al mismo chat.

Sin esta feature, 001 solo funciona como función Go probada por test —
no hay forma de que el usuario le mande un mensaje real al bot.

## Por qué

Es el pegamento que hace que el producto exista de verdad: Telegram →
parser → Notion, de punta a punta, como dice `agents.md`.

## Criterios de aceptación

- [ ] `go run ./cmd/bot` arranca y empieza a hacer long-polling contra
      la Telegram Bot API usando `TELEGRAM_BOT_TOKEN`.
- [ ] Un mensaje de texto enviado al bot dispara el usecase de la
      feature 001 y el bot responde en el mismo chat (el contenido
      exacto de la respuesta lo define la feature 004).
- [ ] Mensajes que no son texto (foto, sticker, audio, documento) se
      ignoran sin que el proceso crashee.
- [ ] Si `getUpdates` falla temporalmente (red, 5xx de Telegram), el
      bot reintenta con backoff en vez de terminar el proceso.
- [ ] `Ctrl+C` / señal de terminación hace un apagado limpio (deja de
      hacer polling, no dispara panics).

## Fuera de alcance

- Modo webhook (requiere dominio público con HTTPS) — long-polling es
  suficiente para uso personal, decisión ya anotada en
  `.claude/rules/telegram-rules.md`.
- Comandos slash de Telegram (`/start`, `/help`) — nada en los
  requisitos actuales los necesita.
- Persistir el offset de `getUpdates` entre reinicios del proceso — ver
  `plan.md`, riesgo aceptado para v1.
- Quién puede hablarle al bot — eso es la feature 003.
- Qué le responde el bot exactamente — eso es la feature 004.
