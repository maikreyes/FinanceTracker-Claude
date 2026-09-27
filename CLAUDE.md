# CLAUDE.md — Memoria del proyecto

Se lee automático al abrir sesión en esta carpeta.

- **[`AGENT.md`](AGENT.md)** (raíz, nunca se mueve) — manual estático de
  arquitectura, stack y estructura del proyecto.
- **[`.claude/rules/`](.claude/rules/)** — reglas operativas y de
  comportamiento, modulares (comunicación, flujo de trabajo, estilo de
  código, seguridad, bitácora, más las específicas de Telegram/Notion),
  pensadas para cambiar con el tiempo.
- **[`.claude/history/`](.claude/history/)** — una entrada de bitácora
  por archivo, nombrado `YYYY-MM-DD-feature-que_se_realizo.md` (ver
  `.claude/rules/bitacora.md`).
- **[`.claude/plans/`](.claude/plans/)** — planes de trabajo del
  proyecto.
- **[`specs/constitution/`](specs/constitution/)** — misión, roadmap y
  tech-stack del producto (referencia fija).
- **[`specs/features/`](specs/features/)** — una carpeta por historia
  de usuario, con `spec.md` / `plan.md` / `tasks.md`.

Este archivo es la bitácora corta: la última entrada completa + un
índice de todas las anteriores.

## Resumen rápido del proyecto

Finance Tracker: Telegram to Notion. Bot en Go (módulo `finance-tracker`)
que lee mensajes en lenguaje natural desde Telegram, los parsea en
`transaction.Transaction` (gastos e ingresos, calcado del esquema real de
la data source "Expenses" en Notion) y los escribe vía la API de Notion.
Arquitectura por feature con puertos, el patrón de `maikreyes/Triggo/WebHook`:
`api/` (función de Vercel) y `cmd/` (long-polling local) cablean
`pkg/<feature>/{handler,model,services}` a través de `pkg/ports`. Se
despliega en Vercel con estado en Upstash Redis (ver `DEPLOY.md`). Flujo de
trabajo spec-driven: historias de usuario → `specs/features/NNN-nombre-feature/`
antes de tocar código. Ramas/commits/PRs los maneja el usuario, nunca
Claude. Detalle completo en `AGENT.md`.

## Índice de bitácora (más reciente primero)

⭐ = entrada clave, releer antes de tocar esa parte del sistema.

