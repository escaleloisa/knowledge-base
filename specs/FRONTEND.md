# Web Frontend — Spec

## Overview

A Next.js (TypeScript + React) web application that serves as the primary user interface for the knowledge base. It connects to all existing backend services to provide a complete experience: browsing collections, viewing/editing markdown notes, full-text search, tag exploration, backlink navigation, and file import.

## Goals

- Provide a clean, demo-ready UI for the knowledge base
- Render markdown notes with full formatting and syntax highlighting
- Enable users to organize notes via collections
- Expose all backend capabilities: search, tags, backlinks
- Support importing `.md` files directly from the browser
- Deployable via Docker alongside the existing services

## Tech Stack

| Component | Technology |
|-----------|-----------|
| Framework | Next.js 14 (App Router) |
| Language | TypeScript |
| Styling | Tailwind CSS |
| Markdown rendering | react-markdown + remark-gfm + rehype-highlight |
| HTTP client | fetch (native) |
| State management | React Server Components + client hooks (no heavy state lib) |
| Containerization | Docker (multi-stage, standalone output) |

## Architecture

```mermaid
graph LR
    subgraph Frontend (Next.js - port 3000)
        direction TB
        SIDEBAR[Sidebar\nCollections + Tags]
        NOTELIST[Note List]
        VIEWER[Markdown Viewer]
        EDITOR[Note Editor]
        SEARCH[Search UI]
        IMPORT[File Import]
    end

    subgraph Backend Services
        NS[Note Service\nport 8080]
        SA[Search API\nport 8081]
        TA[Tag Aggregator\nport 8082]
    end

    SIDEBAR --> NS
    SIDEBAR --> TA
    NOTELIST --> NS
    VIEWER --> NS
    EDITOR --> NS
    SEARCH --> SA
    IMPORT --> NS
```

## Pages & Routes

| Route | Page | Description |
|-------|------|-------------|
| `/` | Home/Dashboard | Overview: recent notes, collection list, tag cloud |
| `/collections` | Collections List | All collections with note counts |
| `/collections/[id]` | Collection View | Notes within a specific collection |
| `/notes/[id]` | Note Viewer | Rendered markdown with backlinks panel |
| `/notes/[id]/edit` | Note Editor | Edit note content, tags, collections |
| `/notes/new` | New Note | Create a new note |
| `/search` | Search Results | Full-text search with filters |
| `/import` | Import | Upload .md files |

## Features

### 1. Collection Sidebar

- Lists all user-created collections with note counts
- Click to filter/view notes in that collection
- "New Collection" button — inline create
- Right-click or menu to rename/delete
- API: `GET /api/collections`, `POST /api/collections`

### 2. Markdown Viewer

- Renders note content as formatted HTML
- Supports GitHub Flavored Markdown (tables, task lists, strikethrough)
- Syntax highlighting for code blocks (all common languages)
- Renders `[[wiki-links]]` as clickable internal links to other notes
- Responsive layout with proper typography
- Libraries: `react-markdown`, `remark-gfm`, `rehype-highlight`

### 3. Note Editor

- Textarea/code editor for writing markdown
- Live preview (split pane or toggle)
- Edit title, content, and tags
- Assign/remove note from collections
- Save triggers `PUT /api/notes/:id` → Kafka event → all consumers update
- API: `PUT /api/notes/:id`, `GET /api/notes/:id/collections`, `POST /api/collections/:id/notes`

### 4. Full-Text Search

- Search bar in header (always accessible)
- Results page with highlighted matching content
- Filter by tags (multi-select)
- Filter by date range
- Relevance scoring displayed
- API: `GET /api/search?q=&tags=&from=&to=`

### 5. Tag Explorer

- Tag cloud or list showing all tags with counts
- Click a tag to filter notes by that tag
- Visual indication of tag popularity (size/color based on count)
- API: `GET /api/tags`

### 6. Backlinks Panel

- Shown on the note viewer page
- Lists all notes that link TO the current note via `[[wiki-links]]`
- Clickable to navigate to linking notes
- Shows the context/snippet where the link appears
- API: `GET /api/notes/:id/backlinks`

### 7. File Import

- Dedicated import page + drag-and-drop zone
- Accept single or multiple `.md` files
- Client-side parsing:
  - Extract title from first `# heading` or filename
  - File content becomes note content
  - Optional: detect tags from frontmatter (if present)
