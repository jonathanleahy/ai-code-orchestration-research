package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"app/store"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

var db store.Store

func init() {
	db = store.NewInMemoryStore()
}

func main() {
	http.HandleFunc("/health", handleHealth)
	http.HandleFunc("/", handleDashboard)
	http.HandleFunc("/client/", handleClientRouter)
	http.HandleFunc("/invoice/", handleInvoiceRouter)
	http.HandleFunc("/api/clients", handleClientsAPI)
	http.HandleFunc("/api/invoices", handleInvoicesAPI)
	http.HandleFunc("/api/activities", handleActivitiesAPI)

	log.Println("Starting CRM server at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func handleDashboard(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	query := r.URL.Query().Get("q")
	var clients []*store.Client

	if query != "" {
		var err error
		clients, err = db.SearchClients(query)
		if err != nil {
			http.Error(w, "Search failed", http.StatusInternalServerError)
			return
		}
	} else {
		clients = getAllClients()
	}

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, dashboardHTML(clients, query))
}

func handleClientRouter(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/client/")
	parts := strings.Split(path, "/")

	if len(parts) == 0 || parts[0] == "" {
		http.NotFound(w, r)
		return
	}

	clientID := parts[0]

	if len(parts) == 1 {
		if r.Method == http.MethodGet {
			handleClientProfile(w, r, clientID)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
		return
	}

	if len(parts) >= 2 && parts[1] == "invoice" && len(parts) >= 3 && parts[2] == "new" {
		if r.Method == http.MethodGet {
			handleNewInvoiceForm(w, r, clientID)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
		return
	}

	http.NotFound(w, r)
}

func handleInvoiceRouter(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/invoice/")
	parts := strings.Split(path, "/")

	if len(parts) == 0 || parts[0] == "" {
		http.NotFound(w, r)
		return
	}

	invoiceID := parts[0]

	if len(parts) == 1 {
		if r.Method == http.MethodGet {
			handleInvoiceDetail(w, r, invoiceID)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
		return
	}

	if len(parts) >= 2 && parts[1] == "print" {
		if r.Method == http.MethodGet {
			handleInvoicePrint(w, r, invoiceID)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
		return
	}

	http.NotFound(w, r)
}

func handleClientsAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodPost:
		handleCreateClient(w, r)
	case http.MethodGet:
		handleGetClients(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Method not allowed"})
	}
}

func handleInvoicesAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	path := strings.TrimPrefix(r.URL.Path, "/api/invoices")
	if path == "" || path == "/" {
		if r.Method == http.MethodPost {
			handleCreateInvoice(w, r)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Method not allowed"})
		}
		return
	}

	parts := strings.Split(strings.TrimPrefix(path, "/"), "/")
	if len(parts) < 1 || parts[0] == "" {
		http.NotFound(w, r)
		return
	}

	invoiceID := parts[0]

	if len(parts) >= 2 {
		action := parts[1]
		if r.Method == http.MethodPost {
			switch action {
			case "mark-sent":
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
			case "mark-paid":
				db.MarkInvoicePaid(invoiceID)
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
			case "void":
				db.VoidInvoice(invoiceID)
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
			default:
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(ErrorResponse{Error: "Action not found"})
			}
		}
		return
	}

	if r.Method == http.MethodPost {
		handleUpdateInvoice(w, r, invoiceID)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Method not allowed"})
	}
}

func handleActivitiesAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodPost {
		handleCreateActivity(w, r)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Method not allowed"})
	}
}

func handleCreateClient(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)

	client := &store.Client{
		ID:      fmt.Sprintf("c%d", time.Now().UnixNano()),
		Name:    r.FormValue("name"),
		Email:   r.FormValue("email"),
		Phone:   r.FormValue("phone"),
		Company: r.FormValue("company"),
		Address: r.FormValue("address"),
	}

	if err := db.CreateClient(client); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to create client"})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(client)
}

func handleGetClients(w http.ResponseWriter, r *http.Request) {
	clients := getAllClients()
	json.NewEncoder(w).Encode(clients)
}

