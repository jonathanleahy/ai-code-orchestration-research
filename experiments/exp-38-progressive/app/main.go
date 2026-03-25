package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"app/store"
)

var s *store.Store

func init() {
	s = store.NewStore()
	// Seed with sample data
	s.AddClient("John Doe", "john@example.com", "555-0101")
	s.AddClient("Jane Smith", "jane@example.com", "555-0102")
}

func main() {
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

	log.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func handleDashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	clients := s.ListClients()

	html := `<!DOCTYPE html>
<html>
<head>
	<meta charset="utf-8">
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<title>CRM Dashboard</title>
	<style>
		* { margin: 0; padding: 0; box-sizing: border-box; }
		body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif; background: #f8f9fa; color: #333; }
		.container { max-width: 1200px; margin: 0 auto; padding: 20px; }
		header { margin-bottom: 30px; }
		h1 { font-size: 28px; margin-bottom: 10px; }
		.section { background: white; border-radius: 8px; padding: 20px; margin-bottom: 20px; box-shadow: 0 1px 3px rgba(0,0,0,0.1); }
		table { width: 100%; border-collapse: collapse; }
		th { text-align: left; padding: 12px; border-bottom: 1px solid #ddd; font-weight: 600; font-size: 14px; background: #f0f1f3; }
		td { padding: 12px; border-bottom: 1px solid #eee; }
		tr:hover { background: #f8f9fa; }
		a { color: #0066cc; text-decoration: none; }
		a:hover { text-decoration: underline; }
		.form-group { margin-bottom: 15px; }
		label { display: block; font-weight: 600; margin-bottom: 5px; font-size: 14px; }
		input { padding: 8px 12px; border: 1px solid #ddd; border-radius: 4px; font-size: 14px; width: 100%; }
		input:focus { outline: none; border-color: #0066cc; box-shadow: 0 0 0 3px rgba(0, 102, 204, 0.1); }
		button { padding: 10px 20px; background: #0066cc; color: white; border: none; border-radius: 4px; cursor: pointer; font-weight: 600; font-size: 14px; }
		button:hover { background: #0052a3; }
		.form-row { display: grid; grid-template-columns: 1fr 1fr 1fr; gap: 15px; }
		@media (max-width: 768px) { .form-row { grid-template-columns: 1fr; } }
		.empty { text-align: center; color: #999; padding: 40px 20px; }
	</style>
</head>
<body>
	<div class="container">
		<header>
			<h1>CRM Dashboard</h1>
			<p>Manage your clients</p>
		</header>

		<div class="section">
			<h2 style="margin-bottom: 20px;">Add New Client</h2>
			<form id="addClientForm">
				<div class="form-row">
					<div class="form-group">
						<label for="name">Name</label>
						<input type="text" id="name" name="name" required>
					</div>
					<div class="form-group">
						<label for="email">Email</label>
						<input type="email" id="email" name="email" required>
					</div>
					<div class="form-group">
						<label for="phone">Phone</label>
						<input type="text" id="phone" name="phone" required>
					</div>
				</div>
				<button type="submit">Add Client</button>
			</form>
		</div>

		<div class="section">
			<h2 style="margin-bottom: 20px;">Clients List</h2>
			<div style="margin-bottom: 20px;">
				<input type="text" id="searchBox" placeholder="Search clients by name or email..." style="width: 100%; max-width: 400px;">
			</div>
`

	if len(clients) == 0 {
		html += `<div class="empty">No clients yet. Add one to get started!</div>`
	} else {
		html += `<table id="clientsTable">
				<thead>
					<tr>
						<th>Name</th>
						<th>Email</th>
						<th>Phone</th>
						<th>Created</th>
						<th>Action</th>
					</tr>
				</thead>
				<tbody>`
		for _, c := range clients {
			created := c.CreatedAt.Format("Jan 02, 2006")
			html += `<tr>
					<td>` + c.Name + `</td>
					<td>` + c.Email + `</td>
					<td>` + c.Phone + `</td>
					<td>` + created + `</td>
					<td><a href="/client/` + c.ID + `">View</a></td>
				</tr>`
		}
		html += `</tbody>
			</table>`
	}

	html += `
		</div>
	</div>

	<script>
		var form = document.getElementById('addClientForm');
		form.addEventListener('submit', function(e) {
			e.preventDefault();

			var formData = new FormData(form);
			var name = document.getElementById('name').value;
			var email = document.getElementById('email').value;
			var phone = document.getElementById('phone').value;

			fetch('/api/clients', {
				method: 'POST',
				body: formData
			})
			.then(function(response) { return response.json(); })
			.then(function(data) {
				if (data.id) {
					alert('Client added: ' + data.name);
					form.reset();
					setTimeout(function() { location.reload(); }, 500);
				}
			})
			.catch(function(err) { console.error('Error:', err); });
		});

		var searchBox = document.getElementById('searchBox');
		if (searchBox) {
			searchBox.addEventListener('input', function() {
				var query = this.value.trim();
				if (query === '') {
					location.reload();
					return;
				}
				fetch('/api/clients?search=' + encodeURIComponent(query))
					.then(function(response) { return response.json(); })
					.then(function(data) {
						var table = document.getElementById('clientsTable');
						if (!table) return;
						var tbody = table.querySelector('tbody');
						tbody.innerHTML = '';
						if (data && data.length > 0) {
							data.forEach(function(c) {
								var created = new Date(c.created_at).toLocaleDateString('en-US', { year: 'numeric', month: 'short', day: '2-digit' });
								var row = '<tr><td>' + c.name + '</td><td>' + c.email + '</td><td>' + c.phone + '</td><td>' + created + '</td><td><a href="/client/' + c.id + '">View</a></td></tr>';
								tbody.innerHTML += row;
							});
						} else {
							tbody.innerHTML = '<tr><td colspan="5" style="text-align: center; color: #999;">No clients found</td></tr>';
						}
					})
					.catch(function(err) { console.error('Error:', err); });
			});
		}
	</script>
</body>
</html>`

	fmt.Fprint(w, html)
}

