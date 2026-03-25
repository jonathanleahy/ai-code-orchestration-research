# AI Code Orchestration Playbook

**How to go from a one-line idea to a running, tested, reviewed product.**

Based on 44 experiments. Total research cost: ~$12.

**Key numbers:** $0.34-0.96/app | 20 reviewers at $0.047 (no diminishing returns) | Progressive enhancement: ZERO regressions | Playwright: 20s, catches what 26 tests + 10 reviewers miss | gosec finds 12 issues FREE | AI pen test: 4 High for $0.02

---

## The Pipeline

```
1. BRIEF                    → One-line idea
2. PERSONA DISCOVERY         → 3 personas, interviews, MVP ($0.01)
3. PRE-CODE REVIEW (7)       → Dev, Product, QA, Market, Domain, SaaS UX, Simplicity ($0.035)
4. REVISE                    → Incorporate all feedback ($0.005)
5. GO TYPES                  → Exact type signatures ($0.007)
6. BUILD (progressive)       → One feature at a time, verify after each ($0.10-0.20 per iteration)
7. POST-CODE REVIEW (3)      → Code Architecture, Accessibility, OWASP ($0.015)
8. APPLY FIXES + REBUILD     → Fix issues, re-test ($0.05-0.20)
9. PLAYWRIGHT                → Personas click through the app (essential)
10. SHIP
```

**Total cost: $0.40-0.96 per product.**

---

## Pre-Code Reviewers (7)

Run BEFORE any code is written. Review the spec/wireframes.

| # | Reviewer | Checks | Cost | Exp |
|---|---------|--------|------|-----|
| 1 | **Dev Architect** | Architecture, cost, $5 VPS viability, security | $0.005 | 25 |
| 2 | **Product Owner** | Every CRUD action has a button/form. Every persona journey has a matching screen. | $0.006 | 30 |
| 3 | **QA Engineer** | Empty forms, duplicates, delete with children, empty states, network errors, back navigation | $0.004 | 30 |
| 4 | **Market Analyst** | vs HubSpot/Zoho Free. Differentiation. Would someone pay $15/month? | $0.005 | 30 |
| 5 | **Domain Expert** | Workflow-specific. For invoicing: line items, print, send, pay, void, due date, overdue. Auto-generate from brief. | $0.006 | 32 |
| 6 | **SaaS UX Designer** | Breadcrumbs, toasts (not alerts), confirmations on delete, back buttons, empty states with CTA, responsive, search, consistent button colors | $0.005 | 37 |
| 7 | **Constrained Simplicity** | Cannot cut features. Can simplify HOW. Address=textarea not 5 fields. Print=browser print not PDF library. | $0.005 | 36 |

**Total: $0.035. Each finds unique issues.**

---

## Post-Code Reviewers (3)

Run AFTER code is built. Review actual HTML/Go code.

| # | Reviewer | Checks | Cost | Exp |
|---|---------|--------|------|-----|
| 8 | **Code Architecture** | Separation of concerns, DRY, handler size <50 lines, error consistency | $0.005 | 37 |
| 9 | **Accessibility** | WCAG 2.1 AA, ARIA labels, form labels, keyboard nav, color contrast, heading hierarchy | $0.005 | 31 |
| 10 | **OWASP Security** | XSS (html.EscapeString), CSRF, CORS, input validation, no secrets in source | $0.005 | 31 |

**Total: $0.015. IMPORTANT: re-run ALL tests after applying fixes (Exp 32: CSP fix broke inline JS).**

---

## Build Method: TDD + Progressive Enhancement

### TDD (Tests First) — Exp 48
Write tests BEFORE code. Code is built to pass the tests.

| Metric | TDD (Exp 48) | Code-First (Exp 38) |
|--------|-------------|---------------------|
| Coverage | **90.3%** | 57.4% |
| Tests | 34/34 | 32/32 |
| Server lines | 491 | 894 |
| Cost | **$0.24** | $0.82 |

**TDD produces 33% higher coverage at 71% less cost.**

