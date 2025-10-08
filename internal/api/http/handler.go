package handlers

import (
	"Noty/internal/service"
	"Noty/pkg/auth"
	"log/slog"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Handler struct {
	UserService  *service.UserService
	NoteService  *service.NoteService
	Logger       *slog.Logger
	TokenManager *auth.Manager
}

func NewHandler(userS *service.UserService, noteS *service.NoteService,
	logger *slog.Logger, TM *auth.Manager) *Handler {
	return &Handler{
		UserService:  userS,
		NoteService:  noteS,
		Logger:       logger,
		TokenManager: TM,
	}
}

func (h *Handler) GetRouter() chi.Router {
	r := chi.NewRouter()

	r.Use(h.MWLog())

	r.Get("/docs/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	r.Route("/api/v1", func(r chi.Router) {

		r.Post("/sign-in", h.PostSignUp)
		r.Post("/login", h.PostLogin)
		// r.Get("/refresh", h.GetRefresh)

		r.With(h.jwtAuthMiddleware).Group(func(r chi.Router) {
			r.Post("/notes", h.PostNewNote)
			r.With(h.GetMeta).Get("/notes", h.GetNotes)
			r.Get("/notes/{note_id}", h.GetNoteByID)
			r.Put("/notes/{note_id}", h.PutUpdateNote)
			r.Delete("/notes/{note_id}", h.DeleteNote)
		})
	})

	return r
}
