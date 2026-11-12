package handlers

import (
	"context"
	"technovilochka/internal/platform/bot/transport"
)

type System struct {
	client transport.BotClient
}

func New(client transport.BotClient) *System {
	return &System{
		client: client,
	}
}

func (s *System) Help(_ context.Context, update *transport.Update) {
	chatID := update.Message.ChatID
	resp := "Список доступных команд:\n"

	for _, command := range transport.Commands {
		resp += string(command) + "\n"
	}

	s.client.SendMessage(chatID, resp)
}
