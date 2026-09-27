# Estilo de código

- **Antes de implementar código, siempre cargar la skill `/antislop`**
  (núcleo) y, para código y comentarios, `/antislop-code`. Aplica a
  cualquier tarea que escriba o modifique código (features, fixes,
  refactors, tests), no solo a las grandes. Cargarla antes de escribir
  la primera línea, no al final como revisión. Si la tarea también toca
  texto de cara al usuario (mensajes del bot, copy), sumar
  `/antislop-copywriting`. Sus reglas se suman a las de este archivo; si
  chocan, gana la más estricta.
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
- Arquitectura por feature con puertos (patrón de Triggo/WebHook):
  `pkg/<feature>/{handler,model,services}` + `pkg/ports/`. `handler` y
  `services` dependen de `pkg/ports` y de `model`, nunca entre sí ni
  de otra feature; `model` solo importa la stdlib y otros `model`;
  `pkg/ports` solo importa `model`. Solo `pkg/app`, `api/` y `cmd/`
  conocen las implementaciones concretas y las cablean.
- Errores: envolver con contexto (`fmt.Errorf("...: %w", err)`), nunca
  tragarlos en silencio.
