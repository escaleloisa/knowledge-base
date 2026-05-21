FROM golang:1.22-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /search-indexer ./cmd/search-indexer

FROM alpine:3.20
COPY --from=build /search-indexer /search-indexer
CMD ["/search-indexer"]
