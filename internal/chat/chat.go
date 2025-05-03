package chat

import (
	"time"
)

type Chat struct {
	ID        int
	CreatedAt time.Time
	UserID1   [36]byte
	UserID2   [36]byte
	Messages  map[int]Message
}
