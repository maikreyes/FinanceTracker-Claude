# 003 · Restringir acceso al bot — Tareas

_Checklist accionable derivada del `plan.md`. Tareas pequeñas y concretas; marca `[x]` al completarlas._

- [ ] `.env.example`: agregar `TELEGRAM_ALLOWED_CHAT_IDS`.
- [ ] `.claude/rules/seguridad.md` y `.claude/rules/telegram-rules.md`:
      documentar la nueva variable y la regla fail-closed.
- [ ] `cmd/bot/main.go`: parseo de `TELEGRAM_ALLOWED_CHAT_IDS` al
      arrancar + log de advertencia si queda vacía.
- [ ] `isAllowed(chatID int64, allowed []int64) bool` + chequeo en el
      loop de polling antes de despachar al usecase.
- [ ] Prueba manual: mensaje desde un chat no autorizado no crea nada
      en Notion y no genera respuesta.
- [ ] Averiguar y guardar (en `.env` local, no versionado) el
      `chat_id` real del usuario — instrucción: hablarle al bot
      `@userinfobot` o revisar el `chat.id` del primer update recibido
      en los logs.
- [ ] `go build ./...`, `go vet ./...` en verde.
- [ ] Validar contra los criterios de aceptación de `spec.md`.
- [ ] Mover la feature a "Hecho" en `../../constitution/roadmap.md`.
