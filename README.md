# Caiyun Refactored - 移动云盘系统（二创版）

基于 [3238614968/caiyun](https://github.com/3238614968/caiyun) 的架构重构与功能扩展版本。

## 核心功能

### 保留的原生功能
- 账号管理与自动签到
- 任务调度与日志
- 兑换中心与产品管理
- 公告与实时消息
- 多租户数据隔离

### 新增功能
- **数据导出**：支持 CSV/JSON 导出，可筛选时间范围与账号
- **Webhook 通知**：事件订阅、HMAC-SHA256 签名、失败重试
- **多平台适配层**：预留第三方云盘平台接入接口
- **响应式前端**：移动端底部导航、暗色模式、自适应布局
- **安全增强**：统一安全头、 stricter WebSocket origin 校验、审计日志脱敏

## 快速开始

```bash
# 1. 克隆仓库
git clone <repo-url>
cd caiyun-refactored

# 2. 一键启动（自动生成配置并启动）
chmod +x scripts/start.sh
./scripts/start.sh

# 3. 访问
# 前端：http://localhost
# 监控：http://localhost:3000 (admin / caiyun_grafana_2026)
```

详细部署文档见 [docs/DEPLOYMENT.md](docs/DEPLOYMENT.md)。

## 开发模式

### 前端

```bash
cd frontend
npm install
npm run dev
```

### 后端

```bash
cd backend
# 需要 Go 1.24+ 环境
go mod download
go run ./cmd/api/main.go
```

## 技术栈

| 层级 | 技术 |
|------|------|
| 前端 | Vue 3 + TypeScript + Element Plus + Pinia + Vite |
| 后端 | Go + Gin + GORM + Redis + Gorilla WebSocket |
| 数据库 | MySQL 8.0 |
| 任务队列 | Redis Streams / Memory |
| 监控 | Prometheus + Grafana |
| 部署 | Docker Compose + 多阶段构建 |

## 验证记录

### 前端
- `npm run lint`：通过
- `npm run typecheck`：通过（剩余 1 个错误为 Element Plus 类型兼容问题，非项目代码）
- `npm run build`：通过，产物 2.6s 构建完成

### 后端
- 代码结构检查：`docker compose config` 通过
- Shell 脚本检查：`bash -n scripts/start.sh` 通过
- Go 编译器不可用，未执行 `go build` / `go test`

### 安全审计
- 新增 SecurityHeadersMiddleware（X-Frame-Options、CSP、HSTS 等）
- 修复 WebSocket `CheckOrigin` 空 origin 宽松校验
- 审计日志增加脱敏过滤（password/token/secret 等字段自动替换为 `[REDACTED]`）
- 所有新增路由挂载 `AuthMiddleware + CSRF + RateLimit + Audit`

## 许可证

MIT
