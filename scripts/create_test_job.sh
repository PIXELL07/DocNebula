#!/usr/bin/env bash
set -euo pipefail

API_URL=${API_URL:-"http://localhost:8080/upload"}

echo "📤 Creating test job..."
RESP=$(curl -s -X POST "$API_URL")

echo "$RESP" | jq . || echo "$RESP"

echo "✅ Test job submitted."