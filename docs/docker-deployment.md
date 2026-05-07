# Docker Deployment Guide

## Quick Start

```bash
cd knowledge-base

# Build and start all containers
docker compose up --build -d

# Verify containers are running
docker compose ps

# Make a request
curl -s -X POST http://localhost:8080/api/notes \
  -H "Content-Type: application/json" \
  -d '{"title":"My First Note","content":"Hello world","tags":["test"]}' | python3 -m json.tool

# Stop everything
docker compose down
```

---

## What Happens When You Run `docker compose up --build -d`

```
docker compose up --build -d
        │
        ▼
┌─────────────────────────────────────────────────────────┐
│ 1. Reads docker-compose.yml                             │
│ 2. Builds note-service image from Dockerfile            │
│ 3. Pulls postgres:16-alpine from Docker Hub             │
│ 4. Creates a Docker network (knowledge-base_default)    │
│ 5. Creates a volume (knowledge-base_pgdata)             │
│ 6. Starts postgres container                            │
│ 7. Runs healthcheck (pg_isready) every 5s               │
│ 8. Once healthy → starts note-service container         │
└─────────────────────────────────────────────────────────┘
```

### Flags

| Flag | Purpose |
|------|---------|
| `--build` | Rebuild images from Dockerfiles (picks up code changes) |
| `-d` | Detached mode — runs in background, gives you your terminal back |

Without `-d`, logs stream to your terminal and `Ctrl+C` stops everything.

---

## File Explanations

### docker-compose.yml

```yaml
services:
  # ─── Container 1: PostgreSQL Database ───────────────────────────────
  postgres:
    image: postgres:16-alpine          # Use official Postgres image from Docker Hub
                                       # "alpine" variant = smaller image (~80MB vs ~400MB)
    environment:
      POSTGRES_USER: postgres          # DB superuser name
      POSTGRES_PASSWORD: postgres      # DB password (don't use this in production!)
      POSTGRES_DB: knowledgebase       # Create this database on first start
    ports:
      - "5432:5432"                    # host_port:container_port
                                       # Lets you connect from your machine (e.g., pgAdmin, psql)
    volumes:
      - pgdata:/var/lib/postgresql/data
        # Named volume — persists DB data across container restarts
        # Without this, data is lost every time you run "docker compose down"

      - ./scripts/init.sql:/docker-entrypoint-initdb.d/init.sql
        # Postgres runs any .sql file in /docker-entrypoint-initdb.d/ on FIRST start
        # This creates our "notes" table automatically

    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
        # pg_isready checks if Postgres is accepting connections
      interval: 5s                     # Check every 5 seconds
      timeout: 3s                      # Fail if no response in 3s
      retries: 5                       # Mark unhealthy after 5 failures

  # ─── Container 2: Note Service (your Go app) ───────────────────────
  note-service:
    build:
      context: .                       # Build context = project root (so it can access all files)
      dockerfile: deployments/docker/note-service/Dockerfile
    environment:
      PORT: "8080"
      DATABASE_URL: "postgres://postgres:postgres@postgres:5432/knowledgebase?sslmode=disable"
        #                                         ^^^^^^^^
        #                                         "postgres" here = the service name above
        #                                         Docker DNS resolves container names automatically
        #                                         NOT localhost — containers are isolated
    ports:
      - "8080:8080"                    # Expose API on localhost:8080
    depends_on:
      postgres:
        condition: service_healthy     # Wait for Postgres healthcheck to pass
                                       # Without this, note-service might start before DB is ready

# ─── Volumes ─────────────────────────────────────────────────────────
volumes:
  pgdata:                              # Declares the named volume
                                       # Docker manages where it's stored on disk
                                       # Survives "docker compose down" but NOT "docker compose down -v"
```

### Dockerfile (deployments/docker/note-service/Dockerfile)

```dockerfile
# ─── Stage 1: Build ──────────────────────────────────────────────────
# Uses the full Go toolchain image (~300MB) to compile your code
FROM golang:1.22-alpine AS build
WORKDIR /app

# Copy dependency files first (for Docker layer caching)
# If go.mod/go.sum haven't changed, Docker reuses the cached layer
# This means "go mod download" doesn't re-run on every code change
COPY go.mod go.sum ./
RUN go mod download

# Copy all source code
COPY . .

# Compile to a single static binary
# Output: /note-service (one file, no dependencies)
RUN go build -o /note-service ./cmd/note-service

# ─── Stage 2: Run ────────────────────────────────────────────────────
# Uses a minimal Alpine image (~7MB) — no Go compiler, no source code
# Only contains: the compiled binary + basic OS utilities
FROM alpine:3.20

# Copy ONLY the binary from the build stage
COPY --from=build /note-service /note-service

# Document which port the app uses (informational only, doesn't publish it)
EXPOSE 8080

# Run the binary when the container starts
CMD ["/note-service"]

# ─── Result ──────────────────────────────────────────────────────────
# Final image ≈ 15MB (Alpine 7MB + your binary 8MB)
# vs if you used golang image directly ≈ 300MB+
# This is why multi-stage builds exist — small, secure production images
```

