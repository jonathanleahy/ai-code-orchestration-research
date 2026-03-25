package store

import (
	"sync"
	"time"
)

type Client struct {
	ID        string
	Name      string
	Email     string
	CreatedAt time.Time
}

type LineItem struct {
	Description string
	Amount      float64
}

type Invoice struct {
	ID        string
	ClientID  string
	Status    string // draft, sent, paid, void
	Total     float64
	LineItems []LineItem
	CreatedAt time.Time
}

type Store struct {
	mu       sync.RWMutex
	clients  map[string]*Client
	invoices map[string]*Invoice
	nextID   int
}

func NewStore() *Store {
	return &Store{
		clients:  make(map[string]*Client),
		invoices: make(map[string]*Invoice),
		nextID:   1000,
	}
}

func (s *Store) genID() string {
	s.nextID++
	return string(rune(s.nextID))
}

func (s *Store) CreateClient(c *Client) {
	s.mu.Lock()
	defer s.mu.Unlock()
	c.ID = s.genID()
	c.CreatedAt = time.Now()
	s.clients[c.ID] = c
}

func (s *Store) GetClient(id string) (*Client, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	c, ok := s.clients[id]
	return c, ok
}

func (s *Store) ListClients() []*Client {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var clients []*Client
	for _, c := range s.clients {
		clients = append(clients, c)
	}
	return clients
}

func (s *Store) UpdateClient(id string, c *Client) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.clients[id]; !ok {
		return false
	}
	c.ID = id
	s.clients[id] = c
	return true
}

func (s *Store) DeleteClient(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.clients[id]; !ok {
		return false
	}
	delete(s.clients, id)
	return true
}

func (s *Store) CreateInvoice(inv *Invoice) {
	s.mu.Lock()
	defer s.mu.Unlock()
	inv.ID = s.genID()
	inv.CreatedAt = time.Now()
	s.invoices[inv.ID] = inv
}

func (s *Store) GetInvoice(id string) (*Invoice, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	inv, ok := s.invoices[id]
	return inv, ok
}

func (s *Store) ListInvoicesByClient(clientID string) []*Invoice {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var invoices []*Invoice
	for _, inv := range s.invoices {
		if inv.ClientID == clientID {
			invoices = append(invoices, inv)
		}
	}
	return invoices
}

func (s *Store) UpdateInvoice(id string, inv *Invoice) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.invoices[id]; !ok {
		return false
	}
	inv.ID = id
	s.invoices[id] = inv
	return true
}

func (s *Store) DeleteInvoice(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.invoices[id]; !ok {
		return false
	}
	delete(s.invoices, id)
	return true
}
