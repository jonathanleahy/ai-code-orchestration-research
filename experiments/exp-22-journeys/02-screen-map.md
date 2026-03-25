# Screen Map

### Screen Map

---

#### **S01 - Dashboard**
- **URL**: `/`
- **Purpose**: Display all configured checks and their current status for quick overview.
- **Reached from**: 
  - S06 (Incident Detail)
  - S03 (Check Detail)
  - S05 (New Incident)
  - S02 (Add Check)
- **Components**:
  - Header with navigation
  - Service cards (each showing status indicator, name, last checked time)
  - "Add Check" button
  - Empty state message ("No checks configured yet") if no checks exist
- **Data shown**:
  - List of all configured checks
  - Status indicators (green/yellow/red)
  - Last checked timestamp
- **Actions available**:
  - Click on any check card to view details
  - Click "Add Check" to create a new check
- **Navigation**:
  - To S02 (Add Check)
  - To S03 (Check Detail)
  - To S04 (Incidents)

---

#### **S02 - Add Check**
- **URL**: `/admin/checks/new`
- **Purpose**: Allow admin users to configure a new monitoring check.
- **Reached from**: 
  - S01 (Dashboard)
- **Components**:
  - Modal dialog
  - Input fields: URL, Name, Check Type
  - Dropdown for Check Type
  - Submit button ("Create Check")
  - Cancel button
  - Validation error messages
- **Data shown**:
  - Form fields for URL, name, and check type
- **Actions available**:
  - Fill in form fields
  - Submit form
- **Navigation**:
  - Back to S01 (Dashboard) on success or cancel
  - Error state shows validation messages

---

#### **S03 - Check Detail**
- **URL**: `/checks/:id`
- **Purpose**: Show detailed information about a specific check including timeline, incidents, and activity log.
- **Reached from**:
  - S01 (Dashboard)
- **Components**:
  - Header with back button
  - Status badge
  - Timeline chart
  - Recent incidents list
  - Activity log
  - “Create Incident” button
- **Data shown**:
  - Check name, URL, status
  - Timeline chart of historical data
  - List of recent incidents
  - Activity log
- **Actions available**:
  - Click "Create Incident"
  - Go back to dashboard
- **Navigation**:
  - To S04 (Incidents)
  - Back to S01 (Dashboard)

---

#### **S04 - Incidents List**
- **URL**: `/incidents`
- **Purpose**: Display a list of all incidents, sorted by most recent.
- **Reached from**:
  - S01 (Dashboard)
  - S06 (Incident Detail)
- **Components**:
  - Header with navigation
  - Incident cards (title, status, time, severity)
  - "Create Incident" button
  - Empty state message ("No incidents found") if none exist
- **Data shown**:
  - List of incidents with key details
- **Actions available**:
  - Click on an incident to view details
  - Click "Create Incident"
- **Navigation**:
  - To S05 (New Incident)
  - To S06 (Incident Detail)
  - Back to S01 (Dashboard)

---

#### **S05 - New Incident**
- **URL**: `/incidents/new`
- **Purpose**: Allow SREs to create a new incident with title, description, and severity.
- **Reached from**:
  - S01 (Dashboard)
  - S03 (Check Detail)
  - S04 (Incidents List)
- **Components**:
  - Modal dialog
  - Input fields: Title, Description, Severity Level
  - Submit button ("Create Incident")
  - Cancel button
  - Validation error messages
- **Data shown**:
  - Form fields for incident details
- **Actions available**:
  - Fill in form fields
  - Submit form
- **Navigation**:
  - Back to S04 (Incidents List) on success or cancel
  - Redirects to S06 (Incident Detail) on successful creation

---

#### **S06 - Incident Detail**
- **URL**: `/incidents/:id`
- **Purpose**: Show full details of an incident including timeline, updates, and resolution status.
- **Reached from**:
  - S04 (Incidents List)
  - S05 (New Incident)
- **Components**:
  - Header with back button
  - Incident title and status
  - Timeline of events
  - Update form
  - Status update dropdown
  - Resolution notes
- **Data shown**:
  - Incident title, description, severity, status
  - Timeline of updates and status changes
- **Actions available**:
  - Add update to timeline
  - Update incident status
  - Go back to dashboard or incidents list
- **Navigation**:
  - To S01 (Dashboard)
  - To S04 (Incidents List)
  - Back to S05 (New Incident) if needed

---

#### **S07 - Public Status Page**
- **URL**: `/status`
- **Purpose**: Allow external subscribers to view service status and subscribe to updates.
- **Reached from**: 
  - Direct access via URL
- **Components**:
  - Header with navigation
  - Service cards (status, name, last checked)
  - "Subscribe to Updates" button per service
- **Data shown**:
  - List of services with current status indicators
- **Actions available**:
  - Click on service card to view details
  - Click "Subscribe to Updates"
- **Navigation**:
  - To S08 (Service Detail)
  - Back to S07 (Public Status Page)

---

#### **S08 - Service Detail**
- **URL**: `/status/:serviceId`
- **Purpose**: Show detailed status of a specific service including recent incidents and subscription options.
- **Reached from**:
  - S07 (Public Status Page)
- **Components**:
  - Header with back button
  - Service name and current status
  - Recent incidents list
  - Subscription form
  - "Subscribe to Updates" button
- **Data shown**:
  - Service name, current status
  - List of recent incidents
  - Subscription form
- **Actions available**:
  - Subscribe to updates
  - Go back to public status page
- **Navigation**:
  - To S09 (Subscription Confirmation)
  - Back to S07 (Public Status Page)

---

#### **S09 - Subscription Confirmation**
- **URL**: `/status/:serviceId/subscribe`
- **Purpose**: Confirm successful subscription to service updates.
- **Reached from**:
  - S08 (Service Detail)
- **Components**:
  - Success message
  - "Subscribed" badge
  - Back to service detail or home
- **Data shown**:
  - Confirmation of subscription
- **Actions available**:
  - Return to service detail
  - Return to public status page
- **Navigation**:
  - Back to S08 (Service Detail)
  - Back to S07 (Public Status Page)

---

### Screen Flow Diagram

```
S01 (Dashboard) → S02 (Add Check)
                → S03 (Check Detail)
                → S04 (Incidents)
S04 (Incidents) → S05 (New Incident)
                → S06 (Incident Detail)
S07 (Public Status) → S08 (Service Detail)
                    → S09 (Subscription Confirmation)
S06 (Incident Detail) → S01 (Dashboard)
S05 (New Incident) → S06 (Incident Detail)
S03 (Check Detail) → S01 (Dashboard)
S08 (Service Detail) → S07 (Public Status)
S09 (Subscription Confirmation) → S08 (Service Detail)
S02 (Add Check) → S01 (Dashboard)
```

This structure ensures that all user journeys are supported, with clear navigation paths and appropriate error handling for edge cases like invalid inputs or empty states.