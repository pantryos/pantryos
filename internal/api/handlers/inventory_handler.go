// Package handlers provides HTTP request handlers for the application's API endpoints.
// This package includes handlers for inventory management, menu items, deliveries,
// and inventory snapshots. All handlers require authentication and operate within
// the context of user accounts and organizations.
package handlers

import (
	"net/http"
	"strconv"

	"github.com/mnadev/stok/internal/database"
	"github.com/mnadev/stok/internal/models"

	"github.com/gin-gonic/gin"
)

// InventoryHandler handles HTTP requests related to inventory management operations.
// It provides CRUD operations for inventory items, menu items, deliveries, and snapshots.
// All operations are scoped to the authenticated user's account for security.
type InventoryHandler struct {
	// service provides access to the business logic layer for database operations
	service *database.Service
}

// NewInventoryHandler creates a new InventoryHandler instance with the provided database connection.
// This function initializes the handler with a database service that handles all
// business logic and data access operations.
//
// Parameters:
//   - db: The database connection to use for operations
//
// Returns:
//   - *InventoryHandler: A new handler instance ready to handle HTTP requests
func NewInventoryHandler(db *database.DB) *InventoryHandler {
	return &InventoryHandler{service: database.NewService(db)}
}

// Inventory Item Handlers

// GetInventoryItems retrieves all inventory items for the authenticated user's account.
// This endpoint requires authentication and returns items scoped to the user's account.
// The response includes all inventory items with their current stock levels and details.
//
// Authentication: Required (JWT token in Authorization header)
// Authorization: User must be authenticated and have access to the account
//
// Response:
//   - 200 OK: List of inventory items for the user's account
//   - 401 Unauthorized: User not authenticated
//   - 404 Not Found: User not found in database
//   - 500 Internal Server Error: Database or service error
//
// Security notes:
//   - Validates user authentication
//   - Scopes results to user's account only
//   - Returns appropriate error codes for different failure scenarios
func (h *InventoryHandler) GetInventoryItems(c *gin.Context) {
	// Extract user ID from context (set by AuthMiddleware)
	userID, _ := c.Get("userID")

	// Retrieve user details from database
	user, err := h.service.GetUser(userID.(int))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Get current stock levels from latest snapshot
	itemsWithStock, err := h.service.GetInventoryItemsWithCurrentStock(user.AccountID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch inventory items with stock levels"})
		return
	}

	c.JSON(http.StatusOK, itemsWithStock)
}

// GetLowStockItems retrieves inventory items that are currently low on stock.
// This endpoint requires authentication and returns items scoped to the user's account
// that have current stock levels below their minimum stock levels.
//
// Authentication: Required (JWT token in Authorization header)
// Authorization: User must be authenticated and have access to the account
//
// Response:
//   - 200 OK: List of low stock inventory items
//   - 401 Unauthorized: User not authenticated
//   - 404 Not Found: User not found in database
//   - 500 Internal Server Error: Database or service error
//
// Security notes:
//   - Validates user authentication
//   - Scopes results to user's account only
//   - Returns appropriate error codes for different failure scenarios
func (h *InventoryHandler) GetLowStockItems(c *gin.Context) {
	// Extract user ID from context (set by AuthMiddleware)
	userID, _ := c.Get("userID")

	// Retrieve user details from database
	user, err := h.service.GetUser(userID.(int))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Get all inventory items with current stock levels
	itemsWithStock, err := h.service.GetInventoryItemsWithCurrentStock(user.AccountID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch inventory items with stock levels"})
		return
	}

	// Filter for low stock items
	var lowStockItems []database.InventoryItemWithStock
	for _, item := range itemsWithStock {
		if item.CurrentStock < item.MinStockLevel {
			lowStockItems = append(lowStockItems, item)
		}
	}

	c.JSON(http.StatusOK, lowStockItems)
}

// CreateInventoryItem creates a new inventory item for the authenticated user's account.
// This endpoint accepts inventory item details in JSON format and creates a new
// item in the database. The account ID is automatically set from the authenticated user.
//
// Authentication: Required (JWT token in Authorization header)
// Authorization: User must be authenticated and have access to the account
//
// Request Body: JSON object with inventory item details
//   - name: Item name (string)
//   - unit: Unit of measurement (string, e.g., "kg", "liters")
//   - cost_per_unit: Cost per unit (float64)
//   - preferred_vendor: Preferred supplier (string, optional)
//   - min_stock_level: Minimum stock level for alerts (float64, optional)
//   - max_stock_level: Maximum stock level (float64, optional)
//
// Response:
//   - 201 Created: Inventory item created successfully
//   - 400 Bad Request: Invalid request body or validation error
//   - 401 Unauthorized: User not authenticated
//   - 404 Not Found: User not found in database
//   - 500 Internal Server Error: Database or service error
//
// Security notes:
//   - Validates user authentication
//   - Automatically sets account ID from authenticated user
//   - Validates JSON request body format
//   - Returns detailed error messages for debugging
func (h *InventoryHandler) CreateInventoryItem(c *gin.Context) {
	// Extract user ID from context (set by AuthMiddleware)
	userID, _ := c.Get("userID")

	// Retrieve user details from database
	user, err := h.service.GetUser(userID.(int))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Parse and validate the JSON request body
	var item models.InventoryItem
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	// Set account ID from authenticated user to ensure proper scoping
	item.AccountID = user.AccountID

	// Create the inventory item in the database
	err = h.service.CreateInventoryItem(&item)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create inventory item: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"item": item})
}

