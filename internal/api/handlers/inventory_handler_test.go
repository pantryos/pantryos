package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mnadev/stok/internal/database"
	"github.com/mnadev/stok/internal/models"
	"github.com/mnadev/stok/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupInventoryTestHandler(t *testing.T) (*gin.Engine, *database.Service, *models.User, *models.Account) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create test database
	db, cleanup := database.SetupTestDB(t)
	defer cleanup()

	// Create service
	service := database.NewService(db)

	// Create test organization and account
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

	// Create test user
	hashedPassword, err := utils.HashPassword("password123")
	require.NoError(t, err)

	user := &models.User{
		Email:     "test@example.com",
		Password:  hashedPassword,
		AccountID: account.ID,
		Role:      "user",
	}
	err = service.CreateUser(user)
	require.NoError(t, err)

	// Create router
	router := gin.New()

	// Create handler
	handler := NewInventoryHandler(db)

	// Setup routes with authentication middleware
	api := router.Group("/api/v1")
	{
		// Add authentication middleware to all routes
		api.Use(func(c *gin.Context) {
			// For testing, we'll set the userID from the request header
			userIDStr := c.GetHeader("X-Test-User-ID")
			if userIDStr != "" {
				if userID, err := strconv.Atoi(userIDStr); err == nil {
					c.Set("userID", userID)
				}
			}
			c.Next()
		})

		// Inventory items routes
		inventory := api.Group("/inventory")
		{
			inventory.GET("/items", handler.GetInventoryItems)
			inventory.POST("/items", handler.CreateInventoryItem)
			inventory.GET("/items/:id", handler.GetInventoryItem)
			inventory.PUT("/items/:id", handler.UpdateInventoryItem)
			inventory.DELETE("/items/:id", handler.DeleteInventoryItem)
			inventory.GET("/vendor/:vendor", handler.GetInventoryItemsByVendor)
		}

		// Menu items routes
		menu := api.Group("/menu")
		{
			menu.GET("/items", handler.GetMenuItems)
			menu.POST("/items", handler.CreateMenuItem)
		}

		// Deliveries routes
		deliveries := api.Group("/deliveries")
		{
			deliveries.GET("", handler.GetDeliveries)
			deliveries.POST("", handler.LogDelivery)
			deliveries.GET("/vendor/:vendor", handler.GetDeliveriesByVendor)
		}

		// Snapshots routes
		snapshots := api.Group("/snapshots")
		{
			snapshots.GET("", handler.GetInventorySnapshots)
			snapshots.POST("", handler.CreateInventorySnapshot)
		}
	}

	return router, service, user, account
}

// Helper function to create authenticated request
func createAuthenticatedRequest(method, path string, body interface{}, userID int) (*http.Request, *httptest.ResponseRecorder) {
	var jsonData []byte
	var err error

	if body != nil {
		jsonData, err = json.Marshal(body)
		if err != nil {
			panic(err)
		}
	}

	req, _ := http.NewRequest(method, path, bytes.NewBuffer(jsonData))
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	// Set test user ID for authentication
	req.Header.Set("X-Test-User-ID", strconv.Itoa(userID))

	return req, httptest.NewRecorder()
}

// Test Inventory Items

func TestGetInventoryItems(t *testing.T) {
	router, service, user, _ := setupInventoryTestHandler(t)

	t.Run("Get Empty Inventory Items", func(t *testing.T) {
		req, w := createAuthenticatedRequest("GET", "/api/v1/inventory/items", nil, user.ID)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Contains(t, response, "items")
		items := response["items"].([]interface{})
		assert.Empty(t, items)
	})

	t.Run("Get Inventory Items with Data", func(t *testing.T) {
		// Create test inventory item
		item := &models.InventoryItem{
			AccountID:       user.AccountID,
			Name:            "Coffee Beans",
			Unit:            "kg",
			CostPerUnit:     15.50,
			PreferredVendor: "Coffee Supply Co.",
			MinStockLevel:   5.0,
			MaxStockLevel:   50.0,
		}
		err := service.CreateInventoryItem(item)
		require.NoError(t, err)

		req, w := createAuthenticatedRequest("GET", "/api/v1/inventory/items", nil, user.ID)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Contains(t, response, "items")
		items := response["items"].([]interface{})
		assert.Len(t, items, 1)
	})
}

