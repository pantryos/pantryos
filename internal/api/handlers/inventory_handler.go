// Package handlers provides HTTP request handlers for the application's API endpoints.
// This package includes handlers for inventory management, menu items, deliveries,
// and inventory snapshots. All handlers require authentication and operate within
// the context of user accounts and organizations.
package handlers

import (
	helpers "github.com/mnadev/pantryos/internal/api/helper"
	"net/http"
	"strconv"

	"github.com/mnadev/pantryos/internal/database"
	"github.com/mnadev/pantryos/internal/models"

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
	userIDInterface, exists := c.Get("userID")
	if !exists {
		errDetails := helpers.APIError{
			Code:    "UNAUTHORIZED",
			Details: "User ID not found in request context. Ensure token is valid.",
		}
		helpers.Error(c.Writer, http.StatusUnauthorized, "Authentication token is missing or invalid.", errDetails)
		return
	}

	userID, ok := userIDInterface.(int)
	if !ok {
		errDetails := helpers.APIError{
			Code:    "INTERNAL_SERVER_ERROR",
			Details: "User ID in context is not of a valid type.",
		}
		helpers.Error(c.Writer, http.StatusInternalServerError, "An internal server error occurred.", errDetails)
		return
	}

	// Retrieve user details from the database.
	user, err := h.service.GetUser(userID)
	if err != nil {
		errDetails := helpers.APIError{
			Code:    "USER_NOT_FOUND",
			Details: err.Error(),
		}
		helpers.Error(c.Writer, http.StatusNotFound, "User associated with token not found.", errDetails)
		return
	}

	// Get current stock levels from the latest snapshot.
	itemsWithStock, err := h.service.GetInventoryItemsWithCurrentStock(user.AccountID)
	if err != nil {
		errDetails := helpers.APIError{
			Code:    "DB_FETCH_FAILED",
			Details: err.Error(),
		}
		helpers.Error(c.Writer, http.StatusInternalServerError, "Failed to retrieve inventory data.", errDetails)
		return
	}

	helpers.Success(c.Writer, http.StatusOK, "Inventory items retrieved successfully.", itemsWithStock)
}

