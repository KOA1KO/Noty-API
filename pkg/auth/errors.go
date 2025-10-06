package auth

import "errors"

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrWrongMethod  = errors.New("wrong method")
)
