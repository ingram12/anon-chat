package handlers

import (
	"anon-chat/internal/api"
	"anon-chat/internal/config"
	"anon-chat/internal/pow"
	"anon-chat/internal/token"
	"anon-chat/internal/users"
	"net/http"

	"github.com/labstack/echo/v4"
)

func SolveFirstChallenge(
	ctx echo.Context,
	cfg *config.Config,
	storage *users.UserStorage,
	rotatingToken *token.RotatingToken,
) error {
	var req api.SolveFirstChallengeRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	if storage.IsUserExist(req.Token) {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "challenge already solved"})
	}

	userToken := req.Token
	globalToken, err := rotatingToken.GetRotatingToken()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": "global token generation failed"})
	}

	isVerified := pow.VerifyChallenge(userToken, globalToken, req.Challenge, cfg.TokenSecretKey)
	if !isVerified {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "challenge verification failed"})
	}

	if !pow.VerifyChallengeNonce(req.Challenge, req.Nonce, int(req.Difficulty)) {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "invalid PoW nonce"})
	}

	newUserToken := token.GenerateUserToken()
	newGlobalToken, err := rotatingToken.GetRotatingToken()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	newChallenge := pow.GenerateChallenge(newUserToken, newGlobalToken, cfg.TokenSecretKey)

	user, err := storage.CreateUser(userToken, newChallenge)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	resp := api.SolveFirstChallengeResponse{
		UserId:     string(user.ID[:]),
		Challenge:  newChallenge,
		Difficulty: int32(user.Difficulty),
		Token:      newUserToken,
	}
	return ctx.JSON(http.StatusOK, resp)
}
