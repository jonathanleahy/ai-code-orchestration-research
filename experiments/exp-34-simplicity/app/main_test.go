package main

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"app/store"
)

func initTestStore() {
	s = store.NewStore()

	clients := []*store.Client{
		{ID: "1", Name: "Acme Corp", Email: "contact@acme.com"},
		{ID: "2", Name: "TechStart Inc", Email: "hello@techstart.com"},
	}
	for _, c := range clients {
		s.AddClient(c)
	}

	projects := []*store.Project{
		{ID: "p1", ClientID: "1", Title: "Website Redesign"},
		{ID: "p2", ClientID: "2", Title: "Mobile App"},
	}
	for _, p := range projects {
		s.AddProject(p)
	}

	activities := []*store.Activity{
		{ID: "a1", ClientID: "1", Type: "meeting", Content: "Kickoff meeting held", Timestamp: "2026-03-20T10:00:00Z"},
		{ID: "a2", ClientID: "1", Type: "email", Content: "Quote sent", Timestamp: "2026-03-19T14:30:00Z"},
		{ID: "a3", ClientID: "2", Type: "meeting", Content: "Scope review", Timestamp: "2026-03-18T09:00:00Z"},
	}
	for _, a := range activities {
		s.AddActivity(a)
	}

	invoices := []*store.Invoice{
		{ID: "inv1", ProjectID: "p1", Amount: 5000.00},
		{ID: "inv2", ProjectID: "p1", Amount: 3000.00},
		{ID: "inv3", ProjectID: "p2", Amount: 12000.00},
	}
	for _, inv := range invoices {
		s.AddInvoice(inv)
	}
}

func setupMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", handleHealth)
	mux.HandleFunc("/", handleHome)
	mux.HandleFunc("/api/clients", handleClientsAPI)
	mux.HandleFunc("/api/activities", handleActivitiesAPI)
	mux.HandleFunc("/api/invoices", handleInvoicesAPI)
	return mux
}

func TestHealthEndpoint(t *testing.T) {
	initTestStore()
	mux := setupMux()

	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	if ct := w.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected content-type application/json, got %s", ct)
	}

	var result map[string]string
	if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if result["status"] != "ok" {
		t.Errorf("expected status 'ok', got %q", result["status"])
	}
}

func TestHomePageGET(t *testing.T) {
	initTestStore()
	mux := setupMux()

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	body := w.Body.String()
	if !strings.Contains(body, "Client Manager") {
		t.Error("response should contain 'Client Manager'")
	}

	if !strings.Contains(body, "Acme Corp") {
		t.Error("response should contain seeded client 'Acme Corp'")
	}

	if !strings.Contains(body, "Add New Client") {
		t.Error("response should contain form section")
	}
}

func TestHomePagePOSTAddClient(t *testing.T) {
	initTestStore()
	mux := setupMux()

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("name", "New Client")
	writer.WriteField("email", "new@example.com")
	writer.Close()

	req := httptest.NewRequest("POST", "/", body)
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

	if client.Name != "New Client" {
		t.Errorf("expected name 'New Client', got %q", client.Name)
	}

	if client.Email != "new@example.com" {
		t.Errorf("expected email 'new@example.com', got %q", client.Email)
	}

	// Verify client was added to store
	retrieved := s.GetClient(client.ID)
	if retrieved == nil {
		t.Error("client was not added to store")
	}
}

