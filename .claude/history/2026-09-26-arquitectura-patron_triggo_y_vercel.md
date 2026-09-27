# 2026-09-26 — Reorganización al patrón de Triggo/WebHook + despliegue en Vercel

## Contexto

El usuario pidió que todo el código siguiera la arquitectura de
`github.com/maikreyes/Triggo/tree/main/WebHook` y que el proyecto quedara
listo para subirse a Vercel. También aprobó agregar `godotenv` (primera
dependencia externa), lo que de paso resuelve el
"TELEGRAM_BOT_TOKEN no está seteado" al correr `go run` sin exportar `.env`.

## El patrón (leído del repo de referencia)

`api/` (funciones de Vercel), `cmd/main.go` (servidor local con godotenv),
`pkg/config`, `pkg/ports` (interfaces, un archivo por dominio) y
`pkg/<feature>/{handler,model,services}`: `handler` orquesta usando puertos,
`model` son DTOs en subpaquetes, `services` implementa (un archivo por
operación, `services.go` con struct + `NewServices`). Solo `api/` cablea.

## Mapeo aplicado

| Antes | Ahora |
|---|---|
| `internal/domain` + tipos/errores de `usecase` | `pkg/transaction/model/{transaction,register,summary,receipt,failure}` |
| `usecase` (Registrar, ParseMessage) | `pkg/transaction/services` |
| `internal/bot` | `pkg/telegram/handler` (métodos de `Handler`) + `pkg/telegram/model/{update,keyboard,command,chatstate}` |
| `infrastructure/telegram` | `pkg/telegram/services` |
| `infrastructure/notion` (`Client` + `CategoryResolver`) | `pkg/notion/services` (un solo `Services`) |
| `infrastructure/groq` | `pkg/groq/services` |
| `bot/memorystore` | `pkg/store/services` (`Memory` y `Upstash`) |
| `internal/app` | `pkg/config` + `pkg/app` |
| `cmd/bot` | `cmd/main.go` (`go run ./cmd`) |

Desviaciones respecto de Triggo: `pkg/app` existe porque `api/` y `cmd/`
comparten el cableado (Triggo lo repite dentro de `api/`), la config se
valida por modo (`CheckTelegram`, `CheckWebhook`) en vez de `log.Fatal` en
`NewConfig`, y `initApp` usa `sync.Once` en vez de recablear en cada request.
`TransactionType` pasó a `transaction.Type` para no tartamudear con el
nombre del paquete.

## Lo que se agregó para Vercel (cierra la feature 012)

- `api/webhook.go` (`Handler`) + `Handler.WebhookHandler`: secreto en
  `X-Telegram-Bot-Api-Secret-Token` comparado en tiempo constante (falla
  cerrado sin secreto), responde 200 aunque el procesamiento falle (Telegram
  reintenta todo lo que no sea 2xx), y descarta reenvíos con
  `Store.MarkUpdateSeen` (`SET NX EX`) para no registrar un movimiento dos
  veces.
- `pkg/store/services/upstash.go`: REST de Upstash sin SDK; `Take` usa
  `GETDEL`; el TTL de cada clave es lo que le queda al estado.
- `SetWebhook`/`DeleteWebhook`, `cmd/setwebhook`, `vercel.json`
  (`maxDuration` 60), `.env.example`, `DEPLOY.md`.

## Verificación

Se dejó el árbol viejo intacto y verde mientras se construía el nuevo, se
comparó el inventario de tests (ninguno se perdió, solo `TestStore_*` pasó a
`TestMemory_*`) y recién entonces se borró `internal/`. `go build`, `go vet`,
`go test -race ./...` y `gofmt` limpios.

## Pendiente

Nunca se validó contra Vercel ni Upstash reales, ni con el bot corriendo
tras la reorganización. Antes de desplegar conviene un `go run ./cmd` y
probar los flujos de siempre. No se tocó el boilerplate HTTP duplicado de
`notion`.
