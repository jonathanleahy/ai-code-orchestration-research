# Experiment 40b: Scripted Playwright Tests

## Approach
Write actual Playwright test scripts (TypeScript), run with `npx playwright test`.
Deterministic, fast (20s), catches regressions.

## Results: 7/10 pass, 3 fail

| # | Journey | Result | Issue |
|---|---------|--------|-------|
| 1 | View Dashboard | PASS | |
| 2 | Add Client | FAIL | Uses alert() not toast — test can't verify |
| 3 | View Client Profile | FAIL | Test bug: link text "View" ≠ h1 text |
| 4 | Client has tabs | PASS | |
| 5 | Edit Client | PASS | |
| 6 | Delete button exists | PASS | |
| 7 | Activity Tab | PASS | |
| 8 | Invoices Tab | PASS | |
| 9 | Search | PASS | |
| 10 | Empty State | **FAIL — BUG** | Search for nonexistent shows 500 in page |

## Key Findings

1. **Scripted tests work perfectly** — 20 seconds, deterministic, real browser
2. **Found a real bug**: empty search state triggers error display in the page
3. **2 test bugs**: Add Client uses alert() (test needs dialog handler), View Profile link text mismatch
4. **vs Exp 40 (AI-driven)**: scripted is reliable, AI-driven was all UNKNOWN

## Two Approaches (both valuable)

| Approach | When | Cost | Catches |
|----------|------|------|---------|
| **Static scripted** | After every build | Free (20s) | Regressions |
| **AI testing agent** | After every release | ~$0.10 | Exploratory bugs |
