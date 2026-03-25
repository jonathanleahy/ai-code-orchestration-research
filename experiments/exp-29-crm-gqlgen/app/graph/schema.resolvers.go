package graph

// THIS CODE WILL BE UPDATED WITH SCHEMA CHANGES. PREVIOUS IMPLEMENTATION FOR SCHEMA CHANGES WILL BE KEPT IN THE COMMENT SECTION. IMPLEMENTATION FOR UNCHANGED SCHEMA WILL BE KEPT.

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"
)

type Resolver struct {
	mu                   sync.RWMutex
	clients              map[string]*Client
	activities           map[string]*Activity
	invoices             map[string]*Invoice
	clientIDCounter      int64
	activityIDCounter    int64
	invoiceIDCounter     int64
}

// CreateClient is the resolver for the createClient field.
func (r *mutationResolver) CreateClient(ctx context.Context, name string, company string, email string, phone string, status ClientStatus) (*Client, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.clientIDCounter++
	id := fmt.Sprintf("%d", r.clientIDCounter)
	client := &Client{
		ID:        id,
		Name:      name,
		Company:   company,
		Email:     email,
		Phone:     phone,
		Status:    status,
		CreatedAt: time.Now(),
	}
	r.clients[id] = client
	return client, nil
}

// UpdateClient is the resolver for the updateClient field.
func (r *mutationResolver) UpdateClient(ctx context.Context, id string, name *string, company *string, email *string, phone *string, status *ClientStatus) (*Client, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	client, ok := r.clients[id]
	if !ok {
		return nil, nil
	}

	if name != nil {
		client.Name = *name
	}
	if company != nil {
		client.Company = *company
	}
	if email != nil {
		client.Email = *email
	}
	if phone != nil {
		client.Phone = *phone
	}
	if status != nil {
		client.Status = *status
	}

	return client, nil
}

// DeleteClient is the resolver for the deleteClient field.
func (r *mutationResolver) DeleteClient(ctx context.Context, id string) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.clients[id]; !ok {
		return false, nil
	}

	delete(r.clients, id)

	var activityIDs []string
	for actID, activity := range r.activities {
		if activity.ClientID == id {
			activityIDs = append(activityIDs, actID)
		}
	}
	for _, actID := range activityIDs {
		delete(r.activities, actID)
	}

	var invoiceIDs []string
	for invID, invoice := range r.invoices {
		if invoice.ClientID == id {
			invoiceIDs = append(invoiceIDs, invID)
		}
	}
	for _, invID := range invoiceIDs {
		delete(r.invoices, invID)
	}

	return true, nil
}

// CreateActivity is the resolver for the createActivity field.
func (r *mutationResolver) CreateActivity(ctx context.Context, clientID string, typeArg ActivityType, title string, description string) (*Activity, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.activityIDCounter++
	id := fmt.Sprintf("%d", r.activityIDCounter)
	activity := &Activity{
		ID:          id,
		ClientID:    clientID,
		Type:        typeArg,
		Title:       title,
		Description: description,
		Timestamp:   time.Now(),
	}
	r.activities[id] = activity
	return activity, nil
}

// DeleteActivity is the resolver for the deleteActivity field.
func (r *mutationResolver) DeleteActivity(ctx context.Context, id string) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.activities[id]; !ok {
		return false, nil
	}

	delete(r.activities, id)
	return true, nil
}

// CreateInvoice is the resolver for the createInvoice field.
func (r *mutationResolver) CreateInvoice(ctx context.Context, clientID string, number string, amount float64, description string, dueDate time.Time) (*Invoice, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.invoiceIDCounter++
	id := fmt.Sprintf("%d", r.invoiceIDCounter)
	invoice := &Invoice{
		ID:          id,
		ClientID:    clientID,
		Number:      number,
		Amount:      amount,
		Description: description,
		Status:      "draft",
		CreatedAt:   time.Now(),
		DueDate:     dueDate,
	}
	r.invoices[id] = invoice
	return invoice, nil
}

// UpdateInvoiceStatus is the resolver for the updateInvoiceStatus field.
func (r *mutationResolver) UpdateInvoiceStatus(ctx context.Context, id string, status string) (*Invoice, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	invoice, ok := r.invoices[id]
	if !ok {
		return nil, nil
	}

	invoice.Status = status
	return invoice, nil
}

// DeleteInvoice is the resolver for the deleteInvoice field.
func (r *mutationResolver) DeleteInvoice(ctx context.Context, id string) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.invoices[id]; !ok {
		return false, nil
	}

	delete(r.invoices, id)
	return true, nil
}

// Clients is the resolver for the clients field.
func (r *queryResolver) Clients(ctx context.Context, search *string) ([]*Client, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var results []*Client
	for _, client := range r.clients {
		if search != nil {
			searchStr := strings.ToLower(*search)
			if !strings.Contains(strings.ToLower(client.Name), searchStr) &&
				!strings.Contains(strings.ToLower(client.Company), searchStr) {
				continue
			}
		}
		results = append(results, client)
	}
	return results, nil
}

// Client is the resolver for the client field.
func (r *queryResolver) Client(ctx context.Context, id string) (*Client, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.clients[id], nil
}

// Activities is the resolver for the activities field.
func (r *queryResolver) Activities(ctx context.Context, clientID string) ([]*Activity, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var results []*Activity
	for _, activity := range r.activities {
		if activity.ClientID == clientID {
			results = append(results, activity)
		}
	}
	return results, nil
}

// Invoices is the resolver for the invoices field.
func (r *queryResolver) Invoices(ctx context.Context, clientID string) ([]*Invoice, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var results []*Invoice
	for _, invoice := range r.invoices {
		if invoice.ClientID == clientID {
			results = append(results, invoice)
		}
	}
	return results, nil
}

// Invoice is the resolver for the invoice field.
func (r *queryResolver) Invoice(ctx context.Context, id string) (*Invoice, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.invoices[id], nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

// NewResolver creates a new Resolver with initialized maps
func NewResolver() *Resolver {
	return &Resolver{
		clients:    make(map[string]*Client),
		activities: make(map[string]*Activity),
		invoices:   make(map[string]*Invoice),
	}
}
