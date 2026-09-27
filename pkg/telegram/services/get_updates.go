package services

import (
	"context"
	"net/url"
	"strconv"

	"finance-tracker/pkg/telegram/model/update"
)

// GetUpdates hace long-polling nativo de Telegram (timeout=30s del lado
// del servidor, no sondeo agresivo del cliente).
func (s *Services) GetUpdates(ctx context.Context, offset int) ([]update.Update, error) {
	params := url.Values{}
	params.Set("offset", strconv.Itoa(offset))
	params.Set("timeout", "30")

	var updates []update.Update
	if err := s.call(ctx, "getUpdates", params, &updates); err != nil {
		return nil, err
	}
	return updates, nil
}
