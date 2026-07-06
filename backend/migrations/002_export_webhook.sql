-- 15. 数据导出与 Webhook 相关表

-- 15.1 导出历史表
CREATE TABLE IF NOT EXISTS `export_history` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `user_id` INT UNSIGNED NOT NULL COMMENT '导出用户ID',
    `export_type` VARCHAR(50) NOT NULL COMMENT '导出类型: accounts/tasks/exchange/records',
    `format` VARCHAR(10) NOT NULL COMMENT '导出格式: csv/json',
    `filters` JSON COMMENT '筛选条件',
    `file_path` VARCHAR(255) COMMENT '文件存储路径',
    `file_size` BIGINT COMMENT '文件大小(字节)',
    `status` VARCHAR(20) NOT NULL DEFAULT 'pending' COMMENT '状态: pending/success/failed',
    `error_msg` VARCHAR(500) COMMENT '错误信息',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    INDEX `idx_export_history_user` (`user_id`),
    INDEX `idx_export_history_type` (`export_type`),
    INDEX `idx_export_history_status` (`status`),
    INDEX `idx_export_history_created` (`created_at`),
    CONSTRAINT `fk_export_history_user` FOREIGN KEY (`user_id`) REFERENCES `users`(`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 15.2 Webhook 端点表
CREATE TABLE IF NOT EXISTS `webhook_endpoints` (
    `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
    `user_id` INT UNSIGNED NOT NULL COMMENT '所属用户ID',
    `name` VARCHAR(100) NOT NULL COMMENT '端点名称',
    `url` VARCHAR(500) NOT NULL COMMENT 'Webhook URL',
    `events` JSON NOT NULL COMMENT '订阅事件列表',
    `secret` VARCHAR(100) COMMENT '签名密钥',
    `headers` JSON COMMENT '自定义请求头',
    `is_active` TINYINT(1) NOT NULL DEFAULT 1 COMMENT '是否启用',
    `last_triggered_at` DATETIME COMMENT '最后触发时间',
    `fail_count` INT NOT NULL DEFAULT 0 COMMENT '连续失败次数',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    INDEX `idx_webhook_endpoints_user` (`user_id`),
    INDEX `idx_webhook_endpoints_active` (`is_active`),
    CONSTRAINT `fk_webhook_endpoints_user` FOREIGN KEY (`user_id`) REFERENCES `users`(`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 15.3 Webhook 投递日志表
CREATE TABLE IF NOT EXISTS `webhook_deliveries` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `endpoint_id` INT UNSIGNED NOT NULL COMMENT '端点ID',
    `user_id` INT UNSIGNED NOT NULL COMMENT '用户ID',
    `event_type` VARCHAR(50) NOT NULL COMMENT '事件类型',
    `payload` JSON COMMENT '请求体',
    `status_code` INT COMMENT 'HTTP状态码',
    `response_body` TEXT COMMENT '响应体',
    `error_msg` VARCHAR(500) COMMENT '错误信息',
    `duration_ms` INT COMMENT '耗时(毫秒)',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    INDEX `idx_webhook_deliveries_endpoint` (`endpoint_id`),
    INDEX `idx_webhook_deliveries_user` (`user_id`),
    INDEX `idx_webhook_deliveries_created` (`created_at`),
    CONSTRAINT `fk_webhook_deliveries_endpoint` FOREIGN KEY (`endpoint_id`) REFERENCES `webhook_endpoints`(`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_webhook_deliveries_user` FOREIGN KEY (`user_id`) REFERENCES `users`(`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
