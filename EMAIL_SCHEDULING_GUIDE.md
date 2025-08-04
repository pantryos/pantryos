# Email Scheduling Guide

## Overview

PantryOS now supports configurable email scheduling that allows users to control when and how often they receive automated emails. Users can turn emails on/off, customize schedules, and manage different types of email notifications.

## Email Types Available

1. **Weekly Stock Report** - Inventory status and stock levels
2. **Weekly Supply Chain Report** - Supply chain analytics and reorder recommendations  
3. **Low Stock Alert** - Immediate notifications for items running low

## How to Manage Email Schedules

### Via API Endpoints

#### Get All Email Schedules
```bash
GET /api/v1/accounts/{accountID}/email-schedules
```

#### Get Specific Email Schedule
```bash
GET /api/v1/accounts/{accountID}/email-schedules/{emailType}
```

#### Create New Email Schedule
```bash
POST /api/v1/accounts/{accountID}/email-schedules
Content-Type: application/json

{
  "email_type": "weekly_stock_report",
  "frequency": "weekly",
  "day_of_week": 1,
  "time_of_day": "09:00",
  "is_active": true
}
```

#### Update Email Schedule
```bash
PUT /api/v1/accounts/{accountID}/email-schedules/{emailType}
Content-Type: application/json

{
  "email_type": "weekly_stock_report",
  "frequency": "weekly",
  "day_of_week": 2,
  "time_of_day": "14:00",
  "is_active": false
}
```

#### Toggle Email Schedule (Turn On/Off)
```bash
PATCH /api/v1/accounts/{accountID}/email-schedules/{emailType}/toggle
```

#### Delete Email Schedule
```bash
DELETE /api/v1/accounts/{accountID}/email-schedules/{emailType}
```

### Via Frontend Component

The `EmailSettings` component provides a user-friendly interface for managing email schedules:

```tsx
import EmailSettings from './components/EmailSettings';

// In your app
<EmailSettings />
```

## How to Turn Off Emails

### Method 1: Toggle Active Status
Set `is_active: false` in the email schedule. This keeps the schedule but disables it:

```json
{
  "email_type": "weekly_stock_report",
  "frequency": "weekly",
  "day_of_week": 1,
  "time_of_day": "09:00",
  "is_active": false
}
```

### Method 2: Delete the Schedule
Completely remove the email schedule using the DELETE endpoint.

### Method 3: Use the Toggle Endpoint
Quickly turn on/off using the toggle endpoint:

```bash
# Turn off
PATCH /api/v1/accounts/1/email-schedules/weekly_stock_report/toggle

# Turn back on
PATCH /api/v1/accounts/1/email-schedules/weekly_stock_report/toggle
```

## Schedule Configuration Options

### Frequency
- `weekly` - Send once per week
- `daily` - Send every day
- `monthly` - Send once per month

### Day of Week (for weekly frequency)
- `0` - Sunday
- `1` - Monday
- `2` - Tuesday
- `3` - Wednesday
- `4` - Thursday
- `5` - Friday
- `6` - Saturday

### Time of Day
Format: `HH:MM` (24-hour format)
Examples:
- `09:00` - 9:00 AM
- `14:30` - 2:30 PM
- `18:00` - 6:00 PM

## Default Behavior

If no email schedule is configured for an account:
- **Weekly Stock Reports**: Sent every Monday at 9:00 AM
- **Weekly Supply Chain Reports**: Sent every Tuesday at 9:00 AM
- **Low Stock Alerts**: Sent every 12 hours when items are low

## Examples

### Turn Off Weekly Stock Reports
```bash
curl -X PATCH \
  http://localhost:8080/api/v1/accounts/1/email-schedules/weekly_stock_report/toggle \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### Change Weekly Report to Wednesday at 2 PM
```bash
curl -X PUT \
  http://localhost:8080/api/v1/accounts/1/email-schedules/weekly_stock_report \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "email_type": "weekly_stock_report",
    "frequency": "weekly",
    "day_of_week": 3,
    "time_of_day": "14:00",
    "is_active": true
  }'
```

### Create Custom Low Stock Alert Schedule
```bash
curl -X POST \
  http://localhost:8080/api/v1/accounts/1/email-schedules \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "email_type": "low_stock_alert",
    "frequency": "daily",
    "time_of_day": "08:00",
    "is_active": true
  }'
```

## Security

- All endpoints require authentication
- Users can only manage email schedules for their own account
- Account access is validated for all operations

## Error Handling

The API returns appropriate HTTP status codes:
- `200 OK` - Success
- `400 Bad Request` - Invalid request data
- `401 Unauthorized` - Missing or invalid authentication
- `403 Forbidden` - Access denied
- `404 Not Found` - Email schedule not found
- `409 Conflict` - Schedule already exists
- `500 Internal Server Error` - Server error

## Monitoring

Email schedules include:
- `last_sent_at` - Timestamp of last email sent
- `created_at` - When the schedule was created
- `updated_at` - When the schedule was last modified

This helps track email delivery and schedule changes. 