func handleClientProfile(w http.ResponseWriter, r *http.Request, clientID string) {
	client, _ := db.GetClient(clientID)
	if client == nil {
		http.NotFound(w, r)
		return
	}

	activities, _ := db.GetClientActivities(clientID)
	invoices, _ := db.GetClientInvoices(clientID)

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, clientProfileHTML(client, activities, invoices))
}

func handleNewInvoiceForm(w http.ResponseWriter, r *http.Request, clientID string) {
	client, _ := db.GetClient(clientID)
	if client == nil {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, newInvoiceFormHTML(client))
}

func handleCreateInvoice(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)

	clientID := r.FormValue("client_id")
	dueDate, _ := time.Parse("2006-01-02", r.FormValue("due_date"))

	var items []store.InvoiceItem
	itemCount := 0
	for i := 0; i < 100; i++ {
		desc := r.FormValue("item_desc_" + strconv.Itoa(i))
		if desc == "" {
			continue
		}
		qty, _ := strconv.Atoi(r.FormValue("item_qty_" + strconv.Itoa(i)))
		rate, _ := strconv.ParseFloat(r.FormValue("item_rate_" + strconv.Itoa(i)), 64)

		items = append(items, store.InvoiceItem{
			Description: desc,
			Quantity:    qty,
			Rate:        rate,
		})
		itemCount++
	}

	invoice := &store.Invoice{
		ID:       fmt.Sprintf("inv%d", time.Now().UnixNano()),
		ClientID: clientID,
		Status:   "draft",
		Items:    items,
		DueDate:  dueDate,
	}

	if err := db.CreateInvoice(invoice); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to create invoice"})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(invoice)
}

func handleUpdateInvoice(w http.ResponseWriter, r *http.Request, invoiceID string) {
	r.ParseMultipartForm(10 << 20)

	invoice, _ := db.GetInvoice(invoiceID)
	if invoice == nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Invoice not found"})
		return
	}

	dueDate, _ := time.Parse("2006-01-02", r.FormValue("due_date"))
	invoice.DueDate = dueDate

	var items []store.InvoiceItem
	for i := 0; i < 100; i++ {
		desc := r.FormValue("item_desc_" + strconv.Itoa(i))
		if desc == "" {
			continue
		}
		qty, _ := strconv.Atoi(r.FormValue("item_qty_" + strconv.Itoa(i)))
		rate, _ := strconv.ParseFloat(r.FormValue("item_rate_" + strconv.Itoa(i)), 64)

		items = append(items, store.InvoiceItem{
			Description: desc,
			Quantity:    qty,
			Rate:        rate,
		})
	}

	invoice.Items = items

	if err := db.UpdateInvoice(invoiceID, invoice); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to update invoice"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(invoice)
}

func handleCreateActivity(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)

	activity := &store.Activity{
		ID:       fmt.Sprintf("act%d", time.Now().UnixNano()),
		ClientID: r.FormValue("client_id"),
		Type:     r.FormValue("type"),
		Notes:    r.FormValue("description"),
	}

	if err := db.CreateActivity(activity); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to create activity"})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(activity)
}

func handleInvoiceDetail(w http.ResponseWriter, r *http.Request, invoiceID string) {
	invoice, _ := db.GetInvoice(invoiceID)
	if invoice == nil {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(invoice)
}

func handleInvoicePrint(w http.ResponseWriter, r *http.Request, invoiceID string) {
	invoice, _ := db.GetInvoice(invoiceID)
	if invoice == nil {
		http.NotFound(w, r)
		return
	}

	client, _ := db.GetClient(invoice.ClientID)

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, invoicePrintHTML(invoice, client))
}

func getAllClients() []*store.Client {
	var results []*store.Client
	clients, _ := db.SearchClients("")
	return append(results, clients...)
}

// HTML Templates

