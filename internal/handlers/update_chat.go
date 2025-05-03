package handlers

import (
	"anon-chat/internal/api"
	"anon-chat/internal/chat"
	"anon-chat/internal/users"
	"net/http"

	"github.com/labstack/echo/v4"
)

func UpdateChat(ctx echo.Context, userID string, storage *users.UserStorage, chatStorage *chat.Storage) error {
	user, exist := storage.GetUser(userID)
	if !exist {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "User not found"})
	}

	chatID := user.ChatID
	if chatID == 0 {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "User not in chat 2"})
	}

	chat, err := chatStorage.GetChat(chatID)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "chat not found"})
	}

	if !chat.IsUserInChat(user.ID) {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "User not in chat"})
	}

	messages, err := chatStorage.GetPeerMessages(chatID, user.ID)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "failed to get messages"})
	}

	respMessages := make([]api.ChatMessage, len(messages))
	for i, msg := range messages {
		respMessages[i] = api.ChatMessage{
			Message:   msg.Message,
			Timestamp: &msg.CreatedAt,
		}
	}

	resp := api.UpdateChatResponse{
		Status:   "active",
		Messages: respMessages,
	}
	return ctx.JSON(http.StatusOK, resp)
}
