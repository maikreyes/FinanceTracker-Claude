# 012 · Plan técnico

## Paquete nuevo: `internal/bot`

Hoy toda la lógica de despacho y los `handle*` viven en `cmd/bot/`
(package `main`), invisible para un segundo entry point. Se mueven a
`internal/bot` (package `bot`), que ambos `cmd/bot/main.go` y
`api/webhook.go` importan.

### `internal/bot/state.go`

```go
type PendingConversation struct {
    TxType  domain.TransactionType
    Step    ConversationStep
    Concept string
    Amount  float64
    Method  domain.PaymentMethod
}

type ChatState struct {
    PendingPhotoFileID string                `json:"pending_photo_file_id,omitempty"`
    Conversation       *PendingConversation  `json:"conversation,omitempty"`
    PendingCategoryAdd bool                  `json:"pending_category_add,omitempty"`
}

// Store persiste el ChatState de cada chat. Get de un chat sin estado
// devuelve el ChatState cero, no error — "sin estado pendiente" es un
// resultado válido, no una falla.
type Store interface {
    Get(chatID int64) (ChatState, error)
    Set(chatID int64, state ChatState) error
    Delete(chatID int64) error
}
```

Nota: los tres mapas actuales eran mutuamente excluyentes por
construcción del dispatcher (activar una espera cancelaba las otras
dos) — unificarlos en un solo `ChatState` por chat es fiel al
comportamiento ya validado, no un cambio de reglas de negocio.

### `internal/bot/memorystore/memorystore.go`

Implementación en memoria (`map[int64]bot.ChatState` + no hace falta
mutex si se sigue llamando solo desde un loop secuencial, como hoy).
Es la que sigue usando `cmd/bot/main.go`.

### `internal/bot/dispatch.go`

```go
type Deps struct {
    Bot        *telegram.Bot
    Registrar  *usecase.Registrar
    Summarizer usecase.MonthlySummarizer
    Finder     usecase.LatestTransactionFinder
    Categories usecase.CategoryResolver
    Store      Store
    AllowedChatIDs []int64
}

func Dispatch(ctx context.Context, deps Deps, update telegram.Update) error
```

Contiene exactamente el `switch` que hoy está inline en
`runPollingLoop` (autorización, foto, comandos, conversación
pendiente, `resumen`, texto normal) y todos los `handle*`/`*Reply` que
hoy son funciones privadas de `cmd/bot`. Se mueven tal cual, cambiando
las lecturas/escrituras de los tres mapas por `deps.Store.Get/Set/Delete`.

## `cmd/bot/main.go`

Se reduce a: armar `Deps` (con `memorystore.New()` como `Store`),
correr `runPollingLoop` (que ahora solo hace `GetUpdates` + llama a
`bot.Dispatch` por cada update — el resto de la lógica ya no vive
acá).

## `api/webhook.go`

```go
package api

func Handler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost { w.WriteHeader(http.StatusMethodNotAllowed); return }
    if r.Header.Get("X-Telegram-Bot-Api-Secret-Token") != os.Getenv("TELEGRAM_WEBHOOK_SECRET") {
        w.WriteHeader(http.StatusForbidden)
        return
    }
    var update telegram.Update
    if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
        w.WriteHeader(http.StatusBadRequest)
        return
    }
    deps := buildDeps() // wiring equivalente al de cmd/bot/main.go, con upstash.New() como Store
    if err := bot.Dispatch(r.Context(), deps, update); err != nil {
        log.Printf("error despachando update: %v", err)
    }
    w.WriteHeader(http.StatusOK) // Telegram reintenta si no hay 200 — devolver 200 aunque el
                                  // despacho interno haya fallado y ya se le haya avisado al
                                  // usuario con un mensaje de error (mismo criterio que hoy:
                                  // errores de negocio no son errores de transporte).
}
```

`buildDeps()` construye los mismos clientes que hoy arma `main()`
(notion.Client, notion.CategoryResolver, groq.Client, telegram.Bot) —
leyendo las mismas env vars. Se re-arman en cada invocación (una
función serverless no reusa estado entre requests de forma confiable);
son clientes HTTP livianos, no hay costo real en reconstruirlos.

## `internal/infrastructure/upstash/store.go`

Cliente HTTP propio (mismo patrón que
`internal/infrastructure/notion` — sin SDK externo), implementa
`bot.Store`:

- `Get`: `GET {UPSTASH_REDIS_REST_URL}/get/chat:{chatID}`, header
  `Authorization: Bearer {UPSTASH_REDIS_REST_TOKEN}`. Respuesta
  `{"result": "<json o null>"}` — `null` = sin estado, devolver
  `ChatState{}` sin error.
