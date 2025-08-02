package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mnadev/stok/internal/database"
	"github.com/mnadev/stok/internal/models"
	"github.com/mnadev/stok/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestHandler(t *testing.T) (*gin.Engine, *database.Service, func()) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create test database
	db, cleanup := database.SetupTestDBLegacy(t)

	// Create service
	service := database.NewService(db)

	// Create router
	router := gin.New()

	// Create handler
	handler := NewAuthHandler(db)

	// Setup routes
	router.POST("/register", handler.Register)
	router.POST("/login", handler.Login)

	return router, service, cleanup
}

func TestRegister(t *testing.T) {
	router, service, cleanup := setupTestHandler(t)
	defer cleanup()

	t.Run("Successful Registration", func(t *testing.T) {
		// Create organization and account first
		org := &models.Organization{
			Name:        "Test Corp",
			Description: "Test organization",
		}
		err := service.CreateOrganization(org)
		require.NoError(t, err)

		account := &models.Account{
			OrganizationID: &org.ID,
			Name:           "Test Shop",
			Status:         "active",
		}
		err = service.CreateAccount(account)
		require.NoError(t, err)

		// Create a test user to act as the inviter
		inviter := &models.User{
			Email:     "admin@test.com",
			Password:  "hashedpassword",
			AccountID: account.ID,
			Role:      "admin",
		}
		err = service.CreateUser(inviter)
		require.NoError(t, err)

		// Create an invitation for the test user
		invitation := &models.AccountInvitation{
			AccountID: account.ID,
			Email:     "test@example.com",
			InvitedBy: inviter.ID,
			Status:    models.AccountInvitationStatusPending,
			ExpiresAt: time.Now().AddDate(0, 0, 7), // 7 days from now
		}
		err = service.CreateInvitation(invitation)
		require.NoError(t, err)

		// Test registration
		registerData := map[string]interface{}{
			"email":    "test@example.com",
			"password": "password123",
		}

		jsonData, _ := json.Marshal(registerData)
		req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Contains(t, response, "message")
		assert.Equal(t, "User registered successfully", response["message"])
	})

	t.Run("Registration without Invitation", func(t *testing.T) {
		registerData := map[string]interface{}{
			"email":    "test2@example.com",
			"password": "password123",
		}

		jsonData, _ := json.Marshal(registerData)
		req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Contains(t, response, "error")
		assert.Contains(t, response["error"], "No invitation found")
	})

	t.Run("Registration with Invalid Email", func(t *testing.T) {
		registerData := map[string]interface{}{
			"email":    "invalid-email",
			"password": "password123",
		}

		jsonData, _ := json.Marshal(registerData)
		req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Contains(t, response, "error")
		assert.Contains(t, response["error"], "Invalid request body")
	})

	t.Run("Registration with Weak Password", func(t *testing.T) {
		registerData := map[string]interface{}{
			"email":    "test3@example.com",
			"password": "123", // Too short
		}

		jsonData, _ := json.Marshal(registerData)
		req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Contains(t, response, "error")
		assert.Contains(t, response["error"], "Invalid request body")
	})
}

func TestLogin(t *testing.T) {
	router, service, cleanup := setupTestHandler(t)
	defer cleanup()

	t.Run("Successful Login", func(t *testing.T) {
		// Create organization and account first
		org := &models.Organization{
			Name:        "Login Test Corp",
			Description: "Test organization",
		}
		err := service.CreateOrganization(org)
		require.NoError(t, err)

		account := &models.Account{
			OrganizationID: &org.ID,
			Name:           "Login Test Shop",
			Status:         "active",
		}
		err = service.CreateAccount(account)
		require.NoError(t, err)

		// Hash password before creating user
		hashedPassword, err := utils.HashPassword("password123")
		require.NoError(t, err)

		// Create user with hashed password
		user := &models.User{
			Email:     "login@example.com",
			Password:  hashedPassword,
			AccountID: account.ID,
		}
		err = service.CreateUser(user)
		require.NoError(t, err)

		// Test login
		loginData := map[string]interface{}{
			"email":    "login@example.com",
			"password": "password123",
		}

		jsonData, _ := json.Marshal(loginData)
		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Contains(t, response, "token")
		assert.Contains(t, response, "user")
	})

	t.Run("Login with Invalid Credentials", func(t *testing.T) {
		loginData := map[string]interface{}{
			"email":    "nonexistent@example.com",
			"password": "wrongpassword",
		}

		jsonData, _ := json.Marshal(loginData)
		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Contains(t, response, "error")
		assert.Contains(t, response["error"], "Invalid credentials")
	})

	t.Run("Login with Wrong Password", func(t *testing.T) {
		// Create organization and account first
		org := &models.Organization{
			Name:        "Login Test Corp 2",
			Description: "Test organization",
		}
		err := service.CreateOrganization(org)
		require.NoError(t, err)

		account := &models.Account{
			OrganizationID: &org.ID,
			Name:           "Login Test Shop 2",
			Status:         "active",
		}
		err = service.CreateAccount(account)
		require.NoError(t, err)

		// Hash password before creating user
		hashedPassword, err := utils.HashPassword("password123")
		require.NoError(t, err)

		// Create user with hashed password
		user := &models.User{
			Email:     "login2@example.com",
			Password:  hashedPassword,
			AccountID: account.ID,
		}
		err = service.CreateUser(user)
		require.NoError(t, err)

		// Test login with wrong password
		loginData := map[string]interface{}{
			"email":    "login2@example.com",
			"password": "wrongpassword",
		}

		jsonData, _ := json.Marshal(loginData)
		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Contains(t, response, "error")
		assert.Contains(t, response["error"], "Invalid credentials")
	})
}
