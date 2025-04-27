package users

import (
	"anon-chat-backend/api/proto"
	"anon-chat-backend/internal/config"
	"anon-chat-backend/internal/pow"
	"anon-chat-backend/internal/token"
	"context"

	"github.com/google/uuid"
)

type grpcServer struct {
	proto.UnimplementedUserServiceServer
	storage *UserStorage
	config  *config.Config
}

func NewGRPCServer(cfg *config.Config) *grpcServer {
	return &grpcServer{
		storage: NewUserStorage(),
		config:  cfg,
	}
}

func (s *grpcServer) GetFirstChallenge(ctx context.Context, req *proto.GetFirstChallengeRequest) (*proto.GetFirstChallengeResponse, error) {
	challenge, err := pow.GenerateChallenge(100) // TODO: make difficulty configurable
	if err != nil {
		return nil, err
	}

	userId := uuid.New().String()

	token, err := token.GenerateToken(challenge.Challenge+userId, s.config.TokenSecretKey)
	if err != nil {
		return nil, err
	}

	return &proto.GetFirstChallengeResponse{
		Challenge:  challenge.Challenge,
		Token:      token,
		Difficulty: int32(challenge.Difficulty),
		UserId:     userId,
	}, nil
}

func (s *grpcServer) SolveFirstChallenge(ctx context.Context, req *proto.SolveFirstChallengeRequest) (*proto.SolveFirstChallengeResponse, error) {
	if err := token.VerifyToken(req.Challenge+req.UserId, req.Token, s.config.TokenSecretKey); err != nil {
		return nil, err
	}

	if !pow.VerifySolution(req.Challenge, req.Solution, int(req.Difficulty)) {
		return nil, ErrInvalidSolution
	}

	newChallenge, err := pow.GenerateChallenge(100) // TODO: make difficulty configurable
	if err != nil {
		return nil, err
	}

	user, err := s.storage.CreateUser(req.UserId, newChallenge.Challenge, int(req.Difficulty))
	if err != nil {
		return nil, err
	}

	return &proto.SolveFirstChallengeResponse{
		UserId:     user.Id,
		Challenge:  user.CurrentChallenge,
		Difficulty: int32(user.Difficulty),
	}, nil
}

func (s *grpcServer) RegisterUser(ctx context.Context, req *proto.RegisterUserRequest) (*proto.RegisterUserResponse, error) {
	user, exists := s.storage.GetUser(req.UserId)
	if !exists {
		return nil, ErrUserNotFound
	}

	if user.CurrentChallenge != req.Challenge {
		return nil, ErrInvalidChallenge
	}

	if !pow.VerifySolution(req.Challenge, req.Solution, 1) {
		return nil, ErrInvalidSolution
	}

	user.Nickname = req.Nickname
	user.Tags = req.Tags
	user.PublicKey = req.PublicKey
	user.IsRegistered = true
	s.storage.UpdateUser(user)

	return &proto.RegisterUserResponse{
		UserId:  user.Id,
		Success: true,
		Message: "User registered successfully",
	}, nil
}
