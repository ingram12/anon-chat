package users

import (
	"sync"
	"time"
)

type UserStorage struct {
	Mu                    sync.RWMutex
	users                 map[string]User
	userInactivityTimeout time.Duration
}

func NewUserStorage(userInactivityTimeout time.Duration) *UserStorage {
	return &UserStorage{
		users:                 make(map[string]User),
		userInactivityTimeout: userInactivityTimeout,
	}
}

func (s *UserStorage) CreateUser(user User) User {
	s.users[user.ID] = user
	return user
}

func (s *UserStorage) GetUser(userID string) (User, bool) {
	user, exists := s.users[userID]
	return user, exists
}

func (s *UserStorage) GetUserLocked(userID string) (User, bool) {
	s.Mu.RLock()
	defer s.Mu.RUnlock()
	user, exists := s.users[userID]
	return user, exists
}

func (s *UserStorage) IsUserExist(userID string) bool {
	_, exists := s.GetUser(userID)
	return exists
}

func (s *UserStorage) UpdateUser(user User) {
	_, exist := s.users[user.ID]
	if !exist {
		return
	}
	s.users[user.ID] = user
}

func (s *UserStorage) UpdateLastActivityLocked(userID string) bool {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	user, exists := s.users[userID]
	if !exists {
		return false
	}
	user.LastActivity = time.Now()
	s.users[userID] = user
	return true
}

func (s *UserStorage) RemoveInactiveUsers() {
	timeNow := time.Now()
	for id, user := range s.users {
		if timeNow.Sub(user.LastActivity) > s.userInactivityTimeout {
			delete(s.users, id)
		}
	}
}
