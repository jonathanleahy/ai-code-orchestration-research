package main

import (
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"app/store"
)

var db *store.Store

func main() {
	db = store.NewStore()

	// Routes
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/", dashboardHandler)
	http.HandleFunc("/clients", clientsHandler)
	http.HandleFunc("/clients/new", createClientHandler)
	http.HandleFunc("/clients/api/create", apiCreateClientHandler)
	http.HandleFunc("/clients/api/list", apiListClientsHandler)
	http.HandleFunc("/clients/api/update", apiUpdateClientHandler)
	http.HandleFunc("/clients/api/delete", apiDeleteClientHandler)
	http.HandleFunc("/clients/api/balance", clientBalanceHandler)
	http.HandleFunc("/clients/profile", clientProfileHandler)
	http.HandleFunc("/activities", activitiesHandler)
	http.HandleFunc("/activities/api/create", apiCreateActivityHandler)
	http.HandleFunc("/activities/api/delete", apiDeleteActivityHandler)
	http.HandleFunc("/invoices", invoicesHandler)
	http.HandleFunc("/invoices/api/create", apiCreateInvoiceHandler)
	http.HandleFunc("/invoices/api/update-status", apiUpdateInvoiceStatusHandler)
	http.HandleFunc("/invoices/api/calculate-totals", apiCalculateTotalsHandler)
	http.HandleFunc("/invoices/api/delete", apiDeleteInvoiceHandler)
	http.HandleFunc("/invoices/print", invoicePrintHandler)

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// healthHandler returns a health check response
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// dashboardHandler shows the main dashboard with client list
func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	search := r.URL.Query().Get("search")
	clients := db.ListClients()

	// Filter by search
	var filtered []store.Client
	for _, c := range clients {
		if search == "" || strings.Contains(strings.ToLower(c.Name), strings.ToLower(search)) ||
			strings.Contains(strings.ToLower(c.Email), strings.ToLower(search)) {
			filtered = append(filtered, c)
		}
	}

	// Sort by name
	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].Name < filtered[j].Name
	})

	htmlContent := dashboardHTML(filtered, search)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'self' 'unsafe-inline'; script-src 'self'")
	fmt.Fprint(w, htmlContent)
}

// clientProfileHandler shows client details with tabs
func clientProfileHandler(w http.ResponseWriter, r *http.Request) {
	clientID := r.URL.Query().Get("id")
	if clientID == "" {
		http.Error(w, "Client ID required", http.StatusBadRequest)
		return
	}

	client, err := db.GetClient(clientID)
	if err != nil {
		http.Error(w, "Client not found", http.StatusNotFound)
		return
	}

	tab := r.URL.Query().Get("tab")
	if tab == "" {
		tab = "details"
	}

	// Get related data
	invoices := db.ListInvoices()
	var clientInvoices []store.Invoice
	for _, inv := range invoices {
		if inv.ClientID == clientID {
			clientInvoices = append(clientInvoices, inv)
		}
	}

	activities := db.ListActivities()
	var clientActivities []store.Activity
	for _, act := range activities {
		if act.ClientID == clientID {
			clientActivities = append(clientActivities, act)
		}
	}

	// Sort activities by date descending
	sort.Slice(clientActivities, func(i, j int) bool {
		return clientActivities[i].Date.After(clientActivities[j].Date)
	})

	// Sort invoices by date descending
	sort.Slice(clientInvoices, func(i, j int) bool {
		return clientInvoices[i].InvoiceDate.After(clientInvoices[j].InvoiceDate)
	})

	htmlContent := clientProfileHTML(client, tab, clientInvoices, clientActivities)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'self' 'unsafe-inline'; script-src 'self'")
	fmt.Fprint(w, htmlContent)
}

// clientsHandler handles client list operations
func clientsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		clients := db.ListClients()
		sort.Slice(clients, func(i, j int) bool {
			return clients[i].Name < clients[j].Name
		})
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(clients)
	}
}

// activitiesHandler handles activity list operations
func activitiesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		activities := db.ListActivities()
		sort.Slice(activities, func(i, j int) bool {
			return activities[i].Date.After(activities[j].Date)
		})
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(activities)
	}
}

// createClientHandler shows the create client form
func createClientHandler(w http.ResponseWriter, r *http.Request) {
	htmlContent := createClientFormHTML()
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'self' 'unsafe-inline'; script-src 'self'")
	fmt.Fprint(w, htmlContent)
}

// apiCreateClientHandler creates a client via API
func apiCreateClientHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form", http.StatusBadRequest)
		return
	}

	// Validation
	name := strings.TrimSpace(r.FormValue("name"))
	email := strings.TrimSpace(r.FormValue("email"))
	phone := strings.TrimSpace(r.FormValue("phone"))

	if name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	if email != "" && !isValidEmail(email) {
		http.Error(w, "Invalid email format", http.StatusBadRequest)
		return
	}

	addr := &store.Address{
		Street:  strings.TrimSpace(r.FormValue("street")),
		City:    strings.TrimSpace(r.FormValue("city")),
		State:   strings.TrimSpace(r.FormValue("state")),
		ZipCode: strings.TrimSpace(r.FormValue("zipcode")),
		Country: strings.TrimSpace(r.FormValue("country")),
	}

	client := store.Client{
		Name:    name,
		Email:   email,
		Phone:   phone,
		Address: addr,
	}

	if err := db.CreateClient(client); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok", "message": "Client created successfully"})
}

