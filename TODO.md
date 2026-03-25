# Research TODO

## Completed (49 experiments)
- [x] Exp 1-18: Code generation (prompts, models, auto-fix, V-Model, parser)
- [x] Exp 19-21: Multi-service + different apps
- [x] Exp 22-23: Product design (journeys, personas, wireframes)
- [x] Exp 24-26: Wireframes→code, full pipeline, test layers
- [x] Exp 27-28: Fully automated pipeline
- [x] Exp 29: gqlgen GraphQL
- [x] Exp 30-32: Reviewer panels (4, 8, domain expert)
- [x] Exp 33: Add feature (CSV export, zero regression)
- [x] Exp 34-36: Simplicity agent, Playwright, domain+simplicity
- [x] Exp 37: 10-reviewer panel with SaaS UX
- [x] Exp 38: Progressive enhancement (ZERO regressions)
- [x] Exp 39: 20 reviewers (NO diminishing returns)
- [x] Exp 40-40b: Playwright (AI-driven failed, scripted works)
- [x] Exp 41-41d: AI testing agent, adversarial, pen test
- [x] Exp 42: AI pen test agent (4 High)
- [x] Exp 43: Static security (gosec, govulncheck)
- [x] Exp 44: Frontend + API security (27/37)
- [x] Exp 45: Chaos agent (25/25 survived)
- [x] Exp 46: Screenshot → product (Instatus clone)
- [x] Exp 47: Code metrics (57.4% coverage, complexity max 9)
- [x] Exp 48: TDD (90.3% coverage, 33% higher than code-first)
- [x] Exp 49: Mutation testing (89% score)

## Next Experiments

### Exp 50: GDPR/Privacy Reviewer
- [ ] Pre-code reviewer checks: right to deletion, export, consent, minimisation, retention
- [ ] Post-code reviewer checks: PII in logs, encrypted storage, PII in URLs, error message leaks
- [ ] Privacy pen tester: enumerate IDs, access other users' data, PII in HTML source

### Exp 51: Multi-User Concurrent Testing (Two Browsers)
- [ ] Two Playwright browsers running simultaneously
- [ ] User A adds client while User B deletes a client
- [ ] User A views invoice while User B marks it paid
- [ ] User A searches while User B creates 10 clients rapidly
- [ ] Verify: no crashes, data consistency, no stale reads

### Exp 52: SaaS Multi-Tenant Reviewer
- [ ] Data isolation: org_id on every query, user A can't see user B's data
- [ ] Auth on every route: no unauthenticated endpoints
- [ ] Billing hooks: where does Stripe go?
- [ ] Team features: invite, roles (admin/member/viewer)
- [ ] Rate limiting per tenant
- [ ] Subdomain/custom domain support

### Exp 53: Design System Pipeline (new spike)
- [ ] Screenshot/URL → multimodal AI reads the visual
- [ ] Extracts: colors, fonts, spacing, border-radius, shadows, layout grid
- [ ] Generates: design-tokens.json (CSS variables)
- [ ] Generates: component library (button, card, form, nav, table)
- [ ] Feeds into: every product uses these tokens
- [ ] Storybook or HTML catalog for visual reference

### Exp 54: Website Clone Pipeline (new spike)
- [ ] Multimodal reads live site (screenshots of pages)
- [ ] Extracts: content, structure, navigation, forms
- [ ] Produces: matching clone with same layout and content
- [ ] Feeds design tokens from Exp 53

### Exp 55: Docker + Deploy
- [ ] Take a passing app → generate Dockerfile
- [ ] Build Docker image
- [ ] Deploy to VPS
- [ ] Verify running at deployed URL

### Exp 56: User Feedback Loop
- [ ] Personas USE the app (via Playwright)
- [ ] Each reports what works, what's broken, what's missing
- [ ] Pipeline updates the app based on feedback
- [ ] Re-test → repeat 3 times
- [ ] Measure: does the app improve with each iteration?

