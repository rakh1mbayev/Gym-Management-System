package main

import (
	"database/sql"
	"github.com/nats-io/nats.go"
	"github.com/rakh1mbayev/Gym-Management-System/inventory_service/internal/subscriber"
	"github.com/redis/go-redis/v9"
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
	dsn := "postgres://postgres:12345678@localhost:5432/database?sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	natsConn, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}
	defer natsConn.Close()
	log.Println("Successfully connected to NATS server")

	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	repo := postgres.NewProductRepository(db, rdb)

	sub := subscriber.NewNatsSubscriber(natsConn, repo)
	if err := sub.Subscribe(); err != nil {
		log.Fatalf("Failed to subscribe to order.created: %v", err)
	}

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
