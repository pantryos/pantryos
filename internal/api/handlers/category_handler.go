// Package handlers provides HTTP request handlers for the application's API endpoints.
// This package includes handlers for category management operations.
// All handlers require authentication and operate within the context of user accounts.
package handlers

import (
	"net/http"
	"strconv"

	"github.com/mnadev/stok/internal/database"
	"github.com/mnadev/stok/internal/models"

	"github.com/gin-gonic/gin"
)

// CategoryHandler handles HTTP requests related to category management operations.
// It provides CRUD operations for categories that help organize inventory items
// and menu items for better management and reporting.
type CategoryHandler struct {
	// service provides access to the business logic layer for database operations
	service *database.Service
}

// NewCategoryHandler creates a new CategoryHandler instance with the provided database connection.
// This function initializes the handler with a database service that handles all
// business logic and data access operations.
//
// Parameters:
//   - db: The database connection to use for operations
//
// Returns:
//   - *CategoryHandler: A new handler instance ready to handle HTTP requests
func NewCategoryHandler(db *database.DB) *CategoryHandler {
	return &CategoryHandler{service: database.NewService(db)}
}

// GetCategories retrieves all categories for the authenticated user's account.
// This endpoint requires authentication and returns categories scoped to the user's account.
// The response includes all categories with their details and status.
//
// Authentication: Required (JWT token in Authorization header)
// Authorization: User must be authenticated and have access to the account
//
// Response:
//   - 200 OK: List of categories for the user's account
//   - 401 Unauthorized: User not authenticated
//   - 404 Not Found: User not found in database
//   - 500 Internal Server Error: Database or service error
//
// Security notes:
//   - Validates user authentication
//   - Scopes results to user's account only
//   - Returns appropriate error codes for different failure scenarios
func (h *CategoryHandler) GetCategories(c *gin.Context) {
	// Extract user ID from context (set by AuthMiddleware)
	userID, _ := c.Get("userID")

	// Retrieve user details from database
	user, err := h.service.GetUser(userID.(int))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Get all categories for the account
	categories, err := h.service.GetCategoriesByAccount(user.AccountID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch categories"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"categories": categories})
}

// GetActiveCategories retrieves only active categories for the authenticated user's account.
// This endpoint requires authentication and returns only active categories scoped to the user's account.
// This is useful for UI dropdowns and active item categorization.
//
// Authentication: Required (JWT token in Authorization header)
// Authorization: User must be authenticated and have access to the account
//
// Response:
//   - 200 OK: List of active categories for the user's account
//   - 401 Unauthorized: User not authenticated
//   - 404 Not Found: User not found in database
//   - 500 Internal Server Error: Database or service error
//
// Security notes:
//   - Validates user authentication
//   - Scopes results to user's account only
//   - Returns appropriate error codes for different failure scenarios
func (h *CategoryHandler) GetActiveCategories(c *gin.Context) {
	// Extract user ID from context (set by AuthMiddleware)
	userID, _ := c.Get("userID")

	// Retrieve user details from database
	user, err := h.service.GetUser(userID.(int))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Get active categories for the account
	categories, err := h.service.GetActiveCategoriesByAccount(user.AccountID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch active categories"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"categories": categories})
}

// CreateCategory creates a new category for the authenticated user's account.
// This endpoint accepts category details in JSON format and creates a new
// category in the database. The account ID is automatically set from the authenticated user.
//
// Authentication: Required (JWT token in Authorization header)
// Authorization: User must be authenticated and have access to the account
//
// Request Body: JSON object with category details
//   - name: Category name (string, required)
//   - description: Category description (string, optional)
//   - color: Hex color code for UI display (string, optional, default: #6B7280)
//   - is_active: Whether the category is active (bool, optional, default: true)
//
// Response:
//   - 201 Created: Category created successfully
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
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	// Extract user ID from context (set by AuthMiddleware)
	userID, _ := c.Get("userID")

	// Retrieve user details from database
	user, err := h.service.GetUser(userID.(int))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Parse and validate the JSON request body
	var category models.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	// Set account ID from authenticated user to ensure proper scoping
	category.AccountID = user.AccountID

	// Create the category in the database
	err = h.service.CreateCategory(&category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create category: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"category": category})
}

// GetCategory retrieves a specific category by ID.
// This endpoint requires authentication and validates that the category belongs
// to the authenticated user's account before returning it.
//
// Authentication: Required (JWT token in Authorization header)
// Authorization: User must be authenticated and have access to the account
//
// URL Parameters:
//   - id: The category ID to retrieve (integer)
//
// Response:
//   - 200 OK: Category details
//   - 400 Bad Request: Invalid category ID format
//   - 401 Unauthorized: User not authenticated
//   - 403 Forbidden: Category does not belong to user's account
//   - 404 Not Found: Category not found or user not found
//
// Security notes:
//   - Validates user authentication
//   - Validates category ID format
//   - Ensures category belongs to user's account (authorization)
//   - Returns appropriate error codes for different scenarios
func (h *CategoryHandler) GetCategory(c *gin.Context) {
	userID, _ := c.Get("userID")
	user, err := h.service.GetUser(userID.(int))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	category, err := h.service.GetCategory(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	// Check if category belongs to user's account
	if category.AccountID != user.AccountID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"category": category})
}

