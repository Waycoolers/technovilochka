package engine

import (
	"context"
	"log/slog"
	"strings"
	"technovilochka/internal/bot/state"
	"technovilochka/internal/bot/transport"
)

type Router struct {
	commands  map[transport.Command]Handler
	callbacks map[transport.Callback]Handler
	fallback  Handler
	states    map[state.State]Handler
	sm        *state.Machine
}

type Handler func(ctx context.Context, update *transport.Update)

func NewRouter(sm *state.Machine) *Router {
	return &Router{
		sm:        sm,
		commands:  make(map[transport.Command]Handler),
		callbacks: make(map[transport.Callback]Handler),
		states:    make(map[state.State]Handler),
	}
}

func (r *Router) RegisterCommand(cmd transport.Command, handler Handler) {
	r.commands[cmd] = handler
}

func (r *Router) RegisterCallback(callback transport.Callback, handler Handler) {
	r.callbacks[callback] = handler
}

func (r *Router) RegisterFallback(handler Handler) {
	r.fallback = handler
}

func (r *Router) RegisterState(s state.State, handler Handler) {
	r.states[s] = handler
}

func (r *Router) handleUpdate(ctx context.Context, update *transport.Update) {
	if update.Message != nil {
		r.handleMessage(ctx, update)
		return
	}
	if update.CallbackQuery != nil {
		r.handleCallback(ctx, update)
		return
	}

	slog.Warn("Received unexpected update", "update", update)
}

func (r *Router) handleMessage(ctx context.Context, update *transport.Update) {
	chatID := update.Message.ChatID
	text := update.Message.Text
	userState := r.sm.GetStep(chatID)

	if h, ok := r.commands[transport.Command(text)]; ok {
		r.sm.SetStep(chatID, state.Empty)
		h(ctx, update)
		return
	}

	if userState != "" {
		if h, ok := r.states[userState]; ok {
			h(ctx, update)
			return
		}
		r.sm.SetStep(chatID, state.Empty)
	}

	if r.fallback != nil {
		r.fallback(ctx, update)
		return
	}

	slog.Warn("Received unexpected message", "message", update.Message)
}

func (r *Router) handleCallback(ctx context.Context, update *transport.Update) {
	data := update.CallbackQuery.Data
	prefix := strings.SplitN(data, ":", 2)[0]

	if h, ok := r.callbacks[transport.Callback(prefix)]; ok {
		h(ctx, update)
		return
	}

	if r.fallback != nil {
		r.fallback(ctx, update)
		return
	}

	slog.Warn("Received unexpected callback", "callback", update.CallbackQuery)
}