func TestCreateInventoryItem(t *testing.T) {
	router, _, user, _ := setupInventoryTestHandler(t)

	t.Run("Create Valid Inventory Item", func(t *testing.T) {
		itemData := map[string]interface{}{
			"name":             "Milk",
			"unit":             "liters",
			"cost_per_unit":    2.50,
			"preferred_vendor": "Local Dairy",
			"min_stock_level":  10.0,
			"max_stock_level":  100.0,
		}

		req, w := createAuthenticatedRequest("POST", "/api/v1/inventory/items", itemData, user.ID)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Contains(t, response, "item")
		item := response["item"].(map[string]interface{})
		assert.Equal(t, "Milk", item["name"])
		assert.Equal(t, "liters", item["unit"])
		assert.Equal(t, float64(2.50), item["cost_per_unit"])
	})

	t.Run("Create Inventory Item with Minimal Data", func(t *testing.T) {
		itemData := map[string]interface{}{
			"name": "Test Item", // Provide a valid name
			"unit": "kg",
		}

		req, w := createAuthenticatedRequest("POST", "/api/v1/inventory/items", itemData, user.ID)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Contains(t, response, "item")
		item := response["item"].(map[string]interface{})
		assert.Equal(t, "Test Item", item["name"]) // Valid name is preserved
		assert.Equal(t, "kg", item["unit"])
	})
}

func TestGetInventoryItem(t *testing.T) {
	router, service, user, _ := setupInventoryTestHandler(t)

	t.Run("Get Existing Inventory Item", func(t *testing.T) {
		// Create test inventory item
		item := &models.InventoryItem{
			AccountID:       user.AccountID,
			Name:            "Sugar",
			Unit:            "kg",
			CostPerUnit:     1.20,
			PreferredVendor: "Sweet Supplies",
			MinStockLevel:   2.0,
			MaxStockLevel:   20.0,
		}
		err := service.CreateInventoryItem(item)
		require.NoError(t, err)

		req, w := createAuthenticatedRequest("GET", "/api/v1/inventory/items/"+strconv.Itoa(item.ID), nil, user.ID)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Contains(t, response, "item")
		retrievedItem := response["item"].(map[string]interface{})
		assert.Equal(t, "Sugar", retrievedItem["name"])
	})

	t.Run("Get Non-Existent Inventory Item", func(t *testing.T) {
		req, w := createAuthenticatedRequest("GET", "/api/v1/inventory/items/999", nil, user.ID)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Contains(t, response, "error")
		assert.Contains(t, response["error"], "not found")
	})

	t.Run("Get Inventory Item with Invalid ID", func(t *testing.T) {
		req, w := createAuthenticatedRequest("GET", "/api/v1/inventory/items/invalid", nil, user.ID)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Contains(t, response, "error")
		assert.Contains(t, response["error"], "Invalid item ID")
	})
}

func TestUpdateInventoryItem(t *testing.T) {
	router, service, user, _ := setupInventoryTestHandler(t)

	t.Run("Update Existing Inventory Item", func(t *testing.T) {
		// Create test inventory item
		item := &models.InventoryItem{
			AccountID:       user.AccountID,
			Name:            "Flour",
			Unit:            "kg",
			CostPerUnit:     2.00,
			PreferredVendor: "Bakery Supplies",
			MinStockLevel:   5.0,
			MaxStockLevel:   50.0,
		}
		err := service.CreateInventoryItem(item)
		require.NoError(t, err)

		updateData := map[string]interface{}{
			"name":             "Premium Flour",
			"unit":             "kg",
			"cost_per_unit":    2.50,
			"preferred_vendor": "Premium Bakery Supplies",
			"min_stock_level":  10.0,
			"max_stock_level":  100.0,
		}

		req, w := createAuthenticatedRequest("PUT", "/api/v1/inventory/items/"+strconv.Itoa(item.ID), updateData, user.ID)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Contains(t, response, "item")
		updatedItem := response["item"].(map[string]interface{})
		assert.Equal(t, "Premium Flour", updatedItem["name"])
		assert.Equal(t, float64(2.50), updatedItem["cost_per_unit"])
	})
}

