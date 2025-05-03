package chat

import (
	"time"
)

type Message struct {
	ID          int
	CreatedAt   time.Time
	UserID      [36]byte
	Message     string
	IsDelivered bool
}
