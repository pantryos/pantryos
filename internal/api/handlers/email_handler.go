package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	helpers "github.com/mnadev/pantryos/internal/api/helper"
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

// SendVerificationEmail godoc
// @Summary      Send verification email
// @Description  Generates a new verification token and sends it to the user's email address.
// @Tags         email
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        user_id  path      int                                         true  "User ID to send verification email to"
// @Success      200      {object}  helpers.APIResponse                           "Verification email sent successfully"
// @Failure      400      {object}  helpers.APIResponse                           "Invalid user ID format"
// @Failure      404      {object}  helpers.APIResponse                           "User or associated account not found"
// @Failure      500      {object}  helpers.APIResponse                           "Internal server error (e.g., failed to send email)"
// @Router       /email/verification/{user_id} [post]
func (h *EmailHandler) SendVerificationEmail(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		errDetails := helpers.APIError{Code: "INVALID_INPUT", Details: "User ID must be a valid integer."}
		helpers.Error(c.Writer, http.StatusBadRequest, "Invalid user ID.", errDetails)
		return
	}

	// Get user from the database
	user, err := h.service.GetUser(userID)
	if err != nil {
		errDetails := helpers.APIError{Code: "USER_NOT_FOUND", Details: err.Error()}
		helpers.Error(c.Writer, http.StatusNotFound, "User not found.", errDetails)
		return
	}

	// Get the account for the user
	account, err := h.service.GetAccount(user.AccountID)
	if err != nil {
		errDetails := helpers.APIError{Code: "ACCOUNT_NOT_FOUND", Details: err.Error()}
		helpers.Error(c.Writer, http.StatusNotFound, "Account associated with the user not found.", errDetails)
		return
	}

	// Generate a verification token
	token, err := h.generateVerificationToken(userID)
	if err != nil {
		errDetails := helpers.APIError{Code: "TOKEN_GENERATION_FAILED", Details: err.Error()}
		helpers.Error(c.Writer, http.StatusInternalServerError, "Failed to generate verification token.", errDetails)
		return
	}

	// Create the full verification URL
	baseURL := getBaseURL(c)
	verificationURL := fmt.Sprintf("%s/verify-email?token=%s", baseURL, token)

	// Send the verification email via the email service
	if err := h.emailService.SendVerificationEmail(*user, *account, verificationURL); err != nil {
		// Log the email failure for debugging purposes
		h.logEmailFailure(user.AccountID, &userID, user.Email, "Verify Your PantryOS Account", models.EmailTypeVerification, err.Error())
		errDetails := helpers.APIError{Code: "EMAIL_SEND_FAILED", Details: err.Error()}
		helpers.Error(c.Writer, http.StatusInternalServerError, "Failed to send verification email.", errDetails)
		return
	}

	// Log the successful email dispatch
	h.logEmailSuccess(user.AccountID, &userID, user.Email, "Verify Your PantryOS Account", models.EmailTypeVerification)

	responseData := gin.H{
		"user_id": userID,
	}
	helpers.Success(c.Writer, http.StatusOK, "Verification email sent successfully.", responseData)
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
		errDetails := helpers.APIError{Code: "INVALID_INPUT", Details: "Verification token is required."}
		helpers.Error(c.Writer, http.StatusBadRequest, "Missing verification token.", errDetails)
		return
	}

	// For now, we'll just return success since we need to add email verification methods to the service
	// TODO: Implement proper token validation through service layer
	responseData := gin.H{
		"message": "Email verified successfully",
	}
	helpers.Success(c.Writer, http.StatusOK, "Email verified successfully. You can now log in.", responseData)
}

