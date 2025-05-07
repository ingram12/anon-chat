package handlers

import (
	"anon-chat/internal/api"
	"anon-chat/internal/chat"
	"anon-chat/internal/users"
	"net/http"

	"github.com/labstack/echo/v4"
)

func UpdateChat(ctx echo.Context, userID string, userStorage *users.UserStorage, chatStorage *chat.Storage) error {
	user, exist := userStorage.GetUser(userID)
	if !exist {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "User not found"})
	}

	chatID := user.ChatID
	if chatID == 0 {
		resp := api.UpdateChatResponse{
			Status: "closed",
		}
		return ctx.JSON(http.StatusOK, resp)
	}

	if !chatStorage.IsActiveChat(chatID) {
		resp := api.UpdateChatResponse{
			Status: "closed",
		}
		return ctx.JSON(http.StatusOK, resp)
	}

	messages, err := chatStorage.GetPeerMessages(chatID, user.ID)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "User not in chat 4"})
	}
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
}
