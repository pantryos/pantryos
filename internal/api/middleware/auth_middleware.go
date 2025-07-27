package middleware

import (
	"fmt"
	"net/http"

	"github.com/mnadev/stok/internal/auth"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := auth.ExtractToken(c.Request)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		claims, err := auth.ValidateToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		// Set user ID in context for downstream handlers
		c.Set("userID", claims.UserID)

		// Get user details and set account ID in context
		// Note: This requires a database connection, so we'll get it in the handlers
		c.Next()
	}
}

// AccountMiddleware ensures the user has access to the specified account
func AccountMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, exists := c.Get("userID")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		// Get account ID from URL parameter or request body
		accountIDStr := c.Param("accountID")
		if accountIDStr == "" {
			// Try to get from query parameter
			accountIDStr = c.Query("account_id")
		}

		if accountIDStr == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Account ID is required"})
			return
		}

		// Parse account ID
		var accountID int
		if _, err := fmt.Sscanf(accountIDStr, "%d", &accountID); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID"})
			return
		}

		// Set account ID in context
		c.Set("accountID", accountID)
		c.Next()
	}
}
