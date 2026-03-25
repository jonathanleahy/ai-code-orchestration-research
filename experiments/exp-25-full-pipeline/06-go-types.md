# Go Type Specification

```go
package store

import (
	"time"
)

// Invoice represents an invoice with auto-calculated totals
type Invoice struct {
	ID            string    `json:"id"`
	ClientID      string    `json:"client_id"`
	Number        string    `json:"number"`
	IssueDate     time.Time `json:"issue_date"`
	DueDate       time.Time `json:"due_date"`
	Status        string    `json:"status"` // "draft", "sent", "paid", "overdue"
	LineItems     []LineItem `json:"line_items"`
	TaxRate       float64   `json:"tax_rate"`
	Subtotal      float64   `json:"subtotal"`
	TaxAmount     float64   `json:"tax_amount"`
	Total         float64   `json:"total"`
	Notes         string    `json:"notes"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// Client represents a client with payment history
type Client struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	Phone        string    `json:"phone"`
	Address      string    `json:"address"`
	PaymentTerms int       `json:"payment_terms"` // in days
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// LineItem represents a line item in an invoice
type LineItem struct {
	ID          string  `json:"id"`
	Description string  `json:"description"`
	Quantity    int     `json:"quantity"`
	UnitPrice   float64 `json:"unit_price"`
	Total       float64 `json:"total"`
}

// Store manages invoices and clients with file-based persistence
type Store struct {
	// Behavior: All data is persisted to JSON files in the configured directory
	// Behavior: Backup strategy is implemented with daily rotation (30-day retention)
	// Behavior: Data validation is performed on all CRUD operations
	// Behavior: Error handling includes retry logic and logging
	// Behavior: GDPR compliance is enforced through data minimization and deletion capabilities
}

// Constants for status values
const (
	StatusDraft    = "draft"
	StatusSent     = "sent"
	StatusPaid     = "paid"
	StatusOverdue  = "overdue"
)

// Constants for backup retention
const (
	BackupRetentionDays = 30
)

// NewStore creates a new Store instance
// Behavior: Initializes the store with file-based persistence
// Behavior: Creates data directories if they don't exist
// Behavior: Loads existing data from disk into memory cache
func NewStore(dataDir string) (*Store, error) {
	// Implementation would initialize the store with file persistence
	// and load existing data
	return &Store{}, nil
}

// InvoiceStore interface defines the operations for invoices
type InvoiceStore interface {
	// CreateInvoice creates a new invoice
	// Behavior: Generates unique ID and sets timestamps
	// Behavior: Calculates subtotal, tax, and total based on line items and tax rate
	// Behavior: Validates input data before saving
	CreateInvoice(inv *Invoice) error

	// GetInvoice retrieves an invoice by ID
	// Behavior: Returns error if invoice not found
	GetInvoice(id string) (*Invoice, error)

	// UpdateInvoice updates an existing invoice
	// Behavior: Updates timestamps
	// Behavior: Recalculates totals
	// Behavior: Validates input data
	UpdateInvoice(inv *Invoice) error

	// DeleteInvoice removes an invoice by ID
	// Behavior: Returns error if invoice not found
	DeleteInvoice(id string) error

	// ListInvoices returns all invoices
	// Behavior: Optionally filters by status or date range
	ListInvoices() ([]*Invoice, error)

	// GetInvoicesByClient returns invoices for a specific client
	GetInvoicesByClient(clientID string) ([]*Invoice, error)
}

// ClientStore interface defines the operations for clients
type ClientStore interface {
	// CreateClient creates a new client
	// Behavior: Generates unique ID and sets timestamps
	// Behavior: Validates input data before saving
	CreateClient(client *Client) error

	// GetClient retrieves a client by ID
	// Behavior: Returns error if client not found
	GetClient(id string) (*Client, error)

	// UpdateClient updates an existing client
	// Behavior: Updates timestamps
	// Behavior: Validates input data
	UpdateClient(client *Client) error

	// DeleteClient removes a client by ID
	// Behavior: Returns error if client not found
	// Behavior: Also deletes associated invoices
	DeleteClient(id string) error

	// ListClients returns all clients
	ListClients() ([]*Client, error)
}

// Store implements both InvoiceStore and ClientStore interfaces
var _ InvoiceStore = (*Store)(nil)
var _ ClientStore = (*Store)(nil)

// CreateInvoice creates a new invoice
func (s *Store) CreateInvoice(inv *Invoice) error {
	// Implementation would save invoice to JSON file
	// and update in-memory cache
	return nil
}

// GetInvoice retrieves an invoice by ID
func (s *Store) GetInvoice(id string) (*Invoice, error) {
	// Implementation would retrieve invoice from JSON file or cache
	return &Invoice{}, nil
}

// UpdateInvoice updates an existing invoice
func (s *Store) UpdateInvoice(inv *Invoice) error {
	// Implementation would update invoice in JSON file
	// and refresh in-memory cache
	return nil
}

// DeleteInvoice removes an invoice by ID
func (s *Store) DeleteInvoice(id string) error {
	// Implementation would delete invoice file
	// and remove from in-memory cache
	return nil
}

// ListInvoices returns all invoices
func (s *Store) ListInvoices() ([]*Invoice, error) {
	// Implementation would read all invoice files
	// and return in-memory cache data
	return []*Invoice{}, nil
}

// GetInvoicesByClient returns invoices for a specific client
func (s *Store) GetInvoicesByClient(clientID string) ([]*Invoice, error) {
	// Implementation would filter invoices by client ID
	return []*Invoice{}, nil
}

// CreateClient creates a new client
func (s *Store) CreateClient(client *Client) error {
	// Implementation would save client to JSON file
	// and update in-memory cache
	return nil
}

// GetClient retrieves a client by ID
func (s *Store) GetClient(id string) (*Client, error) {
	// Implementation would retrieve client from JSON file or cache
	return &Client{}, nil
}

// UpdateClient updates an existing client
func (s *Store) UpdateClient(client *Client) error {
	// Implementation would update client in JSON file
	// and refresh in-memory cache
	return nil
}

// DeleteClient removes a client by ID
func (s *Store) DeleteClient(id string) error {
	// Implementation would delete client file
	// and remove from in-memory cache
	// Also deletes associated invoices
	return nil
}

// ListClients returns all clients
func (s *Store) ListClients() ([]*Client, error) {
	// Implementation would read all client files
	// and return in-memory cache data
	return []*Client{}, nil
}

// BackupData creates a backup of all data
// Behavior: Creates timestamped backup files with rotation
// Behavior: Maintains 30-day retention policy
func (s *Store) BackupData() error {
	// Implementation would create backup files
	// and handle rotation of old backups
	return nil
}

// ValidateInvoice validates an invoice for data integrity
// Behavior: Checks required fields and calculates totals
func (s *Store) ValidateInvoice(inv *Invoice) error {
	// Implementation would validate invoice data
	return nil
}

// ValidateClient validates a client for data integrity
// Behavior: Checks required fields
func (s *Store) ValidateClient(client *Client) error {
	// Implementation would validate client data
	return nil
}
```