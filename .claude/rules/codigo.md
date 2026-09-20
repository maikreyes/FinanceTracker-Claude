# Estilo de código

- Comentarios solo si son necesarios; por defecto no comentar. Un
  comentario se justifica solo cuando explica un WHY no obvio
  (restricción oculta, workaround, efecto secundario no evidente) —
  nunca para decir QUÉ hace el código si el nombre ya lo dice.
- El comentario se redacta desde el contexto técnico del código mismo,
  nunca desde la tarea que lo originó. **Jamás** escribir "esto es para
  la Fase 3" o "agregado para el setup inicial".
  - Correcto: `// token leído de env var porque Telegram lo exige por
    request, no se cachea en memoria`
  - Incorrecto: `// implementado para la Fase 3 del setup`
- Nombre de paquete = nombre de carpeta (convención idiomática de Go).
- Dirección de dependencia: `infrastructure -> usecase -> domain`.
  `internal/domain/` sin imports fuera de la stdlib; `internal/usecase/`
  nunca importa tipos concretos de `internal/infrastructure/`.
- Errores: envolver con contexto (`fmt.Errorf("...: %w", err)`), nunca
  tragarlos en silencio.
