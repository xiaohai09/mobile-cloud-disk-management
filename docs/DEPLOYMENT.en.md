# Deployment Guide

[![Docker](https://img.shields.io/badge/Docker-20.10%2B-blue.svg)](https://www.docker.com/)
[![Docker Compose](https://img.shields.io/badge/Docker_Compose-2.0%2B-blue.svg)](https://docs.docker.com/compose/)
[![Go](https://img.shields.io/badge/Go-1.25%2B-blue.svg)](https://go.dev/dl/)
[![Node.js](https://img.shields.io/badge/Node.js-20%2B-green.svg)](https://nodejs.org/)

> This guide covers the complete process from environment preparation, one-click deployment, to production environment optimization.

---

## 📑 Table of Contents

- [System Requirements](#system-requirements)
- [Architecture Overview](#architecture-overview)
- [Quick Deployment](#quick-deployment)
- [Environment Variables](#environment-variables)
- [Service Access Points](#service-access-points)
- [Deployment Verification](#deployment-verification)
- [Production Configuration](#production-configuration)
- [Monitoring and Alerting](#monitoring-and-alerting)
- [Backup and Recovery](#backup-and-recovery)
- [Updates and Upgrades](#updates-and-upgrades)
- [Troubleshooting](#troubleshooting)
- [Frequently Asked Questions](#frequently-asked-questions)

---

## System Requirements

### Minimum Configuration (Development/Testing)

| Resource | Minimum | Recommended |
|----------|---------|-------------|
| CPU | 2 cores | 4 cores |
| Memory | 4 GB | 8 GB |
| Disk | 10 GB | 20 GB SSD |
| OS | Linux x86_64 | Ubuntu 22.04 LTS / Debian 12 |
| Network | Public IP + ports 80/443 | Public IP + domain + SSL certificate |

### Software Dependencies

| Software | Version Requirement | Purpose |
|----------|---------------------|---------|
| Docker Engine | 20.10+ | Container runtime |
| Docker Compose | 2.0+ | Multi-container orchestration |
| Git | 2.0+ | Code cloning |
| curl | Any | Health checks |

---

## Architecture Overview

> **Important**: This project uses a **multi-container architecture**. Single-image one-click deployment is not supported.
> A complete setup requires at least: MySQL, Redis, backend API, Worker, and frontend Nginx.

### Architecture Diagram

```
┌─────────────────────────────────────────────────────────────┐
│                    Docker Compose Network                     │
│                                                              │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐  │
│  │   MySQL     │  │    Redis    │  │   Adminer (Optional) │  │
│  │   8.0       │  │   7.x       │  │  Database Admin UI   │  │
│  │  :3306      │  │  :6379      │  │  :8081               │  │
│  └─────────────┘  └─────────────┘  └─────────────────────┘  │
│        ▲                 ▲                                    │
│        │                 │                                    │
│  ┌─────┴─────────────────┴─────────────────────────────┐   │
│  │           API Service (Go / Gin)                      │   │
│  │          :8080 - REST API + WebSocket                 │   │
│  └──────────────────────────────────────────────────────┘   │
│        ▲                                                    │
│        │                                                    │
│  ┌─────┴─────────────────┐  ┌──────────────────────────┐   │
│  │   Worker Service      │  │   Frontend (Nginx)       │   │
│  │   (Task Scheduling)   │  │   :80 - Vue3 SPA         │   │
│  └───────────────────────┘  └──────────────────────────┘   │
│                                                              │
│  ┌──────────────────────────────────────────────────────┐   │
│  │  Prometheus + Grafana (Monitoring)                    │   │
│  │  :9090 (Prometheus) :3000 (Grafana)                   │   │
│  └──────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
```

### Service Description

| Service | Container Name | Port | Description |
|---------|---------------|------|-------------|
| MySQL | caiyun-mysql | 3306 | Primary database, persistent storage |
| Redis | caiyun-redis | 6379 | Cache, sessions, queues |
| API | caiyun-api | 8080 | RESTful API + WebSocket |
| Worker | caiyun-worker | - | Background task scheduling |
| Frontend | caiyun-frontend | 80 | Vue3 Single Page Application |
| Grafana | caiyun-grafana | 3000 | Monitoring dashboard |
| Prometheus | caiyun-prometheus | 9090 | Metrics collection |

---

## Quick Deployment

### Method 1: Clone Repository One-Click Startup (Recommended)

```bash
# 1. Clone repository
git clone https://github.com/xiaohai09/mobile-cloud-disk-management.git
cd mobile-cloud-disk-management

# 2. Copy environment configuration
cp .env.example .env

# 3. Start services
docker compose up -d

# 4. Verify services
curl http://localhost/api/health
```

### Method 2: Using scripts/start.sh

```bash
git clone https://github.com/xiaohai09/mobile-cloud-disk-management.git
cd mobile-cloud-disk-management
chmod +x scripts/start.sh
./scripts/start.sh
```

The `start.sh` script will automatically:
- Generate `.env` configuration file (with random secrets)
- Pull/build Docker images
- Start all services
- Perform health checks
- Automatically create default admin account
- Print access URLs and default credentials

### Method 3: Docker Compose Manual Startup

```bash
# 1. Copy environment variable template
cp .env.example .env

# 2. Modify sensitive configuration (passwords, JWT secret, etc.)
vim .env

# 3. Start services
docker compose up -d

# 4. View logs
docker compose logs -f
```

---

## Environment Variables

### Complete Environment Variable List

| Variable | Description | Default Value | Required |
|----------|-------------|---------------|----------|
| `MYSQL_ROOT_PASSWORD` | MySQL root password | Randomly generated | ✅ |
| `MYSQL_PASSWORD` | MySQL application user password | Randomly generated | ✅ |
| `MYSQL_DATABASE` | Database name | `caiyun` | ❌ |
| `REDIS_PASSWORD` | Redis password | Randomly generated | ✅ |
| `JWT_SECRET` | JWT signing key | Randomly generated (64 chars) | ✅ |
| `JWT_ALGORITHM` | JWT algorithm | `HS256` | ❌ |
| `JWT_ACCESS_TTL` | Access Token TTL | `15m` | ❌ |
| `JWT_REFRESH_TTL` | Refresh Token TTL | `168h` | ❌ |
| `WORKER_MONITOR_TOKEN` | Worker monitor token | Randomly generated | ✅ |
| `GRAFANA_ADMIN_PASSWORD` | Grafana admin password | Randomly generated | ✅ |
| `CORS_ALLOWED_ORIGINS` | CORS allowed origins | `http://localhost:80` | ❌ |
| `DEFAULT_ADMIN_USERNAME` | Default admin username | `admin` | ❌ |
| `DEFAULT_ADMIN_PASSWORD` | Default admin password | `admin123` | ❌ |

### Configuration Examples

#### Database Configuration

```bash
# MySQL configuration
MYSQL_ROOT_PASSWORD=your_secure_root_password
MYSQL_PASSWORD=your_secure_app_password
MYSQL_DATABASE=caiyun
```

#### Cache Configuration

```bash
# Redis configuration
REDIS_PASSWORD=your_secure_redis_password
```

#### JWT Authentication Configuration

```bash
# JWT configuration (MUST change in production)
JWT_SECRET=your_very_long_and_random_secret_key_here_at_least_64_chars
JWT_ALGORITHM=HS256
JWT_ACCESS_TTL=15m
JWT_REFRESH_TTL=168h
```

#### Monitoring Configuration

```bash
# Grafana configuration
GRAFANA_ADMIN_PASSWORD=your_secure_grafana_password
```

#### CORS Configuration

```bash
# Production environment example
CORS_ALLOWED_ORIGINS=https://your-domain.com,https://www.your-domain.com
```

---

## Service Access Points

| Service | URL | Description |
|---------|-----|-------------|
| Frontend | http://localhost | User interface |
| API Service | http://localhost:8080 | RESTful API |
| API Health Check | http://localhost:8080/health | Health status endpoint |
| Grafana | http://localhost:3000 | Monitoring dashboard |
| Prometheus | http://localhost:9090 | Metrics query |
| Adminer (Optional) | http://localhost:8081 | Database management |

### Default Admin Account

| Field | Value |
|-------|-------|
| Username | `admin` |
| Password | `admin123` |

> **Security Note**: Change the default password immediately after first login. In production, configure strong passwords via environment variables.

---

## Deployment Verification

### Health Checks

```bash
# Check API health status
curl http://localhost:8080/health

# Expected response
{"status":"ok"}
```

### Container Status Check

```bash
# View all container status
docker compose ps

# Expected output
NAME                    SERVICE         STATUS
caiyun-api              api             running
caiyun-frontend         frontend        running
caiyun-mysql            mysql           running
caiyun-redis            redis           running
caiyun-worker           worker          running
caiyun-grafana          grafana         running
caiyun-prometheus       prometheus      running
```

### Log Inspection

```bash
# View all service logs
docker compose logs -f

# View specific service logs
docker compose logs -f api
docker compose logs -f frontend
docker compose logs -f worker
```

---

## Production Configuration

### Recommended Architecture

For production, use standalone MySQL/Redis instances instead of Docker Compose built-in services.

#### Independent Database Configuration

```yaml
# docker-compose.prod.yml example
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

### SSL/TLS Configuration

HTTPS is mandatory for production:

```bash
# Using Let's Encrypt + Caddy
# 1. Install Caddy
curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/gpg.key' | sudo apt-key add -
curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/debian.deb.txt' | sudo tee /etc/apt/sources.list.d/caddy-stable.list

# 2. Install and configure
sudo apt update && sudo apt install caddy

# 3. Configure Caddyfile
cat > /etc/caddy/Caddyfile << 'EOF'
your-domain.com {
    reverse_proxy localhost:80
    encode gzip
    tls your-email@example.com
}
EOF

# 4. Start Caddy
sudo systemctl enable --now caddy
```

### Security Hardening Checklist

- [ ] Change all default passwords
- [ ] Use strong JWT Secret (at least 64 random characters)
- [ ] Configure firewall rules (only open 80/443/22)
- [ ] Enable Redis password authentication
- [ ] Restrict MySQL remote access
- [ ] Regular database backups
- [ ] Enable audit logging
- [ ] Configure fail2ban for brute force protection

---

## Monitoring and Alerting

### Prometheus Configuration

Prometheus is built-in and automatically collects:

- HTTP request volume, latency, error rate
- WebSocket connection count
- Database query performance
- Task queue depth
- System resource utilization

### Grafana Dashboards

Access http://localhost:3000 for monitoring:

- **System Overview**: CPU, memory, disk, network
- **API Performance**: Request latency, error rate, throughput
- **Database**: Connection count, query performance, slow queries
- **Business Metrics**: Task success rate, export volume, Webhook delivery rate

### Alert Rules

Configure the following alert rules in Grafana:

```yaml
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

## Backup and Recovery

### Database Backup

```bash
# Manual backup
docker compose exec mysql mysqldump -u root -p caiyun > backup_$(date +%Y%m%d).sql

# Automated backup script (cron)
cat > /etc/cron.d/caiyun-backup << 'EOF'
0 2 * * * root docker compose exec mysql mysqldump -u root -p caiyun > /backup/caiyun_$(date +\%Y\%m\%d).sql
EOF
```

### Data Recovery

```bash
# Restore database
docker compose exec mysql mysql -u root -p caiyun < backup_20250101.sql

# Restore Redis
docker compose exec redis redis-cli --pass your_redis_password FLUSHALL
# Then restore from backup
```

### Backup Strategy

| Data Type | Frequency | Retention | Storage Location |
|-----------|-----------|-----------|------------------|
| MySQL | Daily full backup | 30 days | Off-site storage |
| Redis | Hourly snapshot | 7 days | Local + off-site |
| Configuration | Every change | Permanent | Git version control |
| Docker images | Every build | Permanent | GHCR |

---

## Updates and Upgrades

### Standard Update Process

```bash
# 1. Pull latest code
git pull origin main

# 2. Review changes
git log --oneline -5

# 3. Backup data (IMPORTANT!)
docker compose exec mysql mysqldump -u root -p caiyun > backup_$(date +%Y%m%d).sql

# 4. Rebuild images
docker compose build

# 5. Run database migrations (if any)
docker compose run --rm api ./migrate

# 6. Restart services (zero downtime)
docker compose up -d --force-recreate

# 7. Verify service status
docker compose ps
curl http://localhost/api/health

# 8. Clean old images
docker image prune -f
```

### Rollback Process

```bash
# 1. Rollback to previous version
git checkout HEAD~1

# 2. Rebuild
docker compose build

# 3. Restart services
docker compose up -d --force-recreate
```

---

## Troubleshooting

### Port Conflicts

```bash
# Check port usage
sudo netstat -tlnp | grep :80

# Modify docker-compose.yml port mapping
ports:
  - "8081:80"  # host:container
```

### Container Startup Failures

```bash
# View container logs
docker compose logs <service-name>

# Common causes:
# 1. Environment variables not configured
# 2. Port already in use
# 3. Insufficient disk space
# 4. Permission issues
```

### Database Connection Failures

```bash
# Check MySQL status
docker compose exec mysql mysql -u root -p -e "SELECT 1"

# Check network connectivity
docker compose exec api ping mysql

# Check environment variables
docker compose exec api env | grep DB_
```

### Frontend Build Failures

```bash
# Check Node.js version
docker compose exec frontend node --version

# Clean cache and rebuild
docker compose build --no-cache frontend
```

---

## Frequently Asked Questions

**Q: Why is cloning the repository recommended over single-image deployment?**

A: This project contains multiple services (MySQL, Redis, API, Worker, Frontend) with custom initialization scripts and environment configuration. A single image cannot cover the complete multi-service architecture and configuration flexibility.

**Q: How to upgrade to the latest version?**

A: Refer to the [Updates and Upgrades](#updates-and-upgrades) section. Always test in a staging environment before upgrading production.

**Q: What is the data backup strategy?**

A: Daily full backups of MySQL and hourly Redis snapshots are recommended. Detailed configuration in [Backup and Recovery](#backup-and-recovery).

**Q: How to configure HTTPS?**

A: Let's Encrypt + Caddy/Nginx is recommended. Detailed configuration in [Production Configuration](#production-configuration).

**Q: Is Kubernetes deployment supported?**

A: The current version is based on Docker Compose. Kubernetes Helm Chart support is planned.

**Q: How to monitor system status?**

A: Built-in Prometheus + Grafana monitoring. Access http://localhost:3000 for dashboards. Detailed configuration in [Monitoring and Alerting](#monitoring-and-alerting).

For more questions, please refer to the [User Guide](USER_GUIDE.md).

---

## 📞 Technical Support

- **Documentation Issues**: [Submit Issue](https://github.com/xiaohai09/mobile-cloud-disk-management/issues)
- **Security Issues**: [Security Advisories](https://github.com/xiaohai09/mobile-cloud-disk-management/security/advisories)
- **Discussions**: [GitHub Discussions](https://github.com/xiaohai09/mobile-cloud-disk-management/discussions)

---

*Last updated: 2026-07-07*
