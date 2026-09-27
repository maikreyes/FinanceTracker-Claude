# 2026-09-20 — Spec 006: adjuntar foto de recibo al registro

## Contexto

Usuario pidió una spec nueva: poder mandar una foto de un recibo y que
quede cargada al registro. `Receipt` ya existe en el schema real de
Notion (documentado desde la spec 001) pero se había dejado fuera de
alcance por no tener forma de adjuntar un archivo desde el formato de
texto.

Se preguntó cómo se manda la foto junto con el movimiento — dos
opciones: foto con caption (un solo mensaje atómico) vs. texto primero
y foto después (requiere estado entre mensajes). El usuario eligió
foto con caption.

## Qué se hizo

Creada `specs/features/006-adjuntar-recibo-foto/` con
`spec.md`/`plan.md`/`tasks.md`:

- Gramática: la foto se manda con el mensaje normal
  (`Tipo Concepto Monto MedioDePago Categoria`) como **caption** de la
  foto, no como texto aparte.
- Decisión técnica clave: los `file_id` de Telegram no son URLs
  permanentes — hay que descargar la foto de Telegram
  (`getFile`/`DownloadFile`) y subirla de verdad al **File Upload API**
  de Notion (`file_uploads` → `upload_url` → `status: uploaded`) antes
  de poder referenciarla en `Receipt`. No se puede guardar el link de
  Telegram directo (expira).
- `domain.Transaction.ReceiptURLs` se renombra a `ReceiptFileIDs` en el
  plan — el nombre viejo era impreciso (son IDs de `file_upload` de
  Notion, no URLs).
- Nuevo método `Registrar.RegisterWithPhoto` (separado de `Register`,
  no lo toca) + interfaz `ReceiptUploader`.
- Foto sin caption → error antes de tocar Notion (no sube el archivo
  si no hay a qué movimiento asociarlo).
- Riesgos anotados: límite de tamaño del File Upload API (sin asumir
  el número, se verifica al implementar), foto subida pero
  `CreateExpense` falla después (file upload huérfano, aceptado para
  v1, sin rollback manual).
- Fuera de alcance: adjuntar a un movimiento ya existente, álbumes de
  varias fotos, documentos/PDF, editar/reemplazar un `Receipt` ya
  puesto.
- `specs/constitution/roadmap.md`: 006 agregada, marcada independiente
  de 003/004/005 (solo depende de 001/002, ya implementadas) — se
  puede hacer en cualquier momento, no tiene que respetar el orden de
  las demás.

## Estado al cierre de sesión

Spec, plan y tasks de 006 completos. **Sin implementar** — según el
flujo pactado, falta OK explícito del usuario.
