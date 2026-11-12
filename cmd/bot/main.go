package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"technovilochka/internal/platform/bot/client"
	"technovilochka/internal/platform/bot/engine"
	"technovilochka/internal/platform/bot/handlers"
	"technovilochka/internal/platform/bot/responder"
	"technovilochka/internal/platform/bot/state"
	"technovilochka/internal/platform/bot/transport"
	"technovilochka/internal/platform/config"
	"technovilochka/internal/platform/storage"
	"technovilochka/internal/user"
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

	s, err := storage.New(cfg.DBConfig)
	if err != nil {
		slog.Error("Error connecting to database", "error", err)
		os.Exit(1)
	}

	err = s.Migrate()
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

	r := responder.New(tgClient)

	bot := engine.NewBot(tgClient, router)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	systemHandlers := handlers.New(tgClient)
	router.RegisterCommand(transport.Help, systemHandlers.Help)
	router.RegisterMessageFallback(r.ReplyUnknownMessage)
	router.RegisterCallbackFallback(r.ReplyUnknownCallback)

	user.New(router, tgClient, r, s.DB())

	bot.Run(ctx)

	bot.Stop()
	s.Close()
}
