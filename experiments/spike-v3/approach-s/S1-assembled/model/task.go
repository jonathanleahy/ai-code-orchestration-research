package model

import (
	"fmt"
	"sync"
	"time"
)

// Status represents the current state of a task
type Status string

const (
	StatusTodo  Status = "TODO"
	StatusDoing Status = "DOING"
	StatusDone  Status = "DONE"
)

// Task represents a task in the task board
type Task struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      Status `json:"status"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

// Store provides thread-safe in-memory storage for tasks
type Store struct {
	mu     sync.RWMutex
	tasks  map[string]*Task
	nextID int
}

// NewStore creates a new Store instance
func NewStore() *Store {
	return &Store{
		tasks:  make(map[string]*Task),
		nextID: 1,
	}
}

// Create adds a new task to the store
func (s *Store) Create(title, description string) (*Task, error) {
	if title == "" {
		return nil, fmt.Errorf("title cannot be empty")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().UTC().Format(time.RFC3339)
	task := &Task{
		ID:          fmt.Sprintf("%d", s.nextID),
		Title:       title,
		Description: description,
		Status:      StatusTodo,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	s.tasks[task.ID] = task
	s.nextID++

	return task, nil
}

// Get retrieves a task by ID
func (s *Store) Get(id string) (*Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	task, exists := s.tasks[id]
	if !exists {
		return nil, fmt.Errorf("task not found")
	}

	return task, nil
}

// List returns all tasks or filters by status if provided
func (s *Store) List(status *Status) []*Task {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []*Task
	for _, task := range s.tasks {
		if status == nil || task.Status == *status {
			result = append(result, task)
		}
	}

	return result
}

// Update modifies an existing task's fields
func (s *Store) Update(id string, title, description *string, status *Status) (*Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	task, exists := s.tasks[id]
	if !exists {
		return nil, fmt.Errorf("task not found")
	}

	if title != nil {
		task.Title = *title
	}
	if description != nil {
		task.Description = *description
	}
	if status != nil {
		task.Status = *status
	}

	task.UpdatedAt = time.Now().UTC().Format(time.RFC3339)

	return task, nil
}

// Delete removes a task from the store
func (s *Store) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.tasks[id]; !exists {
		return fmt.Errorf("task not found")
	}

	delete(s.tasks, id)
	return nil
}