package handlers

import (
	"anon-chat/internal/api"
	"anon-chat/internal/config"
	"anon-chat/internal/pow"
	"anon-chat/internal/token"
	"anon-chat/internal/users"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func SolveFirstChallenge(
	ctx echo.Context,
	cfg *config.Config,
	userStorage *users.UserStorage,
	rotatingToken *token.RotatingToken,
) error {
	var req api.SolveFirstChallengeRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	userStorage.Mu.Lock()
	defer userStorage.Mu.Unlock()

	if userStorage.IsUserExist(req.Token) {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "challenge already solved"})
	}

	userToken := req.Token
	globalToken := rotatingToken.GetRotatingToken()

	isVerified := pow.VerifyFirstChallenge(userToken, globalToken, req.Challenge, cfg.TokenSecretKey)
	if !isVerified {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "challenge verification failed"})
	}

	if !pow.VerifyChallengeNonce(req.Challenge, req.Nonce, int(req.Difficulty)) {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "invalid PoW nonce"})
	}

	newChallenge := pow.GenerateChallenge()
	timeNow := time.Now()

	user := userStorage.CreateUser(
		users.User{
			ID:               userToken,
			CreatedAt:        timeNow,
			LastActivity:     timeNow,
			CurrentChallenge: newChallenge,
			IsRegistered:     false,
			Difficulty:       int(req.Difficulty/3 + 1), // Increase difficulty for the next challenge
		},
	)

	resp := api.SolveFirstChallengeResponse{
		UserId:     user.ID,
		Challenge:  newChallenge,
		Difficulty: int32(user.Difficulty),
	}
	return ctx.JSON(http.StatusOK, resp)
}
