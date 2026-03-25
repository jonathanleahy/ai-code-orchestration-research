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

func main() {
	s = store.NewStore()

	// Seed data
	seedData()

	mux := http.NewServeMux()
	mux.HandleFunc("/health", handleHealth)
	mux.HandleFunc("/", handleHome)
	mux.HandleFunc("/api/clients", handleClientsAPI)
	mux.HandleFunc("/api/activities", handleActivitiesAPI)
	mux.HandleFunc("/api/invoices", handleInvoicesAPI)

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func seedData() {
	clients := []*store.Client{
		{ID: "1", Name: "Acme Corp", Email: "contact@acme.com"},
		{ID: "2", Name: "TechStart Inc", Email: "hello@techstart.com"},
		{ID: "3", Name: "Global Services", Email: "info@globalsvcs.com"},
	}
	for _, c := range clients {
		s.AddClient(c)
	}

	projects := []*store.Project{
		{ID: "p1", ClientID: "1", Title: "Website Redesign"},
		{ID: "p2", ClientID: "1", Title: "API Integration"},
		{ID: "p3", ClientID: "2", Title: "Mobile App"},
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
		{ID: "inv3", ProjectID: "p2", Amount: 8500.00},
		{ID: "inv4", ProjectID: "p3", Amount: 12000.00},
	}
	for _, inv := range invoices {
		s.AddInvoice(inv)
	}
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		r.ParseMultipartForm(32 << 20)
		name := r.FormValue("name")
		email := r.FormValue("email")

		if name == "" || email == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "name and email required"})
			return
		}

		id := fmt.Sprintf("%d", time.Now().UnixNano())
		c := &store.Client{ID: id, Name: name, Email: email}
		s.AddClient(c)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(c)
		return
	}

	if strings.HasPrefix(r.URL.Path, "/client/") {
		clientID := strings.TrimPrefix(r.URL.Path, "/client/")
		handleClientDetail(w, r, clientID)
		return
	}

	renderHome(w)
}

