package main

import (
	"auth_service/internal/delivery/http"
	"auth_service/internal/repository/postgres"
	"auth_service/internal/usecase"
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	db, err := sql.Open("postgres", "postgres://postgres:12345678@localhost:5432/auth_db?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	userRepo := postgres.NewUserRepository(db)
	authUC := usecase.NewAuthUsecase(userRepo)

	router := gin.Default()
	http.RegisterAuthRoutes(router, authUC)

	router.Run(":8083")
}
