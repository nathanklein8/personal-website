# Debug Admin App - Investigation Status

## Problem
- Admin app renders components strangely
- Admin app has hydration issues, specifically with mode-toggle
- Changes made to tailwind CSS aren't always reflected in the dev server

## Findings So Far

### 1. Shared Library Definition (`packages/shared/package.json`)
**Status: PARTIALLY CORRECT**
- Package name: `@nk/shared` ✓
- Exports map defined for all components ✓
- Peer dependencies declared ✓
- **ISSUE**: The `index.ts` export is empty (just a comment). No barrel exports.
- **ISSUE**: `embla-carousel-svelte` is NOT listed as a peer dependency, but the carousel component imports it.

### 2. Admin App Configuration
**Status: CORRECT**
- `svelte.config.js`: Alias `@nk/shared` → `../packages/shared/src` ✓
- `vite.config.ts`: Alias `@nk/shared` → `../packages/shared/src` ✓
- `package.json`: Has `@nk/shared: "file:../packages/shared"` in devDependencies ✓
- Has all peer dependencies as regular dependencies ✓
- `embla-carousel-svelte` not in admin deps — not needed (admin doesn't use carousel) ✓

### 3. Main App Configuration
**Status: MOSTLY CORRECT**
- `svelte.config.js`: Alias `@nk/shared` → `./packages/shared/src` ✓
- `vite.config.ts`: Alias `@nk/shared` → `./packages/shared/src` ✓
- Root `package.json`: Has workspaces config ✓
- Has `@nk/shared: "file:packages/shared"` in devDependencies ✓

### 4. CSS / Tailwind CSS Issue
**Status: CONFIRMED ROOT CAUSE**
- Both apps use Tailwind CSS v4 (CSS-based config in `app.css`)
- Tailwind CSS v4 auto-scans for class names, but only within the project's source tree
- Shared library components in `packages/shared/src/components/ui/` contain Tailwind classes
- **CONFIRMED**: Tailwind v4 does NOT scan files outside the admin source tree when building admin
- Both `app.css` files are identical between main app and admin ✓

### 5. Hydration Warnings (Mode Toggle)
**Status: CONFIRMED ISSUE**
- `mode-toggle.svelte` imports `mode-watcher` which detects OS theme preference
- `ModeWatcher` component in `+layout.svelte` runs client-side only
- Server doesn't know the user's OS theme preference → mismatch between server-rendered HTML and client-hydrated DOM
- **ROOT CAUSE**: `mode-watcher` library inherently causes hydration mismatches when used server-side

### 6. Carousel Context Import Issue
**Status: BUG FOUND**
- `packages/shared/src/components/ui/carousel/context.ts` imports `$lib/utils.js`
- `$lib` is app-specific - it resolves differently in the main app vs admin app
- The shared library's carousel component should import `@nk/shared/utils` instead
- **NOTE**: Admin doesn't use the carousel, so this doesn't directly affect admin, but it's a bug for the main app

### 7. Dockerfile for Admin
**Status: CORRECT**
- Builds from root context to resolve workspace dependencies ✓
- Copies `packages/` and `admin/` ✓
- Runs `bun install` for workspace resolution ✓
- Copies built artifacts and packages to runner image ✓

### 8. Dockerfile for Main App
**Status: CORRECT**
- Similar approach, copies all workspace dirs ✓
- Copies `packages` to runner image ✓

## Key Suspects for CSS Issues

1. **Tailwind CSS v4 scanning - CONFIRMED ROOT CAUSE**:
   - Main app CSS output: **306 unique classes** (includes shared component classes like `bg-primary`, `text-primary-foreground`, `bg-card`, `rounded-xl`, `shadow-sm`, `border-border`)
   - Admin app CSS output: **51 unique classes** (MISSING all shared component classes)
   - Tailwind v4 auto-scans files only within the project's source tree
   - For the main app, `packages/shared/src` is within the project root
   - For the admin app, `packages/shared/src` is at `../../packages/shared/src` relative to `admin/src/` — **outside the admin source tree**
   - The symlinked package in `admin/node_modules/@nk/shared` is also ignored (Tailwind v4 skips `node_modules`)
   - `tailwind-variants` dynamic classes are NOT the issue — the problem is purely directory scanning scope

## Root Cause Summary

The shared library components (`button.svelte`, `card.svelte`, `dialog.svelte`, etc.) contain Tailwind CSS utility classes. When the main app builds, Tailwind v4 auto-scans the project and picks up these classes because `packages/shared/src` is within the project root.

When the admin app builds, Tailwind v4 does NOT scan `packages/shared/src` because:
1. The admin app is built from the `admin/` directory
2. `packages/shared/src` is at `../../packages/shared/src` relative to `admin/src/` — outside the admin source tree
3. The symlinked package in `admin/node_modules/@nk/shared` points to `node_modules`, which Tailwind v4 explicitly ignores (via `.gitignore` heuristics)

The `@source` directive in CSS can tell Tailwind v4 to scan additional directories.

## Next Steps

1. **Add `@source` directive to admin's `app.css`** to scan `../../packages/shared/src`
2. **Fix the `$lib/utils.js` import in carousel context.ts** — should use `@nk/shared/utils` instead
3. **Fix hydration warning in mode-toggle** — defer client-side logic with `onMount` or `client:only`
4. **Consider adding `embla-carousel-svelte` as a peer dependency** in shared package (currently missing)
5. **Consider adding `@source` to main app's `app.css`** as well for consistency, even though it works currently
