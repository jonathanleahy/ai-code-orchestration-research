# Text Wireframes

# Screen Wireframes

## S01 - Dashboard

```
┌─────────────────────────────────────────────────────────────────────┐
│ Dashboard                          [Add Check] [Notifications] [User] │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐               │
│  │  Service A   │  │  Service B   │  │  Service C   │               │
│  │   [Green]    │  │   [Red]      │  │   [Yellow]   │               │
│  │   (2 min)    │  │   (15 min)   │  │   (5 min)    │               │
│  │   http://... │  │   http://... │  │   http://... │               │
│  └──────────────┘  └──────────────┘  └──────────────┘               │
│                                                                     │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐               │
│  │  Service D   │  │  Service E   │  │  Service F   │               │
│  │   [Green]    │  │   [Green]    │  │   [Red]      │               │
│  │   (1 min)    │  │   (8 min)    │  │   (20 min)   │               │
│  │   http://... │  │   http://... │  │   http://... │               │
│  └──────────────┘  └──────────────┘  └──────────────┘               │
│                                                                     │
│                                                                     │
│  [No checks configured yet]                                         │
│                                                                     │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
```

**Responsive Notes:**
- On mobile: Cards stack vertically
- Sidebar collapses to hamburger menu
- Buttons become full-width

---

## S02 - Add Check

```
┌─────────────────────────────────────────────────────────────────────┐
│ Add New Check                                                       │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│  ┌───────────────────────────────────────────────────────────────┐  │
│  │  Name: [________________________________________________]   │  │
│  │  URL:  [________________________________________________]   │  │
│  │                                                                 │  │
│  │  Check Type: [Dropdown ▼]                                      │  │
│  │    ▼ HTTP Status Code                                          │  │
│  │    ▼ DNS Resolution                                            │  │
│  │    ▼ SSL Certificate                                           │  │
│  │    ▼ Ping                                                      │  │
│  │                                                                 │  │
│  │  [Create Check] [Cancel]                                        │  │
│  └───────────────────────────────────────────────────────────────┘  │
│                                                                     │
│  [Error: Please enter a valid URL]                                  │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
```

**Responsive Notes:**
- Form fields stack vertically on mobile
- Buttons stack on small screens
- Dropdown expands to full width

---

## S03 - Check Detail

```
┌─────────────────────────────────────────────────────────────────────┐
│ [Back] Check Detail                                                 │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│  ┌───────────────────────────────────────────────────────────────┐  │
│  │  Service Name: API Gateway                                    │  │
│  │  Status: [Green]                                              │  │
│  │  URL: https://api.example.com/v1/status                       │  │
│  │  Last Checked: 2 minutes ago                                  │  │
│  │  [Create Incident]                                            │  │
│  └───────────────────────────────────────────────────────────────┘  │
│                                                                     │
│  ┌───────────────────────────────────────────────────────────────┐  │
│  │  Timeline Chart (Line Graph)                                  │  │
│  │  [100%] [95%] [90%] [85%] [80%] [75%] [70%] [65%] [60%]      │  │
│  │  [00:00] [01:00] [02:00] [03:00] [04:00] [05:00] [06:00]     │  │
│  └───────────────────────────────────────────────────────────────┘  │
│                                                                     │
│  ┌───────────────────────────────────────────────────────────────┐  │
│  │  Recent Incidents                                             │  │
│  │  ┌─────────────────────────────────────────────────────────┐  │  │
│  │  │  Incident #12345                                           │  │  │
│  │  │  Status: Resolved                                         │  │  │
│  │  │  Time: 2 hours ago                                        │  │  │
│  │  │  Severity: Medium                                         │  │  │
│  │  │  [View Details]                                           │  │  │
│  │  └─────────────────────────────────────────────────────────┘  │  │
│  │  ┌─────────────────────────────────────────────────────────┐  │  │
│  │  │  Incident #12346                                           │  │  │
│  │  │  Status: Open                                             │  │  │
│  │  │  Time: 1 hour ago                                         │  │  │
│  │  │  Severity: High                                           │  │  │
│  │  │  [View Details]                                           │  │  │
│  │  └─────────────────────────────────────────────────────────┘  │  │
│  └───────────────────────────────────────────────────────────────┘  │
│                                                                     │
│  ┌───────────────────────────────────────────────────────────────┐  │
│  │  Activity Log                                                 │  │
│  │  [10:30 AM] Service check completed successfully              │  │
│  │  [10:25 AM] Service check failed (500 error)                  │  │
│  │  [10:20 AM] Incident #12345 created                           │  │
│  └───────────────────────────────────────────────────────────────┘  │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
```

