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
	db, err := sql.Open("postgres", "postgres://postgres:Bogdan20041@localhost:5432/inventory_db?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Connected to PostgreSQL successfully!")

	productRepo := postgres.NewProductRepository(db)
	productUC := usecase.NewProductUsecase(productRepo)

	router := gin.Default()
	http.RegisterProductRoutes(router, productUC)

	router.GET("/", func(c *gin.Context) {
		c.File("./frontend/products.html")
	})

	router.StaticFile("/products.html", "./frontend/products.html")

	router.Static("/frontend", "./frontend")
	router.Static("/css", "./frontend/css")
	router.Static("/js", "./frontend/js")

	router.Run(":8080")
}
