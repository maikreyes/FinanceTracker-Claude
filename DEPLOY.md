# Despliegue en Vercel

El bot corre como una función serverless (`api/webhook.go`): Telegram le manda un POST por cada update y el estado pendiente de cada chat vive en Upstash Redis, no en memoria.

## 1. Upstash

1. Crea una base de datos Redis en [upstash.com](https://upstash.com).
2. Copia `UPSTASH_REDIS_REST_URL` y `UPSTASH_REDIS_REST_TOKEN` de la pestaña REST API.

## 2. Secreto del webhook

Genera un valor aleatorio y guárdalo como `TELEGRAM_WEBHOOK_SECRET`. Telegram lo manda en cada request y el endpoint rechaza con 403 cualquier request que no lo traiga.

```
openssl rand -hex 32
```

## 3. Vercel

1. Importa el repositorio en Vercel (framework: Other). `vercel.json` ya fija `maxDuration` en 60 s para `api/webhook.go`, que cubre las llamadas a Notion y Groq.
2. En Settings, Environment Variables, carga todas las variables de `.env.example`:
   `TELEGRAM_BOT_TOKEN`, `TELEGRAM_ALLOWED_CHAT_IDS`, `TELEGRAM_WEBHOOK_SECRET`, `NOTION_API_KEY`, `NOTION_DATABASE_ID`, `NOTION_CATEGORY_DATA_SOURCE_ID`, `GROQ_API_KEY`, `UPSTASH_REDIS_REST_URL`, `UPSTASH_REDIS_REST_TOKEN`.
3. Despliega. Las variables nuevas solo aplican a despliegues posteriores.

## 4. Apuntar Telegram al despliegue

Con el mismo `.env` en tu máquina (mismo token y mismo secreto):

```
go run ./cmd/setwebhook -url https://<tu-proyecto>.vercel.app/api/webhook
```

Para comprobarlo:

```
curl "https://api.telegram.org/bot$TELEGRAM_BOT_TOKEN/getWebhookInfo"
```

`last_error_message` debe estar vacío. Escribe al bot y revisa los logs de la función en Vercel si no responde.

## Volver a desarrollo local

Con un webhook activo, `getUpdates` falla, así que primero hay que borrarlo:

```
go run ./cmd/setwebhook -delete
go run ./cmd
```

Para volver a producción, repite el paso 4.

## Cosas a tener en cuenta

- `TELEGRAM_ALLOWED_CHAT_IDS` vacío deja al bot sin responder a nadie.
- Si falta alguna variable, la función responde 500 y deja el motivo en los logs, sin mostrarlo en la respuesta.
- Telegram reintenta un update si tarda en responder. El bot descarta los repetidos por `update_id`, así que un reintento no registra el movimiento dos veces.
- El estado pendiente (flujo guiado, foto esperando Ingreso/Egreso) vence a los 30 minutos.
