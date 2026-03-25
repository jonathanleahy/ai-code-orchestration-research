package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"app/store"
)

func TestHealthEndpoint(t *testing.T) {
	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	handleHealth(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `{"status":"ok"}`
	if w.Body.String() != expected+"\n" {
		t.Errorf("handler returned unexpected body: got %v want %v", w.Body.String(), expected)
	}
}

func TestDashboardEndpoint(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	handleDashboard(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	if ct := w.Header().Get("Content-Type"); ct != "text/html; charset=utf-8" {
		t.Errorf("handler returned wrong content type: got %v", ct)
	}

	if w.Body.Len() == 0 {
		t.Error("handler returned empty body")
	}
}

func TestCreateClient(t *testing.T) {
	// Reset store for clean test
	s = store.NewStore()

	clientData := `{"name":"Test Client","email":"test@example.com","phone":"123-456-7890","preferred_payment":"stripe"}`
	req := httptest.NewRequest("POST", "/api/clients", bytes.NewBufferString(clientData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handleClientsAPI(w, req)

	if status := w.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	var response store.Client
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Errorf("failed to decode response: %v", err)
	}

	if response.Name != "Test Client" || response.Email != "test@example.com" {
		t.Errorf("response client data incorrect: got %+v", response)
	}

	if response.ID == "" {
		t.Error("response client missing ID")
	}
}

func TestCreateClientDuplicate(t *testing.T) {
	s = store.NewStore()

	clientData := `{"name":"Duplicate Client","email":"dup1@example.com","phone":"123-456-7890"}`
	req := httptest.NewRequest("POST", "/api/clients", bytes.NewBufferString(clientData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handleClientsAPI(w, req)

	if status := w.Code; status != http.StatusCreated {
		t.Fatalf("first create failed: got %v", status)
	}

	// Try to create duplicate
	clientData2 := `{"name":"Duplicate Client","email":"dup2@example.com","phone":"999-999-9999"}`
	req2 := httptest.NewRequest("POST", "/api/clients", bytes.NewBufferString(clientData2))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	handleClientsAPI(w2, req2)

	if status := w2.Code; status != http.StatusBadRequest {
		t.Errorf("duplicate create should fail: got %v", status)
	}
}

func TestListClients(t *testing.T) {
	s = store.NewStore()

	// Create two clients
	client1 := &store.Client{Name: "Alice", Email: "alice@example.com"}
	client2 := &store.Client{Name: "Bob", Email: "bob@example.com"}
	s.CreateClient(client1)
	s.CreateClient(client2)

	req := httptest.NewRequest("GET", "/api/clients", nil)
	w := httptest.NewRecorder()
	handleClientsAPI(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v", status)
	}

	var clients []*store.Client
	if err := json.NewDecoder(w.Body).Decode(&clients); err != nil {
		t.Errorf("failed to decode response: %v", err)
	}

	if len(clients) != 2 {
		t.Errorf("expected 2 clients, got %d", len(clients))
	}
}

func TestGetClient(t *testing.T) {
	s = store.NewStore()

	client := &store.Client{Name: "Test", Email: "test@example.com"}
	if err := s.CreateClient(client); err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	req := httptest.NewRequest("GET", "/api/client/"+client.ID, nil)
	w := httptest.NewRecorder()
	handleClientAPI(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v", status)
	}

	var response store.Client
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Errorf("failed to decode response: %v", err)
	}

	if response.ID != client.ID || response.Name != "Test" {
		t.Errorf("response data incorrect: got %+v", response)
	}
}

func TestGetClientNotFound(t *testing.T) {
	s = store.NewStore()

	req := httptest.NewRequest("GET", "/api/client/nonexistent", nil)
	w := httptest.NewRecorder()
	handleClientAPI(w, req)

	if status := w.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}
}

func TestUpdateClient(t *testing.T) {
	s = store.NewStore()

	client := &store.Client{Name: "Original", Email: "original@example.com", Phone: "123-456-7890"}
	if err := s.CreateClient(client); err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	updateData := `{"name":"Updated","email":"updated@example.com","phone":"999-999-9999"}`
	req := httptest.NewRequest("PATCH", "/api/client/"+client.ID, bytes.NewBufferString(updateData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handleClientAPI(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v", status)
	}

	var response store.Client
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Errorf("failed to decode response: %v", err)
	}

	if response.Name != "Updated" || response.Email != "updated@example.com" {
		t.Errorf("response data incorrect: got %+v", response)
	}
}

func TestUpdateClientNotFound(t *testing.T) {
	s = store.NewStore()

	updateData := `{"name":"Updated","email":"updated@example.com"}`
	req := httptest.NewRequest("PATCH", "/api/client/nonexistent", bytes.NewBufferString(updateData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handleClientAPI(w, req)

	if status := w.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v", status)
	}
}

func TestDeleteClient(t *testing.T) {
	s = store.NewStore()

	client := &store.Client{Name: "To Delete", Email: "delete@example.com"}
	if err := s.CreateClient(client); err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	req := httptest.NewRequest("DELETE", "/api/client/"+client.ID, nil)
	w := httptest.NewRecorder()
	handleClientAPI(w, req)

	if status := w.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNoContent)
	}

	// Verify client is deleted
	_, err := s.GetClient(client.ID)
	if err == nil {
		t.Error("client should be deleted")
	}
}

func TestDeleteClientNotFound(t *testing.T) {
	s = store.NewStore()

	req := httptest.NewRequest("DELETE", "/api/client/nonexistent", nil)
	w := httptest.NewRecorder()
	handleClientAPI(w, req)

	if status := w.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v", status)
	}
}

func TestCreateActivity(t *testing.T) {
	s = store.NewStore()

	client := &store.Client{Name: "Activity Test", Email: "activity@example.com"}
	s.CreateClient(client)

	activityData := `{"client_id":"` + client.ID + `","type":"call","details":"Initial consultation"}`
	req := httptest.NewRequest("POST", "/api/activities", bytes.NewBufferString(activityData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handleActivitiesAPI(w, req)

	if status := w.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	var response store.Activity
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Errorf("failed to decode response: %v", err)
	}

	if response.ClientID != client.ID || response.Type != "call" {
		t.Errorf("response data incorrect: got %+v", response)
	}

	if response.ID == "" {
		t.Error("response activity missing ID")
	}
}

func TestListActivities(t *testing.T) {
	s = store.NewStore()

	client := &store.Client{Name: "List Test", Email: "list@example.com"}
	s.CreateClient(client)

	activity1 := &store.Activity{ClientID: client.ID, Type: "email", Details: "Email sent"}
	activity2 := &store.Activity{ClientID: client.ID, Type: "call", Details: "Call made"}
	s.CreateActivity(activity1)
	s.CreateActivity(activity2)

	req := httptest.NewRequest("GET", "/api/activities", nil)
	w := httptest.NewRecorder()
	handleActivitiesAPI(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v", status)
	}

	var activities []*store.Activity
	if err := json.NewDecoder(w.Body).Decode(&activities); err != nil {
		t.Errorf("failed to decode response: %v", err)
	}

	if len(activities) != 2 {
		t.Errorf("expected 2 activities, got %d", len(activities))
	}
}

func TestGetActivity(t *testing.T) {
	s = store.NewStore()

	client := &store.Client{Name: "Get Activity Test", Email: "getactivity@example.com"}
	s.CreateClient(client)

	activity := &store.Activity{ClientID: client.ID, Type: "meeting", Details: "Client meeting"}
	if err := s.CreateActivity(activity); err != nil {
		t.Fatalf("failed to create activity: %v", err)
	}

	req := httptest.NewRequest("GET", "/api/activity/"+activity.ID, nil)
	w := httptest.NewRecorder()
	handleActivityAPI(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v", status)
	}

	var response store.Activity
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Errorf("failed to decode response: %v", err)
	}

	if response.ID != activity.ID || response.Type != "meeting" {
		t.Errorf("response data incorrect: got %+v", response)
	}
}

func TestGetActivityNotFound(t *testing.T) {
	s = store.NewStore()

	req := httptest.NewRequest("GET", "/api/activity/nonexistent", nil)
	w := httptest.NewRecorder()
	handleActivityAPI(w, req)

	if status := w.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}
}

func TestUpdateActivity(t *testing.T) {
	s = store.NewStore()

	client := &store.Client{Name: "Update Activity", Email: "updateactivity@example.com"}
	s.CreateClient(client)

	activity := &store.Activity{ClientID: client.ID, Type: "email", Details: "Original detail"}
	s.CreateActivity(activity)

	updateData := `{"client_id":"` + client.ID + `","type":"call","details":"Updated detail"}`
	req := httptest.NewRequest("PATCH", "/api/activity/"+activity.ID, bytes.NewBufferString(updateData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handleActivityAPI(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v", status)
	}

	var response store.Activity
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Errorf("failed to decode response: %v", err)
	}

	if response.Type != "call" || response.Details != "Updated detail" {
		t.Errorf("response data incorrect: got %+v", response)
	}
}

func TestUpdateActivityNotFound(t *testing.T) {
	s = store.NewStore()

	updateData := `{"client_id":"xyz","type":"call","details":"Updated"}`
	req := httptest.NewRequest("PATCH", "/api/activity/nonexistent", bytes.NewBufferString(updateData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handleActivityAPI(w, req)

	if status := w.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v", status)
	}
}

func TestDeleteActivity(t *testing.T) {
	s = store.NewStore()

	client := &store.Client{Name: "Delete Activity", Email: "deleteactivity@example.com"}
	s.CreateClient(client)

	activity := &store.Activity{ClientID: client.ID, Type: "note", Details: "To delete"}
	s.CreateActivity(activity)

	req := httptest.NewRequest("DELETE", "/api/activity/"+activity.ID, nil)
	w := httptest.NewRecorder()
	handleActivityAPI(w, req)

	if status := w.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNoContent)
	}

	// Verify activity is deleted
	_, err := s.GetActivity(activity.ID)
	if err == nil {
		t.Error("activity should be deleted")
	}
}

func TestDeleteActivityNotFound(t *testing.T) {
	s = store.NewStore()

	req := httptest.NewRequest("DELETE", "/api/activity/nonexistent", nil)
	w := httptest.NewRecorder()
	handleActivityAPI(w, req)

	if status := w.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v", status)
	}
}

