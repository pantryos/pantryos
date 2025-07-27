package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// CountsMap represents a map of inventory item IDs to quantities
// This is used for storing inventory snapshots and other quantity mappings
// It implements JSON serialization for database storage
type CountsMap map[int]float64

// Value implements the driver.Valuer interface for JSON serialization
// This allows the CountsMap to be stored as JSON in the database
func (c CountsMap) Value() (driver.Value, error) {
	if c == nil {
		return nil, nil
	}
	bytes, err := json.Marshal(c)
	return string(bytes), err
}

// Scan implements the sql.Scanner interface for JSON deserialization
// This allows the CountsMap to be retrieved from JSON stored in the database
func (c *CountsMap) Scan(value interface{}) error {
	if value == nil {
		*c = nil
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case string:
		bytes = []byte(v)
	case []byte:
		bytes = v
	default:
		return errors.New("cannot scan non-string value into CountsMap")
	}

	return json.Unmarshal(bytes, c)
}

// Organization represents a parent entity that can contain multiple accounts
// This is the top-level entity in the multi-tenant architecture
// Each organization can have multiple business locations (accounts)
type Organization struct {
	ID          int       `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string    `json:"name" gorm:"not null"`
	Description string    `json:"description"`
	Type        string    `json:"type" gorm:"not null;default:'multi_location'"` // single_location, multi_location, enterprise
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	// Note: Foreign key relationships are handled in application logic for ramsql compatibility
}

// User represents a system user with authentication and authorization
// Users belong to specific accounts and have roles that determine their permissions
// The password field is omitted from JSON responses for security
type User struct {
	ID        int       `json:"id" gorm:"primaryKey;autoIncrement"`
	Email     string    `json:"email" gorm:"uniqueIndex;not null"`
	Password  string    `json:"-" gorm:"not null"` // Omit from JSON responses for security
	AccountID int       `json:"account_id" gorm:"not null;index"`
	Role      string    `json:"role" gorm:"not null;default:'user'"` // user, manager, admin, org_admin
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	// Note: Foreign key relationships are handled in application logic for ramsql compatibility
}

// Account represents a business location that can exist independently or within an organization
// This could be a standalone coffee shop or part of a larger chain
// Each account has its own inventory, menu items, and users
type Account struct {
	ID             int       `json:"id" gorm:"primaryKey;autoIncrement"`
	OrganizationID *int      `json:"organization_id" gorm:"index"` // Optional - null for standalone businesses
	Name           string    `json:"name" gorm:"not null"`         // e.g., "Main Street Coffee Shop"
	Location       string    `json:"location"`                     // e.g., "123 Main St, City, State"
	Phone          string    `json:"phone"`
	Email          string    `json:"email"`
	BusinessType   string    `json:"business_type" gorm:"not null;default:'single_location'"` // single_location, multi_location, enterprise
	Status         string    `json:"status" gorm:"not null;default:'active'"`                 // active, inactive, suspended
	CreatedAt      time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	// Note: Foreign key relationships are handled in application logic for ramsql compatibility
}

// InventoryItem represents a physical item that can be tracked in inventory
// Each item belongs to a specific account and has stock level management
// Items can be ingredients, supplies, or any consumable resource
type InventoryItem struct {
	ID              int     `json:"id" gorm:"primaryKey;autoIncrement"`
	AccountID       int     `json:"account_id" gorm:"not null;index"`
	Name            string  `json:"name" gorm:"not null"`
	Unit            string  `json:"unit" gorm:"not null"` // e.g., "kg", "liters", "pieces"
	CostPerUnit     float64 `json:"cost_per_unit" gorm:"not null;default:0"`
	PreferredVendor string  `json:"preferred_vendor" gorm:"default:''"` // Default supplier for this item
	MinStockLevel   float64 `json:"min_stock_level" gorm:"default:0"`   // Alert when stock goes below this
	MaxStockLevel   float64 `json:"max_stock_level" gorm:"default:0"`   // Don't order more than this
	MinWeeksStock   float64 `json:"min_weeks_stock" gorm:"default:2"`   // Minimum weeks of stock to maintain
	MaxWeeksStock   float64 `json:"max_weeks_stock" gorm:"default:8"`   // Maximum weeks of stock to maintain
	// Note: Foreign key relationships are handled in application logic for ramsql compatibility
}

// MenuItem represents a product that can be sold to customers
// Menu items are organized by categories and have pricing information
// They can be linked to inventory items through recipes
type MenuItem struct {
	ID        int     `json:"id" gorm:"primaryKey;autoIncrement"`
	AccountID int     `json:"account_id" gorm:"not null;index"`
	Name      string  `json:"name" gorm:"not null"`
	Price     float64 `json:"price" gorm:"not null;default:0"`
	Category  string  `json:"category"` // e.g., "drinks", "food", "desserts"
	// Note: Foreign key relationships are handled in application logic for ramsql compatibility
}

// RecipeIngredient links menu items to their required inventory items
// This allows for automatic inventory tracking when menu items are sold
// The quantity field specifies how much of the inventory item is needed per menu item
type RecipeIngredient struct {
	ID              int     `json:"id" gorm:"primaryKey;autoIncrement"`
	MenuItemID      int     `json:"menu_item_id" gorm:"not null;index"`
	InventoryItemID int     `json:"inventory_item_id" gorm:"not null;index"`
	Quantity        float64 `json:"quantity" gorm:"not null;default:0"`
	// Note: Foreign key relationships are handled in application logic for ramsql compatibility
}

// InventorySnapshot represents a point-in-time snapshot of inventory levels
// This is used for historical tracking and reporting
// The Counts field stores a map of inventory item IDs to their quantities at the snapshot time
type InventorySnapshot struct {
	ID        int       `json:"id" gorm:"primaryKey;autoIncrement"`
	AccountID int       `json:"account_id" gorm:"not null;index"`
	Timestamp time.Time `json:"timestamp" gorm:"not null;index"`
	Counts    CountsMap `json:"counts" gorm:"type:text"` // map[InventoryItemID]quantity - stored as JSON
	// Note: Foreign key relationships are handled in application logic for ramsql compatibility
}

// Delivery represents a shipment of inventory items from a vendor
// This tracks when items are received and their associated costs
// Used for inventory replenishment and cost tracking
type Delivery struct {
	ID              int       `json:"id" gorm:"primaryKey;autoIncrement"`
	AccountID       int       `json:"account_id" gorm:"not null;index"`
	InventoryItemID int       `json:"inventory_item_id" gorm:"not null;index"`
	Vendor          string    `json:"vendor" gorm:"not null"` // e.g., "Coffee Supply Co.", "Local Dairy"
	Quantity        float64   `json:"quantity" gorm:"not null;default:0"`
	DeliveryDate    time.Time `json:"delivery_date" gorm:"not null;index"`
	Cost            float64   `json:"cost" gorm:"not null;default:0"`
	// Note: Foreign key relationships are handled in application logic for ramsql compatibility
}

// Order represents a purchase order for inventory items
// Orders go through various statuses from pending to delivered
// They can be created by users and approved by managers
type Order struct {
	ID           int       `json:"id" gorm:"primaryKey;autoIncrement"`
	AccountID    int       `json:"account_id" gorm:"not null;index"`
	Status       string    `json:"status" gorm:"not null;default:'pending'"` // pending, approved, ordered, delivered, cancelled
	OrderDate    time.Time `json:"order_date" gorm:"not null"`
	ExpectedDate time.Time `json:"expected_date"`
	TotalCost    float64   `json:"total_cost" gorm:"not null;default:0"`
	Notes        string    `json:"notes"`
	CreatedBy    int       `json:"created_by" gorm:"not null"`
	ApprovedBy   *int      `json:"approved_by"`
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	// Note: Foreign key relationships are handled in application logic for ramsql compatibility
}

// OrderItem represents a specific item within an order
// Each order can contain multiple items from different vendors
// This tracks the quantity, cost, and vendor for each item
type OrderItem struct {
	ID              int     `json:"id" gorm:"primaryKey;autoIncrement"`
	OrderID         int     `json:"order_id" gorm:"not null;index"`
	InventoryItemID int     `json:"inventory_item_id" gorm:"not null;index"`
	Quantity        float64 `json:"quantity" gorm:"not null;default:0"`
	UnitCost        float64 `json:"unit_cost" gorm:"not null;default:0"`
	TotalCost       float64 `json:"total_cost" gorm:"not null;default:0"`
	Vendor          string  `json:"vendor" gorm:"not null"`
	Notes           string  `json:"notes"`
	// Note: Foreign key relationships are handled in application logic for ramsql compatibility
}

// OrderRequest represents a request for inventory items that needs approval
// This allows for a workflow where users can request items that need manager approval
// Requests can have different priorities and deadlines
type OrderRequest struct {
	ID          int       `json:"id" gorm:"primaryKey;autoIncrement"`
	AccountID   int       `json:"account_id" gorm:"not null;index"`
	Status      string    `json:"status" gorm:"not null;default:'pending'"`  // pending, approved, rejected, fulfilled
	Priority    string    `json:"priority" gorm:"not null;default:'normal'"` // low, normal, high, urgent
	RequestDate time.Time `json:"request_date" gorm:"not null"`
	NeededBy    time.Time `json:"needed_by"`
	Notes       string    `json:"notes"`
	CreatedBy   int       `json:"created_by" gorm:"not null"`
	ApprovedBy  *int      `json:"approved_by"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	// Note: Foreign key relationships are handled in application logic for ramsql compatibility
}

