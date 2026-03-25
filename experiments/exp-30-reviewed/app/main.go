package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net/http"
	"sort"
	"strings"
	"app/store"
)

var s *store.Store

func init() {
	s = store.NewStore()
}

// securityMiddleware adds security headers to all responses
func securityMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Content Security Policy to prevent XSS attacks
		w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'; img-src 'self' data:")
		// Prevent clickjacking
		w.Header().Set("X-Frame-Options", "DENY")
		// Prevent MIME sniffing
		w.Header().Set("X-Content-Type-Options", "nosniff")
		next(w, r)
	}
}

func main() {
	http.HandleFunc("/health", handleHealth)
	http.HandleFunc("/", securityMiddleware(handleDashboard))
	http.HandleFunc("/client/", securityMiddleware(handleClientPage))

	// API endpoints
	http.HandleFunc("/api/clients", handleClientsAPI)
	http.HandleFunc("/api/client/", handleClientAPI)
	http.HandleFunc("/api/clients/export", handleExportClientsAPI)
	http.HandleFunc("/api/activities", handleActivitiesAPI)
	http.HandleFunc("/api/activity/", handleActivityAPI)
	http.HandleFunc("/api/invoices", handleInvoicesAPI)
	http.HandleFunc("/api/invoice/", handleInvoiceAPI)

	log.Println("Server starting at :8080")
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

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	clients := s.ListClients()

	sort.Slice(clients, func(i, j int) bool {
		return clients[i].Name < clients[j].Name
	})

	html := getDashboardHTML(clients)
	fmt.Fprint(w, html)
}

func handleClientPage(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.NotFound(w, r)
		return
	}

	clientID := parts[2]
	client, err := s.GetClient(clientID)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Get activities for this client
	allActivities := s.ListActivities()
	var activities []*store.Activity
	for _, a := range allActivities {
		if a.ClientID == clientID {
			activities = append(activities, a)
		}
	}
	sort.Slice(activities, func(i, j int) bool {
		return activities[i].Timestamp.After(activities[j].Timestamp)
	})

	// Get invoices for this client
	allInvoices := s.ListInvoices()
	var invoices []*store.Invoice
	for _, inv := range allInvoices {
		if inv.ClientID == clientID {
			invoices = append(invoices, inv)
		}
	}
	sort.Slice(invoices, func(i, j int) bool {
		return invoices[i].CreatedAt.After(invoices[j].CreatedAt)
	})

	html := getClientPageHTML(client, activities, invoices)
	fmt.Fprint(w, html)
}

func handleClientsAPI(w http.ResponseWriter, r *http.Request) {
	// TODO: Add rate limiting (e.g., max 100 requests per minute per IP) in production
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodPost {
		var client store.Client
		if err := json.NewDecoder(r.Body).Decode(&client); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "invalid request"})
			return
		}

		if err := s.CreateClient(&client); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(client)
		return
	}

	if r.Method == http.MethodGet {
		clients := s.ListClients()
		sort.Slice(clients, func(i, j int) bool {
			return clients[i].Name < clients[j].Name
		})
		json.NewEncoder(w).Encode(clients)
		return
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
}

func handleClientAPI(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		http.NotFound(w, r)
		return
	}

	clientID := parts[3]
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodGet {
		client, err := s.GetClient(clientID)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "client not found"})
			return
		}
		json.NewEncoder(w).Encode(client)
		return
	}

	if r.Method == http.MethodPatch {
		var client store.Client
		if err := json.NewDecoder(r.Body).Decode(&client); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "invalid request"})
			return
		}

		if err := s.UpdateClient(clientID, &client); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		json.NewEncoder(w).Encode(client)
		return
	}

	if r.Method == http.MethodDelete {
		if err := s.DeleteClient(clientID); err != nil {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "client not found"})
			return
		}

		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
}

func handleActivitiesAPI(w http.ResponseWriter, r *http.Request) {
	// TODO: Add rate limiting and authentication in production
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodPost {
		var activity store.Activity
		if err := json.NewDecoder(r.Body).Decode(&activity); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "invalid request"})
			return
		}

		if err := s.CreateActivity(&activity); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(activity)
		return
	}

	if r.Method == http.MethodGet {
		activities := s.ListActivities()
		json.NewEncoder(w).Encode(activities)
		return
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
}

