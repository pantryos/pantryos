package scheduler

import (
	"fmt"
	"log"
	"time"

	"github.com/mnadev/pantryos/internal/database"
	"github.com/mnadev/pantryos/internal/email"
	"github.com/mnadev/pantryos/internal/models"
)

// Scheduler handles automated tasks like sending weekly stock reports
type Scheduler struct {
	db           *database.DB
	service      *database.Service
	emailService *email.EmailService
	stopChan     chan bool
}

// NewScheduler creates a new scheduler instance
func NewScheduler(db *database.DB) *Scheduler {
	return &Scheduler{
		db:           db,
		service:      database.NewService(db),
		emailService: email.NewEmailService(),
		stopChan:     make(chan bool),
	}
}

// Start starts the scheduler
func (s *Scheduler) Start() {
	log.Println("Starting email scheduler...")

	// Start weekly stock report scheduler
	go s.scheduleWeeklyStockReports()

	// Start weekly supply chain report scheduler
	go s.scheduleWeeklySupplyChainReports()

	// Start low stock alert scheduler
	go s.scheduleLowStockAlerts()

	log.Println("Email scheduler started successfully")
}

// Stop stops the scheduler
func (s *Scheduler) Stop() {
	log.Println("Stopping email scheduler...")
	close(s.stopChan)
}

// scheduleWeeklyStockReports schedules weekly stock report emails
func (s *Scheduler) scheduleWeeklyStockReports() {
	ticker := time.NewTicker(24 * time.Hour) // Check every 24 hours
	defer ticker.Stop()

	for {
		select {
		case <-s.stopChan:
			return
		case <-ticker.C:
			s.sendWeeklyStockReports()
		}
	}
}

// scheduleLowStockAlerts schedules low stock alert emails
func (s *Scheduler) scheduleLowStockAlerts() {
	ticker := time.NewTicker(12 * time.Hour) // Check every 12 hours
	defer ticker.Stop()

	for {
		select {
		case <-s.stopChan:
			return
		case <-ticker.C:
			s.sendLowStockAlerts()
		}
	}
}

// scheduleWeeklySupplyChainReports schedules weekly supply chain report emails
func (s *Scheduler) scheduleWeeklySupplyChainReports() {
	ticker := time.NewTicker(24 * time.Hour) // Check every 24 hours
	defer ticker.Stop()

	for {
		select {
		case <-s.stopChan:
			return
		case <-ticker.C:
			s.sendWeeklySupplyChainReports()
		}
	}
}

// sendWeeklyStockReports sends weekly stock reports to all accounts
func (s *Scheduler) sendWeeklyStockReports() {
	log.Println("Checking for weekly stock reports to send...")

	// Get all accounts
	accounts, err := s.service.GetAllAccounts()
	if err != nil {
		log.Printf("Failed to get accounts for weekly stock reports: %v", err)
		return
	}

	for _, account := range accounts {
		// Check if it's time to send weekly report for this account
		if s.shouldSendWeeklyReport(account.ID) {
			s.sendWeeklyStockReportForAccount(account)
		}
	}
}

// sendLowStockAlerts sends low stock alerts to all accounts
func (s *Scheduler) sendLowStockAlerts() {
	log.Println("Checking for low stock alerts to send...")

	// Get all accounts
	accounts, err := s.service.GetAllAccounts()
	if err != nil {
		log.Printf("Failed to get accounts for low stock alerts: %v", err)
		return
	}

	for _, account := range accounts {
		s.sendLowStockAlertForAccount(account)
	}
}

// sendWeeklySupplyChainReports sends weekly supply chain reports to all accounts
func (s *Scheduler) sendWeeklySupplyChainReports() {
	log.Println("Checking for weekly supply chain reports to send...")

	// Get all accounts
	accounts, err := s.service.GetAllAccounts()
	if err != nil {
		log.Printf("Failed to get accounts for weekly supply chain reports: %v", err)
		return
	}

	for _, account := range accounts {
		// Check if it's time to send weekly supply chain report for this account
		if s.shouldSendWeeklySupplyChainReport(account.ID) {
			s.sendWeeklySupplyChainReportForAccount(account)
		}
	}
}

// shouldSendWeeklyReport checks if it's time to send a weekly report for an account
func (s *Scheduler) shouldSendWeeklyReport(accountID int) bool {
	// For now, we'll send weekly reports every Monday at 9 AM
	// TODO: Make this configurable per account
	now := time.Now()

	// Check if it's Monday and between 9-10 AM
	if now.Weekday() == time.Monday && now.Hour() == 9 {
		return true
	}

	return false
}

// shouldSendWeeklySupplyChainReport checks if it's time to send a weekly supply chain report for an account
func (s *Scheduler) shouldSendWeeklySupplyChainReport(accountID int) bool {
	// For now, we'll send weekly supply chain reports every Tuesday at 9 AM
	// TODO: Make this configurable per account
	now := time.Now()

	// Check if it's Tuesday and between 9-10 AM
	if now.Weekday() == time.Tuesday && now.Hour() == 9 {
		return true
	}

	return false
}

