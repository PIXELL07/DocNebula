#!/usr/bin/env bash
set -euo pipefail

PROJECT="DocNebula"

echo "🚀 Starting DocFlow local stack..."
docker compose -f deployments/docker-compose.yml up -d

echo "⏳ Waiting for Postgres to be ready..."
until docker exec ${PROJECT}-postgres pg_isready -U postgres >/dev/null 2>&1; do
  sleep 2
done

echo "🗄 Running DB migrations..."
bash scripts/migrate_db.sh

echo "Local environment ready."