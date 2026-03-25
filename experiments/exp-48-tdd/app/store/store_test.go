package store

import "testing"

func TestAddClient(t *testing.T) {
	store := NewStore()

	// Happy path
	client := store.AddClient("John Doe", "john@example.com", "123-456-7890")
	if client == nil {
		t.Fatal("AddClient should return a client")
	}
	if client.Name != "John Doe" {
		t.Errorf("Expected name 'John Doe', got '%s'", client.Name)
	}
	if client.Email != "john@example.com" {
		t.Errorf("Expected email 'john@example.com', got '%s'", client.Email)
	}
	if client.Phone != "123-456-7890" {
		t.Errorf("Expected phone '123-456-7890', got '%s'", client.Phone)
	}

	// Empty name
	client = store.AddClient("", "john@example.com", "123-456-7890")
	if client != nil {
		t.Fatal("AddClient should reject empty name")
	}
}

func TestGetClient(t *testing.T) {
	store := NewStore()

	// Not found
	client := store.GetClient(1)
	if client != nil {
		t.Fatal("GetClient should return nil for non-existent client")
	}

	// Found
	client = store.AddClient("John Doe", "john@example.com", "123-456-7890")
	if client == nil {
		t.Fatal("Failed to add client")
	}

	foundClient := store.GetClient(client.ID)
	if foundClient == nil {
		t.Fatal("GetClient should return the client")
	}
	if foundClient.ID != client.ID {
		t.Errorf("Expected client ID %d, got %d", client.ID, foundClient.ID)
	}
}

func TestListClients(t *testing.T) {
	store := NewStore()

	// Empty list
	clients := store.ListClients()
	if len(clients) != 0 {
		t.Fatal("ListClients should return empty slice for no clients")
	}

	// With clients
	store.AddClient("John Doe", "john@example.com", "123-456-7890")
	store.AddClient("Jane Smith", "jane@example.com", "098-765-4321")

	clients = store.ListClients()
	if len(clients) != 2 {
		t.Errorf("Expected 2 clients, got %d", len(clients))
	}
}

func TestUpdateClient(t *testing.T) {
	store := NewStore()

	// Not found
	client := store.UpdateClient(1, "John Doe", "john@example.com", "123-456-7890")
	if client != nil {
		t.Fatal("UpdateClient should return nil for non-existent client")
	}

	// Happy path
	client = store.AddClient("John Doe", "john@example.com", "123-456-7890")
	if client == nil {
		t.Fatal("Failed to add client")
	}

	updatedClient := store.UpdateClient(client.ID, "Jane Smith", "jane@example.com", "098-765-4321")
	if updatedClient == nil {
		t.Fatal("UpdateClient should return the updated client")
	}
	if updatedClient.Name != "Jane Smith" {
		t.Errorf("Expected name 'Jane Smith', got '%s'", updatedClient.Name)
	}
	if updatedClient.Email != "jane@example.com" {
		t.Errorf("Expected email 'jane@example.com', got '%s'", updatedClient.Email)
	}
	if updatedClient.Phone != "098-765-4321" {
		t.Errorf("Expected phone '098-765-4321', got '%s'", updatedClient.Phone)
	}
}

func TestDeleteClient(t *testing.T) {
	store := NewStore()

	// Not found
	result := store.DeleteClient(1)
	if result {
		t.Fatal("DeleteClient should return false for non-existent client")
	}

	// Happy path
	client := store.AddClient("John Doe", "john@example.com", "123-456-7890")
	if client == nil {
		t.Fatal("Failed to add client")
	}

	result = store.DeleteClient(client.ID)
	if !result {
		t.Fatal("DeleteClient should return true for existing client")
	}

	// Verify deletion
	deletedClient := store.GetClient(client.ID)
	if deletedClient != nil {
		t.Fatal("Client should be deleted")
	}
}

