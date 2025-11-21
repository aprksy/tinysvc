# TinySvc Architecture

## Overview

TinySvc is built following **Clean Architecture** principles, with a strong emphasis on **SOLID** design and **Dependency Injection through Composition**.

## Architecture Layers

### 1. Domain Layer (`internal/domain/`)

The **innermost layer** containing business entities and rules. This layer:
- Has **NO dependencies** on other layers
- Contains pure business logic
- Defines domain models and validation rules

**Files:**
- `paste.go`: Core paste entity and validation logic
- `errors.go`: Domain-specific errors

**Key Principles:**
```go
// Domain entities are self-contained
type Paste struct {
    ID         string
    Content    string
    IsMarkdown bool
    ExpiresAt  *time.Time
    CreatedAt  time.Time
}

// Business logic belongs to entities
func (p *Paste) IsExpired() bool {
    if p.ExpiresAt == nil {
        return false
    }
    return time.Now().After(*p.ExpiresAt)
}
```

### 2. Repository Layer (`internal/repository/`)

Defines **interfaces** for data persistence. This layer:
- Contains **only interfaces**, not implementations
- Allows domain/usecase layers to remain independent of storage details
- Enables easy swapping of storage backends

**Key Principle - Dependency Inversion:**
```go
// Interface defined in repository package
type PasteRepository interface {
    Create(ctx context.Context, paste *domain.Paste) error
    GetByID(ctx context.Context, id string) (*domain.Paste, error)
    Delete(ctx context.Context, id string) error
    DeleteExpired(ctx context.Context) (int64, error)
}

// Implementation lives in infrastructure layer
// Use cases depend on interface, not concrete implementation
```

### 3. Use Case Layer (`internal/usecase/`)

Contains **application business logic**. This layer:
- Orchestrates domain entities
- Depends only on repository **interfaces**
- Implements core application workflows

**Key Principle - Dependency Injection:**
```go
type PasteService interface {
    CreatePaste(ctx context.Context, req domain.PasteCreateRequest) (*domain.Paste, error)
    GetPaste(ctx context.Context, id string) (*domain.Paste, error)
    DeletePaste(ctx context.Context, id string) error
    CleanupExpired(ctx context.Context) (int64, error)
}

// Service depends on interface, not concrete implementation
type pasteService struct {
    repo repository.PasteRepository // Interface, not *sqlite.Repository
}

// Dependency injected through constructor
func NewPasteService(repo repository.PasteRepository) PasteService {
    return &pasteService{repo: repo}
}
```

### 4. Delivery Layer (`internal/delivery/http/`)

Handles **HTTP communication**. This layer:
- Translates HTTP requests to use case calls
- Converts use case responses to HTTP responses
- Handles HTTP-specific concerns (routing, middleware, serialization)
- Depends on use case **interfaces**

**Key Principle - Single Responsibility:**
```go
type PasteHandler struct {
    pasteService usecase.PasteService // Interface dependency
}

// Each handler focuses only on HTTP translation
func (h *PasteHandler) CreatePaste(w http.ResponseWriter, r *http.Request) {
    var req domain.PasteCreateRequest
    
    // HTTP deserialization
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        respondError(w, http.StatusBadRequest, "Invalid request body")
        return
    }
    
    // Delegate to use case
    paste, err := h.pasteService.CreatePaste(r.Context(), req)
    
    // HTTP serialization
    if err != nil {
        // Error handling
        return
    }
    
    respondJSON(w, http.StatusCreated, paste)
}
```

### 5. Infrastructure Layer (`internal/infrastructure/`)

Contains **external dependencies and implementations**. This layer:
- Implements repository interfaces
- Handles configuration
- Manages database connections
- Is the **outermost layer**

**Key Principle - Plugin Architecture:**
```go
// SQLite implementation of repository interface
type pasteRepository struct {
    db *sql.DB
}

func NewPasteRepository(db *sql.DB) repository.PasteRepository {
    return &pasteRepository{db: db}
}

// Returns interface, hiding implementation details
// Can be swapped with PostgreSQL, MongoDB, etc. without changing other layers
```

---

## Dependency Flow

