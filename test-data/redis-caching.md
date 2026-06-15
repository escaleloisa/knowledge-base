---
tags: [redis, caching, performance]
---

# Redis Caching Strategies

Redis is an in-memory data store used for caching, session management, and real-time analytics.

## Data Structures

| Structure | Use Case | Example |
|-----------|----------|---------|
| String | Simple cache | `SET user:123 "{json}"` |
| Hash | Object fields | `HSET user:123 name "John"` |
| List | Queues, feeds | `LPUSH notifications msg` |
| Set | Unique items | `SADD online_users user1` |
| Sorted Set | Rankings, tags | `ZADD tags:all 15 "docker"` |

## Tag Aggregation Pattern

```go
// When a note is created with tags ["go", "docker"]
pipe := rdb.Pipeline()
for _, tag := range note.Tags {
    pipe.Incr(ctx, "tag:"+tag)
    pipe.ZIncrBy(ctx, "tags:all", 1, tag)
}
pipe.Exec(ctx)
```

## Cache Patterns

### Cache-Aside (Lazy Loading)

```go
func GetUser(id string) *User {
    // Try cache first
    cached := redis.Get("user:" + id)
    if cached != nil {
        return cached
    }
    // Cache miss: load from DB
    user := db.GetUser(id)
    redis.Set("user:"+id, user, 5*time.Minute)
    return user
}
```

### Write-Through

Write to cache and DB simultaneously. Ensures consistency.

### TTL Strategy

- Hot data: 5 minutes
- Warm data: 1 hour  
- Cold data: Don't cache

See [[golang-patterns]] for Go Redis client usage and [[kafka-events]] for event-driven cache invalidation.
