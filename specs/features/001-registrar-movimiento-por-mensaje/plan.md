# 001 · Registrar movimiento por mensaje de Telegram — Plan

_Cómo se implementa lo descrito en `spec.md`. Debe respetar la `constitution/`._

## Enfoque

Parser posicional simple (no NLP libre) en `internal/usecase/`, contra
tablas de palabras clave fijas para `Tipo` y `Payment Method`, y contra
un `CategoryResolver` (interfaz definida en usecase, implementada en
`internal/infrastructure/notion/`) que consulta la data source
`Category` real en vez de hardcodear la lista de categorías — si el
usuario agrega una categoría nueva en Notion, el bot la reconoce sin
tocar código.

## Implementación

1. `internal/usecase/parse_transaction.go`: función `ParseMessage(text
   string) (domain.Transaction, error)`.
   - Tokeniza por espacios.
   - Token 0 → `Tipo` (case-insensitive `egreso`/`ingreso`). Si no
     matchea, error.
   - Busca el primer token que parsea como número → `Monto`. Todo lo
     que hay entre `Tipo` y `Monto` → `Concepto` (puede ser vacío).
   - De los tokens después de `Monto`: intenta matchear, desde el
     final, un nombre de categoría conocido (1-2 palabras) — si
     matchea, lo separa como `Categoria` y lo saca de la lista.
   - Lo que queda después de sacar `Categoria` → `MedioDePago`,
     matcheado contra la tabla de `spec.md`. Si no matchea ningún medio
     de pago conocido, error.
   - Si no hay `Concepto`, usa `Tipo` capitalizado como `Name`.
2. `internal/usecase/category_resolver.go` (o similar): interfaz
   `CategoryResolver` con `Resolve(name string) (categoryID string, ok
   bool, err error)`. Usecase depende de la interfaz, no de Notion.
3. `internal/infrastructure/notion/category_resolver.go`: implementa
   `CategoryResolver` consultando la data source `Category`
   (`e1e6ba98-2ca3-825b-9c56-872d36f73df0`) vía la REST API — trae
   nombre → ID, cacheable en memoria con TTL corto (evita 1 request de
   Notion por mensaje si el usuario manda varios seguidos).
4. `internal/usecase` orquesta: `ParseMessage` → resolver categoría (si
   vino) → `notion.Client.CreateExpense` → responder al usuario por
   Telegram (éxito o error claro).
5. `internal/infrastructure/telegram/`: handler de mensaje entrante
   invoca el usecase, no el parser directamente (usecase es el punto de
   entrada de negocio, ver `AGENT.md`).

## Decisiones

- **Parser posicional, no NLP** — el formato del usuario ya es
  estructurado (`Tipo Concepto Monto MedioDePago Categoria`); un parser
  de reglas es más predecible y fácil de debuggear que NLP libre para
  este caso de uso personal.
- **Categorías resueltas contra Notion en vivo, no hardcodeadas** — se
  descartó una lista fija en código porque el usuario ya las administra
  en Notion (`Category` data source) y espera que cambien ahí sin tocar
  el bot.
- **Categoría/medio de pago desconocidos → error, no se crea la
  página** — se descartó "crear categoría nueva silenciosamente" o
  "ignorar el campo": los datos financieros del usuario no deben quedar
  mal clasificados sin que se dé cuenta.
- **Sin fecha en el mensaje → fecha de recepción** — más simple, y
  ningún ejemplo del usuario incluye fecha.

## Riesgos

- **Ambigüedad si `Concepto` o `Categoria` contienen un número o una
  palabra clave de medio de pago** (ej. concepto "banco" literal) — el
  parser posicional puede confundirlo. Mitigación: documentado como
  limitación conocida, no bloqueante para v1; se revisa si aparece en
  uso real.
- **Categoría con espacio (`Seguridad Social`) vs. una sola palabra** —
  el matching debe probar 2 palabras antes que 1 al final del mensaje,
  para no cortar `Seguridad Social` en `Social` suelto.
- **Cache de categorías desincronizada** si el usuario crea una
  categoría nueva en Notion y manda un mensaje antes de que expire el
  TTL — mitigación: TTL corto (ej. 5 min), no crítico para un bot
  personal de uso esporádico.
