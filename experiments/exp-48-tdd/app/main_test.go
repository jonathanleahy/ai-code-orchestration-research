package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"app/store"
)

func setupServer() *http.ServeMux {
	db = store.NewStore()

	mux := http.NewServeMux()

	mux.HandleFunc("/health", handleHealth)
	mux.HandleFunc("/", handleDashboard)

	// Clients
	mux.HandleFunc("/api/clients", handleClients)
	mux.HandleFunc("/api/clients/", handleClientDetail)

	// Activities
	mux.HandleFunc("/api/activities", handleActivities)

	// Invoices
	mux.HandleFunc("/api/invoices", handleInvoices)

	return mux
}

func makeRequest(mux *http.ServeMux, method, path string, body io.Reader) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, body)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	return w
}

func TestHealthCheck(t *testing.T) {
	mux := setupServer()
	w := makeRequest(mux, "GET", "/health", nil)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var result map[string]string
	json.NewDecoder(w.Body).Decode(&result)
	if result["status"] != "ok" {
		t.Errorf("expected status 'ok', got %v", result["status"])
	}
}

func TestDashboard(t *testing.T) {
	mux := setupServer()
	w := makeRequest(mux, "GET", "/", nil)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "text/html" {
		t.Errorf("expected text/html, got %s", contentType)
	}

	body := w.Body.String()
	if len(body) == 0 {
		t.Errorf("expected non-empty HTML body")
	}
	if !bytes.Contains(w.Body.Bytes(), []byte("Client Management Dashboard")) {
		t.Errorf("expected dashboard title in HTML")
	}
}

func TestGetClientsEmpty(t *testing.T) {
	mux := setupServer()
	w := makeRequest(mux, "GET", "/api/clients", nil)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var clients []*store.Client
	json.NewDecoder(w.Body).Decode(&clients)
	if clients == nil || len(clients) != 0 {
		t.Errorf("expected empty clients list")
	}
}

