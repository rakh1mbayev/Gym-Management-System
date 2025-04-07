package http

import (
	"github.com/gin-gonic/gin"
	"inventory_service/internal/usecase"
)

func RegisterProductRoutes(r *gin.Engine, uc *usecase.ProductUsecase) {
	h := &Handler{UC: uc}

	r.POST("/products", h.CreateProduct)
	r.GET("/products/:id", h.GetProduct)
	r.GET("/products", h.ListProducts)
	r.PATCH("/products/:id", h.UpdateProduct)
	r.DELETE("/products/:id", h.DeleteProduct)
}
