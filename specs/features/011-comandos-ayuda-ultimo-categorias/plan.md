# 011 · Plan técnico

## Capas tocadas

### `internal/usecase/ports.go`

- `LatestTransaction` struct (display-only, no reusa
  `domain.Transaction` para evitar parsear `Date` de vuelta a
  `time.Time` sin necesidad — `/ultimo` solo muestra, no reescribe):
  `PageID, Name string, Amount float64, Type domain.TransactionType,
  PaymentMethod domain.PaymentMethod, CategoryName, Date string`.
- `LatestTransactionFinder` interface: `FindLatest() (LatestTransaction,
  bool, error)` — el `bool` distingue "no hay ningún movimiento
  todavía" de un error real.
- `CategoryResolver` gana un método nuevo: `Create(name string)
  (categoryID string, err error)`. Va en la misma interfaz que
  `Resolve`/`ListNames` porque las tres operan sobre la misma data
  source `Category` y las implementa el mismo tipo concreto
  (`notion.CategoryResolver`) — no se justifica partir en dos
  interfaces separadas para un método más.

### `internal/infrastructure/notion/latest.go` (nuevo)

- `(*Client) FindLatest()` — `POST /data_sources/{id}/query` con
  `sorts: [{"timestamp":"created_time","direction":"descending"}]`,
  `page_size: 1`. Si `Category` tiene una relación, un segundo `GET
  /pages/{category_page_id}` resuelve el nombre para mostrar (no hay
  forma de pedirle a Notion el título de la relación en la misma query
  de `Expenses`). Es una sola llamada extra, en un comando de uso
  esporádico — no un hot path.

### `internal/infrastructure/notion/category_resolver.go`

- `(*CategoryResolver) Create(name string) (string, error)`: primero
  llama a `Resolve` (reusa el caché) — si ya existe, devuelve ese ID
  sin crear nada (evita duplicados por mayúsculas/espacios). Si no
  existe, `POST /pages` con `parent.data_source_id = r.dataSourceID` y
  `properties.Name.title`. Tras crear, invalida el caché
  (`r.cache = nil`) para que el próximo `Resolve`/`ListNames` la vea
  sin esperar el TTL de 5 min.

### `cmd/bot/main.go`

- Nuevo estado en memoria: `pendingCategoryAdd map[int64]bool`
  (chat_id → "esperando el nombre de la categoría nueva"), mismo
  patrón sin-mutex que `pendingPhotos`/`pendingConversations`.
- `matchesCommand(text, command string) bool` — helper genérico
  (case-insensitive, tolera sufijo `@BotUsername`) que reemplaza la
  lógica que hoy vive ad-hoc en `parseStartCommand` para los comandos
  sin parámetros (`/ayuda`, `/ultimo`, `/categorias`).
  `parseStartCommand` se deja como está (ya tiene su propia lógica y
  tests, no vale la pena refactorizarlo a costa de tocar tests que ya
  pasan).
- `handleAyuda`, `handleUltimo`, `handleCategoriasCommand` (manda los
  botones), `handleCategoriaCallback` (procesa `cat:listar` /
  `cat:agregar`), `handleCategoryAddReply` (toma el siguiente texto
  como nombre y llama a `Create`).
- `handleCallbackQuery` se reparte por prefijo de `cq.Data`:
  `tipo:` (008, sin cambios) vs. `cat:` (nuevo). Un callback
  desconocido solo se responde (`AnswerCallbackQuery`) para no dejar
  el botón "cargando".
- `editMsg` helper (análogo a `send`, ya existente) para no repetir el
  `if err := ...EditMessageText(...); err != nil { log... }` en cada
  handler nuevo.
- Dispatcher: orden de prioridad dentro del `switch` de `Message`:
  1. Foto → cancela `pendingConversations` y `pendingCategoryAdd`.
  2. `isStartCommand` (`/egreso`/`/ingreso`) → cancela
     `pendingCategoryAdd`, arranca/reinicia conversación (como hoy).
  3. `matchesCommand(..., "/ayuda")` / `/ultimo` / `/categorias` →
     cancelan ambos estados pendientes, ejecutan el comando.
  4. `pendingCategoryAdd[chatID]` → `handleCategoryAddReply`.
  5. `pendingConversations[chatID] != nil` → como hoy.
  6. `resumen` → cancela ambos, como hoy.
  7. Texto normal → cancela ambos, como hoy.
- `main()`: agrega `/ayuda`, `/ultimo`, `/categorias` a
  `SetMyCommands`. `runPollingLoop` gana dos parámetros:
  `finder usecase.LatestTransactionFinder` (satisfecho por `writer`,
  igual que `summarizer`) y `pendingCategoryAdd map[int64]bool`.
  `handleCallbackQuery` gana el parámetro `categories
  usecase.CategoryResolver` y `pendingCategoryAdd`.

## Por qué no tocar `Registrar`

`/categorias` no registra transacciones — es gestión de categorías,
ortogonal a `Registrar`. Se pasa `resolver` (ya existe en `main()`,
implementa `usecase.CategoryResolver`) directo a los handlers nuevos,
igual que `summarizer` ya se pasa por separado de `registrar` desde la
feature 005.

## Tests

- `internal/usecase/ports.go`: no hay lógica propia que testear (son
  tipos), pero `fakeCategoryResolver` en
  `register_transaction_test.go` necesita el método `Create` nuevo
  para seguir compilando (interfaz ampliada).
- `internal/infrastructure/notion/latest_test.go`: `httptest.Server`,
  caso con relación de categoría (dos requests: query + get page),
  caso sin categoría, caso sin resultados (`bool = false`).
- `internal/infrastructure/notion/category_resolver_test.go` (nuevo si
  no existe, o se agrega a un archivo de test ya presente si lo hay):
  `Create` con nombre nuevo (crea), `Create` con nombre existente (no
  crea, reusa ID).
- `cmd/bot/main_test.go`: `matchesCommand` (tabla, análoga a
  `TestParseStartCommand`).
- `cmd/bot/reply_test.go`: `ayudaReply` (si se extrae a función,
  análogo a `resumenReply`) y `ultimoReply`.

## Validación manual

Con el bot corriendo (una sola instancia verificada por `ps aux`):
`/ayuda` → texto correcto. `/ultimo` → datos reales del último
movimiento en Notion. `/categorias` → botones → "Listar" → nombres
reales; volver a `/categorias` → "Agregar" → nombre de prueba →
confirmar en Notion que la página nueva existe en `Category`; repetir
con el mismo nombre → confirmar que no se duplicó.