func handleActivityAPI(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		http.NotFound(w, r)
		return
	}

	activityID := parts[3]
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodGet {
		activity, err := s.GetActivity(activityID)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "activity not found"})
			return
		}
		json.NewEncoder(w).Encode(activity)
		return
	}

	if r.Method == http.MethodPatch {
		var activity store.Activity
		if err := json.NewDecoder(r.Body).Decode(&activity); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "invalid request"})
			return
		}

		if err := s.UpdateActivity(activityID, &activity); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		json.NewEncoder(w).Encode(activity)
		return
	}

	if r.Method == http.MethodDelete {
		if err := s.DeleteActivity(activityID); err != nil {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "activity not found"})
			return
		}

		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
}

func handleInvoicesAPI(w http.ResponseWriter, r *http.Request) {
	// TODO: Add rate limiting and authentication in production
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodPost {
		var invoice store.Invoice
		if err := json.NewDecoder(r.Body).Decode(&invoice); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "invalid request"})
			return
		}

		if err := s.CreateInvoice(&invoice); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(invoice)
		return
	}

	if r.Method == http.MethodGet {
		invoices := s.ListInvoices()
		json.NewEncoder(w).Encode(invoices)
		return
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
}

func handleInvoiceAPI(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		http.NotFound(w, r)
		return
	}

	invoiceID := parts[3]
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodGet {
		invoice, err := s.GetInvoice(invoiceID)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "invoice not found"})
			return
		}
		json.NewEncoder(w).Encode(invoice)
		return
	}

	if r.Method == http.MethodPatch {
		var invoice store.Invoice
		if err := json.NewDecoder(r.Body).Decode(&invoice); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "invalid request"})
			return
		}

		if err := s.UpdateInvoice(invoiceID, &invoice); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		json.NewEncoder(w).Encode(invoice)
		return
	}

	if r.Method == http.MethodDelete {
		if err := s.DeleteInvoice(invoiceID); err != nil {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "invoice not found"})
			return
		}

		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
}

func handleExportClientsAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", "attachment; filename=\"clients.csv\"")

	clients := s.ListClients()
	sort.Slice(clients, func(i, j int) bool {
		return clients[i].Name < clients[j].Name
	})

	writer := csv.NewWriter(w)
	defer writer.Flush()

	// Write header
	writer.Write([]string{"ID", "Name", "Company", "Email", "Phone", "Status"})

	// Write data rows
	for _, c := range clients {
		writer.Write([]string{c.ID, c.Name, "", c.Email, c.Phone, "Active"})
	}
}

