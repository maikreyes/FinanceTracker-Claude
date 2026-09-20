# 2026-09-19 — Andamiaje Go + Clean Architecture

## Decisiones

- Lenguaje: Go (module `finance-tracker`, go 1.23).
- Patrón: Clean Architecture (domain / usecase / infrastructure).
- Punto de entrada único: `cmd/bot/main.go`.
- Entidad `domain.Transaction`: `Concepto string`, `Monto float64`,
  `Tipo TransactionType`, `Fecha time.Time`.
- Infra separada por proveedor: `internal/infrastructure/telegram/` y
  `internal/infrastructure/notion/`, cada uno con su propio archivo de
  reglas en `.claude/rules/`.
- Usecase no debe importar tipos concretos de infraestructura (regla de
  dependencia hacia adentro).
- Config vía variables de entorno: `TELEGRAM_BOT_TOKEN`,
  `NOTION_API_KEY`, `NOTION_DATABASE_ID`. Documentadas en `.env.example`,
  nunca committeadas (`.env` en `.gitignore`).
- Repo git inicializado localmente (no existía antes de esta sesión).

## Estado al cierre de sesión

- Fases 1-3 completadas: andamiaje `.claude/`, `go.mod`, `.gitignore`,
  `.env.example`, estructura de directorios y stubs de paquete creados.
- Sin lógica funcional: clients de Telegram/Notion y parser son stubs
  (`TODO`). Sin dependencias externas en `go.mod` todavía.
- Pendiente: validación de arquitectura por el usuario antes de
  implementar lógica de negocio o clientes de API.
