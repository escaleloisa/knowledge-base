# Full Presentation Script

Total time: ~10-15 minutes (adjust based on interviewer engagement)

---

"I'll walk you through this in three parts: first I'll give you a quick overview of what the project is and the architecture using the README. Then I'll show you a live demo of the application running on EC2. And if you'd like, I can go into a code deep dive showing how the services are structured and how the event-driven flow works. Feel free to stop me and ask questions at any point."


## Part 1: README Walkthrough (~4 minutes)

*[Open README.md on screen]*

### Title & Intro

"This is a team knowledge base I built. The goal is to let teams organize their markdown documentation, search across all of it, and automatically find connections between documents using wiki-links. I built it as a learning exercise — specifically to get hands-on with Go, Kafka, Docker, and distributed systems design."

### Architecture

*[Scroll to Mermaid diagram]*

"The architecture is the core of this project. The frontend talks to three separate APIs. The Note Service is the main entry point — it handles CRUD operations and stores notes in Postgres.

The key design choice: every write operation publishes an event to Kafka. Three independent consumers pick up that event:
- One indexes the content into Elasticsearch for search
- One updates tag counts in Redis
- One parses wiki-links and builds a backlink graph

These consumers don't know about each other. They're completely decoupled. If I want to add a new feature — say, sending Slack notifications — I just write a new consumer. No changes to anything existing."

### Tech Stack

