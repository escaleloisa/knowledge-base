# Interview Demo Script

## Opening (1 minute)

"This is a self-hosted team knowledge base I built. The idea is teams can organize their markdown documentation into collections, search across everything with full-text search, and automatically discover connections between documents through wiki-links.

I built it with an event-driven microservices architecture — Go for the backend, Kafka for async messaging, Elasticsearch for search, Redis for real-time tag aggregation, PostgreSQL for persistence, and a Next.js frontend. Everything runs in Docker containers, deployed on AWS EC2."

---

## Architecture Overview (2-3 minutes)

*Open the SPEC.md or draw on whiteboard*

"Let me walk you through the architecture:

**The write path:**
When a user creates or updates a note, the Note Service writes to Postgres, then publishes an event to Kafka. That single event triggers three consumers in parallel:

1. **Search Indexer** — indexes the content into Elasticsearch for full-text search
2. **Tag Aggregator** — updates tag counts in Redis (sorted sets for instant ranking)
3. **Backlink Extractor** — parses wiki-links like `[[Docker Basics]]` and stores the graph in Postgres

**The read path:**
The frontend talks to three different APIs depending on what it needs — the Note Service for CRUD, the Search API for full-text queries, and the Tag Aggregator for the tag cloud.

**Why this architecture?**
- Services are decoupled — the Note Service doesn't know or care about search indexing
- Each consumer can be scaled independently
- If Elasticsearch goes down, note creation still works — events queue in Kafka
- Adding a new feature (like notifications) is just adding another Kafka consumer"

---

## Live Demo (5-7 minutes)

### 1. Import files

*Navigate to /import*

"Let me import some engineering documentation. I'll drag in these markdown files..."

*Drag 5-6 .md files, add tags like 'engineering, devops', select a collection*

"Notice it extracts titles from the first heading, and I can apply tags to all files at once. These files also have wiki-links between them — we'll see those resolve in a moment."

### 2. Browse collections

*Click the collection in the sidebar*

"Collections are a first-class entity — users create them, organize notes into them. A note can be in multiple collections. Think of it like Spotify playlists for documentation."

### 3. View a note (markdown rendering)

*Click on a note like 'Kafka Event-Driven Architecture'*

"The content renders as formatted markdown — code blocks with syntax highlighting, tables, headings. And notice at the bottom — the backlinks panel shows which other notes reference this one through wiki-links. That happened automatically via the Backlink Extractor consumer."

### 4. Search

*Navigate to /search, type 'kubernetes'*

"Full-text search powered by Elasticsearch. It searches across titles and content, highlights matching text, and I can filter by tags. The title field is boosted 2x so title matches rank higher."

### 5. Tags

*Point to sidebar tag cloud*

"Tags update in real-time through Kafka events. When I created those notes with tags, the Tag Aggregator consumed the events and updated Redis. This sorted set gives me 'all tags ranked by popularity' in under a millisecond."

---

## Technical Deep Dive (3-5 minutes, if asked)

### Code structure

*Open the repo*

"Each service follows the same layered pattern:
- **Handler** — HTTP request/response, validation
- **Service** — business logic, Kafka event publishing
- **Repository** — database queries

I use Go's standard library `net/http` with the new Go 1.22 routing (method + path patterns). No heavy frameworks."

### Event-driven flow

*Open the note-service service.go*

"When a note is created, the service writes to Postgres, then publishes a Kafka event with the full note data. The consumers are independent — they each have their own consumer group, so they process at their own pace."

### Collections design

"Collections are a separate first-class entity with a many-to-many relationship to notes. I chose this over using tags-as-collections because users need explicit control — rename, delete, organize. Tags are metadata for search; collections are intentional organization."

### Testing approach

"I have handler-level unit tests with mocked service interfaces. The interface pattern lets me test HTTP behavior without needing a database. For integration testing, I use docker-compose to spin up the full stack."

---

## Questions I'm Ready For

**Q: Why microservices for something this size?**
"In production this could be a monolith. I chose microservices because I wanted to practice containerization — each service has its own Dockerfile and multi-stage build — and orchestrating multiple containers with Docker Compose. It also let me learn event-driven communication patterns with Kafka, where services are decoupled and can scale independently."

**Q: Why Kafka instead of a simpler message queue?**
"Kafka gives us event replay, consumer groups, and ordering guarantees. If I need to rebuild the search index, I can replay all events from the beginning. A simpler queue like RabbitMQ would work too, but Kafka better demonstrates event sourcing patterns."

**Q: Why not use Postgres for search?**
"Postgres has `ILIKE` and full-text search (`tsvector`), and that's fine for small scale. Elasticsearch gives us relevance scoring (BM25), field boosting, highlighting, and it scales horizontally for read-heavy workloads."

**Q: What would you add next?**
"Authentication and team-based permissions. Users belong to teams, collections have visibility levels (private, team, public), and notes inherit permissions. I'd add JWT auth with a middleware on all endpoints."

**Q: How do you handle failures?**
"Kafka consumers have at-least-once delivery — if a consumer crashes, it re-reads from the last committed offset. The trade-off is occasional duplicates, which I handle with `ON CONFLICT DO NOTHING` in the database. For production I'd add a dead letter topic for messages that fail repeatedly."

**Q: How is this deployed?**
"Docker Compose on a single EC2 instance for the demo. For production, I'd move to Kubernetes — each service as a Deployment, Kafka via Strimzi operator, Elasticsearch via ECK operator, and the Note Importer as a CronJob."

---

## Closing (30 seconds)

"This project covers the full lifecycle — design, implementation, testing, containerization, deployment. It demonstrates event-driven architecture, microservices communication, and how to build a system where adding new features is just adding a new Kafka consumer."

---

## Pre-Demo Checklist

- [ ] EC2 instance is running
- [ ] `docker compose ps` shows all services healthy
- [ ] Flush Redis if there are orphaned tags: `docker compose exec redis redis-cli FLUSHALL`
- [ ] Delete old test data if needed
- [ ] Have test .md files ready to import
- [ ] Browser open to http://<EC2-IP>:3000
- [ ] Code editor open to show architecture
