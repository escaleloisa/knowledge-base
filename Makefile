.PHONY: build test lint run docker-up docker-down docker-rebuild logs clean

# Build all Go services
build:
	@go build ./...

# Run unit tests
test:
	@go test -v ./...

# Run tests with coverage
test-coverage:
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html

# Run the note-service locally (requires Postgres + Kafka)
run:
	@go run ./cmd/note-service

# Start frontend dev server
web-dev:
	@cd web && npm run dev

# Docker commands
docker-up:
	@docker compose up -d --build

docker-down:
	@docker compose down

docker-rebuild:
	@docker compose up -d --build note-service search-api tag-aggregator backlink-extractor web-ui

docker-logs:
	@docker compose logs -f

# View logs for a specific service (usage: make logs s=note-service)
logs:
	@docker compose logs -f $(s)

# Reset all data (dangerous)
reset-data:
	@docker compose exec redis redis-cli FLUSHALL
	@docker compose exec postgres psql -U postgres -d knowledgebase -c "DELETE FROM backlinks;"
	@docker compose exec postgres psql -U postgres -d knowledgebase -c "DELETE FROM collection_notes;"
	@docker compose exec postgres psql -U postgres -d knowledgebase -c "DELETE FROM collections;"
	@docker compose exec postgres psql -U postgres -d knowledgebase -c "DELETE FROM notes;"
	@echo "All data cleared."

# Clean build artifacts
clean:
	@rm -rf bin/ coverage.out coverage.html web/.next web/node_modules