- 2026-09-26 · [Reorganización al patrón Triggo/WebHook + despliegue en Vercel (Upstash, webhook, setwebhook)](.claude/history/2026-09-26-arquitectura-patron_triggo_y_vercel.md) ⭐ (última, ver completa abajo)
- 2026-09-26 · [Correcciones de la revisión thermos sobre 012 (paso 1): token en logs, ChatState, bot desacoplado, monto/fecha](.claude/history/2026-09-26-bot-correcciones_revision_thermos.md) ⭐
- 2026-09-20 · [Feature 012 (paso 1): internal/bot + memorystore, dispatch unificado](.claude/history/2026-09-20-feature012-paso1-internal_bot.md) ⭐
- 2026-09-20 · [Feature 011: comandos /ayuda, /ultimo y /categorias — 001-011 completas](.claude/history/2026-09-20-feature011-comandos_ayuda_ultimo_categorias.md) ⭐
- 2026-09-20 · [Feature 010: comandos /ingreso y /egreso con flujo guiado — 001-010 completas](.claude/history/2026-09-20-feature010-comandos_flujo_guiado.md) ⭐
- 2026-09-20 · [Feature 009: corrección del signo del balance (Ingreso - Egreso, no al revés)](.claude/history/2026-09-20-feature009-correccion_signo_balance.md) ⭐
- 2026-09-20 · [Feature 009: balance global en resumen + backfill de Tipo — 001-009 completas](.claude/history/2026-09-20-feature009-balance_y_backfill.md) ⭐
- 2026-09-20 · [Feature 008 implementada: preguntar Tipo con botones inline — 001-008 completas](.claude/history/2026-09-20-feature008-botones_tipo.md) ⭐
- 2026-09-20 · [Feature 007 validada en vivo: modelo Groq original dado de baja, corregido a qwen/qwen3.8-27b — 001-007 completas](.claude/history/2026-09-20-feature007-validada_en_vivo.md) ⭐
- 2026-09-20 · [Feature 007 implementada (Groq, no Gemini): interpretar recibo con IA, falta validar en vivo](.claude/history/2026-09-20-feature007-cambio_proveedor_groq.md) ⭐
- 2026-09-20 · [Feature 006 implementada: adjuntar foto de recibo, File Upload API validado en vivo](.claude/history/2026-09-20-feature006-implementada.md) ⭐
- 2026-09-20 · [Spec 007: interpretar recibo con IA — Google Gemini, sin caption, depende de 006](.claude/history/2026-09-20-spec007-interpretar_recibo_ia.md) ⭐
- 2026-09-20 · [Feature 005 implementada: resumen de gastos — 001-005 completas](.claude/history/2026-09-20-feature005-implementada.md) ⭐
- 2026-09-20 · [Feature 004 implementada: confirmación y manejo de errores, validada en vivo](.claude/history/2026-09-20-feature004-implementada.md) ⭐
- 2026-09-20 · [Feature 003 implementada: whitelist de acceso, validada en vivo (autorizado + no autorizado)](.claude/history/2026-09-20-feature003-implementada.md) ⭐
- 2026-09-20 · [Spec 006: adjuntar foto de recibo al registro (File Upload API de Notion)](.claude/history/2026-09-20-spec006-adjuntar_recibo_foto.md) ⭐
- 2026-09-20 · [Feature 002 implementada: transporte Telegram, probado con mensajes reales](.claude/history/2026-09-20-feature002-implementada.md) ⭐
- 2026-09-20 · [Feature 001 implementada: ParseMessage, CategoryResolver real, Registrar, validado end-to-end](.claude/history/2026-09-20-feature001-implementada.md) ⭐
- 2026-09-20 · [Specs 002-005: roadmap completo del MVP (transporte Telegram, acceso, errores/confirmación, resumen)](.claude/history/2026-09-20-specs002_005-roadmap_completo.md) ⭐
- 2026-09-20 · [Spec 001: Registrar movimiento por mensaje de Telegram (gramática, mapeo medio de pago, resolución de categoría)](.claude/history/2026-09-20-spec001-registrar_movimiento.md) ⭐
- 2026-09-20 · [Schema real de Notion: agregada propiedad Type (Ingreso/Egreso), domain.Transaction y cliente actualizados](.claude/history/2026-09-20-notion-schema-add-type.md) ⭐
- 2026-09-20 · [Cliente Notion real validado end-to-end: token de integración real, escritura confirmada](.claude/history/2026-09-20-notion-client-real-end-to-end.md) ⭐
- 2026-09-19 · [Validación con request real contra Notion (MCP): parent.data_source_id confirmado, cliente actualizado](.claude/history/2026-09-19-notion-real-request-validation.md) ⭐
- 2026-09-19 · [Notion MCP: mapeo del esquema real de "Expenses" a domain.Transaction + cliente HTTP](.claude/history/2026-09-19-notion-mcp_mapping.md) ⭐
- 2026-09-19 · [specs/: replicar carpeta de FDC/DESARROLLO (constitution/ + features/NNN-nombre-feature/)](.claude/history/2026-09-19-specs-replicar_carpeta_fdc.md) ⭐
- 2026-09-19 · [Documentación: adaptar patrón de FDC/DESARROLLO (AGENT.md, CLAUDE.md, .claude/rules/)](.claude/history/2026-09-19-documentacion-adaptar_patron_fdc.md) ⭐
- 2026-09-19 · [Andamiaje Go + Clean Architecture: go.mod, estructura de directorios, stubs](.claude/history/2026-09-19-setup-andamiaje_go_clean_architecture.md) ⭐

## Última entrada completa

### 2026-09-26 — Reorganización al patrón Triggo/WebHook + Vercel

Todo el código pasó de `domain/usecase/infrastructure` al patrón de
`maikreyes/Triggo/WebHook`: `api/` y `cmd/` como entradas, y
`pkg/<feature>/{handler,model,services}` conectados solo por `pkg/ports`
(`transaction`, `telegram`, `notion`, `groq`, `store`, más `config` y
`app`). Sin cambios de comportamiento: se comparó el inventario de tests
contra el árbol viejo antes de borrarlo.

Se cerró lo que faltaba de la 012: `api/webhook.go` (secreto en tiempo
constante, siempre 200, dedupe por `update_id`), store de Upstash por
REST (`GETDEL` para consumir una vez), `SetWebhook`/`DeleteWebhook`,
`cmd/setwebhook`, `vercel.json`, `.env.example` y `DEPLOY.md`. Se agregó
`godotenv` (aprobada) para que `go run ./cmd` cargue `.env` solo.

Estado: `go build`, `go vet`, `go test -race ./...` y `gofmt` limpios.
**Falta** probar `go run ./cmd` contra Telegram real y desplegar en
Vercel con cuenta propia (nunca se validó contra Vercel ni Upstash
reales). Nuevos comandos: `go run ./cmd`, `go run ./cmd/setwebhook`.
