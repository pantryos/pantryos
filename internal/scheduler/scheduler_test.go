package scheduler

import (
	"testing"
	"time"

	"github.com/mnadev/pantryos/internal/database"
	"github.com/mnadev/pantryos/internal/models"
)

func TestShouldSendWeeklyReport(t *testing.T) {
	// Initialize test database using SQLite
	db, cleanup := database.SetupTestDBLegacy(t)
	defer cleanup()

	// Create scheduler
	scheduler := NewScheduler(db)

	// Test case 1: No schedule exists (should use default behavior)
	// This should return true only on Monday at 9 AM
	now := time.Now()
	shouldSend := scheduler.shouldSendWeeklyReport(1) // account ID 1

	if now.Weekday() == time.Monday && now.Hour() == 9 {
		if !shouldSend {
			t.Error("Expected shouldSend to be true on Monday at 9 AM with no schedule")
		}
	} else {
		if shouldSend {
			t.Error("Expected shouldSend to be false when not Monday at 9 AM with no schedule")
		}
	}

	// Test case 2: Create a custom schedule
	schedule := &models.EmailSchedule{
		AccountID: 1,
		EmailType: models.EmailTypeWeeklyReport,
		Frequency: "weekly",
		DayOfWeek: &[]int{3}[0], // Wednesday
		TimeOfDay: "14:00",      // 2 PM
		IsActive:  true,
	}

	if err := scheduler.service.CreateEmailSchedule(schedule); err != nil {
		t.Fatalf("Failed to create test schedule: %v", err)
	}

	// Test with custom schedule
	// Should return true only on Wednesday at 2 PM
	shouldSend = scheduler.shouldSendWeeklyReport(1)

	if now.Weekday() == time.Wednesday && now.Hour() == 14 {
		if !shouldSend {
			t.Error("Expected shouldSend to be true on Wednesday at 2 PM with custom schedule")
		}
	} else {
		if shouldSend {
			t.Error("Expected shouldSend to be false when not Wednesday at 2 PM with custom schedule")
		}
	}

	// Test case 3: Inactive schedule
	schedule.IsActive = false
	if err := scheduler.service.UpdateEmailSchedule(schedule); err != nil {
		t.Fatalf("Failed to update test schedule: %v", err)
	}

	shouldSend = scheduler.shouldSendWeeklyReport(1)
	if shouldSend {
		t.Error("Expected shouldSend to be false with inactive schedule")
	}
}

func TestShouldSendWeeklySupplyChainReport(t *testing.T) {
	// Initialize test database using SQLite
	db, cleanup := database.SetupTestDBLegacy(t)
	defer cleanup()

	// Create scheduler
	scheduler := NewScheduler(db)

	// Test case 1: No schedule exists (should use default behavior)
	// This should return true only on Tuesday at 9 AM
	now := time.Now()
	shouldSend := scheduler.shouldSendWeeklySupplyChainReport(1) // account ID 1

	if now.Weekday() == time.Tuesday && now.Hour() == 9 {
		if !shouldSend {
			t.Error("Expected shouldSend to be true on Tuesday at 9 AM with no schedule")
		}
	} else {
		if shouldSend {
			t.Error("Expected shouldSend to be false when not Tuesday at 9 AM with no schedule")
		}
	}

	// Test case 2: Create a custom schedule
	schedule := &models.EmailSchedule{
		AccountID: 1,
		EmailType: models.EmailTypeWeeklySupplyChain,
		Frequency: "weekly",
		DayOfWeek: &[]int{4}[0], // Thursday
		TimeOfDay: "16:00",      // 4 PM
		IsActive:  true,
	}

	if err := scheduler.service.CreateEmailSchedule(schedule); err != nil {
		t.Fatalf("Failed to create test schedule: %v", err)
	}

	// Test with custom schedule
	// Should return true only on Thursday at 4 PM
	shouldSend = scheduler.shouldSendWeeklySupplyChainReport(1)

	if now.Weekday() == time.Thursday && now.Hour() == 16 {
		if !shouldSend {
			t.Error("Expected shouldSend to be true on Thursday at 4 PM with custom schedule")
		}
	} else {
		if shouldSend {
			t.Error("Expected shouldSend to be false when not Thursday at 4 PM with custom schedule")
		}
	}
}

func TestCreateDefaultEmailSchedules(t *testing.T) {
	// Initialize test database using SQLite
	db, cleanup := database.SetupTestDBLegacy(t)
	defer cleanup()

	// Create scheduler
	scheduler := NewScheduler(db)

	// Test creating default schedules for a new account
	accountID := 999 // Use a non-existent account ID

	// Initially, there should be no schedules
	schedules, err := scheduler.service.GetEmailSchedulesByAccount(accountID)
	if err != nil {
		t.Fatalf("Failed to get schedules: %v", err)
	}
	if len(schedules) != 0 {
		t.Errorf("Expected 0 schedules initially, got %d", len(schedules))
	}

	// Create default schedules
	if err := scheduler.createDefaultEmailSchedules(accountID); err != nil {
		t.Fatalf("Failed to create default schedules: %v", err)
	}

	// Check that schedules were created
	schedules, err = scheduler.service.GetEmailSchedulesByAccount(accountID)
	if err != nil {
		t.Fatalf("Failed to get schedules after creation: %v", err)
	}
	if len(schedules) != 2 {
		t.Errorf("Expected 2 schedules after creation, got %d", len(schedules))
	}

	// Verify the schedules are correct
	weeklyStockFound := false
	weeklySupplyChainFound := false

	for _, schedule := range schedules {
		if schedule.EmailType == models.EmailTypeWeeklyReport {
			weeklyStockFound = true
			if schedule.DayOfWeek == nil || *schedule.DayOfWeek != 1 {
				t.Error("Weekly stock schedule should be on Monday (day 1)")
			}
			if schedule.TimeOfDay != "09:00" {
				t.Error("Weekly stock schedule should be at 09:00")
			}
		} else if schedule.EmailType == models.EmailTypeWeeklySupplyChain {
			weeklySupplyChainFound = true
			if schedule.DayOfWeek == nil || *schedule.DayOfWeek != 2 {
				t.Error("Weekly supply chain schedule should be on Tuesday (day 2)")
			}
			if schedule.TimeOfDay != "09:00" {
				t.Error("Weekly supply chain schedule should be at 09:00")
			}
		}
	}

	if !weeklyStockFound {
		t.Error("Weekly stock schedule not found")
	}
	if !weeklySupplyChainFound {
		t.Error("Weekly supply chain schedule not found")
	}
}
