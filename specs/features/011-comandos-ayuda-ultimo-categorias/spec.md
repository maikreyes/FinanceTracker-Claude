# 011 · Comandos /ayuda, /ultimo y /categorias

**Estado:** implementado ✅ (validado en vivo: /ayuda, /categorias listar+agregar, /ultimo — este último con reintento por un timeout de red puntual en el primer intento)

## Qué hace

Tres comandos nuevos de Telegram:

- **`/ayuda`** — recuerda el formato de un mensaje de una línea, los
  comandos disponibles y un ejemplo. Respuesta fija, sin estado.
- **`/ultimo`** — muestra el último movimiento registrado en Notion
  (el más reciente por fecha de creación de la página), con sus datos
  principales y el link a la página.
- **`/categorias`** — manda dos botones inline: **Listar** (muestra
  las categorías reales que ya existen en Notion) y **Agregar
  categoría** (pregunta el nombre por texto y crea la página nueva en
  la data source `Category`). Mismo patrón de botones que la feature
  008 (Ingreso/Egreso sobre una foto sin caption).

## Por qué

El usuario preguntó qué otros comandos convenía tener antes de
desplegar a Vercel. Se recomendaron estos tres por ser los de mayor
valor con menor esfuerzo (no requieren estado nuevo salvo el paso de
"esperando nombre de categoría"). El usuario pidió implementarlos, más
que `/categorias` incluya los botones listar/agregar con el mismo
patrón visual que ya usa el bot.

## Flujo — /ayuda

```
Usuario: /ayuda
Bot:     Formato de un mensaje: tipo [concepto] monto medio_de_pago [categoria]
         Ejemplo: egreso mecato 13000 efectivo comida
         Tipos válidos: Egreso, Ingreso
         Medios de pago válidos: efectivo, tarjeta de credito, tarjeta de debito, banco, transferencia
         También podés:
         - Mandar una foto de un recibo (con o sin texto)
         - /egreso o /ingreso: registrar paso a paso
         - resumen: balance y totales del mes
         - /ultimo: ver el último movimiento
         - /categorias: ver o agregar categorías
```

## Flujo — /ultimo

```
Usuario: /ultimo
Bot:     Último movimiento:
         Tipo: Egreso
         Concepto: mecato
         Monto: 13000
         Medio de pago: Cash
         Categoría: Comida
         Fecha: 2026-09-20
         Notion: https://notion.so/...
```

Si no hay ningún movimiento registrado todavía: "Todavía no hay
movimientos registrados."

## Flujo — /categorias

```
Usuario: /categorias
Bot:     ¿Qué querés hacer? [Listar categorías] [Agregar categoría]

--- toca "Listar categorías" ---
Bot (edita el mensaje): Categorías:
         Ahorro
         Comida
         Transporte

--- toca "Agregar categoría" ---
Bot (edita el mensaje): ¿Cómo se llama la nueva categoría?
Usuario: Mascotas
Bot:     Categoría "Mascotas" creada.
```

Si el nombre ya existe (case-insensitive), no se crea una página
duplicada — se reusa la categoría existente y se confirma igual (el
usuario no necesita saber si ya existía o no).

## Criterios de aceptación

- [ ] `/ayuda` responde siempre el mismo texto fijo, sin llamar a
      Notion.
- [ ] `/ultimo` consulta Notion y muestra el movimiento más reciente
      (por fecha de creación de la página), con Categoría resuelta a
      nombre (no ID) cuando exista.
- [ ] `/ultimo` sin movimientos registrados responde con el mensaje de
      "no hay movimientos", no un error.
- [ ] `/categorias` manda los dos botones; "Listar" edita el mensaje
      con la lista real (o "Sin categorías todavía." si está vacía).
- [ ] `/categorias` → "Agregar categoría" edita el mensaje pidiendo el
      nombre, y el siguiente mensaje de texto del chat se interpreta
      como ese nombre (no como un movimiento nuevo ni otro comando).
- [ ] Agregar una categoría con un nombre que ya existe
      (case-insensitive) no crea una página duplicada.
- [ ] Los tres comandos aparecen en el menú de comandos de Telegram.
- [ ] Una foto, `resumen`, `/egreso`/`/ingreso`, o un mensaje de una
      línea mientras el bot espera el nombre de una categoría nueva
      cancela esa espera (mismo criterio que las conversaciones
      guiadas de la feature 010).

## Fuera de alcance

- Editar o eliminar una categoría existente.
- Editar o eliminar el último movimiento (`/ultimo` es de solo
  lectura).
- Paginar `/categorias` si la lista es muy larga — se manda completa
  en un solo mensaje.
- Confirmación explícita antes de crear la categoría ("¿segura?") — se
  crea directo con el nombre que mandó el usuario.
