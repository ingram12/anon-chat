package handlers

import (
	"anon-chat/internal/api"
	"anon-chat/internal/chat"
	"anon-chat/internal/users"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func WaitForChat(
	ctx echo.Context,
	userID string,
	userStorage *users.UserStorage,
	chatStorage *chat.Storage,
	waitingQueue *users.WaitingQueue,
) error {
	user, exist := userStorage.GetUser(userID)
	if !exist {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "User not found"})
	}

	waitingQueue.AddUser(user.ID)
	defer waitingQueue.RemoveUser(user.ID)

	waitChan := make(chan int, 1)

	go func() {
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Request().Context().Done():
				return
			case <-ticker.C:
				waitingQueue.TryMatch(chatStorage, userStorage)
				user, _ := userStorage.GetUser(userID)
				if user.ChatID != 0 {
					waitChan <- user.ChatID
					return
				}
			}
		}
	}()

	select {
	case chatID := <-waitChan:
		user, exist := userStorage.GetUser(userID)
		if !exist {
			return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "User not found"})
		}

		chat, err := chatStorage.GetChat(chatID)
		if err == nil && chat.IsUserInChat(user.ID) {
			peerID := chat.GetPeerID(user.ID)
			peerUser, exist := userStorage.GetUser(users.BytesToString(peerID))
			if exist && chat.IsUserInChat(peerUser.ID) {
				resp := api.WaitForChatResponse{
					Status:        "assigned",
					PeerPublicKey: &peerUser.PublicKey,
					Nickname:      &peerUser.Nickname,
				}
				return ctx.JSON(http.StatusOK, resp)
			}
		}
	case <-time.After(10 * time.Second): // TODO: make it configurable
	case <-ctx.Request().Context().Done():
	}

	resp := api.WaitForChatResponse{
		Status:        "waiting",
		PeerPublicKey: nil,
		Nickname:      nil,
	}
	return ctx.JSON(http.StatusOK, resp)
}
