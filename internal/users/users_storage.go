package users

import (
	"sync"
	"time"
)

type UserStorage struct {
	Mu                    sync.RWMutex
	Users                 map[string]User
	userInactivityTimeout time.Duration
}

func NewUserStorage(userInactivityTimeout time.Duration) *UserStorage {
	return &UserStorage{
		Users:                 make(map[string]User),
		userInactivityTimeout: userInactivityTimeout,
	}
}

func (s *UserStorage) CreateUser(user User) User {
	s.Users[user.ID] = user
	return user
}

func (s *UserStorage) GetUser(userID string) (User, bool) {
	user, exists := s.Users[userID]
	return user, exists
}

func (s *UserStorage) GetUserLocked(userID string) (User, bool) {
	s.Mu.RLock()
	defer s.Mu.RUnlock()
	user, exists := s.Users[userID]
	return user, exists
}

func (s *UserStorage) IsUserExist(userID string) bool {
	_, exists := s.GetUser(userID)
	return exists
}

func (s *UserStorage) UpdateUser(user User) {
	_, exist := s.Users[user.ID]
	if !exist {
		return
	}
	s.Users[user.ID] = user
}

func (s *UserStorage) UpdateLastActivityLocked(userID string) bool {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	user, exists := s.Users[userID]
	if !exists {
		return false
	}
	user.LastActivity = time.Now()
	s.Users[userID] = user
	return true
}

func (s *UserStorage) RemoveInactiveUsers() {
	timeNow := time.Now()
	for id, user := range s.Users {
		if timeNow.Sub(user.LastActivity) > s.userInactivityTimeout {
			delete(s.Users, id)
		}
	}
}
