package handlers

import (
	"anon-chat/internal/api"
	"anon-chat/internal/chat"
	"anon-chat/internal/users"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type waitResult struct {
	ChatID int
	Closed bool
}

func UpdateChat(ctx echo.Context, userID string, userStorage *users.UserStorage, chatStorage *chat.Storage) error {
	if !userStorage.UpdateLastActivityLocked(userID) {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "User not found"})
	}

	waitChan := make(chan waitResult, 1)

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
					select {
					case waitChan <- waitResult{Closed: true}:
					case <-ctx.Request().Context().Done():
					}
					return
				}

				chatID := user.ChatID
				if !chatStorage.IsActiveChat(chatID) {
					chatStorage.Mu.RUnlock()
					userStorage.Mu.RUnlock()
					select {
					case waitChan <- waitResult{Closed: true}:
					case <-ctx.Request().Context().Done():
					}
					return
				}

				if chatStorage.HasNewMessages(chatID, user.ID) {
					chatStorage.Mu.RUnlock()
					userStorage.Mu.RUnlock()
					select {
					case waitChan <- waitResult{ChatID: chatID}:
					case <-ctx.Request().Context().Done():
					}
					return
				}

				chatStorage.Mu.RUnlock()
				userStorage.Mu.RUnlock()
			}
		}
	}()

	select {
	case result := <-waitChan:
		if result.Closed {
			return ctx.JSON(http.StatusOK, api.UpdateChatResponse{Status: "closed"})
		}

		userStorage.Mu.RLock()
		chatStorage.Mu.Lock()
		defer userStorage.Mu.RUnlock()
		defer chatStorage.Mu.Unlock()

		user, exist := userStorage.GetUser(userID)
		if !exist {
			return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "User not found"})
		}

		if user.ChatID != result.ChatID || !chatStorage.IsActiveChat(result.ChatID) {
			return ctx.JSON(http.StatusOK, api.UpdateChatResponse{Status: "closed"})
		}

		messages, _ := chatStorage.GetPeerMessages(result.ChatID, user.ID)
		err := chatStorage.RemovePeerMessages(result.ChatID, user.ID) // TODO: add aprove delivery message
		if err != nil {
			log.Printf("[UpdateChat] error removing messages: %v", err)
		}

		respMessages := make([]api.ChatMessage, len(messages))
		for i, msg := range messages {
			respMessages[i] = api.ChatMessage{
				Message:   msg.Message,
				Timestamp: msg.Timestamp,
			}
		}

		return ctx.JSON(http.StatusOK, api.UpdateChatResponse{
			Status:   "active",
			Messages: respMessages,
		})
	case <-time.After(15 * time.Second): // TODO: move to config
		return ctx.JSON(http.StatusOK, api.UpdateChatResponse{Status: "active"})
	case <-ctx.Request().Context().Done():
		return ctx.NoContent(http.StatusRequestTimeout)
	}
}
