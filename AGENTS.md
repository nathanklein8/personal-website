# Personal Website Project - Agent Guide

## Project Overview

A personal portfolio website built with **SvelteKit** (frontend) and **Go** (backend), deployed via Docker Compose. Features a landing page, projects showcase, photography gallery, and admin panel for content management.

---

## Tech Stack

### Frontend
- **SvelteKit 2.x** - Full-stack framework
- **Svelte 5** - Runtime with runes
- **TypeScript** - Type safety
- **Tailwind CSS 4** - Utility-first styling
- **Shadcn/ui** - UI component library (custom implementation)
- **Lucide icons** - Icon set
- **Embla Carousel** - Photo carousel

### Backend
- **Go 1.22+** - Backend API
- **chi router** - HTTP router
- **PostgreSQL 18** - Database (via pgx driver)
- **Atlas** - Database schema management
- **imaging** - Image resizing for thumbnails/medium previews
- **goexif** - EXIF metadata extraction from photos
- **Layered Architecture** - Model/Repository/Service/Controller pattern

### DevOps
- **Docker Compose** - Container orchestration
- **Multi-stage Docker builds** - Optimized images

---

## Project Structure

```
.
├── admin/                 # Admin panel (SvelteKit)
│   ├── src/
│   │   ├── routes/       # Admin pages (landing, projects, photos)
│   │   ├── lib/
│   │   │   ├── components/
│   │   │   ├── server/
│   │   │   └── utils.ts
│   └── vite.config.ts
├── backend/              # Go API server
│   ├── models/          # Data structures with JSON tags
│   │   ├── project.go
│   │   ├── photo.go
│   │   └── landingcard.go
│   ├── repository/      # Database access layer (SQL queries)
│   │   ├── project_repository.go
│   │   ├── photo_repository.go
│   │   └── landingcard_repository.go
│   ├── service/         # Business logic & validation
│   │   ├── errors.go
│   │   ├── project_service.go
│   │   ├── photo_service.go
│   │   └── landingcard_service.go
│   ├── routes/          # HTTP controllers
│   │   ├── health.go
│   │   ├── landingcard.go
│   │   ├── photos.go
│   │   └── projects.go
│   ├── server/          # Database connection
│   ├── migrations/      # Atlas migrations
│   └── schema.hcl       # Database schema
├── photos/              # Photo library (mounted to container)
│   ├── {year}/          # Year directories
│   │   └── {event}/     # Event directories
│   │       └── *.jpg    # Photo files
├── src/                 # Main app (SvelteKit)
│   ├── routes/         # Public pages (landing, projects, photography)
│   ├── lib/
│   │   ├── components/ # Shared components
│   │   ├── server/     # Backend client
│   │   └── utils.ts
├── compose.yaml         # Docker Compose configuration
└── .env.example         # Environment template
```

---

## Environment Variables

```bash
ENV=dev                    # Environment: dev, test, prod
APP_PORT=3000             # Frontend port
BACKEND_PORT=8123         # Backend API port
POSTGRES_USER=user        # Database user
POSTGRES_PASSWORD=password
POSTGRES_DB=devdb         # Database name
POSTGRES_PORT=5432        # Database port
API_URL=http://localhost:${BACKEND_PORT}
DATABASE_URL=...          # Full connection string
PHOTO_LIBRARY=/photos     # Local path mounted to /photos in container (read-only)
```

---

## Database Schema

### landing_card (single row, id=1)
- bio (text)
- email (text)
- linkedin (text)
- github (text)
- skills (jsonb) - array of [category, skill] pairs

### projects
- id (serial, pk)
- icon (text, not null)
- title (text, not null)
- description (text, not null)
- technologies (text) - JSON array stored as text
- deployment_link (text, nullable)
- image (text, nullable)
- alt_text (text, nullable)

### photos
- id (serial, pk)
- title (text, not null)
- file_path (text, not null) - just the filename, e.g. "IMGP5178.jpg"
- alt_text (text, nullable)
- date_taken (text, nullable) - extracted from EXIF, format "YYYY-MM-DD"
- location (text, nullable) - GPS coordinates "lat, lon"
- camera (text, nullable) - extracted from EXIF
- lens (text, nullable) - extracted from EXIF
- aperture (text, nullable) - extracted from EXIF, format "f/X.X"
- shutter_speed (text, nullable) - extracted from EXIF, format "1/Xs" or "Xs"
- iso (text, nullable) - extracted from EXIF
- visible (boolean, default true)
- sort_order (integer, default 0)
- source_path (text, not null) - full relative path in photo library, e.g. "2026/1:15 Cobb Hill Park/IMGP5178.jpg"
- thumbnail_path (text, nullable) - relative path in thumbnails volume, e.g. "2026/1:15 Cobb Hill Park/IMGP5178.jpg_thumb.jpg"
- medium_path (text, nullable) - relative path in thumbnails volume, e.g. "2026/1:15 Cobb Hill Park/IMGP5178.jpg_med.jpg"

