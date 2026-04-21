#!/usr/bin/env bash
# Wait for GitHub Actions workflow to complete and show results.
# Auto-detects which workflow(s) to watch based on changed files.
# Usage: ship-wait.sh <branch-name>
set -euo pipefail

BRANCH="$1"
WORKFLOWS=("frontend.yml" "backend.yml" "admin.yml" "database.yml")

# ── Detect which workflows to watch ────────────────────────────────
changed_files=$(git diff --name-only origin/test...HEAD 2>/dev/null || echo "")

watch_frontend=false
watch_backend=false
watch_admin=false
watch_database=false

if echo "$changed_files" | grep -qE '^(src/|static/|vite\.config\.|package\.json|\.env)'; then
  watch_frontend=true
fi
if echo "$changed_files" | grep -qE '^(backend/)'; then
  watch_backend=true
fi
if echo "$changed_files" | grep -qE '^(admin/)'; then
  watch_admin=true
fi
if echo "$changed_files" | grep -qE '^(backend/schema\.hcl|backend/migrations/)'; then
  watch_database=true
fi

# If nothing detected, watch all
if ! $watch_frontend && ! $watch_backend && ! $watch_admin && ! $watch_database; then
  watch_frontend=true
  watch_backend=true
  watch_admin=true
  watch_database=true
fi

# ── Wait for each relevant workflow ────────────────────────────────
declare -A run_ids

for wf in "${WORKFLOWS[@]}"; do
  case "$wf" in
    frontend.yml)    $watch_frontend || continue ;;
    backend.yml)     $watch_backend || continue ;;
    admin.yml)       $watch_admin || continue ;;
    database.yml)    $watch_database || continue ;;
  esac

  echo "Watching $wf..."
  run_id=$(gh run list --workflow "$wf" --branch "$BRANCH" --status in_progress --json databaseId --jq '.[0].databaseId // empty' 2>/dev/null)

  if [[ -z "$run_id" ]]; then
    echo "⚠️  No in-progress run found for $wf"
    continue
  fi

  echo "  Run $run_id started..."
  run_ids["$wf"]="$run_id"

  # Wait for completion
  gh run watch "$run_id" --exit-status 2>/dev/null || true
done

# ── Show results ───────────────────────────────────────────────────
echo ""
echo "=== Results ==="
for wf in "${!run_ids[@]}"; do
  run_id="${run_ids[$wf]}"
  status=$(gh run view "$run_id" --json status --jq '.status' 2>/dev/null || echo "unknown")
  conclusion=$(gh run view "$run_id" --json conclusion --jq '.conclusion' 2>/dev/null || echo "unknown")
  echo "$wf: $status / $conclusion"
  if [[ "$conclusion" == "failure" ]]; then
    echo "  Failed steps:"
    gh run view "$run_id" --log-failed 2>/dev/null || true
  fi
done