func TestCreateInvoice(t *testing.T) {
	s = store.NewStore()

	client := &store.Client{Name: "Invoice Test", Email: "invoice@example.com"}
	s.CreateClient(client)

	invoiceData := `{"client_id":"` + client.ID + `","amount":1000.50,"status":"unpaid","due_date":"2026-04-25T00:00:00Z"}`
	req := httptest.NewRequest("POST", "/api/invoices", bytes.NewBufferString(invoiceData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handleInvoicesAPI(w, req)

	if status := w.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	var response store.Invoice
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Errorf("failed to decode response: %v", err)
	}

	if response.ClientID != client.ID || response.Amount != 1000.50 || response.Status != "unpaid" {
		t.Errorf("response data incorrect: got %+v", response)
	}

	if response.ID == "" {
		t.Error("response invoice missing ID")
	}
}

func TestListInvoices(t *testing.T) {
	s = store.NewStore()

	client := &store.Client{Name: "List Invoice", Email: "listinvoice@example.com"}
	s.CreateClient(client)

	invoice1 := &store.Invoice{ClientID: client.ID, Amount: 500.00, Status: "paid"}
	invoice2 := &store.Invoice{ClientID: client.ID, Amount: 1500.00, Status: "unpaid"}
	s.CreateInvoice(invoice1)
	s.CreateInvoice(invoice2)

	req := httptest.NewRequest("GET", "/api/invoices", nil)
	w := httptest.NewRecorder()
	handleInvoicesAPI(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v", status)
	}

	var invoices []*store.Invoice
	if err := json.NewDecoder(w.Body).Decode(&invoices); err != nil {
		t.Errorf("failed to decode response: %v", err)
	}

	if len(invoices) != 2 {
		t.Errorf("expected 2 invoices, got %d", len(invoices))
	}
}

func TestGetInvoice(t *testing.T) {
	s = store.NewStore()

	client := &store.Client{Name: "Get Invoice", Email: "getinvoice@example.com"}
	s.CreateClient(client)

	invoice := &store.Invoice{ClientID: client.ID, Amount: 750.00, Status: "due"}
	if err := s.CreateInvoice(invoice); err != nil {
		t.Fatalf("failed to create invoice: %v", err)
	}

	req := httptest.NewRequest("GET", "/api/invoice/"+invoice.ID, nil)
	w := httptest.NewRecorder()
	handleInvoiceAPI(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v", status)
	}

	var response store.Invoice
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Errorf("failed to decode response: %v", err)
	}

	if response.ID != invoice.ID || response.Amount != 750.00 {
		t.Errorf("response data incorrect: got %+v", response)
	}
}

