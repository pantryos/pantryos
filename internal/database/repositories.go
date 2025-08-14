package database

import (
	"errors"
	"time"

	"github.com/mnadev/pantryos/internal/models"
	"gorm.io/gorm"
)

// Repository interfaces for better testability
type OrganizationRepository interface {
	Create(organization *models.Organization) error
	GetByID(id int) (*models.Organization, error)
	Update(organization *models.Organization) error
	Delete(id int) error
	GetAccounts(organizationID int) ([]models.Account, error)
}

type AccountRepository interface {
	Create(account *models.Account) error
	GetByID(id int) (*models.Account, error)
	GetByOrganizationID(organizationID int) ([]models.Account, error)
	GetStandalone() ([]models.Account, error) // Get accounts without organization
	GetAll() ([]models.Account, error)
	Update(account *models.Account) error
	Delete(id int) error
}

type UserRepository interface {
	Create(user *models.User) error
	GetByID(id int) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	GetByAccountID(accountID int) ([]models.User, error)
	GetByOrganizationID(organizationID int) ([]models.User, error)
	Update(user *models.User) error
	Delete(id int) error
}

type InventoryItemRepository interface {
	Create(item *models.InventoryItem) error
	GetByID(id int) (*models.InventoryItem, error)
	GetByAccountID(accountID int) ([]models.InventoryItem, error)
	GetByVendor(accountID int, vendor string) ([]models.InventoryItem, error)
	GetLowStockItems(accountID int) ([]models.InventoryItem, error)
	Update(item *models.InventoryItem) error
	Delete(id int) error
}

type MenuItemRepository interface {
	Create(item *models.MenuItem) error
	GetByID(id int) (*models.MenuItem, error)
	GetByAccountID(accountID int) ([]models.MenuItem, error)
	GetByCategory(accountID int, category string) ([]models.MenuItem, error)
	Update(item *models.MenuItem) error
	Delete(id int) error
	GetWithIngredients(id int) (*models.MenuItem, error)
}

type DeliveryRepository interface {
	Create(delivery *models.Delivery) error
	GetByID(id int) (*models.Delivery, error)
	GetByAccountID(accountID int) ([]models.Delivery, error)
	GetByVendor(accountID int, vendor string) ([]models.Delivery, error)
	GetByDateRange(accountID int, startDate, endDate time.Time) ([]models.Delivery, error)
	Update(delivery *models.Delivery) error
	Delete(id int) error
	GetByAccountIDAfterDate(id int, timestamp time.Time) ([]models.Delivery, error)
}

type SaleRepository interface {
	GetByAccountIDAfterDate(accountID int, afterDate time.Time) ([]models.Sale, error)
}

type RecipeRepository interface {
	GetIngredientsByMenuItemID(menuItemID uint) ([]models.RecipeIngredient, error)
}

type InventorySnapshotRepository interface {
	Create(snapshot *models.InventorySnapshot) error
	GetByID(id int) (*models.InventorySnapshot, error)
	GetByAccountID(accountID int) ([]models.InventorySnapshot, error)
	GetLatestByAccountID(accountID int) (*models.InventorySnapshot, error)
	GetByDateRange(accountID int, startDate, endDate time.Time) ([]models.InventorySnapshot, error)
	Update(snapshot *models.InventorySnapshot) error
	Delete(id int) error
}

type OrderRepository interface {
	Create(order *models.Order) error
	GetByID(id int) (*models.Order, error)
	GetByAccountID(accountID int) ([]models.Order, error)
	GetByStatus(accountID int, status string) ([]models.Order, error)
	Update(order *models.Order) error
	Delete(id int) error
	GetWithItems(id int) (*models.Order, error)
}

type OrderRequestRepository interface {
	Create(request *models.OrderRequest) error
	GetByID(id int) (*models.OrderRequest, error)
	GetByAccountID(accountID int) ([]models.OrderRequest, error)
	GetByStatus(accountID int, status string) ([]models.OrderRequest, error)
	Update(request *models.OrderRequest) error
	Delete(id int) error
	GetWithItems(id int) (*models.OrderRequest, error)
}