// RequestItem represents a specific item within an order request
// Each request can contain multiple items with different priorities
// This tracks the quantity needed and the reason for the request
type RequestItem struct {
	ID              int     `json:"id" gorm:"primaryKey;autoIncrement"`
	OrderRequestID  int     `json:"order_request_id" gorm:"not null;index"`
	InventoryItemID int     `json:"inventory_item_id" gorm:"not null;index"`
	Quantity        float64 `json:"quantity" gorm:"not null;default:0"`
	Reason          string  `json:"reason"` // e.g., "low stock", "new menu item", "special event"
	Priority        string  `json:"priority" gorm:"not null;default:'normal'"`
	// Note: Foreign key relationships are handled in application logic for ramsql compatibility
}

// AccountInvitation represents an invitation for a user to join an account
// This allows account admins to invite users by email without requiring them to know account IDs
// The invitation system provides better security and user experience
type AccountInvitation struct {
	ID         int        `json:"id" gorm:"primaryKey;autoIncrement"`
	AccountID  int        `json:"account_id" gorm:"not null;index"`
	Email      string     `json:"email" gorm:"not null;index"`
	InvitedBy  int        `json:"invited_by" gorm:"not null"`               // User ID who sent the invitation
	Status     string     `json:"status" gorm:"not null;default:'pending'"` // pending, accepted, expired, revoked
	InvitedAt  time.Time  `json:"invited_at" gorm:"autoCreateTime"`
	AcceptedAt *time.Time `json:"accepted_at"`
	ExpiresAt  time.Time  `json:"expires_at" gorm:"not null"` // Invitation expiration date
	CreatedAt  time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	// Note: Foreign key relationships are handled in application logic for ramsql compatibility
}

// Invitation status constants
const (
	AccountInvitationStatusPending  = "pending"
	AccountInvitationStatusAccepted = "accepted"
	AccountInvitationStatusExpired  = "expired"
	AccountInvitationStatusRevoked  = "revoked"
)

// Business type constants
const (
	BusinessTypeSingleLocation = "single_location" // Standalone business (no organization)
	BusinessTypeMultiLocation  = "multi_location"  // Multiple locations under one organization
	BusinessTypeEnterprise     = "enterprise"      // Large enterprise with complex hierarchy
)
