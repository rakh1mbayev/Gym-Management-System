package main

import (
	"database/sql"
	"github.com/nats-io/nats.go"
	"log"
	"net"

	_ "github.com/lib/pq" // Import the PostgreSQL driver
	rpc "github.com/rakh1mbayev/Gym-Management-System/order_service/internal/delivery/grpc"
	"github.com/rakh1mbayev/Gym-Management-System/order_service/internal/repository/postgres"
	"github.com/rakh1mbayev/Gym-Management-System/order_service/internal/usecase"
	nat "github.com/rakh1mbayev/Gym-Management-System/order_service/pkg/nats"
	"github.com/rakh1mbayev/Gym-Management-System/order_service/proto/orderpb"
	"google.golang.org/grpc"
)

func main() {
	db, err := sql.Open("postgres", "postgres://postgres:12345678@localhost:5432/database?sslmode=disable")
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

	publisher := nat.NewNatsPublisher(natsConn)
	repo := postgres.NewOrderRepository(db)
	uc := usecase.NewOrderUsecase(repo, publisher)
	srv := rpc.NewOrderServiceServer(uc)

	lis, err := net.Listen("tcp", ":8082")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	orderpb.RegisterOrderServiceServer(grpcServer, srv)

	log.Println("Order Service gRPC server started on port 8082")

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
