# Finance Tracker: Telegram to Notion

Bot que lee mensajes en lenguaje natural desde Telegram, los parsea en
transacciones financieras y las escribe en una base de datos de Notion.

## Stack

- Lenguaje: Go 1.23+
- Módulo: `finance-tracker` (`go.mod`, raíz del repo)
- Sin dependencias externas todavía (solo stdlib)
- Arquitectura: Clean Architecture (`domain` / `usecase` / `infrastructure`)
- Punto de entrada: `cmd/bot/main.go`

## Comandos

- `go build ./...` — compila todo el módulo
- `go run ./cmd/bot` — corre el bot en local
- `go vet ./...` — chequeo estático
- `go test ./...` — tests (aún no hay ninguno escrito)
- `gofmt -l .` — lista archivos mal formateados

## Estructura del proyecto

- `cmd/bot/` — punto de entrada único (`main.go`). Wiring de
  dependencias (inyecta infraestructura concreta en los usecases) y
  arranque del bot.
- `internal/domain/` — entidad núcleo `Transaction{Name, Amount, Type,
  Date, CategoryID, PaymentMethod, Notes, ReceiptURLs}`, calcada 1:1
  del esquema real de la data source "Expenses" de Notion (ver
  `.claude/rules/notion-rules.md`). `Type` = `Egreso`/`Ingreso`. Sin
  imports fuera de la stdlib.
- `internal/usecase/` — lógica de negocio: parseo de mensajes en lenguaje
  natural a `domain.Transaction`, orquestación del guardado en Notion.
  Depende solo de `domain` y de interfaces (ports), nunca de tipos
  concretos de infraestructura.
- `internal/infrastructure/telegram/` — adaptador de la Telegram Bot API
  (recibe updates, los convierte en llamadas a usecase). Reglas propias
  en `.claude/rules/telegram-rules.md`.
- `internal/infrastructure/notion/` — adaptador de la Notion API (mapea
  `Transaction` a propiedades de página, crea páginas). Reglas propias
  en `.claude/rules/notion-rules.md`.

## Flujo de datos

```
Telegram message -> internal/infrastructure/telegram (recibe update)
                  -> internal/usecase (parsea texto -> domain.Transaction)
                  -> internal/infrastructure/notion (mapea -> crea página)
                  -> Base de datos de Notion
```

## Convenciones

- Dirección de dependencia: `infrastructure -> usecase -> domain`. El
  dominio no depende de nada externo.
- Nombre de paquete = nombre de carpeta.
- Infraestructura implementa interfaces definidas por usecase (ports),
  nunca al revés — usecase no debe importar infraestructura.
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
  explicar por qué.
- No implementar los clientes reales de Telegram/Notion ni el parser de
  mensajes hasta que se indique explícitamente — hoy son stubs
  (`TODO`) a propósito.

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
