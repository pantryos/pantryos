package database

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/mnadev/pantryos/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	driverOnce sync.Once
)

// TestDBConfig holds configuration for test database setup
type TestDBConfig struct {
	MaxRetries       int
	RetryDelay       time.Duration
	MigrationTimeout time.Duration
	LogLevel         logger.LogLevel
}

// DefaultTestDBConfig returns a default configuration for test database setup
func DefaultTestDBConfig() *TestDBConfig {
	return &TestDBConfig{
		MaxRetries:       3,
		RetryDelay:       100 * time.Millisecond,
		MigrationTimeout: 30 * time.Second,
		LogLevel:         logger.Error,
	}
}

// SetupTestDB creates a new SQLite in-memory database for testing with robust error handling
// This function is used by all tests in this package and includes comprehensive error handling,
// retry logic, and graceful degradation mechanisms.
//
// Parameters:
//   - t: The testing context
//   - config: Optional configuration for database setup (uses defaults if nil)
//
// Returns:
//   - *DB: The initialized database connection
//   - func(): Cleanup function to be called after tests
//   - error: Any error that occurred during setup
func SetupTestDB(t *testing.T, config *TestDBConfig) (*DB, func(), error) {
	// Validate input parameters
	if t == nil {
		return nil, nil, fmt.Errorf("testing.T cannot be nil")
	}

	// Use default configuration if none provided
	if config == nil {
		config = DefaultTestDBConfig()
	}

	// Validate configuration
	if err := validateTestDBConfig(config); err != nil {
		return nil, nil, fmt.Errorf("invalid test database configuration: %w", err)
	}

	// Create database connection with retry logic
	gormDB, err := createTestDBConnection(config)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create test database connection: %w", err)
	}

	// Run migrations with timeout and retry logic
	if err := runTestMigrations(gormDB, config); err != nil {
		return nil, nil, fmt.Errorf("failed to run test migrations: %w", err)
	}

	db := &DB{DB: gormDB}

	// Return cleanup function with error handling
	cleanup := func() {
		if err := cleanupTestDB(db); err != nil {
			// Log cleanup errors but don't fail the test
			t.Logf("Warning: Test database cleanup failed: %v", err)
		}
	}

	return db, cleanup, nil
}

// createTestDBConnection creates a GORM database connection with retry logic
func createTestDBConnection(config *TestDBConfig) (*gorm.DB, error) {
	var gormDB *gorm.DB
	var lastErr error

	for attempt := 0; attempt <= config.MaxRetries; attempt++ {
		if attempt > 0 {
			time.Sleep(config.RetryDelay)
		}

		gormDB, lastErr = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: true,
			Logger:                                   logger.Default.LogMode(config.LogLevel),
		})

		if lastErr == nil {
			return gormDB, nil
		}
	}

	return nil, fmt.Errorf("failed to create database connection after %d attempts: %w", config.MaxRetries+1, lastErr)
}

// runTestMigrations runs database migrations with timeout and error handling
func runTestMigrations(gormDB *gorm.DB, config *TestDBConfig) error {
	ctx, cancel := context.WithTimeout(context.Background(), config.MigrationTimeout)
	defer cancel()

	// Define models to migrate with validation
	models := []interface{}{
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
		&models.EmailSchedule{},
	}

	// Run migrations with context
	errChan := make(chan error, 1)
	go func() {
		errChan <- gormDB.WithContext(ctx).AutoMigrate(models...)
	}()

	select {
	case err := <-errChan:
		if err != nil {
			return fmt.Errorf("migration failed: %w", err)
		}
	case <-ctx.Done():
		return fmt.Errorf("migration timed out after %v", config.MigrationTimeout)
	}

	return nil
}

// cleanupTestDB performs cleanup operations for the test database
func cleanupTestDB(db *DB) error {
	if db == nil || db.DB == nil {
		return nil
	}

	// Close the database connection gracefully
	sqlDB, err := db.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("failed to close database connection: %w", err)
	}

	return nil
}

// validateTestDBConfig validates the test database configuration
func validateTestDBConfig(config *TestDBConfig) error {
	if config == nil {
		return fmt.Errorf("configuration cannot be nil")
	}

	if config.MaxRetries < 0 {
		return fmt.Errorf("max retries must be non-negative")
	}

	if config.RetryDelay < 0 {
		return fmt.Errorf("retry delay must be non-negative")
	}

	if config.MigrationTimeout <= 0 {
		return fmt.Errorf("migration timeout must be positive")
	}

	return nil
}

