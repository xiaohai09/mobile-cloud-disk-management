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
    RANDOM_PASS=$(openssl rand -hex 32)
    RANDOM_JWT=$(openssl rand -hex 32)
    RANDOM_WORKER=$(openssl rand -hex 32)
    RANDOM_GRAFANA=$(openssl rand -hex 24)

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
    echo "Pulling production images from GHCR..."
    IMAGE_TAG="${IMAGE_TAG:-latest}" docker compose -f docker-compose.prod.yml pull
else
    echo "Building images locally..."
    docker compose build --no-cache
fi

# Start services
echo "Starting caiyun services..."
if [ "$MODE" = "prod" ]; then
    IMAGE_TAG="${IMAGE_TAG:-latest}" docker compose -f docker-compose.prod.yml up -d
else
    docker compose up -d
fi

# Wait for services to be ready
echo "Waiting for services to be healthy..."
sleep 5

# Check health
echo ""
echo "=========================================="
echo "  Service Status"
echo "=========================================="
if [ "$MODE" = "prod" ]; then
    IMAGE_TAG="${IMAGE_TAG:-latest}" docker compose -f docker-compose.prod.yml ps --format "table {{.Name}}\t{{.Status}}\t{{.Ports}}"
else
    docker compose ps --format "table {{.Name}}\t{{.Status}}\t{{.Ports}}"
fi

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
echo "  Default Admin Account"
echo "=========================================="
echo "  Username: ${DEFAULT_ADMIN_USERNAME:-admin}"
echo "  Password: ${DEFAULT_ADMIN_PASSWORD:-admin123}"
echo ""
echo "Note: If the admin account already exists, the container"
echo "entrypoint will skip creation on subsequent starts."
echo "=========================================="

echo ""
echo "=========================================="
echo "  Next Steps"
echo "=========================================="
echo "1. Visit http://localhost to access the frontend"
echo "2. Login with default admin account:"
echo "   Username: ${DEFAULT_ADMIN_USERNAME:-admin}"
echo "   Password: ${DEFAULT_ADMIN_PASSWORD:-admin123}"
echo "3. Check logs:"
if [ "$MODE" = "prod" ]; then
    echo "   docker compose -f docker-compose.prod.yml logs -f [service]"
else
    echo "   docker compose logs -f [service]"
fi
echo "4. Stop:"
if [ "$MODE" = "prod" ]; then
    echo "   docker compose -f docker-compose.prod.yml down"
else
    echo "   docker compose down"
fi
echo "=========================================="
