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
`domain.Transaction` (gastos, calcado del esquema real de la data source
"Expenses" en Notion) y los escribe vía la API de Notion. Clean
Architecture: `internal/domain` / `internal/usecase` /
`internal/infrastructure/{telegram,notion}`, entrada única en
`cmd/bot/main.go`. Flujo de trabajo spec-driven: historias de usuario →
`specs/features/NNN-nombre-feature/` antes de tocar código. Ramas/commits/
PRs los maneja el usuario, nunca Claude. Detalle completo en `AGENT.md`.

## Índice de bitácora (más reciente primero)

⭐ = entrada clave, releer antes de tocar esa parte del sistema.

- 2026-09-20 · [Feature 001 implementada: ParseMessage, CategoryResolver real, Registrar, validado end-to-end](.claude/history/2026-09-20-feature001-implementada.md) ⭐ (última, ver completa abajo)
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

### 2026-09-20 — Feature 001 implementada: registrar movimiento

Usuario confirmó ("ok") empezar a codear por la 001. Implementado:

- `internal/usecase/errors.go` — errores tipados
  (`ErrInvalidFormat`, `ErrUnrecognizedType`,
  `ErrUnrecognizedPaymentMethod`, `ErrUnrecognizedCategory`,
  `ErrWriteFailed`).
- `internal/usecase/ports.go` — interfaces `CategoryResolver` /
  `ExpenseWriter` (usecase no depende de `infrastructure/notion`).
- `internal/usecase/parse_transaction.go` — `ParseMessage`. **Cambio de
  diseño vs. el `plan.md` original:** en vez de matchear la categoría
  contra una lista conocida desde el final del mensaje, se consume
  primero el medio de pago (vocabulario fijo y chico) y todo lo que
  sobra es `CategoryName` cruda sin resolver — mismo comportamiento
  observable, pero `ParseMessage` queda puro (sin I/O a Notion). Cubre
  `Seguridad Social` (2 palabras) sin caso especial.
- `internal/usecase/register_transaction.go` — `Registrar.Register`
  orquesta parse → resolver categoría → `CreateExpense`.
- `internal/infrastructure/notion/category_resolver.go` —
  `CategoryResolver` real: query paginada a la data source `Category`,
  cache en memoria 5 min.
- Nueva env var `NOTION_CATEGORY_DATA_SOURCE_ID` (valor real
  `e1e6ba98-2ca3-825b-9c56-872d36f73df0`) en `.env`/`.env.example`.
- Tests: 7 ejemplos reales de `spec.md` + casos de error + tests de
  `Registrar` con fakes — todos en verde.
- **Validación real end-to-end** (smoke test temporal, borrado
  después): `Registrar` + `notion.Client` + `notion.CategoryResolver`
  reales contra el Notion del usuario — categoría `comida` resuelta de
  verdad, páginas creadas y archivadas.

Feature 001 marcada **implementado ✅**, movida a "Hecho" en
`specs/constitution/roadmap.md`. Es código Go funcional e invocable,
pero sin transporte — nadie puede hablarle al bot por Telegram hasta
que se implemente 002.

Previo, misma sesión: las 5 features especificadas (002-005) y, antes,
spec 001 + la propiedad `Type` agregada al schema real de Notion (ver
entradas anteriores del índice).

Estado: `go build ./...`, `go vet ./...`, `go test ./...` limpios.
Siguiente paso natural: implementar 002 (transporte Telegram).