type AccountInvitationRepository interface {
	Create(invitation *models.AccountInvitation) error
	GetByID(id int) (*models.AccountInvitation, error)
	GetByEmail(email string) (*models.AccountInvitation, error)
	GetByAccountID(accountID int) ([]models.AccountInvitation, error)
	GetPendingByEmail(email string) (*models.AccountInvitation, error)
	Update(invitation *models.AccountInvitation) error
	Delete(id int) error
	DeleteByEmailAndAccount(email string, accountID int) error
}

type CategoryRepository interface {
	Create(category *models.Category) error
	GetByID(id int) (*models.Category, error)
	GetByAccountID(accountID int) ([]models.Category, error)
	GetActiveByAccountID(accountID int) ([]models.Category, error)
	Update(category *models.Category) error
	Delete(id int) error
}

type EmailScheduleRepository interface {
	Create(schedule *models.EmailSchedule) error
	GetByID(id int) (*models.EmailSchedule, error)
	GetByAccountID(accountID int) ([]models.EmailSchedule, error)
	GetByAccountIDAndType(accountID int, emailType string) (*models.EmailSchedule, error)
	GetActiveByAccountID(accountID int) ([]models.EmailSchedule, error)
	Update(schedule *models.EmailSchedule) error
	Delete(id int) error
	UpdateLastSentAt(id int, lastSentAt time.Time) error
}

// Repository implementations
type organizationRepository struct {
	db *DB
}

func NewOrganizationRepository(db *DB) OrganizationRepository {
	return &organizationRepository{db: db}
}

func (r *organizationRepository) Create(organization *models.Organization) error {
	organization.CreatedAt = time.Now()
	organization.UpdatedAt = time.Now()
	return r.db.Create(organization).Error
}

func (r *organizationRepository) GetByID(id int) (*models.Organization, error) {
	var organization models.Organization
	// Use Find instead of First to avoid LIMIT clause that ramsql doesn't support
	err := r.db.Where("id = ?", id).Find(&organization).Error
	if err != nil {
		return nil, err
	}
	if organization.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &organization, nil
}

func (r *organizationRepository) Update(organization *models.Organization) error {
	organization.UpdatedAt = time.Now()
	return r.db.Save(organization).Error
}

func (r *organizationRepository) Delete(id int) error {
	return r.db.Delete(&models.Organization{}, id).Error
}

func (r *organizationRepository) GetAccounts(organizationID int) ([]models.Account, error) {
	var accounts []models.Account
	err := r.db.Where("organization_id = ?", organizationID).Find(&accounts).Error
	return accounts, err
}

// Account repository implementation
type accountRepository struct {
	db *DB
}

func NewAccountRepository(db *DB) AccountRepository {
	return &accountRepository{db: db}
}

func (r *accountRepository) Create(account *models.Account) error {
	account.CreatedAt = time.Now()
	account.UpdatedAt = time.Now()
	return r.db.Create(account).Error
}

func (r *accountRepository) GetByID(id int) (*models.Account, error) {
	var account models.Account
	// Use Find instead of First to avoid LIMIT clause that ramsql doesn't support
	err := r.db.Where("id = ?", id).Find(&account).Error
	if err != nil {
		return nil, err
	}
	if account.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &account, nil
}

func (r *accountRepository) GetByOrganizationID(organizationID int) ([]models.Account, error) {
	var accounts []models.Account
	err := r.db.Where("organization_id = ?", organizationID).Find(&accounts).Error
	return accounts, err
}

func (r *accountRepository) GetStandalone() ([]models.Account, error) {
	var accounts []models.Account
	err := r.db.Where("organization_id IS NULL").Find(&accounts).Error
	return accounts, err
}

func (r *accountRepository) GetAll() ([]models.Account, error) {
	var accounts []models.Account
	err := r.db.Find(&accounts).Error
	return accounts, err
}

func (r *accountRepository) Update(account *models.Account) error {
	account.UpdatedAt = time.Now()
	return r.db.Save(account).Error
}

