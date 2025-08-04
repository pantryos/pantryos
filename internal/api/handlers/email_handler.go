package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mnadev/pantryos/internal/database"
	"github.com/mnadev/pantryos/internal/email"
	"github.com/mnadev/pantryos/internal/models"
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
		h.logEmailFailure(user.AccountID, &userID, user.Email, "Verify Your PantryOS Account", models.EmailTypeVerification, err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send verification email"})
		return
	}

	// Log successful email
	h.logEmailSuccess(user.AccountID, &userID, user.Email, "Verify Your PantryOS Account", models.EmailTypeVerification)

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

// SendWeeklySupplyChainReport sends a weekly supply chain report email
// @Summary Send weekly supply chain report
// @Description Send a weekly supply chain report email to all users in an account
// @Tags email
// @Accept json
// @Produce json
// @Param account_id path int true "Account ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /email/weekly-supply-chain/{account_id} [post]
func (h *EmailHandler) SendWeeklySupplyChainReport(c *gin.Context) {
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

	// Generate supply chain report data
	supplyChainData, err := h.generateSupplyChainReportData(accountID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate supply chain report"})
		return
	}

	// Send weekly supply chain report
	if err := h.emailService.SendWeeklySupplyChainReport(*account, users, supplyChainData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send weekly supply chain report"})
		return
	}

	// Log successful email sending for each user
	for _, user := range users {
		h.logEmailSuccess(accountID, &user.ID, user.Email, fmt.Sprintf("Weekly Supply Chain Report - %s", account.Name), models.EmailTypeWeeklySupplyChain)
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Weekly supply chain report sent successfully",
		"account_id":  accountID,
		"users_count": len(users),
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

// generateSupplyChainReportData generates supply chain report data for an account
func (h *EmailHandler) generateSupplyChainReportData(accountID int) (*email.SupplyChainData, error) {
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

	// Get recent deliveries for vendor information
	recentDeliveries, err := h.service.GetDeliveriesByAccount(accountID)
	if err != nil {
		return nil, err
	}

	// Generate supply chain report data
	supplyChainData := &email.SupplyChainData{
		ReportDate: time.Now(),
		TotalItems: len(items),
		Items:      make([]email.SupplyChainItemData, 0, len(items)),
	}

	var totalValue float64
	var lowStockCount, outOfStockCount, criticalCount int
	var estimatedReorders float64

	for _, item := range items {
		currentStock := 0.0
		if counts, exists := latestSnapshot.Counts[item.ID]; exists {
			currentStock = counts
		}

		// Calculate item value
		itemValue := currentStock * item.CostPerUnit
		totalValue += itemValue

		// Determine status and calculate supply chain metrics
		status := "normal"
		reorderQuantity := 0.0
		daysUntilStockout := 999 // Default to high number

		if currentStock <= 0 {
			status = "out"
			outOfStockCount++
			reorderQuantity = item.MaxStockLevel
		} else if currentStock <= item.MinStockLevel*0.5 {
			status = "critical"
			criticalCount++
			reorderQuantity = item.MaxStockLevel - currentStock
			daysUntilStockout = int(currentStock / (item.MinStockLevel / 7)) // Rough estimate
		} else if currentStock <= item.MinStockLevel {
			status = "low"
			lowStockCount++
			reorderQuantity = item.MaxStockLevel - currentStock
			daysUntilStockout = int(currentStock / (item.MinStockLevel / 7)) // Rough estimate
		}

		// Calculate estimated reorder cost
		estimatedReorders += reorderQuantity * item.CostPerUnit

		// Get last delivery date for this item
		var lastDeliveryDate *time.Time
		for _, delivery := range recentDeliveries {
			if delivery.InventoryItemID == item.ID {
				if lastDeliveryDate == nil || delivery.DeliveryDate.After(*lastDeliveryDate) {
					lastDeliveryDate = &delivery.DeliveryDate
				}
			}
		}

		// Get category name
		categoryName := ""
		if item.CategoryID != nil {
			category, err := h.service.GetCategory(*item.CategoryID)
			if err == nil {
				categoryName = category.Name
			}
		}

		supplyChainData.Items = append(supplyChainData.Items, email.SupplyChainItemData{
			ID:                item.ID,
			Name:              item.Name,
			Category:          categoryName,
			CurrentStock:      currentStock,
			MinStock:          item.MinStockLevel,
			MaxStock:          item.MaxStockLevel,
			Unit:              item.Unit,
			Status:            status,
			PreferredVendor:   item.PreferredVendor,
			CostPerUnit:       item.CostPerUnit,
			ReorderQuantity:   reorderQuantity,
			DaysUntilStockout: daysUntilStockout,
			LastDeliveryDate:  lastDeliveryDate,
		})
	}

	supplyChainData.TotalValue = totalValue
	supplyChainData.LowStockItems = lowStockCount
	supplyChainData.OutOfStockItems = outOfStockCount
	supplyChainData.CriticalItems = criticalCount
	supplyChainData.EstimatedReorders = estimatedReorders

	return supplyChainData, nil
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

// EmailScheduleRequest represents the request body for creating/updating email schedules
type EmailScheduleRequest struct {
	EmailType  string `json:"email_type" binding:"required"`  // "weekly_stock_report", "weekly_supply_chain_report", "low_stock_alert"
	Frequency  string `json:"frequency" binding:"required"`   // "weekly", "daily", "monthly"
	DayOfWeek  *int   `json:"day_of_week"`                    // 0-6 (Sunday-Saturday) for weekly
	DayOfMonth *int   `json:"day_of_month"`                   // 1-31 for monthly
	TimeOfDay  string `json:"time_of_day" binding:"required"` // "09:00", "18:30"
	IsActive   bool   `json:"is_active"`                      // Whether the schedule is active
}

// GetEmailSchedules handles GET /api/accounts/:accountID/email-schedules
// Returns all email schedules for an account
func GetEmailSchedules(c *gin.Context) {
	// Get account ID from URL parameter
	accountIDStr := c.Param("accountID")
	accountID, err := strconv.Atoi(accountIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID"})
		return
	}

	// Get user from context (set by auth middleware)
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Validate account access
	service := c.MustGet("service").(*database.Service)
	if err := service.ValidateAccountAccess(accountID, user.(*models.User).AccountID); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// Get email schedules for the account
	schedules, err := service.GetEmailSchedulesByAccount(accountID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get email schedules"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"account_id": accountID,
		"schedules":  schedules,
	})
}

// GetEmailSchedule handles GET /api/accounts/:accountID/email-schedules/:emailType
// Returns a specific email schedule for an account
func GetEmailSchedule(c *gin.Context) {
	// Get account ID from URL parameter
	accountIDStr := c.Param("accountID")
	accountID, err := strconv.Atoi(accountIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID"})
		return
	}

	// Get email type from URL parameter
	emailType := c.Param("emailType")

	// Get user from context (set by auth middleware)
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Validate account access
	service := c.MustGet("service").(*database.Service)
	if err := service.ValidateAccountAccess(accountID, user.(*models.User).AccountID); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// Get specific email schedule
	schedule, err := service.GetEmailScheduleByAccountAndType(accountID, emailType)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Email schedule not found"})
		return
	}

	c.JSON(http.StatusOK, schedule)
}

// CreateEmailSchedule handles POST /api/accounts/:accountID/email-schedules
// Creates a new email schedule for an account
func CreateEmailSchedule(c *gin.Context) {
	// Get account ID from URL parameter
	accountIDStr := c.Param("accountID")
	accountID, err := strconv.Atoi(accountIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID"})
		return
	}

	// Get user from context (set by auth middleware)
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Validate account access
	service := c.MustGet("service").(*database.Service)
	if err := service.ValidateAccountAccess(accountID, user.(*models.User).AccountID); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// Parse request body
	var req EmailScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Validate email type
	if !isValidEmailType(req.EmailType) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email type"})
		return
	}

	// Check if schedule already exists
	existingSchedule, err := service.GetEmailScheduleByAccountAndType(accountID, req.EmailType)
	if err == nil && existingSchedule != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email schedule already exists for this type"})
		return
	}

	// Create new email schedule
	schedule := &models.EmailSchedule{
		AccountID:  accountID,
		EmailType:  req.EmailType,
		Frequency:  req.Frequency,
		DayOfWeek:  req.DayOfWeek,
		DayOfMonth: req.DayOfMonth,
		TimeOfDay:  req.TimeOfDay,
		IsActive:   req.IsActive,
	}

	if err := service.CreateEmailSchedule(schedule); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create email schedule"})
		return
	}

	c.JSON(http.StatusCreated, schedule)
}

