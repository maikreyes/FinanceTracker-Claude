package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"finance-tracker/pkg/telegram/model/command"
)

// SetMyCommands registra el menú de comandos que Telegram muestra al
// escribir "/" en el chat. Es solo descubribilidad: el bot igual
// reconoce los comandos por texto aunque esto falle.
func (s *Services) SetMyCommands(ctx context.Context, commands []command.BotCommand) error {
	list := make([]map[string]string, 0, len(commands))
	for _, c := range commands {
		list = append(list, map[string]string{"command": c.Command, "description": c.Description})
	}
	body, err := json.Marshal(list)
	if err != nil {
		return fmt.Errorf("telegram: marshal commands: %w", err)
	}

	params := url.Values{}
	params.Set("commands", string(body))
	return s.call(ctx, "setMyCommands", params, nil)
}
