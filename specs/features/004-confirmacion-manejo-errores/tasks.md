# 004 · Confirmación y manejo de errores en el chat — Tareas

_Checklist accionable derivada del `plan.md`. Tareas pequeñas y concretas; marca `[x]` al completarlas._

- [ ] `internal/usecase/errors.go`: tipos de error
      (`ErrUnrecognizedType`, `ErrUnrecognizedPaymentMethod`,
      `ErrUnrecognizedCategory`, `ErrWriteFailed`).
- [ ] `CategoryResolver.ListNames()` agregado a la interfaz e
      implementación real (notion).
- [ ] Orquestador de 001 devuelve estos errores en los casos que
      correspondan.
- [ ] `replyFor(...)` en la capa de Telegram: un `case` por tipo de
      error + caso de éxito con link a Notion.
- [ ] Log de `ErrWriteFailed.Err` completo a stderr, mensaje genérico
      al chat.
- [ ] Tests unitarios de `replyFor` para cada tipo de error + éxito.
- [ ] `go build ./...`, `go vet ./...`, `go test ./...` en verde.
- [ ] Validar contra los criterios de aceptación de `spec.md`.
- [ ] Mover la feature a "Hecho" en `../../constitution/roadmap.md`.
