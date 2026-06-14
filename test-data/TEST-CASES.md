# Demo Test Cases

Use these to walk through the demo and verify all features work.

## Test Data Files

| File | Tags | Wiki-links |
|------|------|-----------|
| docker-basics.md | docker, containers | [[kubernetes-intro]] |
| kafka-events.md | go, kafka | [[docker-basics]], [[golang-patterns]] |
| kubernetes-intro.md | k8s, docker | [[docker-basics]], [[kafka-events]] |
| golang-patterns.md | go, patterns | [[kafka-events]], [[kubernetes-intro]] |
| redis-caching.md | redis, caching, performance | [[golang-patterns]], [[kafka-events]] |
| elasticsearch-search.md | elasticsearch, search | [[kafka-events]] |

---

## 1. Collections CRUD

### Create collections
- [ ] Create "DevOps" collection → expect 201 with id
- [ ] Create "Go Programming" collection → expect 201
- [ ] Create "Data Stores" collection → expect 201
- [ ] Try creating with empty name → expect 400

### List collections
- [ ] GET /api/collections → shows all 3 with note_count = 0

### Update collection
- [ ] Rename "DevOps" to "Infrastructure" → expect 200

### Delete collection
- [ ] Delete "Data Stores" → expect 204
- [ ] Verify notes are NOT deleted

---

## 2. Notes CRUD

### Create notes via API
- [ ] POST /api/notes with title "Docker Basics", content from docker-basics.md, tags: ["docker", "containers"]
- [ ] POST with empty title → expect 400
- [ ] POST with empty content → expect 400

### Get/List notes
- [ ] GET /api/notes → paginated list
- [ ] GET /api/notes/{id} → full note content
- [ ] GET /api/notes?limit=2&offset=0 → only 2 results

### Update note
- [ ] PUT /api/notes/{id} → update tags
- [ ] Verify updated_at changes

### Delete note
- [ ] DELETE /api/notes/{id} → 204
- [ ] GET same id → 404

---

## 3. File Import (Frontend)

### Single file import
- [ ] Navigate to /import
- [ ] Drag docker-basics.md into drop zone
- [ ] Verify title extracted as "Docker Basics"
- [ ] Click Import → note created successfully

### Bulk import
- [ ] Select all 6 .md files at once
- [ ] Verify all titles extracted correctly
- [ ] Select "Infrastructure" collection
- [ ] Import → all 6 created and added to collection

### Frontmatter parsing
- [ ] Import redis-caching.md → tags auto-detected from frontmatter (redis, caching, performance)

---

## 4. Collection-Note Relationships

### Add notes to collections
- [ ] Add docker-basics + kubernetes-intro to "Infrastructure" collection
- [ ] Add golang-patterns + kafka-events to "Go Programming" collection
- [ ] Verify note_count updates on collection list

### View collection notes
- [ ] GET /api/collections/{id}/notes → shows correct notes
- [ ] Navigate to /collections/{id} in UI → shows note cards

### Remove note from collection
- [ ] Remove kubernetes-intro from "Infrastructure"
- [ ] Verify note still exists in /api/notes (not deleted)
- [ ] Collection note_count decreases

---

## 5. Search (Elasticsearch)

### Full-text search
- [ ] Search "docker" → returns docker-basics and kubernetes-intro
- [ ] Search "consumer group" → returns kafka-events
- [ ] Verify highlights in results

### Tag filtering
- [ ] Search "patterns" with tag filter "go" → returns golang-patterns
- [ ] Search with non-existent tag → 0 results

### Empty/invalid search
- [ ] Search with empty query → 400 error

---

## 6. Tags (Tag Aggregator)

### List tags
- [ ] GET /api/tags → sorted by count descending
- [ ] Verify "go" appears (used in multiple notes)
- [ ] Verify counts are accurate

### Tag cloud in UI
- [ ] Sidebar shows tags with counts
- [ ] Click a tag → navigates to search filtered by that tag

---

## 7. Backlinks (Backlink Extractor)

### Wiki-link detection
- [ ] Create note with content containing `[[Docker Basics]]`
- [ ] Wait for backlink extractor to process (Kafka event)
- [ ] GET /api/notes/{docker-basics-id}/backlinks → shows the linking note

### Backlinks panel in UI
- [ ] View docker-basics note → backlinks panel shows notes that reference it
- [ ] kubernetes-intro and kafka-events should appear as backlinks

---

## 8. Frontend Navigation

### Layout
- [ ] Header shows: logo, search, + New Note
- [ ] Sidebar shows: collections with counts, tags

### Note viewer
- [ ] Markdown renders correctly (headings, code blocks, tables)
- [ ] Syntax highlighting works in code blocks
- [ ] Wiki-links are clickable

### Note editor
- [ ] Create new note → textarea + preview toggle
- [ ] Edit existing note → pre-fills title, content, tags
- [ ] Save → redirects to viewer

---

## 9. Demo Script (Interview Flow)

1. **Show architecture** → explain microservices, Kafka events, Elasticsearch
2. **Import files** → drag 6 .md files, show collection picker
3. **Browse collections** → show sidebar, click into a collection
4. **View note** → show rendered markdown, code highlighting
5. **Search** → type "kubernetes", show highlighted results
6. **Show backlinks** → click a note, show "linked from" panel
7. **Create note** → write new markdown, add tags
8. **Show tags** → tag cloud updates with new tags
9. **Show docker-compose** → explain the infrastructure
10. **Show code** → walk through handler → service → repo layers
