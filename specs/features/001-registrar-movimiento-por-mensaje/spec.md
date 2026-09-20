# 001 · Registrar movimiento por mensaje de Telegram

**Estado:** implementado ✅ (parser + resolución de categoría + escritura en Notion; falta 002 para exponerlo por Telegram de verdad)

## Qué hace

El usuario le manda al bot un mensaje de texto corto con un formato
fijo y el bot lo registra como una fila nueva en la data source
`Expenses` de Notion (gasto o ingreso).

Formato del mensaje (posicional, en este orden):

```
<Tipo> [Concepto] <Monto> <MedioDePago> [Categoria]
```

Ejemplos reales dados por el usuario:

```
Egreso Gasolina 5000 tarjeta credito
Egreso Gasolina 5000 tarjeta debito
Egreso mecato 13000 efectivo
egreso mecato 13000 tranferencia
ingreso 3000000 banco
ingreso 12000 efectivo
egreso mecato 13000 efectivo comida
```

- `Tipo`: `egreso` o `ingreso` (case-insensitive) → `domain.TransactionType`.
- `Concepto`: 0 o más palabras entre `Tipo` y `Monto` → `Name`. Si no
  hay ninguna (como en los ejemplos de ingreso), se usa `Tipo`
  capitalizado como `Name` por defecto (ej. "Ingreso").
- `Monto`: el primer token numérico del mensaje → `Amount`. Siempre
  positivo; el signo lo da `Type`, no el número.
- `MedioDePago`: 1 o 2 palabras justo después del monto, mapeadas por
  palabra clave (ver tabla) → `Payment Method`.
- `Categoria`: opcional, 1 o más palabras al final, deben coincidir
  (case-insensitive) con el nombre exacto de una categoría existente en
  la data source `Category` de Notion → `Category` (relation).

### Mapeo de medio de pago (palabra clave → opción de Notion)

| Palabra(s) en el mensaje | `Payment Method` en Notion |
|---|---|
| `efectivo` | `Cash` |
| `tarjeta credito` / `tarjeta de credito` | `Credit Card` |
| `tarjeta debito` / `tarjeta de debito` | `Debit Card` |
| `banco` | `Bank` |
| `tranferencia` / `transferencia` | `Bank` |

### Categorías válidas (leídas de Notion, no hardcodeadas en el spec)

Al momento de escribir esta spec: `Entretenimiento`, `Subcripciones`,
`Seguridad Social`, `Comida`, `Ahorro`, `Servicios`. La lista puede
cambiar en Notion — el resolver de categoría consulta la data source
`Category` en vez de usar una lista fija en código (ver `plan.md`).

## Por qué

Es el flujo principal del producto: registrar un movimiento financiero
sin abrir Notion, desde Telegram, en una sola línea de texto.

## Criterios de aceptación

- [ ] Un mensaje con el formato completo (`Tipo Concepto Monto
      MedioDePago Categoria`) crea una página en `Expenses` con todas
      las propiedades correctas.
- [ ] Un mensaje sin `Categoria` (como los primeros 6 ejemplos) crea la
      página igual, con `Category` vacío.
- [ ] Un mensaje sin `Concepto` (ej. `ingreso 3000000 banco`) usa
      `Tipo` capitalizado como `Name`.
- [ ] `Tipo`, `MedioDePago` y `Categoria` no distinguen mayúsculas ni
      minúsculas (`Egreso` = `egreso`).
- [ ] Un `MedioDePago` no reconocido → el bot responde con error claro,
      no crea la página.
- [ ] Una `Categoria` que no coincide con ninguna categoría real de
      Notion → el bot responde con error claro (lista las categorías
      válidas), no crea la página.
- [ ] Un mensaje sin `Tipo` válido (`egreso`/`ingreso`) al inicio → el
      bot lo ignora o responde que no reconoce el formato (no revienta).
- [ ] `Date` se registra como la fecha en que Telegram recibió el
      mensaje (no hay fecha en el formato del mensaje).

## Fuera de alcance

- Fecha explícita en el mensaje (ej. "ayer", "15 de marzo") — se usa
  siempre la fecha de recepción del mensaje.
- Adjuntar foto de comprobante (`Receipt`) — no está en el formato del
  mensaje.
- Crear una categoría nueva si no existe — el usuario debe usar una
  categoría ya creada en Notion.
- Edición o borrado de un movimiento ya registrado, desde Telegram.
- Restricción de quién puede usar el bot (chat_id whitelist) — decisión
  de seguridad pendiente, se resuelve en la implementación del
  transporte de Telegram, no en el parser.
