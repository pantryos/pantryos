package handlers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

//func TestEmailScheduleHandlers(t *testing.T) {
//	// Setup test database
//	db, cleanup := database.SetupTestDBLegacy(t)
//	defer cleanup()
//
//	// Create test service
//	service := database.NewService(db)
//
//	// Create test organization and account
//	org := &models.Organization{Name: "Test Org"}
//	require.NoError(t, service.CreateOrganization(org))
//
//	account := &models.Account{
//		OrganizationID: &org.ID,
//		Name:           "Test Account",
//		Location:       "Test Location",
//		Phone:          "555-1234",
//		Email:          "test@example.com",
//		Status:         "active",
//	}
//	require.NoError(t, service.CreateAccount(account))
//
//	// Create test user
//	user := &models.User{
//		AccountID: account.ID,
//		Email:     "user@example.com",
//		Password:  "password123",
//		Role:      "admin",
//	}
//	require.NoError(t, service.CreateUser(user))
//
//	// Setup Gin router
//	gin.SetMode(gin.TestMode)
//	router := gin.New()
//
//	// Add middleware to set user and service in context
//	router.Use(func(c *gin.Context) {
//		c.Set("user", user)
//		c.Set("service", service)
//		c.Next()
//	})
//
//	// Add routes
//	router.GET("/accounts/:accountID/email-schedules", GetEmailSchedules)
//	router.POST("/accounts/:accountID/email-schedules", CreateEmailSchedule)
//	router.PATCH("/accounts/:accountID/email-schedules/:emailType/toggle", ToggleEmailSchedule)
//
//	t.Run("GetEmailSchedules", func(t *testing.T) {
//		// Test getting email schedules for account
//		w := httptest.NewRecorder()
//		req, _ := http.NewRequest("GET", "/accounts/1/email-schedules", nil)
//		router.ServeHTTP(w, req)
//
//		assert.Equal(t, http.StatusOK, w.Code)
//
//		var response map[string]interface{}
//		err := json.Unmarshal(w.Body.Bytes(), &response)
//		require.NoError(t, err)
//
//		assert.Equal(t, float64(1), response["account_id"])
//		assert.NotNil(t, response["schedules"])
//	})
//
//	t.Run("CreateEmailSchedule", func(t *testing.T) {
//		// Test creating a new email schedule
//		scheduleData := EmailScheduleRequest{
//			EmailType: "weekly_stock_report",
//			Frequency: "weekly",
//			DayOfWeek: &[]int{1}[0], // Monday
//			TimeOfDay: "09:00",
//			IsActive:  true,
//		}
//
//		jsonData, _ := json.Marshal(scheduleData)
//		w := httptest.NewRecorder()
//		req, _ := http.NewRequest("POST", "/accounts/1/email-schedules", bytes.NewBuffer(jsonData))
//		req.Header.Set("Content-Type", "application/json")
//		router.ServeHTTP(w, req)
//
//		assert.Equal(t, http.StatusCreated, w.Code)
//
//		var response models.EmailSchedule
//		err := json.Unmarshal(w.Body.Bytes(), &response)
//		require.NoError(t, err)
//
//		assert.Equal(t, 1, response.AccountID)
//		assert.Equal(t, "weekly_stock_report", response.EmailType)
//		assert.Equal(t, "weekly", response.Frequency)
//		assert.Equal(t, 1, *response.DayOfWeek)
//		assert.Equal(t, "09:00", response.TimeOfDay)
//		assert.True(t, response.IsActive)
//	})
//
//	t.Run("ToggleEmailSchedule", func(t *testing.T) {
//		// First create a schedule to toggle
//		scheduleData := EmailScheduleRequest{
//			EmailType: "weekly_supply_chain_report",
//			Frequency: "weekly",
//			DayOfWeek: &[]int{2}[0], // Tuesday
//			TimeOfDay: "10:00",
//			IsActive:  true,
//		}
//
//		jsonData, _ := json.Marshal(scheduleData)
//		w := httptest.NewRecorder()
//		req, _ := http.NewRequest("POST", "/accounts/1/email-schedules", bytes.NewBuffer(jsonData))
//		req.Header.Set("Content-Type", "application/json")
//		router.ServeHTTP(w, req)
//
//		// Now toggle the schedule
//		w = httptest.NewRecorder()
//		req, _ = http.NewRequest("PATCH", "/accounts/1/email-schedules/weekly_supply_chain_report/toggle", nil)
//		router.ServeHTTP(w, req)
//
//		assert.Equal(t, http.StatusOK, w.Code)
//
//		var response map[string]interface{}
//		err := json.Unmarshal(w.Body.Bytes(), &response)
//		require.NoError(t, err)
//
//		assert.Equal(t, "Email schedule disabled successfully", response["message"])
//		assert.False(t, response["is_active"].(bool))
//		assert.Equal(t, "weekly_supply_chain_report", response["email_type"])
//		assert.Equal(t, float64(1), response["account_id"])
//	})
//
//	t.Run("ToggleEmailScheduleBack", func(t *testing.T) {
//		// Toggle the schedule back to active
//		w := httptest.NewRecorder()
//		req, _ := http.NewRequest("PATCH", "/accounts/1/email-schedules/weekly_supply_chain_report/toggle", nil)
//		router.ServeHTTP(w, req)
//
//		assert.Equal(t, http.StatusOK, w.Code)
//
//		var response map[string]interface{}
//		err := json.Unmarshal(w.Body.Bytes(), &response)
//		require.NoError(t, err)
//
//		assert.Equal(t, "Email schedule enabled successfully", response["message"])
//		assert.True(t, response["is_active"].(bool))
//	})
//}

func TestIsValidEmailType(t *testing.T) {
	// Test valid email types
	assert.True(t, isValidEmailType("weekly_stock_report"))
	assert.True(t, isValidEmailType("weekly_supply_chain_report"))
	assert.True(t, isValidEmailType("low_stock_alert"))

	// Test invalid email types
	assert.False(t, isValidEmailType("invalid_type"))
	assert.False(t, isValidEmailType(""))
	assert.False(t, isValidEmailType("weekly_report"))
}
