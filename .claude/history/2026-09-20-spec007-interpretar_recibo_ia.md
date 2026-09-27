# 2026-09-20 — Spec 007: interpretar recibo con IA (sin caption)

## Contexto

Con las 5 features del MVP (001-005) implementadas y validadas, y 006
(foto de recibo) especificada pero sin implementar, el usuario pidió
una feature nueva: que un modelo de IA gratuito interprete la foto del
recibo, para no tener que escribir el caption que pide 006.

Se preguntó qué proveedor de IA gratuito usar — opciones: Google
Gemini, Groq, OpenRouter (modelos `:free`), o local con Ollama. El
usuario eligió Google Gemini.

## Qué se hizo

Creada `specs/features/007-interpretar-recibo-con-ia/` con
`spec.md`/`plan.md`/`tasks.md`. Depende de 006 (reusa su
infraestructura de descarga de foto de Telegram).

- **Relación con 006:** el caption manual sigue siendo el camino por
  defecto (determinístico, gratis, sin dependencia externa) — la IA
  es el *fallback* solo cuando la foto llega sin caption. Esto
  reemplaza el criterio de 006 de "foto sin caption → pedir que se
  reenvíe" una vez que 007 esté implementada.
- Campos que extrae la IA: `Monto` (obligatorio), `Concepto`,
  `Categoría` (elige entre las categorías reales del usuario, pasadas
  en el prompt — no inventa categorías nuevas), `Medio de pago` (solo
  si es claro en la imagen). `Tipo` queda fijo en `Egreso` — un recibo
  fotografiado es casi siempre un gasto, soportar ingresos por esta
  vía queda fuera de alcance.
- Decisión técnica: `generationConfig.response_mime_type:
  "application/json"` de la API de Gemini para forzar salida JSON
  estructurada, en vez de parsear texto libre de un LLM.
- Sin monto identificado → error, cae de vuelta al camino de caption
  manual — no se registra nada a medias.
- La respuesta de confirmación debe decir explícitamente que fue
  interpretado por IA (`RegisterResult.Interpreted`), para que el
  usuario sepa que debe revisar los datos — distinto nivel de
  confianza que un registro por texto exacto.
- Nueva env var `GEMINI_API_KEY`, nuevo paquete
  `internal/infrastructure/gemini/`, nueva regla
  `.claude/rules/gemini-rules.md` (a crear al implementar).
- Riesgos anotados: cuota del tier gratuito de Gemini, alucinación del
  modelo (mitigada solo por el disclaimer, no por validación
  automática — no hay forma de validar un monto sin el recibo físico),
  tamaño de imagen en base64 inline, nombre del modelo Gemini puede
  cambiar con el tiempo (confirmar al implementar).
- `specs/constitution/roadmap.md`: 006 y 007 quedan como mejoras
  encima del MVP (001-005 ya completo), 007 explícitamente después de
  006 por la dependencia.

## Estado al cierre de sesión

Spec, plan y tasks de 007 completos. **Sin implementar** — depende de
que 006 se implemente primero (comparten infraestructura de descarga
de foto), y de OK explícito del usuario para empezar a codear
cualquiera de las dos.
