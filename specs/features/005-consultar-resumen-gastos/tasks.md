# 005 · Consultar resumen de gastos — Tareas

_Checklist accionable derivada del `plan.md`. Tareas pequeñas y concretas; marca `[x]` al completarlas._

- [ ] `notion.Client.SumThisMonth`: query con filtro de fecha +
      paginación + suma por `Type`.
- [ ] Interfaz `MonthlySummarizer` en usecase + implementación real.
- [ ] Manejo explícito de zona horaria (Colombia, `America/Bogota`) al
      calcular "primer día del mes".
- [ ] Detección del mensaje `resumen` en la capa de Telegram +
      formateo de la respuesta.
- [ ] Test de paginación de `SumThisMonth` (simular 2+ páginas).
- [ ] Test de mes sin movimientos (devuelve 0, 0, sin error).
- [ ] `go build ./...`, `go vet ./...`, `go test ./...` en verde.
- [ ] Validar contra los criterios de aceptación de `spec.md`.
- [ ] Mover la feature a "Hecho" en `../../constitution/roadmap.md`.