// createTestOrganization creates a test organization with comprehensive error handling
func createTestOrganization(t *testing.T, service *Service, name string) (*models.Organization, error) {
	// Validate input parameters
	if t == nil {
		return nil, fmt.Errorf("testing.T cannot be nil")
	}
	if service == nil {
		return nil, fmt.Errorf("service cannot be nil")
	}
	if name == "" {
		return nil, fmt.Errorf("organization name cannot be empty")
	}

	// Sanitize and validate organization name
	if len(name) > 255 {
		return nil, fmt.Errorf("organization name too long (max 255 characters)")
	}

	org := &models.Organization{
		Name:        name,
		Description: "Test organization",
		Type:        "multi_location", // Set explicit type for consistency
	}

	if err := service.CreateOrganization(org); err != nil {
		return nil, fmt.Errorf("failed to create test organization '%s': %w", name, err)
	}

	// Verify the organization was created successfully
	if org.ID == 0 {
		return nil, fmt.Errorf("organization created but ID is zero")
	}

	return org, nil
}

// createTestAccount creates a test account under an organization with error handling
func createTestAccount(t *testing.T, service *Service, orgID int, name string) (*models.Account, error) {
	// Validate input parameters
	if t == nil {
		return nil, fmt.Errorf("testing.T cannot be nil")
	}
	if service == nil {
		return nil, fmt.Errorf("service cannot be nil")
	}
	if name == "" {
		return nil, fmt.Errorf("account name cannot be empty")
	}
	if orgID <= 0 {
		return nil, fmt.Errorf("organization ID must be positive")
	}

	// Validate that the organization exists
	org, err := service.GetOrganization(orgID)
	if err != nil {
		return nil, fmt.Errorf("failed to validate organization ID %d: %w", orgID, err)
	}
	if org == nil {
		return nil, fmt.Errorf("organization with ID %d does not exist", orgID)
	}

	// Sanitize and validate account name
	if len(name) > 255 {
		return nil, fmt.Errorf("account name too long (max 255 characters)")
	}

	account := &models.Account{
		OrganizationID: &orgID,
		Name:           name,
		Location:       "123 Test St",
		Phone:          "555-1234",
		Email:          name + "@test.com",
		Status:         "active",
		BusinessType:   "single_location",
	}

	if err := service.CreateAccount(account); err != nil {
		return nil, fmt.Errorf("failed to create test account '%s': %w", name, err)
	}

	// Verify the account was created successfully
	if account.ID == 0 {
		return nil, fmt.Errorf("account created but ID is zero")
	}

	return account, nil
}

// createTestStandaloneAccount creates a test standalone account with error handling
func createTestStandaloneAccount(t *testing.T, service *Service, name string) (*models.Account, error) {
	// Validate input parameters
	if t == nil {
		return nil, fmt.Errorf("testing.T cannot be nil")
	}
	if service == nil {
		return nil, fmt.Errorf("service cannot be nil")
	}
	if name == "" {
		return nil, fmt.Errorf("account name cannot be empty")
	}

	// Sanitize and validate account name
	if len(name) > 255 {
		return nil, fmt.Errorf("account name too long (max 255 characters)")
	}

	account := &models.Account{
		OrganizationID: nil, // Standalone account
		Name:           name,
		Location:       "123 Test St",
		Phone:          "555-1234",
		Email:          name + "@test.com",
		Status:         "active",
		BusinessType:   models.BusinessTypeSingleLocation,
	}

	if err := service.CreateAccount(account); err != nil {
		return nil, fmt.Errorf("failed to create test standalone account '%s': %w", name, err)
	}

	// Verify the account was created successfully
	if account.ID == 0 {
		return nil, fmt.Errorf("standalone account created but ID is zero")
	}

	return account, nil
}

// createTestUser creates a test user under an account with error handling
func createTestUser(t *testing.T, service *Service, accountID int, email, role string) (*models.User, error) {
	// Validate input parameters
	if t == nil {
		return nil, fmt.Errorf("testing.T cannot be nil")
	}
	if service == nil {
		return nil, fmt.Errorf("service cannot be nil")
	}
	if accountID <= 0 {
		return nil, fmt.Errorf("account ID must be positive")
	}
	if email == "" {
		return nil, fmt.Errorf("email cannot be empty")
	}
	if role == "" {
		return nil, fmt.Errorf("role cannot be empty")
	}

	// Validate email format (basic validation)
	if len(email) > 255 || !isValidEmailFormat(email) {
		return nil, fmt.Errorf("invalid email format: %s", email)
	}

	// Validate role
	validRoles := []string{"user", "manager", "admin", "org_admin"}
	if !isValidRoleInList(role, validRoles) {
		return nil, fmt.Errorf("invalid role '%s'. Valid roles: %v", role, validRoles)
	}

	// Validate that the account exists
	account, err := service.GetAccount(accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to validate account ID %d: %w", accountID, err)
	}
	if account == nil {
		return nil, fmt.Errorf("account with ID %d does not exist", accountID)
	}

	user := &models.User{
		AccountID: accountID,
		Email:     email,
		Password:  "hashedpassword",
		Role:      role,
	}

	if err := service.CreateUser(user); err != nil {
		return nil, fmt.Errorf("failed to create test user '%s': %w", email, err)
	}

	// Verify the user was created successfully
	if user.ID == 0 {
		return nil, fmt.Errorf("user created but ID is zero")
	}

	return user, nil
}