### Mutation Testing (Exp 49)
Intentionally change code, verify tests catch it. Score: **89%** (8/9 caught).
Only gap: no test checks timestamps (CreatedAt set to zero → tests still pass).

### Progressive Enhancement
Don't build everything at once. Add one feature at a time.

```
Iteration 1: Client CRUD (add, list, view)        → test → PASS
Iteration 2: + Edit + Delete                       → test → PASS (old tests too)
Iteration 3: + Activity log                        → test → PASS (old tests too)
Iteration 4: + Invoices                            → test → PASS (old tests too)
Iteration 5: + Search + print                      → test → PASS (old tests too)
```

**Why:** Exp 38 — 5 iterations, 32/32 tests, ZERO regressions. Exp 32 — one-shot with 10 features, store tests failed (too complex).

**Rule:** Each iteration must pass ALL existing tests before proceeding.

---

## Playwright Testing (Non-Negotiable)

Run after EVERY build. Simulates real users clicking through the app.

**Bugs only Playwright catches (Exp 32, 37):**
| Bug | API Tests | Reviewers | Playwright |
|-----|----------|-----------|-----------|
| CSP blocks inline JS | PASS | Approved | CAUGHT |
| Plain text error → JSON parse fail | PASS | Approved | CAUGHT |
| ParseForm ignores multipart FormData | PASS | Approved | CAUGHT |
| Modal overlay blocks all clicks | PASS | Approved | CAUGHT |

**4 bugs, all with passing tests and approved reviews. Only clicking found them.**

---

## Go Code Rules (from experiments)

| Rule | Why | Exp |
|------|-----|-----|
| Exact type signatures in spec | Prevents V-Model type mismatches (0% → 93%) | 17 |
| `r.ParseMultipartForm` not `r.ParseForm` | FormData sends multipart, ParseForm ignores it | 32 |
| Return JSON errors, not `http.Error` plain text | JS `.json()` fails on plain text | 32 |
| No `Content-Security-Policy` with inline JS | CSP `script-src 'self'` blocks embedded `<script>` | 32 |
| String concatenation in JS, not template literals | Go backticks can't contain backticks | 15, 18 |
| `goimports -w .` before testing | Fixes 40% of model errors for free | 4, 16 |
| `fix-address-of-const.py` after goimports | Fixes `&constant` errors (shared blind spot) | 16 |
| One store per test function, no shared state | Prevents test order dependencies | 38 |
| Address as single string, not struct | Simpler, fewer fields, less to break | 36 |
| Invoice status as string, not state machine | "draft"/"sent"/"paid"/"void" — simple and correct | 36 |

---

## Security Testing Stack (Exp 41-44)

### Static Analysis (after code gen, FREE, instant)
```bash
cd app/
gosec ./...                    # XSS taint, unbounded parsing, hardcoded creds
staticcheck ./...              # Code correctness, performance
govulncheck ./...              # Known CVEs in dependencies
go vet ./...                   # Suspicious constructs
```

| Tool | Findings (Exp 43) | Catches |
|------|-------------------|---------|
| gosec | 12 issues | XSS via taint (3), unbounded form parsing (3), unhandled errors (4), no timeout (1) |
| staticcheck | 0 | Code correctness |
| govulncheck | 13 CVEs | Go stdlib vulnerabilities (need upgrade) |
| go vet | 0 | Printf format errors, unreachable code |

### Dynamic Analysis (against live server)
| Tool | Findings | Cost | Catches |
|------|----------|------|---------|
| Frontend HTML checker (44) | 1 fail | FREE | Missing lang, inline handlers, eval() |
| Security headers (44) | 6 missing | FREE | X-Frame-Options, CSP, HSTS, Referrer-Policy |
| REST API checker (44) | 3 fails | FREE | No pagination, rate limiting, Cache-Control |
| Stored XSS (44) | SAFE | FREE | Unescaped user input in HTML |
| Deep adversarial (41d) | 11 issues | FREE | XSS, SQLi, path traversal, unicode, methods |
| AI pen test agent (42) | 4 High | $0.02 | AI-planned attacks, chained findings |
| AI testing agent (41c) | 0 bugs | $0.035 | Exploratory — agent decides what to test |

