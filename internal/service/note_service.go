package service

import (
	"Noty/internal/config"
	"Noty/internal/models"
	"Noty/internal/storage"
	"Noty/internal/storage/pg"
	"Noty/pkg/auth"
	"Noty/pkg/logger/sl"
	"context"
	"errors"
	"log/slog"
)

type NoteService struct {
	Config         *config.Config
	Logger         *slog.Logger
	NoteRepository storage.NoteStorage
	TokenManager   auth.TokenManager
}

func NewNoteService(NoteRep storage.NoteStorage, TM auth.TokenManager,
	logger *slog.Logger, cfg *config.Config) *NoteService {
	return &NoteService{
		Config:         cfg,
		Logger:         logger,
		NoteRepository: NoteRep,
		TokenManager:   TM,
	}
}

// CreateNewNote returns noteID(in response) and error.
func (s *NoteService) CreateNewNote(ctx context.Context, input NoteInput, userID int) (models.ResponseNoteID, error) {
	const op = "internal.service.note_service.CreateNewNote"
	log := s.Logger.With(
		sl.Operation(op),
	)

	noteID, err := s.NoteRepository.CreateNewNote(ctx, userID, input.Title, input.Content)
	if err != nil {
		if errors.Is(err, pg.ErrNotFound) {
			log.Debug("note not found")
			return models.ResponseNoteID{}, ErrNoteNotFound
		}
		log.Error(err.Error())
		return models.ResponseNoteID{}, ErrInternal
	}
	return models.ResponseNoteID{
		NewNoteID: noteID,
	}, nil
}

// Func returns notes of user and error.
func (s *NoteService) GetNotes(ctx context.Context, userID int, meta models.NoteMeta) (*[]models.Note, error) {
	const op = "internal.service.note_service.GetNotes"
	log := s.Logger.With(
		sl.Operation(op),
	)
	notes, err := s.NoteRepository.GetAllNotes(ctx, userID, meta)

	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return notes, nil
}

// Return note and error
func (s *NoteService) GetNote(ctx context.Context, userID, noteID int) (models.NoteResponse, error) {
	const op = "internal.service.note_service.GetNote"
	log := s.Logger.With(
		sl.Operation(op),
	)

	note, err := s.NoteRepository.GetNote(ctx, noteID, userID)
	if err != nil {
		if errors.Is(err, pg.ErrNotFound) {
			log.Debug("note not found")
			return models.NoteResponse{}, ErrNoteNotFound
		} else if errors.Is(err, pg.ErrForbidden) {
			log.Debug("forbidden")
			return models.NoteResponse{}, ErrForbidden
		}
		log.Error(err.Error())
		return models.NoteResponse{}, err
	}

	return models.NoteResponse{
		ID:        note.ID,
		Title:     note.Title,
		Content:   note.Content,
		CreatedAt: note.CreatedAt,
		UpdatedAt: note.UpdatedAt,
	}, nil
}

func (s *NoteService) UpdateNote(ctx context.Context, userID, noteID int, title, content string) (models.NoteResponse, error) {
	const op = "internal.service.note_service.UpdateNote"
	log := s.Logger.With(
		sl.Operation(op),
	)

	note, err := s.NoteRepository.UpdateNote(ctx, noteID, userID, title, content)
	if err != nil {
		if errors.Is(err, pg.ErrNotFound) {
			log.Debug("invalid note id")
			return models.NoteResponse{}, ErrNoteNotFound
		} else if errors.Is(err, pg.ErrForbidden) {
			log.Debug("forbidden")
			return models.NoteResponse{}, ErrForbidden
		}
		log.Error(err.Error())
		return models.NoteResponse{}, err
	}

	return models.NoteResponse{
		ID:        note.ID,
		Title:     note.Title,
		Content:   note.Content,
		CreatedAt: note.CreatedAt,
		UpdatedAt: note.UpdatedAt,
	}, err
}

func (s *NoteService) DeleteNote(ctx context.Context, userID, noteID int) error {
	const op = "internal.service.note_service.DeleteNote"
	log := s.Logger.With(
		sl.Operation(op),
	)

	err := s.NoteRepository.DeleteNote(ctx, noteID, userID)
	if err != nil {
		if errors.Is(err, pg.ErrNotFound) {
			log.Debug("invalid note id")
			return ErrNoteNotFound
		} else if errors.Is(err, pg.ErrForbidden) {
			log.Debug("forbidden")
			return ErrForbidden
		}
		log.Error(err.Error())
		return err
	}
	return nil
}
