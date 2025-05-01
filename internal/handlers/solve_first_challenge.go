package handlers

import (
	"anon-chat/internal/api"
	"anon-chat/internal/config"
	"anon-chat/internal/pow"
	"anon-chat/internal/token"
	"anon-chat/internal/users"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

var (
	ErrInvalidSolution = errors.New("invalid solution")
	ErrInvalidToken    = errors.New("invalid token")
)

func SolveFirstChallenge(ctx echo.Context, config *config.Config, storage *users.UserStorage) error {
	var req api.SolveFirstChallengeRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	err := token.VerifyToken(req.Key, req.Challenge, config.TokenSecretKey)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	if !pow.VerifySolution(req.Challenge, req.Solution, int(req.Difficulty)) {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": ErrInvalidSolution.Error()})
	}

	newChallenge, err := pow.RandomKey() // TODO: make difficulty configurable
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	user, err := storage.CreateUser(newChallenge, int(req.Difficulty))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	resp := api.SolveFirstChallengeResponse{
		UserId: user.ID,
	}
	return ctx.JSON(http.StatusOK, resp)
}
