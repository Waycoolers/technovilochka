package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"technovilochka/internal/bot/client"
	"technovilochka/internal/bot/engine"
	"technovilochka/internal/bot/handlers"
	"technovilochka/internal/bot/state"
	"technovilochka/internal/bot/transport"
	"technovilochka/internal/config"
	"technovilochka/internal/storage"
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

	db, err := storage.New(cfg.DBConfig)
	if err != nil {
		slog.Error("Error connecting to database", "error", err)
		os.Exit(1)
	}

	err = db.Migrate()
	if err != nil {
		slog.Error("Error migrating database", "error", err)
		os.Exit(1)
	}

	sm := state.New()

	router := engine.NewRouter(sm)

	tgClient, err := client.NewClient(cfg.BotConfig)
	if err != nil {
		slog.Error("Error creating telegram client", "error", err)
		os.Exit(1)
	}

	bot := engine.NewBot(tgClient, router)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	systemHandlers := handlers.New(tgClient)
	router.RegisterCommand(transport.Start, systemHandlers.Start)
	router.RegisterCommand(transport.Help, systemHandlers.Help)

	bot.Run(ctx)

	<-ctx.Done()
	bot.Stop()
	db.Close()
}
