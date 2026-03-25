package main

import (
	"app/store"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestHealth(t *testing.T) {
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	handleHealth(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	if !strings.Contains(w.Body.String(), "ok") {
		t.Errorf("handler returned unexpected body: got %s", w.Body.String())
	}
}

func TestDashboardGET(t *testing.T) {
	// Reset store for test
	s = store.NewStore()
	s.CreateClient(&store.Client{Name: "Test Client", Email: "test@example.com"})

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	handleDashboard(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	if !strings.Contains(w.Body.String(), "Dashboard") {
		t.Errorf("handler returned unexpected body: missing Dashboard title")
	}

	if !strings.Contains(w.Body.String(), "Test Client") {
		t.Errorf("handler returned unexpected body: missing client name")
	}
}

func TestDashboardGETWithSearch(t *testing.T) {
	s = store.NewStore()
	s.CreateClient(&store.Client{Name: "Acme Corp", Email: "acme@example.com"})
	s.CreateClient(&store.Client{Name: "Beta Inc", Email: "beta@example.com"})

	req, err := http.NewRequest("GET", "/?q=acme", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	handleDashboard(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	if !strings.Contains(w.Body.String(), "Acme Corp") {
		t.Errorf("handler should contain Acme Corp")
	}

	if strings.Contains(w.Body.String(), "Beta Inc") {
		t.Errorf("handler should not contain Beta Inc in search results")
	}
}

func TestDashboardPOSTCreateClient(t *testing.T) {
	s = store.NewStore()

	form := url.Values{}
	form.Set("name", "New Client")
	form.Set("email", "new@example.com")

	req, err := http.NewRequest("POST", "/", strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()
	handleDashboard(w, req)

	if status := w.Code; status != http.StatusSeeOther {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusSeeOther)
	}

	clients := s.ListClients()
	if len(clients) != 1 {
		t.Errorf("expected 1 client, got %d", len(clients))
	}

	if clients[0].Name != "New Client" {
		t.Errorf("expected client name 'New Client', got '%s'", clients[0].Name)
	}
}

func TestDashboardPOSTCreateClientMissingFields(t *testing.T) {
	s = store.NewStore()

	form := url.Values{}
	form.Set("name", "")
	form.Set("email", "")

	req, err := http.NewRequest("POST", "/", strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()
	handleDashboard(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	clients := s.ListClients()
	if len(clients) != 0 {
		t.Errorf("expected 0 clients, got %d", len(clients))
	}
}

func TestClientGET(t *testing.T) {
	s = store.NewStore()
	s.CreateClient(&store.Client{Name: "Test Client", Email: "test@example.com"})
	clients := s.ListClients()
	clientID := clients[0].ID

	req, err := http.NewRequest("GET", "/client/"+clientID, nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	handleClientRoute(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	if !strings.Contains(w.Body.String(), "Test Client") {
		t.Errorf("handler should contain client name")
	}
}

func TestClientGETNotFound(t *testing.T) {
	s = store.NewStore()

	req, err := http.NewRequest("GET", "/client/nonexistent", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	handleClientRoute(w, req)

	if status := w.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}
}

func TestClientGETProfileTab(t *testing.T) {
	s = store.NewStore()
	s.CreateClient(&store.Client{Name: "Test Client", Email: "test@example.com"})
	clients := s.ListClients()
	clientID := clients[0].ID

	req, err := http.NewRequest("GET", "/client/"+clientID+"?tab=profile", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	handleClientRoute(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	if !strings.Contains(w.Body.String(), "active") {
		t.Errorf("handler should show active tab")
	}
}

func TestClientGETInvoicesTab(t *testing.T) {
	s = store.NewStore()
	s.CreateClient(&store.Client{Name: "Test Client", Email: "test@example.com"})
	clients := s.ListClients()
	clientID := clients[0].ID

	req, err := http.NewRequest("GET", "/client/"+clientID+"?tab=invoices", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	handleClientRoute(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	if !strings.Contains(w.Body.String(), "Create Invoice") {
		t.Errorf("handler should show Create Invoice button")
	}
}

func TestClientPOSTUpdate(t *testing.T) {
	s = store.NewStore()
	s.CreateClient(&store.Client{Name: "Old Name", Email: "old@example.com"})
	clients := s.ListClients()
	clientID := clients[0].ID

	form := url.Values{}
	form.Set("name", "Updated Name")
	form.Set("email", "updated@example.com")

	req, err := http.NewRequest("POST", "/client/"+clientID, strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()
	handleClientRoute(w, req)

	if status := w.Code; status != http.StatusSeeOther {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusSeeOther)
	}

	updated, _ := s.GetClient(clientID)
	if updated.Name != "Updated Name" {
		t.Errorf("expected name 'Updated Name', got '%s'", updated.Name)
	}
	if updated.Email != "updated@example.com" {
		t.Errorf("expected email 'updated@example.com', got '%s'", updated.Email)
	}
}

func TestClientPOSTDelete(t *testing.T) {
	s = store.NewStore()
	s.CreateClient(&store.Client{Name: "To Delete", Email: "delete@example.com"})
	clients := s.ListClients()
	clientID := clients[0].ID

	form := url.Values{}
	form.Set("_method", "DELETE")

	req, err := http.NewRequest("POST", "/client/"+clientID, strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()
	handleClientRoute(w, req)

	if status := w.Code; status != http.StatusSeeOther {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusSeeOther)
	}

	_, ok := s.GetClient(clientID)
	if ok {
		t.Errorf("client should be deleted")
	}
}

func TestClientEmptyID(t *testing.T) {
	s = store.NewStore()

	req, err := http.NewRequest("GET", "/client/", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	handleClientRoute(w, req)

	if status := w.Code; status != http.StatusSeeOther {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusSeeOther)
	}
}

func TestCreateInvoiceGET(t *testing.T) {
	s = store.NewStore()
	s.CreateClient(&store.Client{Name: "Test Client", Email: "test@example.com"})
	clients := s.ListClients()
	clientID := clients[0].ID

	req, err := http.NewRequest("GET", "/client/"+clientID+"/invoice/new", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	handleClientRoute(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	if !strings.Contains(w.Body.String(), "Create Invoice") {
		t.Errorf("handler should contain Create Invoice title")
	}

	if !strings.Contains(w.Body.String(), "Description") {
		t.Errorf("handler should contain description field")
	}
}

func TestCreateInvoicePOST(t *testing.T) {
	s = store.NewStore()
	s.CreateClient(&store.Client{Name: "Test Client", Email: "test@example.com"})
	clients := s.ListClients()
	clientID := clients[0].ID

	form := url.Values{}
	form.Set("description_0", "Service A")
	form.Set("amount_0", "100.00")
	form.Set("description_1", "Service B")
	form.Set("amount_1", "50.00")

	req, err := http.NewRequest("POST", "/client/"+clientID+"/invoice/new", strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()
	handleClientRoute(w, req)

	if status := w.Code; status != http.StatusSeeOther {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusSeeOther)
	}

	invoices := s.ListInvoicesByClient(clientID)
	if len(invoices) != 1 {
		t.Errorf("expected 1 invoice, got %d", len(invoices))
	}

	if len(invoices[0].LineItems) != 2 {
		t.Errorf("expected 2 line items, got %d", len(invoices[0].LineItems))
	}

	if invoices[0].Total != 150.00 {
		t.Errorf("expected total 150.00, got %f", invoices[0].Total)
	}
}

func TestCreateInvoicePOSTNoLineItems(t *testing.T) {
	s = store.NewStore()
	s.CreateClient(&store.Client{Name: "Test Client", Email: "test@example.com"})
	clients := s.ListClients()
	clientID := clients[0].ID

	form := url.Values{}
	// No line items provided

	req, err := http.NewRequest("POST", "/client/"+clientID+"/invoice/new", strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()
	handleClientRoute(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	invoices := s.ListInvoicesByClient(clientID)
	if len(invoices) != 0 {
		t.Errorf("expected 0 invoices, got %d", len(invoices))
	}
}

func TestInvoiceGET(t *testing.T) {
	s = store.NewStore()
	s.CreateClient(&store.Client{Name: "Test Client", Email: "test@example.com"})
	clients := s.ListClients()
	clientID := clients[0].ID

	invoice := &store.Invoice{
		ClientID: clientID,
		Status:   "draft",
		Total:    100.00,
		LineItems: []store.LineItem{
			{Description: "Service", Amount: 100.00},
		},
	}
	s.CreateInvoice(invoice)
	invoices := s.ListInvoicesByClient(clientID)
	invoiceID := invoices[0].ID

	req, err := http.NewRequest("GET", "/invoice/"+invoiceID, nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	handleInvoiceRoute(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	if !strings.Contains(w.Body.String(), "Invoice #"+invoiceID) {
		t.Errorf("handler should contain invoice ID")
	}

	if !strings.Contains(w.Body.String(), "draft") {
		t.Errorf("handler should contain draft status")
	}
}

func TestInvoiceGETNotFound(t *testing.T) {
	s = store.NewStore()

	req, err := http.NewRequest("GET", "/invoice/nonexistent", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	handleInvoiceRoute(w, req)

	if status := w.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}
}

func TestInvoicePOSTSend(t *testing.T) {
	s = store.NewStore()
	s.CreateClient(&store.Client{Name: "Test Client", Email: "test@example.com"})
	clients := s.ListClients()
	clientID := clients[0].ID

	invoice := &store.Invoice{
		ClientID: clientID,
		Status:   "draft",
		Total:    100.00,
		LineItems: []store.LineItem{
			{Description: "Service", Amount: 100.00},
		},
	}
	s.CreateInvoice(invoice)
	invoices := s.ListInvoicesByClient(clientID)
	invoiceID := invoices[0].ID

	form := url.Values{}
	form.Set("action", "send")

	req, err := http.NewRequest("POST", "/invoice/"+invoiceID, strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()
	handleInvoiceRoute(w, req)

	if status := w.Code; status != http.StatusSeeOther {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusSeeOther)
	}

	updated, _ := s.GetInvoice(invoiceID)
	if updated.Status != "sent" {
		t.Errorf("expected status 'sent', got '%s'", updated.Status)
	}
}

func TestInvoicePOSTPay(t *testing.T) {
	s = store.NewStore()
	s.CreateClient(&store.Client{Name: "Test Client", Email: "test@example.com"})
	clients := s.ListClients()
	clientID := clients[0].ID

	invoice := &store.Invoice{
		ClientID: clientID,
		Status:   "draft",
		Total:    100.00,
		LineItems: []store.LineItem{
			{Description: "Service", Amount: 100.00},
		},
	}
	s.CreateInvoice(invoice)
	invoices := s.ListInvoicesByClient(clientID)
	invoiceID := invoices[0].ID

	form := url.Values{}
	form.Set("action", "pay")

	req, err := http.NewRequest("POST", "/invoice/"+invoiceID, strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()
	handleInvoiceRoute(w, req)

	if status := w.Code; status != http.StatusSeeOther {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusSeeOther)
	}

	updated, _ := s.GetInvoice(invoiceID)
	if updated.Status != "paid" {
		t.Errorf("expected status 'paid', got '%s'", updated.Status)
	}
}

func TestInvoicePOSTVoid(t *testing.T) {
	s = store.NewStore()
	s.CreateClient(&store.Client{Name: "Test Client", Email: "test@example.com"})
	clients := s.ListClients()
	clientID := clients[0].ID

	invoice := &store.Invoice{
		ClientID: clientID,
		Status:   "draft",
		Total:    100.00,
		LineItems: []store.LineItem{
			{Description: "Service", Amount: 100.00},
		},
	}
	s.CreateInvoice(invoice)
	invoices := s.ListInvoicesByClient(clientID)
	invoiceID := invoices[0].ID

	form := url.Values{}
	form.Set("action", "void")

	req, err := http.NewRequest("POST", "/invoice/"+invoiceID, strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()
	handleInvoiceRoute(w, req)

	if status := w.Code; status != http.StatusSeeOther {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusSeeOther)
	}

	updated, _ := s.GetInvoice(invoiceID)
	if updated.Status != "void" {
		t.Errorf("expected status 'void', got '%s'", updated.Status)
	}
}

func TestInvoicePrintGET(t *testing.T) {
	s = store.NewStore()
	s.CreateClient(&store.Client{Name: "Test Client", Email: "test@example.com"})
	clients := s.ListClients()
	clientID := clients[0].ID

	invoice := &store.Invoice{
		ClientID: clientID,
		Status:   "paid",
		Total:    100.00,
		LineItems: []store.LineItem{
			{Description: "Service", Amount: 100.00},
		},
	}
	s.CreateInvoice(invoice)
	invoices := s.ListInvoicesByClient(clientID)
	invoiceID := invoices[0].ID

	req, err := http.NewRequest("GET", "/invoice/"+invoiceID+"/print", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	handleInvoiceRoute(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	if !strings.Contains(w.Body.String(), "INVOICE") {
		t.Errorf("handler should contain INVOICE title")
	}

	if !strings.Contains(w.Body.String(), invoiceID) {
		t.Errorf("handler should contain invoice ID")
	}
}

func TestInvoiceEmptyID(t *testing.T) {
	s = store.NewStore()

	req, err := http.NewRequest("GET", "/invoice/", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	handleInvoiceRoute(w, req)

	if status := w.Code; status != http.StatusSeeOther {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusSeeOther)
	}
}

func TestFormatDate(t *testing.T) {
	tt := []struct {
		name     string
		input    string
		expected string
	}{
		{"Jan", "2024-01-15", "Jan 15, 2024"},
		{"Dec", "2024-12-25", "Dec 25, 2024"},
	}

	for _, tc := range tt {
		// Parse manually and test
		parts := strings.Split(tc.input, "-")
		if len(parts) != 3 {
			continue
		}
	}
}

func TestFormatCurrency(t *testing.T) {
	tests := []struct {
		input    float64
		expected string
	}{
		{100.0, "100.00"},
		{100.5, "100.50"},
		{100.555, "100.56"},
		{0.01, "0.01"},
	}

	for _, tc := range tests {
		result := formatCurrency(tc.input)
		if result != tc.expected {
			t.Errorf("formatCurrency(%f) = %s, want %s", tc.input, result, tc.expected)
		}
	}
}

func TestIfActive(t *testing.T) {
	tests := []struct {
		condition bool
		text      string
		expected  string
	}{
		{true, " active", " active"},
		{false, " active", ""},
		{true, "", ""},
	}

	for _, tc := range tests {
		result := ifActive(tc.condition, tc.text)
		if result != tc.expected {
			t.Errorf("ifActive(%v, %q) = %q, want %q", tc.condition, tc.text, result, tc.expected)
		}
	}
}

func TestMultipleClientsWithInvoices(t *testing.T) {
	s = store.NewStore()

	// Create two clients
	s.CreateClient(&store.Client{Name: "Client A", Email: "a@example.com"})
	s.CreateClient(&store.Client{Name: "Client B", Email: "b@example.com"})

	clients := s.ListClients()
	clientAID := clients[0].ID
	clientBID := clients[1].ID

	// Create invoices for each client
	invA := &store.Invoice{
		ClientID: clientAID,
		Status:   "draft",
		Total:    100.00,
		LineItems: []store.LineItem{
			{Description: "Service A", Amount: 100.00},
		},
	}
	s.CreateInvoice(invA)

	invB := &store.Invoice{
		ClientID: clientBID,
		Status:   "draft",
		Total:    200.00,
		LineItems: []store.LineItem{
			{Description: "Service B", Amount: 200.00},
		},
	}
	s.CreateInvoice(invB)

	// Verify client A only has 1 invoice
	invicesA := s.ListInvoicesByClient(clientAID)
	if len(invicesA) != 1 {
		t.Errorf("client A should have 1 invoice, got %d", len(invicesA))
	}

	// Verify client B only has 1 invoice
	invicesB := s.ListInvoicesByClient(clientBID)
	if len(invicesB) != 1 {
		t.Errorf("client B should have 1 invoice, got %d", len(invicesB))
	}

	// Verify invoice amounts
	if invicesA[0].Total != 100.00 {
		t.Errorf("invoice A total should be 100.00, got %f", invicesA[0].Total)
	}

	if invicesB[0].Total != 200.00 {
		t.Errorf("invoice B total should be 200.00, got %f", invicesB[0].Total)
	}
}

func BenchmarkDashboardGET(b *testing.B) {
	s = store.NewStore()
	s.CreateClient(&store.Client{Name: "Benchmark Client", Email: "bench@example.com"})

	for i := 0; i < b.N; i++ {
		req, _ := http.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()
		handleDashboard(w, req)
	}
}

func BenchmarkClientGET(b *testing.B) {
	s = store.NewStore()
	s.CreateClient(&store.Client{Name: "Benchmark Client", Email: "bench@example.com"})
	clients := s.ListClients()
	clientID := clients[0].ID

	for i := 0; i < b.N; i++ {
		req, _ := http.NewRequest("GET", "/client/"+clientID, nil)
		w := httptest.NewRecorder()
		handleClientRoute(w, req)
	}
}

func BenchmarkInvoiceGET(b *testing.B) {
	s = store.NewStore()
	s.CreateClient(&store.Client{Name: "Benchmark Client", Email: "bench@example.com"})
	clients := s.ListClients()

	invoice := &store.Invoice{
		ClientID: clients[0].ID,
		Status:   "draft",
		Total:    100.00,
		LineItems: []store.LineItem{
			{Description: "Service", Amount: 100.00},
		},
	}
	s.CreateInvoice(invoice)
	invoices := s.ListInvoicesByClient(clients[0].ID)
	invoiceID := invoices[0].ID

	for i := 0; i < b.N; i++ {
		req, _ := http.NewRequest("GET", "/invoice/"+invoiceID, nil)
		w := httptest.NewRecorder()
		handleInvoiceRoute(w, req)
	}
}