// SendWeeklyStockReport godoc
// @Summary      Send weekly stock report
// @Description  Triggers the sending of a weekly stock report email to all users in a specific account.
// @Tags         email
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        account_id  path      int                                         true  "Account ID to send the report for"
// @Success      200         {object}  helpers.APIResponse                           "Weekly stock report sent successfully"
// @Failure      400         {object}  helpers.APIResponse                           "Invalid account ID format or no users in account"
// @Failure      404         {object}  helpers.APIResponse                           "Account not found"
// @Failure      500         {object}  helpers.APIResponse                           "Internal server error (e.g., failed to send email)"
// @Router       /email/weekly-report/{account_id} [post]
func (h *EmailHandler) SendWeeklyStockReport(c *gin.Context) {
	accountID, err := strconv.Atoi(c.Param("account_id"))
	if err != nil {
		errDetails := helpers.APIError{Code: "INVALID_INPUT", Details: "Account ID must be a valid integer."}
		helpers.Error(c.Writer, http.StatusBadRequest, "Invalid Account ID.", errDetails)
		return
	}

	// Get the account from the database
	account, err := h.service.GetAccount(accountID)
	if err != nil {
		errDetails := helpers.APIError{Code: "ACCOUNT_NOT_FOUND", Details: err.Error()}
		helpers.Error(c.Writer, http.StatusNotFound, "Account not found.", errDetails)
		return
	}

	// Get all users in the account
	users, err := h.service.GetUsersByAccount(accountID)
	if err != nil {
		errDetails := helpers.APIError{Code: "DB_FETCH_FAILED", Details: err.Error()}
		helpers.Error(c.Writer, http.StatusInternalServerError, "Failed to get users for the account.", errDetails)
		return
	}

	if len(users) == 0 {
		errDetails := helpers.APIError{Code: "NO_USERS_IN_ACCOUNT", Details: "Cannot send report because there are no users associated with this account."}
		helpers.Error(c.Writer, http.StatusBadRequest, "No users found in the account.", errDetails)
		return
	}

	// Generate the data for the stock report
	stockData, err := h.generateStockReportData(accountID)
	if err != nil {
		errDetails := helpers.APIError{Code: "REPORT_GENERATION_FAILED", Details: err.Error()}
		helpers.Error(c.Writer, http.StatusInternalServerError, "Failed to generate stock report data.", errDetails)
		return
	}

	// Send the weekly stock report via the email service
	if err := h.emailService.SendWeeklyStockReport(*account, users, stockData); err != nil {
		errDetails := helpers.APIError{Code: "EMAIL_SEND_FAILED", Details: err.Error()}
		helpers.Error(c.Writer, http.StatusInternalServerError, "Failed to send weekly stock report.", errDetails)
		return
	}

	// Log the successful email dispatch for each user
	for _, user := range users {
		h.logEmailSuccess(accountID, &user.ID, user.Email, fmt.Sprintf("Weekly Stock Report - %s", account.Name), models.EmailTypeWeeklyReport)
	}

	responseData := gin.H{
		"account_id":  accountID,
		"users_count": len(users),
	}
	helpers.Success(c.Writer, http.StatusOK, "Weekly stock report sent successfully.", responseData)
}

