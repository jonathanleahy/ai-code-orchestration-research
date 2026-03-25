# Freelancer CRM - MVP Design

## Screen 1: Dashboard (Client Overview + Recent Activity)
```
┌─────────────────────────────────────────────────────────┐
│  Freelancer CRM Dashboard                               │
├─────────────────────────────────────────────────────────┤
│  Welcome Sarah!                                         │
│                                                         │
│  [+] New Client   [+] New Project   [+] New Invoice     │
│                                                         │
│  ┌─────────────────────────────────────────────────────┐ │
│  │  CLIENTS (3)                                        │ │
│  │  ┌──────────────────────┐  ┌──────────────────────┐ │ │
│  │  │  🧑‍💼 Client A       │  │  🧑‍💼 Client B       │ │ │
│  │  │  5 projects         │  │  3 projects         │ │ │
│  │  │  Last contact: 2d   │  │  Last contact: 1w   │ │ │
│  │  └──────────────────────┘  └──────────────────────┘ │ │
│  │  ┌──────────────────────┐  ┌──────────────────────┐ │ │
│  │  │  🧑‍💼 Client C       │  │  🧑‍💼 Client D       │ │ │
│  │  │  2 projects         │  │  1 project          │ │ │
│  │  │  Last contact: 3d   │  │  Last contact: 3w   │ │ │
│  │  └──────────────────────┘  └──────────────────────┘ │ │
│  └─────────────────────────────────────────────────────┘ │
│                                                         │
│  ┌─────────────────────────────────────────────────────┐ │
│  │  RECENT ACTIVITY                                    │ │
│  │  • Client A: Sent invoice #INV-001                  │ │
│  │  • Client B: Updated project timeline               │ │
│  │  • Client C: Added note on content revision         │ │
│  │  • Client D: Scheduled meeting                      │ │
│  └─────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────┘
```

## Screen 2: Client Profile & Project Management
```
┌─────────────────────────────────────────────────────────┐
│  Client Profile: Client A                               │
├─────────────────────────────────────────────────────────┤
│  🧑‍💼 Client A - Web Development                          │
│  Email: clienta@example.com                              │
│  Phone: +1 (555) 123-4567                               │
│  Location: New York, NY                                 │
│                                                         │
│  ┌─────────────────────────────────────────────────────┐ │
│  │  PROJECTS                                            │ │
│  │  ┌─────────────────────────────────────────────────┐ │ │
│  │  │  📝 Website Redesign - Active                   │ │ │
│  │  │  Due: 2024-06-15                                 │ │ │
│  │  │  Status: In Progress                             │ │ │
│  │  │  Last Updated: 2024-05-20                        │ │ │
│  │  └─────────────────────────────────────────────────┘ │ │
│  │  ┌─────────────────────────────────────────────────┐ │ │
│  │  │  📝 Landing Page - Completed                     │ │ │
│  │  │  Due: 2024-04-10                                 │ │ │
│  │  │  Status: Completed                               │ │ │
│  │  │  Last Updated: 2024-04-10                        │ │ │
│  │  └─────────────────────────────────────────────────┘ │ │
│  └─────────────────────────────────────────────────────┘ │
│                                                         │
│  [+] Add Project   [+] Add Communication   [+] Invoice  │
└─────────────────────────────────────────────────────────┘
```

