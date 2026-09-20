# 2026-09-20 — Schema real de Notion: agregada propiedad Type (Ingreso/Egreso)

## Contexto

Pedí al usuario ejemplos reales de mensajes de Telegram para diseñar
el parser. Ejemplos recibidos:

```
Egreso Gasolina 5000 tarjeta credito
Egreso Gasolina 5000 tarjeta debito
Egreso mecato 13000 efectivo
egreso mecato 13000 tranferencia
ingreso 3000000 banco
ingreso 12000 efectivo
```

Dos de los seis son **ingresos**. El schema real de `Expenses`
(validado la sesión anterior) no tenía ningún campo de tipo — era solo
gastos, sin dónde representar un ingreso. Contradicción real entre lo
que el usuario espera registrar y lo que el schema permite.

## Decisión

Se presentaron 3 opciones al usuario: (a) agregar propiedad `Type` al
schema real, (b) soportar solo egresos por ahora, (c) usar el signo de
`Amount` sin tocar el schema. Eligió (a).

## Qué se hizo

1. `notion-update-data-source` sobre la data source `Expenses`
   (`3e36ba98-2ca3-8396-9ede-07441cc01bf8`): `ADD COLUMN "Type"
   SELECT('Egreso':red, 'Ingreso':green)`. Cambio real sobre el Notion
   del usuario, no un mock.
2. `domain.Transaction` (`internal/domain/transaction.go`): nuevo campo
   `Type TransactionType`, tipo con constantes `TransactionTypeExpense
   = "Egreso"` / `TransactionTypeIncome = "Ingreso"` (coinciden
   exactamente con los nombres de opción de Notion, igual que
   `PaymentMethod`). `Amount` se documenta como siempre positivo — el
   signo lo da `Type`, no el número.
3. `internal/infrastructure/notion/client.go`:
   `mapTransactionToProperties` mapea `Type` como `select` (mismo
   patrón que `Payment Method`).
4. Revalidado con el cliente Go real: dos páginas de prueba (una
   `Egreso`, una `Ingreso`), creadas y confirmadas con `CreateExpense`,
   luego archivadas (`PATCH .../pages/{id}` `{"archived": true}`).
   `cmd/notion_smoketest/` temporal, borrado después.

## Qué NO se hizo (deliberado)

- No se hizo backfill de `Type` en las filas históricas (existentes
  antes de este cambio) — quedan con `Type` vacío. No fue pedido y
  toca datos financieros reales existentes del usuario; tocarlos sin
  pedido explícito sería alcance no solicitado.
- No se creó todavía el parser ni la spec `001` — este cambio fue solo
  de schema/dominio/cliente, a raíz del hallazgo. El parser sigue
  pendiente de las demás decisiones de producto (categoría, mapeo de
  "tranferencia" a `Payment Method`, etc.).

## Estado al cierre de sesión

`go build ./...` y `go vet ./...` limpios. Schema real de Notion y
`domain.Transaction`/cliente alineados de nuevo. Pendiente: diseñar y
documentar la gramática del mensaje de Telegram (`spec.md` de la
feature 001) a partir de los 6 ejemplos dados.