// createTestInventoryItem creates a test inventory item with error handling
func createTestInventoryItem(t *testing.T, service *Service, accountID int, name string) (*models.InventoryItem, error) {
	// Validate input parameters
	if t == nil {
		return nil, fmt.Errorf("testing.T cannot be nil")
	}
	if service == nil {
		return nil, fmt.Errorf("service cannot be nil")
	}
	if accountID <= 0 {
		return nil, fmt.Errorf("account ID must be positive")
	}
	if name == "" {
		return nil, fmt.Errorf("item name cannot be empty")
	}

	// Validate that the account exists
	account, err := service.GetAccount(accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to validate account ID %d: %w", accountID, err)
	}
	if account == nil {
		return nil, fmt.Errorf("account with ID %d does not exist", accountID)
	}

	// Sanitize and validate item name
	if len(name) > 255 {
		return nil, fmt.Errorf("item name too long (max 255 characters)")
	}

	item := &models.InventoryItem{
		AccountID:       accountID,
		Name:            name,
		Unit:            "kg",
		CostPerUnit:     10.0,
		PreferredVendor: "Test Vendor",
		MinStockLevel:   5.0,
		MaxStockLevel:   50.0,
	}

	if err := service.CreateInventoryItem(item); err != nil {
		return nil, fmt.Errorf("failed to create test inventory item '%s': %w", name, err)
	}

	// Verify the item was created successfully
	if item.ID == 0 {
		return nil, fmt.Errorf("inventory item created but ID is zero")
	}

	return item, nil
}

// createTestMenuItem creates a test menu item with error handling
func createTestMenuItem(t *testing.T, service *Service, accountID int, name string) (*models.MenuItem, error) {
	// Validate input parameters
	if t == nil {
		return nil, fmt.Errorf("testing.T cannot be nil")
	}
	if service == nil {
		return nil, fmt.Errorf("service cannot be nil")
	}
	if accountID <= 0 {
		return nil, fmt.Errorf("account ID must be positive")
	}
	if name == "" {
		return nil, fmt.Errorf("menu item name cannot be empty")
	}

	// Validate that the account exists
	account, err := service.GetAccount(accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to validate account ID %d: %w", accountID, err)
	}
	if account == nil {
		return nil, fmt.Errorf("account with ID %d does not exist", accountID)
	}

	// Sanitize and validate menu item name
	if len(name) > 255 {
		return nil, fmt.Errorf("menu item name too long (max 255 characters)")
	}

	item := &models.MenuItem{
		AccountID: accountID,
		Name:      name,
		Price:     5.0,
		Category:  "drinks",
	}

	if err := service.CreateMenuItem(item); err != nil {
		return nil, fmt.Errorf("failed to create test menu item '%s': %w", name, err)
	}

	// Verify the item was created successfully
	if item.ID == 0 {
		return nil, fmt.Errorf("menu item created but ID is zero")
	}

	return item, nil
}

// createTestDelivery creates a test delivery with error handling
func createTestDelivery(t *testing.T, service *Service, accountID, itemID int) (*models.Delivery, error) {
	// Validate input parameters
	if t == nil {
		return nil, fmt.Errorf("testing.T cannot be nil")
	}
	if service == nil {
		return nil, fmt.Errorf("service cannot be nil")
	}
	if accountID <= 0 {
		return nil, fmt.Errorf("account ID must be positive")
	}
	if itemID <= 0 {
		return nil, fmt.Errorf("inventory item ID must be positive")
	}

	// Validate that the account exists
	account, err := service.GetAccount(accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to validate account ID %d: %w", accountID, err)
	}
	if account == nil {
		return nil, fmt.Errorf("account with ID %d does not exist", accountID)
	}

	// Validate that the inventory item exists
	item, err := service.GetInventoryItem(itemID)
	if err != nil {
		return nil, fmt.Errorf("failed to validate inventory item ID %d: %w", itemID, err)
	}
	if item == nil {
		return nil, fmt.Errorf("inventory item with ID %d does not exist", itemID)
	}

	// Validate that the item belongs to the account
	if item.AccountID != accountID {
		return nil, fmt.Errorf("inventory item %d does not belong to account %d", itemID, accountID)
	}

	delivery := &models.Delivery{
		AccountID:       accountID,
		InventoryItemID: itemID,
		Vendor:          "Test Vendor",
		Quantity:        10.0,
		DeliveryDate:    time.Now(),
		Cost:            100.0,
	}

	if err := service.CreateDelivery(delivery); err != nil {
		return nil, fmt.Errorf("failed to create test delivery: %w", err)
	}

	// Verify the delivery was created successfully
	if delivery.ID == 0 {
		return nil, fmt.Errorf("delivery created but ID is zero")
	}

	return delivery, nil
}