```
main.go (Composition Root)
   │
   ├─> Infrastructure (SQLite, Config)
   │       │
   │       └─> Repository Interface Implementation
   │
   ├─> Use Cases
   │       │
   │       └─> Depends on Repository Interface
   │
   └─> HTTP Handlers
           │
           └─> Depends on Use Case Interface
```

**Key Point:** Dependencies point **inward**, toward the domain layer.

---

## SOLID Principles Applied

### Single Responsibility Principle (SRP)
Each component has **one reason to change**:
- `PasteHandler`: Changes only if HTTP interface changes
- `PasteService`: Changes only if business logic changes
- `PasteRepository`: Changes only if storage mechanism changes

### Open/Closed Principle (OCP)
System is **open for extension, closed for modification**:
- Add new storage backend by implementing `PasteRepository` interface
- Add new delivery mechanism (gRPC, CLI) without modifying use cases
- Add new services without modifying existing handlers

### Liskov Substitution Principle (LSP)
Any implementation can be substituted for its interface:
```go
// Can swap SQLite with any other implementation
var repo repository.PasteRepository
repo = sqlite.NewPasteRepository(db)
// OR
repo = postgres.NewPasteRepository(db)
// OR
repo = memory.NewPasteRepository()
```

### Interface Segregation Principle (ISP)
Interfaces are **small and focused**:
```go
// Small, focused interfaces
type PasteRepository interface {
    Create(ctx context.Context, paste *domain.Paste) error
    GetByID(ctx context.Context, id string) (*domain.Paste, error)
    Delete(ctx context.Context, id string) error
    DeleteExpired(ctx context.Context) (int64, error)
}

// Not a bloated interface with methods clients don't need
```

### Dependency Inversion Principle (DIP)
High-level modules don't depend on low-level modules. Both depend on abstractions:
```go
// Use case depends on interface (abstraction)
type pasteService struct {
    repo repository.PasteRepository // Interface, not concrete type
}

// Concrete implementation provided at runtime
service := usecase.NewPasteService(sqlite.NewPasteRepository(db))
```

---

## Composition Root (`cmd/server/main.go`)

The **only place** where concrete types are wired together:

```go
func run() error {
    // 1. Initialize infrastructure
    db, err := sqlite.InitDB(cfg.Database.Path)
    
    // 2. Create repository implementations
    pasteRepo := sqlite.NewPasteRepository(db)
    
    // 3. Inject into use cases
    pasteService := usecase.NewPasteService(pasteRepo)
    ipService := usecase.NewIPService()
    
    // 4. Inject into handlers
    router := httpdelivery.NewRouter(pasteService, ipService)
    
    // 5. Start server
    return srv.ListenAndServe()
}
```

**Benefits:**
- All dependencies flow from main → outward
- Easy to see entire dependency graph
- Simple to swap implementations for testing
- Clear separation of concerns

---

## Adding New Features

### Example: Adding a "URL Shortener" Service

#### 1. Define Domain Entity (`internal/domain/url.go`)
```go
type ShortURL struct {
    ID        string
    LongURL   string
    ShortCode string
    CreatedAt time.Time
}
```

#### 2. Define Repository Interface (`internal/repository/url_repository.go`)
```go
type URLRepository interface {
    Create(ctx context.Context, url *domain.ShortURL) error
    GetByCode(ctx context.Context, code string) (*domain.ShortURL, error)
}
```

#### 3. Implement Repository (`internal/infrastructure/persistence/sqlite/url_repository.go`)
```go
type urlRepository struct {
    db *sql.DB
}

func NewURLRepository(db *sql.DB) repository.URLRepository {
    return &urlRepository{db: db}
}
```

#### 4. Create Use Case (`internal/usecase/url_service.go`)
```go
type URLService interface {
    ShortenURL(ctx context.Context, longURL string) (*domain.ShortURL, error)
    GetURL(ctx context.Context, code string) (*domain.ShortURL, error)
}

type urlService struct {
    repo repository.URLRepository
}

func NewURLService(repo repository.URLRepository) URLService {
    return &urlService{repo: repo}
}
```

#### 5. Add HTTP Handler (`internal/delivery/http/url_handler.go`)
```go
type URLHandler struct {
    urlService usecase.URLService
}

func NewURLHandler(urlService usecase.URLService) *URLHandler {
    return &URLHandler{urlService: urlService}
}
```

