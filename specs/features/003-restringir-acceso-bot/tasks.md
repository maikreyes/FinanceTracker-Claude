# 003 · Restringir acceso al bot — Tareas

_Checklist accionable derivada del `plan.md`. Tareas pequeñas y concretas; marca `[x]` al completarlas._

- [x] `.env.example`: agregar `TELEGRAM_ALLOWED_CHAT_IDS`.
- [x] `.claude/rules/seguridad.md` y `.claude/rules/telegram-rules.md`:
      documentar la nueva variable y la regla fail-closed.
- [x] `cmd/bot/main.go`: parseo de `TELEGRAM_ALLOWED_CHAT_IDS` al
      arrancar + log de advertencia si queda vacía
      (`parseAllowedChatIDs`).
- [x] `isAllowed(chatID int64, allowed []int64) bool` + chequeo en el
      loop de polling antes de despachar al usecase.
- [x] Prueba manual: mensaje desde un chat no autorizado no crea nada
      en Notion y no genera respuesta — validado en vivo (log:
      `mensaje ignorado de chat no autorizado: 1196865164`).
- [x] Averiguar y guardar (en `.env` local, no versionado) el
      `chat_id` real del usuario — `1196865164`, capturado de los logs
      de la prueba manual de la feature 002.
- [x] `go build ./...`, `go vet ./...` en verde (+ tests unitarios de
      `isAllowed`/`parseAllowedChatIDs` en `cmd/bot/main_test.go`).
- [x] Validar contra los criterios de aceptación de `spec.md` — ambos
      caminos probados con el bot corriendo de verdad: chat autorizado
      registra en Notion, chat no autorizado se ignora sin respuesta.
- [x] Mover la feature a "Hecho" en `../../constitution/roadmap.md`.
