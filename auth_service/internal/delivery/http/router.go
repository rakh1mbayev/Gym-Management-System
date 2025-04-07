package http

import (
	"auth_service/internal/usecase"
	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(r *gin.Engine, uc *usecase.AuthUsecase) {
	h := &Handler{UC: uc}

	r.POST("/register", h.Register)
	r.POST("/login", h.Login)
}
