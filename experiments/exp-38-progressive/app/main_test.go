package main

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"app/store"
)

func setupTest() *http.ServeMux {
	s = store.NewStore()
	s.AddClient("Test User", "test@example.com", "555-0001")

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", handleHealth)
	mux.HandleFunc("GET /", handleDashboard)
	mux.HandleFunc("GET /client/{id}", handleClientDetail)
	mux.HandleFunc("POST /api/clients", handleCreateClient)
	mux.HandleFunc("GET /api/clients", handleListClients)
	mux.HandleFunc("POST /api/clients/{id}/update", handleUpdateClient)
	mux.HandleFunc("POST /api/clients/{id}/delete", handleDeleteClient)
	mux.HandleFunc("POST /api/clients/{id}/activities", handleCreateActivity)
	mux.HandleFunc("GET /api/clients/{id}/activities", handleListActivities)
	mux.HandleFunc("POST /api/clients/{id}/invoices", handleCreateInvoice)
	mux.HandleFunc("GET /api/clients/{id}/invoices", handleListInvoices)
	mux.HandleFunc("POST /api/invoices/{id}/pay", handleMarkInvoicePaid)
	mux.HandleFunc("GET /invoice/{id}/print", handlePrintInvoice)

	return mux
}

func TestHealthEndpoint(t *testing.T) {
	mux := setupTest()
	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var resp map[string]string
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp["status"] != "ok" {
		t.Errorf("expected status 'ok', got %s", resp["status"])
	}
}

func TestDashboardEndpoint(t *testing.T) {
	mux := setupTest()
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	if w.Header().Get("Content-Type") != "text/html; charset=utf-8" {
		t.Errorf("expected Content-Type text/html, got %s", w.Header().Get("Content-Type"))
	}

	if !bytes.Contains(w.Body.Bytes(), []byte("CRM Dashboard")) {
		t.Error("expected 'CRM Dashboard' in response")
	}

	if !bytes.Contains(w.Body.Bytes(), []byte("Add New Client")) {
		t.Error("expected 'Add New Client' in response")
	}

	if !bytes.Contains(w.Body.Bytes(), []byte("Test User")) {
		t.Error("expected test client in response")
	}
}

func TestListClientsAPI(t *testing.T) {
	mux := setupTest()
	req := httptest.NewRequest("GET", "/api/clients", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	if w.Header().Get("Content-Type") != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", w.Header().Get("Content-Type"))
	}

	var clients []*store.Client
	if err := json.NewDecoder(w.Body).Decode(&clients); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(clients) < 1 {
		t.Error("expected at least 1 client")
	}

	if clients[0].Name != "Test User" {
		t.Errorf("expected 'Test User', got %s", clients[0].Name)
	}
}

func TestCreateClientAPI(t *testing.T) {
	mux := setupTest()

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("name", "New Client")
	writer.WriteField("email", "new@example.com")
	writer.WriteField("phone", "555-0002")
	writer.Close()

	req := httptest.NewRequest("POST", "/api/clients", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", w.Code)
	}

	var client store.Client
	if err := json.NewDecoder(w.Body).Decode(&client); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if client.Name != "New Client" {
		t.Errorf("expected name 'New Client', got %s", client.Name)
	}

	if client.Email != "new@example.com" {
		t.Errorf("expected email 'new@example.com', got %s", client.Email)
	}

	if client.Phone != "555-0002" {
		t.Errorf("expected phone '555-0002', got %s", client.Phone)
	}

	if client.ID == "" {
		t.Error("expected non-empty client ID")
	}
}

func TestCreateClientMissingFields(t *testing.T) {
	mux := setupTest()

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("name", "Incomplete")
	writer.Close()

	req := httptest.NewRequest("POST", "/api/clients", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}

	var resp map[string]string
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp["error"] != "missing fields" {
		t.Errorf("expected error 'missing fields', got %s", resp["error"])
	}
}

