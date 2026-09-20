# Bitácora (`.claude/history/`)

- Un archivo por entrada, nombre exacto:
  `YYYY-MM-DD-feature-que_se_realizo.md`. `feature` = slug corto del
  módulo o feature tocado (ej. `setup`, `domain`, `telegram`, `notion`,
  `usecase`); `que_se_realizo` = qué se hizo, snake_case corto.
- Contenido: decisiones, hallazgos y contexto que no está en el código
  ni en `AGENT.md` — para que una sesión nueva arranque con el mismo
  contexto. No reescribir historial ya registrado; una entrada nueva por
  evento que valga la pena recordar.
- `CLAUDE.md` (raíz) solo guarda: índice de una línea por entrada (con
  link a su archivo en `.claude/history/`), marcando ⭐ las claves para
  releer, y la última entrada completa. Al agregar una entrada nueva:
  crear su archivo, agregar la línea al índice, mover la entrada
  completa anterior a solo el índice (ya vive en su archivo) y poner la
  nueva como "Última entrada completa".
