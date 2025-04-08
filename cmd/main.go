package main

import (
	"Gym-Management-System/internal/router"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize Gin router
	r := gin.Default()

	// Setup routes with reverse proxy and JWT authentication
	router.SetupRoutes(r)

	// Start API Gateway on port 8080
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to run server: ", err)
	}
}
