Here are the **ASCII wireframes** for each screen of the **Freelancer CRM MVP**, based on the features and structure you've outlined.

---

## 📋 1. Dashboard – Client List

**URL:** `/dashboard`

```
+-------------------------------------------------------------+
| Freelancer CRM - Dashboard                                  |
|                                                             |
| [+] Add New Client         [Search] [Notifications]         |
|                                                             |
| +----------------+----------------+----------------+   |
| | Client Name    | Email            | Phone          |   |
| +----------------+----------------+----------------+   |
| | John Doe       | john@example.com | 123-456-7890   |   |
| | Jane Smith     | jane@example.com | 098-765-4321   |   |
| | Alex Johnson   | alex@example.com | 555-555-5555   |   |
| +----------------+----------------+----------------+   |
|                                                             |
| [Edit] [Delete] [View Details]                              |
+-------------------------------------------------------------+
```

**Actions Available:**
- Add New Client (button)
- View Details (per client)
- Edit/Delete (per client)

---

## 📝 2. Add Client Form

**URL:** `/clients/new`

```
+-------------------------------------------------------------+
| Add New Client                                              |
|                                                             |
| Name: [_________________________]                           |
| Email: [_________________________]                          |
| Phone: [_________________________]                          |
| Payment Method: [_________________________]                |
|                                                             |
| [Save] [Cancel]                                             |
+-------------------------------------------------------------+
```

**Actions Available:**
- Save (submit form)
- Cancel (go back to dashboard)

---

## 👤 3. Client Profile Page (with Tabs)

**URL:** `/clients/:id`

```
+-------------------------------------------------------------+
| Client Profile - John Doe                                   |
|                                                             |
| [Back to Dashboard] [Edit] [Delete]                         |
|                                                             |
| [Details] [History] [Billing]                               |
|                                                             |
| Details Tab:                                                |
| Name: John Doe                                              |
| Email: john@example.com                                     |
| Phone: 123-456-7890                                         |
| Payment Info: PayPal                                        |
+-------------------------------------------------------------+
```

**Actions Available:**
- Back to Dashboard
- Edit/Delete
- Switch between tabs (Details / History / Billing)

---

## 📄 4. Client Details Tab

```
+-------------------------------------------------------------+
| Client Profile - John Doe - Details                         |
|                                                             |
| Name: John Doe                                              |
| Email: john@example.com                                     |
| Phone: 123-456-7890                                         |
| Payment Info: PayPal                                        |
|                                                             |
| [Edit] [Delete]                                             |
+-------------------------------------------------------------+
```

---

## 🕒 5. Activity History Tab

```
+-------------------------------------------------------------+
| Client Profile - John Doe - Activity History                |
|                                                             |
| +----------------+----------------+----------------+   |
| | Type           | Content          | Timestamp        |   |
| +----------------+----------------+----------------+   |
| | Email          | Follow-up email  | 2025-04-01 10:00 |   |
| | Meeting        | Discussed project| 2025-03-28 14:30 |   |
| | Call           | Left voicemail   | 2025-03-25 09:15 |   |
| +----------------+----------------+----------------+   |
|                                                             |
| [Add New Activity]                                          |
+-------------------------------------------------------------+
```

**Actions Available:**
- Add New Activity (button)
- View/Edit/Delete activities

---

## 💵 6. Billing Tab

```
+-------------------------------------------------------------+
| Client Profile - John Doe - Billing                         |
|                                                             |
| +----------------+--------+----------------+----------+   |
| | Invoice ID     | Amount | Due Date       | Status   |   |
| +----------------+--------+----------------+----------+   |
| | INV-001        | $500   | 2025-04-15     | Unpaid   |   |
| | INV-002        | $300   | 2025-03-30     | Paid     |   |
| +----------------+--------+----------------+----------+   |
|                                                             |
| [Create Invoice]                                            |
+-------------------------------------------------------------+
```

**Actions Available:**
- Create Invoice (button)
- View Invoice Details

---

## ✏️ 7. Edit Client Form

**URL:** `/clients/:id/edit`

```
+-------------------------------------------------------------+
| Edit Client - John Doe                                      |
|                                                             |
| Name: [John Doe]                                            |
| Email: [john@example.com]                                   |
| Phone: [123-456-7890]                                       |
| Payment Method: [PayPal]                                    |
|                                                             |
| [Update] [Cancel]                                           |
+-------------------------------------------------------------+
```

**Actions Available:**
- Update (submit form)
- Cancel (go back to profile)

---

## ❌ 8. Delete Confirmation

**URL:** `/clients/:id/delete`

```
+-------------------------------------------------------------+
| Confirm Delete Client                                       |
|                                                             |
| Are you sure you want to delete John Doe?                   |
| This action cannot be undone.                               |
|                                                             |
| [Yes, Delete] [No, Cancel]                                  |
+-------------------------------------------------------------+
```

**Actions Available:**
- Yes, Delete (confirm deletion)
- No, Cancel (return to client profile)

---

## 🔐 9. Login Page (Optional for MVP)

**URL:** `/login`

```
+-------------------------------------------------------------+
| Freelancer CRM - Login                                      |
|                                                             |
| Username: [_________________________]                       |
| Password: [_________________________]                       |
|                                                             |
| [Login] [Sign Up]                                           |
+-------------------------------------------------------------+
```

**Actions Available:**
- Login
- Sign Up (if needed)

---

## 🧾 10. Invoice Creation Form

**URL:** `/invoices/new`

```
+-------------------------------------------------------------+
| Create New Invoice                                          |
|                                                             |
| Client: [John Doe]                                          |
| Amount: [_______]                                           |
| Due Date: [_______]                                         |
|                                                             |
| [Create Invoice] [Cancel]                                   |
+-------------------------------------------------------------+
```

**Actions Available:**
- Create Invoice (submit)
- Cancel (go back)

---

Let me know if you'd like these exported into a markdown file or want to generate HTML mockups from them.