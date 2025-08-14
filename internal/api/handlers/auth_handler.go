package handlers

import (
	helpers "github.com/mnadev/pantryos/internal/api/helper"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/mnadev/pantryos/internal/auth"
	"github.com/mnadev/pantryos/internal/database"
	"github.com/mnadev/pantryos/internal/models"
	"github.com/mnadev/pantryos/pkg/utils"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service *database.Service
}

func NewAuthHandler(db *database.DB) *AuthHandler {
	return &AuthHandler{service: database.NewService(db)}
}

// RegisterRequest represents the registration request body
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	// AccountID is no longer required - it will be determined from the invitation
}

// LoginRequest represents the login request body
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// CreateInvitationRequest represents the request body for creating an invitation
type CreateInvitationRequest struct {
	Email     string    `json:"email" binding:"required,email"`
	ExpiresAt time.Time `json:"expires_at"`
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user with email and password (account ID determined from invitation)
// @Tags authentication
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Registration details"
// @Success 201 {object} map[string]interface{} "User registered successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body or no invitation found"
// @Failure 409 {object} map[string]interface{} "User already exists"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errDetails := helpers.APIError{Code: "INVALID_INPUT", Details: err.Error()}
		helpers.Error(c.Writer, http.StatusBadRequest, "Invalid request body provided.", errDetails)
		return
	}

	// Check if user already exists
	existingUser, err := h.service.GetUserByEmail(req.Email)
	if err == nil && existingUser != nil {
		errDetails := helpers.APIError{Code: "USER_ALREADY_EXISTS"}
		helpers.Error(c.Writer, http.StatusConflict, "A user with this email already exists.", errDetails)
		return
	}

	// Check if user has a pending invitation
	invitation, err := h.service.GetPendingInvitationByEmail(req.Email)
	if err != nil {
		errDetails := helpers.APIError{Code: "INVITATION_NOT_FOUND"}
		helpers.Error(c.Writer, http.StatusBadRequest, "No invitation found for this email. Please request one from your account administrator.", errDetails)
		return
	}

	// Check if invitation has expired
	if time.Now().After(invitation.ExpiresAt) {
		errDetails := helpers.APIError{Code: "INVITATION_EXPIRED"}
		helpers.Error(c.Writer, http.StatusBadRequest, "Your invitation has expired. Please request a new one from your account administrator.", errDetails)
		return
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		errDetails := helpers.APIError{Code: "INTERNAL_SERVER_ERROR", Details: "Failed to hash password"}
		helpers.Error(c.Writer, http.StatusInternalServerError, "An internal error occurred. Please try again later.", errDetails)
		return
	}

	// Create user with account ID from invitation
	user := &models.User{
		Email:     req.Email,
		Password:  hashedPassword,
		AccountID: invitation.AccountID,
	}

	err = h.service.CreateUser(user)
	if err != nil {
		errDetails := helpers.APIError{Code: "USER_CREATION_FAILED", Details: err.Error()}
		helpers.Error(c.Writer, http.StatusInternalServerError, "Failed to create the user account.", errDetails)
		return
	}

	// Mark invitation as accepted
	invitation.Status = models.AccountInvitationStatusAccepted
	now := time.Now()
	invitation.AcceptedAt = &now
	err = h.service.UpdateInvitation(invitation)
	if err != nil {
		log.Printf("CRITICAL: Failed to update invitation %d after user registration: %v", invitation.ID, err)
	}

	responseData := gin.H{
		"user_id": user.ID,
	}
	helpers.Success(c.Writer, http.StatusCreated, "User registered successfully.", responseData)
}

// Login godoc
// @Summary Login user
// @Description Authenticate user and return JWT token
// @Tags authentication
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login credentials"
// @Success 200 {object} map[string]interface{} "Login successful"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 401 {object} map[string]interface{} "Invalid credentials"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	// Get user by email
	user, err := h.service.GetUserByEmail(req.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Check password
	if !utils.CheckPasswordHash(req.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate JWT token
	token, err := auth.GenerateJWT(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":         user.ID,
			"email":      user.Email,
			"account_id": user.AccountID,
		},
	})
}

