package store

import (
	"fmt"
	"sync"
	"time"
)

type Subscriber struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	WebhookURL string   `json:"webhook_url"`
	Events     []string `json:"events"`
	Active     bool     `json:"active"`
	CreatedAt  string   `json:"created_at"`
}

type DispatchResult struct {
	SubscriberID string `json:"subscriber_id"`
	Success      bool   `json:"success"`
	StatusCode   int    `json:"status_code"`
	Error        string `json:"error,omitempty"`
}

type Store struct {
	mu          sync.RWMutex
	subscribers map[string]*Subscriber
	nextID      int
}

func NewStore() *Store {
	return &Store{
		subscribers: make(map[string]*Subscriber),
		nextID:      1,
	}
}

func (s *Store) AddSubscriber(name, webhookURL string, events []string) (*Subscriber, error) {
	if name == "" {
		return nil, fmt.Errorf("name cannot be empty")
	}
	if webhookURL == "" {
		return nil, fmt.Errorf("webhook URL cannot be empty")
	}

	if events == nil || len(events) == 0 {
		events = []string{"all"}
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	id := fmt.Sprintf("sub-%d", s.nextID)
	s.nextID++

	subscriber := &Subscriber{
		ID:         id,
		Name:       name,
		WebhookURL: webhookURL,
		Events:     events,
		Active:     true,
		CreatedAt:  time.Now().Format(time.RFC3339),
	}

	s.subscribers[id] = subscriber
	return subscriber, nil
}

func (s *Store) ListSubscribers() []*Subscriber {
	s.mu.RLock()
	defer s.mu.RUnlock()

	subscribers := make([]*Subscriber, 0, len(s.subscribers))
	for _, subscriber := range s.subscribers {
		subscribers = append(subscribers, subscriber)
	}
	return subscribers
}

func (s *Store) GetSubscriber(id string) (*Subscriber, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	subscriber, exists := s.subscribers[id]
	if !exists {
		return nil, fmt.Errorf("subscriber with id %s not found", id)
	}
	return subscriber, nil
}

func (s *Store) DeleteSubscriber(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.subscribers[id]; !exists {
		return fmt.Errorf("subscriber with id %s not found", id)
	}
	delete(s.subscribers, id)
	return nil
}

func (s *Store) GetMatchingSubscribers(event string) []*Subscriber {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var matching []*Subscriber
	for _, subscriber := range s.subscribers {
		if !subscriber.Active {
			continue
		}
		for _, e := range subscriber.Events {
			if e == "all" || e == event {
				matching = append(matching, subscriber)
				break
			}
		}
	}
	return matching
}
