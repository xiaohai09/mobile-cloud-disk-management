-- ============================================
-- 数据库索引优化（高频查询覆盖）
-- ============================================

-- accounts：按 user_id + is_active 查询账号列表
ALTER TABLE `accounts`
  ADD INDEX `idx_user_active_v2` (`user_id`, `is_active`);

-- exchange_accounts：按 user_id + is_active 查询兑换账号
ALTER TABLE `exchange_accounts`
  ADD INDEX `idx_user_active_v2` (`user_id`, `is_active`);

-- exchange_tasks：按 status + created_at 分页拉取待执行任务
ALTER TABLE `exchange_tasks`
  ADD INDEX `idx_status_created` (`status`, `created_at`);

-- exchange_tasks：按 user_id + status 查询用户任务列表
ALTER TABLE `exchange_tasks`
  ADD INDEX `idx_user_status` (`user_id`, `status`);

-- exchange_records：按 user_id + created_at 查询兑换记录
ALTER TABLE `exchange_records`
  ADD INDEX `idx_user_created` (`user_id`, `created_at`);

-- task_logs：按 user_id + created_at 查询执行日志
ALTER TABLE `task_logs`
  ADD INDEX `idx_user_created` (`user_id`, `created_at`);

-- task_logs：按 account_id + created_at 查询账号日志
ALTER TABLE `task_logs`
  ADD INDEX `idx_account_created` (`account_id`, `created_at`);
