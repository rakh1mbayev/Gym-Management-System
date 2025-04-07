package router

import (
	"Gym-Management-System/internal/middleware"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
)

// Service endpoints (ideally loaded from config)
const (
	inventoryServiceURL = "http://localhost:8081" // Inventory Service URL
	orderServiceURL     = "http://localhost:8082" // Order Service URL
	authServiceURL      = "http://localhost:8083"
)

// SetupRoutes registers API routes and applies middleware.
func SetupRoutes(r *gin.Engine) {
	// Apply JWT middleware to protected routes
	protected := r.Group("/", middleware.JWTAuthMiddleware())

	// Routes for Inventory Service
	protected.Any("/products/*action", reverseProxy(inventoryServiceURL))
	// Routes for Order Service
	protected.Any("/orders/*action", reverseProxy(orderServiceURL))

	// Public routes for Auth endpoints (if you wish to proxy them or handle locally)
	// For example, the gateway could forward /login and /register to the Auth Service.
	r.POST("/login", reverseProxy(authServiceURL))    // Uncomment and set authServiceURL if needed.
	r.POST("/register", reverseProxy(authServiceURL)) // Uncomment and set authServiceURL if needed.
}

// reverseProxy returns a gin.HandlerFunc that proxies the request to the target service.
func reverseProxy(target string) gin.HandlerFunc {
	return func(c *gin.Context) {
		targetURL, err := url.Parse(target)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid target URL"})
			return
		}

		proxy := httputil.NewSingleHostReverseProxy(targetURL)

		// Use the wildcard parameter if available, otherwise use the original path.
		path := c.Param("action")
		if path == "" {
			path = c.Request.URL.Path
		}
		c.Request.URL.Path = path

		proxy.ServeHTTP(c.Writer, c.Request)
	}
}
