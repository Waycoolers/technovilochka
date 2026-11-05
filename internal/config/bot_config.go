package config

import (
	"errors"
	"log/slog"
	"os"
	"strconv"
)

type BotConfig struct {
	Token          string
	UpdatesTimeout int
}

func loadBotConfig() (*BotConfig, error) {
	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		return nil, errors.New("BOT_TOKEN environment variable not set")
	}

	updatesTimeout := os.Getenv("UPDATES_TIMEOUT")
	if updatesTimeout == "" {
		slog.Warn("UPDATES_TIMEOUT environment variable not set")
		updatesTimeout = "60"
	}
	intUpdatesTimeout, err := strconv.Atoi(updatesTimeout)
	if err != nil {
		return nil, err
	}

	return &BotConfig{
		Token:          token,
		UpdatesTimeout: intUpdatesTimeout,
	}, nil
}