// GetCurrentUser godoc
// @Summary      Get current user information
// @Description  Get information about the currently authenticated user based on the provided bearer token.
// @Tags         authentication
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  helpers.APIResponse{data=GetUserSuccessData}  "Successfully retrieved user information"
// @Failure      401  {object}  helpers.APIResponse                           "Error: User not authenticated"
// @Failure      404  {object}  helpers.APIResponse                           "Error: User not found"
// @Router       /api/v1/me [get]
func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	userIDInterface, exists := c.Get("userID")
	if !exists {
		errDetails := helpers.APIError{Code: "UNAUTHORIZED", Details: "User ID not found in request context."}
		helpers.Error(c.Writer, http.StatusUnauthorized, "User not authenticated.", errDetails)
		return
	}

	// We assume userID from the middleware is of the correct type.
	// A robust implementation would also check the `ok` value from the type assertion.
	userID := userIDInterface.(int)

	user, err := h.service.GetUser(userID)
	if err != nil {
		errDetails := helpers.APIError{Code: "USER_NOT_FOUND", Details: err.Error()}
		helpers.Error(c.Writer, http.StatusNotFound, "User not found.", errDetails)
		return
	}

	// Construct the specific data payload for the success response.
	responseData := helpers.GetUserSuccessData{
		ID:        user.ID,
		Email:     user.Email,
		AccountID: user.AccountID,
		CreatedAt: user.CreatedAt,
	}

	// Use the standard success helper to return the data.
	helpers.Success(c.Writer, http.StatusOK, "Current user retrieved successfully.", responseData)
}

// GetAvailableAccounts godoc
// @Summary Get available accounts for registration
// @Description Retrieve a list of available accounts that users can register for
// @Tags authentication
// @Produce json
// @Success 200 {array} models.Account "List of available accounts"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /auth/accounts [get]
func (h *AuthHandler) GetAvailableAccounts(c *gin.Context) {
	// For now, we'll get all accounts. In a real system, you might want to filter by organization
	// or add pagination for large numbers of accounts
	accounts, err := h.service.GetAllAccounts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve accounts: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, accounts)
}

// CreateInvitation godoc
// @Summary Create a new invitation
// @Description Create an invitation for a user to join an account
// @Tags invitations
// @Accept json
// @Produce json
// @Param account_id path int true "Account ID"
// @Param request body CreateInvitationRequest true "Invitation details"
// @Success 201 {object} map[string]interface{} "Invitation created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 403 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/accounts/{account_id}/invitations [post]
func (h *AuthHandler) CreateInvitation(c *gin.Context) {
	// Get current user ID from context
	userIDInterface, exists := c.Get("userID")
	if !exists {
		errDetails := helpers.APIError{Code: "UNAUTHORIZED", Details: "User ID not found in context."}
		helpers.Error(c.Writer, http.StatusUnauthorized, "User not authenticated.", errDetails)
		return
	}

	userID, ok := userIDInterface.(int)
	if !ok {
		errDetails := helpers.APIError{Code: "INTERNAL_SERVER_ERROR", Details: "User ID in context is not a valid integer."}
		helpers.Error(c.Writer, http.StatusInternalServerError, "An internal server error occurred.", errDetails)
		return
	}

	// Get account ID from URL
	accountIDStr := c.Param("account_id")
	accountID, err := strconv.Atoi(accountIDStr)
	if err != nil {
		errDetails := helpers.APIError{Code: "INVALID_INPUT", Details: "Account ID in URL must be a valid integer."}
		helpers.Error(c.Writer, http.StatusBadRequest, "Invalid Account ID provided.", errDetails)
		return
	}

	// TODO: Add authorization check - ensure user is admin of this account
	// For now, we'll allow any authenticated user to create invitations

	var req CreateInvitationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errDetails := helpers.APIError{Code: "INVALID_INPUT", Details: err.Error()}
		helpers.Error(c.Writer, http.StatusBadRequest, "Invalid request body.", errDetails)
		return
	}

	// Create invitation
	invitation := &models.AccountInvitation{
		AccountID: accountID,
		Email:     req.Email,
		InvitedBy: userID,
		ExpiresAt: req.ExpiresAt,
	}

	err = h.service.CreateInvitation(invitation)
	if err != nil {
		errDetails := helpers.APIError{Code: "DB_INSERT_FAILED", Details: err.Error()}
		helpers.Error(c.Writer, http.StatusInternalServerError, "Failed to create invitation.", errDetails)
		return
	}

	// Use the Success helper for the final response
	responseData := gin.H{
		"invitation_id": invitation.ID,
	}
	helpers.Success(c.Writer, http.StatusCreated, "Invitation created successfully.", responseData)
}

