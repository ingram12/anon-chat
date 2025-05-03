package handlers

import (
	"anon-chat/internal/chat"
	"anon-chat/internal/config"
	"anon-chat/internal/token"
	"anon-chat/internal/users"

	"github.com/labstack/echo/v4"
)

type UserService struct {
	cfg           *config.Config
	storage       *users.UserStorage
	rotatingToken *token.RotatingToken
	chatStorage   *chat.Storage
}

func NewUserService(cfg *config.Config) *UserService {
	userService := UserService{
		cfg:           cfg,
		storage:       users.NewUserStorage(cfg.UserInactivityTimeout),
		rotatingToken: token.NewRotatingToken(cfg.RotatingTokenLifeTime),
		chatStorage:   chat.NewChatStorage(),
	}

	go userService.storage.MatchUsersIntoChats(userService.chatStorage)

	return &userService
}

func (s *UserService) GetFirstChallenge(ctx echo.Context) error {
	return GetFirstChallenge(ctx, s.cfg, s.rotatingToken)
}

func (s *UserService) SolveFirstChallenge(ctx echo.Context) error {
	return SolveFirstChallenge(ctx, s.cfg, s.storage, s.rotatingToken)
}

func (s *UserService) RegisterUser(ctx echo.Context) error {
	return RegisterUser(ctx, s.storage)
}

func (s *UserService) WaitForChat(ctx echo.Context, userID string) error {
	return WaitForChat(ctx, userID, s.storage, s.chatStorage)
}

func (s *UserService) UpdateChat(ctx echo.Context, userID string) error {
	return UpdateChat(ctx, userID, s.storage, s.chatStorage)
}

func (s *UserService) SendChatMessage(ctx echo.Context, userID string) error {
	return SendChatMessage(ctx, userID, s.storage, s.chatStorage)
}
