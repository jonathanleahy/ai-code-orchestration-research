# Screens v1

# MVP Screen Wireframes

## 1. **Login Screen**
**URL Route**: `/login`

**Purpose**: Allow users to authenticate into their account

```
┌─────────────────────────────────────────────────────────────┐
│                        [INVOICE APP]                        │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  [Email]  ┌───────────────────────────────────────┐       │
│           │              example@email.com        │       │
│           └───────────────────────────────────────┘       │
│                                                             │
│  [Password] ┌───────────────────────────────────────┐    │
│             │              ••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••......# MVP Screen Wireframes

## 1. **Login Screen**
**URL Route**: `/login`

**Purpose**: Allow users to authenticate into their account

```
┌─────────────────────────────────────────────────────────────┐
│                        [INVOICE APP]                        │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  [Email]  ┌───────────────────────────────────────┐       │
│           │              example@email.com        │       │
│           └───────────────────────────────────────┘       │
│                                                             │
│  [Password]  ┌───────────────────────────────────────┐   │
│               │              •••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••••......Here are the **screen wireframes** for your **MVP plan**, organized by screen with **ASCII wireframes**, **purpose**, and **interactions**. These are tailored to meet the needs of **Maria (UX Designer)** and **James (Web Developer)**, focusing on **core pain points** like time-consuming invoicing, automation, and security.

---

## 🧭 1. **Login / Signup Screen**  
**URL Route:** `/login`  
**Purpose:** Allow users to log in or sign up for an account.

```
+-------------------------------------------------------+
|                    [LOGO]                             |
|                                                       |
|   [Email Input]                                       |
|   [Password Input]                                    |
|   [Login Button]                                      |
|                                                       |
|   [Sign Up]                                           |
|   [Forgot Password?]                                  |
|                                                       |
|   [Google Login] [GitHub Login]                       |
+-------------------------------------------------------+
```

**Key Interactions:**
- Clicking **Login** → Validates credentials and redirects to dashboard.
- Clicking **Sign Up** → Redirects to `/signup`.
- Clicking **Forgot Password** → Sends reset link.

---

## 🧭 2. **Onboarding / Setup Wizard**  
**URL Route:** `/onboarding`  
**Purpose:** Guide new users through initial setup (name, company info, integration preferences).

```
+-------------------------------------------------------+
|              [Setup Wizard: Step 1/3]                 |
|                                                       |
|   What's your name? [John Doe]                        |
|   Company Name [Creative Solutions LLC]               |
|                                                       |
|   [Next]                                              |
+-------------------------------------------------------+
```

**Key Interactions:**
- Clicking **Next** → Proceeds to step 2.
- Input validation for required fields.

---

## 🧭 3. **Dashboard**  
**URL Route:** `/dashboard`  
**Purpose:** Overview of recent activity, quick actions, and invoice stats.

```
+-------------------------------------------------------+
| [Header: Logo | User Profile | Notifications]       |
|                                                       |
| [Create Invoice] [Import Time] [QuickBooks Sync]      |
|                                                       |
| [Recent Invoices]                                     |
|  - Invoice #1234 [Client: Acme Corp] [Status: Sent]   |
|  - Invoice #1235 [Client: StartupXYZ] [Status: Paid]  |
|  - Invoice #1236 [Client: DesignCo] [Status: Overdue]|
|                                                       |
| [Quick Stats]                                         |
|  - Total Invoices: 24                                 |
|  - Paid This Month: $3,200                            |
|  - Overdue: $1,200                                    |
+-------------------------------------------------------+
```

**Key Interactions:**
- Clicking **Create Invoice** → Redirects to `/invoices/new`.
- Clicking **Invoice #1234** → Opens `/invoices/1234`.

---

## 🧭 4. **Create Invoice (Manual)**  
**URL Route:** `/invoices/new`  
**Purpose:** Allow users to manually create invoices.

```
+-------------------------------------------------------+
| [Header: Logo | User Profile]                         |
|                                                       |
| [Client Info]                                         |
|   [Select Client] [Add New Client]                    |
|   [Email] [Phone]                                     |
|                                                       |
| [Invoice Items]                                       |
|   [Item Name] [Qty] [Rate] [Total]                    |
|   [Add Line Item]                                     |
|                                                       |
| [Notes]                                               |
|   [Optional Notes Field]                              |
|                                                       |
| [Send Invoice] [Save Draft]                           |
+-------------------------------------------------------+
```

**Key Interactions:**
- Clicking **Add Line Item** → Adds a new row.
- Clicking **Send Invoice** → Sends to client and shows success message.
- Clicking **Save Draft** → Saves as draft in DB.

---

## 🧭 5. **View Sent Invoices**  
**URL Route:** `/invoices`  
**Purpose:** List all created invoices with status and actions.

```
+-------------------------------------------------------+
| [Header: Logo | User Profile]                         |
|                                                       |
| [Filter: All | Sent | Paid | Overdue]                 |
|                                                       |
| [Invoice List]                                        |
|  - #1234 [Acme Corp] [Due: 2025-04-10] [Status: Sent] |
|  - #1235 [StartupXYZ] [Due: 2025-04-05] [Status: Paid]|
|  - #1236 [DesignCo] [Due: 2025-03-28] [Status: Overdue]|
|                                                       |
| [Pagination] [New Invoice]                            |
+-------------------------------------------------------+
```