// UpdateEmailSchedule handles PUT /api/accounts/:accountID/email-schedules/:emailType
// Updates an existing email schedule for an account
func UpdateEmailSchedule(c *gin.Context) {
	// Get account ID from URL parameter
	accountIDStr := c.Param("accountID")
	accountID, err := strconv.Atoi(accountIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID"})
		return
	}

	// Get email type from URL parameter
	emailType := c.Param("emailType")

	// Get user from context (set by auth middleware)
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Validate account access
	service := c.MustGet("service").(*database.Service)
	if err := service.ValidateAccountAccess(accountID, user.(*models.User).AccountID); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// Get existing schedule
	schedule, err := service.GetEmailScheduleByAccountAndType(accountID, emailType)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Email schedule not found"})
		return
	}

	// Parse request body
	var req EmailScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Validate email type
	if !isValidEmailType(req.EmailType) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email type"})
		return
	}

	// Update schedule fields
	schedule.EmailType = req.EmailType
	schedule.Frequency = req.Frequency
	schedule.DayOfWeek = req.DayOfWeek
	schedule.DayOfMonth = req.DayOfMonth
	schedule.TimeOfDay = req.TimeOfDay
	schedule.IsActive = req.IsActive

	if err := service.UpdateEmailSchedule(schedule); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update email schedule"})
		return
	}

	c.JSON(http.StatusOK, schedule)
}

