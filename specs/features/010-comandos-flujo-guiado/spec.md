# 010 · Comandos /ingreso y /egreso con flujo guiado

**Estado:** implementado ✅ (validado en vivo: flujo completo paso a paso, registrado correcto en Notion, categoría omitida con "-")

## Qué hace

Dos comandos nuevos de Telegram, `/ingreso` y `/egreso`, que arrancan
una conversación paso a paso en vez de esperar una sola línea con todo
el formato. El bot pregunta un campo a la vez: monto → concepto
(opcional) → medio de pago → categoría (opcional), y al final registra
el movimiento — mismo resultado que mandar el mensaje de una línea
(001), pero de a un dato por vez.

El formato de una sola línea (001) y las fotos (006/007/008) siguen
funcionando exactamente igual — esto es un camino adicional, no un
reemplazo.

## Por qué

El usuario lo pidió explícitamente: quiere poder usar comandos tipo
`/ingreso`/`/egreso` en vez de recordar el formato completo de una
línea.

## Flujo

```
Usuario: /egreso
Bot:     Registrando un Egreso.
         ¿Cuál es el monto?
Usuario: 13000
Bot:     ¿Cuál es el concepto? (mandá "-" para omitirlo)
Usuario: mecato
Bot:     ¿Cuál fue el medio de pago? Opciones: efectivo, tarjeta de
         credito, tarjeta de debito, banco, transferencia
Usuario: efectivo
Bot:     ¿Cuál es la categoría? (mandá "-" para omitirla)
Usuario: comida
Bot:     Registrado: Egreso
         Concepto: mecato
         Monto: 13000
         Medio de pago: Cash
         Categoría: comida
         Notion: https://notion.so/...
```

## Criterios de aceptación

- [ ] `/egreso` (o `/ingreso`) sin nada más arranca la conversación y
      pregunta el monto primero.
- [ ] Monto no numérico → pide que lo reenvíe, no avanza de paso.
- [ ] Concepto: `-` lo omite (usa el nombre del tipo por defecto, igual
      que 001); cualquier otro texto se usa tal cual.
- [ ] Medio de pago no reconocido → pide que lo reenvíe, con la lista
      de opciones válidas, no avanza de paso.
- [ ] Categoría no reconocida → pide que la reenvíe, con la lista de
      categorías reales de Notion, no avanza de paso (a diferencia del
      resto de campos, si esto fallara y abortara la conversación se
      perdería todo lo ya contestado — por eso este paso permite
      reintentar en vez de cancelar).
- [ ] Categoría: `-` la omite.
- [ ] Al completar los 4 pasos, el registro es idéntico (mismos campos
      en Notion) al que produciría el formato de una línea equivalente.
- [ ] Mandar `/egreso` o `/ingreso` de nuevo mientras una conversación
      ya está en curso la reinicia (no hace falta un comando de
      cancelar aparte).
- [ ] Mandar una foto, `resumen`, o un mensaje de una línea completo
      mientras hay una conversación pendiente interrumpe la
      conversación (no queda una pregunta vieja esperando una
      respuesta que nunca va a llegar en ese formato).
- [ ] `/egreso` e `/ingreso` aparecen en el menú de comandos de
      Telegram (al escribir "/" en el chat).

## Fuera de alcance

- Editar una respuesta ya dada dentro de la misma conversación (ej.
  "esperá, el monto era otro") — hay que reiniciar con el comando de
  nuevo.
- Comando de cancelar explícito (`/cancelar`) — reiniciar con
  `/egreso`/`/ingreso` de nuevo ya cumple esa función.
- Adjuntar foto dentro de la conversación guiada — para eso ya existen
  006/007/008 (foto con o sin caption).
- Persistir la conversación entre reinicios del bot — vive en memoria
  del proceso, igual que `pendingPhotos` de la feature 008.
