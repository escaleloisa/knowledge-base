FROM golang:1.22-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /note-service ./cmd/note-service

FROM alpine:3.20
COPY --from=build /note-service /note-service
EXPOSE 8080
CMD ["/note-service"]
