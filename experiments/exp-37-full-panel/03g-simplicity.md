### ✅ Simplified Feature List + Implementation Notes

---

#### **1. Search Clients**
- **Address**: `textarea` input field with `onInput` handler  
- **Implementation**:  
  ```js
  const filteredClients = clients.filter(client => 
    client.name.toLowerCase().includes(searchTerm.toLowerCase())
  );
  ```
- **Status**: `string` (search term)

---

#### **2. Add New Client**
- **Address**: `button` with `onClick` to open modal/form  
- **Implementation**:  
  - Modal with fields: `name`, `email`, `notes`  
  - On submit: push to `clients` array  
- **Status**: `string` (form state)

---

#### **3. Filter Clients**
- **Address**: `select` dropdown with filter options  
- **Implementation**:  
  ```js
  const filteredClients = clients.filter(client => 
    client.status === filterStatus
  );
  ```
- **Status**: `string` (filter value)

---

#### **4. View Client Details**
- **Address**: Click client → show detail in modal or new screen  
- **Implementation**:  
  - Display: `name`, `email`, `notes`, `projects`  
- **Status**: `object` (selected client)

---

#### **5. Create Invoice**
- **Address**: `button` → opens invoice form  
- **Implementation**:  
  - Fields: `client`, `items` (description, qty, rate), `dueDate`  
  - Save to `invoices` array  
- **Status**: `object` (invoice data)

---

#### **6. Invoice Line Items**
- **Address**: Invoice detail screen  
- **Implementation**:  
  - List of items with `description`, `quantity`, `rate`, `amount`  
- **Status**: `array` of item objects

---

#### **7. Invoice Numbering**
- **Address**: Auto-generated on save  
- **Implementation**:  
  ```js
  `INV-${invoices.length + 1}`  
  ```
- **Status**: `string` (invoice number)

---

#### **8. Print Invoice**
- **Address**: `button` → triggers `window.print()`  
- **Implementation**:  
  - Use browser print dialog  
- **Status**: `boolean` (print triggered)

---

#### **9. Send Invoice**
- **Address**: `button` → opens email client or modal  
- **Implementation**:  
  - Pre-fill email with invoice PDF or link  
- **Status**: `string` (email action)

---

#### **10. Mark Invoice as Paid**
- **Address**: `button` in invoice detail  
- **Implementation**:  
  - Update `invoice.status = 'paid'`  
- **Status**: `string` (status update)

---

#### **11. Void Invoice**
- **Address**: `button` in invoice detail  
- **Implementation**:  
  - Update `invoice.status = 'void'`  
- **Status**: `string` (status update)

---

#### **12. Due Date Display**
- **Address**: Invoice detail screen  
- **Implementation**:  
  - Show `dueDate` field as `YYYY-MM-DD`  
- **Status**: `string` (date format)

---

#### **13. Toast Notifications**
- **Address**: Any action (save, delete, etc.)  
- **Implementation**:  
  ```js
  const toast = document.createElement('div');
  toast.textContent = 'Success!';
  document.body.appendChild(toast);
  setTimeout(() => toast.remove(), 3000);
  ```
- **Status**: `string` (message)

---

#### **14. Breadcrumbs**
- **Address**: Top of screen (if needed)  
- **Implementation**:  
  - Simple text path: `Dashboard > Clients > Client Name`  
- **Status**: `string` (path)

---

#### **15. Consistent Layout**
- **Address**: All screens  
- **Implementation**:  
  - Use same padding/margins, font sizes, button styles  
- **Status**: `boolean` (layout consistency)

---

### ✅ Summary: Minimal Features for MVP

| Feature                  | Simplest Implementation |
|--------------------------|--------------------------|
| Search Clients           | `textarea` + `filter()` |
| Add New Client           | `button` + form modal   |
| Filter Clients           | `select` + `filter()`   |
| View Client Details      | Modal or detail screen  |
| Create Invoice           | Form + save to array    |
| Invoice Line Items       | List of objects         |
| Invoice Numbering        | Auto-increment string   |
| Print Invoice            | `window.print()`        |
| Send Invoice             | Email pre-fill          |
| Mark Paid                | Button + status update  |
| Void Invoice             | Button + status update  |
| Due Date Display         | `date` field            |
| Toast Notifications      | `div` + `setTimeout()`  |
| Breadcrumbs              | Text path               |
| Consistent Layout        | CSS + component reuse   |

Let me know if you want this in JSON or code-ready format.