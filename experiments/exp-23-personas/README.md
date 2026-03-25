# Experiment 23: Full Persona Loop — Interview → Design → Review → Iterate

## The Idea
Before designing screens, ask simulated users what they need. Then show them the screens and let them push back. Iterate until they're happy.

## Brief
"Build an invoice generator for freelancers"

## The Loop

```
Brief → Personas → Interviews → MVP → Acceptance → Screens → Review → (Iterate)
```

## Results

| Step | Output | Cost | Key Finding |
|------|--------|------|-------------|
| Personas | 4 diverse users | $0.003 | Freelancer, agency owner, accountant, startup founder |
| Interviews | Needs, frustrations, budget | $0.009 | Recurring invoices was #1 request (not in brief!) |
| MVP Synthesis | Feature priority matrix | $0.005 | 6 must-have, 4 should-have, 3 nice-to-have |
| MVP Acceptance | 2 accept, 2 reject | $0.006 | Agency owner + accountant rejected — needed integrations |
| Screens + Wireframes | Full wireframes | $0.018 | 8+ screens with layouts, forms, empty states |
| Screen Review | **APPROVED (v1)** | $0.010 | Personas happy with screens on first pass |

**Total cost: $0.051**
**Iterations needed: 1**

## Key Findings

### 1. Personas Found Features Not in the Brief
The brief said "invoice generator." The personas asked for:
- **Recurring invoices** — mentioned by 3/4 personas
- **Payment tracking** — "did they pay yet?"
- **QuickBooks integration** — accountant's dealbreaker
- **Multi-currency** — agency owner works with international clients

Without persona interviews, the blueprint would have built a simple PDF generator. With them, it builds an invoicing platform.

### 2. The Rejection Was Valuable
2/4 personas rejected the MVP. Their feedback:
- Agency owner: "No recurring invoices? That's 80% of my billing"
- Accountant: "No QuickBooks export? I can't recommend this to clients"

This feedback shaped the MVP scope before any code was written. Fixing this in the spec costs $0. Fixing it after development costs a sprint.

### 3. Screen Review Passed First Try
Once the MVP incorporated persona feedback, the screens were approved on v1. The interview stage front-loads the disagreements — by the time screens exist, the personas already agreed on what to build.

## Pipeline Stage Recommendation

Add this as a new stage between brainstorm and design:

```
1. Brainstorm → decide what to build
1.5. Research → market/technical validation
NEW: Persona Discovery → who uses it, what they need, acceptance gate
2.5. Design → journeys, screens, wireframes (from persona-validated requirements)
2. Blueprint → spec with exact types (from persona-approved screens)
3. Development → code
```

## Files
- `01-personas.md` — 4 personas with backgrounds, pain points, budgets
- `02-interviews.md` — Full interview transcripts (in character)
- `03-mvp.md` — Feature priority matrix, MVP scope, pricing
- `04-mvp-acceptance.md` — 2 accept, 2 reject with reasons
- `05-screens-v1.md` — Full wireframes for all screens
- `06-screen-review-v1.md` — Persona reviews: APPROVED
