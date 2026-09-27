# Notion API rules

## Contrato de datos (leído vía Notion MCP el 2026-09-19)

Base: página "Expense Tracker" → "Moik Space". Contiene dos data
sources (no databases separadas en el sentido viejo de la API: son
"data sources" de un mismo tipo de página, ver nota de versión de API
más abajo).

### Data source "Expenses" (la que usa este proyecto)

- URL de página: `https://app.notion.com/p/db36ba982ca382e8b43b01cc3c21c031`
- ID de página (database, con guiones): `db36ba98-2ca3-82e8-b43b-01cc3c21c031`
- **ID de data source (usar como `NOTION_DATABASE_ID`):**
  `3e36ba98-2ca3-8396-9ede-07441cc01bf8`

| Propiedad Notion | Tipo Notion | Notas |
|---|---|---|
| `Name` | `title` | Concepto del gasto. Obligatorio (toda página necesita título). |
| `Amount` | `number` (`colombian_peso`) | Monto. Numérico, sin decimales de centavos en la práctica (COP). |
| `Date` | `date` | Fecha del gasto. Formato de visualización `ll`/`LT`, pero la API se escribe en ISO-8601. |
| `Category` | `relation` (a data source `Category`, `limit: 1`) | **No es un select** — es relación a una página existente en la data source `Category` (ver abajo). Escribir requiere el ID de esa página, no un string libre. |
| `Payment Method` | `select` | Opciones fijas: `Credit Card`, `Debit Card`, `Bank`, `Cash`. Cualquier otro valor lo rechaza Notion. |
| `Notes` | `rich_text` | Opcional. |
| `Receipt` | `file` (`files`) | Opcional. Escribir requiere un `file_upload` real (ver sección File Upload API abajo) — no se puede guardar un link externo de Telegram (expira). |
| `Type` | `select` | **Agregada el 2026-09-20** — el schema original no tenía ningún campo de tipo (era solo gastos). Opciones: `Egreso`, `Ingreso`. `Amount` sigue siempre positivo; el signo lo da `Type`, no el número. |
| `Add to Month` | `formula` | **Solo lectura** — Notion la calcula, la API rechaza escrituras en propiedades `formula`. No mapear en el cliente. |

### Data source "Category"

- **ID de data source (usar como `NOTION_CATEGORY_DATA_SOURCE_ID`):**
  `e1e6ba98-2ca3-825b-9c56-872d36f73df0`
- `Name` (`title`) — nombre de la categoría.
- `Expenses` (`relation`, lado inverso de `Category` en Expenses) — no se
  escribe desde este lado.
- `Spent This Month` / `Total Spent` (`rollup`, solo lectura) — no
  mapear.

Para asignar categoría a un gasto nuevo hay que resolver el nombre a un
ID de página existente en esta data source (buscarla, no crearla salvo
que se pida explícitamente) antes de crear la página en `Expenses`.

### Nota de versión de API — validada con request real (2026-09-19)

Este workspace usa "data sources" (feature de Notion de bases
multi-fuente). Validado creando y leyendo de vuelta una página real vía
Notion MCP: `parent: {"type": "data_source_id", "data_source_id":
"3e36ba98-2ca3-8396-9ede-07441cc01bf8"}` funciona — Name, Amount,
Date, Payment Method, Notes y Category (relation) se escribieron y
leyeron correctamente. `Add to Month` se confirmó de solo lectura
(Notion la calcula sola, no acepta escritura).

**Validado end-to-end el 2026-09-20** con el cliente Go real (no MCP):
`Client.CreateExpense` contra la REST API pública, token de integración
real, `apiVersion 2025-09-03`, `parent.data_source_id`. Ver
`.claude/history/2026-09-20-notion-client-real-end-to-end.md`.

**2026-09-20 (más tarde):** ejemplos de mensajes reales del usuario
incluían ingresos ("ingreso 3000000 banco") — el schema original no
podía representarlos. Se agregó la propiedad `Type` (select
Egreso/Ingreso) a la data source `Expenses` real vía
`notion-update-data-source` (`ADD COLUMN "Type" SELECT('Egreso':red,
'Ingreso':green)`), y se re-validó con dos escrituras reales del
cliente Go (una Egreso, una Ingreso). Filas históricas (anteriores a
este cambio) quedan con `Type` vacío — no se hizo backfill masivo, no
fue lo que se pidió. Ver
`.claude/history/2026-09-20-notion-schema-add-type.md`. Sin pendientes
de validación de esquema/cliente.

**2026-09-20 (feature 001 implementada):** `Services.Resolve`
(`pkg/notion/services/categories.go`) resuelve
nombre → ID de página contra esta data source vía `POST
/v1/data_sources/{id}/query`, cacheado en memoria (TTL 5 min).
Validado end-to-end junto con `transaction/services`: mensaje con
categoría real (`comida`) → resuelto → página creada en `Expenses` con
la relación correcta. Página de prueba archivada después.

### File Upload API — implementado el 2026-09-20 (feature 006)

Para adjuntar un archivo real a `Receipt` (no un link externo): dos
pasos.

1. `POST /v1/file_uploads` (body `{}`) → `{"id": ..., "upload_url":
   ...}`.
2. `POST <upload_url>` con `multipart/form-data`, campo `file` con el
   contenido real y su `Content-Type` — queda `status: uploaded`.

El `id` del paso 1 se referencia en la propiedad de la página:
`"Receipt": {"files": [{"type": "file_upload", "file_upload": {"id":
"<id>"}}]}`. Implementado en
`pkg/notion/services/upload_receipt.go`
(`Services.Upload`). Validado end-to-end con una foto real de Telegram:
el archivo queda nativo en Notion (`type: "file"`), no un link.

Límite de tamaño del File Upload API en modo `single_part`: no
confirmado explícitamente todavía (no hizo falta con la foto de prueba
real) — revisar antes de subir archivos grandes.

## Reglas generales

- API key leída solo de `NOTION_API_KEY`. Nunca hardcodear.
- `NOTION_DATABASE_ID` = ID de la data source `Expenses` (ver tabla
  arriba), no el ID de la página contenedora "Expense Tracker".
- `NOTION_CATEGORY_DATA_SOURCE_ID` = ID de la data source `Category`
  (ver tabla arriba) — la usa `Services.Resolve`.
- `pkg/notion/services/` es dueño de todas las llamadas a la
  API de Notion y del mapeo de propiedades. Fuera de ese paquete nadie importa
  el cliente de Notion, solo los puertos de `pkg/ports`.
- `transaction.Transaction` mapea a propiedades de página de Notion en una
  sola función centralizada (`mapTransactionToProperties`) — no
  dispersar el mapeo.
- Rate limit de la API de Notion (~3 req/seg promedio) — batchear o
  limitar una vez exista lógica de escritura en volumen.
- La integración debe compartirse manualmente con la página/data source
  destino desde la UI de Notion antes de que cualquier llamada a la API
  funcione — eso no lo puede hacer el código.
