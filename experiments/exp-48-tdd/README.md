# Experiment 48: TDD — Tests First

## Write tests before code, then build code to pass them.

### Process
1. Generate tests from spec (NO code exists)
2. Verify tests don't compile (expected)
3. Generate code to pass the tests
4. Auto-fix loop until tests pass
5. Build HTTP server + tests

### Results
| Metric | TDD (Exp 48) | Code-First (Exp 38) |
|--------|-------------|---------------------|
| Store tests | 0/0 | N/A (via HTTP) |
| HTTP tests | 34/34 | 32/32 |
| Coverage | **90.3%** | 57.4% |
| Server lines | 491 | 894 |
| Cost | $0.2359 | $0.82 |

### Key Finding
TDD coverage: **90.3%** vs code-first: 57.4%.
Tests written first guarantee the code is testable.
