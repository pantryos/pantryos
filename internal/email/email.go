package email

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"
	"net/smtp"
	"os"
	"time"

	"github.com/mnadev/stok/internal/models"
)

// EmailService handles all email operations
type EmailService struct {
	config *EmailConfig
}

// EmailConfig holds email server configuration
type EmailConfig struct {
	SMTPHost     string
	SMTPPort     string
	SMTPUsername string
	SMTPPassword string
	FromEmail    string
	FromName     string
	UseTLS       bool
}

// EmailData holds data for email templates
type EmailData struct {
	AccountName       string
	UserName          string
	UserEmail         string
	VerificationURL   string
	StockReport       *StockReportData
	SupplyChainReport *SupplyChainData
	LowStockItems     []models.InventoryItem
	ExpiringItems     []models.InventoryItem
}

// StockReportData holds data for weekly stock reports
type StockReportData struct {
	ReportDate      time.Time
	TotalItems      int
	LowStockItems   int
	OutOfStockItems int
	TotalValue      float64
	Items           []StockItemData
}

// StockItemData holds individual item data for stock reports
type StockItemData struct {
	ID           int
	Name         string
	Category     string
	CurrentStock float64
	MinStock     float64
	MaxStock     float64
	Unit         string
	Status       string // "normal", "low", "out"
}

// SupplyChainData holds data for weekly supply chain reports
type SupplyChainData struct {
	ReportDate        time.Time
	TotalItems        int
	LowStockItems     int
	OutOfStockItems   int
	CriticalItems     int
	TotalValue        float64
	EstimatedReorders float64
	Items             []SupplyChainItemData
}

// SupplyChainItemData holds individual item data for supply chain reports
type SupplyChainItemData struct {
	ID                int
	Name              string
	Category          string
	CurrentStock      float64
	MinStock          float64
	MaxStock          float64
	Unit              string
	Status            string // "normal", "low", "critical", "out"
	PreferredVendor   string
	CostPerUnit       float64
	ReorderQuantity   float64
	DaysUntilStockout int
	LastDeliveryDate  *time.Time
}

// NewEmailService creates a new email service with configuration
func NewEmailService() *EmailService {
	config := &EmailConfig{
		SMTPHost:     getEnvOrDefault("SMTP_HOST", "smtp.gmail.com"),
		SMTPPort:     getEnvOrDefault("SMTP_PORT", "587"),
		SMTPUsername: getEnvOrDefault("SMTP_USERNAME", ""),
		SMTPPassword: getEnvOrDefault("SMTP_PASSWORD", ""),
		FromEmail:    getEnvOrDefault("FROM_EMAIL", "noreply@stok.com"),
		FromName:     getEnvOrDefault("FROM_NAME", "Stok Inventory System"),
		UseTLS:       getEnvOrDefault("SMTP_USE_TLS", "true") == "true",
	}

	return &EmailService{
		config: config,
	}
}

// SendVerificationEmail sends account verification email
func (es *EmailService) SendVerificationEmail(user models.User, account models.Account, verificationURL string) error {
	data := EmailData{
		AccountName:     account.Name,
		UserName:        user.Email, // Using email as username for now
		UserEmail:       user.Email,
		VerificationURL: verificationURL,
	}

	subject := "Verify Your Stok Account"
	body, err := es.renderTemplate("verification", data)
	if err != nil {
		return fmt.Errorf("failed to render verification template: %w", err)
	}

	return es.sendEmail(user.Email, subject, body)
}

// SendWeeklyStockReport sends weekly stock report email
func (es *EmailService) SendWeeklyStockReport(account models.Account, users []models.User, stockData *StockReportData) error {
	data := EmailData{
		AccountName: account.Name,
		StockReport: stockData,
	}

	subject := fmt.Sprintf("Weekly Stock Report - %s", account.Name)
	body, err := es.renderTemplate("weekly_stock_report", data)
	if err != nil {
		return fmt.Errorf("failed to render weekly stock report template: %w", err)
	}

	// Send to all users in the account
	for _, user := range users {
		if err := es.sendEmail(user.Email, subject, body); err != nil {
			// Log error but continue with other users
			fmt.Printf("Failed to send weekly stock report to %s: %v\n", user.Email, err)
		}
	}

	return nil
}

