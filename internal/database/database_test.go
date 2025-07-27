package database

import (
	"testing"
	"time"

	"github.com/mnadev/stok/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// This file uses the SetupTestDB function from test_setup.go

func TestDatabaseConnection(t *testing.T) {
	db, cleanup := SetupTestDB(t)
	defer cleanup()

	assert.NotNil(t, db)
}

func TestAccountOperations(t *testing.T) {
	db, cleanup := SetupTestDB(t)
	defer cleanup()

	service := NewService(db)

	// Create organization first
	org := createTestOrganization(t, service, "Test Corp")

	// Test creating an account
	account := &models.Account{
		OrganizationID: org.ID,
		Name:           "Test Coffee Shop",
		Location:       "123 Test St",
		Phone:          "555-1234",
		Email:          "test@coffee.com",
		Status:         "active",
	}

	err := service.CreateAccount(account)
	require.NoError(t, err)
	assert.NotZero(t, account.ID)

	// Test retrieving the account
	retrievedAccount, err := service.GetAccount(account.ID)
	require.NoError(t, err)
	assert.Equal(t, account.Name, retrievedAccount.Name)
	assert.Equal(t, account.ID, retrievedAccount.ID)
	assert.Equal(t, org.ID, retrievedAccount.OrganizationID)
}

func TestUserOperations(t *testing.T) {
	db, cleanup := SetupTestDB(t)
	defer cleanup()

	service := NewService(db)

	// Create organization and account
	org := createTestOrganization(t, service, "User Test Corp")
	account := createTestAccount(t, service, org.ID, "Test Coffee Shop")

	// Test creating a user
	user := &models.User{
		Email:     "test@example.com",
		Password:  "hashedpassword",
		AccountID: account.ID,
		Role:      "manager",
	}

	err := service.CreateUser(user)
	require.NoError(t, err)
	assert.NotZero(t, user.ID)
	assert.Equal(t, "manager", user.Role)

	// Test retrieving user by email
	retrievedUser, err := service.GetUserByEmail("test@example.com")
	require.NoError(t, err)
	assert.Equal(t, user.Email, retrievedUser.Email)
	assert.Equal(t, user.AccountID, retrievedUser.AccountID)
}

func TestInventoryItemOperations(t *testing.T) {
	db, cleanup := SetupTestDB(t)
	defer cleanup()

	service := NewService(db)

	// Create organization and account
	org := createTestOrganization(t, service, "Inventory Test Corp")
	account := createTestAccount(t, service, org.ID, "Test Coffee Shop")

	// Test creating an inventory item
	item := &models.InventoryItem{
		AccountID:       account.ID,
		Name:            "Coffee Beans",
		Unit:            "kg",
		CostPerUnit:     15.50,
		PreferredVendor: "Coffee Supply Co.",
		MinStockLevel:   10.0,
		MaxStockLevel:   50.0,
	}

	err := service.CreateInventoryItem(item)
	require.NoError(t, err)
	assert.NotZero(t, item.ID)

	// Test retrieving the item
	retrievedItem, err := service.GetInventoryItem(item.ID)
	require.NoError(t, err)
	assert.Equal(t, item.Name, retrievedItem.Name)
	assert.Equal(t, item.CostPerUnit, retrievedItem.CostPerUnit)
	assert.Equal(t, item.PreferredVendor, retrievedItem.PreferredVendor)
}

func TestMenuItemOperations(t *testing.T) {
	db, cleanup := SetupTestDB(t)
	defer cleanup()

	service := NewService(db)

	// Create organization and account
	org := createTestOrganization(t, service, "Menu Test Corp")
	account := createTestAccount(t, service, org.ID, "Test Coffee Shop")

	// Test creating a menu item
	menuItem := &models.MenuItem{
		AccountID: account.ID,
		Name:      "Espresso",
		Price:     3.50,
		Category:  "drinks",
	}

	err := service.CreateMenuItem(menuItem)
	require.NoError(t, err)
	assert.NotZero(t, menuItem.ID)

	// Test retrieving the menu item
	retrievedItem, err := service.GetMenuItem(menuItem.ID)
	require.NoError(t, err)
	assert.Equal(t, menuItem.Name, retrievedItem.Name)
	assert.Equal(t, menuItem.Price, retrievedItem.Price)
	assert.Equal(t, menuItem.Category, retrievedItem.Category)
}

func TestDeliveryOperations(t *testing.T) {
	db, cleanup := SetupTestDB(t)
	defer cleanup()

	service := NewService(db)

	// Create organization, account, and inventory item
	org := createTestOrganization(t, service, "Delivery Test Corp")
	account := createTestAccount(t, service, org.ID, "Test Coffee Shop")
	item := createTestInventoryItem(t, service, account.ID, "Coffee Beans")

	// Test creating a delivery
	delivery := &models.Delivery{
		AccountID:       account.ID,
		InventoryItemID: item.ID,
		Vendor:          "Coffee Supply Co.",
		Quantity:        25.0,
		DeliveryDate:    time.Now(),
		Cost:            387.50,
	}

	err := service.CreateDelivery(delivery)
	require.NoError(t, err)
	assert.NotZero(t, delivery.ID)

	// Test retrieving the delivery
	retrievedDelivery, err := service.GetDelivery(delivery.ID)
	require.NoError(t, err)
	assert.Equal(t, delivery.Vendor, retrievedDelivery.Vendor)
	assert.Equal(t, delivery.Quantity, retrievedDelivery.Quantity)
	assert.Equal(t, delivery.Cost, retrievedDelivery.Cost)
}

func TestInventorySnapshotOperations(t *testing.T) {
	db, cleanup := SetupTestDB(t)
	defer cleanup()

	service := NewService(db)

	// Create organization, account, and inventory items
	org := createTestOrganization(t, service, "Snapshot Test Corp")
	account := createTestAccount(t, service, org.ID, "Test Coffee Shop")
	item1 := createTestInventoryItem(t, service, account.ID, "Coffee Beans")
	item2 := createTestInventoryItem(t, service, account.ID, "Milk")

	// Test creating an inventory snapshot
	snapshot := &models.InventorySnapshot{
		AccountID: account.ID,
		Timestamp: time.Now(),
		Counts: map[int]float64{
			item1.ID: 15.5,
			item2.ID: 8.0,
		},
	}

	err := service.CreateInventorySnapshot(snapshot)
	require.NoError(t, err)
	assert.NotZero(t, snapshot.ID)

	// Test retrieving the snapshot
	retrievedSnapshot, err := service.GetInventorySnapshot(snapshot.ID)
	require.NoError(t, err)
	assert.Equal(t, snapshot.AccountID, retrievedSnapshot.AccountID)
	assert.Equal(t, len(snapshot.Counts), len(retrievedSnapshot.Counts))
}

func TestOrganizationOperations(t *testing.T) {
	db, cleanup := SetupTestDB(t)
	defer cleanup()

	service := NewService(db)

	// Test creating an organization
	org := &models.Organization{
		Name:        "Test Coffee Chain",
		Description: "A test coffee chain for testing",
	}

	err := service.CreateOrganization(org)
	require.NoError(t, err)
	assert.NotZero(t, org.ID)

	// Test retrieving the organization
	retrievedOrg, err := service.GetOrganization(org.ID)
	require.NoError(t, err)
	assert.Equal(t, org.Name, retrievedOrg.Name)
	assert.Equal(t, org.Description, retrievedOrg.Description)

	// Test creating accounts under the organization
	account1 := createTestAccount(t, service, org.ID, "Downtown Location")
	account2 := createTestAccount(t, service, org.ID, "Uptown Location")

	// Test getting all accounts for the organization
	accounts, err := service.GetOrganizationAccounts(org.ID)
	require.NoError(t, err)
	assert.Len(t, accounts, 2)

	// Verify account names and IDs
	names := make(map[string]bool)
	ids := make(map[int]bool)
	for _, acc := range accounts {
		names[acc.Name] = true
		ids[acc.ID] = true
	}
	assert.True(t, names["Downtown Location"])
	assert.True(t, names["Uptown Location"])
	assert.True(t, ids[account1.ID])
	assert.True(t, ids[account2.ID])
}

func TestRoleValidation(t *testing.T) {
	db, cleanup := SetupTestDB(t)
	defer cleanup()

	service := NewService(db)

	// Create organization and account
	org := createTestOrganization(t, service, "Role Test Corp")
	account := createTestAccount(t, service, org.ID, "Test Shop")

	// Test valid roles
	validRoles := []string{"user", "manager", "admin", "org_admin"}
	for _, role := range validRoles {
		user := &models.User{
			AccountID: account.ID,
			Email:     role + "@test.com",
			Password:  "hashedpassword",
			Role:      role,
		}
		err := service.CreateUser(user)
		assert.NoError(t, err, "Role %s should be valid", role)
	}

	// Test invalid role
	user := &models.User{
		AccountID: account.ID,
		Email:     "invalid@test.com",
		Password:  "hashedpassword",
		Role:      "invalid_role",
	}
	err := service.CreateUser(user)
	assert.Error(t, err, "Invalid role should be rejected")
}

func TestDefaultRoleAssignment(t *testing.T) {
	db, cleanup := SetupTestDB(t)
	defer cleanup()

	service := NewService(db)

	// Create organization and account
	org := createTestOrganization(t, service, "Default Role Corp")
	account := createTestAccount(t, service, org.ID, "Test Shop")

	// Create user without specifying role
	user := &models.User{
		AccountID: account.ID,
		Email:     "default@test.com",
		Password:  "hashedpassword",
		// Role is empty
	}

	err := service.CreateUser(user)
	require.NoError(t, err)
	assert.Equal(t, "user", user.Role, "Default role should be 'user'")
}
