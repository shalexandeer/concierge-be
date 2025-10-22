-- Role Migration Script
-- This script handles role assignment and foreign key constraint

USE `concierge_be`;

-- First, ensure we have the default roles
INSERT IGNORE INTO `roles` (`id`, `name`, `description`, `created_at`, `updated_at`) VALUES
    (UUID(), 'super_admin', 'Super administrator with full access to all tenants and features', NOW(), NOW()),
    (UUID(), 'tenant_admin', 'Tenant administrator with management access within their tenant', NOW(), NOW()),
    (UUID(), 'user', 'Regular user with view and booking access', NOW(), NOW());

-- Get the user role ID
SET @user_role_id = (SELECT id FROM roles WHERE name = 'user' LIMIT 1);

-- Update all users without a role to have the default 'user' role
UPDATE users SET role_id = @user_role_id WHERE role_id IS NULL OR role_id = '';

-- Now add the foreign key constraint
SET @constraint_exists = (
    SELECT COUNT(*) 
    FROM information_schema.table_constraints 
    WHERE constraint_schema = 'concierge_be' 
    AND table_name = 'users' 
    AND constraint_name = 'fk_users_role'
);

SET @sql = IF(@constraint_exists = 0, 
    'ALTER TABLE `users` ADD CONSTRAINT `fk_users_role` FOREIGN KEY (`role_id`) REFERENCES `roles`(`id`) ON UPDATE CASCADE ON DELETE SET NULL',
    'SELECT "Foreign key constraint already exists" as message'
);

PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Show the results
SELECT 'Migration completed successfully' as status;
SELECT COUNT(*) as total_roles FROM roles;
SELECT COUNT(*) as users_with_roles FROM users WHERE role_id IS NOT NULL AND role_id != '';
