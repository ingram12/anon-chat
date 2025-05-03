package handlers

import (
	"anon-chat/internal/api"
	"anon-chat/internal/chat"
	"anon-chat/internal/users"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func WaitForChat(ctx echo.Context, userID string, storage *users.UserStorage, chatStorage *chat.Storage) error {
	if !storage.IsUserExist(userID) {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "User not found"})
	}

	waitChan := make(chan int, 1)

	go func() {
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Request().Context().Done():
				return
			case <-ticker.C:
				user, _ := storage.GetUser(userID)
				if user.ChatID != 0 {
					waitChan <- user.ChatID
					return
				}
			}
		}
	}()

	select {
	case chatID := <-waitChan:
		user, exist := storage.GetUser(userID)
		if !exist {
			return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "User not found"})
		}

		chat, err := chatStorage.GetChat(chatID)
		if err == nil && chat.IsUserInChat(user.ID) {
			peerID := chat.GetPeerID(user.ID)
			peerUser, exist := storage.GetUser(users.BytesToString(peerID))
			if exist {
				resp := api.WaitForChatResponse{
					Status:        "assigned",
					PeerPublicKey: &peerUser.PublicKey,
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
	}
	return ctx.JSON(http.StatusOK, resp)
}
