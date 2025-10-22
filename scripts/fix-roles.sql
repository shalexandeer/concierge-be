-- Fix roles and user assignments
USE `concierge_be`;

-- First, ensure we have the default roles
INSERT IGNORE INTO `roles` (`id`, `name`, `description`, `created_at`, `updated_at`) VALUES
    ('role-super-admin-001', 'super_admin', 'Super administrator with full access to all tenants and features', NOW(), NOW()),
    ('role-tenant-admin-001', 'tenant_admin', 'Tenant administrator with management access within their tenant', NOW(), NOW()),
    ('role-user-001', 'user', 'Regular user with view and booking access', NOW(), NOW());

-- Update all users without a role to have the default 'user' role
UPDATE users SET role_id = 'role-user-001' WHERE role_id IS NULL OR role_id = '';

-- Show the results
SELECT 'Migration completed successfully' as status;
SELECT COUNT(*) as total_roles FROM roles;
SELECT COUNT(*) as users_with_roles FROM users WHERE role_id IS NOT NULL AND role_id != '';
SELECT id, username, email, role_id FROM users LIMIT 5;