// apiListClientsHandler lists clients
func apiListClientsHandler(w http.ResponseWriter, r *http.Request) {
	clients := db.ListClients()
	sort.Slice(clients, func(i, j int) bool {
		return clients[i].Name < clients[j].Name
	})
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(clients)
}

// apiUpdateClientHandler updates a client
func apiUpdateClientHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form", http.StatusBadRequest)
		return
	}

	clientID := r.FormValue("id")
	name := strings.TrimSpace(r.FormValue("name"))
	email := strings.TrimSpace(r.FormValue("email"))
	phone := strings.TrimSpace(r.FormValue("phone"))

	if name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	if email != "" && !isValidEmail(email) {
		http.Error(w, "Invalid email format", http.StatusBadRequest)
		return
	}

	addr := &store.Address{
		Street:  strings.TrimSpace(r.FormValue("street")),
		City:    strings.TrimSpace(r.FormValue("city")),
		State:   strings.TrimSpace(r.FormValue("state")),
		ZipCode: strings.TrimSpace(r.FormValue("zipcode")),
		Country: strings.TrimSpace(r.FormValue("country")),
	}

	client := store.Client{
		Name:    name,
		Email:   email,
		Phone:   phone,
		Address: addr,
	}

	if err := db.UpdateClient(clientID, client); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok", "message": "Client updated successfully"})
}

// apiDeleteClientHandler deletes a client
func apiDeleteClientHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	clientID := r.URL.Query().Get("id")
	if clientID == "" {
		http.Error(w, "Client ID required", http.StatusBadRequest)
		return
	}

	if err := db.DeleteClient(clientID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok", "message": "Client deleted successfully"})
}

// clientBalanceHandler calculates client balance
func clientBalanceHandler(w http.ResponseWriter, r *http.Request) {
	clientID := r.URL.Query().Get("id")
	if clientID == "" {
		http.Error(w, "Client ID required", http.StatusBadRequest)
		return
	}

	invoices := db.ListInvoices()
	var balance float64
	for _, inv := range invoices {
		if inv.ClientID == clientID {
			if inv.Status == store.InvoiceStatusPaid || inv.Status == store.InvoiceStatusVoid {
				continue
			}
			balance += inv.Total
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"balance": balance})
}

// apiCreateActivityHandler creates an activity
func apiCreateActivityHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form", http.StatusBadRequest)
		return
	}

	clientID := r.FormValue("client_id")
	projectID := r.FormValue("project_id")
	description := strings.TrimSpace(r.FormValue("description"))
	durationStr := r.FormValue("duration")

	if clientID == "" || description == "" || durationStr == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	duration, err := strconv.Atoi(durationStr)
	if err != nil || duration <= 0 {
		http.Error(w, "Invalid duration", http.StatusBadRequest)
		return
	}

	activity := store.Activity{
		ClientID:    clientID,
		ProjectID:   projectID,
		Description: description,
		Duration:    time.Duration(duration) * time.Minute,
		Date:        time.Now(),
	}

	if err := db.CreateActivity(activity); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok", "message": "Activity created successfully"})
}

// apiDeleteActivityHandler deletes an activity
func apiDeleteActivityHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	activityID := r.URL.Query().Get("id")
	if activityID == "" {
		http.Error(w, "Activity ID required", http.StatusBadRequest)
		return
	}

	if err := db.DeleteActivity(activityID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok", "message": "Activity deleted successfully"})
}

// invoicesHandler handles invoice operations
func invoicesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		invoices := db.ListInvoices()
		sort.Slice(invoices, func(i, j int) bool {
			return invoices[i].InvoiceDate.After(invoices[j].InvoiceDate)
		})
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(invoices)
	}
}

