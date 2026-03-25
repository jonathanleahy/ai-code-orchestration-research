### **Changes from Each Reviewer (1 Line Each)**

- **Dev:** "MVP is clean and scalable; focus on missing core features like invoicing and activity logging."
- **QA:** "Missing edge cases and key UI elements like search filters and validation states."
- **Market:** "Missing freelancer-specific features and invoicing automation — need to define a clear value proposition."
- **Domain Expert:** "Invoicing is completely missing — must include creation, numbering, and basic fields."

---

### **Revised Feature List (Must-Have for v1)**

1. **Client Management**
   - Add, Edit, Delete Clients
2. **Activity Logging**
   - Log time spent on projects
3. **Invoice Creation**
   - Create invoices with line items (description, quantity, rate)
   - Automatic invoice numbering (e.g., INV-001)
   - Subtotal, tax, total calculations
4. **Search & Filter**
   - Search clients, projects, invoices
5. **Project Tracking**
   - Associate activities with projects

---

### **Invoicing Workflow (from Domain Expert): What's IN v1**

- **Create Invoice Screen**
  - Fields: Invoice Date, Due Date, Client, Line Items (Description, Quantity, Rate)
  - Auto-calculate Subtotal, Tax, Total
  - Auto-generate Invoice Number (INV-001, etc.)
- **No Advanced Features (yet):**
  - No recurring invoices
  - No PDF export
  - No payment status tracking

---

### **Revised Screen List**

1. **Dashboard**
2. **Client List / Add Client**
3. **Client Profile View**
4. **Activity Log / Add Activity**
5. **Project List / Add Project**
6. **Invoice List / Add Invoice**
7. **Invoice Detail View**
8. **Search Results Page**