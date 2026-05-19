# Knowledge Base — Event-Driven Microservices

A self-hosted knowledge base built with an event-driven microservices architecture. This project serves as a hands-on exploration of Go, Kafka, Elasticsearch, Kubernetes, and containerized service design by building a real application end-to-end.

## Tech Stack

| Component | Technology |
|-----------|-----------|
| Language | Go |
| Database | PostgreSQL |
| Search | Elasticsearch |
| Messaging | Apache Kafka |
| Cache | Redis |
| Containers | Docker |
| Orchestration | Kubernetes (kind) |

## Architecture

Six services communicating via Kafka events:

- **Note Service** — CRUD API, publishes events to Kafka
- **Search Indexer** — Consumes events, indexes into Elasticsearch
- **Search API** — Full-text search with filtering and highlighting
- **Tag Aggregator** — Maintains tag counts in Redis
- **Backlink Extractor** — Detects `[[wiki-links]]` between notes
- **Note Importer** — CronJob that scans a folder for markdown files

## What I Learned

- Designing event-driven systems with Kafka (producer/consumer patterns)
- Multi-stage Docker builds for Go services
- Docker Compose for local multi-service development
- Kubernetes workload types: Deployments, CronJobs, Services
- Data consistency patterns across distributed services
- Tradeoffs between monolithic and microservice architectures

## Running Locally

```bash
# Start all services
docker-compose up -d

# Import markdown files
cp your-notes/*.md import/
```

## Project Structure

```
├── cmd/                  # Service entrypoints
├── internal/             # Service-specific logic
├── pkg/                  # Shared packages (config, kafka)
├── deployments/docker/   # Dockerfiles per service
├── scripts/              # DB init scripts
├── import/               # Drop .md files here for import
├── specs/                # Technical specifications
└── SPEC.md
```

## Status

🚧 In progress — building incrementally as a learning exercise.
