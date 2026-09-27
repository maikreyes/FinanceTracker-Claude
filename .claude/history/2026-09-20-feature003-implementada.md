# 2026-09-20 — Feature 003 implementada: restringir acceso al bot

## Contexto

Con 001 y 002 implementadas (y la spec 006 ya escrita), el usuario
pidió avanzar con 003.

## Qué se implementó

- `cmd/bot/main.go`: `parseAllowedChatIDs(raw string) []int64` (split
  por coma, ignora vacíos, loguea advertencia por cada valor no
  numérico sin abortar el arranque) + `isAllowed(chatID int64, allowed
  []int64) bool`. Chequeo en el loop de polling, antes de
  `handleMessage` — un chat no autorizado se ignora en silencio (sin
  `SendMessage`), solo queda logueado server-side.
- Fail closed: si `TELEGRAM_ALLOWED_CHAT_IDS` queda vacía, log de
  advertencia al arrancar y `isAllowed` nunca es `true`.
- `.env.example`/`.env` real: nueva variable
  `TELEGRAM_ALLOWED_CHAT_IDS`. Valor real usado:
  `1196865164` — el `chat_id` del usuario, capturado de los logs de la
  prueba manual de la feature 002 (no hizo falta pedirlo aparte).
- `cmd/bot/main_test.go`: tests de `isAllowed`/`parseAllowedChatIDs`.
- `.claude/rules/telegram-rules.md` y `.claude/rules/seguridad.md`
  actualizadas con el contrato de la whitelist.

## Validación real (con contratiempos propios, no del código)

Al probar en vivo se dejaron sin querer dos/tres instancias del bot
corriendo en paralelo varias veces seguidas (un `pkill`/`kill` que no
alcanzó a todos los procesos, y un `curl` de debug directo a
`getUpdates` que compitió con el bot) — Telegram devolvió `409
Conflict` repetidas veces y probablemente se "comió" un mensaje de
prueba del usuario sin procesarlo (justo el riesgo de offset-en-memoria
ya documentado en 002, no un bug de 003). Se resolvió verificando
`ps aux` explícitamente antes de cada prueba en vez de asumir que un
`kill` anterior funcionó.

Con una sola instancia confirmada:

1. **Camino autorizado:** mensaje real del chat `1196865164` →
   registrado en Notion (`Name: pizza`), sin log de error (el éxito no
   loguea nada, ver limitación ya anotada para la feature 004).
2. **Camino no autorizado:** bot reiniciado con
   `TELEGRAM_ALLOWED_CHAT_IDS=999999999` (no incluye el chat real) →
   mismo chat mandó otro mensaje → log `mensaje ignorado de chat no
   autorizado: 1196865164`, sin respuesta en Telegram. Confirmado el
   comportamiento fail-closed / silencioso pedido en la spec.
3. Bot reiniciado una última vez con la whitelist real
   (`1196865164`) y detenido limpio al terminar — sin instancias
   corriendo al cierre de la sesión.

## Estado al cierre de sesión

- `go build ./...`, `go vet ./...`, `go test ./...` limpios.
- Feature 003 marcada **implementado ✅**, movida a "Hecho" en
  `specs/constitution/roadmap.md`. El riesgo abierto anotado en la
  entrada de 002 ("el bot acepta cualquier chat") queda resuelto.
- Sin instancias del bot corriendo (verificado con `ps aux`).
- Siguiente paso natural: 004 (confirmación y manejo de errores) o 005
  (resumen) — ambas ya solo dependen de 001-003. 006 (foto de recibo)
  también queda disponible, es independiente.