#### 6. Wire in Main (`cmd/server/main.go`)
```go
func run() error {
    // ... existing setup ...
    
    // Add new repository
    urlRepo := sqlite.NewURLRepository(db)
    
    // Add new service
    urlService := usecase.NewURLService(urlRepo)
    
    // Update router to accept new service
    router := httpdelivery.NewRouter(pasteService, ipService, urlService)
    
    // ... rest of setup ...
}
```

**Notice:** No existing code needs to be modified, only extended!

---

## Testing Strategy

### Unit Testing Use Cases
```go
// Mock repository for testing
type mockPasteRepository struct {
    createFunc func(ctx context.Context, paste *domain.Paste) error
}

func (m *mockPasteRepository) Create(ctx context.Context, paste *domain.Paste) error {
    if m.createFunc != nil {
        return m.createFunc(ctx, paste)
    }
    return nil
}

func TestPasteService_CreatePaste(t *testing.T) {
    mockRepo := &mockPasteRepository{
        createFunc: func(ctx context.Context, paste *domain.Paste) error {
            return nil
        },
    }
    
    service := usecase.NewPasteService(mockRepo)
    
    // Test business logic without database
    paste, err := service.CreatePaste(context.Background(), req)
    // assertions...
}
```

### Integration Testing
```go
func TestSQLitePasteRepository(t *testing.T) {
    // Use in-memory SQLite for fast tests
    db, _ := sql.Open("sqlite3", ":memory:")
    repo := sqlite.NewPasteRepository(db)
    
    // Test actual database operations
    err := repo.Create(context.Background(), paste)
    // assertions...
}
```

---

## Performance Considerations

### For Low-Memory Environments (2GB RAM)

1. **SQLite Configuration**
   - Single connection pool (SQLite works best with 1 connection)
   - Proper indexing on frequently queried columns
   - Regular VACUUM to reclaim space

2. **Request Handling**
   - 60-second timeout prevents resource exhaustion
   - Rate limiting prevents abuse
   - Graceful shutdown ensures clean resource cleanup

3. **Memory Management**
   - Stream large responses instead of loading into memory
   - Set reasonable content size limits (10MB)
   - Periodic cleanup of expired data

4. **Goroutine Management**
   - Limited concurrent connections via http.Server settings
   - Context-based cancellation for all operations
   - Proper cleanup in defer statements

---

## Future Extensibility

The architecture supports easy addition of:

1. **Authentication Layer** (Priority 2)
   ```go
   // Add new middleware
   func (rt *Router) SetupAuthRoutes() http.Handler {
       r := chi.NewRouter()
       r.Use(rt.authMiddleware.Authenticate)
       // Protected routes
   }
   ```

2. **Different Storage Backends**
   ```go
   // Implement PasteRepository interface
   type postgresPasteRepository struct {
       db *pgx.Pool
   }
   
   // Swap in main.go
   pasteRepo := postgres.NewPasteRepository(db)
   ```

3. **Caching Layer**
   ```go
   type cachedPasteRepository struct {
       repo  repository.PasteRepository
       cache Cache
   }
   
   // Decorator pattern - no changes to existing code
   ```

4. **Metrics & Monitoring**
   ```go
   // Add middleware for metrics collection
   r.Use(middleware.Prometheus)
   ```

---

## Design Decisions

### Why SQLite?
- Zero configuration
- Perfect for single-instance deployments
- Low memory footprint
- File-based (easy backups)
- Sufficient performance for personal use

### Why Chi Router?
- Lightweight and fast
- Composable middleware
- Context-aware
- Good community support

### Why Clean Architecture?
- Testability without external dependencies
- Clear separation of concerns
- Easy to understand and maintain
- Supports future growth

---

## Conclusion

TinySvc demonstrates how Clean Architecture and SOLID principles create a maintainable, testable, and extensible system even for small projects. The clear separation of concerns and dependency injection make it easy to:

- Test components in isolation
- Swap implementations without breaking changes
- Add new features without modifying existing code
- Understand the codebase quickly

This architecture scales from a personal utility on a 2GB laptop to a production service handling significant load.
```

---

### **22. `docs/API.md`**

```markdown
# TinySvc API Documentation