func getDashboardHTML(clients []*store.Client) string {
	clientRows := ""
	if len(clients) == 0 {
		clientRows = "<tr><td colspan='5' style='text-align:center;padding:40px;color:#999;'><svg style='width:48px;height:48px;margin-bottom:12px;opacity:0.5;' aria-hidden='true' fill='none' stroke='currentColor' viewBox='0 0 24 24'><circle cx='12' cy='8' r='4'/><path d='M4 20c0-4 3-8 8-8s8 4 8 8'/></svg><p style='margin:0;font-size:16px;'>No clients yet — add your first client to get started</p></td></tr>"
	} else {
		for _, c := range clients {
			total := 0.0
			unpaid := 0.0
			for _, inv := range s.ListInvoices() {
				if inv.ClientID == c.ID {
					total += inv.Amount
					if inv.Status != "paid" {
						unpaid += inv.Amount
					}
				}
			}
			// HTML-escape all user input to prevent XSS
			escapedID := html.EscapeString(c.ID)
			escapedName := html.EscapeString(c.Name)
			escapedEmail := html.EscapeString(c.Email)
			escapedPhone := html.EscapeString(c.Phone)
			clientRows += "<tr role='button' tabindex='0' onclick='goToClient(\"" + escapedID + "\")' onkeypress='if(event.key===\" \"||event.key===\"Enter\")goToClient(\"" + escapedID + "\")' style='cursor:pointer;'><td style='padding:12px;border-bottom:1px solid #eee;'>" + escapedName + "</td><td style='padding:12px;border-bottom:1px solid #eee;'>" + escapedEmail + "</td><td style='padding:12px;border-bottom:1px solid #eee;'>" + escapedPhone + "</td><td style='padding:12px;border-bottom:1px solid #eee;'>$" + fmt.Sprintf("%.2f", unpaid) + " of $" + fmt.Sprintf("%.2f", total) + "</td><td style='padding:12px;border-bottom:1px solid #eee;'>" + c.CreatedAt.Format("Jan 2, 2006") + "</td></tr>"
		}
	}

	return "<!DOCTYPE html><html lang='en'><head><meta charset='UTF-8'><meta name='viewport' content='width=device-width,initial-scale=1'><title>Client Management</title><style>*{margin:0;padding:0;box-sizing:border-box}body{font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',sans-serif;background-color:#f9fafb;color:#1f2937;line-height:1.6}header{background-color:#fff;border-bottom:1px solid #e5e7eb;padding:20px 0;box-shadow:0 1px 2px rgba(0,0,0,0.05)}header .container{max-width:1200px;margin:0 auto;padding:0 20px;display:flex;justify-content:space-between;align-items:center}h1{font-size:28px;font-weight:700;margin-bottom:5px}.header-buttons{display:flex;gap:10px;}button{background-color:#3b82f6;color:white;border:none;padding:10px 16px;border-radius:6px;cursor:pointer;font-size:14px;font-weight:500;transition:background-color 0.2s;outline-offset:2px}button:hover{background-color:#2563eb}button:focus{outline:2px solid #3b82f6;outline-offset:2px}.container{max-width:1200px;margin:0 auto;padding:0 20px}.main-content{margin-top:30px}section{margin-bottom:30px}.search-bar{display:flex;gap:10px;margin-bottom:20px;align-items:center}#searchInput{flex:1;padding:10px 12px;border:1px solid #d1d5db;border-radius:6px;font-size:14px;outline-offset:2px}#searchInput:focus{outline:2px solid #3b82f6}table{width:100%;background-color:white;border-radius:8px;overflow:hidden;box-shadow:0 1px 3px rgba(0,0,0,0.1)}thead{background-color:#f3f4f6;font-weight:600;font-size:13px;color:#6b7280;text-transform:uppercase;letter-spacing:0.5px}th{padding:12px;text-align:left;border-bottom:2px solid #e5e7eb}tbody tr:hover{background-color:#f9fafb}tbody tr:focus{outline:2px solid #3b82f6;outline-offset:-2px}td{padding:12px;border-bottom:1px solid #eee}.modal{display:none;position:fixed;top:0;left:0;right:0;bottom:0;background-color:rgba(0,0,0,0.5);z-index:1000;justify-content:center;align-items:center}.modal.active{display:flex}.modal-content{background-color:white;border-radius:8px;padding:30px;width:100%;max-width:500px;box-shadow:0 10px 40px rgba(0,0,0,0.2)}.modal-header{font-size:20px;font-weight:600;margin-bottom:20px}.modal-body{margin-bottom:20px}.form-group{margin-bottom:15px}label{display:block;margin-bottom:5px;font-weight:500;font-size:14px;color:#374151}input[type='email'],input[type='tel'],input[type='text'],select{width:100%;padding:8px 12px;border:1px solid #d1d5db;border-radius:6px;font-size:14px;outline-offset:2px}input:focus,select:focus{outline:2px solid #3b82f6}input.error{border-color:#ef4444}input.error~.error-msg{color:#ef4444;font-size:12px;margin-top:4px;display:block}.modal-footer{display:flex;gap:10px;justify-content:flex-end}.required::after{content:' *';color:#ef4444}</style></head><body><header role='banner'><div class='container'><h1>Client Management</h1><div class='header-buttons'><button aria-label='Add a new client' onclick='openAddClientModal()'>+ Add Client</button><button aria-label='Export client list as CSV' onclick='downloadCSV()'>Export CSV</button></div></div></header><main role='main'><div class='container'><section aria-label='Clients list'><div class='search-bar'><label for='searchInput' style='margin:0;font-weight:600;font-size:14px;'>Search:</label><input type='text' id='searchInput' placeholder='Search clients by name, email...' onkeyup='filterClients()' aria-label='Search clients'></div><table role='table'><thead><tr><th>Name</th><th>Email</th><th>Phone</th><th>Unpaid Balance</th><th>Added</th></tr></thead><tbody id='clientsTableBody'>" + clientRows + "</tbody></table></section></div></main><div id='addClientModal' class='modal' role='dialog' aria-labelledby='addClientTitle'><div class='modal-content'><div class='modal-header' id='addClientTitle'>Add New Client</div><form onsubmit='handleAddClient(event)'><div class='form-group'><label for='clientName' class='required'>Name</label><input type='text' id='clientName' required aria-required='true'></div><div class='form-group'><label for='clientEmail' class='required'>Email</label><input type='email' id='clientEmail' required aria-required='true'></div><div class='form-group'><label for='clientPhone'>Phone</label><input type='tel' id='clientPhone'></div><div class='form-group'><label for='clientPayment'>Preferred Payment</label><select id='clientPayment'><option value=''>Select method</option><option value='bank'>Bank Transfer</option><option value='stripe'>Stripe</option><option value='paypal'>PayPal</option><option value='cash'>Cash</option></select></div><div class='modal-footer'><button type='button' aria-label='Cancel adding a client' onclick='closeAddClientModal()' style='background-color:#e5e7eb;color:#1f2937;'>Cancel</button><button type='submit' aria-label='Submit form to create a new client'>Create Client</button></div></form></div></div><script>function openAddClientModal(){var modal=document.getElementById('addClientModal');modal.classList.add('active');document.getElementById('clientName').focus()}function closeAddClientModal(){var modal=document.getElementById('addClientModal');modal.classList.remove('active');document.querySelector('form').reset()}function handleAddClient(event){event.preventDefault();var name=document.getElementById('clientName').value;var email=document.getElementById('clientEmail').value;var phone=document.getElementById('clientPhone').value;var payment=document.getElementById('clientPayment').value;if(!name||!email){alert('Name and Email are required');return}var client={name:name,email:email,phone:phone,preferred_payment:payment};fetch('/api/clients',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify(client)}).then(function(response){if(response.status===201){closeAddClientModal();location.reload()}else{alert('Error creating client')}}).catch(function(error){console.error('Error:',error);alert('Error creating client')})}function filterClients(){var input=document.getElementById('searchInput').value.toLowerCase();var table=document.getElementById('clientsTableBody');var rows=table.getElementsByTagName('tr');for(var i=0;i<rows.length;i++){var row=rows[i];if(row.textContent.toLowerCase().indexOf(input)>-1){row.style.display=''}else{row.style.display='none'}}}function downloadCSV(){fetch('/api/clients/export').then(function(response){return response.blob()}).then(function(blob){var url=window.URL.createObjectURL(blob);var a=document.createElement('a');a.href=url;a.download='clients.csv';document.body.appendChild(a);a.click();window.URL.revokeObjectURL(url);document.body.removeChild(a)}).catch(function(error){console.error('Error:',error);alert('Error downloading CSV')})}function goToClient(clientId){window.location.href='/client/'+clientId}</script></body></html>"
}