**Responsive Notes:**
- Timeline chart becomes stacked on mobile
- Incident cards stack vertically
- Activity log scrolls if too long

---

## S04 - Incidents List

```
┌─────────────────────────────────────────────────────────────────────┐
│ Incidents                          [Create Incident] [Filter]      │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│  ┌───────────────────────────────────────────────────────────────┐  │
│  │  Incident #12345                                              │  │
│  │  Status: [Resolved]                                           │  │
│  │  Time: 2 hours ago                                            │  │
│  │  Severity: [Medium]                                           │  │
│  │  Title: API Gateway Timeout                                   │  │
│  │  [View Details]                                               │  │
│  └───────────────────────────────────────────────────────────────┘  │
│                                                                     │
│  ┌───────────────────────────────────────────────────────────────┐  │
│  │  Incident #12346                                              │  │
│  │  Status: [Open]                                               │  │
│  │  Time: 1 hour ago                                             │  │
│  │  Severity: [High]                                             │  │
│  │  Title: Database Connection Failure                           │  │
│  │  [View Details]                                               │  │
│  └───────────────────────────────────────────────────────────────┘  │
│                                                                     │
│  ┌───────────────────────────────────────────────────────────────┐  │
│  │  Incident #12347                                              │  │
│  │  Status: [Investigating]                                      │  │
│  │  Time: 30 minutes ago                                         │  │
│  │  Severity: [Critical]                                         │  │
│  │  Title: Service Outage                                        │  │
│  │  [View Details]                                               │  │
│  └───────────────────────────────────────────────────────────────┘  │
│                                                                     │
│  [No incidents found]                                               │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
```

**Responsive Notes:**
- Incident cards stack vertically on mobile
- Filter dropdown becomes full-width
- "Create Incident" button becomes floating action button

---

## S05 - New Incident

```
┌─────────────────────────────────────────────────────────────────────┐
│ Create New Incident                                                 │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│  ┌───────────────────────────────────────────────────────────────┐  │
│  │  Title: [________________________________________________]   │  │
│  │  Description: [________________________________________]   │  │
│  │  [________________________________________________________]   │  │
│  │  [________________________________________________________]   │  │
│  │                                                                 │  │
│  │  Severity Level: [Dropdown ▼]                                 │  │
│  │    ▼ Low                                                      │  │
│  │    ▼ Medium                                                   │  │
│  │    ▼ High                                                     │  │
│  │    ▼ Critical                                                 │  │
│  │                                                                 │  │
│  │  [Create Incident] [Cancel]                                     │  │
│  └───────────────────────────────────────────────────────────────┘  │
│                                                                     │
│  [Error: Please enter a title]                                      │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
```

**Responsive Notes:**
- Form fields stack vertically on mobile
- Buttons stack on small screens
- Description textarea expands to full width

---

## S06 - Incident Detail

```
┌─────────────────────────────────────────────────────────────────────┐
│ [Back] Incident Details                                             │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│  ┌───────────────────────────────────────────────────────────────┐  │
│  │  Incident #12345                                              │  │
│  │  Status: [Resolved]                                           │  │
│  │  Severity: [High]                                             │  │
│  │  Created: 2 hours ago                                         │  │
│  │  Last Updated: 1 hour ago                                     │  │
│  └───────────────────────────────────────────────────────────────┘  │
│                                                                     │
│  ┌───────────────────────────────────────────────────────────────┐  │
│  │  Description:                                                 │  │
│  │  The API Gateway experienced timeout errors during peak load  │  │
│  │  periods. This was caused by insufficient resource allocation. │  │
│  └───────────────────────────────────────────────────────────────┘  │
│                                                                     │
│  ┌───────────────────────────────────────────────────────────────┐  │
│  │  Timeline of Events                                           │  │
│  │  [10:30 AM] Incident created                                  │  │
│  │  [10:35 AM] Initial investigation started                     │  │
│  │  [10:45 AM] Resource scaling initiated                        │  │
│  │  [11:00 AM] Service restored                                  │  │
│  │  [11:15 AM] Resolution confirmed                              │  │
│  └───────────────────────────────────────────────────────────────┘  │
│                                                                     │
│  ┌───────────────────────────────────────────────────────────────┐  │
│  │  Add Update                                                   │  │
│  │  [________________________________________________________]   │  │
│  │  [Update Status: [Dropdown ▼]]                               │  │
│  │    ▼ Open                                                     │  │
│  │    ▼ Investigating                                            │  │
│  │    ▼ Resolved                                                 │  │
│  │                                                                 │  │
│  │  [Add Update]                                                   │  │
│  └───────────────────────────────────────────────────────────────┘  │
│                                                                     │
│  ┌───────────────────────────────────────────────────────────────┐  │
│  │  Resolution Notes                                             │  │
│  │  The issue was resolved by scaling up the API Gateway nodes.  │  │
│  │  Additional monitoring has been implemented to prevent this.  │  │
│  └───────────────────────────────────────────────────────────────┘  │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
```

