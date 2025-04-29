package handlers

import (
	"anon-chat-backend/internal/api"
	"anon-chat-backend/internal/config"
	"anon-chat-backend/internal/pow"
	"anon-chat-backend/internal/users"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

var (
	ErrUserNotFound     = errors.New("user not found")
	ErrInvalidChallenge = errors.New("invalid challenge")
)

func RegisterUser(ctx echo.Context, config *config.Config, storage *users.UserStorage) error {
	var req api.RegisterUserRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	user, exists := storage.GetUser(req.UserId)
	if !exists {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": ErrUserNotFound.Error()})
	}

	if user.CurrentChallenge != req.Challenge {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": ErrInvalidChallenge.Error()})
	}

	if !pow.VerifySolution(req.Challenge, req.Solution, int(req.Difficulty)) {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": ErrInvalidSolution.Error()})
	}

	user.Nickname = req.Nickname
	user.Tags = req.Tags
	user.PublicKey = req.PublicKey
	user.IsRegistered = true
	storage.UpdateUser(user)

	resp := api.RegisterUserResponse{
		UserId:  user.Id,
		Success: true,
		Message: "User registered successfully",
	}
	return ctx.JSON(http.StatusOK, resp)
}