// UpdateCategory updates an existing category.
// This endpoint requires authentication and validates that the category belongs
// to the authenticated user's account before updating it.
//
// Authentication: Required (JWT token in Authorization header)
// Authorization: User must be authenticated and have access to the account
//
// URL Parameters:
//   - id: The category ID to update (integer)
//
// Request Body: JSON object with updated category details
//   - name: Category name (string, required)
//   - description: Category description (string, optional)
//   - color: Hex color code for UI display (string, optional)
//   - is_active: Whether the category is active (bool, optional)
//
// Response:
//   - 200 OK: Category updated successfully
//   - 400 Bad Request: Invalid request body or validation error
//   - 401 Unauthorized: User not authenticated
//   - 403 Forbidden: Category does not belong to user's account
//   - 404 Not Found: Category not found or user not found
//   - 500 Internal Server Error: Database or service error
//
// Security notes:
//   - Validates user authentication
//   - Validates category ID format
//   - Ensures category belongs to user's account (authorization)
//   - Preserves account ID during update
func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	userID, _ := c.Get("userID")
	user, err := h.service.GetUser(userID.(int))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	// Get existing category
	existingCategory, err := h.service.GetCategory(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	// Check if category belongs to user's account
	if existingCategory.AccountID != user.AccountID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	var category models.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	// Preserve account ID and ID
	category.ID = id
	category.AccountID = user.AccountID

	err = h.service.UpdateCategory(&category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update category: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"category": category})
}

// DeleteCategory deletes a category by ID.
// This endpoint requires authentication and validates that the category belongs
// to the authenticated user's account before deleting it.
//
// Authentication: Required (JWT token in Authorization header)
// Authorization: User must be authenticated and have access to the account
//
// URL Parameters:
//   - id: The category ID to delete (integer)
//
// Response:
//   - 200 OK: Category deleted successfully
//   - 400 Bad Request: Invalid category ID
//   - 401 Unauthorized: User not authenticated
//   - 403 Forbidden: Category does not belong to user's account
//   - 404 Not Found: Category not found or user not found
//   - 500 Internal Server Error: Database or service error
//
// Security notes:
//   - Validates user authentication
//   - Validates category ID format
//   - Ensures category belongs to user's account (authorization)
//   - Prevents deletion if category is still in use
func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	userID, _ := c.Get("userID")
	user, err := h.service.GetUser(userID.(int))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	// Get existing category to check ownership
	existingCategory, err := h.service.GetCategory(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	// Check if category belongs to user's account
	if existingCategory.AccountID != user.AccountID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	err = h.service.DeleteCategory(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete category: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Category deleted successfully"})
}

// GetInventoryItemsByCategory retrieves all inventory items in a specific category.
// This endpoint requires authentication and validates that the category belongs
// to the authenticated user's account before returning the items.
//
// Authentication: Required (JWT token in Authorization header)
// Authorization: User must be authenticated and have access to the account
//
// URL Parameters:
//   - id: The category ID to get items for (integer)
//
// Response:
//   - 200 OK: List of inventory items in the category
//   - 400 Bad Request: Invalid category ID format
//   - 401 Unauthorized: User not authenticated
//   - 403 Forbidden: Category does not belong to user's account
//   - 404 Not Found: Category not found or user not found
//   - 500 Internal Server Error: Database or service error
//
// Security notes:
//   - Validates user authentication
//   - Validates category ID format
//   - Ensures category belongs to user's account (authorization)
//   - Returns appropriate error codes for different scenarios
func (h *CategoryHandler) GetInventoryItemsByCategory(c *gin.Context) {
	userID, _ := c.Get("userID")
	user, err := h.service.GetUser(userID.(int))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	// Verify category exists and belongs to user's account
	category, err := h.service.GetCategory(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	if category.AccountID != user.AccountID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// Get inventory items in the category
	items, err := h.service.GetInventoryItemsByCategory(user.AccountID, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch inventory items by category"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": items})
}

// GetMenuItemsByCategory retrieves all menu items in a specific category.
// This endpoint requires authentication and validates that the category belongs
// to the authenticated user's account before returning the items.
//
// Authentication: Required (JWT token in Authorization header)
// Authorization: User must be authenticated and have access to the account
//
// URL Parameters:
//   - id: The category ID to get items for (integer)
//
// Response:
//   - 200 OK: List of menu items in the category
//   - 400 Bad Request: Invalid category ID format
//   - 401 Unauthorized: User not authenticated
//   - 403 Forbidden: Category does not belong to user's account
//   - 404 Not Found: Category not found or user not found
//   - 500 Internal Server Error: Database or service error
//
// Security notes:
//   - Validates user authentication
//   - Validates category ID format
//   - Ensures category belongs to user's account (authorization)
//   - Returns appropriate error codes for different scenarios
func (h *CategoryHandler) GetMenuItemsByCategory(c *gin.Context) {
	userID, _ := c.Get("userID")
	user, err := h.service.GetUser(userID.(int))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	// Verify category exists and belongs to user's account
	category, err := h.service.GetCategory(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	if category.AccountID != user.AccountID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// Get menu items in the category
	items, err := h.service.GetMenuItemsByCategoryID(user.AccountID, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch menu items by category"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": items})
} 