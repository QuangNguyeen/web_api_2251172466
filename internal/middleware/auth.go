package middleware

import (
	"net/http"
	"strings"

	"library-api/internal/dto"
	"library-api/internal/utils"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates JWT token
func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, dto.APIResponse{
				Success: false,
				Message: "Authorization header required",
			})
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, dto.APIResponse{
				Success: false,
				Message: "Invalid authorization format",
			})
			c.Abort()
			return
		}

		token := parts[1]

		// Validate token
		claims, err := utils.ValidateToken(token, jwtSecret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, dto.APIResponse{
				Success: false,
				Message: "Invalid or expired token",
			})
			c.Abort()
			return
		}

		// Set user info in context
		c.Set("userID", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("isAdmin", claims.IsAdmin)

		c.Next()
	}
}

// AdminMiddleware checks if user is admin
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		isAdmin, exists := c.Get("isAdmin")
		if !exists || !isAdmin.(bool) {
			c.JSON(http.StatusForbidden, dto.APIResponse{
				Success: false,
				Message: "Admin access required",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
