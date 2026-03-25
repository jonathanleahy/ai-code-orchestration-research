# Screen Wireframes

Here are the ASCII wireframes for the **Invoice Generator MVP**, covering the core screens with layout, interactive elements, and sample data. Each screen includes a name, URL route, and a simple ASCII representation.

---

### 1. **Dashboard**
**URL Route:** `/dashboard`

```
+-------------------------------------------------------------+
|                    INVOICE GENERATOR DASHBOARD               |
+-------------------------------------------------------------+
|                                                             |
|  [Create New Invoice]     [View All Invoices]     [Clients] |
|                                                             |
|  [Recent Invoices]                                          |
|  ┌─────────────────────────────────────────────────────────┐ |
|  │ Invoice #INV-001 | Client: Acme Corp        | $1,200.00 │ |
|  │ Due Date: 2025-04-15 | Status: Overdue        |         │ |
|  └─────────────────────────────────────────────────────────┘ |
|  ┌─────────────────────────────────────────────────────────┐ |
|  │ Invoice #INV-002 | Client: Beta Inc         | $850.00  │ |
|  │ Due Date: 2025-04-20 | Status: Paid           |         │ |
|  └─────────────────────────────────────────────────────────┘ |
|                                                             |
|  [View All Invoices] [Settings]                             |
+-------------------------------------------------------------+
```

---

### 2. **Create Invoice**
**URL Route:** `/create-invoice`

```
+-------------------------------------------------------------+
|                    CREATE NEW INVOICE                        |
+-------------------------------------------------------------+
|                                                             |
|  [Client] [Select Client]                                   |
|  [Invoice Date] [2025-04-01]                                |
|  [Due Date] [2025-04-15]                                    |
|                                                             |
|  [Items]                                                   |
|  ┌─────────────────────────────────────────────────────────┐ |
|  │ Description     | Quantity | Unit Price | Total     │ |
|  │ Design Work     | 10       | $100.00    | $1,000.00 │ |
|  │ Travel          | 2        | $50.00     | $100.00   │ |
|  └─────────────────────────────────────────────────────────┘ |
|                                                             |
|  [Add Item] [Remove Item]                                   |
|                                                             |
|  [Tax Rate] [8.5%]                                          |
|  [Notes] [Project completion milestone]                     |
|                                                             |
|  [Preview] [Save Draft] [Send Invoice]                      |
+-------------------------------------------------------------+
```

---

### 3. **Invoice List**
**URL Route:** `/invoices`

```
+-------------------------------------------------------------+
|                        INVOICE LIST                          |
+-------------------------------------------------------------+
|                                                             |
|  [Search] [Filter by Status] [Sort by Date]                 |
|                                                             |
|  ┌─────────────────────────────────────────────────────────┐ |
|  │ Invoice #INV-001 | Client: Acme Corp        | $1,200.00 │ |
|  │ Due Date: 2025-04-15 | Status: Overdue        |         │ |
|  └─────────────────────────────────────────────────────────┘ |
|  ┌─────────────────────────────────────────────────────────┐ |
|  │ Invoice #INV-002 | Client: Beta Inc         | $850.00  │ |
|  │ Due Date: 2025-04-20 | Status: Paid           |         │ |
|  └─────────────────────────────────────────────────────────┘ |
|  ┌─────────────────────────────────────────────────────────┐ |
|  │ Invoice #INV-003 | Client: Gamma Ltd        | $3,200.00│ |
|  │ Due Date: 2025-04-25 | Status: Unpaid         |         │ |
|  └─────────────────────────────────────────────────────────┘ |
|                                                             |
|  [New Invoice] [Export CSV]                                 |
+-------------------------------------------------------------+
```

---

### 4. **Invoice Detail / Preview**
**URL Route:** `/invoice/INV-001`

