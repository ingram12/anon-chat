package users

import (
	"time"
)

type User struct {
	ID               [36]byte
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

func (u *User) CalcDifficalty() int {
	if !u.IsRegistered {
		return 300
	}

	return 100
}

func (u *User) GetUserID() string {
	return string(u.ID[:])
}

func (u *User) SetUserID(userID string) {
	var idBytes [36]byte
	copy(idBytes[:], userID)
	u.ID = idBytes
}
