# 004 · Confirmación y manejo de errores en el chat — Tareas

_Checklist accionable derivada del `plan.md`. Tareas pequeñas y concretas; marca `[x]` al completarlas._

- [x] `internal/usecase/errors.go`: tipos de error
      (`ErrUnrecognizedType`, `ErrUnrecognizedPaymentMethod`,
      `ErrUnrecognizedCategory` con campo `Valid []string`,
      `ErrWriteFailed`) — ya existían desde 001, se extendió
      `ErrUnrecognizedCategory` con la lista de categorías válidas.
- [x] `CategoryResolver.ListNames()` agregado a la interfaz e
      implementación real (`notion.CategoryResolver`, con casing
      original preservado y orden alfabético).
- [x] `Registrar.Register` devuelve `RegisterResult` (page ID +
      transacción + nombre de categoría) en vez de solo el page ID, y
      arma `ErrUnrecognizedCategory` con `Valid` poblado (best effort).
- [x] `cmd/bot/main.go`: `successReply` + `errorReply` — un `case` por
      tipo de error + caso de éxito con link a Notion.
      `usecase.ValidTypes()`/`ValidPaymentMethods()` nuevos, para no
      hardcodear las listas en dos lugares.
- [x] `ErrWriteFailed` se loguea completo a stderr, mensaje genérico al
      chat (`cmd/bot/reply_test.go` verifica que no se filtre).
- [x] Tests unitarios de `errorReply`/`successReply` para cada tipo de
      error + éxito (`cmd/bot/reply_test.go`) + tests actualizados de
      `Registrar` (`internal/usecase/register_transaction_test.go`).
- [x] `go build ./...`, `go vet ./...`, `go test ./...` en verde.
- [x] Validar contra los criterios de aceptación de `spec.md` —
      probado con el bot corriendo de verdad: mensaje con medio de
      pago inválido (`bitcoin`) respondido con el error correcto
      (`medio de pago no reconocido: "bitcoin"`), y mensaje válido con
      categoría registrado y confirmado en Notion.
- [x] Mover la feature a "Hecho" en `../../constitution/roadmap.md`.
