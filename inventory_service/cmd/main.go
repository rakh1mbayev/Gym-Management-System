package main

import (
	"database/sql"
	"inventory_service/internal/delivery/http"
	_ "inventory_service/internal/domain"
	"inventory_service/internal/repository/postgres"
	"inventory_service/internal/usecase"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	db, err := sql.Open("postgres", "postgres://user:pass@localhost:5432/inventory_db?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	productRepo := postgres.NewProductRepository(db)
	productUC := usecase.NewProductUsecase(productRepo)

	router := gin.Default()
	http.RegisterProductRoutes(router, productUC)

	router.Run(":8080")
}
