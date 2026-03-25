package store

import (
	"fmt"
	"sync"
	"time"
)

type CheckStatus string

const (
	StatusUp      CheckStatus = "up"
	StatusDown    CheckStatus = "down"
	StatusUnknown CheckStatus = "unknown"
)

type Check struct {
	ID              string      `json:"id"`
	Name            string      `json:"name"`
	URL             string      `json:"url"`
	IntervalSeconds int         `json:"interval_seconds"`
	Status          CheckStatus `json:"status"`
	LatencyMs       int         `json:"latency_ms"`
	LastCheckedAt   string      `json:"last_checked_at"`
	CreatedAt       string      `json:"created_at"`
}

type Result struct {
	Status     CheckStatus `json:"status"`
	LatencyMs  int         `json:"latency_ms"`
	StatusCode int         `json:"status_code"`
	CheckedAt  string      `json:"checked_at"`
}

type Store struct {
	mu      sync.RWMutex
	checks  map[string]*Check
	results map[string][]*Result
	nextID  int
}

func NewStore() *Store {
	return &Store{
		checks:  make(map[string]*Check),
		results: make(map[string][]*Result),
		nextID:  1,
	}
}

func (s *Store) CreateCheck(name, url string, intervalSeconds int) (*Check, error) {
	if name == "" {
		return nil, fmt.Errorf("name cannot be empty")
	}
	if url == "" {
		return nil, fmt.Errorf("url cannot be empty")
	}

	if intervalSeconds < 10 {
		intervalSeconds = 60
	}

	check := &Check{
		ID:              fmt.Sprintf("check-%d", s.nextID),
		Name:            name,
		URL:             url,
		IntervalSeconds: intervalSeconds,
		Status:          StatusUnknown,
		LatencyMs:       0,
		LastCheckedAt:   "",
		CreatedAt:       time.Now().Format(time.RFC3339),
	}

	s.mu.Lock()
	s.checks[check.ID] = check
	s.nextID++
	s.mu.Unlock()

	return check, nil
}

func (s *Store) GetCheck(id string) (*Check, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	check, exists := s.checks[id]
	if !exists {
		return nil, fmt.Errorf("check not found")
	}

	return check, nil
}

func (s *Store) ListChecks() []*Check {
	s.mu.RLock()
	defer s.mu.RUnlock()

	checks := make([]*Check, 0, len(s.checks))
	for _, check := range s.checks {
		checks = append(checks, check)
	}

	return checks
}

func (s *Store) DeleteCheck(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, exists := s.checks[id]
	if !exists {
		return fmt.Errorf("check not found")
	}

	delete(s.checks, id)
	delete(s.results, id)

	return nil
}

func (s *Store) RecordResult(checkID string, result *Result) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	check, exists := s.checks[checkID]
	if !exists {
		return fmt.Errorf("check not found")
	}

	check.Status = result.Status
	check.LatencyMs = result.LatencyMs
	check.LastCheckedAt = result.CheckedAt

	results := s.results[checkID]
	results = append(results, result)

	if len(results) > 100 {
		results = results[len(results)-100:]
	}

	s.results[checkID] = results

	return nil
}

func (s *Store) GetResults(checkID string, limit int) []*Result {
	s.mu.RLock()
	defer s.mu.RUnlock()

	results, exists := s.results[checkID]
	if !exists {
		return []*Result{}
	}

	if limit <= 0 {
		return []*Result{}
	}

	if limit > len(results) {
		limit = len(results)
	}

	start := len(results) - limit
	return results[start:]
}
