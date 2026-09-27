package services

import (
	"context"
	"net/url"
)

// SetWebhook apunta el bot a webhookURL. Telegram manda secretToken en el
// header X-Telegram-Bot-Api-Secret-Token de cada update, que es lo que
// permite al endpoint distinguir a Telegram de cualquiera que adivine la
// URL. Mientras haya un webhook activo, getUpdates falla.
func (s *Services) SetWebhook(ctx context.Context, webhookURL, secretToken string) error {
	params := url.Values{}
	params.Set("url", webhookURL)
	params.Set("secret_token", secretToken)
	params.Set("allowed_updates", `["message","callback_query"]`)
	return s.call(ctx, "setWebhook", params, nil)
}

// DeleteWebhook desactiva el webhook para poder volver a long-polling.
// dropPending descarta los updates que Telegram tenía encolados.
func (s *Services) DeleteWebhook(ctx context.Context, dropPending bool) error {
	params := url.Values{}
	if dropPending {
		params.Set("drop_pending_updates", "true")
	}
	return s.call(ctx, "deleteWebhook", params, nil)
}
