package handlers

import (
	"anon-chat-backend/internal/api"
	"anon-chat-backend/internal/config"
	"anon-chat-backend/internal/pow"
	"anon-chat-backend/internal/token"
	"anon-chat-backend/internal/users"
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

	if err := token.VerifyToken(req.Challenge+req.UserId, req.Token, config.TokenSecretKey); err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	if !pow.VerifySolution(req.Challenge, req.Solution, int(req.Difficulty)) {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": ErrInvalidSolution.Error()})
	}

	newChallenge, err := pow.GenerateChallenge(100) // TODO: make difficulty configurable
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	user, err := storage.CreateUser(req.UserId, newChallenge.Challenge, int(req.Difficulty))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	resp := api.SolveFirstChallengeResponse{
		UserId:     user.Id,
		Challenge:  user.CurrentChallenge,
		Difficulty: int32(user.Difficulty),
	}
	return ctx.JSON(http.StatusOK, resp)
}
