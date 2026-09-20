# 004 · Confirmación y manejo de errores en el chat — Plan

_Cómo se implementa lo descrito en `spec.md`. Debe respetar la `constitution/`._

## Enfoque

El usecase de 001 devuelve errores tipados (sentinel errors o tipos de
error concretos), sin saber nada de Telegram. La capa de transporte
(002) traduce esos errores a texto para el chat — así el mapeo de
copy vive en un solo lugar y el usecase se puede testear sin mockear
Telegram.

## Implementación

1. `internal/usecase/errors.go`: errores tipados —
   `ErrUnrecognizedType`, `ErrUnrecognizedPaymentMethod{Got string}`,
   `ErrUnrecognizedCategory{Got string}`, `ErrWriteFailed{Err error}`
   (wrapea el error real de Notion, se loguea completo pero no se
   manda al chat).
2. `internal/usecase.CategoryResolver` (definida en la feature 001):
   agregar `ListNames() ([]string, error)` — para poder listar
   categorías válidas en el mensaje de error, no solo resolver una.
3. `internal/usecase/register_transaction.go` (orquestador de 001):
   devuelve estos errores en vez de errores genéricos.
4. `internal/infrastructure/telegram/` (o `cmd/bot/main.go`): función
   `replyFor(result RegisterResult, err error) string` — switch con
   `errors.As` sobre los tipos de 1, arma el texto exacto por caso.
5. Texto de éxito: usa el `pageID` devuelto por `CreateExpense`,
   arma el link como `https://notion.so/<pageID sin guiones>`.

## Decisiones

- **Errores tipados en el usecase, copy en el transporte** — mantiene
  la regla `usecase` no importa nada de `infrastructure/telegram`
  (`AGENT.md`). El usecase solo describe *qué* pasó, no *cómo se lo
  decimos al usuario por Telegram*.
- **Categorías del mensaje de error se listan en vivo** — mismo
  principio que la resolución de categoría en 001: nunca hardcodear la
  lista en el código, para que agregar una categoría en Notion no
  requiera tocar strings de error.
- **Sin exponer errores internos de Notion al chat** — el error real
  (`ErrWriteFailed.Err`) se loguea a stdout/stderr del proceso, el chat
  solo ve un mensaje genérico. Evita filtrar detalles de
  infraestructura a quien esté viendo el chat.
- **Sin emojis por default** — no fue pedido explícitamente; se puede
  agregar después si el usuario lo pide.

## Riesgos

- **Mensaje de error de categoría muy largo** si hay muchas categorías
  en Notion — aceptable mientras la lista sea corta (hoy son 6); si
  crece mucho, se puede truncar o sugerir la más parecida (fuera de
  alcance por ahora).
