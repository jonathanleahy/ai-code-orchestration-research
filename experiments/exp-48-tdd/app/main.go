package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"app/store"
)

var db *store.Store

func main() {
	db = store.NewStore()

	mux := http.NewServeMux()

	mux.HandleFunc("/health", handleHealth)
	mux.HandleFunc("/", handleDashboard)

	// Clients
	mux.HandleFunc("/api/clients", handleClients)
	mux.HandleFunc("/api/clients/", handleClientDetail)

	// Activities
	mux.HandleFunc("/api/activities", handleActivities)

	// Invoices
	mux.HandleFunc("/api/invoices", handleInvoices)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Println("Server listening on :8080")
	log.Fatal(server.ListenAndServe())
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func handleDashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(dashboardHTML))
}

func handleClients(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodGet {
		clients := db.ListClients()
		if clients == nil {
			clients = []*store.Client{}
		}
		json.NewEncoder(w).Encode(clients)
		return
	}

	if r.Method == http.MethodPost {
		r.ParseMultipartForm(10 << 20)
		name := r.FormValue("name")
		email := r.FormValue("email")
		phone := r.FormValue("phone")

		if name == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "name is required"})
			return
		}

		client := db.AddClient(name, email, phone)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(client)
		return
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
	json.NewEncoder(w).Encode(map[string]string{"error": "method not allowed"})
}

func handleClientDetail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/api/clients/"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid client id"})
		return
	}

	if r.Method == http.MethodGet {
		client := db.GetClient(id)
		if client == nil {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "client not found"})
			return
		}
		json.NewEncoder(w).Encode(client)
		return
	}

	if r.Method == http.MethodPut {
		r.ParseMultipartForm(10 << 20)
		name := r.FormValue("name")
		email := r.FormValue("email")
		phone := r.FormValue("phone")

		if name == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "name is required"})
			return
		}

		client := db.UpdateClient(id, name, email, phone)
		if client == nil {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "client not found"})
			return
		}
		json.NewEncoder(w).Encode(client)
		return
	}

	if r.Method == http.MethodDelete {
		ok := db.DeleteClient(id)
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "client not found"})
			return
		}
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
	json.NewEncoder(w).Encode(map[string]string{"error": "method not allowed"})
}

func handleActivities(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodPost {
		r.ParseMultipartForm(10 << 20)
		clientID, err := strconv.Atoi(r.FormValue("client_id"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "invalid client_id"})
			return
		}

		actType := r.FormValue("type")
		desc := r.FormValue("description")

		activity := db.AddActivity(clientID, actType, desc)
		if activity == nil {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "client not found"})
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(activity)
		return
	}

	if r.Method == http.MethodGet {
		clientID := r.URL.Query().Get("client_id")
		if clientID == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "client_id is required"})
			return
		}

		id, err := strconv.Atoi(clientID)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "invalid client_id"})
			return
		}

		activities := db.GetActivities(id)
		if activities == nil {
			activities = []*store.Activity{}
		}
		json.NewEncoder(w).Encode(activities)
		return
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
	json.NewEncoder(w).Encode(map[string]string{"error": "method not allowed"})
}

func handleInvoices(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodGet {
		invoices := db.ListInvoices()
		if invoices == nil {
			invoices = []*store.Invoice{}
		}
		json.NewEncoder(w).Encode(invoices)
		return
	}

	if r.Method == http.MethodPost {
		r.ParseMultipartForm(10 << 20)
		clientID, err := strconv.Atoi(r.FormValue("client_id"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "invalid client_id"})
			return
		}

		number := r.FormValue("number")
		amount, err := strconv.ParseFloat(r.FormValue("amount"), 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "invalid amount"})
			return
		}

		status := r.FormValue("status")
		dueDate := r.FormValue("due_date")

		invoice := db.CreateInvoice(clientID, number, amount, status, dueDate)
		if invoice == nil {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "client not found or invalid due date"})
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(invoice)
		return
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
	json.NewEncoder(w).Encode(map[string]string{"error": "method not allowed"})
}

