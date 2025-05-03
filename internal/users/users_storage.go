package users

import (
	"anon-chat/internal/chat"
	"log"
	"math/rand"
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
	availableUsers := make([]User, 0, 2)
	for _, user := range s.users {
		if user.ChatID == 0 && user.IsRegistered {
			log.Printf("User %s is available, chatId %d\n", user.GetUserID(), user.ChatID)
			availableUsers = append(availableUsers, user)
		}
	}
	return availableUsers
}

func (s *UserStorage) MatchUsersIntoChats(chatStorage *chat.Storage) {
	ticker := time.NewTicker(300 * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		s.mu.Lock()

		users := s.GetUsersWithoutChat()

		if len(users) < 2 {
			s.mu.Unlock()
			continue
		}

		// Shuffle users randomly
		rand.Shuffle(len(users), func(i, j int) {
			users[i], users[j] = users[j], users[i]
		})

		timeNow := time.Now()
		// Create pairs and assign chats
		for i := 0; i < len(users)-1; i += 2 {
			user1 := users[i]
			user2 := users[i+1]

			chat, err := chatStorage.CreateChat(user1.ID, user2.ID)
			if err != nil {
				continue
			}

			user1.ChatID = chat.ID
			user1.LastActivity = timeNow
			s.users[user1.ID] = user1

			user2.ChatID = chat.ID
			user2.LastActivity = timeNow
			s.users[user2.ID] = user2

			log.Printf("Matched users %s and %s into chat %d\n", user1.GetUserID(), user2.GetUserID(), chat.ID)
		}

		s.mu.Unlock()
	}
}
