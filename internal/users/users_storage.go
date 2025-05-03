package users

import (
	"log"
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

	idBytes := StringToBytes(userID)

	user, exists := s.users[idBytes]
	return user, exists
}

func (s *UserStorage) IsUserExist(userID string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	idBytes := StringToBytes(userID)

	_, exists := s.users[idBytes]
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

func (s *UserStorage) DeleteInactiveUsers() {
	s.mu.Lock()
	defer s.mu.Unlock()

	timeNow := time.Now()
	for id, user := range s.users {
		if timeNow.Sub(user.LastActivity) > s.userInactivityTimeout {
			delete(s.users, id)
		}
	}
}

func (s *UserStorage) GetUsersWithoutChat() []User {
	s.mu.RLock()
	defer s.mu.RUnlock()

	availableUsers := make([]User, 0, 2)
	for _, user := range s.users {
		if user.ChatID == 0 && user.IsRegistered {
			log.Printf("User %s is available, chatId %d\n", user.GetUserID(), user.ChatID)
			availableUsers = append(availableUsers, user)
		}
	}
	return availableUsers
}
