package users

import (
	"anon-chat-backend/api/proto"
	"anon-chat-backend/internal/config"
	"anon-chat-backend/internal/pow"
	"anon-chat-backend/internal/token"
	"context"
	"log"
	"testing"
)

func TestGRPCGetFirstChallenge(t *testing.T) {
	log.Println("Starting TestGRPCGetFirstChallenge")

	// Create test config
	cfg := config.NewConfig()
	log.Printf("Created test config with secret key: %s", cfg.TokenSecretKey)

	// Create gRPC server
	server := NewGRPCServer(cfg)
	log.Println("Created gRPC server instance")

	// Create test context
	ctx := context.Background()

	// Call GetFirstChallenge
	resp, err := server.GetFirstChallenge(ctx, &proto.GetFirstChallengeRequest{})
	if err != nil {
		t.Fatalf("GetFirstChallenge failed: %v", err)
	}

	log.Printf("Response:")
	log.Printf("  Challenge: %s", resp.Challenge)
	log.Printf("  Token: %s", resp.Token)
	log.Printf("  UserId: %s", resp.UserId)
	log.Printf("  Difficulty: %d", resp.Difficulty)

	// Verify challenge is not empty
	if resp.Challenge == "" {
		t.Error("Expected non-empty challenge")
	}

	// Verify token is valid
	if err := token.VerifyToken(resp.Challenge+resp.UserId, resp.Token, cfg.TokenSecretKey); err != nil {
		t.Errorf("Generated token is invalid: %v", err)
	} else {
		log.Println("Token verification successful")
	}

	log.Println("TestGRPCGetFirstChallenge completed")
}

func TestGRPCSolveFirstChallenge(t *testing.T) {
	log.Println("Starting TestGRPCSolveFirstChallenge")

	// Create test config
	cfg := config.NewConfig()
	log.Printf("Created test config with secret key: %s", cfg.TokenSecretKey)

	// Create gRPC server
	server := NewGRPCServer(cfg)
	log.Println("Created gRPC server instance")

	// Create test context
	ctx := context.Background()

	// First get a challenge
	challengeResp, err := server.GetFirstChallenge(ctx, &proto.GetFirstChallengeRequest{})
	if err != nil {
		t.Fatalf("GetFirstChallenge failed: %v", err)
	}

	solution, err := pow.SolveChallenge(challengeResp.Challenge, int(challengeResp.Difficulty))
	if err != nil {
		t.Fatalf("SolveChallenge failed: %v", err)
	}

	// Create solve request
	solveReq := &proto.SolveFirstChallengeRequest{
		Challenge:  challengeResp.Challenge,
		Solution:   solution,
		Token:      challengeResp.Token,
		Difficulty: challengeResp.Difficulty,
		UserId:     challengeResp.UserId,
	}

	// Call SolveFirstChallenge
	solveResp, err := server.SolveFirstChallenge(ctx, solveReq)
	if err != nil {
		t.Fatalf("SolveFirstChallenge failed: %v", err)
	}

	log.Printf("Response:")
	log.Printf("  UserID: %s", solveResp.UserId)
	log.Printf("  CurrentChallenge: %s", solveResp.Challenge)

	// Verify response
	if solveResp.UserId == "" {
		t.Error("Expected non-empty user ID")
	}
	if solveResp.Challenge == "" {
		t.Error("Expected non-empty current challenge")
	}

	log.Println("TestGRPCSolveFirstChallenge completed")
}
