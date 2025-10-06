package pg

import (
	"Noty/internal/models"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

func scanNotes(rows *sql.Rows) ([]models.Note, error) {
	notesList := make([]models.Note, 0)
	var err error
	for rows.Next() {
		var note models.Note

		err = rows.Scan(
			&note.ID,
			&note.Title,
			&note.Content,
			&note.CreatedAt,
			&note.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		notesList = append(notesList, note)
	}
	return notesList, rows.Err()
}

func (s *Storage) CreateNewNote(ctx context.Context, userID int, title, content string) (int, error) {
	const op = "storage.pg.notes.CreateNewNote"

	query := `
		INSERT INTO notes (user_id, title, content) VALUES ($1, $2, $3) RETURNING id;
	`
	var noteID int
	err := s.db.QueryRowContext(ctx, query, userID, title, content).Scan(&noteID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, fmt.Errorf("%s: %w", op, ErrNotFound)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return noteID, nil
}

func (s *Storage) GetAllNotes(ctx context.Context, userID int, meta models.NoteMeta) (*[]models.Note, error) {
	const op = "storage.pg.notes.GetAllNotes"

	dir := meta.Sort
	query := fmt.Sprintf(`
		SELECT id, title, content, created_at, updated_at FROM notes 
		WHERE user_id = $1 ORDER BY created_at %s
		LIMIT $2 OFFSET $3;
	`, dir)
	rows, err := s.db.QueryContext(ctx, query, userID, meta.Limit, meta.Offset)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	result, err := scanNotes(rows)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, ErrNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &result, nil
}

func (s *Storage) GetNote(ctx context.Context, noteID, userID int) (models.Note, error) {
	const op = "storage.pg.notes.GetNote"

	query := `
		SELECT id, user_id, title, content, created_at, updated_at
		FROM notes WHERE id = $1 AND user_id = $2;
	`
	var note models.Note
	err := s.db.QueryRowContext(ctx, query, noteID, userID).Scan(
		&note.ID,
		&note.UserID,
		&note.Title,
		&note.Content,
		&note.CreatedAt,
		&note.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		if err := s.Validate(ctx, userID, noteID); errors.Is(err, ErrForbidden) {
			return models.Note{}, fmt.Errorf("%s: %w", op, ErrForbidden)
		}
		return models.Note{}, fmt.Errorf("%s: %w", op, ErrNotFound)
	}
	if err != nil {
		return models.Note{}, fmt.Errorf("%s: %w", op, err)
	}

	return note, nil
}

// UpdateNote requests update note by noteID.
func (s *Storage) UpdateNote(ctx context.Context, noteID, userID int, title, content string) (models.Note, error) {
	const op = "storage.pg.notes.UpdateNote"

	query := `
		UPDATE notes
		SET title = $1,
			content = $2
		WHERE id = $3 AND user_id = $4
		RETURNING id, title, content, created_at, updated_at;
	`
	var note models.Note
	err := s.db.QueryRowContext(ctx, query, title, content, noteID, userID).Scan(
		&note.ID,
		&note.Title,
		&note.Content,
		&note.CreatedAt,
		&note.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		if err := s.Validate(ctx, userID, noteID); errors.Is(err, ErrForbidden) {
			return models.Note{}, fmt.Errorf("%s: %w", op, ErrForbidden)
		}
		return models.Note{}, fmt.Errorf("%s: %w", op, ErrNotFound)
	}
	if err != nil {
		return models.Note{}, fmt.Errorf("%s: %w", op, err)
	}

	return note, nil
}

func (s *Storage) DeleteNote(ctx context.Context, noteID int, userID int) error {
	const op = "storage.pg.notes.DeleteNote"

	query := `
		DELETE FROM notes WHERE id = $1 AND user_id = $2 RETURNING id;
	`
	var returnedID int
	err := s.db.QueryRowContext(ctx, query, noteID, userID).Scan(&returnedID)
	if errors.Is(err, sql.ErrNoRows) {
		if err := s.Validate(ctx, userID, noteID); errors.Is(err, ErrForbidden) {
			return fmt.Errorf("%s: %w", op, ErrForbidden)
		}
		return fmt.Errorf("%s: %w", op, ErrNotFound)
	}
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
