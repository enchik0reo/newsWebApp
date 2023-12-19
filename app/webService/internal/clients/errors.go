package clients

import "errors"

var (
	ErrInvalidValue        = errors.New("invalid credentials")
	ErrInvalidToken        = errors.New("invalid token")
	ErrUserDoesntExists    = errors.New("user doesn't exist")
	ErrUserExists          = errors.New("user already exists")
	ErrTokenExpired        = errors.New("token expired")
	ErrSessionNotFound     = errors.New("session not found")
	ErrNoPublishedArticles = errors.New("there are no published articles")
)
