package client

import (
	"bytes"
	"errors"
	"log/slog"
	"net/url"
	"strings"
	"testing"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const testToken = "123456:SECRET-token"

func errWithToken() error {
	return &url.Error{
		Op:  "Post",
		URL: "https://api.telegram.org/bot" + testToken + "/sendMessage",
		Err: errors.New("i/o timeout"),
	}
}

func TestConvertUpdateWithoutOptionalFields(t *testing.T) {
	msg := convertUpdate(tgbotapi.Update{
		Message: &tgbotapi.Message{Chat: &tgbotapi.Chat{ID: 1}, Text: "hi"},
	})
	if msg.Message != nil || msg.CallbackQuery != nil {
		t.Fatalf("message without From must be skipped, got %+v", msg)
	}

	cb := convertUpdate(tgbotapi.Update{
		CallbackQuery: &tgbotapi.CallbackQuery{Data: "x"},
	})
	if cb.Message != nil || cb.CallbackQuery != nil {
		t.Fatalf("callback without Message must be skipped, got %+v", cb)
	}
}

func TestRedactTokenFromError(t *testing.T) {
	err := redactToken(errWithToken(), testToken)
	if strings.Contains(err.Error(), testToken) {
		t.Fatalf("token leaked into error: %v", err)
	}
}

func TestTokenRedactingLogger(t *testing.T) {
	var buf bytes.Buffer
	prev := slog.Default()
	slog.SetDefault(slog.New(slog.NewTextHandler(&buf, nil)))
	t.Cleanup(func() { slog.SetDefault(prev) })

	tokenRedactingLogger{token: testToken}.Println(errWithToken())

	if out := buf.String(); out == "" || strings.Contains(out, testToken) {
		t.Fatalf("expected redacted log line, got %q", out)
	}
}
