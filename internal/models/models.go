package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"gorm.io/gorm"
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
	gorm.Model
	Name        string `json:"name" gorm:"not null"`
	Description string `json:"description"`
	Type        string `json:"type" gorm:"not null;default:'multi_location'"` // single_location, multi_location, enterprise

	// Relationships
	Accounts []Account `json:"accounts" gorm:"foreignKey:OrganizationID"`
}

// User represents a system user with authentication and authorization
// Users belong to specific accounts and have roles that determine their permissions
// The password field is omitted from JSON responses for security
type User struct {
	gorm.Model
	Email      string `json:"email" gorm:"uniqueIndex;not null"`
	Password   string `json:"-" gorm:"not null"` // Omit from JSON responses for security
	AccountID  uint   `json:"account_id" gorm:"not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	IsVerified bool   `json:"isVerified" gorm:"not null;default:false"`
	Role       string `json:"role" gorm:"not null;default:'user'"` // user, manager, admin, org_admin

	// Relationships
	Account Account `json:"account" gorm:"foreignKey:AccountID"`
}

// Account represents a business location that can exist independently or within an organization
// This could be a standalone coffee shop or part of a larger chain
// Each account has its own inventory, menu items, and users
type Account struct {
	gorm.Model
	OrganizationID *uint  `json:"organization_id" gorm:"index;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"` // Optional - null for standalone businesses
	Name           string `json:"name" gorm:"not null"`                                                       // e.g., "Main Street Coffee Shop"
	Location       string `json:"location"`                                                                   // e.g., "123 Main St, City, State"
	Phone          string `json:"phone"`
	Email          string `json:"email"`
	BusinessType   string `json:"business_type" gorm:"not null;default:'single_location'"` // single_location, multi_location, enterprise
	Status         string `json:"status" gorm:"not null;default:'active'"`                 // active, inactive, suspended

	// Relationships
	Organization   *Organization       `json:"organization" gorm:"foreignKey:OrganizationID"`
	Users          []User              `json:"users" gorm:"foreignKey:AccountID"`
	Categories     []Category          `json:"categories" gorm:"foreignKey:AccountID"`
	Units          []Unit              `json:"units" gorm:"foreignKey:AccountID"`
	Vendors        []Vendor            `json:"vendors" gorm:"foreignKey:AccountID"`
	InventoryItems []InventoryItem     `json:"inventory_items" gorm:"foreignKey:AccountID"`
	MenuItems      []MenuItem          `json:"menu_items" gorm:"foreignKey:AccountID"`
	Sales          []Sale              `json:"sales" gorm:"foreignKey:AccountID"`
	Orders         []Order             `json:"orders" gorm:"foreignKey:AccountID"`
	OrderRequests  []OrderRequest      `json:"order_requests" gorm:"foreignKey:AccountID"`
	Deliveries     []Delivery          `json:"deliveries" gorm:"foreignKey:AccountID"`
	Snapshots      []InventorySnapshot `json:"snapshots" gorm:"foreignKey:AccountID"`
	Invitations    []AccountInvitation `json:"invitations" gorm:"foreignKey:AccountID"`
	EmailLogs      []EmailLog          `json:"email_logs" gorm:"foreignKey:AccountID"`
	EmailSchedules []EmailSchedule     `json:"email_schedules" gorm:"foreignKey:AccountID"`
}

// Category represents a classification for inventory items and menu items
// Categories help organize items for better management and reporting
// Each category belongs to a specific account for proper scoping
type Category struct {
	gorm.Model
	AccountID   uint   `json:"account_id" gorm:"not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Name        string `json:"name" gorm:"not null"`
	Description string `json:"description"`
	Color       string `json:"color" gorm:"default:'#6B7280'"` // Hex color for UI display
	IsActive    bool   `json:"is_active" gorm:"not null;default:true"`

	// Relationships
	Account        Account         `json:"account" gorm:"foreignKey:AccountID"`
	InventoryItems []InventoryItem `json:"inventory_items" gorm:"foreignKey:CategoryID"`
	MenuItems      []MenuItem      `json:"menu_items" gorm:"foreignKey:CategoryID"`
}

