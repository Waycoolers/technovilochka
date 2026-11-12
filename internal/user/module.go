package user

import (
	"technovilochka/internal/platform/bot/engine"
	"technovilochka/internal/platform/bot/responder"
	"technovilochka/internal/platform/bot/transport"
	"technovilochka/internal/user/adapter"
	"technovilochka/internal/user/delivery"
	"technovilochka/internal/user/usecase"

	"github.com/jmoiron/sqlx"
)

func New(router *engine.Router, client transport.BotClient, r *responder.Responder, db *sqlx.DB) {
	repo := adapter.NewRepo(db)

	service := usecase.NewService(repo)

	handler := delivery.NewUserHandler(client, service, r)

	router.RegisterCommand(transport.Start, handler.Start)
}
