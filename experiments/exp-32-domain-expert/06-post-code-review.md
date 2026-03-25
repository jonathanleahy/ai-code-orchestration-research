### **UI/UX Expert: Layout, hierarchy, forms, feedback. Top 5 issues.**

1. **Tab navigation lacks keyboard support**  
   Add `tabindex="0"` and `onkeydown` handlers for tab switching.

2. **No visual feedback on tab hover or click**  
   Enhance `.tab:hover` and `.tab.active` styles with transitions or more distinct borders.

3. **Content area has no loading or empty state indicators**  
   Add a placeholder or spinner when content is loading or empty.

4. **Form fields (if any) lack labels or proper structure**  
   Ensure all inputs have associated `<label>` elements with `for` attributes.

5. **Inconsistent spacing and typography hierarchy**  
   Use consistent margin/padding and font weights for headers and body text.

---

### **Accessibility: ARIA labels, form labels, keyboard nav, contrast. Top 5 issues.**

1. **Missing ARIA roles and labels for interactive elements**  
   Add `role="tab"` and `aria-selected="true"` to tabs.

2. **No semantic HTML structure for tabs**  
   Use `<nav>` and `<section>` tags to improve screen reader support.

3. **Low color contrast in some UI elements**  
   Adjust colors to meet WCAG 2.1 AA contrast ratio (4.5:1 minimum).

4. **No keyboard focus indicators for tabs**  
   Ensure focus styles are visible and accessible via tab key.

5. **Invoice print view lacks accessible heading structure**  
   Use `<h1>` for main title and `<h2>` for subheadings like "Invoice #".

---

### **OWASP Security: XSS, CSRF, CSP, auth, injection. Top 5 issues.**

1. **Direct HTML injection in client data (Name, Address)**  
   Sanitize all user-provided data using `html.EscapeString()` before inserting into HTML.

2. **No CSRF protection for form submissions**  
   Implement CSRF tokens in forms and validate them server-side.

3. **Inline JavaScript in invoice print view**  
   Move `window.print()` to an external script or use event listeners instead.

4. **Potential XSS via invoice number or line item descriptions**  
   Escape dynamic content like `invoice.InvoiceNumber`, `item.Description` using `html.EscapeString()`.

5. **Lack of Content Security Policy (CSP)**  
   Add a strict CSP header to prevent inline script execution and XSS.

--- 

Let me know if you'd like sample fixes or code snippets for any of these!