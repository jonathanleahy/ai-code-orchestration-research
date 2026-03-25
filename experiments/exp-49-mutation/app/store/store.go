package store

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

type Client struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	CreatedAt time.Time `json:"created_at"`
}

type Activity struct {
	ID          string    `json:"id"`
	ClientID    string    `json:"client_id"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type Invoice struct {
	ID          string    `json:"id"`
	ClientID    string    `json:"client_id"`
	Number      string    `json:"number"`
	Amount      float64   `json:"amount"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	DueDate     time.Time `json:"due_date"`
	CreatedAt   time.Time `json:"created_at"`
}

type Store struct {
	mu           sync.RWMutex
	clients      map[string]*Client
	activities   map[string][]*Activity
	invoices     map[string][]*Invoice
	counter      int
	actCounter   int
	invCounter   int
}

func NewStore() *Store {
	return &Store{
		clients:    make(map[string]*Client),
		activities: make(map[string][]*Activity),
		invoices:   make(map[string][]*Invoice),
		counter:    0,
		actCounter: 0,
		invCounter: 0,
	}
}

func (s *Store) AddClient(name, email, phone string) *Client {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.counter++
	id := time.Now().Format("20060102150405") + "-" + fmt.Sprintf("%04d", s.counter)

	client := &Client{
		ID:        id,
		Name:      name,
		Email:     email,
		Phone:     phone,
		CreatedAt: time.Now(),
	}

	s.clients[id] = client
	return client
}

func (s *Store) GetClient(id string) *Client {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.clients[id]
}

func (s *Store) ListClients() []*Client {
	s.mu.RLock()
	defer s.mu.RUnlock()

	clients := make([]*Client, 0, len(s.clients))
	for _, c := range s.clients {
		clients = append(clients, c)
	}
	return clients
}

func (s *Store) SearchClients(query string) []*Client {
	s.mu.RLock()
	defer s.mu.RUnlock()

	clients := make([]*Client, 0)
	query = strings.ToLower(query)
	for _, c := range s.clients {
		if strings.Contains(strings.ToLower(c.Name), query) || strings.Contains(strings.ToLower(c.Email), query) {
			clients = append(clients, c)
		}
	}
	return clients
}

func (s *Store) UpdateClient(id, name, email, phone string) *Client {
	s.mu.Lock()
	defer s.mu.Unlock()

	client := s.clients[id]
	if client == nil {
		return nil
	}

	client.Name = name
	client.Email = email
	client.Phone = phone
	return client
}

func (s *Store) DeleteClient(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.clients[id]; exists {
		delete(s.clients, id)
		return true
	}
	return false
}

func (s *Store) AddActivity(clientID, actType, description string) *Activity {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.clients[clientID]; !exists {
		return nil
	}

	s.actCounter++
	id := time.Now().Format("20060102150405") + "-act-" + fmt.Sprintf("%04d", s.actCounter)

	activity := &Activity{
		ID:          id,
		ClientID:    clientID,
		Type:        actType,
		Description: description,
		CreatedAt:   time.Now(),
	}

	s.activities[clientID] = append(s.activities[clientID], activity)
	return activity
}

func (s *Store) GetActivities(clientID string) []*Activity {
	s.mu.RLock()
	defer s.mu.RUnlock()

	activities := s.activities[clientID]
	return activities
}

func (s *Store) CreateInvoice(clientID, number, description string, amount float64, dueDate time.Time) *Invoice {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.clients[clientID]; !exists {
		return nil
	}

	s.invCounter++
	id := time.Now().Format("20060102150405") + "-inv-" + fmt.Sprintf("%04d", s.invCounter)

	invoice := &Invoice{
		ID:          id,
		ClientID:    clientID,
		Number:      number,
		Amount:      amount,
		Description: description,
		Status:      "draft",
		DueDate:     dueDate,
		CreatedAt:   time.Now(),
	}

	s.invoices[clientID] = append(s.invoices[clientID], invoice)
	return invoice
}

func (s *Store) ListInvoices(clientID string) []*Invoice {
	s.mu.RLock()
	defer s.mu.RUnlock()

	invoices := s.invoices[clientID]
	return invoices
}

func (s *Store) GetInvoice(id string) *Invoice {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, invoiceList := range s.invoices {
		for _, inv := range invoiceList {
			if inv.ID == id {
				return inv
			}
		}
	}
	return nil
}

func (s *Store) MarkInvoicePaid(id string) *Invoice {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, invoiceList := range s.invoices {
		for _, inv := range invoiceList {
			if inv.ID == id {
				inv.Status = "paid"
				return inv
			}
		}
	}
	return nil
}
