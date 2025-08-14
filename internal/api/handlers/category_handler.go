// Package handlers provides HTTP request handlers for the application's API endpoints.
// This package includes handlers for category management operations.
// All handlers require authentication and operate within the context of user accounts.
package handlers

import (
	helpers "github.com/mnadev/pantryos/internal/api/helper"
	"net/http"
	"strconv"

	"github.com/mnadev/pantryos/internal/database"
	"github.com/mnadev/pantryos/internal/models"

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
//
//	All responses are wrapped in the standard APIResponse structure.
//	- Success: { "success": true, "message": "...", "data": [...] }
//	- Error:   { "success": false, "message": "...", "error": { "code": "...", "details": "..." } }
//
// Status Codes:
//   - 200 OK: Categories retrieved successfully. The 'data' field contains a list of categories.
//   - 401 Unauthorized: User not authenticated.
//   - 404 Not Found: User associated with token not found in the database.
//   - 500 Internal Server Error: Database or other service error.
func (h *CategoryHandler) GetCategories(c *gin.Context) {
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

	// Get all categories for the account
	categories, err := h.service.GetCategoriesByAccount(user.AccountID)
	if err != nil {
		errDetails := helpers.APIError{Code: "DB_FETCH_FAILED", Details: err.Error()}
		helpers.Error(c.Writer, http.StatusInternalServerError, "Failed to fetch categories.", errDetails)
		return
	}

	// Return a 200 OK response with the list of categories in the data field.
	helpers.Success(c.Writer, http.StatusOK, "Categories retrieved successfully.", categories)
}

// GetActiveCategories retrieves only active categories for the authenticated user's account.
// This endpoint requires authentication and returns only active categories scoped to the user's account.
// This is useful for UI dropdowns and active item categorization.
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
//   - 200 OK: Active categories retrieved successfully. The 'data' field contains a list of active categories.
//   - 401 Unauthorized: User not authenticated.
//   - 404 Not Found: User associated with token not found in the database.
//   - 500 Internal Server Error: Database or other service error.
func (h *CategoryHandler) GetActiveCategories(c *gin.Context) {
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

	// Get active categories for the account
	categories, err := h.service.GetActiveCategoriesByAccount(user.AccountID)
	if err != nil {
		errDetails := helpers.APIError{Code: "DB_FETCH_FAILED", Details: err.Error()}
		helpers.Error(c.Writer, http.StatusInternalServerError, "Failed to fetch active categories.", errDetails)
		return
	}

	// Return a 200 OK response with the list of active categories in the data field.
	helpers.Success(c.Writer, http.StatusOK, "Active categories retrieved successfully.", categories)
}

// CreateCategory creates a new category for the authenticated user's account.
// This endpoint accepts category details in JSON format and creates a new
// category in the database. The account ID is automatically set from the authenticated user.
//
// Authentication: Required (JWT token in Authorization header)
// Authorization: User must be authenticated and have access to the account
//
// Response:
//
//	All responses are wrapped in the standard APIResponse structure.
//	- Success: { "success": true, "message": "...", "data": ... }
//	- Error:   { "success": false, "message": "...", "error": { "code": "...", "details": "..." } }
//
// Status Codes:
//   - 201 Created: Category created successfully. The 'data' field contains the new category.
//   - 400 Bad Request: Invalid request body or validation error.
//   - 401 Unauthorized: User not authenticated.
//   - 404 Not Found: User associated with token not found in the database.
//   - 500 Internal Server Error: Database or other service error.
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
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
	var category models.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		errDetails := helpers.APIError{Code: "INVALID_INPUT", Details: err.Error()}
		helpers.Error(c.Writer, http.StatusBadRequest, "Invalid request body.", errDetails)
		return
	}

	// Set account ID from the authenticated user to ensure proper scoping
	category.AccountID = user.AccountID

	// Create the category in the database
	err = h.service.CreateCategory(&category)
	if err != nil {
		// Consider checking for a unique constraint violation to return a 409 Conflict status.
		errDetails := helpers.APIError{Code: "DB_INSERT_FAILED", Details: err.Error()}
		helpers.Error(c.Writer, http.StatusInternalServerError, "Failed to create category.", errDetails)
		return
	}

	// Return a 201 Created response with the newly created category.
	helpers.Success(c.Writer, http.StatusCreated, "Category created successfully.", category)
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
//
//	All responses are wrapped in the standard APIResponse structure.
//	- Success: { "success": true, "message": "...", "data": {...} }
//	- Error:   { "success": false, "message": "...", "error": { "code": "...", "details": "..." } }
//
// Status Codes:
//   - 200 OK: Category retrieved successfully. The 'data' field contains the category object.
//   - 400 Bad Request: Invalid category ID format in URL.
//   - 401 Unauthorized: User not authenticated.
//   - 403 Forbidden: The requested category does not belong to the user's account.
//   - 404 Not Found: The user or the category could not be found.
func (h *CategoryHandler) GetCategory(c *gin.Context) {
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

	// Parse and validate the category ID from the URL parameter
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		errDetails := helpers.APIError{Code: "INVALID_INPUT", Details: "Category ID must be a valid integer."}
		helpers.Error(c.Writer, http.StatusBadRequest, "Invalid Category ID.", errDetails)
		return
	}

	// Retrieve the category from the database
	category, err := h.service.GetCategory(id)
	if err != nil {
		errDetails := helpers.APIError{Code: "CATEGORY_NOT_FOUND", Details: err.Error()}
		helpers.Error(c.Writer, http.StatusNotFound, "Category not found.", errDetails)
		return
	}

	// Authorization check: Ensure the category belongs to the user's account
	if category.AccountID != user.AccountID {
		errDetails := helpers.APIError{Code: "FORBIDDEN", Details: "You do not have permission to access this category."}
		helpers.Error(c.Writer, http.StatusForbidden, "Access denied.", errDetails)
		return
	}

	// Return a 200 OK response with the category object in the data field.
	helpers.Success(c.Writer, http.StatusOK, "Category retrieved successfully.", category)
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
// Response:
//
//	All responses are wrapped in the standard APIResponse structure.
//	- Success: { "success": true, "message": "...", "data": {...} }
//	- Error:   { "success": false, "message": "...", "error": { "code": "...", "details": "..." } }
//
// Status Codes:
//   - 200 OK: Category updated successfully. The 'data' field contains the updated category object.
//   - 400 Bad Request: Invalid category ID format or invalid request body.
//   - 401 Unauthorized: User not authenticated.
//   - 403 Forbidden: The category does not belong to the user's account.
//   - 404 Not Found: The user or the category could not be found.
//   - 500 Internal Server Error: Database or other service error.
func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
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

	// Parse and validate the category ID from the URL parameter
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		errDetails := helpers.APIError{Code: "INVALID_INPUT", Details: "Category ID must be a valid integer."}
		helpers.Error(c.Writer, http.StatusBadRequest, "Invalid Category ID.", errDetails)
		return
	}

	// Get existing category to verify ownership
	existingCategory, err := h.service.GetCategory(id)
	if err != nil {
		errDetails := helpers.APIError{Code: "CATEGORY_NOT_FOUND", Details: err.Error()}
		helpers.Error(c.Writer, http.StatusNotFound, "Category not found.", errDetails)
		return
	}

	// Authorization check: Ensure the category belongs to the user's account
	if existingCategory.AccountID != user.AccountID {
		errDetails := helpers.APIError{Code: "FORBIDDEN", Details: "You do not have permission to modify this category."}
		helpers.Error(c.Writer, http.StatusForbidden, "Access denied.", errDetails)
		return
	}

	// Parse and validate the JSON request body with the updates
	var category models.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		errDetails := helpers.APIError{Code: "INVALID_INPUT", Details: err.Error()}
		helpers.Error(c.Writer, http.StatusBadRequest, "Invalid request body.", errDetails)
		return
	}

	// Preserve the original ID and AccountID to prevent them from being changed.
	category.ID = id
	category.AccountID = user.AccountID

	// Update the category in the database
	err = h.service.UpdateCategory(&category)
	if err != nil {
		errDetails := helpers.APIError{Code: "DB_UPDATE_FAILED", Details: err.Error()}
		helpers.Error(c.Writer, http.StatusInternalServerError, "Failed to update category.", errDetails)
		return
	}

	// Return a 200 OK response with the updated category object.
	helpers.Success(c.Writer, http.StatusOK, "Category updated successfully.", category)
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
//
//	All responses are wrapped in the standard APIResponse structure.
//	- Success: { "success": true, "message": "...", "data": null }
//	- Error:   { "success": false, "message": "...", "error": { "code": "...", "details": "..." } }
//
// Status Codes:
//   - 200 OK: Category deleted successfully.
//   - 400 Bad Request: Invalid category ID format.
//   - 401 Unauthorized: User not authenticated.
//   - 403 Forbidden: The category does not belong to the user's account.
//   - 404 Not Found: The user or the category could not be found.
//   - 500 Internal Server Error: Database or service error.
func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
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

	// Parse and validate the category ID from the URL parameter
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		errDetails := helpers.APIError{Code: "INVALID_INPUT", Details: "Category ID must be a valid integer."}
		helpers.Error(c.Writer, http.StatusBadRequest, "Invalid Category ID.", errDetails)
		return
	}

	// Get existing category to check ownership
	existingCategory, err := h.service.GetCategory(id)
	if err != nil {
		errDetails := helpers.APIError{Code: "CATEGORY_NOT_FOUND", Details: err.Error()}
		helpers.Error(c.Writer, http.StatusNotFound, "Category not found.", errDetails)
		return
	}

	// Authorization check: Ensure the category belongs to the user's account
	if existingCategory.AccountID != user.AccountID {
		errDetails := helpers.APIError{Code: "FORBIDDEN", Details: "You do not have permission to delete this category."}
		helpers.Error(c.Writer, http.StatusForbidden, "Access denied.", errDetails)
		return
	}

	// Attempt to delete the category
	err = h.service.DeleteCategory(id)
	if err != nil {
		// The service might return a specific error if the category is in use.
		// This could be handled here to return a 409 Conflict status.
		errDetails := helpers.APIError{Code: "DB_DELETE_FAILED", Details: err.Error()}
		helpers.Error(c.Writer, http.StatusInternalServerError, "Failed to delete category.", errDetails)
		return
	}

	// Return a 200 OK response with a success message.
	helpers.Success(c.Writer, http.StatusOK, "Category deleted successfully.", nil)
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
//
//	All responses are wrapped in the standard APIResponse structure.
//	- Success: { "success": true, "message": "...", "data": [...] }
//	- Error:   { "success": false, "message": "...", "error": { "code": "...", "details": "..." } }
//
// Status Codes:
//   - 200 OK: Inventory items retrieved successfully. The 'data' field contains a list of items.
//   - 400 Bad Request: Invalid category ID format.
//   - 401 Unauthorized: User not authenticated.
//   - 403 Forbidden: The requested category does not belong to the user's account.
//   - 404 Not Found: The user or the category could not be found.
//   - 500 Internal Server Error: Database or other service error.
func (h *CategoryHandler) GetInventoryItemsByCategory(c *gin.Context) {
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

	// Parse and validate the category ID from the URL parameter
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		errDetails := helpers.APIError{Code: "INVALID_INPUT", Details: "Category ID must be a valid integer."}
		helpers.Error(c.Writer, http.StatusBadRequest, "Invalid Category ID.", errDetails)
		return
	}

	// Verify category exists and belongs to the user's account
	category, err := h.service.GetCategory(id)
	if err != nil {
		errDetails := helpers.APIError{Code: "CATEGORY_NOT_FOUND", Details: err.Error()}
		helpers.Error(c.Writer, http.StatusNotFound, "Category not found.", errDetails)
		return
	}

	if category.AccountID != user.AccountID {
		errDetails := helpers.APIError{Code: "FORBIDDEN", Details: "You do not have permission to access this category."}
		helpers.Error(c.Writer, http.StatusForbidden, "Access denied.", errDetails)
		return
	}

	// Get inventory items in the specified category
	items, err := h.service.GetInventoryItemsByCategory(user.AccountID, id)
	if err != nil {
		errDetails := helpers.APIError{Code: "DB_FETCH_FAILED", Details: err.Error()}
		helpers.Error(c.Writer, http.StatusInternalServerError, "Failed to fetch inventory items.", errDetails)
		return
	}

	// Return a 200 OK response with the list of items in the data field.
	helpers.Success(c.Writer, http.StatusOK, "Inventory items retrieved successfully.", items)
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
//
//	All responses are wrapped in the standard APIResponse structure.
//	- Success: { "success": true, "message": "...", "data": [...] }
//	- Error:   { "success": false, "message": "...", "error": { "code": "...", "details": "..." } }
//
// Status Codes:
//   - 200 OK: Menu items retrieved successfully. The 'data' field contains a list of menu items.
//   - 400 Bad Request: Invalid category ID format.
//   - 401 Unauthorized: User not authenticated.
//   - 403 Forbidden: The requested category does not belong to the user's account.
//   - 404 Not Found: The user or the category could not be found.
//   - 500 Internal Server Error: Database or other service error.
func (h *CategoryHandler) GetMenuItemsByCategory(c *gin.Context) {
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

	// Parse and validate the category ID from the URL parameter
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		errDetails := helpers.APIError{Code: "INVALID_INPUT", Details: "Category ID must be a valid integer."}
		helpers.Error(c.Writer, http.StatusBadRequest, "Invalid Category ID.", errDetails)
		return
	}

	// Verify category exists and belongs to the user's account
	category, err := h.service.GetCategory(id)
	if err != nil {
		errDetails := helpers.APIError{Code: "CATEGORY_NOT_FOUND", Details: err.Error()}
		helpers.Error(c.Writer, http.StatusNotFound, "Category not found.", errDetails)
		return
	}

	if category.AccountID != user.AccountID {
		errDetails := helpers.APIError{Code: "FORBIDDEN", Details: "You do not have permission to access this category."}
		helpers.Error(c.Writer, http.StatusForbidden, "Access denied.", errDetails)
		return
	}

	// Get menu items in the specified category
	items, err := h.service.GetMenuItemsByCategoryID(user.AccountID, id)
	if err != nil {
		errDetails := helpers.APIError{Code: "DB_FETCH_FAILED", Details: err.Error()}
		helpers.Error(c.Writer, http.StatusInternalServerError, "Failed to fetch menu items.", errDetails)
		return
	}

	// Return a 200 OK response with the list of menu items in the data field.
	helpers.Success(c.Writer, http.StatusOK, "Menu items retrieved successfully.", items)
}