// SendLowStockAlert sends low stock alert email
func (es *EmailService) SendLowStockAlert(account models.Account, users []models.User, lowStockItems []models.InventoryItem) error {
	data := EmailData{
		AccountName:   account.Name,
		LowStockItems: lowStockItems,
	}

	subject := fmt.Sprintf("Low Stock Alert - %s", account.Name)
	body, err := es.renderTemplate("low_stock_alert", data)
	if err != nil {
		return fmt.Errorf("failed to render low stock alert template: %w", err)
	}

	// Send to all users in the account
	for _, user := range users {
		if err := es.sendEmail(user.Email, subject, body); err != nil {
			// Log error but continue with other users
			fmt.Printf("Failed to send low stock alert to %s: %v\n", user.Email, err)
		}
	}

	return nil
}

// SendWeeklySupplyChainReport sends weekly supply chain report email
func (es *EmailService) SendWeeklySupplyChainReport(account models.Account, users []models.User, supplyChainData *SupplyChainData) error {
	data := EmailData{
		AccountName:       account.Name,
		SupplyChainReport: supplyChainData,
	}

	subject := fmt.Sprintf("Weekly Supply Chain Report - %s", account.Name)
	body, err := es.renderTemplate("weekly_supply_chain_report", data)
	if err != nil {
		return fmt.Errorf("failed to render weekly supply chain report template: %w", err)
	}

	// Send to all users in the account
	for _, user := range users {
		if err := es.sendEmail(user.Email, subject, body); err != nil {
			// Log error but continue with other users
			fmt.Printf("Failed to send weekly supply chain report to %s: %v\n", user.Email, err)
		}
	}

	return nil
}