---

## API Endpoints

### Landing Card
- `GET /api/landingcard` - Get bio, social links, skills
- `POST/PUT /api/landingcard` - Update (id=1)

### Projects
- `GET /api/projects` - List all projects
- `POST/PUT /api/projects` - Create/update (id=0 inserts, non-zero updates)
- `DELETE /api/projects/{id}` - Delete

### Photos - File Explorer Browsing
- `GET /api/photos/available` - List year directories from photo library
- `GET /api/photos/available/{year}` - List event directories within a year
- `GET /api/photos/available/{year}/{event}` - List .jpg filenames in an event
- `GET /api/photos/available/{year}/{event}/{filename}/preview` - Serve thumbnail or medium preview (auto-generates if missing)

### Photos - CRUD
- `GET /api/photos` - List all photos in database (ordered by sort_order)
- `POST /api/photos` - Add photo from file explorer (accepts filename + title + sortOrder; reads file, extracts EXIF, generates thumbnails, inserts DB)
- `PUT /api/photos/{id}` - Update title, sortOrder, visible
- `DELETE /api/photos/{id}` - Delete photo record and generated thumbnails
- `POST /api/photos/regenerate-thumbnails` - Iterate all DB entries, regenerate all thumbnails and medium previews

**Important:** All routes are prefixed with `/api/photos` via `r.Mount("/api/photos", routes.PhotoRoutes(s))` in `main.go`.

---

## Photo System Architecture

### Directory Structure
```
photos/                          # Local photo library
├── 2025/                        # Year
│   ├── 7:11 Killiney/           # Event
│   │   ├── DSCN0390.NRW
│   │   └── DSCN0390_EDIT.jpg
│   └── 7:13 Phoenix Park/       # Event
│       └── IMGP4647.DNG
└── 2026/
    └── 1:15 Cobb Hill Park/     # Event
        └── IMGP5178.jpg
```

### Docker Volumes
| Volume | Mount | Purpose |
|--------|-------|---------|
| `PHOTO_LIBRARY` | `/photos:ro` | Read-only mount of local photo library |
| `photos-thumbnails` | `/thumbnails` | Read-write volume for generated thumbnails/medium previews |

### Thumbnail Generation
When adding a photo (`POST /api/photos`), the system:
1. Reads the source file from `/photos/{year}/{event}/{filename}`
2. Extracts EXIF metadata (date, camera, lens, aperture, shutter speed, ISO, GPS)
3. Generates a thumbnail (200px wide) at `/thumbnails/{year}/{event}/{filename}_thumb.jpg`
4. Generates a medium preview (800px wide) at `/thumbnails/{year}/{event}/{filename}_med.jpg`
5. Inserts a DB row with all metadata and paths

### Preview Serving
The preview endpoint (`GET /available/{year}/{event}/{filename}/preview`) auto-generates thumbnails on first request if they don't exist. The filename determines size:
- `filename.jpg` → medium preview (default)
- `filename.jpg_thumb.jpg` → thumbnail
- `filename.jpg_med.jpg` → medium preview

---

## Key Files to Know

### Frontend Routes
- `src/routes/+page.svelte` - Main landing page
- `src/routes/+page.server.ts` - Server-side data loading
- `src/routes/projects/+page.svelte` - Projects page
- `src/routes/photography/+page.svelte` - Photo gallery
- `src/routes/hiking/+page.svelte` - Placeholder (under construction)

### Shared Components
- `src/lib/components/landing-card.svelte` - Bio/social display
- `src/lib/components/project-card.svelte` - Project showcase
- `src/lib/components/photo-info.svelte` - Photo details dialog
- `src/lib/components/error-card.svelte` - Error display
- `src/lib/components/Header.svelte` - Navigation

### UI Components (`src/lib/components/ui/`)
- button, card, dialog, carousel, sheet, tooltip, scroll-area, separator, badge, alert, dropdown-menu

### Backend Models
- `backend/models/photo.go` - Photo, AddPhotoRequest, PhotoUpdateRequest, DirectoryEntry
- `backend/models/project.go` - Project data structure
- `backend/models/landingcard.go` - Landing card data structure

