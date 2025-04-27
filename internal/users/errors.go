package users

import "errors"

var (
	ErrInvalidToken     = errors.New("invalid token")
	ErrInvalidSolution  = errors.New("invalid solution")
	ErrUserNotFound     = errors.New("user not found")
	ErrInvalidChallenge = errors.New("invalid challenge")
)
