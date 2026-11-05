package config

import (
	"log/slog"
	"os"
)

type Config struct {
	LogLVL    string
	BotConfig *BotConfig
	DBConfig  *DBConfig
}

func Load() (*Config, error) {
	logLVL := os.Getenv("LOG_LVL")
	if logLVL == "" {
		slog.Warn("LOG_LVL environment variable not set")
		logLVL = "warn"
	}

	bot, err := loadBotConfig()
	if err != nil {
		return nil, err
	}
	db, err := loadDBConfig()
	if err != nil {
		return nil, err
	}

	return &Config{
		LogLVL:    logLVL,
		BotConfig: bot,
		DBConfig:  db,
	}, nil
}