*[Scroll briefly — don't read the table]*

"Go backend with the standard library, no frameworks. Next.js frontend. Postgres, Elasticsearch, Kafka, Redis — each chosen for what it's best at. Everything containerized with Docker, deployed on EC2."

### Project Structure

*[Scroll to project structure, or open file explorer in editor]*

"It follows a standard Go project layout. `cmd/` has the entrypoint for each service — each one is its own binary. `internal/` has the business logic, separated by service.

The note-service follows a clean architecture pattern: handler for HTTP concerns, service for business logic and Kafka event publishing, repository for database queries. Each layer only knows about the one below it, and I use interfaces so I can mock the service layer and unit test handlers without a database.

`pkg/` has shared code — Kafka event types, data models, middleware. And `web/` is the Next.js frontend, completely separate from the backend."

### Design Decisions

*[Scroll to Design Decisions]*

"A few things I want to highlight: the event-driven approach gives us failure isolation — if search goes down, users can still create notes, events queue in Kafka until it recovers. I chose different data stores for different access patterns — Postgres for structured data, Elasticsearch for text search, Redis for counters — rather than forcing everything into one database.

And I'll be upfront — at this scale, microservices is over-engineered. A monolith would work fine. I chose this architecture because I wanted to practice containerization — each service has its own Dockerfile, its own multi-stage build — and learn how to orchestrate multiple containers with Docker Compose. It gave me a real reason to use Kafka for event-driven communication."

---

## Part 2: Specs Overview (~1-2 minutes)

*[Open the specs/ folder in editor or scroll to Specs section in README]*

"I have design specs for each major feature. Let me quickly walk through them."

### SPEC.md

*[Open briefly]*

"This is the original architecture doc — I wrote this before writing any code. It has the data models, API contracts, Kafka event schemas, and a phased implementation plan. It shows how I broke the project into incremental steps."

### COLLECTIONS.md

*[Open briefly]*

"When I added collections, I specced it out first — the database schema with a many-to-many join table, all the API endpoints, request/response formats. Collections are separate from tags by design: tags are metadata for search, collections are intentional user-driven organization."

### FRONTEND.md

*[Open briefly]*

"The frontend spec defines all the pages, which backend service each feature talks to, and how the UI components connect. It also covers the file import flow — how we parse markdown files client-side and POST them to the existing note API."

### DEPLOYMENT.md

*[Mention, don't open]*

"And I have a deployment guide — step by step instructions for spinning up on EC2 with Docker."

### docs/ (API & Database)

*[Open docs/API.md or docs/DATABASE.md briefly]*

"I also have reference docs — a full API reference with request/response examples for every endpoint, and a database schema doc with an ER diagram showing how the tables relate. These are for quick reference, not design discussion."

*[Transition]*

"Which brings us to the live demo..."

---

## Part 3: Live Demo (~5-7 minutes)

*[Switch to browser — http://EC2-IP:3000]*

"Let me show you it running."

### Import files

*[Navigate to /import]*

"I'll import some engineering docs. These are real markdown files with headings, code blocks, tables, and wiki-links between them."

*[Drag .md files in, add tags, select collection, click Import]*

"It extracts titles from the first heading, I can bulk-tag everything, and assign to a collection in one action."

### Browse collections

*[Click collection in sidebar]*

"Collections show up in the sidebar with note counts. Click in, see the notes. Simple organization."

### View a note

*[Click a note — show the rendered markdown]*

"The markdown renders with syntax highlighting, tables, everything. And look at the bottom — the backlinks panel. Other notes that reference this one through wiki-links show up automatically. That's the Backlink Extractor consumer doing its work via Kafka."

### Search

*[Navigate to /search, type a query]*

"Full-text search. It searches across titles and content, highlights the matching text, shows relevance scores. I can filter by tags too. This is Elasticsearch under the hood with BM25 scoring."

### Tags

*[Point to sidebar]*

"The tag cloud updates in real-time through Kafka. When those notes were created with tags, the Tag Aggregator consumed the events and updated Redis. This sorted set returns all tags ranked by usage in under a millisecond."

### Create a note

*[Click New Note, write something with a wiki-link]*

"I can create notes directly — write markdown, add tags, save. If I include a wiki-link to an existing note, the backlink extractor will detect it within seconds."

---

## Part 4: Code Deep Dive (3-5 minutes)

*[Open code editor — file explorer visible]*

### Entry point

*[Open cmd/note-service/main.go]*

"Each service starts here. It's simple — load config from environment variables, connect to Postgres, create a Kafka writer, then wire up the layers: repository takes the DB connection, service takes the repository and Kafka writer, handler takes the service, routes register the handler on the HTTP mux. Then we start listening with CORS middleware wrapping everything."

### The layered pattern

*[Open internal/note-service/handler/handler.go]*

"The handler is the HTTP layer. It decodes the request body, validates input — like checking title and content aren't empty — then calls the service. It doesn't know about the database or Kafka. It just calls `h.svc.Create()` and returns JSON."

*[Open internal/note-service/service/service.go]*

"The service is where things get interesting. When a note is created, it does two things: writes to Postgres via the repository, then publishes a Kafka event. This is the single point where the event-driven flow starts. One write, one event, three consumers react independently."

*[Scroll to publishEvent function]*

"The event contains the full note — ID, title, content, tags, timestamp. It's serialized as JSON and written to the `notes` Kafka topic. The key is the note ID, which ensures all events for the same note go to the same partition — maintaining order."

*[Open internal/note-service/repository/repository.go]*

"The repository is pure SQL. Parameterized queries, nothing fancy. It uses `RETURNING` to get back the full row after insert — so Postgres generates the UUID and timestamps, and we get them in one round trip instead of two queries."

### Interface and testing

*[Open internal/note-service/handler/service_interface.go]*

"The handler depends on this interface, not the concrete service. That's what makes unit testing possible without a database."

*[Open internal/note-service/handler/collection_handler_test.go]*

"In the tests, I create a mock that implements the interface. Each test sets up a function for the specific method being tested. Then I use `httptest` to send real HTTP requests to the handler and assert on status codes and response bodies. No Docker, no database, runs in milliseconds."

### Event consumers

*[Open internal/search-indexer/indexer.go]*

"The search indexer is a Kafka consumer in a loop. It reads a message, unmarshals the event, and depending on the type — created, updated, or deleted — it either PUTs the document to Elasticsearch or DELETEs it. The PUT uses the note ID as the document ID, so it's idempotent — processing the same event twice just overwrites with the same data."

*[Open internal/backlink-extractor/extractor.go]*

"The backlink extractor does the same Kafka consume loop, but its logic is different. It uses a regex to find all `[[wiki-link]]` patterns in the content, then queries Postgres to find notes with matching titles — case-insensitive, with hyphen-to-space normalization. When it finds a match, it inserts into the backlinks table. On every update, it clears existing backlinks for that note and re-extracts — so edits are always reflected."

*[Open internal/tag-aggregator/consumer/consumer.go]*

"The tag aggregator is the simplest consumer. On create, it increments the count for each tag in a Redis sorted set. On delete, it decrements and removes tags that hit zero. The sorted set gives us 'all tags ranked by popularity' in a single Redis command."

### Docker

*[Open docker-compose.yml]*

"Everything is orchestrated here. Infrastructure services use pre-built images — Postgres, Kafka, Elasticsearch, Redis. My services use `build:` with multi-stage Dockerfiles. Health checks and `depends_on` with conditions ensure the right startup order — Postgres must be healthy before note-service starts, Kafka must be healthy before consumers start."

*[Open deployments/docker/note-service/Dockerfile]*

"Multi-stage build. First stage compiles the Go binary. Second stage copies just the binary into a minimal `alpine` image. The final image is around 15-20MB instead of hundreds. Each service has its own Dockerfile following this pattern."

### Frontend

*[Open web/src/lib/api.ts]*

"The frontend has a centralized API client. It talks to three different backend services on different ports. Each function is typed with TypeScript interfaces matching the backend's JSON responses."

*[Open web/src/components/MarkdownViewer.tsx]*

"The markdown viewer uses react-markdown with plugins for GitHub Flavored Markdown and syntax highlighting. It also preprocesses wiki-links — replacing `[[text]]` with actual Next.js links so they're clickable and navigate within the app."

### Closing

"That's the full codebase — from entry point through the layered architecture, event publishing, async consumers, all the way to the frontend. Every piece has a clear responsibility."


---

## Closing

"That's the full project — architecture design, event-driven backend, working frontend, containerized deployment with CI. Happy to dive deeper into any specific part."