func dashboardHTML(clients []*store.Client, query string) string {
	html := `<!DOCTYPE html>
<html>
<head>
	<title>CRM Dashboard</title>
	<style>
		* { margin: 0; padding: 0; box-sizing: border-box; }
		body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif; background: #f8f9fa; color: #333; }
		.header { background: white; padding: 20px 40px; border-bottom: 1px solid #e0e0e0; }
		.header h1 { font-size: 24px; margin-bottom: 10px; }
		.search-bar { display: flex; gap: 10px; margin-bottom: 20px; }
		.search-bar input { flex: 1; padding: 10px; border: 1px solid #ddd; border-radius: 4px; }
		.search-bar button { padding: 10px 20px; background: #007bff; color: white; border: none; border-radius: 4px; cursor: pointer; }
		.container { max-width: 1200px; margin: 40px auto; padding: 0 20px; }
		.grid { display: grid; grid-template-columns: 1fr 1fr; gap: 30px; margin-bottom: 40px; }
		.card { background: white; padding: 30px; border-radius: 8px; box-shadow: 0 1px 3px rgba(0,0,0,0.1); }
		.card h2 { margin-bottom: 20px; font-size: 18px; }
		.form-group { margin-bottom: 15px; }
		.form-group label { display: block; margin-bottom: 5px; font-weight: 500; }
		.form-group input, .form-group textarea { width: 100%; padding: 8px; border: 1px solid #ddd; border-radius: 4px; font-family: inherit; }
		.form-group textarea { resize: vertical; min-height: 80px; }
		.form-group button { padding: 10px 20px; background: #28a745; color: white; border: none; border-radius: 4px; cursor: pointer; }
		.clients-list { list-style: none; }
		.client-item { padding: 12px; border-bottom: 1px solid #eee; }
		.client-item:last-child { border-bottom: none; }
		.client-link { color: #007bff; text-decoration: none; }
		.client-link:hover { text-decoration: underline; }
		.client-email { font-size: 13px; color: #666; }
	</style>
</head>
<body>
	<div class="header">
		<h1>CRM Dashboard</h1>
		<div class="search-bar">
			<form style="display: flex; gap: 10px; flex: 1;">
				<input type="text" name="q" placeholder="Search by name or email..." value="` + query + `">
				<button type="submit">Search</button>
			</form>
		</div>
	</div>

	<div class="container">
		<div class="grid">
			<div class="card">
				<h2>Add New Client</h2>
				<form id="addClientForm">
					<div class="form-group">
						<label>Name</label>
						<input type="text" name="name" required>
					</div>
					<div class="form-group">
						<label>Email</label>
						<input type="email" name="email" required>
					</div>
					<div class="form-group">
						<label>Phone</label>
						<input type="tel" name="phone">
					</div>
					<div class="form-group">
						<label>Company</label>
						<input type="text" name="company">
					</div>
					<div class="form-group">
						<label>Address</label>
						<textarea name="address"></textarea>
					</div>
					<div class="form-group">
						<button type="submit">Add Client</button>
					</div>
				</form>
			</div>

			<div class="card">
				<h2>Clients (` + strconv.Itoa(len(clients)) + `)</h2>
				<ul class="clients-list">`

	for _, client := range clients {
		html += `<li class="client-item">
					<a href="/client/` + client.ID + `" class="client-link">` + client.Name + `</a>
					<div class="client-email">` + client.Email + `</div>
				</li>`
	}

	html += `</ul>
			</div>
		</div>
	</div>

	<script>
		document.getElementById('addClientForm').addEventListener('submit', function(e) {
			e.preventDefault();
			var form = new FormData(this);
			fetch('/api/clients', {
				method: 'POST',
				body: form
			}).then(r => r.json()).then(d => {
				alert('Client added');
				window.location.reload();
			}).catch(e => alert('Error: ' + e));
		});
	</script>
</body>
</html>`
	return html
}

