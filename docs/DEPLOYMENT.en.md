# Deployment Guide

## System Requirements

- Docker 20.10+
- Docker Compose v2.0+
- Memory: 4GB+ (8GB recommended)
- Disk: 20GB+ available space

## Quick Deployment

### Method 1: One-Click Startup (Recommended)

```bash
cd /root/caiyun-refactored
chmod +x scripts/start.sh
./scripts/start.sh
```

The script will automatically:
- Generate `.env` configuration file (with random secrets)
- Pull/build Docker images
- Start all services (MySQL, Redis, API, Worker, Frontend, Grafana)
- Perform health checks

### Method 2: Docker Compose Manual Startup

```bash
# 1. Copy environment template
cp .env.example .env

# 2. Modify sensitive config (passwords, JWT secret, etc.)
vim .env

# 3. Start services
docker compose up -d

# 4. View logs
docker compose logs -f
```

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `MYSQL_ROOT_PASSWORD` | MySQL root password | Randomly generated |
| `MYSQL_PASSWORD` | MySQL application user password | Randomly generated |
| `MYSQL_DATABASE` | Database name | `caiyun` |
| `REDIS_PASSWORD` | Redis password | Randomly generated |
| `JWT_SECRET` | JWT signing key | Randomly generated (64 chars) |
| `JWT_ALGORITHM` | JWT algorithm | `HS256` |
| `JWT_ACCESS_TTL` | Access Token TTL | `15m` |
| `JWT_REFRESH_TTL` | Refresh Token TTL | `168h` |
| `WORKER_MONITOR_TOKEN` | Worker monitor token | Randomly generated |
| `GRAFANA_ADMIN_PASSWORD` | Grafana admin password | Randomly generated |
| `CORS_ALLOWED_ORIGINS` | CORS allowed origins | `http://localhost:80` |

## Service Access

| Service | URL | Username | Password |
|---------|-----|----------|----------|
| Frontend | http://localhost | - | - |
| API | http://localhost:8080 | - | - |
| Grafana | http://localhost:3000 | admin | Randomly generated (see .env) |

## Common Issues

### 1. Port Conflict

If ports 80/8080/3000 are occupied, modify port mappings in `docker-compose.yml`:

```yaml
ports:
  - "8081:8080"  # host:container
```

### 2. Database Initialization Failure

Check MySQL container logs:

```bash
docker compose logs mysql
```

Ensure `config/mysql/init.sql` exists and is readable.

### 3. Frontend Build Failure

Ensure Node.js version >= 18:

```bash
node --version
```

### 4. Backend Startup Failure

Check Go version and environment variables:

```bash
go version
cat .env | grep JWT
```

## Production Recommendations

- Use standalone MySQL/Redis instances (instead of Docker Compose)
- Configure SSL/TLS certificates (Let's Encrypt + Caddy/Nginx)
- Enable Redis persistence (AOF + RDB)
- Configure log rotation (logrotate)
- Regular database backups
- Use GitHub Actions for automatic Docker image builds

## Update Deployment

```bash
# 1. Pull latest code
git pull origin main

# 2. Rebuild images
docker compose build

# 3. Restart services (zero downtime)
docker compose up -d --force-recreate

# 4. Verify service status
docker compose ps
```
