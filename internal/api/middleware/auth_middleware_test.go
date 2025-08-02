package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mnadev/pantryos/internal/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	return router
}

func TestAuthMiddleware(t *testing.T) {
	t.Run("Valid Token", func(t *testing.T) {
		router := setupTestRouter()

		// Generate a valid JWT token
		userID := 123
		token, err := auth.GenerateJWT(userID)
		require.NoError(t, err)

		// Add middleware to router
		router.Use(AuthMiddleware())

		// Add a test endpoint
		router.GET("/test", func(c *gin.Context) {
			userIDFromContext, exists := c.Get("userID")
			assert.True(t, exists)
			assert.Equal(t, userID, userIDFromContext)
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		// Create request with valid token
		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Missing Authorization Header", func(t *testing.T) {
		router := setupTestRouter()
		router.Use(AuthMiddleware())

		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		req, _ := http.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Invalid Authorization Header Format", func(t *testing.T) {
		router := setupTestRouter()
		router.Use(AuthMiddleware())

		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "InvalidFormat token123")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Invalid Token", func(t *testing.T) {
		router := setupTestRouter()
		router.Use(AuthMiddleware())

		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer invalid-token")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Empty Token", func(t *testing.T) {
		router := setupTestRouter()
		router.Use(AuthMiddleware())

		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer ")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Case Insensitive Bearer", func(t *testing.T) {
		router := setupTestRouter()

		// Generate a valid JWT token
		userID := 456
		token, err := auth.GenerateJWT(userID)
		require.NoError(t, err)

		router.Use(AuthMiddleware())
		router.GET("/test", func(c *gin.Context) {
			userIDFromContext, exists := c.Get("userID")
			assert.True(t, exists)
			assert.Equal(t, userID, userIDFromContext)
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "bearer "+token) // lowercase "bearer"
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestAccountMiddleware(t *testing.T) {
	t.Run("Valid Account ID in URL Parameter", func(t *testing.T) {
		router := setupTestRouter()

		// Set up user ID in context (simulating AuthMiddleware)
		router.Use(func(c *gin.Context) {
			c.Set("userID", 123)
			c.Next()
		})

		router.Use(AccountMiddleware())

		router.GET("/accounts/:accountID/test", func(c *gin.Context) {
			accountID, exists := c.Get("accountID")
			assert.True(t, exists)
			assert.Equal(t, 456, accountID)
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		req, _ := http.NewRequest("GET", "/accounts/456/test", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Valid Account ID in Query Parameter", func(t *testing.T) {
		router := setupTestRouter()

		// Set up user ID in context (simulating AuthMiddleware)
		router.Use(func(c *gin.Context) {
			c.Set("userID", 123)
			c.Next()
		})

		router.Use(AccountMiddleware())

		router.GET("/test", func(c *gin.Context) {
			accountID, exists := c.Get("accountID")
			assert.True(t, exists)
			assert.Equal(t, 789, accountID)
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		req, _ := http.NewRequest("GET", "/test?account_id=789", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("No User ID in Context", func(t *testing.T) {
		router := setupTestRouter()
		router.Use(AccountMiddleware())

		router.GET("/accounts/:accountID/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		req, _ := http.NewRequest("GET", "/accounts/456/test", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Missing Account ID", func(t *testing.T) {
		router := setupTestRouter()

		// Set up user ID in context (simulating AuthMiddleware)
		router.Use(func(c *gin.Context) {
			c.Set("userID", 123)
			c.Next()
		})

		router.Use(AccountMiddleware())

		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		req, _ := http.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Invalid Account ID Format", func(t *testing.T) {
		router := setupTestRouter()

		// Set up user ID in context (simulating AuthMiddleware)
		router.Use(func(c *gin.Context) {
			c.Set("userID", 123)
			c.Next()
		})

		router.Use(AccountMiddleware())

		router.GET("/accounts/:accountID/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		req, _ := http.NewRequest("GET", "/accounts/invalid/test", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Zero Account ID", func(t *testing.T) {
		router := setupTestRouter()

		// Set up user ID in context (simulating AuthMiddleware)
		router.Use(func(c *gin.Context) {
			c.Set("userID", 123)
			c.Next()
		})

		router.Use(AccountMiddleware())

		router.GET("/accounts/:accountID/test", func(c *gin.Context) {
			accountID, exists := c.Get("accountID")
			assert.True(t, exists)
			assert.Equal(t, 0, accountID)
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		req, _ := http.NewRequest("GET", "/accounts/0/test", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Large Account ID", func(t *testing.T) {
		router := setupTestRouter()

		// Set up user ID in context (simulating AuthMiddleware)
		router.Use(func(c *gin.Context) {
			c.Set("userID", 123)
			c.Next()
		})

		router.Use(AccountMiddleware())

		router.GET("/accounts/:accountID/test", func(c *gin.Context) {
			accountID, exists := c.Get("accountID")
			assert.True(t, exists)
			assert.Equal(t, 999999999, accountID)
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		req, _ := http.NewRequest("GET", "/accounts/999999999/test", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestMiddlewareChain(t *testing.T) {
	t.Run("Auth and Account Middleware Together", func(t *testing.T) {
		router := setupTestRouter()

		// Generate a valid JWT token
		userID := 123
		token, err := auth.GenerateJWT(userID)
		require.NoError(t, err)

		// Add both middlewares
		router.Use(AuthMiddleware())
		router.Use(AccountMiddleware())

		router.GET("/accounts/:accountID/test", func(c *gin.Context) {
			userIDFromContext, exists := c.Get("userID")
			assert.True(t, exists)
			assert.Equal(t, userID, userIDFromContext)

			accountID, exists := c.Get("accountID")
			assert.True(t, exists)
			assert.Equal(t, 456, accountID)

			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		req, _ := http.NewRequest("GET", "/accounts/456/test", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Auth Middleware Fails, Account Middleware Not Reached", func(t *testing.T) {
		router := setupTestRouter()

		// Add both middlewares
		router.Use(AuthMiddleware())
		router.Use(AccountMiddleware())

		router.GET("/accounts/:accountID/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		req, _ := http.NewRequest("GET", "/accounts/456/test", nil)
		// No Authorization header
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestMiddlewareErrorResponses(t *testing.T) {
	t.Run("Auth Middleware Error Response Format", func(t *testing.T) {
		router := setupTestRouter()
		router.Use(AuthMiddleware())

		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		req, _ := http.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		// The response should be JSON with an error field
		assert.Contains(t, w.Header().Get("Content-Type"), "application/json")
	})

	t.Run("Account Middleware Error Response Format", func(t *testing.T) {
		router := setupTestRouter()

		// Set up user ID in context (simulating AuthMiddleware)
		router.Use(func(c *gin.Context) {
			c.Set("userID", 123)
			c.Next()
		})

		router.Use(AccountMiddleware())

		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		req, _ := http.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		// The response should be JSON with an error field
		assert.Contains(t, w.Header().Get("Content-Type"), "application/json")
	})
}
