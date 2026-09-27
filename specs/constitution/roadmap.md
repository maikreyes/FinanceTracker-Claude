# Roadmap

_Orden y estado de las features. Es la vista de "qué hay hecho, qué toca ahora y qué viene". Cada entrada apunta a su carpeta en `features/`._

## Hecho ✅

_Features completadas, en orden de implementación._

1. **Andamiaje inicial** — módulo Go (`finance-tracker`), estructura
   Clean Architecture (`domain`/`usecase`/`infrastructure`), stubs de
   paquete, `.gitignore`, `.env.example`. Documentación replicada del
   patrón `moveeng/FDC/DESARROLLO`: `AGENT.md`, `CLAUDE.md`,
   `.claude/rules/`, `.claude/history/`, `.claude/plans/`,
   `specs/constitution/`, `specs/features/NNN-nombre-feature/`. Sin
   feature numerada todavía — es setup, no una historia de usuario.
2. **[001 · Registrar movimiento por mensaje de Telegram](../features/001-registrar-movimiento-por-mensaje/spec.md)** —
   parser posicional (`Tipo Concepto Monto MedioDePago Categoria`) +
   resolución de categoría contra Notion (`CategoryResolver`) +
   escritura en `Expenses` (`Registrar`). Validado con tests unitarios
   y un smoke test real end-to-end contra el Notion del usuario. Falta
   002 para exponerlo por Telegram de verdad — hoy solo es código Go
   invocable, sin transporte.
3. **[002 · Transporte Telegram](../features/002-transporte-telegram/spec.md)** —
   long-polling real con `net/http` puro, wiring completo en
   `cmd/bot/main.go`. Probado con mensajes reales de Telegram (uno
   inválido, uno completo con categoría) — ambos casos funcionaron.
4. **[003 · Restringir acceso al bot](../features/003-restringir-acceso-bot/spec.md)** —
   whitelist `TELEGRAM_ALLOWED_CHAT_IDS` (fail closed), chequeo en
   `cmd/bot/main.go` antes de despachar al usecase. Validado en vivo:
   chat autorizado registra en Notion, chat no autorizado se ignora en
   silencio (sin respuesta). El bot ya no acepta cualquier chat.
5. **[004 · Confirmación y manejo de errores en el chat](../features/004-confirmacion-manejo-errores/spec.md)** —
   `successReply`/`errorReply` en `cmd/bot/main.go`, errores tipados
   con `ErrUnrecognizedCategory.Valid` poblado en vivo contra Notion.
   Validado con el bot corriendo: error de medio de pago inválido
   respondido correcto, registro exitoso con categoría confirmado.
6. **[005 · Consultar resumen de gastos](../features/005-consultar-resumen-gastos/spec.md)** —
   `SumThisMonth` (query con filtro de fecha + paginación, suma por
   `Type`, zona horaria `America/Bogota` explícita) + comando `resumen`
   en el bot. Validado en vivo: respondió `Egresos: 0, Ingresos: 0`,
   confirmado independientemente contra Notion como el resultado
   correcto (filas del mes históricas sin `Type`, o páginas de prueba
   ya archivadas por el usuario).
7. **[006 · Adjuntar foto de recibo al registro](../features/006-adjuntar-recibo-foto/spec.md)** —
   foto con caption = mensaje normal, la foto se sube al File Upload
   API de Notion (`Client.Upload`, dos pasos: crear + enviar multipart)
   y se referencia en `Receipt`. Validado en vivo: foto sin caption →
   pide reenviar; foto con caption → registrada con el archivo real
   adjunto (`recibo.jpg`, `type: "file"`, no un link externo).
8. **[007 · Interpretar recibo con IA (sin caption)](../features/007-interpretar-recibo-con-ia/spec.md)** —
   depende de 006. Foto **sin** caption → Groq (tier gratuito; se
   descartó Google Gemini porque pedía vincular facturación)
   interpreta el recibo y extrae monto/concepto/categoría/medio de
   pago. Caption manual sigue teniendo prioridad cuando existe.
   Validado en vivo con un recibo real (recarga Nequi/Bancolombia):
   monto/concepto/categoría correctos, medio de pago vacío por
   incertidumbre (no inventó). Modelo vigente en Groq:
   `qwen/qwen3.8-27b` (el original del plan estaba dado de baja).
