# Scripts

This directory contains utility scripts for PantryOS.

## Migration Scripts

### `migrate_user_structure.sql`

Database migration script to transition from the old user-account structure to the new flexible structure.

**What it does:**
- Migrates existing user-account relationships to the new `UserAccount` table
- Creates default permissions for existing roles
- Sets up organization relationships where needed
- Adds performance indexes
- Creates data integrity constraints
- Provides verification queries

**Usage:**
```bash
# Run the migration script
mysql -u username -p database_name < scripts/migrate_user_structure.sql
```

**Important Notes:**
- Run this script after the new tables are created but before the application is updated
- Always backup your database before running migrations
- Test in a staging environment first
- The script includes verification queries to ensure successful migration

**Verification:**
After running the migration, execute the verification queries at the end of the script to ensure:
- All users have at least one account relationship
- Users have exactly one primary account
- No duplicate primary accounts exist

For more details, see the main documentation in `USER_ACCOUNT_STRUCTURE.md`.
