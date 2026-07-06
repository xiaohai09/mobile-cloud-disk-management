#!/usr/bin/env bash
set -euo pipefail

# caiyun one-click startup script
# Usage: ./scripts/start.sh [dev|prod]
#   dev  - start with docker compose (local build)
#   prod - pull prebuilt images and start

MODE="${1:-dev}"
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
cd "$PROJECT_ROOT"

echo "=========================================="
echo "  caiyun startup script"
echo "  Mode: $MODE"
echo "  Root: $PROJECT_ROOT"
echo "=========================================="

# Check Docker
if ! command -v docker >/dev/null 2>&1; then
    echo "ERROR: Docker is not installed" >&2
    exit 1
fi

if ! docker compose version >/dev/null 2>&1; then
    echo "ERROR: Docker Compose is not installed" >&2
    exit 1
fi

# Setup environment file
if [ ! -f .env ]; then
    echo "Creating .env from .env.example..."
    cp .env.example .env

    # Generate strong random passwords
    RANDOM_PASS=$(openssl rand -base64 32 | tr -d '/+=' | head -c 32)
    RANDOM_JWT=$(openssl rand -base64 32 | tr -d '/+=' | head -c 32)
    RANDOM_WORKER=$(openssl rand -base64 32 | tr -d '/+=' | head -c 32)
    RANDOM_GRAFANA=$(openssl rand -base64 24 | tr -d '/+=' | head -c 24)

    # Replace placeholders
    if [[ "$(uname)" == "Darwin" ]]; then
        # macOS
        sed -i '' "s/replace_with_strong_root_password/${RANDOM_PASS}/" .env
        sed -i '' "s/replace_with_strong_app_password/${RANDOM_PASS}/" .env
        sed -i '' "s/replace_with_strong_redis_password/${RANDOM_PASS}/" .env
        sed -i '' "s/replace_with_32_chars_random_secret/${RANDOM_JWT}/" .env
        sed -i '' "s/replace_with_random_monitor_token/${RANDOM_WORKER}/" .env
        sed -i '' "s/replace_with_strong_grafana_password/${RANDOM_GRAFANA}/" .env
    else
        # Linux
        sed -i "s/replace_with_strong_root_password/${RANDOM_PASS}/" .env
        sed -i "s/replace_with_strong_app_password/${RANDOM_PASS}/" .env
        sed -i "s/replace_with_strong_redis_password/${RANDOM_PASS}/" .env
        sed -i "s/replace_with_32_chars_random_secret/${RANDOM_JWT}/" .env
        sed -i "s/replace_with_random_monitor_token/${RANDOM_WORKER}/" .env
        sed -i "s/replace_with_strong_grafana_password/${RANDOM_GRAFANA}/" .env
    fi

    echo "✓ .env created with strong passwords"
else
    echo "✓ .env already exists"
fi

# Build or pull images
if [ "$MODE" = "prod" ]; then
    echo "Pulling production images..."
    docker compose pull
else
    echo "Building images locally..."
    docker compose build --no-cache
fi

# Start services
echo "Starting caiyun services..."
docker compose up -d

# Wait for services to be ready
echo "Waiting for services to be healthy..."
sleep 5

# Check health
echo ""
echo "=========================================="
echo "  Service Status"
echo "=========================================="
docker compose ps --format "table {{.Name}}\t{{.Status}}\t{{.Ports}}"

echo ""
echo "=========================================="
echo "  Access URLs"
echo "=========================================="
echo "Frontend:  http://localhost"
echo "API:       http://localhost:8080"
echo "API Health: http://localhost:8080/health"
echo "Grafana:   http://localhost:3000"
echo ""

# Try to access frontend
if curl -sf http://localhost/ >/dev/null 2>&1; then
    echo "✓ Frontend is accessible"
else
    echo "⚠ Frontend is not yet accessible (may need a few more seconds)"
fi

if curl -sf http://localhost:8080/health >/dev/null 2>&1; then
    echo "✓ API is healthy"
else
    echo "⚠ API is not yet healthy (may need a few more seconds)"
fi

echo ""
echo "=========================================="
echo "  Next Steps"
echo "=========================================="
echo "1. Visit http://localhost to access the frontend"
echo "2. Register an admin account via the frontend"
echo "3. Check logs: docker compose logs -f [service]"
echo "4. Stop:     docker compose down"
echo "=========================================="
