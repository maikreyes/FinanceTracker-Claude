# 2026-09-19 — Notion MCP: mapeo del esquema real

## Contexto

Notion MCP recién conectado (`claude mcp add --transport http notion
https://mcp.notion.com/mcp`, autorizado por OAuth). Sin nombre/ID de
base de datos dado por el usuario, se buscó en el workspace
(`notion-search` sin resultados por keyword; `notion-list-recent-pages`
sí encontró la página "Expense Tracker").

## Hallazgo: el schema real NO tiene campo Income/Expense

"Expense Tracker" (bajo "Moik Space") contiene dos data sources:

- **Expenses** (la que usa este proyecto) — data source ID
  `3e36ba98-2ca3-8396-9ede-07441cc01bf8`, página contenedora
  `db36ba98-2ca3-82e8-b43b-01cc3c21c031`.
  - `Name` (title), `Amount` (number, `colombian_peso`), `Date` (date),
    `Category` (**relation**, no select, `limit: 1`, apunta a la data
    source `Category`), `Payment Method` (select: Credit Card/Debit
    Card/Bank/Cash), `Notes` (rich_text), `Receipt` (file), `Add to
    Month` (**formula, solo lectura** — Notion la calcula sola).
- **Category** — data source ID `e1e6ba98-2ca3-825b-9c56-872d36f73df0`.
  `Name` (title), `Expenses` (relation inversa), `Spent This
  Month`/`Total Spent` (rollup, solo lectura).

Esquema completo documentado como contrato de datos en
`.claude/rules/notion-rules.md`.

## Decisiones

- `domain.Transaction` (`internal/domain/transaction.go`) reescrito
  **estrictamente** según este schema: `Name`, `Amount`, `Date`,
  `CategoryID`, `PaymentMethod`, `Notes`, `ReceiptURLs`. Se eliminó el
  campo `Tipo`/`TransactionType` (Income/Expense) de la Fase 1 — no
  existe en el schema real, este tracker es solo de gastos. Si el
  usuario quiere soportar ingresos más adelante, es una decisión de
  producto nueva a definir explícitamente (feature nueva en
  `specs/features/`), no algo que se pueda inferir del schema actual.
- `CategoryID` guarda el ID de página de Notion (no el nombre) porque
  `Category` es una relación, no un select — asignar categoría requiere
  resolver nombre → ID contra la data source `Category` antes de
  escribir. Ese resolver todavía no existe (queda para cuando se
  implemente el usecase real).
- `PaymentMethod` como tipo enum de Go con las 4 opciones exactas del
  `select` de Notion — cualquier otro valor lo rechazaría la API igual,
  se refleja en el tipo para que el error se vea en compilación, no en
  runtime.
- "Add to Month" (formula) deliberadamente fuera del struct y fuera del
  mapeo del cliente — la API de Notion rechaza escrituras en
  propiedades `formula`.
- `internal/infrastructure/notion/client.go`: esqueleto de cliente HTTP
  con `net/http` + `encoding/json` de la stdlib, sin SDK externo (no
  hay SDK oficial de Notion en Go). `NewClient(apiKey, dataSourceID)`,
  `CreateExpense(domain.Transaction) error`, y
  `mapTransactionToProperties` como función única y centralizada de
  mapeo (regla de `notion-rules.md`). Sin reintentos ni rate limiting
  todavía — pendiente, como ya estaba anotado en las reglas.
- Pendiente de validar con una request real: si esta versión de la API
  (`2022-06-28`) acepta `parent.database_id` con el ID de la data
  source para este workspace multi-source, o si hace falta
  `parent.data_source_id` de una versión más nueva. Anotado como TODO
  en el código y en `notion-rules.md`.
- `AGENT.md` y `specs/constitution/tech-stack.md` actualizados para que
  no queden con el struct viejo (`Concepto/Monto/Tipo/Fecha`).

## Estado al cierre de sesión

- `go build ./...` y `go vet ./...` limpios con el struct nuevo y el
  cliente de Notion.
- Falta que el usuario ponga el valor real de `NOTION_DATABASE_ID`
  (`3e36ba98-2ca3-8396-9ede-07441cc01bf8`) y `NOTION_API_KEY` en su
  `.env` local (no versionado) antes de poder probar contra la API real.
- Sin lógica de Telegram todavía — sigue como stub, a la espera de
  indicación explícita (siguiente paso natural tras validar este
  struct).
