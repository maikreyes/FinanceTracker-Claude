# Telegram API rules

- Bot token read from `TELEGRAM_BOT_TOKEN` env var only. Never hardcode.
- `pkg/telegram/services/` owns all Telegram Bot API calls.
  The handler talks to it only through `ports.Messenger`.
- Two transports share the same `Handler.Dispatch`: long-polling
  (`go run ./cmd`, local dev, in-memory store) and webhook
  (`api/webhook.go` on Vercel, Upstash store). Don't couple parser or
  handler logic to either.
- **Webhook y long-polling son excluyentes:** con un webhook activo
  `getUpdates` falla. Para volver a desarrollo local hay que borrarlo
  primero (`go run ./cmd/setwebhook -delete`); para volver a producción,
  `go run ./cmd/setwebhook -url https://<proyecto>.vercel.app/api/webhook`.
- **El webhook exige `X-Telegram-Bot-Api-Secret-Token`** igual a
  `TELEGRAM_WEBHOOK_SECRET` (comparación en tiempo constante, fail closed
  si no está seteado) y responde 200 aunque el procesamiento falle, para
  que Telegram no reintente un update que ya registró un movimiento;
  `MarkUpdateSeen` descarta los reenvíos por `update_id`.
- Incoming update text is untrusted input: validate/sanitize before
  passing to the transaction services.
- Respect Telegram rate limits (~30 msg/sec global, 1 msg/sec per chat)
  once send logic exists.
- **Solo una instancia del bot corriendo a la vez por token.** Telegram
  rechaza `getUpdates` con `409 Conflict` si hay más de un consumidor
  para el mismo `TELEGRAM_BOT_TOKEN` — confirmado en vivo el
  2026-09-20. No es un bug del código, es la API; matar la instancia
  vieja antes de arrancar una nueva (`ps aux | grep "[c]md"` para
  verificar, no asumir que un `pkill` anterior funcionó).
- **Whitelist de acceso (`TELEGRAM_ALLOWED_CHAT_IDS`)** — lista de
  `chat_id` separados por coma, parseada una sola vez al arrancar
  (`pkg/config`, `ParseAllowedChatIDs`). Fail closed: vacía o mal seteada = el bot no
  procesa ningún mensaje (no fail open). Un mensaje de un chat no
  autorizado se ignora en silencio — no se responde nada, para no
  confirmarle a un desconocido que el bot existe y funciona.
- **El token va en la URL de cada request** (`/bot<TOKEN>/método`), y
  `net/http` devuelve un `*url.Error` que imprime la URL completa. Todo
  error de red del cliente de Telegram se pasa por `withoutURL` antes de
  envolverlo — sin eso, cualquier timeout escribe el token en los logs.
  Nunca envolver con `%w` el error crudo de `httpClient.Do` en este
  paquete.
- **Estado pendiente por chat (`bot.ChatState`):** un chat espera una
  sola cosa a la vez (`Kind`), vence a los 30 min y se consume con
  `Store.Take` (atómico) — nunca `Get` + `Delete`, para que un doble
  toque o una reentrega de Telegram no registre dos veces. Si `Store.Get`
  falla, el dispatcher avisa y no sigue: continuar con un estado vacío y
  guardarlo pisaría el estado real.
