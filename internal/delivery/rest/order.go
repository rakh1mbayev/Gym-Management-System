package rest

import (
	"Gym-Management-System/internal/grpc"
	"Gym-Management-System/pkg/proto/orderpb"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	client grpc.OrderGRPCClient
}

func NewOrderHandler(client grpc.OrderGRPCClient) *OrderHandler {
	return &OrderHandler{client: client}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req orderpb.OrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Call the gRPC service to create the order
	res, err := h.client.CreateOrder(c, &req)
	if err != nil {
		log.Printf("Error creating order: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"order_id": res.OrderId})
}

func (h *OrderHandler) GetOrder(c *gin.Context) {
	orderID := c.Param("id")

	// Call the gRPC service to get the order details
	order, err := h.client.GetOrder(c, &orderpb.GetOrderRequest{OrderId: orderID})
	if err != nil {
		log.Printf("Error getting order: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get order"})
		return
	}

	c.JSON(http.StatusOK, order)
}

func (h *OrderHandler) UpdateOrderStatus(c *gin.Context) {
	orderID := c.Param("id") // Get orderID from the URL parameter

	var req orderpb.UpdateOrderStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Ensure that the orderID is set in the request
	req.OrderId = orderID // Set the orderID in the request

	// Call the gRPC service to update the order status
	_, err := h.client.UpdateOrderStatus(c, &req)
	if err != nil {
		log.Printf("Error updating order status: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Order status updated successfully"})
}

func (h *OrderHandler) ListOrders(c *gin.Context) {
	userID := c.DefaultQuery("user_id", "")

	// Call the gRPC service to list orders
	orders, err := h.client.ListOrders(c, &orderpb.OrderListRequest{UserId: userID})
	if err != nil {
		log.Printf("Error listing orders: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list orders"})
		return
	}

	c.JSON(http.StatusOK, orders)
}
