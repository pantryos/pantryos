package database

import (
	"errors"
	"time"

	"github.com/mnadev/stok/internal/models"
)

// Service combines all repositories and provides business logic
// This layer handles validation, business rules, and orchestrates operations
// between different entities while maintaining data consistency
type Service struct {
	organizations      OrganizationRepository
	accounts           AccountRepository
	users              UserRepository
	inventoryItems     InventoryItemRepository
	menuItems          MenuItemRepository
	deliveries         DeliveryRepository
	inventorySnapshots InventorySnapshotRepository
}

// NewService creates a new database service with all repositories
// This initializes all the repository interfaces needed for business operations
func NewService(db *DB) *Service {
	return &Service{
		organizations:      NewOrganizationRepository(db),
		accounts:           NewAccountRepository(db),
		users:              NewUserRepository(db),
		inventoryItems:     NewInventoryItemRepository(db),
		menuItems:          NewMenuItemRepository(db),
		deliveries:         NewDeliveryRepository(db),
		inventorySnapshots: NewInventorySnapshotRepository(db),
	}
}

// Organization operations
// These methods handle the top-level entity in the multi-tenant architecture

// CreateOrganization creates a new organization
// Organizations are the parent entities that contain multiple business accounts
func (s *Service) CreateOrganization(organization *models.Organization) error {
	return s.organizations.Create(organization)
}

// GetOrganization retrieves an organization by its ID
func (s *Service) GetOrganization(id int) (*models.Organization, error) {
	return s.organizations.GetByID(id)
}

// UpdateOrganization updates an existing organization
func (s *Service) UpdateOrganization(organization *models.Organization) error {
	return s.organizations.Update(organization)
}

// DeleteOrganization deletes an organization if it has no accounts
// This prevents orphaned data and maintains referential integrity
func (s *Service) DeleteOrganization(id int) error {
	// Check if organization has any accounts
	accounts, err := s.organizations.GetAccounts(id)
	if err != nil {
		return err
	}
	if len(accounts) > 0 {
		return errors.New("cannot delete organization with existing accounts")
	}
	return s.organizations.Delete(id)
}

// GetOrganizationAccounts retrieves all accounts belonging to an organization
func (s *Service) GetOrganizationAccounts(organizationID int) ([]models.Account, error) {
	return s.accounts.GetByOrganizationID(organizationID)
}

// Account operations
// These methods handle business locations within organizations

// CreateAccount creates a new account under an organization
// Validates that the parent organization exists before creating the account
func (s *Service) CreateAccount(account *models.Account) error {
	// Validate that the organization exists
	_, err := s.organizations.GetByID(account.OrganizationID)
	if err != nil {
		return errors.New("invalid organization ID")
	}
	return s.accounts.Create(account)
}

// GetAccount retrieves an account by its ID
func (s *Service) GetAccount(id int) (*models.Account, error) {
	return s.accounts.GetByID(id)
}

// GetAccountsByOrganization retrieves all accounts for a given organization
func (s *Service) GetAccountsByOrganization(organizationID int) ([]models.Account, error) {
	return s.accounts.GetByOrganizationID(organizationID)
}

// UpdateAccount updates an existing account
func (s *Service) UpdateAccount(account *models.Account) error {
	return s.accounts.Update(account)
}

// DeleteAccount deletes an account if it has no users
// This prevents orphaned data and maintains referential integrity
func (s *Service) DeleteAccount(id int) error {
	// Check if account has any users
	users, err := s.users.GetByAccountID(id)
	if err != nil {
		return err
	}
	if len(users) > 0 {
		return errors.New("cannot delete account with existing users")
	}
	return s.accounts.Delete(id)
}

// User operations
// These methods handle user management and authentication

// CreateUser creates a new user under an account
// Validates that the parent account exists and assigns a default role if none specified
func (s *Service) CreateUser(user *models.User) error {
	// Validate that the account exists
	_, err := s.accounts.GetByID(user.AccountID)
	if err != nil {
		return errors.New("invalid account ID")
	}

	// Set default role if not specified
	if user.Role == "" {
		user.Role = "user"
	}

	// Validate role
	if !isValidRole(user.Role) {
		return errors.New("invalid user role")
	}

	return s.users.Create(user)
}

// GetUser retrieves a user by their ID
func (s *Service) GetUser(id int) (*models.User, error) {
	return s.users.GetByID(id)
}

// GetUserByEmail retrieves a user by their email address
// Used for authentication and login operations
func (s *Service) GetUserByEmail(email string) (*models.User, error) {
	return s.users.GetByEmail(email)
}

// GetUsersByAccount retrieves all users belonging to a specific account
func (s *Service) GetUsersByAccount(accountID int) ([]models.User, error) {
	return s.users.GetByAccountID(accountID)
}

// GetUsersByOrganization retrieves all users across all accounts in an organization
// This is useful for organization-wide user management
func (s *Service) GetUsersByOrganization(organizationID int) ([]models.User, error) {
	return s.users.GetByOrganizationID(organizationID)
}

