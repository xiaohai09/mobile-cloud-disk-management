-- ============================================
-- RBAC 权限体系基础表（P1-3）
-- ============================================

-- 权限定义表
CREATE TABLE IF NOT EXISTS `permissions` (
  `id` INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  `name` VARCHAR(100) NOT NULL COMMENT '权限名称',
  `code` VARCHAR(100) NOT NULL UNIQUE COMMENT '权限标识符，如 user:read',
  `module` VARCHAR(50) NOT NULL COMMENT '所属模块',
  `description` VARCHAR(255) NULL COMMENT '描述',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  INDEX `idx_module` (`module`),
  INDEX `idx_code` (`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='权限定义表';

-- 角色权限映射表
CREATE TABLE IF NOT EXISTS `role_permissions` (
  `id` INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  `role` VARCHAR(20) NOT NULL COMMENT '角色名称',
  `permission_id` INT UNSIGNED NOT NULL COMMENT '权限ID',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  UNIQUE KEY `uk_role_permission` (`role`, `permission_id`),
  INDEX `idx_role` (`role`),
  CONSTRAINT `fk_role_permissions_permission` FOREIGN KEY (`permission_id`) REFERENCES `permissions` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='角色权限映射表';

-- 用户权限表（用户特有权限，覆盖角色权限）
CREATE TABLE IF NOT EXISTS `user_permissions` (
  `id` INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  `user_id` INT UNSIGNED NOT NULL COMMENT '用户ID',
  `permission_id` INT UNSIGNED NOT NULL COMMENT '权限ID',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  UNIQUE KEY `uk_user_permission` (`user_id`, `permission_id`),
  INDEX `idx_user_id` (`user_id`),
  CONSTRAINT `fk_user_permissions_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_user_permissions_permission` FOREIGN KEY (`permission_id`) REFERENCES `permissions` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户权限表';

-- 插入基础权限数据
INSERT IGNORE INTO `permissions` (`code`, `name`, `module`, `description`) VALUES
-- 用户模块
('user:read', '查看用户', 'user', '查看用户基本信息'),
('user:write', '编辑用户', 'user', '编辑用户信息'),
('user:delete', '删除用户', 'user', '删除用户账号'),
('user:role', '管理角色', 'user', '分配用户角色'),
('user:permission', '管理权限', 'user', '分配用户权限'),

-- 账号模块
('account:read', '查看账号', 'account', '查看兑换账号'),
('account:write', '编辑账号', 'account', '编辑兑换账号'),
('account:delete', '删除账号', 'account', '删除兑换账号'),

-- 任务模块
('task:read', '查看任务', 'task', '查看任务配置'),
('task:write', '编辑任务', 'task', '编辑任务配置'),
('task:execute', '执行任务', 'task', '手动执行任务'),

-- 兑换模块
('exchange:read', '查看兑换记录', 'exchange', '查看兑换记录'),
('exchange:write', '编辑兑换记录', 'exchange', '编辑兑换记录'),

-- 导出模块
('export:read', '查看导出', 'export', '查看导出记录'),
('export:write', '管理导出', 'export', '创建/删除导出'),

-- Webhook模块
('webhook:read', '查看Webhook', 'webhook', '查看Webhook配置'),
('webhook:write', '管理Webhook', 'webhook', '创建/编辑/删除Webhook'),

-- 系统模块
('system:config', '系统配置', 'system', '管理系统配置'),
('system:log', '查看日志', 'system', '查看系统日志');

-- 插入角色权限映射
-- admin 角色拥有所有权限
INSERT IGNORE INTO `role_permissions` (`role`, `permission_id`)
SELECT 'admin', id FROM `permissions`;

-- user 角色只有基础查看权限
INSERT IGNORE INTO `role_permissions` (`role`, `permission_id`)
SELECT 'user', id FROM `permissions` WHERE code IN (
  'user:read', 'account:read', 'task:read', 'exchange:read', 'export:read', 'webhook:read'
);