// Unit represents a unit of measurement that can be used for inventory items
// This allows for consistent unit handling across the system
type Unit struct {
	gorm.Model
	AccountID  uint   `json:"account_id" gorm:"not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Name       string `json:"name" gorm:"not null;uniqueIndex"`           // e.g., "gallon", "pound", "piece"
	Symbol     string `json:"symbol" gorm:"not null"`                     // e.g., "gal", "lb", "pc"
	Type       string `json:"type" gorm:"not null"`                       // "volume", "weight", "count", "length"
	IsBaseUnit bool   `json:"is_base_unit" gorm:"not null;default:false"` // Whether this is a base unit for conversions

	// Relationships
	Account           Account            `json:"account" gorm:"foreignKey:AccountID"`
	FromConversions   []UnitConversion   `json:"from_conversions" gorm:"foreignKey:FromUnitID"`
	ToConversions     []UnitConversion   `json:"to_conversions" gorm:"foreignKey:ToUnitID"`
	InventoryItems    []InventoryItem    `json:"inventory_items" gorm:"foreignKey:BaseUnitID"`
	InventoryVariants []InventoryVariant `json:"inventory_variants" gorm:"foreignKey:BaseUnitID"`
	PackageVariants   []InventoryVariant `json:"package_variants" gorm:"foreignKey:PackageUnitID"`
}

// UnitConversion represents conversion factors between different units
// This allows for automatic conversion between units (e.g., 1 gallon = 3.785 liters)
type UnitConversion struct {
	gorm.Model
	AccountID      uint    `json:"account_id" gorm:"not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	FromUnitID     uint    `json:"from_unit_id" gorm:"not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	ToUnitID       uint    `json:"to_unit_id" gorm:"not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	ConversionRate float64 `json:"conversion_rate" gorm:"not null"` // How many 'to' units equal 1 'from' unit
	IsActive       bool    `json:"is_active" gorm:"not null;default:true"`

	// Relationships
	Account  Account `json:"account" gorm:"foreignKey:AccountID"`
	FromUnit Unit    `json:"from_unit" gorm:"foreignKey:FromUnitID"`
	ToUnit   Unit    `json:"to_unit" gorm:"foreignKey:ToUnitID"`
}

// InventoryItem represents a base inventory item that can have multiple variants
// This is the core item definition (e.g., "Milk", "Coffee Beans", "Sugar")
type InventoryItem struct {
	gorm.Model
	AccountID     uint    `json:"account_id" gorm:"not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Name          string  `json:"name" gorm:"not null"`                                                             // e.g., "Milk", "Coffee Beans"
	Description   string  `json:"description"`                                                                      // Detailed description
	BaseUnitID    uint    `json:"base_unit_id" gorm:"not null;index;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"` // Primary unit for this item (e.g., gallon for milk)
	CategoryID    *uint   `json:"category_id" gorm:"index;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`           // Optional category assignment
	MinStockLevel float64 `json:"min_stock_level" gorm:"default:0"`                                                 // Alert when stock goes below this (in base units)
	MaxStockLevel float64 `json:"max_stock_level" gorm:"default:0"`                                                 // Don't order more than this (in base units)
	MinWeeksStock float64 `json:"min_weeks_stock" gorm:"default:2"`                                                 // Minimum weeks of stock to maintain
	MaxWeeksStock float64 `json:"max_weeks_stock" gorm:"default:8"`                                                 // Maximum weeks of stock to maintain
	WastageRate   float64 `json:"wastage_rate" gorm:"default:0"`                                                    // Wastage rate as a percentage
	IsActive      bool    `json:"is_active" gorm:"not null;default:true"`

	// Relationships
	Account           Account            `json:"account" gorm:"foreignKey:AccountID"`
	BaseUnit          Unit               `json:"base_unit" gorm:"foreignKey:BaseUnitID"`
	Category          *Category          `json:"category" gorm:"foreignKey:CategoryID"`
	Variants          []InventoryVariant `json:"variants" gorm:"foreignKey:InventoryItemID"`
	RecipeIngredients []RecipeIngredient `json:"recipe_ingredients" gorm:"foreignKey:InventoryVariantID"`
}

