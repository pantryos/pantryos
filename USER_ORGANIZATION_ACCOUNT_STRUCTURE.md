# User, Organization, and Account Structure Redesign

## Overview

This document outlines the redesigned user, organization, and account structure for PantryOS. The new architecture provides greater flexibility, better scalability, and improved permission management while maintaining backward compatibility where possible.

## Key Changes

### 1. Decoupled User-Account Relationship

**Before:**
- Users were tightly coupled to a single account
- Users had a direct `AccountID` field
- Users could only work within one business location

**After:**
- Users are decoupled from accounts through a many-to-many relationship
- Users can belong to multiple accounts with different roles
- Users have a `UserAccount` relationship table

### 2. Enhanced Role and Permission System

**Before:**
- Simple role-based access control (user, manager, admin, org_admin)
- Roles were tied to accounts only
- Limited permission granularity

**After:**
- Fine-grained permission system with `Permission` and `RolePermission` tables
- Roles can have different scopes (account, organization, system)
- Support for organization-wide roles through `UserOrganization`

### 3. Improved Organization Management

**Before:**
- Organizations were primarily for grouping accounts
- Limited organization-level functionality

**After:**
- Organizations can have their own users with organization-wide roles
- Organization invitations and management
- Better support for enterprise scenarios

## New Data Model

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

#### Organization
```go
type Organization struct {
    ID          int       `json:"id"`
    Name        string    `json:"name"`
    Description string    `json:"description"`
    Type        string    `json:"type"` // single_location, multi_location, enterprise
    Status      string    `json:"status"` // active, inactive, suspended
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}
```

#### Account
```go
type Account struct {
    ID             int       `json:"id"`
    OrganizationID *int      `json:"organization_id"` // Optional - null for standalone
    Name           string    `json:"name"`
    Location       string    `json:"location"`
    Phone          string    `json:"phone"`
    Email          string    `json:"email"`
    BusinessType   string    `json:"business_type"`
    Status         string    `json:"status"`
    CreatedAt      time.Time `json:"created_at"`
    UpdatedAt      time.Time `json:"updated_at"`
}
```

### Relationship Tables

