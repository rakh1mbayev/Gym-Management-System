package main

import (
	"database/sql"
	"log"
	"net"

	_ "github.com/lib/pq" // Import the PostgreSQL driver
	rpc "github.com/rakh1mbayev/Gym-Management-System/order_service/internal/delivery/grpc"
	"github.com/rakh1mbayev/Gym-Management-System/order_service/internal/repository/postgres"
	"github.com/rakh1mbayev/Gym-Management-System/order_service/internal/usecase"
	"github.com/rakh1mbayev/Gym-Management-System/order_service/proto/orderpb"
	"google.golang.org/grpc"
)

func main() {
	// Set up the Postgres connection
	db, err := sql.Open("postgres", "postgres://postgres:12345678@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Verify that the database connection is established
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// Initialize the repository and usecase
	repo := postgres.NewOrderRepository(db) // Use the actual db init
	uc := usecase.NewOrderUsecase(repo)

	// Set up the gRPC server
	srv := rpc.NewOrderServiceServer(uc)

	// Create a listener on port 50052
	lis, err := net.Listen("tcp", ":8082")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	orderpb.RegisterOrderServiceServer(grpcServer, srv)

	// Log server startup
	log.Println("Order Service gRPC server started on port 8082")

	// Serve the gRPC server
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
