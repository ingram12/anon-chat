package chat

import (
	"testing"
	"time"
)

func TestChat(t *testing.T) {
	const user1 = "user1"
	const user2 = "user2"
	const user3 = "user3"

	now := time.Now()
	chat := &Chat{
		ID:        1,
		CreatedAt: now,
		UserID1:   user1,
		UserID2:   user2,
		User1Messages: []Message{
			{Timestamp: now, Message: "hello from user1", IsDelivered: true},
		},
		User2Messages: []Message{
			{Timestamp: now, Message: "hello from user2", IsDelivered: true},
		},
	}

	t.Run("IsUserInChat", func(t *testing.T) {
		if !chat.IsUserInChat(user1) {
			t.Error("user1 should be in chat")
		}
		if !chat.IsUserInChat(user2) {
			t.Error("user2 should be in chat")
		}
		if chat.IsUserInChat(user3) {
			t.Error("user3 should not be in chat")
		}
	})

	t.Run("IsActive", func(t *testing.T) {
		if !chat.IsActive() {
			t.Error("chat should be active with both users")
		}

		chat.UserID1 = ""
		if chat.IsActive() {
			t.Error("chat should not be active with only one user")
		}

		chat.UserID2 = ""
		if chat.IsActive() {
			t.Error("chat should not be active with no users")
		}
	})

	t.Run("GetPeerID", func(t *testing.T) {
		chat.UserID1 = user1
		chat.UserID2 = user2

		if peerID := chat.GetPeerID(user1); peerID != user2 {
			t.Errorf("Expected peer of user1 to be user2, got %s", peerID)
		}
		if peerID := chat.GetPeerID(user2); peerID != user1 {
			t.Errorf("Expected peer of user2 to be user1, got %s", peerID)
		}
	})

	t.Run("GetPeerMessages", func(t *testing.T) {
		messages := chat.GetPeerMessages(user1)
		if len(messages) != 1 || messages[0].Message != "hello from user2" {
			t.Error("user1 should see user2's messages")
		}

		messages = chat.GetPeerMessages(user2)
		if len(messages) != 1 || messages[0].Message != "hello from user1" {
			t.Error("user2 should see user1's messages")
		}
	})
}
