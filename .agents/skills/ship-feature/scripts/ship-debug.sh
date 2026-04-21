#!/usr/bin/env bash
# Debugging helpers for failed deployments.
# Usage: ship-debug.sh <command> [args...]
#
# Commands:
#   logs <service> [filter]  - Show container logs (optionally filter for errors)
#   restart <service>        - Restart a single service
#   db-check                 - Verify database accessibility
#   db-query <sql>           - Run a raw SQL query against the database
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
source "$SCRIPT_DIR/ship-env.sh"

PROJECT="${DEPLOY_PROJECT_NAME:-pwsite-test}"
VALID_SERVICES="app backend admin db"

usage() {
  echo "Usage: ship-debug.sh <command> [args...]"
  echo ""
  echo "Commands:"
  echo "  logs <service> [filter]  - Show container logs"
  echo "  restart <service>        - Restart a single service"
  echo "  db-check                 - Verify database accessibility"
  echo "  db-query <sql>           - Run a raw SQL query"
  exit 1
}

is_valid_service() {
  local svc="$1"
  for valid in $VALID_SERVICES; do
    [[ "$svc" == "$valid" ]] && return 0
  done
  echo "Invalid service: $svc (valid: $VALID_SERVICES)"
  exit 1
}

case "${1:-}" in
  logs)
    SERVICE="${2:?Usage: ship-debug.sh logs <service> [filter]}"
    is_valid_service "$SERVICE"
    FILTER="${3:-}"
    if [[ -n "$FILTER" ]]; then
      echo "→ Logs for $SERVICE (filtered: $FILTER)"
      docker compose --project-name "$PROJECT" logs "$SERVICE" 2>&1 | grep -i "$FILTER"
    else
      echo "→ Logs for $SERVICE"
      docker compose --project-name "$PROJECT" logs -f "$SERVICE"
    fi
    ;;

  restart)
    SERVICE="${2:?Usage: ship-debug.sh restart <service>}"
    is_valid_service "$SERVICE"
    echo "→ Restarting $SERVICE..."
    docker compose --project-name "$PROJECT" up -d --no-deps "$SERVICE"
    echo "Done."
    ;;

  db-check)
    echo "→ Checking database..."
    docker compose --project-name "$PROJECT" exec db pg_isready
    ;;

  db-query)
    SQL="${2:?Usage: ship-debug.sh db-query <sql>}"
    echo "→ Running query..."
    docker compose --project-name "$PROJECT" exec db psql -U "$POSTGRES_USER" -d "$POSTGRES_DB" -c "$SQL"
    ;;

  *)
    usage
    ;;
esac
