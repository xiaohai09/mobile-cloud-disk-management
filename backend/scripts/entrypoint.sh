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

echo "[entrypoint] Entrypoint completed, starting API server..."
echo "=========================================="
echo "  Mobile Cloud Disk Management System"
echo "=========================================="
echo "Frontend:  http://localhost"
echo "API:       http://localhost:8080"
echo "API Health: http://localhost:8080/health"
echo "Grafana:   http://localhost:3000"
echo "=========================================="

# Execute the main command
exec "$@"