// GetLowStockItems retrieves inventory items that are currently low on stock.
// This endpoint requires authentication and returns items scoped to the user's account
// that have current stock levels below their minimum stock levels.
//
// Authentication: Required (JWT token in Authorization header)
// Authorization: User must be authenticated and have access to the account
//
// Response:
//
//	All responses are wrapped in the standard APIResponse structure.
//	- Success: { "success": true, "message": "...", "data": [...] }
//	- Error:   { "success": false, "message": "...", "error": { "code": "...", "details": "..." } }
//
// Status Codes:
//   - 200 OK: Low stock items retrieved successfully. The 'data' field contains a list of items.
//   - 401 Unauthorized: User not authenticated.
//   - 404 Not Found: User associated with token not found in the database.
//   - 500 Internal Server Error: Database or other service error.
func (h *InventoryHandler) GetLowStockItems(c *gin.Context) {
	// Extract user ID from context (set by AuthMiddleware)
	userIDInterface, exists := c.Get("userID")
	if !exists {
		errDetails := helpers.APIError{Code: "UNAUTHORIZED", Details: "User ID not found in request context."}
		helpers.Error(c.Writer, http.StatusUnauthorized, "User not authenticated.", errDetails)
		return
	}
	userID := userIDInterface.(int)

	// Retrieve user details from database
	user, err := h.service.GetUser(userID)
	if err != nil {
		errDetails := helpers.APIError{Code: "USER_NOT_FOUND", Details: err.Error()}
		helpers.Error(c.Writer, http.StatusNotFound, "User not found.", errDetails)
		return
	}

	// Get all inventory items with their current stock levels
	itemsWithStock, err := h.service.GetInventoryItemsWithCurrentStock(user.AccountID)
	if err != nil {
		errDetails := helpers.APIError{Code: "DB_FETCH_FAILED", Details: err.Error()}
		helpers.Error(c.Writer, http.StatusInternalServerError, "Failed to fetch inventory data.", errDetails)
		return
	}

	// Filter for low stock items in memory
	var lowStockItems []database.InventoryItemWithStock
	for _, item := range itemsWithStock {
		if item.CurrentStock < item.MinStockLevel {
			lowStockItems = append(lowStockItems, item)
		}
	}

	// Return a 200 OK response with the list of low stock items.
	helpers.Success(c.Writer, http.StatusOK, "Low stock items retrieved successfully.", lowStockItems)
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
//
//	All responses are wrapped in the standard APIResponse structure.
//	- Success: { "success": true, "message": "...", "data": ... }
//	- Error:   { "success": false, "message": "...", "error": { "code": "...", "details": "..." } }
//
// Status Codes:
//   - 201 Created: Inventory item created successfully. The 'data' field contains the new item.
//   - 400 Bad Request: Invalid request body or validation error.
//   - 401 Unauthorized: User not authenticated.
//   - 404 Not Found: User associated with token not found in the database.
//   - 500 Internal Server Error: Database or other service error.
func (h *InventoryHandler) CreateInventoryItem(c *gin.Context) {
	// Extract user ID from context (set by AuthMiddleware)
	userIDInterface, exists := c.Get("userID")
	if !exists {
		errDetails := helpers.APIError{Code: "UNAUTHORIZED", Details: "User ID not found in request context."}
		helpers.Error(c.Writer, http.StatusUnauthorized, "User not authenticated.", errDetails)
		return
	}
	userID := userIDInterface.(int)

	// Retrieve user details from database to get the AccountID
	user, err := h.service.GetUser(userID)
	if err != nil {
		errDetails := helpers.APIError{Code: "USER_NOT_FOUND", Details: err.Error()}
		helpers.Error(c.Writer, http.StatusNotFound, "User not found.", errDetails)
		return
	}

	// Parse and validate the JSON request body
	var item models.InventoryItem
	if err := c.ShouldBindJSON(&item); err != nil {
		errDetails := helpers.APIError{Code: "INVALID_INPUT", Details: err.Error()}
		helpers.Error(c.Writer, http.StatusBadRequest, "Invalid request body.", errDetails)
		return
	}

	// Set account ID from authenticated user to ensure proper scoping
	item.AccountID = user.AccountID

	// Create the inventory item in the database
	err = h.service.CreateInventoryItem(&item)
	if err != nil {
		// Consider checking for specific database errors, like a unique constraint violation,
		// which might warrant a 409 Conflict status code.
		errDetails := helpers.APIError{Code: "DB_INSERT_FAILED", Details: err.Error()}
		helpers.Error(c.Writer, http.StatusInternalServerError, "Failed to create inventory item.", errDetails)
		return
	}

	// Return a 201 Created response with the newly created item in the data field.
	helpers.Success(c.Writer, http.StatusCreated, "Inventory item created successfully.", item)
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
//
//	All responses are wrapped in the standard APIResponse structure.
//	- Success: { "success": true, "message": "...", "data": {...} }
//	- Error:   { "success": false, "message": "...", "error": { "code": "...", "details": "..." } }
//
// Status Codes:
//   - 200 OK: Inventory item retrieved successfully. The 'data' field contains the item object.
//   - 400 Bad Request: Invalid item ID format in URL.
//   - 401 Unauthorized: User not authenticated.
//   - 403 Forbidden: The requested item does not belong to the user's account.
//   - 404 Not Found: The user or the item could not be found.
func (h *InventoryHandler) GetInventoryItem(c *gin.Context) {
	// Extract user ID from context (set by AuthMiddleware)
	userIDInterface, exists := c.Get("userID")
	if !exists {
		errDetails := helpers.APIError{Code: "UNAUTHORIZED", Details: "User ID not found in request context."}
		helpers.Error(c.Writer, http.StatusUnauthorized, "User not authenticated.", errDetails)
		return
	}
	userID := userIDInterface.(int)

	// Retrieve user details from database
	user, err := h.service.GetUser(userID)
	if err != nil {
		errDetails := helpers.APIError{Code: "USER_NOT_FOUND", Details: err.Error()}
		helpers.Error(c.Writer, http.StatusNotFound, "User not found.", errDetails)
		return
	}

	// Parse and validate the item ID from the URL parameter
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		errDetails := helpers.APIError{Code: "INVALID_INPUT", Details: "Item ID must be a valid integer."}
		helpers.Error(c.Writer, http.StatusBadRequest, "Invalid Item ID.", errDetails)
		return
	}

	// Retrieve the inventory item from the database
	item, err := h.service.GetInventoryItem(id)
	if err != nil {
		errDetails := helpers.APIError{Code: "ITEM_NOT_FOUND", Details: err.Error()}
		helpers.Error(c.Writer, http.StatusNotFound, "Inventory item not found.", errDetails)
		return
	}

	// Authorization check: Ensure the item belongs to the user's account
	if item.AccountID != user.AccountID {
		errDetails := helpers.APIError{Code: "FORBIDDEN", Details: "You do not have permission to access this item."}
		helpers.Error(c.Writer, http.StatusForbidden, "Access denied.", errDetails)
		return
	}

	// Return a 200 OK response with the item object in the data field.
	helpers.Success(c.Writer, http.StatusOK, "Inventory item retrieved successfully.", item)
}

// UpdateInventoryItem godoc
// @Summary      Update inventory item
// @Description  Update an existing inventory item by its ID. The user must be authenticated and own the item.
// @Tags         inventory
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path      int                     true  "Inventory Item ID"
// @Param        item  body      models.InventoryItem    true  "Updated inventory item details"
// @Success      200   {object}  helpers.APIResponse{data=models.InventoryItem}  "Inventory item updated successfully"
// @Failure      400   {object}  helpers.APIResponse                           "Invalid request body or item ID"
// @Failure      401   {object}  helpers.APIResponse                           "User not authenticated"
// @Failure      403   {object}  helpers.APIResponse                           "Access denied"
// @Failure      404   {object}  helpers.APIResponse                           "Item or user not found"
// @Failure      500   {object}  helpers.APIResponse                           "Internal server error"
// @Router       /api/v1/inventory/items/{id} [put]
func (h *InventoryHandler) UpdateInventoryItem(c *gin.Context) {
	// Extract user ID from context (set by AuthMiddleware)
	userIDInterface, exists := c.Get("userID")
	if !exists {
		errDetails := helpers.APIError{Code: "UNAUTHORIZED", Details: "User ID not found in request context."}
		helpers.Error(c.Writer, http.StatusUnauthorized, "User not authenticated.", errDetails)
		return
	}
	userID := userIDInterface.(int)

	// Retrieve user details from database
	user, err := h.service.GetUser(userID)
	if err != nil {
		errDetails := helpers.APIError{Code: "USER_NOT_FOUND", Details: err.Error()}
		helpers.Error(c.Writer, http.StatusNotFound, "User not found.", errDetails)
		return
	}

	// Parse and validate the item ID from the URL parameter
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		errDetails := helpers.APIError{Code: "INVALID_INPUT", Details: "Item ID must be a valid integer."}
		helpers.Error(c.Writer, http.StatusBadRequest, "Invalid Item ID.", errDetails)
		return
	}

	// Get existing item to verify ownership
	existingItem, err := h.service.GetInventoryItem(id)
	if err != nil {
		errDetails := helpers.APIError{Code: "ITEM_NOT_FOUND", Details: err.Error()}
		helpers.Error(c.Writer, http.StatusNotFound, "Inventory item not found.", errDetails)
		return
	}

	// Authorization check: Ensure the item belongs to the user's account
	if existingItem.AccountID != user.AccountID {
		errDetails := helpers.APIError{Code: "FORBIDDEN", Details: "You do not have permission to modify this item."}
		helpers.Error(c.Writer, http.StatusForbidden, "Access denied.", errDetails)
		return
	}

	// Parse and validate the JSON request body with the updates
	var item models.InventoryItem
	if err := c.ShouldBindJSON(&item); err != nil {
		errDetails := helpers.APIError{Code: "INVALID_INPUT", Details: err.Error()}
		helpers.Error(c.Writer, http.StatusBadRequest, "Invalid request body.", errDetails)
		return
	}

	// Preserve the original ID and AccountID to prevent them from being changed.
	item.ID = id
	item.AccountID = user.AccountID

	// Update the inventory item in the database
	err = h.service.UpdateInventoryItem(&item)
	if err != nil {
		errDetails := helpers.APIError{Code: "DB_UPDATE_FAILED", Details: err.Error()}
		helpers.Error(c.Writer, http.StatusInternalServerError, "Failed to update inventory item.", errDetails)
		return
	}

	// Return a 200 OK response with the updated item object.
	helpers.Success(c.Writer, http.StatusOK, "Inventory item updated successfully.", item)
}

