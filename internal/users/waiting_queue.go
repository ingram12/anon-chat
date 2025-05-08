package users

import (
	"anon-chat/internal/chat"
	"errors"
	"sync"
	"time"
)

type WaitingQueue struct {
	Mu         sync.RWMutex
	users      map[string]time.Time
	isMatching bool
}

func NewWaitingQueue() *WaitingQueue {
	return &WaitingQueue{
		users:      make(map[string]time.Time),
		isMatching: false,
	}
}

func (wq *WaitingQueue) AddUser(userID string) {
	wq.users[userID] = time.Now()
}

func (wq *WaitingQueue) RemoveUser(userID string) {
	delete(wq.users, userID)
}

func (wq *WaitingQueue) AddUserLocked(userID string) {
	wq.Mu.Lock()
	defer wq.Mu.Unlock()
	wq.AddUser(userID)
}

func (wq *WaitingQueue) RemoveUserLocked(userID string) {
	wq.Mu.Lock()
	defer wq.Mu.Unlock()
	wq.RemoveUser(userID)
}

func (wq *WaitingQueue) GetLen() int {
	return len(wq.users)
}

func (wq *WaitingQueue) GetTwoRandomUsers() (string, string, error) {
	if len(wq.users) < 2 {
		return "", "", errors.New("not enough users")
	}

	userIDs := make([]string, 0, 2)
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
	wq.Mu.Lock()
	if wq.isMatching {
		wq.Mu.Unlock()
		return // Already matching
	}
	wq.isMatching = true
	wq.Mu.Unlock()

	go func() {
		defer func() {
			wq.Mu.Lock()
			wq.isMatching = false
			wq.Mu.Unlock()
		}()

		MatchUsers(userStorage, chatStorage, wq)
	}()
}
