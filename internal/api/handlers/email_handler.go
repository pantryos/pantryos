package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mnadev/stok/internal/database"
	"github.com/mnadev/stok/internal/email"
	"github.com/mnadev/stok/internal/models"
)

// EmailHandler handles email-related API endpoints
type EmailHandler struct {
	service      *database.Service
	emailService *email.EmailService
}

// NewEmailHandler creates a new email handler
func NewEmailHandler(db *database.DB) *EmailHandler {
	return &EmailHandler{
		service:      database.NewService(db),
		emailService: email.NewEmailService(),
	}
}

// SendVerificationEmail sends a verification email to a user
// @Summary Send verification email
// @Description Send a verification email to a user
// @Tags email
// @Accept json
// @Produce json
// @Param user_id path int true "User ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /email/verification/{user_id} [post]
func (h *EmailHandler) SendVerificationEmail(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Get user from database
	user, err := h.service.GetUser(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Get account for the user
	account, err := h.service.GetAccount(user.AccountID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		return
	}

	// Generate verification token
	token, err := h.generateVerificationToken(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate verification token"})
		return
	}

	// Create verification URL
	baseURL := getBaseURL(c)
	verificationURL := fmt.Sprintf("%s/verify-email?token=%s", baseURL, token)

	// Send verification email
	if err := h.emailService.SendVerificationEmail(*user, *account, verificationURL); err != nil {
		// Log email failure
		h.logEmailFailure(user.AccountID, &userID, user.Email, "Verify Your Stok Account", models.EmailTypeVerification, err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send verification email"})
		return
	}

	// Log successful email
	h.logEmailSuccess(user.AccountID, &userID, user.Email, "Verify Your Stok Account", models.EmailTypeVerification)

	c.JSON(http.StatusOK, gin.H{
		"message": "Verification email sent successfully",
		"user_id": userID,
	})
}

// VerifyEmail verifies a user's email using a token
// @Summary Verify email
// @Description Verify a user's email using a verification token
// @Tags email
// @Accept json
// @Produce json
// @Param token query string true "Verification token"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /email/verify [get]
func (h *EmailHandler) VerifyEmail(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Verification token is required"})
		return
	}

	// For now, we'll just return success since we need to add email verification methods to the service
	// TODO: Implement proper token validation through service layer
	c.JSON(http.StatusOK, gin.H{
		"message": "Email verified successfully",
	})
}

// SendWeeklyStockReport sends a weekly stock report email
// @Summary Send weekly stock report
// @Description Send a weekly stock report email to all users in an account
// @Tags email
// @Accept json
// @Produce json
// @Param account_id path int true "Account ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /email/weekly-report/{account_id} [post]
func (h *EmailHandler) SendWeeklyStockReport(c *gin.Context) {
	accountID, err := strconv.Atoi(c.Param("account_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID"})
		return
	}

	// Get account
	account, err := h.service.GetAccount(accountID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		return
	}

	// Get all users in the account
	users, err := h.service.GetUsersByAccount(accountID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get users"})
		return
	}

	if len(users) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No users found in account"})
		return
	}

	// Generate stock report data
	stockData, err := h.generateStockReportData(accountID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate stock report"})
		return
	}

	// Send weekly stock report
	if err := h.emailService.SendWeeklyStockReport(*account, users, stockData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send weekly stock report"})
		return
	}

	// Log successful email sending for each user
	for _, user := range users {
		h.logEmailSuccess(accountID, &user.ID, user.Email, fmt.Sprintf("Weekly Stock Report - %s", account.Name), models.EmailTypeWeeklyReport)
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Weekly stock report sent successfully",
		"account_id":  accountID,
		"users_count": len(users),
	})
}

// SendLowStockAlert sends a low stock alert email
// @Summary Send low stock alert
// @Description Send a low stock alert email to all users in an account
// @Tags email
// @Accept json
// @Produce json
// @Param account_id path int true "Account ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /email/low-stock-alert/{account_id} [post]
func (h *EmailHandler) SendLowStockAlert(c *gin.Context) {
	accountID, err := strconv.Atoi(c.Param("account_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID"})
		return
	}

	// Get account
	account, err := h.service.GetAccount(accountID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		return
	}

	// Get all users in the account
	users, err := h.service.GetUsersByAccount(accountID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get users"})
		return
	}

	if len(users) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No users found in account"})
		return
	}

	// Get low stock items
	lowStockItems, err := h.service.GetLowStockItems(accountID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get inventory items"})
		return
	}

	if len(lowStockItems) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"message":    "No low stock items found",
			"account_id": accountID,
		})
		return
	}

	// Send low stock alert
	if err := h.emailService.SendLowStockAlert(*account, users, lowStockItems); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send low stock alert"})
		return
	}

	// Log successful email sending for each user
	for _, user := range users {
		h.logEmailSuccess(accountID, &user.ID, user.Email, fmt.Sprintf("Low Stock Alert - %s", account.Name), models.EmailTypeLowStockAlert)
	}

	c.JSON(http.StatusOK, gin.H{
		"message":               "Low stock alert sent successfully",
		"account_id":            accountID,
		"users_count":           len(users),
		"low_stock_items_count": len(lowStockItems),
	})
}

