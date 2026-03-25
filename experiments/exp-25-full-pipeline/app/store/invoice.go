package store

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Invoice represents an invoice with auto-calculated totals
type Invoice struct {
	ID        string     `json:"id"`
	ClientID  string     `json:"client_id"`
	Number    string     `json:"number"`
	IssueDate time.Time  `json:"issue_date"`
	DueDate   time.Time  `json:"due_date"`
	Status    string     `json:"status"` // "draft", "sent", "paid", "overdue"
	LineItems []LineItem `json:"line_items"`
	TaxRate   float64    `json:"tax_rate"`
	Subtotal  float64    `json:"subtotal"`
	TaxAmount float64    `json:"tax_amount"`
	Total     float64    `json:"total"`
	Notes     string     `json:"notes"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
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
	dataDir string
}

// Constants for status values
const (
	StatusDraft   = "draft"
	StatusSent    = "sent"
	StatusPaid    = "paid"
	StatusOverdue = "overdue"
)

// Constants for backup retention
const (
	BackupRetentionDays = 30
)

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

	// DeleteInvoice removes an invoice
	// Behavior: Returns error if invoice not found
	DeleteInvoice(id string) error

	// ListInvoices returns all invoices
	ListInvoices() ([]*Invoice, error)

	// GetInvoicesByClient returns invoices for a specific client
	GetInvoicesByClient(clientID string) ([]*Invoice, error)

	// GetInvoiceByNumber returns an invoice by its number
	GetInvoiceByNumber(number string) (*Invoice, error)
}

// ClientStore interface defines the operations for clients
type ClientStore interface {
	// CreateClient creates a new client
	CreateClient(client *Client) error

	// GetClient retrieves a client by ID
	GetClient(id string) (*Client, error)

	// UpdateClient updates an existing client
	UpdateClient(client *Client) error

	// DeleteClient removes a client
	DeleteClient(id string) error

	// ListClients returns all clients
	ListClients() ([]*Client, error)
}

// NewStore creates a new Store instance
// Behavior: Initializes the store with file-based persistence
// Behavior: Creates data directories if they don't exist
// Behavior: Loads existing data from disk into memory cache
func NewStore(dataDir string) (*Store, error) {
	store := &Store{
		dataDir: dataDir,
	}

	// Create data directories if they don't exist
	err := os.MkdirAll(filepath.Join(dataDir, "invoices"), 0755)
	if err != nil {
		return nil, fmt.Errorf("failed to create invoices directory: %w", err)
	}

	err = os.MkdirAll(filepath.Join(dataDir, "clients"), 0755)
	if err != nil {
		return nil, fmt.Errorf("failed to create clients directory: %w", err)
	}

	err = os.MkdirAll(filepath.Join(dataDir, "backups"), 0755)
	if err != nil {
		return nil, fmt.Errorf("failed to create backups directory: %w", err)
	}

	return store, nil
}

// CreateInvoice creates a new invoice
// Behavior: Generates unique ID and sets timestamps
// Behavior: Calculates subtotal, tax, and total based on line items and tax rate
// Behavior: Validates input data before saving
func (s *Store) CreateInvoice(inv *Invoice) error {
	if inv == nil {
		return fmt.Errorf("invoice cannot be nil")
	}

	// Set default status if not provided
	if inv.Status == "" {
		inv.Status = StatusDraft
	}

	// Set timestamps
	now := time.Now()
	inv.CreatedAt = now
	inv.UpdatedAt = now

	// Generate unique ID if not provided
	if inv.ID == "" {
		inv.ID = generateID()
	}

	// Generate invoice number if not provided
	if inv.Number == "" {
		inv.Number = generateInvoiceNumber()
	}

	// Validate invoice
	if err := s.validateInvoice(inv); err != nil {
		return fmt.Errorf("invalid invoice: %w", err)
	}

	// Calculate totals
	s.calculateTotals(inv)

	// Save to file
	filePath := filepath.Join(s.dataDir, "invoices", inv.ID+".json")
	data, err := json.MarshalIndent(inv, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal invoice: %w", err)
	}

	err = ioutil.WriteFile(filePath, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write invoice file: %w", err)
	}

	return nil
}

// GetInvoice retrieves an invoice by ID
// Behavior: Returns error if invoice not found
func (s *Store) GetInvoice(id string) (*Invoice, error) {
	if id == "" {
		return nil, fmt.Errorf("invoice ID cannot be empty")
	}

	filePath := filepath.Join(s.dataDir, "invoices", id+".json")
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("invoice not found")
		}
		return nil, fmt.Errorf("failed to read invoice file: %w", err)
	}

	var inv Invoice
	err = json.Unmarshal(data, &inv)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal invoice: %w", err)
	}

	return &inv, nil
}

// UpdateInvoice updates an existing invoice
// Behavior: Updates timestamps
// Behavior: Recalculates totals
// Behavior: Validates input data
func (s *Store) UpdateInvoice(inv *Invoice) error {
	if inv == nil || inv.ID == "" {
		return fmt.Errorf("invoice cannot be nil or have empty ID")
	}

	// Validate invoice
	if err := s.validateInvoice(inv); err != nil {
		return fmt.Errorf("invalid invoice: %w", err)
	}

	// Set updated timestamp
	inv.UpdatedAt = time.Now()

	// Calculate totals
	s.calculateTotals(inv)

	// Save to file
	filePath := filepath.Join(s.dataDir, "invoices", inv.ID+".json")
	data, err := json.MarshalIndent(inv, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal invoice: %w", err)
	}

	err = ioutil.WriteFile(filePath, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write invoice file: %w", err)
	}

	return nil
}

// DeleteInvoice removes an invoice
// Behavior: Returns error if invoice not found
func (s *Store) DeleteInvoice(id string) error {
	if id == "" {
		return fmt.Errorf("invoice ID cannot be empty")
	}

	filePath := filepath.Join(s.dataDir, "invoices", id+".json")
	err := os.Remove(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("invoice not found")
		}
		return fmt.Errorf("failed to delete invoice file: %w", err)
	}

	return nil
}

// ListInvoices returns all invoices
func (s *Store) ListInvoices() ([]*Invoice, error) {
	var invoices []*Invoice

	dirPath := filepath.Join(s.dataDir, "invoices")
	entries, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read invoices directory: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".json") {
			filePath := filepath.Join(dirPath, entry.Name())
			data, err := ioutil.ReadFile(filePath)
			if err != nil {
				return nil, fmt.Errorf("failed to read invoice file %s: %w", entry.Name(), err)
			}

			var inv Invoice
			err = json.Unmarshal(data, &inv)
			if err != nil {
				return nil, fmt.Errorf("failed to unmarshal invoice file %s: %w", entry.Name(), err)
			}

			invoices = append(invoices, &inv)
		}
	}

	return invoices, nil
}

// GetInvoicesByClient returns invoices for a specific client
func (s *Store) GetInvoicesByClient(clientID string) ([]*Invoice, error) {
	if clientID == "" {
		return nil, fmt.Errorf("client ID cannot be empty")
	}

	allInvoices, err := s.ListInvoices()
	if err != nil {
		return nil, err
	}

	var clientInvoices []*Invoice
	for _, inv := range allInvoices {
		if inv.ClientID == clientID {
			clientInvoices = append(clientInvoices, inv)
		}
	}

	return clientInvoices, nil
}

// GetInvoiceByNumber returns an invoice by its number
func (s *Store) GetInvoiceByNumber(number string) (*Invoice, error) {
	if number == "" {
		return nil, fmt.Errorf("invoice number cannot be empty")
	}

	allInvoices, err := s.ListInvoices()
	if err != nil {
		return nil, err
	}

	for _, inv := range allInvoices {
		if inv.Number == number {
			return inv, nil
		}
	}

	return nil, fmt.Errorf("invoice with number %s not found", number)
}

// CreateClient creates a new client
func (s *Store) CreateClient(client *Client) error {
	if client == nil {
		return fmt.Errorf("client cannot be nil")
	}

	// Set timestamps
	now := time.Now()
	client.CreatedAt = now
	client.UpdatedAt = now

	// Generate unique ID if not provided
	if client.ID == "" {
		client.ID = generateID()
	}

	// Validate client
	if err := s.validateClient(client); err != nil {
		return fmt.Errorf("invalid client: %w", err)
	}

	// Save to file
	filePath := filepath.Join(s.dataDir, "clients", client.ID+".json")
	data, err := json.MarshalIndent(client, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal client: %w", err)
	}

	err = ioutil.WriteFile(filePath, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write client file: %w", err)
	}

	return nil
}

// GetClient retrieves a client by ID
func (s *Store) GetClient(id string) (*Client, error) {
	if id == "" {
		return nil, fmt.Errorf("client ID cannot be empty")
	}

	filePath := filepath.Join(s.dataDir, "clients", id+".json")
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("client not found")
		}
		return nil, fmt.Errorf("failed to read client file: %w", err)
	}

	var client Client
	err = json.Unmarshal(data, &client)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal client: %w", err)
	}

	return &client, nil
}

// UpdateClient updates an existing client
func (s *Store) UpdateClient(client *Client) error {
	if client == nil || client.ID == "" {
		return fmt.Errorf("client cannot be nil or have empty ID")
	}

	// Validate client
	if err := s.validateClient(client); err != nil {
		return fmt.Errorf("invalid client: %w", err)
	}

	// Set updated timestamp
	client.UpdatedAt = time.Now()

	// Save to file
	filePath := filepath.Join(s.dataDir, "clients", client.ID+".json")
	data, err := json.MarshalIndent(client, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal client: %w", err)
	}

	err = ioutil.WriteFile(filePath, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write client file: %w", err)
	}

	return nil
}

// DeleteClient removes a client
func (s *Store) DeleteClient(id string) error {
	if id == "" {
		return fmt.Errorf("client ID cannot be empty")
	}

	filePath := filepath.Join(s.dataDir, "clients", id+".json")
	err := os.Remove(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("client not found")
		}
		return fmt.Errorf("failed to delete client file: %w", err)
	}

	return nil
}

// ListClients returns all clients
func (s *Store) ListClients() ([]*Client, error) {
	var clients []*Client

	dirPath := filepath.Join(s.dataDir, "clients")
	entries, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read clients directory: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".json") {
			filePath := filepath.Join(dirPath, entry.Name())
			data, err := ioutil.ReadFile(filePath)
			if err != nil {
				return nil, fmt.Errorf("failed to read client file %s: %w", entry.Name(), err)
			}

			var client Client
			err = json.Unmarshal(data, &client)
			if err != nil {
				return nil, fmt.Errorf("failed to unmarshal client file %s: %w", entry.Name(), err)
			}

			clients = append(clients, &client)
		}
	}

	return clients, nil
}

func (s *Store) validateInvoice(inv *Invoice) error {
	if inv.ClientID == "" {
		return fmt.Errorf("client ID is required")
	}

	if inv.Number == "" {
		return fmt.Errorf("invoice number is required")
	}

	if inv.IssueDate.IsZero() {
		return fmt.Errorf("issue date is required")
	}

	if inv.DueDate.IsZero() {
		return fmt.Errorf("due date is required")
	}

	if inv.TaxRate < 0 {
		return fmt.Errorf("tax rate cannot be negative")
	}

	for _, item := range inv.LineItems {
		if item.Quantity <= 0 {
			return fmt.Errorf("line item quantity must be positive")
		}
		if item.UnitPrice < 0 {
			return fmt.Errorf("line item unit price cannot be negative")
		}
	}

	return nil
}

func (s *Store) validateClient(client *Client) error {
	if client.Name == "" {
		return fmt.Errorf("client name is required")
	}

	if client.Email == "" {
		return fmt.Errorf("client email is required")
	}

	return nil
}

func (s *Store) calculateTotals(inv *Invoice) {
	// Calculate subtotal
	inv.Subtotal = 0
	for _, item := range inv.LineItems {
		item.Total = float64(item.Quantity) * item.UnitPrice
		inv.Subtotal += item.Total
	}

	// Calculate tax amount
	inv.TaxAmount = inv.Subtotal * (inv.TaxRate / 100)

	// Calculate total
	inv.Total = inv.Subtotal + inv.TaxAmount
}

func generateID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

func generateInvoiceNumber() string {
	return fmt.Sprintf("INV-%d", time.Now().Unix())
}