// UpdateUser updates an existing user
// Validates the role if it's being changed
func (s *Service) UpdateUser(user *models.User) error {
	if user.Role != "" && !isValidRole(user.Role) {
		return errors.New("invalid user role")
	}
	return s.users.Update(user)
}

// DeleteUser deletes a user by their ID
func (s *Service) DeleteUser(id int) error {
	return s.users.Delete(id)
}

// isValidRole checks if a role is valid
// Defines the allowed roles in the system
func isValidRole(role string) bool {
	validRoles := []string{"user", "manager", "admin", "org_admin"}
	for _, validRole := range validRoles {
		if role == validRole {
			return true
		}
	}
	return false
}

// Inventory operations
// These methods handle inventory item management

// CreateInventoryItem creates a new inventory item
// Validates that the parent account exists
func (s *Service) CreateInventoryItem(item *models.InventoryItem) error {
	// Validate that the account exists
	_, err := s.accounts.GetByID(item.AccountID)
	if err != nil {
		return errors.New("invalid account ID")
	}
	return s.inventoryItems.Create(item)
}

// GetInventoryItem retrieves an inventory item by its ID
func (s *Service) GetInventoryItem(id int) (*models.InventoryItem, error) {
	return s.inventoryItems.GetByID(id)
}

// GetInventoryItemsByAccount retrieves all inventory items for a specific account
func (s *Service) GetInventoryItemsByAccount(accountID int) ([]models.InventoryItem, error) {
	return s.inventoryItems.GetByAccountID(accountID)
}

// GetInventoryItemsByVendor retrieves inventory items from a specific vendor
func (s *Service) GetInventoryItemsByVendor(accountID int, vendor string) ([]models.InventoryItem, error) {
	return s.inventoryItems.GetByVendor(accountID, vendor)
}

// GetLowStockItems retrieves inventory items that are below their minimum stock level
// Useful for generating reorder alerts
func (s *Service) GetLowStockItems(accountID int) ([]models.InventoryItem, error) {
	return s.inventoryItems.GetLowStockItems(accountID)
}

// UpdateInventoryItem updates an existing inventory item
func (s *Service) UpdateInventoryItem(item *models.InventoryItem) error {
	return s.inventoryItems.Update(item)
}

// DeleteInventoryItem deletes an inventory item by its ID
func (s *Service) DeleteInventoryItem(id int) error {
	return s.inventoryItems.Delete(id)
}

// Menu operations
// These methods handle menu item management

// CreateMenuItem creates a new menu item
// Validates that the parent account exists
func (s *Service) CreateMenuItem(item *models.MenuItem) error {
	// Validate that the account exists
	_, err := s.accounts.GetByID(item.AccountID)
	if err != nil {
		return errors.New("invalid account ID")
	}
	return s.menuItems.Create(item)
}

// GetMenuItem retrieves a menu item by its ID
func (s *Service) GetMenuItem(id int) (*models.MenuItem, error) {
	return s.menuItems.GetByID(id)
}

// GetMenuItemWithIngredients retrieves a menu item along with its recipe ingredients
func (s *Service) GetMenuItemWithIngredients(id int) (*models.MenuItem, error) {
	return s.menuItems.GetWithIngredients(id)
}

// GetMenuItemsByAccount retrieves all menu items for a specific account
func (s *Service) GetMenuItemsByAccount(accountID int) ([]models.MenuItem, error) {
	return s.menuItems.GetByAccountID(accountID)
}

// GetMenuItemsByCategory retrieves menu items filtered by category
func (s *Service) GetMenuItemsByCategory(accountID int, category string) ([]models.MenuItem, error) {
	return s.menuItems.GetByCategory(accountID, category)
}

// UpdateMenuItem updates an existing menu item
func (s *Service) UpdateMenuItem(item *models.MenuItem) error {
	return s.menuItems.Update(item)
}

// DeleteMenuItem deletes a menu item by its ID
func (s *Service) DeleteMenuItem(id int) error {
	return s.menuItems.Delete(id)
}

// Delivery operations
// These methods handle delivery tracking and inventory replenishment

// CreateDelivery creates a new delivery record
// Validates that both the account and inventory item exist
func (s *Service) CreateDelivery(delivery *models.Delivery) error {
	// Validate that the account exists
	_, err := s.accounts.GetByID(delivery.AccountID)
	if err != nil {
		return errors.New("invalid account ID")
	}

	// Validate that the inventory item exists
	_, err = s.inventoryItems.GetByID(delivery.InventoryItemID)
	if err != nil {
		return errors.New("invalid inventory item ID")
	}

	return s.deliveries.Create(delivery)
}

// GetDelivery retrieves a delivery by its ID
func (s *Service) GetDelivery(id int) (*models.Delivery, error) {
	return s.deliveries.GetByID(id)
}

