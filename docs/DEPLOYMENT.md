# 部署指南

## 系统要求

- Docker 20.10+
- Docker Compose v2.0+
- 内存：4GB+（推荐 8GB）
- 磁盘：20GB+ 可用空间

## 快速部署

### 方式一：一键启动（推荐）

```bash
cd /root/caiyun-refactored
chmod +x scripts/start.sh
./scripts/start.sh
```

脚本会自动完成：
- 生成 `.env` 配置文件（含随机 secrets）
- 拉取/构建 Docker 镜像
- 启动所有服务（MySQL、Redis、API、Worker、Frontend、Grafana）
- 执行健康检查

### 方式二：Docker Compose 手动启动

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

## 环境变量说明

| 变量 | 说明 | 默认值 |
|------|------|--------|
| `MYSQL_ROOT_PASSWORD` | MySQL root 密码 | 随机生成 |
| `MYSQL_PASSWORD` | MySQL 应用用户密码 | 随机生成 |
| `MYSQL_DATABASE` | 数据库名 | `caiyun` |
| `REDIS_PASSWORD` | Redis 密码 | 随机生成 |
| `JWT_SECRET` | JWT 签名密钥 | 随机生成（64位） |
| `JWT_ALGORITHM` | JWT 算法 | `HS256` |
| `JWT_ACCESS_TTL` | Access Token 有效期 | `15m` |
| `JWT_REFRESH_TTL` | Refresh Token 有效期 | `168h` |
| `WORKER_MONITOR_TOKEN` | Worker 监控令牌 | 随机生成 |
| `GRAFANA_ADMIN_PASSWORD` | Grafana 管理员密码 | 随机生成 |
| `CORS_ALLOWED_ORIGINS` | CORS 允许源 | `http://localhost:80` |

## 服务访问地址

| 服务 | 地址 | 账号 | 密码 |
|------|------|------|------|
| 前端 | http://localhost | - | - |
| API | http://localhost:8080 | - | - |
| Grafana | http://localhost:3000 | admin | 随机生成（见 .env） |

## 常见问题

### 1. 端口冲突

如果 80/8080/3000 端口被占用，修改 `docker-compose.yml` 中的端口映射：

```yaml
ports:
  - "8081:8080"  # 主机端口:容器端口
```

### 2. 数据库初始化失败

检查 MySQL 容器日志：

```bash
docker compose logs mysql
```

确保 `config/mysql/init.sql` 文件存在且可读。

### 3. 前端构建失败

确保 Node.js 版本 >= 18：

```bash
node --version
```

### 4. 后端启动失败

检查 Go 版本和环境变量：

```bash
go version
cat .env | grep JWT
```

## 生产环境建议

- 使用独立的 MySQL/Redis 实例（而非 Docker Compose）
- 配置 SSL/TLS 证书（Let's Encrypt + Caddy/Nginx）
- 启用 Redis 持久化（AOF + RDB）
- 配置日志轮转（logrotate）
- 定期备份数据库
- 使用 GitHub Actions 自动构建 Docker 镜像

## 更新部署

```bash
# 1. 拉取最新代码
git pull origin main

# 2. 重新构建镜像
docker compose build

# 3. 重启服务（零停机）
docker compose up -d --force-recreate

# 4. 验证服务状态
docker compose ps
```