func TestGetInvoiceNotFound(t *testing.T) {
	s = store.NewStore()

	req := httptest.NewRequest("GET", "/api/invoice/nonexistent", nil)
	w := httptest.NewRecorder()
	handleInvoiceAPI(w, req)

	if status := w.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}
}

func TestUpdateInvoice(t *testing.T) {
	s = store.NewStore()

	client := &store.Client{Name: "Update Invoice", Email: "updateinvoice@example.com"}
	s.CreateClient(client)

	invoice := &store.Invoice{ClientID: client.ID, Amount: 500.00, Status: "unpaid"}
	s.CreateInvoice(invoice)

	updateData := `{"client_id":"` + client.ID + `","amount":750.00,"status":"paid","due_date":"2026-04-25T00:00:00Z"}`
	req := httptest.NewRequest("PATCH", "/api/invoice/"+invoice.ID, bytes.NewBufferString(updateData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handleInvoiceAPI(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v", status)
	}

	var response store.Invoice
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Errorf("failed to decode response: %v", err)
	}

	if response.Amount != 750.00 || response.Status != "paid" {
		t.Errorf("response data incorrect: got %+v", response)
	}
}

func TestUpdateInvoiceNotFound(t *testing.T) {
	s = store.NewStore()

	updateData := `{"client_id":"xyz","amount":1000,"status":"paid","due_date":"2026-04-25T00:00:00Z"}`
	req := httptest.NewRequest("PATCH", "/api/invoice/nonexistent", bytes.NewBufferString(updateData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handleInvoiceAPI(w, req)

	if status := w.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v", status)
	}
}

func TestDeleteInvoice(t *testing.T) {
	s = store.NewStore()

	client := &store.Client{Name: "Delete Invoice", Email: "deleteinvoice@example.com"}
	s.CreateClient(client)

	invoice := &store.Invoice{ClientID: client.ID, Amount: 1000.00, Status: "unpaid"}
	s.CreateInvoice(invoice)

	req := httptest.NewRequest("DELETE", "/api/invoice/"+invoice.ID, nil)
	w := httptest.NewRecorder()
	handleInvoiceAPI(w, req)

	if status := w.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNoContent)
	}

	// Verify invoice is deleted
	_, err := s.GetInvoice(invoice.ID)
	if err == nil {
		t.Error("invoice should be deleted")
	}
}

