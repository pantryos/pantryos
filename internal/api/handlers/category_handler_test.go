package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mnadev/pantryos/internal/database"
	"github.com/mnadev/pantryos/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCategoryHandler_CreateCategory(t *testing.T) {
	// Setup test database
	db, cleanup := database.SetupTestDBLegacy(t)
	defer cleanup()

	service := database.NewService(db)

	// Create test account
	account := &models.Account{
		Name:         "Test Coffee Shop",
		Location:     "123 Test St",
		BusinessType: "single_location",
		Status:       "active",
	}
	err := service.CreateAccount(account)
	require.NoError(t, err)

	// Create test user
	user := &models.User{
		Email:     "test@example.com",
		Password:  "password123",
		AccountID: account.ID,
		Role:      "admin",
	}
	err = service.CreateUser(user)
	require.NoError(t, err)

	// Setup router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	categoryHandler := NewCategoryHandler(db)

	// Add auth middleware mock
	router.Use(func(c *gin.Context) {
		c.Set("userID", user.ID)
		c.Next()
	})

	// Add routes
	router.POST("/categories", categoryHandler.CreateCategory)

	// Test creating a category
	categoryData := models.Category{
		Name:        "Coffee Beans",
		Description: "Various types of coffee beans",
		Color:       "#8B4513",
		IsActive:    true,
	}

	jsonData, _ := json.Marshal(categoryData)
	req, _ := http.NewRequest("POST", "/categories", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Verify the category was created
	categories, err := service.GetCategoriesByAccount(account.ID)
	assert.NoError(t, err)
	assert.Len(t, categories, 1)
	assert.Equal(t, "Coffee Beans", categories[0].Name)
	assert.Equal(t, account.ID, categories[0].AccountID)
}

func TestCategoryHandler_GetCategories(t *testing.T) {
	// Setup test database
	db, cleanup := database.SetupTestDBLegacy(t)
	defer cleanup()

	service := database.NewService(db)

	// Create test account
	account := &models.Account{
		Name:         "Test Coffee Shop",
		Location:     "123 Test St",
		BusinessType: "single_location",
		Status:       "active",
	}
	err := service.CreateAccount(account)
	require.NoError(t, err)

	// Create test user
	user := &models.User{
		Email:     "test@example.com",
		Password:  "password123",
		AccountID: account.ID,
		Role:      "admin",
	}
	err = service.CreateUser(user)
	require.NoError(t, err)

	// Create test categories
	category1 := &models.Category{
		AccountID:   account.ID,
		Name:        "Coffee Beans",
		Description: "Various types of coffee beans",
		Color:       "#8B4513",
		IsActive:    true,
	}
	err = service.CreateCategory(category1)
	require.NoError(t, err)

	category2 := &models.Category{
		AccountID:   account.ID,
		Name:        "Dairy",
		Description: "Milk and cream products",
		Color:       "#FFFFFF",
		IsActive:    true,
	}
	err = service.CreateCategory(category2)
	require.NoError(t, err)

	// Setup router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	categoryHandler := NewCategoryHandler(db)

	// Add auth middleware mock
	router.Use(func(c *gin.Context) {
		c.Set("userID", user.ID)
		c.Next()
	})

	// Add routes
	router.GET("/categories", categoryHandler.GetCategories)

	// Test getting categories
	req, _ := http.NewRequest("GET", "/categories", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	categories := response["data"].([]interface{})
	assert.Len(t, categories, 2)
}

func TestCategoryHandler_CreateInventoryItemWithCategory(t *testing.T) {
	// Setup test database
	db, cleanup := database.SetupTestDBLegacy(t)
	defer cleanup()

	service := database.NewService(db)

	// Create test account
	account := &models.Account{
		Name:         "Test Coffee Shop",
		Location:     "123 Test St",
		BusinessType: "single_location",
		Status:       "active",
	}
	err := service.CreateAccount(account)
	require.NoError(t, err)

	// Create test user
	user := &models.User{
		Email:     "test@example.com",
		Password:  "password123",
		AccountID: account.ID,
		Role:      "admin",
	}
	err = service.CreateUser(user)
	require.NoError(t, err)

	// Create test category
	category := &models.Category{
		AccountID:   account.ID,
		Name:        "Coffee Beans",
		Description: "Various types of coffee beans",
		Color:       "#8B4513",
		IsActive:    true,
	}
	err = service.CreateCategory(category)
	require.NoError(t, err)

	// Setup router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	inventoryHandler := NewInventoryHandler(db)

	// Add auth middleware mock
	router.Use(func(c *gin.Context) {
		c.Set("userID", user.ID)
		c.Next()
	})

	// Add routes
	router.POST("/inventory/items", inventoryHandler.CreateInventoryItem)

	// Test creating an inventory item with category
	itemData := models.InventoryItem{
		Name:            "Arabica Coffee Beans",
		Unit:            "kg",
		CostPerUnit:     15.99,
		PreferredVendor: "Coffee Supply Co.",
		MinStockLevel:   10.0,
		MaxStockLevel:   50.0,
		CategoryID:      &category.ID,
	}

	jsonData, _ := json.Marshal(itemData)
	req, _ := http.NewRequest("POST", "/inventory/items", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	// Verify the item was created with the category
	items, err := service.GetInventoryItemsByAccount(account.ID)
	assert.NoError(t, err)
	assert.Len(t, items, 1)
	assert.Equal(t, "Arabica Coffee Beans", items[0].Name)
	assert.Equal(t, category.ID, *items[0].CategoryID)
}
