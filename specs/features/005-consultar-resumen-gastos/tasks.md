# 005 · Consultar resumen de gastos — Tareas

_Checklist accionable derivada del `plan.md`. Tareas pequeñas y concretas; marca `[x]` al completarlas._

- [x] `notion.Client.SumThisMonth`: query con filtro de fecha (`Date
      on_or_after` primer día del mes) + paginación + suma por `Type`
      (`internal/infrastructure/notion/summary.go`).
- [x] Interfaz `MonthlySummarizer` en usecase (`ports.go`) +
      implementación real (`notion.Client` la satisface directamente).
- [x] Manejo explícito de zona horaria (`America/Bogota`, con
      `_ "time/tzdata"` embebido para no depender del SO) al calcular
      "primer día del mes".
- [x] Detección del mensaje `resumen` (case-insensitive, trim) en
      `cmd/bot/main.go` + `handleResumen`/`resumenReply`.
- [x] Test de paginación de `SumThisMonth` (`summary_test.go`, dos
      páginas simuladas con `httptest.Server`) — se hizo `apiBaseURL`
      configurable (`const` → `var`) para poder apuntarlo al server de
      prueba.
- [x] Test de mes sin movimientos (devuelve 0, 0, sin error) + filas
      sin `Type` (históricas) correctamente excluidas de la suma.
- [x] `go build ./...`, `go vet ./...`, `go test ./...` en verde.
- [x] Validar contra los criterios de aceptación de `spec.md` —
      probado con el bot corriendo de verdad: `resumen` respondió
      `Egresos: 0, Ingresos: 0`. Verificado independientemente contra
      Notion (query directa) que ese es el resultado correcto: las 30
      filas del mes son históricas (sin `Type`, de antes de esta
      sesión) y las páginas de prueba ("pizza") que se habían creado
      durante las pruebas de 002-004 ya estaban archivadas por el
      usuario — no quedaba ningún movimiento real de septiembre con
      `Type` seteado al momento de la prueba.
- [x] Mover la feature a "Hecho" en `../../constitution/roadmap.md`.