**Key Interactions:**
- Clicking **#1234** → Opens `/invoices/1234`.
- Clicking **New Invoice** → Redirects to `/invoices/new`.

---

## 🧭 6. **Invoice Detail View**  
**URL Route:** `/invoices/:id`  
**Purpose:** Show full invoice details and actions.

```
+-------------------------------------------------------+
| [Header: Logo | User Profile]                         |
|                                                       |
| [Invoice #1234]                                       |
|   Client: Acme Corp                                   |
|   Email: contact@acme.com                             |
|   Date: 2025-03-15 | Due: 2025-04-10                     |
|                                                       |
| [Line Items]                                          |
|   - Design Review [10 hrs @ $150] = $1,500            |
|   - Wireframing [5 hrs @ $150] = $750                 |
|                                                       |
| [Total: $2,250]                                       |
|                                                       |
| [Actions]                                             |
|   [Send Invoice] [Download PDF] [Edit] [Delete]       |
+-------------------------------------------------------+
```

**Key Interactions:**
- Clicking **Send Invoice** → Sends via email.
- Clicking **Download PDF** → Generates and downloads invoice.
- Clicking **Edit** → Redirects to `/invoices/:id/edit`.

---

## 🧭 7. **Client-Facing Invoice View**  
**URL Route:** `/invoice/:id/public`  
**Purpose:** Secure public view of invoice for clients.

```
+-------------------------------------------------------+
| [Company Logo]                                        |
|                                                       |
| Invoice #1234                                         |
| Date: 2025-03-15 | Due: 2025-04-10                     |
|                                                       |
| [Client Info]                                         |
|   Acme Corp                                           |
|   contact@acme.com                                    |
|                                                       |
| [Line Items]                                          |
|   - Design Review [10 hrs @ $150] = $1,500            |
|   - Wireframing [5 hrs @ $150] = $750                 |
|                                                       |
| [Total: $2,250]                                       |
|                                                       |
| [Payment Options]                                     |
|   [Pay Now with Stripe]                               |
|   [View Invoice PDF]                                  |
+-------------------------------------------------------+
```

**Key Interactions:**
- Clicking **Pay Now** → Redirects to payment gateway (Stripe).
- Clicking **View PDF** → Downloads invoice.

---

## 🧭 8. **Settings / Profile**  
**URL Route:** `/settings`  
**Purpose:** Manage account settings, integrations, and security.

```
+-------------------------------------------------------+
| [Header: Logo | User Profile]                         |
|                                                       |
| [Profile Settings]                                    |
|   Name: John Doe                                      |
|   Email: john@example.com                             |
|   Company: Creative Solutions LLC                     |
|   [Update Profile]                                    |
|                                                       |
| [Integrations]                                        |
|   QuickBooks Online [Connected]                       |
|   [Disconnect]                                        |
|                                                       |
| [Security]                                            |
|   Two-Factor Auth [Enabled]                           |
|   [Manage]                                            |
|                                                       |
| [Billing]                                             |
|   Plan: Basic ($15/month)                             |
|   [Upgrade]                                           |
+-------------------------------------------------------+
```

**Key Interactions:**
- Clicking **Update Profile** → Saves changes.
- Clicking **Disconnect QuickBooks** → Removes integration.
- Clicking **Manage 2FA** → Opens 2FA settings.

---

## 🧭 9. **Recurring Invoices (Basic)**  
**URL Route:** `/recurring`  
**Purpose:** Manage recurring invoices (e.g., retainer clients).

```
+-------------------------------------------------------+
| [Header: Logo | User Profile]                         |
|                                                       |
| [Recurring Invoices]                                  |
|  - Monthly Retainer [Acme Corp] [Status: Active]      |
|  - Weekly Support [StartupXYZ] [Status: Paused]       |
|                                                       |
| [New Recurring Invoice]                               |
+-------------------------------------------------------+
```

**Key Interactions:**
- Clicking **New Recurring Invoice** → Redirects to `/recurring/new`.
- Clicking **Pause/Resume** → Updates status.

---

## 🧭 10. **Empty State - No Invoices Yet**  
**URL Route:** `/invoices` (when empty)  
**Purpose:** Encourage creation of first invoice.

```
+-------------------------------------------------------+
| [Header: Logo | User Profile]                         |
|                                                       |
| [No Invoices Found]                                   |
|                                                       |
|   You haven't created any invoices yet.               |
|   [Create Your First Invoice]                         |
|                                                       |
|   [Import Time from Toggl]                            |
+-------------------------------------------------------+
```

**Key Interactions:**
- Clicking **Create Your First Invoice** → Redirects to `/invoices/new`.

---

## 🧭 11. **Mobile View (Responsive)**  
**URL Route:** `/mobile` *(simulated)*  
**Purpose:** Show how the app looks on mobile devices.

```
+-------------------------------------------------------+
| [Menu Button] [Logo]                                  |
|                                                       |
| [Create Invoice]                                      |
|                                                       |
| [Recent Invoices]                                     |
|   - Invoice #1234 [Client: Acme Corp]                 |
|   - Invoice #1235 [Client: StartupXYZ]                |
|                                                       |
| [Footer Nav]                                          |
| [Dashboard] [Invoices] [Settings]                     |
+-------------------------------------------------------+
```

**Key Interactions:**
- Clicking **Menu Button** → Opens sidebar menu.
- Clicking **Invoices** → Navigates to `/invoices`.

---

Let me know if you'd like these as **Figma wireframes** or **clickable prototypes** next!