9. **[008 · Preguntar Ingreso/Egreso con botones](../features/008-preguntar-tipo-con-botones/spec.md)** —
   reemplaza el `Tipo: Egreso` fijo de 007: foto sin caption → el bot
   pregunta con botones inline ("Ingreso"/"Egreso", `callback_query`
   de Telegram) antes de interpretar con IA — el usuario no escribe
   nada. Foto guardada en memoria por `chat_id` mientras se espera la
   respuesta. Validado en vivo: botones aparecieron, toque de
   "Ingreso" registró correcto (`Type: Ingreso`) sin errores.
10. **[009 · Balance global en resumen + backfill de Tipo](../features/009-balance-global-resumen/spec.md)** —
    `resumen` ahora incluye `Balance = Ingreso - Egreso` ("cuánto
    tengo disponible" — primer intento fue al revés, corregido en el
    mismo turno) y cuántos movimientos son de cada tipo, no solo la
    suma. Backfill real: las 31 páginas históricas sin `Type`
    (anteriores a que existiera la propiedad) se marcaron como
    `Egreso`. Validado en vivo: `Egresos: 2338155 (30 movimientos),
    Ingresos: 0 (0 movimientos)`, confirmado por el usuario contra un
    cálculo independiente (`Balance` con la fórmula final:
    `-2338155`).
11. **[010 · Comandos /ingreso y /egreso con flujo guiado](../features/010-comandos-flujo-guiado/spec.md)** —
    dos comandos nuevos que arrancan una conversación paso a paso
    (monto → concepto → medio de pago → categoría) en vez de una sola
    línea con todo el formato; registrados en el menú de comandos de
    Telegram (`SetMyCommands`). `RegisterFromFields` arma la
    transacción directo de los campos ya separados (sin reparsear un
    string reconstruido). Solo el paso de categoría permite reintentar
    sin perder lo ya contestado. Validado en vivo: flujo completo,
    registrado en Notion con todos los campos correctos, categoría
    omitida con "-" como se pidió.
12. **[011 · Comandos /ayuda, /ultimo y /categorias](../features/011-comandos-ayuda-ultimo-categorias/spec.md)** —
    `/ayuda` (texto fijo de formato y comandos), `/ultimo` (último
    movimiento registrado, resolviendo el nombre de la categoría) y
    `/categorias` (botones "Listar"/"Agregar categoría", mismo patrón
    visual que los botones Ingreso/Egreso de 008; agregar no duplica
    si el nombre ya existe). Validado en vivo: `/ayuda` y
    `/categorias` (listar + agregar) correctos a la primera; `/ultimo`
    falló por timeout de red en el primer intento (no un bug —
    reintento exitoso, movimiento correcto).

## Siguiente 🔜

_Las 11 features especificadas (001-011) están implementadas._

- **[012 · Migración de arquitectura para desplegar en Vercel](../features/012-migracion-arquitectura-vercel/spec.md)** —
  en progreso 🚧. Paso 1 listo: lógica de despacho movida a
  `internal/bot` (`Dispatch`, `Store`), `memorystore` para desarrollo
  local, validado en vivo (mismo comportamiento que antes del
  refactor). Faltan: store de Upstash Redis, webhook de Vercel,
  `cmd/setwebhook`, documentación de despliegue. El despliegue real a
  Vercel queda fuera de lo que se puede validar en esta sesión
  (requiere cuenta de Vercel/Upstash del usuario).

## Backlog / ideas 💡

_Sin comprometer ni ordenar del todo. Ideas que respetan la constitución._

- **<Nombre>** — <qué aportaría>.
- **<Nombre>** — <qué aportaría>.

> Cada feature nueva se crea como `features/NNN-nombre-feature/` con `spec.md`, `plan.md` y `tasks.md` antes de tocar código.
