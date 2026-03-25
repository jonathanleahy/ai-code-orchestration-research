package store

import (
	"reflect"
	"testing"
)

func TestAddSubscriber(t *testing.T) {
	store := NewStore()
	subscriber, err := store.AddSubscriber("Test", "http://example.com/webhook", []string{"event1", "event2"})

	if err != nil {
		t.Fatalf("AddSubscriber failed: %v", err)
	}

	if subscriber.Name != "Test" {
		t.Errorf("Expected name 'Test', got '%s'", subscriber.Name)
	}

	if subscriber.WebhookURL != "http://example.com/webhook" {
		t.Errorf("Expected webhook URL 'http://example.com/webhook', got '%s'", subscriber.WebhookURL)
	}

	if !reflect.DeepEqual(subscriber.Events, []string{"event1", "event2"}) {
		t.Errorf("Expected events ['event1', 'event2'], got %v", subscriber.Events)
	}

	if !subscriber.Active {
		t.Error("Expected subscriber to be active")
	}

	if subscriber.ID == "" {
		t.Error("Expected subscriber ID to be set")
	}
}

func TestAddSubscriberDefaultEvents(t *testing.T) {
	store := NewStore()
	subscriber, err := store.AddSubscriber("Test", "http://example.com/webhook", nil)

	if err != nil {
		t.Fatalf("AddSubscriber failed: %v", err)
	}

	if !reflect.DeepEqual(subscriber.Events, []string{"all"}) {
		t.Errorf("Expected default events ['all'], got %v", subscriber.Events)
	}
}

func TestAddSubscriberEmptyName(t *testing.T) {
	store := NewStore()
	_, err := store.AddSubscriber("", "http://example.com/webhook", []string{"event1"})

	if err == nil {
		t.Error("Expected error for empty name")
	}
}

func TestAddSubscriberEmptyURL(t *testing.T) {
	store := NewStore()
	_, err := store.AddSubscriber("Test", "", []string{"event1"})

	if err == nil {
		t.Error("Expected error for empty URL")
	}
}

func TestListSubscribers(t *testing.T) {
	store := NewStore()

	subscribers := store.ListSubscribers()
	if len(subscribers) != 0 {
		t.Errorf("Expected 0 subscribers, got %d", len(subscribers))
	}

	store.AddSubscriber("Test1", "http://example.com/webhook1", []string{"event1"})
	store.AddSubscriber("Test2", "http://example.com/webhook2", []string{"event2"})

	subscribers = store.ListSubscribers()
	if len(subscribers) != 2 {
		t.Errorf("Expected 2 subscribers, got %d", len(subscribers))
	}
}

func TestDeleteSubscriber(t *testing.T) {
	store := NewStore()
	subscriber, _ := store.AddSubscriber("Test", "http://example.com/webhook", []string{"event1"})

	err := store.DeleteSubscriber(subscriber.ID)
	if err != nil {
		t.Fatalf("DeleteSubscriber failed: %v", err)
	}

	_, err = store.GetSubscriber(subscriber.ID)
	if err == nil {
		t.Error("Expected error when getting deleted subscriber")
	}
}

func TestGetMatchingAll(t *testing.T) {
	store := NewStore()
	store.AddSubscriber("Test1", "http://example.com/webhook1", []string{"all"})
	store.AddSubscriber("Test2", "http://example.com/webhook2", []string{"event1"})
	store.AddSubscriber("Test3", "http://example.com/webhook3", []string{"event2"})

	matching := store.GetMatchingSubscribers("event1")
	if len(matching) != 1 {
		t.Errorf("Expected 1 matching subscriber, got %d", len(matching))
	}

	if matching[0].Name != "Test1" {
		t.Errorf("Expected Test1 to match")
	}
}

func TestGetMatchingSpecific(t *testing.T) {
	store := NewStore()
	store.AddSubscriber("Test1", "http://example.com/webhook1", []string{"all"})
	store.AddSubscriber("Test2", "http://example.com/webhook2", []string{"event1"})
	store.AddSubscriber("Test3", "http://example.com/webhook3", []string{"event2"})

	matching := store.GetMatchingSubscribers("event2")
	if len(matching) != 2 {
		t.Errorf("Expected 2 matching subscribers, got %d", len(matching))
	}

	names := make(map[string]bool)
	for _, sub := range matching {
		names[sub.Name] = true
	}

	if !names["Test1"] || !names["Test3"] {
		t.Error("Expected both Test1 and Test3 to match")
	}
}