// DeleteInventoryItem godoc
// @Summary      Delete inventory item
// @Description  Delete an inventory item by its ID. The user must be authenticated and own the item.
// @Tags         inventory
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Inventory Item ID"
// @Success      200  {object}  helpers.APIResponse "Item deleted successfully"
// @Failure      400  {object}  helpers.APIResponse "Invalid item ID"
// @Failure      401  {object}  helpers.APIResponse "User not authenticated"
// @Failure      403  {object}  helpers.APIResponse "Access denied"
// @Failure      404  {object}  helpers.APIResponse "Item or user not found"
// @Failure      500  {object}  helpers.APIResponse "Internal server error"
// @Router       /api/v1/inventory/items/{id} [delete]
func (h *InventoryHandler) DeleteInventoryItem(c *gin.Context) {
	// Extract user ID from context (set by AuthMiddleware)
	userIDInterface, exists := c.Get("userID")
	if !exists {
		errDetails := helpers.APIError{Code: "UNAUTHORIZED", Details: "User ID not found in request context."}
		helpers.Error(c.Writer, http.StatusUnauthorized, "User not authenticated.", errDetails)
		return
	}
	userID := userIDInterface.(int)

	// Retrieve user details from database
	user, err := h.service.GetUser(userID)
	if err != nil {
		errDetails := helpers.APIError{Code: "USER_NOT_FOUND", Details: err.Error()}
		helpers.Error(c.Writer, http.StatusNotFound, "User not found.", errDetails)
		return
	}

	// Parse and validate the item ID from the URL parameter
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		errDetails := helpers.APIError{Code: "INVALID_INPUT", Details: "Item ID must be a valid integer."}
		helpers.Error(c.Writer, http.StatusBadRequest, "Invalid Item ID.", errDetails)
		return
	}

	// Get existing item to verify ownership
	existingItem, err := h.service.GetInventoryItem(id)
	if err != nil {
		errDetails := helpers.APIError{Code: "ITEM_NOT_FOUND", Details: err.Error()}
		helpers.Error(c.Writer, http.StatusNotFound, "Inventory item not found.", errDetails)
		return
	}

	// Authorization check: Ensure the item belongs to the user's account
	if existingItem.AccountID != user.AccountID {
		errDetails := helpers.APIError{Code: "FORBIDDEN", Details: "You do not have permission to delete this item."}
		helpers.Error(c.Writer, http.StatusForbidden, "Access denied.", errDetails)
		return
	}

	// Attempt to delete the inventory item
	err = h.service.DeleteInventoryItem(id)
	if err != nil {
		// The service layer might return a specific error if the item is in use
		// (e.g., part of a recipe), which could be handled here to return a 409 Conflict.
		errDetails := helpers.APIError{Code: "DB_DELETE_FAILED", Details: err.Error()}
		helpers.Error(c.Writer, http.StatusInternalServerError, "Failed to delete inventory item.", errDetails)
		return
	}

	// Return a 200 OK response with a success message.
	helpers.Success(c.Writer, http.StatusOK, "Inventory item deleted successfully.", nil)
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
// @Summary      Log a new delivery
// @Description  Log an incoming inventory delivery to update stock levels.
// @Tags         deliveries
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        delivery  body      models.Delivery               true  "Delivery details"
// @Success      201       {object}  helpers.APIResponse{data=models.Delivery}  "Delivery logged successfully"
// @Failure      400       {object}  helpers.APIResponse                           "Invalid request body"
// @Failure      401       {object}  helpers.APIResponse                           "User not authenticated"
// @Failure      404       {object}  helpers.APIResponse                           "User not found"
// @Failure      500       {object}  helpers.APIResponse                           "Internal server error"
// @Router       /api/v1/deliveries [post]
func (h *InventoryHandler) LogDelivery(c *gin.Context) {
	// Extract user ID from context (set by AuthMiddleware)
	userIDInterface, exists := c.Get("userID")
	if !exists {
		errDetails := helpers.APIError{Code: "UNAUTHORIZED", Details: "User ID not found in request context."}
		helpers.Error(c.Writer, http.StatusUnauthorized, "User not authenticated.", errDetails)
		return
	}
	userID := userIDInterface.(int)

	// Retrieve user details from database
	user, err := h.service.GetUser(userID)
	if err != nil {
		errDetails := helpers.APIError{Code: "USER_NOT_FOUND", Details: err.Error()}
		helpers.Error(c.Writer, http.StatusNotFound, "User not found.", errDetails)
		return
	}

	// Parse and validate the JSON request body
	var delivery models.Delivery
	if err := c.ShouldBindJSON(&delivery); err != nil {
		errDetails := helpers.APIError{Code: "INVALID_INPUT", Details: err.Error()}
		helpers.Error(c.Writer, http.StatusBadRequest, "Invalid request body.", errDetails)
		return
	}

	// Set account ID from the authenticated user to ensure proper scoping
	delivery.AccountID = user.AccountID

	// Create the delivery record in the database
	err = h.service.CreateDelivery(&delivery)
	if err != nil {
		// The service layer should validate if the inventory_item_id exists.
		// A more specific error could be returned here (e.g., 400 Bad Request).
		errDetails := helpers.APIError{Code: "DB_INSERT_FAILED", Details: err.Error()}
		helpers.Error(c.Writer, http.StatusInternalServerError, "Failed to log delivery.", errDetails)
		return
	}

	// Return a 201 Created response with the new delivery object.
	helpers.Success(c.Writer, http.StatusCreated, "Delivery logged successfully.", delivery)
}