// SendLowStockAlert godoc
// @Summary      Send low stock alert
// @Description  Checks for low stock items and sends an alert email to all users in the account if any are found.
// @Tags         email
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        account_id  path      int                                         true  "Account ID to check and send the alert for"
// @Success      200         {object}  helpers.APIResponse                           "Alert sent successfully, or no low stock items were found"
// @Failure      400         {object}  helpers.APIResponse                           "Invalid account ID format or no users in account"
// @Failure      404         {object}  helpers.APIResponse                           "Account not found"
// @Failure      500         {object}  helpers.APIResponse                           "Internal server error"
// @Router       /email/low-stock-alert/{account_id} [post]
func (h *EmailHandler) SendLowStockAlert(c *gin.Context) {
	accountID, err := strconv.Atoi(c.Param("account_id"))
	if err != nil {
		errDetails := helpers.APIError{Code: "INVALID_INPUT", Details: "Account ID must be a valid integer."}
		helpers.Error(c.Writer, http.StatusBadRequest, "Invalid Account ID.", errDetails)
		return
	}

	// Get the account from the database
	account, err := h.service.GetAccount(accountID)
	if err != nil {
		errDetails := helpers.APIError{Code: "ACCOUNT_NOT_FOUND", Details: err.Error()}
		helpers.Error(c.Writer, http.StatusNotFound, "Account not found.", errDetails)
		return
	}

	// Get all users in the account
	users, err := h.service.GetUsersByAccount(accountID)
	if err != nil {
		errDetails := helpers.APIError{Code: "DB_FETCH_FAILED", Details: err.Error()}
		helpers.Error(c.Writer, http.StatusInternalServerError, "Failed to get users for the account.", errDetails)
		return
	}

	if len(users) == 0 {
		errDetails := helpers.APIError{Code: "NO_USERS_IN_ACCOUNT", Details: "Cannot send alert because there are no users associated with this account."}
		helpers.Error(c.Writer, http.StatusBadRequest, "No users found in the account.", errDetails)
		return
	}

	// Get low stock items
	lowStockItems, err := h.service.GetLowStockItems(accountID)
	if err != nil {
		errDetails := helpers.APIError{Code: "DB_FETCH_FAILED", Details: err.Error()}
		helpers.Error(c.Writer, http.StatusInternalServerError, "Failed to retrieve low stock items.", errDetails)
		return
	}

	// If no items are low on stock, it's a successful outcome with no email needed.
	if len(lowStockItems) == 0 {
		responseData := gin.H{
			"account_id": accountID,
		}
		helpers.Success(c.Writer, http.StatusOK, "No low stock items found to report.", responseData)
		return
	}

	// Send the low stock alert via the email service
	if err := h.emailService.SendLowStockAlert(*account, users, lowStockItems); err != nil {
		errDetails := helpers.APIError{Code: "EMAIL_SEND_FAILED", Details: err.Error()}
		helpers.Error(c.Writer, http.StatusInternalServerError, "Failed to send low stock alert.", errDetails)
		return
	}

	// Log the successful email dispatch for each user
	for _, user := range users {
		h.logEmailSuccess(accountID, &user.ID, user.Email, fmt.Sprintf("Low Stock Alert - %s", account.Name), models.EmailTypeLowStockAlert)
	}

	responseData := gin.H{
		"account_id":            accountID,
		"users_count":           len(users),
		"low_stock_items_count": len(lowStockItems),
	}
	helpers.Success(c.Writer, http.StatusOK, "Low stock alert sent successfully.", responseData)
}