func TestClientDetailPage(t *testing.T) {
	mux := setupTest()
	clients := s.ListClients()
	if len(clients) == 0 {
		t.Fatal("no clients available for testing")
	}

	clientID := clients[0].ID
	req := httptest.NewRequest("GET", "/client/"+clientID, nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	if w.Header().Get("Content-Type") != "text/html; charset=utf-8" {
		t.Errorf("expected Content-Type text/html, got %s", w.Header().Get("Content-Type"))
	}

	if !bytes.Contains(w.Body.Bytes(), []byte("Test User")) {
		t.Error("expected 'Test User' in response")
	}

	if !bytes.Contains(w.Body.Bytes(), []byte("Back to Dashboard")) {
		t.Error("expected 'Back to Dashboard' link in response")
	}
}

func TestClientDetailNotFound(t *testing.T) {
	mux := setupTest()
	req := httptest.NewRequest("GET", "/client/nonexistent-id", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}

	if !bytes.Contains(w.Body.Bytes(), []byte("not found")) {
		t.Error("expected 'not found' message in response")
	}
}

func TestUpdateClientAPI(t *testing.T) {
	mux := setupTest()
	clients := s.ListClients()
	if len(clients) == 0 {
		t.Fatal("no clients available for testing")
	}

	clientID := clients[0].ID

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("name", "Updated Name")
	writer.WriteField("email", "updated@example.com")
	writer.WriteField("phone", "555-9999")
	writer.Close()

	req := httptest.NewRequest("POST", "/api/clients/"+clientID+"/update", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var client store.Client
	if err := json.NewDecoder(w.Body).Decode(&client); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if client.Name != "Updated Name" {
		t.Errorf("expected name 'Updated Name', got %s", client.Name)
	}

	if client.Email != "updated@example.com" {
		t.Errorf("expected email 'updated@example.com', got %s", client.Email)
	}

	if client.Phone != "555-9999" {
		t.Errorf("expected phone '555-9999', got %s", client.Phone)
	}
}

func TestUpdateClientNotFound(t *testing.T) {
	mux := setupTest()

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("name", "Updated Name")
	writer.WriteField("email", "updated@example.com")
	writer.WriteField("phone", "555-9999")
	writer.Close()

	req := httptest.NewRequest("POST", "/api/clients/nonexistent-id/update", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}

	var resp map[string]string
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp["error"] != "client not found" {
		t.Errorf("expected error 'client not found', got %s", resp["error"])
	}
}

func TestUpdateClientMissingFields(t *testing.T) {
	mux := setupTest()
	clients := s.ListClients()
	if len(clients) == 0 {
		t.Fatal("no clients available for testing")
	}

	clientID := clients[0].ID

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("name", "Updated Name")
	writer.Close()

	req := httptest.NewRequest("POST", "/api/clients/"+clientID+"/update", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}

	var resp map[string]string
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp["error"] != "missing fields" {
		t.Errorf("expected error 'missing fields', got %s", resp["error"])
	}
}

func TestDeleteClientAPI(t *testing.T) {
	mux := setupTest()
	clients := s.ListClients()
	if len(clients) == 0 {
		t.Fatal("no clients available for testing")
	}

	clientID := clients[0].ID

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.Close()

	req := httptest.NewRequest("POST", "/api/clients/"+clientID+"/delete", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var resp map[string]bool
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if !resp["success"] {
		t.Error("expected success: true")
	}

	if s.GetClient(clientID) != nil {
		t.Error("expected client to be deleted")
	}
}

func TestDeleteClientNotFound(t *testing.T) {
	mux := setupTest()

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.Close()

	req := httptest.NewRequest("POST", "/api/clients/nonexistent-id/delete", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}

	var resp map[string]string
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp["error"] != "client not found" {
		t.Errorf("expected error 'client not found', got %s", resp["error"])
	}
}

func TestCreateActivityAPI(t *testing.T) {
	mux := setupTest()
	clients := s.ListClients()
	if len(clients) == 0 {
		t.Fatal("no clients available for testing")
	}

	clientID := clients[0].ID

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("type", "call")
	writer.WriteField("description", "Called to discuss project")
	writer.Close()

	req := httptest.NewRequest("POST", "/api/clients/"+clientID+"/activities", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", w.Code)
	}

	var activity store.Activity
	if err := json.NewDecoder(w.Body).Decode(&activity); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if activity.Type != "call" {
		t.Errorf("expected type 'call', got %s", activity.Type)
	}

	if activity.Description != "Called to discuss project" {
		t.Errorf("expected description 'Called to discuss project', got %s", activity.Description)
	}

	if activity.ClientID != clientID {
		t.Errorf("expected client_id %s, got %s", clientID, activity.ClientID)
	}

	if activity.ID == "" {
		t.Error("expected non-empty activity ID")
	}
}

func TestCreateActivityClientNotFound(t *testing.T) {
	mux := setupTest()

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("type", "email")
	writer.WriteField("description", "Sent email")
	writer.Close()

	req := httptest.NewRequest("POST", "/api/clients/nonexistent-id/activities", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}

	var resp map[string]string
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp["error"] != "client not found" {
		t.Errorf("expected error 'client not found', got %s", resp["error"])
	}
}

func TestCreateActivityMissingFields(t *testing.T) {
	mux := setupTest()
	clients := s.ListClients()
	if len(clients) == 0 {
		t.Fatal("no clients available for testing")
	}

	clientID := clients[0].ID

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("type", "meeting")
	writer.Close()

	req := httptest.NewRequest("POST", "/api/clients/"+clientID+"/activities", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}

	var resp map[string]string
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp["error"] != "missing fields" {
		t.Errorf("expected error 'missing fields', got %s", resp["error"])
	}
}

func TestListActivitiesAPI(t *testing.T) {
	mux := setupTest()
	clients := s.ListClients()
	if len(clients) == 0 {
		t.Fatal("no clients available for testing")
	}

	clientID := clients[0].ID

	s.AddActivity(clientID, "call", "Initial call")
	s.AddActivity(clientID, "email", "Follow up email")

	req := httptest.NewRequest("GET", "/api/clients/"+clientID+"/activities", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	if w.Header().Get("Content-Type") != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", w.Header().Get("Content-Type"))
	}

	var activities []*store.Activity
	if err := json.NewDecoder(w.Body).Decode(&activities); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(activities) != 2 {
		t.Errorf("expected 2 activities, got %d", len(activities))
	}

	if activities[0].Type != "call" {
		t.Errorf("expected first activity type 'call', got %s", activities[0].Type)
	}

	if activities[1].Type != "email" {
		t.Errorf("expected second activity type 'email', got %s", activities[1].Type)
	}
}

func TestListActivitiesClientNotFound(t *testing.T) {
	mux := setupTest()

	req := httptest.NewRequest("GET", "/api/clients/nonexistent-id/activities", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}

	var resp map[string]string
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp["error"] != "client not found" {
		t.Errorf("expected error 'client not found', got %s", resp["error"])
	}
}

func TestListActivitiesEmpty(t *testing.T) {
	mux := setupTest()
	clients := s.ListClients()
	if len(clients) == 0 {
		t.Fatal("no clients available for testing")
	}

	clientID := clients[0].ID

	req := httptest.NewRequest("GET", "/api/clients/"+clientID+"/activities", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var activities []*store.Activity
	if err := json.NewDecoder(w.Body).Decode(&activities); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(activities) != 0 {
		t.Errorf("expected 0 activities, got %d", len(activities))
	}
}

func TestCreateInvoiceAPI(t *testing.T) {
	mux := setupTest()
	clients := s.ListClients()
	if len(clients) == 0 {
		t.Fatal("no clients available for testing")
	}

	clientID := clients[0].ID

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("number", "INV-001")
	writer.WriteField("amount", "100.50")
	writer.WriteField("description", "Services rendered")
	writer.WriteField("due_date", "2026-04-25")
	writer.Close()

	req := httptest.NewRequest("POST", "/api/clients/"+clientID+"/invoices", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", w.Code)
	}

	var invoice store.Invoice
	if err := json.NewDecoder(w.Body).Decode(&invoice); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if invoice.Number != "INV-001" {
		t.Errorf("expected invoice number 'INV-001', got %s", invoice.Number)
	}

	if invoice.Amount != 100.50 {
		t.Errorf("expected amount 100.50, got %f", invoice.Amount)
	}

	if invoice.Status != "draft" {
		t.Errorf("expected status 'draft', got %s", invoice.Status)
	}

	if invoice.ClientID != clientID {
		t.Errorf("expected client_id %s, got %s", clientID, invoice.ClientID)
	}

	if invoice.ID == "" {
		t.Error("expected non-empty invoice ID")
	}
}

func TestCreateInvoiceClientNotFound(t *testing.T) {
	mux := setupTest()

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("number", "INV-001")
	writer.WriteField("amount", "100.00")
	writer.WriteField("description", "Services")
	writer.WriteField("due_date", "2026-04-25")
	writer.Close()

	req := httptest.NewRequest("POST", "/api/clients/nonexistent-id/invoices", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}

	var resp map[string]string
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp["error"] != "client not found" {
		t.Errorf("expected error 'client not found', got %s", resp["error"])
	}
}

func TestCreateInvoiceMissingFields(t *testing.T) {
	mux := setupTest()
	clients := s.ListClients()
	if len(clients) == 0 {
		t.Fatal("no clients available for testing")
	}

	clientID := clients[0].ID

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("number", "INV-001")
	writer.Close()

	req := httptest.NewRequest("POST", "/api/clients/"+clientID+"/invoices", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}

	var resp map[string]string
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp["error"] != "missing fields" {
		t.Errorf("expected error 'missing fields', got %s", resp["error"])
	}
}

func TestCreateInvoiceInvalidAmount(t *testing.T) {
	mux := setupTest()
	clients := s.ListClients()
	if len(clients) == 0 {
		t.Fatal("no clients available for testing")
	}

	clientID := clients[0].ID

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("number", "INV-001")
	writer.WriteField("amount", "invalid")
	writer.WriteField("description", "Services")
	writer.WriteField("due_date", "2026-04-25")
	writer.Close()

	req := httptest.NewRequest("POST", "/api/clients/"+clientID+"/invoices", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}

	var resp map[string]string
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp["error"] != "invalid amount" {
		t.Errorf("expected error 'invalid amount', got %s", resp["error"])
	}
}

func TestListInvoicesAPI(t *testing.T) {
	mux := setupTest()
	clients := s.ListClients()
	if len(clients) == 0 {
		t.Fatal("no clients available for testing")
	}

	clientID := clients[0].ID

	dueDate := "2026-04-25"
	s.CreateInvoice(clientID, "INV-001", "First invoice", 100.00, parseDate(dueDate))
	s.CreateInvoice(clientID, "INV-002", "Second invoice", 200.00, parseDate(dueDate))

	req := httptest.NewRequest("GET", "/api/clients/"+clientID+"/invoices", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	if w.Header().Get("Content-Type") != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", w.Header().Get("Content-Type"))
	}

	var invoices []*store.Invoice
	if err := json.NewDecoder(w.Body).Decode(&invoices); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(invoices) != 2 {
		t.Errorf("expected 2 invoices, got %d", len(invoices))
	}

	if invoices[0].Number != "INV-001" {
		t.Errorf("expected first invoice number 'INV-001', got %s", invoices[0].Number)
	}

	if invoices[1].Number != "INV-002" {
		t.Errorf("expected second invoice number 'INV-002', got %s", invoices[1].Number)
	}
}

func TestListInvoicesClientNotFound(t *testing.T) {
	mux := setupTest()

	req := httptest.NewRequest("GET", "/api/clients/nonexistent-id/invoices", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}

	var resp map[string]string
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp["error"] != "client not found" {
		t.Errorf("expected error 'client not found', got %s", resp["error"])
	}
}

func TestListInvoicesEmpty(t *testing.T) {
	mux := setupTest()
	clients := s.ListClients()
	if len(clients) == 0 {
		t.Fatal("no clients available for testing")
	}

	clientID := clients[0].ID

	req := httptest.NewRequest("GET", "/api/clients/"+clientID+"/invoices", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var invoices []*store.Invoice
	if err := json.NewDecoder(w.Body).Decode(&invoices); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(invoices) != 0 {
		t.Errorf("expected 0 invoices, got %d", len(invoices))
	}
}

func TestMarkInvoicePaidAPI(t *testing.T) {
	mux := setupTest()
	clients := s.ListClients()
	if len(clients) == 0 {
		t.Fatal("no clients available for testing")
	}

	clientID := clients[0].ID
	dueDate := "2026-04-25"
	invoice := s.CreateInvoice(clientID, "INV-001", "Test invoice", 150.00, parseDate(dueDate))

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.Close()

	req := httptest.NewRequest("POST", "/api/invoices/"+invoice.ID+"/pay", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var respInvoice store.Invoice
	if err := json.NewDecoder(w.Body).Decode(&respInvoice); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if respInvoice.Status != "paid" {
		t.Errorf("expected status 'paid', got %s", respInvoice.Status)
	}
}

func TestMarkInvoicePaidNotFound(t *testing.T) {
	mux := setupTest()

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.Close()

	req := httptest.NewRequest("POST", "/api/invoices/nonexistent-id/pay", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}

	var resp map[string]string
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp["error"] != "invoice not found" {
		t.Errorf("expected error 'invoice not found', got %s", resp["error"])
	}
}

func parseDate(dateStr string) time.Time {
	t, _ := time.Parse("2006-01-02", dateStr)
	return t
}

func TestSearchClientsAPI(t *testing.T) {
	mux := setupTest()
	s.AddClient("Alice Johnson", "alice@example.com", "555-0003")
	s.AddClient("Bob Smith", "bob@example.com", "555-0004")

	req := httptest.NewRequest("GET", "/api/clients?search=Alice", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	if w.Header().Get("Content-Type") != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", w.Header().Get("Content-Type"))
	}

	var clients []*store.Client
	if err := json.NewDecoder(w.Body).Decode(&clients); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(clients) != 1 {
		t.Errorf("expected 1 client, got %d", len(clients))
	}

	if clients[0].Name != "Alice Johnson" {
		t.Errorf("expected name 'Alice Johnson', got %s", clients[0].Name)
	}
}

func TestSearchClientsByEmail(t *testing.T) {
	mux := setupTest()
	s.AddClient("Charlie Brown", "charlie@example.com", "555-0005")
	s.AddClient("Diana Prince", "diana@example.com", "555-0006")

	req := httptest.NewRequest("GET", "/api/clients?search=diana@example.com", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var clients []*store.Client
	if err := json.NewDecoder(w.Body).Decode(&clients); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(clients) != 1 {
		t.Errorf("expected 1 client, got %d", len(clients))
	}

	if clients[0].Email != "diana@example.com" {
		t.Errorf("expected email 'diana@example.com', got %s", clients[0].Email)
	}
}

func TestSearchClientsEmpty(t *testing.T) {
	mux := setupTest()

	req := httptest.NewRequest("GET", "/api/clients?search=nonexistent", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var clients []*store.Client
	if err := json.NewDecoder(w.Body).Decode(&clients); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(clients) != 0 {
		t.Errorf("expected 0 clients, got %d", len(clients))
	}
}

func TestPrintInvoicePage(t *testing.T) {
	mux := setupTest()
	clients := s.ListClients()
	if len(clients) == 0 {
		t.Fatal("no clients available for testing")
	}

	clientID := clients[0].ID
	dueDate := "2026-04-25"
	invoice := s.CreateInvoice(clientID, "INV-001", "Test invoice", 150.00, parseDate(dueDate))

	req := httptest.NewRequest("GET", "/invoice/"+invoice.ID+"/print", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	if w.Header().Get("Content-Type") != "text/html; charset=utf-8" {
		t.Errorf("expected Content-Type text/html, got %s", w.Header().Get("Content-Type"))
	}

	if !bytes.Contains(w.Body.Bytes(), []byte("INV-001")) {
		t.Error("expected invoice number in response")
	}

	if !bytes.Contains(w.Body.Bytes(), []byte("Test User")) {
		t.Error("expected client name in response")
	}

	if !bytes.Contains(w.Body.Bytes(), []byte("150.00")) {
		t.Error("expected invoice amount in response")
	}

	if !bytes.Contains(w.Body.Bytes(), []byte("Print Invoice")) {
		t.Error("expected print button in response")
	}
}

func TestPrintInvoiceNotFound(t *testing.T) {
	mux := setupTest()

	req := httptest.NewRequest("GET", "/invoice/nonexistent-id/print", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}

	if !bytes.Contains(w.Body.Bytes(), []byte("not found")) {
		t.Error("expected 'not found' message in response")
	}
}