// apiCreateInvoiceHandler creates an invoice with line items
func apiCreateInvoiceHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form", http.StatusBadRequest)
		return
	}

	clientID := r.FormValue("client_id")
	invoiceNumber := strings.TrimSpace(r.FormValue("invoice_number"))
	taxStr := r.FormValue("tax")

	if clientID == "" || invoiceNumber == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	tax := 0.0
	if taxStr != "" {
		t, err := strconv.ParseFloat(taxStr, 64)
		if err != nil {
			http.Error(w, "Invalid tax amount", http.StatusBadRequest)
			return
		}
		tax = t
	}

	// Parse line items
	var lineItems []store.LineItem
	descriptions := r.Form["line_description"]
	quantities := r.Form["line_quantity"]
	rates := r.Form["line_rate"]

	for i := 0; i < len(descriptions); i++ {
		if descriptions[i] == "" {
			continue
		}

		qty, err := strconv.Atoi(quantities[i])
		if err != nil || qty <= 0 {
			http.Error(w, "Invalid quantity", http.StatusBadRequest)
			return
		}

		rate, err := strconv.ParseFloat(rates[i], 64)
		if err != nil || rate <= 0 {
			http.Error(w, "Invalid rate", http.StatusBadRequest)
			return
		}

		lineItems = append(lineItems, store.LineItem{
			Description: descriptions[i],
			Quantity:    qty,
			Rate:        rate,
			Subtotal:    float64(qty) * rate,
		})
	}

	if len(lineItems) == 0 {
		http.Error(w, "At least one line item required", http.StatusBadRequest)
		return
	}

	// Calculate subtotal
	subtotal := 0.0
	for _, item := range lineItems {
		subtotal += item.Subtotal
	}

	invoice := store.Invoice{
		ClientID:      clientID,
		InvoiceNumber: invoiceNumber,
		InvoiceDate:   time.Now(),
		DueDate:       time.Now().AddDate(0, 0, 30),
		LineItems:     lineItems,
		Subtotal:      subtotal,
		Tax:           tax,
		Total:         subtotal + tax,
		Status:        store.InvoiceStatusDraft,
	}

	if err := db.CreateInvoice(invoice); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok", "message": "Invoice created successfully"})
}

// apiUpdateInvoiceStatusHandler updates invoice status
func apiUpdateInvoiceStatusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form", http.StatusBadRequest)
		return
	}

	invoiceID := r.FormValue("id")
	status := r.FormValue("status")

	if invoiceID == "" || status == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	if err := db.SetInvoiceStatus(invoiceID, store.InvoiceStatus(status)); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok", "message": "Invoice status updated successfully"})
}

// apiCalculateTotalsHandler recalculates invoice totals
func apiCalculateTotalsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	invoiceID := r.URL.Query().Get("id")
	if invoiceID == "" {
		http.Error(w, "Invoice ID required", http.StatusBadRequest)
		return
	}

	if err := db.CalculateInvoiceTotals(invoiceID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok", "message": "Totals recalculated"})
}

// apiDeleteInvoiceHandler deletes an invoice
func apiDeleteInvoiceHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	invoiceID := r.URL.Query().Get("id")
	if invoiceID == "" {
		http.Error(w, "Invoice ID required", http.StatusBadRequest)
		return
	}

	if err := db.DeleteInvoice(invoiceID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok", "message": "Invoice deleted successfully"})
}

// invoicePrintHandler shows the printable invoice view
func invoicePrintHandler(w http.ResponseWriter, r *http.Request) {
	invoiceID := r.URL.Query().Get("id")
	if invoiceID == "" {
		http.Error(w, "Invoice ID required", http.StatusBadRequest)
		return
	}

	invoice, err := db.GetInvoice(invoiceID)
	if err != nil {
		http.Error(w, "Invoice not found", http.StatusNotFound)
		return
	}

	client, err := db.GetClient(invoice.ClientID)
	if err != nil {
		http.Error(w, "Client not found", http.StatusNotFound)
		return
	}

	htmlContent := invoicePrintHTML(invoice, client)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'self' 'unsafe-inline'; script-src 'self'")
	fmt.Fprint(w, htmlContent)
}

