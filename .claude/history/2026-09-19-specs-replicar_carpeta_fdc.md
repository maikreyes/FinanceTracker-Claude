# 2026-09-19 — specs/: replicar carpeta de FDC/DESARROLLO

## Contexto

Usuario pidió replicar también `specs/` de
`moveeng/FDC/DESARROLLO/` (workflow spec-driven: `constitution/` +
`features/NNN-nombre-feature/`). Revierte la decisión de la entrada
anterior ("documentacion-adaptar_patron_fdc"), que había excluido
deliberadamente este sistema por no haber sido pedido — ahora sí se
pidió explícitamente.

## Decisiones

- `specs/constitution/mission.md` — copiado tal cual como template
  (placeholders `<>` sin llenar), igual que en la referencia: necesita
  input de negocio del usuario (para quién, principios, qué NO es), no
  algo que se pueda inventar por contexto de código.
- `specs/constitution/roadmap.md` — llenado con estado real: única
  entrada en "Hecho" es el andamiaje inicial (sin número de feature,
  es setup); "Siguiente" queda abierto a que el usuario defina la
  primera historia de usuario real (001).
- `specs/constitution/tech-stack.md` — llenado, derivado 1:1 de
  `AGENT.md` (stack Go, estructura, comandos, modelo de dominio,
  convenciones, límites duros).
- `specs/features/NNN-nombre-feature/{spec,plan,tasks}.md` — copiados
  tal cual como templates (idénticos a la referencia, son genéricos por
  diseño, no específicos de Spring/Angular).
- `AGENT.md` y `.claude/rules/flujo_de_trabajo.md` actualizados: el
  flujo de trabajo ahora empieza con "historias de usuario →
  specs/features/NNN-.../ antes de tocar código", igual que FDC.
- No se creó ninguna feature numerada real (`001-...`) — no hay
  historia de usuario definida todavía; el usuario debe darla primero.

## Estado al cierre de sesión

- `specs/` completo y alineado al patrón FDC. `mission.md` pendiente de
  llenar (requiere input de negocio del usuario, igual que en la
  referencia). Código sigue stub, sin lógica de negocio.
