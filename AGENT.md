# Finance Tracker: Telegram to Notion

Bot que lee mensajes en lenguaje natural desde Telegram, los parsea en
transacciones financieras y las escribe en una base de datos de Notion.

## Stack

- Lenguaje: Go 1.23+
- Módulo: `finance-tracker` (`go.mod`, raíz del repo)
- Dependencias externas: solo `github.com/joho/godotenv` (cargar `.env`
  en `cmd/`); todo lo demás es stdlib
- Arquitectura: por feature con puertos, el mismo patrón que
  `maikreyes/Triggo/WebHook` (`api/`, `cmd/`, `pkg/<feature>/{handler,model,services}`,
  `pkg/ports`)
- Puntos de entrada: `cmd/main.go` (long-polling, local) y
  `api/webhook.go` (función serverless de Vercel)

## Comandos

- `go build ./...` — compila todo el módulo
- `go run ./cmd` — corre el bot en local (long-polling, carga `.env`)
- `go run ./cmd/setwebhook -url https://<proyecto>.vercel.app/api/webhook`
  — apunta el webhook de Telegram a Vercel (`-delete` lo borra)
- `go vet ./...` — chequeo estático
- `go test -race ./...` — tests
- `gofmt -l .` — lista archivos mal formateados

## Estructura del proyecto

- `api/webhook.go` — `Handler`, la función serverless de Vercel. Cablea
  las dependencias una vez por instancia (`sync.Once`) con el store de
  Upstash y delega en `Handler.WebhookHandler`. `vercel.json` fija su
  `maxDuration`.
- `cmd/main.go` — long-polling para desarrollo local, con el store en
  memoria. `cmd/setwebhook/` — comando de un solo uso para
  configurar o borrar el webhook.
- `pkg/config/` — `Config` leído del entorno (`NewConfig`) y sus
  validaciones por modo (`CheckTelegram`, `CheckWebhook`).
- `pkg/ports/` — las interfaces, un archivo por dominio (`telegram.go`,
  `notion.go`, `groq.go`, `store.go`, `transaction.go`). Lo único de lo
  que dependen `handler` y `services`.
- `pkg/transaction/` — la lógica de negocio. `model/` (subpaquetes
  `transaction`, `register`, `summary`, `receipt`, `failure`) y
  `services/` (parseo de mensajes y registro de movimientos).
- `pkg/telegram/` — `model/` (`update`, `keyboard`, `command`,
  `chatstate`), `services/` (cliente de la Bot API, implementa
  `ports.Messenger`) y `handler/` (`Dispatch`, flujo guiado, botones,
  `WebhookHandler`).
- `pkg/notion/services/` — cliente de la Notion API (movimientos,
  recibos, resumen, último movimiento, categorías). Reglas en
  `.claude/rules/notion-rules.md`.
- `pkg/groq/services/` — interpreta la foto de un recibo con IA.
- `pkg/store/services/` — `ChatStore` en memoria (local) y en Upstash
  Redis (producción).
- `pkg/app/` — raíz de composición compartida por `api/` y `cmd/`.

## Flujo de datos

```
Telegram -> api/webhook.go | cmd/main.go
         -> pkg/telegram/handler (autoriza, decide, guarda estado en ChatStore)
         -> pkg/transaction/services (parsea, resuelve categoría, registra)
         -> pkg/notion/services (crea la página)  |  pkg/groq/services (foto)
```

## Convenciones

- Dependencias: `handler` y `services` dependen de `pkg/ports` y de
  `model`, nunca entre sí ni de otra feature. `model` solo importa la
  stdlib y otros `model`. Solo `pkg/app`, `api/` y `cmd/` conocen las
  implementaciones concretas.
- Cada `services/` tiene `services.go` (struct + `NewServices`) y un
  archivo por operación.
- Nombre de paquete = nombre de carpeta.
- Secretos solo vía variables de entorno (ver `.env.example`), nunca
  hardcodeados. Ver `.claude/rules/seguridad.md`.
- Comentarios: solo si explican un WHY no obvio. Ver
  `.claude/rules/codigo.md`.

## No hagas

- No crear ramas, no `git commit`, no `push`, no PR's. Git es del
  usuario; el trabajo termina en dejar el working tree listo y un
  resumen claro de qué cambió.
- Ningún dato sensible hardcodeado (token de Telegram, API key de
  Notion, database id) — siempre variable de entorno.
- No agregar dependencias externas (`go get`) sin avisar antes y
  explicar por qué (hoy solo `godotenv`).
- No desplegar a Vercel ni correr `cmd/setwebhook` contra el bot real
  sin que el usuario lo pida: toca su Telegram y su Notion reales.

## Flujo de trabajo

1. Historias de usuario → `specs/features/NNN-nombre-feature/`
   (`spec.md`, `plan.md`, `tasks.md`) antes de tocar código. Crear la
   spec no implica empezar a implementar.
2. No implementar hasta que el usuario lo indique explícitamente.
3. Cambios no triviales (tocan el contrato entre capas, agregan una
   dependencia, cambian la entidad `Transaction`): proponer un plan
   corto y esperar OK antes de tocar código.
4. Cambios chicos y acotados a un solo paquete: proceder directo y
   reportar al final.
5. Gate antes de dar algo por completado: `go build ./...`, `go vet
   ./...` y `go test ./...` (cuando haya tests) deben pasar.
6. Si no hay 80% de seguridad sobre una decisión de negocio o de
   arquitectura, preguntar — no asumir alcance no pedido.

## Documentación

- `CLAUDE.md` (raíz) — bitácora corta: resumen rápido + índice de
  `.claude/history/` + última entrada completa. Se lee automático al
  abrir sesión en esta carpeta.
- `.claude/rules/` — reglas operativas modulares (comunicación, flujo de
  trabajo, estilo de código, seguridad, bitácora, más las específicas de
  Telegram/Notion), pensadas para cambiar con el tiempo.
- `.claude/history/` — una entrada de bitácora por archivo, nombrado
  `YYYY-MM-DD-feature-que_se_realizo.md` (ver
  `.claude/rules/bitacora.md`).
- `.claude/plans/` — planes de trabajo del proyecto.
- `specs/constitution/` — misión, roadmap y tech-stack del producto
  (referencia fija; `mission.md` pendiente de llenar, necesita input de
  negocio del usuario).
- `specs/features/NNN-nombre-feature/` — una carpeta por historia de
  usuario, con `spec.md` / `plan.md` / `tasks.md`.