func (r *accountRepository) Delete(id int) error {
	return r.db.Delete(&models.Account{}, id).Error
}

// User repository implementation
type userRepository struct {
	db *DB
}

func NewUserRepository(db *DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *models.User) error {
	user.CreatedAt = time.Now()
	return r.db.Create(user).Error
}

func (r *userRepository) GetByID(id int) (*models.User, error) {
	var user models.User
	// Use Find instead of First to avoid LIMIT clause that ramsql doesn't support
	err := r.db.Where("id = ?", id).Find(&user).Error
	if err != nil {
		return nil, err
	}
	if user.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &user, nil
}

func (r *userRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	// Use Find instead of First to avoid LIMIT clause that ramsql doesn't support
	err := r.db.Where("email = ?", email).Find(&user).Error
	if err != nil {
		return nil, err
	}
	if user.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &user, nil
}

func (r *userRepository) GetByAccountID(accountID int) ([]models.User, error) {
	var users []models.User
	err := r.db.Where("account_id = ?", accountID).Find(&users).Error
	return users, err
}

func (r *userRepository) GetByOrganizationID(organizationID int) ([]models.User, error) {
	var users []models.User
	err := r.db.Joins("JOIN accounts ON users.account_id = accounts.id").
		Where("accounts.organization_id = ?", organizationID).
		Find(&users).Error
	return users, err
}

func (r *userRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) Delete(id int) error {
	return r.db.Delete(&models.User{}, id).Error
}

// Inventory item repository implementation
type inventoryItemRepository struct {
	db *DB
}

func NewInventoryItemRepository(db *DB) InventoryItemRepository {
	return &inventoryItemRepository{db: db}
}

func (r *inventoryItemRepository) Create(item *models.InventoryItem) error {
	return r.db.Create(item).Error
}

func (r *inventoryItemRepository) GetByID(id int) (*models.InventoryItem, error) {
	var item models.InventoryItem
	// Use Find instead of First to avoid LIMIT clause that ramsql doesn't support
	err := r.db.Where("id = ?", id).Find(&item).Error
	if err != nil {
		return nil, err
	}
	if item.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &item, nil
}

func (r *inventoryItemRepository) GetByAccountID(accountID int) ([]models.InventoryItem, error) {
	var items []models.InventoryItem
	err := r.db.Where("account_id = ?", accountID).Find(&items).Error
	return items, err
}

func (r *inventoryItemRepository) GetByVendor(accountID int, vendor string) ([]models.InventoryItem, error) {
	var items []models.InventoryItem
	err := r.db.Where("account_id = ? AND preferred_vendor = ?", accountID, vendor).Find(&items).Error
	return items, err
}

func (r *inventoryItemRepository) GetLowStockItems(accountID int) ([]models.InventoryItem, error) {
	var items []models.InventoryItem
	err := r.db.Where("account_id = ? AND min_stock_level > 0", accountID).Find(&items).Error
	return items, err
}

func (r *inventoryItemRepository) Update(item *models.InventoryItem) error {
	return r.db.Save(item).Error
}

func (r *inventoryItemRepository) Delete(id int) error {
	return r.db.Delete(&models.InventoryItem{}, id).Error
}

// Menu item repository implementation
type menuItemRepository struct {
	db *DB
}

func NewMenuItemRepository(db *DB) MenuItemRepository {
	return &menuItemRepository{db: db}
}

func (r *menuItemRepository) Create(item *models.MenuItem) error {
	return r.db.Create(item).Error
}

func (r *menuItemRepository) GetByID(id int) (*models.MenuItem, error) {
	var item models.MenuItem
	// Use Find instead of First to avoid LIMIT clause that ramsql doesn't support
	err := r.db.Where("id = ?", id).Find(&item).Error
	if err != nil {
		return nil, err
	}
	if item.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &item, nil
}

func (r *menuItemRepository) GetByAccountID(accountID int) ([]models.MenuItem, error) {
	var items []models.MenuItem
	err := r.db.Where("account_id = ?", accountID).Find(&items).Error
	return items, err
}

