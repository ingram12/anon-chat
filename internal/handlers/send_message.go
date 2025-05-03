package handlers

import (
	"anon-chat/internal/api"
	"anon-chat/internal/chat"
	"anon-chat/internal/users"
	"net/http"

	"github.com/labstack/echo/v4"
)

func SendChatMessage(ctx echo.Context, userID string, storage *users.UserStorage, chatStorage *chat.Storage) error {
	var req api.SendChatMessageRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	user, exists := storage.GetUser(userID)
	if !exists {
		return ctx.JSON(http.StatusNotFound, echo.Map{"error": "User not found"})
	}

	if user.ChatID == 0 {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "User not in chat"})
	}

	timestamp, err := chatStorage.AddMessage(user.ChatID, user.ID, req.Message)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	resp := api.SendChatMessageResponse{
		Success:   true,
		Timestamp: timestamp,
	}
	return ctx.JSON(http.StatusOK, resp)
}
