package store

import (
	"testing"
	"time"
)

func TestClientCRUD(t *testing.T) {
	store := NewInMemoryStore()

	// Create client
	client := &Client{
		ID:      "client1",
		Name:    "John Doe",
		Email:   "john@example.com",
		Company: "Acme Corp",
	}
	err := store.CreateClient(client)
	if err != nil {
		t.Fatalf("CreateClient failed: %v", err)
	}

	// Get client
	retrieved, err := store.GetClient("client1")
	if err != nil {
		t.Fatalf("GetClient failed: %v", err)
	}
	if retrieved == nil {
		t.Fatal("GetClient returned nil")
	}
	if retrieved.Name != "John Doe" {
		t.Errorf("Expected name 'John Doe', got '%s'", retrieved.Name)
	}

	// Update client
	client.Name = "Jane Doe"
	err = store.UpdateClient("client1", client)
	if err != nil {
		t.Fatalf("UpdateClient failed: %v", err)
	}
	retrieved, err = store.GetClient("client1")
	if err != nil {
		t.Fatalf("GetClient after update failed: %v", err)
	}
	if retrieved.Name != "Jane Doe" {
		t.Errorf("Expected name 'Jane Doe', got '%s'", retrieved.Name)
	}

	// Delete client
	err = store.DeleteClient("client1")
	if err != nil {
		t.Fatalf("DeleteClient failed: %v", err)
	}
	retrieved, err = store.GetClient("client1")
	if err != nil {
		t.Fatalf("GetClient after delete failed: %v", err)
	}
	if retrieved != nil {
		t.Error("Client should be deleted")
	}
}

func TestActivityCRUD(t *testing.T) {
	store := NewInMemoryStore()

	// Create activity
	activity := &Activity{
		ID:       "activity1",
		ClientID: "client1",
		Type:     "meeting",
		Notes:    "Team meeting",
	}
	err := store.CreateActivity(activity)
	if err != nil {
		t.Fatalf("CreateActivity failed: %v", err)
	}

	// Get activity
	retrieved, err := store.GetActivity("activity1")
	if err != nil {
		t.Fatalf("GetActivity failed: %v", err)
	}
	if retrieved == nil {
		t.Fatal("GetActivity returned nil")
	}
	if retrieved.Type != "meeting" {
		t.Errorf("Expected type 'meeting', got '%s'", retrieved.Type)
	}

	// Update activity
	activity.Notes = "Updated meeting notes"
	err = store.UpdateActivity("activity1", activity)
	if err != nil {
		t.Fatalf("UpdateActivity failed: %v", err)
	}
	retrieved, err = store.GetActivity("activity1")
	if err != nil {
		t.Fatalf("GetActivity after update failed: %v", err)
	}
	if retrieved.Notes != "Updated meeting notes" {
		t.Errorf("Expected notes 'Updated meeting notes', got '%s'", retrieved.Notes)
	}

	// Delete activity
	err = store.DeleteActivity("activity1")
	if err != nil {
		t.Fatalf("DeleteActivity failed: %v", err)
	}
	retrieved, err = store.GetActivity("activity1")
	if err != nil {
		t.Fatalf("GetActivity after delete failed: %v", err)
	}
	if retrieved != nil {
		t.Error("Activity should be deleted")
	}
}