### Browser Testing (Playwright)
| Test Suite | Findings | Time |
|------------|----------|------|
| User journeys (40b) | 3 fails | 20s |
| Mobile viewport (40b) | 2 fails (table overflow, button off-screen) | 5s |
| Console errors (40b) | 1 fail | 5s |

---

## Cost Breakdown

| Component | Cost | % of Total |
|-----------|------|-----------|
| Pre-code reviews (7-20) | $0.035-0.10 | 5% |
| Post-code reviews (3) | $0.015 | 2% |
| Store code (Qwen3-30B) | $0.01-0.02 | 2% |
| HTTP server (claude -p) | $0.15-0.25 | 35% |
| HTTP tests (claude -p) | $0.10-0.20 | 25% |
| Security testing | $0.00-0.05 | 5% |
| Playwright | FREE | 0% |
| Static analysis | FREE | 0% |
| **Total** | **$0.40-0.96** | |

**$0.05 of reviews + FREE security testing saves $5.00 of post-launch fixes.**

---

## Tools Created

| Tool | Purpose | Location |
|------|---------|----------|
| parse-blocks-v2.py | Hardened file parser (8 formats) | scripts/dev-spike-v3/ |
| fix-address-of-const.py | Go &constant auto-fix | scripts/dev-spike-v3/ |
| exp27-fully-automated.py | Full pipeline script | scripts/dev-spike-v3/ |
| exp30-full-reviewed.py | 4-reviewer pipeline | scripts/dev-spike-v3/ |
| exp37-full-panel.py | 10-reviewer pipeline | scripts/dev-spike-v3/ |
| exp38-progressive.py | Progressive enhancement | scripts/dev-spike-v3/ |

---

## 20 Reviewers — No Diminishing Returns (Exp 39)

At $0.005/reviewer, 20 reviewers cost $0.047. First 10 found 80 issues, last 10 found 82.
**Every reviewer contributes. There are no diminishing returns at 20.**

Extra reviewers beyond the core 10:
- Freelance Designer ("does this solve MY problem?")
- Agency Owner ("does this scale to 200 clients?")
- First-Time User ("what do I do first?")
- Power User ("500 clients, 50 invoices/month")
- Mobile User ("can I add a quick note on my phone?")
- Competitor User ("why switch from HubSpot?")
- Data Privacy Officer ("GDPR? Data export?")
- Billing Specialist ("partial payments? Multi-currency?")
- Onboarding Specialist ("CSV import? Getting-started wizard?")
- Growth Hacker ("viral loop? Free tier hook?")

---

## Adversarial Testing (Exp 41b)

20 adversarial tests: XSS, SQL injection, empty inputs, 1000-char names, unicode, negative amounts.
**Zero bugs found — the app handles all adversarial inputs correctly.** Cost: $0.003.

---

## Testing Strategy

| Layer | When | Cost | Catches |
|-------|------|------|---------|
| Go unit tests | After store build | $0.01-0.02 | Logic bugs |
| HTTP integration tests | After server build | $0.10-0.20 | API contract bugs |
| Scripted Playwright | After every build | FREE (20s) | UI regressions |
| Adversarial testing | Before release | $0.003 | Security holes |
| Manual Playwright | On demand | — | Complex interaction bugs |

---

## For Dark Factory Integration

1. **Blueprint stage** → add exact Go types (already done: IR-EXP17)
2. **New stage: Persona Discovery** → between brainstorm and blueprint
3. **New stage: Pre-code Review** → 7-20 reviewers before development ($0.035-0.10)
4. **Development stage** → progressive enhancement (one feature at a time)
5. **New stage: Post-code Review** → 3 reviewers after development
6. **Testing stage** → Playwright scripts + adversarial tests
7. **Auto-fix loop** → goimports + gofmt + fix-address-of-const after every code generation
8. **SaaS UX patterns** → breadcrumbs, toasts, confirmations in every server build prompt
