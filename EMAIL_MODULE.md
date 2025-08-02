# Email Module Documentation

The email module provides comprehensive email functionality for the Stok inventory management system, including account verification, weekly stock reports, and low stock alerts.

## Features

### 1. Account Verification Emails
- Send verification emails to new users
- Secure token-based verification system
- 24-hour expiration for security
- Beautiful HTML email templates

### 2. Weekly Stock Reports
- Automated weekly stock reports sent to all users in an account
- Comprehensive inventory summary with current stock levels
- Visual status indicators (normal, low, out of stock)
- Total inventory value calculations

### 3. Low Stock Alerts
- Automated alerts when items fall below minimum stock levels
- Sent to all users in the affected account
- Immediate notification system for inventory management

### 4. Email Scheduling
- Automated scheduler for weekly reports (Mondays at 9 AM)
- Automated low stock alerts (every 12 hours)
- Configurable scheduling system

## Configuration

### Environment Variables

Set the following environment variables to configure email functionality:

```bash
# SMTP Configuration
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your-email@gmail.com
SMTP_PASSWORD=your-app-password
SMTP_USE_TLS=true

# Email Sender Configuration
FROM_EMAIL=noreply@stok.com
FROM_NAME=Stok Inventory System
```

### Gmail Setup

To use Gmail for sending emails:

1. Enable 2-factor authentication on your Gmail account
2. Generate an App Password:
   - Go to Google Account settings
   - Security → 2-Step Verification → App passwords
   - Generate a password for "Mail"
3. Use the generated password as `SMTP_PASSWORD`

## API Endpoints

### Account Verification

#### Send Verification Email
```http
POST /api/v1/email/verification/{user_id}
```

Sends a verification email to the specified user.

**Response:**
```json
{
  "message": "Verification email sent successfully",
  "user_id": 123
}
```

#### Verify Email
```http
GET /api/v1/email/verify?token={verification_token}
```

Verifies a user's email using the provided token.

**Response:**
```json
{
  "message": "Email verified successfully",
  "user_id": 123
}
```

### Stock Reports

#### Send Weekly Stock Report
```http
POST /api/v1/email/weekly-report/{account_id}
```

Manually triggers a weekly stock report for the specified account.

**Response:**
```json
{
  "message": "Weekly stock report sent successfully",
  "account_id": 1,
  "users_count": 3
}
```

#### Send Low Stock Alert
```http
POST /api/v1/email/low-stock-alert/{account_id}
```

Manually triggers a low stock alert for the specified account.

**Response:**
```json
{
  "message": "Low stock alert sent successfully",
  "account_id": 1,
  "users_count": 3,
  "low_stock_items_count": 2
}
```

## Email Templates

### Account Verification Template
- Professional welcome message
- Clear verification button
- Fallback verification link
- 24-hour expiration notice

### Weekly Stock Report Template
- Summary dashboard with key metrics
- Detailed inventory table
- Color-coded status indicators
- Total inventory value

### Low Stock Alert Template
- Urgent action required message
- List of low stock items
- Current stock levels
- Minimum stock requirements

## Automated Scheduling

The email scheduler runs automatically and handles:

### Weekly Stock Reports
- **Schedule:** Every Monday at 9:00 AM
- **Recipients:** All users in each account
- **Content:** Complete inventory status report

### Low Stock Alerts
- **Schedule:** Every 12 hours
- **Recipients:** All users in accounts with low stock items
- **Content:** List of items requiring reorder

## Database Models

### EmailVerificationToken
```go
type EmailVerificationToken struct {
    ID        int       `json:"id"`
    UserID    int       `json:"user_id"`
    Token     string    `json:"token"`
    Type      string    `json:"type"`
    ExpiresAt time.Time `json:"expires_at"`
    UsedAt    *time.Time `json:"used_at"`
    CreatedAt time.Time `json:"created_at"`
}
```

### EmailLog
```go
type EmailLog struct {
    ID          int       `json:"id"`
    AccountID   int       `json:"account_id"`
    UserID      *int      `json:"user_id"`
    ToEmail     string    `json:"to_email"`
    Subject     string    `json:"subject"`
    EmailType   string    `json:"email_type"`
    Status      string    `json:"status"`
    ErrorMsg    string    `json:"error_msg"`
    SentAt      time.Time `json:"sent_at"`
}
```

### EmailSchedule
```go
type EmailSchedule struct {
    ID          int       `json:"id"`
    AccountID   int       `json:"account_id"`
    EmailType   string    `json:"email_type"`
    Frequency   string    `json:"frequency"`
    DayOfWeek   *int      `json:"day_of_week"`
    DayOfMonth  *int      `json:"day_of_month"`
    TimeOfDay   string    `json:"time_of_day"`
    IsActive    bool      `json:"is_active"`
    LastSentAt  *time.Time `json:"last_sent_at"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}
```

## Usage Examples

### Sending a Verification Email
```go
emailService := email.NewEmailService()
user := models.User{Email: "user@example.com"}
account := models.Account{Name: "Coffee Shop"}
verificationURL := "https://app.example.com/verify?token=abc123"

err := emailService.SendVerificationEmail(user, account, verificationURL)
if err != nil {
    log.Printf("Failed to send verification email: %v", err)
}
```

### Sending a Weekly Stock Report
```go
emailService := email.NewEmailService()
account := models.Account{Name: "Coffee Shop"}
users := []models.User{...}
stockData := &email.StockReportData{...}

err := emailService.SendWeeklyStockReport(account, users, stockData)
if err != nil {
    log.Printf("Failed to send weekly stock report: %v", err)
}
```

### Sending a Low Stock Alert
```go
emailService := email.NewEmailService()
account := models.Account{Name: "Coffee Shop"}
users := []models.User{...}
lowStockItems := []models.InventoryItem{...}

err := emailService.SendLowStockAlert(account, users, lowStockItems)
if err != nil {
    log.Printf("Failed to send low stock alert: %v", err)
}
```

## Testing

Run the email module tests:

```bash
go test ./internal/email
```

## Security Considerations

1. **Token Security:** Verification tokens are cryptographically secure and expire after 24 hours
2. **SMTP Security:** TLS encryption is used for all email communications
3. **Access Control:** Email endpoints require authentication
4. **Rate Limiting:** Consider implementing rate limiting for email endpoints
5. **Logging:** All email activities are logged for audit purposes

## Troubleshooting

### Common Issues

1. **SMTP Authentication Failed**
   - Verify SMTP credentials
   - Check if 2FA is enabled for Gmail
   - Use App Password instead of regular password

2. **Emails Not Sending**
   - Check SMTP configuration
   - Verify network connectivity
   - Check email logs for error messages

3. **Scheduler Not Running**
   - Verify the scheduler is started in main.go
   - Check application logs for scheduler messages

### Debug Mode

Enable debug logging by setting the log level:

```bash
export GIN_MODE=debug
```

## Future Enhancements

1. **Email Templates:** Add more customizable email templates
2. **Scheduling:** Allow per-account email scheduling configuration
3. **Notifications:** Add SMS notifications for critical alerts
4. **Analytics:** Email open/click tracking
5. **Bulk Operations:** Support for bulk email operations
6. **Template Editor:** Web-based email template editor 