// GetInventoryItem retrieves a specific inventory item by ID.
// This endpoint requires authentication and validates that the item belongs
// to the authenticated user's account before returning it.
//
// Authentication: Required (JWT token in Authorization header)
// Authorization: User must be authenticated and have access to the account
//
// URL Parameters:
//   - id: The inventory item ID to retrieve (integer)
//
// Response:
//   - 200 OK: Inventory item details
//   - 400 Bad Request: Invalid item ID format
//   - 401 Unauthorized: User not authenticated
//   - 403 Forbidden: Item does not belong to user's account
//   - 404 Not Found: Item not found or user not found
//
// Security notes:
//   - Validates user authentication
//   - Validates item ID format
//   - Ensures item belongs to user's account (authorization)
//   - Returns appropriate error codes for different scenarios
func (h *InventoryHandler) GetInventoryItem(c *gin.Context) {
	userID, _ := c.Get("userID")
	user, err := h.service.GetUser(userID.(int))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}

	item, err := h.service.GetInventoryItem(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Inventory item not found"})
		return
	}

	// Check if item belongs to user's account
	if item.AccountID != user.AccountID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"item": item})
}

// UpdateInventoryItem godoc
// @Summary Update inventory item
// @Description Update an existing inventory item
// @Tags inventory
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Inventory item ID"
// @Param item body models.InventoryItem true "Updated inventory item details"
// @Success 200 {object} map[string]interface{} "Inventory item updated"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "User not authenticated"
// @Failure 403 {object} map[string]interface{} "Access denied"
// @Failure 404 {object} map[string]interface{} "Item not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/inventory/items/{id} [put]
func (h *InventoryHandler) UpdateInventoryItem(c *gin.Context) {
	userID, _ := c.Get("userID")
	user, err := h.service.GetUser(userID.(int))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}

	// Get existing item
	existingItem, err := h.service.GetInventoryItem(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Inventory item not found"})
		return
	}

	// Check if item belongs to user's account
	if existingItem.AccountID != user.AccountID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	var item models.InventoryItem
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	// Preserve account ID and ID
	item.ID = id
	item.AccountID = user.AccountID

	err = h.service.UpdateInventoryItem(&item)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update inventory item"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"item": item})
}

// DeleteInventoryItem godoc
// @Summary Delete inventory item
// @Description Delete an inventory item by ID
// @Tags inventory
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Inventory item ID"
// @Success 200 {object} map[string]interface{} "Item deleted successfully"
// @Failure 400 {object} map[string]interface{} "Invalid item ID"
// @Failure 401 {object} map[string]interface{} "User not authenticated"
// @Failure 403 {object} map[string]interface{} "Access denied"
// @Failure 404 {object} map[string]interface{} "Item not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/inventory/items/{id} [delete]
func (h *InventoryHandler) DeleteInventoryItem(c *gin.Context) {
	userID, _ := c.Get("userID")
	user, err := h.service.GetUser(userID.(int))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}

	// Get existing item to check ownership
	existingItem, err := h.service.GetInventoryItem(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Inventory item not found"})
		return
	}

	// Check if item belongs to user's account
	if existingItem.AccountID != user.AccountID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	err = h.service.DeleteInventoryItem(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete inventory item"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Inventory item deleted successfully"})
}

// Menu Item Handlers

// GetMenuItems godoc
// @Summary Get all menu items
// @Description Retrieve all menu items for the authenticated user's account
// @Tags menu
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "List of menu items"
// @Failure 401 {object} map[string]interface{} "User not authenticated"
// @Failure 404 {object} map[string]interface{} "User not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/menu/items [get]
func (h *InventoryHandler) GetMenuItems(c *gin.Context) {
	userID, _ := c.Get("userID")
	user, err := h.service.GetUser(userID.(int))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	items, err := h.service.GetMenuItemsByAccount(user.AccountID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch menu items"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": items})
}

// CreateMenuItem godoc
// @Summary Create a new menu item
// @Description Create a new menu item for the authenticated user's account
// @Tags menu
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param item body models.MenuItem true "Menu item details"
// @Success 201 {object} map[string]interface{} "Menu item created"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 401 {object} map[string]interface{} "User not authenticated"
// @Failure 404 {object} map[string]interface{} "User not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/menu/items [post]
func (h *InventoryHandler) CreateMenuItem(c *gin.Context) {
	userID, _ := c.Get("userID")
	user, err := h.service.GetUser(userID.(int))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var item models.MenuItem
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	// Set account ID from authenticated user
	item.AccountID = user.AccountID

	err = h.service.CreateMenuItem(&item)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create menu item: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"item": item})
}