- Preview before import (show parsed title + content snippet)
- Let user select target collection for imported notes
- Progress indicator for bulk imports
- API: `POST /api/notes` (per file), then `POST /api/collections/:id/notes`

### 8. Note List

- Paginated list of notes (default: newest first)
- Shows title, first line of content, tags, date
- Filter by collection or tag
- Sort by date created, date updated, or title
- API: `GET /api/notes?limit=&offset=`, `GET /api/collections/:id/notes`

## Backend Service Integration

| Frontend Feature | Backend Service | Endpoint |
|-----------------|----------------|----------|
| Note CRUD | Note Service (8080) | `/api/notes`, `/api/notes/:id` |
| Collections | Note Service (8080) | `/api/collections`, `/api/collections/:id/notes` |
| Full-text search | Search API (8081) | `/api/search?q=&tags=&from=&to=` |
| Tag cloud/list | Tag Aggregator (8082) | `/api/tags` |
| Backlinks | Note Service (8080) | `/api/notes/:id/backlinks` |
| File import | Note Service (8080) | `POST /api/notes` |

## API Client Layer

Centralized API client to handle all backend communication:

```typescript
// lib/api.ts
const NOTE_SERVICE = process.env.NEXT_PUBLIC_NOTE_SERVICE_URL || 'http://localhost:8080';
const SEARCH_API = process.env.NEXT_PUBLIC_SEARCH_API_URL || 'http://localhost:8081';
const TAG_API = process.env.NEXT_PUBLIC_TAG_API_URL || 'http://localhost:8082';

export const api = {
  notes: {
    list: (limit?: number, offset?: number) => fetch(`${NOTE_SERVICE}/api/notes?limit=${limit}&offset=${offset}`),
    get: (id: string) => fetch(`${NOTE_SERVICE}/api/notes/${id}`),
    create: (data: CreateNoteRequest) => fetch(`${NOTE_SERVICE}/api/notes`, { method: 'POST', body: JSON.stringify(data) }),
    update: (id: string, data: UpdateNoteRequest) => fetch(`${NOTE_SERVICE}/api/notes/${id}`, { method: 'PUT', body: JSON.stringify(data) }),
    delete: (id: string) => fetch(`${NOTE_SERVICE}/api/notes/${id}`, { method: 'DELETE' }),
    backlinks: (id: string) => fetch(`${NOTE_SERVICE}/api/notes/${id}/backlinks`),
  },
  collections: {
    list: () => fetch(`${NOTE_SERVICE}/api/collections`),
    get: (id: string) => fetch(`${NOTE_SERVICE}/api/collections/${id}`),
    create: (data: CreateCollectionRequest) => fetch(`${NOTE_SERVICE}/api/collections`, { method: 'POST', body: JSON.stringify(data) }),
    update: (id: string, data: UpdateCollectionRequest) => fetch(`${NOTE_SERVICE}/api/collections/${id}`, { method: 'PUT', body: JSON.stringify(data) }),
    delete: (id: string) => fetch(`${NOTE_SERVICE}/api/collections/${id}`, { method: 'DELETE' }),
    notes: (id: string, limit?: number, offset?: number) => fetch(`${NOTE_SERVICE}/api/collections/${id}/notes?limit=${limit}&offset=${offset}`),
    addNotes: (id: string, noteIds: string[]) => fetch(`${NOTE_SERVICE}/api/collections/${id}/notes`, { method: 'POST', body: JSON.stringify({ note_ids: noteIds }) }),
    removeNote: (id: string, noteId: string) => fetch(`${NOTE_SERVICE}/api/collections/${id}/notes/${noteId}`, { method: 'DELETE' }),
  },
  search: {
    query: (q: string, tags?: string, from?: string, to?: string) => fetch(`${SEARCH_API}/api/search?q=${q}&tags=${tags || ''}&from=${from || ''}&to=${to || ''}`),
  },
  tags: {
    list: () => fetch(`${TAG_API}/api/tags`),
  },
};
```

## UI Layout

