# 009 · Balance global + backfill — Plan

_Documenta cómo se implementó (spec escrita después del código, ver `spec.md`)._

## Implementación

1. `internal/usecase/ports.go`: `MonthlySummarizer.SumThisMonth`
   cambia de firma — devuelve `(MonthlySummary, error)` en vez de
   `(egreso, ingreso float64, err error)`. `MonthlySummary{Egreso,
   Ingreso float64, CountEgreso, CountIngreso int}` + método
   `Balance() float64` (`Ingreso - Egreso` — se implementó primero
   `Egreso - Ingreso` literal del primer pedido, corregido en el mismo
   turno cuando el usuario aclaró que lo quería al revés).
2. `internal/infrastructure/notion/summary.go`: `SumThisMonth` cuenta
   filas además de sumar `Amount` (mismo loop de paginación de 005, un
   `CountEgreso++`/`CountIngreso++` junto a la suma). Ahora importa
   `usecase` (mismo patrón ya usado por `groq`: infra puede importar
   usecase, no al revés).
3. `cmd/bot/main.go`: `resumenReply` recibe `usecase.MonthlySummary`
   en vez de dos `float64` sueltos, arma el texto con `Balance`,
   `Egresos (N movimiento(s))`, `Ingresos (N movimiento(s))`.
4. **Backfill (operación única, no código de la app):** query real
   contra Notion (`filter: {"property": "Type", "select": {"is_empty":
   true}}`) para listar las páginas sin `Type` — dio 31. Loop de
   `curl PATCH /v1/pages/{id}` con `{"properties": {"Type": {"select":
   {"name": "Egreso"}}}}` por cada una. Verificado después: la misma
   query vuelve a dar 0 resultados.

## Decisiones

- **`Balance = Ingreso - Egreso`** ("cuánto tengo disponible") — el
  pedido inicial fue literal `Egreso - Ingreso`; al mostrar el
  resultado, el usuario aclaró que lo quería al revés. Corregido en el
  mismo turno.
- **Conteo agregado al mismo `SumThisMonth`**, no una llamada aparte —
  ya se está iterando cada página para sumar `Amount`, contar es
  gratis en el mismo loop.
- **Backfill vía `curl` directo, no un comando/flag de la app** —
  mismo patrón que la limpieza de páginas de prueba en sesiones
  anteriores: una operación de mantenimiento puntual sobre datos
  reales, no algo que el bot deba poder hacer en producción (no hay
  ningún comando de Telegram para "reclasificar todo").
- **Las 31 filas se marcaron todas como `Egreso` sin excepción** —
  algunas (`Prestamo Sara`, `Prestamo Oscarin`, `Préstamo Camilo`) son
  ambiguas por nombre, pero el usuario pidió explícitamente "todos los
  registros" sin condicionar — no se interpretó ni se dejó ninguna
  afuera por duda propia.

## Riesgos

- **Backfill es irreversible sin volver a editar a mano** — no se guardó
  qué filas tenían `Type` vacío antes (la lista quedó en el historial
  de la sesión, no en un archivo aparte). Si alguna de las 31 en
  realidad era un ingreso, hay que corregirla manualmente en Notion.
- **Convención de signo del balance** — ya corregida y confirmada por
  el usuario (`Ingreso - Egreso`), sin riesgo abierto.
