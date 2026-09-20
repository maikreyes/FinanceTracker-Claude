# 005 · Consultar resumen de gastos

**Estado:** propuesta

## Qué hace

El usuario manda el mensaje `resumen` (case-insensitive) y el bot
responde con el total de Egresos y el total de Ingresos del mes
calendario actual, sin necesidad de abrir Notion.

## Por qué

Consultar el estado del mes es el segundo caso de uso más obvio de un
tracker financiero por chat, después de registrar movimientos.

## Criterios de aceptación

- [ ] Mensaje `resumen` (o `Resumen`, `RESUMEN` — case-insensitive) →
      el bot responde con dos números: total Egresos y total Ingresos
      del mes calendario en curso (fecha del servidor).
- [ ] Si no hay movimientos este mes → responde 0 para ambos, no
      error.
- [ ] Solo responde a chats autorizados (reusa la whitelist de la
      feature 003).
- [ ] El cálculo usa los datos reales de la data source `Expenses`
      (no un valor cacheado desactualizado).

## Fuera de alcance

- Resumen por categoría o por medio de pago.
- Resumen de meses anteriores o rango de fechas custom.
- Gráficos o visualizaciones — solo texto.
- Otros comandos de consulta (ej. "últimos 5 movimientos") — se evalúan
  como features aparte si se necesitan.
