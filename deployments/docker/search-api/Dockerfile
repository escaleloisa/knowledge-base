FROM golang:1.22-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /search-api ./cmd/search-api

FROM alpine:3.20
COPY --from=build /search-api /search-api
EXPOSE 8081
CMD ["/search-api"]
