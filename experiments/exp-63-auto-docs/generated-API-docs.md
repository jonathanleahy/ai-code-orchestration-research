# API Documentation

## GET /health
**Description**: Health check endpoint to verify the server is running
**Request**: No body required
**Response**: 
- Status Code: 200 OK
- Body: `{"status": "ok"}`
**Errors**: None

## GET /
**Description**: Returns the CRM dashboard HTML page
**Request**: No body required
**Response**: 
- Status Code: 200 OK
- Body: HTML content of the dashboard
**Errors**: None

## GET /client/{id}
**Description**: Retrieve details of a specific client by ID
**Request**: No body required
**Response**: 
- Status Code: 200 OK
- Body: Client details in JSON format
- Example: `{"id": "1", "name": "John Doe", "email": "john@example.com", "phone": "555-0101"}`
**Errors**: 
- 404 Not Found: Client not found

## POST /api/clients
**Description**: Create a new client
**Request**: 
- Body: JSON with client details
- Required fields: name, email, phone
- Example: `{"name": "Alice Johnson", "email": "alice@example.com", "phone": "555-0103"}`
**Response**: 
- Status Code: 201 Created
- Body: Created client details in JSON format
- Example: `{"id": "3", "name": "Alice Johnson", "email": "alice@example.com", "phone": "555-0103"}`
**Errors**: 
- 400 Bad Request: Invalid request body

## GET /api/clients
**Description**: List all clients
**Request**: No body required
**Response**: 
- Status Code: 200 OK
- Body: Array of client objects in JSON format
- Example: `[{"id": "1", "name": "John Doe", "email": "john@example.com", "phone": "555-0101"}, {"id": "2", "name": "Jane Smith", "email": "jane@example.com", "phone": "555-0102"}]`
**Errors**: None

## POST /api/clients/{id}/update
**Description**: Update an existing client's information
**Request**: 
- Body: JSON with updated client fields
- Required fields: name, email, phone
- Example: `{"name": "John Smith", "email": "john.smith@example.com", "phone": "555-0104"}`
**Response**: 
- Status Code: 200 OK
- Body: Updated client details in JSON format
- Example: `{"id": "1", "name": "John Smith", "email": "john.smith@example.com", "phone": "555-0104"}`
**Errors**: 
- 400 Bad Request: Invalid request body
- 404 Not Found: Client not found

## POST /api/clients/{id}/delete
**Description**: Delete a client by ID
**Request**: No body required
**Response**: 
- Status Code: 200 OK
- Body: Success message in JSON format
- Example: `{"message": "Client deleted successfully"}`
**Errors**: 
- 404 Not Found: Client not found

## POST /api/clients/{id}/activities
**Description**: Create a new activity for a client
**Request**: 
- Body: JSON with activity details
- Required fields: description, date
- Example: `{"description": "Follow-up call", "date": "2023-05-15"}`
**Response**: 
- Status Code: 201 Created
- Body: Created activity details in JSON format
- Example: `{"id": "1", "clientId": "1", "description": "Follow-up call", "date": "2023-05-15"}`
**Errors**: 
- 400 Bad Request: Invalid request body
- 404 Not Found: Client not found

## GET /api/clients/{id}/activities
**Description**: List all activities for a specific client
**Request**: No body required
**Response**: 
- Status Code: 200 OK
- Body: Array of activity objects in JSON format
- Example: `[{"id": "1", "clientId": "1", "description": "Follow-up call", "date": "2023-05-15"}]`
**Errors**: 
- 404 Not Found: Client not found

## POST /api/clients/{id}/invoices
**Description**: Create a new invoice for a client
**Request**: 
- Body: JSON with invoice details
- Required fields: amount, dueDate
- Example: `{"amount": 1000, "dueDate": "2023-06-15"}`
**Response**: 
- Status Code: 201 Created
- Body: Created invoice details in JSON format
- Example: `{"id": "1", "clientId": "1", "amount": 1000, "dueDate": "2023-06-15", "paid": false}`
**Errors**: 
- 400 Bad Request: Invalid request body
- 404 Not Found: Client not found

## GET /api/clients/{id}/invoices
**Description**: List all invoices for a specific client
**Request**: No body required
**Response**: 
- Status Code: 200 OK
- Body: Array of invoice objects in JSON format
- Example: `[{"id": "1", "clientId": "1", "amount": 1000, "dueDate": "2023-06-15", "paid": false}]`
**Errors**: 
- 404 Not Found: Client not found

## POST /api/invoices/{id}/pay
**Description**: Mark an invoice as paid
**Request**: No body required
**Response**: 
- Status Code: 200 OK
- Body: Success message in JSON format
- Example: `{"message": "Invoice marked as paid"}`
**Errors**: 
- 404 Not Found: Invoice not found

## GET /invoice/{id}/print
**Description**: Generate a printable version of an invoice
**Request**: No body required
**Response**: 
- Status Code: 200 OK
- Body: HTML content for printing the invoice
**Errors**: 
- 404 Not Found: Invoice not found