func TestAddActivity(t *testing.T) {
	store := NewStore()

	// Client not found
	activity := store.AddActivity(1, "call", "Called client")
	if activity != nil {
		t.Fatal("AddActivity should return nil for non-existent client")
	}

	// Happy path
	client := store.AddClient("John Doe", "john@example.com", "123-456-7890")
	if client == nil {
		t.Fatal("Failed to add client")
	}

	activity = store.AddActivity(client.ID, "call", "Called client")
	if activity == nil {
		t.Fatal("AddActivity should return the activity")
	}
	if activity.ClientID != client.ID {
		t.Errorf("Expected client ID %d, got %d", client.ID, activity.ClientID)
	}
	if activity.Type != "call" {
		t.Errorf("Expected activity type 'call', got '%s'", activity.Type)
	}
	if activity.Description != "Called client" {
		t.Errorf("Expected description 'Called client', got '%s'", activity.Description)
	}
}

func TestGetActivities(t *testing.T) {
	store := NewStore()

	// No activities
	activities := store.GetActivities(1)
	if len(activities) != 0 {
		t.Fatal("GetActivities should return empty slice for no activities")
	}

	// With activities
	client := store.AddClient("John Doe", "john@example.com", "123-456-7890")
	if client == nil {
		t.Fatal("Failed to add client")
	}

	store.AddActivity(client.ID, "call", "Called client")
	store.AddActivity(client.ID, "email", "Sent email")

	activities = store.GetActivities(client.ID)
	if len(activities) != 2 {
		t.Errorf("Expected 2 activities, got %d", len(activities))
	}
}

func TestCreateInvoice(t *testing.T) {
	store := NewStore()

	// Client not found
	invoice := store.CreateInvoice(1, 100.0, "Service", "2023-12-31")
	if invoice != nil {
		t.Fatal("CreateInvoice should return nil for non-existent client")
	}

	// Happy path
	client := store.AddClient("John Doe", "john@example.com", "123-456-7890")
	if client == nil {
		t.Fatal("Failed to add client")
	}

	invoice = store.CreateInvoice(client.ID, 100.0, "Service", "2023-12-31")
	if invoice == nil {
		t.Fatal("CreateInvoice should return the invoice")
	}
	if invoice.ClientID != client.ID {
		t.Errorf("Expected client ID %d, got %d", client.ID, invoice.ClientID)
	}
	if invoice.Amount != 100.0 {
		t.Errorf("Expected amount 100.0, got %f", invoice.Amount)
	}
	if invoice.Number != "INV-001" {
		t.Errorf("Expected invoice number 'INV-001', got '%s'", invoice.Number)
	}
	if invoice.Status != "draft" {
		t.Errorf("Expected status 'draft', got '%s'", invoice.Status)
	}
}

func TestListInvoices(t *testing.T) {
	store := NewStore()

	// No invoices
	invoices := store.ListInvoices(1)
	if len(invoices) != 0 {
		t.Fatal("ListInvoices should return empty slice for no invoices")
	}

	// With invoices
	client := store.AddClient("John Doe", "john@example.com", "123-456-7890")
	if client == nil {
		t.Fatal("Failed to add client")
	}

	store.CreateInvoice(client.ID, 100.0, "Service", "2023-12-31")
	store.CreateInvoice(client.ID, 200.0, "Another Service", "2023-12-31")

	invoices = store.ListInvoices(client.ID)
	if len(invoices) != 2 {
		t.Errorf("Expected 2 invoices, got %d", len(invoices))
	}
}

func TestMarkInvoicePaid(t *testing.T) {
	store := NewStore()

	// Not found
	invoice := store.MarkInvoicePaid(1)
	if invoice != nil {
		t.Fatal("MarkInvoicePaid should return nil for non-existent invoice")
	}

	// Happy path
	client := store.AddClient("John Doe", "john@example.com", "123-456-7890")
	if client == nil {
		t.Fatal("Failed to add client")
	}

	invoice = store.CreateInvoice(client.ID, 100.0, "Service", "2023-12-31")
	if invoice == nil {
		t.Fatal("Failed to create invoice")
	}

	paidInvoice := store.MarkInvoicePaid(invoice.ID)
	if paidInvoice == nil {
		t.Fatal("MarkInvoicePaid should return the paid invoice")
	}
	if paidInvoice.Status != "paid" {
		t.Errorf("Expected status 'paid', got '%s'", paidInvoice.Status)
	}
}
