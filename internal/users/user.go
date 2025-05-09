package users

import (
	"time"
)

type User struct {
	ID               string
	Nickname         string
	Tags             []string
	PublicKey        string
	CurrentChallenge string
	Difficulty       int
	CreatedAt        time.Time
	LastActivity     time.Time
	IsRegistered     bool
	ChatID           int // ID of the chat in which the user is currently participating
}