func (r *menuItemRepository) GetByCategory(accountID int, category string) ([]models.MenuItem, error) {
	var items []models.MenuItem
	err := r.db.Where("account_id = ? AND category = ?", accountID, category).Find(&items).Error
	return items, err
}

func (r *menuItemRepository) Update(item *models.MenuItem) error {
	return r.db.Save(item).Error
}

func (r *menuItemRepository) Delete(id int) error {
	return r.db.Delete(&models.MenuItem{}, id).Error
}

func (r *menuItemRepository) GetWithIngredients(id int) (*models.MenuItem, error) {
	var item models.MenuItem
	err := r.db.Preload("Ingredients").First(&item, id).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

// Delivery repository implementation
type deliveryRepository struct {
	db *DB
}

func NewDeliveryRepository(db *DB) DeliveryRepository {
	return &deliveryRepository{db: db}
}

func (r *deliveryRepository) Create(delivery *models.Delivery) error {
	return r.db.Create(delivery).Error
}

func (r *deliveryRepository) GetByID(id int) (*models.Delivery, error) {
	var delivery models.Delivery
	// Use Find instead of First to avoid LIMIT clause that ramsql doesn't support
	err := r.db.Where("id = ?", id).Find(&delivery).Error
	if err != nil {
		return nil, err
	}
	if delivery.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &delivery, nil
}

func (r *deliveryRepository) GetByAccountID(accountID int) ([]models.Delivery, error) {
	var deliveries []models.Delivery
	err := r.db.Where("account_id = ?", accountID).Order("delivery_date DESC").Find(&deliveries).Error
	return deliveries, err
}

func (r *deliveryRepository) GetByVendor(accountID int, vendor string) ([]models.Delivery, error) {
	var deliveries []models.Delivery
	err := r.db.Where("account_id = ? AND vendor = ?", accountID, vendor).Order("delivery_date DESC").Find(&deliveries).Error
	return deliveries, err
}

func (r *deliveryRepository) GetByDateRange(accountID int, startDate, endDate time.Time) ([]models.Delivery, error) {
	var deliveries []models.Delivery
	err := r.db.Where("account_id = ? AND delivery_date BETWEEN ? AND ?",
		accountID, startDate, endDate).Order("delivery_date DESC").Find(&deliveries).Error
	return deliveries, err
}

func (r *deliveryRepository) Update(delivery *models.Delivery) error {
	return r.db.Save(delivery).Error
}

func (r *deliveryRepository) Delete(id int) error {
	return r.db.Delete(&models.Delivery{}, id).Error
}

// GetByAccountIDAfterDate retrieves all deliveries for a specific account
// that occurred after the given timestamp. This is crucial for the hybrid
// stock calculation model.
func (r *deliveryRepository) GetByAccountIDAfterDate(accountID int, afterDate time.Time) ([]models.Delivery, error) {
	var deliveries []models.Delivery

	// Query finds deliveries for the account where the delivery_date is more recent than the provided timestamp.
	err := r.db.Where("account_id = ? AND delivery_date > ?", accountID, afterDate).
		Order("delivery_date asc"). // Ordering chronologically is good practice for processing transactions.
		Find(&deliveries).Error

	if err != nil {
		return nil, err
	}

	return deliveries, nil
}

// Inventory snapshot repository implementation
type inventorySnapshotRepository struct {
	db *DB
}

func NewInventorySnapshotRepository(db *DB) InventorySnapshotRepository {
	return &inventorySnapshotRepository{db: db}
}

func (r *inventorySnapshotRepository) Create(snapshot *models.InventorySnapshot) error {
	snapshot.Timestamp = time.Now()
	return r.db.Create(snapshot).Error
}

func (r *inventorySnapshotRepository) GetByID(id int) (*models.InventorySnapshot, error) {
	var snapshot models.InventorySnapshot
	// Use Find instead of First to avoid LIMIT clause that ramsql doesn't support
	err := r.db.Where("id = ?", id).Find(&snapshot).Error
	if err != nil {
		return nil, err
	}
	if snapshot.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &snapshot, nil
}

func (r *inventorySnapshotRepository) GetByAccountID(accountID int) ([]models.InventorySnapshot, error) {
	var snapshots []models.InventorySnapshot
	err := r.db.Where("account_id = ?", accountID).Order("timestamp DESC").Find(&snapshots).Error
	return snapshots, err
}

func (r *inventorySnapshotRepository) GetLatestByAccountID(accountID int) (*models.InventorySnapshot, error) {
	var snapshots []models.InventorySnapshot
	// Use Find instead of First to avoid LIMIT clause that ramsql doesn't support
	err := r.db.Where("account_id = ?", accountID).Order("timestamp DESC").Find(&snapshots).Error
	if err != nil {
		return nil, err
	}
	if len(snapshots) == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &snapshots[0], nil
}

func (r *inventorySnapshotRepository) GetByDateRange(accountID int, startDate, endDate time.Time) ([]models.InventorySnapshot, error) {
	var snapshots []models.InventorySnapshot
	err := r.db.Where("account_id = ? AND timestamp BETWEEN ? AND ?",
		accountID, startDate, endDate).Order("timestamp DESC").Find(&snapshots).Error
	return snapshots, err
}

func (r *inventorySnapshotRepository) Update(snapshot *models.InventorySnapshot) error {
	return r.db.Save(snapshot).Error
}

func (r *inventorySnapshotRepository) Delete(id int) error {
	return r.db.Delete(&models.InventorySnapshot{}, id).Error
}

// sale
type saleRepository struct {
	db *DB
}

func NewSaleRepository(db *DB) SaleRepository {
	return &saleRepository{db: db}
}

func (r *saleRepository) GetByAccountIDAfterDate(accountID int, afterDate time.Time) ([]models.Sale, error) {
	var sales []models.Sale

	err := r.db.Preload("Items").
		Where("account_id = ? AND sale_date > ?", accountID, afterDate).
		Order("sale_date asc").
		Find(&sales).Error

	if err != nil {
		return nil, err
	}

	return sales, nil
}

// Order repository implementation
type orderRepository struct {
	db *DB
}

type recipeRepository struct {
	db *DB
}

func NewRecipeRepository(db *DB) RecipeRepository {
	return &recipeRepository{db: db}
}

func (r *recipeRepository) GetIngredientsByMenuItemID(menuItemID uint) ([]models.RecipeIngredient, error) {
	var ingredients []models.RecipeIngredient

	err := r.db.Where("menu_item_id = ?", menuItemID).Find(&ingredients).Error
	if err != nil {
		return nil, err
	}

	return ingredients, nil
}

func NewOrderRepository(db *DB) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) Create(order *models.Order) error {
	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()
	return r.db.Create(order).Error
}