func getClientPageHTML(client *store.Client, activities []*store.Activity, invoices []*store.Invoice) string {
	// HTML-escape client data to prevent XSS
	escapedClientID := html.EscapeString(client.ID)
	escapedClientName := html.EscapeString(client.Name)
	escapedClientEmail := html.EscapeString(client.Email)
	escapedClientPhone := html.EscapeString(client.Phone)
	escapedClientPayment := html.EscapeString(client.PreferredPayment)

	// Calculate totals
	totalInvoiced := 0.0
	totalPaid := 0.0
	totalUnpaid := 0.0
	for _, inv := range invoices {
		totalInvoiced += inv.Amount
		if inv.Status == "paid" {
			totalPaid += inv.Amount
		} else {
			totalUnpaid += inv.Amount
		}
	}

	// Build activity timeline
	activityHTML := ""
	if len(activities) == 0 {
		activityHTML = "<div style='text-align:center;padding:40px;color:#999;'><p style='margin:0;'>No activities yet — add your first activity</p></div>"
	} else {
		for _, a := range activities {
			escapedActivityID := html.EscapeString(a.ID)
			escapedType := html.EscapeString(a.Type)
			escapedDetails := html.EscapeString(a.Details)
			activityHTML += "<div style='padding:15px;border-bottom:1px solid #eee;display:flex;justify-content:space-between;align-items:start;'><div><p style='margin:0 0 4px 0;font-weight:500;color:#1f2937;'>" + escapedType + "</p><p style='margin:0 0 8px 0;font-size:13px;color:#6b7280;'>" + a.Timestamp.Format("Jan 2, 2006 at 3:04 PM") + "</p><p style='margin:0;font-size:14px;color:#374151;'>" + escapedDetails + "</p></div><button aria-label='Delete activity' onclick='deleteActivity(\"" + escapedActivityID + "\")' style='background-color:#ef4444;padding:4px 8px;font-size:12px;'>Delete</button></div>"
		}
	}

	// Build invoices table
	invoiceHTML := ""
	if len(invoices) == 0 {
		invoiceHTML = "<tr><td colspan='5' style='text-align:center;padding:40px;color:#999;'>No invoices yet</td></tr>"
	} else {
		for _, inv := range invoices {
			statusColor := "#10b981"
			if inv.Status == "unpaid" {
				statusColor = "#ef4444"
			} else if inv.Status == "due" {
				statusColor = "#f59e0b"
			}
			escapedInvID := html.EscapeString(inv.ID)
			escapedInvStatus := html.EscapeString(inv.Status)
			invoiceHTML += "<tr><td style='padding:12px;border-bottom:1px solid #eee;'>" + escapedInvID[0:8] + "...</td><td style='padding:12px;border-bottom:1px solid #eee;'>$" + fmt.Sprintf("%.2f", inv.Amount) + "</td><td style='padding:12px;border-bottom:1px solid #eee;'><span style='background-color:" + statusColor + ";color:white;padding:4px 8px;border-radius:4px;font-size:12px;'>" + escapedInvStatus + "</span></td><td style='padding:12px;border-bottom:1px solid #eee;'>" + inv.DueDate.Format("Jan 2, 2006") + "</td><td style='padding:12px;border-bottom:1px solid #eee;'><button aria-label='Delete invoice' onclick='deleteInvoice(\"" + escapedInvID + "\")' style='background-color:#ef4444;padding:4px 8px;font-size:12px;'>Delete</button></td></tr>"
		}
	}

	return "<!DOCTYPE html><html lang='en'><head><meta charset='UTF-8'><meta name='viewport' content='width=device-width,initial-scale=1'><title>" + escapedClientName + " - Client Details</title><style>*{margin:0;padding:0;box-sizing:border-box}body{font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',sans-serif;background-color:#f9fafb;color:#1f2937;line-height:1.6}header{background-color:#fff;border-bottom:1px solid #e5e7eb;padding:20px 0;box-shadow:0 1px 2px rgba(0,0,0,0.05)}header .container{max-width:1200px;margin:0 auto;padding:0 20px}h1{font-size:24px;font-weight:600;margin-bottom:10px}a{color:#3b82f6;text-decoration:none}a:hover{text-decoration:underline}.container{max-width:1200px;margin:0 auto;padding:0 20px}.header-actions{display:flex;gap:10px;margin-bottom:30px}.header-actions button{background-color:#3b82f6;color:white;border:none;padding:10px 16px;border-radius:6px;cursor:pointer;font-size:14px;font-weight:500}.header-actions button.delete{background-color:#ef4444}.tabs{display:flex;gap:20px;border-bottom:2px solid #e5e7eb;margin:30px 0;}.tabs button{background:none;border:none;padding:10px 0;cursor:pointer;font-size:14px;font-weight:600;color:#6b7280;border-bottom:3px solid transparent;margin-bottom:-2px}.tabs button.active{color:#3b82f6;border-bottom-color:#3b82f6}.tab-content{display:none}.tab-content.active{display:block}.card{background-color:white;border-radius:8px;padding:20px;box-shadow:0 1px 3px rgba(0,0,0,0.1);margin-bottom:20px}.card h2{font-size:18px;font-weight:600;margin-bottom:15px}.info-row{display:grid;grid-template-columns:150px 1fr;gap:20px;margin-bottom:12px}.info-label{font-weight:500;color:#6b7280}.form-group{margin-bottom:15px}label{display:block;margin-bottom:5px;font-weight:500;font-size:14px}input,select{width:100%;padding:8px 12px;border:1px solid #d1d5db;border-radius:6px;font-size:14px}.modal{display:none;position:fixed;top:0;left:0;right:0;bottom:0;background-color:rgba(0,0,0,0.5);z-index:1000;justify-content:center;align-items:center}.modal.active{display:flex}.modal-content{background-color:white;border-radius:8px;padding:30px;width:100%;max-width:500px;box-shadow:0 10px 40px rgba(0,0,0,0.2)}.modal-header{font-size:18px;font-weight:600;margin-bottom:20px}.modal-footer{display:flex;gap:10px;justify-content:flex-end;margin-top:20px}button[type='submit']{background-color:#3b82f6;color:white;border:none;padding:10px 16px;border-radius:6px;cursor:pointer;font-weight:500}button[type='submit']:hover{background-color:#2563eb}table{width:100%;background-color:white;border-radius:8px;overflow:hidden}thead{background-color:#f3f4f6;font-weight:600;font-size:13px;color:#6b7280}th{padding:12px;text-align:left;border-bottom:2px solid #e5e7eb}tbody tr:hover{background-color:#f9fafb}td{padding:12px;border-bottom:1px solid #eee}.required::after{content:' *';color:#ef4444}</style></head><body><header><div class='container'><h1>" + escapedClientName + "</h1><p style='color:#6b7280;margin:0;'><a href='/'>← Back to Clients</a></p></div></header><div class='container'><div class='header-actions'><button onclick='openEditModal()'>Edit Client</button><button class='delete' onclick='openDeleteModal()'>Delete Client</button></div><div class='tabs'><button class='tab-btn active' onclick='switchTab(\"details\")'>Details</button><button class='tab-btn' onclick='switchTab(\"history\")'>History</button><button class='tab-btn' onclick='switchTab(\"billing\")'>Billing</button></div><div id='details' class='tab-content active'><div class='card'><h2>Client Information</h2><div class='info-row'><span class='info-label'>Name</span><span>" + escapedClientName + "</span></div><div class='info-row'><span class='info-label'>Email</span><span><a href='mailto:" + escapedClientEmail + "'>" + escapedClientEmail + "</a></span></div><div class='info-row'><span class='info-label'>Phone</span><span>" + escapedClientPhone + "</span></div><div class='info-row'><span class='info-label'>Payment Method</span><span>" + escapedClientPayment + "</span></div><div class='info-row'><span class='info-label'>Added</span><span>" + client.CreatedAt.Format("January 2, 2006") + "</span></div><div class='info-row'><span class='info-label'>Last Updated</span><span>" + client.UpdatedAt.Format("January 2, 2006") + "</span></div></div></div><div id='history' class='tab-content'><div class='card'><h2 style='margin-bottom:20px;'>Activity Timeline</h2><div style='margin-bottom:20px;'><button onclick='openAddActivityModal()' style='background-color:#10b981;color:white;border:none;padding:10px 16px;border-radius:6px;cursor:pointer;font-weight:500;'>+ Add Activity</button></div>" + activityHTML + "</div></div><div id='billing' class='tab-content'><div class='card'><h2 style='margin-bottom:20px;'>Billing Summary</h2><div style='display:grid;grid-template-columns:repeat(3,1fr);gap:20px;margin-bottom:30px;'><div style='padding:15px;background-color:#f3f4f6;border-radius:6px;'><p style='margin:0 0 5px 0;font-size:12px;color:#6b7280;font-weight:600;text-transform:uppercase;'>Total Invoiced</p><p style='margin:0;font-size:24px;font-weight:600;color:#1f2937;'>$" + fmt.Sprintf("%.2f", totalInvoiced) + "</p></div><div style='padding:15px;background-color:#dbeafe;border-radius:6px;'><p style='margin:0 0 5px 0;font-size:12px;color:#1e40af;font-weight:600;text-transform:uppercase;'>Paid</p><p style='margin:0;font-size:24px;font-weight:600;color:#1e40af;'>$" + fmt.Sprintf("%.2f", totalPaid) + "</p></div><div style='padding:15px;background-color:#fee2e2;border-radius:6px;'><p style='margin:0 0 5px 0;font-size:12px;color:#7f1d1d;font-weight:600;text-transform:uppercase;'>Unpaid</p><p style='margin:0;font-size:24px;font-weight:600;color:#7f1d1d;'>$" + fmt.Sprintf("%.2f", totalUnpaid) + "</p></div></div><div><button onclick='openAddInvoiceModal()' style='background-color:#10b981;color:white;border:none;padding:10px 16px;border-radius:6px;cursor:pointer;font-weight:500;margin-bottom:20px;'>+ Create Invoice</button></div><table><thead><tr><th>Invoice ID</th><th>Amount</th><th>Status</th><th>Due Date</th><th>Action</th></tr></thead><tbody>" + invoiceHTML + "</tbody></table></div></div></div><div id='editModal' class='modal'><div class='modal-content'><div class='modal-header'>Edit Client</div><form onsubmit='handleEditClient(event)'><div class='form-group'><label class='required'>Name</label><input type='text' id='editName' value='" + escapedClientName + "' required></div><div class='form-group'><label class='required'>Email</label><input type='email' id='editEmail' value='" + escapedClientEmail + "' required></div><div class='form-group'><label>Phone</label><input type='tel' id='editPhone' value='" + escapedClientPhone + "'></div><div class='form-group'><label>Preferred Payment</label><select id='editPayment'><option value=''>Select method</option><option value='bank'>Bank Transfer</option><option value='stripe'>Stripe</option><option value='paypal'>PayPal</option><option value='cash'>Cash</option></select></div><div class='modal-footer'><button type='button' onclick='closeEditModal()' style='background-color:#e5e7eb;color:#1f2937;'>Cancel</button><button type='submit'>Save Changes</button></div></form></div></div><div id='deleteModal' class='modal'><div class='modal-content'><div class='modal-header'>Delete Client</div><p style='margin-bottom:20px;'>Are you sure you want to delete <strong>" + escapedClientName + "</strong>? This action cannot be undone.</p><div class='modal-footer'><button onclick='closeDeleteModal()' style='background-color:#e5e7eb;color:#1f2937;'>Cancel</button><button onclick='confirmDelete()' style='background-color:#ef4444;color:white;'>Delete</button></div></div></div><div id='addActivityModal' class='modal'><div class='modal-content'><div class='modal-header'>Add Activity</div><form onsubmit='handleAddActivity(event)'><div class='form-group'><label class='required'>Type</label><select id='activityType' required><option value=''>Select type</option><option value='call'>Call</option><option value='email'>Email</option><option value='meeting'>Meeting</option><option value='note'>Note</option></select></div><div class='form-group'><label class='required'>Details</label><textarea id='activityDetails' style='width:100%;padding:8px 12px;border:1px solid #d1d5db;border-radius:6px;font-family:inherit;' rows='4' required></textarea></div><div class='modal-footer'><button type='button' onclick='closeAddActivityModal()' style='background-color:#e5e7eb;color:#1f2937;'>Cancel</button><button type='submit'>Add Activity</button></div></form></div></div><div id='addInvoiceModal' class='modal'><div class='modal-content'><div class='modal-header'>Create Invoice</div><form onsubmit='handleAddInvoice(event)'><div class='form-group'><label class='required'>Amount</label><input type='number' id='invoiceAmount' step='0.01' min='0' required></div><div class='form-group'><label class='required'>Status</label><select id='invoiceStatus' required><option value='unpaid'>Unpaid</option><option value='paid'>Paid</option><option value='due'>Due</option></select></div><div class='form-group'><label class='required'>Due Date</label><input type='date' id='invoiceDueDate' required></div><div class='modal-footer'><button type='button' onclick='closeAddInvoiceModal()' style='background-color:#e5e7eb;color:#1f2937;'>Cancel</button><button type='submit'>Create Invoice</button></div></form></div></div><script>function switchTab(tabName){var tabs=document.querySelectorAll('.tab-content');tabs.forEach(function(tab){tab.classList.remove('active')});document.getElementById(tabName).classList.add('active');var buttons=document.querySelectorAll('.tab-btn');buttons.forEach(function(btn){btn.classList.remove('active')});event.target.classList.add('active')}function openEditModal(){document.getElementById('editPayment').value='" + escapedClientPayment + "';document.getElementById('editModal').classList.add('active')}function closeEditModal(){document.getElementById('editModal').classList.remove('active')}function handleEditClient(event){event.preventDefault();var name=document.getElementById('editName').value;var email=document.getElementById('editEmail').value;var phone=document.getElementById('editPhone').value;var payment=document.getElementById('editPayment').value;var updates={name:name,email:email,phone:phone,preferred_payment:payment};fetch('/api/client/" + escapedClientID + "',{method:'PATCH',headers:{'Content-Type':'application/json'},body:JSON.stringify(updates)}).then(function(response){if(response.ok){location.reload()}else{alert('Error updating client')}}).catch(function(error){console.error('Error:',error);alert('Error updating client')})}function openDeleteModal(){document.getElementById('deleteModal').classList.add('active')}function closeDeleteModal(){document.getElementById('deleteModal').classList.remove('active')}function confirmDelete(){fetch('/api/client/" + escapedClientID + "',{method:'DELETE'}).then(function(response){if(response.status===204){window.location.href='/'}else{alert('Error deleting client')}}).catch(function(error){console.error('Error:',error);alert('Error deleting client')})}function openAddActivityModal(){document.getElementById('addActivityModal').classList.add('active')}function closeAddActivityModal(){document.getElementById('addActivityModal').classList.remove('active');document.querySelector('#addActivityModal form').reset()}function handleAddActivity(event){event.preventDefault();var type=document.getElementById('activityType').value;var details=document.getElementById('activityDetails').value;var activity={client_id:'" + escapedClientID + "',type:type,details:details};fetch('/api/activities',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify(activity)}).then(function(response){if(response.status===201){closeAddActivityModal();location.reload()}else{alert('Error creating activity')}}).catch(function(error){console.error('Error:',error);alert('Error creating activity')})}function deleteActivity(activityId){if(confirm('Delete this activity?')){fetch('/api/activity/'+activityId,{method:'DELETE'}).then(function(response){if(response.status===204){location.reload()}else{alert('Error deleting activity')}}).catch(function(error){console.error('Error:',error);alert('Error deleting activity')})}}function openAddInvoiceModal(){document.getElementById('addInvoiceModal').classList.add('active')}function closeAddInvoiceModal(){document.getElementById('addInvoiceModal').classList.remove('active');document.querySelector('#addInvoiceModal form').reset()}function handleAddInvoice(event){event.preventDefault();var amount=parseFloat(document.getElementById('invoiceAmount').value);var status=document.getElementById('invoiceStatus').value;var dueDate=document.getElementById('invoiceDueDate').value;var invoice={client_id:'" + escapedClientID + "',amount:amount,status:status,due_date:new Date(dueDate+'T00:00:00Z').toISOString()};fetch('/api/invoices',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify(invoice)}).then(function(response){if(response.status===201){closeAddInvoiceModal();location.reload()}else{alert('Error creating invoice')}}).catch(function(error){console.error('Error:',error);alert('Error creating invoice')})}function deleteInvoice(invoiceId){if(confirm('Delete this invoice?')){fetch('/api/invoice/'+invoiceId,{method:'DELETE'}).then(function(response){if(response.status===204){location.reload()}else{alert('Error deleting invoice')}}).catch(function(error){console.error('Error:',error);alert('Error deleting invoice')})}}</script></body></html>"
}
