package main

import (
	"database/sql"
	"log"
	"order-service/internal/delivery/http"
	"order-service/internal/repository/postgres"
	"order-service/internal/usecase"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	db, err := sql.Open("postgres", "postgres://user:pass@localhost:5432/order_db?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	orderRepo := postgres.NewOrderRepository(db)
	orderUC := usecase.NewOrderUsecase(orderRepo)

	router := gin.Default()
	http.RegisterOrderRoutes(router, orderUC)

	router.Run(":8081")
}