// Helper functions for validation

// isValidEmailFormat performs basic email format validation
func isValidEmailFormat(email string) bool {
	if len(email) == 0 || len(email) > 255 {
		return false
	}

	// Basic email validation - check for @ symbol and domain
	atIndex := -1
	for i, char := range email {
		if char == '@' {
			if atIndex != -1 {
				return false // Multiple @ symbols
			}
			atIndex = i
		}
	}

	if atIndex == -1 || atIndex == 0 || atIndex == len(email)-1 {
		return false
	}

	// Check for domain part
	domain := email[atIndex+1:]
	if len(domain) == 0 || len(domain) > 255 {
		return false
	}

	return true
}

// isValidRoleInList checks if the role is in the list of valid roles
func isValidRoleInList(role string, validRoles []string) bool {
	for _, validRole := range validRoles {
		if role == validRole {
			return true
		}
	}
	return false
}

// Backward-compatible wrapper functions for existing tests
// These maintain the original API while the new functions provide better error handling

// SetupTestDBLegacy provides backward compatibility for existing tests
func SetupTestDBLegacy(t *testing.T) (*DB, func()) {
	// Use default configuration for backward compatibility
	config := DefaultTestDBConfig()

	// Create database connection with retry logic
	gormDB, err := createTestDBConnection(config)
	if err != nil {
		t.Fatalf("Failed to create test database connection: %v", err)
	}

	// Run migrations with timeout and retry logic
	if err := runTestMigrations(gormDB, config); err != nil {
		t.Fatalf("Failed to run test migrations: %v", err)
	}

	db := &DB{DB: gormDB}

	// Return cleanup function with error handling
	cleanup := func() {
		if err := cleanupTestDB(db); err != nil {
			// Log cleanup errors but don't fail the test
			t.Logf("Warning: Test database cleanup failed: %v", err)
		}
	}

	return db, cleanup
}

// createTestOrganizationLegacy provides backward compatibility
func createTestOrganizationLegacy(t *testing.T, service *Service, name string) *models.Organization {
	org, err := createTestOrganization(t, service, name)
	if err != nil {
		t.Fatalf("Failed to create test organization: %v", err)
	}
	return org
}

// createTestAccountLegacy provides backward compatibility
func createTestAccountLegacy(t *testing.T, service *Service, orgID int, name string) *models.Account {
	account, err := createTestAccount(t, service, orgID, name)
	if err != nil {
		t.Fatalf("Failed to create test account: %v", err)
	}
	return account
}

// createTestStandaloneAccountLegacy provides backward compatibility
func createTestStandaloneAccountLegacy(t *testing.T, service *Service, name string) *models.Account {
	account, err := createTestStandaloneAccount(t, service, name)
	if err != nil {
		t.Fatalf("Failed to create test standalone account: %v", err)
	}
	return account
}

// createTestUserLegacy provides backward compatibility
func createTestUserLegacy(t *testing.T, service *Service, accountID int, email, role string) *models.User {
	user, err := createTestUser(t, service, accountID, email, role)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}
	return user
}

// createTestInventoryItemLegacy provides backward compatibility
func createTestInventoryItemLegacy(t *testing.T, service *Service, accountID int, name string) *models.InventoryItem {
	item, err := createTestInventoryItem(t, service, accountID, name)
	if err != nil {
		t.Fatalf("Failed to create test inventory item: %v", err)
	}
	return item
}

// createTestMenuItemLegacy provides backward compatibility
func createTestMenuItemLegacy(t *testing.T, service *Service, accountID int, name string) *models.MenuItem {
	item, err := createTestMenuItem(t, service, accountID, name)
	if err != nil {
		t.Fatalf("Failed to create test menu item: %v", err)
	}
	return item
}

// createTestDeliveryLegacy provides backward compatibility
func createTestDeliveryLegacy(t *testing.T, service *Service, accountID, itemID int) *models.Delivery {
	delivery, err := createTestDelivery(t, service, accountID, itemID)
	if err != nil {
		t.Fatalf("Failed to create test delivery: %v", err)
	}
	return delivery
}
