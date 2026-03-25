package store

import (
	"testing"
	"time"
)

func TestSendInvoice(t *testing.T) {
	store := NewStore()

	invoice := Invoice{
		ID:            "invoice1",
		InvoiceNumber: "INV-001",
		ClientID:      "client1",
		InvoiceDate:   time.Now(),
		DueDate:       time.Now().Add(30 * 24 * time.Hour),
		LineItems: []LineItem{
			{
				ID:          "item1",
				Description: "Service",
				Quantity:    1,
				Rate:        100.0,
				Subtotal:    100.0,
			},
		},
		Subtotal:  100.0,
		Tax:       0,
		Total:     100.0,
		Status:    InvoiceStatusDraft,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := store.CreateInvoice(invoice)
	if err != nil {
		t.Fatalf("Failed to create invoice: %v", err)
	}

	err = store.SendInvoice(invoice.ID)
	if err != nil {
		t.Fatalf("Failed to send invoice: %v", err)
	}

	updatedInvoice, err := store.GetInvoice(invoice.ID)
	if err != nil {
		t.Fatalf("Failed to get invoice: %v", err)
	}

	if updatedInvoice.Status != InvoiceStatusSent {
		t.Errorf("Expected status sent, got %s", updatedInvoice.Status)
	}
}