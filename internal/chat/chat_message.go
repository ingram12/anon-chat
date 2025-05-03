package chat

import (
	"time"
)

type Message struct {
	CreatedAt   time.Time
	Message     string
	IsDelivered bool
}
