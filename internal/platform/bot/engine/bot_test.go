package engine

import (
	"context"
	"technovilochka/internal/platform/bot/state"
	"technovilochka/internal/platform/bot/transport"
	"testing"
	"time"
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

func TestRunStopsWhenContextCancelled(t *testing.T) {
	updates := make(chan transport.Update) // канал остаётся открытым, как при обычной работе
	bot := NewBot(&fakeClient{updates: updates}, NewRouter(state.New()))

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		defer close(done)
		bot.Run(ctx)
	}()

	cancel()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("Run did not stop after context cancellation")
	}
}
