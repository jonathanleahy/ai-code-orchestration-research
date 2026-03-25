package store

import (
	"sync"
	"time"
)

// Client represents a client in the CRM
type Client struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Address   Address   `json:"address"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Address represents a client's address
type Address struct {
	Street  string `json:"street"`
	City    string `json:"city"`
	State   string `json:"state"`
	ZipCode string `json:"zip_code"`
	Country string `json:"country"`
}

// Invoice represents an invoice
type Invoice struct {
	ID          string    `json:"id"`
	ClientID    string    `json:"client_id"`
	Items       []Item    `json:"items"`
	TotalAmount float64   `json:"total_amount"`
	IssueDate   time.Time `json:"issue_date"`
	DueDate     time.Time `json:"due_date"`
	Status      string    `json:"status"` // "draft", "sent", "paid"
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Item represents an item on an invoice
type Item struct {
	Description string  `json:"description"`
	Quantity    int     `json:"quantity"`
	Price       float64 `json:"price"`
	Total       float64 `json:"total"`
}

// Comment represents a note/comment on a client
type Comment struct {
	ID        string    `json:"id"`
	ClientID  string    `json:"client_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

// HistoryEntry represents a client interaction log
type HistoryEntry struct {
	ID        string    `json:"id"`
	ClientID  string    `json:"client_id"`
	Note      string    `json:"note"`
	CreatedAt time.Time `json:"created_at"`
}

// Store interface defines the methods for data persistence
type Store interface {
	// Client methods
	GetClient(id string) (*Client, error)
	GetClients() ([]*Client, error)
	CreateClient(client *Client) error
	UpdateClient(id string, client *Client) error
	DeleteClient(id string) error
	SearchClients(query string) ([]*Client, error)

	// Invoice methods
	GetInvoice(id string) (*Invoice, error)
	GetInvoices() ([]*Invoice, error)
	CreateInvoice(invoice *Invoice) error
	UpdateInvoice(id string, invoice *Invoice) error
	DeleteInvoice(id string) error

	// Comment methods
	GetComments(clientID string) ([]*Comment, error)
	CreateComment(comment *Comment) error
	DeleteComment(id string) error

	// History methods
	GetHistory(clientID string) ([]*HistoryEntry, error)
	CreateHistory(entry *HistoryEntry) error

	// Persistence
	Save() error
	Load() error
}

// store implements the Store interface
type store struct {
	clients  map[string]*Client
	invoices map[string]*Invoice
	comments map[string]*Comment
	history  map[string]*HistoryEntry
	mu       sync.RWMutex
}

// NewStore creates a new in-memory store
func NewStore() Store {
	return &store{
		clients:  make(map[string]*Client),
		invoices: make(map[string]*Invoice),
		comments: make(map[string]*Comment),
		history:  make(map[string]*HistoryEntry),
	}
}

// GetClient retrieves a client by ID
func (s *store) GetClient(id string) (*Client, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	client, exists := s.clients[id]
	if !exists {
		return nil, nil
	}
	return client, nil
}

// GetClients retrieves all clients
func (s *store) GetClients() ([]*Client, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var clients []*Client
	for _, client := range s.clients {
		clients = append(clients, client)
	}
	return clients, nil
}

// CreateClient creates a new client
func (s *store) CreateClient(client *Client) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.clients[client.ID] = client
	return nil
}

// UpdateClient updates an existing client
func (s *store) UpdateClient(id string, client *Client) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.clients[id]; !exists {
		return nil
	}
	client.ID = id
	client.UpdatedAt = time.Now()
	s.clients[id] = client
	return nil
}

// DeleteClient deletes a client
func (s *store) DeleteClient(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.clients, id)
	return nil
}

// SearchClients searches clients by name or email
func (s *store) SearchClients(query string) ([]*Client, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var results []*Client
	for _, client := range s.clients {
		if client.Name == query || client.Email == query {
			results = append(results, client)
		}
	}
	return results, nil
}

// GetInvoice retrieves an invoice by ID
func (s *store) GetInvoice(id string) (*Invoice, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	invoice, exists := s.invoices[id]
	if !exists {
		return nil, nil
	}
	return invoice, nil
}

// GetInvoices retrieves all invoices
func (s *store) GetInvoices() ([]*Invoice, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var invoices []*Invoice
	for _, invoice := range s.invoices {
		invoices = append(invoices, invoice)
	}
	return invoices, nil
}

// CreateInvoice creates a new invoice
func (s *store) CreateInvoice(invoice *Invoice) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.invoices[invoice.ID] = invoice
	return nil
}

// UpdateInvoice updates an existing invoice
func (s *store) UpdateInvoice(id string, invoice *Invoice) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.invoices[id]; !exists {
		return nil
	}
	invoice.ID = id
	invoice.UpdatedAt = time.Now()
	s.invoices[id] = invoice
	return nil
}

// DeleteInvoice deletes an invoice
func (s *store) DeleteInvoice(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.invoices, id)
	return nil
}

// GetComments retrieves all comments for a client
func (s *store) GetComments(clientID string) ([]*Comment, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var comments []*Comment
	for _, comment := range s.comments {
		if comment.ClientID == clientID {
			comments = append(comments, comment)
		}
	}
	return comments, nil
}

// CreateComment creates a new comment
func (s *store) CreateComment(comment *Comment) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.comments[comment.ID] = comment
	return nil
}

// DeleteComment deletes a comment
func (s *store) DeleteComment(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.comments, id)
	return nil
}

// GetHistory retrieves all history entries for a client
func (s *store) GetHistory(clientID string) ([]*HistoryEntry, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var history []*HistoryEntry
	for _, entry := range s.history {
		if entry.ClientID == clientID {
			history = append(history, entry)
		}
	}
	return history, nil
}

// CreateHistory creates a new history entry
func (s *store) CreateHistory(entry *HistoryEntry) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.history[entry.ID] = entry
	return nil
}

// Save persists the data (stub for future implementation)
func (s *store) Save() error {
	// In-memory store doesn't need to save
	return nil
}

// Load loads the data (stub for future implementation)
func (s *store) Load() error {
	// In-memory store doesn't need to load
	return nil
}