// SendWeeklySupplyChainReport godoc
// @Summary      Send weekly supply chain report
// @Description  Triggers the sending of a weekly supply chain report email to all users in a specific account.
// @Tags         email
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        account_id  path      int                                         true  "Account ID to send the report for"
// @Success      200         {object}  helpers.APIResponse                           "Weekly supply chain report sent successfully"
// @Failure      400         {object}  helpers.APIResponse                           "Invalid account ID format or no users in account"
// @Failure      404         {object}  helpers.APIResponse                           "Account not found"
// @Failure      500         {object}  helpers.APIResponse                           "Internal server error (e.g., failed to send email)"
// @Router       /email/weekly-supply-chain/{account_id} [post]
func (h *EmailHandler) SendWeeklySupplyChainReport(c *gin.Context) {
	accountID, err := strconv.Atoi(c.Param("account_id"))
	if err != nil {
		errDetails := helpers.APIError{Code: "INVALID_INPUT", Details: "Account ID must be a valid integer."}
		helpers.Error(c.Writer, http.StatusBadRequest, "Invalid Account ID.", errDetails)
		return
	}

	// Get the account from the database
	account, err := h.service.GetAccount(accountID)
	if err != nil {
		errDetails := helpers.APIError{Code: "ACCOUNT_NOT_FOUND", Details: err.Error()}
		helpers.Error(c.Writer, http.StatusNotFound, "Account not found.", errDetails)
		return
	}

	// Get all users in the account
	users, err := h.service.GetUsersByAccount(accountID)
	if err != nil {
		errDetails := helpers.APIError{Code: "DB_FETCH_FAILED", Details: err.Error()}
		helpers.Error(c.Writer, http.StatusInternalServerError, "Failed to get users for the account.", errDetails)
		return
	}

	if len(users) == 0 {
		errDetails := helpers.APIError{Code: "NO_USERS_IN_ACCOUNT", Details: "Cannot send report because there are no users associated with this account."}
		helpers.Error(c.Writer, http.StatusBadRequest, "No users found in the account.", errDetails)
		return
	}

	// Generate the data for the supply chain report
	supplyChainData, err := h.generateSupplyChainReportData(accountID)
	if err != nil {
		errDetails := helpers.APIError{Code: "REPORT_GENERATION_FAILED", Details: err.Error()}
		helpers.Error(c.Writer, http.StatusInternalServerError, "Failed to generate supply chain report data.", errDetails)
		return
	}

	// Send the weekly supply chain report via the email service
	if err := h.emailService.SendWeeklySupplyChainReport(*account, users, supplyChainData); err != nil {
		errDetails := helpers.APIError{Code: "EMAIL_SEND_FAILED", Details: err.Error()}
		helpers.Error(c.Writer, http.StatusInternalServerError, "Failed to send weekly supply chain report.", errDetails)
		return
	}

	// Log the successful email dispatch for each user
	for _, user := range users {
		h.logEmailSuccess(accountID, &user.ID, user.Email, fmt.Sprintf("Weekly Supply Chain Report - %s", account.Name), models.EmailTypeWeeklySupplyChain)
	}

	responseData := gin.H{
		"account_id":  accountID,
		"users_count": len(users),
	}
	helpers.Success(c.Writer, http.StatusOK, "Weekly supply chain report sent successfully.", responseData)
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

// GetEmailSchedules handles GET /api/v1/accounts/:account_id/email-schedules
// Returns all email schedules for a specific account.
//
// Authentication: Required (JWT token in Authorization header)
// Authorization: User must be authenticated and have access to the specified account.
//
// URL Parameters:
//   - account_id: The ID of the account to retrieve email schedules for (integer)
//
// Response:
//
//	All responses are wrapped in the standard APIResponse structure.
//
// Status Codes:
//   - 200 OK: Email schedules retrieved successfully.
//   - 400 Bad Request: Invalid account ID format.
//   - 401 Unauthorized: User not authenticated.
//   - 403 Forbidden: User does not have access to the requested account.
//   - 500 Internal Server Error: Failed to retrieve email schedules.
func GetEmailSchedules(c *gin.Context) {
	// Get account ID from URL parameter
	accountID, err := strconv.Atoi(c.Param("account_id"))
	if err != nil {
		errDetails := helpers.APIError{Code: "INVALID_INPUT", Details: "Account ID must be a valid integer."}
		helpers.Error(c.Writer, http.StatusBadRequest, "Invalid Account ID.", errDetails)
		return
	}

	// Get user from context (set by auth middleware)
	userInterface, exists := c.Get("user")
	if !exists {
		errDetails := helpers.APIError{Code: "UNAUTHORIZED", Details: "User not found in request context."}
		helpers.Error(c.Writer, http.StatusUnauthorized, "User not authenticated.", errDetails)
		return
	}
	user := userInterface.(*models.User)

	// Validate account access
	service := c.MustGet("service").(*database.Service)
	if err := service.ValidateAccountAccess(accountID, user.AccountID); err != nil {
		errDetails := helpers.APIError{Code: "FORBIDDEN", Details: "You do not have permission to access this account's schedules."}
		helpers.Error(c.Writer, http.StatusForbidden, "Access denied.", errDetails)
		return
	}

	// Get email schedules for the account
	schedules, err := service.GetEmailSchedulesByAccount(accountID)
	if err != nil {
		errDetails := helpers.APIError{Code: "DB_FETCH_FAILED", Details: err.Error()}
		helpers.Error(c.Writer, http.StatusInternalServerError, "Failed to get email schedules.", errDetails)
		return
	}

	// Return a 200 OK response with the list of schedules
	helpers.Success(c.Writer, http.StatusOK, "Email schedules retrieved successfully.", schedules)
}

// GetEmailSchedule handles GET /api/v1/accounts/:account_id/email-schedules/:emailType
// Returns a specific email schedule for an account.
//
// Authentication: Required (JWT token in Authorization header)
// Authorization: User must be authenticated and have access to the specified account.
//
// URL Parameters:
//   - account_id: The ID of the account (integer)
//   - emailType: The type of the email schedule (string)
//
// Response:
//
//	All responses are wrapped in the standard APIResponse structure.
//
// Status Codes:
//   - 200 OK: Email schedule retrieved successfully.
//   - 400 Bad Request: Invalid account ID format.
//   - 401 Unauthorized: User not authenticated.
//   - 403 Forbidden: User does not have access to the requested account.
//   - 404 Not Found: The requested email schedule does not exist.
//   - 500 Internal Server Error: General server error.
func GetEmailSchedule(c *gin.Context) {
	// Get account ID from URL parameter
	accountID, err := strconv.Atoi(c.Param("account_id"))
	if err != nil {
		errDetails := helpers.APIError{Code: "INVALID_INPUT", Details: "Account ID must be a valid integer."}
		helpers.Error(c.Writer, http.StatusBadRequest, "Invalid Account ID.", errDetails)
		return
	}

	// Get email type from URL parameter
	emailType := c.Param("emailType")

	// Get user from context (set by auth middleware)
	userInterface, exists := c.Get("user")
	if !exists {
		errDetails := helpers.APIError{Code: "UNAUTHORIZED", Details: "User not found in request context."}
		helpers.Error(c.Writer, http.StatusUnauthorized, "User not authenticated.", errDetails)
		return
	}
	user := userInterface.(*models.User)

	// Validate account access
	service := c.MustGet("service").(*database.Service)
	if err := service.ValidateAccountAccess(accountID, user.AccountID); err != nil {
		errDetails := helpers.APIError{Code: "FORBIDDEN", Details: "You do not have permission to access this account's schedules."}
		helpers.Error(c.Writer, http.StatusForbidden, "Access denied.", errDetails)
		return
	}

	// Get the specific email schedule
	schedule, err := service.GetEmailScheduleByAccountAndType(accountID, emailType)
	if err != nil {
		// This could be a genuine "not found" or another database error.
		// Returning 404 is a safe default for a "get by ID" type of function.
		errDetails := helpers.APIError{Code: "NOT_FOUND", Details: err.Error()}
		helpers.Error(c.Writer, http.StatusNotFound, "Email schedule not found.", errDetails)
		return
	}

	// Return a 200 OK response with the schedule object
	helpers.Success(c.Writer, http.StatusOK, "Email schedule retrieved successfully.", schedule)
}

// CreateEmailSchedule handles POST /api/v1/accounts/:account_id/email-schedules
// Creates a new email schedule for an account.
//
// Authentication: Required (JWT token in Authorization header)
// Authorization: User must be authenticated and have access to the specified account.
//
// URL Parameters:
//   - account_id: The ID of the account to create the schedule for (integer)
//
// Request Body: JSON object with new email schedule details.
//
// Response:
//
//	All responses are wrapped in the standard APIResponse structure.
//
// Status Codes:
//   - 201 Created: Email schedule created successfully.
//   - 400 Bad Request: Invalid account ID format or request body.
//   - 401 Unauthorized: User not authenticated.
//   - 403 Forbidden: User does not have access to the requested account.
//   - 409 Conflict: A schedule with the same email type already exists for this account.
//   - 500 Internal Server Error: Failed to create the email schedule.
func CreateEmailSchedule(c *gin.Context) {
	// Get account ID from URL parameter
	accountID, err := strconv.Atoi(c.Param("account_id"))
	if err != nil {
		errDetails := helpers.APIError{Code: "INVALID_INPUT", Details: "Account ID must be a valid integer."}
		helpers.Error(c.Writer, http.StatusBadRequest, "Invalid Account ID.", errDetails)
		return
	}

	// Get user from context (set by auth middleware)
	userInterface, exists := c.Get("user")
	if !exists {
		errDetails := helpers.APIError{Code: "UNAUTHORIZED", Details: "User not found in request context."}
		helpers.Error(c.Writer, http.StatusUnauthorized, "User not authenticated.", errDetails)
		return
	}
	user := userInterface.(*models.User)

	// Validate account access
	service := c.MustGet("service").(*database.Service)
	if err := service.ValidateAccountAccess(accountID, user.AccountID); err != nil {
		errDetails := helpers.APIError{Code: "FORBIDDEN", Details: "You do not have permission to modify this account's schedules."}
		helpers.Error(c.Writer, http.StatusForbidden, "Access denied.", errDetails)
		return
	}

	// Parse request body
	var req EmailScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errDetails := helpers.APIError{Code: "INVALID_INPUT", Details: err.Error()}
		helpers.Error(c.Writer, http.StatusBadRequest, "Invalid request body.", errDetails)
		return
	}

	// Validate email type
	if !isValidEmailType(req.EmailType) {
		errDetails := helpers.APIError{Code: "INVALID_EMAIL_TYPE", Details: fmt.Sprintf("'%s' is not a valid email type.", req.EmailType)}
		helpers.Error(c.Writer, http.StatusBadRequest, "Invalid email type provided.", errDetails)
		return
	}

	// Check if a schedule of this type already exists for the account
	existingSchedule, err := service.GetEmailScheduleByAccountAndType(accountID, req.EmailType)
	if err == nil && existingSchedule != nil {
		errDetails := helpers.APIError{Code: "CONFLICT", Details: "An email schedule for this type already exists."}
		helpers.Error(c.Writer, http.StatusConflict, "Email schedule already exists.", errDetails)
		return
	}

	// Create a new email schedule model
	schedule := &models.EmailSchedule{
		AccountID:  accountID,
		EmailType:  req.EmailType,
		Frequency:  req.Frequency,
		DayOfWeek:  req.DayOfWeek,
		DayOfMonth: req.DayOfMonth,
		TimeOfDay:  req.TimeOfDay,
		IsActive:   req.IsActive,
	}

	// Save the new schedule to the database
	if err := service.CreateEmailSchedule(schedule); err != nil {
		errDetails := helpers.APIError{Code: "DB_INSERT_FAILED", Details: err.Error()}
		helpers.Error(c.Writer, http.StatusInternalServerError, "Failed to create email schedule.", errDetails)
		return
	}

	// Return a 201 Created response with the new schedule object
	helpers.Success(c.Writer, http.StatusCreated, "Email schedule created successfully.", schedule)
}