// GetInvitationsByAccount godoc
// @Summary Get invitations for an account
// @Description Retrieve all invitations for a specific account
// @Tags invitations
// @Produce json
// @Param account_id path int true "Account ID"
// @Success 200 {array} models.AccountInvitation "List of invitations"
// @Failure 400 {object} map[string]interface{} "Invalid account ID"
// @Failure 403 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/accounts/{account_id}/invitations [get]
func (h *AuthHandler) GetInvitationsByAccount(c *gin.Context) {
	accountIDStr := c.Param("account_id")
	accountID, err := strconv.Atoi(accountIDStr)
	if err != nil {
		errDetails := helpers.APIError{Code: "INVALID_INPUT", Details: "Account ID in URL must be a valid integer."}
		helpers.Error(c.Writer, http.StatusBadRequest, "Invalid Account ID provided.", errDetails)
		return
	}

	// TODO: Add authorization check - ensure the current user is an admin of this account.

	// Retrieve invitations from the service layer
	invitations, err := h.service.GetInvitationsByAccount(accountID)
	if err != nil {
		errDetails := helpers.APIError{Code: "DB_FETCH_FAILED", Details: err.Error()}
		helpers.Error(c.Writer, http.StatusInternalServerError, "Failed to retrieve invitations.", errDetails)
		return
	}

	helpers.Success(c.Writer, http.StatusOK, "Invitations retrieved successfully.", invitations)
}

// DeleteInvitation godoc
// @Summary Delete an invitation
// @Description Delete/revoke an invitation
// @Tags invitations
// @Produce json
// @Param account_id path int true "Account ID"
// @Param invitation_id path int true "Invitation ID"
// @Success 200 {object} map[string]interface{} "Invitation deleted successfully"
// @Failure 400 {object} map[string]interface{} "Invalid invitation ID"
// @Failure 403 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/accounts/{account_id}/invitations/{invitation_id} [delete]
func (h *AuthHandler) DeleteInvitation(c *gin.Context) {
	// Get invitation ID from URL
	invitationIDStr := c.Param("invitation_id")
	invitationID, err := strconv.Atoi(invitationIDStr)
	if err != nil {
		errDetails := helpers.APIError{Code: "INVALID_INPUT", Details: "Invitation ID must be a valid integer."}
		helpers.Error(c.Writer, http.StatusBadRequest, "Invalid Invitation ID.", errDetails)
		return
	}

	// TODO: Add authorization check to ensure the user is an admin of this account.

	// Attempt to delete the invitation
	err = h.service.DeleteInvitation(invitationID)
	if err != nil {
		errDetails := helpers.APIError{Code: "DB_DELETE_FAILED", Details: err.Error()}
		helpers.Error(c.Writer, http.StatusInternalServerError, "Failed to delete invitation.", errDetails)
		return
	}

	helpers.Success(c.Writer, http.StatusOK, "Invitation deleted successfully.", nil)
}
