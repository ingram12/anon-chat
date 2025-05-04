package handlers

import (
	"anon-chat/internal/api"
	"anon-chat/internal/chat"
	"anon-chat/internal/users"
	"net/http"

	"github.com/labstack/echo/v4"
)

func QuitChat(ctx echo.Context, userID string, userStorage *users.UserStorage, chatStorage *chat.Storage) error {
	user, exist := userStorage.GetUser(userID)
	if !exist {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "User not found"})
	}

	chatID := user.ChatID
	if chatID == 0 {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "User not in chat 2"})
	}

	err := chatStorage.QuitChat(chatID, user.ID)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "failed to quit chat"})
	}
	user.ChatID = 0
	userStorage.UpdateUser(user)

	resp := api.QuitChatResponse{
		Success: true,
	}
	return ctx.JSON(http.StatusOK, resp)
}
