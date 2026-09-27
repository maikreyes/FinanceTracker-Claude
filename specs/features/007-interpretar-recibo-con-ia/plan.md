# 007 · Interpretar recibo con IA (sin caption) — Plan

_Cómo se implementa lo descrito en `spec.md`. Debe respetar la `constitution/`._

## Enfoque

Adaptador de infraestructura para Groq (mismo patrón que
`internal/infrastructure/notion/` y `telegram/`: `net/http` puro, sin
SDK), detrás de una interfaz `ReceiptInterpreter` en usecase, agnóstica
de proveedor. Un método `Registrar.RegisterFromPhotoAI` orquesta:
interpretar con Groq → resolver categoría (si la IA sugirió una,
best-effort) → subir la foto (`ReceiptUploader` de 006) →
`CreateExpense`. El caption manual (006) sigue siendo el camino por
defecto; esta feature solo entra cuando no hay caption.

**Nota:** se evaluó Google Gemini primero (ver `spec.md` y
`.claude/rules/groq-rules.md`), pero su flujo de creación de API key
pedía vincular facturación — se cambió a Groq porque no lo pedía.

## Implementación

1. `internal/infrastructure/groq/client.go`:
   - `Client{apiKey, model string, httpClient *http.Client}`,
     `NewClient(apiKey string) *Client` — modelo fijo en
     `defaultModel` (`llama-3.2-11b-vision-preview` al momento de
     implementar; los modelos "preview" de Groq cambian seguido,
     confirmar el vigente si falla).
   - `Interpret(imageData []byte, categoryNames []string)
     (usecase.InterpretedReceipt, error)` — `POST
     https://api.groq.com/openai/v1/chat/completions` (API compatible
     con el formato de OpenAI). Body: `messages` con un mensaje de rol
     `user`, `content` como array con un part de texto (el prompt) y
     un part `image_url` con la imagen en base64 (`data:image/jpeg;
     base64,...`).
     - El prompt pide extraer `monto`, `concepto`, `categoria`,
       `medio_de_pago`.
     - Da la lista real de categorías (`categoryNames`) y las 4
       opciones válidas de medio de pago, y pide elegir solo entre
       esas o dejar el campo `null`.
     - `response_format: {"type": "json_object"}` fuerza salida JSON
       (modo estructurado de la API estilo OpenAI que Groq soporta) —
       evita parsear texto libre. El prompt menciona "JSON"
       explícitamente (requisito de ese modo).
   - La función devuelve directamente `usecase.InterpretedReceipt`
     (el paquete `groq` importa `usecase` — dirección de dependencia
     permitida: `infrastructure -> usecase`). `Amount` es puntero para
     distinguir "no encontrado" (`nil`) de `0`.
2. `internal/usecase/ports.go`: interfaz `ReceiptInterpreter` +
   struct `InterpretedReceipt{Amount *float64, Concept, CategoryName,
   PaymentMethod string}` — agnóstica de proveedor, no menciona Groq.
3. `internal/usecase/errors.go`: `ErrReceiptAmountNotFound` (la IA no
   encontró un monto — camino distinto a una falla de API) y
   `ErrReceiptInterpretationFailed{Err error}` (falla de red/cuota/API,
   no expone el error real al chat, mismo patrón que `ErrWriteFailed`).
4. `internal/usecase/register_transaction.go`:
   - `RegisterResult` gana un campo `Interpreted bool`.
   - `Registrar` gana un campo `interpreter ReceiptInterpreter`
     (`NewRegistrar` pasa a recibir 4 argumentos).
   - `Registrar.RegisterFromPhotoAI(photo []byte, filename string)
     (RegisterResult, error)`: llama `interpreter.Interpret`, valida
     `Amount != nil` (si no, `ErrReceiptAmountNotFound`), intenta
     resolver la categoría sugerida (si la IA la dio) — **a diferencia
     de `Register`, si no matchea no bloquea el registro**, solo se
     deja vacía (la IA no es el usuario escribiendo a mano). Valida el
     `medio_de_pago` contra las 4 opciones exactas
     (`matchPaymentMethodName`); si no matchea, también queda vacío.
     Sube la foto (`ReceiptUploader`, de 006), arma la
     `domain.Transaction` con `Type: domain.TransactionTypeExpense`
     fijo, llama `CreateExpense`.
