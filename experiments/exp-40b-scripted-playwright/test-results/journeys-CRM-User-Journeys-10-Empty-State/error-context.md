# Page snapshot

```yaml
- generic [ref=e2]:
  - banner [ref=e3]:
    - heading "CRM Dashboard" [level=1] [ref=e4]
    - paragraph [ref=e5]: Manage your clients
  - generic [ref=e6]:
    - heading "Add New Client" [level=2] [ref=e7]
    - generic [ref=e8]:
      - generic [ref=e9]:
        - generic [ref=e10]:
          - generic [ref=e11]: Name
          - textbox "Name" [ref=e12]
        - generic [ref=e13]:
          - generic [ref=e14]: Email
          - textbox "Email" [ref=e15]
        - generic [ref=e16]:
          - generic [ref=e17]: Phone
          - textbox "Phone" [ref=e18]
      - button "Add Client" [ref=e19] [cursor=pointer]
  - generic [ref=e20]:
    - heading "Clients List" [level=2] [ref=e21]
    - textbox "Search clients by name or email..." [active] [ref=e23]: xyznonexistent99999
    - table [ref=e24]:
      - rowgroup [ref=e25]:
        - row "Name Email Phone Created Action" [ref=e26]:
          - columnheader "Name" [ref=e27]
          - columnheader "Email" [ref=e28]
          - columnheader "Phone" [ref=e29]
          - columnheader "Created" [ref=e30]
          - columnheader "Action" [ref=e31]
      - rowgroup [ref=e32]:
        - row "No clients found" [ref=e33]:
          - cell "No clients found" [ref=e34]
```