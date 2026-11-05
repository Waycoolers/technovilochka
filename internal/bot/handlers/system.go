package handlers

import (
	"context"
	"technovilochka/internal/bot/transport"
)

type System struct {
	client transport.BotClient
}

func New(client transport.BotClient) *System {
	return &System{
		client: client,
	}
}

func (s *System) Start(_ context.Context, update *transport.Update) {
	chatID := update.Message.ChatID
	text := "Техновилочка приветствует тебя!"

	s.client.SendMessage(chatID, text)
}

func (s *System) Help(_ context.Context, update *transport.Update) {
	chatID := update.Message.ChatID
	text := "Список доступных команд:\n"

	for _, command := range transport.Commands {
		text += string(command) + "\n"
	}

	s.client.SendMessage(chatID, text)
}