5. `cmd/bot/main.go`:
   - `handlePhotoMessage` descarga la foto siempre; si `Caption` tiene
     texto → `RegisterWithPhoto` (006); si está vacío → nuevo camino,
     `RegisterFromPhotoAI`.
   - `successReply`: si `result.Interpreted`, antepone una línea
     "Interpretado por IA — revisá los datos en Notion si algo no
     coincide.".
   - `errorReply`: casos nuevos para `ErrReceiptAmountNotFound`
     ("No pude leer un monto en la foto. Reenviala con el texto del
     movimiento como caption.") y `ErrReceiptInterpretationFailed`
     ("No se pudo interpretar la foto en este momento. Intenta de
     nuevo o mandala con caption.").
6. `.env.example`/`.env`: nueva variable `GROQ_API_KEY`.
7. `.claude/rules/groq-rules.md`: API key nunca hardcodeada, contrato
   de la interfaz, contexto de por qué Groq y no Gemini.

## Decisiones

- **Caption manual tiene prioridad sobre la IA** — determinístico,
  gratis, sin dependencia externa; la IA es el fallback para cuando el
  usuario no quiere escribir el caption, no un reemplazo.
- **Groq en vez de Gemini** — Gemini pedía vincular facturación para
  sacar la API key (aunque el uso dentro del tier gratuito siguiera
  siendo $0); Groq no lo pidió. `ReceiptInterpreter` quedó agnóstica
  de proveedor desde el diseño, así que este cambio no tocó `usecase`
  ni `cmd/bot` más allá del wiring (nombre del paquete + env var).
- **`response_format: json_object` en vez de parsear texto libre** —
  más confiable que pedirle a un LLM que "responda en JSON" sin
  forzarlo, evita parseos frágiles.
- **Tipo fijo `Egreso`** — un recibo fotografiado es casi siempre un
  gasto; soportar ingresos por esta vía es un salto de alcance
  innecesario para v1.
- **Monto ausente → error, no registro parcial** — un movimiento sin
  monto no tiene sentido; mejor pedir que se reenvíe con caption que
  guardar un registro incompleto.
- **Categoría/medio de pago sugeridos por IA que no matchean → se
  dejan vacíos, no bloquean** — distinto del camino de texto manual
  (`Register`), donde una categoría inválida sí es un error: ahí la
  escribió el usuario; acá es la IA la que pudo haberse equivocado.
- **Disclaimer explícito de "interpretado por IA"** — la confianza en
  los datos de un registro por caption manual y uno interpretado por
  IA no es la misma; el usuario debe poder distinguirlos sin abrir
  Notion.

## Riesgos

- **Límites de tasa/cuota del tier gratuito de Groq** — confirmar el
  número vigente al operar el bot (cambia con el tiempo). Para un bot
  personal de uso esporádico no debería ser problema, pero el error de
  cuota agotada debe distinguirse de una falla genérica en los logs
  (aunque el chat vea el mismo mensaje genérico).
- **Alucinación del modelo** — puede "leer" un monto o categoría
  incorrectos con confianza. Mitigado parcialmente por el disclaimer
  en la respuesta, no por validación automática (no hay forma de
  validar un monto contra la realidad sin el recibo físico).
- **Tamaño de imagen** — fotos grandes en base64 inline pueden acercarse
  al límite de tamaño de request de la API de Groq. Si pasa, error
  claro en vez de timeout silencioso — confirmar el límite exacto al
  operar con fotos pesadas.
- **Nombre del modelo cambia con el tiempo** — los modelos "preview"
  con soporte de visión de Groq se retiran/renombran más seguido que
  en otros proveedores; el nombre queda en una única constante fácil
  de actualizar (`defaultModel`), no hardcodeado en múltiples lugares.
