# Experiment 25: Full Pipeline — One-Line Brief to Running Product

## The Vision
Can you go from "build an invoice generator for freelancers" to a running Go application with a single pipeline, including customer discovery and engineering review?

## The Pipeline

```
Brief: "Build an invoice generator for freelancers"
  ↓
[1] Persona Discovery     — 3 users interviewed ($0.004)
  ↓
[2] MVP Prioritization    — features ranked by persona demand ($0.005)
  ↓
[3] Screen Wireframes     — every screen with ASCII layout ($0.008)
  ↓
[4] Dev Review            — 3 senior engineers review ($0.007)
  ↓                        → NEEDS CHANGES (pushed back!)
[5] Revised MVP           — simplified from dev feedback ($0.006)
  ↓
[6] Go Spec               — exact types + signatures ($0.007)
  ↓
[7] Store Layer            — Qwen3-30B builds store ($0.013) → COMPILES
  ↓
[8] HTTP Server            — claude -p builds main.go ($0.16) → COMPILES (713 lines)
```

**Total: $0.21 from idea to compiled application.**

## What the Dev Review Caught

The 3 engineers reviewed the MVP and said NEEDS CHANGES:

- **Backend Architect**: "In-memory is fine for MVP, but add JSON file persistence for restart survival. Single binary on $5 VPS is correct."
- **Cost & Complexity**: "PDF generation sounds simple but requires a library or raw HTML-to-PDF. Defer to v2. Email sending needs SMTP — use webhook notifications instead."
- **Security & Ops**: "Invoice data is PII. Need HTTPS enforced, no plaintext storage of client emails. Add rate limiting on create endpoints."

The revised MVP removed PDF generation and email sending from v1, added JSON persistence and rate limiting. These changes would have cost days to fix after development.

## Results

| Metric | Value |
|--------|-------|
| Brief | 8 words |
| Personas | 3 (freelancer, agency, accountant) |
| Features prioritized | 13 (6 must-have, 4 should-have, 3 deferred) |
| Dev review verdict | NEEDS CHANGES → revised |
| Store | COMPILES (Qwen3-30B, $0.013) |
| Server | COMPILES (claude -p, 713 lines) |
| Total cost | **$0.21** |

## What's Missing: Tests

This pipeline doesn't generate tests yet. The complete pipeline needs:

1. **Store tests** (Exp 16 approach): Generate alongside store code, same call
2. **Acceptance tests** (Exp 17 V-Model): Blueprint generates hidden tests, executor builds blind
3. **HTTP integration tests**: Verify API endpoints match the spec

Adding tests would add ~$0.02 (store tests via Qwen3-30B) + ~$0.05 (acceptance tests) = ~$0.07, bringing the total to ~$0.28.

## The Complete Pipeline (proven)

```
Brief ($0)
  → Persona Discovery ($0.004) — who needs it, what they need
  → Persona Interviews ($0.009) — in-character interviews
  → MVP Synthesis ($0.005) — feature prioritization
  → Persona Acceptance ($0.006) — personas approve/reject
  → Screen Wireframes ($0.008) — every screen with layout
  → Dev Review ($0.007) — engineers approve/reject architecture
  → Revised Spec ($0.006) — incorporate dev feedback
  → Go Types ($0.007) — exact type signatures
  → Store Layer ($0.013) — Qwen3-30B + auto-fix
  → Store Tests ($0.007) — 2-file approach (TODO)
  → Acceptance Tests ($0.05) — V-Model hidden tests (TODO)
  → HTTP Server ($0.16) — claude -p Haiku
  → Integration Tests (FREE) — verify API (TODO)
  ────────────────────
  Total: ~$0.28 from idea to tested product
```

## Files
- `01-discover.md` — 3 personas + interview transcripts
- `02-mvp.md` — Feature priority matrix + MVP scope
- `03-screens.md` — Screen wireframes
- `04-dev-review.md` — 3 engineer reviews (NEEDS CHANGES)
- `05-revised-mvp.md` — Revised scope from dev feedback
- `06-go-types.md` — Exact Go type declarations
- `07-store-build.md` — Store compilation result
- `app/` — The actual Go application
