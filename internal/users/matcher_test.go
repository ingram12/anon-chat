package users

import (
	"anon-chat/internal/chat"
	"testing"
	"time"
)

func TestMatchUsers(t *testing.T) {
	userStorage := NewUserStorage(5 * time.Minute)
	chatStorage := chat.NewChatStorage()
	waitingQueue := NewWaitingQueue()

	// Create test users
	user1 := User{
		ID:           "user1",
		Nickname:     "User 1",
		PublicKey:    "key1",
		CreatedAt:    time.Now(),
		LastActivity: time.Now(),
		IsRegistered: true,
		ChatID:       0,
	}
	user2 := User{
		ID:           "user2",
		Nickname:     "User 2",
		PublicKey:    "key2",
		CreatedAt:    time.Now(),
		LastActivity: time.Now(),
		IsRegistered: true,
		ChatID:       0,
	}

	t.Run("MatchTwoUsers", func(t *testing.T) {
		userStorage.CreateUser(user1)
		userStorage.CreateUser(user2)
		waitingQueue.AddUser(user1.ID)
		waitingQueue.AddUser(user2.ID)

		MatchUsers(userStorage, chatStorage, waitingQueue)

		// Verify users were matched
		updatedUser1, _ := userStorage.GetUser(user1.ID)
		updatedUser2, _ := userStorage.GetUser(user2.ID)

		if updatedUser1.ChatID == 0 {
			t.Error("User1 should be assigned to a chat")
		}
		if updatedUser2.ChatID == 0 {
			t.Error("User2 should be assigned to a chat")
		}
		if updatedUser1.ChatID != updatedUser2.ChatID {
			t.Error("Users should be assigned to the same chat")
		}

		// Verify users were removed from waiting queue
		if waitingQueue.GetLen() != 0 {
			t.Error("Waiting queue should be empty after matching")
		}
	})

	t.Run("NoMatchWithSingleUser", func(t *testing.T) {
		// Reset storages
		userStorage = NewUserStorage(5 * time.Minute)
		chatStorage = chat.NewChatStorage()
		waitingQueue = NewWaitingQueue()

		userStorage.CreateUser(user1)
		waitingQueue.AddUser(user1.ID)

		MatchUsers(userStorage, chatStorage, waitingQueue)

		// Verify no matching occurred
		updatedUser1, _ := userStorage.GetUser(user1.ID)
		if updatedUser1.ChatID != 0 {
			t.Error("Single user should not be matched")
		}
		if waitingQueue.GetLen() != 1 {
			t.Error("Single user should remain in waiting queue")
		}
	})

	t.Run("SkipUsersAlreadyInChat", func(t *testing.T) {
		// Reset storages
		userStorage = NewUserStorage(5 * time.Minute)
		chatStorage = chat.NewChatStorage()
		waitingQueue = NewWaitingQueue()

		user1.ChatID = 1 // User already in chat
		userStorage.CreateUser(user1)
		userStorage.CreateUser(user2)
		waitingQueue.AddUser(user1.ID)
		waitingQueue.AddUser(user2.ID)

		MatchUsers(userStorage, chatStorage, waitingQueue)

		// Verify user1 was skipped and removed from queue
		if waitingQueue.GetLen() != 1 {
			t.Error("Only the user already in chat should be removed from queue")
		}
		updatedUser2, _ := userStorage.GetUser(user2.ID)
		if updatedUser2.ChatID != 0 {
			t.Error("User2 should not be matched since user1 was invalid")
		}
	})
}