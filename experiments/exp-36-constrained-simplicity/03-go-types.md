Here's the **exact Go types** for the CRM, keeping them simple but complete as per your requirements.

---

### ✅ `store/types.go`

```go
package store

import (
	"time"
)

// Client represents a client in the CRM
type Client struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Company   string    `json:"company"`
	Address   string    `json:"address"` // Single textarea field
	Notes     string    `json:"notes"`   // Textarea field
	CreatedAt time.Time `json:"created_at"`
}

// Activity represents an activity logged for a client
type Activity struct {
	ID          string    `json:"id"`
	ClientID    string    `json:"client_id"`
	Type        string    `json:"type"` // "Call", "Email", "Meeting", "Note"
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

// LineItem represents a line item in an invoice
type LineItem struct {
	Description string  `json:"description"`
	Quantity    int     `json:"quantity"`
	Rate        float64 `json:"rate"`
	Amount      float64 `json:"amount"`
}

// Invoice represents an invoice
type Invoice struct {
	ID        string      `json:"id"`
	ClientID  string      `json:"client_id"`
	Number    string      `json:"number"`
	Items     []LineItem  `json:"items"`
	Status    string      `json:"status"` // "draft", "sent", "paid", "void"
	DueDate   time.Time   `json:"due_date"`
	CreatedAt time.Time   `json:"created_at"`
}
```

---

### ✅ `store/store.go`

```go
package store

import (
	"time"
)

type Store interface {
	// Client CRUD
	CreateClient(client *Client) error
	GetClient(id string) (*Client, error)
	UpdateClient(id string, client *Client) error
	DeleteClient(id string) error
	SearchClients(query string) ([]*Client, error)

	// Activity CRUD
	CreateActivity(activity *Activity) error
	GetActivity(id string) (*Activity, error)
	UpdateActivity(id string, activity *Activity) error
	DeleteActivity(id string) error
	GetClientActivities(clientID string) ([]*Activity, error)

	// Invoice CRUD
	CreateInvoice(invoice *Invoice) error
	GetInvoice(id string) (*Invoice, error)
	UpdateInvoice(id string, invoice *Invoice) error
	DeleteInvoice(id string) error
	GetClientInvoices(clientID string) ([]*Invoice, error)

	// Invoice Status Updates
	MarkInvoicePaid(id string) error
	VoidInvoice(id string) error
}
```

---

### ✅ Example in-memory implementation (optional)

If you're using an in-memory store for simplicity:

```go
package store

import (
	"sync"
	"time"
)

type InMemoryStore struct {
	clients   map[string]*Client
	activities map[string]*Activity
	invoices  map[string]*Invoice
	mu        sync.RWMutex
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		clients:   make(map[string]*Client),
		activities: make(map[string]*Activity),
		invoices:  make(map[string]*Invoice),
	}
}

// Implement all Store methods here...
```

Let me know if you want the full in-memory implementation or a version using a database like SQLite or PostgreSQL.