package main

import (
	"database/sql"
	"log"
	"net"

	_ "github.com/lib/pq"
	rpc "github.com/rakh1mbayev/Gym-Management-System/inventory_service/internal/delivery/grpc"
	"github.com/rakh1mbayev/Gym-Management-System/inventory_service/internal/repository/postgres"
	"github.com/rakh1mbayev/Gym-Management-System/inventory_service/internal/usecase"
	"github.com/rakh1mbayev/Gym-Management-System/inventory_service/proto/inventorypb"
	"google.golang.org/grpc"
)

func main() {

	// Set up the Postgres connection
	db, err := sql.Open("postgres", "postgres://postgres:12345678@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	repo := postgres.NewProductRepository(db) // Assume initialized with DB connection
	uc := usecase.NewProductUsecase(repo)
	server := rpc.NewInventoryServer(uc)

	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	inventorypb.RegisterInventoryServiceServer(s, server)

	log.Println("gRPC Inventory Service started on port :8081")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
