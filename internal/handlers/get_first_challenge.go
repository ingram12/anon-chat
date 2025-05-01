package handlers

import (
	"anon-chat/internal/api"
	"anon-chat/internal/config"
	"anon-chat/internal/pow"
	"anon-chat/internal/token"
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetFirstChallenge(
	ctx echo.Context,
	cfg *config.Config,
	rotatingToken *token.RotatingToken,
) error {
	userToken, err := token.GenerateUserToken()
	if err != nil {
		return err
	}

	globalToken, err := rotatingToken.GetRotatingToken()
	if err != nil {
		return err
	}

	challenge := pow.GenerateChallenge(userToken, globalToken, cfg.TokenSecretKey)

	resp := api.GetFirstChallengeResponse{
		Challenge:  challenge,
		Token:      userToken,
		Difficulty: int32(30), // TODO: make it configurable
	}

	return ctx.JSON(http.StatusOK, resp)
}
