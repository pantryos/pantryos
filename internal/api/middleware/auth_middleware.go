// Package middleware provides HTTP middleware functions for the application.
// This package includes authentication and authorization middleware that can be
// used with the Gin web framework to secure API endpoints.
package middleware

import (
	"fmt"
	"net/http"

	"github.com/mnadev/stok/internal/auth"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware creates a Gin middleware function that validates JWT tokens
// and extracts user information from the Authorization header.
//
// This middleware:
//   - Extracts the JWT token from the Authorization header
//   - Validates the token signature and expiration
//   - Sets the user ID in the Gin context for downstream handlers
//   - Aborts the request with 401 Unauthorized if authentication fails
//
// Usage:
//
//	router.Use(AuthMiddleware())
//	router.GET("/protected", handler)
//
// Returns:
//   - gin.HandlerFunc: A middleware function that can be used with Gin
//
// Security notes:
//   - Validates JWT tokens using HMAC-SHA256
//   - Checks token expiration automatically
//   - Provides clear error messages for debugging
//   - Sets user ID in context for authorization checks
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract JWT token from Authorization header
		tokenString, err := auth.ExtractToken(c.Request)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		// Validate the JWT token and extract claims
		claims, err := auth.ValidateToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		// Set user ID in context for downstream handlers to access
		c.Set("userID", claims.UserID)

		// Continue to the next middleware or handler
		// Note: Account ID is retrieved in handlers to avoid database calls in middleware
		c.Next()
	}
}

// AccountMiddleware creates a Gin middleware function that validates account access
// by extracting and validating account IDs from request parameters.
//
// This middleware:
//   - Ensures the user is authenticated (userID exists in context)
//   - Extracts account ID from URL parameters or query parameters
//   - Validates the account ID format
//   - Sets the account ID in the Gin context for downstream handlers
//   - Aborts the request with appropriate error codes if validation fails
//
// The middleware looks for account ID in the following order:
//  1. URL parameter named "accountID"
//  2. Query parameter named "account_id"
//
// Usage:
//
//	router.Use(AccountMiddleware())
//	router.GET("/accounts/:accountID/items", handler)
//
// Returns:
//   - gin.HandlerFunc: A middleware function that can be used with Gin
//
// Security notes:
//   - Requires user authentication (userID in context)
//   - Validates account ID format to prevent injection attacks
//   - Provides clear error messages for debugging
//   - Supports both URL and query parameter extraction
func AccountMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if user is authenticated (userID should be set by AuthMiddleware)
		_, exists := c.Get("userID")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		// Try to get account ID from URL parameter first
		accountIDStr := c.Param("accountID")
		if accountIDStr == "" {
			// Fallback to query parameter if URL parameter is not available
			accountIDStr = c.Query("account_id")
		}

		// Validate that account ID is provided
		if accountIDStr == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Account ID is required"})
			return
		}

		// Parse account ID string to integer
		var accountID int
		if _, err := fmt.Sscanf(accountIDStr, "%d", &accountID); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID"})
			return
		}

		// Set account ID in context for downstream handlers to access
		c.Set("accountID", accountID)
		c.Next()
	}
}
