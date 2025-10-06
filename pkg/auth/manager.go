package auth

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenManager interface {
	NewJWT(userID string, ttl time.Duration) (string, error)
	Parse(tokenStr string, dest jwt.Claims) error
	NewRefreshToken() (string, error)
}

type Manager struct {
	signingKey string
}

func NewManager(signingKey string) (*Manager, error) {
	if signingKey == "" {
		return nil, fmt.Errorf("empty signing key")
	}
	return &Manager{signingKey: signingKey}, nil
}

func (m *Manager) NewJWT(userID string, ttl time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	})

	return token.SignedString([]byte(m.signingKey))
}

func (m *Manager) Parse(tokenStr string, claims jwt.Claims) error {
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (any, error) {
		return []byte(m.signingKey), nil
	})
	if err != nil {
		return ErrInvalidToken
	}

	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return ErrWrongMethod
	}

	if !token.Valid {
		return ErrInvalidToken
	}

	return nil
}

func (m *Manager) NewRefreshToken() (string, error) {
	b := make([]byte, 32)

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	if _, err := r.Read(b); err != nil {
		return "", fmt.Errorf("%w", err)
	}

	return fmt.Sprintf("%x", b), nil
}
