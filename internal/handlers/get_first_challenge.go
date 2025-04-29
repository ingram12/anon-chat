package handlers

import (
	"anon-chat-backend/internal/api"
	"anon-chat-backend/internal/config"
	"anon-chat-backend/internal/pow"
	"anon-chat-backend/internal/token"
	"anon-chat-backend/internal/users"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func GetFirstChallenge(ctx echo.Context, config *config.Config, storage *users.UserStorage) error {
	challenge, err := pow.GenerateChallenge(100) // TODO: make difficulty configurable
	if err != nil {
		return err
	}

	userId := uuid.New().String()

	token, err := token.GenerateToken(challenge.Challenge+userId, config.TokenSecretKey)
	if err != nil {
		return err
	}

	resp := api.GetFirstChallengeResponse{
		Challenge:  challenge.Challenge,
		Token:      token,
		Difficulty: int32(challenge.Difficulty),
		UserId:     userId,
	}
	return ctx.JSON(http.StatusOK, resp)
}
