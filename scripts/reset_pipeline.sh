#!/usr/bin/env bash
set -euo pipefail

PROJECT="docflow"

echo "🧹 Resetting Redis queues..."
docker exec ${PROJECT}-redis redis-cli FLUSHALL >/dev/null

echo "🧹 Clearing job retries..."
docker exec -i ${PROJECT}-postgres psql -U postgres <<'EOF'
UPDATE jobs SET retry_count = 0;
EOF

echo "Pipeline reset complete."