func TestDeleteInvoiceNotFound(t *testing.T) {
	s = store.NewStore()

	req := httptest.NewRequest("DELETE", "/api/invoice/nonexistent", nil)
	w := httptest.NewRecorder()
	handleInvoiceAPI(w, req)

	if status := w.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v", status)
	}
}

func TestClientAPIBadRequest(t *testing.T) {
	s = store.NewStore()

	invalidData := `invalid json`
	req := httptest.NewRequest("POST", "/api/clients", bytes.NewBufferString(invalidData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handleClientsAPI(w, req)

	if status := w.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}

func TestActivityAPIBadRequest(t *testing.T) {
	s = store.NewStore()

	invalidData := `not valid json`
	req := httptest.NewRequest("POST", "/api/activities", bytes.NewBufferString(invalidData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handleActivitiesAPI(w, req)

	if status := w.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}

func TestInvoiceAPIBadRequest(t *testing.T) {
	s = store.NewStore()

	invalidData := `malformed`
	req := httptest.NewRequest("POST", "/api/invoices", bytes.NewBufferString(invalidData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handleInvoicesAPI(w, req)

	if status := w.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}

func TestClientsAPIMethodNotAllowed(t *testing.T) {
	s = store.NewStore()

	req := httptest.NewRequest("DELETE", "/api/clients", nil)
	w := httptest.NewRecorder()
	handleClientsAPI(w, req)

	if status := w.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusMethodNotAllowed)
	}
}

func TestActivitiesAPIMethodNotAllowed(t *testing.T) {
	s = store.NewStore()

	req := httptest.NewRequest("DELETE", "/api/activities", nil)
	w := httptest.NewRecorder()
	handleActivitiesAPI(w, req)

	if status := w.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusMethodNotAllowed)
	}
}

func TestInvoicesAPIMethodNotAllowed(t *testing.T) {
	s = store.NewStore()

	req := httptest.NewRequest("PUT", "/api/invoices", nil)
	w := httptest.NewRecorder()
	handleInvoicesAPI(w, req)

	if status := w.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusMethodNotAllowed)
	}
}

