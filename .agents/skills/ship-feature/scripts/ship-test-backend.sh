#!/usr/bin/env bash
# Backend API testing via SSH into the deployment server.
# Usage: ship-test-backend.sh <ssh-user>@<host> [backend-port]
#   backend-port defaults to BACKEND_PORT from GitHub env vars
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
source "$SCRIPT_DIR/ship-env.sh"

SSH_TARGET="$1"
BACKEND_PORT="${2:-${BACKEND_PORT:-8123}}"

echo "=== Backend API Testing ==="
echo "Server: $SSH_TARGET"
echo "Backend port: $BACKEND_PORT"
echo ""

# ── Health check ───────────────────────────────────────────────────
echo "→ Health check..."
ssh "$SSH_TARGET" "curl -s http://localhost:${BACKEND_PORT}/api/health"
echo ""

# ── List photos ────────────────────────────────────────────────────
echo "→ List photos..."
ssh "$SSH_TARGET" "curl -s http://localhost:${BACKEND_PORT}/api/photos | head -c 500"
echo ""

# ── List projects ──────────────────────────────────────────────────
echo "→ List projects..."
ssh "$SSH_TARGET" "curl -s http://localhost:${BACKEND_PORT}/api/projects | head -c 500"
echo ""

# ── Landing card ───────────────────────────────────────────────────
echo "→ Landing card..."
ssh "$SSH_TARGET" "curl -s http://localhost:${BACKEND_PORT}/api/landingcard | head -c 500"
echo ""

echo "Done."
