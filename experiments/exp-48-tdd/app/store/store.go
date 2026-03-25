package store

import (
	"sync"
	"time"
)

type Client struct {
	ID        int
	Name      string
	Email     string
	Phone     string
	CreatedAt time.Time
}

type Activity struct {
	ID          int
	ClientID    int
	Type        string
	Description string
	CreatedAt   time.Time
}

type Invoice struct {
	ID        int
	ClientID  int
	Number    string
	Amount    float64
	Status    string
	DueDate   time.Time
	CreatedAt time.Time
}

type Store struct {
	clients    map[int]*Client
	activities map[int]*Activity
	invoices   map[int]*Invoice
	clientID   int
	activityID int
	invoiceID  int
	mutex      sync.RWMutex
}

func NewStore() *Store {
	return &Store{
		clients:    make(map[int]*Client),
		activities: make(map[int]*Activity),
		invoices:   make(map[int]*Invoice),
		clientID:   1,
		activityID: 1,
		invoiceID:  1,
	}
}

func (s *Store) AddClient(name, email, phone string) *Client {
	if name == "" {
		return nil
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	client := &Client{
		ID:        s.clientID,
		Name:      name,
		Email:     email,
		Phone:     phone,
		CreatedAt: time.Now(),
	}

	s.clients[s.clientID] = client
	s.clientID++

	return client
}

func (s *Store) GetClient(id int) *Client {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	client, exists := s.clients[id]
	if !exists {
		return nil
	}

	return client
}

func (s *Store) ListClients() []*Client {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	clients := make([]*Client, 0, len(s.clients))
	for _, client := range s.clients {
		clients = append(clients, client)
	}

	return clients
}

func (s *Store) UpdateClient(id int, name, email, phone string) *Client {
	if name == "" {
		return nil
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	client, exists := s.clients[id]
	if !exists {
		return nil
	}

	client.Name = name
	client.Email = email
	client.Phone = phone

	return client
}

func (s *Store) DeleteClient(id int) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	_, exists := s.clients[id]
	if !exists {
		return false
	}

	delete(s.clients, id)
	return true
}

func (s *Store) AddActivity(clientID int, actType, desc string) *Activity {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	_, exists := s.clients[clientID]
	if !exists {
		return nil
	}

	activity := &Activity{
		ID:          s.activityID,
		ClientID:    clientID,
		Type:        actType,
		Description: desc,
		CreatedAt:   time.Now(),
	}

	s.activities[s.activityID] = activity
	s.activityID++

	return activity
}

func (s *Store) GetActivities(clientID int) []*Activity {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	activities := make([]*Activity, 0)
	for _, activity := range s.activities {
		if activity.ClientID == clientID {
			activities = append(activities, activity)
		}
	}

	return activities
}

func (s *Store) CreateInvoice(clientID int, number string, amount float64, status string, dueDate string) *Invoice {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	_, exists := s.clients[clientID]
	if !exists {
		return nil
	}

	// Parse due date
	dueTime, err := time.Parse("2006-01-02", dueDate)
	if err != nil {
		return nil
	}

	invoice := &Invoice{
		ID:        s.invoiceID,
		ClientID:  clientID,
		Number:    number,
		Amount:    amount,
		Status:    status,
		DueDate:   dueTime,
		CreatedAt: time.Now(),
	}

	s.invoices[s.invoiceID] = invoice
	s.invoiceID++

	return invoice
}

func (s *Store) ListInvoices() []*Invoice {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	invoices := make([]*Invoice, 0, len(s.invoices))
	for _, invoice := range s.invoices {
		invoices = append(invoices, invoice)
	}

	return invoices
}
