# 011 · Tareas

- [x] `internal/usecase/ports.go`: `LatestTransaction`,
      `LatestTransactionFinder`, `CategoryResolver.Create`.
- [x] `internal/infrastructure/notion/latest.go`: `Client.FindLatest`.
- [x] `internal/infrastructure/notion/category_resolver.go`:
      `CategoryResolver.Create` (no duplica, invalida caché).
- [x] `cmd/bot/main.go`: `pendingCategoryAdd`, `matchesCommand`,
      `handleAyuda`, `handleUltimo`, `handleCategoriasCommand`,
      `handleCategoriaCallback`, `handleCategoryAddReply`, `editMsg`,
      `handleCallbackQuery` repartido por prefijo, dispatcher
      actualizado, `SetMyCommands` con los 3 comandos nuevos.
- [x] `internal/usecase/register_transaction_test.go`:
      `fakeCategoryResolver.Create`.
- [x] Tests: `latest_test.go`, `Create` en category_resolver,
      `matchesCommand`, `ultimoReply`.
- [x] `go build ./...`, `go vet ./...`, `go test ./...` en verde,
      `gofmt -l .` vacío.
- [x] Prueba manual real con el bot corriendo: `/ayuda` (texto
      correcto), `/categorias` → Listar (categorías reales) y →
      Agregar (categoría nueva creada, confirmado por el usuario),
      `/ultimo` (falló por timeout de red en el primer intento —no es
      un bug de código—, reintento exitoso mostró el movimiento
      correcto). El caso de nombre duplicado en `/categorias` quedó
      cubierto por test unitario
      (`TestCategoryResolver_Create_CategoriaExistenteNoDuplica`), no
      se repitió manualmente contra el Notion real.
- [x] Mover la feature a "Hecho" en `../../constitution/roadmap.md`.
- [x] Entrada de bitácora + actualizar `CLAUDE.md`.
