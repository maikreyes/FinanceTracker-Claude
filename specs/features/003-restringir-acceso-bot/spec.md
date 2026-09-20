# 003 · Restringir acceso al bot

**Estado:** propuesta

## Qué hace

El bot solo procesa mensajes de uno o más `chat_id` autorizados
(whitelist). Un mensaje de cualquier otro chat se ignora — no se
parsea, no se toca Notion.

## Por qué

El bot escribe directo en el Notion financiero real del usuario. Sin
esto, cualquiera que encuentre el username del bot en Telegram podría
registrar movimientos falsos en su base de datos personal.

## Criterios de aceptación

- [ ] Nueva variable de entorno `TELEGRAM_ALLOWED_CHAT_IDS` (lista de
      IDs separados por coma) documentada en `.env.example`.
- [ ] Mensaje de un `chat_id` en la whitelist → se procesa normal
      (feature 001 + 002).
- [ ] Mensaje de un `chat_id` fuera de la whitelist → se ignora, el bot
      **no responde nada** (no confirma a un desconocido que el bot
      existe o funciona).
- [ ] Si `TELEGRAM_ALLOWED_CHAT_IDS` no está seteada o está vacía → el
      bot no procesa ningún mensaje (fail closed, no fail open) y deja
      un log de advertencia al arrancar.
- [ ] La whitelist se lee una sola vez al arrancar el proceso (no hace
      falta que sea dinámica).

## Fuera de alcance

- Múltiples usuarios con roles distintos — esto es un bot personal, la
  whitelist es solo "quién puede hablarle", no hay permisos por rol.
- Administrar la whitelist desde el propio chat de Telegram (comando
  para agregar/quitar IDs) — se edita a mano en `.env`.
- Rate limiting o detección de abuso más allá del whitelist simple.