func clientProfileHTML(client *store.Client, activities []*store.Activity, invoices []*store.Invoice) string {
	html := `<!DOCTYPE html>
<html>
<head>
	<title>` + client.Name + ` - Client Profile</title>
	<style>
		* { margin: 0; padding: 0; box-sizing: border-box; }
		body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif; background: #f8f9fa; color: #333; }
		.header { background: white; padding: 20px 40px; border-bottom: 1px solid #e0e0e0; }
		.header h1 { font-size: 24px; margin-bottom: 10px; }
		.container { max-width: 1000px; margin: 40px auto; padding: 0 20px; }
		.card { background: white; padding: 30px; border-radius: 8px; box-shadow: 0 1px 3px rgba(0,0,0,0.1); margin-bottom: 30px; }
		.card h2 { margin-bottom: 20px; }
		.buttons { display: flex; gap: 10px; margin-bottom: 20px; }
		.btn { padding: 10px 20px; border: none; border-radius: 4px; cursor: pointer; }
		.btn-primary { background: #007bff; color: white; }
		.btn-danger { background: #dc3545; color: white; }
		.btn-success { background: #28a745; color: white; }
		.tabs { display: flex; gap: 10px; margin-bottom: 20px; border-bottom: 1px solid #eee; }
		.tab { padding: 10px 15px; cursor: pointer; border-bottom: 2px solid transparent; }
		.tab.active { border-bottom-color: #007bff; color: #007bff; }
		.tab-content { display: none; }
		.tab-content.active { display: block; }
		.info-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 20px; margin-bottom: 20px; }
		.info-item label { font-weight: 500; display: block; margin-bottom: 5px; }
		.info-item div { color: #666; }
		.form-group { margin-bottom: 15px; }
		.form-group label { display: block; margin-bottom: 5px; font-weight: 500; }
		.form-group input, .form-group textarea { width: 100%; padding: 8px; border: 1px solid #ddd; border-radius: 4px; }
		.form-group textarea { resize: vertical; min-height: 80px; }
		.activity-item { padding: 15px; border-left: 3px solid #007bff; background: #f8f9fa; margin-bottom: 10px; border-radius: 4px; }
		.activity-type { font-weight: 500; color: #007bff; }
		.activity-date { font-size: 12px; color: #999; }
		.invoice-item { padding: 10px; border-bottom: 1px solid #eee; display: flex; justify-content: space-between; align-items: center; }
		.invoice-item:last-child { border-bottom: none; }
		.hidden { display: none; }
	</style>
</head>
<body>
	<div class="header">
		<h1>` + client.Name + `</h1>
		<div class="buttons">
			<button class="btn btn-primary" onclick="toggleEditForm()">Edit</button>
			<button class="btn btn-danger" onclick="deleteClient()">Delete</button>
		</div>
	</div>

	<div class="container">
		<div class="card">
			<div id="viewInfo">
				<h2>Client Information</h2>
				<div class="info-grid">
					<div class="info-item">
						<label>Email</label>
						<div>` + client.Email + `</div>
					</div>
					<div class="info-item">
						<label>Phone</label>
						<div>` + client.Phone + `</div>
					</div>
					<div class="info-item">
						<label>Company</label>
						<div>` + client.Company + `</div>
					</div>
					<div class="info-item">
						<label>Address</label>
						<div>` + client.Address + `</div>
					</div>
				</div>
			</div>

			<div id="editForm" class="hidden">
				<h2>Edit Client</h2>
				<form id="updateClientForm">
					<div class="form-group">
						<label>Name</label>
						<input type="text" name="name" value="` + client.Name + `" required>
					</div>
					<div class="form-group">
						<label>Email</label>
						<input type="email" name="email" value="` + client.Email + `" required>
					</div>
					<div class="form-group">
						<label>Phone</label>
						<input type="tel" name="phone" value="` + client.Phone + `">
					</div>
					<div class="form-group">
						<label>Company</label>
						<input type="text" name="company" value="` + client.Company + `">
					</div>
					<div class="form-group">
						<label>Address</label>
						<textarea name="address">` + client.Address + `</textarea>
					</div>
					<div style="display: flex; gap: 10px;">
						<button type="submit" class="btn btn-success">Save</button>
						<button type="button" class="btn" onclick="toggleEditForm()" style="background: #6c757d; color: white;">Cancel</button>
					</div>
				</form>
			</div>
		</div>

		<div class="card">
			<div class="tabs">
				<div class="tab active" onclick="switchTab('activity')">Activity</div>
				<div class="tab" onclick="switchTab('invoices')">Invoices</div>
			</div>

			<div id="activity" class="tab-content active">
				<h3>Activity Timeline</h3>
				<div style="margin-bottom: 20px;">`

	for _, activity := range activities {
		html += `<div class="activity-item">
					<div>
						<div class="activity-type">` + activity.Type + `</div>
						<div>` + activity.Notes + `</div>
						<div class="activity-date">` + activity.CreatedAt.Format("2006-01-02 15:04") + `</div>
					</div>
				</div>`
	}

	html += `</div>
				<h3>Add Activity</h3>
				<form id="addActivityForm">
					<div class="form-group">
						<label>Type</label>
						<select name="type" style="width: 100%; padding: 8px; border: 1px solid #ddd; border-radius: 4px;">
							<option value="call">Call</option>
							<option value="email">Email</option>
							<option value="meeting">Meeting</option>
							<option value="note">Note</option>
						</select>
					</div>
					<div class="form-group">
						<label>Description</label>
						<textarea name="description"></textarea>
					</div>
					<button type="submit" class="btn btn-success">Add Activity</button>
				</form>
			</div>

			<div id="invoices" class="tab-content">
				<h3>Invoices</h3>
				<div style="margin-bottom: 20px;">`

	for _, invoice := range invoices {
		total := 0.0
		for _, item := range invoice.Items {
			total += item.Amount
		}
		html += `<div class="invoice-item">
					<div>
						<strong>Invoice ` + invoice.ID + `</strong> - Status: ` + invoice.Status + ` - $` + fmt.Sprintf("%.2f", total) + `
					</div>
					<div>
						<a href="/invoice/` + invoice.ID + `/print" class="btn btn-primary" style="padding: 5px 10px; text-decoration: none; color: white; font-size: 12px;">Print</a>
					</div>
				</div>`
	}

	html += `</div>
				<a href="/client/` + client.ID + `/invoice/new" class="btn btn-success">Create Invoice</a>
			</div>
		</div>
	</div>

	<script>
		function toggleEditForm() {
			document.getElementById('viewInfo').classList.toggle('hidden');
			document.getElementById('editForm').classList.toggle('hidden');
		}

		function switchTab(tab) {
			document.querySelectorAll('.tab-content').forEach(e => e.classList.remove('active'));
			document.querySelectorAll('.tab').forEach(e => e.classList.remove('active'));
			document.getElementById(tab).classList.add('active');
			event.target.classList.add('active');
		}

		function deleteClient() {
			if (confirm('Are you sure you want to delete this client?')) {
				fetch('/api/clients/` + client.ID + `', { method: 'DELETE' }).then(() => {
					alert('Client deleted');
					window.location.href = '/';
				});
			}
		}

		document.getElementById('updateClientForm').addEventListener('submit', function(e) {
			e.preventDefault();
			var form = new FormData(this);
			fetch('/api/clients/` + client.ID + `', {
				method: 'POST',
				body: form
			}).then(r => r.json()).then(d => {
				alert('Client updated');
				window.location.reload();
			}).catch(e => alert('Error: ' + e));
		});

		document.getElementById('addActivityForm').addEventListener('submit', function(e) {
			e.preventDefault();
			var form = new FormData(this);
			form.append('client_id', '` + client.ID + `');
			fetch('/api/activities', {
				method: 'POST',
				body: form
			}).then(r => r.json()).then(d => {
				alert('Activity added');
				window.location.reload();
			}).catch(e => alert('Error: ' + e));
		});
	</script>
</body>
</html>`
	return html
}

