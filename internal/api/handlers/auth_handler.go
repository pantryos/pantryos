package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/mnadev/stok/internal/auth"
	"github.com/mnadev/stok/internal/database"
	"github.com/mnadev/stok/internal/models"
	"github.com/mnadev/stok/pkg/utils"

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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	// Check if user already exists
	existingUser, err := h.service.GetUserByEmail(req.Email)
	if err == nil && existingUser != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User with this email already exists"})
		return
	}

	// Check if user has a pending invitation
	invitation, err := h.service.GetPendingInvitationByEmail(req.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No invitation found for this email. Please contact your account administrator."})
		return
	}

	// Check if invitation has expired
	if time.Now().After(invitation.ExpiresAt) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invitation has expired. Please contact your account administrator for a new invitation."})
		return
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
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
		// Check for specific validation errors
		if err.Error() == "invalid account ID" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account ID"})
			return
		}
		if err.Error() == "invalid user role" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user role"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user: " + err.Error()})
		return
	}

	// Mark invitation as accepted
	invitation.Status = "accepted"
	now := time.Now()
	invitation.AcceptedAt = &now
	err = h.service.UpdateInvitation(invitation)
	if err != nil {
		// Log the error but don't fail the registration
		// The user is already created successfully
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"user_id": user.ID,
	})
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
// @Summary Get current user information
// @Description Get information about the currently authenticated user
// @Tags authentication
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "User information"
// @Failure 401 {object} map[string]interface{} "User not authenticated"
// @Failure 404 {object} map[string]interface{} "User not found"
// @Router /api/v1/me [get]
func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	user, err := h.service.GetUser(userID.(int))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":         user.ID,
			"email":      user.Email,
			"account_id": user.AccountID,
			"created_at": user.CreatedAt,
		},
	})
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
	// Get current user
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Get account ID from URL
	accountIDStr := c.Param("account_id")
	accountID, err := strconv.Atoi(accountIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID"})
		return
	}

	// TODO: Add authorization check - ensure user is admin of this account
	// For now, we'll allow any authenticated user to create invitations

	var req CreateInvitationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	// Create invitation
	invitation := &models.AccountInvitation{
		AccountID: accountID,
		Email:     req.Email,
		InvitedBy: userID.(int),
		ExpiresAt: req.ExpiresAt,
	}

	err = h.service.CreateInvitation(invitation)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create invitation: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":       "Invitation created successfully",
		"invitation_id": invitation.ID,
	})
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
	// Get account ID from URL
	accountIDStr := c.Param("account_id")
	accountID, err := strconv.Atoi(accountIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID"})
		return
	}

	// TODO: Add authorization check - ensure user is admin of this account
	// For now, we'll allow any authenticated user to view invitations

	invitations, err := h.service.GetInvitationsByAccount(accountID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve invitations: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, invitations)
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid invitation ID"})
		return
	}

	// TODO: Add authorization check - ensure user is admin of this account
	// For now, we'll allow any authenticated user to delete invitations

	err = h.service.DeleteInvitation(invitationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete invitation: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Invitation deleted successfully"})
}
