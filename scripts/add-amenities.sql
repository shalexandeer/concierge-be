-- Migration to add amenities and amenities_categories tables

-- 设施分类表 (tenant-scoped)
CREATE TABLE IF NOT EXISTS `amenities_categories` (
    `id` VARCHAR(36) NOT NULL COMMENT '分类ID',
    `tenant_id` VARCHAR(36) NOT NULL COMMENT '租户ID',
    `name` VARCHAR(100) NOT NULL COMMENT '分类名称',
    `description` TEXT COMMENT '描述',
    `created_at` DATETIME(3) DEFAULT NULL COMMENT '创建时间',
    `updated_at` DATETIME(3) DEFAULT NULL COMMENT '更新时间',
    `deleted_at` DATETIME(3) DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    KEY `idx_tenant_id` (`tenant_id`),
    KEY `idx_deleted_at` (`deleted_at`),
    CONSTRAINT `fk_amenities_categories_tenant` FOREIGN KEY (`tenant_id`) REFERENCES `tenants`(`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='设施分类表';

-- 设施表 (tenant-scoped)
CREATE TABLE IF NOT EXISTS `amenities` (
    `id` VARCHAR(36) NOT NULL COMMENT '设施ID',
    `tenant_id` VARCHAR(36) NOT NULL COMMENT '租户ID',
    `category_id` VARCHAR(36) NOT NULL COMMENT '分类ID',
    `item_name` VARCHAR(100) NOT NULL COMMENT '设施名称',
    `description` TEXT COMMENT '描述',
    `stock` INT NOT NULL DEFAULT 0 COMMENT '库存数量',
    `minimum_stock` INT NOT NULL DEFAULT 0 COMMENT '最小库存',
    `available` TINYINT(1) DEFAULT 1 COMMENT '是否可用',
    `created_at` DATETIME(3) DEFAULT NULL COMMENT '创建时间',
    `updated_at` DATETIME(3) DEFAULT NULL COMMENT '更新时间',
    `deleted_at` DATETIME(3) DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    KEY `idx_tenant_id` (`tenant_id`),
    KEY `idx_category_id` (`category_id`),
    KEY `idx_deleted_at` (`deleted_at`),
    CONSTRAINT `fk_amenities_tenant` FOREIGN KEY (`tenant_id`) REFERENCES `tenants`(`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_amenities_category` FOREIGN KEY (`category_id`) REFERENCES `amenities_categories`(`id`) ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='设施表';
