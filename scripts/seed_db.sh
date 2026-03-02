#!/usr/bin/env bash
set -euo pipefail

PROJECT="DocNebula"

echo "🌱 Seeding database with sample job..."

docker exec -i ${PROJECT}-postgres psql -U postgres <<'EOF'
INSERT INTO jobs (id, status, retry_count, idempotency_key)
VALUES ('seed-job-1', 'UPLOADED', 0, 'seed-key-1')
ON CONFLICT (idempotency_key) DO NOTHING;
EOF

echo "Seed complete."