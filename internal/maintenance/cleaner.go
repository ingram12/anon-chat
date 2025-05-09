package maintenance

import (
	"anon-chat/internal/chat"
	"anon-chat/internal/users"
	"log"
	"time"
)

type Cleaner struct {
	UserStorage  *users.UserStorage
	WaitingQueue *users.WaitingQueue
	ChatStorage  *chat.Storage
	TTL          time.Duration
	TickInterval time.Duration
}

func NewCleaner(
	userStorage *users.UserStorage,
	waitingQueue *users.WaitingQueue,
	chatStorage *chat.Storage,
	ttl time.Duration,
	interval time.Duration,
) *Cleaner {
	return &Cleaner{
		UserStorage:  userStorage,
		WaitingQueue: waitingQueue,
		ChatStorage:  chatStorage,
		TTL:          ttl,
		TickInterval: interval,
	}
}

func (c *Cleaner) Start() {
	go func() {
		ticker := time.NewTicker(c.TickInterval)
		defer ticker.Stop()

		for range ticker.C {
			c.cleanup()
		}
	}()
}

func (c *Cleaner) cleanup() {
	now := time.Now()
	removed := 0

	c.UserStorage.Mu.Lock()
	for id, user := range c.UserStorage.Users {
		if now.Sub(user.LastActivity) > c.TTL {
			delete(c.UserStorage.Users, id)
			c.WaitingQueue.RemoveUserLocked(user.ID)
			if user.ChatID != 0 {
				err := c.ChatStorage.QuitChatLocked(user.ChatID, user.ID)
				if err != nil {
					log.Printf("[Cleaner] error quitting chat for user %s: %v", user.ID, err)
				}
			}
			removed++
		}
	}
	c.UserStorage.Mu.Unlock()

	c.ChatStorage.RemoveInactiveChatsLocked()

	log.Printf("[Cleaner] removed %d inactive users", removed)
}