func (r *orderRepository) GetByID(id int) (*models.Order, error) {
	var order models.Order
	// Use Find instead of First to avoid LIMIT clause that ramsql doesn't support
	err := r.db.Where("id = ?", id).Find(&order).Error
	if err != nil {
		return nil, err
	}
	if order.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &order, nil
}

func (r *orderRepository) GetByAccountID(accountID int) ([]models.Order, error) {
	var orders []models.Order
	err := r.db.Where("account_id = ?", accountID).Order("created_at DESC").Find(&orders).Error
	return orders, err
}

func (r *orderRepository) GetByStatus(accountID int, status string) ([]models.Order, error) {
	var orders []models.Order
	err := r.db.Where("account_id = ? AND status = ?", accountID, status).Order("created_at DESC").Find(&orders).Error
	return orders, err
}

func (r *orderRepository) Update(order *models.Order) error {
	order.UpdatedAt = time.Now()
	return r.db.Save(order).Error
}

func (r *orderRepository) Delete(id int) error {
	return r.db.Delete(&models.Order{}, id).Error
}

func (r *orderRepository) GetWithItems(id int) (*models.Order, error) {
	var order models.Order
	err := r.db.Preload("OrderItems").First(&order, id).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// Order request repository implementation
type orderRequestRepository struct {
	db *DB
}

func NewOrderRequestRepository(db *DB) OrderRequestRepository {
	return &orderRequestRepository{db: db}
}

func (r *orderRequestRepository) Create(request *models.OrderRequest) error {
	request.CreatedAt = time.Now()
	request.UpdatedAt = time.Now()
	return r.db.Create(request).Error
}

func (r *orderRequestRepository) GetByID(id int) (*models.OrderRequest, error) {
	var request models.OrderRequest
	// Use Find instead of First to avoid LIMIT clause that ramsql doesn't support
	err := r.db.Where("id = ?", id).Find(&request).Error
	if err != nil {
		return nil, err
	}
	if request.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &request, nil
}

func (r *orderRequestRepository) GetByAccountID(accountID int) ([]models.OrderRequest, error) {
	var requests []models.OrderRequest
	err := r.db.Where("account_id = ?", accountID).Order("created_at DESC").Find(&requests).Error
	return requests, err
}

func (r *orderRequestRepository) GetByStatus(accountID int, status string) ([]models.OrderRequest, error) {
	var requests []models.OrderRequest
	err := r.db.Where("account_id = ? AND status = ?", accountID, status).Order("created_at DESC").Find(&requests).Error
	return requests, err
}

func (r *orderRequestRepository) Update(request *models.OrderRequest) error {
	request.UpdatedAt = time.Now()
	return r.db.Save(request).Error
}

func (r *orderRequestRepository) Delete(id int) error {
	return r.db.Delete(&models.OrderRequest{}, id).Error
}

func (r *orderRequestRepository) GetWithItems(id int) (*models.OrderRequest, error) {
	var request models.OrderRequest
	err := r.db.Preload("RequestItems").First(&request, id).Error
	if err != nil {
		return nil, err
	}
	return &request, nil
}

// Account invitation repository implementation
type accountInvitationRepository struct {
	db *DB
}

func NewAccountInvitationRepository(db *DB) AccountInvitationRepository {
	return &accountInvitationRepository{db: db}
}

func (r *accountInvitationRepository) Create(invitation *models.AccountInvitation) error {
	invitation.CreatedAt = time.Now()
	invitation.UpdatedAt = time.Now()
	return r.db.Create(invitation).Error
}

func (r *accountInvitationRepository) GetByID(id int) (*models.AccountInvitation, error) {
	var invitation models.AccountInvitation
	// Use Find instead of First to avoid LIMIT clause that ramsql doesn't support
	err := r.db.Where("id = ?", id).Find(&invitation).Error
	if err != nil {
		return nil, err
	}
	if invitation.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &invitation, nil
}

func (r *accountInvitationRepository) GetByEmail(email string) (*models.AccountInvitation, error) {
	var invitation models.AccountInvitation
	// Use Find instead of First to avoid LIMIT clause that ramsql doesn't support
	err := r.db.Where("email = ?", email).Find(&invitation).Error
	if err != nil {
		return nil, err
	}
	if invitation.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &invitation, nil
}

func (r *accountInvitationRepository) GetByAccountID(accountID int) ([]models.AccountInvitation, error) {
	var invitations []models.AccountInvitation
	err := r.db.Where("account_id = ?", accountID).Find(&invitations).Error
	return invitations, err
}

func (r *accountInvitationRepository) GetPendingByEmail(email string) (*models.AccountInvitation, error) {
	var invitation models.AccountInvitation
	// Use Find instead of First to avoid LIMIT clause that ramsql doesn't support
	err := r.db.Where("email = ? AND status = ?", email, "pending").Find(&invitation).Error
	if err != nil {
		return nil, err
	}
	if invitation.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &invitation, nil
}

func (r *accountInvitationRepository) Update(invitation *models.AccountInvitation) error {
	invitation.UpdatedAt = time.Now()
	return r.db.Save(invitation).Error
}

func (r *accountInvitationRepository) Delete(id int) error {
	return r.db.Delete(&models.AccountInvitation{}, id).Error
}

func (r *accountInvitationRepository) DeleteByEmailAndAccount(email string, accountID int) error {
	return r.db.Where("email = ? AND account_id = ?", email, accountID).Delete(&models.AccountInvitation{}).Error
}

// Category repository implementation
type categoryRepository struct {
	db *DB
}

func NewCategoryRepository(db *DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) Create(category *models.Category) error {
	category.CreatedAt = time.Now()
	category.UpdatedAt = time.Now()
	return r.db.Create(category).Error
}

func (r *categoryRepository) GetByID(id int) (*models.Category, error) {
	var category models.Category
	// Use Find instead of First to avoid LIMIT clause that ramsql doesn't support
	err := r.db.Where("id = ?", id).Find(&category).Error
	if err != nil {
		return nil, err
	}
	if category.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &category, nil
}

func (r *categoryRepository) GetByAccountID(accountID int) ([]models.Category, error) {
	var categories []models.Category
	err := r.db.Where("account_id = ?", accountID).Find(&categories).Error
	return categories, err
}

func (r *categoryRepository) GetActiveByAccountID(accountID int) ([]models.Category, error) {
	var categories []models.Category
	err := r.db.Where("account_id = ? AND is_active = true", accountID).Find(&categories).Error
	return categories, err
}

func (r *categoryRepository) Update(category *models.Category) error {
	category.UpdatedAt = time.Now()
	return r.db.Save(category).Error
}

func (r *categoryRepository) Delete(id int) error {
	return r.db.Delete(&models.Category{}, id).Error
}

// Email schedule repository implementation
type emailScheduleRepository struct {
	db *DB
}

func NewEmailScheduleRepository(db *DB) EmailScheduleRepository {
	return &emailScheduleRepository{db: db}
}

func (r *emailScheduleRepository) Create(schedule *models.EmailSchedule) error {
	schedule.CreatedAt = time.Now()
	schedule.UpdatedAt = time.Now()
	return r.db.Create(schedule).Error
}

func (r *emailScheduleRepository) GetByID(id int) (*models.EmailSchedule, error) {
	var schedule models.EmailSchedule
	// Use Find instead of First to avoid LIMIT clause that ramsql doesn't support
	err := r.db.Where("id = ?", id).Find(&schedule).Error
	if err != nil {
		return nil, err
	}
	if schedule.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &schedule, nil
}

func (r *emailScheduleRepository) GetByAccountID(accountID int) ([]models.EmailSchedule, error) {
	var schedules []models.EmailSchedule
	err := r.db.Where("account_id = ?", accountID).Order("created_at DESC").Find(&schedules).Error
	return schedules, err
}

func (r *emailScheduleRepository) GetByAccountIDAndType(accountID int, emailType string) (*models.EmailSchedule, error) {
	var schedule models.EmailSchedule
	err := r.db.Where("account_id = ? AND email_type = ?", accountID, emailType).First(&schedule).Error
	if err != nil {
		return nil, err
	}
	if schedule.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &schedule, nil
}

func (r *emailScheduleRepository) GetActiveByAccountID(accountID int) ([]models.EmailSchedule, error) {
	var schedules []models.EmailSchedule
	err := r.db.Where("account_id = ? AND is_active = true", accountID).Order("created_at DESC").Find(&schedules).Error
	return schedules, err
}

func (r *emailScheduleRepository) Update(schedule *models.EmailSchedule) error {
	schedule.UpdatedAt = time.Now()
	return r.db.Save(schedule).Error
}

func (r *emailScheduleRepository) Delete(id int) error {
	return r.db.Delete(&models.EmailSchedule{}, id).Error
}

func (r *emailScheduleRepository) UpdateLastSentAt(id int, lastSentAt time.Time) error {
	return r.db.Model(&models.EmailSchedule{}).Where("id = ?", id).Update("last_sent_at", lastSentAt).Error
}

// Business logic functions
func (db *DB) GetInventoryVariance(accountID int, startDate, endDate time.Time) (map[int]float64, error) {
	// This will be implemented when we add Toast POS integration
	// For now, return empty map
	return make(map[int]float64), nil
}

// Helper function to check if a record belongs to an account
func (db *DB) validateAccountAccess(accountID int, userAccountID int) error {
	if accountID != userAccountID {
		return errors.New("access denied: account mismatch")
	}
	return nil
}
