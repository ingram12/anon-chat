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
	user, exist := storage.GetUser(userID)

	if !exist {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "user not found"})
	}

	if user.ChatID != 0 {
		chat, err := chatStorage.GetChat(user.ChatID)

		peerPublicKey := string(chat.UserID1[:])
		if err != nil {
			resp := api.WaitForChatResponse{
				Status:        "assigned",
				PeerPublicKey: &peerPublicKey,
			}
			return ctx.JSON(http.StatusOK, resp)
		}

		user.ChatID = 0
		storage.UpdateUser(user)
	}

	waitChan := make(chan int, 1)

	go func() {
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()

		for range ticker.C {
			user, _ := storage.GetUser(userID)

			if user.ChatID != 0 {
				waitChan <- user.ChatID
				return
			}
		}
	}()

	select {
	case chatID := <-waitChan:
		chat, err := chatStorage.GetChat(chatID)

		peerPublicKey := string(chat.UserID1[:])
		if err != nil {
			resp := api.WaitForChatResponse{
				Status:        "assigned",
				PeerPublicKey: &peerPublicKey,
			}
			return ctx.JSON(http.StatusOK, resp)
		}

		user.ChatID = 0
		storage.UpdateUser(user)

		resp := api.WaitForChatResponse{
			Status:        "waiting",
			PeerPublicKey: nil,
		}
		return ctx.JSON(http.StatusOK, resp)
	case <-time.After(4 * time.Second):
		resp := api.WaitForChatResponse{
			Status:        "waiting",
			PeerPublicKey: nil,
		}
		return ctx.JSON(http.StatusOK, resp)
	case <-ctx.Request().Context().Done():
		resp := api.WaitForChatResponse{
			Status:        "waiting",
			PeerPublicKey: nil,
		}
		return ctx.JSON(http.StatusOK, resp)
	}
}