### Backend Repository
- `backend/repository/photo_repository.go` - Photo SQL queries + thumbnail path updates

### Backend Service
- `backend/service/photo_service.go` - Directory browsing, EXIF extraction, thumbnail generation, CRUD
- `backend/service/errors.go` - APIError type definition

### Backend Routes
- `backend/routes/photos.go` - All photo API routes (file explorer + CRUD)

---

## Backend Architecture

The backend follows a **Model → Repository → Service → Controller** pattern:

| Layer | Files | Responsibility |
|-------|-------|----------------|
| **Model** | `backend/models/*.go` | Data structures with JSON tags, no business logic |
| **Repository** | `backend/repository/*.go` | Database access layer with SQL queries, returns models or `sql.ErrNoRows` |
| **Service** | `backend/service/*.go` | Business logic, validation, file I/O, returns `APIError` for validation failures |
| **Controller** | `backend/routes/*.go` | HTTP routing, request/response handling, uses service layer |

### APIError Type
```go
type APIError struct {
    Status  int    `json:"-"`
    Message string `json:"message"`
}
```
Used by service layer for validation errors; controllers check for this type and return appropriate HTTP status codes.

---

## Development Commands

```bash
# Start dev servers
npm run dev              # Frontend + backend (via Docker)

# Database operations
npm run db:start         # Start PostgreSQL container
npm run db:diff          # Generate migration SQL
npm run db:apply         # Apply migrations

# Code quality
npm run check            # Type check
npm run lint             # Prettier + ESLint
npm run format           # Format code

# Build
npm run build            # Production build
npm run preview          # Preview production build
```

### Admin Panel
```bash
cd admin
npm run dev              # Admin panel dev server
```

### Docker Compose
```bash
docker compose up -d     # Start all services
docker compose down      # Stop all services
docker compose logs -f   # View logs
```

### Docker Volume Management
```bash
# Access thumbnail volume
docker compose exec backend ls -R /thumbnails/

# Access photo library (read-only)
docker compose exec backend ls -R /photos/
```

---

## Atlas Migration Gotchas

- **Checksum errors after deleting migration files**: If you delete a migration file, run `cd backend && atlas migrate hash` to re-hash the directory before running `db:diff`.
- **Dev database vs actual database**: `db:diff` uses `--dev-url "docker://postgres/18/dev"` which spins up a temporary dev database. If you drop a table from the actual DB but the dev database still has it, the diff will generate `ALTER TABLE` instead of `CREATE TABLE`. Fix by dropping from dev DB or rewriting the migration manually.
- **NOT NULL column additions**: Adding a `NOT NULL` column to an existing table with data requires either providing a default value or ensuring the table is empty.

---

## Testing & Debugging

### Testing the Backend
```bash
# Health check
curl http://localhost:8123/api/health

# List years
curl http://localhost:8123/api/photos/available

# List events in a year
curl "http://localhost:8123/api/photos/available/2026"

# List photos in an event
curl "http://localhost:8123/api/photos/available/2026/1:15%20Cobb%20Hill%20Park"

# Add a photo
curl -X POST http://localhost:8123/api/photos \
  -H "Content-Type: application/json" \
  -d '{"filename":"2026/1:15 Cobb Hill Park/IMGP5178.jpg","title":"Test","sortOrder":1}'

# View all photos
curl http://localhost:8123/api/photos

# Preview image
curl -o /dev/null -w "%{http_code} %{size_download}" \
  "http://localhost:8123/api/photos/available/2026/1:15%20Cobb%20Hill%20Park/IMGP5178.jpg/preview"
```

### Viewing Logs
```bash
docker compose logs -f backend
docker compose logs -f app
docker compose logs -f db
```

### Database Access
```bash
docker compose exec db psql -U user -d devdb
```

---

## Notes

- **Photo library**: Local `photos/` directory mounted read-only to `/photos` in the backend container. Structured as `year/event/filename.jpg`.
- **Thumbnails**: Generated on-demand when adding photos or previewing. Stored in a separate Docker volume `photos-thumbnails` at `/thumbnails`.
- **EXIF extraction**: Automatic on photo add. Camera model, lens, aperture, shutter speed, ISO, date, and GPS coordinates are extracted. Some older cameras (e.g. COOLPIX P340) may have aperture parsing issues.
- **Authentication**: Admin panel currently unauthenticated.
- **Database**: PostgreSQL 18 with Atlas schema management.
- **Styling**: Tailwind CSS 4 with OKLCH color space, dark mode via CSS custom properties.
- **Image sizes**: Thumbnails are 200px wide, medium previews are 800px wide, both using Lanczos resampling.