func TestHomePagePOSTMissingFields(t *testing.T) {
	initTestStore()
	mux := setupMux()

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("name", "Missing Email")
	writer.Close()

	req := httptest.NewRequest("POST", "/", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}

	var errResp map[string]string
	if err := json.NewDecoder(w.Body).Decode(&errResp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if !strings.Contains(errResp["error"], "required") {
		t.Errorf("expected error message about required fields, got %q", errResp["error"])
	}
}

func TestClientDetailPage(t *testing.T) {
	initTestStore()
	mux := setupMux()

	req := httptest.NewRequest("GET", "/client/1", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	body := w.Body.String()
	if !strings.Contains(body, "Acme Corp") {
		t.Error("response should contain client name")
	}

	if !strings.Contains(body, "contact@acme.com") {
		t.Error("response should contain client email")
	}

	if !strings.Contains(body, "Activity") {
		t.Error("response should contain Activity tab")
	}

	if !strings.Contains(body, "Invoices") {
		t.Error("response should contain Invoices tab")
	}
}

func TestClientDetailPageNotFound(t *testing.T) {
	initTestStore()
	mux := setupMux()

	req := httptest.NewRequest("GET", "/client/nonexistent", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}

	body := w.Body.String()
	if !strings.Contains(body, "not found") {
		t.Errorf("expected 'not found' message, got %q", body)
	}
}

func TestClientsAPI(t *testing.T) {
	initTestStore()
	mux := setupMux()

	req := httptest.NewRequest("GET", "/api/clients", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	if ct := w.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected content-type application/json, got %s", ct)
	}

	var clients []*store.Client
	if err := json.NewDecoder(w.Body).Decode(&clients); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(clients) != 2 {
		t.Errorf("expected 2 clients, got %d", len(clients))
	}

	found := false
	for _, c := range clients {
		if c.Name == "Acme Corp" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected to find 'Acme Corp' in clients")
	}
}

func TestActivitiesAPI(t *testing.T) {
	initTestStore()
	mux := setupMux()

	req := httptest.NewRequest("GET", "/api/activities?clientId=1", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var activities []*store.Activity
	if err := json.NewDecoder(w.Body).Decode(&activities); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(activities) != 2 {
		t.Errorf("expected 2 activities for client 1, got %d", len(activities))
	}

	for _, a := range activities {
		if a.ClientID != "1" {
			t.Errorf("expected clientId 1, got %s", a.ClientID)
		}
	}
}

func TestActivitiesAPIEmptyResult(t *testing.T) {
	initTestStore()
	mux := setupMux()

	req := httptest.NewRequest("GET", "/api/activities?clientId=nonexistent", nil)
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
		t.Errorf("expected empty activities, got %d", len(activities))
	}
}

func TestInvoicesAPI(t *testing.T) {
	initTestStore()
	mux := setupMux()

	req := httptest.NewRequest("GET", "/api/invoices?clientId=1", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var invoices []*store.Invoice
	if err := json.NewDecoder(w.Body).Decode(&invoices); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(invoices) != 2 {
		t.Errorf("expected 2 invoices for client 1, got %d", len(invoices))
	}

	totalAmount := 0.0
	for _, inv := range invoices {
		totalAmount += inv.Amount
	}

	expected := 8000.0
	if totalAmount != expected {
		t.Errorf("expected total amount %f, got %f", expected, totalAmount)
	}
}

func TestInvoicesAPIMultipleClients(t *testing.T) {
	initTestStore()
	mux := setupMux()

	// Client 1 has projects p1 with 2 invoices
	req1 := httptest.NewRequest("GET", "/api/invoices?clientId=1", nil)
	w1 := httptest.NewRecorder()
	mux.ServeHTTP(w1, req1)

	var invoices1 []*store.Invoice
	json.NewDecoder(w1.Body).Decode(&invoices1)

	// Client 2 has project p2 with 1 invoice
	req2 := httptest.NewRequest("GET", "/api/invoices?clientId=2", nil)
	w2 := httptest.NewRecorder()
	mux.ServeHTTP(w2, req2)

	var invoices2 []*store.Invoice
	json.NewDecoder(w2.Body).Decode(&invoices2)

	if len(invoices1) != 2 {
		t.Errorf("expected 2 invoices for client 1, got %d", len(invoices1))
	}

	if len(invoices2) != 1 {
		t.Errorf("expected 1 invoice for client 2, got %d", len(invoices2))
	}
}

func TestResponseHeaders(t *testing.T) {
	initTestStore()
	mux := setupMux()

	tests := []struct {
		path           string
		expectedCT     string
		expectedStatus int
	}{
		{"/health", "application/json", http.StatusOK},
		{"/api/clients", "application/json", http.StatusOK},
		{"/api/activities?clientId=1", "application/json", http.StatusOK},
		{"/api/invoices?clientId=1", "application/json", http.StatusOK},
	}

	for _, tt := range tests {
		req := httptest.NewRequest("GET", tt.path, nil)
		w := httptest.NewRecorder()

		mux.ServeHTTP(w, req)

		if w.Code != tt.expectedStatus {
			t.Errorf("[%s] expected status %d, got %d", tt.path, tt.expectedStatus, w.Code)
		}

		if ct := w.Header().Get("Content-Type"); ct != tt.expectedCT {
			t.Errorf("[%s] expected content-type %s, got %s", tt.path, tt.expectedCT, ct)
		}
	}
}

func TestHTMLResponses(t *testing.T) {
	initTestStore()
	mux := setupMux()

	tests := []struct {
		path          string
		expectedTexts []string
	}{
		{"/", []string{"<!DOCTYPE html", "Client Manager", "Add New Client"}},
		{"/client/1", []string{"<!DOCTYPE html", "Acme Corp", "Activity", "Invoices"}},
	}

	for _, tt := range tests {
		req := httptest.NewRequest("GET", tt.path, nil)
		w := httptest.NewRecorder()

		mux.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("[%s] expected status 200, got %d", tt.path, w.Code)
		}

		body := w.Body.String()
		for _, expectedText := range tt.expectedTexts {
			if !strings.Contains(body, expectedText) {
				t.Errorf("[%s] expected body to contain %q", tt.path, expectedText)
			}
		}
	}
}
