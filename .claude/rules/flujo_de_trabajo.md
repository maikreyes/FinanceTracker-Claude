# Flujo de trabajo

1. Historias de usuario → `specs/features/NNN-nombre-feature/`
   (`spec.md`, `plan.md`, `tasks.md`) antes de tocar código. Crear la
   spec no implica empezar a implementar.
2. No implementar hasta que el usuario lo indique explícitamente.
3. Plan primero, siempre, para tareas no triviales (tocan el contrato
   entre capas, agregan una dependencia externa, cambian la entidad
   `Transaction`) — proponer plan corto, esperar OK. Cambios chicos y
   acotados a un solo paquete: proceder directo y reportar al final.
4. Una spec/tarea a la vez. Al terminar, resumir qué cambió y por qué.
5. Gate de tests antes de dar algo por completado: `go build ./...`,
   `go vet ./...` y `go test ./...` (cuando haya tests) deben pasar en
   los paquetes tocados.
6. Git es del usuario: nunca crear ramas, `git commit`, `push` ni PR's.
   El trabajo termina en dejar el working tree listo y un resumen claro.
7. No agregar dependencias externas (`go get`) sin avisar antes y
   explicar por qué.
8. Si no hay 80% de seguridad, preguntar — no inventar reglas de negocio
   ni asumir alcance no pedido, sobre todo en el mapeo de propiedades de
   Notion y en el formato esperado de los mensajes de Telegram.
