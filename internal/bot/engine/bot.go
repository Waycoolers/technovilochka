package engine

import (
	"context"
	"log/slog"
	"runtime/debug"
	"technovilochka/internal/bot/transport"
)

type Bot struct {
	client transport.BotClient
	router *Router
}

func NewBot(c transport.BotClient, r *Router) *Bot {
	return &Bot{
		client: c,
		router: r,
	}
}

func (b *Bot) Run(ctx context.Context) {
	slog.Info("Bot is running")

	updateChan := b.client.GetUpdatesChan()

	for update := range updateChan {
		if update.Message != nil || update.CallbackQuery != nil {
			b.handleUpdate(ctx, &update)
		}
	}
}

// handleUpdate гасит панику одного хендлера, иначе она останавливает бота для всех пользователей.
func (b *Bot) handleUpdate(ctx context.Context, update *transport.Update) {
	defer func() {
		r := recover()
		if r != nil {
			slog.Error("Panic while handling update", "panic", r, "stack", string(debug.Stack()))
		}
	}()
	slog.Info("Received update", "message", update.Message)
	b.router.handleUpdate(ctx, update)
}

func (b *Bot) Stop() {
	if b.client != nil {
		slog.Info("Bot is stopping")
		b.client.StopReceivingUpdates()
	}
}
