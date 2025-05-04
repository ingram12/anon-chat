package users

import (
	"anon-chat/internal/chat"
	"errors"
	"sync"
	"time"
)

type WaitingQueue struct {
	mu         sync.RWMutex
	users      map[[36]byte]time.Time
	isMatching bool
}

func NewWaitingQueue() *WaitingQueue {
	return &WaitingQueue{
		users:      make(map[[36]byte]time.Time),
		isMatching: false,
	}
}

func (wq *WaitingQueue) AddUser(userID [36]byte) {
	wq.mu.Lock()
	defer wq.mu.Unlock()
	wq.users[userID] = time.Now()
}

func (wq *WaitingQueue) RemoveUser(userID [36]byte) {
	wq.mu.Lock()
	defer wq.mu.Unlock()
	delete(wq.users, userID)
}

func (wq *WaitingQueue) GetLen() int {
	wq.mu.RLock()
	defer wq.mu.RUnlock()
	return len(wq.users)
}

func (wq *WaitingQueue) GetTwoRandomUsers() ([36]byte, [36]byte, error) {
	wq.mu.RLock()
	defer wq.mu.RUnlock()

	if len(wq.users) < 2 {
		return [36]byte{}, [36]byte{}, errors.New("not enough users")
	}

	userIDs := make([][36]byte, 0, 2)
	count := 0
	for id := range wq.users {
		userIDs = append(userIDs, id)
		count++
		if count >= 2 {
			break
		}
	}

	return userIDs[0], userIDs[1], nil
}

func (wq *WaitingQueue) TryMatch(chatStorage *chat.Storage, userStorage *UserStorage) {
	wq.mu.Lock()
	if wq.isMatching {
		wq.mu.Unlock()
		return // Already matching
	}
	wq.isMatching = true
	wq.mu.Unlock()

	go func() {
		defer func() {
			wq.mu.Lock()
			wq.isMatching = false
			wq.mu.Unlock()
		}()

		MatchUsers(userStorage, chatStorage, wq)
	}()
}
