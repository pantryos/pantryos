package email

import (
	"testing"
	"time"

	"github.com/mnadev/stok/internal/models"
)

func TestNewEmailService(t *testing.T) {
	service := NewEmailService()
	if service == nil {
		t.Fatal("Expected email service to be created")
	}

	if service.config == nil {
		t.Fatal("Expected email config to be initialized")
	}
}

func TestEmailTemplates(t *testing.T) {
	service := NewEmailService()

	// Test verification template
	data := EmailData{
		AccountName:     "Test Coffee Shop",
		UserName:        "test@example.com",
		UserEmail:       "test@example.com",
		VerificationURL: "https://example.com/verify?token=test123",
	}

	body, err := service.renderTemplate("verification", data)
	if err != nil {
		t.Fatalf("Failed to render verification template: %v", err)
	}

	if body == "" {
		t.Fatal("Expected non-empty email body")
	}

	// Test weekly stock report template
	stockData := &StockReportData{
		ReportDate:      time.Now(),
		TotalItems:      5,
		LowStockItems:   2,
		OutOfStockItems: 1,
		TotalValue:      150.50,
		Items: []StockItemData{
			{
				ID:           1,
				Name:         "Coffee Beans",
				Category:     "Coffee",
				CurrentStock: 10.5,
				MinStock:     5.0,
				MaxStock:     20.0,
				Unit:         "kg",
				Status:       "normal",
			},
		},
	}

	data.StockReport = stockData
	body, err = service.renderTemplate("weekly_stock_report", data)
	if err != nil {
		t.Fatalf("Failed to render weekly stock report template: %v", err)
	}

	if body == "" {
		t.Fatal("Expected non-empty email body")
	}

	// Test low stock alert template
	lowStockItems := []models.InventoryItem{
		{
			ID:            1,
			AccountID:     1,
			Name:          "Coffee Beans",
			Unit:          "kg",
			CostPerUnit:   15.0,
			MinStockLevel: 5.0,
		},
	}

	data.LowStockItems = lowStockItems
	body, err = service.renderTemplate("low_stock_alert", data)
	if err != nil {
		t.Fatalf("Failed to render low stock alert template: %v", err)
	}

	if body == "" {
		t.Fatal("Expected non-empty email body")
	}
}

func TestGetEmailTemplate(t *testing.T) {
	// Test verification template
	template := getEmailTemplate("verification")
	if template == "" {
		t.Fatal("Expected non-empty verification template")
	}

	// Test weekly stock report template
	template = getEmailTemplate("weekly_stock_report")
	if template == "" {
		t.Fatal("Expected non-empty weekly stock report template")
	}

	// Test low stock alert template
	template = getEmailTemplate("low_stock_alert")
	if template == "" {
		t.Fatal("Expected non-empty low stock alert template")
	}

	// Test non-existent template
	template = getEmailTemplate("non_existent")
	if template != "" {
		t.Fatal("Expected empty template for non-existent template name")
	}
}

func TestGetEnvOrDefault(t *testing.T) {
	// Test with default value
	result := getEnvOrDefault("NON_EXISTENT_ENV", "default_value")
	if result != "default_value" {
		t.Fatalf("Expected 'default_value', got '%s'", result)
	}
}
