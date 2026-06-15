# Kafka Event-Driven Architecture

Apache Kafka is a distributed event streaming platform used for high-throughput, fault-tolerant messaging between microservices.

## Why Kafka?

- **Decoupling**: Services don't need to know about each other
- **Scalability**: Partitions allow parallel processing
- **Durability**: Events are persisted to disk
- **Replay**: Consumers can re-read events from any offset

## Key Concepts

### Topics and Partitions

```
Topic: "notes"
├── Partition 0: [event1, event4, event7]
├── Partition 1: [event2, event5, event8]
└── Partition 2: [event3, event6, event9]
```

### Consumer Groups

Multiple consumers in a group share the work. Each partition is consumed by exactly one consumer in the group.

## Event Schema

```json
{
  "type": "note.created",
  "note_id": "uuid-here",
  "title": "My Note",
  "content": "# Hello World",
  "tags": ["go", "kafka"],
  "timestamp": "2026-06-13T10:00:00Z"
}
```

## Patterns

1. **Event Sourcing**: Store all changes as events
2. **CQRS**: Separate read and write models
3. **Saga**: Coordinate distributed transactions

See [[docker-basics]] for containerizing Kafka and [[golang-patterns]] for Go consumer implementations.