// generateVerificationToken creates a new verification token for a user
func (h *EmailHandler) generateVerificationToken(userID int) (string, error) {
	// Generate random token
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", err
	}
	token := hex.EncodeToString(tokenBytes)

	// TODO: Add CreateEmailVerificationToken method to service
	// For now, we'll just return the token
	// verificationToken := models.EmailVerificationToken{
	// 	UserID:    userID,
	// 	Token:     token,
	// 	Type:      models.TokenTypeEmailVerification,
	// 	ExpiresAt: time.Now().Add(24 * time.Hour), // 24 hour expiration
	// }
	return token, nil
}

// generateStockReportData generates stock report data for an account
func (h *EmailHandler) generateStockReportData(accountID int) (*email.StockReportData, error) {
	// Get all inventory items for the account
	items, err := h.service.GetInventoryItemsByAccount(accountID)
	if err != nil {
		return nil, err
	}

	// Get latest inventory snapshot
	latestSnapshot, err := h.service.GetLatestInventorySnapshot(accountID)
	if err != nil {
		return nil, err
	}

	// Generate stock report data
	stockData := &email.StockReportData{
		ReportDate: time.Now(),
		TotalItems: len(items),
		Items:      make([]email.StockItemData, 0, len(items)),
	}

	var totalValue float64
	var lowStockCount, outOfStockCount int

	for _, item := range items {
		currentStock := 0.0
		if counts, exists := latestSnapshot.Counts[item.ID]; exists {
			currentStock = counts
		}

		// Calculate item value
		itemValue := currentStock * item.CostPerUnit
		totalValue += itemValue

		// Determine status
		status := "normal"
		if currentStock <= 0 {
			status = "out"
			outOfStockCount++
		} else if currentStock <= item.MinStockLevel {
			status = "low"
			lowStockCount++
		}

		// Get category name
		categoryName := ""
		if item.CategoryID != nil {
			category, err := h.service.GetCategory(*item.CategoryID)
			if err == nil {
				categoryName = category.Name
			}
		}

		stockData.Items = append(stockData.Items, email.StockItemData{
			ID:           item.ID,
			Name:         item.Name,
			Category:     categoryName,
			CurrentStock: currentStock,
			MinStock:     item.MinStockLevel,
			MaxStock:     item.MaxStockLevel,
			Unit:         item.Unit,
			Status:       status,
		})
	}

	stockData.TotalValue = totalValue
	stockData.LowStockItems = lowStockCount
	stockData.OutOfStockItems = outOfStockCount

	return stockData, nil
}

// logEmailSuccess logs a successful email send
func (h *EmailHandler) logEmailSuccess(accountID int, userID *int, toEmail, subject, emailType string) {
	emailLog := models.EmailLog{
		AccountID: accountID,
		UserID:    userID,
		ToEmail:   toEmail,
		Subject:   subject,
		EmailType: emailType,
		Status:    models.EmailStatusSent,
	}

	// TODO: Add CreateEmailLog method to service
	// For now, we'll just ignore the logging
	_ = emailLog
}

// logEmailFailure logs a failed email send
func (h *EmailHandler) logEmailFailure(accountID int, userID *int, toEmail, subject, emailType, errorMsg string) {
	emailLog := models.EmailLog{
		AccountID: accountID,
		UserID:    userID,
		ToEmail:   toEmail,
		Subject:   subject,
		EmailType: emailType,
		Status:    models.EmailStatusFailed,
		ErrorMsg:  errorMsg,
	}

	// TODO: Add CreateEmailLog method to service
	// For now, we'll just ignore the logging
	_ = emailLog
}

// getBaseURL gets the base URL for the application
func getBaseURL(c *gin.Context) string {
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	return fmt.Sprintf("%s://%s", scheme, c.Request.Host)
}
