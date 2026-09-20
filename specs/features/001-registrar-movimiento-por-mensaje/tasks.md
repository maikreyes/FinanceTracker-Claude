# 001 · Registrar movimiento por mensaje de Telegram — Tareas

_Checklist accionable derivada del `plan.md`. Tareas pequeñas y concretas; marca `[x]` al completarlas._

- [x] `internal/usecase/parse_transaction.go`: `ParseMessage` con
      tokenización, `Tipo`, `Concepto`, `Monto`, `MedioDePago`.
- [x] Tabla de palabras clave de `Payment Method` (`efectivo`, `tarjeta
      credito`/`tarjeta de credito`, `tarjeta debito`/`tarjeta de
      debito`, `banco`, `tranferencia`/`transferencia`).
- [x] `internal/usecase/ports.go`: interfaz `CategoryResolver`.
- [x] `internal/infrastructure/notion/category_resolver.go`:
      implementación real contra la data source `Category`, con cache
      en memoria con TTL.
- [x] Matching de `Categoria`: implementado como "lo que sobra después
      de consumir el medio de pago", no como matching contra una lista
      — funciona para `Seguridad Social` (2 palabras) sin caso especial
      (test `TestParseMessage_SeguridadSocialDosPalabras`).
- [x] Orquestación en usecase (`Registrar.Register`): parse → resolver
      categoría → `notion.Client.CreateExpense` (vía interfaz
      `ExpenseWriter`).
- [x] Manejo de errores: `ErrUnrecognizedType`,
      `ErrUnrecognizedPaymentMethod`, `ErrUnrecognizedCategory`,
      `ErrWriteFailed`, `ErrInvalidFormat` — no crea página en ningún
      caso de error.
- [x] Tests unitarios de `ParseMessage` con los 7 ejemplos reales de
      `spec.md` + casos de error + tests de `Registrar` con fakes.
- [x] `go build ./...`, `go vet ./...`, `go test ./...` en verde.
- [x] Validar contra los criterios de aceptación de `spec.md` — además
      de los tests, se corrió un smoke test real end-to-end
      (`Registrar` + `notion.Client` + `notion.CategoryResolver` reales)
      contra el Notion del usuario: categoría `comida` resuelta
      correctamente, página creada, luego archivada.
- [x] Mover la feature a "Hecho" en `../../constitution/roadmap.md`.

## Fuera de esta feature (bloquea la integración completa, no el parser)

- [ ] Handler de Telegram (`internal/infrastructure/telegram/`) que
      invoque el usecase — feature 002.
- [ ] Restricción de quién puede usar el bot (chat_id whitelist) —
      feature 003.
