### ADR-001: Use Go Standard Library Only

**Status**: Accepted

**Context**: The CRM needs to be lightweight and avoid external dependencies that could introduce security vulnerabilities, compatibility issues, or increase deployment complexity. The project scope is limited to basic CRM functionality without complex features that would require specialized libraries.

**Decision**: Use only Go's standard library for all functionality, including HTTP routing, JSON handling, and concurrency management.

**Consequences**: 
- Good: Reduced dependency management overhead, smaller binary size, no external security vulnerabilities, easier deployment and maintenance
- Bad: Missing out on battle-tested features and optimizations that popular frameworks provide, more boilerplate code, limited built-in features for complex operations

**Alternatives Considered**: 
- Using Gin or Echo frameworks: Rejected because they would add unnecessary complexity and dependencies for a simple CRM with minimal features
- Using Fiber: Rejected for the same reasons as above, plus additional dependency maintenance overhead

### ADR-002: In-Memory Store Instead of Database

**Status**: Accepted

**Context**: The CRM is designed for small-scale use cases where data persistence isn't critical, and the focus is on rapid development and simple deployment. The application is intended to be lightweight and run in single-instance scenarios.

**Decision**: Use in-memory storage with Go's built-in data structures (maps, slices) to store CRM data, with no persistent storage layer.

**Consequences**: 
- Good: Faster development, simpler architecture, no database setup or maintenance, better performance for read-heavy operations
- Bad: Data loss on application restart, no persistence across deployments, limited scalability for large datasets, not suitable for production environments requiring data durability

**Alternatives Considered**: 
- Using SQLite: Rejected because it would add complexity and dependencies while providing minimal benefit for a simple CRM
- Using PostgreSQL/MySQL: Rejected because it would significantly increase deployment complexity and infrastructure requirements

### ADR-003: Embedded HTML Instead of Template Engine

**Status**: Accepted

**Context**: The CRM has simple HTML output requirements with minimal dynamic content. Using a template engine would add unnecessary complexity and dependencies for a small application with straightforward UI needs.

**Decision**: Embed HTML directly in Go code using string literals and string formatting for dynamic content, avoiding template engines entirely.

**Consequences**: 
- Good: No additional dependencies, simpler build process, faster startup time, easier debugging, more control over HTML generation
- Bad: HTML mixing with Go code makes maintenance harder, more error-prone for complex HTML structures, harder to separate concerns, less maintainable for large UI components

**Alternatives Considered**: 
- Using html/template package: Rejected because it adds complexity and dependencies for simple HTML generation
- Using Pongo2 or similar template engines: Rejected due to dependency overhead and complexity for basic HTML rendering needs

### ADR-004: Use sync.RWMutex for Concurrency Control

**Status**: Accepted

**Context**: The CRM needs to handle concurrent access to shared data structures safely while maintaining good performance. The application is expected to handle multiple simultaneous requests for CRUD operations.

**Decision**: Use Go's sync.RWMutex for read-heavy operations where multiple readers can access data concurrently while ensuring exclusive access for writes.

**Consequences**: 
- Good: Better performance for read-heavy workloads, proper synchronization for concurrent access, standard library solution with no additional dependencies
- Bad: More complex code for managing read/write locks, potential for deadlocks if not used carefully, requires careful consideration of lock granularity

**Alternatives Considered**: 
- Using sync.Mutex: Rejected because it would serialize all operations, reducing performance for read-heavy scenarios
- Using channels for coordination: Rejected because it would add unnecessary complexity for simple data access patterns

### ADR-005: Sequential IDs Instead of UUIDs

**Status**: Accepted

**Context**: The CRM application is designed for internal use with a simple ID generation requirement. Sequential IDs are sufficient for the application's needs and provide better readability and simplicity.

**Decision**: Generate sequential integer IDs for CRM entities instead of using UUIDs or other random ID generators.

**Consequences**: 
- Good: Simpler ID generation logic, more readable IDs for debugging, smaller ID size, easier to understand and maintain
- Bad: Not suitable for distributed systems, potential collision issues in high-concurrency scenarios, less secure than UUIDs, not globally unique

**Alternatives Considered**: 
- Using UUIDs: Rejected because they add complexity to ID generation and storage, increase ID size, and are overkill for a simple internal CRM
- Using snowflake IDs: Rejected due to increased complexity and dependency on external libraries for ID generation

### ADR-006: Single Binary Deployment

**Status**: Accepted

**Context**: The CRM application is designed to be deployed in simple environments where minimizing deployment complexity is crucial. The application should be easily distributable and runnable without complex setup procedures.

**Decision**: Package the entire application as a single binary that can be executed directly without requiring additional dependencies or complex deployment procedures.

**Consequences**: 
- Good: Simple deployment process, no dependency management issues, easy distribution, single file for backup and version control
- Bad: Larger binary size due to including standard library, harder to update individual components, no ability to share libraries between applications, limited flexibility for complex deployment scenarios

**Alternatives Considered**: 
- Using Docker containers: Rejected because it adds complexity to deployment and requires Docker installation
- Using multiple binaries or shared libraries: Rejected because it increases deployment complexity and maintenance overhead