- `Set`: `POST {UPSTASH_REDIS_REST_URL}/set/chat:{chatID}/{value}/EX/3600`
  (o el comando `SET` vía el endpoint de pipeline, lo que resulte más
  simple de mapear 1:1 sin SDK) — TTL de 1 hora, configurable con una
  constante.
- `Delete`: `GET {UPSTASH_REDIS_REST_URL}/del/chat:{chatID}`.

apiBaseURL leído de env var directamente (no hay caso de test contra
Notion real acá, así que no hace falta el patrón `var` + swap salvo
para los tests con `httptest.Server`, que sí lo van a necesitar igual
que en `notion.Client`).

## `internal/infrastructure/telegram/client.go`

- `SetWebhook(ctx context.Context, url, secretToken string) error` —
  `POST .../setWebhook` con `url` y `secret_token`.
- `DeleteWebhook(ctx context.Context) error` — `POST .../deleteWebhook`.

## `cmd/setwebhook/main.go`

Comando de un solo uso, no un servicio: lee `TELEGRAM_BOT_TOKEN` y
`TELEGRAM_WEBHOOK_SECRET` de env, la URL de destino y si hay que borrar
en vez de setear de un flag (`--url` / `--delete`). Se corre a mano
después de cada despliegue a Vercel, o para volver a desarrollo local.

## Refactors de compatibilidad

- `cmd/bot/main_test.go` / `reply_test.go`: las funciones que testean
  (`matchesCommand`, `*Reply`, `parseStartCommand`, etc.) se mueven a
  `internal/bot` junto con el código que testean — los archivos de
  test migran con ellas. `cmd/bot` se queda solo con lo específico de
  long-polling (`parseAllowedChatIDs`, `isAllowed` si siguen ahí, el
  loop en sí).
- Todo lo que hoy son fakes/httptest en `internal/usecase` y
  `internal/infrastructure/notion` no cambia — esta feature no toca
  contratos de negocio, solo transporte y estado.

## Tests

- `internal/bot/dispatch_test.go`: fakes de `Deps` (reusar los fakes
  de `usecase` donde aplique, más un `Store` fake en memoria) —
  reproducir como tests de tabla los mismos casos que hoy están
  repartidos en `cmd/bot/main_test.go`/`reply_test.go`, más casos
  nuevos de prioridad de `ChatState` (una foto llega con
  `Conversation` activa → se borra; etc., ya cubierto hoy pero ahora
  contra el `Store` unificado).
- `internal/bot/memorystore/memorystore_test.go`: Get/Set/Delete
  triviales.
- `internal/infrastructure/upstash/store_test.go`: `httptest.Server`
  simulando las respuestas de la REST API de Upstash — Get con y sin
  valor previo, Set con TTL, Delete, error de red.
- `api/webhook_test.go`: request sin secret → 403; request con secret
  correcto y body válido → 200 y `Dispatch` fue invocado (verificable
  con un `Store`/`Registrar` fake); body inválido → 400.
- `internal/infrastructure/telegram/client_test.go` (si no existe
  todavía, se crea): `SetWebhook`/`DeleteWebhook` contra
  `httptest.Server`.

## Validación

- Automatizada (esta sesión): toda la batería de tests de arriba, más
  `go build ./...` (incluye compilar `api/webhook.go`, que Vercel
  compilará igual con `go build` estándar) y una prueba manual real
  con `cmd/bot/main.go` (long-polling, store en memoria) igual que las
  features 001-011, para confirmar que el refactor a `internal/bot` no
  rompió nada del comportamiento ya validado.
- **No automatizada en esta sesión** (requiere cuenta de Vercel/Upstash
  del usuario): desplegar a Vercel, setear env vars ahí, correr
  `cmd/setwebhook` apuntando a la URL real, mandar un mensaje real por
  Telegram y confirmar que llega vía webhook. Se deja documentado como
  pasos a seguir, no como criterio de aceptación de esta sesión.

## Env vars nuevas (`.env.example`)

```
TELEGRAM_WEBHOOK_SECRET=
UPSTASH_REDIS_REST_URL=
UPSTASH_REDIS_REST_TOKEN=
```

## Orden de implementación sugerido

1. `internal/bot` (state + dispatch), moviendo código de `cmd/bot` sin
   cambiar comportamiento — validar con los tests migrados +
   prueba manual real de `cmd/bot` antes de seguir.
2. `internal/bot/memorystore`.
3. `internal/infrastructure/upstash`.
4. `internal/infrastructure/telegram`: `SetWebhook`/`DeleteWebhook`.
5. `api/webhook.go`.
6. `cmd/setwebhook`.
7. `.env.example` + documentación de despliegue.
