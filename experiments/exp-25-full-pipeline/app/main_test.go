package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"invoicegen/store"
)

func TestMain(m *testing.M) {
	// Use a temporary directory for test data
	tmpDir, err := os.MkdirTemp("", "invoicegen-test-")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create temp dir: %v\n", err)
		os.Exit(1)
	}
	defer os.RemoveAll(tmpDir)

	// Re-initialize the global store with test directory
	var initErr error
	invoiceStore, initErr = store.NewStore(tmpDir)
	if initErr != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize test store: %v\n", initErr)
		os.Exit(1)
	}

	code := m.Run()
	os.Exit(code)
}

func setupTestServer() *httptest.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", handleHealth)
	mux.HandleFunc("/", handleDashboard)
	mux.HandleFunc("/invoices/new", handleNewInvoiceForm)
	mux.HandleFunc("POST /api/invoices", handleCreateInvoice)
	mux.HandleFunc("GET /api/invoices", handleListInvoices)
	mux.HandleFunc("GET /api/invoices/{id}", handleGetInvoice)
	mux.HandleFunc("PATCH /api/invoices/{id}", handleUpdateInvoice)
	mux.HandleFunc("DELETE /api/invoices/{id}", handleDeleteInvoice)
	mux.HandleFunc("POST /api/invoices/{id}/send", handleSendInvoice)
	mux.HandleFunc("GET /api/invoices/{id}/preview", handlePreviewInvoice)

	return httptest.NewServer(mux)
}