func TestDeleteInventoryItem(t *testing.T) {
	router, service, user, _ := setupInventoryTestHandler(t)

	t.Run("Delete Existing Inventory Item", func(t *testing.T) {
		// Create test inventory item
		item := &models.InventoryItem{
			AccountID:       user.AccountID,
			Name:            "Tea Leaves",
			Unit:            "kg",
			CostPerUnit:     8.00,
			PreferredVendor: "Tea Suppliers",
			MinStockLevel:   1.0,
			MaxStockLevel:   10.0,
		}
		err := service.CreateInventoryItem(item)
		require.NoError(t, err)

		req, w := createAuthenticatedRequest("DELETE", "/api/v1/inventory/items/"+strconv.Itoa(item.ID), nil, user.ID)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Contains(t, response, "message")
		assert.Contains(t, response["message"], "deleted")
	})
}

// Test Menu Items

func TestGetMenuItems(t *testing.T) {
	router, service, user, _ := setupInventoryTestHandler(t)

	t.Run("Get Empty Menu Items", func(t *testing.T) {
		req, w := createAuthenticatedRequest("GET", "/api/v1/menu/items", nil, user.ID)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Contains(t, response, "items")
		items := response["items"].([]interface{})
		assert.Empty(t, items)
	})

	t.Run("Get Menu Items with Data", func(t *testing.T) {
		// Create test menu item
		item := &models.MenuItem{
			AccountID: user.AccountID,
			Name:      "Cappuccino",
			Price:     4.50,
			Category:  "drinks",
		}
		err := service.CreateMenuItem(item)
		require.NoError(t, err)

		req, w := createAuthenticatedRequest("GET", "/api/v1/menu/items", nil, user.ID)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Contains(t, response, "items")
		items := response["items"].([]interface{})
		assert.Len(t, items, 1)
	})
}

func TestCreateMenuItem(t *testing.T) {
	router, _, user, _ := setupInventoryTestHandler(t)

	t.Run("Create Valid Menu Item", func(t *testing.T) {
		itemData := map[string]interface{}{
			"name":     "Latte",
			"price":    5.00,
			"category": "drinks",
		}

		req, w := createAuthenticatedRequest("POST", "/api/v1/menu/items", itemData, user.ID)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Contains(t, response, "item")
		item := response["item"].(map[string]interface{})
		assert.Equal(t, "Latte", item["name"])
		assert.Equal(t, float64(5.00), item["price"])
		assert.Equal(t, "drinks", item["category"])
	})
}

// Test Deliveries

func TestLogDelivery(t *testing.T) {
	router, service, user, _ := setupInventoryTestHandler(t)

	t.Run("Log Valid Delivery", func(t *testing.T) {
		// Create test inventory item first
		item := &models.InventoryItem{
			AccountID:       user.AccountID,
			Name:            "Coffee Beans",
			Unit:            "kg",
			CostPerUnit:     15.50,
			PreferredVendor: "Coffee Supply Co.",
		}
		err := service.CreateInventoryItem(item)
		require.NoError(t, err)

		deliveryData := map[string]interface{}{
			"inventory_item_id": item.ID,
			"vendor":            "Coffee Supply Co.",
			"quantity":          25.0,
			"delivery_date":     time.Now().Format(time.RFC3339),
			"cost":              387.50,
		}

		req, w := createAuthenticatedRequest("POST", "/api/v1/deliveries", deliveryData, user.ID)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Contains(t, response, "delivery")
		delivery := response["delivery"].(map[string]interface{})
		assert.Equal(t, float64(25.0), delivery["quantity"])
		assert.Equal(t, "Coffee Supply Co.", delivery["vendor"])
	})
}

