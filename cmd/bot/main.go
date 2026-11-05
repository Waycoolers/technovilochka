package main

import (
	"log/slog"
	"os"
	"technovilochka/internal/config"
	"technovilochka/pkg/logger"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load(".env")

	jsonHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	slog.SetDefault(slog.New(jsonHandler))

	cfg, err := config.Load()
	if err != nil {
		slog.Error("Error loading config", "error", err)
		os.Exit(1)
	}

	logger.Init(cfg.LogLVL)
}
