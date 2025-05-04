package chat

import (
	"time"
)

type Chat struct {
	ID            int
	CreatedAt     time.Time
	UserID1       string
	UserID2       string
	User1Messages []Message
	User2Messages []Message
}

func (c *Chat) IsUserInChat(userID string) bool {
	return c.UserID1 == userID || c.UserID2 == userID
}

func (c *Chat) GetPeerID(userID string) string {
	if c.UserID1 == userID {
		return c.UserID2
	}
	return c.UserID1
}

func (c *Chat) GetPeerMessages(userID string) []Message {
	if c.UserID1 == userID {
		return c.User2Messages
	}
	return c.User1Messages
}