func TestHealthEndpoint(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	resp, err := http.Get(server.URL + "/health")
	if err != nil {
		t.Fatalf("Failed to GET /health: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	var result map[string]string
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		t.Fatalf("Failed to decode health response: %v", err)
	}

	if result["status"] != "ok" {
		t.Errorf("Expected status 'ok', got '%s'", result["status"])
	}
}

func TestDashboardEndpoint(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	resp, err := http.Get(server.URL + "/")
	if err != nil {
		t.Fatalf("Failed to GET /: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	if !strings.Contains(resp.Header.Get("Content-Type"), "text/html") {
		t.Errorf("Expected Content-Type text/html, got %s", resp.Header.Get("Content-Type"))
	}

	var buf bytes.Buffer
	buf.ReadFrom(resp.Body)
	content := buf.String()

	if !strings.Contains(content, "Invoice Generator") {
		t.Error("Expected HTML content to contain 'Invoice Generator'")
	}
}

func TestCreateInvoice(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	invoiceData := store.Invoice{
		ClientID:  "client123",
		IssueDate: time.Now(),
		DueDate:   time.Now().AddDate(0, 0, 30),
		Status:    store.StatusDraft,
		TaxRate:   10.0,
		LineItems: []store.LineItem{
			{
				Description: "Service A",
				Quantity:    1,
				UnitPrice:   100.0,
			},
		},
		Notes: "Test invoice",
	}

	body, _ := json.Marshal(invoiceData)
	resp, err := http.Post(
		server.URL+"/api/invoices",
		"application/json",
		bytes.NewReader(body),
	)
	if err != nil {
		t.Fatalf("Failed to POST /api/invoices: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", resp.StatusCode)
	}

	var result store.Invoice
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if result.ID == "" {
		t.Error("Expected invoice ID to be generated")
	}

	if result.Number == "" {
		t.Error("Expected invoice number to be generated")
	}

	if result.ClientID != "client123" {
		t.Errorf("Expected ClientID 'client123', got '%s'", result.ClientID)
	}

	if result.Total != 110.0 { // 100 + 10% tax
		t.Errorf("Expected total 110.0, got %f", result.Total)
	}

	if !strings.Contains(resp.Header.Get("Content-Type"), "application/json") {
		t.Errorf("Expected Content-Type application/json, got %s", resp.Header.Get("Content-Type"))
	}
}

func TestListInvoices(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	// Create a test invoice first
	invoiceData := store.Invoice{
		ClientID:  "client123",
		IssueDate: time.Now(),
		DueDate:   time.Now().AddDate(0, 0, 30),
		Status:    store.StatusDraft,
		TaxRate:   0,
		LineItems: []store.LineItem{
			{
				Description: "Service A",
				Quantity:    1,
				UnitPrice:   100.0,
			},
		},
	}

	body, _ := json.Marshal(invoiceData)
	createResp, _ := http.Post(
		server.URL+"/api/invoices",
		"application/json",
		bytes.NewReader(body),
	)
	createResp.Body.Close()

	// Now list invoices
	resp, err := http.Get(server.URL + "/api/invoices")
	if err != nil {
		t.Fatalf("Failed to GET /api/invoices: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	var result []*store.Invoice
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(result) == 0 {
		t.Error("Expected at least one invoice in the list")
	}

	if result[0].ClientID != "client123" {
		t.Errorf("Expected ClientID 'client123', got '%s'", result[0].ClientID)
	}
}

func TestGetInvoiceByID(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	// Create a test invoice
	invoiceData := store.Invoice{
		ClientID:  "client456",
		IssueDate: time.Now(),
		DueDate:   time.Now().AddDate(0, 0, 30),
		Status:    store.StatusDraft,
		TaxRate:   5.0,
		LineItems: []store.LineItem{
			{
				Description: "Service B",
				Quantity:    2,
				UnitPrice:   50.0,
			},
		},
	}

	body, _ := json.Marshal(invoiceData)
	createResp, _ := http.Post(
		server.URL+"/api/invoices",
		"application/json",
		bytes.NewReader(body),
	)

	var createdInvoice store.Invoice
	json.NewDecoder(createResp.Body).Decode(&createdInvoice)
	createResp.Body.Close()

	// Get the invoice by ID
	resp, err := http.Get(server.URL + "/api/invoices/" + createdInvoice.ID)
	if err != nil {
		t.Fatalf("Failed to GET /api/invoices/{id}: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	var result store.Invoice
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if result.ID != createdInvoice.ID {
		t.Errorf("Expected ID '%s', got '%s'", createdInvoice.ID, result.ID)
	}

	if result.ClientID != "client456" {
		t.Errorf("Expected ClientID 'client456', got '%s'", result.ClientID)
	}

	if result.Total != 105.0 { // 100 + 5% tax
		t.Errorf("Expected total 105.0, got %f", result.Total)
	}
}

func TestGetNonexistentInvoice(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	resp, err := http.Get(server.URL + "/api/invoices/nonexistent")
	if err != nil {
		t.Fatalf("Failed to GET /api/invoices/nonexistent: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", resp.StatusCode)
	}

	var result map[string]string
	json.NewDecoder(resp.Body).Decode(&result)

	if result["error"] != "Invoice not found" {
		t.Errorf("Expected error 'Invoice not found', got '%s'", result["error"])
	}
}

func TestUpdateInvoice(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	// Create a test invoice
	invoiceData := store.Invoice{
		ClientID:  "client789",
		IssueDate: time.Now(),
		DueDate:   time.Now().AddDate(0, 0, 30),
		Status:    store.StatusDraft,
		TaxRate:   0,
		LineItems: []store.LineItem{
			{
				Description: "Service C",
				Quantity:    1,
				UnitPrice:   200.0,
			},
		},
	}

	body, _ := json.Marshal(invoiceData)
	createResp, _ := http.Post(
		server.URL+"/api/invoices",
		"application/json",
		bytes.NewReader(body),
	)

	var createdInvoice store.Invoice
	json.NewDecoder(createResp.Body).Decode(&createdInvoice)
	createResp.Body.Close()

	// Update the invoice
	updateData := map[string]interface{}{
		"status": store.StatusSent,
		"notes":  "Updated notes",
	}

	updateBody, _ := json.Marshal(updateData)
	req, _ := http.NewRequest(
		"PATCH",
		server.URL+"/api/invoices/"+createdInvoice.ID,
		bytes.NewReader(updateBody),
	)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to PATCH /api/invoices/{id}: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	var result store.Invoice
	json.NewDecoder(resp.Body).Decode(&result)

	if result.Status != store.StatusSent {
		t.Errorf("Expected status '%s', got '%s'", store.StatusSent, result.Status)
	}

	if result.Notes != "Updated notes" {
		t.Errorf("Expected notes 'Updated notes', got '%s'", result.Notes)
	}
}

func TestDeleteInvoice(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	// Create a test invoice
	invoiceData := store.Invoice{
		ClientID:  "client999",
		IssueDate: time.Now(),
		DueDate:   time.Now().AddDate(0, 0, 30),
		Status:    store.StatusDraft,
		TaxRate:   0,
		LineItems: []store.LineItem{
			{
				Description: "Service D",
				Quantity:    1,
				UnitPrice:   150.0,
			},
		},
	}

	body, _ := json.Marshal(invoiceData)
	createResp, _ := http.Post(
		server.URL+"/api/invoices",
		"application/json",
		bytes.NewReader(body),
	)

	var createdInvoice store.Invoice
	json.NewDecoder(createResp.Body).Decode(&createdInvoice)
	createResp.Body.Close()

	// Delete the invoice
	req, _ := http.NewRequest("DELETE", server.URL+"/api/invoices/"+createdInvoice.ID, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to DELETE /api/invoices/{id}: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	var result map[string]string
	json.NewDecoder(resp.Body).Decode(&result)

	if result["status"] != "deleted" {
		t.Errorf("Expected status 'deleted', got '%s'", result["status"])
	}

	// Verify the invoice is actually deleted
	getResp, _ := http.Get(server.URL + "/api/invoices/" + createdInvoice.ID)
	if getResp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected 404 after deletion, got %d", getResp.StatusCode)
	}
	getResp.Body.Close()
}

func TestDeleteNonexistentInvoice(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	req, _ := http.NewRequest("DELETE", server.URL+"/api/invoices/nonexistent", nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to DELETE /api/invoices/nonexistent: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", resp.StatusCode)
	}

	var result map[string]string
	json.NewDecoder(resp.Body).Decode(&result)

	if result["error"] != "Invoice not found" {
		t.Errorf("Expected error 'Invoice not found', got '%s'", result["error"])
	}
}

func TestPreviewInvoice(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	// Create a test invoice
	invoiceData := store.Invoice{
		ClientID:  "client_preview",
		IssueDate: time.Now(),
		DueDate:   time.Now().AddDate(0, 0, 30),
		Status:    store.StatusDraft,
		TaxRate:   0,
		LineItems: []store.LineItem{
			{
				Description: "Service Preview",
				Quantity:    1,
				UnitPrice:   100.0,
			},
		},
	}

	body, _ := json.Marshal(invoiceData)
	createResp, _ := http.Post(
		server.URL+"/api/invoices",
		"application/json",
		bytes.NewReader(body),
	)

	var createdInvoice store.Invoice
	json.NewDecoder(createResp.Body).Decode(&createdInvoice)
	createResp.Body.Close()

	// Get the preview
	resp, err := http.Get(server.URL + "/api/invoices/" + createdInvoice.ID + "/preview")
	if err != nil {
		t.Fatalf("Failed to GET /api/invoices/{id}/preview: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	if !strings.Contains(resp.Header.Get("Content-Type"), "text/html") {
		t.Errorf("Expected Content-Type text/html, got %s", resp.Header.Get("Content-Type"))
	}

	var buf bytes.Buffer
	buf.ReadFrom(resp.Body)
	content := buf.String()

	if !strings.Contains(content, createdInvoice.Number) {
		t.Errorf("Expected HTML to contain invoice number '%s'", createdInvoice.Number)
	}

	if !strings.Contains(content, "Service Preview") {
		t.Error("Expected HTML to contain line item description")
	}
}

func TestSendInvoice(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	// Create a test invoice
	invoiceData := store.Invoice{
		ClientID:  "client_send",
		IssueDate: time.Now(),
		DueDate:   time.Now().AddDate(0, 0, 30),
		Status:    store.StatusDraft,
		TaxRate:   0,
		LineItems: []store.LineItem{
			{
				Description: "Service Send",
				Quantity:    1,
				UnitPrice:   100.0,
			},
		},
	}

	body, _ := json.Marshal(invoiceData)
	createResp, _ := http.Post(
		server.URL+"/api/invoices",
		"application/json",
		bytes.NewReader(body),
	)

	var createdInvoice store.Invoice
	json.NewDecoder(createResp.Body).Decode(&createdInvoice)
	createResp.Body.Close()

	// Send the invoice
	req, _ := http.NewRequest(
		"POST",
		server.URL+"/api/invoices/"+createdInvoice.ID+"/send",
		nil,
	)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to POST /api/invoices/{id}/send: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	var result store.Invoice
	json.NewDecoder(resp.Body).Decode(&result)

	if result.Status != store.StatusSent {
		t.Errorf("Expected status '%s', got '%s'", store.StatusSent, result.Status)
	}
}

func TestInvoiceWithMultipleLineItems(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	invoiceData := store.Invoice{
		ClientID:  "client_multi",
		IssueDate: time.Now(),
		DueDate:   time.Now().AddDate(0, 0, 30),
		Status:    store.StatusDraft,
		TaxRate:   10.0,
		LineItems: []store.LineItem{
			{
				Description: "Service A",
				Quantity:    2,
				UnitPrice:   100.0,
			},
			{
				Description: "Service B",
				Quantity:    3,
				UnitPrice:   50.0,
			},
		},
	}

	body, _ := json.Marshal(invoiceData)
	resp, err := http.Post(
		server.URL+"/api/invoices",
		"application/json",
		bytes.NewReader(body),
	)
	if err != nil {
		t.Fatalf("Failed to POST /api/invoices: %v", err)
	}
	defer resp.Body.Close()

	var result store.Invoice
	json.NewDecoder(resp.Body).Decode(&result)

	// Subtotal: (2*100) + (3*50) = 200 + 150 = 350
	// Tax: 350 * 0.10 = 35
	// Total: 350 + 35 = 385
	expectedTotal := 385.0
	if result.Total != expectedTotal {
		t.Errorf("Expected total %f, got %f", expectedTotal, result.Total)
	}

	if len(result.LineItems) != 2 {
		t.Errorf("Expected 2 line items, got %d", len(result.LineItems))
	}
}

func TestCreateInvalidInvoice(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	// Invalid invoice with no client ID
	invoiceData := store.Invoice{
		IssueDate: time.Now(),
		DueDate:   time.Now().AddDate(0, 0, 30),
		Status:    store.StatusDraft,
		TaxRate:   0,
	}

	body, _ := json.Marshal(invoiceData)
	resp, err := http.Post(
		server.URL+"/api/invoices",
		"application/json",
		bytes.NewReader(body),
	)
	if err != nil {
		t.Fatalf("Failed to POST /api/invoices: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status 400 for invalid invoice, got %d", resp.StatusCode)
	}

	var result map[string]string
	json.NewDecoder(resp.Body).Decode(&result)

	if result["error"] == "" {
		t.Error("Expected error message in response")
	}
}
