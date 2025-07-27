package database

import (
	"testing"

	"github.com/mnadev/stok/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOrganizationStructure(t *testing.T) {
	t.Run("Create Organization", func(t *testing.T) {
		db, cleanup := SetupTestDB(t)
		defer cleanup()

		service := NewService(db)

		org := &models.Organization{
			Name:        "Coffee Chain Corp",
			Description: "A chain of coffee shops",
		}

		err := service.CreateOrganization(org)
		assert.NoError(t, err)
		assert.NotZero(t, org.ID)
		assert.NotZero(t, org.CreatedAt)
		assert.NotZero(t, org.UpdatedAt)
	})

	t.Run("Create Account Under Organization", func(t *testing.T) {
		db, cleanup := SetupTestDB(t)
		defer cleanup()

		service := NewService(db)

		// First create an organization
		org := &models.Organization{
			Name:        "Test Coffee Corp",
			Description: "Test organization",
		}
		err := service.CreateOrganization(org)
		require.NoError(t, err)

		// Create account under the organization
		account := &models.Account{
			OrganizationID: &org.ID,
			Name:           "Downtown Coffee Shop",
			Location:       "123 Main St, City, State",
			Phone:          "555-1234",
			Email:          "downtown@coffee.com",
			Status:         "active",
		}

		err = service.CreateAccount(account)
		assert.NoError(t, err)
		assert.NotZero(t, account.ID)
		assert.Equal(t, org.ID, *account.OrganizationID)
	})

	t.Run("Create User Under Account", func(t *testing.T) {
		db, cleanup := SetupTestDB(t)
		defer cleanup()

		service := NewService(db)

		// Create organization
		org := &models.Organization{
			Name: "User Test Corp",
		}
		err := service.CreateOrganization(org)
		require.NoError(t, err)

		// Create account
		account := &models.Account{
			OrganizationID: &org.ID,
			Name:           "Test Shop",
		}
		err = service.CreateAccount(account)
		require.NoError(t, err)

		// Create user
		user := &models.User{
			AccountID: account.ID,
			Email:     "manager@test.com",
			Password:  "hashedpassword",
			Role:      "manager",
		}

		err = service.CreateUser(user)
		assert.NoError(t, err)
		assert.NotZero(t, user.ID)
		assert.Equal(t, account.ID, user.AccountID)
		assert.Equal(t, "manager", user.Role)
	})

	t.Run("Get Organization Accounts", func(t *testing.T) {
		db, cleanup := SetupTestDB(t)
		defer cleanup()

		service := NewService(db)

		// Create organization
		org := &models.Organization{
			Name: "Multi Account Corp",
		}
		err := service.CreateOrganization(org)
		require.NoError(t, err)

		// Create multiple accounts
		account1 := &models.Account{
			OrganizationID: &org.ID,
			Name:           "Location 1",
		}
		account2 := &models.Account{
			OrganizationID: &org.ID,
			Name:           "Location 2",
		}

		err = service.CreateAccount(account1)
		require.NoError(t, err)
		err = service.CreateAccount(account2)
		require.NoError(t, err)

		// Get all accounts for the organization
		accounts, err := service.GetOrganizationAccounts(org.ID)
		assert.NoError(t, err)
		assert.Len(t, accounts, 2)

		// Verify account names
		names := make(map[string]bool)
		for _, acc := range accounts {
			names[acc.Name] = true
		}
		assert.True(t, names["Location 1"])
		assert.True(t, names["Location 2"])
	})

	t.Run("Get Users By Organization", func(t *testing.T) {
		db, cleanup := SetupTestDB(t)
		defer cleanup()

		service := NewService(db)

		// Create organization
		org := &models.Organization{
			Name: "User Org Test",
		}
		err := service.CreateOrganization(org)
		require.NoError(t, err)

		// Create accounts
		account1 := &models.Account{
			OrganizationID: &org.ID,
			Name:           "Shop 1",
		}
		account2 := &models.Account{
			OrganizationID: &org.ID,
			Name:           "Shop 2",
		}

		err = service.CreateAccount(account1)
		require.NoError(t, err)
		err = service.CreateAccount(account2)
		require.NoError(t, err)

		// Create users in different accounts
		user1 := &models.User{
			AccountID: account1.ID,
			Email:     "user1@test.com",
			Password:  "hashedpassword",
			Role:      "user",
		}
		user2 := &models.User{
			AccountID: account2.ID,
			Email:     "user2@test.com",
			Password:  "hashedpassword",
			Role:      "manager",
		}

		err = service.CreateUser(user1)
		require.NoError(t, err)
		err = service.CreateUser(user2)
		require.NoError(t, err)

		// Get all users in the organization
		users, err := service.GetUsersByOrganization(org.ID)
		assert.NoError(t, err)
		assert.Len(t, users, 2)

		// Verify user emails
		emails := make(map[string]bool)
		for _, user := range users {
			emails[user.Email] = true
		}
		assert.True(t, emails["user1@test.com"])
		assert.True(t, emails["user2@test.com"])
	})

	t.Run("Role Validation", func(t *testing.T) {
		db, cleanup := SetupTestDB(t)
		defer cleanup()

		service := NewService(db)

		// Create organization and account
		org := &models.Organization{Name: "Role Test Corp"}
		err := service.CreateOrganization(org)
		require.NoError(t, err)

		account := &models.Account{
			OrganizationID: &org.ID,
			Name:           "Test Shop",
		}
		err = service.CreateAccount(account)
		require.NoError(t, err)

		// Test valid roles
		validRoles := []string{"user", "manager", "admin", "org_admin"}
		for _, role := range validRoles {
			user := &models.User{
				AccountID: account.ID,
				Email:     role + "@test.com",
				Password:  "hashedpassword",
				Role:      role,
			}
			err = service.CreateUser(user)
			assert.NoError(t, err, "Role %s should be valid", role)
		}

		// Test invalid role
		user := &models.User{
			AccountID: account.ID,
			Email:     "invalid@test.com",
			Password:  "hashedpassword",
			Role:      "invalid_role",
		}
		err = service.CreateUser(user)
		assert.Error(t, err, "Invalid role should be rejected")
	})

	t.Run("Default Role Assignment", func(t *testing.T) {
		db, cleanup := SetupTestDB(t)
		defer cleanup()

		service := NewService(db)

		// Create organization and account
		org := &models.Organization{Name: "Default Role Corp"}
		err := service.CreateOrganization(org)
		require.NoError(t, err)

		account := &models.Account{
			OrganizationID: &org.ID,
			Name:           "Test Shop",
		}
		err = service.CreateAccount(account)
		require.NoError(t, err)

		// Create user without specifying role
		user := &models.User{
			AccountID: account.ID,
			Email:     "default@test.com",
			Password:  "hashedpassword",
			// Role is empty
		}

		err = service.CreateUser(user)
		assert.NoError(t, err)
		assert.Equal(t, "user", user.Role, "Default role should be 'user'")
	})

	t.Run("Organization Access Validation", func(t *testing.T) {
		db, cleanup := SetupTestDB(t)
		defer cleanup()

		service := NewService(db)

		// Create two organizations
		org1 := &models.Organization{Name: "Org 1"}
		org2 := &models.Organization{Name: "Org 2"}

		err := service.CreateOrganization(org1)
		require.NoError(t, err)
		err = service.CreateOrganization(org2)
		require.NoError(t, err)

		// Create accounts
		account1 := &models.Account{OrganizationID: &org1.ID, Name: "Shop 1"}
		account2 := &models.Account{OrganizationID: &org2.ID, Name: "Shop 2"}

		err = service.CreateAccount(account1)
		require.NoError(t, err)
		err = service.CreateAccount(account2)
		require.NoError(t, err)

		// Create users
		user1 := &models.User{AccountID: account1.ID, Email: "user1@test.com", Password: "hashedpassword"}
		user2 := &models.User{AccountID: account2.ID, Email: "user2@test.com", Password: "hashedpassword"}

		err = service.CreateUser(user1)
		require.NoError(t, err)
		err = service.CreateUser(user2)
		require.NoError(t, err)

		// Test access validation
		hasAccess, err := service.ValidateOrganizationAccess(org1.ID, user1.ID)
		assert.NoError(t, err, "User should have access to their organization")
		assert.True(t, hasAccess, "User should have access to their organization")

		hasAccess, err = service.ValidateOrganizationAccess(org1.ID, user2.ID)
		assert.Error(t, err, "User should not have access to different organization")
		assert.False(t, hasAccess, "User should not have access to different organization")
	})

	t.Run("Organization Admin Check", func(t *testing.T) {
		db, cleanup := SetupTestDB(t)
		defer cleanup()

		service := NewService(db)

		// Create organization and account
		org := &models.Organization{Name: "Admin Test Corp"}
		err := service.CreateOrganization(org)
		require.NoError(t, err)

		account := &models.Account{
			OrganizationID: &org.ID,
			Name:           "Test Shop",
		}
		err = service.CreateAccount(account)
		require.NoError(t, err)

		// Create regular user
		user := &models.User{
			AccountID: account.ID,
			Email:     "user@test.com",
			Password:  "hashedpassword",
			Role:      "user",
		}
		err = service.CreateUser(user)
		require.NoError(t, err)

		// Create org admin
		admin := &models.User{
			AccountID: account.ID,
			Email:     "admin@test.com",
			Password:  "hashedpassword",
			Role:      "org_admin",
		}
		err = service.CreateUser(admin)
		require.NoError(t, err)

		// Test admin check
		isAdmin, err := service.IsOrganizationAdmin(user.ID)
		assert.NoError(t, err)
		assert.False(t, isAdmin, "Regular user should not be org admin")

		isAdmin, err = service.IsOrganizationAdmin(admin.ID)
		assert.NoError(t, err)
		assert.True(t, isAdmin, "Org admin should be recognized as admin")
	})

	t.Run("Delete Protection", func(t *testing.T) {
		db, cleanup := SetupTestDB(t)
		defer cleanup()

		service := NewService(db)

		// Create organization
		org := &models.Organization{Name: "Delete Test Corp"}
		err := service.CreateOrganization(org)
		require.NoError(t, err)

		// Create account
		account := &models.Account{
			OrganizationID: &org.ID,
			Name:           "Test Shop",
		}
		err = service.CreateAccount(account)
		require.NoError(t, err)

		// Create user
		user := &models.User{
			AccountID: account.ID,
			Email:     "user@test.com",
			Password:  "hashedpassword",
		}
		err = service.CreateUser(user)
		require.NoError(t, err)

		// Try to delete account with users (should fail)
		err = service.DeleteAccount(account.ID)
		assert.Error(t, err, "Should not be able to delete account with users")

		// Delete user first
		err = service.DeleteUser(user.ID)
		assert.NoError(t, err)

		// Now delete account (should succeed)
		err = service.DeleteAccount(account.ID)
		assert.NoError(t, err)

		// Create a new account to test organization deletion protection
		account2 := &models.Account{
			OrganizationID: &org.ID,
			Name:           "Test Shop 2",
		}
		err = service.CreateAccount(account2)
		require.NoError(t, err)

		// Try to delete organization with accounts (should fail)
		err = service.DeleteOrganization(org.ID)
		assert.Error(t, err, "Should not be able to delete organization with accounts")

		// Delete the second account first
		err = service.DeleteAccount(account2.ID)
		assert.NoError(t, err)

		// Now delete organization (should succeed)
		err = service.DeleteOrganization(org.ID)
		assert.NoError(t, err)
	})
}

