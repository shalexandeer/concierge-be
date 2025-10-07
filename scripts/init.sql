-- 创建数据库
CREATE DATABASE IF NOT EXISTS `gin_boilerplate` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE `gin_boilerplate`;

-- 用户表
CREATE TABLE IF NOT EXISTS `users` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '用户ID',
    `username` VARCHAR(50) NOT NULL COMMENT '用户名',
    `email` VARCHAR(100) NOT NULL COMMENT '邮箱',
    `password` VARCHAR(255) NOT NULL COMMENT '密码（加密后）',
    `full_name` VARCHAR(100) DEFAULT NULL COMMENT '全名',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` DATETIME DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_username` (`username`),
    UNIQUE KEY `idx_email` (`email`),
    KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';

-- 插入示例数据（可选）
-- INSERT INTO `users` (`username`, `email`, `password`, `full_name`)
-- VALUES ('admin', 'admin@example.com', '$2a$10$xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx', 'Administrator');

-- 租户表
CREATE TABLE IF NOT EXISTS `tenants` (
    `id` VARCHAR(36) NOT NULL COMMENT '租户ID',
    `name` VARCHAR(100) NOT NULL COMMENT '租户名称',
    `description` TEXT COMMENT '描述',
    `domain` VARCHAR(100) COMMENT '域名',
    `is_active` TINYINT(1) DEFAULT 1 COMMENT '是否激活',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` DATETIME DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_domain` (`domain`),
    KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='租户表';

-- 用户租户关联表
CREATE TABLE IF NOT EXISTS `user_tenants` (
    `id` VARCHAR(36) NOT NULL COMMENT 'ID',
    `user_id` VARCHAR(36) NOT NULL COMMENT '用户ID',
    `tenant_id` VARCHAR(36) NOT NULL COMMENT '租户ID',
    `role` VARCHAR(50) DEFAULT 'member' COMMENT '角色',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_tenant_id` (`tenant_id`),
    FOREIGN KEY (`user_id`) REFERENCES `users`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenants`(`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户租户关联表';

-- 设施分类表 (tenant-scoped)
CREATE TABLE IF NOT EXISTS `amenities_categories` (
    `id` VARCHAR(36) NOT NULL COMMENT '分类ID',
    `tenant_id` VARCHAR(36) NOT NULL COMMENT '租户ID',
    `name` VARCHAR(100) NOT NULL COMMENT '分类名称',
    `description` TEXT COMMENT '描述',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` DATETIME DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    KEY `idx_tenant_id` (`tenant_id`),
    KEY `idx_deleted_at` (`deleted_at`),
    FOREIGN KEY (`tenant_id`) REFERENCES `tenants`(`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='设施分类表';

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
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` DATETIME DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    KEY `idx_tenant_id` (`tenant_id`),
    KEY `idx_category_id` (`category_id`),
    KEY `idx_deleted_at` (`deleted_at`),
    FOREIGN KEY (`tenant_id`) REFERENCES `tenants`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`category_id`) REFERENCES `amenities_categories`(`id`) ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='设施表';
