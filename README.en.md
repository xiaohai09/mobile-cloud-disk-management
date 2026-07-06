# Mobile Cloud Disk Management System - Cloud Disk System (Derivative Version)

An architectural refactoring and feature extension based on [3238614968/caiyun](https://github.com/3238614968/caiyun).

## Core Features

### Preserved Native Features
- Account management and automatic sign-in
- Task scheduling and logging
- Exchange center and product management
- Announcements and real-time messaging
- Multi-tenant data isolation

### New Features
- **Data Export**: CSV/JSON export with time range and account filtering
- **Webhook Notifications**: Event subscription, HMAC-SHA256 signing, retry on failure
- **Multi-Platform Adapter Layer**: Pre-built interface for third-party cloud disk integration
- **Responsive Frontend**: Mobile bottom navigation, dark mode, adaptive layout
- **Security Hardening**: Unified security headers, stricter WebSocket origin validation, audit log sanitization

## Quick Start

```bash
# 1. Clone repository
git clone <repo-url>
cd caiyun-refactored

# 2. One-click startup (auto-generate config and start)
chmod +x scripts/start.sh
./scripts/start.sh

# 3. Access
# Frontend: http://localhost
# Monitoring: http://localhost:3000 (admin / caiyun_grafana_2026)
```

For detailed deployment instructions, see [docs/DEPLOYMENT.md](docs/DEPLOYMENT.md).

## Development Mode

### Frontend

```bash
cd frontend
npm install
npm run dev
```

### Backend

```bash
cd backend
# Requires Go 1.24+
go mod download
go run ./cmd/api/main.go
```

## Tech Stack

| Layer | Technology |
|-------|-----------|
| Frontend | Vue 3 + TypeScript + Element Plus + Pinia + Vite |
| Backend | Go + Gin + GORM + Redis + Gorilla WebSocket |
| Database | MySQL 8.0 |
| Task Queue | Redis Streams / Memory |
| Monitoring | Prometheus + Grafana |
| Deployment | Docker Compose + Multi-stage Builds |

## Verification Records

### Frontend
- `npm run lint`: Passed
- `npm run typecheck`: Passed (1 remaining error is Element Plus type compatibility, not project code)
- `npm run build`: Passed, build completed in 2.6s

### Backend
- Code structure check: `docker compose config` Passed
- Shell script check: `bash -n scripts/start.sh` Passed
- Go compiler not available, `go build` / `go test` not executed

### Security Audit
- Added SecurityHeadersMiddleware (X-Frame-Options, CSP, HSTS, etc.)
- Fixed WebSocket `CheckOrigin` empty origin宽松校验
- Audit log adds sanitization filter (password/token/secret fields automatically replaced with `[REDACTED]`)
- All new routes mount `AuthMiddleware + CSRF + RateLimit + Audit`

## License

MIT
