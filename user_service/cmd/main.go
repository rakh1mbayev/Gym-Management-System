package main

import (
	"database/sql"
	"fmt"
	mailNats "github.com/rakh1mbayev/Gym-Management-System/mail_service/pkg/nats"
	rpc "github.com/rakh1mbayev/Gym-Management-System/user_service/internal/delivery/grpc"
	"github.com/rakh1mbayev/Gym-Management-System/user_service/internal/delivery/grpc_client"
	userNats "github.com/rakh1mbayev/Gym-Management-System/user_service/internal/nats"
	"github.com/rakh1mbayev/Gym-Management-System/user_service/internal/repository/postgres"
	"github.com/rakh1mbayev/Gym-Management-System/user_service/internal/usecase"
	"github.com/rakh1mbayev/Gym-Management-System/user_service/proto/userpb"
	"log"
	"net"
	_ "os"

	_ "github.com/lib/pq" // Postgres driver
	"google.golang.org/grpc"
	_ "google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

func main() {

	// Set up the Postgres connection
	db, err := sql.Open("postgres", "postgres://postgres:12345678@localhost:5432/database?sslmode=disable")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	conn, err := grpc.Dial("localhost:8084", grpc.WithInsecure()) // adjust port
	if err != nil {
		log.Fatalf("Failed to connect to mail service: %v", err)
	}
	defer conn.Close()

	natsConn, err := mailNats.ConnectNATS() // Or move ConnectNATS to a shared pkg
	if err != nil {
		log.Fatal("Failed to connect to NATS:", err)
	}
	defer mailNats.Close()

	natsPublisher := userNats.NewNatsPublisher(natsConn)

	// Initialize repositories and use cases
	userRepo := postgres.NewUserRepository(db)
	mailClient := grpc_client.NewMailClient(conn)
	userUsecase := usecase.NewUserUsecase(userRepo, mailClient, natsPublisher)

	// Initialize the gRPC server and register the User Service
	grpcServer := grpc.NewServer()
	userServiceServer := rpc.NewUserServiceServer(userUsecase)
	userpb.RegisterUserServiceServer(grpcServer, userServiceServer)

	// Register reflection service on gRPC server (optional, for testing with gRPC CLI)
	reflection.Register(grpcServer)

	// Start the gRPC server
	port := ":8083" // You can change the port if necessary
	fmt.Printf("Starting gRPC server on port %s...\n", port)

	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to start listener: %v", err)
	}

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to start gRPC server: %v", err)
	}
}
