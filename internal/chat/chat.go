package chat

import (
	"time"
)

type Chat struct {
	ID            int
	CreatedAt     time.Time
	UserID1       [36]byte
	UserID2       [36]byte
	User1Messages []Message
	User2Messages []Message
}

func (c *Chat) IsUserInChat(userID [36]byte) bool {
	return c.UserID1 == userID || c.UserID2 == userID
}

func (c *Chat) GetPeerID(userID [36]byte) [36]byte {
	if c.UserID1 == userID {
		return c.UserID2
	}
	return c.UserID1
}

func (c *Chat) GetPeerMessages(userID [36]byte) []Message {
	if c.UserID1 == userID {
		return c.User2Messages
	}
	return c.User1Messages
}