func TestClientAPIMethodNotAllowed(t *testing.T) {
	s = store.NewStore()

	client := &store.Client{Name: "Test", Email: "test@example.com"}
	s.CreateClient(client)

	req := httptest.NewRequest("POST", "/api/client/"+client.ID, nil)
	w := httptest.NewRecorder()
	handleClientAPI(w, req)

	if status := w.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusMethodNotAllowed)
	}
}

func TestActivityAPIMethodNotAllowed(t *testing.T) {
	s = store.NewStore()

	client := &store.Client{Name: "Test", Email: "test@example.com"}
	s.CreateClient(client)

	activity := &store.Activity{ClientID: client.ID, Type: "call", Details: "Test"}
	s.CreateActivity(activity)

	req := httptest.NewRequest("PUT", "/api/activity/"+activity.ID, nil)
	w := httptest.NewRecorder()
	handleActivityAPI(w, req)

	if status := w.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusMethodNotAllowed)
	}
}

func TestInvoiceAPIMethodNotAllowed(t *testing.T) {
	s = store.NewStore()

	client := &store.Client{Name: "Test", Email: "test@example.com"}
	s.CreateClient(client)

	invoice := &store.Invoice{ClientID: client.ID, Amount: 100.00}
	s.CreateInvoice(invoice)

	req := httptest.NewRequest("PUT", "/api/invoice/"+invoice.ID, nil)
	w := httptest.NewRecorder()
	handleInvoiceAPI(w, req)

	if status := w.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusMethodNotAllowed)
	}
}

func TestDashboardWith404(t *testing.T) {
	s = store.NewStore()

	req := httptest.NewRequest("GET", "/invalid", nil)
	w := httptest.NewRecorder()
	handleDashboard(w, req)

	if status := w.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}
}

func TestClientPageNotFound(t *testing.T) {
	s = store.NewStore()

	req := httptest.NewRequest("GET", "/client/nonexistent", nil)
	w := httptest.NewRecorder()
	handleClientPage(w, req)

	if status := w.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}
}