// Vendor represents a supplier or vendor that provides inventory items
// This allows for better vendor management and relationship tracking
type Vendor struct {
	gorm.Model
	AccountID    uint   `json:"account_id" gorm:"not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Name         string `json:"name" gorm:"not null"`                  // e.g., "Coffee Supply Co.", "Local Dairy"
	ContactName  string `json:"contact_name"`                          // Primary contact person
	Email        string `json:"email"`                                 // Contact email
	Phone        string `json:"phone"`                                 // Contact phone
	Website      string `json:"website"`                               // Vendor website
	Address      string `json:"address"`                               // Physical address
	PaymentTerms string `json:"payment_terms" gorm:"default:'Net 30'"` // Payment terms (e.g., "Net 30", "COD")
	IsActive     bool   `json:"is_active" gorm:"not null;default:true"`

	// Relationships
	Account           Account            `json:"account" gorm:"foreignKey:AccountID"`
	PreferredVariants []InventoryVariant `json:"preferred_variants" gorm:"foreignKey:PreferredVendorID"`
	Deliveries        []Delivery         `json:"deliveries" gorm:"foreignKey:VendorID"`
	OrderItems        []OrderItem        `json:"order_items" gorm:"foreignKey:VendorID"`
}

// InventoryVariant represents a specific variant of an inventory item
// This handles different brands, sizes, or packaging of the same base item
type InventoryVariant struct {
	gorm.Model
	AccountID         uint    `json:"account_id" gorm:"not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	InventoryItemID   uint    `json:"inventory_item_id" gorm:"not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Brand             string  `json:"brand" gorm:"not null"`                                                               // e.g., "Organic Valley", "Starbucks"
	VariantName       string  `json:"variant_name"`                                                                        // e.g., "2% Milk", "Dark Roast"
	SKU               string  `json:"sku" gorm:"uniqueIndex"`                                                              // Stock Keeping Unit
	PackageUnitID     uint    `json:"package_unit_id" gorm:"not null;index;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"` // Unit of the package (e.g., "box", "bottle")
	PackageSize       float64 `json:"package_size" gorm:"not null"`                                                        // Size of the package in package units
	BaseUnitID        uint    `json:"base_unit_id" gorm:"not null;index;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`    // Unit of individual items (e.g., "gallon", "pound")
	ItemsPerPackage   float64 `json:"items_per_package" gorm:"not null;default:1"`                                         // How many base units per package
	CostPerPackage    float64 `json:"cost_per_package" gorm:"not null;default:0"`                                          // Cost for the entire package
	CostPerBaseUnit   float64 `json:"cost_per_base_unit" gorm:"not null;default:0"`                                        // Calculated cost per base unit
	PreferredVendorID *uint   `json:"preferred_vendor_id" gorm:"index;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`      // Default supplier for this variant
	Link              string  `json:"link"`                                                                                // URL to product page, supplier catalog, or other reference
	IsActive          bool    `json:"is_active" gorm:"not null;default:true"`

	// Relationships
	Account           Account            `json:"account" gorm:"foreignKey:AccountID"`
	InventoryItem     InventoryItem      `json:"inventory_item" gorm:"foreignKey:InventoryItemID"`
	PackageUnit       Unit               `json:"package_unit" gorm:"foreignKey:PackageUnitID"`
	BaseUnit          Unit               `json:"base_unit" gorm:"foreignKey:BaseUnitID"`
	PreferredVendor   *Vendor            `json:"preferred_vendor" gorm:"foreignKey:PreferredVendorID"`
	Lots              []InventoryLot     `json:"lots" gorm:"foreignKey:InventoryVariantID"`
	Stock             InventoryStock     `json:"stock" gorm:"foreignKey:InventoryVariantID"`
	Deliveries        []Delivery         `json:"deliveries" gorm:"foreignKey:InventoryVariantID"`
	OrderItems        []OrderItem        `json:"order_items" gorm:"foreignKey:InventoryVariantID"`
	RequestItems      []RequestItem      `json:"request_items" gorm:"foreignKey:InventoryVariantID"`
	RecipeIngredients []RecipeIngredient `json:"recipe_ingredients" gorm:"foreignKey:InventoryVariantID"`
}

