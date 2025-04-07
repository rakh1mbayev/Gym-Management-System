package main

import (
	"database/sql"
	"log"
	"order_serivce/internal/delivery/http"
	"order_serivce/internal/repository/postgres"
	"order_serivce/internal/usecase"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	db, err := sql.Open("postgres", "postgres://postgres:12345678@localhost:5432/order_db?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	orderRepo := postgres.NewOrderRepository(db)
	orderUC := usecase.NewOrderUsecase(orderRepo)

	router := gin.Default()
	http.RegisterOrderRoutes(router, orderUC)

	router.Run(":8082")
}
