package store

import (
	"sync"
	"time"
)

type Store interface {
	// Client CRUD
	CreateClient(client *Client) error
	GetClient(id string) (*Client, error)
	UpdateClient(id string, client *Client) error
	DeleteClient(id string) error
	SearchClients(query string) ([]*Client, error)

	// Activity CRUD
	CreateActivity(activity *Activity) error
	GetActivity(id string) (*Activity, error)
	UpdateActivity(id string, activity *Activity) error
	DeleteActivity(id string) error
	GetClientActivities(clientID string) ([]*Activity, error)

	// Invoice CRUD
	CreateInvoice(invoice *Invoice) error
	GetInvoice(id string) (*Invoice, error)
	UpdateInvoice(id string, invoice *Invoice) error
	DeleteInvoice(id string) error
	GetClientInvoices(clientID string) ([]*Invoice, error)

	// Invoice Status Updates
	MarkInvoicePaid(id string) error
	VoidInvoice(id string) error
}

type InMemoryStore struct {
	clients    map[string]*Client
	activities map[string]*Activity
	invoices   map[string]*Invoice
	mu         sync.RWMutex
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		clients:    make(map[string]*Client),
		activities: make(map[string]*Activity),
		invoices:   make(map[string]*Invoice),
	}
}

// Client CRUD
func (s *InMemoryStore) CreateClient(client *Client) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	client.CreatedAt = time.Now()
	s.clients[client.ID] = client
	return nil
}

func (s *InMemoryStore) GetClient(id string) (*Client, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	client, exists := s.clients[id]
	if !exists {
		return nil, nil
	}
	return client, nil
}

func (s *InMemoryStore) UpdateClient(id string, client *Client) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, exists := s.clients[id]
	if !exists {
		return nil
	}

	client.CreatedAt = time.Now()
	s.clients[id] = client
	return nil
}

func (s *InMemoryStore) DeleteClient(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.clients, id)
	return nil
}

func (s *InMemoryStore) SearchClients(query string) ([]*Client, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var results []*Client
	for _, client := range s.clients {
		if client.Name == query || client.Email == query || client.Company == query {
			results = append(results, client)
		}
	}
	return results, nil
}

// Activity CRUD
func (s *InMemoryStore) CreateActivity(activity *Activity) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	activity.CreatedAt = time.Now()
	s.activities[activity.ID] = activity
	return nil
}

func (s *InMemoryStore) GetActivity(id string) (*Activity, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	activity, exists := s.activities[id]
	if !exists {
		return nil, nil
	}
	return activity, nil
}

func (s *InMemoryStore) UpdateActivity(id string, activity *Activity) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, exists := s.activities[id]
	if !exists {
		return nil
	}

	activity.CreatedAt = time.Now()
	s.activities[id] = activity
	return nil
}

func (s *InMemoryStore) DeleteActivity(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.activities, id)
	return nil
}

func (s *InMemoryStore) GetClientActivities(clientID string) ([]*Activity, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var activities []*Activity
	for _, activity := range s.activities {
		if activity.ClientID == clientID {
			activities = append(activities, activity)
		}
	}
	return activities, nil
}

// Invoice CRUD
func (s *InMemoryStore) CreateInvoice(invoice *Invoice) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	invoice.CreatedAt = time.Now()
	for i := range invoice.Items {
		invoice.Items[i].Amount = float64(invoice.Items[i].Quantity) * invoice.Items[i].Rate
	}
	s.invoices[invoice.ID] = invoice
	return nil
}

func (s *InMemoryStore) GetInvoice(id string) (*Invoice, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	invoice, exists := s.invoices[id]
	if !exists {
		return nil, nil
	}
	return invoice, nil
}

func (s *InMemoryStore) UpdateInvoice(id string, invoice *Invoice) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, exists := s.invoices[id]
	if !exists {
		return nil
	}

	invoice.CreatedAt = time.Now()
	for i := range invoice.Items {
		invoice.Items[i].Amount = float64(invoice.Items[i].Quantity) * invoice.Items[i].Rate
	}
	s.invoices[id] = invoice
	return nil
}

func (s *InMemoryStore) DeleteInvoice(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.invoices, id)
	return nil
}

func (s *InMemoryStore) GetClientInvoices(clientID string) ([]*Invoice, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var invoices []*Invoice
	for _, invoice := range s.invoices {
		if invoice.ClientID == clientID {
			invoices = append(invoices, invoice)
		}
	}
	return invoices, nil
}

// Invoice Status Updates
func (s *InMemoryStore) MarkInvoicePaid(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	invoice, exists := s.invoices[id]
	if !exists {
		return nil
	}

	invoice.Status = "paid"
	return nil
}

func (s *InMemoryStore) VoidInvoice(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	invoice, exists := s.invoices[id]
	if !exists {
		return nil
	}

	invoice.Status = "void"
	return nil
}
