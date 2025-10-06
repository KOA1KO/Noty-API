package app

import (
	handlers "Noty/internal/api/http"
	"Noty/internal/config"
	"Noty/internal/service"
	"Noty/internal/storage/pg"
	"Noty/pkg/auth"
	"log/slog"
)

func New(cfg *config.Config, log *slog.Logger, storage *pg.Storage, manager *auth.Manager) *handlers.Handler {
	userService := service.NewUserService(storage, manager, log, cfg)
	noteService := service.NewNoteService(storage, manager, log, cfg)
	return handlers.NewHandler(userService, noteService, log, manager)
}
