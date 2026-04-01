# Go Conventions

Sources: Uber Go Style Guide, Google Go Style Guide, Effective Go

## Commands

```
go build ./...
go test ./...
go test ./... -race
go vet ./...
golangci-lint run
```

## Project Structure

```
cmd/
  myapp/
    main.go          — Entrypoint only, no logic
internal/
  handler/           — HTTP handlers
  service/           — Business logic
  repository/        — Data access
  model/             — Domain types
pkg/                 — Public library code (if any)
```

## Error Handling

```go
// Good — wrap with context
if err != nil {
    return fmt.Errorf("parse config %s: %w", path, err)
}

// Good — sentinel errors for expected conditions
var ErrNotFound = errors.New("not found")

if errors.Is(err, ErrNotFound) {
    // handle missing resource
}

// Good — custom error types for rich context
type ValidationError struct {
    Field   string
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation: %s %s", e.Field, e.Message)
}

// Bad — discarding errors
json.Unmarshal(data, &config)

// Bad — generic messages
return fmt.Errorf("something went wrong")
```

## Interfaces

```go
// Good — small interfaces, defined by consumer
type Reader interface {
    Read(p []byte) (n int, err error)
}

// Good — accept interfaces, return structs
func NewService(repo UserRepository) *Service { ... }

// Bad — large interfaces defined by the implementer
type UserManager interface {
    Create(u User) error
    Update(u User) error
    Delete(id string) error
    Find(id string) (User, error)
    List(filter Filter) ([]User, error)
    // ... 10 more methods
}
```

## Goroutines

```go
// Good — always handle goroutine lifecycle
g, ctx := errgroup.WithContext(ctx)
g.Go(func() error {
    return processItems(ctx, items)
})
if err := g.Wait(); err != nil {
    return err
}

// Good — pass context, respect cancellation
func fetch(ctx context.Context, url string) ([]byte, error) {
    req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
    // ...
}

// Bad — fire and forget
go doSomething()

// Bad — goroutine without cancellation path
go func() {
    for { process() }
}()
```

## Structs

```go
// Good — use field names in literals
srv := &Server{
    Addr:    ":8080",
    Handler: mux,
}

// Good — zero values should be useful
type Buffer struct {
    data []byte  // nil slice is valid empty buffer
}

// Bad — positional initialization
srv := &Server{":8080", mux}

// Bad — constructor for every struct (use zero value + setters if needed)
```

## Testing

```go
// Good — table-driven tests
func TestParseSize(t *testing.T) {
    tests := []struct {
        name  string
        input string
        want  int64
        err   bool
    }{
        {"bytes", "100B", 100, false},
        {"kilobytes", "1KB", 1024, false},
        {"invalid", "abc", 0, true},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := ParseSize(tt.input)
            if tt.err {
                require.Error(t, err)
                return
            }
            require.NoError(t, err)
            assert.Equal(t, tt.want, got)
        })
    }
}

// Good — test doubles via interfaces, not mocking frameworks
type fakeStore struct{ users map[string]User }
func (f *fakeStore) Get(id string) (User, error) { ... }

// Bad — testing private functions directly
// Bad — mocking everything
```

## Style

- Packages: short, lowercase, singular: `user` not `users` or `userService`
- Exported names: `PascalCase`
- Unexported names: `camelCase`
- Acronyms: all caps — `HTTPHandler`, `userID`
- Getters: `Name()` not `GetName()`
- Line length: no hard limit, but break at natural points
- Comments: on exported symbols, start with the name: `// Server handles incoming requests.`
- No `init()` unless absolutely necessary
- No `panic` in library code
