Here's a **Go type definition** for your **Revised Freelancer CRM MVP**, structured with clean, reusable types that support:

- **CRUD operations**
- **Search functionality**
- **Invoice workflow**
- **SaaS UX patterns** (toasts, confirmations, breadcrumbs, etc.)

---

### 🧱 **Go Types Overview**

We'll define the following core types:

1. `Address` – string
2. `Status` – string (e.g., `"active"`, `"completed"`, `"pending"`)
3. `Client`
4. `Project`
5. `Invoice`
6. `LineItem`
7. `User` (optional but useful for multi-user support)
8. `DashboardData` – for dashboard overview

---

### ✅ **Go Code: Types**

```go
package crm

// Address represents a client's address
type Address string

// Status represents a project or invoice status
type Status string

const (
	StatusActive   Status = "active"
	StatusPending  Status = "pending"
	StatusCompleted Status = "completed"
	StatusPaid     Status = "paid"
)

// Client represents a client
type Client struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Address   Address   `json:"address"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Project represents a project
type Project struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	ClientID    string    `json:"client_id"`
	Status      Status    `json:"status"`
	DueDate     time.Time `json:"due_date"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// LineItem represents an item in an invoice
type LineItem struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	Quantity    int    `json:"quantity"`
	Rate        float64 `json:"rate"`
	Amount      float64 `json:"amount"`
}

// Invoice represents an invoice
type Invoice struct {
	ID           string     `json:"id"`
	InvoiceNumber string    `json:"invoice_number"`
	ClientID     string     `json:"client_id"`
	LineItems    []LineItem `json:"line_items"`
	Status       Status     `json:"status"`
	IssueDate    time.Time  `json:"issue_date"`
	DueDate      time.Time  `json:"due_date"`
	TotalAmount  float64    `json:"total_amount"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// User represents a user (optional for multi-user support)
type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// DashboardData represents data shown on the dashboard
type DashboardData struct {
	TotalClients   int `json:"total_clients"`
	TotalProjects  int `json:"total_projects"`
	TotalInvoices  int `json:"total_invoices"`
	RecentClients  []Client `json:"recent_clients"`
	UpcomingInvoices []Invoice `json:"upcoming_invoices"`
	ActiveProjects []Project `json:"active_projects"`
}
```

---

### 🧠 **SaaS UX Patterns in Go (Conceptual)**

While these are not direct Go types, they are **patterns** that can be implemented in your Go backend or frontend logic:

| Pattern | Implementation |
|--------|----------------|
| **Breadcrumbs** | Store navigation path in context or session |
| **Toasts** | Return toast messages in API responses or via WebSocket |
| **Confirmations** | Validate delete requests with confirmation flag |
| **Back Buttons** | Use navigation state or redirect logic |
| **Empty States** | Return empty list or message in API responses |

---

### 🔄 **CRUD + Search Example (Conceptual)**

You can define interfaces or services like:

```go
type ClientService interface {
	CreateClient(client *Client) error
	GetClient(id string) (*Client, error)
	UpdateClient(id string, client *Client) error
	DeleteClient(id string) error
	SearchClients(query string) ([]Client, error)
}

type InvoiceService interface {
	CreateInvoice(invoice *Invoice) error
	GetInvoice(id string) (*Invoice, error)
	UpdateInvoice(id string, invoice *Invoice) error
	DeleteInvoice(id string) error
	MarkAsPaid(id string) error
	SendInvoice(id string) error
}
```

---

### 🧪 **Example Usage**

```go
client := &Client{
	ID:        "c1",
	Name:      "John Doe",
	Email:     "john@example.com",
	Address:   "123 Main St",
	CreatedAt: time.Now(),
	UpdatedAt: time.Now(),
}

invoice := &Invoice{
	ID:            "i1",
	InvoiceNumber: "INV-001",
	ClientID:      "c1",
	LineItems: []LineItem{
		{
			ID:          "li1",
			Description: "Web Design",
			Quantity:    10,
			Rate:        50.0,
			Amount:      500.0,
		},
	},
	Status:    StatusPending,
	IssueDate: time.Now(),
	DueDate:   time.Now().AddDate(0, 1, 0),
	TotalAmount: 500.0,
	CreatedAt: time.Now(),
	UpdatedAt: time.Now(),
}
```

---

Let me know if you'd like:
- Full CRUD handlers (e.g., in Gin or Echo)
- JSON schema validation
- Database schema (PostgreSQL/SQLite)
- Frontend type definitions (TypeScript)