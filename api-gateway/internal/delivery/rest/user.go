package rest

import (
	"fmt"
	"github.com/rakh1mbayev/Gym-Management-System/user_service/proto/userpb"
	"golang.org/x/net/context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rakh1mbayev/Gym-Management-System/api-gateway/internal/grpc"
)

type UserHandler struct {
	client grpc.UserGRPCClient
}

func NewUserHandler(client grpc.UserGRPCClient) *UserHandler {
	return &UserHandler{client: client}
}

func (h *UserHandler) RegisterUser(c *gin.Context) {
	var req userpb.CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.client.RegisterUser(c, &req)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, res)
}

func (h *UserHandler) AuthenticateUser(c *gin.Context) {
	var req userpb.AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.client.AuthenticateUser(c, &req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Return only the token
	c.JSON(http.StatusOK, gin.H{"token": res.GetToken()})
}

func (h *UserHandler) GetUserProfile(c *gin.Context) {
	id := c.Param("id")

	res, err := h.client.GetUserProfile(c, &userpb.GetRequest{UserId: int64(parseID(id))})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func parseID(s string) int32 {
	var id int
	fmt.Sscan(s, &id)
	return int32(id)
}
func (h *UserHandler) ConfirmEmail(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "token is required"})
		return
	}

	req := &userpb.ConfirmEmailRequest{Token: token}
	resp, err := h.client.ConfirmEmail(context.Background(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !resp.Success {
		c.JSON(http.StatusBadRequest, gin.H{"error": resp.GetMessage()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": resp.GetMessage()})
}
