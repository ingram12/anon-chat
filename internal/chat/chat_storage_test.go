package chat

import (
	"testing"
)

func TestChatStorage(t *testing.T) {
	storage := NewChatStorage()

	const user1 = "user1"
	const user2 = "user2"
	const user3 = "user3"

	t.Run("CreateChat", func(t *testing.T) {
		chat := storage.CreateChat(user1, user2)
		if chat.ID != 1 {
			t.Errorf("Expected first chat ID to be 1, got %d", chat.ID)
		}
		if chat.UserID1 != user1 || chat.UserID2 != user2 {
			t.Error("Chat created with wrong user IDs")
		}
		if len(chat.User1Messages) != 0 || len(chat.User2Messages) != 0 {
			t.Error("New chat should have no messages")
		}
	})

	t.Run("GetChat", func(t *testing.T) {
		chat, err := storage.GetChat(1)
		if err != nil {
			t.Fatalf("Failed to get chat: %v", err)
		}
		if chat.ID != 1 {
			t.Errorf("Got wrong chat ID: %d", chat.ID)
		}

		_, err = storage.GetChat(999)
		if err != ErrChatNotFound {
			t.Error("Expected ErrChatNotFound for non-existent chat")
		}
	})

	t.Run("AddMessage", func(t *testing.T) {
		timestamp, err := storage.AddMessage(1, user1, "hello")
		if err != nil {
			t.Fatalf("Failed to add message: %v", err)
		}
		if timestamp.IsZero() {
			t.Error("Expected non-zero timestamp")
		}

		chat, _ := storage.GetChat(1)
		if len(chat.User1Messages) != 1 {
			t.Error("Message not added to user1's messages")
		}
		if chat.User1Messages[0].Message != "hello" {
			t.Error("Wrong message content stored")
		}

		// Test adding message for non-existent chat
		_, err = storage.AddMessage(999, user1, "hello")
		if err != ErrChatNotFound {
			t.Error("Expected ErrChatNotFound")
		}

		// Test adding message for non-participant
		_, err = storage.AddMessage(1, user3, "hello")
		if err == nil {
			t.Error("Expected error when adding message from non-participant")
		}
	})

	t.Run("GetPeerMessages", func(t *testing.T) {
		messages, err := storage.GetPeerMessages(1, user2)
		if err != nil {
			t.Fatalf("Failed to get peer messages: %v", err)
		}
		if len(messages) != 1 {
			t.Error("Expected one message")
		}
		if messages[0].Message != "hello" {
			t.Error("Wrong message content retrieved")
		}
	})

	t.Run("RemovePeerMessages", func(t *testing.T) {
		err := storage.RemovePeerMessages(1, user2)
		if err != nil {
			t.Fatalf("Failed to remove peer messages: %v", err)
		}

		messages, _ := storage.GetPeerMessages(1, user2)
		if len(messages) != 0 {
			t.Error("Messages should have been removed")
		}

		err = storage.RemovePeerMessages(999, user1)
		if err != ErrChatNotFound {
			t.Error("Expected ErrChatNotFound")
		}
	})

	t.Run("QuitChat", func(t *testing.T) {
		err := storage.QuitChat(1, user1)
		if err != nil {
			t.Fatalf("Failed to quit chat: %v", err)
		}

		chat, _ := storage.GetChat(1)
		if chat.UserID1 != "" {
			t.Error("user1 should have been removed from chat")
		}
		if chat.UserID2 != user2 {
			t.Error("user2 should still be in chat")
		}

		if storage.IsActiveChat(1) {
			t.Error("Chat should not be active after user quit")
		}

		err = storage.QuitChat(999, user1)
		if err != ErrChatNotFound {
			t.Error("Expected ErrChatNotFound")
		}
	})

	t.Run("IsUserInChat", func(t *testing.T) {
		if storage.IsUserInChat(1, user1) {
			t.Error("user1 should not be in chat after quitting")
		}
		if !storage.IsUserInChat(1, user2) {
			t.Error("user2 should still be in chat")
		}
		if storage.IsUserInChat(1, user3) {
			t.Error("user3 was never in chat")
		}
	})

	t.Run("IsActiveChat", func(t *testing.T) {
		if storage.IsActiveChat(1) {
			t.Error("Chat should not be active with only one user")
		}

		err := storage.QuitChat(1, user2)
		if err != nil {
			t.Fatalf("Failed to quit chat: %v", err)
		}

		if storage.IsActiveChat(1) {
			t.Error("Chat should not be active with no users")
		}
	})
}
