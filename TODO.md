# Research TODO

## Completed (29 experiments)
- [x] Exp 1-18: Code generation (prompts, models, auto-fix, V-Model, parser)
- [x] Exp 19-21: Multi-service + different apps (StatusPulse, URL shortener)
- [x] Exp 22-23: Product design (journeys, personas, wireframes)
- [x] Exp 24-26: Wireframes→code, full pipeline, test layers
- [x] Exp 27-28: Fully automated pipeline (bookmark manager, CRM)
- [x] Exp 29: gqlgen GraphQL (iterative pipeline for code generation)
- [x] Exp 30: 4-reviewer panel (dev, product, QA, market) — RUNNING

## Next Experiments

### Exp 31: Personas USE the Running App (Playwright)
- [ ] Build app from brief
- [ ] Start the server
- [ ] Use Playwright to simulate each persona doing their journey
- [ ] Persona 1: "I tried to add a client and the button didn't work"
- [ ] Persona 2: "I searched for 'Acme' and got no results even though I just added it"
- [ ] Report what works and what doesn't — test the PRODUCT, not the code
- [ ] Measures: journey completion rate per persona

### Exp 32: Regression Testing (Change a Feature)
- [ ] Take a working app with passing tests
- [ ] Change one thing: rename a field, change a default, alter a validation
- [ ] Run the test suite — do the tests catch the change?
- [ ] Measures: what % of intentional changes are caught by tests
- [ ] Tests test quality, not just test count

### Exp 33: Add Feature to Existing App
- [ ] Take the CRM from Exp 28/30
- [ ] Brief: "Add CSV export for client list"
- [ ] Tests whether the pipeline handles ENHANCEMENT, not just greenfield
- [ ] Must not break existing tests
- [ ] Measures: feature added + all old tests still pass

### Exp 34: Docker Build + Deploy
- [ ] Take a passing app
- [ ] Generate Dockerfile
- [ ] Build Docker image
- [ ] Deploy to VPS (or local Docker)
- [ ] Verify it runs at the deployed URL
- [ ] Measures: time from passing tests to live URL

### Exp 35: Screenshot → Product
- [ ] Input: screenshot of a competitor (HubSpot, Instatus, etc.)
- [ ] Claude reads the screenshot
- [ ] Extracts: layout, features, navigation, style
- [ ] Generates wireframes + spec + code matching the screenshot
- [ ] Measures: how close is the result to the original?

### Exp 36: Minimum Viable Pipeline (How Cheap?)
- [ ] Strip everything: 1 persona, no wireframes, just types → code → tests
- [ ] What's the absolute floor cost for a working app?
- [ ] Compare: full pipeline ($0.50) vs minimal pipeline ($?)
- [ ] At what point does removing steps reduce quality?

### Exp 37: Multi-Language (TypeScript/Node.js)
- [ ] Same pipeline but target TypeScript instead of Go
- [ ] Uses: npm, Express/Fastify, Vitest
- [ ] Measures: does the approach generalize beyond Go?

### Exp 38: Real Database (PostgreSQL)
- [ ] Replace in-memory store with PostgreSQL
- [ ] Generate migrations, connection pooling, queries
- [ ] Measures: how much more complex is the pipeline for persistent storage?

### Exp 39: AI Code Review (Model Reviews Model)
- [ ] One model writes the code
- [ ] Different model reviews it (security, quality, patterns)
- [ ] Reviewer can REQUEST CHANGES → code model fixes
- [ ] Like pr-review.sh but for AI-generated code
- [ ] Measures: what issues does the reviewer catch?

### Exp 40: Progressive Enhancement (MVP → V2 → V3)
- [ ] Build MVP from brief
- [ ] Personas use it → feedback
- [ ] Add features based on feedback (not a new brief)
- [ ] Repeat 3 times
- [ ] Measures: does the app improve with each iteration? Do tests keep passing?