Complete API reference for TinySvc endpoints.

## Base URL

```
http://localhost:8080/api/v1
```

## Authentication

Currently, all endpoints are **public** and do not require authentication (Priority 1 features).

Future Priority 2 features will implement OAuth authentication (Google, GitHub).

## Rate Limiting

- **Rate**: 10 requests per second per IP
- **Burst**: 20 requests
- **Response**: `429 Too Many Requests` when exceeded

## Common Response Codes

| Code | Description |
|------|-------------|
| 200 | Success |
| 201 | Created |
| 204 | No Content (successful deletion) |
| 400 | Bad Request (validation error) |
| 404 | Not Found |
| 410 | Gone (resource expired) |
| 413 | Payload Too Large (>10MB) |
| 429 | Too Many Requests (rate limit exceeded) |
| 500 | Internal Server Error |

---

## Endpoints

### Health Check

#### GET `/health`

Check if the service is running.

**Response:**
```json
{
  "status": "ok"
}
```

---

### IP Detection

#### GET `/ip`

Get your public IP address.

**Response:**
```json
{
  "ip": "203.0.113.42"
}
```

**Headers Checked (in order):**
1. `CF-Connecting-IP` (Cloudflare)
2. `X-Forwarded-For`
3. `X-Real-IP`
4. `RemoteAddr` (fallback)

**Example:**
```bash
curl http://localhost:8080/api/v1/ip
```

---

### Pastebin

#### POST `/paste`

Create a new paste.

**Request Body:**
```json
{
  "content": "string (required, max 10MB)",
  "is_markdown": "boolean (optional, default: false)",
  "expiry_days": "integer (optional)"
}
```

**Expiry Days Options:**
- `null` or `0`: Default (30 days)
- Positive number: Custom expiration
- Negative number (e.g., `-1`): Never expires

**Example: Simple Text Paste**
```bash
curl -X POST http://localhost:8080/api/v1/paste \
  -H "Content-Type: application/json" \
  -d '{
    "content": "Hello, world!",
    "is_markdown": false,
    "expiry_days": 30
  }'
```

**Example: Markdown Paste**
```bash
curl -X POST http://localhost:8080/api/v1/paste \
  -H "Content-Type: application/json" \
  -d '{
    "content": "# Title\n\n- Item 1\n- Item 2",
    "is_markdown": true,
    "expiry_days": 7
  }'
```

**Example: Permanent Paste**
```bash
curl -X POST http://localhost:8080/api/v1/paste \
  -H "Content-Type: application/json" \
  -d '{
    "content": "This never expires",
    "is_markdown": false,
    "expiry_days": -1
  }'
```

**Success Response (201):**
```json
{
  "id": "abc12345",
  "content": "Hello, world!",
  "is_markdown": false,
  "expires_at": "2024-02-15T10:30:00Z",
  "created_at": "2024-01-15T10:30:00Z"
}
```

**Error Responses:**

*400 Bad Request:*
```json
{
  "error": "content cannot be empty"
}
```

*413 Payload Too Large:*
```json
{
  "error": "content exceeds 10MB limit"
}
```

---

#### GET `/paste/{id}`

Retrieve a paste by ID.

**Parameters:**
- `id` (path): Paste identifier

**Example:**
```bash
curl http://localhost:8080/api/v1/paste/abc12345
```

**Success Response (200):**
```json
{
  "id": "abc12345",
  "content": "Hello, world!",
  "is_markdown": false,
  "expires_at": "2024-02-15T10:30:00Z",
  "created_at": "2024-01-15T10:30:00Z"
}
```

**Error Responses:**

*404 Not Found:*
```json
{
  "error": "Paste not found"
}
```

*410 Gone (Expired):*
```json
{
  "error": "Paste has expired"
}
```

---

#### DELETE `/paste/{id}`

Delete a paste.

**Parameters:**
- `id` (path): Paste identifier

**Example:**
```bash
curl -X DELETE http://localhost:8080/api/v1/paste/abc12345
```

**Success Response:**
- **204 No Content** (empty body)

**Error Responses:**

*404 Not Found:*
```json
{
  "error": "Paste not found"
}
```