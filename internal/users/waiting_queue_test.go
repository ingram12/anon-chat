package users

import (
	"anon-chat/internal/chat"
	"fmt"
	"sync"
	"testing"
	"time"
)

func f(format string, a ...interface{}) string {
	return fmt.Sprintf(format, a...)
}

func TestWaitingQueue(t *testing.T) {
	t.Run("Basic Queue Operations", func(t *testing.T) {
		wq := NewWaitingQueue()

		if wq.GetLen() != 0 {
			t.Error("New queue should be empty")
		}

		wq.AddUser("user1")
		if wq.GetLen() != 1 {
			t.Error("Queue should have one user")
		}

		wq.RemoveUser("user1")
		if wq.GetLen() != 0 {
			t.Error("Queue should be empty after removal")
		}
	})

	t.Run("Concurrent Access", func(t *testing.T) {
		wq := NewWaitingQueue()
		var wg sync.WaitGroup
		numGoroutines := 10

		wg.Add(numGoroutines)
		// Concurrent additions
		for i := 0; i < numGoroutines; i++ {
			go func(id int) {
				defer wg.Done()
				wq.AddUserLocked(f("user%d", id))
			}(i)
		}
		wg.Wait()

		wg.Add(numGoroutines)
		// Concurrent removals
		for i := 0; i < numGoroutines; i++ {
			go func(id int) {
				defer wg.Done()
				wq.RemoveUserLocked(f("user%d", id))
			}(i)
		}
		wg.Wait()

		if wq.GetLen() != 0 {
			t.Error("Queue should be empty after all operations")
		}
	})

	t.Run("GetTwoRandomUsers", func(t *testing.T) {
		wq := NewWaitingQueue()

		// Test with empty queue
		_, _, err := wq.GetTwoRandomUsers()
		if err == nil {
			t.Error("Expected error when getting users from empty queue")
		}

		// Test with one user
		wq.AddUser("user1")
		_, _, err = wq.GetTwoRandomUsers()
		if err == nil {
			t.Error("Expected error when getting users from queue with single user")
		}

		// Test with two users
		wq.AddUser("user2")
		user1, user2, err := wq.GetTwoRandomUsers()
		if err != nil {
			t.Error("Unexpected error when getting two users")
		}
		if user1 == "" || user2 == "" {
			t.Error("Both users should be non-empty")
		}
		if user1 == user2 {
			t.Error("Should return different users")
		}
	})

	t.Run("TryMatch", func(t *testing.T) {
		wq := NewWaitingQueue()
		userStorage := NewUserStorage(5 * time.Minute)
		chatStorage := chat.NewChatStorage()

		user1 := User{
			ID:           "user1",
			IsRegistered: true,
		}
		user2 := User{
			ID:           "user2",
			IsRegistered: true,
		}

		userStorage.CreateUser(user1)
		userStorage.CreateUser(user2)
		wq.AddUser(user1.ID)
		wq.AddUser(user2.ID)

		// Test concurrent TryMatch calls
		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			defer wg.Done()
			wq.TryMatch(chatStorage, userStorage)
		}()
		go func() {
			defer wg.Done()
			wq.TryMatch(chatStorage, userStorage)
		}()

		wg.Wait()
		// Give the matching goroutine time to complete
		time.Sleep(10 * time.Millisecond)

		// Verify only one matching operation occurred
		updatedUser1, _ := userStorage.GetUser(user1.ID)
		updatedUser2, _ := userStorage.GetUser(user2.ID)
		if updatedUser1.ChatID == 0 || updatedUser2.ChatID == 0 {
			t.Error("Users should be matched")
		}
		if updatedUser1.ChatID != updatedUser2.ChatID {
			t.Error("Users should be in the same chat")
		}
	})
}
