package store

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestAcceptCreateInvoice(t *testing.T) {
	store, cleanup := setupTestStore(t)
	defer cleanup()

	inv := &Invoice{
		ClientID:    "client1",
		Number:      "INV-001",
		IssueDate:   time.Now(),
		DueDate:     time.Now().Add(30 * 24 * time.Hour),
		Status:      StatusDraft,
		LineItems:   []LineItem{{Description: "Service", Quantity: 1, UnitPrice: 100.0}},
		TaxRate:     0.1,
		Notes:       "Test invoice",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := store.CreateInvoice(inv)
	if err != nil {
		t.Fatalf("CreateInvoice failed: %v", err)
	}

	if inv.ID == "" {
		t.Error("Invoice ID should be generated")
	}

	if inv.Subtotal != 100.0 {
		t.Errorf("Expected subtotal 100.0, got %f", inv.Subtotal)
	}

	if inv.TaxAmount != 10.0 {
		t.Errorf("Expected tax amount 10.0, got %f", inv.TaxAmount)
	}

	if inv.Total != 110.0 {
		t.Errorf("Expected total 110.0, got %f", inv.Total)
	}
}

func TestAcceptCreateInvoiceMissingClient(t *testing.T) {
	store, cleanup := setupTestStore(t)
	defer cleanup()

	inv := &Invoice{
		ClientID:  "nonexistent",
		Number:    "INV-002",
		IssueDate: time.Now(),
		DueDate:   time.Now().Add(30 * 24 * time.Hour),
		Status:    StatusDraft,
		LineItems: []LineItem{{Description: "Service", Quantity: 1, UnitPrice: 100.0}},
		TaxRate:   0.1,
	}

	err := store.CreateInvoice(inv)
	if err == nil {
		t.Error("Expected error for missing client")
	}
}

func TestAcceptGetInvoice(t *testing.T) {
	store, cleanup := setupTestStore(t)
	defer cleanup()

	inv := &Invoice{
		ClientID:    "client1",
		Number:      "INV-003",
		IssueDate:   time.Now(),
		DueDate:     time.Now().Add(30 * 24 * time.Hour),
		Status:      StatusDraft,
		LineItems:   []LineItem{{Description: "Service", Quantity: 1, UnitPrice: 100.0}},
		TaxRate:     0.1,
		Notes:       "Test invoice",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := store.CreateInvoice(inv)
	if err != nil {
		t.Fatalf("CreateInvoice failed: %v", err)
	}

	fetched, err := store.GetInvoice(inv.ID)
	if err != nil {
		t.Fatalf("GetInvoice failed: %v", err)
	}

	if fetched.ID != inv.ID {
		t.Errorf("Expected ID %s, got %s", inv.ID, fetched.ID)
	}

	if fetched.Number != inv.Number {
		t.Errorf("Expected number %s, got %s", inv.Number, fetched.Number)
	}

	if fetched.ClientID != inv.ClientID {
		t.Errorf("Expected client ID %s, got %s", inv.ClientID, fetched.ClientID)
	}
}

func TestAcceptListInvoices(t *testing.T) {
	store, cleanup := setupTestStore(t)
	defer cleanup()

	inv1 := &Invoice{
		ClientID:    "client1",
		Number:      "INV-004",
		IssueDate:   time.Now(),
		DueDate:     time.Now().Add(30 * 24 * time.Hour),
		Status:      StatusDraft,
		LineItems:   []LineItem{{Description: "Service", Quantity: 1, UnitPrice: 100.0}},
		TaxRate:     0.1,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	inv2 := &Invoice{
		ClientID:    "client1",
		Number:      "INV-005",
		IssueDate:   time.Now(),
		DueDate:     time.Now().Add(30 * 24 * time.Hour),
		Status:      StatusDraft,
		LineItems:   []LineItem{{Description: "Service", Quantity: 1, UnitPrice: 100.0}},
		TaxRate:     0.1,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := store.CreateInvoice(inv1)
	if err != nil {
		t.Fatalf("CreateInvoice failed: %v", err)
	}

	err = store.CreateInvoice(inv2)
	if err != nil {
		t.Fatalf("CreateInvoice failed: %v", err)
	}

	invoices, err := store.ListInvoices()
	if err != nil {
		t.Fatalf("ListInvoices failed: %v", err)
	}

	if len(invoices) != 2 {
		t.Errorf("Expected 2 invoices, got %d", len(invoices))
	}
}

func TestAcceptUpdateInvoice(t *testing.T) {
	store, cleanup := setupTestStore(t)
	defer cleanup()

	inv := &Invoice{
		ClientID:    "client1",
		Number:      "INV-006",
		IssueDate:   time.Now(),
		DueDate:     time.Now().Add(30 * 24 * time.Hour),
		Status:      StatusDraft,
		LineItems:   []LineItem{{Description: "Service", Quantity: 1, UnitPrice: 100.0}},
		TaxRate:     0.1,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := store.CreateInvoice(inv)
	if err != nil {
		t.Fatalf("CreateInvoice failed: %v", err)
	}

	originalID := inv.ID
	inv.Status = StatusSent
	inv.Number = "INV-006-updated"

	err = store.UpdateInvoice(inv)
	if err != nil {
		t.Fatalf("UpdateInvoice failed: %v", err)
	}

	fetched, err := store.GetInvoice(originalID)
	if err != nil {
		t.Fatalf("GetInvoice failed: %v", err)
	}

	if fetched.Status != StatusSent {
		t.Errorf("Expected status %s, got %s", StatusSent, fetched.Status)
	}

	if fetched.Number != "INV-006-updated" {
		t.Errorf("Expected number %s, got %s", "INV-006-updated", fetched.Number)
	}
}

func TestAcceptDeleteInvoice(t *testing.T) {
	store, cleanup := setupTestStore(t)
	defer cleanup()

	inv := &Invoice{
		ClientID:    "client1",
		Number:      "INV-007",
		IssueDate:   time.Now(),
		DueDate:     time.Now().Add(30 * 24 * time.Hour),
		Status:      StatusDraft,
		LineItems:   []LineItem{{Description: "Service", Quantity: 1, UnitPrice: 100.0}},
		TaxRate:     0.1,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := store.CreateInvoice(inv)
	if err != nil {
		t.Fatalf("CreateInvoice failed: %v", err)
	}

	err = store.DeleteInvoice(inv.ID)
	if err != nil {
		t.Fatalf("DeleteInvoice failed: %v", err)
	}

	_, err = store.GetInvoice(inv.ID)
	if err == nil {
		t.Error("Expected error when getting deleted invoice")
	}
}

func setupTestStore(t *testing.T) (*Store, func()) {
	tempDir, err := os.MkdirTemp("", "invoice_store_test")
	if err != nil {
		t.Fatal(err)
	}

	store, err := NewStore(tempDir)
	if err != nil {
		t.Fatal(err)
	}

	// Create a client for testing
	client := &Client{
		ID:           "client1",
		Name:         "Test Client",
		Email:        "test@example.com",
		Phone:        "123-456-7890",
		Address:      "123 Test St",
		PaymentTerms: 30,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Create a dummy client file for testing
	clientPath := filepath.Join(tempDir, "clients", "client1.json")
	err = os.MkdirAll(filepath.Dir(clientPath), 0755)
	if err != nil {
		t.Fatal(err)
	}
	
	// Write client data to file manually to simulate existing client
	clientData := `{"id":"client1","name":"Test Client","email":"test@example.com","phone":"123-456-7890","address":"123 Test St","payment_terms":30,"created_at":"` + client.CreatedAt.Format(time.RFC3339) + `","updated_at":"` + client.UpdatedAt.Format(time.RFC3339) + `"}`
	err = os.WriteFile(clientPath, []byte(clientData), 0644)
	if err != nil {
		t.Fatal(err)
	}

	return store, func() {
		os.RemoveAll(tempDir)
	}
}