# PostgreSQL Migration

# Go CRM Migration Plan: From In-Memory to PostgreSQL

## Migration Steps

### 1. Schema Design

```sql
-- Clients table
CREATE TABLE clients (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    phone VARCHAR(50),
    company VARCHAR(255),
    address TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Activities table
CREATE TABLE activities (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    client_id UUID REFERENCES clients(id) ON DELETE CASCADE,
    type VARCHAR(100) NOT NULL,
    description TEXT,
    scheduled_at TIMESTAMP WITH TIME ZONE,
    completed_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Invoices table
CREATE TABLE invoices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    client_id UUID REFERENCES clients(id) ON DELETE CASCADE,
    number VARCHAR(100) UNIQUE NOT NULL,
    amount DECIMAL(10,2) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'draft',
    issued_at TIMESTAMP WITH TIME ZONE,
    due_at TIMESTAMP WITH TIME ZONE,
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Indexes for performance
CREATE INDEX idx_clients_email ON clients(email);
CREATE INDEX idx_clients_company ON clients(company);
CREATE INDEX idx_activities_client_id ON activities(client_id);
CREATE INDEX idx_activities_scheduled_at ON activities(scheduled_at);
CREATE INDEX idx_invoices_client_id ON invoices(client_id);
CREATE INDEX idx_invoices_status ON invoices(status);
CREATE INDEX idx_invoices_issued_at ON invoices(issued_at);
```

### 2. Connection Pooling

```go
// config.go
package config

import (
    "database/sql"
    "fmt"
    "log"
    "os"
    "time"
)

type DatabaseConfig struct {
    Host     string
    Port     string
    User     string
    Password string
    Name     string
    SSLMode  string
    PoolMin  int
    PoolMax  int
}

func (c *DatabaseConfig) DSN() string {
    return fmt.Sprintf(
        "host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
        c.Host, c.Port, c.User, c.Password, c.Name, c.SSLMode,
    )
}

func (c *DatabaseConfig) GetPoolConfig() *sql.DB {
    db, err := sql.Open("pgx", c.DSN())
    if err != nil {
        log.Fatal("Failed to open database:", err)
    }

    db.SetMaxOpenConns(c.PoolMax)
    db.SetMaxIdleConns(c.PoolMin)
    db.SetConnMaxLifetime(5 * time.Minute)

    return db
}
```

### 3. Migration Tool

```go
// migrations/000001_create_clients_table.up.sql
CREATE TABLE clients (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    phone VARCHAR(50),
    company VARCHAR(255),
    address TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_clients_email ON clients(email);
CREATE INDEX idx_clients_company ON clients(company);

// migrations/000002_create_activities_table.up.sql
CREATE TABLE activities (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    client_id UUID REFERENCES clients(id) ON DELETE CASCADE,
    type VARCHAR(100) NOT NULL,
    description TEXT,
    scheduled_at TIMESTAMP WITH TIME ZONE,
    completed_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_activities_client_id ON activities(client_id);
CREATE INDEX idx_activities_scheduled_at ON activities(scheduled_at);
```

### 4. Repository Pattern

```go
// repository/interface.go
package repository

import (
    "context"
    "time"
)

type Client struct {
    ID        string    `json:"id"`
    Name      string    `json:"name"`
    Email     string    `json:"email"`
    Phone     string    `json:"phone"`
    Company   string    `json:"company"`
    Address   string    `json:"address"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

