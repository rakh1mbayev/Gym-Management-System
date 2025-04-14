package main

import (
	"log"
	_ "net/http"
	"os"

	_ "Gym-Management-System/internal/middleware"
	"Gym-Management-System/internal/router"
	_ "github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

func main() {
	// Initialize gRPC connection
	grpcServerURL := "localhost:8080"                          // Update to match the actual gRPC server address
	conn, err := grpc.Dial(grpcServerURL, grpc.WithInsecure()) // Use WithInsecure for non-TLS connections
	if err != nil {
		log.Fatalf("failed to connect to gRPC server: %v", err)
		os.Exit(1)
	}
	defer conn.Close()

	// Setup Gin router
	r := router.SetupRoutes(conn, "superSecret")

	// Apply middleware if needed (e.g., JWT Authentication)
	// r.Use(middleware.JWTAuthMiddleware()) // Uncomment if you want to apply JWT middleware

	// Start Gin server on port 8080 (can be updated based on your preference)
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("failed to start the server: %v", err)
	}
}
