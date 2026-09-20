# 002 · Transporte Telegram (long-polling) — Plan

_Cómo se implementa lo descrito en `spec.md`. Debe respetar la `constitution/`._

## Enfoque

Cliente HTTP de Telegram con stdlib (`net/http`), mismo patrón que
`internal/infrastructure/notion/client.go` — sin SDK externo, sin
agregar dependencias a `go.mod` sin avisar (regla de
`.claude/rules/flujo_de_trabajo.md`). Loop de long-polling secuencial
en `cmd/bot/main.go`, que es el único punto de entrada del proceso.

## Implementación

1. `internal/infrastructure/telegram/client.go`: reemplaza el stub
   actual. `Bot` struct con `token string`, `httpClient *http.Client`.
   - `GetUpdates(ctx context.Context, offset int) ([]Update, error)` —
     llama `GET /bot<token>/getUpdates?offset=...&timeout=30` (long
     poll nativo de Telegram, evita sondeo agresivo).
   - `SendMessage(ctx context.Context, chatID int64, text string)
     error` — llama `POST /bot<token>/sendMessage`.
   - Tipos `Update`, `Message` mínimos (solo los campos que se usan:
     `update_id`, `message.chat.id`, `message.text`).
2. `cmd/bot/main.go`: reemplaza el stub actual.
   - Carga `TELEGRAM_BOT_TOKEN`, `NOTION_API_KEY`, `NOTION_DATABASE_ID`
     de env vars (sin librería de `.env` — el proceso espera que ya
     estén exportadas, o se cargan a mano en dev con `source .env`;
     ver riesgo abajo).
   - Wiring: `notion.NewClient` → usecase de 001 → `telegram.Bot`.
   - Loop: `for { updates := bot.GetUpdates(ctx, offset); despacha cada
     mensaje al usecase; offset = último update_id + 1 }`.
   - Apagado limpio con `signal.NotifyContext(os.Interrupt,
     syscall.SIGTERM)`.
3. Backoff simple ante error de `GetUpdates`: espera corta (ej. 2s,
   luego backoff exponencial acotado) antes de reintentar, no termina
   el proceso.

## Decisiones

- **`net/http` puro, sin SDK de Telegram** — consistencia con el
  cliente de Notion ya escrito, cero dependencias nuevas sin avisar.
- **Long-polling, no webhook** — ya documentado como preferencia en
  `.claude/rules/telegram-rules.md`, no requiere infraestructura
  pública.
- **Offset solo en memoria, no persistido** — si el proceso se
  reinicia, puede reprocesar el último mensaje o perder mensajes
  recibidos mientras estaba caído. Aceptado para v1 (bot personal, bajo
  volumen); revisar si se vuelve un problema real.
- **Sin librería `.env`** — el proyecto no tiene una todavía y agregar
  una es una dependencia nueva (avisar antes). En dev, cargar `.env` a
  mano (`set -a && source .env && set +a`) antes de `go run
  ./cmd/bot`, como ya se hizo para los smoke tests de Notion.

## Riesgos

- **Reproceso de mensajes tras un reinicio** (offset no persistido) —
  podría re-registrar un gasto duplicado en Notion. Mitigación v1:
  ninguna automática; el usuario puede borrar el duplicado a mano si
  pasa. Si se vuelve frecuente, se persiste el offset en un archivo
  local (feature aparte).
- **Rate limits de Telegram** si el long-poll falla y se reintenta muy
  rápido — mitigado con backoff.
