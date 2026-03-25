package store

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestInvoiceStore(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "store_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	store, err := NewStore(tempDir)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}

	// Test CreateInvoice (valid)
	t.Run("CreateInvoice valid", func(t *testing.T) {
		invoice := &Invoice{
			ClientID:  "client1",
			Number:    "INV-001",
			IssueDate: time.Now(),
			DueDate:   time.Now().Add(30 * 24 * time.Hour),
			Status:    StatusDraft,
			LineItems: []LineItem{
				{Description: "Service 1", Quantity: 1, UnitPrice: 100.0},
				{Description: "Service 2", Quantity: 2, UnitPrice: 50.0},
			},
			TaxRate:   0.1,
			Notes:     "Test invoice",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		err := store.CreateInvoice(invoice)
		if err != nil {
			t.Fatalf("CreateInvoice failed: %v", err)
		}

		if invoice.ID == "" {
			t.Error("Invoice ID should be generated")
		}

		if invoice.Subtotal != 200.0 {
			t.Errorf("Expected subtotal 200.0, got %f", invoice.Subtotal)
		}

		if invoice.TaxAmount != 20.0 {
			t.Errorf("Expected tax amount 20.0, got %f", invoice.TaxAmount)
		}

		if invoice.Total != 220.0 {
			t.Errorf("Expected total 220.0, got %f", invoice.Total)
		}
	})

	// Test CreateInvoice (missing required fields)
	t.Run("CreateInvoice missing required fields", func(t *testing.T) {
		invoice := &Invoice{
			ClientID:  "",
			Number:    "",
			IssueDate: time.Time{},
			DueDate:   time.Time{},
			Status:    "",
			LineItems: []LineItem{},
			TaxRate:   0.0,
		}

		err := store.CreateInvoice(invoice)
		if err == nil {
			t.Error("CreateInvoice should fail with missing required fields")
		}
	})

	// Test GetInvoice by ID
	t.Run("GetInvoice by ID", func(t *testing.T) {
		invoice := &Invoice{
			ClientID:  "client1",
			Number:    "INV-002",
			IssueDate: time.Now(),
			DueDate:   time.Now().Add(30 * 24 * time.Hour),
			Status:    StatusDraft,
			LineItems: []LineItem{
				{Description: "Service 1", Quantity: 1, UnitPrice: 100.0},
			},
			TaxRate:   0.1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		err := store.CreateInvoice(invoice)
		if err != nil {
			t.Fatalf("CreateInvoice failed: %v", err)
		}

		retrieved, err := store.GetInvoice(invoice.ID)
		if err != nil {
			t.Fatalf("GetInvoice failed: %v", err)
		}

		if retrieved.ID != invoice.ID {
			t.Errorf("Expected ID %s, got %s", invoice.ID, retrieved.ID)
		}

		if retrieved.Number != invoice.Number {
			t.Errorf("Expected number %s, got %s", invoice.Number, retrieved.Number)
		}
	})

	// Test GetInvoice not found
	t.Run("GetInvoice not found", func(t *testing.T) {
		_, err := store.GetInvoice("nonexistent")
		if err == nil {
			t.Error("GetInvoice should return error for non-existent invoice")
		}
	})

	// Test ListInvoices
	t.Run("ListInvoices", func(t *testing.T) {
		invoices, err := store.ListInvoices()
		if err != nil {
			t.Fatalf("ListInvoices failed: %v", err)
		}

		if len(invoices) != 2 {
			t.Errorf("Expected 2 invoices, got %d", len(invoices))
		}
	})

	// Test UpdateInvoice
	t.Run("UpdateInvoice", func(t *testing.T) {
		invoice := &Invoice{
			ClientID:  "client1",
			Number:    "INV-003",
			IssueDate: time.Now(),
			DueDate:   time.Now().Add(30 * 24 * time.Hour),
			Status:    StatusDraft,
			LineItems: []LineItem{
				{Description: "Service 1", Quantity: 1, UnitPrice: 100.0},
			},
			TaxRate:   0.1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		err := store.CreateInvoice(invoice)
		if err != nil {
			t.Fatalf("CreateInvoice failed: %v", err)
		}

		// Update the invoice
		updatedStatus := StatusSent
		invoice.Status = updatedStatus
		invoice.LineItems = append(invoice.LineItems, LineItem{
			Description: "Service 2",
			Quantity:    2,
			UnitPrice:   50.0,
		})

		err = store.UpdateInvoice(invoice)
		if err != nil {
			t.Fatalf("UpdateInvoice failed: %v", err)
		}

		retrieved, err := store.GetInvoice(invoice.ID)
		if err != nil {
			t.Fatalf("GetInvoice failed: %v", err)
		}

		if retrieved.Status != updatedStatus {
			t.Errorf("Expected status %s, got %s", updatedStatus, retrieved.Status)
		}

		if retrieved.Subtotal != 200.0 {
			t.Errorf("Expected subtotal 200.0, got %f", retrieved.Subtotal)
		}
	})

	// Test DeleteInvoice
	t.Run("DeleteInvoice", func(t *testing.T) {
		invoice := &Invoice{
			ClientID:  "client1",
			Number:    "INV-004",
			IssueDate: time.Now(),
			DueDate:   time.Now().Add(30 * 24 * time.Hour),
			Status:    StatusDraft,
			LineItems: []LineItem{
				{Description: "Service 1", Quantity: 1, UnitPrice: 100.0},
			},
			TaxRate:   0.1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		err := store.CreateInvoice(invoice)
		if err != nil {
			t.Fatalf("CreateInvoice failed: %v", err)
		}

		err = store.DeleteInvoice(invoice.ID)
		if err != nil {
			t.Fatalf("DeleteInvoice failed: %v", err)
		}

		_, err = store.GetInvoice(invoice.ID)
		if err == nil {
			t.Error("GetInvoice should return error after deletion")
		}
	})

	// Test DeleteInvoice not found
	t.Run("DeleteInvoice not found", func(t *testing.T) {
		err := store.DeleteInvoice("nonexistent")
		if err == nil {
			t.Error("DeleteInvoice should return error for non-existent invoice")
		}
	})

	// Test GetInvoicesByClient
	t.Run("GetInvoicesByClient", func(t *testing.T) {
		invoices, err := store.GetInvoicesByClient("client1")
		if err != nil {
			t.Fatalf("GetInvoicesByClient failed: %v", err)
		}

		if len(invoices) != 3 {
			t.Errorf("Expected 3 invoices for client1, got %d", len(invoices))
		}
	})

	// Test GetInvoiceByNumber
	t.Run("GetInvoiceByNumber", func(t *testing.T) {
		invoice, err := store.GetInvoiceByNumber("INV-001")
		if err != nil {
			t.Fatalf("GetInvoiceByNumber failed: %v", err)
		}

		if invoice.Number != "INV-001" {
			t.Errorf("Expected number INV-001, got %s", invoice.Number)
		}
	})

	// Test GetInvoiceByNumber not found
	t.Run("GetInvoiceByNumber not found", func(t *testing.T) {
		_, err := store.GetInvoiceByNumber("INV-999")
		if err == nil {
			t.Error("GetInvoiceByNumber should return error for non-existent invoice number")
		}
	})
}

func TestClientStore(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "store_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	store, err := NewStore(tempDir)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}

	// Test CreateClient
	t.Run("CreateClient", func(t *testing.T) {
		client := &Client{
			Name:         "Test Client",
			Email:        "test@example.com",
			Phone:        "123-456-7890",
			Address:      "123 Test St",
			PaymentTerms: 30,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		err := store.CreateClient(client)
		if err != nil {
			t.Fatalf("CreateClient failed: %v", err)
		}

		if client.ID == "" {
			t.Error("Client ID should be generated")
		}
	})

	// Test GetClient
	t.Run("GetClient", func(t *testing.T) {
		client := &Client{
			Name:         "Test Client 2",
			Email:        "test2@example.com",
			Phone:        "098-765-4321",
			Address:      "456 Test Ave",
			PaymentTerms: 15,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		err := store.CreateClient(client)
		if err != nil {
			t.Fatalf("CreateClient failed: %v", err)
		}

		retrieved, err := store.GetClient(client.ID)
		if err != nil {
			t.Fatalf("GetClient failed: %v", err)
		}

		if retrieved.ID != client.ID {
			t.Errorf("Expected ID %s, got %s", client.ID, retrieved.ID)
		}

		if retrieved.Name != client.Name {
			t.Errorf("Expected name %s, got %s", client.Name, retrieved.Name)
		}
	})

	// Test GetClient not found
	t.Run("GetClient not found", func(t *testing.T) {
		_, err := store.GetClient("nonexistent")
		if err == nil {
			t.Error("GetClient should return error for non-existent client")
		}
	})

	// Test UpdateClient
	t.Run("UpdateClient", func(t *testing.T) {
		client := &Client{
			Name:         "Test Client 3",
			Email:        "test3@example.com",
			Phone:        "555-555-5555",
			Address:      "789 Test Blvd",
			PaymentTerms: 60,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		err := store.CreateClient(client)
		if err != nil {
			t.Fatalf("CreateClient failed: %v", err)
		}

		// Update the client
		updatedEmail := "updated@example.com"
		client.Email = updatedEmail

		err = store.UpdateClient(client)
		if err != nil {
			t.Fatalf("UpdateClient failed: %v", err)
		}

		retrieved, err := store.GetClient(client.ID)
		if err != nil {
			t.Fatalf("GetClient failed: %v", err)
		}

		if retrieved.Email != updatedEmail {
			t.Errorf("Expected email %s, got %s", updatedEmail, retrieved.Email)
		}
	})

	// Test DeleteClient
	t.Run("DeleteClient", func(t *testing.T) {
		client := &Client{
			Name:         "Test Client 4",
			Email:        "test4@example.com",
			Phone:        "444-444-4444",
			Address:      "101 Test Ln",
			PaymentTerms: 45,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		err := store.CreateClient(client)
		if err != nil {
			t.Fatalf("CreateClient failed: %v", err)
		}

		err = store.DeleteClient(client.ID)
		if err != nil {
			t.Fatalf("DeleteClient failed: %v", err)
		}

		_, err = store.GetClient(client.ID)
		if err == nil {
			t.Error("GetClient should return error after deletion")
		}
	})

	// Test DeleteClient not found
	t.Run("DeleteClient not found", func(t *testing.T) {
		err := store.DeleteClient("nonexistent")
		if err == nil {
			t.Error("DeleteClient should return error for non-existent client")
		}
	})

	// Test ListClients (after one was deleted above, expect 3)
	t.Run("ListClients", func(t *testing.T) {
		clients, err := store.ListClients()
		if err != nil {
			t.Fatalf("ListClients failed: %v", err)
		}

		if len(clients) != 3 {
			t.Errorf("Expected 3 clients, got %d", len(clients))
		}
	})
}

func TestStoreInitialization(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "store_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Test that directories are created
	_, err = NewStore(tempDir)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}

	// Check if directories exist
	invoicesDir := filepath.Join(tempDir, "invoices")
	clientsDir := filepath.Join(tempDir, "clients")

	if _, err := os.Stat(invoicesDir); os.IsNotExist(err) {
		t.Error("Invoices directory should exist")
	}

	if _, err := os.Stat(clientsDir); os.IsNotExist(err) {
		t.Error("Clients directory should exist")
	}

	// Test with invalid directory
	_, err = NewStore("/invalid/path/that/does/not/exist")
	if err == nil {
		t.Error("NewStore should fail with invalid path")
	}
}
