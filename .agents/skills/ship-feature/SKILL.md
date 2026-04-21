---
name: ship-feature
description: End-to-end deployment workflow for this project. Use when pushing code changes to test or prod environments and verifying the deployment worked. Covers git workflow, deployment verification via gh CLI, Playwright UI testing, backend API testing, and Docker log inspection.
---

# Ship Feature

## Prerequisites

Before starting, ensure you have:
- `gh` CLI authenticated (`gh auth status`)
- SSH access to the deployment server (for backend testing and Docker logs)
- Playwright MCP configured in your MCP client (`@playwright/mcp@latest --headless --caps=testing,storage`)
- The test environment URL (ask the user if unknown)

## Workflow

### 1. Ship to test (or prod)

Run the automated ship script. It creates a branch, commits, merges, and waits for CI:

```bash
# Ship to test only (default)
.agents/skills/ship-feature/scripts/ship-feature.sh <branch-name>

# Ship to test AND prod
.agents/skills/ship-feature/scripts/ship-feature.sh <branch-name> --to-prod
```

The script:
- Creates a new branch with all current changes
- Commits with the branch name as the message (single line, no emojis)
- Merges into the test branch (and optionally main)
- Auto-detects which GitHub Actions workflows need to run based on changed files
- Waits for CI and reports results

### 2. Test changes

Run the appropriate test script based on what changed.

#### 2a. UI changes (Playwright MCP)

If the frontend or admin changed:

```bash
.agents/skills/ship-feature/scripts/ship-test-ui.sh <test-environment-url>
```

This prints the manual Playwright steps. Then:

1. Connect to the test environment via Playwright MCP
2. Navigate to the test environment URL
3. Take a snapshot: `browser_snapshot`
4. Test the specific changes
5. For critical flows, take screenshots: `browser_take_screenshot`

Key pages:
- Landing page: main URL
- Projects page: `/projects`
- Photography page: `/photography`
- Admin panel: admin URL (if admin changes)

Use the `testing` capability to add assertions:
- `browser_verify_element_visible`
- `browser_verify_text_visible`

#### 2b. Backend changes (curl via SSH)

```bash
.agents/skills/ship-feature/scripts/ship-test-backend.sh <ssh-user>@<host> [backend-port]
```

The backend port is auto-loaded from GitHub env vars (test/prod environment or repo-level). Pass a port explicitly to override.

This runs health check, list photos, list projects, and landing card endpoints against the server.

### 3. Debug failures

If CI failed or tests failed:

#### View failed CI runs

```bash
# List recent failures on a branch
gh run list --status failure --branch <branch>

# View a specific run
gh run view <run-id> --verbose
gh run view <run-id> --log-failed

# Rerun failed jobs
gh run rerun <run-id> --failed
```

#### Docker logs & restarts

```bash
# Show logs for a service (optionally filter for errors)
.agents/skills/ship-feature/scripts/ship-debug.sh logs <service> [filter]
# e.g.
.agents/skills/ship-feature/scripts/ship-debug.sh logs app error

# Restart a single service
.agents/skills/ship-feature/scripts/ship-debug.sh restart <service>
```

Valid services: `app` (frontend), `backend`, `admin`, `db`

#### Database checks

```bash
# Verify database is accessible
.agents/skills/ship-feature/scripts/ship-debug.sh db-check

# Run a raw SQL query
.agents/skills/ship-feature/scripts/ship-debug.sh db-query "SELECT count(*) FROM photos;"
```

#### Playwright debugging

If UI tests fail:
- `browser_take_screenshot` — capture page state
- `browser_console_messages` — check for JS errors
- `browser_network_requests` — check for failed API calls
- `browser_evaluate` — run diagnostic JS: `page.evaluate(() => document.body.innerHTML)`
