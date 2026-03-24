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
	mu      sync.RWMutex
	tasks   map[string]*Task
	nextID  int
}

func NewStore() *Store {
	return &Store{
		tasks:  make(map[string]*Task),
		nextID: 1,
	}
}

func (s *Store) Create(title, description string) (*Task, error) {
	if title == "" {
		return nil, fmt.Errorf("title cannot be empty")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	id := fmt.Sprintf("%d", s.nextID)
	s.nextID++

	now := time.Now().UTC().Format(time.RFC3339)
	task := &Task{
		ID:          id,
		Title:       title,
		Description: description,
		Status:      StatusTodo,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	s.tasks[id] = task
	return task, nil
}

func (s *Store) Get(id string) (*Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	task, exists := s.tasks[id]
	if !exists {
		return nil, fmt.Errorf("task with id %s not found", id)
	}

	return task, nil
}

func (s *Store) List(status *Status) []*Task {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []*Task

	for _, task := range s.tasks {
		if status == nil || *status == task.Status {
			result = append(result, task)
		}
	}

	return result
}

func (s *Store) Update(id string, title, description *string, status *Status) (*Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	task, exists := s.tasks[id]
	if !exists {
		return nil, fmt.Errorf("task with id %s not found", id)
	}

	now := time.Now().UTC().Format(time.RFC3339)
	if title != nil {
		if *title == "" {
			return nil, fmt.Errorf("title cannot be empty")
		}
		task.Title = *title
	}
	if description != nil {
		task.Description = *description
	}
	if status != nil {
		task.Status = *status
	}
	task.UpdatedAt = now

	return task, nil
}

func (s *Store) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.tasks[id]; !exists {
		return fmt.Errorf("task with id %s not found", id)
	}

	delete(s.tasks, id)
	return nil
}