package pg

import (
	"Noty/internal/models"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
)

func (s *Storage) CreateUser(ctx context.Context, username, passwordHash string) (int, error) {
	const op = "storage.pg.users.CreateUser"

	query := `
		INSERT INTO users (username, password_hash) VALUES ($1, $2) RETURNING id;
	`
	var userID int
	err := s.db.QueryRowContext(ctx, query, username, passwordHash).Scan(&userID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return 0, fmt.Errorf("%s: %w", op, ErrAlredyExists)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return userID, nil
}

func (s *Storage) GetUser(username string) (models.User, error) {
	const op = "storage.pg.users.GetUser"

	query := `
		SELECT id, password_hash, user_session FROM users WHERE username = $1
	`
	var (
		User models.User
		raw  []byte
	)
	err := s.db.QueryRow(query, username).Scan(&User.ID, &User.Password, &raw)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, fmt.Errorf("%s: %w", op, ErrNotFound)
		}
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}
	if err := json.Unmarshal(raw, &User.Session); err != nil {
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}
	return User, nil
}

func (s *Storage) SetSession(ctx context.Context, userID int, session models.Session) error {
	const op = "internal.storage.pg.sessions.SetSession"
	query := `
		UPDATE users SET user_session = $1 WHERE id = $2;
	`
	ses, err := json.Marshal(session)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = s.db.ExecContext(ctx, query, ses, userID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