type Activity struct {
    ID          string    `json:"id"`
    ClientID    string    `json:"client_id"`
    Type        string    `json:"type"`
    Description string    `json:"description"`
    ScheduledAt time.Time `json:"scheduled_at"`
    CompletedAt time.Time `json:"completed_at"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

type Invoice struct {
    ID        string    `json:"id"`
    ClientID  string    `json:"client_id"`
    Number    string    `json:"number"`
    Amount    float64   `json:"amount"`
    Status    string    `json:"status"`
    IssuedAt  time.Time `json:"issued_at"`
    DueAt     time.Time `json:"due_at"`
    Notes     string    `json:"notes"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

type Store interface {
    // Client operations
    CreateClient(ctx context.Context, client *Client) error
    GetClient(ctx context.Context, id string) (*Client, error)
    UpdateClient(ctx context.Context, client *Client) error
    DeleteClient(ctx context.Context, id string) error
    ListClients(ctx context.Context, limit, offset int) ([]*Client, error)
    
    // Activity operations
    CreateActivity(ctx context.Context, activity *Activity) error
    GetActivity(ctx context.Context, id string) (*Activity, error)
    UpdateActivity(ctx context.Context, activity *Activity) error
    DeleteActivity(ctx context.Context, id string) error
    ListActivities(ctx context.Context, clientID string, limit, offset int) ([]*Activity, error)
    
    // Invoice operations
    CreateInvoice(ctx context.Context, invoice *Invoice) error
    GetInvoice(ctx context.Context, id string) (*Invoice, error)
    UpdateInvoice(ctx context.Context, invoice *Invoice) error
    DeleteInvoice(ctx context.Context, id string) error
    ListInvoices(ctx context.Context, clientID string, limit, offset int) ([]*Invoice, error)
    
    // Health check
    Ping(ctx context.Context) error
}

// memory_store.go
package repository

import (
    "context"
    "sync"
    "time"
)

type MemoryStore struct {
    mu     sync.RWMutex
    clients map[string]*Client
    activities map[string]*Activity
    invoices map[string]*Invoice
}

func NewMemoryStore() *MemoryStore {
    return &MemoryStore{
        clients: make(map[string]*Client),
        activities: make(map[string]*Activity),
        invoices: make(map[string]*Invoice),
    }
}

// ... implement all methods using in-memory maps

// postgres_store.go
package repository

import (
    "context"
    "database/sql"
    "fmt"
    "time"
)

type PostgresStore struct {
    db *sql.DB
}

func NewPostgresStore(db *sql.DB) *PostgresStore {
    return &PostgresStore{db: db}
}

func (s *PostgresStore) Ping(ctx context.Context) error {
    return s.db.PingContext(ctx)
}

func (s *PostgresStore) CreateClient(ctx context.Context, client *Client) error {
    query := `
        INSERT INTO clients (id, name, email, phone, company, address, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
    
    _, err := s.db.ExecContext(ctx, query, 
        client.ID, client.Name, client.Email, client.Phone, 
        client.Company, client.Address, client.CreatedAt, client.UpdatedAt)
    return err
}

// ... implement all other methods with prepared statements
```

### 5. Transaction Handling

```go
// transaction.go
package repository

import (
    "context"
    "database/sql"
)

func (s *PostgresStore) WithTransaction(ctx context.Context, fn func(*sql.Tx) error) error {
    tx, err := s.db.BeginTx(ctx, nil)
    if err != nil {
        return err
    }
    
    defer func() {
        if p := recover(); p != nil {
            tx.Rollback()
            panic(p)
        }
    }()
    
    if err := fn(tx); err != nil {
        if rbErr := tx.Rollback(); rbErr != nil {
            return fmt.Errorf("transaction failed and rollback failed: %v", rbErr)
        }
        return err
    }
    
    return tx.Commit()
}

// Example usage in business logic
func (s *PostgresStore) CreateClientWithActivity(ctx context.Context, client *Client, activity *Activity) error {
    return s.WithTransaction(ctx, func(tx *sql.Tx) error {
        // Create client
        client.CreatedAt = time.Now()
        client.UpdatedAt = time.Now()
        if err := s.createClientTx(ctx, tx, client); err != nil {
            return err
        }
        
        // Create activity
        activity.ClientID = client.ID
        activity.CreatedAt = time.Now()
        activity.UpdatedAt = time.Now()
        if err := s.createActivityTx(ctx, tx, activity); err != nil {
            return err
        }
        
        return nil
    })
}
```

### 6. Query Patterns

```go
// postgres_store.go (continued)
func (s *PostgresStore) GetClient(ctx context.Context, id string) (*Client, error) {
    query := `
        SELECT id, name, email, phone, company, address, created_at, updated_at
        FROM clients WHERE id = $1`
    
    var client Client
    err := s.db.QueryRowContext(ctx, query, id).Scan(
        &client.ID, &client.Name, &client.Email, &client.Phone,
        &client.Company, &client.Address, &client.CreatedAt, &client.UpdatedAt)
    
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("client not found: %s", id)
        }
        return nil, err
    }
    
    return &client, nil
}

func (s *PostgresStore) ListClients(ctx context.Context, limit, offset int) ([]*Client, error) {
    query := `
        SELECT id, name, email, phone, company, address, created_at, updated_at
        FROM clients
        ORDER BY created_at DESC
        LIMIT $1 OFFSET $2`
    
    rows, err := s.db.QueryContext(ctx, query, limit, offset)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var clients []*Client
    for rows.Next() {
        var client Client
        if err := rows.Scan(
            &client.ID, &client.Name, &client.Email, &client.Phone,
            &client.Company, &client.Address, &client.CreatedAt, &client.UpdatedAt); err != nil {
            return nil, err
        }
        clients = append(clients, &client)
    }
    
    return clients, rows.Err()
}
```

## Go Code Changes

### Interfaces to Define

```go
// core interfaces
type ClientService interface {
    CreateClient(ctx context.Context, client *Client) error
    GetClient(ctx context.Context, id string) (*Client, error)
    UpdateClient(ctx context.Context, client *Client) error
    DeleteClient(ctx context.Context, id string) error
    ListClients(ctx context.Context, limit, offset int) ([]*Client, error)
}

type ActivityService interface {
    CreateActivity(ctx context.Context, activity *Activity) error
    GetActivity(ctx context.Context, id string) (*Activity, error)
    UpdateActivity(ctx context.Context, activity *Activity) error
    DeleteActivity(ctx context.Context, id string) error
    ListActivities(ctx context.Context, clientID string, limit, offset int) ([]*Activity, error)
}

type InvoiceService interface {
    CreateInvoice(ctx context.Context, invoice *Invoice) error
    GetInvoice(ctx context.Context, id string) (*Invoice, error)
    UpdateInvoice(ctx context.Context, invoice *Invoice) error
    DeleteInvoice(ctx context.Context, id string) error
    ListInvoices(ctx context.Context, clientID string, limit, offset int) ([]*Invoice, error)
}
```

### Swapping MemoryStore → PostgresStore

```go
// main.go
func main() {
    // Configuration
    dbConfig := &config.DatabaseConfig{
        Host:     getEnv("DB_HOST", "localhost"),
        Port:     getEnv("DB_PORT", "5432"),
        User:     getEnv("DB_USER", "postgres"),
        Password: getEnv("DB_PASSWORD", "password"),
        Name:     getEnv("DB_NAME", "crm"),
        SSLMode:  getEnv("DB_SSLMODE", "disable"),
        PoolMin:  5,
        PoolMax:  20,
    }
    
    // Choose store based on environment
    var store repository.Store
    if os.Getenv("USE_MEMORY_STORE") == "true" {
        store = repository.NewMemoryStore()
    } else {
        db := dbConfig.GetPoolConfig()
        store = repository.NewPostgresStore(db)
    }
    
    // Initialize services with the store
    clientService := service.NewClientService(store)
    activityService := service.NewActivityService(store)
    invoiceService := service.NewInvoiceService(store)
    
    // ... rest of application setup
}
```

### Connection String Configuration

```go
// config.go (extended)
func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

func GetDatabaseConfig() *DatabaseConfig {
    return &DatabaseConfig{
        Host:     getEnv("DB_HOST", "localhost"),
        Port:     getEnv("DB_PORT", "5432"),
        User:     getEnv("DB_USER", "postgres"),
        Password: getEnv("DB_PASSWORD", "password"),
        Name:     getEnv("DB_NAME", "crm"),
        SSLMode:  getEnv("DB_SSLMODE", "disable"),
        PoolMin:  getEnvInt("DB_POOL_MIN", 5),
        PoolMax:  getEnvInt("DB_POOL_MAX", 20),
    }
}

func getEnvInt(key string, defaultValue int) int {
    if value := os.Getenv(key); value != "" {
        if i, err := strconv.Atoi(value); err == nil {
            return i
        }
    }
    return defaultValue
}
```

### Health Check

```go
// health.go
package health

import (
    "context"
    "database/sql"
    "time"
)

type HealthChecker struct {
    db *sql.DB
}

func NewHealthChecker(db *sql.DB) *HealthChecker {
    return &HealthChecker{db: db}
}

func (h *HealthChecker) Check(ctx context.Context) error {
    // Create a context with timeout
    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()
    
    // Ping database
    if err := h.db.PingContext(ctx); err != nil {
        return err
    }
    
    return nil
}

// In your HTTP handler
func healthHandler(w http.ResponseWriter, r *http.Request) {
    if err := healthChecker.Check(r.Context()); err != nil {
        http.Error(w, "Database unhealthy", http.StatusServiceUnavailable)
        return
    }
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("OK"))
}
```

## Testing Strategy

### TestContainers for Integration Tests

```go
// integration_test.go
package integration

import (
    "context"
    "database/sql"
    "os"
    "testing"
    "time"

    "github.com/testcontainers/testcontainers-go"
    "github.com/testcontainers/testcontainers-go/modules/postgres"
    "github.com/testcontainers/testcontainers-go/wait"
    _ "github.com/jackc/pgx/v4/stdlib"
)

func TestPostgresStore(t *testing.T) {
    ctx := context.Background()
    
    // Start PostgreSQL container
    postgresC, err := postgres.RunContainer(ctx,
        testcontainers.WithImage("postgres:15"),
        postgres.WithDatabase("testdb"),
        postgres.WithUsername("testuser"),
        postgres.WithPassword("testpass"),
        testcontainers.WithWaitStrategy(
            wait.ForLog("database system is ready to accept connections").
                WithOccurrence(1).
                WithStartupTimeout(30*time.Second)),
    )
    if err != nil {
        t.Fatal(err)
    }
    defer func() {
        if err := postgresC.Terminate(ctx); err != nil {
            t.Fatalf("failed to terminate container: %v", err)
        }
    }()
    
    // Get connection string
    connStr, err := postgresC.ConnectionString(ctx)
    if err != nil {
        t.Fatal(err)
    }
    
    // Connect to database
    db, err := sql.Open("pgx", connStr)
    if err != nil {
        t.Fatal(err)
    }
    defer db.Close()
    
    // Run your tests
    testStore(t, db)
}

func testStore(t *testing.T, db *sql.DB) {
    store := repository.NewPostgresStore(db)
    
    // Test ping
    ctx := context.Background()
    if err := store.Ping(ctx); err != nil {
        t.Fatal("Store ping failed:", err)
    }
    
    // Test CRUD operations
    // ... your test cases
}
```

### Keep MemoryStore for Unit Tests

```go
// unit_test.go
package repository

import (
    "testing"
)

func TestMemoryStore(t *testing.T) {
    store := NewMemoryStore()
    
    // Test with memory store
    testStore(t, store)
}

func testStore(t *testing.T, store Store) {
    // Common test logic that works with any Store implementation
    ctx := context.Background()
    
    // Create client
    client := &Client{
        ID:        "test-id",
        Name:      "Test Client",
        Email:     "test@example.com",
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }
    
    if err := store.CreateClient(ctx, client); err != nil {
        t.Fatal("CreateClient failed:", err)
    }
    
    // Get client
    retrieved, err := store.GetClient(ctx, "test-id")
    if err != nil {
        t.Fatal("GetClient failed:", err)
    }
    
    if retrieved.Name != "Test Client" {
       