package store

import (
	"testing"
	"time"
)

func TestGetClient(t *testing.T) {
	store := NewStore()

	// Test getting non-existent client
	client, err := store.GetClient("non-existent")
	if err != nil {
		t.Errorf("GetClient returned error: %v", err)
	}
	if client != nil {
		t.Errorf("GetClient should return nil for non-existent client, got: %v", client)
	}

	// Test getting existing client
	clientID := "test-client"
	clientData := &Client{
		ID:        clientID,
		Name:      "Test Client",
		Email:     "test@example.com",
		Phone:     "123-456-7890",
		Address:   Address{Street: "123 Main St", City: "Test City", State: "TS", ZipCode: "12345", Country: "USA"},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	store.CreateClient(clientData)

	result, err := store.GetClient(clientID)
	if err != nil {
		t.Errorf("GetClient returned error: %v", err)
	}
	if result == nil {
		t.Errorf("GetClient should return client, got nil")
	}
	if result.ID != clientID {
		t.Errorf("GetClient returned wrong client ID: got %s, want %s", result.ID, clientID)
	}
}

func TestGetClients(t *testing.T) {
	store := NewStore()

	// Test getting empty list
	clients, err := store.GetClients()
	if err != nil {
		t.Errorf("GetClients returned error: %v", err)
	}
	if len(clients) != 0 {
		t.Errorf("GetClients should return empty slice, got %d items", len(clients))
	}

	// Test getting multiple clients
	client1 := &Client{
		ID:        "client1",
		Name:      "Client 1",
		Email:     "client1@example.com",
		Phone:     "111-111-1111",
		Address:   Address{Street: "111 First St", City: "City 1", State: "S1", ZipCode: "11111", Country: "USA"},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	client2 := &Client{
		ID:        "client2",
		Name:      "Client 2",
		Email:     "client2@example.com",
		Phone:     "222-222-2222",
		Address:   Address{Street: "222 Second St", City: "City 2", State: "S2", ZipCode: "22222", Country: "USA"},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	store.CreateClient(client1)
	store.CreateClient(client2)

	clients, err = store.GetClients()
	if err != nil {
		t.Errorf("GetClients returned error: %v", err)
	}
	if len(clients) != 2 {
		t.Errorf("GetClients should return 2 clients, got %d", len(clients))
	}

	// Check that both clients are returned
	clientIDs := make(map[string]bool)
	for _, client := range clients {
		clientIDs[client.ID] = true
	}

	if !clientIDs["client1"] {
		t.Errorf("GetClients should return client1")
	}
	if !clientIDs["client2"] {
		t.Errorf("GetClients should return client2")
	}
}

func TestCreateClient(t *testing.T) {
	store := NewStore()

	client := &Client{
		ID:        "test-client",
		Name:      "Test Client",
		Email:     "test@example.com",
		Phone:     "123-456-7890",
		Address:   Address{Street: "123 Main St", City: "Test City", State: "TS", ZipCode: "12345", Country: "USA"},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := store.CreateClient(client)
	if err != nil {
		t.Errorf("CreateClient returned error: %v", err)
	}

	// Verify client was created
	result, err := store.GetClient("test-client")
	if err != nil {
		t.Errorf("GetClient returned error: %v", err)
	}
	if result == nil {
		t.Errorf("Client was not created")
	}

	// Verify fields match
	if result.Name != "Test Client" {
		t.Errorf("Client name mismatch: got %s, want Test Client", result.Name)
	}
	if result.Email != "test@example.com" {
		t.Errorf("Client email mismatch: got %s, want test@example.com", result.Email)
	}
}

func TestUpdateClient(t *testing.T) {
	store := NewStore()

	// Create initial client
	initialClient := &Client{
		ID:        "test-client",
		Name:      "Original Name",
		Email:     "original@example.com",
		Phone:     "111-111-1111",
		Address:   Address{Street: "111 Original St", City: "Original City", State: "OR", ZipCode: "11111", Country: "USA"},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	store.CreateClient(initialClient)

	// Update client
	updatedClient := &Client{
		Name:      "Updated Name",
		Email:     "updated@example.com",
		Phone:     "222-222-2222",
		Address:   Address{Street: "222 Updated St", City: "Updated City", State: "UP", ZipCode: "22222", Country: "USA"},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := store.UpdateClient("test-client", updatedClient)
	if err != nil {
		t.Errorf("UpdateClient returned error: %v", err)
	}

	// Verify update
	result, err := store.GetClient("test-client")
	if err != nil {
		t.Errorf("GetClient returned error: %v", err)
	}
	if result == nil {
		t.Errorf("Client was not found after update")
	}

	if result.Name != "Updated Name" {
		t.Errorf("Client name not updated: got %s, want Updated Name", result.Name)
	}
	if result.Email != "updated@example.com" {
		t.Errorf("Client email not updated: got %s, want updated@example.com", result.Email)
	}
	if result.Phone != "222-222-2222" {
		t.Errorf("Client phone not updated: got %s, want 222-222-2222", result.Phone)
	}

	// Verify that ID is preserved
	if result.ID != "test-client" {
		t.Errorf("Client ID was changed: got %s, want test-client", result.ID)
	}

	// Verify that UpdatedAt was updated
	if result.UpdatedAt.IsZero() {
		t.Errorf("Client UpdatedAt was not set")
	}
}

func TestDeleteClient(t *testing.T) {
	store := NewStore()

	// Create client
	client := &Client{
		ID:        "test-client",
		Name:      "Test Client",
		Email:     "test@example.com",
		Phone:     "123-456-7890",
		Address:   Address{Street: "123 Main St", City: "Test City", State: "TS", ZipCode: "12345", Country: "USA"},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	store.CreateClient(client)

	// Verify client exists
	result, err := store.GetClient("test-client")
	if err != nil {
		t.Errorf("GetClient returned error: %v", err)
	}
	if result == nil {
		t.Errorf("Client was not created")
	}

	// Delete client
	err = store.DeleteClient("test-client")
	if err != nil {
		t.Errorf("DeleteClient returned error: %v", err)
	}

	// Verify client was deleted
	result, err = store.GetClient("test-client")
	if err != nil {
		t.Errorf("GetClient returned error: %v", err)
	}
	if result != nil {
		t.Errorf("Client was not deleted")
	}
}

func TestGetInvoice(t *testing.T) {
	store := NewStore()

	// Test getting non-existent invoice
	invoice, err := store.GetInvoice("non-existent")
	if err != nil {
		t.Errorf("GetInvoice returned error: %v", err)
	}
	if invoice != nil {
		t.Errorf("GetInvoice should return nil for non-existent invoice, got: %v", invoice)
	}

	// Test getting existing invoice
	invoiceID := "test-invoice"
	invoiceData := &Invoice{
		ID:          invoiceID,
		ClientID:    "test-client",
		Items:       []Item{{Description: "Test Item", Quantity: 1, Price: 100.0, Total: 100.0}},
		TotalAmount: 100.0,
		IssueDate:   time.Now(),
		DueDate:     time.Now().Add(30 * 24 * time.Hour),
		Status:      "draft",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	store.CreateInvoice(invoiceData)

	result, err := store.GetInvoice(invoiceID)
	if err != nil {
		t.Errorf("GetInvoice returned error: %v", err)
	}
	if result == nil {
		t.Errorf("GetInvoice should return invoice, got nil")
	}
	if result.ID != invoiceID {
		t.Errorf("GetInvoice returned wrong invoice ID: got %s, want %s", result.ID, invoiceID)
	}
}

func TestGetInvoices(t *testing.T) {
	store := NewStore()

	// Test getting empty list
	invoices, err := store.GetInvoices()
	if err != nil {
		t.Errorf("GetInvoices returned error: %v", err)
	}
	if len(invoices) != 0 {
		t.Errorf("GetInvoices should return empty slice, got %d items", len(invoices))
	}

	// Test getting multiple invoices
	invoice1 := &Invoice{
		ID:          "invoice1",
		ClientID:    "client1",
		Items:       []Item{{Description: "Item 1", Quantity: 1, Price: 100.0, Total: 100.0}},
		TotalAmount: 100.0,
		IssueDate:   time.Now(),
		DueDate:     time.Now().Add(30 * 24 * time.Hour),
		Status:      "draft",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	invoice2 := &Invoice{
		ID:          "invoice2",
		ClientID:    "client2",
		Items:       []Item{{Description: "Item 2", Quantity: 2, Price: 50.0, Total: 100.0}},
		TotalAmount: 100.0,
		IssueDate:   time.Now(),
		DueDate:     time.Now().Add(30 * 24 * time.Hour),
		Status:      "sent",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	store.CreateInvoice(invoice1)
	store.CreateInvoice(invoice2)

	invoices, err = store.GetInvoices()
	if err != nil {
		t.Errorf("GetInvoices returned error: %v", err)
	}
	if len(invoices) != 2 {
		t.Errorf("GetInvoices should return 2 invoices, got %d", len(invoices))
	}

	// Check that both invoices are returned
	invoiceIDs := make(map[string]bool)
	for _, invoice := range invoices {
		invoiceIDs[invoice.ID] = true
	}

	if !invoiceIDs["invoice1"] {
		t.Errorf("GetInvoices should return invoice1")
	}
	if !invoiceIDs["invoice2"] {
		t.Errorf("GetInvoices should return invoice2")
	}
}

func TestCreateInvoice(t *testing.T) {
	store := NewStore()

	invoice := &Invoice{
		ID:          "test-invoice",
		ClientID:    "test-client",
		Items:       []Item{{Description: "Test Item", Quantity: 1, Price: 100.0, Total: 100.0}},
		TotalAmount: 100.0,
		IssueDate:   time.Now(),
		DueDate:     time.Now().Add(30 * 24 * time.Hour),
		Status:      "draft",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := store.CreateInvoice(invoice)
	if err != nil {
		t.Errorf("CreateInvoice returned error: %v", err)
	}

	// Verify invoice was created
	result, err := store.GetInvoice("test-invoice")
	if err != nil {
		t.Errorf("GetInvoice returned error: %v", err)
	}
	if result == nil {
		t.Errorf("Invoice was not created")
	}

	// Verify fields match
	if result.ClientID != "test-client" {
		t.Errorf("Invoice client ID mismatch: got %s, want test-client", result.ClientID)
	}
	if result.TotalAmount != 100.0 {
		t.Errorf("Invoice total amount mismatch: got %f, want 100.0", result.TotalAmount)
	}
	if result.Status != "draft" {
		t.Errorf("Invoice status mismatch: got %s, want draft", result.Status)
	}
}

func TestUpdateInvoice(t *testing.T) {
	store := NewStore()

	// Create initial invoice
	initialInvoice := &Invoice{
		ID:          "test-invoice",
		ClientID:    "test-client",
		Items:       []Item{{Description: "Original Item", Quantity: 1, Price: 100.0, Total: 100.0}},
		TotalAmount: 100.0,
		IssueDate:   time.Now(),
		DueDate:     time.Now().Add(30 * 24 * time.Hour),
		Status:      "draft",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	store.CreateInvoice(initialInvoice)

	// Update invoice
	updatedInvoice := &Invoice{
		ClientID:    "updated-client",
		Items:       []Item{{Description: "Updated Item", Quantity: 2, Price: 50.0, Total: 100.0}},
		TotalAmount: 100.0,
		IssueDate:   time.Now(),
		DueDate:     time.Now().Add(60 * 24 * time.Hour),
		Status:      "sent",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := store.UpdateInvoice("test-invoice", updatedInvoice)
	if err != nil {
		t.Errorf("UpdateInvoice returned error: %v", err)
	}

	// Verify update
	result, err := store.GetInvoice("test-invoice")
	if err != nil {
		t.Errorf("GetInvoice returned error: %v", err)
	}
	if result == nil {
		t.Errorf("Invoice was not found after update")
	}

	if result.ClientID != "updated-client" {
		t.Errorf("Invoice client ID not updated: got %s, want updated-client", result.ClientID)
	}
	if result.Status != "sent" {
		t.Errorf("Invoice status not updated: got %s, want sent", result.Status)
	}

	// Verify that ID is preserved
	if result.ID != "test-invoice" {
		t.Errorf("Invoice ID was changed: got %s, want test-invoice", result.ID)
	}

	// Verify that UpdatedAt was updated
	if result.UpdatedAt.IsZero() {
		t.Errorf("Invoice UpdatedAt was not set")
	}
}

func TestDeleteInvoice(t *testing.T) {
	store := NewStore()

	// Create invoice
	invoice := &Invoice{
		ID:          "test-invoice",
		ClientID:    "test-client",
		Items:       []Item{{Description: "Test Item", Quantity: 1, Price: 100.0, Total: 100.0}},
		TotalAmount: 100.0,
		IssueDate:   time.Now(),
		DueDate:     time.Now().Add(30 * 24 * time.Hour),
		Status:      "draft",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	store.CreateInvoice(invoice)

	// Verify invoice exists
	result, err := store.GetInvoice("test-invoice")
	if err != nil {
		t.Errorf("GetInvoice returned error: %v", err)
	}
	if result == nil {
		t.Errorf("Invoice was not created")
	}

	// Delete invoice
	err = store.DeleteInvoice("test-invoice")
	if err != nil {
		t.Errorf("DeleteInvoice returned error: %v", err)
	}

	// Verify invoice was deleted
	result, err = store.GetInvoice("test-invoice")
	if err != nil {
		t.Errorf("GetInvoice returned error: %v", err)
	}
	if result != nil {
		t.Errorf("Invoice was not deleted")
	}
}

func TestGetComments(t *testing.T) {
	store := NewStore()

	// Test getting comments for non-existent client
	comments, err := store.GetComments("non-existent-client")
	if err != nil {
		t.Errorf("GetComments returned error: %v", err)
	}
	if len(comments) != 0 {
		t.Errorf("GetComments should return empty slice, got %d items", len(comments))
	}

	// Test getting comments for existing client
	comment1 := &Comment{
		ID:        "comment1",
		ClientID:  "test-client",
		Content:   "First comment",
		CreatedAt: time.Now(),
	}

	comment2 := &Comment{
		ID:        "comment2",
		ClientID:  "test-client",
		Content:   "Second comment",
		CreatedAt: time.Now(),
	}

	store.CreateComment(comment1)
	store.CreateComment(comment2)

	comments, err = store.GetComments("test-client")
	if err != nil {
		t.Errorf("GetComments returned error: %v", err)
	}
	if len(comments) != 2 {
		t.Errorf("GetComments should return 2 comments, got %d", len(comments))
	}

	// Check that both comments are returned
	commentIDs := make(map[string]bool)
	for _, comment := range comments {
		commentIDs[comment.ID] = true
	}

	if !commentIDs["comment1"] {
		t.Errorf("GetComments should return comment1")
	}
	if !commentIDs["comment2"] {
		t.Errorf("GetComments should return comment2")
	}
}

func TestCreateComment(t *testing.T) {
	store := NewStore()

	comment := &Comment{
		ID:        "test-comment",
		ClientID:  "test-client",
		Content:   "Test comment content",
		CreatedAt: time.Now(),
	}

	err := store.CreateComment(comment)
	if err != nil {
		t.Errorf("CreateComment returned error: %v", err)
	}

	// Verify comment was created
	result, err := store.GetComments("test-client")
	if err != nil {
		t.Errorf("GetComments returned error: %v", err)
	}
	if len(result) != 1 {
		t.Errorf("Should have 1 comment, got %d", len(result))
	}

	if result[0].Content != "Test comment content" {
		t.Errorf("Comment content mismatch: got %s, want Test comment content", result[0].Content)
	}
}

func TestDeleteComment(t *testing.T) {
	store := NewStore()

	// Create comment
	comment := &Comment{
		ID:        "test-comment",
		ClientID:  "test-client",
		Content:   "Test comment content",
		CreatedAt: time.Now(),
	}
	store.CreateComment(comment)

	// Verify comment exists
	comments, err := store.GetComments("test-client")
	if err != nil {
		t.Errorf("GetComments returned error: %v", err)
	}
	if len(comments) != 1 {
		t.Errorf("Should have 1 comment, got %d", len(comments))
	}

	// Delete comment
	err = store.DeleteComment("test-comment")
	if err != nil {
		t.Errorf("DeleteComment returned error: %v", err)
	}

	// Verify comment was deleted
	comments, err = store.GetComments("test-client")
	if err != nil {
		t.Errorf("GetComments returned error: %v", err)
	}
	if len(comments) != 0 {
		t.Errorf("Should have 0 comments after deletion, got %d", len(comments))
	}
}

func TestGetHistory(t *testing.T) {
	store := NewStore()

	// Test getting history for non-existent client
	history, err := store.GetHistory("non-existent-client")
	if err != nil {
		t.Errorf("GetHistory returned error: %v", err)
	}
	if len(history) != 0 {
		t.Errorf("GetHistory should return empty slice, got %d items", len(history))
	}

	// Test getting history for existing client
	entry1 := &HistoryEntry{
		ID:        "entry1",
		ClientID:  "test-client",
		Note:      "First history entry",
		CreatedAt: time.Now(),
	}

	entry2 := &HistoryEntry{
		ID:        "entry2",
		ClientID:  "test-client",
		Note:      "Second history entry",
		CreatedAt: time.Now(),
	}

	store.CreateHistory(entry1)
	store.CreateHistory(entry2)

	history, err = store.GetHistory("test-client")
	if err != nil {
		t.Errorf("GetHistory returned error: %v", err)
	}
	if len(history) != 2 {
		t.Errorf("GetHistory should return 2 entries, got %d", len(history))
	}

	// Check that both entries are returned
	entryIDs := make(map[string]bool)
	for _, entry := range history {
		entryIDs[entry.ID] = true
	}

	if !entryIDs["entry1"] {
		t.Errorf("GetHistory should return entry1")
	}
	if !entryIDs["entry2"] {
		t.Errorf("GetHistory should return entry2")
	}
}

func TestCreateHistory(t *testing.T) {
	store := NewStore()

	entry := &HistoryEntry{
		ID:        "test-entry",
		ClientID:  "test-client",
		Note:      "Test history entry",
		CreatedAt: time.Now(),
	}

	err := store.CreateHistory(entry)
	if err != nil {
		t.Errorf("CreateHistory returned error: %v", err)
	}

	// Verify entry was created
	result, err := store.GetHistory("test-client")
	if err != nil {
		t.Errorf("GetHistory returned error: %v", err)
	}
	if len(result) != 1 {
		t.Errorf("Should have 1 history entry, got %d", len(result))
	}

	if result[0].Note != "Test history entry" {
		t.Errorf("History entry note mismatch: got %s, want Test history entry", result[0].Note)
	}
}

func TestSaveAndLoad(t *testing.T) {
	store := NewStore()

	// Create some test data
	client := &Client{
		ID:        "test-client",
		Name:      "Test Client",
		Email:     "test@example.com",
		Phone:     "123-456-7890",
		Address:   Address{Street: "123 Main St", City: "Test City", State: "TS", ZipCode: "12345", Country: "USA"},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	invoice := &Invoice{
		ID:          "test-invoice",
		ClientID:    "test-client",
		Items:       []Item{{Description: "Test Item", Quantity: 1, Price: 100.0, Total: 100.0}},
		TotalAmount: 100.0,
		IssueDate:   time.Now(),
		DueDate:     time.Now().Add(30 * 24 * time.Hour),
		Status:      "draft",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	comment := &Comment{
		ID:        "test-comment",
		ClientID:  "test-client",
		Content:   "Test comment",
		CreatedAt: time.Now(),
	}

	entry := &HistoryEntry{
		ID:        "test-entry",
		ClientID:  "test-client",
		Note:      "Test history entry",
		CreatedAt: time.Now(),
	}

	store.CreateClient(client)
	store.CreateInvoice(invoice)
	store.CreateComment(comment)
	store.CreateHistory(entry)

	// Test Save (should not error)
	err := store.Save()
	if err != nil {
		t.Errorf("Save returned error: %v", err)
	}

	// Test Load (should not error)
	err = store.Load()
	if err != nil {
		t.Errorf("Load returned error: %v", err)
	}

	// Verify data is still there
	resultClient, err := store.GetClient("test-client")
	if err != nil {
		t.Errorf("GetClient returned error: %v", err)
	}
	if resultClient == nil {
		t.Errorf("Client was not found after save/load")
	}

	resultInvoice, err := store.GetInvoice("test-invoice")
	if err != nil {
		t.Errorf("GetInvoice returned error: %v", err)
	}
	if resultInvoice == nil {
		t.Errorf("Invoice was not found after save/load")
	}

	resultComments, err := store.GetComments("test-client")
	if err != nil {
		t.Errorf("GetComments returned error: %v", err)
	}
	if len(resultComments) != 1 {
		t.Errorf("Should have 1 comment after save/load, got %d", len(resultComments))
	}

	resultHistory, err := store.GetHistory("test-client")
	if err != nil {
		t.Errorf("GetHistory returned error: %v", err)
	}
	if len(resultHistory) != 1 {
		t.Errorf("Should have 1 history entry after save/load, got %d", len(resultHistory))
	}
}
