# 2026-09-19 — Documentación: adaptar patrón de FDC/DESARROLLO

## Contexto

Usuario pidió replicar la estructura de memoria/contexto usada en
`/Users/michaelestivenreyesescobar/dev/moveeng/FDC/DESARROLLO/`
(`AGENT.md`, `CLAUDE.md`, `.claude/rules/`, `.claude/history/`,
`.claude/plans/`), adaptándola a este proyecto en Go.

## Decisiones

- `agents.md` (raíz, minúscula) reemplazado por `AGENT.md` (mayúscula,
  igual que la referencia) — manual estático de arquitectura, stack,
  comandos, estructura, convenciones, "no hagas" y flujo de trabajo.
- `CLAUDE.md` reescrito como bitácora corta: resumen rápido + índice de
  `.claude/history/` + última entrada completa (antes era solo un
  router liviano con status).
- `.claude/rules/` ganó 5 archivos modulares calcados del patrón de FDC
  y adaptados a Go: `comunicacion.md`, `codigo.md`,
  `flujo_de_trabajo.md`, `seguridad.md`, `bitacora.md`. Se mantienen
  además `telegram-rules.md` y `notion-rules.md` (específicas de cada
  API externa) porque este proyecto, a diferencia de FDC, integra dos
  APIs de terceros desde el día uno — no es una desviación del patrón,
  es una extensión justificada por dominio.
- `flujo_de_trabajo.md` NO replica el sistema `specs/features/NNN-...`
  de FDC (spec-driven con Jira): no fue pedido para este proyecto y
  agregarlo sin pedido sería alcance no solicitado. Se mantiene sí la
  regla equivalente más simple: plan corto antes de cambios no
  triviales, sin implementar hasta indicación explícita.
- `.claude/history/` renombrado a la convención exacta de FDC:
  `YYYY-MM-DD-feature-que_se_realizo.md` (antes:
  `2026-09-19-01-setup-arquitectura.md`, ahora:
  `2026-09-19-setup-andamiaje_go_clean_architecture.md`).
- `.claude/plans/` creado vacío, igual que en la referencia.
- No se replicó `.claude/commands/` (comandos custom tipo `/boton`) —
  la referencia tiene uno específico de su dominio (Spring/Angular);
  no hay pedido de un comando equivalente para este proyecto todavía.

## Estado al cierre de sesión

- Go 1.27.1 instalado (`brew install go`), `go build ./...` y `go vet
  ./...` pasan limpio sobre el stub actual.
- Documentación/reglas alineadas al patrón FDC. Código sigue siendo
  stub (`TODO`) — sin lógica de negocio, sin clientes reales de
  Telegram/Notion. Eso sigue pendiente de indicación explícita del
  usuario.
