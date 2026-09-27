# 009 · Balance global + backfill — Tareas

_Escrita después de implementar (ver spec.md). Marcadas según lo que ya se hizo._

- [x] `internal/usecase/ports.go`: `MonthlySummary` + `Balance()`.
- [x] `internal/infrastructure/notion/summary.go`: cuenta filas por
      tipo además de sumar.
- [x] `internal/infrastructure/notion/summary_test.go`: actualizado a
      la nueva firma, con `Balance()` verificado.
- [x] `cmd/bot/main.go`: `resumenReply` con Balance + conteos.
- [x] `cmd/bot/reply_test.go`: actualizado a la nueva firma.
- [x] Backfill real: 31 páginas sin `Type` → `Egreso`. Verificado
      después: 0 páginas sin `Type` restantes.
- [x] `go build ./...`, `go vet ./...`, `go test ./...` en verde.
- [x] Prueba manual real: `resumen` con el bot corriendo → confirmado
      por el usuario que coincide con el cálculo independiente:
      `Balance: 2338155, Egresos: 2338155 (30 movimientos), Ingresos:
      0 (0 movimientos)`.
- [x] Validar contra los criterios de aceptación de `spec.md`.
- [x] Mover la feature a "Hecho" en `../../constitution/roadmap.md`.

## Hallazgo durante la prueba manual

Al reiniciar el bot tras un corte de red (timeouts de `getUpdates`
recuperados con backoff, no un bug), el offset en memoria volvió a 0 y
reprocesó un mensaje viejo pendiente en la cola de Telegram — mismo
riesgo ya documentado en `specs/features/002-transporte-telegram/plan.md`
("offset solo en memoria"). No afectó la prueba, solo generó un log de
error esperado.
