package handlers

import (
	"anon-chat/internal/config"
	"anon-chat/internal/token"
	"anon-chat/internal/users"

	"github.com/labstack/echo/v4"
)

type UserService struct {
	config        *config.Config
	storage       *users.UserStorage
	rotatingToken *token.RotatingToken
}

func NewUserService(config *config.Config) *UserService {
	return &UserService{
		config:        config,
		storage:       users.NewUserStorage(),
		rotatingToken: token.NewRotatingToken(config.SecondLifeTime),
	}
}

func (s *UserService) GetFirstChallenge(ctx echo.Context) error {
	return GetFirstChallenge(ctx, s.config, s.rotatingToken)
}

func (s *UserService) SolveFirstChallenge(ctx echo.Context) error {
	return SolveFirstChallenge(ctx, s.config, s.storage, s.rotatingToken)
}

func (s *UserService) RegisterUser(ctx echo.Context) error {
	return RegisterUser(ctx, s.storage)
}