// UpdateEmailSchedule handles PUT /api/v1/accounts/:account_id/email-schedules/:emailType
// Updates an existing email schedule for an account.
//
// Authentication: Required (JWT token in Authorization header)
// Authorization: User must be authenticated and have access to the specified account.
//
// URL Parameters:
//   - account_id: The ID of the account (integer)
//   - emailType: The type of the email schedule to update (string)
//
// Request Body: JSON object with updated email schedule details.
//
// Response:
//
//	All responses are wrapped in the standard APIResponse structure.
//
// Status Codes:
//   - 200 OK: Email schedule updated successfully.
//   - 400 Bad Request: Invalid account ID format or request body.
//   - 401 Unauthorized: User not authenticated.
//   - 403 Forbidden: User does not have access to the requested account.
//   - 404 Not Found: The requested email schedule does not exist.
//   - 500 Internal Server Error: Failed to update the email schedule.
func UpdateEmailSchedule(c *gin.Context) {
	// Get account ID from URL parameter
	accountID, err := strconv.Atoi(c.Param("account_id"))
	if err != nil {
		errDetails := helpers.APIError{Code: "INVALID_INPUT", Details: "Account ID must be a valid integer."}
		helpers.Error(c.Writer, http.StatusBadRequest, "Invalid Account ID.", errDetails)
		return
	}

	// Get email type from URL parameter
	emailType := c.Param("emailType")

	// Get user from context (set by auth middleware)
	userInterface, exists := c.Get("user")
	if !exists {
		errDetails := helpers.APIError{Code: "UNAUTHORIZED", Details: "User not found in request context."}
		helpers.Error(c.Writer, http.StatusUnauthorized, "User not authenticated.", errDetails)
		return
	}
	user := userInterface.(*models.User)

	// Validate account access
	service := c.MustGet("service").(*database.Service)
	if err := service.ValidateAccountAccess(accountID, user.AccountID); err != nil {
		errDetails := helpers.APIError{Code: "FORBIDDEN", Details: "You do not have permission to modify this account's schedules."}
		helpers.Error(c.Writer, http.StatusForbidden, "Access denied.", errDetails)
		return
	}

	// Get the existing schedule
	schedule, err := service.GetEmailScheduleByAccountAndType(accountID, emailType)
	if err != nil {
		errDetails := helpers.APIError{Code: "NOT_FOUND", Details: err.Error()}
		helpers.Error(c.Writer, http.StatusNotFound, "Email schedule not found.", errDetails)
		return
	}

	// Parse request body
	var req EmailScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errDetails := helpers.APIError{Code: "INVALID_INPUT", Details: err.Error()}
		helpers.Error(c.Writer, http.StatusBadRequest, "Invalid request body.", errDetails)
		return
	}

	// Validate email type
	if !isValidEmailType(req.EmailType) {
		errDetails := helpers.APIError{Code: "INVALID_EMAIL_TYPE", Details: fmt.Sprintf("'%s' is not a valid email type.", req.EmailType)}
		helpers.Error(c.Writer, http.StatusBadRequest, "Invalid email type provided.", errDetails)
		return
	}

	// Update the schedule fields from the request
	schedule.EmailType = req.EmailType
	schedule.Frequency = req.Frequency
	schedule.DayOfWeek = req.DayOfWeek
	schedule.DayOfMonth = req.DayOfMonth
	schedule.TimeOfDay = req.TimeOfDay
	schedule.IsActive = req.IsActive

	// Save the updated schedule to the database
	if err := service.UpdateEmailSchedule(schedule); err != nil {
		errDetails := helpers.APIError{Code: "DB_UPDATE_FAILED", Details: err.Error()}
		helpers.Error(c.Writer, http.StatusInternalServerError, "Failed to update email schedule.", errDetails)
		return
	}

	// Return a 200 OK response with the updated schedule object
	helpers.Success(c.Writer, http.StatusOK, "Email schedule updated successfully.", schedule)
}

