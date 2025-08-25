# Inventory System Redesign: Units and Variants

## Overview

The inventory system has been redesigned to handle complex real-world scenarios where:
- Items come in packages (e.g., a box of 4 items, each 1 gallon)
- Different brands/variants of the same base item exist
- Unit conversions are needed between package units and individual units

## New Model Structure

### 1. Unit System

#### Unit Model
Represents units of measurement that can be used across the system.

```go
type Unit struct {
    ID         int       `json:"id"`
    AccountID  int       `json:"account_id"`
    Name       string    `json:"name"`        // e.g., "gallon", "pound", "piece"
    Symbol     string    `json:"symbol"`      // e.g., "gal", "lb", "pc"
    Type       string    `json:"type"`        // "volume", "weight", "count", "length"
    IsBaseUnit bool      `json:"is_base_unit"` // Whether this is a base unit for conversions
}
```

#### UnitConversion Model
Handles conversion factors between different units.

```go
type UnitConversion struct {
    ID             int     `json:"id"`
    FromUnitID     int     `json:"from_unit_id"`
    ToUnitID       int     `json:"to_unit_id"`
    ConversionRate float64 `json:"conversion_rate"` // How many 'to' units equal 1 'from' unit
}
```

### 2. Inventory Hierarchy

#### InventoryItem (Base Item)
Represents the core item definition (e.g., "Milk", "Coffee Beans", "Sugar").

```go
type InventoryItem struct {
    ID            int     `json:"id"`
    Name          string  `json:"name"`           // e.g., "Milk", "Coffee Beans"
    Description   string  `json:"description"`
    BaseUnitID    int     `json:"base_unit_id"`   // Primary unit (e.g., gallon for milk)
    CategoryID    *int    `json:"category_id"`
    MinStockLevel float64 `json:"min_stock_level"` // In base units
    MaxStockLevel float64 `json:"max_stock_level"` // In base units
}
```

#### Vendor (Supplier)
Represents a supplier or vendor that provides inventory items.

```go
type Vendor struct {
    ID           int     `json:"id"`
    AccountID    int     `json:"account_id"`
    Name         string  `json:"name"`              // e.g., "Coffee Supply Co.", "Local Dairy"
    ContactName  string  `json:"contact_name"`      // Primary contact person
    Email        string  `json:"email"`             // Contact email
    Phone        string  `json:"phone"`             // Contact phone
    Website      string  `json:"website"`           // Vendor website
    Address      string  `json:"address"`           // Physical address
    PaymentTerms string  `json:"payment_terms"`     // Payment terms (e.g., "Net 30", "COD")
    IsActive     bool    `json:"is_active"`
}
```

#### InventoryVariant (Specific Variant)
Represents a specific brand, size, or packaging of a base item.

```go
type InventoryVariant struct {
    ID                int     `json:"id"`
    InventoryItemID   int     `json:"inventory_item_id"`
    Brand             string  `json:"brand"`              // e.g., "Organic Valley", "Starbucks"
    VariantName       string  `json:"variant_name"`       // e.g., "2% Milk", "Dark Roast"
    SKU               string  `json:"sku"`                // Stock Keeping Unit
    PackageUnitID     int     `json:"package_unit_id"`    // Unit of package (e.g., "box", "bottle")
    PackageSize       float64 `json:"package_size"`       // Size in package units
    BaseUnitID        int     `json:"base_unit_id"`       // Unit of individual items (e.g., "gallon")
    ItemsPerPackage   float64 `json:"items_per_package"`  // How many base units per package
    CostPerPackage    float64 `json:"cost_per_package"`   // Cost for entire package
    CostPerBaseUnit   float64 `json:"cost_per_base_unit"` // Calculated cost per base unit
    PreferredVendorID *int    `json:"preferred_vendor_id"` // Default supplier for this variant
    Link              string  `json:"link"`               // URL to product page or catalog
}
```

#### InventoryLot (Batch Tracking)
Tracks individual batches of inventory received from deliveries, enabling FIFO/FEFO management.

