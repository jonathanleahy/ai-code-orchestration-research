package store

import (
	"sync"
	"time"
)

type Client struct {
	ID    string
	Name  string
	Email string
}

type Project struct {
	ID       string
	ClientID string
	Title    string
}

type Communication struct {
	ID        string
	ClientID  string
	Subject   string
	Content   string
	Timestamp string
}

type Invoice struct {
	ID        string
	ProjectID string
	Amount    float64
}

type Activity struct {
	ID        string
	ClientID  string
	Type      string
	Content   string
	Timestamp string
}

type Store struct {
	clients        map[string]*Client
	projects       map[string]*Project
	communications map[string]*Communication
	invoices       map[string]*Invoice
	activities     map[string]*Activity
	mu             sync.RWMutex
}

func NewStore() *Store {
	return &Store{
		clients:        make(map[string]*Client),
		projects:       make(map[string]*Project),
		communications: make(map[string]*Communication),
		invoices:       make(map[string]*Invoice),
		activities:     make(map[string]*Activity),
	}
}

func (s *Store) AddClient(c *Client) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.clients[c.ID] = c
}

func (s *Store) GetClient(id string) *Client {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.clients[id]
}

func (s *Store) AddProject(p *Project) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.projects[p.ID] = p
}

func (s *Store) GetProject(id string) *Project {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.projects[id]
}

func (s *Store) AddCommunication(c *Communication) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.communications[c.ID] = c
}

func (s *Store) GetCommunication(id string) *Communication {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.communications[id]
}

func (s *Store) AddInvoice(i *Invoice) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.invoices[i.ID] = i
}

func (s *Store) GetInvoice(id string) *Invoice {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.invoices[id]
}

func (s *Store) AddActivity(a *Activity) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.activities[a.ID] = a
}

func (s *Store) GetActivity(id string) *Activity {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.activities[id]
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

func (s *Store) ListProjects() []*Project {
	s.mu.RLock()
	defer s.mu.RUnlock()
	projects := make([]*Project, 0, len(s.projects))
	for _, p := range s.projects {
		projects = append(projects, p)
	}
	return projects
}

func (s *Store) ListCommunications() []*Communication {
	s.mu.RLock()
	defer s.mu.RUnlock()
	communications := make([]*Communication, 0, len(s.communications))
	for _, c := range s.communications {
		communications = append(communications, c)
	}
	return communications
}

func (s *Store) ListInvoices() []*Invoice {
	s.mu.RLock()
	defer s.mu.RUnlock()
	invoices := make([]*Invoice, 0, len(s.invoices))
	for _, i := range s.invoices {
		invoices = append(invoices, i)
	}
	return invoices
}

func (s *Store) ListActivities() []*Activity {
	s.mu.RLock()
	defer s.mu.RUnlock()
	activities := make([]*Activity, 0, len(s.activities))
	for _, a := range s.activities {
		activities = append(activities, a)
	}
	return activities
}

func (s *Store) Now() string {
	return time.Now().Format(time.RFC3339)
}
