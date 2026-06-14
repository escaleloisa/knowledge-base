# Collections Feature — Spec

## Overview

Collections are a first-class entity that allows users to manually organize notes into named groups. Unlike tags (which are metadata for search/filtering), collections are user-created folders for intentional organization — think Spotify playlists for your knowledge base.

## Goals

- Let users create, rename, and delete collections
- Allow assigning notes to one or more collections
- Provide a clear organizational layer separate from tags
- Support the frontend sidebar navigation

## Data Model

### collections (PostgreSQL)

#### Schema

| Column | Data Type | Mandatory | Description |
|--------|:---------:|:---------:|-------------|
| id | UUID | Y | Unique identifier, auto-generated |
| name | VARCHAR(255) | Y | Collection name, user-defined |
| description | TEXT | N | Optional description. Default: empty string |
| created_at | TIMESTAMPTZ | Y | Creation timestamp. Default: NOW() |
| updated_at | TIMESTAMPTZ | Y | Last modified timestamp. Default: NOW() |

#### Indexes

| Name | Type | Column(s) | Description |
|------|------|-----------|-------------|
| idx_collections_name | B-tree | name | Fast lookup by name |
| idx_collections_created_at | B-tree | created_at DESC | Chronological listing |

### collection_notes (PostgreSQL — join table)

#### Schema

| Column | Data Type | Mandatory | Description |
|--------|:---------:|:---------:|-------------|
| collection_id | UUID | Y | FK → collections(id) ON DELETE CASCADE |
| note_id | UUID | Y | FK → notes(id) ON DELETE CASCADE |
| added_at | TIMESTAMPTZ | Y | When the note was added. Default: NOW() |

Primary key: (collection_id, note_id)

#### Indexes

| Name | Type | Column(s) | Description |
|------|------|-----------|-------------|
| idx_collection_notes_note | B-tree | note_id | Fast "which collections is this note in?" query |

#### Query Patterns

```sql
-- List all collections
SELECT * FROM collections ORDER BY name ASC

-- Get notes in a collection (paginated, newest added first)
SELECT n.* FROM notes n
JOIN collection_notes cn ON cn.note_id = n.id
WHERE cn.collection_id = $1
ORDER BY cn.added_at DESC
LIMIT $2 OFFSET $3

-- Get collections a note belongs to
SELECT c.* FROM collections c
JOIN collection_notes cn ON cn.collection_id = c.id
WHERE cn.note_id = $1

-- Count notes per collection
SELECT c.id, c.name, COUNT(cn.note_id) as note_count
FROM collections c
LEFT JOIN collection_notes cn ON cn.collection_id = c.id
GROUP BY c.id, c.name
ORDER BY c.name ASC
```

## API Endpoints

All endpoints live on the **Note Service** (port 8080).

### Collection CRUD

| Method | Path | Description |
|--------|------|-------------|
| POST | /api/collections | Create a collection |
| GET | /api/collections | List all collections (with note counts) |
| GET | /api/collections/:id | Get collection details |
| PUT | /api/collections/:id | Update collection name/description |
| DELETE | /api/collections/:id | Delete collection (notes are NOT deleted) |

### Collection-Note Relationships

| Method | Path | Description |
|--------|------|-------------|
| GET | /api/collections/:id/notes | List notes in a collection (paginated) |
| POST | /api/collections/:id/notes | Add note(s) to collection `{note_ids: [...]}` |
| DELETE | /api/collections/:id/notes/:noteId | Remove a note from collection |
| GET | /api/notes/:id/collections | Get collections a note belongs to |

## Request/Response Schemas

### Create Collection

```json
// POST /api/collections
// Request
{ "name": "Engineering Docs", "description": "Technical documentation for the team" }

// Response (201)
{
  "id": "uuid",
  "name": "Engineering Docs",
  "description": "Technical documentation for the team",
  "created_at": "2026-06-13T...",
  "updated_at": "2026-06-13T..."
}
```

### List Collections (with counts)

```json
// GET /api/collections
// Response (200)
[
  { "id": "uuid", "name": "Engineering Docs", "description": "...", "note_count": 12, "created_at": "...", "updated_at": "..." },
  { "id": "uuid", "name": "Meeting Notes", "description": "...", "note_count": 5, "created_at": "...", "updated_at": "..." }
]
```

### Add Notes to Collection

```json
// POST /api/collections/:id/notes
// Request
{ "note_ids": ["uuid-1", "uuid-2", "uuid-3"] }

// Response (200)
{ "added": 3 }
```

### Remove Note from Collection

```
// DELETE /api/collections/:id/notes/:noteId
// Response (204 No Content)
```

## Migration Script

```sql
CREATE TABLE IF NOT EXISTS collections (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name        VARCHAR(255) NOT NULL,
    description TEXT DEFAULT '',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_collections_name ON collections (name);
CREATE INDEX IF NOT EXISTS idx_collections_created_at ON collections (created_at DESC);

CREATE TABLE IF NOT EXISTS collection_notes (
    collection_id UUID NOT NULL REFERENCES collections(id) ON DELETE CASCADE,
    note_id       UUID NOT NULL REFERENCES notes(id) ON DELETE CASCADE,
    added_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (collection_id, note_id)
);

CREATE INDEX IF NOT EXISTS idx_collection_notes_note ON collection_notes (note_id);
```

## Implementation Tasks

1. Add migration script to `scripts/`
2. Create `pkg/models/collection.go` with structs
3. Create `internal/note-service/repository/collection_repo.go`
4. Create `internal/note-service/service/collection_service.go`
5. Create `internal/note-service/handler/collection_handler.go`
6. Register collection routes in `internal/note-service/routes/routes.go`
7. Update `scripts/init.sql` with new tables

## Design Decisions

- **Collections live in Note Service** — no need for a separate microservice; they're tightly coupled to notes.
- **Many-to-many relationship** — a note can belong to multiple collections (like a file with symlinks in multiple folders).
- **Deleting a collection does NOT delete its notes** — notes are independent entities.
- **No nesting** — collections are flat for simplicity. Can add sub-collections later if needed.
- **Separate from tags** — tags are for search/filtering metadata; collections are for user-driven organization.
