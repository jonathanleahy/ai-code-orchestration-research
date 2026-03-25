# Revised MVP

# Revised MVP Plan Based on Senior Engineer Feedback

## Changes Made (what changed and why)

1. **Removed PDF generation** from MVP scope to reduce complexity and development time
2. **Added file-based persistence** for critical data to ensure data durability and prevent loss on restart
3. **Added basic authentication** to secure the single-user interface
4. **Simplified email delivery** to basic SMTP configuration instead of complex templating
5. **Added backup strategy** for persisted data
6. **Added GDPR compliance considerations** for handling personal data
7. **Maintained QuickBooks integration** as a key differentiator
8. **Kept core features** that were validated by all engineers

## Revised MVP Scope

### IN:
- Basic invoice creation with auto-calculation
- Client database with payment history tracking
- Automated overdue invoice follow-ups
- Simple template system
- Single-user interface with basic authentication
- QuickBooks integration (read-only for now)
- Email delivery functionality (basic SMTP)
- Basic client management (add/edit/delete)
- File-based persistence for invoices and clients
- Backup strategy for persisted data
- GDPR-compliant data handling

### OUT:
- PDF generation
- Team collaboration features
- Project-based invoicing with time tracking
- Payment processing integrations
- Advanced reporting
- Mobile app
- Complex permission systems
- Multi-user authentication
- Advanced tax compliance features
- Email template engine (complex)
- Slack/communication integrations

## Revised Architecture

### Core Architecture:
- **Single Go binary** (as originally planned)
- **File-based persistence** using JSON format for invoices and clients
- **In-memory cache** for active operations
- **Basic authentication** for user access
- **Configuration file** for settings (tax rates, default templates, SMTP settings)

### Data Storage Strategy:
- **Primary storage:** JSON files on disk for invoices and clients
- **Backup strategy:** Daily backups with rotation (30-day retention)
- **Error handling:** Robust file I/O with retry logic and error logging
- **Data validation:** Input sanitization and validation for all data

### Security Measures:
- **Authentication:** Basic username/password authentication
- **Data encryption:** At-rest encryption for sensitive data
- **GDPR compliance:** Data minimization, user consent handling, data deletion capabilities
- **Access controls:** Role-based access (single user for MVP)
- **Audit logging:** Basic logging of critical operations

### Deployment:
- **Hosting:** $5/month VPS (as originally planned)
- **Monitoring:** Basic health checks and error logging
- **Backup:** Automated daily backups with rotation
- **Deployment:** Simple binary deployment with configuration management

## Features Removed/Simplified

### Removed:
1. **PDF generation** - Too complex for MVP, deferred to v2
2. **Email template engine** - Simplified to basic SMTP configuration
3. **Complex permission systems** - Removed for single-user focus
4. **Multi-user authentication** - Simplified to single-user with basic auth
5. **Advanced reporting** - Deferred to v2
6. **Mobile app support** - Not prioritized for MVP

### Simplified:
1. **Email delivery** - Basic SMTP configuration instead of complex email system
2. **Data persistence** - File-based instead of database (with backup strategy)
3. **Authentication** - Basic auth instead of complex auth system
4. **Tax compliance** - Simplified to basic tax calculation with configurable rates

## Features Added (from dev recommendations)

### Added Features:
1. **File-based persistence** for invoices and clients with proper error handling
2. **Basic authentication** for user access security
3. **Backup strategy** with daily backups and rotation
4. **Configuration file** for settings management
5. **GDPR compliance framework** for data handling
6. **Data encryption** for sensitive information
7. **Audit logging** for critical operations
8. **Error handling** for file I/O operations
9. **Input validation** for all user data
10. **Health monitoring** for application stability

### Technical Implementation Details:
- **Storage format:** JSON files with structured data organization
- **Backup rotation:** 30-day retention with daily incremental backups
- **Configuration management:** YAML/JSON config file with environment overrides
- **Security:** Password hashing, session management, input sanitization
- **Logging:** Structured logging with error tracking and monitoring
- **Monitoring:** Basic uptime and error rate monitoring

This revised MVP maintains the core value proposition while addressing all security, operational, and technical concerns raised by the senior engineers, ensuring a solid foundation for future development while remaining focused on delivering immediate value to users.