// InventoryLot represents a specific batch of inventory received from a delivery
// This enables FIFO/FEFO inventory management and accurate cost tracking
type InventoryLot struct {
	gorm.Model
	AccountID          uint       `json:"account_id" gorm:"not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	InventoryVariantID uint       `json:"inventory_variant_id" gorm:"not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	DeliveryID         uint       `json:"delivery_id" gorm:"index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"` // Reference to the delivery that created this lot
	InitialQuantity    float64    `json:"initial_quantity" gorm:"not null"`                                      // Initial quantity in base units
	CurrentQuantity    float64    `json:"current_quantity" gorm:"not null"`                                      // Current quantity in base units
	CostPerBaseUnit    float64    `json:"cost_per_base_unit" gorm:"not null"`                                    // Cost for this specific batch
	ReceivedAt         time.Time  `json:"received_at" gorm:"not null"`                                           // When this lot was received
	ExpirationDate     *time.Time `json:"expiration_date"`                                                       // Expiration date for perishable goods
	IsActive           bool       `json:"is_active" gorm:"not null;default:true"`                                // Whether this lot is still active

	// Relationships
	Account          Account          `json:"account" gorm:"foreignKey:AccountID"`
	InventoryVariant InventoryVariant `json:"inventory_variant" gorm:"foreignKey:InventoryVariantID"`
	Delivery         Delivery         `json:"delivery" gorm:"foreignKey:DeliveryID"`
}

// InventoryStock represents the current stock level for a specific variant
// This is now a calculated view of all active lots for a variant
type InventoryStock struct {
	gorm.Model
	AccountID           uint      `json:"account_id" gorm:"not null;index"`
	InventoryVariantID  uint      `json:"inventory_variant_id" gorm:"not null;uniqueIndex"`
	QuantityInBaseUnits float64   `json:"quantity_in_base_units" gorm:"not null;default:0"` // Calculated sum of all active lots
	LastUpdated         time.Time `json:"last_updated" gorm:"autoUpdateTime"`
	// Note: Foreign key relationships are handled in application logic for ramsql compatibility
}

// CalculateBaseUnits calculates the base units from package units
// This is a helper method to ensure consistency in calculations
func (iv *InventoryVariant) CalculateBaseUnits(packageQuantity float64) float64 {
	return packageQuantity * iv.ItemsPerPackage
}

// CalculatePackageUnits calculates the package units from base units
// This is a helper method to ensure consistency in calculations
func (iv *InventoryVariant) CalculatePackageUnits(baseQuantity float64) float64 {
	if iv.ItemsPerPackage == 0 {
		return 0
	}
	return baseQuantity / iv.ItemsPerPackage
}

// UpdateStock updates the stock levels for this variant
// This method sets the base unit quantity directly
func (is *InventoryStock) UpdateStock(baseQuantity float64) {
	is.QuantityInBaseUnits = baseQuantity
}

// UpdateStockFromPackages updates the stock levels from package quantity
// This method converts package quantity to base units and stores it
func (is *InventoryStock) UpdateStockFromPackages(packageQuantity float64, variant *InventoryVariant) {
	is.QuantityInBaseUnits = variant.CalculateBaseUnits(packageQuantity)
}

