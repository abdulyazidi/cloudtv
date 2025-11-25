package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/abdulyazidi/cloudtv/backend/internal/auth"
	"github.com/abdulyazidi/cloudtv/backend/internal/db/sqlc"
	pb "github.com/abdulyazidi/cloudtv/backend/pb/auth"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	fmt.Println("cloudTV :)")
	ctx := context.Background()
	if err := godotenv.Load(); err != nil {
		log.Fatal("Could not load .env file", err)
	}
	jwtSecret := os.Getenv("JWT_HMAC_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_HMAC_SECRET is missing")
	}

	log.Println("Connecting to database...")
	connString := "postgres://postgres:postgres@localhost:5432/cloudtv"
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	queries := sqlc.New(pool)
	authService := auth.NewService(queries, []byte(jwtSecret))
	log.Println("✓ Auth service created")

	authHandler := auth.NewHandler(authService)
	log.Println("✓ Auth handler created")

	grpcServer := grpc.NewServer()
	log.Println("✓ gRPC server created")

	pb.RegisterAuthServiceServer(grpcServer, authHandler)
	log.Println("✓ Auth service registered")
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen on port 50051: %v", err)
	}

	log.Println("🚀 cloudTV gRPC server listening on :50051")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
