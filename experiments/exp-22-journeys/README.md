# Experiment 22: One-Line Brief → User Journeys → Screen Wireframes

## The Problem

Exp 21 built a 4-service status page that compiles and runs. But the UI was guessed — claude -p made a single dark-themed HTML page with no admin interface, no forms, no navigation. A developer tool, not a product.

**The gap: we can build code cheaply, but what code should we build?**

## The Experiment

From a one-line brief, generate the full product design in 4 steps:

```
"Build a public status page for monitoring website uptime,
 managing incidents, and notifying subscribers."
         ↓
Step 1: User Personas + Journeys     ($0.004, 5K chars)
         ↓
Step 2: Screen Map (every route)     ($0.009, 6.6K chars)
         ↓
Step 3: Text Wireframes (ASCII)      ($0.019, 21K chars)
         ↓
Step 4: Compare with guessed UI      ($0.006, 4K chars)
```

**Total cost: $0.038 on the cheapest model.**

## What It Generated

### Personas (4)
- Site Reliability Engineer (daily monitoring, incident response)
- Product Manager (stakeholder communication)
- Customer Support Agent (customer-facing status info)
- External Subscriber (self-service status checks)

### Journeys (5)
- Initial setup and configuration
- Daily monitoring routine
- Incident response workflow
- Subscriber notification experience
- Incident resolution and postmortem

### Screens (8+)
- S01: Dashboard (service cards, status overview)
- S02: Add Check (form with URL, name, type)
- S03: Check Detail (timeline, incidents, activity)
- S04: Incidents List (filter by status, severity)
- S05: New Incident (creation form)
- S06: Incident Detail (timeline, updates, resolution)
- S07: Subscribers (manage webhook endpoints)
- S08: Public Status Page (external-facing)

### Wireframes
Full ASCII wireframes for every screen with:
- Layout structure (header, sidebar, content)
- Interactive elements ([Add Check], [Create Incident])
- Sample data (realistic service names, timestamps)
- Empty states ("No checks configured yet")
- Responsive notes (mobile behaviour)

## The Comparison

| Dimension | Guessed (Exp 21) | Journey-Driven (Exp 22) |
|-----------|-----------------|------------------------|
| Screens | 1 | 8+ |
| Admin UI | None | Full CRUD forms |
| Navigation | None | Complete flow |
| Empty states | None | Defined |
| Mobile | None | Responsive notes |
| Would pay $20/mo | No | Yes |
| Code delta | ~50 lines HTML | ~200 lines HTML |

**The journey-driven design produces a product. The guessed design produces a demo.**

## Key Finding

**$0.038 of design work determines whether the $0.021 of code produces a product or a toy.**

The cheapest model can generate useful user journeys, screen maps, and wireframes from a one-line brief. These artifacts feed directly into the blueprint spec, giving the development stage concrete screens to build instead of guessing.

## Pipeline Integration

```
Brief (one line)
  → Exp 22: Journeys + Screens + Wireframes ($0.04)
  → Blueprint: Spec with exact types + wireframes
  → Development: Builds from wireframes (knows every screen)
  → Testing: Verifies every journey path
  → Result: Product, not demo
```

Total cost to go from idea to complete design: **$0.04**
Total cost to build the code: **$0.02**
**Total cost from idea to running product: ~$0.06**

## Files

- `01-user-journeys.md` — 4 personas, 5 journeys with edge cases
- `02-screen-map.md` — 8+ screens with components, data, navigation
- `03-wireframes.md` — Full ASCII wireframes for every screen
- `04-comparison.md` — Side-by-side comparison with Exp 21