// GetDisplayName returns a human-readable name for the variant
// This combines brand, variant name, and package information
func (iv *InventoryVariant) GetDisplayName() string {
	if iv.VariantName != "" {
		return iv.Brand + " " + iv.VariantName
	}
	return iv.Brand
}

// GetFullContactInfo returns a formatted string with vendor contact information
// This is useful for displaying vendor details in reports or UI
func (v *Vendor) GetFullContactInfo() string {
	info := v.Name
	if v.ContactName != "" {
		info += " (" + v.ContactName + ")"
	}
	if v.Phone != "" {
		info += " - " + v.Phone
	}
	if v.Email != "" {
		info += " - " + v.Email
	}
	return info
}

// IsPreferredVendorFor checks if this vendor is the preferred vendor for a given variant
// This is a helper method for vendor relationship management
func (v *Vendor) IsPreferredVendorFor(variant *InventoryVariant) bool {
	return variant.PreferredVendorID != nil && uint(*variant.PreferredVendorID) == v.ID
}

// GetQuantityInPackages calculates the package quantity from base units
// This is a computed property that can be used when package quantities are needed
func (is *InventoryStock) GetQuantityInPackages(variant *InventoryVariant) float64 {
	return variant.CalculatePackageUnits(is.QuantityInBaseUnits)
}

// GetQuantityInPackages calculates the package quantity from base units
// This is a computed property that can be used when package quantities are needed
func (d *Delivery) GetQuantityInPackages(variant *InventoryVariant) float64 {
	return variant.CalculatePackageUnits(d.QuantityInBaseUnits)
}

// GetQuantityInPackages calculates the package quantity from base units
// This is a computed property that can be used when package quantities are needed
func (oi *OrderItem) GetQuantityInPackages(variant *InventoryVariant) float64 {
	return variant.CalculatePackageUnits(oi.QuantityInBaseUnits)
}

// GetQuantityInPackages calculates the package quantity from base units
// This is a computed property that can be used when package quantities are needed
func (ri *RequestItem) GetQuantityInPackages(variant *InventoryVariant) float64 {
	return variant.CalculatePackageUnits(ri.QuantityInBaseUnits)
}

// GetQuantityInPackages calculates the package quantity from base units
// This is a computed property that can be used when package quantities are needed
func (il *InventoryLot) GetQuantityInPackages(variant *InventoryVariant) float64 {
	return variant.CalculatePackageUnits(il.CurrentQuantity)
}

// IsExpired checks if this lot has expired
// This is useful for FEFO (First-Expired, First-Out) inventory management
func (il *InventoryLot) IsExpired() bool {
	if il.ExpirationDate == nil {
		return false
	}
	return time.Now().After(*il.ExpirationDate)
}

// DaysUntilExpiration returns the number of days until expiration
// Returns negative values for expired lots
func (il *InventoryLot) DaysUntilExpiration() int {
	if il.ExpirationDate == nil {
		return 999999 // No expiration date, treat as very far in the future
	}
	return int(time.Until(*il.ExpirationDate).Hours() / 24)
}

// ConsumeQuantity reduces the current quantity by the specified amount
// This method ensures the quantity doesn't go below zero
func (il *InventoryLot) ConsumeQuantity(amount float64) bool {
	if il.CurrentQuantity < amount {
		return false // Not enough quantity available
	}
	il.CurrentQuantity -= amount

	// Mark as inactive if completely consumed
	if il.CurrentQuantity <= 0 {
		il.IsActive = false
	}
	return true
}

// GetTotalValue calculates the total value of this lot
// This is useful for inventory valuation
func (il *InventoryLot) GetTotalValue() float64 {
	return il.CurrentQuantity * il.CostPerBaseUnit
}