```
┌─────────────────────────────────────────────────────────────┐
│  Header: Logo / Search Bar / New Note Button                │
├────────────┬────────────────────────────────────────────────┤
│            │                                                │
│  Sidebar   │  Main Content Area                            │
│            │                                                │
│  Collections│  - Note Viewer (rendered markdown)            │
│  - Eng Docs │  - Note Editor (textarea + preview)          │
│  - Meeting  │  - Note List (paginated)                     │
│  - Personal │  - Search Results                            │
│            │  - Import UI                                   │
│  Tags      │                                                │
│  - go (12) │  ┌──────────────────────────────────────────┐ │
│  - k8s (8) │  │  Backlinks Panel (bottom or right side)  │ │
│  - docker  │  │  "Linked from: Note A, Note B"           │ │
│            │  └──────────────────────────────────────────┘ │
└────────────┴────────────────────────────────────────────────┘
```

## Wiki-Link Rendering

Notes can contain `[[note-title]]` references. The frontend should:

1. Parse `[[...]]` patterns in markdown content
2. Resolve them to actual note IDs (by title lookup or from backlinks data)
3. Render them as clickable `<Link>` components pointing to `/notes/[id]`
4. Style them distinctly (e.g., dotted underline, different color)
5. Show tooltip on hover with note title/preview

## File Import Flow

```
User Action                          System
───────────                          ──────
1. Navigate to /import
2. Drag .md file(s) into drop zone
3. Files read client-side             → FileReader API
4. Parse each file:
   - Title from first # or filename
   - Content = file body
   - Tags from YAML frontmatter (optional)
5. Show preview table                 → User reviews parsed notes
6. User selects target collection
7. Click "Import"
8. For each file:                     → POST /api/notes {title, content, tags}
                                      → POST /api/collections/:id/notes {note_ids}
9. Show success/failure per file
10. Redirect to collection view
```

## Docker Setup

### Dockerfile

```dockerfile
FROM node:20-alpine AS deps
WORKDIR /app
COPY package.json package-lock.json ./
RUN npm ci

FROM node:20-alpine AS builder
WORKDIR /app
COPY --from=deps /app/node_modules ./node_modules
COPY . .
RUN npm run build

FROM node:20-alpine AS runner
WORKDIR /app
ENV NODE_ENV=production
RUN addgroup --system --gid 1001 nodejs
RUN adduser --system --uid 1001 nextjs
COPY --from=builder /app/public ./public
COPY --from=builder --chown=nextjs:nodejs /app/.next/standalone ./
COPY --from=builder --chown=nextjs:nodejs /app/.next/static ./.next/static
USER nextjs
EXPOSE 3000
ENV PORT=3000
CMD ["node", "server.js"]
```

### docker-compose addition

```yaml
  web-ui:
    build:
      context: ./web
      dockerfile: Dockerfile
    environment:
      NEXT_PUBLIC_NOTE_SERVICE_URL: "http://note-service:8080"
      NEXT_PUBLIC_SEARCH_API_URL: "http://search-api:8081"
      NEXT_PUBLIC_TAG_API_URL: "http://tag-aggregator:8082"
    ports:
      - "3000:3000"
    depends_on:
      - note-service
      - search-api
      - tag-aggregator
```

## CORS Configuration

The Go backend services need CORS headers added for the frontend to call them from the browser. Add middleware to each HTTP service:

```go
func corsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }
        next.ServeHTTP(w, r)
    })
}
```

## Implementation Tasks

### Backend (Go)

1. Add CORS middleware to Note Service, Search API, and Tag Aggregator
2. Implement `GET /api/notes/:id/backlinks` endpoint (specified in original spec but not yet coded)
3. Implement all collection endpoints (see COLLECTIONS.md)

### Frontend (Next.js)

1. Initialize Next.js project with TypeScript + Tailwind in `web/` directory
2. Configure `next.config.ts` with `output: 'standalone'`
3. Create API client layer (`lib/api.ts`)
4. Build layout: header + sidebar + main content area
5. Build collection sidebar component
6. Build note list component (paginated)
7. Build markdown viewer with wiki-link support
8. Build note editor (create + edit)
9. Build search page with filters
10. Build tag explorer component
11. Build backlinks panel
12. Build file import page (drag-and-drop + preview)
13. Add Dockerfile for production build
14. Update `docker-compose.yml` with web-ui service

## Non-Goals (for now)

- Real-time collaboration / WebSocket updates
- Offline support / PWA
- User authentication (single-user system)
- Mobile-responsive design (desktop-first for demo)
- Dark mode (nice-to-have, not required)
- WYSIWYG editor (markdown textarea is sufficient)
