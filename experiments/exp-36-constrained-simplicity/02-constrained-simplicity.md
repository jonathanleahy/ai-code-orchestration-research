### V1 Features (must-have, simplified)

#### CLIENT MANAGEMENT
- **Add Client**  
  - Form with fields: Name, Email, Phone, Company, Address, Notes  
  - One textarea for Address (not 5 separate fields)  
  - Save to local storage or in-memory DB  

- **Edit Client**  
  - Reuse Add Client form with pre-filled data  
  - Check if ID exists → update instead of create  

- **Delete Client**  
  - Delete button next to each client in list  
  - Confirm dialog before deletion  

- **Search Clients**  
  - Simple text input that filters client list by name/email/company using `strings.Contains`  

- **View Client Profile**  
  - Click client → show full profile in a modal or new page  

- **Phone Number Field**  
  - Text input field  

- **Email Address Field**  
  - Text input field  

- **Address Field**  
  - Single textarea field  

- **Company Field**  
  - Text input field  

- **Notes Field**  
  - Textarea field  

#### ACTIVITY HISTORY
- **Log Call**  
  - Simple form with dropdown for activity type (Call), description textarea  

- **Log Email**  
  - Simple form with dropdown for activity type (Email), description textarea  

- **Log Meeting**  
  - Simple form with dropdown for activity type (Meeting), description textarea  

- **Log Note**  
  - Simple form with dropdown for activity type (Note), description textarea  

- **View Timeline**  
  - List of all activities for a client in chronological order  

- **Edit Activity**  
  - Reuse log form with pre-filled values  

- **Delete Activity**  
  - Delete button next to each activity  

#### INVOICING
- **Create Invoice**  
  - Form with: Client selection dropdown, Due Date, Line Items (description, quantity, price)  

- **Line Items**  
  - Reusable form section for adding line items (description, quantity, price)  

- **Mark Sent**  
  - Button to update invoice status to "Sent"  

- **Mark Paid**  
  - Button to update invoice status to "Paid"  

- **Print View**  
  - Use browser print (`@media print` CSS)  

- **Due Date**  
  - Date picker input  

---

### V2 Features (deferred)

- **Delete Client**  
- **Delete Activity**  
- **Void Invoice**  
- **Address Field** (as separate fields)  
- **Search Clients** (full-text search)  
- **Activity Timeline** (advanced filtering or grouping)  

---

### Screen Plan (minimum screens needed)

1. **Dashboard**  
   - Client list with search  
   - Add new client form  
   - Quick access to create invoice or log activity  

2. **Client Profile Page**  
   - View all client details  
   - View timeline of activities  
   - Edit client info  

3. **Activity Log Page**  
   - Form to log new activity  
   - View timeline of activities for a client  

4. **Invoice Page**  
   - Create new invoice  
   - View invoice details  
   - Mark as sent/paid  
   - Print view  

---

### Implementation Notes (simplest way for each)

| Feature | Implementation |
|--------|----------------|
| **Address Field** | Single textarea field |
| **Edit Client** | Reuse add form with ID check |
| **Invoice Print** | Use browser print with CSS `@media print` |
| **Search Clients** | Simple `strings.Contains` filter on name/email/company |
| **Activity Log** | Simple form with type dropdown + description textarea |
| **Mark as Paid** | Update status string (no payment processing) |
| **Overdue** | Compare due date to current date |
| **Client List** | In-memory array or localStorage |
| **Invoice Line Items** | Reusable form section with add/remove buttons |
| **Dashboard** | Single-page layout with client list + form + search |
| **Client Profile** | Modal or new page with full details |
| **Activity Timeline** | List of activities sorted by date |
| **Delete Client/Activity** | Confirm dialog + remove from array/localStorage |
| **Void Invoice** | Not in v1 |
| **Full-text Search** | Not in v1 |
| **Separate Address Fields** | Not in v1 |
| **Advanced Timeline Filters** | Not in v1 |