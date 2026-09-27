# 006 · Adjuntar foto de recibo al registro — Tareas

_Checklist accionable derivada del `plan.md`. Tareas pequeñas y concretas; marca `[x]` al completarlas._

- [x] `internal/infrastructure/telegram/client.go`: `Message.Photo`,
      `Message.Caption`, `PhotoSize`, `Bot.GetFile`, `Bot.DownloadFile`.
- [x] `internal/usecase/ports.go`: interfaz `ReceiptUploader`.
- [x] `internal/infrastructure/notion/receipt_uploader.go`:
      implementación real contra el File Upload API de Notion (crear →
      `upload_url` → enviar multipart → `status: uploaded`).
- [x] `internal/domain/transaction.go`: `ReceiptURLs` →
      `ReceiptFileIDs`.
- [x] `internal/infrastructure/notion/client.go`:
      `mapTransactionToProperties` mapea `Receipt` con `file_upload`.
- [x] `internal/usecase/register_transaction.go`:
      `Registrar.RegisterWithPhoto` (+ `resolveCategory` extraído como
      helper compartido con `Register`, para no duplicar la lógica).
- [x] `cmd/bot/main.go`: `handlePhotoMessage` — detecta `Message.Photo`,
      usa `Caption`, descarga + sube + registra; foto sin caption →
      `noCaptionReply` sin llamar al usecase.
- [x] Tests unitarios de `RegisterWithPhoto` con fakes (`fakeReceiptUploader`)
      — éxito con categoría, falla de subida, caption inválido (no
      intenta subir la foto si el caption ya es inválido).
- [x] Prueba manual real: foto sin caption → el bot pidió reenviar;
      foto con caption válido → registrada en Notion con `Receipt`
      conteniendo el archivo real (`recibo.jpg`, tipo `file`, no un
      link externo).
- [x] `go build ./...`, `go vet ./...`, `go test ./...` en verde.
- [x] Validar contra los criterios de aceptación de `spec.md`.
- [x] Mover la feature a "Hecho" en `../../constitution/roadmap.md`.

## Nota para 007

El "límite de tamaño del File Upload API de Notion" (riesgo anotado en
`plan.md`) no se validó explícitamente — no hizo falta con la foto de
prueba real. Si 007 sube fotos más pesadas (celulares modernos),
revisar el límite real al implementar.
