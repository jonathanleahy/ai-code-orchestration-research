package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"app/store"
	"github.com/google/uuid"
)

var s store.Store

func init() {
	s = store.NewStore()
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", handleHealth)
	mux.HandleFunc("/", handleDashboard)
	mux.HandleFunc("/api/clients", handleClients)
	mux.HandleFunc("/api/clients/", handleClientByID)
	mux.HandleFunc("/api/invoices", handleInvoices)
	mux.HandleFunc("/api/invoices/", handleInvoiceByID)
	mux.HandleFunc("/api/comments", handleComments)
	mux.HandleFunc("/api/comments/", handleCommentByID)
	mux.HandleFunc("/api/history", handleHistory)

	log.Println("Server starting at :8080")
	log.Fatal(http.ListenAndServe(":8080", corsMiddleware(mux)))
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
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
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	clients, _ := s.GetClients()
	invoices, _ := s.GetInvoices()
	comments, _ := s.GetComments("")
	history, _ := s.GetHistory("")

	clientsJSON, _ := json.Marshal(clients)
	invoicesJSON, _ := json.Marshal(invoices)
	commentsJSON, _ := json.Marshal(comments)
	historyJSON, _ := json.Marshal(history)

	html := `<!DOCTYPE html>
<html>
<head>
	<title>CRM Dashboard</title>
	<style>
		body { font-family: sans-serif; margin: 20px; background: #f5f5f5; }
		.container { max-width: 1200px; margin: 0 auto; }
		h1 { color: #333; }
		.section { background: white; padding: 20px; margin-bottom: 20px; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
		h2 { color: #0066cc; border-bottom: 2px solid #0066cc; padding-bottom: 10px; }
		.form-group { margin-bottom: 15px; }
		label { display: block; margin-bottom: 5px; font-weight: bold; color: #333; }
		input, textarea, select { width: 100%; padding: 8px; border: 1px solid #ddd; border-radius: 4px; box-sizing: border-box; }
		textarea { resize: vertical; min-height: 80px; }
		button { background: #0066cc; color: white; padding: 10px 20px; border: none; border-radius: 4px; cursor: pointer; font-size: 14px; }
		button:hover { background: #0052a3; }
		.list-item { padding: 10px; border: 1px solid #ddd; margin-bottom: 10px; border-radius: 4px; display: flex; justify-content: space-between; align-items: center; }
		.list-item-content { flex: 1; }
		.list-item button { margin-left: 10px; padding: 5px 10px; font-size: 12px; }
		.grid { display: grid; grid-template-columns: 1fr 1fr; gap: 20px; }
		.delete-btn { background: #cc0000; }
		.delete-btn:hover { background: #990000; }
		.error { color: red; padding: 10px; background: #ffe6e6; border-radius: 4px; margin-bottom: 10px; }
		.success { color: green; padding: 10px; background: #e6ffe6; border-radius: 4px; margin-bottom: 10px; }
	</style>
</head>
<body>
	<div class="container">
		<h1>CRM Dashboard</h1>

		<div class="grid">
			<div class="section">
				<h2>Clients</h2>
				<div id="clientMessage"></div>
				<form id="clientForm">
					<div class="form-group">
						<label>Name:</label>
						<input type="text" name="name" required>
					</div>
					<div class="form-group">
						<label>Email:</label>
						<input type="email" name="email" required>
					</div>
					<div class="form-group">
						<label>Phone:</label>
						<input type="text" name="phone">
					</div>
					<div class="form-group">
						<label>Street:</label>
						<input type="text" name="street">
					</div>
					<div class="form-group">
						<label>City:</label>
						<input type="text" name="city">
					</div>
					<div class="form-group">
						<label>State:</label>
						<input type="text" name="state">
					</div>
					<div class="form-group">
						<label>Zip Code:</label>
						<input type="text" name="zip_code">
					</div>
					<div class="form-group">
						<label>Country:</label>
						<input type="text" name="country">
					</div>
					<button type="submit">Create Client</button>
				</form>
				<div id="clientsList"></div>
			</div>

			<div class="section">
				<h2>Invoices</h2>
				<div id="invoiceMessage"></div>
				<form id="invoiceForm">
					<div class="form-group">
						<label>Client ID:</label>
						<input type="text" name="client_id" required>
					</div>
					<div class="form-group">
						<label>Item Description:</label>
						<input type="text" name="item_description" required>
					</div>
					<div class="form-group">
						<label>Item Quantity:</label>
						<input type="number" name="item_quantity" value="1" required>
					</div>
					<div class="form-group">
						<label>Item Price:</label>
						<input type="number" name="item_price" step="0.01" required>
					</div>
					<div class="form-group">
						<label>Status:</label>
						<select name="status">
							<option>draft</option>
							<option>sent</option>
							<option>paid</option>
						</select>
					</div>
					<button type="submit">Create Invoice</button>
				</form>
				<div id="invoicesList"></div>
			</div>
		</div>

		<div class="section">
			<h2>Comments</h2>
			<div id="commentMessage"></div>
			<form id="commentForm">
				<div class="form-group">
					<label>Client ID:</label>
					<input type="text" name="client_id" required>
				</div>
				<div class="form-group">
					<label>Comment:</label>
					<textarea name="content" required></textarea>
				</div>
				<button type="submit">Add Comment</button>
			</form>
			<div id="commentsList"></div>
		</div>

		<div class="section">
			<h2>History</h2>
			<div id="historyMessage"></div>
			<form id="historyForm">
				<div class="form-group">
					<label>Client ID:</label>
					<input type="text" name="client_id" required>
				</div>
				<div class="form-group">
					<label>Note:</label>
					<textarea name="note" required></textarea>
				</div>
				<button type="submit">Add History Entry</button>
			</form>
			<div id="historyList"></div>
		</div>
	</div>

	<script>
		var baseURL = window.location.origin + '/api';
		var initialData = {
			clients: ` + string(clientsJSON) + `,
			invoices: ` + string(invoicesJSON) + `,
			comments: ` + string(commentsJSON) + `,
			history: ` + string(historyJSON) + `
		};

		function showMessage(elementId, message, isError) {
			var el = document.getElementById(elementId);
			el.className = isError ? 'error' : 'success';
			el.textContent = message;
			setTimeout(function() { el.textContent = ''; el.className = ''; }, 3000);
		}

		function loadClients() {
			fetch(baseURL + '/clients').then(function(r) { return r.json(); }).then(function(clients) {
				var html = '';
				if (clients && clients.length > 0) {
					for (var i = 0; i < clients.length; i++) {
						var c = clients[i];
						html = html + '<div class="list-item"><div class="list-item-content"><strong>' + (c.name || 'N/A') + '</strong> - ' + (c.email || 'N/A') + '</div>';
						html = html + '<button class="delete-btn" onclick="deleteClient(\'' + c.id + '\')">Delete</button></div>';
					}
				}
				document.getElementById('clientsList').innerHTML = html;
			}).catch(function(e) { console.error(e); });
		}

		function loadInvoices() {
			fetch(baseURL + '/invoices').then(function(r) { return r.json(); }).then(function(invoices) {
				var html = '';
				if (invoices && invoices.length > 0) {
					for (var i = 0; i < invoices.length; i++) {
						var inv = invoices[i];
						html = html + '<div class="list-item"><div class="list-item-content"><strong>' + inv.id + '</strong> - Client: ' + inv.client_id + ' - ' + inv.status + ' - $' + inv.total_amount + '</div>';
						html = html + '<button class="delete-btn" onclick="deleteInvoice(\'' + inv.id + '\')">Delete</button></div>';
					}
				}
				document.getElementById('invoicesList').innerHTML = html;
			}).catch(function(e) { console.error(e); });
		}

		function loadComments() {
			fetch(baseURL + '/comments').then(function(r) { return r.json(); }).then(function(comments) {
				var html = '';
				if (comments && comments.length > 0) {
					for (var i = 0; i < comments.length; i++) {
						var cm = comments[i];
						html = html + '<div class="list-item"><div class="list-item-content"><strong>' + cm.id + '</strong> - Client: ' + cm.client_id + ' - ' + cm.content + '</div>';
						html = html + '<button class="delete-btn" onclick="deleteComment(\'' + cm.id + '\')">Delete</button></div>';
					}
				}
				document.getElementById('commentsList').innerHTML = html;
			}).catch(function(e) { console.error(e); });
		}

		function loadHistory() {
			fetch(baseURL + '/history').then(function(r) { return r.json(); }).then(function(history) {
				var html = '';
				if (history && history.length > 0) {
					for (var i = 0; i < history.length; i++) {
						var h = history[i];
						html = html + '<div class="list-item"><div class="list-item-content"><strong>' + h.id + '</strong> - Client: ' + h.client_id + ' - ' + h.note + '</div></div>';
					}
				}
				document.getElementById('historyList').innerHTML = html;
			}).catch(function(e) { console.error(e); });
		}

		function deleteClient(id) {
			fetch(baseURL + '/clients/' + id, { method: 'DELETE' }).then(function() {
				showMessage('clientMessage', 'Client deleted', false);
				loadClients();
			}).catch(function(e) { showMessage('clientMessage', 'Error: ' + e.message, true); });
		}

		function deleteInvoice(id) {
			fetch(baseURL + '/invoices/' + id, { method: 'DELETE' }).then(function() {
				showMessage('invoiceMessage', 'Invoice deleted', false);
				loadInvoices();
			}).catch(function(e) { showMessage('invoiceMessage', 'Error: ' + e.message, true); });
		}

		function deleteComment(id) {
			fetch(baseURL + '/comments/' + id, { method: 'DELETE' }).then(function() {
				showMessage('commentMessage', 'Comment deleted', false);
				loadComments();
			}).catch(function(e) { showMessage('commentMessage', 'Error: ' + e.message, true); });
		}

		document.getElementById('clientForm').addEventListener('submit', function(e) {
			e.preventDefault();
			var form = new FormData(this);
			var clientData = {
				id: 'client-' + Date.now(),
				name: form.get('name'),
				email: form.get('email'),
				phone: form.get('phone'),
				address: {
					street: form.get('street'),
					city: form.get('city'),
					state: form.get('state'),
					zip_code: form.get('zip_code'),
					country: form.get('country')
				},
				created_at: new Date().toISOString(),
				updated_at: new Date().toISOString()
			};
			fetch(baseURL + '/clients', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify(clientData)
			}).then(function() {
				showMessage('clientMessage', 'Client created', false);
				document.getElementById('clientForm').reset();
				loadClients();
			}).catch(function(e) { showMessage('clientMessage', 'Error: ' + e.message, true); });
		});

		document.getElementById('invoiceForm').addEventListener('submit', function(e) {
			e.preventDefault();
			var form = new FormData(this);
			var qty = parseInt(form.get('item_quantity'));
			var price = parseFloat(form.get('item_price'));
			var invoiceData = {
				id: 'inv-' + Date.now(),
				client_id: form.get('client_id'),
				items: [{
					description: form.get('item_description'),
					quantity: qty,
					price: price,
					total: qty * price
				}],
				total_amount: qty * price,
				issue_date: new Date().toISOString(),
				due_date: new Date(Date.now() + 30*24*60*60*1000).toISOString(),
				status: form.get('status'),
				created_at: new Date().toISOString(),
				updated_at: new Date().toISOString()
			};
			fetch(baseURL + '/invoices', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify(invoiceData)
			}).then(function() {
				showMessage('invoiceMessage', 'Invoice created', false);
				document.getElementById('invoiceForm').reset();
				loadInvoices();
			}).catch(function(e) { showMessage('invoiceMessage', 'Error: ' + e.message, true); });
		});

		document.getElementById('commentForm').addEventListener('submit', function(e) {
			e.preventDefault();
			var form = new FormData(this);
			var commentData = {
				id: 'comment-' + Date.now(),
				client_id: form.get('client_id'),
				content: form.get('content'),
				created_at: new Date().toISOString()
			};
			fetch(baseURL + '/comments', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify(commentData)
			}).then(function() {
				showMessage('commentMessage', 'Comment created', false);
				document.getElementById('commentForm').reset();
				loadComments();
			}).catch(function(e) { showMessage('commentMessage', 'Error: ' + e.message, true); });
		});

		document.getElementById('historyForm').addEventListener('submit', function(e) {
			e.preventDefault();
			var form = new FormData(this);
			var historyData = {
				id: 'hist-' + Date.now(),
				client_id: form.get('client_id'),
				note: form.get('note'),
				created_at: new Date().toISOString()
			};
			fetch(baseURL + '/history', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify(historyData)
			}).then(function() {
				showMessage('historyMessage', 'History entry created', false);
				document.getElementById('historyForm').reset();
				loadHistory();
			}).catch(function(e) { showMessage('historyMessage', 'Error: ' + e.message, true); });
		});

		loadClients();
		loadInvoices();
		loadComments();
		loadHistory();
	</script>
</body>
</html>`

	w.Write([]byte(html))
}

