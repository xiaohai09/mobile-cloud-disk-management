# 云飞云盘系统（二创版）

[![CI Status](https://github.com/xiaohai09/caiyun-refactored/actions/workflows/ci.yml/badge.svg)](https://github.com/xiaohai09/caiyun-refactored/actions/workflows/ci.yml)
[![Docker Build](https://github.com/xiaohai09/caiyun-refactored/actions/workflows/publish.yml/badge.svg)](https://github.com/xiaohai09/caiyun-refactored/actions/workflows/publish.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Language](https://img.shields.io/badge/Language-Chinese-blue.svg)](#)
[![English](https://img.shields.io/badge/README-English-green.svg)](README.en.md)

> 基于 [3238614968/caiyun](https://github.com/3238614968/caiyun) 的架构重构与功能扩展版本。

## 📋 文档导航（多语言）

| 语言 | 文档 | 说明 |
|------|------|------|
| 🇨🇳 **中文** | [README.md](README.md) | 项目概览与快速开始（默认） |
| 🇺🇸 English | [README.en.md](README.en.md) | Project overview and quick start |
| 🇨🇳 **中文** | [部署指南](docs/DEPLOYMENT.md) | Docker 部署、环境配置、故障排查 |
| 🇺🇸 English | [Deployment Guide](docs/DEPLOYMENT.en.md) | Docker deployment, config, troubleshooting |
| 🇨🇳 **中文** | [API 文档](docs/API.md) | 认证、导出 API、Webhook API |
| 🇺🇸 English | [API Documentation](docs/API.en.md) | Auth, export API, webhook API |
| 🇨🇳 **中文** | [用户指南](docs/USER_GUIDE.md) | 角色说明、功能模块、FAQ |
| 🇺🇸 English | [User Guide](docs/USER_GUIDE.en.md) | Roles, features, FAQ |
| 🇨🇳 **中文** | [测试报告](docs/TEST_REPORT.md) | 构建、E2E、Docker、安全审计 |
| 🇺🇸 English | [Test Report](docs/TEST_REPORT.en.md) | Build, E2E, Docker, security audit |

## ✨ 核心功能

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

## 🚀 快速开始

```bash
# 1. 克隆仓库
git clone https://github.com/xiaohai09/caiyun-refactored.git
cd caiyun-refactored

# 2. 一键启动（自动生成配置并启动）
chmod +x scripts/start.sh
./scripts/start.sh

# 3. 访问
# 前端：http://localhost
# 监控：http://localhost:3000 (admin / 见 .env)
```

详细部署文档见 [docs/DEPLOYMENT.md](docs/DEPLOYMENT.md)。

## 💻 开发模式

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

## 🛠 技术栈

| 层级 | 技术 |
|------|------|
| 前端 | Vue 3 + TypeScript + Element Plus + Pinia + Vite |
| 后端 | Go + Gin + GORM + Redis + Gorilla WebSocket |
| 数据库 | MySQL 8.0 |
| 任务队列 | Redis Streams / Memory |
| 监控 | Prometheus + Grafana |
| 部署 | Docker Compose + 多阶段构建 |

## ✅ 验证记录

### 前端
- `npm run lint`：通过
- `npm run typecheck`：通过
- `npm run build`：通过

### 后端
- `go build ./cmd/api`：通过
- `go build ./cmd/worker`：通过
- `go test ./...`：通过

### 安全审计
- 新增 SecurityHeadersMiddleware（X-Frame-Options、CSP、HSTS 等）
- 修复 WebSocket `CheckOrigin` 空 origin 宽松校验
- 审计日志增加脱敏过滤
- 所有新增路由挂载 `AuthMiddleware + CSRF + RateLimit + Audit`

## 📝 许可证

MIT
