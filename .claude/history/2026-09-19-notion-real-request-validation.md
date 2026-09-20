# 2026-09-19 — Validación con request real contra Notion

## Contexto

Tras mapear el struct `domain.Transaction` y escribir el esqueleto del
cliente Go, el usuario pidió validar con una request real contra
Notion antes de seguir, en vez de confiar solo en el esquema leído.

## Qué se hizo

Usando las herramientas de Notion MCP (no el cliente Go, que necesita
un token de integración que el usuario todavía no tiene configurado):

1. `notion-query-data-sources` sobre la data source `Category`
   (`e1e6ba98-2ca3-825b-9c56-872d36f73df0`) para conseguir un ID de
   categoría real (`Ahorro`,
   `3b66ba98-2ca3-8062-8da4-c05ef1f7b35e`).
2. `notion-create-pages` con `parent: {"type": "data_source_id",
   "data_source_id": "3e36ba98-2ca3-8396-9ede-07441cc01bf8"}` (la data
   source `Expenses`) y properties: `Name`, `Amount`, `date:Date:start`,
   `Payment Method`, `Notes`, `Category` (relation, como array con el
   ID de página). Página creada:
   `https://app.notion.com/p/3e16ba982ca381468ec7e7740d218d23`.
3. `notion-fetch` de esa misma página para confirmar que todo se
   guardó y se lee de vuelta correcto, incluida la relación `Category`
   resuelta a la página `Ahorro`.

## Resultado

- Confirmado: el parent correcto es `data_source_id`, no
  `database_id` plano — este workspace es multi-source database.
- Confirmado: `Category` como relation acepta el ID de página en un
  array; `Payment Method` acepta el nombre de opción como string;
  `Date` se escribe como `date:Date:start` (formato MCP) — en la REST
  pública equivalente es `{"date": {"start": "..."}}`, ya así en
  `mapTransactionToProperties`.
- Confirmado: `Add to Month` (formula) es de solo lectura, Notion la
  calcula — no se intentó escribir y no hizo falta.
- `internal/infrastructure/notion/client.go` actualizado:
  `apiVersion` a `2025-09-03` (versión de la API pública que soporta
  data sources) y `parent` construido como `{"type": "data_source_id",
  "data_source_id": ...}`.

## Limitación de esta validación

Se hizo vía Notion MCP, que usa auth OAuth de usuario — **no** es lo
mismo que la REST API pública (`api.notion.com/v1`) con un token de
integración (`NOTION_API_KEY`), que es lo que de verdad usa
`client.go`. Sigue pendiente una prueba end-to-end con el cliente Go
real: el usuario necesita crear una integración en
notion.so/my-integrations, compartirla manualmente con la data source
`Expenses`, y poner el token en su `.env` local.

## Limpieza

La página de prueba (`TEST - Finance Tracker Go client`) quedó creada
en la data source `Expenses` real del usuario — no hay tool de MCP
disponible en esta sesión para borrar una página individual (solo
trash de data source completa, que no aplica). El usuario debe
borrarla manualmente desde Notion si no la quiere ahí.

## Estado al cierre de sesión

`go build ./...` y `go vet ./...` limpios con el cliente actualizado.
Mapeo de propiedades validado contra datos reales. Falta: prueba
end-to-end con token de integración real, y borrar la página de
prueba.