func handleClients(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodPost:
		var client store.Client
		if err := json.NewDecoder(r.Body).Decode(&client); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		if client.ID == "" {
			client.ID = uuid.New().String()
		}
		if client.CreatedAt.IsZero() {
			client.CreatedAt = time.Now()
		}
		if client.UpdatedAt.IsZero() {
			client.UpdatedAt = time.Now()
		}
		if err := s.CreateClient(&client); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(client)

	case http.MethodGet:
		clients, err := s.GetClients()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		if clients == nil {
			clients = []*store.Client{}
		}
		json.NewEncoder(w).Encode(clients)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func handleClientByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := strings.TrimPrefix(r.URL.Path, "/api/clients/")
	if id == "" || id == "/api/clients" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid client id"})
		return
	}

	switch r.Method {
	case http.MethodGet:
		client, err := s.GetClient(id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		if client == nil {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "client not found"})
			return
		}
		json.NewEncoder(w).Encode(client)

	case http.MethodPatch:
		var client store.Client
		if err := json.NewDecoder(r.Body).Decode(&client); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		if err := s.UpdateClient(id, &client); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		client.ID = id
		json.NewEncoder(w).Encode(client)

	case http.MethodDelete:
		if err := s.DeleteClient(id); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		w.WriteHeader(http.StatusNoContent)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func handleInvoices(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodPost:
		var invoice store.Invoice
		if err := json.NewDecoder(r.Body).Decode(&invoice); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		if invoice.ID == "" {
			invoice.ID = uuid.New().String()
		}
		if invoice.CreatedAt.IsZero() {
			invoice.CreatedAt = time.Now()
		}
		if invoice.UpdatedAt.IsZero() {
			invoice.UpdatedAt = time.Now()
		}
		if err := s.CreateInvoice(&invoice); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(invoice)

	case http.MethodGet:
		invoices, err := s.GetInvoices()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		if invoices == nil {
			invoices = []*store.Invoice{}
		}
		json.NewEncoder(w).Encode(invoices)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func handleInvoiceByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := strings.TrimPrefix(r.URL.Path, "/api/invoices/")
	if id == "" || id == "/api/invoices" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid invoice id"})
		return
	}

	switch r.Method {
	case http.MethodGet:
		invoice, err := s.GetInvoice(id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		if invoice == nil {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "invoice not found"})
			return
		}
		json.NewEncoder(w).Encode(invoice)

	case http.MethodPatch:
		var invoice store.Invoice
		if err := json.NewDecoder(r.Body).Decode(&invoice); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		if err := s.UpdateInvoice(id, &invoice); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		invoice.ID = id
		json.NewEncoder(w).Encode(invoice)

	case http.MethodDelete:
		if err := s.DeleteInvoice(id); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		w.WriteHeader(http.StatusNoContent)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func handleComments(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodPost:
		var comment store.Comment
		if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		if comment.ID == "" {
			comment.ID = uuid.New().String()
		}
		if comment.CreatedAt.IsZero() {
			comment.CreatedAt = time.Now()
		}
		if err := s.CreateComment(&comment); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(comment)

	case http.MethodGet:
		comments, err := s.GetComments("")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		if comments == nil {
			comments = []*store.Comment{}
		}
		json.NewEncoder(w).Encode(comments)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func handleCommentByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := strings.TrimPrefix(r.URL.Path, "/api/comments/")
	if id == "" || id == "/api/comments" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid comment id"})
		return
	}

	switch r.Method {
	case http.MethodDelete:
		if err := s.DeleteComment(id); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		w.WriteHeader(http.StatusNoContent)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func handleHistory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodPost:
		var entry store.HistoryEntry
		if err := json.NewDecoder(r.Body).Decode(&entry); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		if entry.ID == "" {
			entry.ID = uuid.New().String()
		}
		if entry.CreatedAt.IsZero() {
			entry.CreatedAt = time.Now()
		}
		if err := s.CreateHistory(&entry); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(entry)

	case http.MethodGet:
		history, err := s.GetHistory("")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		if history == nil {
			history = []*store.HistoryEntry{}
		}
		json.NewEncoder(w).Encode(history)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
