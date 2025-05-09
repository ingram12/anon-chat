package handlers

import (
	"anon-chat/internal/chat"
	"anon-chat/internal/config"
	"anon-chat/internal/maintenance"
	"anon-chat/internal/token"
	"anon-chat/internal/users"

	"github.com/labstack/echo/v4"
)

type Server struct {
	cfg           *config.Config
	userStorage   *users.UserStorage
	rotatingToken *token.RotatingToken
	chatStorage   *chat.Storage
	waitingQueue  *users.WaitingQueue
}

func NewServer(cfg *config.Config) *Server {
	server := Server{
		cfg:           cfg,
		userStorage:   users.NewUserStorage(cfg.UserInactivityTimeout),
		rotatingToken: token.NewRotatingToken(cfg.RotatingTokenLifeTime),
		waitingQueue:  users.NewWaitingQueue(),
		chatStorage:   chat.NewChatStorage(),
	}

	server.StartCleaner()
	return &server
}

func (s *Server) GetFirstChallenge(ctx echo.Context) error {
	return GetFirstChallenge(ctx, s.cfg, s.rotatingToken)
}

func (s *Server) SolveFirstChallenge(ctx echo.Context) error {
	return SolveFirstChallenge(ctx, s.cfg, s.userStorage, s.rotatingToken)
}

func (s *Server) RegisterUser(ctx echo.Context) error {
	return RegisterUser(ctx, s.userStorage)
}

func (s *Server) WaitForChat(ctx echo.Context, userID string) error {
	return WaitForChat(ctx, userID, s.userStorage, s.chatStorage, s.waitingQueue)
}

func (s *Server) UpdateChat(ctx echo.Context, userID string) error {
	return UpdateChat(ctx, userID, s.userStorage, s.chatStorage)
}

func (s *Server) SendChatMessage(ctx echo.Context, userID string) error {
	return SendChatMessage(ctx, userID, s.userStorage, s.chatStorage)
}

func (s *Server) QuitChat(ctx echo.Context, userID string) error {
	return QuitChat(ctx, userID, s.userStorage, s.chatStorage)
}

func (s *Server) StartCleaner() {
	cleaner := maintenance.NewCleaner(
		s.userStorage,
		s.waitingQueue,
		s.chatStorage,
		s.cfg.UserInactivityTimeout,
		s.cfg.RotatingTokenLifeTime,
	)
	cleaner.Start()
}
