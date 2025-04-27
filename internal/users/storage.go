package users

import (
	"sync"
	"time"
)

type User struct {
	Id               string    `json:"user_id"`
	Nickname         string    `json:"nickname,omitempty"`
	Tags             []string  `json:"tags,omitempty"`
	PublicKey        string    `json:"public_key,omitempty"`
	CurrentChallenge string    `json:"current_challenge,omitempty"`
	Difficulty       int       `json:"difficulty,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
	IsRegistered     bool      `json:"is_registered"`
}

type UserStorage struct {
	mu    sync.RWMutex
	users map[string]*User
}

func NewUserStorage() *UserStorage {
	return &UserStorage{
		users: make(map[string]*User),
	}
}

func (s *UserStorage) CreateUser(userId string, challenge string, difficulty int) (*User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	user := &User{
		Id:               userId,
		CurrentChallenge: challenge,
		Difficulty:       difficulty,
		CreatedAt:        time.Now(),
		IsRegistered:     false,
	}

	s.users[userId] = user
	return user, nil
}

func (s *UserStorage) GetUser(userID string) (*User, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, exists := s.users[userID]
	return user, exists
}

func (s *UserStorage) UpdateUser(user *User) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.users[user.Id] = user
}
