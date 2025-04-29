package handlers

import (
	"anon-chat/internal/api"
	"anon-chat/internal/config"
	"anon-chat/internal/pow"
	"anon-chat/internal/token"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func GetFirstChallenge(ctx echo.Context, config *config.Config) error {
	challenge, err := pow.GenerateChallenge(100) // TODO: make difficulty configurable
	if err != nil {
		return err
	}

	userID := uuid.New().String()

	token, err := token.GenerateToken(challenge.Challenge+userID, config.TokenSecretKey)
	if err != nil {
		return err
	}

	resp := api.GetFirstChallengeResponse{
		Challenge:  challenge.Challenge,
		Token:      token,
		Difficulty: int32(challenge.Difficulty),
		UserId:     userID,
	}
	return ctx.JSON(http.StatusOK, resp)
}
