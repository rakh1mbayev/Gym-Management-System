package main

import (
	"database/sql"
	"log"
	"order_service/internal/delivery/http"
	"order_service/internal/repository/postgres"
	"order_service/internal/usecase"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	db, err := sql.Open("postgres", "postgres://postgres:Bogdan20041@localhost:5432/order_db?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Connected to PostgreSQL successfully!")

	orderRepo := postgres.NewOrderRepository(db)
	orderUC := usecase.NewOrderUsecase(orderRepo)

	router := gin.Default()
	http.RegisterOrderRoutes(router, orderUC)

	router.Run(":8082")
}
