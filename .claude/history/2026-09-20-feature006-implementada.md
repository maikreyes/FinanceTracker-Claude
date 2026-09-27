# 2026-09-20 — Feature 006 implementada: adjuntar foto de recibo

## Contexto

Con las 5 features del MVP (001-005) implementadas y 007 ya
especificada (depende de esta), el usuario pidió empezar por 006.

## Qué se implementó

- `internal/infrastructure/telegram/client.go`: `Message.Photo
  []PhotoSize`, `Message.Caption string`, `PhotoSize{FileID, Width,
  Height}`, `Bot.GetFile` (resuelve `file_path` interno de Telegram),
  `Bot.DownloadFile` (descarga el contenido real).
- `internal/domain/transaction.go`: `ReceiptURLs` renombrado a
  `ReceiptFileIDs` — son IDs de `file_upload` de Notion, no URLs.
- `internal/infrastructure/notion/receipt_uploader.go`:
  `Client.Upload` implementa el File Upload API de Notion en dos pasos
  (crear `file_upload` → enviar el contenido real vía
  `multipart/form-data` al `upload_url`). Implementa
  `usecase.ReceiptUploader`.
- `internal/infrastructure/notion/client.go`:
  `mapTransactionToProperties` mapea `Receipt` como `{"files":
  [{"type": "file_upload", "file_upload": {"id": ...}}]}`.
- `internal/usecase/ports.go`: interfaz `ReceiptUploader`.
- `internal/usecase/register_transaction.go`: `Registrar` gana un
  campo `uploader` (constructor `NewRegistrar` ahora recibe 3
  argumentos, no 2) y el método `RegisterWithPhoto(text string, photo
  []byte, filename, contentType string)`. Lógica de resolución de
  categoría extraída a un helper privado `resolveCategory`, compartido
  entre `Register` y `RegisterWithPhoto` (antes estaba duplicada).
- `cmd/bot/main.go`: `handlePhotoMessage` — foto con caption usa el
  caption como texto normal (mismo parser de 001) y adjunta la foto;
  foto sin caption responde pidiendo que se reenvíe con caption
  (`noCaptionReply`) sin registrar nada. El dispatcher del loop de
  polling se reescribió como `switch` (foto / `resumen` / texto
  normal / ninguno) en vez de una cadena de `if`s.
- Tests: `register_transaction_test.go` con `fakeReceiptUploader` —
  éxito con categoría y subida, falla de subida, caption inválido (no
  intenta subir la foto si el caption ya es inválido).
- `.claude/rules/notion-rules.md`: nueva sección documentando el
  contrato del File Upload API.

## Validación real (con el bot corriendo)

Una sola instancia confirmada por `ps aux` antes de empezar. Dos
mensajes reales:

1. Foto **sin** caption → el bot respondió pidiendo que se reenvíe con
   el texto del movimiento (sin log de error, es un camino esperado,
   no una falla).
2. La misma foto **con** caption válido → registrada. Confirmado vía
   query directa a Notion: la página nueva (`Name: consignación,
   Amount: 120, Type: Ingreso`) tiene `Receipt` con 1 archivo
   (`recibo.jpg`, `type: "file"`) — un archivo nativo de Notion, no un
   link externo que fuera a expirar.

Apagado limpio confirmado (`bot detenido`), sin instancias corriendo
al cerrar.

## Estado al cierre de sesión

- `go build ./...`, `go vet ./...`, `go test ./...` limpios
  (`gofmt` corrido sobre los archivos tocados).
- Feature 006 marcada **implementado ✅**, movida a "Hecho" en
  `specs/constitution/roadmap.md`.
- Nota dejada en `tasks.md` de 006 para 007: el límite de tamaño del
  File Upload API en modo `single_part` no se confirmó explícitamente
  (no hizo falta con la foto de prueba) — revisar si 007 sube fotos
  más pesadas.
- Siguiente paso natural: **007** (interpretar recibo con IA,
  Gemini) — ahora que 006 (su dependencia) está implementada.