#### UserAccount (Many-to-Many)
```go
type UserAccount struct {
    ID        int       `json:"id"`
    UserID    int       `json:"user_id"`
    AccountID int       `json:"account_id"`
    Role      string    `json:"role"` // user, manager, admin
    IsPrimary bool      `json:"is_primary"` // Primary account for the user
    Status    string    `json:"status"` // active, inactive, suspended
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

#### UserOrganization (Many-to-Many)
```go
type UserOrganization struct {
    ID             int       `json:"id"`
    UserID         int       `json:"user_id"`
    OrganizationID int       `json:"organization_id"`
    Role           string    `json:"role"` // member, admin, owner
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

## Use Cases and Scenarios

### 1. Single Location Business
- **Setup**: One account, no organization
- **Users**: Direct UserAccount relationships
- **Roles**: user, manager, admin
- **Example**: Independent coffee shop

### 2. Multi-Location Chain
- **Setup**: One organization, multiple accounts
- **Users**: Can belong to multiple accounts with different roles
- **Roles**: user, manager, admin (account-level) + member, admin, owner (org-level)
- **Example**: Coffee chain with multiple locations

### 3. Enterprise Organization
- **Setup**: One organization, many accounts, complex hierarchy
- **Users**: Organization-wide roles + account-specific roles
- **Roles**: Full role hierarchy with fine-grained permissions
- **Example**: Large restaurant group with corporate structure

## Benefits of New Structure

### 1. Flexibility
- Users can work across multiple locations
- Easy to transfer users between accounts
- Support for temporary assignments

### 2. Scalability
- Better support for growing organizations
- Easier to add new locations
- Improved performance with proper indexing

### 3. Security
- Fine-grained permission control
- Better audit trails
- Role-based access control with scopes

### 4. User Experience
- Single sign-on across multiple accounts
- Consistent user profiles
- Better invitation and onboarding flow

## Migration Strategy

### Phase 1: Database Schema Updates
1. Add new tables (UserAccount, UserOrganization, Permission, RolePermission)
2. Add new fields to existing tables
3. Create indexes for performance

### Phase 2: Data Migration
1. Migrate existing user-account relationships to UserAccount table
2. Set up default permissions for existing roles
3. Create organization relationships where needed

### Phase 3: API Updates
1. Update authentication to support multiple accounts
2. Modify user management endpoints
3. Add new permission-based authorization

### Phase 4: Frontend Updates
1. Update user interface for account switching
2. Implement permission-based UI components
3. Add organization management features

## Implementation Guidelines

### 1. Authentication Flow
```go
// New authentication flow
func (h *AuthHandler) Login(c *gin.Context) {
    // ... existing validation ...
    
    // Get user's accounts
    userAccounts, err := h.service.GetUserAccounts(user.ID)
    if err != nil {
        // Handle error
    }
    
    // Return user with account list
    response := gin.H{
        "token": token,
        "user": gin.H{
            "id": user.ID,
            "email": user.Email,
            "first_name": user.FirstName,
            "last_name": user.LastName,
            "accounts": userAccounts,
        },
    }
}
```

### 2. Permission Checking
```go
// New permission checking
func (h *Handler) CheckPermission(userID, accountID int, permission string) bool {
    // Get user's permissions for this account
    permissions, err := h.service.GetUserPermissions(userID, accountID)
    if err != nil {
        return false
    }
    
    // Check if user has the required permission
    for _, perm := range permissions {
        if perm.Name == permission {
            return true
        }
    }
    
    return false
}
```

### 3. Account Switching
```go
// New account switching endpoint
func (h *AuthHandler) SwitchAccount(c *gin.Context) {
    userID := c.GetInt("userID")
    accountID := c.Param("accountID")
    
    // Verify user has access to this account
    hasAccess, err := h.service.ValidateUserAccountAccess(userID, accountID)
    if err != nil || !hasAccess {
        c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
        return
    }
    
    // Generate new token with account context
    token, err := auth.GenerateJWTWithAccount(userID, accountID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to switch account"})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{"token": token})
}
```

## Backward Compatibility

### 1. API Compatibility
- Maintain existing endpoints where possible
- Add new endpoints for new functionality
- Use versioning for breaking changes

### 2. Data Compatibility
- Preserve existing user-account relationships
- Provide migration scripts
- Support both old and new data structures during transition

### 3. Frontend Compatibility
- Gradual UI updates
- Feature flags for new functionality
- Fallback to old behavior when needed

## Security Considerations

### 1. Permission Inheritance
- Organization permissions don't automatically grant account permissions
- Explicit permission assignment required
- Audit trail for all permission changes

### 2. Account Isolation
- Users can only access accounts they're explicitly assigned to
- Cross-account data access requires explicit permissions
- Proper validation at all API endpoints

### 3. Token Security
- JWT tokens include account context
- Token expiration and refresh mechanisms
- Secure token storage and transmission

## Testing Strategy

### 1. Unit Tests
- Test new permission system
- Test user-account relationships
- Test role validation

### 2. Integration Tests
- Test authentication flows
- Test permission checking
- Test account switching

### 3. End-to-End Tests
- Test complete user workflows
- Test multi-account scenarios
- Test organization management

## Performance Considerations

### 1. Database Optimization
- Proper indexing on relationship tables
- Efficient permission queries
- Caching for frequently accessed permissions

### 2. API Optimization
- Batch permission checking
- Efficient user account queries
- Proper pagination for large datasets

### 3. Frontend Optimization
- Lazy loading of account data
- Efficient permission checking
- Optimistic UI updates

## Future Enhancements

### 1. Advanced Permissions
- Time-based permissions
- Conditional permissions
- Permission templates

### 2. Organization Features
- Organization hierarchies
- Department management
- Cross-organization reporting

### 3. User Management
- Bulk user operations
- User provisioning workflows
- Advanced user profiles

## Conclusion

The new user, organization, and account structure provides a solid foundation for PantryOS's growth and scalability. The flexible permission system and multi-account support will enable the platform to serve a wide range of business scenarios, from single-location operations to large enterprise organizations.

The migration strategy ensures a smooth transition while maintaining backward compatibility, and the comprehensive testing approach will ensure reliability and security throughout the implementation process.
