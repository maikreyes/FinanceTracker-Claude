# 2026-09-20 — Specs 002-005: roadmap completo del MVP

## Contexto

Con la 001 lista, el usuario pidió especificar el resto de las
features antes de tocar código, "de la misma manera" (spec/plan/tasks,
sin implementar). Se propusieron 4 candidatas basadas en lo que había
quedado explícitamente fuera de la 001 ("fuera de esta feature" en su
`tasks.md`) más una idea nueva de consulta; el usuario las aprobó
todas.

## Features creadas

- **002 · Transporte Telegram** — long-polling real con `net/http`
  puro (mismo patrón que el cliente de Notion, sin SDK externo, sin
  dependencias nuevas sin avisar). Conecta el usecase de 001 con
  Telegram de verdad. Decisión: offset de `getUpdates` solo en
  memoria para v1 (riesgo aceptado: reproceso tras reinicio).
- **003 · Restringir acceso al bot** — whitelist de `chat_id` vía
  `TELEGRAM_ALLOWED_CHAT_IDS`. Decisiones: fail closed si no está
  configurada, no responder nada a chats no autorizados (no confirmar
  que el bot existe), chequeo en el transporte, no en el usecase.
- **004 · Confirmación y manejo de errores en el chat** — errores
  tipados en el usecase (`ErrUnrecognizedType`,
  `ErrUnrecognizedPaymentMethod`, `ErrUnrecognizedCategory`,
  `ErrWriteFailed`), copy vive en la capa de Telegram. Categorías del
  mensaje de error se listan en vivo (agrega `CategoryResolver.
  ListNames()` a la interfaz de 001). Sin exponer errores internos de
  Notion al chat.
- **005 · Consultar resumen de gastos** — comando `resumen`, suma
  `Amount` por `Type` client-side sobre una query filtrada por fecha
  (no usa la fórmula `Add to Month`, que es por fila, no agregable
  server-side). Nueva interfaz `MonthlySummarizer`. Riesgo anotado:
  zona horaria (`America/Bogota` explícito, no `time.Now()` a secas).

## Roadmap

`specs/constitution/roadmap.md` actualizado: orden de implementación
001 → 002 → 003 → 004 → 005 (cada una depende de que la anterior
exista para ser usable de punta a punta).

## Estado al cierre de sesión

5 features completamente especificadas (spec/plan/tasks), **ninguna
implementada todavía** — según el flujo pactado, falta OK explícito
del usuario para empezar a codear. `go build ./...` / `go vet ./...`
siguen limpios (no se tocó código en esta entrada, solo `specs/`).