func TestInvoiceCRUD(t *testing.T) {
	store := NewInMemoryStore()

	// Create invoice
	invoice := &Invoice{
		ID: "invoice1",
		Items: []InvoiceItem{
			{
				Description: "Service 1",
				Quantity:    5,
				Rate:        100.0,
			},
		},
		ClientID: "client1",
	}
	err := store.CreateInvoice(invoice)
	if err != nil {
		t.Fatalf("CreateInvoice failed: %v", err)
	}

	// Verify items were calculated
	if len(invoice.Items) != 1 {
		t.Fatalf("Expected 1 item, got %d", len(invoice.Items))
	}
	if invoice.Items[0].Amount != 500.0 {
		t.Errorf("Expected amount 500.0, got %f", invoice.Items[0].Amount)
	}

	// Get invoice
	retrieved, err := store.GetInvoice("invoice1")
	if err != nil {
		t.Fatalf("GetInvoice failed: %v", err)
	}
	if retrieved == nil {
		t.Fatal("GetInvoice returned nil")
	}
	if len(retrieved.Items) != 1 {
		t.Errorf("Expected 1 item, got %d", len(retrieved.Items))
	}
	if retrieved.Items[0].Amount != 500.0 {
		t.Errorf("Expected amount 500.0, got %f", retrieved.Items[0].Amount)
	}

	// Update invoice
	invoice.Items[0].Description = "Updated Service"
	err = store.UpdateInvoice("invoice1", invoice)
	if err != nil {
		t.Fatalf("UpdateInvoice failed: %v", err)
	}
	retrieved, err = store.GetInvoice("invoice1")
	if err != nil {
		t.Fatalf("GetInvoice after update failed: %v", err)
	}
	if retrieved.Items[0].Description != "Updated Service" {
		t.Errorf("Expected description 'Updated Service', got '%s'", retrieved.Items[0].Description)
	}

	// Delete invoice
	err = store.DeleteInvoice("invoice1")
	if err != nil {
		t.Fatalf("DeleteInvoice failed: %v", err)
	}
	retrieved, err = store.GetInvoice("invoice1")
	if err != nil {
		t.Fatalf("GetInvoice after delete failed: %v", err)
	}
	if retrieved != nil {
		t.Error("Invoice should be deleted")
	}
}

func TestInvoiceStatusUpdates(t *testing.T) {
	store := NewInMemoryStore()

	// Create invoice
	invoice := &Invoice{
		ID: "invoice1",
		Items: []InvoiceItem{
			{
				Description: "Service 1",
				Quantity:    5,
				Rate:        100.0,
			},
		},
		ClientID: "client1",
		Status:   "draft",
	}
	err := store.CreateInvoice(invoice)
	if err != nil {
		t.Fatalf("CreateInvoice failed: %v", err)
	}

	// Mark as paid
	err = store.MarkInvoicePaid("invoice1")
	if err != nil {
		t.Fatalf("MarkInvoicePaid failed: %v", err)
	}
	retrieved, err := store.GetInvoice("invoice1")
	if err != nil {
		t.Fatalf("GetInvoice after mark paid failed: %v", err)
	}
	if retrieved.Status != "paid" {
		t.Errorf("Expected status 'paid', got '%s'", retrieved.Status)
	}

	// Void invoice
	err = store.VoidInvoice("invoice1")
	if err != nil {
		t.Fatalf("VoidInvoice failed: %v", err)
	}
	retrieved, err = store.GetInvoice("invoice1")
	if err != nil {
		t.Fatalf("GetInvoice after void failed: %v", err)
	}
	if retrieved.Status != "void" {
		t.Errorf("Expected status 'void', got '%s'", retrieved.Status)
	}
}

