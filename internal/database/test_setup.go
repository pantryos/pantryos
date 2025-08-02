package database

import (
	"sync"
	"testing"
	"time"

	"github.com/mnadev/stok/internal/models"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	driverOnce sync.Once
)

// SetupTestDB creates a new SQLite in-memory database for testing
// This function is used by all tests in this package
func SetupTestDB(t *testing.T) (*DB, func()) {
	// Open connection to SQLite in-memory
	// Each test gets its own isolated database

	// Create GORM DB with SQLite in-memory
	gormDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		// Disable features that we don't need for testing
		DisableForeignKeyConstraintWhenMigrating: true,
		Logger:                                   logger.Default.LogMode(logger.Error), // Reduce log noise
	})
	require.NoError(t, err)

	// Run migrations
	err = gormDB.AutoMigrate(
		&models.Organization{},
		&models.Account{},
		&models.User{},
		&models.Category{},
		&models.InventoryItem{},
		&models.MenuItem{},
		&models.RecipeIngredient{},
		&models.InventorySnapshot{},
		&models.Delivery{},
		&models.AccountInvitation{},
	)
	require.NoError(t, err)

	db := &DB{DB: gormDB}

	// Return cleanup function (SQLite in-memory doesn't need explicit cleanup)
	cleanup := func() {
		// No cleanup needed for in-memory SQLite
	}

	return db, cleanup
}

// createTestOrganization creates a test organization
func createTestOrganization(t *testing.T, service *Service, name string) *models.Organization {
	org := &models.Organization{
		Name:        name,
		Description: "Test organization",
	}
	err := service.CreateOrganization(org)
	require.NoError(t, err)
	return org
}

// createTestAccount creates a test account under an organization
func createTestAccount(t *testing.T, service *Service, orgID int, name string) *models.Account {
	account := &models.Account{
		OrganizationID: &orgID,
		Name:           name,
		Location:       "123 Test St",
		Phone:          "555-1234",
		Email:          name + "@test.com",
		Status:         "active",
	}
	err := service.CreateAccount(account)
	require.NoError(t, err)
	return account
}

// createTestStandaloneAccount creates a test standalone account (no organization)
func createTestStandaloneAccount(t *testing.T, service *Service, name string) *models.Account {
	account := &models.Account{
		OrganizationID: nil, // Standalone account
		Name:           name,
		Location:       "123 Test St",
		Phone:          "555-1234",
		Email:          name + "@test.com",
		Status:         "active",
		BusinessType:   models.BusinessTypeSingleLocation,
	}
	err := service.CreateAccount(account)
	require.NoError(t, err)
	return account
}

// createTestUser creates a test user under an account
func createTestUser(t *testing.T, service *Service, accountID int, email, role string) *models.User {
	user := &models.User{
		AccountID: accountID,
		Email:     email,
		Password:  "hashedpassword",
		Role:      role,
	}
	err := service.CreateUser(user)
	require.NoError(t, err)
	return user
}

// createTestInventoryItem creates a test inventory item
func createTestInventoryItem(t *testing.T, service *Service, accountID int, name string) *models.InventoryItem {
	item := &models.InventoryItem{
		AccountID:       accountID,
		Name:            name,
		Unit:            "kg",
		CostPerUnit:     10.0,
		PreferredVendor: "Test Vendor",
		MinStockLevel:   5.0,
		MaxStockLevel:   50.0,
	}
	err := service.CreateInventoryItem(item)
	require.NoError(t, err)
	return item
}

// createTestMenuItem creates a test menu item
func createTestMenuItem(t *testing.T, service *Service, accountID int, name string) *models.MenuItem {
	item := &models.MenuItem{
		AccountID: accountID,
		Name:      name,
		Price:     5.0,
		Category:  "drinks",
	}
	err := service.CreateMenuItem(item)
	require.NoError(t, err)
	return item
}

// createTestDelivery creates a test delivery
func createTestDelivery(t *testing.T, service *Service, accountID, itemID int) *models.Delivery {
	delivery := &models.Delivery{
		AccountID:       accountID,
		InventoryItemID: itemID,
		Vendor:          "Test Vendor",
		Quantity:        10.0,
		DeliveryDate:    time.Now(),
		Cost:            100.0,
	}
	err := service.CreateDelivery(delivery)
	require.NoError(t, err)
	return delivery
}
