package chat

import (
	"time"
)

type Message struct {
	Timestamp   time.Time
	Message     string
	IsDelivered bool
}
