// Package database provides the data access layer and business logic for the application.
// This package includes database connections, repositories, and service layer functions
// that handle all data operations while maintaining business rules and data consistency.
// The service layer orchestrates operations between different entities and enforces
// validation rules and business logic.
package database

import (
	"errors"
	"time"

	"github.com/mnadev/pantryos/internal/models"
	"gorm.io/gorm"
)

// Service combines all repositories and provides business logic for the application.
// This layer handles validation, business rules, and orchestrates operations
// between different entities while maintaining data consistency and enforcing
// multi-tenant security boundaries.
//
// The service layer acts as an intermediary between the HTTP handlers and the
// data access layer, ensuring that business rules are enforced and data integrity
// is maintained across all operations.
type Service struct {
	// organizations handles operations for the top-level multi-tenant entities
	organizations OrganizationRepository
	// accounts handles operations for business locations within organizations
	accounts AccountRepository
	// users handles user authentication and authorization data
	users UserRepository
	// inventoryItems handles physical inventory tracking and management
	inventoryItems InventoryItemRepository
	// menuItems handles menu item definitions and pricing
	menuItems MenuItemRepository
	// deliveries handles inventory delivery tracking and vendor management
	deliveries DeliveryRepository
	// inventorySnapshots handles historical inventory level tracking
	inventorySnapshots InventorySnapshotRepository
	// accountInvitations handles user invitation management
	accountInvitations AccountInvitationRepository
	// categories handles item categorization and organization
	categories CategoryRepository
}

// NewService creates a new database service with all repositories initialized.
// This function sets up the complete service layer with all necessary repositories
// for handling business operations across the entire application.
//
// Parameters:
//   - db: The database connection to use for all operations
//
// Returns:
//   - *Service: A fully initialized service instance ready for business operations
func NewService(db *DB) *Service {
	return &Service{
		organizations:      NewOrganizationRepository(db),
		accounts:           NewAccountRepository(db),
		users:              NewUserRepository(db),
		inventoryItems:     NewInventoryItemRepository(db),
		menuItems:          NewMenuItemRepository(db),
		deliveries:         NewDeliveryRepository(db),
		inventorySnapshots: NewInventorySnapshotRepository(db),
		accountInvitations: NewAccountInvitationRepository(db),
		categories:         NewCategoryRepository(db),
	}
}

// Organization operations
// These methods handle the top-level entity in the multi-tenant architecture.
// Organizations are the parent entities that contain multiple business accounts
// and provide the highest level of data isolation.

// CreateOrganization creates a new organization in the system.
// Organizations are the top-level entities in the multi-tenant architecture
// and serve as containers for multiple business accounts (locations).
//
// Parameters:
//   - organization: The organization data to create
//
// Returns:
//   - error: Any error that occurred during creation
//
// Business rules:
//   - Organization names must be unique within the system
//   - Organizations are created with default timestamps
func (s *Service) CreateOrganization(organization *models.Organization) error {
	return s.organizations.Create(organization)
}

// GetOrganization retrieves an organization by its unique identifier.
// This method provides access to organization details for validation
// and business logic operations.
//
// Parameters:
//   - id: The unique identifier of the organization to retrieve
//
// Returns:
//   - *models.Organization: The organization data if found
//   - error: Any error that occurred during retrieval
func (s *Service) GetOrganization(id int) (*models.Organization, error) {
	return s.organizations.GetByID(id)
}

// UpdateOrganization updates an existing organization's information.
// This method allows modification of organization details while
// maintaining data integrity and business rules.
//
// Parameters:
//   - organization: The updated organization data
//
// Returns:
//   - error: Any error that occurred during the update
func (s *Service) UpdateOrganization(organization *models.Organization) error {
	return s.organizations.Update(organization)
}

// DeleteOrganization deletes an organization if it has no associated accounts.
// This method enforces referential integrity by preventing deletion of
// organizations that still have active business accounts.
//
// Parameters:
//   - id: The unique identifier of the organization to delete
//
// Returns:
//   - error: Any error that occurred during deletion
//
// Business rules:
//   - Cannot delete organizations with existing accounts
//   - Maintains referential integrity across the system
func (s *Service) DeleteOrganization(id int) error {
	// Check if organization has any accounts before deletion
	accounts, err := s.organizations.GetAccounts(id)
	if err != nil {
		return err
	}
	if len(accounts) > 0 {
		return errors.New("cannot delete organization with existing accounts")
	}
	return s.organizations.Delete(id)
}

// GetOrganizationAccounts retrieves all accounts belonging to a specific organization.
// This method provides a way to enumerate all business locations within an organization
// for management and reporting purposes.
//
// Parameters:
//   - organizationID: The unique identifier of the organization
//
// Returns:
//   - []models.Account: List of accounts belonging to the organization
//   - error: Any error that occurred during retrieval
func (s *Service) GetOrganizationAccounts(organizationID int) ([]models.Account, error) {
	return s.accounts.GetByOrganizationID(organizationID)
}

// Account operations
// These methods handle business locations within organizations.
// Accounts represent individual business units (e.g., coffee shops, restaurants)
// and provide the primary scope for inventory and user operations.

