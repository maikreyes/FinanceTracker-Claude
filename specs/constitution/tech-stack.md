# Tech stack y convenciones

_Cómo está construido el proyecto y las reglas que todo el código debe respetar. Es la referencia técnica que ningún plan de feature debería contradecir._

## Tecnologías

- **Lenguaje:** Go 1.23+
- **Framework / runtime:** ninguno — stdlib de Go, salvo
  `github.com/joho/godotenv` (carga `.env` solo en `cmd/`, aprobada).
- **Base de datos:** ninguna propia — Notion es la base de datos del
  producto, vía su API HTTP. Estado de chat (`ChatState`) en memoria en
  local, en Upstash Redis (REST) en producción.
- **Auth:** token de bot de Telegram + API key de Notion + API key de
  Groq, todos leídos de variables de entorno.
- **Tests:** `go test -race ./...`.
- **Despliegue:** Vercel (función serverless `api/webhook.go`) +
  long-polling local (`cmd/main.go`) — ver `DEPLOY.md`.

## Archivos / módulos clave

Arquitectura por feature con puertos, patrón de
`maikreyes/Triggo/WebHook` (migrada desde Clean Architecture el
2026-09-26, ver `.claude/history/2026-09-26-arquitectura-patron_triggo_y_vercel.md`):

- `api/webhook.go` — `Handler`, función serverless de Vercel. Cablea
  dependencias una vez por instancia (`sync.Once`) con el store de
  Upstash.
- `cmd/main.go` — long-polling para desarrollo local, store en
  memoria. `cmd/setwebhook/` — comando de un solo uso para
  configurar o borrar el webhook.
- `pkg/config/` — `Config` leído del entorno y sus validaciones por
  modo.
- `pkg/ports/` — las interfaces, un archivo por dominio (`telegram.go`,
  `notion.go`, `groq.go`, `store.go`, `transaction.go`). Lo único de lo
  que dependen `handler` y `services`.
- `pkg/transaction/` — lógica de negocio: `model/` (subpaquetes
  `transaction`, `register`, `summary`, `receipt`, `failure`) y
  `services/` (parseo de mensajes y registro de movimientos).
- `pkg/telegram/` — `model/` (`update`, `keyboard`, `command`,
  `chatstate`), `services/` (cliente de la Bot API) y `handler/`
  (`Dispatch`, flujo guiado, botones, `WebhookHandler`).
- `pkg/notion/services/` — cliente de la Notion API (movimientos,
  recibos, resumen, último movimiento, categorías).
- `pkg/groq/services/` — interpreta la foto de un recibo con IA.
- `pkg/store/services/` — `ChatStore` en memoria (local) y en Upstash
  Redis (producción).
- `pkg/app/` — raíz de composición compartida por `api/` y `cmd/`.
- `specs/features/NNN-nombre-feature/` — una carpeta por historia de
  usuario (`spec.md` / `plan.md` / `tasks.md`).

## Comandos

- `go build ./...` — compila todo el módulo
- `go run ./cmd` — corre el bot en local (long-polling, carga `.env`)
- `go run ./cmd/setwebhook -url https://<proyecto>.vercel.app/api/webhook`
  — apunta el webhook de Telegram a Vercel (`-delete` lo borra)
- `go vet ./...` — chequeo estático
- `go test -race ./...` — tests
- `gofmt -l .` — lista archivos mal formateados

## Modelo de datos / dominio

- `transaction.Transaction{Name string, Amount float64, Type Type, Date
  time.Time, CategoryID string, PaymentMethod PaymentMethod, Notes
  string, ReceiptFileIDs []string}` (`pkg/transaction/model/transaction`)
  — calcado 1:1 del esquema real de la data source "Expenses" de Notion
  (ver `.claude/rules/notion-rules.md`). `Amount` siempre positivo; el
  signo lo da `Type`.
- `Type`: `Egreso` | `Ingreso` (propiedad `Type` agregada al schema real
  el 2026-09-20 — originalmente el tracker era solo gastos).
- `PaymentMethod`: `Credit Card` | `Debit Card` | `Bank` | `Cash` (enum
  fijo, igual que el `select` de Notion).
- `CategoryID` referencia una página existente en la data source
  `Category` (relación, no texto libre) — resolver nombre → ID antes de
  crear la página en `Expenses`.
- El mapeo `Transaction` -> propiedades de página de Notion vive
  centralizado en `pkg/notion/services/` (ver
  `.claude/rules/notion-rules.md`) — no se dispersa el mapeo en otros
  paquetes.

## Convenciones

- Dirección de dependencia: `handler` y `services` dependen de
  `pkg/ports` y de `model`, nunca entre sí ni de otra feature; `model`
  solo importa la stdlib y otros `model`; `pkg/ports` solo importa
  `model`. Solo `pkg/app`, `api/` y `cmd/` conocen las implementaciones
  concretas y las cablean.
- Nombre de paquete = nombre de carpeta.
- **Comentarios en código: solo si son necesarios, nunca de contexto de
  tarea.** Por defecto no se comenta. Un comentario solo se justifica
  cuando explica un WHY no obvio (restricción oculta, workaround,
  efecto secundario no evidente) — nunca para decir QUÉ hace el código.
  Se redacta desde el contexto técnico del código mismo, nunca desde la
  spec/tarea que lo originó: nunca escribir "para la spec 001" o
  "agregado para la historia de autenticación".
- Idioma: nombres de dominio (`Transaction`, `Concepto`, `Monto`) en
  español donde ya existen así; identificadores de infraestructura y
  paquetes en inglés (convención idiomática de Go).

## Límites duros

- No crear ramas, no hacer `git commit`, `push` ni Pull Requests — eso
  lo hace el usuario.
- No agregar dependencias externas (`go get`) sin avisar antes y
  explicar por qué.
- No desplegar a Vercel ni correr `cmd/setwebhook` contra el bot real
  sin que el usuario lo pida: toca su Telegram y su Notion reales.
- **Ningún dato sensible va hardcodeado — nunca, ni siquiera como
  default "solo dev".** `TELEGRAM_BOT_TOKEN`, `NOTION_API_KEY`,
  `NOTION_DATABASE_ID`: siempre variable de entorno, nunca un valor
  real en `.go`, `.yaml`, `docker-compose.yml` ni ningún archivo
  versionado. Nombres documentados en `.env.example`; el `.env` real
  nunca se commitea (gitignored).
