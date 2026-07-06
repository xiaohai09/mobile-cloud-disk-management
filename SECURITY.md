# 安全策略

## 报告安全问题

我们非常重视安全性。如果您发现任何安全漏洞，请**不要**在公开 Issue 中报告。

### 报告流程

1. 发送邮件至：**security@nousresearch.com**（替换为实际联系邮箱）
2. 或在 GitHub 上创建 **Private Security Advisory**：
   - 访问 https://github.com/xiaohai09/mobile-cloud-disk-management/security/advisories
   - 点击 "New draft security advisory"
   - 详细描述漏洞、影响范围和修复建议

### 响应时间

- 我们将在 **3 个工作日内** 确认收到您的报告
- 漏洞修复将在 **30 天内** 发布补丁版本
- 修复后将在 Release Notes 中致谢（如您同意）

### 安全修复流程

1. 私下确认漏洞
2. 开发修复补丁
3. 内部测试验证
4. 发布新版本
5. 公开披露（在修复发布后）

## 安全最佳实践

### 部署安全
- 使用强密码和随机生成的 JWT Secret
- 启用 HTTPS/TLS（Let's Encrypt + Caddy/Nginx）
- 限制数据库和 Redis 访问权限
- 定期更新依赖包

### 开发安全
- 所有新增代码需通过安全扫描（Trivy）
- 禁止硬编码密钥、密码、Token
- 使用参数化查询防止 SQL 注入
- 所有路由需挂载认证中间件

## 已知安全机制

| 机制 | 说明 |
|------|------|
| JWT 认证 | Access/Refresh Token 分离，短期 Access Token |
| CSRF 保护 | 所有表单请求携带 CSRF Token |
| 限流 | pre-auth + post-auth 双层限流 |
| 审计日志 | 记录所有敏感操作，密码/token 自动脱敏 |
| 安全头 | X-Frame-Options、CSP、HSTS 等 |
| WebSocket Origin | 严格 origin 校验，拒绝空 origin |
| 非 root 运行 | 容器内使用 app 用户 |

## 依赖安全

- 使用 Dependabot 自动更新依赖
- CI/CD 中集成 Trivy 安全扫描
- 定期运行 `npm audit` 和 `go mod tidy`

## 联系

- **GitHub Security Advisory**: https://github.com/xiaohai09/mobile-cloud-disk-management/security/advisories
- **Email**: security@nousresearch.com（替换为实际邮箱）
