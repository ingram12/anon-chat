package handlers

import (
	"anon-chat/internal/api"
	"anon-chat/internal/pow"
	"anon-chat/internal/users"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

var (
	ErrUserNotFound          = errors.New("user not found")
	ErrInvalidChallenge      = errors.New("invalid challenge")
	ErrInvalidSolution       = errors.New("invalid solution")
	ErrUserAlreadyRegistered = errors.New("user already registered")
)

func RegisterUser(ctx echo.Context, storage *users.UserStorage) error {
	var req api.RegisterUserRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	user, exists := storage.GetUser(req.UserId)
	if !exists {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": ErrUserNotFound.Error()})
	}

	if user.IsRegistered {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": ErrUserAlreadyRegistered.Error()})
	}

	if user.CurrentChallenge != req.Challenge || user.Difficulty != int(req.Difficulty) {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": ErrInvalidChallenge.Error()})
	}

	if !pow.VerifyChallengeNonce(req.Challenge, req.Nonce, int(req.Difficulty)) {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": ErrInvalidSolution.Error()})
	}

	user.Nickname = req.Nickname
	user.Tags = req.Tags
	user.PublicKey = req.PublicKey
	user.IsRegistered = true
	storage.UpdateUser(user)

	resp := api.RegisterUserResponse{
		UserId:  string(user.ID[:]),
		Success: true,
		Message: "User registered successfully",
	}
	return ctx.JSON(http.StatusOK, resp)
}
