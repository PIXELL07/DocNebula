#!/usr/bin/env bash
set -euo pipefail

PROJECT="docflow"
THRESHOLD_SEC=${THRESHOLD_SEC:-30}

echo "Checking worker heartbeats (threshold=${THRESHOLD_SEC}s)..."

OUT=$(docker exec -i ${PROJECT}-postgres psql -U postgres -At <<'EOF'
SELECT worker_id,
       EXTRACT(EPOCH FROM (NOW() - last_seen))::INT AS seconds_since_seen
FROM worker_heartbeats;
EOF
)

if [ -z "$OUT" ]; then
  echo "⚠️  No workers reporting heartbeats."
  exit 0
fi

while IFS='|' read -r wid secs; do
  if [ "$secs" -gt "$THRESHOLD_SEC" ]; then
    echo "❌ $wid is STALE (${secs}s)"
  else
    echo "✅ $wid healthy (${secs}s)"
  fi
done <<< "$OUT"