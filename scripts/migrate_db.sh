#!/usr/bin/env bash
set -euo pipefail

PROJECT="DocNebula"

echo "🗄 Applying database migrations..."

docker exec -i ${PROJECT}-postgres psql -U postgres <<'EOF'
CREATE TABLE IF NOT EXISTS jobs (
    id TEXT PRIMARY KEY,
    status TEXT,
    retry_count INT DEFAULT 0,
    idempotency_key TEXT UNIQUE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS worker_heartbeats (
    worker_id TEXT PRIMARY KEY,
    last_seen TIMESTAMP
);
EOF

echo "Migration complete."