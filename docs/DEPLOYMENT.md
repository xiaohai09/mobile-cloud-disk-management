# Caiyun Refactored - 部署文档

## 系统要求

- Docker Engine >= 20.10
- Docker Compose >= 2.0
- 至少 2GB 可用内存
- 至少 5GB 可用磁盘空间

## 快速部署

### 1. 克隆项目

```bash
git clone https://github.com/your-org/caiyun-refactored.git
cd caiyun-refactored
```

### 2. 一键启动

```bash
chmod +x scripts/start.sh
./scripts/start.sh
```

脚本会自动：
- 生成强随机密码并写入 `.env`
- 构建 Docker 镜像
- 启动所有服务
- 执行健康检查

### 3. 访问服务

- 前端界面：http://localhost
- API 服务：http://localhost:8080
- Grafana 监控：http://localhost:3000

## 环境变量

| 变量 | 说明 | 默认值 |
|------|------|--------|
| `ENV` | 运行环境 | `production` |
| `DB_HOST` | MySQL 主机 | `mysql` |
| `DB_PORT` | MySQL 端口 | `3306` |
| `DB_USER` | MySQL 用户 | `caiyun` |
| `DB_PASSWORD` | MySQL 密码 | 自动生成 |
| `REDIS_PASSWORD` | Redis 密码 | 自动生成 |
| `JWT_SECRET` | JWT 签名密钥 | 自动生成 |
| `PORT` | API 端口 | `8080` |

## 服务说明

| 服务 | 镜像 | 端口 | 说明 |
|------|------|------|------|
| mysql | mysql:8.0 | 127.0.0.1:3306 | 主数据库 |
| redis | redis:7.0 | 127.0.0.1:6379 | 缓存与会话 |
| backend-api | caiyun-api | 8080 | Go API 服务 |
| backend-worker | caiyun-worker | - | 后台任务 worker |
| frontend | caiyun-frontend | 80 | Vue 前端 |
| grafana | grafana/grafana | 3000 | 监控面板 |

## 常用命令

```bash
# 查看服务状态
docker compose ps

# 查看日志
docker compose logs -f backend-api
docker compose logs -f backend-worker

# 重启服务
docker compose restart backend-api

# 停止服务
docker compose down

# 更新镜像
docker compose pull
docker compose up -d --force-recreate
```

## 安全说明

- MySQL/Redis 仅监听 `127.0.0.1`，不暴露公网
- 所有敏感配置通过环境变量注入
- JWT 使用强随机密钥
- API 启用限流、审计日志、CORS 校验
- 非 root 用户运行容器
