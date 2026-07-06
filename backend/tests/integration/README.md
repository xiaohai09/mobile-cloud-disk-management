# Integration Tests

This folder contains integration tests for the **caiyun** backend. Tests under
`tests/integration` exercise real MySQL and Redis containers via
[testcontainers-go](https://golang.testcontainers.org).

## Prerequisites

- Go 1.25+
- Docker Engine running and accessible from the test process
- Sufficient disk space for two test databases

## Running Integration Tests

Integration tests are **excluded from the default test build** so that they
cannot break local development or CI environments that do not have Docker.
To run them:

```bash
# Linux / macOS
export TEST_USER_PASSWORD="a-secure-test-password"
export JWT_SECRET="a-32-byte-or-longer-secret-for-testing-only"
export JWT_ALGORITHM="HS256"
export DB_HOST=127.0.0.1
export DB_PORT=3306
export REDIS_HOST=127.0.0.1
export REDIS_PORT=6379

# Run only integration tests:
go test -tags=integration ./tests/integration/...
```

### Windows (PowerShell)

```powershell
$env:TEST_USER_PASSWORD = "a-secure-test-password"
$env:JWT_SECRET = "a-32-byte-or-longer-secret-for-testing-only"
$env:JWT_ALGORITHM = "HS256"
$env:DB_HOST = "127.0.0.1"
$env:DB_PORT = "3306"
$env:REDIS_HOST = "127.0.0.1"
$env:REDIS_PORT = "6379"

go test -tags=integration ./tests/integration/...
```

## Environment Variables

| Variable              | Required | Default            | Description                              |
|------------------------|----------|--------------------|------------------------------------------|
| `JWT_ALGORITHM`        | yes      | —                  | `HS256` or `RS256`                       |
| `JWT_SECRET`           | yes*     | —                  | HMAC secret when `JWT_ALGORITHM=HS256`. Must be >= 32 bytes. |
| `JWT_PRIVATE_KEY`      | yes*     | —                  | RSA private key when `JWT_ALGORITHM=RS256` |
| `JWT_PUBLIC_KEY`       | yes*     | —                  | RSA public key when `JWT_ALGORITHM=RS256`  |
| `DB_HOST`              | no       | `localhost`         | MySQL host                               |
| `DB_PORT`              | no       | `3306`              | MySQL port                               |
| `DB_USER`              | no       | `root`              | MySQL user                                |
| `DB_PASSWORD`          | no       | *(empty)*           | MySQL password                            |
| `DB_NAME`              | no       | `caiyun_test`       | Database name                             |
| `REDIS_HOST`           | no       | `localhost`         | Redis host                                |
| `REDIS_PORT`           | no       | `6379`              | Redis port                                |
| `REDIS_PASSWORD`       | no       | *(empty)*           | Redis password                            |
| `TEST_USER_PASSWORD`   | yes      | —                   | Password used for the registered test user. **Never hardcode.** |

> *At least one JWT secret variant must be provided depending on the chosen
> algorithm.

## Test Scope

- `auth_flow_test.go` — register, login, token refresh, and logout against a
  real HTTP server wired to temporary MySQL and Redis containers.

## Safety Notes

- Integration tests spin up real containers; do **not** run them in production
  or restricted environments.
- Tests clean up containers and temporary databases on completion.
