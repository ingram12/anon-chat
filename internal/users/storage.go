package users

import (
	"sync"
	"time"
)

type User struct {
	ID               [36]byte  `json:"user_id"`
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
	users map[[36]byte]*User
}

func NewUserStorage() *UserStorage {
	return &UserStorage{
		users: make(map[[36]byte]*User),
	}
}

func (s *UserStorage) CreateUser(userID, challenge string, difficulty int) (*User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var idBytes [36]byte
	copy(idBytes[:], userID)

	user := &User{
		ID:               idBytes,
		CurrentChallenge: challenge,
		Difficulty:       difficulty,
		CreatedAt:        time.Now(),
		IsRegistered:     false,
	}

	s.users[idBytes] = user
	return user, nil
}

func (s *UserStorage) GetUser(userID string) (*User, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var idBytes [36]byte
	copy(idBytes[:], userID)

	user, exists := s.users[idBytes]
	return user, exists
}

func (s *UserStorage) IsUserExist(userID string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var idBytes [36]byte
	copy(idBytes[:], userID)

	_, exists := s.users[idBytes]
	return exists
}

func (s *UserStorage) UpdateUser(user *User) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.users[user.ID] = user
}
