package ctxKey

import (
	"Noty/internal/models"
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
)

// collisions secure
type ctxKey string

const userIDKey ctxKey = "userID"
const metaKey ctxKey = "meta"

var (
	ErrInvalidSubject = errors.New("неверный subject")
	ErrUnauthorized   = errors.New("пользователь не авторизован")
)

// Getter of UserID. Returns UserID and type assertion.
func GetUserID(ctx context.Context) (jwt.MapClaims, bool) {
	user, ok := ctx.Value(userIDKey).(jwt.MapClaims)
	return user, ok
}

// Setter of UserID. Returns context with UserID.
func WithUserID(ctx context.Context, user jwt.MapClaims) context.Context {
	return context.WithValue(ctx, userIDKey, user)
}

// Getter of Note Meta. Returns NoteMeta model and type assertion.
func GetNoteMeta(ctx context.Context) (models.NoteMeta, bool) {
	meta, ok := ctx.Value(metaKey).(models.NoteMeta)
	return meta, ok
}

// Setter of Note Meta params. Returns context with NoteMeta struct.
func WithNoteMeta(ctx context.Context, meta models.NoteMeta) context.Context {
	return context.WithValue(ctx, metaKey, meta)
}

func GetUserIDFromRequest(r *http.Request) (int, error) {
	user, ok := GetUserID(r.Context())
	if !ok {
		return 0, ErrUnauthorized
	}
	sub, err := user.GetSubject()
	if err != nil {
		return 0, ErrInvalidSubject
	}
	userID, err := strconv.Atoi(sub)
	if err != nil {
		return 0, ErrInvalidSubject
	}
	return userID, nil
}