// sendWeeklyStockReportForAccount sends a weekly stock report for a specific account
func (s *Scheduler) sendWeeklyStockReportForAccount(account models.Account) {
	log.Printf("Sending weekly stock report for account: %s", account.Name)

	// Get all users in the account
	users, err := s.service.GetUsersByAccount(account.ID)
	if err != nil {
		log.Printf("Failed to get users for account %d: %v", account.ID, err)
		return
	}

	if len(users) == 0 {
		log.Printf("No users found for account %d", account.ID)
		return
	}

	// Generate stock report data
	stockData, err := s.generateStockReportData(account.ID)
	if err != nil {
		log.Printf("Failed to generate stock report for account %d: %v", account.ID, err)
		return
	}

	// Send weekly stock report
	if err := s.emailService.SendWeeklyStockReport(account, users, stockData); err != nil {
		log.Printf("Failed to send weekly stock report for account %d: %v", account.ID, err)
		return
	}

	// Log successful email sending for each user
	for _, user := range users {
		s.logEmailSuccess(account.ID, &user.ID, user.Email, fmt.Sprintf("Weekly Stock Report - %s", account.Name), models.EmailTypeWeeklyReport)
	}

	log.Printf("Successfully sent weekly stock report for account: %s", account.Name)
}

// sendLowStockAlertForAccount sends a low stock alert for a specific account
func (s *Scheduler) sendLowStockAlertForAccount(account models.Account) {
	// Get low stock items
	lowStockItems, err := s.service.GetLowStockItems(account.ID)
	if err != nil {
		log.Printf("Failed to get low stock items for account %d: %v", account.ID, err)
		return
	}

	if len(lowStockItems) == 0 {
		return // No low stock items
	}

	// Get all users in the account
	users, err := s.service.GetUsersByAccount(account.ID)
	if err != nil {
		log.Printf("Failed to get users for account %d: %v", account.ID, err)
		return
	}

	if len(users) == 0 {
		return
	}

	// Send low stock alert
	if err := s.emailService.SendLowStockAlert(account, users, lowStockItems); err != nil {
		log.Printf("Failed to send low stock alert for account %d: %v", account.ID, err)
		return
	}

	// Log successful email sending for each user
	for _, user := range users {
		s.logEmailSuccess(account.ID, &user.ID, user.Email, fmt.Sprintf("Low Stock Alert - %s", account.Name), models.EmailTypeLowStockAlert)
	}

	log.Printf("Successfully sent low stock alert for account: %s (%d items)", account.Name, len(lowStockItems))
}

// sendWeeklySupplyChainReportForAccount sends a weekly supply chain report for a specific account
func (s *Scheduler) sendWeeklySupplyChainReportForAccount(account models.Account) {
	log.Printf("Sending weekly supply chain report for account: %s", account.Name)

	// Get all users in the account
	users, err := s.service.GetUsersByAccount(account.ID)
	if err != nil {
		log.Printf("Failed to get users for account %d: %v", account.ID, err)
		return
	}

	if len(users) == 0 {
		log.Printf("No users found for account %d", account.ID)
		return
	}

	// Generate supply chain report data
	supplyChainData, err := s.generateSupplyChainReportData(account.ID)
	if err != nil {
		log.Printf("Failed to generate supply chain report for account %d: %v", account.ID, err)
		return
	}

	// Send weekly supply chain report
	if err := s.emailService.SendWeeklySupplyChainReport(account, users, supplyChainData); err != nil {
		log.Printf("Failed to send weekly supply chain report for account %d: %v", account.ID, err)
		return
	}

	// Log successful email sending for each user
	for _, user := range users {
		s.logEmailSuccess(account.ID, &user.ID, user.Email, fmt.Sprintf("Weekly Supply Chain Report - %s", account.Name), models.EmailTypeWeeklySupplyChain)
	}

	log.Printf("Successfully sent weekly supply chain report for account: %s", account.Name)
}

// generateStockReportData generates stock report data for an account
func (s *Scheduler) generateStockReportData(accountID int) (*email.StockReportData, error) {
	// Get all inventory items for the account
	items, err := s.service.GetInventoryItemsByAccount(accountID)
	if err != nil {
		return nil, err
	}

	// Get latest inventory snapshot
	latestSnapshot, err := s.service.GetLatestInventorySnapshot(accountID)
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
			category, err := s.service.GetCategory(*item.CategoryID)
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
func (s *Scheduler) generateSupplyChainReportData(accountID int) (*email.SupplyChainData, error) {
	// Get all inventory items for the account
	items, err := s.service.GetInventoryItemsByAccount(accountID)
	if err != nil {
		return nil, err
	}

	// Get latest inventory snapshot
	latestSnapshot, err := s.service.GetLatestInventorySnapshot(accountID)
	if err != nil {
		return nil, err
	}

	// Get recent deliveries for vendor information
	recentDeliveries, err := s.service.GetDeliveriesByAccount(accountID)
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
			category, err := s.service.GetCategory(*item.CategoryID)
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
func (s *Scheduler) logEmailSuccess(accountID int, userID *int, toEmail, subject, emailType string) {
	emailLog := models.EmailLog{
		AccountID: accountID,
		UserID:    userID,
		ToEmail:   toEmail,
		Subject:   subject,
		EmailType: emailType,
		Status:    models.EmailStatusSent,
	}

	// TODO: Add CreateEmailLog method to service
	// For now, we'll just log to console
	log.Printf("Email sent successfully: %s to %s", emailType, toEmail)
	_ = emailLog
}
