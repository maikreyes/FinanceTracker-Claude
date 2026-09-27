# 2026-09-20 — Feature 009: corrección del signo del balance

## Contexto

Justo después de cerrar la feature 009 (`Balance = Egreso - Ingreso`,
implementado literal del primer pedido y marcado explícitamente en la
spec como "verificar con el usuario"), se le mostró el resultado al
usuario junto con la pregunta de si esa convención era la que quería.
Respondió que no — la quería al revés.

## Qué se corrigió

- `internal/usecase/ports.go`: `MonthlySummary.Balance()` pasa de
  `Egreso - Ingreso` a **`Ingreso - Egreso`** ("cuánto tengo
  disponible", la convención estándar).
- `cmd/bot/main.go`: comentario de `resumenReply` actualizado.
- Tests ajustados: `internal/infrastructure/notion/summary_test.go`
  (`Balance()` esperado pasa de `-1000` a `1000` para los mismos datos
  de prueba) y `cmd/bot/reply_test.go` (ya no busca el signo negativo
  en el texto).
- `specs/features/009-balance-global-resumen/{spec,plan}.md` y
  `specs/constitution/roadmap.md` actualizados para reflejar la
  fórmula final y dejar registrado que hubo una primera versión
  invertida, corregida en el mismo turno (no se reescribió como si
  siempre hubiera sido así).

No hizo falta tocar el backfill de las 31 filas históricas — el signo
del balance es puramente de presentación en `resumen`, no afecta cómo
se guardan los datos en Notion.

## Estado al cierre de sesión

- `go build ./...`, `go vet ./...`, `go test ./...` limpios.
- Con la fórmula corregida, los mismos números validados antes
  (`Egresos: 2338155`, `Ingresos: 0`) dan `Balance: -2338155` — no se
  volvió a correr el bot en vivo para esta corrección puntual de signo
  (cambio de una línea + tests, sin lógica nueva que valga la pena
  repetir la prueba manual completa).
- Feature 009 sigue **implementado ✅** en el roadmap, con la fórmula
  ya corregida y sin pendientes abiertos.
