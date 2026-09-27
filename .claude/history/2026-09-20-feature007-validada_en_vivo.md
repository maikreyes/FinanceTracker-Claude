# 2026-09-20 — Feature 007 validada en vivo: modelo dado de baja, corregido

## Contexto

Con el código de 007 completo (entrada anterior, cambio de proveedor a
Groq), el usuario consiguió una `GROQ_API_KEY` real
(`console.groq.com/keys`, confirmando que no pedía tarjeta).

## Validación

1. **Smoke test directo del cliente** (`cmd/groq_smoketest/`,
   temporal, borrado después) contra una imagen de recibo sintética
   generada con Pillow (instalado on-the-fly con `pip3 install
   pillow` — no había en el sistema): falló con
   `model_decommissioned` — el modelo del `plan.md`
   (`llama-3.2-11b-vision-preview`) ya no existía en Groq.
2. `GET https://api.groq.com/openai/v1/models` con el `GROQ_API_KEY`
   real, filtrado por `input_modalities` conteniendo `"image"`: único
   modelo vigente con visión en la cuenta, `qwen/qwen3.8-27b`.
   `defaultModel` actualizado en
   `internal/infrastructure/groq/client.go`.
3. Reintentado el smoke test con la imagen sintética: interpretó
   monto exacto (19000), nombre del comercio, y payment method
   (`Cash`, de "PAGO: EFECTIVO" en la imagen) correctamente. La
   categoría varió entre corridas (`Comida`/`Ahorro`) — esperable, el
   nombre de la tienda de prueba ("EL AHORRO") confundía al modelo a
   propósito.
4. **Prueba real con el bot corriendo** (una sola instancia
   verificada): el usuario mandó una foto real de un recibo (recarga
   Nequi/Bancolombia) **sin caption**. Confirmado en Notion vía query
   directa: `Name: RECARGA NEQU BANCOLCOMBIA`, `Amount: 120000`,
   `Type: Egreso`, `Category` resuelta a `Servicios`, `Receipt` con el
   archivo adjunto. `Payment Method` quedó vacío — la IA no estaba
   seguro de cómo se pagó y no inventó, comportamiento correcto según
   `spec.md`.
5. Bot detenido limpio (`bot detenido` en el log), sin instancias
   corriendo al cerrar.

## Estado al cierre de sesión

- `go build ./...`, `go vet ./...`, `go test ./...` limpios.
- Feature 007 marcada **implementado ✅**, movida a "Hecho" en
  `specs/constitution/roadmap.md`.
- `.claude/rules/groq-rules.md` actualizada con el modelo vigente y el
  procedimiento para encontrar uno nuevo si vuelve a pasar
  (`GET /v1/models`, filtrar por `input_modalities`).
- **Las 7 features especificadas (001-007) están implementadas y
  validadas en vivo.** No queda ninguna spec pendiente de implementar
  en `specs/features/`.
