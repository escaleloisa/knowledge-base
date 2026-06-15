# Go Design Patterns

Common patterns and idioms used in Go applications.

## Repository Pattern

```go
type Repository interface {
    Create(ctx context.Context, entity *Entity) error
    Get(ctx context.Context, id string) (*Entity, error)
    Update(ctx context.Context, entity *Entity) error
    Delete(ctx context.Context, id string) error
    List(ctx context.Context, limit, offset int) ([]Entity, error)
}
```

## Dependency Injection

```go
type Service struct {
    repo Repository
    cache Cache
}

func NewService(repo Repository, cache Cache) *Service {
    return &Service{repo: repo, cache: cache}
}
```

## Error Handling

```go
// Sentinel errors
var ErrNotFound = errors.New("not found")

// Wrapping errors
if err != nil {
    return fmt.Errorf("create user: %w", err)
}

// Checking wrapped errors
if errors.Is(err, ErrNotFound) {
    // handle not found
}
```

## Graceful Shutdown

```go
ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
defer cancel()

srv := &http.Server{Addr: ":8080"}
go srv.ListenAndServe()

<-ctx.Done()
srv.Shutdown(context.Background())
```

## Concurrency Patterns

- **Worker Pool**: Limit goroutines with buffered channel
- **Fan-out/Fan-in**: Distribute work, collect results
- **Context cancellation**: Propagate timeouts

See [[kafka-events]] for async messaging patterns and [[kubernetes-intro]] for deployment patterns.
