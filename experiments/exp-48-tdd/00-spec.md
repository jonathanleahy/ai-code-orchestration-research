CRM Store: Client CRUD + Activity Log + Simple Invoices

Types:
  Client: ID, Name, Email, Phone, CreatedAt
  Activity: ID, ClientID, Type (call/email/meeting/note), Description, CreatedAt
  Invoice: ID, ClientID, Number, Amount, Status (draft/sent/paid), DueDate, CreatedAt

Functions:
  NewStore() *Store
  AddClient(name, email, phone) *Client — reject empty name
  GetClient(id) *Client — nil if not found
  ListClients() []*Client
  UpdateClient(id, name, email, phone) *Client — nil if not found
  DeleteClient(id) bool — false if not found
  AddActivity(clientID, actType, desc) *Activity — nil if client not found
  GetActivities(clientID) []*Activity
  CreateInvoice(clientID, amount, desc, dueDate) *Invoice — nil if client not found
  ListInvoices(clientID) []*Invoice
  MarkInvoicePaid(id) *Invoice — nil if not found