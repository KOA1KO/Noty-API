package main

import (
	"Noty/internal/api/app"
	"Noty/internal/config"
	"Noty/internal/storage/pg"
	"Noty/pkg/auth"
	prettyslog "Noty/pkg/logger/pretty_slog"
	"Noty/pkg/logger/sl"
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "Noty/docs"
)

const (
	envLocal = "local"
	envProd  = "prod"
)

// @title        Noty API
// @version      1.0.0
// @description  API для заметок.
// @contact.name @llimd
// @BasePath     /api/v1
// @schemes      http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description "Bearer <access-token>"

func main() {
	cfg := config.MustLoad()
	log := setupLogger(cfg.ENV)

	log.Info("Starting Noty", slog.String("env", cfg.ENV))
	log.Debug("DEBUG enbled")

	storage, err := pg.New(cfg.DSN)
	if err != nil {
		log.Error("DB error")
		panic(err)
	}
	manager, err := auth.NewManager(cfg.SigningKey)
	if err != nil {
		log.Error("manager error")
		panic(err)
	}

	handler := app.New(cfg, log, storage, manager)

	srv := &http.Server{
		Addr:         cfg.HTTP_Server.Address,
		Handler:      handler.GetRouter(),
		ReadTimeout:  cfg.HTTP_Server.Timeout,
		WriteTimeout: cfg.HTTP_Server.Timeout,
		IdleTimeout:  cfg.HTTP_Server.IdleTimeout,
	}

	go func() {
		log.Info("Starting server", slog.String("address", cfg.HTTP_Server.Address))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("Failed to start server")
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	log.Error("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Error("Server forced shutdown", sl.Err(err))
	}

	if err := storage.Close(); err != nil {
		log.Error("db close error", sl.Err(err))
	}
	log.Error("Server stopped gracefully")
}

// create and setup logger by enviroment
func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = setupPrettySlog()
		// log = slog.New(
		// 	slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		// )
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	}
	return log
}

func setupPrettySlog() *slog.Logger {
	opts := prettyslog.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
