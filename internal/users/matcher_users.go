package users

import (
	"anon-chat/internal/chat"
	"log"
	"math/rand"
	"time"
)

func MatchUsersIntoChats(userStorage *UserStorage, chatStorage *chat.Storage) {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		users := userStorage.GetUsersWithoutChat()

		log.Printf("Found %d users without chat\n", len(users))

		if len(users) < 2 {
			continue
		}

		timeNow := time.Now()

		// Shuffle users randomly
		rand.Shuffle(len(users), func(i, j int) {
			users[i], users[j] = users[j], users[i]
		})

		// Create pairs and assign chats
		for i := 0; i < len(users)-1; i += 2 {
			user1 := users[i]
			user2 := users[i+1]

			chat, err := chatStorage.CreateChat(user1.ID, user2.ID)
			if err != nil {
				continue
			}

			user1.ChatID = chat.ID
			user1.LastActivity = timeNow
			userStorage.UpdateUser(user1)

			user2.ChatID = chat.ID
			user2.LastActivity = timeNow
			userStorage.UpdateUser(user2)

			log.Printf("Matched users %s and %s into chat %d\n", user1.GetUserID(), user2.GetUserID(), chat.ID)
		}
	}
}
