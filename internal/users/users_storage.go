package users

import (
	"sync"
	"time"
)

type UserStorage struct {
	mu                    sync.RWMutex
	users                 map[[36]byte]User
	userInactivityTimeout time.Duration
}

func NewUserStorage(userInactivityTimeout time.Duration) *UserStorage {
	return &UserStorage{
		users:                 make(map[[36]byte]User),
		userInactivityTimeout: userInactivityTimeout,
	}
}

func (s *UserStorage) CreateUser(user User) User {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.users[user.ID] = user
	return user
}

func (s *UserStorage) GetUser(userID string) (User, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	user, exists := s.users[StringToBytes(userID)]
	return user, exists
}

func (s *UserStorage) GetUserBytes(userID [36]byte) (User, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	user, exists := s.users[userID]
	return user, exists
}

func (s *UserStorage) IsUserExist(userID string) bool {
	_, exists := s.GetUser(userID)
	return exists
}

func (s *UserStorage) UpdateUser(user User) {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, exist := s.users[user.ID]
	if !exist {
		return
	}

	s.users[user.ID] = user
}

func (s *UserStorage) RemoveInactiveUsers() {
	s.mu.Lock()
	defer s.mu.Unlock()

	timeNow := time.Now()
	for id, user := range s.users {
		if timeNow.Sub(user.LastActivity) > s.userInactivityTimeout {
			delete(s.users, id)
		}
	}
}
