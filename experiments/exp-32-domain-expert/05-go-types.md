```go
package store

import (
    "sync"
    "time"
)

// Address represents a client's address
type Address struct {
    Street  string
    City    string
    State   string
    ZipCode string
    Country string
}

// Client represents a client in the CRM
type Client struct {
    ID        string
    Name      string
    Email     string
    Phone     string
    Address   *Address
    CreatedAt time.Time
    UpdatedAt time.Time
}

// Activity represents time spent on a project
type Activity struct {
    ID        string
    ClientID  string
    ProjectID string
    Description string
    Duration  time.Duration // in minutes
    Date      time.Time
    CreatedAt time.Time
    UpdatedAt time.Time
}

// LineItem represents an item in an invoice
type LineItem struct {
    ID          string
    Description string
    Quantity    int
    Rate        float64
    Subtotal    float64 // calculated
}

// InvoiceStatus represents the status of an invoice
type InvoiceStatus string

const (
    InvoiceStatusDraft     InvoiceStatus = "draft"
    InvoiceStatusSent      InvoiceStatus = "sent"
    InvoiceStatusPaid      InvoiceStatus = "paid"
    InvoiceStatusVoid      InvoiceStatus = "void"
    InvoiceStatusOverdue   InvoiceStatus = "overdue"
)

// Invoice represents an invoice
type Invoice struct {
    ID            string
    InvoiceNumber string
    ClientID      string
    InvoiceDate   time.Time
    DueDate       time.Time
    LineItems     []LineItem
    Subtotal      float64
    Tax           float64
    Total         float64
    Status        InvoiceStatus
    CreatedAt     time.Time
    UpdatedAt     time.Time
}

// Store represents the CRM data store
type Store struct {
    clients   map[string]Client
    activities map[string]Activity
    invoices  map[string]Invoice
    projects  map[string]Project
    mu        sync.RWMutex
}

// Project represents a project
type Project struct {
    ID        string
    ClientID  string
    Name      string
    Description string
    CreatedAt time.Time
    UpdatedAt time.Time
}

// NewStore creates a new Store
func NewStore() *Store {
    return &Store{
        clients:   make(map[string]Client),
        activities: make(map[string]Activity),
        invoices:  make(map[string]Invoice),
        projects:  make(map[string]Project),
    }
}

// Client CRUD operations
func (s *Store) CreateClient(client Client) error {
    s.mu.Lock()
    defer s.mu.Unlock()
    client.ID = generateID()
    client.CreatedAt = time.Now()
    client.UpdatedAt = time.Now()
    s.clients[client.ID] = client
    return nil
}

func (s *Store) GetClient(id string) (*Client, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    client, exists := s.clients[id]
    if !exists {
        return nil, ErrNotFound
    }
    return &client, nil
}

func (s *Store) UpdateClient(id string, client Client) error {
    s.mu.Lock()
    defer s.mu.Unlock()
    existing, exists := s.clients[id]
    if !exists {
        return ErrNotFound
    }
    client.ID = id
    client.CreatedAt = existing.CreatedAt
    client.UpdatedAt = time.Now()
    s.clients[id] = client
    return nil
}

func (s *Store) DeleteClient(id string) error {
    s.mu.Lock()
    defer s.mu.Unlock()
    delete(s.clients, id)
    return nil
}

func (s *Store) ListClients() []Client {
    s.mu.RLock()
    defer s.mu.RUnlock()
    var clients []Client
    for _, client := range s.clients {
        clients = append(clients, client)
    }
    return clients
}

// Activity CRUD operations
func (s *Store) CreateActivity(activity Activity) error {
    s.mu.Lock()
    defer s.mu.Unlock()
    activity.ID = generateID()
    activity.CreatedAt = time.Now()
    activity.UpdatedAt = time.Now()
    s.activities[activity.ID] = activity
    return nil
}

func (s *Store) GetActivity(id string) (*Activity, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    activity, exists := s.activities[id]
    if !exists {
        return nil, ErrNotFound
    }
    return &activity, nil
}

func (s *Store) UpdateActivity(id string, activity Activity) error {
    s.mu.Lock()
    defer s.mu.Unlock()
    existing, exists := s.activities[id]
    if !exists {
        return ErrNotFound
    }
    activity.ID = id
    activity.CreatedAt = existing.CreatedAt
    activity.UpdatedAt = time.Now()
    s.activities[id] = activity
    return nil
}

func (s *Store) DeleteActivity(id string) error {
    s.mu.Lock()
    defer s.mu.Unlock()
    delete(s.activities, id)
    return nil
}

func (s *Store) ListActivities() []Activity {
    s.mu.RLock()
    defer s.mu.RUnlock()
    var activities []Activity
    for _, activity := range s.activities {
        activities = append(activities, activity)
    }
    return activities
}

// Invoice CRUD operations
func (s *Store) CreateInvoice(invoice Invoice) error {
    s.mu.Lock()
    defer s.mu.Unlock()
    invoice.ID = generateID()
    invoice.InvoiceNumber = s.generateInvoiceNumber()
    invoice.CreatedAt = time.Now()
    invoice.UpdatedAt = time.Now()
    s.invoices[invoice.ID] = invoice
    return nil
}

func (s *Store) GetInvoice(id string) (*Invoice, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    invoice, exists := s.invoices[id]
    if !exists {
        return nil, ErrNotFound
    }
    return &invoice, nil
}

func (s *Store) UpdateInvoice(id string, invoice Invoice) error {
    s.mu.Lock()
    defer s.mu.Unlock()
    existing, exists := s.invoices[id]
    if !exists {
        return ErrNotFound
    }
    invoice.ID = id
    invoice.InvoiceNumber = existing.InvoiceNumber
    invoice.CreatedAt = existing.CreatedAt
    invoice.UpdatedAt = time.Now()
    s.invoices[id] = invoice
    return nil
}

func (s *Store) DeleteInvoice(id string) error {
    s.mu.Lock()
    defer s.mu.Unlock()
    delete(s.invoices, id)
    return nil
}

func (s *Store) ListInvoices() []Invoice {
    s.mu.RLock()
    defer s.mu.RUnlock()
    var invoices []Invoice
    for _, invoice := range s.invoices {
        invoices = append(invoices, invoice)
    }
    return invoices
}

// Project CRUD operations
func (s *Store) CreateProject(project Project) error {
    s.mu.Lock()
    defer s.mu.Unlock()
    project.ID = generateID()
    project.CreatedAt = time.Now()
    project.UpdatedAt = time.Now()
    s.projects[project.ID] = project
    return nil
}

func (s *Store) GetProject(id string) (*Project, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    project, exists := s.projects[id]
    if !exists {
        return nil, ErrNotFound
    }
    return &project, nil
}

func (s *Store) UpdateProject(id string, project Project) error {
    s.mu.Lock()
    defer s.mu.Unlock()
    existing, exists := s.projects[id]
    if !exists {
        return ErrNotFound
    }
    project.ID = id
    project.CreatedAt = existing.CreatedAt
    project.UpdatedAt = time.Now()
    s.projects[id] = project
    return nil
}

func (s *Store) DeleteProject(id string) error {
    s.mu.Lock()
    defer s.mu.Unlock()
    delete(s.projects, id)
    return nil
}

func (s *Store) ListProjects() []Project {
    s.mu.RLock()
    defer s.mu.RUnlock()
    var projects []Project
    for _, project := range s.projects {
        projects = append(projects, project)
    }
    return projects
}

// Invoice workflow functions
func (s *Store) CreateInvoiceWithWorkflow(invoice Invoice) error {
    s.mu.Lock()
    defer s.mu.Unlock()
    invoice.ID = generateID()
    invoice.InvoiceNumber = s.generateInvoiceNumber()
    invoice.Status = InvoiceStatusDraft
    invoice.CreatedAt = time.Now()
    invoice.UpdatedAt = time.Now()
    s.invoices[invoice.ID] = invoice
    return nil
}

func (s *Store) SendInvoice(id string) error {
    s.mu.Lock()
    defer s.mu.Unlock()
    invoice, exists := s.invoices[id]
    if !exists {
        return ErrNotFound
    }
    invoice.Status = InvoiceStatusSent
    invoice.UpdatedAt = time.Now()
    s.invoices[id] = invoice
    return nil
}

func (s *Store) MarkAsPaid(id string) error {
    s.mu.Lock()
    defer s.mu.Unlock()
    invoice, exists := s.invoices[id]
    if !exists {
        return ErrNotFound
    }
    invoice.Status = InvoiceStatusPaid
    invoice.UpdatedAt = time.Now()
    s.invoices[id] = invoice
    return nil
}

func (s *Store) VoidInvoice(id string) error {
    s.mu.Lock()
    defer s.mu.Unlock()
    invoice, exists := s.invoices[id]
    if !exists {
        return ErrNotFound
    }
    invoice.Status = InvoiceStatusVoid
    invoice.UpdatedAt = time.Now()
    s.invoices[id] = invoice
    return nil
}

// Helper functions
func (s *Store) generateInvoiceNumber() string {
    // Simple implementation - in real app would query DB for max number
    return "INV-001"
}

func generateID() string {
    // Simple ID generator - in real app would use UUID or database IDs
    return "id-" + time.Now().Format("20060102150405")
}

// Error types
var (
    ErrNotFound = &storeError{"not found"}
)

type storeError struct {
    msg string
}

func (e *storeError) Error() string {
    return e.msg
}

// GetClientBalance returns the total unpaid balance for a client
func (s *Store) GetClientBalance(clientID string) (float64, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    
    var balance float64
    for _, invoice := range s.invoices {
        if invoice.ClientID == clientID && 
           (invoice.Status == InvoiceStatusSent || 
            invoice.Status == InvoiceStatusOverdue) {
            balance += invoice.Total
        }
    }
    return balance, nil
}
```