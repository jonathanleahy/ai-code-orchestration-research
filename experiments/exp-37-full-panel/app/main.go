package main

import (
	"app/store"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var s *store.Store

func init() {
	s = store.NewStore()
	// Add sample data
	s.CreateClient(&store.Client{Name: "Acme Corp", Email: "contact@acme.com"})
	s.CreateClient(&store.Client{Name: "Beta Inc", Email: "hello@beta.com"})
}

func main() {
	http.HandleFunc("/health", handleHealth)
	http.HandleFunc("/", handleDashboard)
	http.HandleFunc("/client/", handleClientRoute)
	http.HandleFunc("/invoice/", handleInvoiceRoute)

	log.Println("Server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func handleDashboard(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseMultipartForm(1024 * 1024)
		name := r.FormValue("name")
		email := r.FormValue("email")
		if name != "" && email != "" {
			s.CreateClient(&store.Client{Name: name, Email: email})
			showToast(w, "Client added successfully")
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
	}

	clients := s.ListClients()
	query := r.URL.Query().Get("q")
	var filtered []*store.Client
	for _, c := range clients {
		if query == "" || strings.Contains(strings.ToLower(c.Name), strings.ToLower(query)) || strings.Contains(strings.ToLower(c.Email), strings.ToLower(query)) {
			filtered = append(filtered, c)
		}
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	html := "<!DOCTYPE html><html><head><title>Dashboard</title><style>" +
		"body{font-family:sans-serif;margin:0;background:#f5f5f5}" +
		".header{background:#fff;padding:20px;border-bottom:1px solid #ddd;display:flex;justify-content:space-between;align-items:center}" +
		".content{max-width:1000px;margin:0 auto;padding:20px}" +
		".btn{padding:8px 16px;border:none;border-radius:4px;cursor:pointer;font-size:14px;text-decoration:none;display:inline-block}" +
		".btn-primary{background:#0066cc;color:#fff}" +
		".btn-primary:hover{background:#0052a3}" +
		".btn-danger{background:#dc2626;color:#fff}" +
		".btn-danger:hover{background:#b91c1c}" +
		"table{width:100%;border-collapse:collapse;background:#fff;border-radius:4px;overflow:hidden}" +
		"th,td{padding:12px;text-align:left;border-bottom:1px solid #ddd}" +
		"th{background:#f9f9f9;font-weight:600}" +
		"tr:hover{background:#fafafa}" +
		".search{margin-bottom:20px}" +
		"input[type=text]{padding:8px 12px;border:1px solid #ddd;border-radius:4px;font-size:14px;width:100%;max-width:300px}" +
		".empty{text-align:center;padding:40px;color:#666}" +
		".empty-btn{margin-top:16px}" +
		".toast{position:fixed;bottom:20px;right:20px;background:#10b981;color:#fff;padding:16px;border-radius:4px;z-index:1000}" +
		"</style></head><body>" +
		"<div class=\"header\"><h1>Dashboard</h1><button class=\"btn btn-primary\" onclick=\"document.getElementById('addForm').style.display='block'\">Add Client</button></div>" +
		"<div class=\"content\">" +
		"<div class=\"search\"><input type=\"text\" id=\"searchInput\" placeholder=\"Search clients...\" onkeyup=\"filterClients()\"></div>"

	if len(filtered) == 0 {
		html += "<div class=\"empty\"><p>No clients yet. Start managing your clients!</p><button class=\"btn btn-primary empty-btn\" onclick=\"document.getElementById('addForm').style.display='block'\">Add your first client</button></div>"
	} else {
		html += "<table><thead><tr><th>Name</th><th>Email</th><th>Created</th><th>Actions</th></tr></thead><tbody>"
		for _, c := range filtered {
			html += "<tr><td><a href=\"/client/" + c.ID + "\" style=\"color:#0066cc;text-decoration:none\">" + c.Name + "</a></td><td>" + c.Email + "</td><td>" + formatDate(c.CreatedAt) + "</td><td><a href=\"/client/" + c.ID + "\" class=\"btn btn-primary\" style=\"font-size:12px\">View</a></td></tr>"
		}
		html += "</tbody></table>"
	}

	html += "<div id=\"addForm\" style=\"display:none;position:fixed;top:0;left:0;right:0;bottom:0;background:rgba(0,0,0,0.5);z-index:999;align-items:center;justify-content:center\">" +
		"<form method=\"POST\" style=\"background:#fff;padding:30px;border-radius:4px;width:100%;max-width:400px\">" +
		"<h2 style=\"margin-top:0\">Add Client</h2>" +
		"<div style=\"margin-bottom:16px\"><label style=\"display:block;margin-bottom:4px;font-weight:500\">Name <span style=\"color:#dc2626\">*</span></label><input type=\"text\" name=\"name\" required style=\"width:100%;padding:8px;border:1px solid #ddd;border-radius:4px;box-sizing:border-box\"></div>" +
		"<div style=\"margin-bottom:16px\"><label style=\"display:block;margin-bottom:4px;font-weight:500\">Email <span style=\"color:#dc2626\">*</span></label><input type=\"email\" name=\"email\" required style=\"width:100%;padding:8px;border:1px solid #ddd;border-radius:4px;box-sizing:border-box\"></div>" +
		"<div style=\"display:flex;gap:8px\"><button type=\"submit\" class=\"btn btn-primary\">Add Client</button><button type=\"button\" class=\"btn\" style=\"background:#e5e5e5\" onclick=\"document.getElementById('addForm').parentElement.style.display='none'\">Cancel</button></div>" +
		"</form></div>" +
		"</div></body>" +
		"<script>" +
		"function filterClients(){var q=document.getElementById('searchInput').value;if(q)window.location.search='q='+encodeURIComponent(q);else window.location.href='/';}" +
		"</script>" +
		"</html>"

	fmt.Fprint(w, html)
}

func handleClientRoute(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/client/"), "/")
	clientID := parts[0]

	if clientID == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	client, ok := s.GetClient(clientID)
	if !ok {
		http.Error(w, "Client not found", http.StatusNotFound)
		return
	}

	if len(parts) > 1 && parts[1] == "invoice" && parts[2] == "new" {
		handleCreateInvoice(w, r, clientID, client)
		return
	}

	if r.Method == "POST" && r.FormValue("_method") == "DELETE" {
		s.DeleteClient(clientID)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if r.Method == "POST" && r.FormValue("_method") != "DELETE" {
		r.ParseMultipartForm(1024 * 1024)
		client.Name = r.FormValue("name")
		client.Email = r.FormValue("email")
		s.UpdateClient(clientID, client)
		showToast(w, "Client updated successfully")
		http.Redirect(w, r, "/client/"+clientID, http.StatusSeeOther)
		return
	}

	tab := r.URL.Query().Get("tab")
	if tab == "" {
		tab = "profile"
	}

	invoices := s.ListInvoicesByClient(clientID)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	html := "<!DOCTYPE html><html><head><title>" + client.Name + "</title><style>" +
		"body{font-family:sans-serif;margin:0;background:#f5f5f5}" +
		".header{background:#fff;padding:20px;border-bottom:1px solid #ddd;display:flex;justify-content:space-between;align-items:center}" +
		".breadcrumb{margin-bottom:20px;font-size:14px}" +
		".breadcrumb a{color:#0066cc;text-decoration:none}" +
		".content{max-width:1000px;margin:0 auto;padding:20px}" +
		".tabs{display:flex;gap:0;border-bottom:1px solid #ddd;margin-bottom:20px}" +
		".tab{padding:12px 16px;border:none;background:none;cursor:pointer;font-size:14px;border-bottom:2px solid transparent;color:#666}" +
		".tab.active{border-bottom-color:#0066cc;color:#0066cc;font-weight:600}" +
		".form-group{margin-bottom:16px}" +
		"label{display:block;margin-bottom:4px;font-weight:500}" +
		"input[type=text],input[type=email],textarea{width:100%;padding:8px;border:1px solid #ddd;border-radius:4px;box-sizing:border-box}" +
		".btn{padding:8px 16px;border:none;border-radius:4px;cursor:pointer;font-size:14px;text-decoration:none;display:inline-block}" +
		".btn-primary{background:#0066cc;color:#fff}" +
		".btn-primary:hover{background:#0052a3}" +
		".btn-danger{background:#dc2626;color:#fff}" +
		".btn-danger:hover{background:#b91c1c}" +
		".back{margin-bottom:20px}" +
		".back a{color:#0066cc;text-decoration:none;font-size:14px}" +
		"table{width:100%;border-collapse:collapse;background:#fff;border-radius:4px;overflow:hidden}" +
		"th,td{padding:12px;text-align:left;border-bottom:1px solid #ddd}" +
		"th{background:#f9f9f9;font-weight:600}" +
		".empty{color:#666;padding:20px;text-align:center}" +
		"</style></head><body>" +
		"<div class=\"header\"><h1>" + client.Name + "</h1></div>" +
		"<div class=\"content\">" +
		"<div class=\"back\"><a href=\"/\">← Back to Dashboard</a></div>" +
		"<div class=\"breadcrumb\"><a href=\"/\">Dashboard</a> > " + client.Name + "</div>"

	html += "<div class=\"tabs\">" +
		"<button class=\"tab" + ifActive(tab == "profile", " active") + "\" onclick=\"window.location.href='/client/" + clientID + "?tab=profile'\">Profile</button>" +
		"<button class=\"tab" + ifActive(tab == "activity", " active") + "\" onclick=\"window.location.href='/client/" + clientID + "?tab=activity'\">Activity</button>" +
		"<button class=\"tab" + ifActive(tab == "invoices", " active") + "\" onclick=\"window.location.href='/client/" + clientID + "?tab=invoices'\">Invoices</button>" +
		"</div>"

	if tab == "profile" {
		html += "<form method=\"POST\" style=\"max-width:600px\">" +
			"<div class=\"form-group\"><label>Name <span style=\"color:#dc2626\">*</span></label><input type=\"text\" name=\"name\" value=\"" + client.Name + "\" required></div>" +
			"<div class=\"form-group\"><label>Email <span style=\"color:#dc2626\">*</span></label><input type=\"email\" name=\"email\" value=\"" + client.Email + "\" required></div>" +
			"<div style=\"display:flex;gap:8px;margin-top:24px\"><button type=\"submit\" class=\"btn btn-primary\">Save Changes</button>" +
			"<button type=\"button\" class=\"btn btn-danger\" onclick=\"if(confirm('Delete this client?')){var f=document.createElement('form');f.method='POST';var m=document.createElement('input');m.name='_method';m.value='DELETE';f.appendChild(m);document.body.appendChild(f);f.submit()}\">Delete Client</button></div>" +
			"</form>"
	} else if tab == "activity" {
		html += "<div style=\"max-width:600px\"><p>No activity yet.</p></div>"
	} else if tab == "invoices" {
		html += "<div style=\"margin-bottom:20px\"><a href=\"/client/" + clientID + "/invoice/new\" class=\"btn btn-primary\">Create Invoice</a></div>"
		if len(invoices) == 0 {
			html += "<div class=\"empty\">No invoices yet. <a href=\"/client/" + clientID + "/invoice/new\">Create one now</a></div>"
		} else {
			html += "<table><thead><tr><th>Invoice</th><th>Amount</th><th>Status</th><th>Created</th><th>Actions</th></tr></thead><tbody>"
			for _, inv := range invoices {
				html += "<tr><td><a href=\"/invoice/" + inv.ID + "\" style=\"color:#0066cc;text-decoration:none\">#" + inv.ID + "</a></td><td>$" + formatCurrency(inv.Total) + "</td><td>" + inv.Status + "</td><td>" + formatDate(inv.CreatedAt) + "</td><td><a href=\"/invoice/" + inv.ID + "\" class=\"btn btn-primary\" style=\"font-size:12px\">View</a></td></tr>"
			}
			html += "</tbody></table>"
		}
	}

	html += "</div></body></html>"
	fmt.Fprint(w, html)
}

func handleCreateInvoice(w http.ResponseWriter, r *http.Request, clientID string, client *store.Client) {
	if r.Method == "POST" {
		r.ParseMultipartForm(1024 * 1024)
		var lineItems []store.LineItem
		var total float64

		for i := 0; i < 10; i++ {
			desc := r.FormValue("description_" + strconv.Itoa(i))
			amtStr := r.FormValue("amount_" + strconv.Itoa(i))
			if desc == "" || amtStr == "" {
				continue
			}
			amt, _ := strconv.ParseFloat(amtStr, 64)
			lineItems = append(lineItems, store.LineItem{Description: desc, Amount: amt})
			total += amt
		}

		if len(lineItems) > 0 {
			inv := &store.Invoice{ClientID: clientID, Status: "draft", Total: total, LineItems: lineItems}
			s.CreateInvoice(inv)
			showToast(w, "Invoice created successfully")
			http.Redirect(w, r, "/invoice/"+inv.ID, http.StatusSeeOther)
			return
		}
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	html := "<!DOCTYPE html><html><head><title>Create Invoice</title><style>" +
		"body{font-family:sans-serif;margin:0;background:#f5f5f5}" +
		".header{background:#fff;padding:20px;border-bottom:1px solid #ddd}" +
		".content{max-width:600px;margin:0 auto;padding:20px}" +
		".back{margin-bottom:20px}" +
		".back a{color:#0066cc;text-decoration:none;font-size:14px}" +
		".breadcrumb{margin-bottom:20px;font-size:14px}" +
		".breadcrumb a{color:#0066cc;text-decoration:none}" +
		".form-group{margin-bottom:16px}" +
		"label{display:block;margin-bottom:4px;font-weight:500}" +
		"input[type=text],input[type=number]{width:100%;padding:8px;border:1px solid #ddd;border-radius:4px;box-sizing:border-box}" +
		".btn{padding:8px 16px;border:none;border-radius:4px;cursor:pointer;font-size:14px}" +
		".btn-primary{background:#0066cc;color:#fff}" +
		".btn-primary:hover{background:#0052a3}" +
		".line-items{margin:20px 0}" +
		".line-item{background:#fff;padding:16px;margin-bottom:12px;border-radius:4px;border:1px solid #ddd}" +
		"</style></head><body>" +
		"<div class=\"header\"><h1>Create Invoice</h1></div>" +
		"<div class=\"content\">" +
		"<div class=\"back\"><a href=\"/client/" + clientID + "\">← Back to Client</a></div>" +
		"<div class=\"breadcrumb\"><a href=\"/\">Dashboard</a> > <a href=\"/client/" + clientID + "\">" + client.Name + "</a> > Create Invoice</div>" +
		"<form method=\"POST\">" +
		"<div class=\"form-group\"><label>Bill To: " + client.Name + "</label></div>" +
		"<div class=\"line-items\" id=\"lineItems\">"

	for i := 0; i < 3; i++ {
		html += "<div class=\"line-item\"><div class=\"form-group\"><label>Description</label><input type=\"text\" name=\"description_" + strconv.Itoa(i) + "\" placeholder=\"Item or service\"></div>" +
			"<div class=\"form-group\"><label>Amount</label><input type=\"number\" name=\"amount_" + strconv.Itoa(i) + "\" step=\"0.01\" placeholder=\"0.00\"></div></div>"
	}

	html += "</div>" +
		"<div style=\"display:flex;gap:8px;margin-top:24px\"><button type=\"submit\" class=\"btn btn-primary\">Create Invoice</button><a href=\"/client/" + clientID + "\" class=\"btn\" style=\"background:#e5e5e5\">Cancel</a></div>" +
		"</form></div></body></html>"

	fmt.Fprint(w, html)
}

func handleInvoiceRoute(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/invoice/"), "/")
	invoiceID := parts[0]

	if invoiceID == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	invoice, ok := s.GetInvoice(invoiceID)
	if !ok {
		http.Error(w, "Invoice not found", http.StatusNotFound)
		return
	}

	client, _ := s.GetClient(invoice.ClientID)

	if len(parts) > 1 && parts[1] == "print" {
		handlePrintInvoice(w, r, invoice, client)
		return
	}

	if r.Method == "POST" {
		action := r.FormValue("action")
		if action == "send" {
			invoice.Status = "sent"
			s.UpdateInvoice(invoiceID, invoice)
		} else if action == "pay" {
			invoice.Status = "paid"
			s.UpdateInvoice(invoiceID, invoice)
		} else if action == "void" {
			invoice.Status = "void"
			s.UpdateInvoice(invoiceID, invoice)
		}
		showToast(w, "Invoice updated")
		http.Redirect(w, r, "/invoice/"+invoiceID, http.StatusSeeOther)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	html := "<!DOCTYPE html><html><head><title>Invoice #" + invoiceID + "</title><style>" +
		"body{font-family:sans-serif;margin:0;background:#f5f5f5}" +
		".header{background:#fff;padding:20px;border-bottom:1px solid #ddd}" +
		".content{max-width:800px;margin:0 auto;padding:20px}" +
		".back{margin-bottom:20px}" +
		".back a{color:#0066cc;text-decoration:none;font-size:14px}" +
		".breadcrumb{margin-bottom:20px;font-size:14px}" +
		".breadcrumb a{color:#0066cc;text-decoration:none}" +
		".invoice-box{background:#fff;padding:30px;border-radius:4px;margin-bottom:20px}" +
		"table{width:100%;border-collapse:collapse}" +
		"th,td{padding:12px;text-align:left;border-bottom:1px solid #ddd}" +
		"th{background:#f9f9f9;font-weight:600}" +
		".total{font-weight:600;font-size:18px}" +
		".btn{padding:8px 16px;border:none;border-radius:4px;cursor:pointer;font-size:14px;display:inline-block;margin-right:8px}" +
		".btn-primary{background:#0066cc;color:#fff}" +
		".btn-primary:hover{background:#0052a3}" +
		".btn-secondary{background:#e5e5e5}" +
		".status-draft{color:#f59e0b}" +
		".status-sent{color:#3b82f6}" +
		".status-paid{color:#10b981}" +
		".status-void{color:#6b7280}" +
		".actions{margin-top:20px}" +
		"</style></head><body>" +
		"<div class=\"header\"><h1>Invoice #" + invoiceID + "</h1></div>" +
		"<div class=\"content\">" +
		"<div class=\"back\"><a href=\"/client/" + invoice.ClientID + "\">← Back to Client</a></div>" +
		"<div class=\"breadcrumb\"><a href=\"/\">Dashboard</a> > <a href=\"/client/" + invoice.ClientID + "\">" + client.Name + "</a> > Invoice #" + invoiceID + "</div>" +
		"<div class=\"invoice-box\">" +
		"<div style=\"display:flex;justify-content:space-between;margin-bottom:20px\"><div><h2 style=\"margin:0\">Invoice</h2></div><div style=\"text-align:right\"><div class=\"status-" + invoice.Status + "\"><strong>Status: " + strings.ToUpper(invoice.Status) + "</strong></div></div></div>" +
		"<div style=\"margin-bottom:20px\"><strong>Bill To:</strong><br>" + client.Name + "<br>" + client.Email + "</div>" +
		"<div style=\"margin-bottom:20px\"><strong>Date:</strong> " + formatDate(invoice.CreatedAt) + "</div>" +
		"<table><thead><tr><th>Description</th><th style=\"text-align:right\">Amount</th></tr></thead><tbody>"

	for _, item := range invoice.LineItems {
		html += "<tr><td>" + item.Description + "</td><td style=\"text-align:right\">$" + formatCurrency(item.Amount) + "</td></tr>"
	}

	html += "</tbody></table>" +
		"<div style=\"text-align:right;margin-top:20px;font-size:18px\"><strong>Total: $" + formatCurrency(invoice.Total) + "</strong></div>" +
		"<div class=\"actions\">" +
		"<form method=\"POST\" style=\"display:inline\">" +
		"<input type=\"hidden\" name=\"action\" value=\"send\"><button type=\"submit\" class=\"btn btn-primary\">Send</button>" +
		"</form>" +
		"<form method=\"POST\" style=\"display:inline\">" +
		"<input type=\"hidden\" name=\"action\" value=\"pay\"><button type=\"submit\" class=\"btn btn-primary\">Mark Paid</button>" +
		"</form>" +
		"<form method=\"POST\" style=\"display:inline\">" +
		"<input type=\"hidden\" name=\"action\" value=\"void\"><button type=\"submit\" class=\"btn btn-secondary\">Void</button>" +
		"</form>" +
		"<a href=\"/invoice/" + invoiceID + "/print\" class=\"btn btn-secondary\">Print</a>" +
		"</div>" +
		"</div></div></body></html>"

	fmt.Fprint(w, html)
}

func handlePrintInvoice(w http.ResponseWriter, r *http.Request, invoice *store.Invoice, client *store.Client) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	html := "<!DOCTYPE html><html><head><title>Invoice #" + invoice.ID + "</title><style>" +
		"body{font-family:serif;margin:0;padding:20px}" +
		"@media print{body{padding:0}}" +
		".invoice{max-width:800px;margin:0 auto;padding:40px;background:#fff}" +
		"h1{margin:0 0 20px 0;font-size:32px}" +
		"table{width:100%;border-collapse:collapse;margin:20px 0}" +
		"th,td{padding:12px;text-align:left;border-bottom:1px solid #000}" +
		".total-row{font-weight:bold;font-size:18px}" +
		"</style></head><body>" +
		"<div class=\"invoice\">" +
		"<h1>INVOICE</h1>" +
		"<div style=\"margin-bottom:30px\"><strong>Invoice Number:</strong> " + invoice.ID + "</div>" +
		"<div style=\"margin-bottom:30px\"><strong>Date:</strong> " + formatDate(invoice.CreatedAt) + "</div>" +
		"<div style=\"margin-bottom:30px\"><strong>Bill To:</strong><br>" + client.Name + "<br>" + client.Email + "</div>" +
		"<table><thead><tr><th>Description</th><th style=\"text-align:right\">Amount</th></tr></thead><tbody>"

	for _, item := range invoice.LineItems {
		html += "<tr><td>" + item.Description + "</td><td style=\"text-align:right\">$" + formatCurrency(item.Amount) + "</td></tr>"
	}

	html += "<tr class=\"total-row\"><td style=\"text-align:right\">TOTAL</td><td style=\"text-align:right\">$" + formatCurrency(invoice.Total) + "</td></tr>" +
		"</tbody></table>" +
		"<div style=\"margin-top:40px;border-top:1px solid #000;padding-top:20px;text-align:center\"><p>Thank you for your business</p></div>" +
		"</div>" +
		"<script>window.print();</script>" +
		"</body></html>"

	fmt.Fprint(w, html)
}

func showToast(w http.ResponseWriter, message string) {
	// Toast is handled via redirect, user can add via response header or client-side
	w.Header().Set("X-Toast-Message", message)
}

func formatDate(t time.Time) string {
	months := []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}
	return months[t.Month()-1] + " " + strconv.Itoa(t.Day()) + ", " + strconv.Itoa(t.Year())
}

func formatCurrency(f float64) string {
	return strconv.FormatFloat(f, 'f', 2, 64)
}

func ifActive(condition bool, text string) string {
	if condition {
		return text
	}
	return ""
}
