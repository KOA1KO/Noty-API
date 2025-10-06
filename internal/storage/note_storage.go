package storage

import (
	"Noty/internal/models"
	"context"
)

type NoteStorage interface {
	CreateNewNote(ctx context.Context, userID int, title, content string) (int, error)
	GetAllNotes(ctx context.Context, userID int, meta models.NoteMeta) (*[]models.Note, error)
	GetNote(ctx context.Context, noteID, userID int) (models.Note, error)
	UpdateNote(ctx context.Context, noteID, userID int, title, content string) (models.Note, error)
	DeleteNote(ctx context.Context, noteID int, userID int) error
}