func TestClientPageSuccess(t *testing.T) {
	s = store.NewStore()

	client := &store.Client{Name: "Page Test", Email: "pagetest@example.com"}
	s.CreateClient(client)

	req := httptest.NewRequest("GET", "/client/"+client.ID, nil)
	w := httptest.NewRecorder()
	handleClientPage(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v", status)
	}

	if ct := w.Header().Get("Content-Type"); ct != "text/html; charset=utf-8" {
		t.Errorf("handler returned wrong content type: got %v", ct)
	}

	body := w.Body.String()
	if body == "" {
		t.Error("handler returned empty body")
	}
}

func TestContentTypeHeaders(t *testing.T) {
	s = store.NewStore()

	tests := []struct {
		name    string
		handler func(http.ResponseWriter, *http.Request)
		req     *http.Request
	}{
		{
			name:    "health endpoint",
			handler: handleHealth,
			req:     httptest.NewRequest("GET", "/health", nil),
		},
		{
			name:    "clients API",
			handler: handleClientsAPI,
			req:     httptest.NewRequest("GET", "/api/clients", nil),
		},
		{
			name:    "activities API",
			handler: handleActivitiesAPI,
			req:     httptest.NewRequest("GET", "/api/activities", nil),
		},
		{
			name:    "invoices API",
			handler: handleInvoicesAPI,
			req:     httptest.NewRequest("GET", "/api/invoices", nil),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			tt.handler(w, tt.req)

			if ct := w.Header().Get("Content-Type"); ct != "application/json" {
				t.Errorf("expected application/json, got %v", ct)
			}
		})
	}
}

func TestEmptyListResponses(t *testing.T) {
	s = store.NewStore()

	tests := []struct {
		name    string
		handler func(http.ResponseWriter, *http.Request)
		path    string
	}{
		{
			name:    "empty clients",
			handler: handleClientsAPI,
			path:    "/api/clients",
		},
		{
			name:    "empty activities",
			handler: handleActivitiesAPI,
			path:    "/api/activities",
		},
		{
			name:    "empty invoices",
			handler: handleInvoicesAPI,
			path:    "/api/invoices",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.path, nil)
			w := httptest.NewRecorder()
			tt.handler(w, req)

			if status := w.Code; status != http.StatusOK {
				t.Errorf("expected 200, got %v", status)
			}

			body := w.Body.String()
			if body != "[]\n" && body != "null\n" {
				// Some implementations may return null instead of []
				t.Errorf("unexpected response body: %v", body)
			}
		})
	}
}

func TestClientEdgeCases(t *testing.T) {
	s = store.NewStore()

	// Test client with all fields empty except name and email
	clientData := `{"name":"Minimal","email":"minimal@example.com"}`
	req := httptest.NewRequest("POST", "/api/clients", bytes.NewBufferString(clientData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handleClientsAPI(w, req)

	if status := w.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v", status)
	}

	var client store.Client
	if err := json.NewDecoder(w.Body).Decode(&client); err != nil {
		t.Errorf("failed to decode response: %v", err)
	}

	if client.Name != "Minimal" || client.Email != "minimal@example.com" {
		t.Errorf("client data incorrect: got %+v", client)
	}
}

func TestInvoiceZeroAmount(t *testing.T) {
	s = store.NewStore()

	client := &store.Client{Name: "Zero Invoice", Email: "zeroinv@example.com"}
	s.CreateClient(client)

	invoiceData := `{"client_id":"` + client.ID + `","amount":0,"status":"paid","due_date":"2026-04-25T00:00:00Z"}`
	req := httptest.NewRequest("POST", "/api/invoices", bytes.NewBufferString(invoiceData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handleInvoicesAPI(w, req)

	if status := w.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v", status)
	}
}

func TestActivityEmptyDetails(t *testing.T) {
	s = store.NewStore()

	client := &store.Client{Name: "Empty Activity", Email: "emptyactivity@example.com"}
	s.CreateClient(client)

	activityData := `{"client_id":"` + client.ID + `","type":"note","details":""}`
	req := httptest.NewRequest("POST", "/api/activities", bytes.NewBufferString(activityData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handleActivitiesAPI(w, req)

	if status := w.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v", status)
	}
}
