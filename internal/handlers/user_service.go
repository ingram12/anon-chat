package handlers

import (
	"anon-chat-backend/internal/config"
	"anon-chat-backend/internal/users"

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
	return GetFirstChallenge(ctx, s.config, s.storage)
}

func (s *UserService) SolveFirstChallenge(ctx echo.Context) error {
	return SolveFirstChallenge(ctx, s.config, s.storage)
}

func (s *UserService) RegisterUser(ctx echo.Context) error {
	return RegisterUser(ctx, s.config, s.storage)
}