// GetDeliveriesByAccount retrieves all deliveries for a specific account
func (s *Service) GetDeliveriesByAccount(accountID int) ([]models.Delivery, error) {
	return s.deliveries.GetByAccountID(accountID)
}

// GetDeliveriesByDateRange retrieves deliveries within a specific date range
func (s *Service) GetDeliveriesByDateRange(accountID int, startDate, endDate time.Time) ([]models.Delivery, error) {
	return s.deliveries.GetByDateRange(accountID, startDate, endDate)
}

// GetDeliveriesByVendor retrieves deliveries from a specific vendor
func (s *Service) GetDeliveriesByVendor(accountID int, vendor string) ([]models.Delivery, error) {
	return s.deliveries.GetByVendor(accountID, vendor)
}

// UpdateDelivery updates an existing delivery
func (s *Service) UpdateDelivery(delivery *models.Delivery) error {
	return s.deliveries.Update(delivery)
}

// DeleteDelivery deletes a delivery by its ID
func (s *Service) DeleteDelivery(id int) error {
	return s.deliveries.Delete(id)
}

// Inventory Snapshot operations
// These methods handle historical inventory tracking

// CreateInventorySnapshot creates a new inventory snapshot
// Validates that the account exists and ensures the snapshot has valid data
func (s *Service) CreateInventorySnapshot(snapshot *models.InventorySnapshot) error {
	// Validate that the account exists
	_, err := s.accounts.GetByID(snapshot.AccountID)
	if err != nil {
		return errors.New("invalid account ID")
	}

	// Ensure timestamp is set
	if snapshot.Timestamp.IsZero() {
		snapshot.Timestamp = time.Now()
	}

	// Validate that counts map is not empty
	if len(snapshot.Counts) == 0 {
		return errors.New("snapshot must contain at least one inventory count")
	}

	return s.inventorySnapshots.Create(snapshot)
}

// GetInventorySnapshot retrieves an inventory snapshot by its ID
func (s *Service) GetInventorySnapshot(id int) (*models.InventorySnapshot, error) {
	return s.inventorySnapshots.GetByID(id)
}

// GetInventorySnapshotsByAccount retrieves all inventory snapshots for a specific account
func (s *Service) GetInventorySnapshotsByAccount(accountID int) ([]models.InventorySnapshot, error) {
	return s.inventorySnapshots.GetByAccountID(accountID)
}

// GetLatestInventorySnapshot retrieves the most recent inventory snapshot for an account
func (s *Service) GetLatestInventorySnapshot(accountID int) (*models.InventorySnapshot, error) {
	return s.inventorySnapshots.GetLatestByAccountID(accountID)
}

// GetInventorySnapshotsByDateRange retrieves inventory snapshots within a specific date range
func (s *Service) GetInventorySnapshotsByDateRange(accountID int, startDate, endDate time.Time) ([]models.InventorySnapshot, error) {
	return s.inventorySnapshots.GetByDateRange(accountID, startDate, endDate)
}

// UpdateInventorySnapshot updates an existing inventory snapshot
func (s *Service) UpdateInventorySnapshot(snapshot *models.InventorySnapshot) error {
	return s.inventorySnapshots.Update(snapshot)
}

// DeleteInventorySnapshot deletes an inventory snapshot by its ID
func (s *Service) DeleteInventorySnapshot(id int) error {
	return s.inventorySnapshots.Delete(id)
}

// GetInventoryVariance calculates the variance in inventory levels between two snapshots
// Useful for identifying inventory discrepancies and losses
func (s *Service) GetInventoryVariance(accountID int, startDate, endDate time.Time) (map[int]float64, error) {
	// This will be implemented when we add Toast POS integration
	// For now, return empty map
	return make(map[int]float64), nil
}

// Access Control operations
// These methods handle user permissions and access validation

// ValidateAccountAccess checks if a user has access to a specific account
// Users can only access resources within their own account
func (s *Service) ValidateAccountAccess(accountID int, userAccountID int) error {
	if accountID != userAccountID {
		return errors.New("access denied: user can only access their own account")
	}
	return nil
}

// ValidateOrganizationAccess checks if a user has access to a specific organization
// Users can access resources within their organization if they have appropriate permissions
func (s *Service) ValidateOrganizationAccess(organizationID int, userID int) (bool, error) {
	// Get the user
	user, err := s.users.GetByID(userID)
	if err != nil {
		return false, err
	}

	// Get the user's account
	account, err := s.accounts.GetByID(user.AccountID)
	if err != nil {
		return false, err
	}

	// Check if the user's account belongs to the organization
	if account.OrganizationID != organizationID {
		return false, errors.New("access denied: user does not belong to this organization")
	}

	return true, nil
}

// IsOrganizationAdmin checks if a user is an organization admin
// Organization admins have elevated permissions across all accounts in their organization
func (s *Service) IsOrganizationAdmin(userID int) (bool, error) {
	user, err := s.users.GetByID(userID)
	if err != nil {
		return false, err
	}

	return user.Role == "org_admin", nil
}