func handleClientDetail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	id := r.PathValue("id")

	client := s.GetClient(id)
	if client == nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, `<html><body><p>Client not found</p><a href="/">Back to Dashboard</a></body></html>`)
		return
	}

	activities := s.GetActivities(id)
	invoices := s.ListInvoices(id)
	created := client.CreatedAt.Format("January 02, 2006 at 3:04 PM")

	html := `<!DOCTYPE html>
<html>
<head>
	<meta charset="utf-8">
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<title>` + client.Name + ` - CRM</title>
	<style>
		* { margin: 0; padding: 0; box-sizing: border-box; }
		body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif; background: #f8f9fa; color: #333; }
		.container { max-width: 800px; margin: 0 auto; padding: 20px; }
		header { margin-bottom: 30px; }
		h1 { font-size: 28px; margin-bottom: 10px; }
		.back-btn { display: inline-block; color: #0066cc; text-decoration: none; margin-bottom: 20px; }
		.back-btn:hover { text-decoration: underline; }
		.tabs { display: flex; gap: 0; margin-bottom: 20px; border-bottom: 1px solid #ddd; }
		.tab-btn { padding: 12px 20px; background: transparent; border: none; cursor: pointer; font-weight: 600; font-size: 14px; color: #666; border-bottom: 3px solid transparent; }
		.tab-btn.active { color: #0066cc; border-bottom-color: #0066cc; }
		.tab-btn:hover { color: #333; }
		.section { background: white; border-radius: 8px; padding: 30px; box-shadow: 0 1px 3px rgba(0,0,0,0.1); margin-bottom: 20px; display: none; }
		.section.active { display: block; }
		.field { margin-bottom: 20px; }
		.field-label { font-weight: 600; color: #666; font-size: 12px; text-transform: uppercase; margin-bottom: 5px; }
		.field-value { font-size: 16px; }
		input[type="text"], input[type="email"], input[type="textarea"], select { padding: 8px 12px; border: 1px solid #ddd; border-radius: 4px; font-size: 14px; width: 100%; }
		input:focus, select:focus, textarea:focus { outline: none; border-color: #0066cc; box-shadow: 0 0 0 3px rgba(0, 102, 204, 0.1); }
		textarea { resize: vertical; min-height: 100px; }
		.button-group { display: flex; gap: 10px; margin-top: 20px; }
		button { padding: 10px 20px; border: none; border-radius: 4px; cursor: pointer; font-weight: 600; font-size: 14px; }
		.btn-primary { background: #0066cc; color: white; }
		.btn-primary:hover { background: #0052a3; }
		.btn-danger { background: #dc3545; color: white; }
		.btn-danger:hover { background: #c82333; }
		.readonly-field { font-size: 16px; padding: 8px 0; }
		.timeline { display: flex; flex-direction: column; gap: 20px; }
		.timeline-item { display: flex; gap: 15px; }
		.timeline-dot { width: 10px; height: 10px; background: #0066cc; border-radius: 50%; margin-top: 6px; flex-shrink: 0; }
		.timeline-content { flex: 1; }
		.timeline-type { display: inline-block; background: #e7f0ff; color: #0066cc; padding: 3px 8px; border-radius: 3px; font-size: 12px; font-weight: 600; margin-bottom: 5px; }
		.timeline-description { font-size: 14px; color: #333; }
		.timeline-date { font-size: 12px; color: #999; margin-top: 5px; }
		.empty-timeline { text-align: center; color: #999; padding: 20px; }
	</style>
</head>
<body>
	<div class="container">
		<a href="/" class="back-btn">← Back to Dashboard</a>
		<header>
			<h1>` + client.Name + `</h1>
		</header>

		<div class="tabs">
			<button class="tab-btn active" onclick="showTab('details')">Details</button>
			<button class="tab-btn" onclick="showTab('activity')">Activity</button>
			<button class="tab-btn" onclick="showTab('invoices')">Invoices</button>
		</div>

		<div id="details" class="section active">
			<form id="editForm">
				<div class="field">
					<div class="field-label">Name</div>
					<input type="text" name="name" value="` + client.Name + `" required>
				</div>
				<div class="field">
					<div class="field-label">Email</div>
					<input type="email" name="email" value="` + client.Email + `" required>
				</div>
				<div class="field">
					<div class="field-label">Phone</div>
					<input type="text" name="phone" value="` + client.Phone + `" required>
				</div>
				<div class="field">
					<div class="field-label">Created</div>
					<div class="readonly-field">` + created + `</div>
				</div>
				<div class="field">
					<div class="field-label">ID</div>
					<div class="readonly-field"><code>` + client.ID + `</code></div>
				</div>
				<div class="button-group">
					<button type="submit" class="btn-primary">Save Changes</button>
					<button type="button" class="btn-danger" onclick="deleteClient()">Delete</button>
				</div>
			</form>
		</div>

		<div id="activity" class="section">
			<h2 style="margin-bottom: 20px;">Log Activity</h2>
			<form id="activityForm">
				<div class="field">
					<div class="field-label">Type</div>
					<select name="type" required>
						<option value="">-- Select Type --</option>
						<option value="call">Call</option>
						<option value="email">Email</option>
						<option value="meeting">Meeting</option>
						<option value="note">Note</option>
					</select>
				</div>
				<div class="field">
					<div class="field-label">Description</div>
					<textarea name="description" required></textarea>
				</div>
				<button type="submit" class="btn-primary">Log Activity</button>
			</form>

			<h2 style="margin-top: 40px; margin-bottom: 20px;">Activity Timeline</h2>
			<div class="timeline" id="activityTimeline">
`

	if len(activities) == 0 {
		html += `<div class="empty-timeline">No activities logged yet.</div>`
	} else {
		for i := len(activities) - 1; i >= 0; i-- {
			a := activities[i]
			timeStr := a.CreatedAt.Format("Jan 02, 2006 at 3:04 PM")
			html += `<div class="timeline-item">
					<div class="timeline-dot"></div>
					<div class="timeline-content">
						<div class="timeline-type">` + strings.ToUpper(a.Type) + `</div>
						<div class="timeline-description">` + a.Description + `</div>
						<div class="timeline-date">` + timeStr + `</div>
					</div>
				</div>`
		}
	}

	html += `
			</div>
		</div>

		<div id="invoices" class="section">
			<h2 style="margin-bottom: 20px;">Create Invoice</h2>
			<form id="invoiceForm">
				<div class="field">
					<div class="field-label">Invoice Number</div>
					<input type="text" name="number" placeholder="e.g. INV-001" required>
				</div>
				<div class="field">
					<div class="field-label">Amount</div>
					<input type="number" name="amount" placeholder="0.00" step="0.01" min="0" required>
				</div>
				<div class="field">
					<div class="field-label">Description</div>
					<textarea name="description" required></textarea>
				</div>
				<div class="field">
					<div class="field-label">Due Date</div>
					<input type="date" name="due_date" required>
				</div>
				<button type="submit" class="btn-primary">Create Invoice</button>
			</form>

			<h2 style="margin-top: 40px; margin-bottom: 20px;">Invoices</h2>`

	if len(invoices) == 0 {
		html += `<div class="empty">No invoices yet.</div>`
	} else {
		html += `<table>
				<thead>
					<tr>
						<th>Number</th>
						<th>Amount</th>
						<th>Description</th>
						<th>Due Date</th>
						<th>Status</th>
						<th>Action</th>
					</tr>
				</thead>
				<tbody>`
		for _, inv := range invoices {
			dueStr := inv.DueDate.Format("Jan 02, 2006")
			statusColor := "background: #e7f0ff; color: #0066cc;"
			if inv.Status == "sent" {
				statusColor = "background: #fff3cd; color: #856404;"
			} else if inv.Status == "paid" {
				statusColor = "background: #d4edda; color: #155724;"
			}
			html += `<tr>
					<td>` + inv.Number + `</td>
					<td>$` + fmt.Sprintf("%.2f", inv.Amount) + `</td>
					<td>` + inv.Description + `</td>
					<td>` + dueStr + `</td>
					<td><span style="padding: 4px 8px; border-radius: 3px; font-size: 12px; font-weight: 600; ` + statusColor + `">` + strings.ToUpper(inv.Status) + `</span></td>
					<td style="display: flex; gap: 10px;">
						<a href="/invoice/` + inv.ID + `/print" target="_blank" style="color: #0066cc; text-decoration: none;">Print</a>`
			if inv.Status != "paid" {
				html += `<button type="button" class="btn-primary" onclick="markInvoicePaid('` + inv.ID + `')" style="padding: 6px 12px; font-size: 12px;">Mark Paid</button>`
			}
			html += `</td>
				</tr>`
		}
		html += `</tbody>
			</table>`
	}

	html += `
		</div>
	</div>

	<script>
		var clientId = '` + client.ID + `';

		function showTab(tabName) {
			var sections = document.querySelectorAll('.section');
			var buttons = document.querySelectorAll('.tab-btn');
			sections.forEach(function(s) { s.classList.remove('active'); });
			buttons.forEach(function(b) { b.classList.remove('active'); });

			document.getElementById(tabName).classList.add('active');
			event.target.classList.add('active');
		}

		document.getElementById('editForm').addEventListener('submit', function(e) {
			e.preventDefault();
			var name = document.querySelector('input[name="name"]').value;
			var email = document.querySelector('input[name="email"]').value;
			var phone = document.querySelector('input[name="phone"]').value;

			var formData = new FormData();
			formData.append('name', name);
			formData.append('email', email);
			formData.append('phone', phone);

			fetch('/api/clients/' + clientId + '/update', {
				method: 'POST',
				body: formData
			})
			.then(function(response) { return response.json(); })
			.then(function(data) {
				if (data.id) {
					alert('Client updated successfully');
					setTimeout(function() { location.reload(); }, 500);
				} else if (data.error) {
					alert('Error: ' + data.error);
				}
			})
			.catch(function(err) { console.error('Error:', err); });
		});

		document.getElementById('activityForm').addEventListener('submit', function(e) {
			e.preventDefault();
			var type = document.querySelector('select[name="type"]').value;
			var description = document.querySelector('textarea[name="description"]').value;

			var formData = new FormData();
			formData.append('type', type);
			formData.append('description', description);

			fetch('/api/clients/' + clientId + '/activities', {
				method: 'POST',
				body: formData
			})
			.then(function(response) { return response.json(); })
			.then(function(data) {
				if (data.id) {
					alert('Activity logged successfully');
					setTimeout(function() { location.reload(); }, 500);
				} else if (data.error) {
					alert('Error: ' + data.error);
				}
			})
			.catch(function(err) { console.error('Error:', err); });
		});

		document.getElementById('invoiceForm').addEventListener('submit', function(e) {
			e.preventDefault();
			var number = document.querySelector('input[name="number"]').value;
			var amount = document.querySelector('input[name="amount"]').value;
			var description = document.querySelector('textarea[name="description"]').value;
			var dueDate = document.querySelector('input[name="due_date"]').value;

			var formData = new FormData();
			formData.append('number', number);
			formData.append('amount', amount);
			formData.append('description', description);
			formData.append('due_date', dueDate);

			fetch('/api/clients/' + clientId + '/invoices', {
				method: 'POST',
				body: formData
			})
			.then(function(response) { return response.json(); })
			.then(function(data) {
				if (data.id) {
					alert('Invoice created successfully');
					setTimeout(function() { location.reload(); }, 500);
				} else if (data.error) {
					alert('Error: ' + data.error);
				}
			})
			.catch(function(err) { console.error('Error:', err); });
		});

		function markInvoicePaid(invoiceId) {
			fetch('/api/invoices/' + invoiceId + '/pay', {
				method: 'POST'
			})
			.then(function(response) { return response.json(); })
			.then(function(data) {
				if (data.id) {
					alert('Invoice marked as paid');
					setTimeout(function() { location.reload(); }, 500);
				} else if (data.error) {
					alert('Error: ' + data.error);
				}
			})
			.catch(function(err) { console.error('Error:', err); });
		}

		function deleteClient() {
			if (confirm('Are you sure you want to delete this client? This action cannot be undone.')) {
				var formData = new FormData();
				fetch('/api/clients/' + clientId + '/delete', {
					method: 'POST',
					body: formData
				})
				.then(function(response) { return response.json(); })
				.then(function(data) {
					if (data.success) {
						alert('Client deleted successfully');
						window.location.href = '/';
					} else if (data.error) {
						alert('Error: ' + data.error);
					}
				})
				.catch(function(err) { console.error('Error:', err); });
			}
		}
	</script>
</body>
</html>`

	fmt.Fprint(w, html)
}

