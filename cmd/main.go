package main

import (
	"log"
	"os"

	"Gym-Management-System/internal/router"
	"google.golang.org/grpc"
)

func main() {
	// 1️⃣ Inventory gRPC (port 8081)
	invConn, err := grpc.Dial("localhost:8081", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect to Inventory gRPC: %v", err)
		os.Exit(1)
	}
	defer invConn.Close()

	// 2️⃣ Order gRPC (port 50052)
	orderConn, err := grpc.Dial("localhost:8082", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect to Order gRPC: %v", err)
		os.Exit(1)
	}
	defer orderConn.Close()

	// 3️⃣ User gRPC (port 8083)
	userConn, err := grpc.Dial("localhost:8083", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect to User gRPC: %v", err)
		os.Exit(1)
	}
	defer userConn.Close()

	// 4️⃣ Wire up the HTTP routes, passing all three conns + your JWT secret
	r := router.SetupRoutes(invConn, orderConn, userConn, "superSecret")

	// 5️⃣ Start the Gateway HTTP server on 8080
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("failed to start HTTP server: %v", err)
	}
}