func TestInventoryWithOrganization(t *testing.T) {
	db, cleanup := SetupTestDB(t)
	defer cleanup()

	service := NewService(db)

	t.Run("Create Inventory Item Under Account", func(t *testing.T) {
		// Setup organization and account
		org := &models.Organization{Name: "Inventory Test Corp"}
		err := service.CreateOrganization(org)
		require.NoError(t, err)

		account := &models.Account{
			OrganizationID: &org.ID,
			Name:           "Test Shop",
		}
		err = service.CreateAccount(account)
		require.NoError(t, err)

		// Create inventory item
		item := &models.InventoryItem{
			AccountID:       account.ID,
			Name:            "Coffee Beans",
			Unit:            "kg",
			CostPerUnit:     15.50,
			PreferredVendor: "Coffee Supply Co.",
			MinStockLevel:   10.0,
			MaxStockLevel:   50.0,
		}

		err = service.CreateInventoryItem(item)
		assert.NoError(t, err)
		assert.NotZero(t, item.ID)
		assert.Equal(t, account.ID, item.AccountID)
	})

	t.Run("Get Low Stock Items", func(t *testing.T) {
		// Setup organization and account
		org := &models.Organization{Name: "Low Stock Test Corp"}
		err := service.CreateOrganization(org)
		require.NoError(t, err)

		account := &models.Account{
			OrganizationID: &org.ID,
			Name:           "Test Shop",
		}
		err = service.CreateAccount(account)
		require.NoError(t, err)

		// Create inventory items with min stock levels
		item1 := &models.InventoryItem{
			AccountID:     account.ID,
			Name:          "Coffee Beans",
			Unit:          "kg",
			MinStockLevel: 10.0,
		}
		item2 := &models.InventoryItem{
			AccountID:     account.ID,
			Name:          "Milk",
			Unit:          "liters",
			MinStockLevel: 5.0,
		}
		item3 := &models.InventoryItem{
			AccountID:     account.ID,
			Name:          "Sugar",
			Unit:          "kg",
			MinStockLevel: 0.0, // No min stock level
		}

		err = service.CreateInventoryItem(item1)
		require.NoError(t, err)
		err = service.CreateInventoryItem(item2)
		require.NoError(t, err)
		err = service.CreateInventoryItem(item3)
		require.NoError(t, err)

		// Get low stock items
		lowStockItems, err := service.GetLowStockItems(account.ID)
		assert.NoError(t, err)
		assert.Len(t, lowStockItems, 2) // Should only include items with min stock level > 0

		// Verify item names
		names := make(map[string]bool)
		for _, item := range lowStockItems {
			names[item.Name] = true
		}
		assert.True(t, names["Coffee Beans"])
		assert.True(t, names["Milk"])
		assert.False(t, names["Sugar"])
	})
}

