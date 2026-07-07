# Mobile Cloud Disk Management System

[![CI Status](https://github.com/xiaohai09/mobile-cloud-disk-management/actions/workflows/ci.yml/badge.svg)](https://github.com/xiaohai09/mobile-cloud-disk-management/actions/workflows/ci.yml/badge)
[![Docker Build](https://github.com/xiaohai09/mobile-cloud-disk-management/actions/workflows/publish.yml/badge.svg)](https://github.com/xiaohai09/mobile-cloud-disk-management/actions/workflows/publish.yml/badge)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.25-blue.svg)](https://go.dev/dl/)
[![Node Version](https://img.shields.io/badge/Node-20-green.svg)](https://nodejs.org/)
[![Docker](https://img.shields.io/badge/Docker-Ready-blue.svg)](https://www.docker.com/)
[![Chinese](https://img.shields.io/badge/Chinese-README-red.svg)](README.md)
[![English](https://img.shields.io/badge/README-English-green.svg)(README.en.md)

> **Enterprise-grade Multi-Cloud Disk Management Platform** — An architectural refactoring and feature extension based on [3238614968/caiyun](https://github.com/3238614968/caiyun), providing enterprise features such as data export, Webhook notifications, responsive frontend, and security hardening.

---

## 📑 Documentation Navigation (Multi-language)

| Language | Document | Description |
|----------|----------|-------------|
| 🇨🇳 **Chinese** | [README.md](README.md) | Project overview and quick start (default) |
| 🇺🇸 **English** | [README.en.md](README.en.md) | Project overview and quick start |
| 🇨🇳 **Chinese** | [Deployment Guide](docs/DEPLOYMENT.md) | Docker deployment, environment config, troubleshooting |
| 🇺🇸 **English** | [Deployment Guide](docs/DEPLOYMENT.en.md) | Docker deployment, config, troubleshooting |
| 🇨🇳 **Chinese** | [API Documentation](docs/API.md) | Authentication, export API, Webhook API |
| 🇺🇸 **English** | [API Documentation](docs/API.en.md) | Auth, export API, webhook API |
| 🇨🇳 **Chinese** | [User Guide](docs/USER_GUIDE.md) | Role descriptions, feature modules, FAQ |
| 🇺🇸 **English** | [User Guide](docs/USER_GUIDE.en.md) | Roles, features, FAQ |
| 🇨🇳 **Chinese** | [Test Report](docs/TEST_REPORT.md) | Build, E2E, Docker, security audit |
| 🇺🇸 **English** | [Test Report](docs/TEST_REPORT.en.md) | Build, E2E, Docker, security audit |

---

## 🎯 Project Overview

### What is Mobile Cloud Disk Management System?

Mobile Cloud Disk Management System is an **enterprise-grade multi-cloud disk account management and automated operations platform** for unified management of multiple cloud disk service accounts, providing automated sign-in, task scheduling, data export, Webhook notifications, and other core capabilities. The system adopts a frontend-backend separation architecture, supports Docker containerized deployment, and is suitable for enterprise internal cloud disk resource management, data analysis, and automated operation scenarios.

### Core Value

- **Unified Management**: Centralized management of multiple cloud disk service accounts with multi-tenant data isolation
- **Automated Operations**: Scheduled tasks, automatic sign-in, redemption center, reducing manual operation costs
- **Data Observability**: CSV/JSON data export, Webhook event push, easy integration into enterprise data flows
- **Security Compliance**: Unified security headers, WebSocket origin validation, audit log sanitization, rate limiting
- **Easy Deployment**: Docker Compose one-click deployment, multi-stage builds, images pushed to GHCR

### Use Cases

- Enterprise internal cloud disk account centralized management and operations
- Automated sign-in and task scheduling
- Data export and third-party system integration
- Monitoring, alerting, and event-driven architecture

---

## ✨ Core Features

### 📊 Account Management
- Unified management of multiple accounts with batch import and grouping
- Automatic sign-in and task scheduling
- Real-time account status monitoring
- Multi-tenant data isolation for data security

### 🔄 Tasks and Automation
- Scheduled task scheduling (Cron expression support)
- Task logs and execution records
- Automatic retry on failure
- Real-time task status notifications

### 🎁 Redemption Center
- Redemption code management and distribution
- Redemption record tracking
- Product management (add, edit, delete)
- Redemption statistics and reports

### 📢 Announcements and Notifications
- System announcement publishing and management
- Real-time message push (WebSocket)
- Announcement reading status tracking

### 📤 Data Export
- CSV/JSON format export support
- Filterable time range and accounts
- Export task queue processing
- Export file lifecycle management

### 🔔 Webhook Notifications
- Event subscription and push
- HMAC-SHA256 signature verification
- Automatic retry on failure (exponential backoff)
- Multiple event type support

### 📱 Responsive Frontend
- Vue 3 + TypeScript modern tech stack
- Element Plus enterprise-grade UI component library
- Mobile bottom navigation, adaptive layout
- Dark mode support

### 🔒 Security Hardening
- Unified security headers (CSP, HSTS, X-Frame-Options)
- WebSocket strict origin validation
- Audit log sanitization filtering
- Unified authentication, CSRF, rate limiting middleware

---

## 🏗️ System Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    Frontend (Vue 3 SPA)                      │
│  ┌────────────┐  ┌────────────┐  ┌────────────┐  ┌────────┐ │
│  │   Account  │  │   Task     │  │ Redemption │  │ Export │ │
│  │ Management │  │  Center    │  │   Center   │  │        │ │
│  └────────────┘  └────────────┘  └────────────┘  └────────┘ │
└─────────────────────────────────────────────────────────────┘
                          │ HTTPS / WSS
                          ▼
┌─────────────────────────────────────────────────────────────┐
│                API Gateway (Go / Gin)                        │
│  ┌───────────────┐  ┌──────────────┐  ┌──────────────────┐  │
│  │ Auth Middleware│  │ Rate Limiter │  │ Security Headers │  │
│  └───────────────┘  └──────────────┘  └──────────────────┘  │
└─────────────────────────────────────────────────────────────┘
                          │
          ┌───────────────┼───────────────┐
          ▼               ▼               ▼
┌─────────────────┐ ┌─────────────────┐ ┌─────────────────┐
│  Worker Service │ │ Export Service  │ │ Webhook Service │
│ (Task Schedule) │ │ (Data Export)   │ │ (Event Notify)  │
└─────────────────┘ └─────────────────┘ └─────────────────┘
          │               │               │
          └───────────────┼───────────────┘
                          ▼
┌─────────────────────────────────────────────────────────────┐
│                    Data Layer                                │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐    │
│  │  MySQL   │  │  Redis   │  │   Log    │  │  Audit   │    │
│  │ (Primary)│  │ (Cache)  │  │ (Log)    │  │ (Audit)  │    │
│  └──────────┘  └──────────┘  └──────────┘  └──────────┘    │
└─────────────────────────────────────────────────────────────┘
```

### Architecture Highlights

- **Frontend-Backend Separation**: Vue 3 SPA + Go backend API, clear responsibilities
- **Microservices**: Independent API, Worker, Export, and Webhook services
- **Containerized Deployment**: Docker multi-stage builds, images pushed to GHCR
- **Observability**: Prometheus metrics + Grafana dashboards
- **Security**: Multi-layer security protection, end-to-end encryption

---

## 🚀 Quick Start

### Prerequisites

- Docker Engine 20.10+
- Docker Compose 2.0+
- 4GB+ available memory
- 10GB+ available disk space

### One-Click Deployment

```bash
# 1. Clone repository
git clone https://github.com/xiaohai09/mobile-cloud-disk-management.git
cd mobile-cloud-disk-management

# 2. Copy environment configuration
cp .env.example .env
# Edit .env file, configure database password, JWT secret, etc.

# 3. Start services
docker compose up -d

# 4. Verify services
curl http://localhost/api/health
```

### Access Points

| Service | URL | Description |
|---------|-----|-------------|
| Frontend | http://localhost | User interface |
| API Service | http://localhost/api | RESTful API |
| Monitoring | http://localhost:3000 | Prometheus + Grafana |
| Health Check | http://localhost/api/health | Health status endpoint |

For detailed deployment instructions, see [Deployment Guide](docs/DEPLOYMENT.md).

---

## 🛠️ Tech Stack

### Frontend

| Technology | Purpose | Version |
|------------|---------|---------|
| Vue 3 | Progressive JavaScript framework | ^3.4.0 |
| TypeScript | Type-safe JavaScript | ^5.3.0 |
| Element Plus | Enterprise-grade UI component library | ^2.4.0 |
| Pinia | State management | ^2.1.0 |
| Vite | Next-generation frontend build tool | ^5.0.0 |
| Axios | HTTP client | ^1.6.0 |

### Backend

| Technology | Purpose | Version |
|------------|---------|---------|
| Go | Programming language | 1.25 |
| Gin | HTTP web framework | ^1.10.0 |
| GORM | ORM library | ^1.25.0 |
| Redis | Cache and session storage | 7.x |
| Gorilla WebSocket | WebSocket support | ^1.5.0 |
| Prometheus | Metrics collection | ^0.47.0 |

### Data Layer

| Technology | Purpose | Version |
|------------|---------|---------|
| MySQL | Primary database | 8.0 |
| Redis | Cache, queue, session | 7.x |

### Infrastructure

| Technology | Purpose |
|------------|---------|
| Docker | Containerization |
| Docker Compose | Multi-container orchestration |
| GitHub Actions | CI/CD |
| GHCR | Image registry |
| Nginx | Reverse proxy (built into frontend image) |

---

## 📁 Project Structure

```
mobile-cloud-disk-management/
├── backend/                    # Backend services
│   ├── cmd/                    # Executable entry points
│   │   ├── api/                # API service entry
│   │   └── worker/             # Worker service entry
│   ├── internal/               # Internal packages
│   │   ├── handler/            # HTTP handlers
│   │   ├── service/            # Business logic layer
│   │   ├── repository/         # Data access layer
│   │   ├── model/              # Data models
│   │   ├── middleware/         # Middleware
│   │   ├── config/             # Configuration management
│   │   └── utils/              # Utilities
│   ├── pkg/                    # Public packages
│   │   ├── auth/               # Authentication
│   │   ├── export/             # Export functionality
│   │   ├── webhook/            # Webhook functionality
│   │   └── security/           # Security utilities
│   ├── migrations/             # Database migrations
│   ├── scripts/                # Script utilities
│   ├── Dockerfile.api          # API service Dockerfile
│   ├── Dockerfile.worker       # Worker service Dockerfile
│   └── go.mod                  # Go module definition
├── frontend/                   # Frontend application
│   ├── src/
│   │   ├── views/              # Page components
│   │   ├── components/         # Shared components
│   │   ├── stores/             # Pinia state management
│   │   ├── router/             # Route configuration
│   │   ├── api/                # API request wrappers
│   │   ├── utils/              # Utilities
│   │   └── assets/             # Static assets
│   ├── public/                 # Public assets
│   ├── Dockerfile              # Frontend Dockerfile
│   ├── nginx.conf              # Nginx configuration
│   └── package.json            # Dependency definition
├── docs/                       # Project documentation
│   ├── DEPLOYMENT.md           # Deployment guide
│   ├── API.md                  # API documentation
│   ├── USER_GUIDE.md           # User guide
│   └── TEST_REPORT.md          # Test report
├── scripts/                    # Project-level scripts
│   └── start.sh                # One-click startup script
├── .github/workflows/          # GitHub Actions workflows
│   ├── ci.yml                  # Continuous integration
│   └── publish.yml             # Image publishing
├── docker-compose.yml          # Docker Compose configuration
├── .env.example                # Environment variable template
└── README.md                   # Project description (Chinese)
```

---

## ⚡ Development Mode

### Environment Requirements

- Go 1.25+
- Node.js 20+
- MySQL 8.0
- Redis 7.x

### Backend Development

```bash
cd backend

# Download dependencies
go mod download

# Start API service (hot reload)
go run ./cmd/api/main.go

# Start Worker service
go run ./cmd/worker/main.go

# Run tests
go test ./... -v -count=1

# Code check
go vet ./...
```

### Frontend Development

```bash
cd frontend

# Install dependencies
npm install

# Start development server
npm run dev

# Code check
npm run lint

# Type check
npm run typecheck

# Production build
npm run build
```

---

## 🧪 Testing

### Frontend Testing

```bash
cd frontend
npm run lint          # ESLint code check
npm run typecheck     # TypeScript type check
npm run build         # Production build verification
```

### Backend Testing

```bash
cd backend
go test ./... -v -count=1 -timeout=5m    # Run all unit tests
go test ./... -coverprofile=coverage.out  # Generate coverage report
go vet ./...                              # Static code analysis
```

### Security Audit

```bash
cd backend
go install github.com/securego/gosec/v2/cmd/gosec@latest
gosec -fmt=json -out=gosec-results.json ./...

go install golang.org/x/vuln/cmd/govulncheck@latest
govulncheck -json ./...
```

---

## 🔒 Security

### Security Features

- **Transport Security**: HTTPS/WSS encrypted communication
- **Authentication & Authorization**: JWT + unified middleware, role-based access control
- **Input Validation**: Strict validation of all input parameters, preventing injection attacks
- **Rate Limiting**: API request rate limiting, preventing brute force attacks
- **Security Headers**: CSP, HSTS, X-Frame-Options, etc.
- **WebSocket Security**: Strict origin validation, preventing unauthorized connections
- **Audit Logging**: Sanitized operation logs for security auditing

### Security Reporting

Found a security vulnerability? Please disclose responsibly via [GitHub Security Advisories](https://github.com/xiaohai09/mobile-cloud-disk-management/security/advisories).

For detailed security audit reports, see [Test Report](docs/TEST_REPORT.md).

---

## 📈 Monitoring and Observability

### Metrics Collection

- Prometheus metrics endpoint: `/metrics`
- Business metrics: request volume, latency, error rate
- System metrics: CPU, memory, disk, network

### Log Management

- Structured JSON logs
- Log levels: DEBUG, INFO, WARN, ERROR
- Audit logs: sanitized, complete operation chain preserved

### Alert Configuration

- Error rate threshold alerts
- Response time P99 alerts
- System resource alerts

---

## 🤝 Contributing

We welcome community contributions! Please follow this process:

1. Fork this repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'feat: add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Contribution Guidelines

- Follow [Conventional Commits](https://www.conventionalcommits.org/) specification
- Ensure all tests pass
- Update relevant documentation
- Merge after code review approval

---

## 🗺️ Roadmap

- [x] Infrastructure refactoring and Dockerization
- [x] Data export functionality (CSV/JSON)
- [x] Webhook notification system
- [x] Responsive frontend interface
- [x] Security hardening and audit logs
- [ ] Multi-platform adapter layer (reserved interfaces)
- [ ] Kubernetes Helm Chart deployment
- [ ] Internationalization (i18n) support
- [ ] Plugin system architecture
- [ ] Performance optimization and caching strategy

---

## ❓ Frequently Asked Questions

**Q: Which cloud disk platforms are supported?**

A: The current version supports mainstream cloud disk platform account management and automated operations. Multi-platform adapter interfaces are reserved for extension. Refer to [User Guide](docs/USER_GUIDE.md) for extension guidelines.

**Q: How to backup data?**

A: It is recommended to regularly back up MySQL database and Redis data. For specific backup strategies, see [Deployment Guide](docs/DEPLOYMENT.md).

**Q: Is cluster deployment supported?**

A: The current version is based on Docker Compose single-node deployment. Kubernetes deployment support is planned.

**Q: How to upgrade versions?**

A: Refer to the upgrade section in [Deployment Guide](docs/DEPLOYMENT.md), remember to back up data and follow migration guidelines.

For more questions, please refer to [User Guide](docs/USER_GUIDE.md).

---

## 📄 License

This project is open source under the [MIT License](LICENSE).

---

## 🙏 Acknowledgements

- Based on [3238614968/caiyun](https://github.com/3238614968/caiyun) architecture refactoring
- Thanks to all contributors for their support and feedback

---

## 📞 Contact

- **Issues**: [GitHub Issues](https://github.com/xiaohai09/mobile-cloud-disk-management/issues)
- **Discussions**: [GitHub Discussions](https://github.com/xiaohai09/mobile-cloud-disk-management/discussions)

---

*Last updated: 2026-07-07*
