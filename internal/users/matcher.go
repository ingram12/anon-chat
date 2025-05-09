package users

import (
	"anon-chat/internal/chat"
	"time"
)

func MatchUsers(userStorage *UserStorage, chatStorage *chat.Storage, waitingQueue *WaitingQueue) {
	userStorage.Mu.Lock()
	waitingQueue.Mu.Lock()
	defer waitingQueue.Mu.Unlock()
	defer userStorage.Mu.Unlock()

	for {
		userID1, userID2, err := waitingQueue.GetTwoRandomUsers()
		if err != nil {
			return // No valid users to match
		}

		user1, exist := userStorage.GetUser(userID1)
		if !exist || user1.ChatID != 0 {
			waitingQueue.RemoveUser(userID1)
			continue
		}
		user2, exist := userStorage.GetUser(userID2)
		if !exist || user2.ChatID != 0 {
			waitingQueue.RemoveUser(userID2)
			continue
		}

		chatStorage.Mu.Lock()
		chat := chatStorage.CreateChat(userID1, userID2)
		chatStorage.Mu.Unlock()

		timeNow := time.Now()

		user1.ChatID = chat.ID
		user1.LastActivity = timeNow
		userStorage.Users[user1.ID] = user1

		user2.ChatID = chat.ID
		user2.LastActivity = timeNow
		userStorage.Users[user2.ID] = user2

		waitingQueue.RemoveUser(userID1)
		waitingQueue.RemoveUser(userID2)
	}
}
