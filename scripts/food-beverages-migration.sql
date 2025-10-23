-- Food & Beverages Migration Script
-- Creates tables for food and beverage management system

-- Create food_beverages_categories table (global, not tenant-scoped)
CREATE TABLE IF NOT EXISTS food_beverages_categories (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    icon VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_deleted_at (deleted_at)
);

-- Create food_beverages table (tenant-scoped)
CREATE TABLE IF NOT EXISTS food_beverages (
    id VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(36) NOT NULL,
    category_id VARCHAR(36) NOT NULL,
    item_name VARCHAR(100) NOT NULL,
    description TEXT,
    price DECIMAL(10,2) NOT NULL,
    preparation_time INT NOT NULL COMMENT 'Preparation time in minutes',
    service_hours_start TIME NULL,
    service_hours_end TIME NULL,
    all_day BOOLEAN DEFAULT FALSE,
    image_path VARCHAR(500),
    is_available BOOLEAN DEFAULT TRUE COMMENT 'Available for ordering',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_category_id (category_id),
    INDEX idx_deleted_at (deleted_at),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    FOREIGN KEY (category_id) REFERENCES food_beverages_categories(id) ON DELETE RESTRICT
);

-- Insert sample categories
INSERT INTO food_beverages_categories (id, name, description, icon) VALUES
    (UUID(), 'Appetizers', 'Starters and small plates', '🥗'),
    (UUID(), 'Main Course', 'Main dishes and entrees', '🍽️'),
    (UUID(), 'Desserts', 'Sweet treats and desserts', '🍰'),
    (UUID(), 'Beverages', 'Drinks and refreshments', '🥤'),
    (UUID(), 'Salads', 'Fresh salads and healthy options', '🥙'),
    (UUID(), 'Soups', 'Hot soups and broths', '🍲'),
    (UUID(), 'Breakfast', 'Breakfast items and brunch', '🍳'),
    (UUID(), 'Snacks', 'Light snacks and finger foods', '🍿')
ON DUPLICATE KEY UPDATE name=name;