// Delivery Handlers

// LogDelivery godoc
// @Summary Log a new delivery
// @Description Log an incoming inventory delivery
// @Tags deliveries
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param delivery body models.Delivery true "Delivery details"
// @Success 201 {object} map[string]interface{} "Delivery logged"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 401 {object} map[string]interface{} "User not authenticated"
// @Failure 404 {object} map[string]interface{} "User not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/deliveries [post]
func (h *InventoryHandler) LogDelivery(c *gin.Context) {
	userID, _ := c.Get("userID")
	user, err := h.service.GetUser(userID.(int))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var delivery models.Delivery
	if err := c.ShouldBindJSON(&delivery); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	// Set account ID from authenticated user
	delivery.AccountID = user.AccountID

	err = h.service.CreateDelivery(&delivery)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to log delivery: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"delivery": delivery})
}

// GetDeliveries godoc
// @Summary Get all deliveries
// @Description Retrieve all deliveries for the authenticated user's account
// @Tags deliveries
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "List of deliveries"
// @Failure 401 {object} map[string]interface{} "User not authenticated"
// @Failure 404 {object} map[string]interface{} "User not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/deliveries [get]
func (h *InventoryHandler) GetDeliveries(c *gin.Context) {
	userID, _ := c.Get("userID")
	user, err := h.service.GetUser(userID.(int))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	deliveries, err := h.service.GetDeliveriesByAccount(user.AccountID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch deliveries"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"deliveries": deliveries})
}

// Snapshot Handlers

// CreateInventorySnapshot godoc
// @Summary Create inventory snapshot
// @Description Create a point-in-time inventory count snapshot
// @Tags snapshots
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param snapshot body models.InventorySnapshot true "Snapshot details"
// @Success 201 {object} map[string]interface{} "Snapshot created"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 401 {object} map[string]interface{} "User not authenticated"
// @Failure 404 {object} map[string]interface{} "User not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/snapshots [post]
func (h *InventoryHandler) CreateInventorySnapshot(c *gin.Context) {
	userID, _ := c.Get("userID")
	user, err := h.service.GetUser(userID.(int))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var snapshot models.InventorySnapshot
	if err := c.ShouldBindJSON(&snapshot); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	// Set account ID from authenticated user
	snapshot.AccountID = user.AccountID

	err = h.service.CreateInventorySnapshot(&snapshot)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create inventory snapshot: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"snapshot": snapshot})
}

// GetInventorySnapshots godoc
// @Summary Get all inventory snapshots
// @Description Retrieve all inventory snapshots for the authenticated user's account
// @Tags snapshots
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "List of snapshots"
// @Failure 401 {object} map[string]interface{} "User not authenticated"
// @Failure 404 {object} map[string]interface{} "User not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/snapshots [get]
func (h *InventoryHandler) GetInventorySnapshots(c *gin.Context) {
	userID, _ := c.Get("userID")
	user, err := h.service.GetUser(userID.(int))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	snapshots, err := h.service.GetInventorySnapshotsByAccount(user.AccountID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch inventory snapshots"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"snapshots": snapshots})
}

// Vendor-based handlers

// GetDeliveriesByVendor godoc
// @Summary Get deliveries by vendor
// @Description Retrieve all deliveries from a specific vendor
// @Tags deliveries
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param vendor path string true "Vendor name"
// @Success 200 {object} map[string]interface{} "List of deliveries by vendor"
// @Failure 400 {object} map[string]interface{} "Vendor parameter required"
// @Failure 401 {object} map[string]interface{} "User not authenticated"
// @Failure 404 {object} map[string]interface{} "User not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/deliveries/vendor/{vendor} [get]
func (h *InventoryHandler) GetDeliveriesByVendor(c *gin.Context) {
	userID, _ := c.Get("userID")
	user, err := h.service.GetUser(userID.(int))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	vendor := c.Param("vendor")
	if vendor == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Vendor parameter is required"})
		return
	}

	deliveries, err := h.service.GetDeliveriesByVendor(user.AccountID, vendor)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch deliveries by vendor"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"deliveries": deliveries})
}

// GetInventoryItemsByVendor godoc
// @Summary Get inventory items by vendor
// @Description Retrieve all inventory items from a specific vendor
// @Tags inventory
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param vendor path string true "Vendor name"
// @Success 200 {object} map[string]interface{} "List of inventory items by vendor"
// @Failure 400 {object} map[string]interface{} "Vendor parameter required"
// @Failure 401 {object} map[string]interface{} "User not authenticated"
// @Failure 404 {object} map[string]interface{} "User not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/inventory/vendor/{vendor} [get]
func (h *InventoryHandler) GetInventoryItemsByVendor(c *gin.Context) {
	userID, _ := c.Get("userID")
	user, err := h.service.GetUser(userID.(int))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	vendor := c.Param("vendor")
	if vendor == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Vendor parameter is required"})
		return
	}

	items, err := h.service.GetInventoryItemsByVendor(user.AccountID, vendor)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch inventory items by vendor"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": items})
}
