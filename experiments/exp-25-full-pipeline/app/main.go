package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"invoicegen/store"
)

var invoiceStore *store.Store

func init() {
	var err error
	invoiceStore, err = store.NewStore("./data")
	if err != nil {
		log.Fatalf("Failed to initialize store: %v", err)
	}
}

func main() {
	mux := http.NewServeMux()

	// Health check
	mux.HandleFunc("/health", handleHealth)

	// Dashboard and forms
	mux.HandleFunc("/", handleDashboard)
	mux.HandleFunc("/invoices/new", handleNewInvoiceForm)

	// API endpoints
	mux.HandleFunc("POST /api/invoices", handleCreateInvoice)
	mux.HandleFunc("GET /api/invoices", handleListInvoices)
	mux.HandleFunc("GET /api/invoices/{id}", handleGetInvoice)
	mux.HandleFunc("PATCH /api/invoices/{id}", handleUpdateInvoice)
	mux.HandleFunc("DELETE /api/invoices/{id}", handleDeleteInvoice)
	mux.HandleFunc("POST /api/invoices/{id}/send", handleSendInvoice)
	mux.HandleFunc("GET /api/invoices/{id}/preview", handlePreviewInvoice)

	log.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
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

	invoices, err := invoiceStore.ListInvoices()
	if err != nil {
		http.Error(w, "Failed to load invoices", http.StatusInternalServerError)
		return
	}

	var totalAmount float64
	var paidAmount float64
	draftCount := 0
	sentCount := 0
	paidCount := 0

	for _, inv := range invoices {
		totalAmount += inv.Total
		if inv.Status == store.StatusPaid {
			paidAmount += inv.Total
			paidCount++
		} else if inv.Status == store.StatusDraft {
			draftCount++
		} else if inv.Status == store.StatusSent {
			sentCount++
		}
	}

	html := `<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>Invoice Generator</title>
	<style>
		* { margin: 0; padding: 0; box-sizing: border-box; }
		body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; background: #f8f9fa; color: #333; }
		.container { max-width: 1200px; margin: 0 auto; padding: 40px 20px; }
		header { margin-bottom: 40px; }
		h1 { font-size: 32px; font-weight: 600; margin-bottom: 8px; }
		.header-actions { display: flex; gap: 12px; margin-top: 20px; }
		.btn { padding: 10px 20px; background: #2563eb; color: white; border: none; border-radius: 6px; cursor: pointer; text-decoration: none; font-size: 14px; font-weight: 500; display: inline-flex; align-items: center; gap: 8px; }
		.btn:hover { background: #1d4ed8; }
		.btn-secondary { background: #6b7280; }
		.btn-secondary:hover { background: #4b5563; }
		.stats { display: grid; grid-template-columns: repeat(auto-fit, minmax(250px, 1fr)); gap: 20px; margin-bottom: 40px; }
		.stat-card { background: white; padding: 24px; border-radius: 8px; border: 1px solid #e5e7eb; box-shadow: 0 1px 3px rgba(0,0,0,0.05); }
		.stat-label { font-size: 13px; font-weight: 600; color: #6b7280; text-transform: uppercase; letter-spacing: 0.5px; margin-bottom: 8px; }
		.stat-value { font-size: 28px; font-weight: 700; color: #111827; }
		.stat-detail { font-size: 12px; color: #9ca3af; margin-top: 8px; }
		.invoices-section { background: white; border-radius: 8px; border: 1px solid #e5e7eb; box-shadow: 0 1px 3px rgba(0,0,0,0.05); overflow: hidden; }
		.section-header { padding: 24px; border-bottom: 1px solid #e5e7eb; }
		.section-header h2 { font-size: 18px; font-weight: 600; }
		table { width: 100%; border-collapse: collapse; }
		th { background: #f9fafb; padding: 12px 24px; text-align: left; font-size: 12px; font-weight: 600; color: #6b7280; text-transform: uppercase; letter-spacing: 0.5px; border-bottom: 1px solid #e5e7eb; }
		td { padding: 16px 24px; border-bottom: 1px solid #e5e7eb; }
		tr:last-child td { border-bottom: none; }
		tr:hover { background: #f9fafb; }
		.status-badge { display: inline-block; padding: 4px 12px; border-radius: 12px; font-size: 12px; font-weight: 500; }
		.status-draft { background: #fef3c7; color: #92400e; }
		.status-sent { background: #dbeafe; color: #1e40af; }
		.status-paid { background: #dcfce7; color: #166534; }
		.status-overdue { background: #fee2e2; color: #991b1b; }
		.invoice-number { font-weight: 500; color: #2563eb; text-decoration: none; }
		.invoice-number:hover { text-decoration: underline; }
		.client-name { font-weight: 500; }
		.amount { text-align: right; font-weight: 500; }
		.actions { text-align: center; }
		.action-btn { padding: 6px 12px; margin: 0 2px; border: none; border-radius: 4px; cursor: pointer; font-size: 12px; text-decoration: none; color: white; background: #6b7280; }
		.action-btn:hover { background: #4b5563; }
		.action-btn-view { background: #3b82f6; }
		.action-btn-view:hover { background: #2563eb; }
		.empty-state { text-align: center; padding: 60px 20px; color: #6b7280; }
		.empty-state-icon { font-size: 48px; margin-bottom: 16px; }
		.empty-state h3 { font-size: 18px; font-weight: 600; margin-bottom: 8px; color: #374151; }
	</style>
</head>
<body>
	<div class="container">
		<header>
			<h1>Invoice Generator</h1>
			<div class="header-actions">
				<a href="/invoices/new" class="btn">+ New Invoice</a>
			</div>
		</header>

		<div class="stats">
			<div class="stat-card">
				<div class="stat-label">Total Invoiced</div>
				<div class="stat-value">$` + fmt.Sprintf("%.2f", totalAmount) + `</div>
				<div class="stat-detail">across ` + fmt.Sprintf("%d", len(invoices)) + ` invoices</div>
			</div>
			<div class="stat-card">
				<div class="stat-label">Paid Amount</div>
				<div class="stat-value">$` + fmt.Sprintf("%.2f", paidAmount) + `</div>
				<div class="stat-detail">from ` + fmt.Sprintf("%d", paidCount) + ` paid invoices</div>
			</div>
			<div class="stat-card">
				<div class="stat-label">Draft</div>
				<div class="stat-value">` + fmt.Sprintf("%d", draftCount) + `</div>
				<div class="stat-detail">awaiting details</div>
			</div>
			<div class="stat-card">
				<div class="stat-label">Sent</div>
				<div class="stat-value">` + fmt.Sprintf("%d", sentCount) + `</div>
				<div class="stat-detail">awaiting payment</div>
			</div>
		</div>

		<div class="invoices-section">
			<div class="section-header">
				<h2>Recent Invoices</h2>
			</div>
			` + renderInvoicesTable(invoices) + `
		</div>
	</div>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(html))
}

func renderInvoicesTable(invoices []*store.Invoice) string {
	if len(invoices) == 0 {
		return `<div class="empty-state">
			<div class="empty-state-icon">📋</div>
			<h3>No invoices yet</h3>
			<p>Create your first invoice to get started</p>
		</div>`
	}

	html := `<table>
		<thead>
			<tr>
				<th>Invoice #</th>
				<th>Client</th>
				<th>Issue Date</th>
				<th>Due Date</th>
				<th>Status</th>
				<th class="amount">Amount</th>
				<th class="actions">Actions</th>
			</tr>
		</thead>
		<tbody>`

	for _, inv := range invoices {
		statusClass := "status-" + strings.ToLower(inv.Status)
		statusTitle := strings.ToUpper(inv.Status[:1]) + inv.Status[1:]
		html += `<tr>
			<td><a href="/api/invoices/` + inv.ID + `/preview" class="invoice-number">` + inv.Number + `</a></td>
			<td class="client-name">` + inv.ClientID + `</td>
			<td>` + inv.IssueDate.Format("Jan 2, 2006") + `</td>
			<td>` + inv.DueDate.Format("Jan 2, 2006") + `</td>
			<td><span class="status-badge ` + statusClass + `">` + statusTitle + `</span></td>
			<td class="amount">$` + fmt.Sprintf("%.2f", inv.Total) + `</td>
			<td class="actions">
				<a href="/api/invoices/` + inv.ID + `/preview" class="action-btn action-btn-view">View</a>
			</td>
		</tr>`
	}

	html += `</tbody></table>`
	return html
}

func handleNewInvoiceForm(w http.ResponseWriter, r *http.Request) {
	html := `<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>Create Invoice</title>
	<style>
		* { margin: 0; padding: 0; box-sizing: border-box; }
		body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; background: #f8f9fa; color: #333; }
		.container { max-width: 900px; margin: 0 auto; padding: 40px 20px; }
		h1 { font-size: 28px; font-weight: 600; margin-bottom: 32px; }
		.form-section { background: white; padding: 28px; border-radius: 8px; border: 1px solid #e5e7eb; margin-bottom: 24px; box-shadow: 0 1px 3px rgba(0,0,0,0.05); }
		.section-title { font-size: 14px; font-weight: 600; text-transform: uppercase; color: #6b7280; margin-bottom: 20px; letter-spacing: 0.5px; }
		.form-group { margin-bottom: 20px; }
		label { display: block; font-size: 14px; font-weight: 500; margin-bottom: 6px; color: #374151; }
		input, textarea, select { width: 100%; padding: 10px 12px; border: 1px solid #d1d5db; border-radius: 6px; font-size: 14px; font-family: inherit; }
		input:focus, textarea:focus, select:focus { outline: none; border-color: #2563eb; box-shadow: 0 0 0 3px rgba(37, 99, 235, 0.1); }
		.form-row { display: grid; grid-template-columns: 1fr 1fr; gap: 20px; }
		textarea { resize: vertical; }
		.line-items { margin-top: 24px; }
		.line-item { display: grid; grid-template-columns: 2fr 1fr 1fr 60px; gap: 12px; align-items: end; margin-bottom: 12px; }
		.line-item input { }
		.remove-btn { padding: 10px 16px; background: #ef4444; color: white; border: none; border-radius: 6px; cursor: pointer; font-size: 14px; }
		.remove-btn:hover { background: #dc2626; }
		.add-line-btn { padding: 10px 16px; background: #10b981; color: white; border: none; border-radius: 6px; cursor: pointer; font-size: 14px; margin-top: 12px; }
		.add-line-btn:hover { background: #059669; }
		.totals { background: #f9fafb; padding: 16px; border-radius: 6px; margin-top: 20px; }
		.total-row { display: flex; justify-content: space-between; margin-bottom: 12px; font-size: 14px; }
		.total-row:last-child { margin-bottom: 0; }
		.total-row.final { font-weight: 600; font-size: 16px; color: #111827; border-top: 2px solid #e5e7eb; padding-top: 12px; }
		.form-actions { display: flex; gap: 12px; }
		.btn { padding: 12px 24px; border: none; border-radius: 6px; cursor: pointer; font-size: 14px; font-weight: 500; }
		.btn-primary { background: #2563eb; color: white; }
		.btn-primary:hover { background: #1d4ed8; }
		.btn-secondary { background: #6b7280; color: white; }
		.btn-secondary:hover { background: #4b5563; }
		.error { color: #dc2626; font-size: 13px; margin-top: 4px; }
	</style>
</head>
<body>
	<div class="container">
		<h1>Create Invoice</h1>

		<form id="invoiceForm">
			<div class="form-section">
				<div class="section-title">Basic Information</div>
				<div class="form-row">
					<div class="form-group">
						<label for="clientId">Client ID *</label>
						<input type="text" id="clientId" name="client_id" required>
					</div>
					<div class="form-group">
						<label for="taxRate">Tax Rate (%) *</label>
						<input type="number" id="taxRate" name="tax_rate" value="0" step="0.01" required>
					</div>
				</div>
				<div class="form-row">
					<div class="form-group">
						<label for="issueDate">Issue Date *</label>
						<input type="date" id="issueDate" name="issue_date" required>
					</div>
					<div class="form-group">
						<label for="dueDate">Due Date *</label>
						<input type="date" id="dueDate" name="due_date" required>
					</div>
				</div>
				<div class="form-group">
					<label for="notes">Notes</label>
					<textarea id="notes" name="notes" rows="4"></textarea>
				</div>
			</div>

			<div class="form-section">
				<div class="section-title">Line Items</div>
				<div id="lineItems">
					<div class="line-item">
						<input type="text" placeholder="Description" class="description">
						<input type="number" placeholder="Qty" class="quantity" value="1" min="1" step="1">
						<input type="number" placeholder="Unit Price" class="unitPrice" value="0" min="0" step="0.01">
						<button type="button" class="remove-btn" style="display:none;">Remove</button>
					</div>
				</div>
				<button type="button" class="add-line-btn" onclick="addLineItem()">+ Add Line Item</button>

				<div class="totals">
					<div class="total-row">
						<span>Subtotal</span>
						<span id="subtotal">$0.00</span>
					</div>
					<div class="total-row">
						<span id="taxLabel">Tax (0%)</span>
						<span id="taxAmount">$0.00</span>
					</div>
					<div class="total-row final">
						<span>Total</span>
						<span id="total">$0.00</span>
					</div>
				</div>
			</div>

			<div class="form-actions">
				<button type="submit" class="btn btn-primary">Create Invoice</button>
				<a href="/" class="btn btn-secondary" style="text-decoration:none;">Cancel</a>
			</div>
		</form>
	</div>

	<script>
		function getLineItems() {
			var items = [];
			var lines = document.querySelectorAll('#lineItems .line-item');
			lines.forEach(function(line) {
				var desc = line.querySelector('.description').value;
				var qty = parseInt(line.querySelector('.quantity').value) || 0;
				var price = parseFloat(line.querySelector('.unitPrice').value) || 0;
				if (desc || qty > 0 || price > 0) {
					items.push({
						description: desc,
						quantity: qty,
						unit_price: price,
						total: qty * price
					});
				}
			});
			return items;
		}

		function calculateTotals() {
			var items = getLineItems();
			var subtotal = 0;
			items.forEach(function(item) {
				subtotal += item.quantity * item.unit_price;
			});
			var taxRate = parseFloat(document.getElementById('taxRate').value) || 0;
			var taxAmount = subtotal * (taxRate / 100);
			var total = subtotal + taxAmount;

			document.getElementById('subtotal').textContent = '$' + subtotal.toFixed(2);
			document.getElementById('taxAmount').textContent = '$' + taxAmount.toFixed(2);
			document.getElementById('taxLabel').textContent = 'Tax (' + taxRate.toFixed(2) + '%)';
			document.getElementById('total').textContent = '$' + total.toFixed(2);
		}

		function addLineItem() {
			var container = document.getElementById('lineItems');
			var newItem = document.createElement('div');
			newItem.className = 'line-item';
			newItem.innerHTML = '<input type="text" placeholder="Description" class="description">' +
				'<input type="number" placeholder="Qty" class="quantity" value="1" min="1" step="1">' +
				'<input type="number" placeholder="Unit Price" class="unitPrice" value="0" min="0" step="0.01">' +
				'<button type="button" class="remove-btn">Remove</button>';
			container.appendChild(newItem);

			var removeBtn = newItem.querySelector('.remove-btn');
			removeBtn.onclick = function() {
				newItem.remove();
				calculateTotals();
			};

			newItem.querySelectorAll('input').forEach(function(input) {
				input.addEventListener('input', calculateTotals);
			});

			calculateTotals();
		}

		function setTodayDate() {
			var today = new Date().toISOString().split('T')[0];
			document.getElementById('issueDate').value = today;
			var dueDate = new Date();
			dueDate.setDate(dueDate.getDate() + 30);
			document.getElementById('dueDate').value = dueDate.toISOString().split('T')[0];
		}

		document.addEventListener('DOMContentLoaded', function() {
			setTodayDate();
			document.getElementById('taxRate').addEventListener('input', calculateTotals);
			document.querySelectorAll('#lineItems input').forEach(function(input) {
				input.addEventListener('input', calculateTotals);
			});
			calculateTotals();
		});

		document.getElementById('invoiceForm').addEventListener('submit', function(e) {
			e.preventDefault();
			var items = getLineItems();
			var payload = {
				client_id: document.getElementById('clientId').value,
				issue_date: new Date(document.getElementById('issueDate').value),
				due_date: new Date(document.getElementById('dueDate').value),
				tax_rate: parseFloat(document.getElementById('taxRate').value) || 0,
				line_items: items,
				notes: document.getElementById('notes').value,
				status: 'draft'
			};

			fetch('/api/invoices', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify(payload)
			})
			.then(function(r) { return r.json(); })
			.then(function(data) {
				if (data.error) {
					alert('Error: ' + data.error);
				} else {
					window.location.href = '/api/invoices/' + data.id + '/preview';
				}
			})
			.catch(function(e) { alert('Failed to create invoice: ' + e); });
		});
	</script>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(html))
}

func handleCreateInvoice(w http.ResponseWriter, r *http.Request) {
	var inv store.Invoice
	err := json.NewDecoder(r.Body).Decode(&inv)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid JSON"})
		return
	}

	err = invoiceStore.CreateInvoice(&inv)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(inv)
}

func handleListInvoices(w http.ResponseWriter, r *http.Request) {
	invoices, err := invoiceStore.ListInvoices()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(invoices)
}

func getInvoiceID(r *http.Request) string {
	id := strings.TrimPrefix(r.URL.Path, "/api/invoices/")
	parts := strings.Split(id, "/")
	return parts[0]
}

func handleGetInvoice(w http.ResponseWriter, r *http.Request) {
	id := getInvoiceID(r)
	inv, err := invoiceStore.GetInvoice(id)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invoice not found"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(inv)
}

func handleUpdateInvoice(w http.ResponseWriter, r *http.Request) {
	id := getInvoiceID(r)
	var updates map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid JSON"})
		return
	}

	inv, err := invoiceStore.GetInvoice(id)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invoice not found"})
		return
	}

	data, _ := json.Marshal(updates)
	json.Unmarshal(data, inv)

	err = invoiceStore.UpdateInvoice(inv)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(inv)
}

func handleDeleteInvoice(w http.ResponseWriter, r *http.Request) {
	id := getInvoiceID(r)
	err := invoiceStore.DeleteInvoice(id)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invoice not found"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "deleted"})
}

func handleSendInvoice(w http.ResponseWriter, r *http.Request) {
	id := getInvoiceID(r)
	inv, err := invoiceStore.GetInvoice(id)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invoice not found"})
		return
	}

	inv.Status = store.StatusSent
	err = invoiceStore.UpdateInvoice(inv)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(inv)
}

func handlePreviewInvoice(w http.ResponseWriter, r *http.Request) {
	id := getInvoiceID(r)
	inv, err := invoiceStore.GetInvoice(id)
	if err != nil {
		http.Error(w, "Invoice not found", http.StatusNotFound)
		return
	}

	html := generateInvoiceHTML(inv)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(html))
}

func generateInvoiceHTML(inv *store.Invoice) string {
	lineItemsHTML := ""
	for _, item := range inv.LineItems {
		lineItemsHTML += `<tr>
			<td>` + item.Description + `</td>
			<td style="text-align: right;">` + fmt.Sprintf("%d", item.Quantity) + `</td>
			<td style="text-align: right;">$` + fmt.Sprintf("%.2f", item.UnitPrice) + `</td>
			<td style="text-align: right;">$` + fmt.Sprintf("%.2f", item.Total) + `</td>
		</tr>`
	}

	html := `<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>Invoice ` + inv.Number + `</title>
	<style>
		* { margin: 0; padding: 0; box-sizing: border-box; }
		body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; background: #f8f9fa; color: #333; padding: 40px 20px; }
		.container { max-width: 900px; margin: 0 auto; background: white; padding: 60px; border: 1px solid #e5e7eb; border-radius: 8px; box-shadow: 0 1px 3px rgba(0,0,0,0.05); }
		header { margin-bottom: 60px; display: flex; justify-content: space-between; align-items: start; border-bottom: 2px solid #e5e7eb; padding-bottom: 30px; }
		.company { }
		.company h1 { font-size: 24px; font-weight: 700; margin-bottom: 4px; }
		.invoice-number { text-align: right; }
		.invoice-number h2 { font-size: 28px; font-weight: 700; color: #2563eb; margin-bottom: 8px; }
		.invoice-meta { font-size: 13px; color: #6b7280; line-height: 1.6; }
		.invoice-meta strong { color: #374151; }
		.details { display: grid; grid-template-columns: 1fr 1fr; gap: 40px; margin-bottom: 40px; }
		.detail-section { }
		.detail-section h3 { font-size: 12px; font-weight: 600; color: #6b7280; text-transform: uppercase; letter-spacing: 0.5px; margin-bottom: 12px; }
		.detail-section p { font-size: 14px; line-height: 1.6; color: #374151; }
		.items { margin-bottom: 40px; }
		.items table { width: 100%; border-collapse: collapse; }
		.items th { background: #f9fafb; padding: 12px; text-align: left; font-size: 12px; font-weight: 600; color: #6b7280; border-bottom: 2px solid #e5e7eb; }
		.items td { padding: 12px; border-bottom: 1px solid #e5e7eb; font-size: 14px; }
		.items tr:last-child td { border-bottom: 2px solid #e5e7eb; }
		.totals { max-width: 400px; margin-left: auto; margin-bottom: 40px; }
		.total-row { display: flex; justify-content: space-between; padding: 12px 0; font-size: 14px; border-bottom: 1px solid #e5e7eb; }
		.total-row.final { font-size: 18px; font-weight: 700; border-bottom: 2px solid #2563eb; padding: 16px 0; color: #2563eb; }
		.status-badge { display: inline-block; padding: 6px 12px; background: #dbeafe; color: #1e40af; border-radius: 4px; font-size: 12px; font-weight: 600; text-transform: uppercase; letter-spacing: 0.5px; margin-bottom: 20px; }
		.notes { background: #f9fafb; padding: 20px; border-radius: 6px; border-left: 4px solid #2563eb; margin-bottom: 40px; }
		.notes h4 { font-size: 12px; font-weight: 600; color: #6b7280; text-transform: uppercase; margin-bottom: 8px; }
		.notes p { font-size: 13px; color: #374151; line-height: 1.6; }
		.actions { display: flex; gap: 12px; padding-top: 20px; border-top: 1px solid #e5e7eb; }
		.btn { padding: 10px 20px; border: none; border-radius: 6px; cursor: pointer; font-size: 14px; font-weight: 500; text-decoration: none; display: inline-block; }
		.btn-primary { background: #2563eb; color: white; }
		.btn-primary:hover { background: #1d4ed8; }
		.btn-secondary { background: #6b7280; color: white; }
		.btn-secondary:hover { background: #4b5563; }
		@media print { body { background: white; padding: 0; } .actions { display: none; } .container { box-shadow: none; border: none; } }
	</style>
</head>
<body>
	<div class="container">
		<header>
			<div class="company">
				<h1>INVOICE</h1>
			</div>
			<div class="invoice-number">
				<h2>` + inv.Number + `</h2>
				<div class="invoice-meta">
					<div><strong>Status:</strong> <span class="status-badge">` + strings.ToUpper(inv.Status) + `</span></div>
					<div><strong>Issue Date:</strong> ` + inv.IssueDate.Format("January 2, 2006") + `</div>
					<div><strong>Due Date:</strong> ` + inv.DueDate.Format("January 2, 2006") + `</div>
				</div>
			</div>
		</header>

		<div class="details">
			<div class="detail-section">
				<h3>Bill To</h3>
				<p>` + inv.ClientID + `</p>
			</div>
			<div class="detail-section">
				<h3>Invoice Details</h3>
				<p>
					<strong>Invoice ID:</strong> ` + inv.ID + `<br>
					<strong>Created:</strong> ` + inv.CreatedAt.Format("January 2, 2006") + `<br>
					<strong>Updated:</strong> ` + inv.UpdatedAt.Format("January 2, 2006") + `
				</p>
			</div>
		</div>

		<div class="items">
			<table>
				<thead>
					<tr>
						<th>Description</th>
						<th style="text-align: right; width: 80px;">Quantity</th>
						<th style="text-align: right; width: 100px;">Unit Price</th>
						<th style="text-align: right; width: 100px;">Amount</th>
					</tr>
				</thead>
				<tbody>
					` + lineItemsHTML + `
				</tbody>
			</table>
		</div>

		<div class="totals">
			<div class="total-row">
				<span>Subtotal</span>
				<span>$` + fmt.Sprintf("%.2f", inv.Subtotal) + `</span>
			</div>
			<div class="total-row">
				<span>Tax (` + fmt.Sprintf("%.2f", inv.TaxRate) + `%)</span>
				<span>$` + fmt.Sprintf("%.2f", inv.TaxAmount) + `</span>
			</div>
			<div class="total-row final">
				<span>Total</span>
				<span>$` + fmt.Sprintf("%.2f", inv.Total) + `</span>
			</div>
		</div>

		` + (func() string {
		if inv.Notes != "" {
			return `<div class="notes">
			<h4>Notes</h4>
			<p>` + inv.Notes + `</p>
		</div>`
		}
		return ""
	}()) + `

		<div class="actions">
			<a href="/" class="btn btn-secondary">Back to Dashboard</a>
			<button class="btn btn-primary" onclick="window.print()">Print</button>
		</div>
	</div>
</body>
</html>`

	return html
}
