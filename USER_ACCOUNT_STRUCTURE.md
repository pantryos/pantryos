# User, Organization, and Account Structure

## Overview

PantryOS supports two distinct business models through a flexible user, organization, and account structure:

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
- Perfect for single coffee shops, restaurants, etc.

### Franchise Model
```
Organization (Franchisor) → Multiple Accounts (Franchise Locations)
├── Franchisor Users (Corporate) → Organization + All Accounts
└── Franchisee Users (Location Owners) → Their Specific Account
```
- Scalable for growing franchise operations
- Flexible access control for different user types
- Brand consistency across locations
- Cross-location reporting and insights

## Data Model

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
    Role           string    `json:"role"` // franchisor, franchise_admin, franchise_support, franchisee
    Status         string    `json:"status"` // active, inactive, suspended
    CreatedAt      time.Time `json:"created_at"`
    UpdatedAt      time.Time `json:"updated_at"`
}
```

### Permission System

#### Permission
```go
type Permission struct {
    ID          int       `json:"id"`
    Name        string    `json:"name"` // e.g., "inventory.read"
    Description string    `json:"description"`
    Resource    string    `json:"resource"` // e.g., "inventory", "users"
    Action      string    `json:"action"` // e.g., "read", "write", "delete"
    CreatedAt   time.Time `json:"created_at"`
}
```

#### RolePermission
```go
type RolePermission struct {
    ID           int       `json:"id"`
    Role         string    `json:"role"` // user, manager, admin, org_admin, org_owner
    PermissionID int       `json:"permission_id"`
    Scope        string    `json:"scope"` // account, organization, system
    CreatedAt    time.Time `json:"created_at"`
}
```

## Role System

### Account-Level Roles
- **owner**: Full access to location (business owner/franchisee)
- **manager**: Manage inventory, orders, employees
- **employee**: Basic operations (view inventory, create orders)

### Organization-Level Roles
- **franchisor**: Full access to all locations and organization
- **franchise_admin**: Manage franchise operations
- **franchise_support**: Support multiple locations
- **franchisee**: Owner of specific franchise location

### Business Types
- **independent**: Standalone business (no organization)
- **franchise**: Franchise location under organization
- **multi_location**: Multiple locations under one organization (non-franchise)
- **enterprise**: Large enterprise with complex hierarchy

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

## Benefits

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

## Migration Strategy

### Database Migration
- Run migration script: `scripts/migrate_user_structure.sql`
- Migrate existing user-account relationships to UserAccount table
- Set up default permissions for existing roles
- Create organization relationships where needed

### Backend Updates
- Update authentication to support multiple accounts
- Implement new permission checking system
- Add account switching functionality
- Update user management endpoints

### Frontend Updates
- Update user interface for account switching
- Implement permission-based UI components
- Add organization management features
- Update user profile and settings

## Security Considerations

### Permission Inheritance
- Organization permissions don't automatically grant account permissions
- Explicit permission assignment required
- Audit trail for all permission changes

### Account Isolation
- Users can only access accounts they're explicitly assigned to
- Cross-account data access requires explicit permissions
- Proper validation at all API endpoints

### Token Security
- JWT tokens include account context
- Token expiration and refresh mechanisms
- Secure token storage and transmission

## Testing Strategy

### Unit Tests
- Test new permission system
- Test user-account relationships
- Test role validation

### Integration Tests
- Test authentication flows
- Test permission checking
- Test account switching

### End-to-End Tests
- Test complete user workflows
- Test multi-account scenarios
- Test organization management

## Performance Considerations

### Database Optimization
- Proper indexing on relationship tables
- Efficient permission queries
- Caching for frequently accessed permissions

### API Optimization
- Batch permission checking
- Efficient user account queries
- Proper pagination for large datasets

### Frontend Optimization
- Lazy loading of account data
- Efficient permission checking
- Optimistic UI updates

## Future Enhancements

### Advanced Permissions
- Time-based permissions
- Conditional permissions
- Permission templates

### Organization Features
- Organization hierarchies
- Department management
- Cross-organization reporting

### User Management
- Bulk user operations
- User provisioning workflows
- Advanced user profiles

## Conclusion

This structure provides the flexibility to serve both business models while maintaining simplicity for independent locations and power for franchise operations. The unified approach ensures a consistent user experience while supporting the specific needs of each business type.

The migration strategy ensures a smooth transition with minimal disruption, and the comprehensive testing approach will ensure reliability and security throughout the implementation process.
