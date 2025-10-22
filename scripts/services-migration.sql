-- Services Migration Script
-- Creates services_categories and services tables

-- Create services_categories table
CREATE TABLE IF NOT EXISTS services_categories (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    icon VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_services_categories_deleted_at (deleted_at)
);

-- Create services table
CREATE TABLE IF NOT EXISTS services (
    id VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(36) NOT NULL,
    category_id VARCHAR(36) NOT NULL,
    service_name VARCHAR(100) NOT NULL,
    description TEXT,
    operating_hours_from TIME NULL,
    operating_hours_to TIME NULL,
    available24_7 BOOLEAN DEFAULT FALSE,
    response_time INT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_services_tenant_id (tenant_id),
    INDEX idx_services_category_id (category_id),
    INDEX idx_services_deleted_at (deleted_at),
    INDEX idx_services_tenant_category (tenant_id, category_id),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    FOREIGN KEY (category_id) REFERENCES services_categories(id) ON DELETE RESTRICT
);

-- Insert some default service categories
INSERT INTO services_categories (id, name, description, icon) VALUES
('sc-001', 'Transportation', 'Transportation and travel services', 'car'),
('sc-002', 'Concierge', 'Personal concierge services', 'bell'),
('sc-003', 'Maintenance', 'Property maintenance and repair services', 'wrench'),
('sc-004', 'Cleaning', 'Cleaning and housekeeping services', 'sparkles'),
('sc-005', 'Security', 'Security and safety services', 'shield'),
('sc-006', 'Entertainment', 'Entertainment and leisure services', 'music'),
('sc-007', 'Food & Beverage', 'Food and beverage services', 'utensils'),
('sc-008', 'Health & Wellness', 'Health and wellness services', 'heart'),
('sc-009', 'Business', 'Business and professional services', 'briefcase'),
('sc-010', 'Emergency', 'Emergency and urgent services', 'phone')
ON DUPLICATE KEY UPDATE name = VALUES(name);
