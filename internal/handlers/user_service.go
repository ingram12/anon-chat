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

type UserService struct {
	config  *config.Config
	storage *users.UserStorage
}

func NewUserService(config *config.Config) *UserService {
	return &UserService{
		config:  config,
		storage: users.NewUserStorage(),
	}
}

func (s *UserService) GetFirstChallenge(ctx echo.Context) error {
	challenge, err := pow.GenerateChallenge(100) // TODO: make difficulty configurable
	if err != nil {
		return err
	}

	userId := uuid.New().String()

	token, err := token.GenerateToken(challenge.Challenge+userId, s.config.TokenSecretKey)
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

func (s *UserService) SolveFirstChallenge(ctx echo.Context) error {
	var req api.SolveFirstChallengeRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	if err := token.VerifyToken(req.Challenge+req.UserId, req.Token, s.config.TokenSecretKey); err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	if !pow.VerifySolution(req.Challenge, req.Solution, int(req.Difficulty)) {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": ErrInvalidSolution.Error()})
	}

	newChallenge, err := pow.GenerateChallenge(100) // TODO: make difficulty configurable
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	user, err := s.storage.CreateUser(req.UserId, newChallenge.Challenge, int(req.Difficulty))
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

func (s *UserService) RegisterUser(ctx echo.Context) error {
	var req api.RegisterUserRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	user, exists := s.storage.GetUser(req.UserId)
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
	s.storage.UpdateUser(user)

	resp := api.RegisterUserResponse{
		UserId:  user.Id,
		Success: true,
		Message: "User registered successfully",
	}
	return ctx.JSON(http.StatusOK, resp)
}