// Utility functions
func isValidEmail(email string) bool {
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func formatCurrency(val float64) string {
	return fmt.Sprintf("$%.2f", val)
}

func formatDate(t time.Time) string {
	return t.Format("2006-01-02")
}

// HTML Templates

func dashboardHTML(clients []store.Client, search string) string {
	clientRows := ""
	if len(clients) == 0 {
		clientRows = `<tr><td colspan="5" style="text-align:center;padding:3em;color:#999;">No clients found. <a href="/clients/new" style="color:#0066cc;">Create one</a></td></tr>`
	} else {
		for _, c := range clients {
			balance := 0.0
			invoices := db.ListInvoices()
			for _, inv := range invoices {
				if inv.ClientID == c.ID && inv.Status != store.InvoiceStatusPaid && inv.Status != store.InvoiceStatusVoid {
					balance += inv.Total
				}
			}
			clientRows += `<tr onclick="window.location='/clients/profile?id=` + html.EscapeString(c.ID) + `'" style="cursor:pointer;"><td>` + html.EscapeString(c.Name) + `</td><td>` + html.EscapeString(c.Email) + `</td><td>` + html.EscapeString(c.Phone) + `</td><td>` + strconv.Itoa(len(db.ListActivities())) + `</td><td>` + formatCurrency(balance) + `</td></tr>`
		}
	}

	return `<!DOCTYPE html>
<html>
<head>
	<meta charset="utf-8">
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<title>CRM Dashboard</title>
	<style>
		* { margin: 0; padding: 0; box-sizing: border-box; }
		body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif; background: #f8f9fa; color: #333; }
		header { background: white; border-bottom: 1px solid #e5e7eb; padding: 1.5rem 2rem; }
		header h1 { font-size: 1.75rem; font-weight: 600; margin: 0; }
		.container { max-width: 1200px; margin: 0 auto; padding: 2rem; }
		.search-box { display: flex; gap: 1rem; margin-bottom: 2rem; }
		.search-box input { flex: 1; padding: 0.75rem 1rem; border: 1px solid #ddd; border-radius: 0.5rem; font-size: 1rem; }
		.search-box button { padding: 0.75rem 1.5rem; background: #0066cc; color: white; border: none; border-radius: 0.5rem; cursor: pointer; font-weight: 500; }
		.search-box button:hover { background: #0052a3; }
		.btn-new { padding: 0.75rem 1.5rem; background: #10b981; color: white; border: none; border-radius: 0.5rem; cursor: pointer; font-weight: 500; text-decoration: none; display: inline-block; }
		.btn-new:hover { background: #059669; }
		table { width: 100%; border-collapse: collapse; background: white; border-radius: 0.5rem; box-shadow: 0 1px 3px rgba(0,0,0,0.1); }
		th { padding: 1rem; text-align: left; background: #f3f4f6; border-bottom: 2px solid #e5e7eb; font-weight: 600; color: #6b7280; }
		td { padding: 1rem; border-bottom: 1px solid #e5e7eb; }
		tr:hover { background: #f9fafb; }
	</style>
</head>
<body>
	<header>
		<h1>CRM Dashboard</h1>
	</header>
	<div class="container">
		<div class="search-box">
			<form method="get" style="display:flex;gap:1rem;width:100%;">
				<input type="text" name="search" placeholder="Search clients by name or email..." value="` + html.EscapeString(search) + `">
				<button type="submit">Search</button>
			</form>
			<a href="/clients/new" class="btn-new">Add Client</a>
		</div>
		<table>
			<thead>
				<tr>
					<th>Name</th>
					<th>Email</th>
					<th>Phone</th>
					<th>Activities</th>
					<th>Outstanding Balance</th>
				</tr>
			</thead>
			<tbody>` + clientRows + `</tbody>
		</table>
	</div>
</body>
</html>`
}

func createClientFormHTML() string {
	return `<!DOCTYPE html>
<html>
<head>
	<meta charset="utf-8">
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<title>Add Client</title>
	<style>
		* { margin: 0; padding: 0; box-sizing: border-box; }
		body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif; background: #f8f9fa; color: #333; }
		header { background: white; border-bottom: 1px solid #e5e7eb; padding: 1.5rem 2rem; }
		header h1 { font-size: 1.75rem; font-weight: 600; }
		.container { max-width: 600px; margin: 2rem auto; padding: 2rem; background: white; border-radius: 0.5rem; box-shadow: 0 1px 3px rgba(0,0,0,0.1); }
		form { display: flex; flex-direction: column; gap: 1.5rem; }
		.form-group { display: flex; flex-direction: column; }
		label { font-weight: 500; margin-bottom: 0.5rem; color: #374151; }
		input, textarea { padding: 0.75rem; border: 1px solid #ddd; border-radius: 0.5rem; font-size: 1rem; font-family: inherit; }
		input:focus, textarea:focus { outline: none; border-color: #0066cc; box-shadow: 0 0 0 3px rgba(0, 102, 204, 0.1); }
		.form-section { border-top: 1px solid #e5e7eb; padding-top: 1.5rem; }
		.form-section h3 { margin-bottom: 1rem; color: #6b7280; font-size: 0.875rem; text-transform: uppercase; }
		.button-group { display: flex; gap: 1rem; }
		button { padding: 0.75rem 1.5rem; background: #0066cc; color: white; border: none; border-radius: 0.5rem; cursor: pointer; font-weight: 500; }
		button:hover { background: #0052a3; }
		.btn-cancel { background: #e5e7eb; color: #333; }
		.btn-cancel:hover { background: #d1d5db; }
		.error { color: #dc2626; font-size: 0.875rem; margin-top: 0.5rem; display: none; }
		.success { color: #10b981; padding: 1rem; background: #f0fdf4; border-radius: 0.5rem; display: none; margin-bottom: 1rem; }
	</style>
</head>
<body>
	<header>
		<h1>Add Client</h1>
	</header>
	<div class="container">
		<div class="success" id="successMsg">Client created successfully! Redirecting...</div>
		<form id="clientForm">
			<div class="form-group">
				<label for="name">Name *</label>
				<input type="text" id="name" name="name" required>
				<span class="error" id="nameError">Name is required</span>
			</div>
			<div class="form-group">
				<label for="email">Email</label>
				<input type="email" id="email" name="email">
				<span class="error" id="emailError">Invalid email format</span>
			</div>
			<div class="form-group">
				<label for="phone">Phone</label>
				<input type="tel" id="phone" name="phone">
			</div>
			<div class="form-section">
				<h3>Address</h3>
				<div class="form-group">
					<label for="street">Street</label>
					<input type="text" id="street" name="street">
				</div>
				<div class="form-group">
					<label for="city">City</label>
					<input type="text" id="city" name="city">
				</div>
				<div class="form-group">
					<label for="state">State</label>
					<input type="text" id="state" name="state">
				</div>
				<div class="form-group">
					<label for="zipcode">Zip Code</label>
					<input type="text" id="zipcode" name="zipcode">
				</div>
				<div class="form-group">
					<label for="country">Country</label>
					<input type="text" id="country" name="country">
				</div>
			</div>
			<div class="button-group">
				<button type="submit">Create Client</button>
				<button type="button" class="btn-cancel" onclick="window.location='/'">Cancel</button>
			</div>
		</form>
	</div>
	<script>
		function validateEmail(email) {
			if (email === '') return true;
			return email.indexOf('@') > -1 && email.indexOf('.') > -1;
		}

		var form = document.getElementById('clientForm');
		form.addEventListener('submit', function(e) {
			e.preventDefault();

			var nameInput = document.getElementById('name');
			var emailInput = document.getElementById('email');
			var name = nameInput.value.trim();
			var email = emailInput.value.trim();

			var valid = true;
			document.getElementById('nameError').style.display = 'none';
			document.getElementById('emailError').style.display = 'none';

			if (name === '') {
				document.getElementById('nameError').style.display = 'block';
				valid = false;
			}

			if (email !== '' && !validateEmail(email)) {
				document.getElementById('emailError').style.display = 'block';
				valid = false;
			}

			if (!valid) return;

			var formData = new FormData(form);
			fetch('/clients/api/create', {
				method: 'POST',
				body: formData
			})
			.then(function(r) { return r.json(); })
			.then(function(data) {
				document.getElementById('successMsg').style.display = 'block';
				setTimeout(function() {
					window.location = '/';
				}, 1500);
			})
			.catch(function(err) {
				alert('Error: ' + err);
			});
		});
	</script>
</body>
</html>`
}

func clientProfileHTML(client *store.Client, tab string, invoices []store.Invoice, activities []store.Activity) string {
	var content string

	switch tab {
	case "history":
		if len(activities) == 0 {
			content = `<div style="padding:3em;text-align:center;color:#999;"><p>No activities recorded yet.</p></div>`
		} else {
			for _, a := range activities {
				content += `<div style="padding:1rem;border-bottom:1px solid #e5e7eb;">
					<div style="font-weight:500;">` + html.EscapeString(a.Description) + `</div>
					<div style="color:#6b7280;font-size:0.875rem;">` + formatDate(a.Date) + ` - ` + fmt.Sprintf("%.0f", a.Duration.Minutes()) + ` minutes</div>
				</div>`
			}
		}

		content += `<div style="padding:2rem;border-top:1px solid #e5e7eb;">
			<h3 style="margin-bottom:1rem;">Add Activity</h3>
			<form id="activityForm">
				<input type="hidden" name="client_id" value="` + html.EscapeString(client.ID) + `">
				<div style="margin-bottom:1rem;">
					<label style="display:block;margin-bottom:0.5rem;font-weight:500;">Description *</label>
					<textarea name="description" required style="width:100%;padding:0.75rem;border:1px solid #ddd;border-radius:0.5rem;"></textarea>
				</div>
				<div style="margin-bottom:1rem;">
					<label style="display:block;margin-bottom:0.5rem;font-weight:500;">Duration (minutes) *</label>
					<input type="number" name="duration" required min="1" style="width:100%;padding:0.75rem;border:1px solid #ddd;border-radius:0.5rem;">
				</div>
				<button type="submit" style="padding:0.75rem 1.5rem;background:#0066cc;color:white;border:none;border-radius:0.5rem;cursor:pointer;">Add Activity</button>
			</form>
			<script>
				var actForm = document.getElementById('activityForm');
				actForm.addEventListener('submit', function(e) {
					e.preventDefault();
					fetch('/activities/api/create', {
						method: 'POST',
						body: new FormData(actForm)
					})
					.then(function(r) { return r.json(); })
					.then(function() {
						window.location = window.location;
					})
					.catch(function(err) { alert('Error: ' + err); });
				});
			</script>
		</div>`

	case "billing":
		if len(invoices) == 0 {
			content = `<div style="padding:3em;text-align:center;color:#999;"><p>No invoices yet.</p></div>`
		} else {
			for _, inv := range invoices {
				statusColor := "#dc2626"
				if inv.Status == store.InvoiceStatusPaid {
					statusColor = "#10b981"
				} else if inv.Status == store.InvoiceStatusVoid {
					statusColor = "#9ca3af"
				}
				content += `<div style="padding:1rem;border-bottom:1px solid #e5e7eb;display:flex;justify-content:space-between;align-items:center;">
					<div>
						<div style="font-weight:500;">` + html.EscapeString(inv.InvoiceNumber) + `</div>
						<div style="color:#6b7280;font-size:0.875rem;">` + formatDate(inv.InvoiceDate) + `</div>
						<div style="color:#6b7280;font-size:0.875rem;">` + formatCurrency(inv.Total) + `</div>
					</div>
					<div style="text-align:right;">
						<span style="background:` + statusColor + `;color:white;padding:0.25rem 0.75rem;border-radius:0.25rem;font-size:0.875rem;font-weight:500;">` + string(inv.Status) + `</span>
						<div style="margin-top:0.5rem;">
							<a href="/invoices/print?id=` + html.EscapeString(inv.ID) + `" style="color:#0066cc;margin-right:1rem;">Print</a>
							<button onclick="updateStatus(this, '` + html.EscapeString(inv.ID) + `', 'sent')" style="background:none;border:none;color:#0066cc;cursor:pointer;margin-right:1rem;">Send</button>
							<button onclick="updateStatus(this, '` + html.EscapeString(inv.ID) + `', 'paid')" style="background:none;border:none;color:#0066cc;cursor:pointer;margin-right:1rem;">Mark Paid</button>
							<button onclick="updateStatus(this, '` + html.EscapeString(inv.ID) + `', 'void')" style="background:none;border:none;color:#dc2626;cursor:pointer;">Void</button>
						</div>
					</div>
				</div>`
			}
		}

		content += `<div style="padding:2rem;border-top:1px solid #e5e7eb;">
			<h3 style="margin-bottom:1rem;">Create Invoice</h3>
			<form id="invoiceForm">
				<input type="hidden" name="client_id" value="` + html.EscapeString(client.ID) + `">
				<div style="margin-bottom:1rem;">
					<label style="display:block;margin-bottom:0.5rem;font-weight:500;">Invoice Number *</label>
					<input type="text" name="invoice_number" required style="width:100%;padding:0.75rem;border:1px solid #ddd;border-radius:0.5rem;">
				</div>
				<div style="margin-bottom:1rem;">
					<label style="display:block;margin-bottom:0.5rem;font-weight:500;">Tax</label>
					<input type="number" name="tax" step="0.01" value="0" style="width:100%;padding:0.75rem;border:1px solid #ddd;border-radius:0.5rem;">
				</div>
				<div id="lineItems">
					<div class="lineItem" style="margin-bottom:1rem;padding:1rem;background:#f9fafb;border-radius:0.5rem;">
						<input type="text" name="line_description" placeholder="Item description" required style="width:100%;margin-bottom:0.5rem;padding:0.75rem;border:1px solid #ddd;border-radius:0.5rem;">
						<div style="display:flex;gap:0.5rem;">
							<input type="number" name="line_quantity" placeholder="Qty" required min="1" style="flex:1;padding:0.75rem;border:1px solid #ddd;border-radius:0.5rem;">
							<input type="number" name="line_rate" placeholder="Rate" required step="0.01" min="0" style="flex:1;padding:0.75rem;border:1px solid #ddd;border-radius:0.5rem;">
						</div>
					</div>
				</div>
				<button type="button" onclick="addLineItem()" style="margin-bottom:1rem;padding:0.75rem 1.5rem;background:#e5e7eb;color:#333;border:none;border-radius:0.5rem;cursor:pointer;">Add Line Item</button>
				<button type="submit" style="padding:0.75rem 1.5rem;background:#0066cc;color:white;border:none;border-radius:0.5rem;cursor:pointer;">Create Invoice</button>
			</form>
			<script>
				function addLineItem() {
					var li = document.createElement('div');
					li.className = 'lineItem';
					li.style.cssText = 'margin-bottom:1rem;padding:1rem;background:#f9fafb;border-radius:0.5rem;';
					li.innerHTML = '<input type="text" name="line_description" placeholder="Item description" required style="width:100%;margin-bottom:0.5rem;padding:0.75rem;border:1px solid #ddd;border-radius:0.5rem;"><div style="display:flex;gap:0.5rem;"><input type="number" name="line_quantity" placeholder="Qty" required min="1" style="flex:1;padding:0.75rem;border:1px solid #ddd;border-radius:0.5rem;"><input type="number" name="line_rate" placeholder="Rate" required step="0.01" min="0" style="flex:1;padding:0.75rem;border:1px solid #ddd;border-radius:0.5rem;"></div>';
					document.getElementById('lineItems').appendChild(li);
				}

				var invForm = document.getElementById('invoiceForm');
				invForm.addEventListener('submit', function(e) {
					e.preventDefault();
					fetch('/invoices/api/create', {
						method: 'POST',
						body: new FormData(invForm)
					})
					.then(function(r) { return r.json(); })
					.then(function() {
						window.location = window.location;
					})
					.catch(function(err) { alert('Error: ' + err); });
				});

				function updateStatus(btn, invId, status) {
					var fd = new FormData();
					fd.append('id', invId);
					fd.append('status', status);
					fetch('/invoices/api/update-status', {
						method: 'POST',
						body: fd
					})
					.then(function(r) { return r.json(); })
					.then(function() {
						window.location = window.location;
					})
					.catch(function(err) { alert('Error: ' + err); });
				}
			</script>
		</div>`

	default: // details
		content = `<div style="padding:2rem;">
			<h3 style="margin-bottom:1rem;">Basic Information</h3>
			<div style="display:grid;grid-template-columns:1fr 1fr;gap:2rem;margin-bottom:2rem;">
				<div>
					<div style="color:#6b7280;font-size:0.875rem;">Name</div>
					<div style="font-weight:500;">` + html.EscapeString(client.Name) + `</div>
				</div>
				<div>
					<div style="color:#6b7280;font-size:0.875rem;">Email</div>
					<div style="font-weight:500;">` + html.EscapeString(client.Email) + `</div>
				</div>
				<div>
					<div style="color:#6b7280;font-size:0.875rem;">Phone</div>
					<div style="font-weight:500;">` + html.EscapeString(client.Phone) + `</div>
				</div>
				<div>
					<div style="color:#6b7280;font-size:0.875rem;">Member Since</div>
					<div style="font-weight:500;">` + formatDate(client.CreatedAt) + `</div>
				</div>
			</div>`

		if client.Address != nil {
			content += `<h3 style="margin-bottom:1rem;">Address</h3>
			<div style="display:grid;grid-template-columns:1fr 1fr;gap:2rem;margin-bottom:2rem;">
				<div>
					<div style="color:#6b7280;font-size:0.875rem;">Street</div>
					<div style="font-weight:500;">` + html.EscapeString(client.Address.Street) + `</div>
				</div>
				<div>
					<div style="color:#6b7280;font-size:0.875rem;">City</div>
					<div style="font-weight:500;">` + html.EscapeString(client.Address.City) + `</div>
				</div>
				<div>
					<div style="color:#6b7280;font-size:0.875rem;">State</div>
					<div style="font-weight:500;">` + html.EscapeString(client.Address.State) + `</div>
				</div>
				<div>
					<div style="color:#6b7280;font-size:0.875rem;">Zip Code</div>
					<div style="font-weight:500;">` + html.EscapeString(client.Address.ZipCode) + `</div>
				</div>
				<div>
					<div style="color:#6b7280;font-size:0.875rem;">Country</div>
					<div style="font-weight:500;">` + html.EscapeString(client.Address.Country) + `</div>
				</div>
			</div>`
		}

		content += `<div style="display:flex;gap:1rem;padding-top:2rem;border-top:1px solid #e5e7eb;">
			<button onclick="deleteClient('` + html.EscapeString(client.ID) + `')" style="padding:0.75rem 1.5rem;background:#dc2626;color:white;border:none;border-radius:0.5rem;cursor:pointer;">Delete Client</button>
		</div>
		<script>
			function deleteClient(id) {
				if (confirm('Are you sure you want to delete this client?')) {
					fetch('/clients/api/delete?id=' + id, { method: 'POST' })
					.then(function(r) { return r.json(); })
					.then(function() {
						window.location = '/';
					})
					.catch(function(err) { alert('Error: ' + err); });
				}
			}
		</script>
		</div>`
	}

	return `<!DOCTYPE html>
<html>
<head>
	<meta charset="utf-8">
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<title>` + client.Name + `</title>
	<style>
		* { margin: 0; padding: 0; box-sizing: border-box; }
		body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif; background: #f8f9fa; color: #333; }
		header { background: white; border-bottom: 1px solid #e5e7eb; padding: 1.5rem 2rem; }
		header h1 { font-size: 1.75rem; font-weight: 600; }
		.container { max-width: 1000px; margin: 0 auto; }
		.tabs { display: flex; background: white; border-bottom: 1px solid #e5e7eb; }
		.tab { padding: 1rem 2rem; cursor: pointer; border-bottom: 3px solid transparent; font-weight: 500; transition: all 0.2s ease; outline: none; }
		.tab:hover { background: #f9fafb; border-bottom-color: #ccc; }
		.tab:focus { outline: 2px solid #0066cc; outline-offset: -2px; }
		.tab.active { border-bottom-color: #0066cc; color: #0066cc; }
		.tab.active:hover { background: white; border-bottom-color: #0066cc; }
		.content { background: white; }
	</style>
</head>
<body>
	<header>
		<h1>` + html.EscapeString(client.Name) + `</h1>
	</header>
	<div class="container">
		<nav class="tabs" role="tablist">
			<div class="tab` + (map[string]string{"details": " active"}[tab]) + `" role="tab" aria-selected="` + (map[string]string{"details": "true"}[tab]) + `" tabindex="0" onclick="window.location='/clients/profile?id=` + html.EscapeString(client.ID) + `&tab=details'" onkeydown="if(event.key==='Enter'||event.key===' ')window.location='/clients/profile?id=` + html.EscapeString(client.ID) + `&tab=details'">Details</div>
			<div class="tab` + (map[string]string{"history": " active"}[tab]) + `" role="tab" aria-selected="` + (map[string]string{"history": "true"}[tab]) + `" tabindex="0" onclick="window.location='/clients/profile?id=` + html.EscapeString(client.ID) + `&tab=history'" onkeydown="if(event.key==='Enter'||event.key===' ')window.location='/clients/profile?id=` + html.EscapeString(client.ID) + `&tab=history'">History</div>
			<div class="tab` + (map[string]string{"billing": " active"}[tab]) + `" role="tab" aria-selected="` + (map[string]string{"billing": "true"}[tab]) + `" tabindex="0" onclick="window.location='/clients/profile?id=` + html.EscapeString(client.ID) + `&tab=billing'" onkeydown="if(event.key==='Enter'||event.key===' ')window.location='/clients/profile?id=` + html.EscapeString(client.ID) + `&tab=billing'">Billing</div>
		</nav>
		<div class="content">
			` + content + `
		</div>
	</div>
</body>
</html>`
}

func invoicePrintHTML(invoice *store.Invoice, client *store.Client) string {
	itemsHTML := ""
	for _, item := range invoice.LineItems {
		itemsHTML += `<tr><td style="padding:0.75rem;border-bottom:1px solid #e5e7eb;">` + html.EscapeString(item.Description) + `</td><td style="padding:0.75rem;border-bottom:1px solid #e5e7eb;text-align:right;">` + fmt.Sprintf("%d", item.Quantity) + `</td><td style="padding:0.75rem;border-bottom:1px solid #e5e7eb;text-align:right;">` + formatCurrency(item.Rate) + `</td><td style="padding:0.75rem;border-bottom:1px solid #e5e7eb;text-align:right;">` + formatCurrency(item.Subtotal) + `</td></tr>`
	}

	clientAddr := ""
	if client.Address != nil && client.Address.Street != "" {
		clientAddr = `<div>` + html.EscapeString(client.Address.Street) + `</div><div>` + html.EscapeString(client.Address.City) + `, ` + html.EscapeString(client.Address.State) + ` ` + html.EscapeString(client.Address.ZipCode) + `</div>`
	}

	return `<!DOCTYPE html>
<html>
<head>
	<meta charset="utf-8">
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<title>Invoice ` + html.EscapeString(invoice.InvoiceNumber) + `</title>
	<style>
		* { margin: 0; padding: 0; box-sizing: border-box; }
		body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif; background: white; color: #333; }
		.container { max-width: 800px; margin: 0 auto; padding: 3rem 2rem; }
		.header { margin-bottom: 3rem; border-bottom: 2px solid #000; padding-bottom: 2rem; }
		.header h1 { font-size: 2rem; margin-bottom: 1rem; }
		.header-info { display: grid; grid-template-columns: 1fr 1fr; gap: 2rem; }
		table { width: 100%; border-collapse: collapse; margin-bottom: 2rem; }
		th { text-align: left; padding: 0.75rem; background: #f3f4f6; font-weight: 600; border-bottom: 2px solid #000; }
		.totals { margin-left: auto; width: 300px; }
		.total-row { display: flex; justify-content: space-between; padding: 0.75rem 0; border-bottom: 1px solid #e5e7eb; }
		.total-row.final { border-bottom: 2px solid #000; border-top: 2px solid #000; font-weight: 600; font-size: 1.25rem; padding: 1rem 0; }
		@media print { body { background: white; } .container { padding: 0; } }
	</style>
</head>
<body>
	<div class="container">
		<div class="header">
			<h1>INVOICE</h1>
			<div class="header-info">
				<div>
					<h2 style="font-size: 1rem; font-weight: 600; margin: 0;">Invoice #</h2>
					<div>` + html.EscapeString(invoice.InvoiceNumber) + `</div>
					<div style="margin-top:0.5rem;"><strong>Date:</strong> ` + formatDate(invoice.InvoiceDate) + `</div>
					<div><strong>Due Date:</strong> ` + formatDate(invoice.DueDate) + `</div>
				</div>
				<div style="text-align:right;">
					<div style="font-weight:600;font-size:1.5rem;">` + html.EscapeString(client.Name) + `</div>
					<div style="margin-top:1rem;color:#6b7280;">
						<div>` + html.EscapeString(client.Email) + `</div>
						<div>` + html.EscapeString(client.Phone) + `</div>
						` + clientAddr + `
					</div>
				</div>
			</div>
		</div>
		<table>
			<thead>
				<tr>
					<th>Description</th>
					<th style="text-align:right;">Qty</th>
					<th style="text-align:right;">Rate</th>
					<th style="text-align:right;">Amount</th>
				</tr>
			</thead>
			<tbody>
				` + itemsHTML + `
			</tbody>
		</table>
		<div class="totals">
			<div class="total-row">
				<span>Subtotal</span>
				<span>` + formatCurrency(invoice.Subtotal) + `</span>
			</div>
			<div class="total-row">
				<span>Tax</span>
				<span>` + formatCurrency(invoice.Tax) + `</span>
			</div>
			<div class="total-row final">
				<span>Total Due</span>
				<span>` + formatCurrency(invoice.Total) + `</span>
			</div>
		</div>
	</div>
	<script>
		window.addEventListener('load', function() {
			window.print();
		});
	</script>
</body>
</html>`
}
