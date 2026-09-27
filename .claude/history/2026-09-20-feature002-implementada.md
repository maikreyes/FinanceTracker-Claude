# 2026-09-20 — Feature 002 implementada: transporte Telegram

## Contexto

Con la 001 implementada y validada, el usuario pidió avanzar con 002.

## Qué se implementó

- `internal/infrastructure/telegram/client.go` — `Bot` con
  `net/http` puro: `GetUpdates` (long-poll nativo, `timeout=30` del
  lado del servidor; `httpClient.Timeout` de 40s para no cortar la
  conexión antes de que Telegram responda), `SendMessage`, tipos
  mínimos `Update`/`Message`/`Chat`, helper `call` genérico que
  decodifica el sobre `{ok, result, description}` de la API.
- `cmd/bot/main.go` — wiring completo: `notion.NewClient` +
  `notion.NewCategoryResolver` + `usecase.NewRegistrar` +
  `telegram.NewBot`; `runPollingLoop` con offset en memoria y backoff
  exponencial acotado (2s → 30s) ante error de `GetUpdates`;
  `signal.NotifyContext` para apagado limpio; `handleMessage` responde
  con el ID de página en éxito o el error crudo en falla (placeholder
  funcional — 004 define el copy real).

## Validación real (con el bot corriendo de verdad)

1. Bot arrancado (`go run ./cmd/bot`) contra el Telegram real del
   usuario.
2. Mensaje corto (formato inválido) → logueado, respondido por
   Telegram, **sin crash**.
3. Mensaje completo (`... 6000 efectivo comida`, tipo Ingreso) →
   confirmado en Notion vía query real: `Name: pizza, Amount: 6000,
   Type: Ingreso, Payment Method: Cash, Category: Comida` — registrado
   correctamente por el usuario real, no por mí. Esa página es un
   movimiento real, no se tocó/borró.
4. Apagado limpio probado con `SIGINT` real: log `bot detenido` tras
   cancelar el contexto en medio de un long-poll.

## Hallazgo operativo

Un `pkill -f "go run ./cmd/bot"` anterior no mató el proceso viejo (el
patrón no machea al binario hijo real de `go run`) — terminé con dos
instancias corriendo, y Telegram devolvió `409 Conflict: terminated by
other getUpdates request`. Confirma que el backoff/reintento funciona
ante error real de la API, pero también es una limitación real a
documentar: **solo una instancia por token**. Agregado a
`.claude/rules/telegram-rules.md`, con la instrucción de verificar con
`ps aux` en vez de confiar en un `pkill` anterior.

## Estado al cierre de sesión

- `go build ./...`, `go vet ./...`, `go test ./...` limpios.
- Feature 002 marcada **implementado ✅**, movida a "Hecho" en
  `specs/constitution/roadmap.md`. Sin instancias del bot corriendo
  (verificado con `ps aux`).
- **Riesgo abierto documentado:** sin 003, el bot acepta cualquier
  chat — no correrlo sin supervisión ni compartir el username hasta
  que 003 esté lista. Siguiente paso natural: implementar 003.