func TestMenuWithOrganization(t *testing.T) {
	db, cleanup := SetupTestDB(t)
	defer cleanup()

	service := NewService(db)

	t.Run("Create Menu Item With Category", func(t *testing.T) {
		// Setup organization and account
		org := &models.Organization{Name: "Menu Test Corp"}
		err := service.CreateOrganization(org)
		require.NoError(t, err)

		account := &models.Account{
			OrganizationID: &org.ID,
			Name:           "Test Shop",
		}
		err = service.CreateAccount(account)
		require.NoError(t, err)

		// Create menu item
		menuItem := &models.MenuItem{
			AccountID: account.ID,
			Name:      "Espresso",
			Price:     3.50,
			Category:  "drinks",
		}

		err = service.CreateMenuItem(menuItem)
		assert.NoError(t, err)
		assert.NotZero(t, menuItem.ID)
		assert.Equal(t, account.ID, menuItem.AccountID)
		assert.Equal(t, "drinks", menuItem.Category)
	})

	t.Run("Get Menu Items By Category", func(t *testing.T) {
		// Setup organization and account
		org := &models.Organization{Name: "Category Test Corp"}
		err := service.CreateOrganization(org)
		require.NoError(t, err)

		account := &models.Account{
			OrganizationID: &org.ID,
			Name:           "Test Shop",
		}
		err = service.CreateAccount(account)
		require.NoError(t, err)

		// Create menu items in different categories
		drink1 := &models.MenuItem{
			AccountID: account.ID,
			Name:      "Espresso",
			Price:     3.50,
			Category:  "drinks",
		}
		drink2 := &models.MenuItem{
			AccountID: account.ID,
			Name:      "Cappuccino",
			Price:     4.50,
			Category:  "drinks",
		}
		food := &models.MenuItem{
			AccountID: account.ID,
			Name:      "Croissant",
			Price:     2.50,
			Category:  "food",
		}

		err = service.CreateMenuItem(drink1)
		require.NoError(t, err)
		err = service.CreateMenuItem(drink2)
		require.NoError(t, err)
		err = service.CreateMenuItem(food)
		require.NoError(t, err)

		// Get drinks
		drinks, err := service.GetMenuItemsByCategory(account.ID, "drinks")
		assert.NoError(t, err)
		assert.Len(t, drinks, 2)

		// Verify drink names
		names := make(map[string]bool)
		for _, drink := range drinks {
			names[drink.Name] = true
		}
		assert.True(t, names["Espresso"])
		assert.True(t, names["Cappuccino"])
		assert.False(t, names["Croissant"])
	})
}