## Screen 3: Project Timeline & Notes
```
┌─────────────────────────────────────────────────────────┐
│  Project: Website Redesign                              │
├─────────────────────────────────────────────────────────┤
│  Status: In Progress   Due: 2024-06-15                  │
│                                                         │
│  ┌─────────────────────────────────────────────────────┐ │
│  │  TIMELINE                                            │ │
│  │  ┌─────────────────────────────────────────────────┐ │ │
│  │  │  📅 2024-05-01: Kickoff Meeting                 │ │ │
│  │  │  📅 2024-05-10: Design Mockups Shared           │ │ │
│  │  │  📅 2024-05-20: Feedback Received               │ │ │
│  │  │  📅 2024-06-01: Final Design Approval           │ │ │
│  │  └─────────────────────────────────────────────────┘ │ │
│  └─────────────────────────────────────────────────────┘ │
│                                                         │
│  ┌─────────────────────────────────────────────────────┐ │
│  │  NOTES                                               │ │
│  │  • Client requested dark mode option               │ │
│  │  • Need to clarify branding guidelines             │ │
│  │  • Meeting scheduled for Friday                    │ │
│  └─────────────────────────────────────────────────────┘ │
│                                                         │
│  [+] Add Note   [+] Add Milestone   [+] View History    │
└─────────────────────────────────────────────────────────┘
```

## Screen 4: Invoicing
```
┌─────────────────────────────────────────────────────────┐
│  Create Invoice                                         │
├─────────────────────────────────────────────────────────┤
│  From Project: Website Redesign                         │
│  Client: Client A                                       │
│                                                         │
│  ┌─────────────────────────────────────────────────────┐ │
│  │  INVOICE DETAILS                                     │ │
│  │  Invoice #: INV-001                                  │ │
│  │  Date: 2024-05-20                                    │ │
│  │  Due Date: 2024-06-20                                │ │
│  │  Status: Draft                                       │ │
│  └─────────────────────────────────────────────────────┘ │
│                                                         │
│  ┌─────────────────────────────────────────────────────┐ │
│  │  ITEMS                                               │ │
│  │  • Design Mockups - $1,200                          │ │
│  │  • Final Design Approval - $800                     │ │
│  │  ┌─────────────────────────────────────────────────┐ │ │
│  │  │  Total: $2,000                                  │ │ │
│  │  └─────────────────────────────────────────────────┘ │ │
│  └─────────────────────────────────────────────────────┘ │
│                                                         │
│  [+] Save Draft   [+] Send Invoice   [+] View History   │
└─────────────────────────────────────────────────────────┘
```

## Screen 5: Communication History
```
┌─────────────────────────────────────────────────────────┐
│  Communication History - Client A                       │
├─────────────────────────────────────────────────────────┤
│  📧 Email: clienta@example.com                          │
│  📞 Phone: +1 (555) 123-4567                           │
│  🗓️ Last Contact: 2024-05-20                            │
│                                                         │
│  ┌─────────────────────────────────────────────────────┐ │
│  │  MESSAGES                                            │ │
│  │  ┌─────────────────────────────────────────────────┐ │ │
│  │  │  📨 2024-05-20: Feedback on mockups             │ │ │
│  │  │  "The color scheme needs to match our brand"   │ │ │
│  │  └─────────────────────────────────────────────────┘ │ │
│  │  ┌─────────────────────────────────────────────────┐ │ │
│  │  │  📨 2024-05-15: Project kickoff                 │ │ │
│  │  │  "Let's schedule a meeting to discuss scope"   │ │ │
│  │  └─────────────────────────────────────────────────┘ │ │
│  │  ┌─────────────────────────────────────────────────┐ │ │
│  │  │  📨 2024-05-10: Design mockups shared          │ │ │
│  │  │  "Please review and let me know what you think"│ │ │
│  │  └─────────────────────────────────────────────────┘ │ │
│  └─────────────────────────────────────────────────────┘ │
│                                                         │
│  [+] New Message   [+] Add Note   [+] View Project      │
└─────────────────────────────────────────────────────────┘
```

## Key Design Decisions:
1. **Combined Screens**: Client Profile + Project Management (Screen 2) and Communication History (Screen 5) are combined to reduce navigation
2. **Minimal Navigation**: All screens accessible from Dashboard
3. **Clear Visual Hierarchy**: Important actions (Add, Create) are prominent
4. **Mobile-Friendly**: Simple layout with clear sections
5. **Consistent Design Language**: Same styling across all screens
6. **Action-Oriented**: Each screen focuses on one primary task with clear next steps