### Exp 57: Multi-Language (TypeScript/Node.js)
- [ ] Same pipeline but Express/Fastify, Vitest, TypeScript
- [ ] Measure: does the approach generalise beyond Go?

### Exp 58: Real Database (PostgreSQL)
- [ ] Replace in-memory with PostgreSQL
- [ ] Generate migrations, connection pooling
- [ ] Measure: how much more complex?

### Exp 59: AI Reviews AI Code (Model Reviews Model)
- [ ] One model writes code
- [ ] Different model reviews (security, quality, patterns)
- [ ] Reviewer can REQUEST CHANGES → code model fixes
- [ ] Like pr-review.sh but for AI-generated code

### Exp 60: Creative Agent (Zara)
- [ ] After MVP is built, creative agent suggests "delight" features
- [ ] "What if invoices had a thank-you note?"
- [ ] "What if the dashboard showed revenue this month?"
- [ ] Measure: does the creative agent improve market reviewer's verdict?

### Exp 61: SaaS Auth (Registration, Login, SSO, MFA)
- [ ] Secure registration: email verification, Argon2id password hashing
- [ ] Login: JWT in httpOnly cookie, session management
- [ ] Forgot password: secure token, email reset flow
- [ ] SSO: OAuth2/OIDC with Google/GitHub
- [ ] MFA/2FA: TOTP (Google Authenticator), backup codes
- [ ] Account lockout after 5 failed attempts
- [ ] Password complexity rules
- [ ] Reviewer: auth security specialist
- [ ] Playwright: test entire auth flow (register → verify → login → 2FA)

### Exp 62: Scripted Regression Suite
- [ ] Collect ALL Playwright tests from every experiment
- [ ] Run full suite after every build (journeys + mobile + console + auth)
- [ ] Each new feature adds its own Playwright tests
- [ ] Gate: build fails if any regression test fails
- [ ] AI agent runs AFTER scripted tests to find NEW bugs

## NEW SPIKE: Documentation Orchestration

### Exp 63: Auto-Generate Project Documentation
- [ ] README.md from code (API endpoints, how to run, env vars)
- [ ] API docs from Go handler signatures (Swagger/OpenAPI style)
- [ ] User guide from wireframes + running app (step-by-step with descriptions)
- [ ] Architecture Decision Records from reviewer feedback
- [ ] Changelog from progressive enhancement iterations

### Exp 64: Review Stage Evidence Trail
- [ ] Each reviewer produces a sign-off document (verdict, items checked, issues found)
- [ ] Checklist format: YES/NO per item with reviewer name + timestamp
- [ ] Aggregate report: "8/10 reviewers APPROVE, 2 REQUEST CHANGES"
- [ ] Issue tracker: every finding logged with severity, owner, resolution status
- [ ] Final gate report: all issues resolved or accepted with justification

### Exp 65: Test Report Generation
- [ ] Test summary: pass/fail count, coverage %, mutation score
- [ ] Per-test-type breakdown: unit, integration, Playwright, adversarial, security
- [ ] Trend: coverage over progressive iterations (iter 1: 40%, iter 5: 90%)
- [ ] Failure analysis: what failed, why, how it was fixed
- [ ] Auto-generated from go test + gosec + Playwright output

### Exp 66: Security Audit Report
- [ ] Consolidated report: gosec + govulncheck + pen test + adversarial + OWASP
- [ ] Per-finding: severity, description, affected code, remediation, status
- [ ] Executive summary: Critical/High/Medium/Low counts
- [ ] Compliance checklist: OWASP Top 10 coverage (checked/not checked)
- [ ] GDPR data map: what PII exists, where, who, retention

### Exp 67: User Documentation
- [ ] Getting started guide (register, first client, first invoice)
- [ ] Feature-by-feature guide with text wireframes as illustrations
- [ ] FAQ generated from persona pain points
- [ ] API reference for integrations
- [ ] Admin guide (deployment, backup, monitoring)