```
+-------------------------------------------------------------+
|                    INVOICE #INV-001                          |
+-------------------------------------------------------------+
|                                                             |
|  [Client: Acme Corp]                                        |
|  [Email: contact@acme.com]                                  |
|  [Address: 123 Main St, City, Country]                      |
|                                                             |
|  [Invoice Date: 2025-04-01] [Due Date: 2025-04-15]           |
|                                                             |
|  [Items]                                                   |
|  ┌─────────────────────────────────────────────────────────┐ |
|  │ Description     | Quantity | Unit Price | Total     │ |
|  │ Design Work     | 10       | $100.00    | $1,000.00 │ |
|  │ Travel          | 2        | $50.00     | $100.00   │ |
|  └─────────────────────────────────────────────────────────┘ |
|                                                             |
|  [Subtotal] $1,100.00                                       |
|  [Tax] $93.50                                               |
|  [Total] $1,193.50                                          |
|                                                             |
|  [Notes] [Project completion milestone]                     |
|                                                             |
|  [Download PDF] [Send Email] [Mark as Paid]                 |
+-------------------------------------------------------------+
```

---

### 5. **Client View**
**URL Route:** `/clients`

```
+-------------------------------------------------------------+
|                        CLIENTS LIST                          |
+-------------------------------------------------------------+
|                                                             |
|  [Search] [Add New Client]                                  |
|                                                             |
|  ┌─────────────────────────────────────────────────────────┐ |
|  │ Client Name: Acme Corp | Contact: John Doe         │ |
|  │ Email: john@acme.com   | Phone: (123) 456-7890     │ |
|  │ Address: 123 Main St, City, Country                     │ |
|  │ Total Invoices: 3 | Paid: 1 | Overdue: 1            │ |
|  └─────────────────────────────────────────────────────────┘ |
|  ┌─────────────────────────────────────────────────────────┐ |
|  │ Client Name: Beta Inc  | Contact: Jane Smith       │ |
|  │ Email: jane@beta.com   | Phone: (987) 654-3210     │ |
|  │ Address: 456 Oak Ave, City, Country                     │ |
|  │ Total Invoices: 2 | Paid: 2 | Overdue: 0            │ |
|  └─────────────────────────────────────────────────────────┘ |
|                                                             |
|  [View Client Details] [Edit] [Delete]                       |
+-------------------------------------------------------------+
```

---

### 6. **Settings**
**URL Route:** `/settings`

```
+-------------------------------------------------------------+
|                        SETTINGS                              |
+-------------------------------------------------------------+
|                                                             |
|  [General Settings] [Email Templates] [Tax Rates]           |
|                                                             |
|  [Email Settings]                                           |
|  [SMTP Host] [smtp.example.com]                             |
|  [Port] [587]                                               |
|  [Username] [user@example.com]                              |
|  [Password] [••••••••]                                       |
|                                                             |
|  [Invoice Templates]                                        |
|  [Default Template] [Custom Template]                       |
|                                                             |
|  [QuickBooks Integration]                                   |
|  [Status: Connected] [Disconnect]                           |
|                                                             |
|  [Save Settings]                                             |
+-------------------------------------------------------------+
```

---

### Empty States

#### Empty Invoice List
```
+-------------------------------------------------------------+
|                        INVOICE LIST                          |
+-------------------------------------------------------------+
|                                                             |
|  [Search] [Filter by Status] [Sort by Date]                 |
|                                                             |
|  [No invoices found. Create your first one!]                |
|  [New Invoice]                                              |
+-------------------------------------------------------------+
```

#### Empty Client List
```
+-------------------------------------------------------------+
|                        CLIENTS LIST                          |
+-------------------------------------------------------------+
|                                                             |
|  [Search] [Add New Client]                                  |
|                                                             |
|  [No clients found. Add your first client!]                 |
|  [Add New Client]                                           |
+-------------------------------------------------------------+
```

---

Let me know if you'd like these exported in a file format or want to add more screens like a login or client detail view.