# 007 · Interpretar recibo con IA (sin caption) — Tareas

_Checklist accionable derivada del `plan.md`. Tareas pequeñas y concretas; marca `[x]` al completarlas._

- [x] Evaluar Google Gemini → descartado (pedía vincular facturación
      para sacar la API key) → Groq elegido en su lugar.
- [x] `internal/infrastructure/groq/client.go`: `Client.Interpret`
      contra `/chat/completions`, `response_format: json_object`,
      prompt con lista real de categorías + medios de pago válidos.
- [x] `internal/usecase/ports.go`: interfaz `ReceiptInterpreter` +
      struct `InterpretedReceipt` (agnóstica de proveedor).
- [x] `internal/usecase/errors.go`: `ErrReceiptAmountNotFound`,
      `ErrReceiptInterpretationFailed`.
- [x] `RegisterResult.Interpreted bool`.
- [x] `Registrar` gana el campo `interpreter` (`NewRegistrar` ahora
      recibe 4 argumentos) + `RegisterFromPhotoAI`.
- [x] `internal/usecase/parse_transaction.go`: `matchPaymentMethodName`
      — valida el `medio_de_pago` sugerido por la IA contra las 4
      opciones reales de Notion.
- [x] `cmd/bot/main.go`: `handlePhotoMessage` bifurca caption vacío →
      `RegisterFromPhotoAI` (antes pedía reenviar).
- [x] `successReply`: disclaimer cuando `Interpreted == true`.
- [x] `errorReply`: casos nuevos de `ErrReceiptAmountNotFound` /
      `ErrReceiptInterpretationFailed`.
- [x] `.env.example`: `GROQ_API_KEY`.
- [x] `.claude/rules/groq-rules.md` nuevo (reemplaza el
      `gemini-rules.md` que se había escrito antes del cambio de
      proveedor).
- [x] Tests unitarios de `RegisterFromPhotoAI` con fakes
      (`fakeReceiptInterpreter`) — completo con categoría y medio de
      pago, sin monto, falla de interpretación, categoría/medio de
      pago no determinados (no bloquean), categoría sugerida que no
      matchea (no bloquea).
- [x] Tests unitarios de `groq.Client.Interpret` con
      `httptest.Server` — respuesta completa (+ prompt incluye las
      categorías reales), sin monto, error de API.
- [x] Prueba manual real: foto de un recibo real (recarga Nequi/
      Bancolombia) sin caption → interpretado y registrado en Notion
      (`Name: RECARGA NEQU BANCOLCOMBIA`, `Amount: 120000`, `Type:
      Egreso`, `Category: Servicios` resuelta, `Receipt` adjunto).
      `Payment Method` quedó vacío — la IA no estaba segura, no
      inventó (comportamiento correcto según spec).
- [x] `go build ./...`, `go vet ./...`, `go test ./...` en verde.
- [x] Validar contra los criterios de aceptación de `spec.md`.
- [x] Mover la feature a "Hecho" en `../../constitution/roadmap.md`.

## Hallazgo durante la prueba manual

El modelo original del `plan.md`
(`llama-3.2-11b-vision-preview`) estaba **dado de baja** en Groq
(`model_decommissioned`) — confirmado con `GET /v1/models` filtrando
por `input_modalities` con `"image"`: al momento de esta prueba, el
único modelo con visión disponible en la cuenta era
`qwen/qwen3.8-27b`. `defaultModel` actualizado en
`internal/infrastructure/groq/client.go`. Confirma el riesgo ya
anotado en `plan.md` ("el nombre del modelo cambia con el tiempo") —
la próxima vez que falle con ese mismo error, repetir el mismo chequeo
contra `/v1/models`.
