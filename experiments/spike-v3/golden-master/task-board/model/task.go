package model

import (
	"fmt"
	"sync"
	"time"
)

type Status string

const (
	StatusTodo  Status = "TODO"
	StatusDoing Status = "DOING"
	StatusDone  Status = "DONE"
)

type Task struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      Status `json:"status"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

type Store struct {
	mu     sync.RWMutex
	tasks  map[string]*Task
	nextID int
}

func NewStore() *Store {
	return &Store{tasks: make(map[string]*Task)}
}

func (s *Store) Create(title, description string) (*Task, error) {
	if title == "" {
		return nil, fmt.Errorf("title is required")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.nextID++
	now := time.Now().UTC().Format(time.RFC3339)
	t := &Task{
		ID:          fmt.Sprintf("%d", s.nextID),
		Title:       title,
		Description: description,
		Status:      StatusTodo,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	s.tasks[t.ID] = t
	return t, nil
}

func (s *Store) Get(id string) (*Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	t, ok := s.tasks[id]
	if !ok {
		return nil, fmt.Errorf("task not found: %s", id)
	}
	return t, nil
}

func (s *Store) List(status *Status) []*Task {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var result []*Task
	for _, t := range s.tasks {
		if status == nil || t.Status == *status {
			result = append(result, t)
		}
	}
	return result
}

func (s *Store) Update(id string, title, description *string, status *Status) (*Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	t, ok := s.tasks[id]
	if !ok {
		return nil, fmt.Errorf("task not found: %s", id)
	}
	if title != nil {
		t.Title = *title
	}
	if description != nil {
		t.Description = *description
	}
	if status != nil {
		t.Status = *status
	}
	t.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	return t, nil
}

func (s *Store) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.tasks[id]; !ok {
		return fmt.Errorf("task not found: %s", id)
	}
	delete(s.tasks, id)
	return nil
}