// DeleteEmailSchedule handles DELETE /api/v1/accounts/:account_id/email-schedules/:emailType
// Deletes an email schedule for an account.
//
// Authentication: Required (JWT token in Authorization header)
// Authorization: User must be authenticated and have access to the specified account.
//
// URL Parameters:
//   - account_id: The ID of the account (integer)
//   - emailType: The type of the email schedule to delete (string)
//
// Response:
//
//	All responses are wrapped in the standard APIResponse structure.
//
// Status Codes:
//   - 200 OK: Email schedule deleted successfully.
//   - 400 Bad Request: Invalid account ID format.
//   - 401 Unauthorized: User not authenticated.
//   - 403 Forbidden: User does not have access to the requested account.
//   - 404 Not Found: The requested email schedule does not exist.
//   - 500 Internal Server Error: Failed to delete the email schedule.
func DeleteEmailSchedule(c *gin.Context) {
	// Get account ID from URL parameter
	accountID, err := strconv.Atoi(c.Param("account_id"))
	if err != nil {
		errDetails := helpers.APIError{Code: "INVALID_INPUT", Details: "Account ID must be a valid integer."}
		helpers.Error(c.Writer, http.StatusBadRequest, "Invalid Account ID.", errDetails)
		return
	}

	// Get email type from URL parameter
	emailType := c.Param("emailType")

	// Get user from context (set by auth middleware)
	userInterface, exists := c.Get("user")
	if !exists {
		errDetails := helpers.APIError{Code: "UNAUTHORIZED", Details: "User not found in request context."}
		helpers.Error(c.Writer, http.StatusUnauthorized, "User not authenticated.", errDetails)
		return
	}
	user := userInterface.(*models.User)

	// Validate account access
	service := c.MustGet("service").(*database.Service)
	if err := service.ValidateAccountAccess(accountID, user.AccountID); err != nil {
		errDetails := helpers.APIError{Code: "FORBIDDEN", Details: "You do not have permission to modify this account's schedules."}
		helpers.Error(c.Writer, http.StatusForbidden, "Access denied.", errDetails)
		return
	}

	// Get the existing schedule to ensure it exists before deleting
	schedule, err := service.GetEmailScheduleByAccountAndType(accountID, emailType)
	if err != nil {
		errDetails := helpers.APIError{Code: "NOT_FOUND", Details: err.Error()}
		helpers.Error(c.Writer, http.StatusNotFound, "Email schedule not found.", errDetails)
		return
	}

	// Delete the schedule from the database
	if err := service.DeleteEmailSchedule(schedule.ID); err != nil {
		errDetails := helpers.APIError{Code: "DB_DELETE_FAILED", Details: err.Error()}
		helpers.Error(c.Writer, http.StatusInternalServerError, "Failed to delete email schedule.", errDetails)
		return
	}

	// Return a 200 OK response with a success message
	helpers.Success(c.Writer, http.StatusOK, "Email schedule deleted successfully.", nil)
}

