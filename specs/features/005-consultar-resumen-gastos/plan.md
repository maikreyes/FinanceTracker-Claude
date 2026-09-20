# 005 · Consultar resumen de gastos — Plan

_Cómo se implementa lo descrito en `spec.md`. Debe respetar la `constitution/`._

## Enfoque

Calcular el resumen en el cliente Go, sumando `Amount` agrupado por
`Type`, sobre las páginas de `Expenses` cuyo `Date` cae en el mes
calendario actual — vía `POST /v1/data_sources/{id}/query` con filtro
de fecha, no usando la fórmula `Add to Month` (esa fórmula ya calcula
"si es de este mes, `Amount`, si no 0" por fila individual, pero sumar
todas las filas igual requiere traerlas y sumarlas del lado del
cliente — la API de Notion no agrega server-side).

## Implementación

1. `internal/infrastructure/notion/client.go`: nuevo método
   `SumThisMonth(ctx context.Context) (egreso, ingreso float64, err
   error)`.
   - `POST /v1/data_sources/{id}/query` con `filter: {"property":
     "Date", "date": {"on_or_after": "<primer día del mes>"}}`,
     paginando con `start_cursor` hasta `has_more: false`.
   - Por cada página: lee `properties.Amount.number` y
     `properties.Type.select.name`, suma al acumulador que
     corresponda.
2. `internal/usecase/`: función `Resumen(ctx) (egreso, ingreso float64,
   err error)` que llama al método de arriba (usecase depende de una
   interfaz `MonthlySummarizer`, no del cliente concreto — mismo patrón
   que `CategoryResolver` de 001).
3. `internal/infrastructure/telegram/` (o `cmd/bot/main.go`): si el
   texto del mensaje (trim + lower) es exactamente `resumen`, llama a
   `Resumen` en vez de `ParseMessage`, y formatea la respuesta.

## Decisiones

- **Suma client-side, no la fórmula `Add to Month`** — la fórmula ya
  existe y ayuda a nivel de fila/UI de Notion, pero no hay forma de
  pedirle a la API un total agregado; hay que traer las filas igual,
  así que se suma `Amount` directo y se ignora la fórmula para este
  cálculo (más simple, un solo campo fuente de verdad).
- **Filtro por `Date` del mes actual en la query, no traer todo y
  filtrar en Go** — menos datos por la red, más rápido con el tiempo
  a medida que crece el historial.
- **Interfaz `MonthlySummarizer` en usecase** — mismo patrón de puertos
  que `CategoryResolver`, para no acoplar el usecase al cliente
  concreto de Notion.

## Riesgos

- **Paginación** — si el mes tiene muchas filas (poco probable para un
  tracker personal), hay que seguir `next_cursor` correctamente o el
  total queda incompleto. Cubrir con un test que simule 2+ páginas.
- **Zona horaria** — "mes calendario actual" depende de la zona
  horaria del proceso; si el bot corre en un servidor con UTC y el
  usuario está en Colombia (UTC-5), un mensaje cerca de medianoche
  podría contarse en el mes equivocado. Mitigación: usar
  `time.Now().In(bogotaLocation)` explícito, no `time.Now()` a secas.