func handleCreateClient(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	r.ParseMultipartForm(1024)

	name := strings.TrimSpace(r.FormValue("name"))
	email := strings.TrimSpace(r.FormValue("email"))
	phone := strings.TrimSpace(r.FormValue("phone"))

	if name == "" || email == "" || phone == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "missing fields"})
		return
	}

	client := s.AddClient(name, email, phone)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(client)
}

func handleListClients(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	query := strings.TrimSpace(r.URL.Query().Get("search"))
	var clients []*store.Client
	if query != "" {
		clients = s.SearchClients(query)
	} else {
		clients = s.ListClients()
	}

	if clients == nil {
		clients = []*store.Client{}
	}
	json.NewEncoder(w).Encode(clients)
}

func handleUpdateClient(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := r.PathValue("id")

	r.ParseMultipartForm(1024)

	name := strings.TrimSpace(r.FormValue("name"))
	email := strings.TrimSpace(r.FormValue("email"))
	phone := strings.TrimSpace(r.FormValue("phone"))

	if name == "" || email == "" || phone == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "missing fields"})
		return
	}

	client := s.UpdateClient(id, name, email, phone)
	if client == nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "client not found"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(client)
}

func handleDeleteClient(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := r.PathValue("id")

	if !s.DeleteClient(id) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "client not found"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

func handleCreateActivity(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := r.PathValue("id")

	if s.GetClient(id) == nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "client not found"})
		return
	}

	r.ParseMultipartForm(1024)

	actType := strings.TrimSpace(r.FormValue("type"))
	description := strings.TrimSpace(r.FormValue("description"))

	if actType == "" || description == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "missing fields"})
		return
	}

	activity := s.AddActivity(id, actType, description)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(activity)
}

