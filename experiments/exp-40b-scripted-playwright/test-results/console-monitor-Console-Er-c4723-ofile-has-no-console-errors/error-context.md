# Page snapshot

```yaml
- generic [ref=e2]:
  - link "← Back to Dashboard" [ref=e3] [cursor=pointer]:
    - /url: /
  - banner [ref=e4]:
    - heading "Test" [level=1] [ref=e5]
  - generic [ref=e6]:
    - button "Details" [ref=e7] [cursor=pointer]
    - button "Activity" [ref=e8] [cursor=pointer]
    - button "Invoices" [active] [ref=e9] [cursor=pointer]
  - generic [ref=e10]:
    - heading "Create Invoice" [level=2] [ref=e11]
    - generic [ref=e12]:
      - generic [ref=e13]:
        - generic [ref=e14]: Invoice Number
        - textbox "e.g. INV-001" [ref=e15]
      - generic [ref=e16]:
        - generic [ref=e17]: Amount
        - spinbutton [ref=e18]
      - generic [ref=e19]:
        - generic [ref=e20]: Description
        - textbox [ref=e21]
      - generic [ref=e22]:
        - generic [ref=e23]: Due Date
        - textbox [ref=e24]
      - button "Create Invoice" [ref=e25] [cursor=pointer]
    - heading "Invoices" [level=2] [ref=e26]
    - generic [ref=e27]: No invoices yet.
```