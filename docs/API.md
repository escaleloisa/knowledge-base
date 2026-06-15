# API Reference

All APIs return JSON. Errors return `{"error": "message"}` with appropriate HTTP status codes.

## Note Service — port 8080

### Create Note

```
POST /api/notes
```

**Request:**
```json
{
  "title": "Docker Basics",
  "content": "# Docker\nLearn containers here.",
  "tags": ["docker", "devops"]
}
```

**Response (201):**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "title": "Docker Basics",
  "content": "# Docker\nLearn containers here.",
  "tags": ["docker", "devops"],
  "created_at": "2026-06-14T05:00:00Z",
  "updated_at": "2026-06-14T05:00:00Z"
}
```

**Errors:**
- `400` — title or content missing

---

### List Notes

```
GET /api/notes?limit=20&offset=0
```

**Query params:**
- `limit` — max results (default: 20, max: 100)
- `offset` — pagination offset

**Response (200):** Array of notes, ordered by `created_at DESC`.

---

### Get Note

```
GET /api/notes/:id
```

**Response (200):** Single note object.
**Errors:** `404` — note not found.

---

### Update Note

```
PUT /api/notes/:id
```

**Request:**
```json
{
  "title": "Updated Title",
  "content": "Updated content",
  "tags": ["new-tag"]
}
```

**Response (200):** Updated note object.
**Errors:** `404` — note not found.

---

### Delete Note

```
DELETE /api/notes/:id
```

**Response:** `204 No Content`
**Errors:** `404` — note not found.

**Side effects:**
- Kafka event published → search index removed, tags decremented, backlinks cleaned up
- Removed from all collections (CASCADE)

---

### Get Backlinks

```
GET /api/notes/:id/backlinks
```

Returns notes that link TO this note via `[[wiki-links]]`.

**Response (200):**
```json
[
  {
    "id": "...",
    "title": "Kubernetes Intro",
    "content": "...",
    "tags": ["k8s"],
    "created_at": "...",
    "updated_at": "..."
  }
]
```

---

### Create Collection

```
POST /api/collections
```

**Request:**
```json
{
  "name": "Engineering Docs",
  "description": "Technical documentation for the team"
}
```

**Response (201):**
```json
{
  "id": "...",
  "name": "Engineering Docs",
  "description": "Technical documentation for the team",
  "created_at": "...",
  "updated_at": "..."
}
```

**Errors:** `400` — name is required.

---

### List Collections

```
GET /api/collections
```

**Response (200):** Array of collections with `note_count` included.

```json
[
  {
    "id": "...",
    "name": "Engineering Docs",
    "description": "...",
    "note_count": 5,
    "created_at": "...",
    "updated_at": "..."
  }
]
```

---

### Get Collection

```
GET /api/collections/:id
```

**Response (200):** Single collection with `note_count`.
**Errors:** `404` — collection not found.

---

### Update Collection

```
PUT /api/collections/:id
```

**Request:**
```json
{
  "name": "New Name",
  "description": "New description"
}
```

**Response (200):** Updated collection.
**Errors:** `404` — collection not found.

---

### Delete Collection

```
DELETE /api/collections/:id
```

**Response:** `204 No Content`
**Errors:** `404` — collection not found.

Notes in the collection are NOT deleted.

---

### List Collection Notes

```
GET /api/collections/:id/notes?limit=20&offset=0
```

**Response (200):** Array of notes in this collection, ordered by `added_at DESC`.

---

### Add Notes to Collection

```
POST /api/collections/:id/notes
```

**Request:**
```json
{
  "note_ids": ["uuid-1", "uuid-2", "uuid-3"]
}
```

**Response (200):**
```json
{
  "added": 3
}
```

Duplicates are silently ignored (ON CONFLICT DO NOTHING).

---

### Remove Note from Collection

```
DELETE /api/collections/:id/notes/:noteId
```

**Response:** `204 No Content`
**Errors:** `404` — note not in collection.

---

## Search API — port 8081

### Search Notes

```
GET /api/search?q=docker&tags=devops&from=2026-01-01&to=2026-12-31
```

**Query params:**
- `q` (required) — search query
- `tags` — comma-separated tag filter
- `from` — date range start (ISO format)
- `to` — date range end

**Response (200):**
```json
{
  "total": 2,
  "results": [
    {
      "id": "...",
      "title": "Docker Basics",
      "content": "...",
      "tags": ["docker"],
      "created_at": "...",
      "updated_at": "...",
      "_score": 4.52,
      "_highlight": {
        "content": ["Learn <em>Docker</em> containers here"]
      }
    }
  ]
}
```

**Scoring:** Title matches are boosted 2x over content matches.
**Errors:** `400` — query parameter `q` is required.

---

## Tag Aggregator — port 8082

### List Tags

```
GET /api/tags
```

**Response (200):**
```json
[
  { "name": "docker", "count": 5 },
  { "name": "go", "count": 3 },
  { "name": "kafka", "count": 2 }
]
```

Sorted by count descending. Updated in real-time via Kafka events.
