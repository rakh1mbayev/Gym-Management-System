package http

import (
	"github.com/gin-gonic/gin"
	"order_serivce/internal/usecase"
)

func RegisterOrderRoutes(r *gin.Engine, uc *usecase.OrderUsecase) {
	h := &Handler{UC: uc}

	r.POST("/orders", h.CreateOrder)
	r.GET("/orders/:id", h.GetOrder)
	r.PATCH("/orders/:id", h.UpdateStatus)
	r.GET("/orders", h.ListUserOrders)
}
