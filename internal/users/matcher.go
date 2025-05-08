package users

import (
	"anon-chat/internal/chat"
	"time"
)

func MatchUsers(userStorage *UserStorage, chatStorage *chat.Storage, waitingQueue *WaitingQueue) {
	userStorage.Mu.Lock()
	chatStorage.Mu.Lock()
	waitingQueue.Mu.Lock()
	defer waitingQueue.Mu.Unlock()
	defer chatStorage.Mu.Unlock()
	defer userStorage.Mu.Unlock()

	for {
		tt := waitingQueue.GetLen()
		if tt < 2 {
			return // Not enough users to match
		}

		userID1, userID2, err := waitingQueue.GetTwoRandomUsers()
		if err != nil {
			return // No valid users to match
		}

		user1, exist := userStorage.GetUserBytes(userID1)
		if !exist || user1.ChatID != 0 {
			waitingQueue.RemoveUser(userID1)
			continue
		}
		user2, exist := userStorage.GetUserBytes(userID2)
		if !exist || user2.ChatID != 0 {
			waitingQueue.RemoveUser(userID2)
			continue
		}

		chat, _ := chatStorage.CreateChat(userID1, userID2)

		timeNow := time.Now()

		user1.ChatID = chat.ID
		user1.LastActivity = timeNow
		userStorage.users[user1.ID] = user1

		user2.ChatID = chat.ID
		user2.LastActivity = timeNow
		userStorage.users[user2.ID] = user2

		waitingQueue.RemoveUser(userID1)
		waitingQueue.RemoveUser(userID2)
	}
}
