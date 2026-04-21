#!/usr/bin/env bash
# Load environment variables from GitHub repository/environment settings.
# Usage: source ship-env.sh [test|prod]
#
# Exports: APP_PORT, BACKEND_PORT, POSTGRES_PORT, etc.
#
# Priority (highest first):
#   1. Environment-scoped vars (GH_ENV_NAME/variables)
#   2. Repository-level vars (GH_REPO/variables)
#   3. .env file (for local dev fallback)
#   4. .env.example (bare minimum fallback)

set -euo pipefail

ENV_NAME="${1:-test}"
REPO="nathanklein8/personal-website"
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/../../../../.." && pwd)"

_load_github_vars() {
  local scope="$1"
  local url=""

  if [[ "$scope" == "env" ]]; then
    url="repos/$REPO/environments/$ENV_NAME/variables"
  else
    url="repos/$REPO/variables"
  fi

  local vars
  vars=$(gh api "$url" --jq '.variables[] | "\(.name)=\(.value)"' 2>/dev/null) || return 0

  while IFS= read -r line; do
    [[ -z "$line" ]] && continue
    local key="${line%%=*}"
    local val="${line#*=}"
    # Only set if not already set (allow overrides from env or .env)
    [[ -z "${!key+x}" ]] && export "$key=$val"
  done <<< "$vars"
}

_load_file() {
  local file="$1"
  [[ -f "$file" ]] || return 0
  while IFS= read -r line; do
    [[ "$line" =~ ^# ]] && continue
    [[ -z "$line" ]] && continue
    local key="${line%%=*}"
    local val="${line#*=}"
    val="${val#\"}"  # strip leading quote
    val="${val%\"}"  # strip trailing quote
    [[ -z "${!key+x}" ]] && export "$key=$val"
  done < "$file"
}

# Load order: GitHub env vars → GitHub repo vars → .env → .env.example
_load_github_vars "env"
_load_github_vars "repo"
_load_file "$REPO_ROOT/.env"
_load_file "$REPO_ROOT/.env.example"