### Makefile

```makefile
.PHONY: build run test clean docker-up docker-down
# .PHONY tells Make these aren't real files — they're commands
# Without it, if a file named "build" existed, Make would skip the command

# Compile the Go binary to bin/ folder
build:
	@go build -o bin/note-service ./cmd/note-service
	# @ = don't print the command itself, just run it

# Build then run locally (without Docker — connects to localhost Postgres)
run: build
	@./bin/note-service

# Run all tests
test:
	@go test -v ./...
	# ./... = all packages recursively

# Remove compiled binaries
clean:
	@rm -rf bin/

# Build images and start all containers in background
docker-up:
	@docker compose up --build -d

# Stop and remove all containers
docker-down:
	@docker compose down
```

Usage:
```bash
make docker-up      # same as: docker compose up --build -d
make docker-down    # same as: docker compose down
make build          # compile locally (no Docker)
make run            # compile + run locally
make test           # run tests
```

---

## Network Diagram

```
┌─── Docker Network: knowledge-base_default ───┐
│                                               │
│  ┌──────────────┐     ┌──────────────────┐   │
│  │ note-service │────→│    postgres       │   │
│  │  :8080       │     │    :5432          │   │
│  └──────┬───────┘     └──────────────────┘   │
│         │                                     │
└─────────┼─────────────────────────────────────┘
          │
          │ port mapping (8080:8080)
          ▼
    Your machine (localhost:8080)
```

- Containers talk to each other by **service name** (`postgres`, not `localhost`)
- Docker creates a private network — containers are isolated from your machine
- `ports:` mapping is what exposes them to localhost

---

## Request Flow

```
curl -X POST localhost:8080/api/notes -d '{"title":"..."}'
              │
              ▼
┌─── Your machine ───────────────────────────────────────────┐
│  Port 8080 is mapped to the note-service container         │
└────────────────────────────┬───────────────────────────────┘
                             │
                             ▼
┌─── note-service container ─────────────────────────────────┐
│  routes.go:      POST /api/notes → handler.Create          │
│  handler.go:     parse JSON body, validate                 │
│  service.go:     call repo.Create()                        │
│  repository.go:  INSERT INTO notes ... RETURNING ...       │
│                        │                                   │
│                        ▼                                   │
│              SQL query sent to "postgres:5432"              │
└────────────────────────┬───────────────────────────────────┘
                         │ (Docker internal network)
                         ▼
┌─── postgres container ─────────────────────────────────────┐
│  Receives SQL, inserts row, returns result                 │
└────────────────────────────────────────────────────────────┘
                         │
                         ▼
              JSON response back to you
```

---

## API Requests

### Create a note

```bash
curl -s -X POST http://localhost:8080/api/notes \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Docker Basics",
    "content": "# Docker\nContainers package your app with its dependencies.",
    "tags": ["docker", "devops"]
  }' | python3 -m json.tool
```

### Get a note by ID

```bash
curl -s http://localhost:8080/api/notes/<note-id> | python3 -m json.tool
```

### List all notes (paginated)

```bash
curl -s "http://localhost:8080/api/notes?limit=10&offset=0" | python3 -m json.tool
```

### Update a note

```bash
curl -s -X PUT http://localhost:8080/api/notes/<note-id> \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Docker Basics (Updated)",
    "content": "# Docker\nUpdated content.",
    "tags": ["docker"]
  }' | python3 -m json.tool
```

### Delete a note

```bash
curl -s -o /dev/null -w "Status: %{http_code}\n" \
  -X DELETE http://localhost:8080/api/notes/<note-id>
```

Returns: `Status: 204` (no content, successfully deleted)

---

## Useful Docker Commands

| Command | What it does |
|---------|-------------|
| `docker compose ps` | See running containers and their status |
| `docker compose logs -f` | Stream logs from all containers (live) |
| `docker compose logs -f note-service` | Stream logs from one service only |
| `docker compose down` | Stop and remove containers |
| `docker compose down -v` | Stop, remove containers AND delete volumes (data) |
| `docker compose up --build -d note-service` | Rebuild and restart only one service |
| `docker exec -it knowledge-base-note-service-1 sh` | Open a shell inside the running container |
| `docker exec -it knowledge-base-postgres-1 psql -U postgres -d knowledgebase` | Open psql inside postgres |
| `docker images` | List all images on your machine |
| `docker system prune` | Clean up unused images/containers/networks |

---

## Common Issues

| Problem | Cause | Fix |
|---------|-------|-----|
| "connection refused" on curl | Container not running | `docker compose ps` to check, `docker compose logs` for errors |
| "database does not exist" | Volume was deleted, init.sql didn't run | `docker compose down -v` then `docker compose up --build -d` |
| Code changes not reflected | Didn't rebuild | Always use `--build` flag or `make docker-up` |
| Port already in use | Another process on 8080 | `lsof -i :8080` to find it, kill it, or change the port in docker-compose.yml |