```go
type InventoryLot struct {
    ID                  int        `json:"id"`
    AccountID           int        `json:"account_id"`
    InventoryVariantID  int        `json:"inventory_variant_id"`
    DeliveryID          int        `json:"delivery_id"`          // Reference to the delivery
    InitialQuantity     float64    `json:"initial_quantity"`     // Initial quantity in base units
    CurrentQuantity     float64    `json:"current_quantity"`     // Current quantity in base units
    CostPerBaseUnit     float64    `json:"cost_per_base_unit"`   // Cost for this specific batch
    ReceivedAt          time.Time  `json:"received_at"`          // When this lot was received
    ExpirationDate      *time.Time `json:"expiration_date"`      // Expiration date for perishable goods
    IsActive            bool       `json:"is_active"`            // Whether this lot is still active
}
```

#### InventoryStock (Calculated Stock)
Represents the calculated sum of all active lots for a variant.

```go
type InventoryStock struct {
    ID                  int     `json:"id"`
    InventoryVariantID  int     `json:"inventory_variant_id"`
    QuantityInBaseUnits float64 `json:"quantity_in_base_units"` // Calculated sum of all active lots
}
```

## Example Usage

### Scenario: Milk Inventory

#### 1. Define Units
```go
// Base unit for milk
gallonUnit := Unit{
    Name:       "gallon",
    Symbol:     "gal",
    Type:       "volume",
    IsBaseUnit: true,
}

// Package unit for milk boxes
boxUnit := Unit{
    Name:       "box",
    Symbol:     "box",
    Type:       "count",
    IsBaseUnit: false,
}
```

#### 2. Define Base Item
```go
milkItem := InventoryItem{
    Name:        "Milk",
    Description: "Dairy milk for coffee and beverages",
    BaseUnitID:  gallonUnit.ID, // Use gallon as base unit
    MinStockLevel: 10.0,        // Minimum 10 gallons
    MaxStockLevel: 50.0,        // Maximum 50 gallons
}
```

#### 3. Define Vendors
```go
// Coffee Supply Company
coffeeSupplyVendor := Vendor{
    Name:         "Coffee Supply Co.",
    ContactName:  "John Smith",
    Email:        "john@coffeesupply.com",
    Phone:        "(555) 123-4567",
    Website:      "https://coffeesupply.com",
    Address:      "123 Coffee St, Bean City, BC 12345",
    PaymentTerms: "Net 30",
    IsActive:     true,
}

// Local Dairy
localDairyVendor := Vendor{
    Name:         "Local Dairy",
    ContactName:  "Mary Johnson",
    Email:        "mary@localdairy.com",
    Phone:        "(555) 987-6543",
    Website:      "https://localdairy.com",
    Address:      "456 Milk Ave, Dairy Town, DT 67890",
    PaymentTerms: "Net 15",
    IsActive:     true,
}
```

#### 4. Define Variants
```go
// Organic Valley 2% Milk - comes in boxes of 4 gallons
organicMilkVariant := InventoryVariant{
    InventoryItemID:   milkItem.ID,
    Brand:             "Organic Valley",
    VariantName:       "2% Milk",
    SKU:               "OV-MILK-2P-4GAL",
    PackageUnitID:     boxUnit.ID,
    PackageSize:       1.0,           // 1 box
    BaseUnitID:        gallonUnit.ID,
    ItemsPerPackage:   4.0,           // 4 gallons per box
    CostPerPackage:    12.00,         // $12 per box
    CostPerBaseUnit:   3.00,          // $3 per gallon
    PreferredVendorID: &localDairyVendor.ID,
    Link:              "https://organicvalley.coop/products/milk/2-percent-milk/",
}

// Starbucks Whole Milk - comes in boxes of 6 gallons
starbucksMilkVariant := InventoryVariant{
    InventoryItemID:   milkItem.ID,
    Brand:             "Starbucks",
    VariantName:       "Whole Milk",
    SKU:               "SB-MILK-WH-6GAL",
    PackageUnitID:     boxUnit.ID,
    PackageSize:       1.0,           // 1 box
    BaseUnitID:        gallonUnit.ID,
    ItemsPerPackage:   6.0,           // 6 gallons per box
    CostPerPackage:    18.00,         // $18 per box
    CostPerBaseUnit:   3.00,          // $3 per gallon
    PreferredVendorID: &localDairyVendor.ID,
    Link:              "https://starbucks.com/products/whole-milk/",
}
```

