-- Migration script to transition from old user-account structure to new flexible structure
-- This script should be run after the new tables are created but before the application is updated

-- Step 1: Create new tables (if not already created by GORM)
-- Note: These tables should already exist from the model updates

-- Step 2: Migrate existing user-account relationships to UserAccount table
-- This preserves the existing relationships while enabling the new structure

INSERT INTO user_accounts (user_id, account_id, role, is_primary, status, created_at, updated_at)
SELECT 
    u.id as user_id,
    u.account_id,
    COALESCE(u.role, 'user') as role,
    true as is_primary, -- Mark existing account as primary
    'active' as status,
    u.created_at,
    NOW() as updated_at
FROM users u
WHERE u.account_id IS NOT NULL;

-- Step 3: Update existing users to have first_name and last_name
-- Extract from email if not available (temporary solution)
UPDATE users 
SET 
    first_name = COALESCE(
        SUBSTRING_INDEX(email, '@', 1), 
        'User'
    ),
    last_name = COALESCE(
        SUBSTRING_INDEX(SUBSTRING_INDEX(email, '@', 1), '.', -1),
        'User'
    ),
    status = 'active',
    updated_at = NOW()
WHERE first_name IS NULL OR first_name = '';

-- Step 4: Create default permissions for existing roles
-- These permissions will be used by the new permission system

-- Account-level permissions
INSERT INTO permissions (name, description, resource, action, created_at) VALUES
('inventory.read', 'Read inventory items', 'inventory', 'read', NOW()),
('inventory.write', 'Create and update inventory items', 'inventory', 'write', NOW()),
('inventory.delete', 'Delete inventory items', 'inventory', 'delete', NOW()),
('categories.read', 'Read categories', 'categories', 'read', NOW()),
('categories.write', 'Create and update categories', 'categories', 'write', NOW()),
('categories.delete', 'Delete categories', 'categories', 'delete', NOW()),
('orders.read', 'Read orders', 'orders', 'read', NOW()),
('orders.write', 'Create and update orders', 'orders', 'write', NOW()),
('orders.delete', 'Delete orders', 'orders', 'delete', NOW()),
('users.read', 'Read users in account', 'users', 'read', NOW()),
('users.write', 'Create and update users in account', 'users', 'write', NOW()),
('users.delete', 'Delete users in account', 'users', 'delete', NOW()),
('reports.read', 'Read reports', 'reports', 'read', NOW()),
('settings.read', 'Read account settings', 'settings', 'read', NOW()),
('settings.write', 'Update account settings', 'settings', 'write', NOW());

-- Organization-level permissions
INSERT INTO permissions (name, description, resource, action, created_at) VALUES
('org.inventory.read', 'Read inventory across organization', 'org_inventory', 'read', NOW()),
('org.inventory.write', 'Manage inventory across organization', 'org_inventory', 'write', NOW()),
('org.users.read', 'Read users across organization', 'org_users', 'read', NOW()),
('org.users.write', 'Manage users across organization', 'org_users', 'write', NOW()),
('org.reports.read', 'Read organization reports', 'org_reports', 'read', NOW()),
('org.settings.read', 'Read organization settings', 'org_settings', 'read', NOW()),
('org.settings.write', 'Update organization settings', 'org_settings', 'write', NOW());

-- Step 5: Create role-permission mappings
-- User role permissions (account scope)
INSERT INTO role_permissions (role, permission_id, scope, created_at)
SELECT 'user', id, 'account', NOW()
FROM permissions 
WHERE name IN (
    'inventory.read',
    'categories.read',
    'orders.read',
    'reports.read'
);

-- Manager role permissions (account scope)
INSERT INTO role_permissions (role, permission_id, scope, created_at)
SELECT 'manager', id, 'account', NOW()
FROM permissions 
WHERE name IN (
    'inventory.read',
    'inventory.write',
    'categories.read',
    'categories.write',
    'orders.read',
    'orders.write',
    'users.read',
    'reports.read',
    'settings.read'
);

-- Admin role permissions (account scope)
INSERT INTO role_permissions (role, permission_id, scope, created_at)
SELECT 'admin', id, 'account', NOW()
FROM permissions 
WHERE name IN (
    'inventory.read',
    'inventory.write',
    'inventory.delete',
    'categories.read',
    'categories.write',
    'categories.delete',
    'orders.read',
    'orders.write',
    'orders.delete',
    'users.read',
    'users.write',
    'users.delete',
    'reports.read',
    'settings.read',
    'settings.write'
);

-- Organization member permissions
INSERT INTO role_permissions (role, permission_id, scope, created_at)
SELECT 'member', id, 'organization', NOW()
FROM permissions 
WHERE name IN (
    'org.inventory.read',
    'org.reports.read'
);

-- Organization admin permissions
INSERT INTO role_permissions (role, permission_id, scope, created_at)
SELECT 'org_admin', id, 'organization', NOW()
FROM permissions 
WHERE name IN (
    'org.inventory.read',
    'org.inventory.write',
    'org.users.read',
    'org.users.write',
    'org.reports.read',
    'org.settings.read',
    'org.settings.write'
);

-- Organization owner permissions
INSERT INTO role_permissions (role, permission_id, scope, created_at)
SELECT 'org_owner', id, 'organization', NOW()
FROM permissions 
WHERE name IN (
    'org.inventory.read',
    'org.inventory.write',
    'org.users.read',
    'org.users.write',
    'org.reports.read',
    'org.settings.read',
    'org.settings.write'
);

