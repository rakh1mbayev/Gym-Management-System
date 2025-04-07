package http

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
	"inventory_service/internal/domain"
	"inventory_service/internal/usecase"
	"net/http"
	"strconv"
	"time"
)

type Handler struct {
	UC *usecase.ProductUsecase
}

func (h *Handler) CreateProduct(c *gin.Context) {
	var p domain.Product
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := h.UC.Create(ctx, &p)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusCreated)
}

func (h *Handler) GetProduct(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	p, err := h.UC.GetByID(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, p)
}

func (h *Handler) ListProducts(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	products, err := h.UC.List(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, products)
}

func (h *Handler) UpdateProduct(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var p domain.Product
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	p.ID = id
	err := h.UC.Update(ctx, &p)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

func (h *Handler) DeleteProduct(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := h.UC.Delete(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
