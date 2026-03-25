package store

import (
	"fmt"
	"sync"
	"time"
)

type IncidentStatus string

const (
	StatusInvestigating IncidentStatus = "investigating"
	StatusIdentified    IncidentStatus = "identified"
	StatusMonitoring    IncidentStatus = "monitoring"
	StatusResolved      IncidentStatus = "resolved"
)

type Severity string

const (
	SeverityMinor    Severity = "minor"
	SeverityMajor    Severity = "major"
	SeverityCritical Severity = "critical"
)

type TimelineEntry struct {
	ID        string         `json:"id"`
	Message   string         `json:"message"`
	Status    IncidentStatus `json:"status"`
	CreatedAt string         `json:"created_at"`
}

type Incident struct {
	ID          string           `json:"id"`
	Title       string           `json:"title"`
	Description string           `json:"description"`
	Status      IncidentStatus   `json:"status"`
	Severity    Severity         `json:"severity"`
	CreatedAt   string           `json:"created_at"`
	UpdatedAt   string           `json:"updated_at"`
	ResolvedAt  string           `json:"resolved_at"`
	Timeline    []*TimelineEntry `json:"timeline"`
}

type Store struct {
	mu        sync.RWMutex
	incidents map[string]*Incident
	nextID    int
	nextTlID  int
}

func NewStore() *Store {
	return &Store{
		incidents: make(map[string]*Incident),
		nextID:    1,
		nextTlID:  1,
	}
}

func (s *Store) Create(title, description string, status IncidentStatus, severity Severity) (*Incident, error) {
	if title == "" {
		return nil, fmt.Errorf("title cannot be empty")
	}

	if status == "" {
		status = StatusInvestigating
	}
	if severity == "" {
		severity = SeverityMinor
	}

	now := time.Now().Format(time.RFC3339)
	id := fmt.Sprintf("inc-%d", s.nextID)
	s.nextID++

	incident := &Incident{
		ID:          id,
		Title:       title,
		Description: description,
		Status:      status,
		Severity:    severity,
		CreatedAt:   now,
		UpdatedAt:   now,
		Timeline: []*TimelineEntry{
			{
				ID:        fmt.Sprintf("tl-%d", s.nextTlID),
				Message:   "Incident created",
				Status:    status,
				CreatedAt: now,
			},
		},
	}

	s.mu.Lock()
	s.incidents[id] = incident
	s.mu.Unlock()

	return incident, nil
}

func (s *Store) Get(id string) (*Incident, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	incident, exists := s.incidents[id]
	if !exists {
		return nil, fmt.Errorf("incident not found")
	}

	return incident, nil
}

func (s *Store) List(status *IncidentStatus, openOnly bool) []*Incident {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []*Incident

	for _, incident := range s.incidents {
		if status != nil && incident.Status != *status {
			continue
		}
		if openOnly && incident.Status == StatusResolved {
			continue
		}
		result = append(result, incident)
	}

	return result
}

func (s *Store) Update(id string, status *IncidentStatus, description *string, message *string) (*Incident, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	incident, exists := s.incidents[id]
	if !exists {
		return nil, fmt.Errorf("incident not found")
	}

	now := time.Now().Format(time.RFC3339)
	if status != nil {
		incident.Status = *status
	}
	if description != nil {
		incident.Description = *description
	}
	if message != nil {
		tlID := fmt.Sprintf("tl-%d", s.nextTlID)
		s.nextTlID++
		incident.Timeline = append(incident.Timeline, &TimelineEntry{
			ID:        tlID,
			Message:   *message,
			Status:    incident.Status,
			CreatedAt: now,
		})
	}
	incident.UpdatedAt = now

	return incident, nil
}

func (s *Store) Resolve(id string, message string) (*Incident, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	incident, exists := s.incidents[id]
	if !exists {
		return nil, fmt.Errorf("incident not found")
	}

	if incident.Status == StatusResolved {
		return nil, fmt.Errorf("incident already resolved")
	}

	now := time.Now().Format(time.RFC3339)
	incident.Status = StatusResolved
	incident.ResolvedAt = now

	tlID := fmt.Sprintf("tl-%d", s.nextTlID)
	s.nextTlID++
	incident.Timeline = append(incident.Timeline, &TimelineEntry{
		ID:        tlID,
		Message:   message,
		Status:    incident.Status,
		CreatedAt: now,
	})

	return incident, nil
}

func (s *Store) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, exists := s.incidents[id]
	if !exists {
		return fmt.Errorf("incident not found")
	}

	delete(s.incidents, id)
	return nil
}