func TestSearchClients(t *testing.T) {
	store := NewInMemoryStore()

	// Create clients
	client1 := &Client{
		ID:      "client1",
		Name:    "John Doe",
		Email:   "john@example.com",
		Company: "Acme Corp",
	}
	client2 := &Client{
		ID:      "client2",
		Name:    "Jane Smith",
		Email:   "jane@example.com",
		Company: "Beta Inc",
	}
	client3 := &Client{
		ID:      "client3",
		Name:    "Bob Johnson",
		Email:   "bob@example.com",
		Company: "Acme Corp",
	}

	store.CreateClient(client1)
	store.CreateClient(client2)
	store.CreateClient(client3)

	// Search by name
	results, err := store.SearchClients("John Doe")
	if err != nil {
		t.Fatalf("SearchClients failed: %v", err)
	}
	if len(results) != 1 {
		t.Errorf("Expected 1 result for name search, got %d", len(results))
	}
	if results[0].ID != "client1" {
		t.Errorf("Expected client1, got %s", results[0].ID)
	}

	// Search by email
	results, err = store.SearchClients("jane@example.com")
	if err != nil {
		t.Fatalf("SearchClients failed: %v", err)
	}
	if len(results) != 1 {
		t.Errorf("Expected 1 result for email search, got %d", len(results))
	}
	if results[0].ID != "client2" {
		t.Errorf("Expected client2, got %s", results[0].ID)
	}

	// Search by company
	results, err = store.SearchClients("Acme Corp")
	if err != nil {
		t.Fatalf("SearchClients failed: %v", err)
	}
	if len(results) != 2 {
		t.Errorf("Expected 2 results for company search, got %d", len(results))
	}

	// Search for non-existent client
	results, err = store.SearchClients("Non Existent")
	if err != nil {
		t.Fatalf("SearchClients failed: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("Expected 0 results for non-existent search, got %d", len(results))
	}
}

func TestGetClientActivities(t *testing.T) {
	store := NewInMemoryStore()

	// Create client
	client := &Client{
		ID:      "client1",
		Name:    "John Doe",
		Email:   "john@example.com",
		Company: "Acme Corp",
	}
	store.CreateClient(client)

	// Create activities
	activity1 := &Activity{
		ID:       "activity1",
		ClientID: "client1",
		Type:     "meeting",
		Notes:    "Team meeting",
	}
	activity2 := &Activity{
		ID:       "activity2",
		ClientID: "client1",
		Type:     "call",
		Notes:    "Follow up call",
	}
	activity3 := &Activity{
		ID:       "activity3",
		ClientID: "client2",
		Type:     "meeting",
		Notes:    "Another meeting",
	}

	store.CreateActivity(activity1)
	store.CreateActivity(activity2)
	store.CreateActivity(activity3)

	// Get client activities
	activities, err := store.GetClientActivities("client1")
	if err != nil {
		t.Fatalf("GetClientActivities failed: %v", err)
	}
	if len(activities) != 2 {
		t.Errorf("Expected 2 activities for client1, got %d", len(activities))
	}

	// Verify activities belong to correct client
	for _, activity := range activities {
		if activity.ClientID != "client1" {
			t.Errorf("Activity %s belongs to wrong client: %s", activity.ID, activity.ClientID)
		}
	}
}

func TestGetClientInvoices(t *testing.T) {
	store := NewInMemoryStore()

	// Create client
	client := &Client{
		ID:      "client1",
		Name:    "John Doe",
		Email:   "john@example.com",
		Company: "Acme Corp",
	}
	store.CreateClient(client)

	// Create invoices
	invoice1 := &Invoice{
		ID:       "invoice1",
		ClientID: "client1",
		Status:   "draft",
	}
	invoice2 := &Invoice{
		ID:       "invoice2",
		ClientID: "client1",
		Status:   "paid",
	}
	invoice3 := &Invoice{
		ID:       "invoice3",
		ClientID: "client2",
		Status:   "draft",
	}

	store.CreateInvoice(invoice1)
	store.CreateInvoice(invoice2)
	store.CreateInvoice(invoice3)

	// Get client invoices
	invoices, err := store.GetClientInvoices("client1")
	if err != nil {
		t.Fatalf("GetClientInvoices failed: %v", err)
	}
	if len(invoices) != 2 {
		t.Errorf("Expected 2 invoices for client1, got %d", len(invoices))
	}

	// Verify invoices belong to correct client
	for _, invoice := range invoices {
		if invoice.ClientID != "client1" {
			t.Errorf("Invoice %s belongs to wrong client: %s", invoice.ID, invoice.ClientID)
		}
	}
}

func TestConcurrentAccess(t *testing.T) {
	store := NewInMemoryStore()

	// Test concurrent client creation
	go func() {
		for i := 0; i < 100; i++ {
			client := &Client{
				ID:    "client" + string(rune(i)),
				Name:  "Client " + string(rune(i)),
				Email: "client" + string(rune(i)) + "@example.com",
			}
			store.CreateClient(client)
		}
	}()

	// Test concurrent client retrieval
	go func() {
		for i := 0; i < 100; i++ {
			store.GetClient("client" + string(rune(i)))
		}
	}()

	// Give goroutines time to complete
	time.Sleep(10 * time.Millisecond)

	// Verify all clients were created
	for i := 0; i < 100; i++ {
		client, err := store.GetClient("client" + string(rune(i)))
		if err != nil {
			t.Fatalf("Failed to retrieve client %d: %v", i, err)
		}
		if client == nil {
			t.Errorf("Client %d was not created", i)
		}
	}
}
