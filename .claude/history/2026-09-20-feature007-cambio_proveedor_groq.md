# 2026-09-20 — Feature 007: implementada, con cambio de proveedor a Groq

## Contexto

Con 006 implementada, el usuario pidió avanzar con 007 (ya
especificada con Google Gemini como proveedor de IA). Al intentar
sacar la API key en aistudio.google.com/apikey, Google le pedía
vincular una cuenta de facturación — aunque el tier gratuito siguiera
siendo $0 dentro de cuota, eso bloqueaba conseguir la key sin tarjeta.
Se preguntó cómo seguir (vincular igual / probar el link directo de AI
Studio / cambiar de proveedor) — el usuario eligió cambiar a **Groq**.

## Qué se implementó

- `internal/infrastructure/groq/client.go` (reemplaza el paquete
  `gemini/` que nunca llegó a probarse con una key real):
  `Client.Interpret` contra `POST
  https://api.groq.com/openai/v1/chat/completions` (API compatible con
  el formato de OpenAI) — imagen como `image_url` con
  `data:image/jpeg;base64,...`, `response_format: {"type":
  "json_object"}` para forzar JSON estructurado. Modelo fijo en
  `defaultModel = "llama-3.2-11b-vision-preview"` (modelos "preview"
  de Groq cambian seguido, confirmar el vigente si falla).
- **La interfaz `usecase.ReceiptInterpreter` y el struct
  `InterpretedReceipt` ya habían quedado agnósticos de proveedor** —
  el cambio de Gemini a Groq no tocó `usecase/ports.go` en su forma,
  solo hubo que corregir un comentario que mencionaba
  `gemini.Client`. Confirma que el diseño de puertos/adaptadores pagó:
  todo el trabajo de `Registrar.RegisterFromPhotoAI`,
  `matchPaymentMethodName`, los errores tipados, y el wiring en
  `cmd/bot/main.go` se hizo una sola vez.
- `internal/usecase/register_transaction.go`: `Registrar` gana el
  campo `interpreter` (`NewRegistrar` ahora recibe 4 argumentos) +
  `RegisterFromPhotoAI(photo []byte, filename string)`. Decisión
  distinta a `Register`/`RegisterWithPhoto`: si la categoría o el
  medio de pago que sugiere la IA no matchean contra Notion, **no
  bloquea el registro** — se dejan vacíos (la IA pudo haberse
  equivocado; el usuario escribiendo a mano sí es un error real).
- `internal/usecase/parse_transaction.go`: `matchPaymentMethodName`
  valida el `medio_de_pago` de la IA contra las 4 opciones reales de
  Notion (en inglés: `Cash`, `Credit Card`, `Debit Card`, `Bank`).
- `internal/usecase/errors.go`: `ErrReceiptAmountNotFound`,
  `ErrReceiptInterpretationFailed`.
- `cmd/bot/main.go`: `handlePhotoMessage` ahora bifurca — con caption
  → `RegisterWithPhoto` (006); sin caption → `RegisterFromPhotoAI`
  (007, ya no pide reenviar). `successReply` antepone disclaimer
  "Interpretado por IA" cuando `Interpreted == true`. `errorReply` con
  los 2 casos nuevos.
- Nueva env var `GROQ_API_KEY` en `.env.example`.
- `.claude/rules/groq-rules.md` (reemplaza el `gemini-rules.md`
  escrito antes del cambio de proveedor — se borró, no se dejó
  desactualizado).
- Tests: `internal/infrastructure/groq/client_test.go`
  (`httptest.Server`: respuesta completa + prompt incluye categorías
  reales, sin monto, error de API) y
  `internal/usecase/register_transaction_test.go` ampliado
  (`RegisterFromPhotoAI`: completo, sin monto, falla de
  interpretación, categoría/medio de pago no determinados no bloquean,
  categoría sugerida que no matchea no bloquea).

## Estado al cierre de sesión

- `go build ./...`, `go vet ./...`, `go test ./...` limpios,
  `gofmt` corrido.
- **Falta la validación en vivo** — no se probó con una
  `GROQ_API_KEY` real todavía (el usuario está sacándola). Spec 007
  marcada "implementado, sin validar en vivo ⏳", no movida a "Hecho"
  en el roadmap hasta que se confirme con una foto real.
- `specs/features/007-interpretar-recibo-con-ia/{spec,plan,tasks}.md`
  actualizados para reflejar Groq (no Gemini) en todos los detalles
  técnicos, no solo un search-and-replace del nombre.
