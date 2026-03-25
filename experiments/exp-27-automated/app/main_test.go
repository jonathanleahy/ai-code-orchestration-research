package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"app/store"
)

func TestHealthEndpoint(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	handleHealth(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var result map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if result["status"] != "ok" {
		t.Errorf("expected status 'ok', got '%s'", result["status"])
	}

	if ct := w.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected Content-Type 'application/json', got '%s'", ct)
	}
}

func TestDashboardEndpoint(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	handleDashboard(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	if ct := w.Header().Get("Content-Type"); ct != "text/html" {
		t.Errorf("expected Content-Type 'text/html', got '%s'", ct)
	}

	body := w.Body.String()
	if !strings.Contains(body, "<!DOCTYPE html>") {
		t.Errorf("expected HTML content")
	}
	if !strings.Contains(body, "CRM Dashboard") {
		t.Errorf("expected 'CRM Dashboard' title in HTML")
	}
}

func TestCreateClient(t *testing.T) {
	s = store.NewStore()

	client := store.Client{
		ID:    "test-client-1",
		Name:  "John Doe",
		Email: "john@example.com",
		Phone: "555-1234",
		Address: store.Address{
			Street:  "123 Main St",
			City:    "Springfield",
			State:   "IL",
			Country: "USA",
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	body, _ := json.Marshal(client)
	req := httptest.NewRequest(http.MethodPost, "/api/clients", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handleClients(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, w.Code)
	}

	var result store.Client
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if result.ID != client.ID || result.Name != client.Name || result.Email != client.Email {
		t.Errorf("response client data doesn't match: %+v", result)
	}
}

func TestListClients(t *testing.T) {
	s = store.NewStore()

	c1 := &store.Client{ID: "c1", Name: "Client 1", Email: "c1@test.com", CreatedAt: time.Now(), UpdatedAt: time.Now()}
	c2 := &store.Client{ID: "c2", Name: "Client 2", Email: "c2@test.com", CreatedAt: time.Now(), UpdatedAt: time.Now()}
	s.CreateClient(c1)
	s.CreateClient(c2)

	req := httptest.NewRequest(http.MethodGet, "/api/clients", nil)
	w := httptest.NewRecorder()

	handleClients(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var results []*store.Client
	if err := json.Unmarshal(w.Body.Bytes(), &results); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("expected 2 clients, got %d", len(results))
	}
}

func TestGetClientByID(t *testing.T) {
	s = store.NewStore()

	client := &store.Client{ID: "test-id", Name: "Test Client", Email: "test@example.com", CreatedAt: time.Now(), UpdatedAt: time.Now()}
	s.CreateClient(client)

	req := httptest.NewRequest(http.MethodGet, "/api/clients/test-id", nil)
	w := httptest.NewRecorder()

	handleClientByID(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var result store.Client
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if result.ID != "test-id" || result.Name != "Test Client" {
		t.Errorf("response doesn't match created client")
	}
}

func TestGetClientNotFound(t *testing.T) {
	s = store.NewStore()

	req := httptest.NewRequest(http.MethodGet, "/api/clients/nonexistent", nil)
	w := httptest.NewRecorder()

	handleClientByID(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
	}

	var result map[string]string
	json.Unmarshal(w.Body.Bytes(), &result)
	if result["error"] != "client not found" {
		t.Errorf("expected error 'client not found', got '%s'", result["error"])
	}
}

func TestUpdateClient(t *testing.T) {
	s = store.NewStore()

	original := &store.Client{ID: "update-test", Name: "Original", Email: "orig@test.com", CreatedAt: time.Now(), UpdatedAt: time.Now()}
	s.CreateClient(original)

	updated := store.Client{Name: "Updated Name", Email: "updated@test.com"}
	body, _ := json.Marshal(updated)
	req := httptest.NewRequest(http.MethodPatch, "/api/clients/update-test", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handleClientByID(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var result store.Client
	json.Unmarshal(w.Body.Bytes(), &result)
	if result.Name != "Updated Name" || result.Email != "updated@test.com" {
		t.Errorf("client not properly updated: %+v", result)
	}
}

func TestDeleteClient(t *testing.T) {
	s = store.NewStore()

	client := &store.Client{ID: "delete-test", Name: "To Delete", Email: "del@test.com", CreatedAt: time.Now(), UpdatedAt: time.Now()}
	s.CreateClient(client)

	req := httptest.NewRequest(http.MethodDelete, "/api/clients/delete-test", nil)
	w := httptest.NewRecorder()

	handleClientByID(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected status %d, got %d", http.StatusNoContent, w.Code)
	}

	retrieved, _ := s.GetClient("delete-test")
	if retrieved != nil {
		t.Errorf("expected client to be deleted")
	}
}

func TestCreateInvoice(t *testing.T) {
	s = store.NewStore()

	invoice := store.Invoice{
		ID:          "inv-1",
		ClientID:    "client-1",
		Items:       []store.Item{{Description: "Service", Quantity: 1, Price: 100.0, Total: 100.0}},
		TotalAmount: 100.0,
		Status:      "draft",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	body, _ := json.Marshal(invoice)
	req := httptest.NewRequest(http.MethodPost, "/api/invoices", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handleInvoices(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, w.Code)
	}

	var result store.Invoice
	json.Unmarshal(w.Body.Bytes(), &result)
	if result.ID != "inv-1" || result.ClientID != "client-1" || result.TotalAmount != 100.0 {
		t.Errorf("invoice data doesn't match: %+v", result)
	}
}

func TestListInvoices(t *testing.T) {
	s = store.NewStore()

	inv1 := &store.Invoice{
		ID:          "inv-1",
		ClientID:    "c1",
		TotalAmount: 100.0,
		Status:      "draft",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	inv2 := &store.Invoice{
		ID:          "inv-2",
		ClientID:    "c1",
		TotalAmount: 200.0,
		Status:      "sent",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	s.CreateInvoice(inv1)
	s.CreateInvoice(inv2)

	req := httptest.NewRequest(http.MethodGet, "/api/invoices", nil)
	w := httptest.NewRecorder()

	handleInvoices(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var results []*store.Invoice
	json.Unmarshal(w.Body.Bytes(), &results)
	if len(results) != 2 {
		t.Errorf("expected 2 invoices, got %d", len(results))
	}
}

func TestGetInvoiceByID(t *testing.T) {
	s = store.NewStore()

	invoice := &store.Invoice{
		ID:          "inv-get",
		ClientID:    "c1",
		TotalAmount: 50.0,
		Status:      "paid",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	s.CreateInvoice(invoice)

	req := httptest.NewRequest(http.MethodGet, "/api/invoices/inv-get", nil)
	w := httptest.NewRecorder()

	handleInvoiceByID(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var result store.Invoice
	json.Unmarshal(w.Body.Bytes(), &result)
	if result.ID != "inv-get" || result.TotalAmount != 50.0 {
		t.Errorf("invoice data doesn't match: %+v", result)
	}
}

func TestUpdateInvoice(t *testing.T) {
	s = store.NewStore()

	original := &store.Invoice{ID: "inv-update", ClientID: "c1", Status: "draft", CreatedAt: time.Now(), UpdatedAt: time.Now()}
	s.CreateInvoice(original)

	updated := store.Invoice{Status: "sent"}
	body, _ := json.Marshal(updated)
	req := httptest.NewRequest(http.MethodPatch, "/api/invoices/inv-update", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handleInvoiceByID(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var result store.Invoice
	json.Unmarshal(w.Body.Bytes(), &result)
	if result.Status != "sent" {
		t.Errorf("invoice status not updated: %s", result.Status)
	}
}

func TestDeleteInvoice(t *testing.T) {
	s = store.NewStore()

	invoice := &store.Invoice{ID: "inv-del", ClientID: "c1", CreatedAt: time.Now(), UpdatedAt: time.Now()}
	s.CreateInvoice(invoice)

	req := httptest.NewRequest(http.MethodDelete, "/api/invoices/inv-del", nil)
	w := httptest.NewRecorder()

	handleInvoiceByID(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected status %d, got %d", http.StatusNoContent, w.Code)
	}

	retrieved, _ := s.GetInvoice("inv-del")
	if retrieved != nil {
		t.Errorf("expected invoice to be deleted")
	}
}

func TestCreateComment(t *testing.T) {
	s = store.NewStore()

	comment := store.Comment{
		ID:        "comment-1",
		ClientID:  "client-1",
		Content:   "This is a test comment",
		CreatedAt: time.Now(),
	}

	body, _ := json.Marshal(comment)
	req := httptest.NewRequest(http.MethodPost, "/api/comments", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handleComments(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, w.Code)
	}

	var result store.Comment
	json.Unmarshal(w.Body.Bytes(), &result)
	if result.ID != "comment-1" || result.Content != "This is a test comment" {
		t.Errorf("comment data doesn't match: %+v", result)
	}
}

func TestListComments(t *testing.T) {
	s = store.NewStore()

	comment1 := &store.Comment{ID: "cm1", ClientID: "", Content: "Comment 1", CreatedAt: time.Now()}
	comment2 := &store.Comment{ID: "cm2", ClientID: "", Content: "Comment 2", CreatedAt: time.Now()}
	comment3 := &store.Comment{ID: "cm3", ClientID: "", Content: "Comment 3", CreatedAt: time.Now()}
	s.CreateComment(comment1)
	s.CreateComment(comment2)
	s.CreateComment(comment3)

	req := httptest.NewRequest(http.MethodGet, "/api/comments", nil)
	w := httptest.NewRecorder()

	handleComments(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var results []*store.Comment
	json.Unmarshal(w.Body.Bytes(), &results)
	if len(results) != 3 {
		t.Errorf("expected 3 comments, got %d", len(results))
	}
}

func TestDeleteComment(t *testing.T) {
	s = store.NewStore()

	comment := &store.Comment{ID: "cm-del", ClientID: "c1", Content: "To delete", CreatedAt: time.Now()}
	s.CreateComment(comment)

	req := httptest.NewRequest(http.MethodDelete, "/api/comments/cm-del", nil)
	w := httptest.NewRecorder()

	handleCommentByID(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected status %d, got %d", http.StatusNoContent, w.Code)
	}

	results, _ := s.GetComments("")
	for _, c := range results {
		if c.ID == "cm-del" {
			t.Errorf("expected comment to be deleted")
		}
	}
}

func TestCreateHistory(t *testing.T) {
	s = store.NewStore()

	entry := store.HistoryEntry{
		ID:        "hist-1",
		ClientID:  "client-1",
		Note:      "Initial contact",
		CreatedAt: time.Now(),
	}

	body, _ := json.Marshal(entry)
	req := httptest.NewRequest(http.MethodPost, "/api/history", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handleHistory(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, w.Code)
	}

	var result store.HistoryEntry
	json.Unmarshal(w.Body.Bytes(), &result)
	if result.ID != "hist-1" || result.Note != "Initial contact" {
		t.Errorf("history entry data doesn't match: %+v", result)
	}
}

func TestListHistory(t *testing.T) {
	s = store.NewStore()

	entry1 := &store.HistoryEntry{ID: "h1", ClientID: "", Note: "Note 1", CreatedAt: time.Now()}
	entry2 := &store.HistoryEntry{ID: "h2", ClientID: "", Note: "Note 2", CreatedAt: time.Now()}
	entry3 := &store.HistoryEntry{ID: "h3", ClientID: "", Note: "Note 3", CreatedAt: time.Now()}
	s.CreateHistory(entry1)
	s.CreateHistory(entry2)
	s.CreateHistory(entry3)

	req := httptest.NewRequest(http.MethodGet, "/api/history", nil)
	w := httptest.NewRecorder()

	handleHistory(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var results []*store.HistoryEntry
	json.Unmarshal(w.Body.Bytes(), &results)
	if len(results) != 3 {
		t.Errorf("expected 3 history entries, got %d", len(results))
	}
}

func TestBadRequestCreateClient(t *testing.T) {
	s = store.NewStore()

	req := httptest.NewRequest(http.MethodPost, "/api/clients", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handleClients(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	var result map[string]string
	json.Unmarshal(w.Body.Bytes(), &result)
	if result["error"] == "" {
		t.Errorf("expected error message in response")
	}
}

func TestInvalidClientID(t *testing.T) {
	s = store.NewStore()

	req := httptest.NewRequest(http.MethodGet, "/api/clients/", nil)
	w := httptest.NewRecorder()

	handleClientByID(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	var result map[string]string
	json.Unmarshal(w.Body.Bytes(), &result)
	if result["error"] != "invalid client id" {
		t.Errorf("expected error 'invalid client id', got '%s'", result["error"])
	}
}

func TestCORSHeaders(t *testing.T) {
	req := httptest.NewRequest(http.MethodOptions, "/api/clients", nil)
	w := httptest.NewRecorder()

	handler := corsMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	if w.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Errorf("expected CORS Access-Control-Allow-Origin header")
	}

	if w.Header().Get("Access-Control-Allow-Methods") != "GET, POST, PATCH, DELETE, OPTIONS" {
		t.Errorf("expected CORS Access-Control-Allow-Methods header")
	}
}

func TestMethodNotAllowed(t *testing.T) {
	s = store.NewStore()

	tests := []struct {
		path    string
		method  string
		handler http.HandlerFunc
	}{
		{"/api/clients", http.MethodDelete, handleClients},
		{"/api/invoices", http.MethodPatch, handleInvoices},
		{"/api/comments/1", http.MethodPatch, handleCommentByID},
		{"/api/history", http.MethodDelete, handleHistory},
	}

	for _, tt := range tests {
		req := httptest.NewRequest(tt.method, tt.path, nil)
		w := httptest.NewRecorder()
		tt.handler(w, req)

		if w.Code != http.StatusMethodNotAllowed {
			t.Errorf("path %s method %s: expected status %d, got %d", tt.path, tt.method, http.StatusMethodNotAllowed, w.Code)
		}
	}
}

func TestInvoiceNotFound(t *testing.T) {
	s = store.NewStore()

	req := httptest.NewRequest(http.MethodGet, "/api/invoices/nonexistent", nil)
	w := httptest.NewRecorder()

	handleInvoiceByID(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
	}

	var result map[string]string
	json.Unmarshal(w.Body.Bytes(), &result)
	if result["error"] != "invoice not found" {
		t.Errorf("expected error 'invoice not found', got '%s'", result["error"])
	}
}

func TestEmptyListResponses(t *testing.T) {
	s = store.NewStore()

	tests := []struct {
		path    string
		handler http.HandlerFunc
	}{
		{"/api/clients", handleClients},
		{"/api/invoices", handleInvoices},
		{"/api/comments", handleComments},
		{"/api/history", handleHistory},
	}

	for _, tt := range tests {
		req := httptest.NewRequest(http.MethodGet, tt.path, nil)
		w := httptest.NewRecorder()
		tt.handler(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("path %s: expected status %d, got %d", tt.path, http.StatusOK, w.Code)
		}

		var result []interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
			t.Errorf("path %s: failed to unmarshal response: %v", tt.path, err)
		}
		if result == nil {
			result = []interface{}{}
		}
		if len(result) != 0 {
			t.Errorf("path %s: expected empty list, got %d items", tt.path, len(result))
		}
	}
}

func TestDashboardNotFoundPath(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/invalid", nil)
	w := httptest.NewRecorder()

	handleDashboard(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}
