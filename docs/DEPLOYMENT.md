# 部署指南

[![Docker](https://img.shields.io/badge/Docker-20.10%2B-blue.svg)](https://www.docker.com/)
[![Docker Compose](https://img.shields.io/badge/Docker_Compose-2.0%2B-blue.svg)](https://docs.docker.com/compose/)
[![Go](https://img.shields.io/badge/Go-1.25%2B-blue.svg)](https://go.dev/dl/)
[![Node.js](https://img.shields.io/badge/Node.js-20%2B-green.svg)](https://nodejs.org/)

> 本指南涵盖从环境准备、一键部署到生产环境优化的完整流程。

---

## 📑 目录

- [系统要求](#系统要求)
- [架构说明](#架构说明)
- [快速部署](#快速部署)
- [环境变量配置](#环境变量配置)
- [服务访问地址](#服务访问地址)
- [验证部署](#验证部署)
- [生产环境配置](#生产环境配置)
- [监控与告警](#监控与告警)
- [备份与恢复](#备份与恢复)
- [更新与升级](#更新与升级)
- [故障排查](#故障排查)
- [常见问题](#常见问题)

---

## 系统要求

### 最低配置（开发/测试）

| 资源 | 最低要求 | 推荐配置 |
|------|---------|---------|
| CPU | 2 核 | 4 核 |
| 内存 | 4 GB | 8 GB |
| 磁盘 | 10 GB | 20 GB SSD |
| 操作系统 | Linux x86_64 | Ubuntu 22.04 LTS / Debian 12 |
| 网络 | 公网 IP + 80/443 端口 | 公网 IP + 域名 + SSL 证书 |

### 软件依赖

| 软件 | 版本要求 | 用途 |
|------|---------|------|
| Docker Engine | 20.10+ | 容器运行时 |
| Docker Compose | 2.0+ | 多容器编排 |
| Git | 2.0+ | 代码克隆 |
| curl | 任意 | 健康检查 |

---

## 架构说明

本项目采用**多容器架构**，不支持单镜像一键拉取部署。完整运行需要以下服务：

```
┌─────────────────────────────────────────────────────────────┐
│                    Docker Compose Network                     │
│                                                              │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐  │
│  │   MySQL     │  │    Redis    │  │   Adminer (可选)     │  │
│  │  8.0        │  │   7.x       │  │  数据库管理界面      │  │
│  │  :3306      │  │  :6379      │  │  :8081              │  │
│  └─────────────┘  └─────────────┘  └─────────────────────┘  │
│        ▲                 ▲                                    │
│        │                 │                                    │
│  ┌─────┴─────────────────┴─────────────────────────────┐   │
│  │              API Service (Go/Gin)                     │   │
│  │         :8080 - REST API + WebSocket                 │   │
│  └──────────────────────────────────────────────────────┘   │
│        ▲                                                    │
│        │                                                    │
│  ┌─────┴─────────────────┐  ┌──────────────────────────┐   │
│  │   Worker Service      │  │   Frontend (Nginx)       │   │
│  │   (任务调度)          │  │   :80 - Vue3 SPA         │   │
│  └───────────────────────┘  └──────────────────────────┘   │
│                                                              │
│  ┌──────────────────────────────────────────────────────┐   │
│  │  Prometheus + Grafana (监控)                          │   │
│  │  :9090 (Prometheus) :3000 (Grafana)                   │   │
│  └──────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
```

### 服务说明

| 服务 | 容器名 | 端口 | 说明 |
|------|--------|------|------|
| MySQL | caiyun-mysql | 3306 | 主数据库，持久化存储 |
| Redis | caiyun-redis | 6379 | 缓存、会话、队列 |
| API | caiyun-api | 8080 | RESTful API + WebSocket |
| Worker | caiyun-worker | - | 后台任务调度 |
| Frontend | caiyun-frontend | 80 | Vue3 单页应用 |
| Grafana | caiyun-grafana | 3000 | 监控仪表板 |
| Prometheus | caiyun-prometheus | 9090 | 指标采集 |

---

## 快速部署

### 方式一：克隆仓库一键启动（推荐）

```bash
# 1. 克隆仓库
git clone https://github.com/xiaohai09/mobile-cloud-disk-management.git
cd mobile-cloud-disk-management

# 2. 复制环境配置
cp .env.example .env

# 3. 启动服务
docker compose up -d

# 4. 验证服务
curl http://localhost/api/health
```

### 方式二：使用 scripts/start.sh

```bash
git clone https://github.com/xiaohai09/mobile-cloud-disk-management.git
cd mobile-cloud-disk-management
chmod +x scripts/start.sh
./scripts/start.sh
```

`start.sh` 脚本会自动完成：
- 生成 `.env` 配置文件（含随机 secrets）
- 拉取/构建 Docker 镜像
- 启动所有服务
- 执行健康检查
- 自动创建默认管理员账号
- 打印访问地址和默认账号密码

### 方式三：Docker Compose 手动启动

```bash
# 1. 复制环境变量模板
cp .env.example .env

# 2. 修改敏感配置（密码、JWT secret 等）
vim .env

# 3. 启动服务
docker compose up -d

# 4. 查看日志
docker compose logs -f
```

---

## 环境变量配置

### 完整环境变量列表

| 变量 | 说明 | 默认值 | 必填 |
|------|------|--------|------|
| `MYSQL_ROOT_PASSWORD` | MySQL root 密码 | 随机生成 | ✅ |
| `MYSQL_PASSWORD` | MySQL 应用用户密码 | 随机生成 | ✅ |
| `MYSQL_DATABASE` | 数据库名 | `caiyun` | ❌ |
| `REDIS_PASSWORD` | Redis 密码 | 随机生成 | ✅ |
| `JWT_SECRET` | JWT 签名密钥 | 随机生成（64位） | ✅ |
| `JWT_ALGORITHM` | JWT 算法 | `HS256` | ❌ |
| `JWT_ACCESS_TTL` | Access Token 有效期 | `15m` | ❌ |
| `JWT_REFRESH_TTL` | Refresh Token 有效期 | `168h` | ❌ |
| `WORKER_MONITOR_TOKEN` | Worker 监控令牌 | 随机生成 | ✅ |
| `GRAFANA_ADMIN_PASSWORD` | Grafana 管理员密码 | 随机生成 | ✅ |
| `CORS_ALLOWED_ORIGINS` | CORS 允许源 | `http://localhost:80` | ❌ |
| `DEFAULT_ADMIN_USERNAME` | 默认管理员用户名 | `admin` | ❌ |
| `DEFAULT_ADMIN_PASSWORD` | 默认管理员密码 | `admin123` | ❌ |

### 环境变量说明

#### 数据库配置

```bash
# MySQL 配置
MYSQL_ROOT_PASSWORD=your_secure_root_password
MYSQL_PASSWORD=your_secure_app_password
MYSQL_DATABASE=caiyun
```

#### 缓存配置

```bash
# Redis 配置
REDIS_PASSWORD=your_secure_redis_password
```

#### JWT 认证配置

```bash
# JWT 配置（生产环境务必修改）
JWT_SECRET=your_very_long_and_random_secret_key_here_at_least_64_chars
JWT_ALGORITHM=HS256
JWT_ACCESS_TTL=15m
JWT_REFRESH_TTL=168h
```

#### 监控配置

```bash
# Grafana 配置
GRAFANA_ADMIN_PASSWORD=your_secure_grafana_password
```

#### CORS 配置

```bash
# 生产环境示例
CORS_ALLOWED_ORIGINS=https://your-domain.com,https://www.your-domain.com
```

---

## 服务访问地址

| 服务 | 地址 | 说明 |
|------|------|------|
| 前端应用 | http://localhost | 用户界面 |
| API 服务 | http://localhost:8080 | RESTful API |
| API 健康检查 | http://localhost:8080/health | 健康状态端点 |
| Grafana 监控 | http://localhost:3000 | 监控仪表板 |
| Prometheus | http://localhost:9090 | 指标查询 |
| Adminer (可选) | http://localhost:8081 | 数据库管理 |

### 默认管理员账号

| 项目 | 值 |
|------|-----|
| 用户名 | `admin` |
| 密码 | `admin123` |

> **安全提示**：首次登录后请立即修改默认密码。生产环境部署时，请通过环境变量配置强密码。

---

## 验证部署

### 健康检查

```bash
# 检查 API 健康状态
curl http://localhost:8080/health

# 预期响应
{"status":"ok"}
```

### 服务状态检查

```bash
# 查看所有容器状态
docker compose ps

# 预期输出
NAME                    SERVICE         STATUS
caiyun-api              api             running
caiyun-frontend         frontend        running
caiyun-mysql            mysql           running
caiyun-redis            redis           running
caiyun-worker           worker          running
```

### 日志检查

```bash
# 查看所有服务日志
docker compose logs -f

# 查看特定服务日志
docker compose logs -f api
docker compose logs -f frontend
docker compose logs -f worker
```

---

## 生产环境配置

### 推荐配置

生产环境建议使用独立的 MySQL/Redis 实例，而非 Docker Compose 内置服务。

#### 独立数据库配置

```yaml
# docker-compose.prod.yml 示例
services:
  api:
    environment:
      DB_HOST: your-mysql-host
      DB_PORT: 3306
      DB_USER: caiyun_app
      DB_PASSWORD: ${MYSQL_PASSWORD}
      DB_NAME: caiyun
      REDIS_HOST: your-redis-host
      REDIS_PORT: 6379
      REDIS_PASSWORD: ${REDIS_PASSWORD}
```

### SSL/TLS 配置

生产环境必须启用 HTTPS：

```bash
# 使用 Let's Encrypt + Caddy
# 1. 安装 Caddy
curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/gpg.key' | sudo apt-key add -
curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/debian.deb.txt' | sudo tee /etc/apt/sources.list.d/caddy-stable.list

# 2. 安装并配置
sudo apt update && sudo apt install caddy

# 3. 配置 Caddyfile
cat > /etc/caddy/Caddyfile << 'EOF'
your-domain.com {
    reverse_proxy localhost:80
    encode gzip
    tls your-email@example.com
}
EOF

# 4. 启动 Caddy
sudo systemctl enable --now caddy
```

### 安全加固

- [ ] 修改所有默认密码
- [ ] 使用强 JWT Secret（至少 64 位随机字符串）
- [ ] 配置防火墙规则（仅开放 80/443/22）
- [ ] 启用 Redis 密码认证
- [ ] 配置 MySQL 远程访问限制
- [ ] 定期备份数据库
- [ ] 启用审计日志
- [ ] 配置 fail2ban 防暴力破解

---

## 监控与告警

### Prometheus 配置

Prometheus 已内置，自动采集以下指标：

- HTTP 请求量、延迟、错误率
- WebSocket 连接数
- 数据库查询性能
- 任务队列深度
- 系统资源使用率

### Grafana 仪表板

访问 http://localhost:3000 查看监控面板：

- **系统概览**：CPU、内存、磁盘、网络
- **API 性能**：请求延迟、错误率、吞吐量
- **数据库**：连接数、查询性能、慢查询
- **业务指标**：任务成功率、导出量、Webhook 送达率

### 告警配置

在 Grafana 中配置以下告警规则：

```yaml
# 示例告警规则
groups:
  - name: caiyun_alerts
    rules:
      - alert: HighErrorRate
        expr: rate(http_requests_total{status=~"5.."}[5m]) > 0.05
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: "High error rate detected"
          
      - alert: HighLatency
        expr: histogram_quantile(0.99, http_request_duration_seconds) > 1
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "High latency detected"
```

---

## 备份与恢复

### 数据库备份

```bash
# 手动备份
docker compose exec mysql mysqldump -u root -p caiyun > backup_$(date +%Y%m%d).sql

# 自动备份脚本（cron）
cat > /etc/cron.d/caiyun-backup << 'EOF'
0 2 * * * root docker compose exec mysql mysqldump -u root -p caiyun > /backup/caiyun_$(date +\%Y\%m\%d).sql
EOF
```

### 数据恢复

```bash
# 恢复数据库
docker compose exec mysql mysql -u root -p caiyun < backup_20250101.sql

# 恢复 Redis
docker compose exec redis redis-cli --pass your_redis_password FLUSHALL
# 然后从备份恢复
```

### 备份策略建议

| 数据类型 | 备份频率 | 保留周期 | 存储位置 |
|---------|---------|---------|---------|
| MySQL | 每日全量备份 | 30 天 | 异地存储 |
| Redis | 每小时快照 | 7 天 | 本地 + 异地 |
| 配置文件 | 每次变更 | 永久 | Git 版本控制 |
| Docker 镜像 | 每次构建 | 永久 | GHCR |

---

## 更新与升级

### 标准更新流程

```bash
# 1. 拉取最新代码
git pull origin main

# 2. 检查变更
git log --oneline -5

# 3. 备份数据（重要！）
docker compose exec mysql mysqldump -u root -p caiyun > backup_$(date +%Y%m%d).sql

# 4. 重新构建镜像
docker compose build

# 5. 执行数据库迁移（如有）
docker compose run --rm api ./migrate

# 6. 重启服务（零停机）
docker compose up -d --force-recreate

# 7. 验证服务状态
docker compose ps
curl http://localhost/api/health

# 8. 清理旧镜像
docker image prune -f
```

### 回滚流程

```bash
# 1. 回滚到上一个版本
git checkout HEAD~1

# 2. 重新构建
docker compose build

# 3. 重启服务
docker compose up -d --force-recreate
```

---

## 故障排查

### 端口冲突

```bash
# 检查端口占用
sudo netstat -tlnp | grep :80

# 修改 docker-compose.yml 端口映射
ports:
  - "8081:80"  # 主机端口:容器端口
```

### 容器启动失败

```bash
# 查看容器日志
docker compose logs <service-name>

# 常见原因：
# 1. 环境变量未配置
# 2. 端口被占用
# 3. 磁盘空间不足
# 4. 权限问题
```

### 数据库连接失败

```bash
# 检查 MySQL 状态
docker compose exec mysql mysql -u root -p -e "SELECT 1"

# 检查网络连通性
docker compose exec api ping mysql

# 检查环境变量
docker compose exec api env | grep DB_
```

### 前端构建失败

```bash
# 检查 Node.js 版本
docker compose exec frontend node --version

# 清理缓存重新构建
docker compose build --no-cache frontend
```

---

## 常见问题

**Q: 为什么推荐克隆仓库部署，而不是单镜像部署？**

A: 本项目包含 MySQL、Redis、API、Worker、Frontend 等多个服务，且包含自定义初始化脚本和环境变量配置。单镜像无法涵盖完整的多服务架构和配置灵活性。

**Q: 如何升级到最新版本？**

A: 参考 [更新与升级](#更新与升级) 章节，建议先在测试环境验证后再升级生产环境。

**Q: 数据备份策略是什么？**

A: 建议每日全量备份 MySQL，每小时 Redis 快照。详细配置见 [备份与恢复](#备份与恢复)。

**Q: 如何配置 HTTPS？**

A: 推荐使用 Let's Encrypt + Caddy/Nginx，详细配置见 [生产环境配置](#生产环境配置)。

**Q: 支持 Kubernetes 部署吗？**

A: 当前版本基于 Docker Compose。Kubernetes Helm Chart 支持正在规划中。

**Q: 如何监控系统运行状态？**

A: 内置 Prometheus + Grafana 监控，访问 http://localhost:3000 查看仪表板。详细配置见 [监控与告警](#监控与告警)。

更多问题请查阅 [用户指南](USER_GUIDE.md)。

---

## 📞 技术支持

- **文档问题**：[提交 Issue](https://github.com/xiaohai09/mobile-cloud-disk-management/issues)
- **安全问题**：[Security Advisories](https://github.com/xiaohai09/mobile-cloud-disk-management/security/advisories)
- **讨论交流**：[GitHub Discussions](https://github.com/xiaohai09/mobile-cloud-disk-management/discussions)

---

*最后更新：2026-07-07*
