package engine

import (
	"context"
	"testing"

	"technovilochka/internal/bot/state"
	"technovilochka/internal/bot/transport"
)

// fakeClient реализует только GetUpdatesChan; остальные методы интерфейса в тесте не вызываются.
type fakeClient struct {
	transport.BotClient
	updates chan transport.Update
}

func (f *fakeClient) GetUpdatesChan() <-chan transport.Update { return f.updates }

func TestRunSurvivesHandlerPanic(t *testing.T) {
	router := NewRouter(state.New())
	handled := 0
	router.RegisterCommand(transport.Start, func(context.Context, *transport.Update) {
		handled++
		panic("boom")
	})

	updates := make(chan transport.Update, 2)
	for range 2 {
		updates <- transport.Update{Message: &transport.Message{ChatID: 1, Text: string(transport.Start)}}
	}
	close(updates)

	NewBot(&fakeClient{updates: updates}, router).Run(context.Background())

	if handled != 2 {
		t.Fatalf("expected both updates to be handled despite the panic, got %d", handled)
	}
}
