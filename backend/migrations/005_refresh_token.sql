-- ============================================
-- 双 Token 认证增强（P1-4）
-- ============================================

-- 为用户表增加 refresh token 黑名单字段（JSON 哈希列表）
CALL AddColumnIfNotExists('users', 'refresh_sessions', 'TEXT NULL COMMENT ''refresh token 哈希黑名单（JSON 数组）'' AFTER `token_version`');

-- 可选：为 refresh token 创建独立存储表（若需更精细控制）
-- 当前采用 user.refresh_sessions JSON 字段存储，简化部署