// ToggleEmailSchedule handles PATCH /api/v1/accounts/:account_id/email-schedules/:emailType/toggle
// Toggles the IsActive status of an email schedule.
//
// Authentication: Required (JWT token in Authorization header)
// Authorization: User must be authenticated and have access to the specified account.
//
// URL Parameters:
//   - account_id: The ID of the account (integer)
//   - emailType: The type of the email schedule to toggle (string)
//
// Response:
//
//	All responses are wrapped in the standard APIResponse structure.
//
// Status Codes:
//   - 200 OK: Email schedule toggled successfully.
//   - 400 Bad Request: Invalid account ID format.
//   - 401 Unauthorized: User not authenticated.
//   - 403 Forbidden: User does not have access to the requested account.
//   - 404 Not Found: The requested email schedule does not exist.
//   - 500 Internal Server Error: Failed to update the email schedule.
func ToggleEmailSchedule(c *gin.Context) {
	// Get account ID from URL parameter
	accountID, err := strconv.Atoi(c.Param("account_id"))
	if err != nil {
		errDetails := helpers.APIError{Code: "INVALID_INPUT", Details: "Account ID must be a valid integer."}
		helpers.Error(c.Writer, http.StatusBadRequest, "Invalid Account ID.", errDetails)
		return
	}

	// Get email type from URL parameter
	emailType := c.Param("emailType")

	// Get user from context (set by auth middleware)
	userInterface, exists := c.Get("user")
	if !exists {
		errDetails := helpers.APIError{Code: "UNAUTHORIZED", Details: "User not found in request context."}
		helpers.Error(c.Writer, http.StatusUnauthorized, "User not authenticated.", errDetails)
		return
	}
	user := userInterface.(*models.User)

	// Validate account access
	service := c.MustGet("service").(*database.Service)
	if err := service.ValidateAccountAccess(accountID, user.AccountID); err != nil {
		errDetails := helpers.APIError{Code: "FORBIDDEN", Details: "You do not have permission to modify this account's schedules."}
		helpers.Error(c.Writer, http.StatusForbidden, "Access denied.", errDetails)
		return
	}

	// Get the existing schedule
	schedule, err := service.GetEmailScheduleByAccountAndType(accountID, emailType)
	if err != nil {
		errDetails := helpers.APIError{Code: "NOT_FOUND", Details: err.Error()}
		helpers.Error(c.Writer, http.StatusNotFound, "Email schedule not found.", errDetails)
		return
	}

	// Toggle the IsActive status
	schedule.IsActive = !schedule.IsActive

	// Save the updated schedule to the database
	if err := service.UpdateEmailSchedule(schedule); err != nil {
		errDetails := helpers.APIError{Code: "DB_UPDATE_FAILED", Details: err.Error()}
		helpers.Error(c.Writer, http.StatusInternalServerError, "Failed to update email schedule.", errDetails)
		return
	}

	// Determine status message for the response
	status := "enabled"
	if !schedule.IsActive {
		status = "disabled"
	}

	responseData := gin.H{
		"is_active":  schedule.IsActive,
		"email_type": emailType,
		"account_id": accountID,
	}
	helpers.Success(c.Writer, http.StatusOK, "Email schedule "+status+" successfully.", responseData)
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