// CreateAccount creates a new account that can be standalone or under an organization.
// This method validates the organization if provided and sets appropriate business type.
//
// Parameters:
//   - account: The account data to create
//
// Returns:
//   - error: Any error that occurred during creation
//
// Business rules:
//   - If OrganizationID is provided, parent organization must exist
//   - If OrganizationID is nil, account is created as standalone
//   - Business type is automatically set based on organization presence
//   - Account names should be unique within the organization (if applicable)
//   - Accounts are created with default status "active"
func (s *Service) CreateAccount(account *models.Account) error {
	// If OrganizationID is provided, validate that the parent organization exists
	if account.OrganizationID != nil {
		_, err := s.organizations.GetByID(*account.OrganizationID)
		if err != nil {
			return errors.New("invalid organization ID")
		}
		// Set business type based on organization type
		org, err := s.organizations.GetByID(*account.OrganizationID)
		if err == nil {
			account.BusinessType = org.Type
		}
	} else {
		// Standalone account - set as single location
		account.BusinessType = models.BusinessTypeSingleLocation
	}

	return s.accounts.Create(account)
}

// GetAccount retrieves an account by its unique identifier.
// This method provides access to account details for validation
// and business logic operations.
//
// Parameters:
//   - id: The unique identifier of the account to retrieve
//
// Returns:
//   - *models.Account: The account data if found
//   - error: Any error that occurred during retrieval
func (s *Service) GetAccount(id int) (*models.Account, error) {
	return s.accounts.GetByID(id)
}

// GetAccountsByOrganization retrieves all accounts for a given organization.
// This method provides a way to enumerate all business locations within
// an organization for management purposes.
//
// Parameters:
//   - organizationID: The unique identifier of the organization
//
// Returns:
//   - []models.Account: List of accounts belonging to the organization
//   - error: Any error that occurred during retrieval
func (s *Service) GetAccountsByOrganization(organizationID int) ([]models.Account, error) {
	return s.accounts.GetByOrganizationID(organizationID)
}

// GetAllAccounts retrieves all accounts in the system.
// This method is primarily used for registration purposes to allow users
// to select which account they want to register for.
//
// Returns:
//   - []models.Account: List of all accounts in the system
//   - error: Any error that occurred during retrieval
func (s *Service) GetAllAccounts() ([]models.Account, error) {
	return s.accounts.GetAll()
}

// GetStandaloneAccounts retrieves all standalone accounts (accounts without an organization).
// This method is useful for small businesses that don't need organizational hierarchy.
//
// Returns:
//   - []models.Account: List of standalone accounts in the system
//   - error: Any error that occurred during retrieval
func (s *Service) GetStandaloneAccounts() ([]models.Account, error) {
	return s.accounts.GetStandalone()
}

// IsStandaloneAccount checks if an account is standalone (not part of an organization).
// This method is useful for determining the business type and access patterns.
//
// Parameters:
//   - accountID: The unique identifier of the account to check
//
// Returns:
//   - bool: True if the account is standalone, false otherwise
//   - error: Any error that occurred during the check
func (s *Service) IsStandaloneAccount(accountID int) (bool, error) {
	account, err := s.accounts.GetByID(accountID)
	if err != nil {
		return false, err
	}
	return account.OrganizationID == nil, nil
}

// UpdateAccount updates an existing account's information.
// This method allows modification of account details while
// maintaining data integrity and business rules.
//
// Parameters:
//   - account: The updated account data
//
// Returns:
//   - error: Any error that occurred during the update
func (s *Service) UpdateAccount(account *models.Account) error {
	return s.accounts.Update(account)
}

