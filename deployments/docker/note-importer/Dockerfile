FROM golang:1.22-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /note-importer ./cmd/note-importer

FROM alpine:3.20
COPY --from=build /note-importer /note-importer
CMD ["/note-importer"]