func newInvoiceFormHTML(client *store.Client) string {
	html := `<!DOCTYPE html>
<html>
<head>
	<title>New Invoice</title>
	<style>
		* { margin: 0; padding: 0; box-sizing: border-box; }
		body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif; background: #f8f9fa; color: #333; }
		.header { background: white; padding: 20px 40px; border-bottom: 1px solid #e0e0e0; }
		.container { max-width: 800px; margin: 40px auto; padding: 0 20px; }
		.card { background: white; padding: 30px; border-radius: 8px; box-shadow: 0 1px 3px rgba(0,0,0,0.1); }
		.form-group { margin-bottom: 15px; }
		.form-group label { display: block; margin-bottom: 5px; font-weight: 500; }
		.form-group input { width: 100%; padding: 8px; border: 1px solid #ddd; border-radius: 4px; }
		.btn { padding: 10px 20px; border: none; border-radius: 4px; cursor: pointer; }
		.btn-primary { background: #007bff; color: white; }
		.btn-success { background: #28a745; color: white; }
		.items-section { margin: 30px 0; }
		.item-row { display: grid; grid-template-columns: 2fr 1fr 1fr auto; gap: 10px; margin-bottom: 10px; align-items: end; }
		.item-row input { padding: 8px; border: 1px solid #ddd; border-radius: 4px; }
		.item-total { padding: 8px; background: #f8f9fa; border-radius: 4px; }
		.summary { background: #f8f9fa; padding: 20px; border-radius: 8px; margin-top: 20px; }
		.summary-row { display: flex; justify-content: space-between; margin-bottom: 10px; }
		.summary-row.total { font-weight: bold; font-size: 18px; border-top: 2px solid #ddd; padding-top: 10px; }
		.buttons { display: flex; gap: 10px; margin-top: 20px; }
	</style>
</head>
<body>
	<div class="header">
		<h1>Create Invoice for ` + client.Name + `</h1>
	</div>

	<div class="container">
		<div class="card">
			<form id="invoiceForm">
				<div class="form-group">
					<label>Due Date</label>
					<input type="date" name="due_date" required>
				</div>

				<div class="items-section">
					<h3>Line Items</h3>
					<div id="itemsContainer"></div>
					<button type="button" class="btn btn-primary" onclick="addItem()">Add Line Item</button>
				</div>

				<div class="summary">
					<div class="summary-row">
						<span>Subtotal</span>
						<span id="subtotal">$0.00</span>
					</div>
					<div class="summary-row total">
						<span>Total</span>
						<span id="total">$0.00</span>
					</div>
				</div>

				<div class="buttons">
					<button type="submit" class="btn btn-success">Create Invoice</button>
					<a href="/client/` + client.ID + `" class="btn" style="background: #6c757d; color: white; text-decoration: none;">Cancel</a>
				</div>
			</form>
		</div>
	</div>

	<script>
		var itemCount = 0;

		function addItem() {
			var container = document.getElementById('itemsContainer');
			var row = document.createElement('div');
			row.className = 'item-row';
			row.id = 'item-' + itemCount;
			row.innerHTML = '<input type="text" name="item_desc_' + itemCount + '" placeholder="Description" required>' +
				'<input type="number" name="item_qty_' + itemCount + '" value="1" min="1" class="qty" required>' +
				'<input type="number" name="item_rate_' + itemCount + '" step="0.01" placeholder="Rate" class="rate" required>' +
				'<button type="button" onclick="removeItem(' + itemCount + ')" class="btn" style="background: #dc3545; color: white;">Remove</button>';
			container.appendChild(row);

			row.querySelectorAll('.qty, .rate').forEach(el => {
				el.addEventListener('change', calculateTotal);
				el.addEventListener('input', calculateTotal);
			});

			itemCount++;
		}

		function removeItem(id) {
			var item = document.getElementById('item-' + id);
			if (item) item.remove();
			calculateTotal();
		}

		function calculateTotal() {
			var total = 0;
			document.querySelectorAll('.item-row').forEach(row => {
				var qty = parseInt(row.querySelector('.qty').value) || 0;
				var rate = parseFloat(row.querySelector('.rate').value) || 0;
				total += qty * rate;
			});
			document.getElementById('subtotal').textContent = '$' + total.toFixed(2);
			document.getElementById('total').textContent = '$' + total.toFixed(2);
		}

		addItem();

		document.getElementById('invoiceForm').addEventListener('submit', function(e) {
			e.preventDefault();
			var form = new FormData(this);
			form.append('client_id', '` + client.ID + `');
			fetch('/api/invoices', {
				method: 'POST',
				body: form
			}).then(r => r.json()).then(d => {
				alert('Invoice created');
				window.location.href = '/client/` + client.ID + `';
			}).catch(e => alert('Error: ' + e));
		});
	</script>
</body>
</html>`
	return html
}

