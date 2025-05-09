package handlers

import (
	"anon-chat/internal/api"
	"anon-chat/internal/chat"
	"anon-chat/internal/users"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type waitChatResult struct {
	ChatID int
}

func WaitForChat(
	ctx echo.Context,
	userID string,
	userStorage *users.UserStorage,
	chatStorage *chat.Storage,
	waitingQueue *users.WaitingQueue,
) error {
	user, exist := userStorage.GetUserLocked(userID)
	if !exist {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "User not found"})
	}

	waitingQueue.AddUserLocked(user.ID)
	defer waitingQueue.RemoveUserLocked(user.ID)

	waitChan := make(chan waitChatResult, 1)

	go func() {
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Request().Context().Done():
				return
			case <-ticker.C:
				waitingQueue.TryMatch(chatStorage, userStorage)

				user, exist := userStorage.GetUserLocked(userID)
				if exist && user.ChatID != 0 {
					select {
					case waitChan <- waitChatResult{ChatID: user.ChatID}:
					case <-ctx.Request().Context().Done():
					}
					return
				}
			}
		}
	}()

	select {
	case result := <-waitChan:
		userStorage.Mu.Lock()
		chatStorage.Mu.Lock()
		defer userStorage.Mu.Unlock()
		defer chatStorage.Mu.Unlock()

		user, exist := userStorage.GetUser(userID)
		if !exist {
			return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "User not found"})
		}

		chat, _ := chatStorage.GetChat(result.ChatID)

		if chat.IsUserInChat(user.ID) {
			peerID := chat.GetPeerID(user.ID)
			peerUser, exist := userStorage.GetUser(peerID)
			if exist && chat.IsUserInChat(peerUser.ID) {
				resp := api.WaitForChatResponse{
					Status:        "assigned",
					PeerPublicKey: &peerUser.PublicKey,
					Nickname:      &peerUser.Nickname,
				}
				return ctx.JSON(http.StatusOK, resp)
			}
		}

	case <-time.After(15 * time.Second): // TODO: make configurable
	case <-ctx.Request().Context().Done():
		return ctx.NoContent(http.StatusRequestTimeout)
	}

	resp := api.WaitForChatResponse{
		Status:        "waiting",
		PeerPublicKey: nil,
		Nickname:      nil,
	}
	return ctx.JSON(http.StatusOK, resp)
}
