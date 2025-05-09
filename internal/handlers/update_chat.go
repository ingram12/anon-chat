package handlers

import (
	"anon-chat/internal/api"
	"anon-chat/internal/chat"
	"anon-chat/internal/users"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func UpdateChat(ctx echo.Context, userID string, userStorage *users.UserStorage, chatStorage *chat.Storage) error {
	userStorage.Mu.RLock()
	_, exist := userStorage.GetUser(userID)
	if !exist {
		userStorage.Mu.RUnlock()
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "User not found"})
	}
	userStorage.Mu.RUnlock()

	waitChan := make(chan int, 1)

	go func() {
		ticker := time.NewTicker(200 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Request().Context().Done():
				return
			case <-ticker.C:
				userStorage.Mu.RLock()
				chatStorage.Mu.RLock()

				user, exist := userStorage.GetUser(userID)
				if !exist {
					chatStorage.Mu.RUnlock()
					userStorage.Mu.RUnlock()
					waitChan <- 0
					return
				}

				chatID := user.ChatID
				if !chatStorage.IsActiveChat(chatID) {
					chatStorage.Mu.RUnlock()
					userStorage.Mu.RUnlock()
					waitChan <- 0
					return
				}

				if chatStorage.HasNewMessages(chatID, user.ID) {
					chatStorage.Mu.RUnlock()
					userStorage.Mu.RUnlock()
					waitChan <- chatID
					return
				}

				chatStorage.Mu.RUnlock()
				userStorage.Mu.RUnlock()
			}
		}
	}()

	select {
	case chatID := <-waitChan:
		if chatID == 0 {
			resp := api.UpdateChatResponse{
				Status: "closed",
			}
			return ctx.JSON(http.StatusOK, resp)
		}

		userStorage.Mu.RLock()
		chatStorage.Mu.Lock()
		defer chatStorage.Mu.Unlock()
		defer userStorage.Mu.RUnlock()

		user, exist := userStorage.GetUser(userID)
		if !exist {
			return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "User not found"})
		}

		if user.ChatID != chatID || !chatStorage.IsActiveChat(chatID) {
			resp := api.UpdateChatResponse{
				Status: "closed",
			}
			return ctx.JSON(http.StatusOK, resp)
		}

		messages, _ := chatStorage.GetPeerMessages(chatID, user.ID)
		_ = chatStorage.RemovePeerMessages(chatID, user.ID) // TODO: check if this is correct

		respMessages := make([]api.ChatMessage, len(messages))
		for i, msg := range messages {
			respMessages[i] = api.ChatMessage{
				Message:   msg.Message,
				Timestamp: msg.Timestamp,
			}
		}

		resp := api.UpdateChatResponse{
			Status:   "active",
			Messages: respMessages,
		}
		return ctx.JSON(http.StatusOK, resp)
	case <-time.After(15 * time.Second): // TODO: make it configurable
	case <-ctx.Request().Context().Done():
	}

	resp := api.UpdateChatResponse{
		Status: "active",
	}
	return ctx.JSON(http.StatusOK, resp)
}
