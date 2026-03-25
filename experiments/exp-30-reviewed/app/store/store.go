package store

import (
	"errors"
	"strconv"
	"sync"
	"time"
)

// Client represents a freelancer's client
type Client struct {
	ID               string    `json:"id"`
	Name             string    `json:"name"`
	Email            string    `json:"email"`
	Phone            string    `json:"phone"`
	PreferredPayment string    `json:"preferred_payment"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// Activity represents an interaction with a client
type Activity struct {
	ID        string    `json:"id"`
	ClientID  string    `json:"client_id"`
	Type      string    `json:"type"` // email, call, meeting, etc.
	Details   string    `json:"details"`
	Timestamp time.Time `json:"timestamp"`
}

// Invoice represents a client invoice
type Invoice struct {
	ID        string    `json:"id"`
	ClientID  string    `json:"client_id"`
	Amount    float64   `json:"amount"`
	Status    string    `json:"status"` // paid, unpaid, due
	DueDate   time.Time `json:"due_date"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Store holds all data in memory with thread safety
type Store struct {
	clients    map[string]*Client
	activities map[string]*Activity
	invoices   map[string]*Invoice
	mutex      sync.RWMutex
}

var (
	ErrClientNotFound   = errors.New("client not found")
	ErrDuplicateClient  = errors.New("client with this name already exists")
	ErrActivityNotFound = errors.New("activity not found")
	ErrInvoiceNotFound  = errors.New("invoice not found")
)

// NewStore creates and returns a new in-memory store
func NewStore() *Store {
	return &Store{
		clients:    make(map[string]*Client),
		activities: make(map[string]*Activity),
		invoices:   make(map[string]*Invoice),
	}
}

// CreateClient adds a new client to the store
func (s *Store) CreateClient(client *Client) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Check for duplicates by name
	for _, c := range s.clients {
		if c.Name == client.Name {
			return ErrDuplicateClient
		}
	}

	client.ID = generateID()
	client.CreatedAt = time.Now()
	client.UpdatedAt = time.Now()
	s.clients[client.ID] = client
	return nil
}

// GetClient retrieves a client by ID
func (s *Store) GetClient(id string) (*Client, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	client, exists := s.clients[id]
	if !exists {
		return nil, ErrClientNotFound
	}
	return client, nil
}

// UpdateClient updates an existing client
func (s *Store) UpdateClient(id string, client *Client) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	existing, exists := s.clients[id]
	if !exists {
		return ErrClientNotFound
	}

	// Check for duplicates by name (excluding self)
	for _, c := range s.clients {
		if c.Name == client.Name && c.ID != id {
			return ErrDuplicateClient
		}
	}

	client.ID = id
	client.CreatedAt = existing.CreatedAt
	client.UpdatedAt = time.Now()
	s.clients[id] = client
	return nil
}

// DeleteClient removes a client by ID
func (s *Store) DeleteClient(id string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, exists := s.clients[id]; !exists {
		return ErrClientNotFound
	}
	delete(s.clients, id)
	return nil
}

// ListClients returns all clients
func (s *Store) ListClients() []*Client {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	clients := make([]*Client, 0, len(s.clients))
	for _, client := range s.clients {
		clients = append(clients, client)
	}
	return clients
}

// CreateActivity adds a new activity to the store
func (s *Store) CreateActivity(activity *Activity) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	activity.ID = generateID()
	activity.Timestamp = time.Now()
	s.activities[activity.ID] = activity
	return nil
}

// GetActivity retrieves an activity by ID
func (s *Store) GetActivity(id string) (*Activity, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	activity, exists := s.activities[id]
	if !exists {
		return nil, ErrActivityNotFound
	}
	return activity, nil
}

// UpdateActivity updates an existing activity
func (s *Store) UpdateActivity(id string, activity *Activity) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	existing, exists := s.activities[id]
	if !exists {
		return ErrActivityNotFound
	}

	activity.ID = id
	activity.Timestamp = existing.Timestamp
	s.activities[id] = activity
	return nil
}

// DeleteActivity removes an activity by ID
func (s *Store) DeleteActivity(id string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, exists := s.activities[id]; !exists {
		return ErrActivityNotFound
	}
	delete(s.activities, id)
	return nil
}

// ListActivities returns all activities
func (s *Store) ListActivities() []*Activity {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	activities := make([]*Activity, 0, len(s.activities))
	for _, activity := range s.activities {
		activities = append(activities, activity)
	}
	return activities
}

// CreateInvoice adds a new invoice to the store
func (s *Store) CreateInvoice(invoice *Invoice) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	invoice.ID = generateID()
	invoice.CreatedAt = time.Now()
	invoice.UpdatedAt = time.Now()
	s.invoices[invoice.ID] = invoice
	return nil
}

// GetInvoice retrieves an invoice by ID
func (s *Store) GetInvoice(id string) (*Invoice, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	invoice, exists := s.invoices[id]
	if !exists {
		return nil, ErrInvoiceNotFound
	}
	return invoice, nil
}

// UpdateInvoice updates an existing invoice
func (s *Store) UpdateInvoice(id string, invoice *Invoice) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	existing, exists := s.invoices[id]
	if !exists {
		return ErrInvoiceNotFound
	}

	invoice.ID = id
	invoice.CreatedAt = existing.CreatedAt
	invoice.UpdatedAt = time.Now()
	s.invoices[id] = invoice
	return nil
}

// DeleteInvoice removes an invoice by ID
func (s *Store) DeleteInvoice(id string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, exists := s.invoices[id]; !exists {
		return ErrInvoiceNotFound
	}
	delete(s.invoices, id)
	return nil
}

// ListInvoices returns all invoices
func (s *Store) ListInvoices() []*Invoice {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	invoices := make([]*Invoice, 0, len(s.invoices))
	for _, invoice := range s.invoices {
		invoices = append(invoices, invoice)
	}
	return invoices
}

// generateID creates a simple unique ID
func generateID() string {
	return time.Now().Format("20060102150405") + "-" + strconv.Itoa(time.Now().Nanosecond())
}
