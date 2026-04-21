#!/usr/bin/env bash
# End-to-end ship feature: create branch, commit, merge, and wait for CI.
# Usage: ship-feature.sh <branch-name> [--to-prod]
set -euo pipefail

BRANCH="${1:?Usage: ship-feature.sh <branch-name> [--to-prod]}"
TO_PROD="${2:-}"
TEST_BRANCH="test"
PROD_BRANCH="main"

# ── 1. Create branch & commit ──────────────────────────────────────
git checkout -b "$BRANCH"
git add -A
git commit -m "$BRANCH"

# ── 2. Merge to test branch ────────────────────────────────────────
git checkout "$TEST_BRANCH"
git merge --no-ff "$BRANCH" -m "Merge $BRANCH into $TEST_BRANCH"
git push origin "$TEST_BRANCH"

if [[ "$TO_PROD" == "--to-prod" ]]; then
  git checkout "$PROD_BRANCH"
  git merge --no-ff "$BRANCH" -m "Merge $BRANCH into $PROD_BRANCH"
  git push origin "$PROD_BRANCH"
fi

# ── 3. Wait for CI ─────────────────────────────────────────────────
echo "Waiting for GitHub Actions..."
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
"$SCRIPT_DIR/ship-wait.sh" "$BRANCH"

echo "✅ Ship complete."
