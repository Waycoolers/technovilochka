package responder

import (
	"context"
	"log/slog"
	"technovilochka/internal/bot/transport"
)

type Responder struct {
	client transport.BotClient
}

func New(client transport.BotClient) *Responder {
	return &Responder{client: client}
}

func (r *Responder) HandleErr(chatID int64, msg string, err error) {
	r.client.SendMessage(chatID, "😔 Что-то пошло не так\nЯ уже стараюсь всё исправить")
	slog.Error("Error with message", "chatID", chatID, "message", msg, "error", err)
}

func (r *Responder) ReplyUnknownCallback(_ context.Context, update *transport.Update) {
	chatID := update.CallbackQuery.ChatID
	r.client.SendMessage(chatID, "🤍 Лучше воспользуйся кнопками - так будет проще")
}

func (r *Responder) ReplyUnknownMessage(_ context.Context, update *transport.Update) {
	chatID := update.Message.ChatID
	r.client.SendMessage(chatID, "🤔 Я пока не понимаю это сообщение\nПопробуй выбрать действие с кнопок ниже")
}
