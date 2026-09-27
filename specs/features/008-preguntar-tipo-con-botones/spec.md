# 008 · Preguntar Ingreso/Egreso con botones en fotos sin caption

**Estado:** implementado ✅ (validado en vivo: botones aparecieron, toque de "Ingreso" interpretó y registró correcto con `Type: Ingreso`)

**Depende de:** 006 y 007 (mismo camino de foto sin caption).

## Qué hace

007 fija `Tipo: Egreso` siempre para fotos sin caption, porque una IA
no puede inferir con confianza si un recibo es un ingreso o un egreso.
Esta feature reemplaza esa suposición: cuando llega una foto sin
caption, el bot primero pregunta "¿Es un ingreso o un egreso?" con dos
botones (Ingreso / Egreso) — el usuario **no escribe nada**, solo toca
uno. Con la respuesta, recién ahí se interpreta la foto con IA (007) y
se registra con el `Tipo` que el usuario eligió.

## Por qué

El usuario lo pidió explícitamente: quiere confirmar el tipo sin tener
que escribir texto, para que el camino de "solo mandar la foto" siga
siendo tan simple como sea posible.

## Criterios de aceptación

- [ ] Foto sin caption → el bot responde con un mensaje con dos
      botones: "Ingreso" y "Egreso". Todavía no interpreta la imagen
      ni registra nada.
- [ ] El usuario toca un botón → el bot interpreta la foto con IA
      (mismo camino de 007: monto obligatorio, categoría/medio de pago
      best-effort) y registra el movimiento con el `Tipo` elegido.
- [ ] El mensaje de los botones se actualiza (edit) para reflejar qué
      se eligió, en vez de quedar con botones tocables sueltos.
- [ ] Si el usuario manda una foto nueva sin caption mientras una
      pregunta anterior seguía sin responder, la pregunta vieja queda
      obsoleta — la respuesta relevante es la de la foto más reciente
      (ver "Fuera de alcance" sobre qué pasa si igual la tocan).
- [ ] Foto **con** caption no cambia — sigue sin preguntar nada, el
      `Tipo` sale del caption como siempre (006).
- [ ] Mensaje de texto normal (sin foto) no cambia — sigue sin
      preguntar nada.

## Fuera de alcance

- Persistir la foto pendiente entre reinicios del bot — vive en
  memoria del proceso; si el bot se reinicia con una pregunta sin
  responder, esa foto se pierde (el usuario tendría que reenviarla).
- Un usuario tocando una pregunta vieja/obsoleta después de mandar una
  foto nueva — comportamiento no garantizado, se documenta como
  limitación conocida, no se valida explícitamente.
- Más de un chat con una foto pendiente en simultáneo — no aplica hoy
  (la whitelist tiene un solo `chat_id`), pero el diseño ya lo soporta
  al guardar la foto pendiente por `chat_id`.
- Botones para otras decisiones (categoría, medio de pago) — solo
  Tipo, que es lo que se pidió.
