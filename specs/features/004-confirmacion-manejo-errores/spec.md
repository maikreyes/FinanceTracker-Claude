# 004 · Confirmación y manejo de errores en el chat

**Estado:** propuesta

## Qué hace

Define exactamente qué le responde el bot al usuario en Telegram: un
mensaje de confirmación cuando el movimiento se registra bien, y un
mensaje de error específico y accionable para cada tipo de falla que
puede producir la feature 001 (parseo) o la escritura en Notion.

## Por qué

La feature 002 necesita "responder algo" — esta feature fija ese
contrato para que el usuario sepa, sin abrir Notion, si su mensaje se
registró y con qué datos, o por qué no.

## Criterios de aceptación

- [ ] Registro exitoso → el bot responde confirmando: tipo (Egreso/
      Ingreso), concepto, monto, medio de pago, categoría (si hubo), y
      un link directo a la página creada en Notion (usa el page ID que
      ya devuelve `notion.Client.CreateExpense`).
- [ ] Mensaje que no empieza con `egreso`/`ingreso` → el bot responde
      explicando el formato esperado, con un ejemplo.
- [ ] Medio de pago no reconocido → el bot responde listando los
      medios de pago válidos (la tabla de `spec.md` de 001).
- [ ] Categoría no reconocida → el bot responde listando las
      categorías válidas (consultadas en vivo contra Notion, no una
      lista fija en el mensaje de error).
- [ ] Falla al escribir en Notion (red, API caída, rate limit) → el
      bot responde con un mensaje de error genérico y claro ("no se
      pudo guardar, intenta de nuevo"), **sin** exponer detalles
      internos (stack trace, token, URL de la API) en el chat.
- [ ] Todos los mensajes son texto plano en español, sin emojis (salvo
      que el usuario pida lo contrario más adelante).

## Fuera de alcance

- Botones inline / teclados de Telegram (`InlineKeyboardMarkup`).
- Internacionalización (todo en español, ya es el idioma del proyecto).
- Reintentos automáticos ante falla de Notion (eso sería una feature de
  resiliencia aparte, no de copy/UX).