#### 5. Track Stock and Deliveries with Batch Tracking
```go
// Delivery from Local Dairy
milkDelivery := Delivery{
    InventoryVariantID:  organicMilkVariant.ID,
    VendorID:            localDairyVendor.ID,
    QuantityInBaseUnits: 20.0,   // 20 gallons (source of truth)
    DeliveryDate:        time.Now(),
    CostPerPackage:      12.00,  // $12 per box
    TotalCost:           60.00,  // $60 total
    ReceivedBy:          userID, // User who received the delivery
}

// Create inventory lot from delivery
milkLot := InventoryLot{
    InventoryVariantID:  organicMilkVariant.ID,
    DeliveryID:          milkDelivery.ID,
    InitialQuantity:     20.0,   // 20 gallons
    CurrentQuantity:     20.0,   // 20 gallons (initially same as initial)
    CostPerBaseUnit:     3.00,   // $3 per gallon
    ReceivedAt:          time.Now(),
    ExpirationDate:      &expirationDate, // Set expiration for perishable goods
    IsActive:            true,
}

// Calculate package quantities when needed
deliveryPackages := milkDelivery.GetQuantityInPackages(&organicMilkVariant) // Returns 5.0 boxes
lotPackages := milkLot.GetQuantityInPackages(&organicMilkVariant)           // Returns 5.0 boxes

// Check expiration status
if milkLot.IsExpired() {
    // Handle expired lot
}
daysUntilExpiry := milkLot.DaysUntilExpiration() // Returns days until expiration

// Consume from lot (FIFO/FEFO)
success := milkLot.ConsumeQuantity(5.0) // Consume 5 gallons
if success {
    // Update inventory stock (sum of all active lots)
    // This would be calculated automatically in the application
}
```

### Scenario: Coffee Beans

#### 1. Define Units
```go
poundUnit := Unit{
    Name:       "pound",
    Symbol:     "lb",
    Type:       "weight",
    IsBaseUnit: true,
}

bagUnit := Unit{
    Name:       "bag",
    Symbol:     "bag",
    Type:       "count",
    IsBaseUnit: false,
}
```

#### 2. Define Base Item
```go
coffeeItem := InventoryItem{
    Name:        "Coffee Beans",
    Description: "Coffee beans for brewing",
    BaseUnitID:  poundUnit.ID,
    MinStockLevel: 20.0,  // Minimum 20 pounds
    MaxStockLevel: 100.0, // Maximum 100 pounds
}
```

#### 3. Define Variants
```go
// Starbucks Dark Roast - comes in 5-pound bags
starbucksDarkVariant := InventoryVariant{
    InventoryItemID: coffeeItem.ID,
    Brand:           "Starbucks",
    VariantName:     "Dark Roast",
    SKU:             "SB-COF-DR-5LB",
    PackageUnitID:   bagUnit.ID,
    PackageSize:     1.0,           // 1 bag
    BaseUnitID:      poundUnit.ID,
    ItemsPerPackage: 5.0,           // 5 pounds per bag
    CostPerPackage:  25.00,         // $25 per bag
    CostPerBaseUnit: 5.00,          // $5 per pound
}

// Local Roaster Medium Roast - comes in 10-pound bags
localMediumVariant := InventoryVariant{
    InventoryItemID: coffeeItem.ID,
    Brand:           "Local Roaster",
    VariantName:     "Medium Roast",
    SKU:             "LR-COF-MR-10LB",
    PackageUnitID:   bagUnit.ID,
    PackageSize:     1.0,           // 1 bag
    BaseUnitID:      poundUnit.ID,
    ItemsPerPackage: 10.0,          // 10 pounds per bag
    CostPerPackage:  45.00,         // $45 per bag
    CostPerBaseUnit: 4.50,          // $4.50 per pound
}
```

## Benefits of This Design

### 1. Flexible Unit Handling
- Supports any unit type (volume, weight, count, length)
- Automatic unit conversions
- Clear distinction between package units and base units

### 2. Brand/Variant Management
- Multiple brands of the same item
- Different package sizes and costs
- Unique SKUs for each variant

### 3. Accurate Cost Tracking
- Cost per package and cost per base unit
- Automatic calculation of base unit costs
- Support for different pricing strategies

### 4. Inventory Management
- Track stock in both package and base units
- Automatic synchronization between units
- Support for partial packages