// sendEmail sends an email using SMTP
func (es *EmailService) sendEmail(to, subject, body string) error {
	if es.config.SMTPUsername == "" || es.config.SMTPPassword == "" {
		return fmt.Errorf("SMTP credentials not configured")
	}

	// Create email message
	message := fmt.Sprintf("From: %s <%s>\r\n", es.config.FromName, es.config.FromEmail)
	message += fmt.Sprintf("To: %s\r\n", to)
	message += fmt.Sprintf("Subject: %s\r\n", subject)
	message += "MIME-Version: 1.0\r\n"
	message += "Content-Type: text/html; charset=UTF-8\r\n"
	message += "\r\n"
	message += body

	// Connect to SMTP server
	auth := smtp.PlainAuth("", es.config.SMTPUsername, es.config.SMTPPassword, es.config.SMTPHost)

	addr := fmt.Sprintf("%s:%s", es.config.SMTPHost, es.config.SMTPPort)

	var err error
	if es.config.UseTLS {
		err = es.sendEmailWithTLS(to, message, addr, auth)
	} else {
		err = smtp.SendMail(addr, auth, es.config.FromEmail, []string{to}, []byte(message))
	}

	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

// sendEmailWithTLS sends email with TLS encryption
func (es *EmailService) sendEmailWithTLS(to, message, addr string, auth smtp.Auth) error {
	// Connect to SMTP server
	conn, err := smtp.Dial(addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	// Start TLS
	if err = conn.StartTLS(&tls.Config{ServerName: es.config.SMTPHost}); err != nil {
		return err
	}

	// Authenticate
	if err = conn.Auth(auth); err != nil {
		return err
	}

	// Send email
	if err = conn.Mail(es.config.FromEmail); err != nil {
		return err
	}

	if err = conn.Rcpt(to); err != nil {
		return err
	}

	writer, err := conn.Data()
	if err != nil {
		return err
	}

	_, err = writer.Write([]byte(message))
	if err != nil {
		return err
	}

	err = writer.Close()
	if err != nil {
		return err
	}

	return conn.Quit()
}

// renderTemplate renders an email template with the given data
func (es *EmailService) renderTemplate(templateName string, data EmailData) (string, error) {
	tmpl, err := template.New(templateName).Parse(getEmailTemplate(templateName))
	if err != nil {
		return "", fmt.Errorf("failed to parse template %s: %w", templateName, err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template %s: %w", templateName, err)
	}

	return buf.String(), nil
}

// getEmailTemplate returns the HTML template for the given template name
func getEmailTemplate(templateName string) string {
	templates := map[string]string{
		"verification": `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Verify Your Account</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background-color: #4F46E5; color: white; padding: 20px; text-align: center; }
        .content { padding: 20px; background-color: #f9f9f9; }
        .button { display: inline-block; padding: 12px 24px; background-color: #4F46E5; color: white; text-decoration: none; border-radius: 5px; margin: 20px 0; }
        .footer { text-align: center; padding: 20px; color: #666; font-size: 12px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Welcome to Stok!</h1>
        </div>
        <div class="content">
            <h2>Verify Your Account</h2>
            <p>Hello,</p>
            <p>Thank you for joining <strong>{{.AccountName}}</strong> on Stok! To complete your account setup, please verify your email address.</p>
            <p>Click the button below to verify your account:</p>
            <div style="text-align: center;">
                <a href="{{.VerificationURL}}" class="button">Verify Account</a>
            </div>
            <p>If the button doesn't work, you can copy and paste this link into your browser:</p>
            <p style="word-break: break-all; color: #666;">{{.VerificationURL}}</p>
            <p>This link will expire in 24 hours for security reasons.</p>
            <p>If you didn't create an account with Stok, please ignore this email.</p>
        </div>
        <div class="footer">
            <p>© 2024 Stok Inventory System. All rights reserved.</p>
        </div>
    </div>
</body>
</html>`,

		"weekly_stock_report": `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Weekly Stock Report</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 800px; margin: 0 auto; padding: 20px; }
        .header { background-color: #4F46E5; color: white; padding: 20px; text-align: center; }
        .content { padding: 20px; background-color: #f9f9f9; }
        .summary { background-color: white; padding: 20px; margin: 20px 0; border-radius: 5px; }
        .summary-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(200px, 1fr)); gap: 20px; margin: 20px 0; }
        .summary-item { text-align: center; padding: 15px; background-color: #f8f9fa; border-radius: 5px; }
        .summary-number { font-size: 24px; font-weight: bold; color: #4F46E5; }
        .table { width: 100%; border-collapse: collapse; margin: 20px 0; }
        .table th, .table td { padding: 12px; text-align: left; border-bottom: 1px solid #ddd; }
        .table th { background-color: #f8f9fa; font-weight: bold; }
        .status-normal { color: #28a745; }
        .status-low { color: #ffc107; }
        .status-out { color: #dc3545; }
        .footer { text-align: center; padding: 20px; color: #666; font-size: 12px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Weekly Stock Report</h1>
            <p>{{.AccountName}}</p>
        </div>
        <div class="content">
            <h2>Stock Report for {{.StockReport.ReportDate.Format "January 2, 2006"}}</h2>
            
            <div class="summary">
                <h3>Summary</h3>
                <div class="summary-grid">
                    <div class="summary-item">
                        <div class="summary-number">{{.StockReport.TotalItems}}</div>
                        <div>Total Items</div>
                    </div>
                    <div class="summary-item">
                        <div class="summary-number">{{.StockReport.LowStockItems}}</div>
                        <div>Low Stock Items</div>
                    </div>
                    <div class="summary-item">
                        <div class="summary-number">{{.StockReport.OutOfStockItems}}</div>
                        <div>Out of Stock</div>
                    </div>
                    <div class="summary-item">
                        <div class="summary-number">${{printf "%.2f" .StockReport.TotalValue}}</div>
                        <div>Total Value</div>
                    </div>
                </div>
            </div>

            <h3>Item Details</h3>
            <table class="table">
                <thead>
                    <tr>
                        <th>Item</th>
                        <th>Category</th>
                        <th>Current Stock</th>
                        <th>Min Stock</th>
                        <th>Max Stock</th>
                        <th>Status</th>
                    </tr>
                </thead>
                <tbody>
                    {{range .StockReport.Items}}
                    <tr>
                        <td>{{.Name}}</td>
                        <td>{{.Category}}</td>
                        <td>{{.CurrentStock}} {{.Unit}}</td>
                        <td>{{.MinStock}} {{.Unit}}</td>
                        <td>{{.MaxStock}} {{.Unit}}</td>
                        <td class="status-{{.Status}}">{{.Status}}</td>
                    </tr>
                    {{end}}
                </tbody>
            </table>
        </div>
        <div class="footer">
            <p>© 2024 Stok Inventory System. All rights reserved.</p>
        </div>
    </div>
</body>
</html>`,

		"low_stock_alert": `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Low Stock Alert</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background-color: #dc3545; color: white; padding: 20px; text-align: center; }
        .content { padding: 20px; background-color: #f9f9f9; }
        .alert { background-color: #fff3cd; border: 1px solid #ffeaa7; padding: 15px; margin: 20px 0; border-radius: 5px; }
        .table { width: 100%; border-collapse: collapse; margin: 20px 0; }
        .table th, .table td { padding: 12px; text-align: left; border-bottom: 1px solid #ddd; }
        .table th { background-color: #f8f9fa; font-weight: bold; }
        .footer { text-align: center; padding: 20px; color: #666; font-size: 12px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>⚠️ Low Stock Alert</h1>
            <p>{{.AccountName}}</p>
        </div>
        <div class="content">
            <div class="alert">
                <h3>Action Required</h3>
                <p>The following items are running low on stock and may need to be reordered soon:</p>
            </div>

            <table class="table">
                <thead>
                    <tr>
                        <th>Item</th>
                        <th>Current Stock</th>
                        <th>Min Stock Level</th>
                        <th>Unit</th>
                    </tr>
                </thead>
                <tbody>
                    {{range .LowStockItems}}
                    <tr>
                        <td>{{.Name}}</td>
                        <td>{{.MinStockLevel}}</td>
                        <td>{{.MinStockLevel}}</td>
                        <td>{{.Unit}}</td>
                    </tr>
                    {{end}}
                </tbody>
            </table>

            <p><strong>Please review these items and place orders as needed to maintain adequate stock levels.</strong></p>
        </div>
        <div class="footer">
            <p>© 2024 Stok Inventory System. All rights reserved.</p>
        </div>
    </div>
</body>
</html>`,

		"weekly_supply_chain_report": `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Weekly Supply Chain Report</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 900px; margin: 0 auto; padding: 20px; }
        .header { background-color: #17a2b8; color: white; padding: 20px; text-align: center; }
        .content { padding: 20px; background-color: #f9f9f9; }
        .summary { background-color: white; padding: 20px; margin: 20px 0; border-radius: 5px; }
        .summary-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(180px, 1fr)); gap: 20px; margin: 20px 0; }
        .summary-item { text-align: center; padding: 15px; background-color: #f8f9fa; border-radius: 5px; }
        .summary-number { font-size: 24px; font-weight: bold; color: #17a2b8; }
        .table { width: 100%; border-collapse: collapse; margin: 20px 0; }
        .table th, .table td { padding: 12px; text-align: left; border-bottom: 1px solid #ddd; }
        .table th { background-color: #f8f9fa; font-weight: bold; }
        .status-normal { color: #28a745; }
        .status-low { color: #ffc107; }
        .status-critical { color: #fd7e14; }
        .status-out { color: #dc3545; }
        .vendor-info { background-color: #e3f2fd; padding: 10px; border-radius: 5px; margin: 10px 0; }
        .footer { text-align: center; padding: 20px; color: #666; font-size: 12px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>📊 Weekly Supply Chain Report</h1>
            <p>{{.AccountName}}</p>
        </div>
        <div class="content">
            <h2>Supply Chain Report for {{.SupplyChainReport.ReportDate.Format "January 2, 2006"}}</h2>
            
            <div class="summary">
                <h3>Supply Chain Summary</h3>
                <div class="summary-grid">
                    <div class="summary-item">
                        <div class="summary-number">{{.SupplyChainReport.TotalItems}}</div>
                        <div>Total Items</div>
                    </div>
                    <div class="summary-item">
                        <div class="summary-number">{{.SupplyChainReport.LowStockItems}}</div>
                        <div>Low Stock Items</div>
                    </div>
                    <div class="summary-item">
                        <div class="summary-number">{{.SupplyChainReport.CriticalItems}}</div>
                        <div>Critical Items</div>
                    </div>
                    <div class="summary-item">
                        <div class="summary-number">{{.SupplyChainReport.OutOfStockItems}}</div>
                        <div>Out of Stock</div>
                    </div>
                    <div class="summary-item">
                        <div class="summary-number">${{printf "%.2f" .SupplyChainReport.TotalValue}}</div>
                        <div>Total Value</div>
                    </div>
                    <div class="summary-item">
                        <div class="summary-number">${{printf "%.2f" .SupplyChainReport.EstimatedReorders}}</div>
                        <div>Est. Reorder Cost</div>
                    </div>
                </div>
            </div>

            <h3>Supply Chain Details</h3>
            <table class="table">
                <thead>
                    <tr>
                        <th>Item</th>
                        <th>Category</th>
                        <th>Current Stock</th>
                        <th>Min Stock</th>
                        <th>Status</th>
                        <th>Vendor</th>
                        <th>Cost/Unit</th>
                        <th>Reorder Qty</th>
                        <th>Days Until Stockout</th>
                    </tr>
                </thead>
                <tbody>
                    {{range .SupplyChainReport.Items}}
                    <tr>
                        <td>{{.Name}}</td>
                        <td>{{.Category}}</td>
                        <td>{{.CurrentStock}} {{.Unit}}</td>
                        <td>{{.MinStock}} {{.Unit}}</td>
                        <td class="status-{{.Status}}">{{.Status}}</td>
                        <td>{{.PreferredVendor}}</td>
                        <td>${{printf "%.2f" .CostPerUnit}}</td>
                        <td>{{.ReorderQuantity}} {{.Unit}}</td>
                        <td>{{.DaysUntilStockout}} days</td>
                    </tr>
                    {{end}}
                </tbody>
            </table>

            <div class="vendor-info">
                <h4>📋 Supply Chain Recommendations</h4>
                <ul>
                    <li><strong>Critical Items:</strong> Items with {{.SupplyChainReport.CriticalItems}} critical status need immediate attention</li>
                    <li><strong>Low Stock Items:</strong> {{.SupplyChainReport.LowStockItems}} items are below minimum stock levels</li>
                    <li><strong>Estimated Reorder Cost:</strong> ${{printf "%.2f" .SupplyChainReport.EstimatedReorders}} for recommended reorders</li>
                    <li><strong>Vendor Management:</strong> Review preferred vendors for optimal pricing and delivery</li>
                </ul>
            </div>
        </div>
        <div class="footer">
            <p>© 2024 Stok Inventory System. All rights reserved.</p>
        </div>
    </div>
</body>
</html>`,
	}

	return templates[templateName]
}

// getEnvOrDefault returns environment variable value or default
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
