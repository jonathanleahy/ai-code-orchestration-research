package store

import (
	"testing"
)

func TestListClients(t *testing.T) {
	s := NewStore()
	clients := s.ListClients()
	if len(clients) != 0 {
		t.Fatalf("expected 0 clients, got %d", len(clients))
	}
}