// DeleteEmailSchedule handles DELETE /api/accounts/:accountID/email-schedules/:emailType
// Deletes an email schedule for an account (effectively turns it off)
func DeleteEmailSchedule(c *gin.Context) {
	// Get account ID from URL parameter
	accountIDStr := c.Param("accountID")
	accountID, err := strconv.Atoi(accountIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID"})
		return
	}

	// Get email type from URL parameter
	emailType := c.Param("emailType")

	// Get user from context (set by auth middleware)
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Validate account access
	service := c.MustGet("service").(*database.Service)
	if err := service.ValidateAccountAccess(accountID, user.(*models.User).AccountID); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// Get existing schedule
	schedule, err := service.GetEmailScheduleByAccountAndType(accountID, emailType)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Email schedule not found"})
		return
	}

	// Delete the schedule
	if err := service.DeleteEmailSchedule(schedule.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete email schedule"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email schedule deleted successfully"})
}

// ToggleEmailSchedule handles PATCH /api/accounts/:accountID/email-schedules/:emailType/toggle
// Toggles the active status of an email schedule (turn on/off)
func ToggleEmailSchedule(c *gin.Context) {
	// Get account ID from URL parameter
	accountIDStr := c.Param("accountID")
	accountID, err := strconv.Atoi(accountIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID"})
		return
	}

	// Get email type from URL parameter
	emailType := c.Param("emailType")

	// Get user from context (set by auth middleware)
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Validate account access
	service := c.MustGet("service").(*database.Service)
	if err := service.ValidateAccountAccess(accountID, user.(*models.User).AccountID); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// Get existing schedule
	schedule, err := service.GetEmailScheduleByAccountAndType(accountID, emailType)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Email schedule not found"})
		return
	}

	// Toggle the active status
	schedule.IsActive = !schedule.IsActive

	if err := service.UpdateEmailSchedule(schedule); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update email schedule"})
		return
	}

	status := "enabled"
	if !schedule.IsActive {
		status = "disabled"
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Email schedule " + status + " successfully",
		"is_active":  schedule.IsActive,
		"email_type": emailType,
		"account_id": accountID,
	})
}

// isValidEmailType validates that the email type is supported
func isValidEmailType(emailType string) bool {
	validTypes := []string{
		models.EmailTypeWeeklyReport,
		models.EmailTypeWeeklySupplyChain,
		models.EmailTypeLowStockAlert,
	}

	for _, validType := range validTypes {
		if emailType == validType {
			return true
		}
	}
	return false
}
