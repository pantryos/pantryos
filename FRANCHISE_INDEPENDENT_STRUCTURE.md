# Franchise vs Independent Location Structure

## Overview

PantryOS needs to support two distinct business models:
1. **Independent Locations** - Single business owners
2. **Franchise Operations** - Multi-location businesses with franchisor oversight

## Business Models

### Independent Location Model
```
User (Owner) → Account (Single Location)
```
- Simple, direct relationship
- Owner has full control over their location
- No organizational hierarchy needed

### Franchise Model
```
Organization (Franchisor) → Multiple Accounts (Franchise Locations)
├── Franchisor Users (Corporate) → Organization + All Accounts
└── Franchisee Users (Location Owners) → Their Specific Account
```

## Refined Data Structure

### Core Entities

#### User
```go
type User struct {
    ID         int       `json:"id"`
    Email      string    `json:"email"`
    Password   string    `json:"-"` // Hidden from JSON
    FirstName  string    `json:"first_name"`
    LastName   string    `json:"last_name"`
    IsVerified bool      `json:"is_verified"`
    Status     string    `json:"status"` // active, inactive, suspended
    CreatedAt  time.Time `json:"created_at"`
    UpdatedAt  time.Time `json:"updated_at"`
}
```

#### Organization (Franchisor)
```go
type Organization struct {
    ID          int       `json:"id"`
    Name        string    `json:"name"` // e.g., "Coffee Chain Corp"
    Description string    `json:"description"`
    Type        string    `json:"type"` // franchise, enterprise, multi_location
    Status      string    `json:"status"` // active, inactive, suspended
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}
```

#### Account (Business Location)
```go
type Account struct {
    ID             int       `json:"id"`
    OrganizationID *int      `json:"organization_id"` // null for independent
    Name           string    `json:"name"` // e.g., "Downtown Coffee Shop"
    Location       string    `json:"location"`
    Phone          string    `json:"phone"`
    Email          string    `json:"email"`
    BusinessType   string    `json:"business_type"` // independent, franchise_location
    Status         string    `json:"status"` // active, inactive, suspended
    CreatedAt      time.Time `json:"created_at"`
    UpdatedAt      time.Time `json:"updated_at"`
}
```

### Relationship Tables

#### UserAccount (User ↔ Account)
```go
type UserAccount struct {
    ID        int       `json:"id"`
    UserID    int       `json:"user_id"`
    AccountID int       `json:"account_id"`
    Role      string    `json:"role"` // owner, manager, employee
    IsPrimary bool      `json:"is_primary"` // Primary account for the user
    Status    string    `json:"status"` // active, inactive, suspended
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

#### UserOrganization (User ↔ Organization)
```go
type UserOrganization struct {
    ID             int       `json:"id"`
    UserID         int       `json:"user_id"`
    OrganizationID int       `json:"organization_id"`
    Role           string    `json:"role"` // franchisor, franchise_admin, franchise_support
    Status         string    `json:"status"` // active, inactive, suspended
    CreatedAt      time.Time `json:"created_at"`
    UpdatedAt      time.Time `json:"updated_at"`
}
```

## Access Patterns

### Independent Location Access
```
User (Owner) → UserAccount (owner role) → Account
```
- Owner has full access to their single location
- Can invite employees with limited roles
- No organizational oversight

### Franchise Access Patterns

#### Franchisor Access
```
User (Franchisor) → UserOrganization (franchisor role) → Organization
User (Franchisor) → UserAccount (owner role) → All Accounts
```
- Can see all franchise locations
- Can access reports across all locations
- Can manage franchise-wide settings

#### Franchisee Access
```
User (Franchisee) → UserAccount (owner role) → Their Specific Account
User (Franchisee) → UserOrganization (franchisee role) → Organization (limited)
```
- Full access to their own location
- Limited access to franchise-wide data
- Can see some shared resources (branding, menus, etc.)

#### Franchise Support Access
```
User (Support) → UserOrganization (franchise_support role) → Organization
User (Support) → UserAccount (manager role) → Multiple Accounts (as needed)
```
- Can help multiple franchise locations
- Limited access to each location
- Can access franchise-wide support tools

## Permission System

### Permission Scopes
- **Account**: Location-specific permissions
- **Organization**: Franchise-wide permissions
- **System**: Platform-wide permissions

### Role Hierarchy

#### Independent Location Roles
- **owner**: Full access to location
- **manager**: Manage inventory, orders, employees
- **employee**: Basic operations (view inventory, create orders)

#### Franchise Roles
- **franchisor**: Full access to all locations and organization
- **franchise_admin**: Manage franchise operations
- **franchise_support**: Support multiple locations
- **franchisee**: Owner of specific location
- **franchise_manager**: Manager of specific location

## Implementation Strategy

### Phase 1: Independent Location Support
1. Simple user → account relationship
2. Basic role system (owner, manager, employee)
3. Single location focus

### Phase 2: Franchise Support
1. Add organization structure
2. Implement cross-location access
3. Add franchise-specific roles and permissions

### Phase 3: Advanced Features
1. Franchise-wide reporting
2. Brand consistency tools
3. Multi-location analytics

## Use Case Examples

### Independent Coffee Shop
```
Owner: john@coffeeshop.com
├── Account: "Downtown Coffee Shop"
├── Role: owner
└── Access: Full control over location
```

### Franchise Chain
```
Organization: "Coffee Chain Corp"
├── Franchisor: corporate@coffeechain.com
│   ├── Organization Role: franchisor
│   └── Access: All locations + franchise management
├── Location 1: "Downtown Franchise"
│   ├── Franchisee: mike@downtown.com
│   │   ├── Account Role: owner
│   │   └── Organization Role: franchisee
│   └── Access: Full control over Downtown location
└── Location 2: "Uptown Franchise"
    ├── Franchisee: sarah@uptown.com
    │   ├── Account Role: owner
    │   └── Organization Role: franchisee
    └── Access: Full control over Uptown location
```

## Benefits of This Approach

### For Independent Locations
- **Simple setup**: Direct user → account relationship
- **Full control**: No organizational overhead
- **Cost-effective**: Pay only for what you need

### For Franchises
- **Scalable**: Easy to add new franchise locations
- **Flexible access**: Different roles for different needs
- **Brand consistency**: Shared resources and standards
- **Insights**: Cross-location reporting and analytics

### For Platform
- **Unified codebase**: Same system serves both models
- **Flexible permissions**: Granular access control
- **Future-proof**: Can evolve from independent to franchise

## Migration Considerations

### From Independent to Franchise
1. Create organization for the franchisor
2. Link existing account to organization
3. Update user roles and permissions
4. Add franchise-specific features

### From Franchise to Independent
1. Remove organization relationships
2. Simplify user roles
3. Maintain account data integrity
4. Update access patterns

## Next Steps

1. **Validate this structure** with potential customers
2. **Implement Phase 1** (independent locations)
3. **Test with real users** before adding franchise complexity
4. **Iterate based on feedback** and actual usage patterns

This structure provides the flexibility to serve both business models while maintaining simplicity for independent locations and power for franchise operations.