// MenuItem represents a product that can be sold to customers
// Menu items are organized by categories and have pricing information
// They can be linked to inventory items through recipes
type MenuItem struct {
	gorm.Model
	AccountID  uint    `json:"account_id" gorm:"not null;index"`
	Name       string  `json:"name" gorm:"not null"`
	Price      float64 `json:"price" gorm:"not null;default:0"`
	Category   string  `json:"category"`                 // e.g., "drinks", "food", "desserts"
	CategoryID *uint   `json:"category_id" gorm:"index"` // Optional category assignment
	// Note: Foreign key relationships are handled in application logic for ramsql compatibility
}

// RecipeIngredient links menu items to their required inventory variants
// This allows for automatic inventory tracking when menu items are sold
// The quantity field specifies how much of the inventory variant is needed per menu item
type RecipeIngredient struct {
	gorm.Model
	MenuItemID          uint    `json:"menu_item_id" gorm:"not null;index"`
	InventoryVariantID  uint    `json:"inventory_variant_id" gorm:"not null;index"`
	QuantityInBaseUnits float64 `json:"quantity_in_base_units" gorm:"not null;default:0"` // Quantity needed in base units per menu item
	// Note: Foreign key relationships are handled in application logic for ramsql compatibility
}

// InventorySnapshot represents a point-in-time snapshot of inventory levels
// This is used for historical tracking and reporting
// The Counts field stores a map of inventory variant IDs to their quantities at the snapshot time
type InventorySnapshot struct {
	gorm.Model
	AccountID int       `json:"account_id" gorm:"not null;index"`
	Timestamp time.Time `json:"timestamp" gorm:"not null;index"`
	Counts    CountsMap `json:"counts" gorm:"type:text"`          // map[InventoryVariantID]quantity_in_base_units - stored as JSON
	CountedBy int       `json:"counted_by" gorm:"not null;index"` // User ID who performed the count
	// Note: Foreign key relationships are handled in application logic for ramsql compatibility
}

// Delivery represents a shipment of inventory variants from a vendor
// This tracks when items are received and their associated costs
// Used for inventory replenishment and cost tracking
type Delivery struct {
	gorm.Model
	AccountID           uint      `json:"account_id" gorm:"not null;index"`
	InventoryVariantID  uint      `json:"inventory_variant_id" gorm:"not null;index"`
	VendorID            uint      `json:"vendor_id" gorm:"not null;index"`                  // Reference to the vendor
	QuantityInBaseUnits float64   `json:"quantity_in_base_units" gorm:"not null;default:0"` // Quantity received in base units (source of truth)
	DeliveryDate        time.Time `json:"delivery_date" gorm:"not null;index"`
	CostPerPackage      float64   `json:"cost_per_package" gorm:"not null;default:0"` // Cost per package
	TotalCost           float64   `json:"total_cost" gorm:"not null;default:0"`       // Total cost for the delivery
	ReceivedBy          uint      `json:"received_by" gorm:"not null;index"`          // User ID who received the delivery
	// Note: Foreign key relationships are handled in application logic for ramsql compatibility
}

type Sale struct {
	gorm.Model
	AccountID    uint       `json:"account_id" gorm:"not null;index"`
	SaleDate     time.Time  `json:"sale_date" gorm:"not null"`
	TotalRevenue float64    `json:"total_revenue"`
	TotalCost    float64    `json:"total_cost"`
	TotalProfit  float64    `json:"total_profit"`
	Notes        string     `json:"notes"`
	Items        []SaleItem `json:"items"`
}

type SaleItem struct {
	gorm.Model
	AccountID   uint     `json:"account_id" gorm:"not null;index"`
	SaleID      uint     `json:"sale_id" gorm:"not null;index"`
	MenuItemID  uint     `json:"menu_item_id" gorm:"not null;index"`
	MenuItem    MenuItem `json:"menu_item"`
	Quantity    int      `json:"quantity"`
	PriceAtSale float64  `json:"price_at_sale"`
	CostAtSale  float64  `json:"cost_at_sale"`
}

