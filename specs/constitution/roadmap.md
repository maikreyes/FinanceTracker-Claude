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

## Siguiente 🔜

_Lo próximo a abordar. Idealmente una sola feature "en curso" a la vez._

Orden pensado: 002 → 003 → 004 → 005 (cada una depende de la anterior
para ser usable de punta a punta). Todas con spec/plan/tasks listos,
ninguna implementada todavía — a la espera de OK explícito.

1. **[002 · Transporte Telegram](../features/002-transporte-telegram/spec.md)** —
   long-polling real, conecta 001 con Telegram de verdad.
2. **[003 · Restringir acceso al bot](../features/003-restringir-acceso-bot/spec.md)** —
   whitelist de `chat_id`, fail closed.
3. **[004 · Confirmación y manejo de errores en el chat](../features/004-confirmacion-manejo-errores/spec.md)** —
   copy exacto de éxito/error que responde el bot.
4. **[005 · Consultar resumen de gastos](../features/005-consultar-resumen-gastos/spec.md)** —
   comando `resumen`, totales Egreso/Ingreso del mes actual.

## Backlog / ideas 💡

_Sin comprometer ni ordenar del todo. Ideas que respetan la constitución._

- **<Nombre>** — <qué aportaría>.
- **<Nombre>** — <qué aportaría>.

> Cada feature nueva se crea como `features/NNN-nombre-feature/` con `spec.md`, `plan.md` y `tasks.md` antes de tocar código.