func TestGetDeliveries(t *testing.T) {
	router, service, user, _ := setupInventoryTestHandler(t)

	t.Run("Get Empty Deliveries", func(t *testing.T) {
		req, w := createAuthenticatedRequest("GET", "/api/v1/deliveries", nil, user.ID)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Contains(t, response, "deliveries")
		deliveries := response["deliveries"].([]interface{})
		assert.Empty(t, deliveries)
	})

	t.Run("Get Deliveries with Data", func(t *testing.T) {
		// Create test inventory item first
		item := &models.InventoryItem{
			AccountID:       user.AccountID,
			Name:            "Milk",
			Unit:            "liters",
			CostPerUnit:     2.50,
			PreferredVendor: "Local Dairy",
		}
		err := service.CreateInventoryItem(item)
		require.NoError(t, err)

		// Create test delivery
		delivery := &models.Delivery{
			AccountID:       user.AccountID,
			InventoryItemID: item.ID,
			Vendor:          "Local Dairy",
			Quantity:        50.0,
			Cost:            125.0,
		}
		err = service.CreateDelivery(delivery)
		require.NoError(t, err)

		req, w := createAuthenticatedRequest("GET", "/api/v1/deliveries", nil, user.ID)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Contains(t, response, "deliveries")
		deliveries := response["deliveries"].([]interface{})
		assert.Len(t, deliveries, 1)
	})
}

func TestGetDeliveriesByVendor(t *testing.T) {
	router, service, user, _ := setupInventoryTestHandler(t)

	t.Run("Get Deliveries by Vendor", func(t *testing.T) {
		// Create test inventory item first
		item := &models.InventoryItem{
			AccountID:       user.AccountID,
			Name:            "Sugar",
			Unit:            "kg",
			CostPerUnit:     1.20,
			PreferredVendor: "Sweet Supplies",
		}
		err := service.CreateInventoryItem(item)
		require.NoError(t, err)

		// Create test delivery
		delivery := &models.Delivery{
			AccountID:       user.AccountID,
			InventoryItemID: item.ID,
			Vendor:          "Sweet Supplies",
			Quantity:        100.0,
			Cost:            120.0,
		}
		err = service.CreateDelivery(delivery)
		require.NoError(t, err)

		req, w := createAuthenticatedRequest("GET", "/api/v1/deliveries/vendor/Sweet%20Supplies", nil, user.ID)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Contains(t, response, "deliveries")
		deliveries := response["deliveries"].([]interface{})
		assert.Len(t, deliveries, 1)
	})
}

func TestCreateInventorySnapshot(t *testing.T) {
	router, service, user, _ := setupInventoryTestHandler(t)

	t.Run("Create Valid Inventory Snapshot", func(t *testing.T) {
		// Create test inventory item first
		item := &models.InventoryItem{
			AccountID:       user.AccountID,
			Name:            "Coffee Beans",
			Unit:            "kg",
			CostPerUnit:     15.50,
			PreferredVendor: "Coffee Supply Co.",
		}
		err := service.CreateInventoryItem(item)
		require.NoError(t, err)

		snapshotData := map[string]interface{}{
			"counts": map[string]interface{}{
				strconv.Itoa(item.ID): 25.5,
			},
		}

		req, w := createAuthenticatedRequest("POST", "/api/v1/snapshots", snapshotData, user.ID)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Contains(t, response, "snapshot")
		snapshot := response["snapshot"].(map[string]interface{})
		assert.Contains(t, snapshot, "counts")
	})
}

