package main

import (
	"anon-chat-backend/api/proto"
	"anon-chat-backend/internal/config"
	"anon-chat-backend/internal/users"
	"log"
	"net"
	"net/http"

	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"google.golang.org/grpc"
)

func main() {
	cfg := config.NewConfig()

	// Create gRPC server
	grpcServer := grpc.NewServer()
	userServer := users.NewGRPCServer(cfg)
	proto.RegisterUserServiceServer(grpcServer, userServer)

	// Wrap gRPC server with gRPC-web
	wrappedGrpc := grpcweb.WrapServer(grpcServer)

	// Create HTTP server with CORS support
	httpServer := &http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Set CORS headers
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, x-grpc-web")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			wrappedGrpc.ServeHTTP(w, r)
		}),
	}

	// Start server
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Println("Server starting on :50051")
	if err := httpServer.Serve(listener); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
