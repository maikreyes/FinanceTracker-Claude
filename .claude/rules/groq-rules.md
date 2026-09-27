# Groq API rules

- API key leída solo de `GROQ_API_KEY`. Nunca hardcodear.
- `pkg/groq/services/` es dueño de todas las llamadas a la
  API de Groq (compatible con el formato de OpenAI:
  `/chat/completions`). Fuera de ese paquete nadie importa el
  cliente de Groq — solo la interfaz `ports.ReceiptInterpreter`.
- **Modelo fijo en una constante** (`defaultModel` en
  `pkg/groq/services/services.go`, valor real al
  2026-09-20: `qwen/qwen3.8-27b`) — los modelos con soporte de visión
  de Groq cambian o se retiran seguido (el modelo original de esta
  feature, `llama-3.2-11b-vision-preview`, ya estaba dado de baja
  antes de la primera prueba real). Si `Interpret` falla con
  `model_decommissioned` o "modelo no encontrado": `GET
  https://api.groq.com/openai/v1/models` con el `GROQ_API_KEY` real,
  filtrar por `input_modalities` conteniendo `"image"`, y actualizar
  `defaultModel` con lo que aparezca.
- `response_format: {"type": "json_object"}` para forzar salida JSON
  estructurada — nunca parsear texto libre de la respuesta del modelo.
  El prompt debe mencionar "JSON" explícitamente (requisito del modo
  `json_object`, heredado de la convención de OpenAI que Groq sigue).
- El prompt le pasa las categorías reales del usuario (consultadas en
  vivo contra Notion) y le pide elegir solo entre esas, o `null` — no
  debe inventar categorías nuevas.
- Solo el monto es obligatorio para registrar un movimiento
  interpretado por IA — categoría y medio de pago que el modelo no
  pudo determinar (o que sugirió pero no matchea contra Notion) se
  dejan vacíos, no bloquean el registro.
- La respuesta al usuario debe decir explícitamente "interpretado por
  IA" — la confianza en estos datos no es la misma que un registro por
  texto exacto.
- **Contexto de por qué Groq y no Gemini:** se eligió Gemini primero,
  pero su flujo de creación de API key pedía vincular facturación en
  la cuenta del usuario (aunque el tier gratuito siguiera siendo $0)
  — se cambió a Groq porque no lo pedía. Ver
  `.claude/history/2026-09-20-feature007-cambio_proveedor_groq.md`.
- Límite de tamaño del File Upload API de Notion en modo `single_part`
  (usado para subir la foto, no específico de Groq): no confirmado
  explícitamente todavía — revisar si se suben fotos pesadas.
