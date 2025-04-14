package main

import (
	"database/sql"
	"log"
	"net"

	"github.com/rakh1mbayev/Gym-Management-System/proto/userpb"
	"user_service/internal/delivery/grpc"
	"user_service/internal/repository/postgres"
	"user_service/internal/usecase"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"
)

func main() {
	// Connect to DB
	db, err := sql.Open("postgres", "postgres://user:pass@localhost:5432/gymdb?sslmode=disable")
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}
	defer db.Close()

	// Init layers
	repo := postgres.NewUserRepo(db)
	uc := usercase.NewUsecase(repo)
	server := &grpc.Server{UC: uc}

	// gRPC Server
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	userpb.RegisterUserServiceServer(s, server)

	log.Println("User service running on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatal("failed to serve:", err)
	}
}
