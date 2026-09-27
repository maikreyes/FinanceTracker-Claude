# Seguridad

- Ningún dato sensible hardcodeado — nunca, ni siquiera como default
  "solo dev". `TELEGRAM_BOT_TOKEN`, `NOTION_API_KEY`,
  `NOTION_DATABASE_ID`: siempre variable de entorno, nunca un valor real
  en `.go`, `.yaml`, `docker-compose.yml` ni ningún archivo versionado.
- Nombres de variables (sin valores) documentados en `.env.example` en
  la raíz; el `.env` real nunca se commitea (ya está en `.gitignore`) ni
  se comparte fuera del equipo.
- Texto entrante de Telegram es input no confiable: validar/sanitizar
  antes de pasarlo a `pkg/transaction/services`.
- El bot escribe directo en el Notion financiero real del usuario —
  acceso restringido por `TELEGRAM_ALLOWED_CHAT_IDS` (whitelist, fail
  closed). Nunca correr `go run ./cmd` sin esa variable seteada
  fuera de una prueba puntual y supervisada.
