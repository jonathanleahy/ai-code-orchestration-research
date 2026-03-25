package store

import "time"

type Client struct {
	ID        string
	Name      string
	Email     string
	Phone     string
	Company   string
	Address   string
	CreatedAt time.Time
}

type Activity struct {
	ID        string
	ClientID  string
	Type      string
	Notes     string
	CreatedAt time.Time
}

type InvoiceItem struct {
	Description string
	Quantity    int
	Rate        float64
	Amount      float64
}

type Invoice struct {
	ID        string
	ClientID  string
	Status    string
	Items     []InvoiceItem
	DueDate   time.Time
	CreatedAt time.Time
}
