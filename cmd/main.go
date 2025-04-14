package main

import (
	"Gym-Management-System/internal/grpc"
	"Gym-Management-System/internal/router"
	"github.com/gin-gonic/gin"
)

func main() {
	grpc_clients.InitClients()
	r := gin.Default()
	router.SetupRoutes(r)
	r.Run(":8080")
}
