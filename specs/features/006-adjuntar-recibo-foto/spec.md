# 006 · Adjuntar foto de recibo al registro

**Estado:** implementado ✅ (validado en vivo: foto sin caption → pide reenviar; foto con caption → registrada con el archivo real adjunto en `Receipt`)

## Qué hace

El usuario manda una foto (de un recibo/comprobante) **con el mensaje
normal como caption** de esa foto (`egreso mecato 13000 efectivo
comida`). El bot registra el movimiento igual que hoy y además adjunta
la foto a la propiedad `Receipt` de esa misma página en Notion.

## Por qué

`Receipt` ya existe en el schema real de Notion (ver
`.claude/rules/notion-rules.md`) pero la feature 001 la dejó
explícitamente fuera de alcance porque el formato de texto no traía
manera de adjuntar un archivo. Esta feature cierra ese hueco sin
inventar un formato nuevo — reusa el caption de la foto como el mismo
mensaje de siempre.

## Criterios de aceptación

- [ ] Un mensaje que es una foto con caption en formato válido
      (`Tipo [Concepto] Monto MedioDePago [Categoria]`) crea la página
      en `Expenses` igual que un mensaje de texto normal, **y** la
      propiedad `Receipt` de esa página queda con la foto adjunta.
- [ ] La foto se sube a Notion de verdad (no un link externo temporal
      de Telegram) — sigue siendo visible desde Notion después de que
      el archivo expire del lado de Telegram.
- [ ] Una foto **sin** caption → el bot responde pidiendo que se
      reenvíe con el texto del movimiento como caption; no crea página.
- [ ] Una foto con caption en formato inválido → mismo manejo de error
      que ya existe para texto (feature 001/004), no crea página ni
      sube la foto.
- [ ] Un mensaje de texto normal (sin foto) sigue funcionando exacto
      igual que antes — esta feature no cambia el camino sin foto.
- [ ] Un mensaje con más de una foto (álbum/media group de Telegram) —
      ver "Fuera de alcance".

## Fuera de alcance

- Adjuntar la foto a un movimiento ya registrado antes (mandar la foto
  por separado, después del texto) — es la otra opción que se
  descartó al definir esta spec; se puede revisar como feature aparte
  si hace falta.
- Álbumes/media groups de Telegram (varias fotos en un solo envío) —
  se soporta una foto por mensaje.
- Documentos/PDFs como comprobante — solo fotos (`message.photo`), no
  `message.document`.
- Editar o reemplazar el `Receipt` de un movimiento ya creado.
- Comprimir/redimensionar la foto antes de subirla — se sube tal cual
  la entrega Telegram (la resolución más alta disponible).
