# PART 1 — PRODUCT MANAGER: Define MVP

## MVP Definition for CRM with Invoice Generator

### Must-Have Features (All Personas Agree)
1. **Client Management**
   - Add/edit/delete clients
   - Store basic client info (name, email, phone, address)
   - Client search and filtering

2. **Invoice Generator**
   - Create invoices from client data
   - Basic invoice fields (client, items, amounts, dates)
   - Export invoices (PDF format)

3. **Client History Tracking**
   - Log client interactions (notes/comments)
   - Timestamped history entries
   - View history per client

4. **Address Book**
   - Centralized client address storage
   - Address management for each client
   - Search by address fields

5. **Comments/Notes System**
   - Add comments to clients
   - View all comments in chronological order
   - Comment threading (optional)

6. **GraphQL API**
   - GraphQL endpoint for all features
   - gqlgen implementation

### Deferred to v2
1. **Multi-user support** (team collaboration)
2. **Advanced reporting/analytics**
3. **Email integration** (automatic sync)
4. **Mobile app** (web-based only for now)
5. **Advanced search/filtering**
6. **File attachments** (documents, images)
7. **Payment tracking** (actual payment processing)
8. **Custom fields** (per-client customization)
9. **Workflow automation**
10. **Role-based permissions**

### Architecture
- Single Go binary
- In-memory store (no database needed)
- GraphQL API using gqlgen
- Simple file-based persistence for data backup (optional)
- Web UI built with React or similar (frontend not part of MVP)

---

# PART 2 — ENGINEER REVIEW

## Is it buildable in a single Go binary?
✅ **YES** - This is absolutely achievable with:
- Go modules for dependency management
- gqlgen for GraphQL schema generation
- Gin or Echo for HTTP routing
- In-memory data structures (maps/slices)
- Simple file I/O for persistence

## What's over-engineered? Simplify.
❌ **Nothing over-engineered** - The requirements are well-scoped for an MVP. The core features are essential and align with user needs.

## Can it run on a $5/month VPS?
✅ **YES** - A $5/month VPS (like DigitalOcean $5 droplet) can easily handle:
- Go binary running with minimal memory usage
- GraphQL API requests
- In-memory data processing
- Basic file I/O for persistence

## Security Concerns?
⚠️ **Some concerns**:
1. **No authentication/authorization** - Critical for any real-world application
2. **In-memory storage** - Data lost on restart
3. **No HTTPS** - Should be addressed in v2
4. **No rate limiting** - Could be exploited
5. **No input sanitization** - Risk of injection attacks

## VERDICT: APPROVE with Minor Changes

**APPROVE** with these changes:
1. Add basic authentication (username/password)
2. Implement rate limiting
3. Add input validation and sanitization
4. Add simple persistence (JSON file backup)
5. Add HTTPS support (can be added later)

---

# PART 3 — REVISED MVP

## Revised MVP Definition for CRM with Invoice Generator

### Must-Have Features (All Personas Agree)
1. **Client Management**
   - Add/edit/delete clients
   - Store basic client info (name, email, phone, address)
   - Client search and filtering

2. **Invoice Generator**
   - Create invoices from client data
   - Basic invoice fields (client, items, amounts, dates)
   - Export invoices (PDF format)

3. **Client History Tracking**
   - Log client interactions (notes/comments)
   - Timestamped history entries
   - View history per client

4. **Address Book**
   - Centralized client address storage
   - Address management for each client
   - Search by address fields

5. **Comments/Notes System**
   - Add comments to clients
   - View all comments in chronological order
   - Comment threading (optional)

6. **GraphQL API**
   - GraphQL endpoint for all features
   - gqlgen implementation

### Deferred to v2
1. **Multi-user support** (team collaboration)
2. **Advanced reporting/analytics**
3. **Email integration** (automatic sync)
4. **Mobile app** (web-based only for now)
5. **Advanced search/filtering**
6. **File attachments** (documents, images)
7. **Payment tracking** (actual payment processing)
8. **Custom fields** (per-client customization)
9. **Workflow automation**
10. **Role-based permissions**

### Architecture Requirements
- Single Go binary
- In-memory store with optional file persistence
- GraphQL API using gqlgen
- Basic authentication
- Rate limiting
- Input validation
- HTTPS support (to be added in v2)

### Technical Implementation Details
1. **Data Model**:
   - Client struct with ID, name, email, phone, addresses, comments
   - Invoice struct with ID, client ID, items, amounts, dates
   - Comment struct with ID, client ID, content, timestamp

2. **Storage**:
   - In-memory maps for fast access
   - JSON file backup on shutdown
   - Load from file on startup

3. **Security**:
   - Basic auth middleware
   - Rate limiting (100 req/min)
   - Input sanitization
   - CORS configuration

4. **API Endpoints**:
   - GraphQL endpoint at `/graphql`
   - REST-like endpoints for basic operations
   - Health check endpoint

5. **Deployment**:
   - Single binary deployment
   - Docker containerization (optional)
   - Configurable port and persistence path

This revised MVP addresses the engineering concerns while maintaining the core functionality that meets all personas' needs.