// Order represents a purchase order for inventory items
// Orders go through various statuses from pending to delivered
// They can be created by users and approved by managers
type Order struct {
	gorm.Model
	AccountID    uint      `json:"account_id" gorm:"not null;index"`
	Status       string    `json:"status" gorm:"not null;default:'pending'"` // pending, approved, ordered, delivered, cancelled
	OrderDate    time.Time `json:"order_date" gorm:"not null"`
	ExpectedDate time.Time `json:"expected_date"`
	TotalCost    float64   `json:"total_cost" gorm:"not null;default:0"`
	Notes        string    `json:"notes"`
	CreatedBy    uint      `json:"created_by" gorm:"not null"`
	ApprovedBy   *uint     `json:"approved_by"`
	// Note: Foreign key relationships are handled in application logic for ramsql compatibility
}

// OrderItem represents a specific variant within an order
// Each order can contain multiple variants from different vendors
// This tracks the quantity, cost, and vendor for each variant
type OrderItem struct {
	gorm.Model
	OrderID             uint    `json:"order_id" gorm:"not null;index"`
	InventoryVariantID  uint    `json:"inventory_variant_id" gorm:"not null;index"`
	VendorID            uint    `json:"vendor_id" gorm:"not null;index"`                  // Reference to the vendor
	QuantityInBaseUnits float64 `json:"quantity_in_base_units" gorm:"not null;default:0"` // Quantity ordered in base units (source of truth)
	CostPerPackage      float64 `json:"cost_per_package" gorm:"not null;default:0"`       // Cost per package
	TotalCost           float64 `json:"total_cost" gorm:"not null;default:0"`             // Total cost for this item
	Notes               string  `json:"notes"`
	// Note: Foreign key relationships are handled in application logic for ramsql compatibility
}

// OrderRequest represents a request for inventory items that needs approval
// This allows for a workflow where users can request items that need manager approval
// Requests can have different priorities and deadlines
type OrderRequest struct {
	gorm.Model
	AccountID   uint      `json:"account_id" gorm:"not null;index"`
	Status      string    `json:"status" gorm:"not null;default:'pending'"`  // pending, approved, rejected, fulfilled
	Priority    string    `json:"priority" gorm:"not null;default:'normal'"` // low, normal, high, urgent
	RequestDate time.Time `json:"request_date" gorm:"not null"`
	NeededBy    time.Time `json:"needed_by"`
	Notes       string    `json:"notes"`
	CreatedBy   uint      `json:"created_by" gorm:"not null"`
	ApprovedBy  *uint     `json:"approved_by"`
	// Note: Foreign key relationships are handled in application logic for ramsql compatibility
}

// RequestItem represents a specific variant within an order request
// Each request can contain multiple variants with different priorities
// This tracks the quantity needed and the reason for the request
type RequestItem struct {
	gorm.Model
	OrderRequestID      uint    `json:"order_request_id" gorm:"not null;index"`
	InventoryVariantID  uint    `json:"inventory_variant_id" gorm:"not null;index"`
	QuantityInBaseUnits float64 `json:"quantity_in_base_units" gorm:"not null;default:0"` // Quantity needed in base units (source of truth)
	Reason              string  `json:"reason"`                                           // e.g., "low stock", "new menu item", "special event"
	Priority            string  `json:"priority" gorm:"not null;default:'normal'"`
	// Note: Foreign key relationships are handled in application logic for ramsql compatibility
}

