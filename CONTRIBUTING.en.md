# Contributing Guide

Thank you for your interest in contributing to **Mobile Cloud Disk Management System**! This guide will help you get started.

## Code of Conduct

- Be respectful and inclusive to all participants
- Accept constructive criticism and focus on project improvement
- Reject harassment, discrimination, or inappropriate behavior
- Follow open source licenses and community standards

## How to Contribute

### Report a Bug

- Create an Issue using the [Bug Report](/.github/ISSUE_TEMPLATE/bug_report.md) template
- Provide a clear title and reproduction steps
- Include relevant logs, screenshots, and version info

### Request a Feature

- Create an Issue using the [Feature Request](/.github/ISSUE_TEMPLATE/feature_request.md) template
- Describe the value and expected behavior
- Include prototypes or mockups if available

### Submit Code

1. Fork this repository
2. Create a feature branch: `git checkout -b feature/amazing-feature`
3. Commit changes: `git commit -m 'feat: add amazing feature'`
4. Push branch: `git push origin feature/amazing-feature`
5. Open a Pull Request

## Development Standards

### Backend (Go)

- Follow Clean Architecture: `domain/application/infrastructure/interfaces`
- Format code with `gofmt`
- Commit message format: `type(scope): description`
- Types: `feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`

### Frontend (Vue3 + TypeScript)

- Format with ESLint + Prettier
- Component names use PascalCase, filenames use kebab-case
- Props/Events must have type definitions
- Commit message format same as backend

### Commit Message Convention

| Type | Description |
|------|-------------|
| `feat` | New feature |
| `fix` | Bug fix |
| `docs` | Documentation changes |
| `style` | Code formatting (no functional change) |
| `refactor` | Refactoring |
| `test` | Test related |
| `chore` | Build/tooling |

Example: `feat(auth): add JWT refresh token support`

## Development Environment

### Requirements

- Go 1.25+
- Node.js 20+
- Docker 20.10+
- Docker Compose v2+

### Local Setup

```bash
# Backend
cd backend
go mod download
go build ./cmd/api
go build ./cmd/worker

# Frontend
cd frontend
npm install
npm run dev
```

### Run Tests

```bash
# Backend tests
cd backend && go test ./...

# Frontend tests
cd frontend && npm run lint && npm run typecheck && npm run test:unit

# E2E tests
cd frontend && npm run e2e
```

## PR Checklist

- [ ] Code follows project standards
- [ ] Tests added/updated
- [ ] Documentation updated if needed
- [ ] All CI checks pass
- [ ] Major functionality manually tested

## Release Process

1. Update `CHANGELOG.md`
2. Create version tag: `git tag v1.0.0 && git push --tags`
3. GitHub Actions will automatically create a Release

## License

This project is licensed under MIT. By contributing, you agree to license your code under the same terms.

## Contact

- Issues: https://github.com/xiaohai09/mobile-cloud-disk-management/issues
- Discussions: https://github.com/xiaohai09/mobile-cloud-disk-management/discussions
