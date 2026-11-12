package delivery

import (
	"context"
	"technovilochka/internal/platform/bot/responder"
	"technovilochka/internal/platform/bot/transport"
	"technovilochka/internal/user/usecase"
)

type UserHandler struct {
	client transport.BotClient
	uc     *usecase.Service
	r      *responder.Responder
}

func NewUserHandler(client transport.BotClient, uc *usecase.Service, r *responder.Responder) *UserHandler {
	return &UserHandler{
		client: client,
		uc:     uc,
		r:      r,
	}
}

func (h *UserHandler) Start(ctx context.Context, update *transport.Update) {
	userID := update.Message.UserID
	chatID := update.Message.ChatID
	username := update.Message.UserName
	firstName := update.Message.FirstName
	text := update.Message.Text

	resp := "Техновилочка приветствует тебя!"

	err := h.uc.Register(ctx, usecase.RegisterInput{
		ID:        userID,
		ChatID:    chatID,
		Username:  username,
		FirstName: firstName,
	})
	if err != nil {
		h.r.HandleErr(chatID, text, err)
		return
	}
	h.client.SendMessage(chatID, resp)
}
