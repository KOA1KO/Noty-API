package models

import "time"

type Session struct {
	ExpiresAt        time.Time `json:"exp"`
	IssuedAt         time.Time `json:"iss"`
	RefreshTokenHash string    `json:"refresh"`
}
