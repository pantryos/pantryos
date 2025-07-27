package handlers

import (
	"net/http"
	"strconv"

	"github.com/mnadev/stok/internal/database"
	"github.com/mnadev/stok/internal/models"

	"github.com/gin-gonic/gin"
)

type InventoryHandler struct {
	service *database.Service
}

func NewInventoryHandler(db *database.DB) *InventoryHandler {
	return &InventoryHandler{service: database.NewService(db)}
}

// Inventory Item Handlers

// GetInventoryItems godoc
// @Summary Get all inventory items
// @Description Retrieve all inventory items for the authenticated user's account
// @Tags inventory
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "List of inventory items"
// @Failure 401 {object} map[string]interface{} "User not authenticated"
// @Failure 404 {object} map[string]interface{} "User not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/inventory/items [get]
func (h *InventoryHandler) GetInventoryItems(c *gin.Context) {
	userID, _ := c.Get("userID")
	user, err := h.service.GetUser(userID.(int))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	items, err := h.service.GetInventoryItemsByAccount(user.AccountID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch inventory items"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": items})
}

// CreateInventoryItem godoc
// @Summary Create a new inventory item
// @Description Create a new inventory item for the authenticated user's account
// @Tags inventory
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param item body models.InventoryItem true "Inventory item details"
// @Success 201 {object} map[string]interface{} "Inventory item created"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 401 {object} map[string]interface{} "User not authenticated"
// @Failure 404 {object} map[string]interface{} "User not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/inventory/items [post]
func (h *InventoryHandler) CreateInventoryItem(c *gin.Context) {
	userID, _ := c.Get("userID")
	user, err := h.service.GetUser(userID.(int))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var item models.InventoryItem
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	// Set account ID from authenticated user
	item.AccountID = user.AccountID

	err = h.service.CreateInventoryItem(&item)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create inventory item: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"item": item})
}

// GetInventoryItem godoc
// @Summary Get inventory item by ID
// @Description Retrieve a specific inventory item by ID
// @Tags inventory
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Inventory item ID"
// @Success 200 {object} map[string]interface{} "Inventory item details"
// @Failure 400 {object} map[string]interface{} "Invalid item ID"
// @Failure 401 {object} map[string]interface{} "User not authenticated"
// @Failure 403 {object} map[string]interface{} "Access denied"
// @Failure 404 {object} map[string]interface{} "Item not found"
// @Router /api/v1/inventory/items/{id} [get]
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