func renderHome(w http.ResponseWriter) {
	clients := s.ListClients()

	html := `<!DOCTYPE html>
<html>
<head>
	<meta charset="utf-8">
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<title>Client Manager</title>
	<style>
		* { margin: 0; padding: 0; box-sizing: border-box; }
		body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif; background: #f9f9f9; color: #333; }
		.container { max-width: 1200px; margin: 0 auto; padding: 20px; }
		.header { margin-bottom: 30px; }
		h1 { font-size: 28px; margin-bottom: 10px; }
		.form-section { background: white; padding: 20px; border-radius: 6px; margin-bottom: 30px; box-shadow: 0 1px 3px rgba(0,0,0,0.1); }
		.form-group { margin-bottom: 15px; }
		label { display: block; font-size: 14px; font-weight: 500; margin-bottom: 5px; }
		input { width: 100%; padding: 8px 12px; border: 1px solid #ddd; border-radius: 4px; font-size: 14px; }
		input:focus { outline: none; border-color: #0066cc; box-shadow: 0 0 0 3px rgba(0,102,204,0.1); }
		button { background: #0066cc; color: white; border: none; padding: 10px 20px; border-radius: 4px; cursor: pointer; font-size: 14px; font-weight: 500; }
		button:hover { background: #0052a3; }
		.clients-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(300px, 1fr)); gap: 15px; }
		.client-card { background: white; padding: 20px; border-radius: 6px; box-shadow: 0 1px 3px rgba(0,0,0,0.1); cursor: pointer; transition: transform 0.2s; }
		.client-card:hover { transform: translateY(-2px); box-shadow: 0 4px 8px rgba(0,0,0,0.15); }
		.client-name { font-size: 16px; font-weight: 600; margin-bottom: 8px; }
		.client-email { font-size: 13px; color: #666; }
		.error { background: #fee; color: #c00; padding: 12px; border-radius: 4px; margin-bottom: 15px; border: 1px solid #fcc; }
	</style>
</head>
<body>
	<div class="container">
		<div class="header">
			<h1>Client Manager</h1>
		</div>

		<div class="form-section">
			<h2 style="font-size: 18px; margin-bottom: 15px;">Add New Client</h2>
			<form id="addClientForm">
				<div class="form-group">
					<label for="name">Name</label>
					<input type="text" id="name" name="name" required>
				</div>
				<div class="form-group">
					<label for="email">Email</label>
					<input type="email" id="email" name="email" required>
				</div>
				<button type="submit">Add Client</button>
			</form>
			<div id="formError" class="error" style="display: none;"></div>
		</div>

		<div>
			<h2 style="font-size: 18px; margin-bottom: 15px;">Clients</h2>
			<div class="clients-grid" id="clientsGrid">
`

	for _, c := range clients {
		html += `				<div class="client-card" onclick="window.location='/client/` + c.ID + `'">
					<div class="client-name">` + c.Name + `</div>
					<div class="client-email">` + c.Email + `</div>
				</div>
`
	}

	html += `			</div>
		</div>
	</div>

	<script>
		document.getElementById('addClientForm').addEventListener('submit', function(e) {
			e.preventDefault();
			var formData = new FormData(this);
			fetch('/', {
				method: 'POST',
				body: formData
			})
			.then(function(resp) { return resp.json(); })
			.then(function(data) {
				if (data.error) {
					document.getElementById('formError').textContent = data.error;
					document.getElementById('formError').style.display = 'block';
				} else {
					window.location.reload();
				}
			})
			.catch(function(err) {
				document.getElementById('formError').textContent = 'Error: ' + err.message;
				document.getElementById('formError').style.display = 'block';
			});
		});
	</script>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, html)
}

func handleClientDetail(w http.ResponseWriter, r *http.Request, clientID string) {
	client := s.GetClient(clientID)
	if client == nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "Client not found")
		return
	}

	html := `<!DOCTYPE html>
<html>
<head>
	<meta charset="utf-8">
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<title>` + client.Name + `</title>
	<style>
		* { margin: 0; padding: 0; box-sizing: border-box; }
		body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif; background: #f9f9f9; color: #333; }
		.container { max-width: 900px; margin: 0 auto; padding: 20px; }
		.header { margin-bottom: 30px; }
		.back-link { display: inline-block; color: #0066cc; text-decoration: none; font-size: 14px; margin-bottom: 15px; }
		.back-link:hover { text-decoration: underline; }
		h1 { font-size: 28px; margin-bottom: 5px; }
		.client-info { font-size: 14px; color: #666; margin-bottom: 30px; }
		.tabs { display: flex; gap: 0; border-bottom: 2px solid #e0e0e0; margin-bottom: 20px; }
		.tab-button { padding: 12px 20px; background: none; border: none; cursor: pointer; font-size: 14px; font-weight: 500; color: #666; border-bottom: 2px solid transparent; margin-bottom: -2px; }
		.tab-button.active { color: #0066cc; border-bottom-color: #0066cc; }
		.tab-content { display: none; }
		.tab-content.active { display: block; }
		.activity-item, .invoice-item { background: white; padding: 15px; margin-bottom: 10px; border-radius: 4px; }
		.activity-type { font-size: 12px; background: #e8f1ff; color: #0066cc; padding: 2px 8px; border-radius: 3px; display: inline-block; margin-bottom: 8px; }
		.activity-content { font-size: 14px; margin-bottom: 5px; }
		.activity-timestamp { font-size: 12px; color: #999; }
		.invoice-header { display: flex; justify-content: space-between; align-items: center; }
		.invoice-amount { font-size: 16px; font-weight: 600; }
	</style>
</head>
<body>
	<div class="container">
		<div class="header">
			<a href="/" class="back-link">← Back to Clients</a>
			<h1>` + client.Name + `</h1>
			<div class="client-info">` + client.Email + `</div>
		</div>

		<div class="tabs">
			<button class="tab-button active" onclick="showTab('activity', this)">Activity</button>
			<button class="tab-button" onclick="showTab('invoices', this)">Invoices</button>
		</div>

		<div id="activity" class="tab-content active"></div>
		<div id="invoices" class="tab-content"></div>
	</div>

	<script>
		function showTab(tabName, btn) {
			document.querySelectorAll('.tab-content').forEach(function(el) { el.classList.remove('active'); });
			document.querySelectorAll('.tab-button').forEach(function(el) { el.classList.remove('active'); });
			document.getElementById(tabName).classList.add('active');
			btn.classList.add('active');
		}

		fetch('/api/activities?clientId=` + clientID + `')
			.then(function(resp) { return resp.json(); })
			.then(function(data) {
				var html = '';
				data.forEach(function(activity) {
					html += '<div class="activity-item">' +
						'<div class="activity-type">' + activity.Type + '</div>' +
						'<div class="activity-content">' + activity.Content + '</div>' +
						'<div class="activity-timestamp">' + activity.Timestamp + '</div>' +
						'</div>';
				});
				if (html === '') html = '<p style="color: #999; font-size: 14px;">No activities</p>';
				document.getElementById('activity').innerHTML = html;
			});

		fetch('/api/invoices?clientId=` + clientID + `')
			.then(function(resp) { return resp.json(); })
			.then(function(data) {
				var html = '';
				data.forEach(function(invoice) {
					html += '<div class="invoice-item">' +
						'<div class="invoice-header">' +
						'<span>Invoice ' + invoice.ID + '</span>' +
						'<span class="invoice-amount">$' + invoice.Amount.toFixed(2) + '</span>' +
						'</div>' +
						'</div>';
				});
				if (html === '') html = '<p style="color: #999; font-size: 14px;">No invoices</p>';
				document.getElementById('invoices').innerHTML = html;
			});
	</script>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, html)
}

func handleClientsAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	clients := s.ListClients()
	json.NewEncoder(w).Encode(clients)
}

func handleActivitiesAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	clientID := r.URL.Query().Get("clientId")
	activities := s.ListActivities()

	var filtered []*store.Activity
	for _, a := range activities {
		if a.ClientID == clientID {
			filtered = append(filtered, a)
		}
	}
	json.NewEncoder(w).Encode(filtered)
}

func handleInvoicesAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	clientID := r.URL.Query().Get("clientId")
	invoices := s.ListInvoices()
	projects := s.ListProjects()

	// Build map of projectID -> clientID
	projectToClient := make(map[string]string)
	for _, p := range projects {
		projectToClient[p.ID] = p.ClientID
	}

	var filtered []*store.Invoice
	for _, inv := range invoices {
		if projectToClient[inv.ProjectID] == clientID {
			filtered = append(filtered, inv)
		}
	}
	json.NewEncoder(w).Encode(filtered)
}
