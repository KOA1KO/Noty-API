package storage

import (
	"Noty/internal/models"
	"context"
)

type UserStorage interface {
	CreateUser(ctx context.Context, username, passwordHash string) (int, error)
	GetUser(username string) (models.User, error)
	SetSession(ctx context.Context, userID int, session models.Session) error
}
