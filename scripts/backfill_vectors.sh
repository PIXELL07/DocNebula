#!/usr/bin/env bash
set -euo pipefail

PROJECT="docflow"
REDIS_CONTAINER="${PROJECT}-redis"

echo "🔁 Backfilling vector stage for completed OCR jobs..."

JOB_IDS=$(docker exec -i ${PROJECT}-postgres psql -U postgres -t -c \
  "SELECT id FROM jobs WHERE status='RUNNING' OR status='UPLOADED';")

for id in $JOB_IDS; do
  id=$(echo "$id" | xargs)
  [ -z "$id" ] && continue

  PAYLOAD=$(jq -n --arg id "$id" '{job_id:$id,attempt:0,ts:(now|todate)}')
  docker exec ${REDIS_CONTAINER} redis-cli LPUSH vector_queue "$PAYLOAD" >/dev/null
  echo "  ↳ queued $id"
done

echo "Backfill complete."