// AccountInvitation represents an invitation for a user to join an account
// This allows account admins to invite users by email without requiring them to know account IDs
// The invitation system provides better security and user experience
type AccountInvitation struct {
	gorm.Model
	AccountID  uint       `json:"account_id" gorm:"not null;index"`
	Email      string     `json:"email" gorm:"not null;index"`
	InvitedBy  uint       `json:"invited_by" gorm:"not null"`               // User ID who sent the invitation
	Status     string     `json:"status" gorm:"not null;default:'pending'"` // pending, accepted, expired, revoked
	InvitedAt  time.Time  `json:"invited_at" gorm:"autoCreateTime"`
	AcceptedAt *time.Time `json:"accepted_at"`
	ExpiresAt  time.Time  `json:"expires_at" gorm:"not null"` // Invitation expiration date
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

// EmailVerificationToken represents a temporary token for email verification
// This allows users to verify their email addresses securely
type EmailVerificationToken struct {
	gorm.Model
	UserID    uint       `json:"user_id" gorm:"not null;index"`
	Token     string     `json:"token" gorm:"not null;uniqueIndex"`
	Type      string     `json:"type" gorm:"not null"` // "email_verification", "password_reset"
	ExpiresAt time.Time  `json:"expires_at" gorm:"not null;index"`
	UsedAt    *time.Time `json:"used_at"`
	// Note: Foreign key relationships are handled in application logic for ramsql compatibility
}

// EmailLog represents a record of sent emails for tracking and debugging
// This helps track email delivery and troubleshoot issues
type EmailLog struct {
	gorm.Model
	AccountID uint      `json:"account_id" gorm:"not null;index"`
	UserID    *uint     `json:"user_id" gorm:"index"` // Optional - null for bulk emails
	ToEmail   string    `json:"to_email" gorm:"not null"`
	Subject   string    `json:"subject" gorm:"not null"`
	EmailType string    `json:"email_type" gorm:"not null"`            // "verification", "weekly_report", "low_stock_alert"
	Status    string    `json:"status" gorm:"not null;default:'sent'"` // sent, failed, pending
	ErrorMsg  string    `json:"error_msg"`                             // Error message if sending failed
	SentAt    time.Time `json:"sent_at" gorm:"autoCreateTime"`
	// Note: Foreign key relationships are handled in application logic for ramsql compatibility
}

// EmailSchedule represents scheduled email tasks
// This allows for automated email sending based on schedules
type EmailSchedule struct {
	gorm.Model
	AccountID  uint       `json:"account_id" gorm:"not null;index"`
	EmailType  string     `json:"email_type" gorm:"not null"`  // "weekly_stock_report", "low_stock_alert"
	Frequency  string     `json:"frequency" gorm:"not null"`   // "weekly", "daily", "monthly"
	DayOfWeek  *uint      `json:"day_of_week"`                 // 0-6 (Sunday-Saturday) for weekly
	DayOfMonth *uint      `json:"day_of_month"`                // 1-31 for monthly
	TimeOfDay  string     `json:"time_of_day" gorm:"not null"` // "09:00", "18:30"
	IsActive   bool       `json:"is_active" gorm:"not null;default:true"`
	LastSentAt *time.Time `json:"last_sent_at"`
	// Note: Foreign key relationships are handled in application logic for ramsql compatibility
}

// Email status constants
const (
	EmailStatusSent    = "sent"
	EmailStatusFailed  = "failed"
	EmailStatusPending = "pending"
)

// Email type constants
const (
	EmailTypeVerification      = "verification"
	EmailTypeWeeklyReport      = "weekly_stock_report"
	EmailTypeWeeklySupplyChain = "weekly_supply_chain_report"
	EmailTypeLowStockAlert     = "low_stock_alert"
	EmailTypePasswordReset     = "password_reset"
	EmailTypeAccountInvite     = "account_invite"
)

// Token type constants
const (
	TokenTypeEmailVerification = "email_verification"
	TokenTypePasswordReset     = "password_reset"
)

// Unit type constants
const (
	UnitTypeVolume = "volume" // e.g., gallon, liter, ounce
	UnitTypeWeight = "weight" // e.g., pound, kilogram, gram
	UnitTypeCount  = "count"  // e.g., piece, box, bottle
	UnitTypeLength = "length" // e.g., meter, foot, inch
)

// Priority constants
const (
	PriorityLow    = "low"
	PriorityNormal = "normal"
	PriorityHigh   = "high"
	PriorityUrgent = "urgent"
)
