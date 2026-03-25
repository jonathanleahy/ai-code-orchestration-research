# MVP Plan

# Invoice Generator MVP Plan

## Must-Have Features (all personas agree)
1. **Auto-calculation of totals and taxes** - Addresses Sarah's error-prone manual data entry and Marcus's inconsistent invoicing
2. **Client database with payment history** - Solves Sarah's tracking issues and Marcus's multiple client management
3. **Automated follow-ups for overdue invoices** - Directly addresses Sarah's payment follow-up frustration and Marcus's cash flow visibility
4. **Basic invoice creation with template support** - Covers all three personas' need for streamlined invoicing

## Should-Have (2/3 personas)
1. **Team collaboration with role-based access** - Marcus's agency needs this, Sarah and Jennifer could benefit
2. **Project-based invoicing with time tracking integration** - Marcus's primary need, Jennifer's compliance needs
3. **Integration with QuickBooks** - Critical for Marcus and Jennifer's existing workflows
4. **Google Drive/Cloud storage integration** - Sarah's preference, Jennifer's backup needs

## Deferred to v2
1. **Stripe payment processing integration** - Jennifer's compliance focus, Marcus's payment needs
2. **Slack/communication integrations** - Marcus's team needs
3. **Advanced reporting and analytics** - All personas could use this later
4. **Mobile app support** - Not prioritized for MVP

## MVP Scope: What's IN and what's OUT

### IN:
- Basic invoice creation with auto-calculation
- Client database with payment history tracking
- Automated overdue invoice follow-ups
- Simple template system
- Single-user interface (no team features)
- QuickBooks integration (read-only for now)
- Email delivery functionality
- Basic client management (add/edit/delete)

### OUT:
- Team collaboration features
- Project-based invoicing with time tracking
- Payment processing integrations
- Advanced reporting
- Mobile app
- Complex permission systems
- Multi-user authentication
- Advanced tax compliance features

## Architecture: Single Go binary, in-memory store, no database
- **Single binary deployment** for easy distribution and maintenance
- **In-memory store** (using Go's map structures) for temporary data persistence
- **No database** for MVP simplicity - data stored in memory during runtime
- **File-based persistence** for critical data (invoices, clients) that can be saved to disk
- **Lightweight and fast** startup times for immediate user experience

## Target price: $15/month
- Aligns with Sarah's $20/month willingness to pay
- Competitive with existing solutions like FreshBooks ($15-25/month)
- Covers basic functionality while remaining affordable for freelancers
- Provides value through time savings (2-3 hours/month for Sarah, 10+ hours/month for Marcus)
- Allows for premium features in future tiers

## MVP Features Breakdown
1. **Core Invoice Generation** - Template-based invoice creation with auto-calculation
2. **Client Management** - Store client information with payment history
3. **Follow-up System** - Automated email reminders for overdue invoices
4. **QuickBooks Sync** - Basic read-only integration for accounting data
5. **Email Delivery** - Direct email sending of invoices
6. **Simple UI** - Clean, intuitive interface for all user types

This MVP addresses the core pain points of all three personas while maintaining simplicity and affordability for early adoption.