const dashboardHTML = `<!DOCTYPE html>
<html>
<head>
<title>Dashboard</title>
<style>
body { font-family: sans-serif; margin: 20px; }
.section { margin: 30px 0; }
table { border-collapse: collapse; width: 100%; }
th, td { border: 1px solid #ddd; padding: 8px; text-align: left; }
th { background-color: #f2f2f2; }
button { padding: 8px 12px; margin: 5px; cursor: pointer; }
input { padding: 6px; margin: 5px; }
.form { background: #f9f9f9; padding: 15px; margin: 15px 0; }
</style>
</head>
<body>
<h1>Client Management Dashboard</h1>

<div class="section">
<h2>Add Client</h2>
<div class="form">
  <input type="text" id="clientName" placeholder="Name">
  <input type="email" id="clientEmail" placeholder="Email">
  <input type="tel" id="clientPhone" placeholder="Phone">
  <button onclick="addClient()">Add Client</button>
</div>
</div>

<div class="section">
<h2>Clients</h2>
<button onclick="loadClients()">Refresh</button>
<table id="clientsTable">
  <thead>
    <tr>
      <th>ID</th>
      <th>Name</th>
      <th>Email</th>
      <th>Phone</th>
      <th>Created</th>
      <th>Actions</th>
    </tr>
  </thead>
  <tbody id="clientsList"></tbody>
</table>
</div>

<div class="section">
<h2>Add Activity</h2>
<div class="form">
  <input type="number" id="actClientId" placeholder="Client ID">
  <input type="text" id="actType" placeholder="Type (call, email, etc)">
  <input type="text" id="actDesc" placeholder="Description">
  <button onclick="addActivity()">Add Activity</button>
</div>
</div>

<div class="section">
<h2>Add Invoice</h2>
<div class="form">
  <input type="number" id="invClientId" placeholder="Client ID">
  <input type="text" id="invNumber" placeholder="Invoice Number">
  <input type="number" id="invAmount" placeholder="Amount" step="0.01">
  <input type="text" id="invStatus" placeholder="Status (draft, sent, paid)">
  <input type="date" id="invDueDate" placeholder="Due Date">
  <button onclick="addInvoice()">Add Invoice</button>
</div>
</div>

<div class="section">
<h2>Invoices</h2>
<button onclick="loadInvoices()">Refresh</button>
<table id="invoicesTable">
  <thead>
    <tr>
      <th>ID</th>
      <th>Client ID</th>
      <th>Number</th>
      <th>Amount</th>
      <th>Status</th>
      <th>Due Date</th>
      <th>Created</th>
    </tr>
  </thead>
  <tbody id="invoicesList"></tbody>
</table>
</div>

<script>
function loadClients() {
  fetch('/api/clients')
    .then(r => r.json())
    .then(data => {
      const tbody = document.getElementById('clientsList');
      tbody.innerHTML = '';
      if (!data || data.length === 0) {
        tbody.innerHTML = '<tr><td colspan="6">No clients</td></tr>';
        return;
      }
      for (let i = 0; i < data.length; i++) {
        const c = data[i];
        const row = '<tr><td>' + c.ID + '</td><td>' + c.Name + '</td><td>' + c.Email + '</td><td>' + c.Phone + '</td><td>' + new Date(c.CreatedAt).toLocaleDateString() + '</td><td><button onclick="deleteClient(' + c.ID + ')">Delete</button></td></tr>';
        tbody.innerHTML = tbody.innerHTML + row;
      }
    })
    .catch(e => alert('Error loading clients: ' + e));
}

function addClient() {
  const name = document.getElementById('clientName').value;
  const email = document.getElementById('clientEmail').value;
  const phone = document.getElementById('clientPhone').value;

  if (!name) {
    alert('Name is required');
    return;
  }

  const fd = new FormData();
  fd.append('name', name);
  fd.append('email', email);
  fd.append('phone', phone);

  fetch('/api/clients', { method: 'POST', body: fd })
    .then(r => r.json())
    .then(data => {
      if (data.error) {
        alert('Error: ' + data.error);
      } else {
        alert('Client added');
        document.getElementById('clientName').value = '';
        document.getElementById('clientEmail').value = '';
        document.getElementById('clientPhone').value = '';
        loadClients();
      }
    })
    .catch(e => alert('Error: ' + e));
}

function deleteClient(id) {
  if (!confirm('Delete client ' + id + '?')) return;
  fetch('/api/clients/' + id, { method: 'DELETE' })
    .then(r => {
      if (r.ok) {
        alert('Client deleted');
        loadClients();
      } else {
        return r.json().then(d => alert('Error: ' + d.error));
      }
    })
    .catch(e => alert('Error: ' + e));
}

function addActivity() {
  const clientId = document.getElementById('actClientId').value;
  const type = document.getElementById('actType').value;
  const desc = document.getElementById('actDesc').value;

  if (!clientId || !type) {
    alert('Client ID and Type are required');
    return;
  }

  const fd = new FormData();
  fd.append('client_id', clientId);
  fd.append('type', type);
  fd.append('description', desc);

  fetch('/api/activities', { method: 'POST', body: fd })
    .then(r => r.json())
    .then(data => {
      if (data.error) {
        alert('Error: ' + data.error);
      } else {
        alert('Activity added');
        document.getElementById('actClientId').value = '';
        document.getElementById('actType').value = '';
        document.getElementById('actDesc').value = '';
      }
    })
    .catch(e => alert('Error: ' + e));
}

function addInvoice() {
  const clientId = document.getElementById('invClientId').value;
  const number = document.getElementById('invNumber').value;
  const amount = document.getElementById('invAmount').value;
  const status = document.getElementById('invStatus').value;
  const dueDate = document.getElementById('invDueDate').value;

  if (!clientId || !number || !amount || !dueDate) {
    alert('All fields are required');
    return;
  }

  const fd = new FormData();
  fd.append('client_id', clientId);
  fd.append('number', number);
  fd.append('amount', amount);
  fd.append('status', status);
  fd.append('due_date', dueDate);

  fetch('/api/invoices', { method: 'POST', body: fd })
    .then(r => r.json())
    .then(data => {
      if (data.error) {
        alert('Error: ' + data.error);
      } else {
        alert('Invoice added');
        document.getElementById('invClientId').value = '';
        document.getElementById('invNumber').value = '';
        document.getElementById('invAmount').value = '';
        document.getElementById('invStatus').value = '';
        document.getElementById('invDueDate').value = '';
        loadInvoices();
      }
    })
    .catch(e => alert('Error: ' + e));
}

function loadInvoices() {
  fetch('/api/invoices')
    .then(r => r.json())
    .then(data => {
      const tbody = document.getElementById('invoicesList');
      tbody.innerHTML = '';
      if (!data || data.length === 0) {
        tbody.innerHTML = '<tr><td colspan="7">No invoices</td></tr>';
        return;
      }
      for (let i = 0; i < data.length; i++) {
        const inv = data[i];
        const row = '<tr><td>' + inv.ID + '</td><td>' + inv.ClientID + '</td><td>' + inv.Number + '</td><td>$' + inv.Amount.toFixed(2) + '</td><td>' + inv.Status + '</td><td>' + new Date(inv.DueDate).toLocaleDateString() + '</td><td>' + new Date(inv.CreatedAt).toLocaleDateString() + '</td></tr>';
        tbody.innerHTML = tbody.innerHTML + row;
      }
    })
    .catch(e => alert('Error loading invoices: ' + e));
}

window.onload = function() {
  loadClients();
  loadInvoices();
};
</script>
</body>
</html>
`
