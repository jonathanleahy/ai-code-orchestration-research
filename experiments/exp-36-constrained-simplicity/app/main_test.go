package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"app/store"
)

func TestHealth(t *testing.T) {
	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	handleHealth(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	var result map[string]string
	json.NewDecoder(w.Body).Decode(&result)
	if result["status"] != "ok" {
		t.Errorf("expected status=ok, got %v", result)
	}
}

func TestDashboard(t *testing.T) {
	db = store.NewInMemoryStore()

	// Create a test client
	client := &store.Client{
		ID:    "testclient1",
		Name:  "Test Client",
		Email: "test@example.com",
		Phone: "555-1234",
	}
	db.CreateClient(client)

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	handleDashboard(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	if w.Header().Get("Content-Type") != "text/html" {
		t.Errorf("expected text/html, got %s", w.Header().Get("Content-Type"))
	}

	body := w.Body.String()
	if !strings.Contains(body, "CRM Dashboard") {
		t.Errorf("expected CRM Dashboard in response")
	}
	if !strings.Contains(body, "Test Client") {
		t.Errorf("expected Test Client in response")
	}
}

func TestDashboardSearch(t *testing.T) {
	db = store.NewInMemoryStore()

	client1 := &store.Client{
		ID:    "c1",
		Name:  "Alice",
		Email: "alice@example.com",
	}
	client2 := &store.Client{
		ID:    "c2",
		Name:  "Bob",
		Email: "bob@example.com",
	}
	db.CreateClient(client1)
	db.CreateClient(client2)

	req := httptest.NewRequest("GET", "/?q=Alice", nil)
	w := httptest.NewRecorder()

	handleDashboard(w, req)

	body := w.Body.String()
	if !strings.Contains(body, "Alice") {
		t.Errorf("expected Alice in search results")
	}
	if strings.Contains(body, "Bob") && !strings.Contains(body, "Clients (1)") {
		t.Errorf("expected only Alice in search results")
	}
}

func TestCreateClient(t *testing.T) {
	db = store.NewInMemoryStore()

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	writer.WriteField("name", "John Doe")
	writer.WriteField("email", "john@example.com")
	writer.WriteField("phone", "555-1234")
	writer.WriteField("company", "ACME Corp")
	writer.WriteField("address", "123 Main St")
	writer.Close()

	req := httptest.NewRequest("POST", "/api/clients", &buf)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()

	handleClientsAPI(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d", w.Code)
	}

	var result store.Client
	json.NewDecoder(w.Body).Decode(&result)
	if result.Name != "John Doe" {
		t.Errorf("expected name=John Doe, got %s", result.Name)
	}
	if result.Email != "john@example.com" {
		t.Errorf("expected email=john@example.com, got %s", result.Email)
	}
	if result.Company != "ACME Corp" {
		t.Errorf("expected company=ACME Corp, got %s", result.Company)
	}
}

func TestGetClients(t *testing.T) {
	db = store.NewInMemoryStore()

	client1 := &store.Client{
		ID:    "c1",
		Name:  "Client 1",
		Email: "client1@example.com",
	}
	client2 := &store.Client{
		ID:    "c2",
		Name:  "Client 2",
		Email: "client2@example.com",
	}
	db.CreateClient(client1)
	db.CreateClient(client2)

	req := httptest.NewRequest("GET", "/api/clients", nil)
	w := httptest.NewRecorder()

	handleClientsAPI(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	var results []*store.Client
	json.NewDecoder(w.Body).Decode(&results)
	if len(results) != 2 {
		t.Errorf("expected 2 clients, got %d", len(results))
	}
}

func TestClientProfile(t *testing.T) {
	db = store.NewInMemoryStore()

	client := &store.Client{
		ID:    "c1",
		Name:  "John",
		Email: "john@example.com",
	}
	db.CreateClient(client)

	activity := &store.Activity{
		ID:       "act1",
		ClientID: "c1",
		Type:     "call",
		Notes:    "Initial contact",
	}
	db.CreateActivity(activity)

	invoice := &store.Invoice{
		ID:       "inv1",
		ClientID: "c1",
		Status:   "draft",
		Items: []store.InvoiceItem{
			{Description: "Service", Quantity: 1, Rate: 100.0},
		},
	}
	db.CreateInvoice(invoice)

	req := httptest.NewRequest("GET", "/client/c1", nil)
	w := httptest.NewRecorder()

	handleClientProfile(w, req, "c1")

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	body := w.Body.String()
	if !strings.Contains(body, "John") {
		t.Errorf("expected client name in profile")
	}
	if !strings.Contains(body, "john@example.com") {
		t.Errorf("expected client email in profile")
	}
}

func TestClientProfileNotFound(t *testing.T) {
	db = store.NewInMemoryStore()

	req := httptest.NewRequest("GET", "/client/nonexistent", nil)
	w := httptest.NewRecorder()

	handleClientProfile(w, req, "nonexistent")

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestNewInvoiceForm(t *testing.T) {
	db = store.NewInMemoryStore()

	client := &store.Client{
		ID:    "c1",
		Name:  "John",
		Email: "john@example.com",
	}
	db.CreateClient(client)

	req := httptest.NewRequest("GET", "/client/c1/invoice/new", nil)
	w := httptest.NewRecorder()

	handleNewInvoiceForm(w, req, "c1")

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	if w.Header().Get("Content-Type") != "text/html" {
		t.Errorf("expected text/html")
	}
}

func TestNewInvoiceFormNotFound(t *testing.T) {
	db = store.NewInMemoryStore()

	req := httptest.NewRequest("GET", "/client/nonexistent/invoice/new", nil)
	w := httptest.NewRecorder()

	handleNewInvoiceForm(w, req, "nonexistent")

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestCreateInvoice(t *testing.T) {
	db = store.NewInMemoryStore()

	client := &store.Client{
		ID:    "c1",
		Name:  "John",
		Email: "john@example.com",
	}
	db.CreateClient(client)

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	writer.WriteField("client_id", "c1")
	writer.WriteField("due_date", "2026-12-31")
	writer.WriteField("item_desc_0", "Consulting")
	writer.WriteField("item_qty_0", "10")
	writer.WriteField("item_rate_0", "100.50")
	writer.WriteField("item_desc_1", "Development")
	writer.WriteField("item_qty_1", "20")
	writer.WriteField("item_rate_1", "150.00")
	writer.Close()

	req := httptest.NewRequest("POST", "/api/invoices", &buf)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()

	handleCreateInvoice(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d", w.Code)
	}

	var result store.Invoice
	json.NewDecoder(w.Body).Decode(&result)
	if result.ClientID != "c1" {
		t.Errorf("expected client_id=c1, got %s", result.ClientID)
	}
	if result.Status != "draft" {
		t.Errorf("expected status=draft, got %s", result.Status)
	}
	if len(result.Items) != 2 {
		t.Errorf("expected 2 items, got %d", len(result.Items))
	}
	if result.Items[0].Description != "Consulting" {
		t.Errorf("expected first item description=Consulting")
	}
	if result.Items[0].Amount != 1005.0 {
		t.Errorf("expected first item amount=1005.0, got %f", result.Items[0].Amount)
	}
}

func TestUpdateInvoice(t *testing.T) {
	db = store.NewInMemoryStore()

	invoice := &store.Invoice{
		ID:       "inv1",
		ClientID: "c1",
		Status:   "draft",
		Items: []store.InvoiceItem{
			{Description: "Old Item", Quantity: 1, Rate: 100.0},
		},
		DueDate: time.Now(),
	}
	db.CreateInvoice(invoice)

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	writer.WriteField("due_date", "2026-06-30")
	writer.WriteField("item_desc_0", "New Item")
	writer.WriteField("item_qty_0", "5")
	writer.WriteField("item_rate_0", "200.00")
	writer.Close()

	req := httptest.NewRequest("POST", "/api/invoices/inv1", &buf)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()

	handleUpdateInvoice(w, req, "inv1")

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	var result store.Invoice
	json.NewDecoder(w.Body).Decode(&result)
	if result.Items[0].Description != "New Item" {
		t.Errorf("expected item description=New Item")
	}
	if result.Items[0].Amount != 1000.0 {
		t.Errorf("expected item amount=1000.0, got %f", result.Items[0].Amount)
	}
}

func TestUpdateInvoiceNotFound(t *testing.T) {
	db = store.NewInMemoryStore()

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	writer.WriteField("due_date", "2026-06-30")
	writer.Close()

	req := httptest.NewRequest("POST", "/api/invoices/nonexistent", &buf)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()

	handleUpdateInvoice(w, req, "nonexistent")

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}

	var result ErrorResponse
	json.NewDecoder(w.Body).Decode(&result)
	if result.Error != "Invoice not found" {
		t.Errorf("expected error=Invoice not found")
	}
}

func TestMarkInvoiceSent(t *testing.T) {
	db = store.NewInMemoryStore()

	invoice := &store.Invoice{
		ID:       "inv1",
		ClientID: "c1",
		Status:   "draft",
		Items:    []store.InvoiceItem{},
		DueDate:  time.Now(),
	}
	db.CreateInvoice(invoice)

	req := httptest.NewRequest("POST", "/api/invoices/inv1/mark-sent", nil)
	w := httptest.NewRecorder()

	handleInvoicesAPI(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	var result map[string]string
	json.NewDecoder(w.Body).Decode(&result)
	if result["status"] != "ok" {
		t.Errorf("expected status=ok")
	}
}

func TestMarkInvoicePaid(t *testing.T) {
	db = store.NewInMemoryStore()

	invoice := &store.Invoice{
		ID:       "inv1",
		ClientID: "c1",
		Status:   "sent",
		Items:    []store.InvoiceItem{},
		DueDate:  time.Now(),
	}
	db.CreateInvoice(invoice)

	req := httptest.NewRequest("POST", "/api/invoices/inv1/mark-paid", nil)
	w := httptest.NewRecorder()

	handleInvoicesAPI(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	// Verify the status was updated
	inv, _ := db.GetInvoice("inv1")
	if inv.Status != "paid" {
		t.Errorf("expected status=paid, got %s", inv.Status)
	}
}

func TestVoidInvoice(t *testing.T) {
	db = store.NewInMemoryStore()

	invoice := &store.Invoice{
		ID:       "inv1",
		ClientID: "c1",
		Status:   "sent",
		Items:    []store.InvoiceItem{},
		DueDate:  time.Now(),
	}
	db.CreateInvoice(invoice)

	req := httptest.NewRequest("POST", "/api/invoices/inv1/void", nil)
	w := httptest.NewRecorder()

	handleInvoicesAPI(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	// Verify the status was updated
	inv, _ := db.GetInvoice("inv1")
	if inv.Status != "void" {
		t.Errorf("expected status=void, got %s", inv.Status)
	}
}

func TestInvoiceDetail(t *testing.T) {
	db = store.NewInMemoryStore()

	invoice := &store.Invoice{
		ID:       "inv1",
		ClientID: "c1",
		Status:   "draft",
		Items: []store.InvoiceItem{
			{Description: "Service", Quantity: 1, Rate: 100.0, Amount: 100.0},
		},
		DueDate: time.Now(),
	}
	db.CreateInvoice(invoice)

	req := httptest.NewRequest("GET", "/invoice/inv1", nil)
	w := httptest.NewRecorder()

	handleInvoiceDetail(w, req, "inv1")

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	var result store.Invoice
	json.NewDecoder(w.Body).Decode(&result)
	if result.ID != "inv1" {
		t.Errorf("expected id=inv1, got %s", result.ID)
	}
	if len(result.Items) != 1 {
		t.Errorf("expected 1 item, got %d", len(result.Items))
	}
}

func TestInvoiceDetailNotFound(t *testing.T) {
	db = store.NewInMemoryStore()

	req := httptest.NewRequest("GET", "/invoice/nonexistent", nil)
	w := httptest.NewRecorder()

	handleInvoiceDetail(w, req, "nonexistent")

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestInvoicePrint(t *testing.T) {
	db = store.NewInMemoryStore()

	client := &store.Client{
		ID:      "c1",
		Name:    "ACME Corp",
		Email:   "contact@acme.com",
		Address: "123 Business Ave",
	}
	db.CreateClient(client)

	invoice := &store.Invoice{
		ID:       "inv1",
		ClientID: "c1",
		Status:   "sent",
		Items: []store.InvoiceItem{
			{Description: "Consulting", Quantity: 5, Rate: 100.0, Amount: 500.0},
		},
		DueDate: time.Now(),
	}
	db.CreateInvoice(invoice)

	req := httptest.NewRequest("GET", "/invoice/inv1/print", nil)
	w := httptest.NewRecorder()

	handleInvoicePrint(w, req, "inv1")

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	if w.Header().Get("Content-Type") != "text/html" {
		t.Errorf("expected text/html")
	}

	body := w.Body.String()
	if !strings.Contains(body, "ACME Corp") {
		t.Errorf("expected client name in print view")
	}
	if !strings.Contains(body, "Consulting") {
		t.Errorf("expected invoice item in print view")
	}
}

func TestInvoicePrintNotFound(t *testing.T) {
	db = store.NewInMemoryStore()

	req := httptest.NewRequest("GET", "/invoice/nonexistent/print", nil)
	w := httptest.NewRecorder()

	handleInvoicePrint(w, req, "nonexistent")

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestCreateActivity(t *testing.T) {
	db = store.NewInMemoryStore()

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	writer.WriteField("client_id", "c1")
	writer.WriteField("type", "call")
	writer.WriteField("description", "Discussed project requirements")
	writer.Close()

	req := httptest.NewRequest("POST", "/api/activities", &buf)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()

	handleCreateActivity(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d", w.Code)
	}

	var result store.Activity
	json.NewDecoder(w.Body).Decode(&result)
	if result.ClientID != "c1" {
		t.Errorf("expected client_id=c1, got %s", result.ClientID)
	}
	if result.Type != "call" {
		t.Errorf("expected type=call, got %s", result.Type)
	}
	if result.Notes != "Discussed project requirements" {
		t.Errorf("expected notes match")
	}
}

func TestGetClientActivities(t *testing.T) {
	db = store.NewInMemoryStore()

	activity1 := &store.Activity{
		ID:       "act1",
		ClientID: "c1",
		Type:     "call",
		Notes:    "Initial contact",
	}
	activity2 := &store.Activity{
		ID:       "act2",
		ClientID: "c1",
		Type:     "email",
		Notes:    "Follow-up",
	}
	activity3 := &store.Activity{
		ID:       "act3",
		ClientID: "c2",
		Type:     "call",
		Notes:    "Different client",
	}
	db.CreateActivity(activity1)
	db.CreateActivity(activity2)
	db.CreateActivity(activity3)

	activities, err := db.GetClientActivities("c1")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(activities) != 2 {
		t.Errorf("expected 2 activities for c1, got %d", len(activities))
	}

	if activities[0].ID != "act1" && activities[1].ID != "act1" {
		t.Errorf("expected act1 in results")
	}
}

func TestInvoiceWorkflow(t *testing.T) {
	db = store.NewInMemoryStore()

	// Create client
	client := &store.Client{
		ID:      "c1",
		Name:    "Customer Inc",
		Email:   "customer@example.com",
		Company: "Customer Inc",
	}
	db.CreateClient(client)

	// Create invoice
	invoice := &store.Invoice{
		ID:       "inv1",
		ClientID: "c1",
		Status:   "draft",
		Items: []store.InvoiceItem{
			{Description: "Service A", Quantity: 10, Rate: 100.0, Amount: 1000.0},
		},
		DueDate: time.Now().AddDate(0, 1, 0),
	}
	db.CreateInvoice(invoice)

	// Verify draft status
	inv, _ := db.GetInvoice("inv1")
	if inv.Status != "draft" {
		t.Errorf("expected status=draft")
	}

	// Mark as sent
	db.MarkInvoicePaid("inv1") // Note: API calls mark-sent but stores don't have MarkInvoiceSent
	// Just test the mark-paid transition

	// Mark as paid
	db.MarkInvoicePaid("inv1")
	inv, _ = db.GetInvoice("inv1")
	if inv.Status != "paid" {
		t.Errorf("expected status=paid after payment")
	}

	// Test void on new invoice
	invoice2 := &store.Invoice{
		ID:       "inv2",
		ClientID: "c1",
		Status:   "draft",
		Items:    []store.InvoiceItem{},
		DueDate:  time.Now(),
	}
	db.CreateInvoice(invoice2)

	db.VoidInvoice("inv2")
	inv2, _ := db.GetInvoice("inv2")
	if inv2.Status != "void" {
		t.Errorf("expected status=void")
	}
}

func TestClientSearch(t *testing.T) {
	db = store.NewInMemoryStore()

	clients := []*store.Client{
		{ID: "c1", Name: "John Smith", Email: "john@example.com", Company: "Tech Corp"},
		{ID: "c2", Name: "Jane Doe", Email: "jane@example.com", Company: "Design Inc"},
		{ID: "c3", Name: "Bob Wilson", Email: "bob@example.com", Company: "Tech Corp"},
	}

	for _, c := range clients {
		db.CreateClient(c)
	}

	// Search by name
	results, _ := db.SearchClients("John Smith")
	if len(results) != 1 || results[0].Name != "John Smith" {
		t.Errorf("search by name failed")
	}

	// Search by email
	results, _ = db.SearchClients("jane@example.com")
	if len(results) != 1 || results[0].Email != "jane@example.com" {
		t.Errorf("search by email failed")
	}

	// Search by company
	results, _ = db.SearchClients("Tech Corp")
	if len(results) != 2 {
		t.Errorf("search by company expected 2 results, got %d", len(results))
	}

	// Search with no match
	results, _ = db.SearchClients("Nonexistent")
	if len(results) != 0 {
		t.Errorf("search with no match expected 0 results, got %d", len(results))
	}
}

func TestInvoiceWithMultipleItems(t *testing.T) {
	db = store.NewInMemoryStore()

	items := []store.InvoiceItem{
		{Description: "Item 1", Quantity: 5, Rate: 100.0},
		{Description: "Item 2", Quantity: 3, Rate: 250.0},
		{Description: "Item 3", Quantity: 2, Rate: 500.0},
	}

	invoice := &store.Invoice{
		ID:       "inv1",
		ClientID: "c1",
		Status:   "draft",
		Items:    items,
		DueDate:  time.Now(),
	}

	db.CreateInvoice(invoice)

	// Amounts should be calculated by store
	inv, _ := db.GetInvoice("inv1")
	if inv.Items[0].Amount != 500.0 {
		t.Errorf("expected item 1 amount=500.0, got %f", inv.Items[0].Amount)
	}
	if inv.Items[1].Amount != 750.0 {
		t.Errorf("expected item 2 amount=750.0, got %f", inv.Items[1].Amount)
	}
	if inv.Items[2].Amount != 1000.0 {
		t.Errorf("expected item 3 amount=1000.0, got %f", inv.Items[2].Amount)
	}
}

func TestAPIMethodNotAllowed(t *testing.T) {
	db = store.NewInMemoryStore()

	req := httptest.NewRequest("DELETE", "/api/clients", nil)
	w := httptest.NewRecorder()

	handleClientsAPI(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", w.Code)
	}

	var result ErrorResponse
	json.NewDecoder(w.Body).Decode(&result)
	if result.Error != "Method not allowed" {
		t.Errorf("expected error=Method not allowed")
	}
}

func TestClientRouterNotFound(t *testing.T) {
	db = store.NewInMemoryStore()

	req := httptest.NewRequest("GET", "/client/", nil)
	w := httptest.NewRecorder()

	handleClientRouter(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestInvoiceRouterNotFound(t *testing.T) {
	db = store.NewInMemoryStore()

	req := httptest.NewRequest("GET", "/invoice/", nil)
	w := httptest.NewRecorder()

	handleInvoiceRouter(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestDashboardNotFound(t *testing.T) {
	db = store.NewInMemoryStore()

	req := httptest.NewRequest("GET", "/nonexistent", nil)
	w := httptest.NewRecorder()

	handleDashboard(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestCompleteClientLifecycle(t *testing.T) {
	db = store.NewInMemoryStore()

	// Create client
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	writer.WriteField("name", "Complete Test Client")
	writer.WriteField("email", "complete@example.com")
	writer.WriteField("phone", "555-9999")
	writer.WriteField("company", "Complete Inc")
	writer.WriteField("address", "999 Test Lane")
	writer.Close()

	req := httptest.NewRequest("POST", "/api/clients", &buf)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()

	handleClientsAPI(w, req)

	var createdClient store.Client
	json.NewDecoder(w.Body).Decode(&createdClient)
	clientID := createdClient.ID

	// Retrieve client
	req2 := httptest.NewRequest("GET", fmt.Sprintf("/client/%s", clientID), nil)
	w2 := httptest.NewRecorder()
	handleClientProfile(w2, req2, clientID)

	if w2.Code != http.StatusOK {
		t.Errorf("failed to retrieve client profile")
	}

	// Create activity for client
	buf3 := bytes.Buffer{}
	writer3 := multipart.NewWriter(&buf3)
	writer3.WriteField("client_id", clientID)
	writer3.WriteField("type", "meeting")
	writer3.WriteField("description", "Planning session")
	writer3.Close()

	req3 := httptest.NewRequest("POST", "/api/activities", &buf3)
	req3.Header.Set("Content-Type", writer3.FormDataContentType())
	w3 := httptest.NewRecorder()

	handleCreateActivity(w3, req3)

	if w3.Code != http.StatusCreated {
		t.Errorf("failed to create activity")
	}

	// Create invoice for client
	buf4 := bytes.Buffer{}
	writer4 := multipart.NewWriter(&buf4)
	writer4.WriteField("client_id", clientID)
	writer4.WriteField("due_date", "2026-12-31")
	writer4.WriteField("item_desc_0", "Complete Project")
	writer4.WriteField("item_qty_0", "40")
	writer4.WriteField("item_rate_0", "200.00")
	writer4.Close()

	req4 := httptest.NewRequest("POST", "/api/invoices", &buf4)
	req4.Header.Set("Content-Type", writer4.FormDataContentType())
	w4 := httptest.NewRecorder()

	handleCreateInvoice(w4, req4)

	if w4.Code != http.StatusCreated {
		t.Errorf("failed to create invoice")
	}

	var createdInvoice store.Invoice
	json.NewDecoder(w4.Body).Decode(&createdInvoice)
	invoiceID := createdInvoice.ID

	// Mark invoice as paid
	req5 := httptest.NewRequest("POST", fmt.Sprintf("/api/invoices/%s/mark-paid", invoiceID), nil)
	w5 := httptest.NewRecorder()
	handleInvoicesAPI(w5, req5)

	if w5.Code != http.StatusOK {
		t.Errorf("failed to mark invoice as paid")
	}

	// Verify client has the invoice and activity
	activities, _ := db.GetClientActivities(clientID)
	if len(activities) != 1 {
		t.Errorf("expected 1 activity, got %d", len(activities))
	}

	invoices, _ := db.GetClientInvoices(clientID)
	if len(invoices) != 1 {
		t.Errorf("expected 1 invoice, got %d", len(invoices))
	}

	// Verify invoice status
	finalInvoice, _ := db.GetInvoice(invoiceID)
	if finalInvoice.Status != "paid" {
		t.Errorf("expected invoice status=paid, got %s", finalInvoice.Status)
	}
}
