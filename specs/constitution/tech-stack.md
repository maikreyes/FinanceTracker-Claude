# Tech stack y convenciones

_Cómo está construido el proyecto y las reglas que todo el código debe respetar. Es la referencia técnica que ningún plan de feature debería contradecir._

## Tecnologías

- **Lenguaje:** Go 1.23+
- **Framework / runtime:** ninguno — stdlib de Go. Librerías cliente para
  Telegram/Notion aún no elegidas ni agregadas (ver "Límites duros").
- **Base de datos:** ninguna propia — Notion es la base de datos del
  producto, vía su API HTTP.
- **Auth:** token de bot de Telegram + API key de Notion, ambos leídos de
  variables de entorno.
- **Tests:** `go test` (stdlib), aún sin tests escritos.
- **Despliegue:** sin definir — solo local por ahora.

## Archivos / módulos clave

- `cmd/bot/` — punto de entrada único, wiring de dependencias.
- `internal/domain/` — entidad `Transaction`.
- `internal/usecase/` — parseo de mensajes en lenguaje natural + orquestación.
- `internal/infrastructure/telegram/` — adaptador de la Telegram Bot API.
- `internal/infrastructure/notion/` — adaptador de la Notion API.
- `specs/features/NNN-nombre-feature/` — una carpeta por historia de
  usuario (`spec.md` / `plan.md` / `tasks.md`).

## Comandos

- `go build ./...` — compila todo el módulo
- `go run ./cmd/bot` — corre el bot en local
- `go vet ./...` — chequeo estático
- `go test ./...` — tests
- `gofmt -l .` — lista archivos mal formateados

## Modelo de datos / dominio

- `domain.Transaction{Name string, Amount float64, Type
  TransactionType, Date time.Time, CategoryID string, PaymentMethod
  PaymentMethod, Notes string, ReceiptURLs []string}` — calcado 1:1 del
  esquema real de la data source "Expenses" de Notion (ver
  `.claude/rules/notion-rules.md`). `Amount` siempre positivo; el signo
  lo da `Type`.
- `TransactionType`: `Egreso` | `Ingreso` (propiedad `Type` agregada al
  schema real el 2026-09-20 — originalmente el tracker era solo
  gastos).
- `PaymentMethod`: `Credit Card` | `Debit Card` | `Bank` | `Cash` (enum
  fijo, igual que el `select` de Notion).
- `CategoryID` referencia una página existente en la data source
  `Category` (relación, no texto libre) — resolver nombre → ID antes de
  crear la página en `Expenses`.
- El mapeo `Transaction` -> propiedades de página de Notion vive
  centralizado en `internal/infrastructure/notion/` (ver
  `.claude/rules/notion-rules.md`) — no se dispersa el mapeo en otros
  paquetes.

## Convenciones

- Dirección de dependencia: `infrastructure -> usecase -> domain`. El
  dominio no depende de nada externo; usecase no importa tipos
  concretos de infraestructura (solo interfaces/ports).
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
- No implementar los clientes reales de Telegram/Notion ni el parser de
  mensajes hasta que se indique explícitamente — hoy son stubs.
- **Ningún dato sensible va hardcodeado — nunca, ni siquiera como
  default "solo dev".** `TELEGRAM_BOT_TOKEN`, `NOTION_API_KEY`,
  `NOTION_DATABASE_ID`: siempre variable de entorno, nunca un valor
  real en `.go`, `.yaml`, `docker-compose.yml` ni ningún archivo
  versionado. Nombres documentados en `.env.example`; el `.env` real
  nunca se commitea (gitignored).
