package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/abdulyazidi/cloudtv/backend/internal/auth"
	"github.com/abdulyazidi/cloudtv/backend/internal/db/sqlc"
	pb "github.com/abdulyazidi/cloudtv/backend/pb/auth"
	"github.com/jackc/pgx/v5"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("cloudTV :)")
	ctx := context.Background()

	log.Println("Connecting to database...")
	connString := "postgres://postgres:postgres@localhost:5432/cloudtv"
	conn, err := pgx.Connect(ctx, connString)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close(ctx)

	queries := sqlc.New(conn)
	authService := auth.NewService(queries)
	log.Println("✓ Auth service created")

	authHandler := auth.NewHandler(authService)
	log.Println("✓ Auth handler created")

	grpcServer := grpc.NewServer()
	log.Println("✓ gRPC server created")

	pb.RegisterAuthServiceServer(grpcServer, authHandler)
	log.Println("✓ Auth service registered")

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen on port 50051: %v", err)
	}

	log.Println("🚀 cloudTV gRPC server listening on :50051")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