// GetDeliveries godoc
// @Summary      Get all deliveries
// @Description  Retrieve all deliveries for the authenticated user's account.
// @Tags         deliveries
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  helpers.APIResponse{data=[]models.Delivery}  "Successfully retrieved list of deliveries"
// @Failure      401  {object}  helpers.APIResponse                            "Error: User not authenticated"
// @Failure      404  {object}  helpers.APIResponse                            "Error: User not found"
// @Failure      500  {object}  helpers.APIResponse                            "Error: Internal server error"
// @Router       /api/v1/deliveries [get]
func (h *InventoryHandler) GetDeliveries(c *gin.Context) {
	// Extract user ID from context (set by AuthMiddleware)
	userIDInterface, exists := c.Get("userID")
	if !exists {
		errDetails := helpers.APIError{Code: "UNAUTHORIZED", Details: "User ID not found in request context."}
		helpers.Error(c.Writer, http.StatusUnauthorized, "User not authenticated.", errDetails)
		return
	}
	userID := userIDInterface.(int)

	// Retrieve user details from the database
	user, err := h.service.GetUser(userID)
	if err != nil {
		errDetails := helpers.APIError{Code: "USER_NOT_FOUND", Details: err.Error()}
		helpers.Error(c.Writer, http.StatusNotFound, "User not found.", errDetails)
		return
	}

	// Get all deliveries for the user's account
	deliveries, err := h.service.GetDeliveriesByAccount(user.AccountID)
	if err != nil {
		errDetails := helpers.APIError{Code: "DB_FETCH_FAILED", Details: err.Error()}
		helpers.Error(c.Writer, http.StatusInternalServerError, "Failed to fetch deliveries.", errDetails)
		return
	}

	// Return a 200 OK response with the list of deliveries in the data field.
	helpers.Success(c.Writer, http.StatusOK, "Deliveries retrieved successfully.", deliveries)
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