func TestGetInventorySnapshots(t *testing.T) {
	router, service, user, _ := setupInventoryTestHandler(t)

	t.Run("Get Empty Inventory Snapshots", func(t *testing.T) {
		req, w := createAuthenticatedRequest("GET", "/api/v1/snapshots", nil, user.ID)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Contains(t, response, "snapshots")
		snapshots := response["snapshots"].([]interface{})
		assert.Empty(t, snapshots)
	})

	t.Run("Get Inventory Snapshots with Data", func(t *testing.T) {
		// Create test inventory item first
		item := &models.InventoryItem{
			AccountID:       user.AccountID,
			Name:            "Tea Leaves",
			Unit:            "kg",
			CostPerUnit:     8.00,
			PreferredVendor: "Tea Suppliers",
		}
		err := service.CreateInventoryItem(item)
		require.NoError(t, err)

		// Create test snapshot
		counts := models.CountsMap{item.ID: 15.0}
		snapshot := &models.InventorySnapshot{
			AccountID: user.AccountID,
			Counts:    counts,
		}
		err = service.CreateInventorySnapshot(snapshot)
		require.NoError(t, err)

		req, w := createAuthenticatedRequest("GET", "/api/v1/snapshots", nil, user.ID)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Contains(t, response, "snapshots")
		snapshots := response["snapshots"].([]interface{})
		assert.Len(t, snapshots, 1)
	})
}

func TestGetInventoryItemsByVendor(t *testing.T) {
	router, service, user, _ := setupInventoryTestHandler(t)

	t.Run("Get Inventory Items by Vendor", func(t *testing.T) {
		// Create test inventory items
		item1 := &models.InventoryItem{
			AccountID:       user.AccountID,
			Name:            "Coffee Beans",
			Unit:            "kg",
			CostPerUnit:     15.50,
			PreferredVendor: "Coffee Supply Co.",
		}
		err := service.CreateInventoryItem(item1)
		require.NoError(t, err)

		item2 := &models.InventoryItem{
			AccountID:       user.AccountID,
			Name:            "Tea Leaves",
			Unit:            "kg",
			CostPerUnit:     8.00,
			PreferredVendor: "Coffee Supply Co.",
		}
		err = service.CreateInventoryItem(item2)
		require.NoError(t, err)

		req, w := createAuthenticatedRequest("GET", "/api/v1/inventory/vendor/Coffee%20Supply%20Co.", nil, user.ID)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Contains(t, response, "items")
		items := response["items"].([]interface{})
		assert.Len(t, items, 2)
	})

	t.Run("Get Inventory Items by Non-Existent Vendor", func(t *testing.T) {
		req, w := createAuthenticatedRequest("GET", "/api/v1/inventory/vendor/NonExistentVendor", nil, user.ID)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Contains(t, response, "items")
		items := response["items"].([]interface{})
		assert.Empty(t, items)
	})
}

// Test Error Cases

func TestInventoryHandlerErrors(t *testing.T) {
	router, _, user, _ := setupInventoryTestHandler(t)

	t.Run("Create Inventory Item with Missing Optional Fields", func(t *testing.T) {
		itemData := map[string]interface{}{
			"name": "Test Item",
			"unit": "kg", // Provide required unit field
			// Missing cost_per_unit - this has default value
		}

		req, w := createAuthenticatedRequest("POST", "/api/v1/inventory/items", itemData, user.ID)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Contains(t, response, "item")
		item := response["item"].(map[string]interface{})
		assert.Equal(t, "Test Item", item["name"])
		assert.Equal(t, "kg", item["unit"])                // Provided unit
		assert.Equal(t, float64(0), item["cost_per_unit"]) // Default 0
	})

	t.Run("Update Non-Existent Inventory Item", func(t *testing.T) {
		updateData := map[string]interface{}{
			"name": "Updated Item",
			"unit": "pieces",
		}

		req, w := createAuthenticatedRequest("PUT", "/api/v1/inventory/items/999", updateData, user.ID)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Contains(t, response, "error")
		assert.Contains(t, response["error"], "not found")
	})

	t.Run("Delete Non-Existent Inventory Item", func(t *testing.T) {
		req, w := createAuthenticatedRequest("DELETE", "/api/v1/inventory/items/999", nil, user.ID)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Contains(t, response, "error")
		assert.Contains(t, response["error"], "not found")
	})
}
