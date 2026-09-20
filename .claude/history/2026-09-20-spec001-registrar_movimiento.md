# 2026-09-20 — Spec 001: Registrar movimiento por mensaje de Telegram

## Contexto

Con el schema de Notion ya soportando `Type` (Ingreso/Egreso, entrada
anterior de bitácora), el usuario confirmó dos cosas más:

1. Las categorías también van en el mensaje, al final: `egreso mecato
   13000 efectivo comida`.
2. El mapeo de medio de pago en español debe traducirse al nombre
   exacto de la opción en Notion (confirmó la propuesta que hice:
   `efectivo`→Cash, `tarjeta credito`→Credit Card, `tarjeta
   debito`→Debit Card, `tranferencia`/`banco`→Bank).

## Qué se hizo

- Se consultó la data source `Category` completa vía MCP
  (`notion-query-data-sources`) para tener la lista real y no
  inventarla: `Entretenimiento`, `Subcripciones`, `Seguridad Social`,
  `Comida`, `Ahorro`, `Servicios`.
- Creada `specs/features/001-registrar-movimiento-por-mensaje/` con
  `spec.md`, `plan.md`, `tasks.md`:
  - Gramática del mensaje: `<Tipo> [Concepto] <Monto> <MedioDePago>
    [Categoria]`, posicional.
  - Tabla de mapeo medio de pago → opción de Notion, documentada en
    `spec.md`.
  - Decisión: categorías se resuelven **contra Notion en vivo**
    (interfaz `CategoryResolver`), no hardcodeadas en el parser — si el
    usuario agrega una categoría nueva en Notion, el bot la reconoce
    sin redeploy.
  - Decisión: `MedioDePago`/`Categoria` no reconocidos → error claro,
    no se crea la página (no se inventa ni se ignora el dato).
  - Decisión: sin `Concepto` (como en los ingresos) → `Name` = `Tipo`
    capitalizado.
  - Riesgo anotado: `Seguridad Social` es una categoría de 2 palabras —
    el matching debe probar 2 palabras antes que 1 al final del
    mensaje.
- `specs/constitution/roadmap.md` actualizado: 001 pasa de "por
  definir" a spec/plan/tasks listos, en "Siguiente".

## Estado al cierre de sesión

Spec, plan y tasks de la feature 001 completos y documentados. **Sin
implementar todavía** — según el flujo de trabajo pactado
(`.claude/rules/flujo_de_trabajo.md`), la creación de la spec no
implica empezar a codear; falta que el usuario dé el OK explícito para
implementar `ParseMessage`, `CategoryResolver` y la orquestación en
usecase.
