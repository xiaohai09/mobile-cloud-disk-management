#!/usr/bin/env bash
set -euo pipefail

echo "[entrypoint] Starting caiyun-api entrypoint..."

# Wait for MySQL to be ready
echo "[entrypoint] Waiting for MySQL..."
for i in $(seq 1 60); do
    if mysqladmin ping -h "${DB_HOST:-mysql}" -P "${DB_PORT:-3306}" \
        -u"${DB_USER:-caiyun_app}" -p"${DB_PASSWORD}" --silent; then
        echo "[entrypoint] MySQL is ready"
        break
    fi
    if [ "$i" -eq 60 ]; then
        echo "[entrypoint] ERROR: MySQL did not become ready in time" >&2
        exit 1
    fi
    sleep 2
done

# Wait for Redis to be ready
echo "[entrypoint] Waiting for Redis..."
for i in $(seq 1 60); do
    if redis-cli -h "${REDIS_HOST:-redis}" -p "${REDIS_PORT:-6379}" \
        -a"${REDIS_PASSWORD}" ping | grep -q PONG; then
        echo "[entrypoint] Redis is ready"
        break
    fi
    if [ "$i" -eq 60 ]; then
        echo "[entrypoint] ERROR: Redis did not become ready in time" >&2
        exit 1
    fi
    sleep 2
done

# Seed default admin user if not exists
ADMIN_USERNAME="${DEFAULT_ADMIN_USERNAME:-admin}"
ADMIN_PASSWORD="${DEFAULT_ADMIN_PASSWORD:-admin123}"

echo "[entrypoint] Checking if default admin user exists..."
USER_EXISTS=$(mysql -N -s -h "${DB_HOST:-mysql}" -P "${DB_PORT:-3306}" \
    -u"${DB_USER:-caiyun_app}" -p"${DB_PASSWORD}" "${DB_NAME:-caiyun}" \
    -e "SELECT COUNT(*) FROM users WHERE username='${ADMIN_USERNAME}';" 2>/dev/null || echo "0")

if [ "${USER_EXISTS}" = "0" ]; then
    echo "[entrypoint] Creating default admin user: ${ADMIN_USERNAME}"
    # bcrypt hash of admin123 (cost 12), verified against golang.org/x/crypto/bcrypt
    ADMIN_HASH='$2b$12$og/0Ll/lXUsCip/g6PcNy.vEZBHAUz8r0W3whcXCRYSGLj1Al4uj6'
    
    mysql -h "${DB_HOST:-mysql}" -P "${DB_PORT:-3306}" \
        -u"${DB_USER:-caiyun_app}" -p"${DB_PASSWORD}" "${DB_NAME:-caiyun}" \
        -e "INSERT INTO users (username, password, email, role, token_version, created_at, updated_at) VALUES ('${ADMIN_USERNAME}', '${ADMIN_HASH}', 'admin@example.com', 'admin', 0, NOW(), NOW()) ON DUPLICATE KEY UPDATE updated_at = NOW();"
    
    echo "[entrypoint] Default admin user created"
else
    echo "[entrypoint] Default admin user already exists, skipping"
fi

echo "[entrypoint] Entrypoint completed, starting API server..."
echo "=========================================="
echo "  Mobile Cloud Disk Management System"
echo "=========================================="
echo "Frontend:  http://localhost"
echo "API:       http://localhost:8080"
echo "API Health: http://localhost:8080/health"
echo "Grafana:   http://localhost:3000"
echo ""
echo "Default Admin Account:"
echo "  Username: ${ADMIN_USERNAME}"
echo "  Password: ${ADMIN_PASSWORD}"
echo "=========================================="

# Execute the main command
exec "$@"
