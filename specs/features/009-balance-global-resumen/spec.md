# 009 · Balance global en el resumen + backfill de Tipo histórico

**Estado:** implementado ✅ (spec escrita después de implementar — el
usuario dio la instrucción completa y clara en un solo mensaje, se
codeó directo; esta spec documenta lo que se construyó, no reemplaza
el flujo normal para el resto de features).

## Qué hace

Dos cambios pedidos juntos por el usuario:

1. **Balance en `resumen`** — el comando `resumen` (005) ahora
   responde con `Balance = Ingresos - Egresos` del mes ("cuánto tengo
   disponible"), además de los totales de Egreso/Ingreso que ya
   mostraba, y cuántos movimientos son de cada tipo (conteo, no solo
   suma).
2. **Backfill histórico** — todas las páginas de `Expenses` que
   existían antes de que se agregara la propiedad `Type` (2026-09-20,
   ver `.claude/rules/notion-rules.md`) y por lo tanto no tenían
   `Type` seteado, se marcaron como `Egreso`.

## Por qué

El usuario lo pidió explícitamente: quiere ver el balance del mes de
un vistazo, no solo los totales separados; y quiere que las filas
históricas (anteriores al bot) cuenten como gastos en ese balance, no
que queden fuera del cálculo para siempre.

## Decisión de signo del balance

Primer pedido, literal: "que sume los egresos y le reste los ingresos"
→ se implementó `Balance = Egreso - Ingreso`. Al mostrárselo al
usuario, aclaró que lo quería al revés — **`Balance = Ingreso -
Egreso`** ("cuánto tengo disponible", la convención estándar). Se
corrigió en el mismo turno, antes de dar la feature por cerrada.

## Criterios de aceptación

- [x] `resumen` responde con `Balance`, `Egresos` (suma + cantidad de
      movimientos) e `Ingresos` (suma + cantidad de movimientos).
- [x] `Balance = Ingreso - Egreso`, calculado sobre el mismo filtro de
      mes calendario que ya usaba 005 (`America/Bogota`).
- [x] Ninguna página de `Expenses` existente en Notion quedó sin
      `Type` después del backfill.
- [x] El backfill no tocó ninguna otra propiedad de las páginas
      (`Name`, `Amount`, `Date`, `Category`, etc. sin cambios) — solo
      `Type`.
- [x] Validado en vivo con el bot corriendo (con la fórmula
      `Egreso - Ingreso`, después corregida): `resumen` respondió
      `Egresos: 2338155 (30 movimientos), Ingresos: 0 (0
      movimientos)`, confirmado por el usuario contra un cálculo
      independiente. Con la fórmula corregida, mismos totales/conteos,
      `Balance: -2338155`.

## Fuera de alcance

- Reclasificar manualmente alguna de las 31 filas históricas que en
  realidad podría ser un ingreso (ej. "Prestamo Sara", "Prestamo
  Oscarin", "Préstamo Camilo" — nombres ambiguos, podrían ser
  préstamos dados o recibidos) — se marcaron todas como `Egreso` sin
  excepción, tal como se pidió literalmente ("coloca **todos** los
  registros... como egresos").
- Balance por categoría o por rango de fechas custom — sigue siendo
  solo el mes calendario actual, como 005.