**Responsive Notes:**
- Timeline scrolls horizontally on mobile
- Update form fields stack vertically
- Status dropdown becomes full-width

---

## S07 - Public Status Page

```
┌─────────────────────────────────────────────────────────────────────┐
│ Public Status Page                 [Subscribe] [Help] [About]       │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐               │
│  │  Service A   │  │  Service B   │  │  Service C   │               │
│  │   [Green]    │  │   [Red]      │  │   [Yellow]   │               │
│  │   (2 min)    │  │   (15 min)   │  │   (5 min)    │               │
│  │   API Gateway│  │   Database   │  │   Frontend   │               │
│  │   [Subscribe]│  │   [Subscribe]│  │   [Subscribe]│               │
│  └──────────────┘  └──────────────┘  └──────────────┘               │
│                                                                     │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐               │
│  │  Service D   │  │  Service E   │  │  Service F   │               │
│  │   [Green]    │  │   [Green]    │  │   [Red]      │               │
│  │   (1 min)    │  │   (8 min)    │  │   (20 min)   │               │
│  │   CDN        │  │   Auth       │  │   Payment    │               │
│  │   [Subscribe]│  │   [Subscribe]│  │   [Subscribe]│               │
│  └──────────────┘  └──────────────┘  └──────────────┘               │
│                                                                     │
│                                                                     │
│  [No services configured yet]                                       │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
```

**Responsive Notes:**
- Cards stack vertically on mobile
- Subscribe buttons become full-width
- Navigation bar collapses to hamburger menu

---

## S08 - Service Detail

```
┌─────────────────────────────────────────────────────────────────────┐
│ [Back] Service Detail                                               │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│  ┌───────────────────────────────────────────────────────────────┐  │
│  │  Service Name: API Gateway                                    │  │
│  │  Status: [Green]                                              │  │
│  │  Last Checked: 2 minutes ago                                  │  │
│  │  URL: https://api.example.com/v1/status                       │  │
│  └───────────────────────────────────────────────────────────────┘  │
│                                                                     │
│  ┌───────────────────────────────────────────────────────────────┐  │
│  │  Recent Incidents                                             │  │
│  │  ┌─────────────────────────────────────────────────────────┐  │  │
│  │  │  Incident #12345                                           │  │  │
│  │  │  Status: Resolved                                         │  │  │
│  │  │  Time: 2 hours ago                                        │  │  │
│  │  │  Severity: Medium                                         │  │  │
│  │  │  [View Details]                                           │  │  │
│  │  └─────────────────────────────────────────────────────────┘  │  │
│  │  ┌─────────────────────────────────────────────────────────┐  │  │
│  │  │  Incident #12346                                           │  │  │
│  │  │  Status: Open                                             │  │  │
│  │  │  Time: 1 hour ago                                         │  │  │
│  │  │  Severity: High                                           │  │  │
│  │  │  [View Details]                                           │  │  │
│  │  └─────────────────────────────────────────────────────────┘  │  │
│  └───────────────────────────────────────────────────────────────┘  │
│                                                                     │
│  ┌───────────────────────────────────────────────────────────────┐  │
│  │  Subscribe to Updates                                         │  │
│  │  Email: [____________________________________________]   │  │
│  │  [Subscribe]                                                  │  │
│  └───────────────────────────────────────────────────────────────┘  │
│                                                                     │
│  [Subscribe to Updates]                                             │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
```

**Responsive Notes:**
- Incident cards stack vertically on mobile
- Subscribe form fields stack vertically
- Back button becomes hamburger menu on mobile

---

## S09 - Subscription Confirmation

```
┌─────────────────────────────────────────────────────────────────────┐
│ Subscription Confirmation                                           │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│  ┌───────────────────────────────────────────────────────────────┐  │
│  │  [Success]                                                     │  │
│  │  You have successfully subscribed to updates for API Gateway.   │  │
│  │                                                                 │  │
│  │  [Subscribed]                                                   │  │
│  │                                                                 │  │
│  │  [Back to Service Detail] [Back to Status Page]               │  │
│  └───────────────────────────────────────────────────────────────┘  │
│                                                                     │
│                                                                     │
│                                                                     │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
```

**Responsive Notes:**
- Buttons stack vertically on mobile
- Success message becomes full-width
- Back buttons become full-width on small screens