### 5. Recipe Integration
- Recipes can specify exact base unit quantities
- Automatic inventory deduction when menu items are sold
- Support for multiple variants in recipes

### 6. Vendor Management
- Centralized vendor information with contact details
- Payment terms tracking for each vendor
- Vendor relationship management
- Easy vendor lookup and reporting

### 7. Data Normalization
- Single source of truth for quantities (base units)
- Eliminates data inconsistency between package and base units
- Package quantities calculated on-demand when needed
- Reduced storage requirements and complexity

### 8. Batch Tracking (FIFO/FEFO)
- Individual lot tracking for each delivery
- Expiration date management for perishable goods
- Accurate cost tracking per batch
- FIFO/FEFO inventory depletion
- Perfect profit margin calculations

### 9. Audit Trail
- User tracking for critical operations (deliveries, inventory counts)
- Complete audit trail for compliance and accountability
- Soft delete functionality for data recovery

### 10. Model Consistency
- All models use `gorm.Model` for consistency
- Automatic `ID`, `CreatedAt`, `UpdatedAt`, and `DeletedAt` fields
- Uniform soft delete capabilities across all entities
- Simplified model definitions and maintenance

### 11. Database-Level Foreign Key Constraints
- **Referential Integrity**: All relationships enforced at database level
- **Cascade Operations**: Proper CASCADE, SET NULL, and RESTRICT behaviors
- **Orphan Prevention**: Automatic prevention of orphaned records
- **Data Corruption Prevention**: Database-level safety net against data inconsistencies
- **Production Ready**: Enterprise-grade data integrity guarantees

## Foreign Key Constraint Strategy

### Cascade Behaviors by Entity Type:

#### **Core Entities (CASCADE on delete)**
- **Account**: When deleted, all related data is removed
- **Organization**: When deleted, all accounts and related data are removed
- **InventoryItem**: When deleted, all variants and related data are removed

#### **Reference Entities (RESTRICT on delete)**
- **Unit**: Cannot be deleted if referenced by inventory items
- **Category**: Cannot be deleted if referenced by inventory items

#### **Optional Relationships (SET NULL on delete)**
- **OrganizationID**: When organization is deleted, account becomes standalone
- **CategoryID**: When category is deleted, items become uncategorized
- **PreferredVendorID**: When vendor is deleted, variant has no preferred vendor

#### **Audit Entities (CASCADE on delete)**
- **Delivery**: When deleted, all related lots are removed
- **Order**: When deleted, all order items are removed
- **Sale**: When deleted, all sale items are removed

## Migration Strategy

### Phase 1: Database Migration
1. Create new tables (Unit, UnitConversion, InventoryVariant, InventoryStock)
2. Migrate existing InventoryItem data to new structure
3. Create default units and variants for existing items

### Phase 2: API Updates
1. Update inventory handlers to work with variants
2. Add unit management endpoints
3. Update delivery and order systems

### Phase 3: Frontend Updates
1. Update inventory management interface
2. Add unit and variant management screens
3. Update reporting and analytics

### Phase 4: Testing and Validation
1. Test unit conversions
2. Validate inventory calculations
3. Ensure data integrity

## Helper Methods

The system includes helper methods for common operations:

```go
// Calculate base units from package units
baseUnits := variant.CalculateBaseUnits(packageQuantity)

// Calculate package units from base units
packageUnits := variant.CalculatePackageUnits(baseQuantity)

// Update stock levels
stock.UpdateStock(packageQuantity, variant)

// Get display name for variant
displayName := variant.GetDisplayName()

// Calculate package quantities when needed
packageQuantity := stock.GetQuantityInPackages(variant)
deliveryPackages := delivery.GetQuantityInPackages(variant)
orderPackages := orderItem.GetQuantityInPackages(variant)
requestPackages := requestItem.GetQuantityInPackages(variant)
lotPackages := lot.GetQuantityInPackages(variant)

// Batch tracking operations
if lot.IsExpired() {
    // Handle expired lot
}
daysUntilExpiry := lot.DaysUntilExpiration()
success := lot.ConsumeQuantity(amount)
totalValue := lot.GetTotalValue()
```

This design provides a robust foundation for handling complex inventory scenarios while maintaining data integrity and providing clear business logic.
