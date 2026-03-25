# Prioritized Fix List

### **Critical (must fix before ship)**  
**1. Broken Access Control & Auth Failures**  
*Fix:* Implement JWT or session-based authentication middleware and protect all API endpoints with `requireAuth`. Add role-based access control for admin operations.  

**2. Security Vulnerabilities (OWASP A01, A07, A03)**  
*Fix:* Add input sanitization, parameterized queries, and authentication to all handlers to prevent unauthorized access and injection attacks.  

---

### **High (fix in v1.1)**  
**1. Missing ARIA Labels & Form Labels**  
*Fix:* Add `aria-label` or `aria-labelledby` to interactive elements and `<label>` tags for all form inputs.  

**2. Missing Semantic HTML & Heading Structure**  
*Fix:* Wrap content in semantic HTML elements (`<header>`, `<main>`, `<section>`, `<footer>`) and ensure proper heading hierarchy (`h1`, `h2`, etc.).  

**3. Keyboard Navigation & Focus Management**  
*Fix:* Add `tabindex="0"` to interactive elements and ensure focus styles are visible.  

**4. Missing Alt Text on Images**  
*Fix:* Add `alt=""` attributes to all `<img>` tags.  

**5. Color Contrast Issues**  
*Fix:* Ensure all text meets 4.5:1 contrast ratio against background colors.  

---

### **Medium (nice to have)**  
**1. Visual Hierarchy & Layout Improvements**  
*Fix:* Add prominent primary button styles and use semantic HTML structure to improve layout and visual hierarchy.  

**2. Empty States & Whitespace**  
*Fix:* Add empty state messages and improve whitespace usage for better readability and UX.  

**3. Typography & Color Consistency**  
*Fix:* Define consistent typography and color palette for UI elements.  

---

### **Low (polish)**  
**1. Form Feedback & Validation**  
*Fix:* Add visual feedback on form submission and validation messages.  

**2. Responsive Design Enhancements**  
*Fix:* Ensure responsive behavior across devices using CSS media queries.  

**3. Interaction Feedback (Hover/Focus States)**  
*Fix:* Add hover and focus states for interactive elements to improve usability.