-- Step 6: Create UserOrganization relationships for existing users
-- This assumes users in accounts that belong to organizations should have org membership
INSERT INTO user_organizations (user_id, organization_id, role, status, created_at, updated_at)
SELECT DISTINCT
    ua.user_id,
    a.organization_id,
    'member' as role,
    'active' as status,
    NOW() as created_at,
    NOW() as updated_at
FROM user_accounts ua
JOIN accounts a ON ua.account_id = a.id
WHERE a.organization_id IS NOT NULL
AND NOT EXISTS (
    SELECT 1 FROM user_organizations uo 
    WHERE uo.user_id = ua.user_id 
    AND uo.organization_id = a.organization_id
);

-- Step 7: Update organization status to active if not set
UPDATE organizations 
SET status = 'active', updated_at = NOW()
WHERE status IS NULL OR status = '';

-- Step 8: Create indexes for performance
-- These indexes will improve query performance for the new structure

-- UserAccount indexes
CREATE INDEX IF NOT EXISTS idx_user_accounts_user_id ON user_accounts(user_id);
CREATE INDEX IF NOT EXISTS idx_user_accounts_account_id ON user_accounts(account_id);
CREATE INDEX IF NOT EXISTS idx_user_accounts_user_account ON user_accounts(user_id, account_id);
CREATE INDEX IF NOT EXISTS idx_user_accounts_status ON user_accounts(status);

-- UserOrganization indexes
CREATE INDEX IF NOT EXISTS idx_user_organizations_user_id ON user_organizations(user_id);
CREATE INDEX IF NOT EXISTS idx_user_organizations_org_id ON user_organizations(organization_id);
CREATE INDEX IF NOT EXISTS idx_user_organizations_user_org ON user_organizations(user_id, organization_id);
CREATE INDEX IF NOT EXISTS idx_user_organizations_status ON user_organizations(status);

-- Permission indexes
CREATE INDEX IF NOT EXISTS idx_permissions_name ON permissions(name);
CREATE INDEX IF NOT EXISTS idx_permissions_resource ON permissions(resource);
CREATE INDEX IF NOT EXISTS idx_permissions_action ON permissions(action);

-- RolePermission indexes
CREATE INDEX IF NOT EXISTS idx_role_permissions_role ON role_permissions(role);
CREATE INDEX IF NOT EXISTS idx_role_permissions_permission_id ON role_permissions(permission_id);
CREATE INDEX IF NOT EXISTS idx_role_permissions_scope ON role_permissions(scope);

-- Step 9: Add constraints to ensure data integrity
-- Note: These constraints help maintain data consistency

-- Ensure users can't have duplicate primary accounts
ALTER TABLE user_accounts 
ADD CONSTRAINT unique_primary_account_per_user 
UNIQUE (user_id, is_primary) 
WHERE is_primary = true;

-- Ensure users can't have duplicate roles in the same account
ALTER TABLE user_accounts 
ADD CONSTRAINT unique_user_account_role 
UNIQUE (user_id, account_id, role);

-- Ensure users can't have duplicate roles in the same organization
ALTER TABLE user_organizations 
ADD CONSTRAINT unique_user_organization_role 
UNIQUE (user_id, organization_id, role);

-- Step 10: Create a view for backward compatibility
-- This view allows existing code to continue working during transition
CREATE OR REPLACE VIEW user_account_view AS
SELECT 
    u.id,
    u.email,
    u.first_name,
    u.last_name,
    u.is_verified,
    u.status,
    u.created_at,
    ua.account_id,
    ua.role,
    ua.is_primary
FROM users u
LEFT JOIN user_accounts ua ON u.id = ua.user_id AND ua.is_primary = true;

-- Step 11: Add audit trail for migration
-- This helps track what was migrated and when
CREATE TABLE IF NOT EXISTS migration_log (
    id INT AUTO_INCREMENT PRIMARY KEY,
    migration_name VARCHAR(255) NOT NULL,
    description TEXT,
    records_migrated INT DEFAULT 0,
    executed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO migration_log (migration_name, description, records_migrated) VALUES
('user_structure_migration', 'Migrated user-account relationships to new flexible structure', 
 (SELECT COUNT(*) FROM user_accounts));

-- Step 12: Verification queries
-- These queries can be run to verify the migration was successful

-- Check that all users have at least one account relationship
SELECT 
    'Users without account relationships' as check_type,
    COUNT(*) as count
FROM users u
WHERE NOT EXISTS (SELECT 1 FROM user_accounts ua WHERE ua.user_id = u.id)

UNION ALL

-- Check that all users have exactly one primary account
SELECT 
    'Users with multiple primary accounts' as check_type,
    COUNT(*) as count
FROM (
    SELECT user_id, COUNT(*) as primary_count
    FROM user_accounts 
    WHERE is_primary = true
    GROUP BY user_id
    HAVING primary_count > 1
) as multiple_primary

UNION ALL

-- Check that all users have exactly one primary account
SELECT 
    'Users without primary account' as check_type,
    COUNT(*) as count
FROM users u
WHERE NOT EXISTS (
    SELECT 1 FROM user_accounts ua 
    WHERE ua.user_id = u.id AND ua.is_primary = true
);

-- Step 13: Cleanup (optional - run after confirming migration is successful)
-- Remove the old account_id column from users table
-- ALTER TABLE users DROP COLUMN account_id;
-- ALTER TABLE users DROP COLUMN role;

-- Note: The cleanup step should only be run after:
-- 1. All application code has been updated to use the new structure
-- 2. All tests are passing
-- 3. The migration has been verified in production-like environment
-- 4. A backup has been created