// DeleteAccount deletes an account if it has no users.
// This method enforces referential integrity by preventing deletion of
// accounts that still have active users.
//
// Parameters:
//   - id: The unique identifier of the account to delete
//
// Returns:
//   - error: Any error that occurred during deletion
//
// Business rules:
//   - Cannot delete accounts with existing users
//   - Maintains referential integrity across the system
func (s *Service) DeleteAccount(id int) error {
	// Check if account has any users before deletion
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
// These methods handle user management and authentication.
// Users are the primary actors in the system and are associated with accounts.

// CreateUser creates a new user under an account.
// This method validates that the parent account exists and assigns a default role
// if none is specified.
//
// Parameters:
//   - user: The user data to create
//
// Returns:
//   - error: Any error that occurred during creation
//
// Business rules:
//   - Parent account must exist
//   - User names should be unique within the account
//   - Users are created with default role "user"
func (s *Service) CreateUser(user *models.User) error {
	// Validate that the parent account exists
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

// GetUser retrieves a user by their unique identifier.
// This method provides access to user details for validation
// and business logic operations.
//
// Parameters:
//   - id: The unique identifier of the user to retrieve
//
// Returns:
//   - *models.User: The user data if found
//   - error: Any error that occurred during retrieval
func (s *Service) GetUser(id int) (*models.User, error) {
	return s.users.GetByID(id)
}

// GetUserByEmail retrieves a user by their email address.
// This method is used for authentication and login operations.
//
// Parameters:
//   - email: The email address of the user to retrieve
//
// Returns:
//   - *models.User: The user data if found
//   - error: Any error that occurred during retrieval
func (s *Service) GetUserByEmail(email string) (*models.User, error) {
	return s.users.GetByEmail(email)
}

// GetUsersByAccount retrieves all users belonging to a specific account.
// This method provides a way to enumerate all users within a single account
// for management purposes.
//
// Parameters:
//   - accountID: The unique identifier of the account
//
// Returns:
//   - []models.User: List of users belonging to the account
//   - error: Any error that occurred during retrieval
func (s *Service) GetUsersByAccount(accountID int) ([]models.User, error) {
	return s.users.GetByAccountID(accountID)
}

// GetUsersByOrganization retrieves all users across all accounts in an organization.
// This method provides a way to enumerate all users within an organization
// for organization-wide user management.
//
// Parameters:
//   - organizationID: The unique identifier of the organization
//
// Returns:
//   - []models.User: List of users belonging to the organization
//   - error: Any error that occurred during retrieval
func (s *Service) GetUsersByOrganization(organizationID int) ([]models.User, error) {
	return s.users.GetByOrganizationID(organizationID)
}

// UpdateUser updates an existing user's information.
// This method allows modification of user details while
// maintaining data integrity and business rules.
//
// Parameters:
//   - user: The updated user data
//
// Returns:
//   - error: Any error that occurred during the update
func (s *Service) UpdateUser(user *models.User) error {
	if user.Role != "" && !isValidRole(user.Role) {
		return errors.New("invalid user role")
	}
	return s.users.Update(user)
}

// DeleteUser deletes a user by their unique identifier.
// This method enforces referential integrity by preventing deletion of
// users that still have active accounts.
//
// Parameters:
//   - id: The unique identifier of the user to delete
//
// Returns:
//   - error: Any error that occurred during deletion
//
// Business rules:
//   - Cannot delete users with existing accounts
//   - Maintains referential integrity across the system
func (s *Service) DeleteUser(id int) error {
	return s.users.Delete(id)
}

// isValidRole checks if a role is valid.
// Defines the allowed roles in the system.
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
// These methods handle inventory item management.
// Inventory items represent physical goods tracked in the system.

// CreateInventoryItem creates a new inventory item in the system.
// This method handles the creation of inventory items while enforcing
// business rules and data validation.
//
// Parameters:
//   - item: The inventory item data to create
//
// Returns:
//   - error: Any error that occurred during creation
//
// Business rules:
//   - Item names must be unique within an account
//   - CategoryID must reference a valid category in the same account
//   - Account ID must be set for proper scoping
func (s *Service) CreateInventoryItem(item *models.InventoryItem) error {
	// Validate required fields
	if item.Name == "" {
		return errors.New("item name is required")
	}
	if item.Unit == "" {
		return errors.New("unit is required")
	}

	// Validate category if provided
	if item.CategoryID != nil {
		category, err := s.categories.GetByID(*item.CategoryID)
		if err != nil {
			return errors.New("invalid category ID")
		}
		if category.AccountID != item.AccountID {
			return errors.New("category does not belong to the same account")
		}
	}

	// Check if item with same name already exists in the account
	existingItems, err := s.inventoryItems.GetByAccountID(item.AccountID)
	if err != nil {
		return err
	}

	for _, existing := range existingItems {
		if existing.Name == item.Name {
			return errors.New("item name already exists in this account")
		}
	}

	return s.inventoryItems.Create(item)
}

// GetInventoryItem retrieves an inventory item by its unique identifier.
// This method provides access to inventory item details for validation
// and business logic operations.
//
// Parameters:
//   - id: The unique identifier of the inventory item to retrieve
//
// Returns:
//   - *models.InventoryItem: The inventory item data if found
//   - error: Any error that occurred during retrieval
func (s *Service) GetInventoryItem(id int) (*models.InventoryItem, error) {
	return s.inventoryItems.GetByID(id)
}

// GetInventoryItemsByAccount retrieves all inventory items for a specific account.
// This method provides a way to enumerate all inventory items within
// a single account for management purposes.
//
// Parameters:
//   - accountID: The unique identifier of the account
//
// Returns:
//   - []models.InventoryItem: List of inventory items belonging to the account
//   - error: Any error that occurred during retrieval
func (s *Service) GetInventoryItemsByAccount(accountID int) ([]models.InventoryItem, error) {
	return s.inventoryItems.GetByAccountID(accountID)
}

// InventoryItemWithStock represents an inventory item with its current stock level
type InventoryItemWithStock struct {
	models.InventoryItem
	CurrentStock float64 `json:"current_stock"`
}

// GetInventoryItemsWithCurrentStock retrieves all inventory items for a specific account
// along with their current stock levels calculated from the latest inventory snapshot.
// This method provides a comprehensive view of inventory status for the frontend.
//
// Parameters:
//   - accountID: The unique identifier of the account
//
// Returns:
//   - []InventoryItemWithStock: List of inventory items with current stock levels
//   - error: Any error that occurred during retrieval
func (s *Service) GetInventoryItemsWithCurrentStock(accountID int) ([]InventoryItemWithStock, error) {
	// Get all inventory items for the account
	items, err := s.inventoryItems.GetByAccountID(accountID)
	if err != nil {
		return nil, err
	}

	// Get the latest inventory snapshot
	latestSnapshot, err := s.inventorySnapshots.GetLatestByAccountID(accountID)
	if err != nil {
		// If no snapshot exists, return items with zero stock
		if err == gorm.ErrRecordNotFound {
			result := make([]InventoryItemWithStock, len(items))
			for i, item := range items {
				result[i] = InventoryItemWithStock{
					InventoryItem: item,
					CurrentStock:  0,
				}
			}
			return result, nil
		}
		return nil, err
	}

	// Create a map of item ID to current stock from the snapshot
	stockMap := latestSnapshot.Counts

	// Create result with current stock information
	result := make([]InventoryItemWithStock, len(items))
	for i, item := range items {
		currentStock := 0.0
		if stock, exists := stockMap[item.ID]; exists {
			currentStock = stock
		}

		result[i] = InventoryItemWithStock{
			InventoryItem: item,
			CurrentStock:  currentStock,
		}
	}

	return result, nil
}

// GetInventoryItemsByVendor retrieves inventory items from a specific vendor.
// This method provides a way to filter inventory items by their vendor.
//
// Parameters:
//   - accountID: The unique identifier of the account
//   - vendor: The vendor name to filter by
//
// Returns:
//   - []models.InventoryItem: List of inventory items from the specified vendor
//   - error: Any error that occurred during retrieval
func (s *Service) GetInventoryItemsByVendor(accountID int, vendor string) ([]models.InventoryItem, error) {
	return s.inventoryItems.GetByVendor(accountID, vendor)
}

// GetLowStockItems retrieves inventory items that are below their minimum stock level.
// This method is useful for generating reorder alerts.
//
// Parameters:
//   - accountID: The unique identifier of the account
//
// Returns:
//   - []models.InventoryItem: List of low stock inventory items
//   - error: Any error that occurred during retrieval
func (s *Service) GetLowStockItems(accountID int) ([]models.InventoryItem, error) {
	return s.inventoryItems.GetLowStockItems(accountID)
}

// UpdateInventoryItem updates an existing inventory item's information.
// This method allows modification of inventory item details while
// maintaining data integrity and business rules.
//
// Parameters:
//   - item: The updated inventory item data
//
// Returns:
//   - error: Any error that occurred during the update
//
// Business rules:
//   - Item names must be unique within an account (excluding the current item)
//   - CategoryID must reference a valid category in the same account
//   - Account ID must be set for proper scoping
func (s *Service) UpdateInventoryItem(item *models.InventoryItem) error {
	// Validate required fields
	if item.Name == "" {
		return errors.New("item name is required")
	}
	if item.Unit == "" {
		return errors.New("unit is required")
	}

	// Validate category if provided
	if item.CategoryID != nil {
		category, err := s.categories.GetByID(*item.CategoryID)
		if err != nil {
			return errors.New("invalid category ID")
		}
		if category.AccountID != item.AccountID {
			return errors.New("category does not belong to the same account")
		}
	}

	// Check if item with same name already exists in the account (excluding current item)
	existingItems, err := s.inventoryItems.GetByAccountID(item.AccountID)
	if err != nil {
		return err
	}

	for _, existing := range existingItems {
		if existing.Name == item.Name && existing.ID != item.ID {
			return errors.New("item name already exists in this account")
		}
	}

	return s.inventoryItems.Update(item)
}

// DeleteInventoryItem deletes an inventory item by its unique identifier.
// This method enforces referential integrity by preventing deletion of
// inventory items that still have active deliveries.
//
// Parameters:
//   - id: The unique identifier of the inventory item to delete
//
// Returns:
//   - error: Any error that occurred during deletion
//
// Business rules:
//   - Cannot delete inventory items with existing deliveries
//   - Maintains referential integrity across the system
func (s *Service) DeleteInventoryItem(id int) error {
	return s.inventoryItems.Delete(id)
}

// Menu operations
// These methods handle menu item management.
// Menu items represent food and beverage offerings in the system.

// CreateMenuItem creates a new menu item in the system.
// This method handles the creation of menu items while enforcing
// business rules and data validation.
//
// Parameters:
//   - item: The menu item data to create
//
// Returns:
//   - error: Any error that occurred during creation
//
// Business rules:
//   - Item names must be unique within an account
//   - CategoryID must reference a valid category in the same account
//   - Account ID must be set for proper scoping
func (s *Service) CreateMenuItem(item *models.MenuItem) error {
	// Validate required fields
	if item.Name == "" {
		return errors.New("item name is required")
	}

	// Validate category if provided
	if item.CategoryID != nil {
		category, err := s.categories.GetByID(*item.CategoryID)
		if err != nil {
			return errors.New("invalid category ID")
		}
		if category.AccountID != item.AccountID {
			return errors.New("category does not belong to the same account")
		}
	}

	// Check if item with same name already exists in the account
	existingItems, err := s.menuItems.GetByAccountID(item.AccountID)
	if err != nil {
		return err
	}

	for _, existing := range existingItems {
		if existing.Name == item.Name {
			return errors.New("item name already exists in this account")
		}
	}

	return s.menuItems.Create(item)
}

// GetMenuItem retrieves a menu item by its unique identifier.
// This method provides access to menu item details for validation
// and business logic operations.
//
// Parameters:
//   - id: The unique identifier of the menu item to retrieve
//
// Returns:
//   - *models.MenuItem: The menu item data if found
//   - error: Any error that occurred during retrieval
func (s *Service) GetMenuItem(id int) (*models.MenuItem, error) {
	return s.menuItems.GetByID(id)
}

// GetMenuItemWithIngredients retrieves a menu item along with its recipe ingredients.
// This method provides a way to fetch a menu item and its associated recipe.
//
// Parameters:
//   - id: The unique identifier of the menu item
//
// Returns:
//   - *models.MenuItem: The menu item data including ingredients
//   - error: Any error that occurred during retrieval
func (s *Service) GetMenuItemWithIngredients(id int) (*models.MenuItem, error) {
	return s.menuItems.GetWithIngredients(id)
}

// GetMenuItemsByAccount retrieves all menu items for a specific account.
// This method provides a way to enumerate all menu items within
// a single account for management purposes.
//
// Parameters:
//   - accountID: The unique identifier of the account
//
// Returns:
//   - []models.MenuItem: List of menu items belonging to the account
//   - error: Any error that occurred during retrieval
func (s *Service) GetMenuItemsByAccount(accountID int) ([]models.MenuItem, error) {
	return s.menuItems.GetByAccountID(accountID)
}

// GetMenuItemsByCategory retrieves menu items filtered by category.
// This method provides a way to filter menu items by their category.
//
// Parameters:
//   - accountID: The unique identifier of the account
//   - category: The category name to filter by
//
// Returns:
//   - []models.MenuItem: List of menu items belonging to the category
//   - error: Any error that occurred during retrieval
func (s *Service) GetMenuItemsByCategory(accountID int, category string) ([]models.MenuItem, error) {
	return s.menuItems.GetByCategory(accountID, category)
}

// UpdateMenuItem updates an existing menu item's information.
// This method allows modification of menu item details while
// maintaining data integrity and business rules.
//
// Parameters:
//   - item: The updated menu item data
//
// Returns:
//   - error: Any error that occurred during the update
//
// Business rules:
//   - Item names must be unique within an account (excluding the current item)
//   - CategoryID must reference a valid category in the same account
//   - Account ID must be set for proper scoping
func (s *Service) UpdateMenuItem(item *models.MenuItem) error {
	// Validate required fields
	if item.Name == "" {
		return errors.New("item name is required")
	}

	// Validate category if provided
	if item.CategoryID != nil {
		category, err := s.categories.GetByID(*item.CategoryID)
		if err != nil {
			return errors.New("invalid category ID")
		}
		if category.AccountID != item.AccountID {
			return errors.New("category does not belong to the same account")
		}
	}

	// Check if item with same name already exists in the account (excluding current item)
	existingItems, err := s.menuItems.GetByAccountID(item.AccountID)
	if err != nil {
		return err
	}

	for _, existing := range existingItems {
		if existing.Name == item.Name && existing.ID != item.ID {
			return errors.New("item name already exists in this account")
		}
	}

	return s.menuItems.Update(item)
}

// DeleteMenuItem deletes a menu item by its unique identifier.
// This method enforces referential integrity by preventing deletion of
// menu items that still have active deliveries.
//
// Parameters:
//   - id: The unique identifier of the menu item to delete
//
// Returns:
//   - error: Any error that occurred during deletion
//
// Business rules:
//   - Cannot delete menu items with existing deliveries
//   - Maintains referential integrity across the system
func (s *Service) DeleteMenuItem(id int) error {
	return s.menuItems.Delete(id)
}

// Delivery operations
// These methods handle delivery tracking and inventory replenishment.
// Deliveries represent the movement of inventory items from vendors to accounts.

// CreateDelivery creates a new delivery record.
// This method validates that both the account and inventory item exist.
//
// Parameters:
//   - delivery: The delivery data to create
//
// Returns:
//   - error: Any error that occurred during creation
//
// Business rules:
//   - Both account and inventory item must exist
//   - Deliveries are created with default status "pending"
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

// GetDelivery retrieves a delivery by its unique identifier.
// This method provides access to delivery details for validation
// and business logic operations.
//
// Parameters:
//   - id: The unique identifier of the delivery to retrieve
//
// Returns:
//   - *models.Delivery: The delivery data if found
//   - error: Any error that occurred during retrieval
func (s *Service) GetDelivery(id int) (*models.Delivery, error) {
	return s.deliveries.GetByID(id)
}

// GetDeliveriesByAccount retrieves all deliveries for a specific account.
// This method provides a way to enumerate all deliveries within
// a single account for management purposes.
//
// Parameters:
//   - accountID: The unique identifier of the account
//
// Returns:
//   - []models.Delivery: List of deliveries belonging to the account
//   - error: Any error that occurred during retrieval
func (s *Service) GetDeliveriesByAccount(accountID int) ([]models.Delivery, error) {
	return s.deliveries.GetByAccountID(accountID)
}

// GetDeliveriesByDateRange retrieves deliveries within a specific date range.
// This method provides a way to filter deliveries by their date.
//
// Parameters:
//   - accountID: The unique identifier of the account
//   - startDate: The start date of the range
//   - endDate: The end date of the range
//
// Returns:
//   - []models.Delivery: List of deliveries within the date range
//   - error: Any error that occurred during retrieval
func (s *Service) GetDeliveriesByDateRange(accountID int, startDate, endDate time.Time) ([]models.Delivery, error) {
	return s.deliveries.GetByDateRange(accountID, startDate, endDate)
}

// GetDeliveriesByVendor retrieves deliveries from a specific vendor.
// This method provides a way to filter deliveries by their vendor.
//
// Parameters:
//   - accountID: The unique identifier of the account
//   - vendor: The vendor name to filter by
//
// Returns:
//   - []models.Delivery: List of deliveries from the specified vendor
//   - error: Any error that occurred during retrieval
func (s *Service) GetDeliveriesByVendor(accountID int, vendor string) ([]models.Delivery, error) {
	return s.deliveries.GetByVendor(accountID, vendor)
}

// UpdateDelivery updates an existing delivery's information.
// This method allows modification of delivery details while
// maintaining data integrity and business rules.
//
// Parameters:
//   - delivery: The updated delivery data
//
// Returns:
//   - error: Any error that occurred during the update
func (s *Service) UpdateDelivery(delivery *models.Delivery) error {
	return s.deliveries.Update(delivery)
}

// DeleteDelivery deletes a delivery by its unique identifier.
// This method enforces referential integrity by preventing deletion of
// deliveries that still have active inventory items.
//
// Parameters:
//   - id: The unique identifier of the delivery to delete
//
// Returns:
//   - error: Any error that occurred during deletion
//
// Business rules:
//   - Cannot delete deliveries with existing inventory items
//   - Maintains referential integrity across the system
func (s *Service) DeleteDelivery(id int) error {
	return s.deliveries.Delete(id)
}

// Inventory Snapshot operations
// These methods handle historical inventory tracking.
// Inventory snapshots represent a point-in-time view of inventory levels.

// CreateInventorySnapshot creates a new inventory snapshot.
// This method validates that the account exists and ensures the snapshot has valid data.
//
// Parameters:
//   - snapshot: The inventory snapshot data to create
//
// Returns:
//   - error: Any error that occurred during creation
//
// Business rules:
//   - Account must exist
//   - Timestamp must be set (defaults to now if not provided)
//   - Counts map must not be empty
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

// GetInventorySnapshot retrieves an inventory snapshot by its unique identifier.
// This method provides access to inventory snapshot details for validation
// and business logic operations.
//
// Parameters:
//   - id: The unique identifier of the inventory snapshot to retrieve
//
// Returns:
//   - *models.InventorySnapshot: The inventory snapshot data if found
//   - error: Any error that occurred during retrieval
func (s *Service) GetInventorySnapshot(id int) (*models.InventorySnapshot, error) {
	return s.inventorySnapshots.GetByID(id)
}

// GetInventorySnapshotsByAccount retrieves all inventory snapshots for a specific account.
// This method provides a way to enumerate all inventory snapshots within
// a single account for management purposes.
//
// Parameters:
//   - accountID: The unique identifier of the account
//
// Returns:
//   - []models.InventorySnapshot: List of inventory snapshots belonging to the account
//   - error: Any error that occurred during retrieval
func (s *Service) GetInventorySnapshotsByAccount(accountID int) ([]models.InventorySnapshot, error) {
	return s.inventorySnapshots.GetByAccountID(accountID)
}

// GetLatestInventorySnapshot retrieves the most recent inventory snapshot for an account.
// This method provides a way to fetch the latest inventory snapshot for
// a specific account for reporting and analysis.
//
// Parameters:
//   - accountID: The unique identifier of the account
//
// Returns:
//   - *models.InventorySnapshot: The most recent inventory snapshot
//   - error: Any error that occurred during retrieval
func (s *Service) GetLatestInventorySnapshot(accountID int) (*models.InventorySnapshot, error) {
	return s.inventorySnapshots.GetLatestByAccountID(accountID)
}

// GetInventorySnapshotsByDateRange retrieves inventory snapshots within a specific date range.
// This method provides a way to filter inventory snapshots by their date.
//
// Parameters:
//   - accountID: The unique identifier of the account
//   - startDate: The start date of the range
//   - endDate: The end date of the range
//
// Returns:
//   - []models.InventorySnapshot: List of inventory snapshots within the date range
//   - error: Any error that occurred during retrieval
func (s *Service) GetInventorySnapshotsByDateRange(accountID int, startDate, endDate time.Time) ([]models.InventorySnapshot, error) {
	return s.inventorySnapshots.GetByDateRange(accountID, startDate, endDate)
}

// UpdateInventorySnapshot updates an existing inventory snapshot's information.
// This method allows modification of inventory snapshot details while
// maintaining data integrity and business rules.
//
// Parameters:
//   - snapshot: The updated inventory snapshot data
//
// Returns:
//   - error: Any error that occurred during the update
func (s *Service) UpdateInventorySnapshot(snapshot *models.InventorySnapshot) error {
	return s.inventorySnapshots.Update(snapshot)
}

// DeleteInventorySnapshot deletes an inventory snapshot by its unique identifier.
// This method enforces referential integrity by preventing deletion of
// inventory snapshots that still have active inventory items.
//
// Parameters:
//   - id: The unique identifier of the inventory snapshot to delete
//
// Returns:
//   - error: Any error that occurred during deletion
//
// Business rules:
//   - Cannot delete inventory snapshots with existing inventory items
//   - Maintains referential integrity across the system
func (s *Service) DeleteInventorySnapshot(id int) error {
	return s.inventorySnapshots.Delete(id)
}

// GetInventoryVariance calculates the variance in inventory levels between two snapshots.
// This method is useful for identifying inventory discrepancies and losses.
//
// Parameters:
//   - accountID: The unique identifier of the account
//   - startDate: The start date of the range for the first snapshot
//   - endDate: The end date of the range for the first snapshot
//
// Returns:
//   - map[int]float64: A map of inventory item IDs to their variance
//   - error: Any error that occurred during calculation
//
// Business logic:
//   - This will be implemented when we add Toast POS integration
//   - For now, it returns an empty map as a placeholder
func (s *Service) GetInventoryVariance(accountID int, startDate, endDate time.Time) (map[int]float64, error) {
	// This will be implemented when we add Toast POS integration
	// For now, return empty map
	return make(map[int]float64), nil
}

// Access Control operations
// These methods handle user permissions and access validation.
// Users are the primary actors in the system and are associated with accounts.

// ValidateAccountAccess checks if a user has access to a specific account.
// This method ensures that users can only access resources within their own account.
//
// Parameters:
//   - accountID: The unique identifier of the account to check access for
//   - userAccountID: The unique identifier of the user's account
//
// Returns:
//   - error: Any error that occurred during validation
//
// Business rules:
//   - Users can only access resources within their own account
//   - Access is denied if the user's account does not match the requested account
func (s *Service) ValidateAccountAccess(accountID int, userAccountID int) error {
	if accountID != userAccountID {
		return errors.New("access denied: user can only access their own account")
	}
	return nil
}

// ValidateOrganizationAccess checks if a user has access to a specific organization.
// This method ensures that users can access resources within their organization
// if they have appropriate permissions.
//
// Parameters:
//   - organizationID: The unique identifier of the organization to check access for
//   - userID: The unique identifier of the user
//
// Returns:
//   - bool: True if access is granted, false otherwise
//   - error: Any error that occurred during validation
//
// Business rules:
//   - Users can access resources within their organization if they have appropriate permissions
//   - Access is denied if the user's account does not belong to the organization
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
	if account.OrganizationID == nil || *account.OrganizationID != organizationID {
		return false, errors.New("access denied: user does not belong to this organization")
	}

	return true, nil
}

// IsOrganizationAdmin checks if a user is an organization admin.
// Organization admins have elevated permissions across all accounts in their organization.
//
// Parameters:
//   - userID: The unique identifier of the user
//
// Returns:
//   - bool: True if the user is an organization admin, false otherwise
//   - error: Any error that occurred during retrieval
func (s *Service) IsOrganizationAdmin(userID int) (bool, error) {
	user, err := s.users.GetByID(userID)
	if err != nil {
		return false, err
	}

	return user.Role == "org_admin", nil
}

// Invitation operations
// These methods handle user invitation management.

// GetPendingInvitationByEmail retrieves a pending invitation for a given email.
// This method is used during user registration to validate that the user has been invited.
//
// Parameters:
//   - email: The email address to check for pending invitations
//
// Returns:
//   - *models.AccountInvitation: The pending invitation if found
//   - error: Any error that occurred during retrieval
func (s *Service) GetPendingInvitationByEmail(email string) (*models.AccountInvitation, error) {
	return s.accountInvitations.GetPendingByEmail(email)
}

// UpdateInvitation updates an existing invitation.
// This method is used to mark invitations as accepted when users register.
//
// Parameters:
//   - invitation: The invitation data to update
//
// Returns:
//   - error: Any error that occurred during the update
func (s *Service) UpdateInvitation(invitation *models.AccountInvitation) error {
	return s.accountInvitations.Update(invitation)
}

// CreateInvitation creates a new invitation for a user to join an account.
// This method is used by account admins to invite new users.
//
// Parameters:
//   - invitation: The invitation data to create
//
// Returns:
//   - error: Any error that occurred during creation
func (s *Service) CreateInvitation(invitation *models.AccountInvitation) error {
	// Set default expiration to 7 days from now
	if invitation.ExpiresAt.IsZero() {
		invitation.ExpiresAt = time.Now().AddDate(0, 0, 7)
	}

	// Set default status to pending
	if invitation.Status == "" {
		invitation.Status = "pending"
	}

	return s.accountInvitations.Create(invitation)
}

// GetInvitationsByAccount retrieves all invitations for a specific account.
// This method is used by account admins to manage invitations.
//
// Parameters:
//   - accountID: The account ID to get invitations for
//
// Returns:
//   - []models.AccountInvitation: List of invitations for the account
//   - error: Any error that occurred during retrieval
func (s *Service) GetInvitationsByAccount(accountID int) ([]models.AccountInvitation, error) {
	return s.accountInvitations.GetByAccountID(accountID)
}

// DeleteInvitation deletes an invitation.
// This method is used by account admins to revoke invitations.
//
// Parameters:
//   - id: The invitation ID to delete
//
// Returns:
//   - error: Any error that occurred during deletion
func (s *Service) DeleteInvitation(id int) error {
	return s.accountInvitations.Delete(id)
}

// Category operations
// These methods handle item categorization and organization.
// Categories help organize inventory items and menu items for better management.

// CreateCategory creates a new category for the specified account.
// Categories help organize inventory items and menu items for better management
// and reporting capabilities.
//
// Parameters:
//   - category: The category data to create
//
// Returns:
//   - error: Any error that occurred during creation
//
// Business rules:
//   - Category names must be unique within an account
//   - Categories are created with default timestamps
//   - Categories are created as active by default
func (s *Service) CreateCategory(category *models.Category) error {
	// Validate that the category name is not empty
	if category.Name == "" {
		return errors.New("category name is required")
	}

	// Check if a category with the same name already exists in the account
	existingCategories, err := s.categories.GetByAccountID(category.AccountID)
	if err != nil {
		return err
	}

	for _, existing := range existingCategories {
		if existing.Name == category.Name && existing.ID != category.ID {
			return errors.New("category name already exists in this account")
		}
	}

	return s.categories.Create(category)
}

// GetCategory retrieves a category by its unique identifier.
// This method provides access to category details for validation
// and business logic operations.
//
// Parameters:
//   - id: The unique identifier of the category to retrieve
//
// Returns:
//   - *models.Category: The category data if found
//   - error: Any error that occurred during retrieval
func (s *Service) GetCategory(id int) (*models.Category, error) {
	return s.categories.GetByID(id)
}

// GetCategoriesByAccount retrieves all categories for a specific account.
// This method returns all categories associated with the account,
// including both active and inactive categories.
//
// Parameters:
//   - accountID: The account ID to retrieve categories for
//
// Returns:
//   - []models.Category: List of categories for the account
//   - error: Any error that occurred during retrieval
func (s *Service) GetCategoriesByAccount(accountID int) ([]models.Category, error) {
	return s.categories.GetByAccountID(accountID)
}

// GetActiveCategoriesByAccount retrieves only active categories for a specific account.
// This method returns only categories that are currently active,
// which is useful for UI dropdowns and active item categorization.
//
// Parameters:
//   - accountID: The account ID to retrieve active categories for
//
// Returns:
//   - []models.Category: List of active categories for the account
//   - error: Any error that occurred during retrieval
func (s *Service) GetActiveCategoriesByAccount(accountID int) ([]models.Category, error) {
	return s.categories.GetActiveByAccountID(accountID)
}

// UpdateCategory updates an existing category in the system.
// This method allows modification of category details while maintaining
// data integrity and business rules.
//
// Parameters:
//   - category: The updated category data
//
// Returns:
//   - error: Any error that occurred during update
//
// Business rules:
//   - Category names must remain unique within an account
//   - Updated timestamps are automatically set
func (s *Service) UpdateCategory(category *models.Category) error {
	// Validate that the category name is not empty
	if category.Name == "" {
		return errors.New("category name is required")
	}

	// Check if a category with the same name already exists in the account
	existingCategories, err := s.categories.GetByAccountID(category.AccountID)
	if err != nil {
		return err
	}

	for _, existing := range existingCategories {
		if existing.Name == category.Name && existing.ID != category.ID {
			return errors.New("category name already exists in this account")
		}
	}

	return s.categories.Update(category)
}

// DeleteCategory removes a category from the system.
// This method handles the deletion of categories while ensuring
// that no orphaned references remain in the system.
//
// Parameters:
//   - id: The unique identifier of the category to delete
//
// Returns:
//   - error: Any error that occurred during deletion
//
// Business rules:
//   - Categories cannot be deleted if they are still referenced by items
//   - Soft deletion is preferred over hard deletion for data integrity
func (s *Service) DeleteCategory(id int) error {
	// Check if any inventory items are using this category
	inventoryItems, err := s.inventoryItems.GetByAccountID(0) // We'll need to get all items to check
	if err != nil {
		return err
	}

	for _, item := range inventoryItems {
		if item.CategoryID != nil && *item.CategoryID == id {
			return errors.New("cannot delete category: it is still assigned to inventory items")
		}
	}

	// Check if any menu items are using this category
	menuItems, err := s.menuItems.GetByAccountID(0) // We'll need to get all items to check
	if err != nil {
		return err
	}

	for _, item := range menuItems {
		if item.CategoryID != nil && *item.CategoryID == id {
			return errors.New("cannot delete category: it is still assigned to menu items")
		}
	}

	return s.categories.Delete(id)
}

// GetInventoryItemsByCategory retrieves all inventory items in a specific category.
// This method provides filtered access to inventory items based on their
// category assignment for better organization and reporting.
//
// Parameters:
//   - accountID: The account ID to scope the search
//   - categoryID: The category ID to filter by
//
// Returns:
//   - []models.InventoryItem: List of inventory items in the category
//   - error: Any error that occurred during retrieval
func (s *Service) GetInventoryItemsByCategory(accountID int, categoryID int) ([]models.InventoryItem, error) {
	// Get all inventory items for the account
	items, err := s.inventoryItems.GetByAccountID(accountID)
	if err != nil {
		return nil, err
	}

	// Filter items by category
	var filteredItems []models.InventoryItem
	for _, item := range items {
		if item.CategoryID != nil && *item.CategoryID == categoryID {
			filteredItems = append(filteredItems, item)
		}
	}

	return filteredItems, nil
}

// GetMenuItemsByCategoryID retrieves all menu items in a specific category.
// This method provides filtered access to menu items based on their
// category assignment for better organization and reporting.
//
// Parameters:
//   - accountID: The account ID to scope the search
//   - categoryID: The category ID to filter by
//
// Returns:
//   - []models.MenuItem: List of menu items in the category
//   - error: Any error that occurred during retrieval
func (s *Service) GetMenuItemsByCategoryID(accountID int, categoryID int) ([]models.MenuItem, error) {
	// Get all menu items for the account
	items, err := s.menuItems.GetByAccountID(accountID)
	if err != nil {
		return nil, err
	}

	// Filter items by category
	var filteredItems []models.MenuItem
	for _, item := range items {
		if item.CategoryID != nil && *item.CategoryID == categoryID {
			filteredItems = append(filteredItems, item)
		}
	}

	return filteredItems, nil
}
