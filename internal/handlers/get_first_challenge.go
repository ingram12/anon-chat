package handlers

import (
	"anon-chat/internal/api"
	"anon-chat/internal/config"
	"anon-chat/internal/pow"
	"anon-chat/internal/token"
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetFirstChallenge(ctx echo.Context, config *config.Config) error {
	key, err := pow.RandomKey()
	if err != nil {
		return err
	}

	challenge, err := token.GenerateToken(key, config.TokenSecretKey)
	if err != nil {
		return err
	}

	resp := api.GetFirstChallengeResponse{
		Challenge:  challenge,
		Key:        key,
		Difficulty: int32(30),
	}

	return ctx.JSON(http.StatusOK, resp)
}
