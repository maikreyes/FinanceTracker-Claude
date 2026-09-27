# 006 · Adjuntar foto de recibo al registro — Plan

_Cómo se implementa lo descrito en `spec.md`. Debe respetar la `constitution/`._

## Enfoque

Telegram entrega la foto como un `file_id` interno, no una URL pública
permanente — hay que (1) descargarla de Telegram y (2) subirla de
verdad a Notion vía su File Upload API antes de poder referenciarla en
la propiedad `Receipt` de la página. No se puede simplemente guardar un
link de Telegram en `Receipt` (expira).

## Implementación

1. `internal/infrastructure/telegram/client.go`:
   - `Message.Photo []PhotoSize` y `Message.Caption string` (Telegram
     manda el texto de una foto en `caption`, no en `text`).
   - `PhotoSize{FileID string, Width, Height int}` — Telegram manda
     varias resoluciones de la misma foto; se usa la de mayor
     resolución (último elemento del array).
   - `Bot.GetFile(ctx, fileID) (filePath string, err error)` — llama
     `getFile`.
   - `Bot.DownloadFile(ctx, filePath) ([]byte, error)` — `GET
     https://api.telegram.org/file/bot<token>/<filePath>`.
2. `internal/usecase/ports.go`: nueva interfaz `ReceiptUploader`:
   `Upload(data []byte, filename, contentType string) (fileUploadID
   string, err error)`.
3. `internal/infrastructure/notion/receipt_uploader.go`: implementa
   `ReceiptUploader` contra el File Upload API de Notion:
   - `POST /v1/file_uploads` (`mode: single_part`) → `{id, upload_url}`.
   - `POST <upload_url>` multipart con el contenido del archivo → queda
     `status: uploaded`.
   - Devuelve el `id` del file upload.
4. `internal/infrastructure/notion/client.go`:
   `mapTransactionToProperties` — si `tx.ReceiptFileIDs` no está vacío,
   agrega `"Receipt": {"files": [{"type": "file_upload",
   "file_upload": {"id": "<id>"}}, ...]}`.
5. `internal/domain/transaction.go`: renombrar `ReceiptURLs []string`
   → `ReceiptFileIDs []string` (son IDs de `file_upload` de Notion, no
   URLs — el nombre viejo era impreciso, se corrige ahora que la
   feature existe de verdad).
6. `internal/usecase/register_transaction.go`: nuevo método
   `Registrar.RegisterWithPhoto(text string, photo []byte, filename,
   contentType string) (pageID string, err error)` — parsea el texto
   igual que `Register`, sube la foto con `ReceiptUploader`, arma
   `tx.ReceiptFileIDs`, sigue el mismo camino de categoría +
   `CreateExpense`. `Register` (sin foto) queda intacto.
7. `cmd/bot/main.go`: en el loop, si `update.Message.Photo` no está
   vacío → usar `Message.Caption` como texto, descargar la foto
   (mayor resolución), llamar `RegisterWithPhoto`. Si no hay `Caption`
   → responder pidiendo que reenvíe con caption, no llamar al usecase.

## Decisiones

- **Se sube la foto a Notion, no se referencia el link de Telegram** —
  los links de archivo de Telegram expiran; guardar uno roto en
  `Receipt` sería peor que no tener nada.
- **Mayor resolución disponible, sin comprimir** — Telegram ya
  comprime razonablemente sus fotos; no vale la pena procesarla más
  para un caso de uso personal.
- **`RegisterWithPhoto` como método separado de `Register`**, no un
  parámetro opcional — mantiene `Register` simple para 001/002/005 que
  no necesitan tocar nada de esto.
- **Caption vacío → error antes de tocar Notion** — evita subir una
  foto sin poder asociarla a ningún movimiento real.

## Riesgos

- **Tamaño del archivo** — fotos de celular pueden pesar varios MB; el
  File Upload API de Notion en modo `single_part` tiene un límite (ver
  documentación de Notion al implementar, no asumido aquí). Si una
  foto lo supera, el error debe ser claro ("la foto es muy pesada"),
  no un fallo genérico.
- **Timeout de descarga/subida** — dos requests HTTP grandes en serie
  (Telegram → bot → Notion) antes de responder; usar un timeout más
  generoso que el de `CreateExpense` para estas dos llamadas
  específicas.
- **Falla a medio camino** (foto subida a Notion pero `CreateExpense`
  falla después) — quedaría un `file_upload` huérfano en Notion sin
  página asociada. Aceptado para v1 (Notion limpia file uploads no
  usados automáticamente después de un tiempo, según su documentación);
  no se implementa rollback manual.
