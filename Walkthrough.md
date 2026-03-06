# Infrastructure Health Dashboard - Code Walkthrough

This document provides an in-depth explanation of how the Go API works. It's designed for developers who are new to building web APIs with Go.

---

## Table of Contents

1. [Project Structure](#project-structure)
2. [The Entry Point: main.go](#the-entry-point-maingo)
3. [Database Connection & Pooling](#database-connection--pooling)
4. [The Repository Pattern](#the-repository-pattern)
5. [HTTP Routing with Gin](#http-routing-with-gin)
6. [Handlers: Processing Requests](#handlers-processing-requests)
7. [Middleware](#middleware)
8. [Models: Data Structures](#models-data-structures)
9. [Request/Response Flow](#requestresponse-flow)
10. [Key Go Concepts Used](#key-go-concepts-used)

---

## Project Structure

```
personal-project/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── internal/                     # Private application code
│   ├── api/
│   │   ├── handlers.go          # HTTP request handlers
│   │   ├── router.go            # Route definitions
│   │   └── middleware.go        # Request/response middleware
│   ├── models/
│   │   ├── resource.go          # Resource data structures
│   │   ├── metric.go            # Metric data structures
│   │   └── health.go            # Health status structures
│   └── repository/
│       ├── mysql.go             # Database connection setup
│       ├── migrations.go        # Auto table creation
│       ├── resource_repo.go     # Resource database operations
│       └── metric_repo.go       # Metric database operations
├── migrations/
│   └── 001_initial.sql          # SQL schema (manual migration)
├── go.mod                        # Go module definition
└── go.sum                        # Dependency checksums
```

### Why This Structure?

- **`cmd/`**: Contains application entry points. Each subdirectory is a separate executable.
- **`internal/`**: Code that's private to this module. Other Go modules cannot import it.
- **`api/`**: Everything related to HTTP handling.
- **`models/`**: Data structures (structs) that represent your domain objects.
- **`repository/`**: Database access layer (following the Repository Pattern).

---

## The Entry Point: main.go

```go
// cmd/server/main.go
package main

import (
    "log"
    "github.com/HarkHorning/infra-health-dashboard/internal/api"
    "github.com/HarkHorning/infra-health-dashboard/internal/repository"
)

func main() {
    // 1. Connect to database
    db, err := repository.NewDBFromEnv()
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }
    defer db.Close()  // Ensures connection closes when main() exits

    // 2. Run migrations (create tables if needed)
    if err := repository.RunMigrations(db); err != nil {
        log.Fatalf("Failed to run migrations: %v", err)
    }

    // 3. Setup router with all routes
    router := api.SetupRouter(db)

    // 4. Start HTTP server
    if err := router.Run(":8080"); err != nil {
        log.Fatalf("Server failed: %v", err)
    }
}
```

### What's Happening Here?

1. **Database Connection**: We create a connection pool to MySQL. This pool manages multiple connections efficiently.

2. **`defer db.Close()`**: The `defer` keyword schedules `db.Close()` to run when `main()` exits. This ensures we clean up database connections even if the program crashes.

3. **Migrations**: Creates database tables if they don't exist.

4. **Router Setup**: Creates the Gin router with all our API routes.

5. **`router.Run(":8080")`**: Starts an HTTP server listening on port 8080. This call **blocks** - the program stays here handling requests until you stop it (Ctrl+C).

---

## Database Connection & Pooling

```go
// internal/repository/mysql.go

// Config holds database settings
type Config struct {
    Host            string
    Port            int
    User            string
    Password        string
    Database        string
    MaxOpenConns    int           // Maximum simultaneous connections
    MaxIdleConns    int           // Connections kept open when idle
    ConnMaxLifetime time.Duration // How long a connection can be reused
}
```

### Connection Pooling Explained

When your API receives many requests, creating a new database connection for each request is slow and wasteful. Instead, we use a **connection pool**:

```
Request 1 ──┐
Request 2 ──┼──► Connection Pool ──► MySQL Database
Request 3 ──┘    (5-25 connections)
```

- **MaxOpenConns (25)**: Never open more than 25 connections. Protects the database from overload.
- **MaxIdleConns (5)**: Keep 5 connections open even when not in use. They're ready for the next request.
- **ConnMaxLifetime (5 min)**: Recycle connections periodically to prevent stale connections.

### Building the Connection String (DSN)

```go
dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local",
    cfg.User,      // root
    cfg.Password,  // devpassword
    cfg.Host,      // localhost
    cfg.Port,      // 3307
    cfg.Database,  // infra_dashboard
)
// Result: "root:devpassword@tcp(localhost:3307)/infra_dashboard?parseTime=true&loc=Local"
```

- **`@tcp(host:port)`**: Connect via TCP protocol
- **`?parseTime=true`**: Automatically convert MySQL timestamps to Go `time.Time`
- **`&loc=Local`**: Use local timezone

### Environment Variables

```go
func ConfigFromEnv() Config {
    cfg := DefaultConfig()

    if host := os.Getenv("DB_HOST"); host != "" {
        cfg.Host = host
    }
    // ... more overrides
}
```

This pattern lets you:
- Use defaults for development
- Override with environment variables in production
- Never hardcode secrets in your code

---

## The Repository Pattern

The **Repository Pattern** abstracts database operations behind a clean interface. Your handlers don't know (or care) how data is stored.

```
Handler                 Repository              Database
   │                        │                      │
   │  GetByID(5)            │                      │
   ├───────────────────────►│                      │
   │                        │  SELECT * FROM...    │
   │                        ├─────────────────────►│
   │                        │     Row data         │
   │                        │◄─────────────────────┤
   │   *Resource            │                      │
   │◄───────────────────────┤                      │
```

### Example: ResourceRepository

```go
// internal/repository/resource_repo.go

type ResourceRepository struct {
    db *sqlx.DB  // Database connection pool
}

// Constructor - creates a new repository with a database connection
func NewResourceRepository(db *sqlx.DB) *ResourceRepository {
    return &ResourceRepository{db: db}
}

// Create inserts a new resource
func (r *ResourceRepository) Create(req models.CreateResourceRequest) (*models.Resource, error) {
    query := `
        INSERT INTO resources (resource_id, resource_type, name, region, tags)
        VALUES (?, ?, ?, ?, ?)
    `

    result, err := r.db.Exec(query, req.ResourceID, req.ResourceType, req.Name, region, tags)
    if err != nil {
        return nil, fmt.Errorf("failed to create resource: %w", err)
    }

    // Get the auto-generated ID
    id, _ := result.LastInsertId()

    // Fetch and return the complete resource
    return r.GetByID(int(id))
}
```

### Key Points:

1. **`(r *ResourceRepository)`**: This is a **method receiver**. It attaches the function to the `ResourceRepository` type.

2. **`?` Placeholders**: Never concatenate user input into SQL! Use `?` placeholders to prevent SQL injection:
   ```go
   // WRONG - SQL injection vulnerability!
   query := "SELECT * FROM users WHERE id = " + userInput

   // CORRECT - Safe parameterized query
   query := "SELECT * FROM users WHERE id = ?"
   db.Get(&user, query, userInput)
   ```

3. **Error Wrapping**: `fmt.Errorf("message: %w", err)` wraps the original error, preserving the stack trace.

### sqlx vs database/sql

We use `sqlx` instead of the standard `database/sql` because it's more convenient:

```go
// Standard library - verbose
rows, _ := db.Query("SELECT id, name FROM resources")
for rows.Next() {
    var id int
    var name string
    rows.Scan(&id, &name)
}

// sqlx - automatic struct mapping
var resources []models.Resource
db.Select(&resources, "SELECT * FROM resources")
```

The `db:"column_name"` struct tags tell sqlx how to map columns:

```go
type Resource struct {
    ID   int    `db:"id"`    // Maps to 'id' column
    Name string `db:"name"`  // Maps to 'name' column
}
```

---

## HTTP Routing with Gin

Gin is a high-performance HTTP framework for Go. It handles:
- Routing (matching URLs to handlers)
- Middleware (code that runs before/after handlers)
- JSON serialization/deserialization
- Parameter parsing

### Setting Up Routes

```go
// internal/api/router.go

func SetupRouter(db *sqlx.DB) *gin.Engine {
    // Create repositories (database access)
    resourceRepo := repository.NewResourceRepository(db)
    metricRepo := repository.NewMetricRepository(db)

    // Create handler with dependencies
    h := NewHandler(resourceRepo, metricRepo)

    // Create router
    router := gin.New()

    // Add middleware (runs on every request)
    router.Use(gin.Recovery())  // Recover from panics
    router.Use(Logger())        // Log requests
    router.Use(CORS())          // Allow cross-origin requests

    // Define routes
    router.GET("/health", h.HealthCheck)

    // Route groups share a common prefix
    v1 := router.Group("/api/v1")
    {
        resources := v1.Group("/resources")
        {
            resources.GET("", h.ListResources)           // GET /api/v1/resources
            resources.POST("", h.CreateResource)         // POST /api/v1/resources
            resources.GET("/:id", h.GetResource)         // GET /api/v1/resources/123
            resources.PUT("/:id", h.UpdateResource)      // PUT /api/v1/resources/123
            resources.DELETE("/:id", h.DeleteResource)   // DELETE /api/v1/resources/123
        }
    }

    return router
}
```

### Route Parameters

`:id` is a **path parameter**. Gin extracts it from the URL:

```
GET /api/v1/resources/42
                      ^^
                      id = "42"
```

Access it in handlers with `c.Param("id")`.

### Query Parameters

Query parameters come after `?` in the URL:

```
GET /api/v1/resources?type=azure_vm&region=eastus
```

Access them with `c.Query("type")` → returns `"azure_vm"`.

---

## Handlers: Processing Requests

Handlers are functions that process HTTP requests and return responses.

```go
// internal/api/handlers.go

// Handler holds dependencies (repositories)
type Handler struct {
    resourceRepo *repository.ResourceRepository
    metricRepo   *repository.MetricRepository
}

// GET /api/v1/resources/:id
func (h *Handler) GetResource(c *gin.Context) {
    // 1. Extract path parameter
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid resource id"})
        return
    }

    // 2. Call repository to get data
    resource, err := h.resourceRepo.GetByID(id)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
        return
    }

    // 3. Return JSON response
    c.JSON(http.StatusOK, resource)
}
```

### The `*gin.Context` Object

`c *gin.Context` is passed to every handler. It contains:
- The HTTP request (`c.Request`)
- Methods to read input (`c.Param`, `c.Query`, `c.ShouldBindJSON`)
- Methods to write output (`c.JSON`, `c.String`, `c.HTML`)

### Reading JSON Request Body

```go
// POST /api/v1/resources
func (h *Handler) CreateResource(c *gin.Context) {
    var req models.CreateResourceRequest

    // Parse JSON body into struct
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // req.ResourceID, req.Name, etc. are now populated
    resource, err := h.resourceRepo.Create(req)
    // ...
}
```

`ShouldBindJSON`:
1. Reads the request body
2. Parses it as JSON
3. Maps fields to the struct
4. Validates required fields (based on `binding:"required"` tags)

### Response Methods

```go
// JSON response
c.JSON(200, gin.H{"status": "ok"})
// Output: {"status":"ok"}

// JSON with struct
c.JSON(200, resource)
// Output: {"id":1,"name":"VM-1",...}

// Set status code
c.Status(204)  // No content

// Abort (stop processing, don't call remaining handlers)
c.AbortWithStatus(401)
```

### `gin.H` - A Shortcut for Maps

```go
gin.H{"key": "value"}
// Is equivalent to:
map[string]interface{}{"key": "value"}
```

---

## Middleware

Middleware are functions that run **before and/or after** your handlers. They're useful for:
- Logging
- Authentication
- CORS headers
- Error recovery

```go
// internal/api/middleware.go

func Logger() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        path := c.Request.URL.Path
        method := c.Request.Method

        // Call the next handler
        c.Next()

        // This runs AFTER the handler completes
        latency := time.Since(start)
        status := c.Writer.Status()

        log.Printf("%s %s %d %v", method, path, status, latency)
    }
}
```

### Middleware Execution Order

```
Request
   │
   ▼
┌──────────────┐
│   Logger()   │ ─── before c.Next()
├──────────────┤
│   CORS()     │ ─── before c.Next()
├──────────────┤
│   Handler    │ ─── actual request processing
├──────────────┤
│   CORS()     │ ─── after c.Next()
├──────────────┤
│   Logger()   │ ─── after c.Next()
└──────────────┘
   │
   ▼
Response
```

### Recovery Middleware

```go
router.Use(gin.Recovery())
```

If a handler panics (crashes), this middleware:
1. Catches the panic
2. Logs it
3. Returns a 500 error
4. Keeps the server running

Without this, a panic would crash your entire server!

---

## Models: Data Structures

Models define the shape of your data. They're used for:
- Database mapping (with `db` tags)
- JSON serialization (with `json` tags)
- Request validation (with `binding` tags)

```go
// internal/models/resource.go

type Resource struct {
    ID           int             `db:"id" json:"id"`
    ResourceID   string          `db:"resource_id" json:"resource_id"`
    ResourceType string          `db:"resource_type" json:"resource_type"`
    Name         string          `db:"name" json:"name"`
    Region       *string         `db:"region" json:"region,omitempty"`
    Tags         json.RawMessage `db:"tags" json:"tags,omitempty"`
    CreatedAt    time.Time       `db:"created_at" json:"created_at"`
    UpdatedAt    time.Time       `db:"updated_at" json:"updated_at"`
}
```

### Struct Tags Explained

| Tag | Purpose | Example |
|-----|---------|---------|
| `db:"column"` | Maps to database column | `db:"resource_id"` → `resource_id` column |
| `json:"field"` | JSON key name | `json:"resource_id"` → `{"resource_id": "..."}` |
| `json:",omitempty"` | Omit if empty/nil | Won't include `null` fields in response |
| `binding:"required"` | Validation - field must be present | Returns error if missing |

### Pointer Types for Nullable Fields

```go
Region *string  // Can be nil (NULL in database)
Name   string   // Cannot be nil (empty string if not set)
```

Using `*string` (pointer to string) allows the field to be `nil`, which maps to SQL `NULL`.

### Request vs Entity Models

```go
// Entity - represents database row
type Resource struct {
    ID        int       `db:"id" json:"id"`
    CreatedAt time.Time `db:"created_at" json:"created_at"`
    // ... all fields
}

// Request - represents API input (no ID, no timestamps)
type CreateResourceRequest struct {
    ResourceID   string `json:"resource_id" binding:"required"`
    ResourceType string `json:"resource_type" binding:"required"`
    Name         string `json:"name" binding:"required"`
    // ... only fields the client provides
}
```

---

## Request/Response Flow

Let's trace a complete request through the system:

### Example: `POST /api/v1/resources`

```
1. Client sends HTTP request
   POST /api/v1/resources
   Content-Type: application/json
   {"resource_id": "vm-001", "resource_type": "azure_vm", "name": "My VM"}

2. Gin router matches the route
   router.POST("/api/v1/resources", h.CreateResource)

3. Middleware runs (Logger, CORS)
   - Logger records start time
   - CORS adds headers

4. Handler executes
   func (h *Handler) CreateResource(c *gin.Context) {
       // Parse JSON into struct
       var req models.CreateResourceRequest
       c.ShouldBindJSON(&req)

       // Call repository
       resource, err := h.resourceRepo.Create(req)

       // Return response
       c.JSON(201, resource)
   }

5. Repository executes SQL
   INSERT INTO resources (resource_id, ...) VALUES (?, ...)

6. Response sent
   HTTP/1.1 201 Created
   Content-Type: application/json
   {"id": 1, "resource_id": "vm-001", ...}

7. Middleware cleanup
   - Logger prints: POST /api/v1/resources 201 5.2ms
```

---

## Key Go Concepts Used

### 1. Structs and Methods

```go
// Struct - a collection of fields
type Handler struct {
    resourceRepo *repository.ResourceRepository
}

// Method - a function attached to a struct
func (h *Handler) GetResource(c *gin.Context) {
    // h is the receiver - access h.resourceRepo here
}
```

### 2. Interfaces (Implicit)

Go interfaces are satisfied implicitly. If your type has the required methods, it implements the interface:

```go
// gin.HandlerFunc is: type HandlerFunc func(*Context)
// Any function with signature func(*gin.Context) satisfies it

func MyHandler(c *gin.Context) { }  // This is a gin.HandlerFunc
```

### 3. Error Handling

Go doesn't have exceptions. Functions return errors as values:

```go
result, err := someFunction()
if err != nil {
    // Handle error
    return err
}
// Use result
```

### 4. Defer

`defer` schedules a function to run when the surrounding function exits:

```go
func doSomething() {
    db, _ := sql.Open(...)
    defer db.Close()  // Will run when doSomething() returns

    // ... use db
}  // db.Close() runs here, even if there was a panic
```

### 5. Goroutines (Behind the Scenes)

Gin handles each request in a separate goroutine (lightweight thread). This is why:
- Multiple requests can be processed simultaneously
- You need a connection **pool** (not a single connection)
- You should be careful with shared state

---

## Testing the API

### Health Check
```powershell
Invoke-RestMethod -Uri http://localhost:8080/health
```

### Create a Resource
```powershell
$body = @{
    resource_id = "vm-001"
    resource_type = "azure_vm"
    name = "Production VM"
} | ConvertTo-Json

Invoke-RestMethod -Uri http://localhost:8080/api/v1/resources `
    -Method Post `
    -ContentType "application/json" `
    -Body $body
```

### List Resources
```powershell
Invoke-RestMethod -Uri http://localhost:8080/api/v1/resources
```

### Add a Metric
```powershell
$metric = @{
    metric_name = "cpu_percent"
    metric_value = 45.2
    unit = "percent"
} | ConvertTo-Json

Invoke-RestMethod -Uri http://localhost:8080/api/v1/resources/1/metrics `
    -Method Post `
    -ContentType "application/json" `
    -Body $metric
```

### Get Metrics
```powershell
Invoke-RestMethod -Uri "http://localhost:8080/api/v1/resources/1/metrics?metric_name=cpu_percent"
```

---

## Common Patterns & Best Practices

### 1. Dependency Injection

Instead of creating dependencies inside functions, pass them in:

```go
// Bad - hard to test, tightly coupled
func GetResource(id int) {
    db := sql.Open(...)  // Creates its own connection
}

// Good - dependencies injected
func NewHandler(repo *ResourceRepository) *Handler {
    return &Handler{resourceRepo: repo}
}
```

### 2. Error Wrapping

Add context to errors as they propagate up:

```go
// Repository
return fmt.Errorf("failed to get resource: %w", err)

// Handler sees: "failed to get resource: sql: no rows"
```

### 3. Empty Slices vs Nil

Return empty slices, not nil, for JSON arrays:

```go
// Returns null in JSON
var resources []Resource
return resources  // nil

// Returns [] in JSON
resources := []Resource{}
return resources  // empty slice
```

---

## Next Steps

Now that you understand the basics, try:

1. **Add a new endpoint**: Create a `/api/v1/resources/:id/health` endpoint
2. **Add validation**: Check that `resource_type` is one of the allowed values
3. **Add pagination**: Limit results with `?limit=10&offset=0`
4. **Add authentication**: Protect endpoints with API keys

Happy coding!
