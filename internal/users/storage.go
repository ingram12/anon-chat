package users

import (
	"sync"
	"time"
)

type UserStorage struct {
	mu    sync.RWMutex
	users map[[36]byte]*User
}

func NewUserStorage() *UserStorage {
	return &UserStorage{
		users: make(map[[36]byte]*User),
	}
}

func (s *UserStorage) CreateUser(userID, challenge string) (*User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var idBytes [36]byte
	copy(idBytes[:], userID)

	user := &User{
		ID:               idBytes,
		CurrentChallenge: challenge,
		CreatedAt:        time.Now(),
		IsRegistered:     false,
	}

	user.Difficulty = user.CalcDifficalty()

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
