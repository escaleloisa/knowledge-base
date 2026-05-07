.PHONY: build run test clean docker-up docker-down

build:
	@go build -o bin/note-service ./cmd/note-service

run: build
	@./bin/note-service

test:
	@go test -v ./...

clean:
	@rm -rf bin/

docker-up:
	@docker compose up --build -d

docker-down:
	@docker compose down
