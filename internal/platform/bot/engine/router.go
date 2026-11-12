package engine

import (
	"context"
	"log/slog"
	"strings"
	"technovilochka/internal/platform/bot/state"
	"technovilochka/internal/platform/bot/transport"
)

type Router struct {
	commands         map[transport.Command]Handler
	callbacks        map[transport.Callback]Handler
	messageFallback  Handler
	callbackFallback Handler
	states           map[state.State]Handler
	sm               *state.Machine
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

func (r *Router) RegisterMessageFallback(handler Handler) {
	r.messageFallback = handler
}

func (r *Router) RegisterCallbackFallback(handler Handler) {
	r.callbackFallback = handler
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
	userID := update.Message.UserID
	text := update.Message.Text
	userState := r.sm.GetStep(userID)

	if h, ok := r.commands[transport.Command(text)]; ok {
		r.sm.SetStep(userID, state.Empty)
		h(ctx, update)
		return
	}

	if userState != "" {
		if h, ok := r.states[userState]; ok {
			h(ctx, update)
			return
		}
		r.sm.SetStep(userID, state.Empty)
	}

	if r.messageFallback != nil {
		r.messageFallback(ctx, update)
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

	if r.callbackFallback != nil {
		r.callbackFallback(ctx, update)
		return
	}

	slog.Warn("Received unexpected callback", "callback", update.CallbackQuery)
}
