# Telegram API rules

- Bot token read from `TELEGRAM_BOT_TOKEN` env var only. Never hardcode.
- `internal/infrastructure/telegram/` owns all Telegram Bot API calls.
  Usecase layer must not import Telegram SDK types directly.
- Prefer long-polling for local dev; webhook is optional future mode —
  don't couple parser logic to either transport.
- Incoming update text is untrusted input: validate/sanitize before
  passing to usecase parser.
- Respect Telegram rate limits (~30 msg/sec global, 1 msg/sec per chat)
  once send logic exists.
