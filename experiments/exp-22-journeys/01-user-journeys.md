# User Journeys

**Brief:** Build a public status page for monitoring website uptime, managing incidents, and notifying subscribers.

### Personas

**Name**: Site Reliability Engineer
**Goal**: Monitor system health, quickly identify and resolve outages, maintain service level agreements
**Frequency**: Multiple times per day, continuous monitoring

**Name**: Product Manager
**Goal**: Stay informed about service status for business planning, communicate with stakeholders, track incident impact
**Frequency**: Several times per day during business hours, daily review of incident history

**Name**: Customer Support Agent
**Goal**: Provide accurate status information to customers, escalate issues when needed, reduce incoming support tickets
**Frequency**: Multiple times per day during shift, peak hours with high customer volume

**Name**: External Subscriber
**Goal**: Stay informed about service availability, receive timely notifications about incidents, understand impact on their usage
**Frequency**: As needed, typically when experiencing service issues or checking status before using services

### User Journeys

**Journey: Initial Setup and Configuration**
1. Admin clicks "Add Check" button in dashboard
2. System displays modal with URL field, name field, and check type dropdown
3. Admin enters "https://api.company.com" in URL field, "API Gateway" in name field, selects "HTTP" from dropdown
4. Admin clicks "Create Check" button
5. System shows green confirmation toast "Check created successfully" and refreshes the main checks list
6. New check appears in list with "unknown" status indicator and "Last checked: Just now" timestamp
7. Success: Check appears in list with proper name and URL, status shows as "unknown"
8. Edge case: If URL is invalid, system shows red error message below URL field and "Create Check" button remains disabled

**Journey: Daily Monitoring Routine**
1. Site Reliability Engineer opens status page in browser
2. System displays main dashboard showing all configured checks with status indicators (green for up, yellow for warning, red for down)
3. Engineer clicks on "API Gateway" check card to view details
4. System shows detailed view with timeline chart, recent incidents, and current status
5. Engineer scrolls through recent activity log to verify no recent issues
6. Success: Engineer confirms all services are operating normally with no active incidents
7. Edge case: If engineer has no configured checks, system shows "No checks configured yet" with "Add Check" CTA

**Journey: Incident Response Workflow**
1. Site Reliability Engineer notices red status indicator on "API Gateway" check
2. Engineer clicks on the check to view detailed incident information
3. System shows incident timeline with timestamps and status changes
4. Engineer clicks "Create Incident" button
5. System displays incident creation modal with fields for title, description, and severity level
6. Engineer enters "API Gateway Outage - High Priority" in title, "All API endpoints returning 500 errors" in description, selects "High" severity
7. Engineer clicks "Create Incident" button
8. System shows green confirmation toast "Incident created successfully" and redirects to incident detail page
9. Engineer adds update to incident with "Investigation in progress, initial root cause identified as database connection pool exhaustion"
10. Success: Incident is created with proper title and description, update is added to timeline
11. Edge case: If engineer tries to create incident without selecting severity, system shows validation error and prevents submission

**Journey: Subscriber Notification Experience**
1. External Subscriber visits public status page
2. System displays main dashboard showing all services with current status indicators
3. Subscriber clicks on "API Gateway" service card
4. System shows service detail page with current status, recent incidents, and notification preferences
5. Subscriber clicks "Subscribe to Updates" button
6. System displays email subscription modal with field for email address
7. Subscriber enters "customer@company.com" in email field and clicks "Subscribe"
8. System shows green confirmation toast "Successfully subscribed to updates for API Gateway"
9. Success: Subscriber receives confirmation email and sees "Subscribed" badge on service card
10. Edge case: If subscriber enters invalid email format, system shows red error message and prevents subscription

**Journey: Incident Resolution and Communication**
1. Site Reliability Engineer resolves the API Gateway issue
2. Engineer navigates to incident detail page
3. Engineer clicks "Update Status" button
4. System displays status update modal with options: "Investigation Complete", "Resolved", "Monitoring"
5. Engineer selects "Resolved" and adds note "Database connection pool issue resolved, all endpoints returning 200 OK"
6. Engineer clicks "Update Incident" button
7. System shows green confirmation toast "Status updated successfully" and refreshes incident timeline
8. System automatically sends notification to all subscribers with resolution details
9. Success: Incident status changes to "Resolved", subscribers receive email notification with resolution details
10. Edge case: If engineer tries to update status without adding resolution notes, system shows validation error and prevents update