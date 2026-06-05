#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "$0")/.." && pwd)"

echo "[dev] root: $ROOT_DIR"

mkdir -p "$ROOT_DIR/backend/certs"

if [ ! -f "$ROOT_DIR/backend/certs/server.crt" ] || [ ! -f "$ROOT_DIR/backend/certs/server.key" ]; then
  echo "[dev] Generating self-signed TLS certs for gRPC (backend/certs)"
  openssl req -x509 -newkey rsa:4096 -nodes -keyout "$ROOT_DIR/backend/certs/server.key" -out "$ROOT_DIR/backend/certs/server.crt" -days 365 -subj "/CN=localhost"
fi

if command -v psql >/dev/null 2>&1; then
  echo "[dev] Applying DB schema (if available)."
  PSQL_URL="${DATABASE_URL:-postgres://gateway:gatewaypass@localhost:5432/gateway?sslmode=disable}"
  echo "[dev] Using DATABASE_URL=$PSQL_URL"
  psql "$PSQL_URL" -f "$ROOT_DIR/backend/internal/adapters/migration/schema.sql" || true
else
  echo "[dev] psql not found — skipping DB migration."
fi

export DATABASE_URL="${DATABASE_URL:-postgres://gateway:gatewaypass@localhost:5432/gateway?sslmode=disable}"
export JWT_SECRET="${JWT_SECRET:-changeme}"
export GRPC_TLS_CERT="${GRPC_TLS_CERT:-$ROOT_DIR/backend/certs/server.crt}"
export GRPC_TLS_KEY="${GRPC_TLS_KEY:-$ROOT_DIR/backend/certs/server.key}"

echo "[dev] Starting gateway (background)..."
(cd "$ROOT_DIR/backend" && go run ./cmd/gateway) &

echo "[dev] Starting frontend (foreground)..."
cd "$ROOT_DIR/frontend"
npm install --no-audit --no-fund
npm run dev
