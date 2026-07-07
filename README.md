# 移动云盘管理系统

[![CI Status](https://github.com/xiaohai09/mobile-cloud-disk-management/actions/workflows/ci.yml/badge.svg)](https://github.com/xiaohai09/mobile-cloud-disk-management/actions/workflows/ci.yml/badge)
[![Docker Build](https://github.com/xiaohai09/mobile-cloud-disk-management/actions/workflows/publish.yml/badge.svg)](https://github.com/xiaohai09/mobile-cloud-disk-management/actions/workflows/publish.yml/badge)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.25-blue.svg)](https://go.dev/dl/)
[![Node Version](https://img.shields.io/badge/Node-20-green.svg)](https://nodejs.org/)
[![Docker](https://img.shields.io/badge/Docker-Ready-blue.svg)](https://www.docker.com/)
[![Language](https://img.shields.io/badge/Language-Chinese-red.svg)](#)
[![English](https://img.shields.io/badge/README-English-green.svg)](README.en.md)

> **企业级移动云盘管理平台** — 基于 [3238614968/caiyun](https://github.com/3238614968/caiyun) 的架构重构与功能扩展版本，提供数据导出、Webhook 通知、响应式前端和安全增强等企业特性。

---

## 📑 文档导航（多语言）

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

---

## 🎯 项目概述

### 什么是移动云盘管理系统？

移动云盘管理系统是一个**企业级多云盘账号管理与自动化运营平台**，用于统一管理多个云盘服务账号，提供自动化签到、任务调度、数据导出、Webhook 通知等核心能力。系统采用前后端分离架构，支持 Docker 容器化部署，适用于企业内部云盘资源管理、数据分析和自动化运营场景。

### 核心价值

- **统一管理**：集中管理多个云盘服务账号，支持多租户数据隔离
- **自动化运营**：定时任务、自动签到、兑换中心，降低人工运维成本
- **数据可观测**：支持 CSV/JSON 数据导出，Webhook 事件推送，便于集成到企业数据流
- **安全合规**：统一安全头、WebSocket 源校验、审计日志脱敏、速率限制
- **易于部署**：Docker Compose 一键部署，支持多阶段构建，镜像推送至 GHCR

### 适用场景

- 企业内部云盘账号集中管理与运营
- 自动化签到与任务调度
- 数据导出与第三方系统集成
- 监控告警与事件驱动架构

---

## ✨ 核心功能

### 📊 账号管理
- 多账号统一管理，支持批量导入与分组
- 自动签到与任务调度
- 账号状态实时监控
- 多租户数据隔离，保障数据安全

### 🔄 任务与自动化
- 定时任务调度（支持 Cron 表达式）
- 任务日志与执行记录
- 失败自动重试机制
- 任务状态实时通知

### 🎁 兑换中心
- 兑换码管理与分发
- 兑换记录追踪
- 产品管理（添加、编辑、删除）
- 兑换统计与报表

### 📢 公告与通知
- 系统公告发布与管理
- 实时消息推送（WebSocket）
- 公告阅读状态追踪

### 📤 数据导出
- 支持 CSV/JSON 格式导出
- 可筛选时间范围与账号
- 导出任务队列化处理
- 导出文件生命周期管理

### 🔔 Webhook 通知
- 事件订阅与推送
- HMAC-SHA256 签名验证
- 失败自动重试（指数退避）
- 支持多种事件类型

### 📱 响应式前端
- Vue 3 + TypeScript 现代化技术栈
- Element Plus 企业级 UI 组件库
- 移动端底部导航，自适应布局
- 暗色模式支持

### 🔒 安全增强
- 统一安全头（CSP、HSTS、X-Frame-Options）
- WebSocket 严格源校验
- 审计日志脱敏过滤
- 统一认证、CSRF、速率限制中间件

---

## 🏗️ 系统架构

```
┌─────────────────────────────────────────────────────────────┐
│                        前端 (Vue 3 SPA)                      │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐    │
│  │ 账号管理 │  │ 任务中心 │  │ 兑换中心 │  │ 数据导出 │    │
│  └──────────┘  └──────────┘  └──────────┘  └──────────┘    │
└─────────────────────────────────────────────────────────────┘
                              │ HTTPS / WSS
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                    API Gateway (Go/Gin)                      │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────────┐  │
│  │ Auth Middleware│  │ Rate Limiter │  │ Security Headers │  │
│  └──────────────┘  └──────────────┘  └──────────────────┘  │
└─────────────────────────────────────────────────────────────┘
                              │
          ┌───────────────────┼───────────────────┐
          ▼                   ▼                   ▼
┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐
│   Worker Service │  │  Export Service │  │  Webhook Service │
│   (任务调度)     │  │  (数据导出)     │  │  (事件通知)     │
└─────────────────┘  └─────────────────┘  └─────────────────┘
          │                   │                   │
          └───────────────────┼───────────────────┘
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                    数据层                                    │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐    │
│  │  MySQL   │  │  Redis   │  │  日志    │  │  审计    │    │
│  │ (主存储) │  │ (缓存)  │  │ (Log)   │  │ (Audit) │    │
│  └──────────┘  └──────────┘  └──────────┘  └──────────┘    │
└─────────────────────────────────────────────────────────────┘
```

### 架构特点

- **前后端分离**：Vue 3 SPA + Go 后端 API，职责清晰
- **微服务化**：API、Worker、Export、Webhook 独立服务
- **容器化部署**：Docker 多阶段构建，镜像推送到 GHCR
- **可观测性**：Prometheus 指标 + Grafana 仪表板
- **安全性**：多层安全防护，从传输到存储的全链路保护

---

## 🚀 快速开始

### 前置要求

- Docker Engine 20.10+
- Docker Compose 2.0+
- 4GB+ 可用内存
- 10GB+ 可用磁盘空间

### 一键部署

```bash
# 1. 克隆仓库
git clone https://github.com/xiaohai09/mobile-cloud-disk-management.git
cd mobile-cloud-disk-management

# 2. 复制环境配置
cp .env.example .env
# 编辑 .env 文件，配置数据库密码、JWT 密钥等

# 3. 启动服务
docker compose up -d

# 4. 验证服务
curl http://localhost/api/health
```

### 访问地址

| 服务 | 地址 | 说明 |
|------|------|------|
| 前端应用 | http://localhost | 用户界面 |
| API 服务 | http://localhost/api | RESTful API |
| 监控面板 | http://localhost:3000 | Prometheus + Grafana |
| 健康检查 | http://localhost/api/health | 健康状态端点 |

详细部署文档见 [部署指南](docs/DEPLOYMENT.md)。

---

## 🛠️ 技术栈

### 前端

| 技术 | 用途 | 版本 |
|------|------|------|
| Vue 3 | 渐进式 JavaScript 框架 | ^3.4.0 |
| TypeScript | 类型安全的 JavaScript | ^5.3.0 |
| Element Plus | 企业级 UI 组件库 | ^2.4.0 |
| Pinia | 状态管理 | ^2.1.0 |
| Vite | 下一代前端构建工具 | ^5.0.0 |
| Axios | HTTP 客户端 | ^1.6.0 |

### 后端

| 技术 | 用途 | 版本 |
|------|------|------|
| Go | 编程语言 | 1.25 |
| Gin | HTTP Web 框架 | ^1.10.0 |
| GORM | ORM 库 | ^1.25.0 |
| Redis | 缓存与会话存储 | 7.x |
| Gorilla WebSocket | WebSocket 支持 | ^1.5.0 |
| Prometheus | 指标采集 | ^0.47.0 |

### 数据层

| 技术 | 用途 | 版本 |
|------|------|------|
| MySQL | 主数据库 | 8.0 |
| Redis | 缓存、队列、会话 | 7.x |

### 基础设施

| 技术 | 用途 |
|------|------|
| Docker | 容器化 |
| Docker Compose | 多容器编排 |
| GitHub Actions | CI/CD |
| GHCR | 镜像仓库 |
| Nginx | 反向代理（前端镜像内置） |

---

## 📁 项目结构

```
mobile-cloud-disk-management/
├── backend/                    # 后端服务
│   ├── cmd/                    # 可执行文件入口
│   │   ├── api/                # API 服务入口
│   │   └── worker/             # Worker 服务入口
│   ├── internal/               # 内部包
│   │   ├── handler/            # HTTP 处理器
│   │   ├── service/            # 业务逻辑层
│   │   ├── repository/         # 数据访问层
│   │   ├── model/              # 数据模型
│   │   ├── middleware/         # 中间件
│   │   ├── config/             # 配置管理
│   │   └── utils/              # 工具函数
│   ├── pkg/                    # 公共包
│   │   ├── auth/               # 认证相关
│   │   ├── export/             # 导出功能
│   │   ├── webhook/            # Webhook 功能
│   │   └── security/           # 安全工具
│   ├── migrations/             # 数据库迁移
│   ├── scripts/                # 脚本工具
│   ├── Dockerfile.api          # API 服务 Dockerfile
│   ├── Dockerfile.worker       # Worker 服务 Dockerfile
│   └── go.mod                  # Go 模块定义
├── frontend/                   # 前端应用
│   ├── src/
│   │   ├── views/              # 页面组件
│   │   ├── components/         # 公共组件
│   │   ├── stores/             # Pinia 状态管理
│   │   ├── router/             # 路由配置
│   │   ├── api/                # API 请求封装
│   │   ├── utils/              # 工具函数
│   │   └── assets/             # 静态资源
│   ├── public/                 # 公共资源
│   ├── Dockerfile              # 前端 Dockerfile
│   ├── nginx.conf              # Nginx 配置
│   └── package.json            # 依赖定义
├── docs/                       # 项目文档
│   ├── DEPLOYMENT.md           # 部署指南
│   ├── API.md                  # API 文档
│   ├── USER_GUIDE.md           # 用户指南
│   └── TEST_REPORT.md          # 测试报告
├── scripts/                    # 项目级脚本
│   └── start.sh                # 一键启动脚本
├── .github/workflows/          # GitHub Actions 工作流
│   ├── ci.yml                  # 持续集成
│   └── publish.yml             # 镜像发布
├── docker-compose.yml          # Docker Compose 配置
├── .env.example                # 环境变量模板
└── README.md                   # 项目说明（本文件）
```

---

## ⚡ 开发模式

### 环境要求

- Go 1.25+
- Node.js 20+
- MySQL 8.0
- Redis 7.x

### 后端开发

```bash
cd backend

# 下载依赖
go mod download

# 启动 API 服务（热重载）
go run ./cmd/api/main.go

# 启动 Worker 服务
go run ./cmd/worker/main.go

# 运行测试
go test ./... -v -count=1

# 代码检查
go vet ./...
```

### 前端开发

```bash
cd frontend

# 安装依赖
npm install

# 启动开发服务器
npm run dev

# 代码检查
npm run lint

# 类型检查
npm run typecheck

# 构建生产包
npm run build
```

---

## 🧪 测试

### 前端测试

```bash
cd frontend
npm run lint          # ESLint 代码检查
npm run typecheck     # TypeScript 类型检查
npm run build         # 生产构建验证
```

### 后端测试

```bash
cd backend
go test ./... -v -count=1 -timeout=5m    # 运行所有单元测试
go test ./... -coverprofile=coverage.out  # 生成覆盖率报告
go vet ./...                              # 代码静态检查
```

### 安全审计

```bash
cd backend
go install github.com/securego/gosec/v2/cmd/gosec@latest
gosec -fmt=json -out=gosec-results.json ./...

go install golang.org/x/vuln/cmd/govulncheck@latest
govulncheck -json ./...
```

---

## 🔒 安全

### 安全特性

- **传输安全**：HTTPS/WSS 强制加密通信
- **认证授权**：JWT + 中间件统一认证，支持角色权限控制
- **输入验证**：所有输入参数严格校验，防止注入攻击
- **速率限制**：API 请求速率限制，防止暴力破解
- **安全头**：CSP、HSTS、X-Frame-Options 等统一安全头
- **WebSocket 安全**：严格源校验，防止未授权连接
- **审计日志**：操作日志脱敏记录，便于安全审计

### 安全报告

发现安全漏洞？请通过 [GitHub Security Advisories](https://github.com/xiaohai09/mobile-cloud-disk-management/security/advisories) 负责任地披露。

详细安全审计报告见 [测试报告](docs/TEST_REPORT.md)。

---

## 📈 监控与可观测性

### 指标采集

- Prometheus 指标端点：`/metrics`
- 业务指标：请求量、延迟、错误率
- 系统指标：CPU、内存、磁盘、网络

### 日志管理

- 结构化 JSON 日志
- 日志级别：DEBUG、INFO、WARN、ERROR
- 审计日志：脱敏处理，保留完整操作链

### 告警配置

- 错误率阈值告警
- 响应时间 P99 告警
- 系统资源告警

---

## 🤝 贡献

我们欢迎社区贡献！请遵循以下流程：

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'feat: add amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 开启 Pull Request

### 贡献指南

- 遵循 [Conventional Commits](https://www.conventionalcommits.org/) 规范
- 确保所有测试通过
- 更新相关文档
- 代码审查通过后合并

---

## 🗺️ 路线图

- [x] 基础架构重构与 Docker 化
- [x] 数据导出功能（CSV/JSON）
- [x] Webhook 通知系统
- [x] 响应式前端界面
- [x] 安全增强与审计日志
- [ ] 多平台适配层（预留接口）
- [ ] Kubernetes Helm Chart 部署
- [ ] 国际化（i18n）支持
- [ ] 插件系统架构
- [ ] 性能优化与缓存策略

---

## ❓ 常见问题

**Q: 系统支持哪些云盘平台？**

A: 当前版本支持主流云盘平台的账号管理与自动化运营。多平台适配层已预留扩展接口，可基于 [docs/USER_GUIDE.md](docs/USER_GUIDE.md) 中的指南进行扩展。

**Q: 如何备份数据？**

A: 建议定期备份 MySQL 数据库和 Redis 数据。具体备份策略见 [部署指南](docs/DEPLOYMENT.md)。

**Q: 支持集群部署吗？**

A: 当前版本基于 Docker Compose 单节点部署。Kubernetes 部署支持正在规划中。

**Q: 如何升级版本？**

A: 参考 [部署指南](docs/DEPLOYMENT.md) 中的升级章节，注意备份数据并遵循迁移指南。

更多问题请查阅 [用户指南](docs/USER_GUIDE.md)。

---

## 📄 许可证

本项目基于 [MIT License](LICENSE) 开源。

---

## 🙏 致谢

- 基于 [3238614968/caiyun](https://github.com/3238614968/caiyun) 架构重构
- 感谢所有贡献者的支持与反馈

---

## 📞 联系方式

- **Issues**: [GitHub Issues](https://github.com/xiaohai09/mobile-cloud-disk-management/issues)
- **Discussions**: [GitHub Discussions](https://github.com/xiaohai09/mobile-cloud-disk-management/discussions)

---

*最后更新：2026-07-07*
