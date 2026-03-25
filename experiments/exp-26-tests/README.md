# Experiment 26: Add Test Layers to Invoice Generator

## The Three Test Layers

| Layer | Generator | Tests | Pass | Cost |
|-------|-----------|-------|------|------|
| Store unit tests | Qwen3-30B | 17 | **16/17** | $0.013 |
| Acceptance tests (V-Model) | Qwen3-30B | 6 | **4/6** | $0.009 |
| HTTP integration tests | claude -p Haiku | 13 | **13/13** | $0.096 |
| **Total** | | **36** | **33/36 (92%)** | **$0.118** |

## Results

### Store Unit Tests (Qwen3-30B, $0.013)
16/17 pass. One test failure: `TestClientStore/ListClients` expects 4 clients but store returns 3 (off-by-one in test setup, not a store bug).

Had one compile error (`err :=` should be `err =` — redeclared variable). The auto-fix didn't catch this; added to the fix list.

### Acceptance Tests — V-Model (Qwen3-30B, $0.009)
4/6 pass. The V-Model caught 2 real interface mismatches:
- `TestAcceptCreateInvoice`: Test calls `Create()` with different arguments than the store expects
- `TestAcceptCreateInvoiceMissingClient`: Same signature mismatch

This is exactly what V-Model acceptance tests are designed to catch — the spec's type signatures didn't perfectly match the implementation.

### HTTP Integration Tests (claude -p, $0.096)
**13/13 pass.** claude -p wrote comprehensive httptest-based tests covering:
- CRUD operations on /api/invoices
- 404 for nonexistent IDs
- Health check
- HTML dashboard

## Key Findings

1. **claude -p writes perfect HTTP tests** — 13/13 because it reads the actual main.go and tests against it
2. **V-Model acceptance tests catch real mismatches** — 2/6 failures were genuine interface discrepancies
3. **Cheap model tests need one manual fix** — redeclared variable is a pattern the auto-fix should handle
4. **92% pass rate across 36 tests** for $0.118 additional cost

## Cost Impact on Full Pipeline

| Pipeline Step | Without Tests | With Tests |
|---------------|--------------|-----------|
| Idea → Code (Exp 25) | $0.21 | $0.21 |
| Test layers (Exp 26) | — | $0.12 |
| **Total** | **$0.21** | **$0.33** |

Adding tests increases cost by 56% but provides 36 automated tests (33 passing).