func invoicePrintHTML(invoice *store.Invoice, client *store.Client) string {
	total := 0.0
	for _, item := range invoice.Items {
		total += item.Amount
	}

	itemsHTML := ""
	for _, item := range invoice.Items {
		itemsHTML += `<tr>
			<td>` + item.Description + `</td>
			<td style="text-align: right;">` + strconv.Itoa(item.Quantity) + `</td>
			<td style="text-align: right;">$` + fmt.Sprintf("%.2f", item.Rate) + `</td>
			<td style="text-align: right;">$` + fmt.Sprintf("%.2f", item.Amount) + `</td>
		</tr>`
	}

	html := `<!DOCTYPE html>
<html>
<head>
	<title>Invoice ` + invoice.ID + `</title>
	<style>
		body { font-family: Arial, sans-serif; max-width: 900px; margin: 20px auto; padding: 20px; background: white; }
		.invoice-header { display: flex; justify-content: space-between; margin-bottom: 40px; }
		.company { font-size: 24px; font-weight: bold; }
		.invoice-title { font-size: 32px; font-weight: bold; color: #333; }
		.details { display: grid; grid-template-columns: 1fr 1fr; gap: 40px; margin-bottom: 40px; }
		.detail-section h4 { margin-bottom: 10px; font-weight: bold; }
		.detail-section p { margin: 5px 0; }
		table { width: 100%; border-collapse: collapse; margin-bottom: 20px; }
		th { background: #f8f9fa; padding: 10px; text-align: left; border-bottom: 2px solid #333; }
		td { padding: 10px; border-bottom: 1px solid #eee; }
		.total-section { display: flex; justify-content: flex-end; margin-top: 30px; }
		.total-box { width: 300px; }
		.total-row { display: flex; justify-content: space-between; padding: 5px 0; }
		.total-row.final { font-weight: bold; font-size: 18px; border-top: 2px solid #333; padding-top: 10px; }
		@media print {
			body { margin: 0; padding: 0; }
			.print-btn { display: none; }
		}
		.print-btn { padding: 10px 20px; background: #007bff; color: white; border: none; border-radius: 4px; cursor: pointer; margin-bottom: 20px; }
	</style>
</head>
<body>
	<button class="print-btn" onclick="window.print()">Print Invoice</button>

	<div class="invoice-header">
		<div class="company">My Company</div>
		<div class="invoice-title">INVOICE</div>
	</div>

	<div class="details">
		<div class="detail-section">
			<h4>Bill To:</h4>
			<p><strong>` + client.Name + `</strong></p>
			<p>` + client.Company + `</p>
			<p>` + client.Email + `</p>
			<p>` + client.Phone + `</p>
			<p>` + client.Address + `</p>
		</div>
		<div class="detail-section">
			<h4>Invoice Details:</h4>
			<p><strong>Invoice #:</strong> ` + invoice.ID + `</p>
			<p><strong>Date:</strong> ` + invoice.CreatedAt.Format("2006-01-02") + `</p>
			<p><strong>Due Date:</strong> ` + invoice.DueDate.Format("2006-01-02") + `</p>
			<p><strong>Status:</strong> ` + invoice.Status + `</p>
		</div>
	</div>

	<table>
		<thead>
			<tr>
				<th>Description</th>
				<th style="text-align: right;">Quantity</th>
				<th style="text-align: right;">Rate</th>
				<th style="text-align: right;">Amount</th>
			</tr>
		</thead>
		<tbody>
			` + itemsHTML + `
		</tbody>
	</table>

	<div class="total-section">
		<div class="total-box">
			<div class="total-row">
				<span>Subtotal:</span>
				<span>$` + fmt.Sprintf("%.2f", total) + `</span>
			</div>
			<div class="total-row final">
				<span>Total:</span>
				<span>$` + fmt.Sprintf("%.2f", total) + `</span>
			</div>
		</div>
	</div>
</body>
</html>`
	return html
}
