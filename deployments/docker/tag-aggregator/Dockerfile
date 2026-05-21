FROM golang:1.22-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /tag-aggregator ./cmd/tag-aggregator

FROM alpine:3.20
COPY --from=build /tag-aggregator /tag-aggregator
EXPOSE 8082
CMD ["/tag-aggregator"]
