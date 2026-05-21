FROM golang:1.22-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /backlink-extractor ./cmd/backlink-extractor

FROM alpine:3.20
COPY --from=build /backlink-extractor /backlink-extractor
CMD ["/backlink-extractor"]
