# 007 · Interpretar recibo con IA (sin caption)

**Estado:** implementado ✅ (validado en vivo con una foto real de
Telegram — monto/concepto/categoría interpretados correctamente,
medio de pago dejado vacío por incertidumbre, sin inventar)

**Depende de:** 006 (usa su misma infraestructura de descarga de foto de
Telegram — `Message.Photo`, `Bot.GetFile`/`DownloadFile`).

## Qué hace

Si el usuario manda una foto **sin caption**, en vez de pedirle que la
reenvíe con el texto del movimiento (comportamiento actual de 006), el
bot le pasa la imagen a un modelo de IA con visión (Groq, tier
gratuito) para que interprete el recibo y extraiga los datos del
movimiento automáticamente. El caption sigue siendo el camino
preferido cuando existe — más rápido, gratis, determinístico, sin
depender de un proveedor externo. La IA es el *fallback* para cuando
no se quiere escribir el caption.

## Por qué

006 resuelve "adjuntar la foto", pero todavía obliga a escribir el
mismo texto de siempre como caption. El usuario quiere mandar solo la
foto cuando le resulte más cómodo.

## Qué extrae la IA del recibo

- **Monto** (obligatorio) — si la IA no logra leer un monto claro, se
  trata como si no se pudiera interpretar el recibo (ver criterios).
- **Concepto** — nombre del comercio o descripción breve del recibo.
- **Categoría** — la IA elige entre las categorías reales del usuario
  (se le pasa la lista actual de Notion en el prompt); si no está
  segura, la deja vacía. No inventa categorías nuevas.
- **Medio de pago** — si el recibo muestra claramente cómo se pagó
  (ej. "VISA", "EFECTIVO"), lo mapea a una de las 4 opciones válidas;
  si no es claro, lo deja vacío.
- **Tipo** — siempre `Egreso`. Un recibo fotografiado es, por
  definición, un gasto; esta feature no interpreta ingresos.

## Criterios de aceptación

- [ ] Foto sin caption → el bot interpreta la imagen con Groq y
      registra el movimiento igual que 006 (con la foto adjunta en
      `Receipt`), en vez de pedir que se reenvíe con caption.
- [ ] Foto **con** caption sigue usando el camino de 006 (texto del
      caption, sin llamar a la IA) — el caption manual tiene
      prioridad, es más confiable.
- [ ] La respuesta de confirmación dice explícitamente que el monto
      /concepto/categoría fueron interpretados por IA, para que el
      usuario sepa que debe revisar los datos en Notion si algo se ve
      raro (a diferencia de un registro por texto, que es exacto).
- [ ] Si la IA no logra identificar un monto en la imagen → el bot
      responde pidiendo que se reenvíe la foto con el caption manual
      (cae de vuelta al camino confiable de 006), no registra nada a
      medias.
- [ ] Si la llamada a la API de Groq falla (red, cuota agotada,
      error del servicio) → error claro, sugiere reenviar con caption
      mientras tanto; no crashea el bot.
- [ ] Categoría o medio de pago que la IA no pudo determinar → el
      movimiento se registra igual, con esos campos vacíos (mismo
      principio que el resto del bot: no inventar, mejor dejar vacío).

## Fuera de alcance

- Interpretar ingresos desde una foto (transferencias, comprobantes de
  pago recibido) — solo egresos vía este camino.
- Corregir o re-interpretar un registro ya hecho por IA.
- Elegir entre múltiples proveedores de IA en runtime — queda fijo en
  Groq para esta feature; cambiar de proveedor sería un cambio de
  `plan.md`, no una opción configurable.
- OCR de recibos en mal estado, borrosos o en otros idiomas — mejor
  esfuerzo, sin garantía de precisión (es justamente el motivo de que
  el mensaje de confirmación avise que fue interpretado por IA).
