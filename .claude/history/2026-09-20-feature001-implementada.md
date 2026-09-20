# 2026-09-20 — Feature 001 implementada: registrar movimiento

## Contexto

Con las 5 specs listas (entrada anterior), el usuario confirmó ("ok")
empezar a codear por la 001.

## Qué se implementó

- `internal/usecase/errors.go` — errores tipados:
  `ErrInvalidFormat`, `ErrUnrecognizedType`,
  `ErrUnrecognizedPaymentMethod`, `ErrUnrecognizedCategory`,
  `ErrWriteFailed` (envuelve error real de Notion, con `Unwrap`).
- `internal/usecase/ports.go` — interfaces `CategoryResolver` y
  `ExpenseWriter`, para que el usecase no dependa de tipos concretos de
  `infrastructure/notion` (regla de dependencia de `AGENT.md`).
- `internal/usecase/parse_transaction.go` — `ParseMessage(text string,
  now time.Time) (ParsedMessage, error)`. Algoritmo:
  1. Token 0 → `Tipo` contra `typeKeywords`.
  2. Primer token parseable como número → `Monto`; todo lo de en medio
     → `Concepto` (vacío → `Tipo` capitalizado, ya viene así del enum).
  3. Después del monto: `matchPaymentMethod` prueba frases de
     `paymentMethodPhrases` de 3, 2 y 1 palabra, en ese orden (para que
     "tarjeta de credito" no se corte en "tarjeta").
  4. Lo que sobra tras consumir el medio de pago = `CategoryName`
     cruda, sin resolver — **decisión de diseño distinta a la del
     `plan.md` original** (que proponía matchear la categoría contra
     una lista conocida desde el final): en vez de eso, se consume el
     medio de pago primero (vocabulario fijo y chico) y todo el resto
     se toma como categoría cruda. Mismo comportamiento observable,
     pero `ParseMessage` no necesita conocer la lista de categorías de
     Notion — se mantiene pura, sin I/O. Cubre `Seguridad Social` (2
     palabras) sin caso especial.
- `internal/usecase/register_transaction.go` — `Registrar.Register`:
  parsea, resuelve categoría (si vino) contra `CategoryResolver`, llama
  `ExpenseWriter.CreateExpense`, mapea errores.
- `internal/infrastructure/notion/category_resolver.go` —
  `CategoryResolver` real: `POST /v1/data_sources/{id}/query` contra la
  data source `Category`, con paginación (`has_more`/`next_cursor`) y
  cache en memoria de 5 min (`sync.Mutex` + TTL).
- Nueva variable de entorno `NOTION_CATEGORY_DATA_SOURCE_ID` (valor
  real: `e1e6ba98-2ca3-825b-9c56-872d36f73df0`) — agregada a
  `.env.example` y al `.env` local real, documentada en
  `.claude/rules/notion-rules.md`.
- Tests: `parse_transaction_test.go` (7 ejemplos reales de `spec.md` +
  caso `Seguridad Social` + 4 casos de error) y
  `register_transaction_test.go` (orquestación con fakes: resuelve
  categoría, sin categoría, categoría desconocida, falla de escritura).

## Validación real end-to-end

Además de los tests unitarios (con fakes), se corrió un smoke test
temporal (`cmd/register_smoketest/`, borrado después) usando
`Registrar` + `notion.Client` + `notion.CategoryResolver` **reales**
contra el Notion del usuario: `"egreso TEST mecato 13000 efectivo
comida"` e `"ingreso 3000000 banco"` — ambas páginas creadas
correctamente (la de "comida" con la relación a Category resuelta de
verdad), luego archivadas (`PATCH .../pages/{id}` `archived: true`).

## Estado al cierre de sesión

- `go build ./...`, `go vet ./...`, `go test ./...` en verde.
- Feature 001 marcada **implementado ✅** en su `spec.md`, tareas
  marcadas `[x]` en `tasks.md` (salvo las 2 explícitamente fuera de
  esta feature: handler de Telegram y whitelist, que son 002/003).
  Movida a "Hecho" en `specs/constitution/roadmap.md`.
- 001 es código Go funcional e invocable, pero **sin transporte
  todavía** — nadie puede hablarle al bot por Telegram hasta 002.
  Siguiente paso natural: implementar 002.
