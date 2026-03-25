# CRM System

A simple Customer Relationship Management (CRM) system built with Go that allows managing clients, activities, and invoices.

## What It Does

This CRM system provides a complete solution for managing customer relationships with the following features:

- **Client Management**: Create, read, update, and delete client records
- **Activity Tracking**: Log activities associated with clients
- **Invoice Management**: Create, track, and mark invoices as paid
- **RESTful API**: Complete HTTP API for all operations
- **Web Dashboard**: Simple HTML dashboard for viewing clients

## Quick Start (how to run)

1. **Prerequisites**: Go 1.19 or higher installed

2. **Clone and Run**:
```bash
git clone <repository-url>
cd <repository-name>
go run main.go
```

3. **Access the Application**:
   - Web Dashboard: http://localhost:8080
   - API Endpoints: http://localhost:8080/api/

## API Endpoints

### Client Management
| Method | Path | Description | Request | Response |
|--------|------|-------------|---------|----------|
| GET | `/api/clients` | List all clients | - | Array of client objects |
| POST | `/api/clients` | Create a new client | `{ "name": "string", "email": "string", "phone": "string" }` | Client object |
| GET | `/api/clients/{id}` | Get client details | - | Client object |
| POST | `/api/clients/{id}/update` | Update client information | `{ "name": "string", "email": "string", "phone": "string" }` | Client object |
| POST | `/api/clients/{id}/delete` | Delete a client | - | Success message |

### Activity Management
| Method | Path | Description | Request | Response |
|--------|------|-------------|---------|----------|
| POST | `/api/clients/{id}/activities` | Add activity for client | `{ "type": "string", "description": "string" }` | Activity object |
| GET | `/api/clients/{id}/activities` | List client activities | - | Array of activity objects |

### Invoice Management
| Method | Path | Description | Request | Response |
|--------|------|-------------|---------|----------|
| POST | `/api/clients/{id}/invoices` | Create invoice for client | `{ "number": "string", "amount": number, "description": "string", "due_date": "datetime" }` | Invoice object |
| GET | `/api/clients/{id}/invoices` | List client invoices | - | Array of invoice objects |
| POST | `/api/invoices/{id}/pay` | Mark invoice as paid | - | Invoice object |

### System Endpoints
| Method | Path | Description | Request | Response |
|--------|------|-------------|---------|----------|
| GET | `/health` | Health check | - | `{"status": "ok"}` |
| GET | `/invoice/{id}/print` | Print invoice | - | Invoice HTML |

## Environment Variables

No environment variables required. The application runs with default settings.

## Architecture

The application follows a layered architecture pattern:

1. **Data Layer**: `store` package handles all data operations with thread-safe maps
2. **Business Logic Layer**: Core logic for managing clients, activities, and invoices
3. **Presentation Layer**: HTTP handlers for API endpoints and web dashboard
4. **Persistence Layer**: In-memory storage with ID generation and timestamping

### Key Components:
- **Store Structure**: Thread-safe map-based storage with read-write mutex
- **ID Generation**: Timestamp-based unique IDs with sequential counters
- **Data Models**: Structured representations of clients, activities, and invoices
- **HTTP Router**: Standard library `http.ServeMux` for routing

## Tech Stack

- **Language**: Go 1.19+
- **Framework**: Standard library HTTP package
- **Storage**: In-memory maps with sync.RWMutex for thread safety
- **Data Format**: JSON for API communication
- **Frontend**: Simple HTML/CSS for dashboard
- **Build Tool**: Go modules for dependency management

The system is designed to be lightweight and efficient, suitable for small to medium-sized CRM needs with a focus on simplicity and performance.