func handleListActivities(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := r.PathValue("id")

	if s.GetClient(id) == nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "client not found"})
		return
	}

	activities := s.GetActivities(id)
	if activities == nil {
		activities = []*store.Activity{}
	}
	json.NewEncoder(w).Encode(activities)
}

func handleCreateInvoice(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := r.PathValue("id")

	if s.GetClient(id) == nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "client not found"})
		return
	}

	r.ParseMultipartForm(1024)

	number := strings.TrimSpace(r.FormValue("number"))
	description := strings.TrimSpace(r.FormValue("description"))
	amountStr := strings.TrimSpace(r.FormValue("amount"))
	dueStr := strings.TrimSpace(r.FormValue("due_date"))

	if number == "" || description == "" || amountStr == "" || dueStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "missing fields"})
		return
	}

	var amount float64
	_, err := fmt.Sscanf(amountStr, "%f", &amount)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid amount"})
		return
	}

	dueDate, err := time.Parse("2006-01-02", dueStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid due date format"})
		return
	}

	invoice := s.CreateInvoice(id, number, description, amount, dueDate)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(invoice)
}

func handleListInvoices(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := r.PathValue("id")

	if s.GetClient(id) == nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "client not found"})
		return
	}

	invoices := s.ListInvoices(id)
	if invoices == nil {
		invoices = []*store.Invoice{}
	}
	json.NewEncoder(w).Encode(invoices)
}

