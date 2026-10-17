#!/usr/bin/env bash
set -e

BASE_URL="${BASE_URL:-http://localhost:8080}"

KEY="test-key"
NAMESPACE="test"
SCOPE="default"

echo "=== SET CACHE ==="

curl -i -X POST "$BASE_URL/v1/cache" \
  -H "Content-Type: application/json" \
  -H "x-cache-key: $KEY" \
  -H "x-cache-namespace: $NAMESPACE" \
  -H "x-cache-scope: $SCOPE" \
  -H "x-cache-ttl: 300" \
  -d '{
    "name": "XFSC",
    "message": "Hello Redis",
    "active": true
  }'

echo
echo "=== GET CACHE ==="

curl -i -X GET "$BASE_URL/v1/cache" \
  -H "x-cache-key: $KEY" \
  -H "x-cache-namespace: $NAMESPACE" \
  -H "x-cache-scope: $SCOPE"

echo