func TestCreateClient(t *testing.T) {
	mux := setupServer()

	data := url.Values{}
	data.Set("name", "John Doe")
	data.Set("email", "john@example.com")
	data.Set("phone", "555-1234")

	req := httptest.NewRequest("POST", "/api/clients", bytes.NewBufferString(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", w.Code)
	}

	var client *store.Client
	json.NewDecoder(w.Body).Decode(&client)
	if client == nil || client.Name != "John Doe" || client.Email != "john@example.com" {
		t.Errorf("expected client with correct data, got %v", client)
	}
}

func TestCreateClientMissingName(t *testing.T) {
	mux := setupServer()

	data := url.Values{}
	data.Set("email", "john@example.com")
	data.Set("phone", "555-1234")

	req := httptest.NewRequest("POST", "/api/clients", bytes.NewBufferString(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}

	var result map[string]string
	json.NewDecoder(w.Body).Decode(&result)
	if result["error"] != "name is required" {
		t.Errorf("expected error message, got %v", result["error"])
	}
}

func TestCreateClientInvalidMethod(t *testing.T) {
	mux := setupServer()
	w := makeRequest(mux, "DELETE", "/api/clients", nil)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status 405, got %d", w.Code)
	}
}

func TestGetClientsWithData(t *testing.T) {
	mux := setupServer()

	// Create a client
	data := url.Values{}
	data.Set("name", "Jane Doe")
	data.Set("email", "jane@example.com")
	data.Set("phone", "555-5678")

	req := httptest.NewRequest("POST", "/api/clients", bytes.NewBufferString(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	// Get clients
	w = makeRequest(mux, "GET", "/api/clients", nil)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var clients []*store.Client
	json.NewDecoder(w.Body).Decode(&clients)
	if len(clients) != 1 || clients[0].Name != "Jane Doe" {
		t.Errorf("expected one client with name Jane Doe, got %v", clients)
	}
}

func TestGetClientDetail(t *testing.T) {
	mux := setupServer()

	// Create a client
	data := url.Values{}
	data.Set("name", "Bob Smith")
	data.Set("email", "bob@example.com")
	data.Set("phone", "555-9999")

	req := httptest.NewRequest("POST", "/api/clients", bytes.NewBufferString(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	// Get client detail
	w = makeRequest(mux, "GET", "/api/clients/1", nil)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var client *store.Client
	json.NewDecoder(w.Body).Decode(&client)
	if client == nil || client.Name != "Bob Smith" {
		t.Errorf("expected Bob Smith, got %v", client)
	}
}

func TestGetClientDetailNotFound(t *testing.T) {
	mux := setupServer()
	w := makeRequest(mux, "GET", "/api/clients/999", nil)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}

	var result map[string]string
	json.NewDecoder(w.Body).Decode(&result)
	if result["error"] != "client not found" {
		t.Errorf("expected client not found error, got %v", result["error"])
	}
}

func TestGetClientDetailInvalidID(t *testing.T) {
	mux := setupServer()
	w := makeRequest(mux, "GET", "/api/clients/invalid", nil)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}

	var result map[string]string
	json.NewDecoder(w.Body).Decode(&result)
	if result["error"] != "invalid client id" {
		t.Errorf("expected invalid client id error, got %v", result["error"])
	}
}

func TestUpdateClient(t *testing.T) {
	mux := setupServer()

	// Create a client
	data := url.Values{}
	data.Set("name", "Alice")
	data.Set("email", "alice@example.com")
	data.Set("phone", "555-1111")

	req := httptest.NewRequest("POST", "/api/clients", bytes.NewBufferString(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	// Update client
	updateData := url.Values{}
	updateData.Set("name", "Alice Updated")
	updateData.Set("email", "alice.new@example.com")
	updateData.Set("phone", "555-2222")

	req = httptest.NewRequest("PUT", "/api/clients/1", bytes.NewBufferString(updateData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var client *store.Client
	json.NewDecoder(w.Body).Decode(&client)
	if client == nil || client.Name != "Alice Updated" || client.Email != "alice.new@example.com" {
		t.Errorf("expected updated client, got %v", client)
	}
}

func TestUpdateClientMissingName(t *testing.T) {
	mux := setupServer()

	// Create a client
	data := url.Values{}
	data.Set("name", "Carol")
	data.Set("email", "carol@example.com")

	req := httptest.NewRequest("POST", "/api/clients", bytes.NewBufferString(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	// Update with missing name
	updateData := url.Values{}
	updateData.Set("email", "carol.new@example.com")

	req = httptest.NewRequest("PUT", "/api/clients/1", bytes.NewBufferString(updateData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestUpdateClientNotFound(t *testing.T) {
	mux := setupServer()

	updateData := url.Values{}
	updateData.Set("name", "David")
	updateData.Set("email", "david@example.com")

	req := httptest.NewRequest("PUT", "/api/clients/999", bytes.NewBufferString(updateData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}
}

func TestDeleteClient(t *testing.T) {
	mux := setupServer()

	// Create a client
	data := url.Values{}
	data.Set("name", "Eve")
	data.Set("email", "eve@example.com")

	req := httptest.NewRequest("POST", "/api/clients", bytes.NewBufferString(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	// Delete client
	w = makeRequest(mux, "DELETE", "/api/clients/1", nil)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", w.Code)
	}

	// Verify it's deleted
	w = makeRequest(mux, "GET", "/api/clients/1", nil)
	if w.Code != http.StatusNotFound {
		t.Errorf("expected deleted client to return 404, got %d", w.Code)
	}
}

func TestDeleteClientNotFound(t *testing.T) {
	mux := setupServer()
	w := makeRequest(mux, "DELETE", "/api/clients/999", nil)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}
}

func TestClientDetailInvalidMethod(t *testing.T) {
	mux := setupServer()

	// Create a client first
	data := url.Values{}
	data.Set("name", "Frank")
	data.Set("email", "frank@example.com")

	req := httptest.NewRequest("POST", "/api/clients", bytes.NewBufferString(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	// Try invalid method
	w = makeRequest(mux, "PATCH", "/api/clients/1", nil)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status 405, got %d", w.Code)
	}
}

func TestCreateActivity(t *testing.T) {
	mux := setupServer()

	// Create a client first
	clientData := url.Values{}
	clientData.Set("name", "Grace")
	clientData.Set("email", "grace@example.com")

	req := httptest.NewRequest("POST", "/api/clients", bytes.NewBufferString(clientData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	// Create activity
	actData := url.Values{}
	actData.Set("client_id", "1")
	actData.Set("type", "call")
	actData.Set("description", "Initial consultation")

	req = httptest.NewRequest("POST", "/api/activities", bytes.NewBufferString(actData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", w.Code)
	}

	var activity *store.Activity
	json.NewDecoder(w.Body).Decode(&activity)
	if activity == nil || activity.Type != "call" || activity.Description != "Initial consultation" {
		t.Errorf("expected activity with correct data, got %v", activity)
	}
}

func TestCreateActivityInvalidClientID(t *testing.T) {
	mux := setupServer()

	actData := url.Values{}
	actData.Set("client_id", "invalid")
	actData.Set("type", "email")
	actData.Set("description", "Follow up")

	req := httptest.NewRequest("POST", "/api/activities", bytes.NewBufferString(actData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}

	var result map[string]string
	json.NewDecoder(w.Body).Decode(&result)
	if result["error"] != "invalid client_id" {
		t.Errorf("expected invalid client_id error, got %v", result["error"])
	}
}

func TestCreateActivityClientNotFound(t *testing.T) {
	mux := setupServer()

	actData := url.Values{}
	actData.Set("client_id", "999")
	actData.Set("type", "email")
	actData.Set("description", "Follow up")

	req := httptest.NewRequest("POST", "/api/activities", bytes.NewBufferString(actData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}

	var result map[string]string
	json.NewDecoder(w.Body).Decode(&result)
	if result["error"] != "client not found" {
		t.Errorf("expected client not found error, got %v", result["error"])
	}
}

func TestGetActivities(t *testing.T) {
	mux := setupServer()

	// Create a client
	clientData := url.Values{}
	clientData.Set("name", "Henry")
	clientData.Set("email", "henry@example.com")

	req := httptest.NewRequest("POST", "/api/clients", bytes.NewBufferString(clientData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	// Create activity
	actData := url.Values{}
	actData.Set("client_id", "1")
	actData.Set("type", "meeting")
	actData.Set("description", "Project kickoff")

	req = httptest.NewRequest("POST", "/api/activities", bytes.NewBufferString(actData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	// Get activities
	w = makeRequest(mux, "GET", "/api/activities?client_id=1", nil)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var activities []*store.Activity
	json.NewDecoder(w.Body).Decode(&activities)
	if len(activities) != 1 || activities[0].Type != "meeting" {
		t.Errorf("expected one meeting activity, got %v", activities)
	}
}

func TestGetActivitiesMissingClientID(t *testing.T) {
	mux := setupServer()
	w := makeRequest(mux, "GET", "/api/activities", nil)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}

	var result map[string]string
	json.NewDecoder(w.Body).Decode(&result)
	if result["error"] != "client_id is required" {
		t.Errorf("expected client_id required error, got %v", result["error"])
	}
}

func TestGetActivitiesInvalidClientID(t *testing.T) {
	mux := setupServer()
	w := makeRequest(mux, "GET", "/api/activities?client_id=invalid", nil)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}

	var result map[string]string
	json.NewDecoder(w.Body).Decode(&result)
	if result["error"] != "invalid client_id" {
		t.Errorf("expected invalid client_id error, got %v", result["error"])
	}
}

func TestActivitiesInvalidMethod(t *testing.T) {
	mux := setupServer()
	w := makeRequest(mux, "DELETE", "/api/activities", nil)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status 405, got %d", w.Code)
	}
}

func TestGetInvoicesEmpty(t *testing.T) {
	mux := setupServer()
	w := makeRequest(mux, "GET", "/api/invoices", nil)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var invoices []*store.Invoice
	json.NewDecoder(w.Body).Decode(&invoices)
	if invoices == nil || len(invoices) != 0 {
		t.Errorf("expected empty invoices list")
	}
}

func TestCreateInvoice(t *testing.T) {
	mux := setupServer()

	// Create a client
	clientData := url.Values{}
	clientData.Set("name", "Iris")
	clientData.Set("email", "iris@example.com")

	req := httptest.NewRequest("POST", "/api/clients", bytes.NewBufferString(clientData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	// Create invoice
	invData := url.Values{}
	invData.Set("client_id", "1")
	invData.Set("number", "INV-001")
	invData.Set("amount", "1500.00")
	invData.Set("status", "sent")
	invData.Set("due_date", "2026-04-25")

	req = httptest.NewRequest("POST", "/api/invoices", bytes.NewBufferString(invData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", w.Code)
	}

	var invoice *store.Invoice
	json.NewDecoder(w.Body).Decode(&invoice)
	if invoice == nil || invoice.Number != "INV-001" || invoice.Amount != 1500.00 {
		t.Errorf("expected invoice with correct data, got %v", invoice)
	}
}

func TestCreateInvoiceInvalidClientID(t *testing.T) {
	mux := setupServer()

	invData := url.Values{}
	invData.Set("client_id", "invalid")
	invData.Set("number", "INV-002")
	invData.Set("amount", "2000.00")
	invData.Set("status", "draft")
	invData.Set("due_date", "2026-05-25")

	req := httptest.NewRequest("POST", "/api/invoices", bytes.NewBufferString(invData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}

	var result map[string]string
	json.NewDecoder(w.Body).Decode(&result)
	if result["error"] != "invalid client_id" {
		t.Errorf("expected invalid client_id error, got %v", result["error"])
	}
}

func TestCreateInvoiceInvalidAmount(t *testing.T) {
	mux := setupServer()

	// Create a client
	clientData := url.Values{}
	clientData.Set("name", "Jack")
	clientData.Set("email", "jack@example.com")

	req := httptest.NewRequest("POST", "/api/clients", bytes.NewBufferString(clientData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	// Create invoice with invalid amount
	invData := url.Values{}
	invData.Set("client_id", "1")
	invData.Set("number", "INV-003")
	invData.Set("amount", "invalid")
	invData.Set("status", "draft")
	invData.Set("due_date", "2026-06-25")

	req = httptest.NewRequest("POST", "/api/invoices", bytes.NewBufferString(invData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}

	var result map[string]string
	json.NewDecoder(w.Body).Decode(&result)
	if result["error"] != "invalid amount" {
		t.Errorf("expected invalid amount error, got %v", result["error"])
	}
}

func TestCreateInvoiceInvalidDueDate(t *testing.T) {
	mux := setupServer()

	// Create a client
	clientData := url.Values{}
	clientData.Set("name", "Kelly")
	clientData.Set("email", "kelly@example.com")

	req := httptest.NewRequest("POST", "/api/clients", bytes.NewBufferString(clientData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	// Create invoice with invalid due date
	invData := url.Values{}
	invData.Set("client_id", "1")
	invData.Set("number", "INV-004")
	invData.Set("amount", "3000.00")
	invData.Set("status", "draft")
	invData.Set("due_date", "invalid-date")

	req = httptest.NewRequest("POST", "/api/invoices", bytes.NewBufferString(invData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404 for invalid due date, got %d", w.Code)
	}
}

func TestCreateInvoiceClientNotFound(t *testing.T) {
	mux := setupServer()

	invData := url.Values{}
	invData.Set("client_id", "999")
	invData.Set("number", "INV-005")
	invData.Set("amount", "4000.00")
	invData.Set("status", "draft")
	invData.Set("due_date", "2026-07-25")

	req := httptest.NewRequest("POST", "/api/invoices", bytes.NewBufferString(invData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}

	var result map[string]string
	json.NewDecoder(w.Body).Decode(&result)
	if result["error"] != "client not found or invalid due date" {
		t.Errorf("expected error, got %v", result["error"])
	}
}

func TestGetInvoicesWithData(t *testing.T) {
	mux := setupServer()

	// Create a client
	clientData := url.Values{}
	clientData.Set("name", "Leo")
	clientData.Set("email", "leo@example.com")

	req := httptest.NewRequest("POST", "/api/clients", bytes.NewBufferString(clientData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	// Create invoice
	invData := url.Values{}
	invData.Set("client_id", "1")
	invData.Set("number", "INV-006")
	invData.Set("amount", "5000.00")
	invData.Set("status", "paid")
	invData.Set("due_date", "2026-03-20")

	req = httptest.NewRequest("POST", "/api/invoices", bytes.NewBufferString(invData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	// Get invoices
	w = makeRequest(mux, "GET", "/api/invoices", nil)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var invoices []*store.Invoice
	json.NewDecoder(w.Body).Decode(&invoices)
	if len(invoices) != 1 || invoices[0].Number != "INV-006" {
		t.Errorf("expected one invoice, got %v", invoices)
	}
}

func TestInvoicesInvalidMethod(t *testing.T) {
	mux := setupServer()
	w := makeRequest(mux, "DELETE", "/api/invoices", nil)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status 405, got %d", w.Code)
	}
}

func TestContentTypeHeaders(t *testing.T) {
	mux := setupServer()

	endpoints := []struct {
		method   string
		path     string
		expected string
	}{
		{"GET", "/health", "application/json"},
		{"GET", "/", "text/html"},
		{"GET", "/api/clients", "application/json"},
		{"GET", "/api/invoices", "application/json"},
	}

	for _, ep := range endpoints {
		w := makeRequest(mux, ep.method, ep.path, nil)
		ct := w.Header().Get("Content-Type")
		if ct != ep.expected {
			t.Errorf("endpoint %s %s: expected Content-Type %s, got %s", ep.method, ep.path, ep.expected, ct)
		}
	}
}

func TestCreateActivityEmptyType(t *testing.T) {
	mux := setupServer()

	// Create a client
	clientData := url.Values{}
	clientData.Set("name", "Mia")
	clientData.Set("email", "mia@example.com")

	req := httptest.NewRequest("POST", "/api/clients", bytes.NewBufferString(clientData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	// Create activity with no validation error in handler, store allows it
	actData := url.Values{}
	actData.Set("client_id", "1")
	actData.Set("type", "")
	actData.Set("description", "No type")

	req = httptest.NewRequest("POST", "/api/activities", bytes.NewBufferString(actData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	// Handler allows empty type, store creates it
	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201 for empty type, got %d", w.Code)
	}
}

func TestIntegrationCreateClientActivityInvoice(t *testing.T) {
	mux := setupServer()

	// Create client
	clientData := url.Values{}
	clientData.Set("name", "Noah")
	clientData.Set("email", "noah@example.com")
	clientData.Set("phone", "555-8888")

	req := httptest.NewRequest("POST", "/api/clients", bytes.NewBufferString(clientData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("failed to create client: %d", w.Code)
	}

	var client *store.Client
	json.NewDecoder(w.Body).Decode(&client)
	clientID := fmt.Sprintf("%d", client.ID)

	// Create activity
	actData := url.Values{}
	actData.Set("client_id", clientID)
	actData.Set("type", "consultation")
	actData.Set("description", "Initial meeting")

	req = httptest.NewRequest("POST", "/api/activities", bytes.NewBufferString(actData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("failed to create activity: %d", w.Code)
	}

	// Create invoice
	invData := url.Values{}
	invData.Set("client_id", clientID)
	invData.Set("number", "INV-NOAH-001")
	invData.Set("amount", "7500.00")
	invData.Set("status", "draft")
	invData.Set("due_date", "2026-04-30")

	req = httptest.NewRequest("POST", "/api/invoices", bytes.NewBufferString(invData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("failed to create invoice: %d", w.Code)
	}

	var invoice *store.Invoice
	json.NewDecoder(w.Body).Decode(&invoice)

	if invoice.Amount != 7500.00 || invoice.Number != "INV-NOAH-001" {
		t.Errorf("expected correct invoice, got %v", invoice)
	}
}