func handleMarkInvoicePaid(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := r.PathValue("id")

	invoice := s.GetInvoice(id)
	if invoice == nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "invoice not found"})
		return
	}

	invoice = s.MarkInvoicePaid(id)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(invoice)
}

func handlePrintInvoice(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	id := r.PathValue("id")

	invoice := s.GetInvoice(id)
	if invoice == nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, `<html><body><p>Invoice not found</p></body></html>`)
		return
	}

	client := s.GetClient(invoice.ClientID)
	if client == nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, `<html><body><p>Client not found</p></body></html>`)
		return
	}

	dueStr := invoice.DueDate.Format("January 02, 2006")
	createdStr := invoice.CreatedAt.Format("January 02, 2006")
	statusBadge := "DRAFT"
	if invoice.Status == "paid" {
		statusBadge = "PAID"
	} else if invoice.Status == "sent" {
		statusBadge = "SENT"
	}

	html := `<!DOCTYPE html>
<html>
<head>
	<meta charset="utf-8">
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<title>Invoice ` + invoice.Number + `</title>
	<style>
		* { margin: 0; padding: 0; box-sizing: border-box; }
		body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif; color: #333; }
		.container { max-width: 800px; margin: 0 auto; padding: 40px; }
		header { margin-bottom: 40px; }
		h1 { font-size: 32px; margin-bottom: 10px; }
		.invoice-info { display: grid; grid-template-columns: 1fr 1fr; gap: 40px; margin-bottom: 40px; }
		.info-section { }
		.info-label { font-weight: 600; font-size: 12px; text-transform: uppercase; color: #666; margin-bottom: 5px; }
		.info-value { font-size: 16px; line-height: 1.6; }
		.divider { border-top: 2px solid #ddd; margin: 40px 0; }
		table { width: 100%; border-collapse: collapse; margin-bottom: 40px; }
		th { text-align: left; padding: 12px; border-bottom: 2px solid #ddd; font-weight: 600; font-size: 14px; background: #f8f9fa; }
		td { padding: 12px; border-bottom: 1px solid #eee; }
		.amount-column { text-align: right; }
		.totals { display: flex; justify-content: flex-end; margin-bottom: 40px; }
		.totals-section { width: 300px; }
		.total-row { display: flex; justify-content: space-between; padding: 8px 0; font-size: 14px; border-bottom: 1px solid #eee; }
		.total-row.final { border-top: 2px solid #333; border-bottom: none; font-weight: 600; font-size: 18px; margin-top: 10px; }
		.status-badge { display: inline-block; padding: 6px 12px; border-radius: 4px; font-weight: 600; font-size: 12px; }
		.status-draft { background: #e7f0ff; color: #0066cc; }
		.status-sent { background: #fff3cd; color: #856404; }
		.status-paid { background: #d4edda; color: #155724; }
		.no-print { display: none; }
		@media print {
			body { padding: 0; margin: 0; }
			.container { padding: 0; }
			.no-print { display: none; }
		}
	</style>
</head>
<body>
	<div class="container">
		<header>
			<h1>Invoice</h1>
			<div style="margin-top: 20px;">
				<div class="info-label">Invoice Number</div>
				<div class="info-value" style="font-size: 24px; font-weight: 600;">` + invoice.Number + `</div>
				<div style="margin-top: 15px;">
					<span class="status-badge status-` + invoice.Status + `">` + statusBadge + `</span>
				</div>
			</div>
		</header>

		<div class="invoice-info">
			<div class="info-section">
				<div class="info-label">Bill To</div>
				<div class="info-value">
					<div>` + client.Name + `</div>
					<div>` + client.Email + `</div>
					<div>` + client.Phone + `</div>
				</div>
			</div>
			<div class="info-section">
				<div class="info-label">Invoice Details</div>
				<div class="info-value">
					<div style="margin-bottom: 8px;"><strong>Invoice Date:</strong> ` + createdStr + `</div>
					<div style="margin-bottom: 8px;"><strong>Due Date:</strong> ` + dueStr + `</div>
					<div><strong>Invoice ID:</strong> <code>` + invoice.ID + `</code></div>
				</div>
			</div>
		</div>

		<div class="divider"></div>

		<table>
			<thead>
				<tr>
					<th>Description</th>
					<th class="amount-column">Amount</th>
				</tr>
			</thead>
			<tbody>
				<tr>
					<td>` + invoice.Description + `</td>
					<td class="amount-column">$` + fmt.Sprintf("%.2f", invoice.Amount) + `</td>
				</tr>
			</tbody>
		</table>

		<div class="totals">
			<div class="totals-section">
				<div class="total-row final">
					<span>Total Due</span>
					<span>$` + fmt.Sprintf("%.2f", invoice.Amount) + `</span>
				</div>
			</div>
		</div>

		<div class="divider"></div>

		<div style="text-align: center; color: #999; font-size: 12px; no-print">
			<p>Thank you for your business!</p>
		</div>

		<div style="text-align: center; margin-top: 40px;" class="no-print">
			<button onclick="window.print()" style="padding: 10px 20px; background: #0066cc; color: white; border: none; border-radius: 4px; cursor: pointer; font-weight: 600; font-size: 14px;">Print Invoice</button>
			<button onclick="history.back()" style="padding: 10px 20px; background: #f0f0f0; color: #333; border: none; border-radius: 4px; cursor: pointer; font-weight: 600; font-size: 14px; margin-left: 10px;">Back</button>
		</div>
	</div>
</body>
</html>`

	fmt.Fprint(w, html)
}
