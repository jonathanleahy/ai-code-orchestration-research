# Experiment 17: V-Model Full Loop

## Hypothesis
Blueprint AI produces spec + hidden acceptance tests. Executor AI builds from spec only. Hidden tests run as surprise verification after build.

## Setup

### 17a (original)
- Single API call for spec + tests → tests not extracted (parser issue)
- Result: 3/3 PASS but compile-only (no acceptance verification)

### 17b (fixed extraction)
- Separate API calls for spec and acceptance tests (guaranteed extraction)
- Acceptance tests extracted in all 5 runs
- Feedback loop: failed test names sent back (not test code)
- 5 runs, 3 loops max

## Results (17b)

| Run | Spec | Tests Extracted | Build Loops | Acceptance | Cost |
|-----|------|----------------|-------------|-----------|------|
| 0 | OK | OK | 3 (all fail) | 0/0 (compile error) | $0.014 |
| 1 | OK | OK | 3 (all fail) | 0/0 (compile error) | $0.016 |
| 2 | OK | OK | 3 (all fail) | 0/0 (compile error) | $0.016 |
| 3 | OK | OK | 3 (all fail) | 0/0 (compile error) | $0.016 |
| 4 | OK | OK | 3 (all fail) | 0/0 (compile error) | $0.016 |

**0/5 — acceptance tests never compile against the generated code.**

## Root Cause

The acceptance tests and the executor produce **incompatible types**:

```
Acceptance tests expect:  StatusTodo = "TODO",  StatusDoing = "DOING"
Executor produces:        StatusTodo = "todo",  StatusInProgress = "in_progress"

Acceptance tests expect:  Update(id, title, desc *string, status *Status)
Executor produces:        Update(id string, updates map[string]interface{})
```

Two separate Qwen3-30B calls interpret the same spec differently. This is the fundamental V-Model challenge: **the spec must be precise enough that two independent implementations agree on types.**

## What The V-Model Caught

This is actually the V-Model **working correctly** — it caught a real integration failure that a compile gate alone would miss. The acceptance tests are a stricter gate than `go build`.

## Fix for Production

The spec must include **exact Go type signatures**, not just descriptions:

```markdown
## BAD (ambiguous)
Status type with constants: Todo, Doing, Done

## GOOD (exact)
type Status string
const (
    StatusTodo  Status = "TODO"
    StatusDoing Status = "DOING"
    StatusDone  Status = "DONE"
)
func (s *Store) Update(id string, title, description *string, status *Status) (*Task, error)
```

When the spec includes exact signatures, both the test writer and code writer produce compatible output.

## Revised V-Model Architecture

```
Blueprint (stronger model) → spec with EXACT type signatures + hidden tests
                                    ↓
Executor (cheap model)    → builds from spec (types are unambiguous)
                                    ↓
Auto-fix                  → goimports + gofmt + &constant
                                    ↓
Acceptance gate           → hidden tests (surprise verification)
                                    ↓
FAIL → feedback (test names only) → retry
PASS → done
```

The key insight: **spec precision determines V-Model success**. Vague specs → type mismatches → compile failures. Exact type signatures → both sides agree → acceptance tests pass.

## 17c: Exact Type Signatures (the fix)

Spec now includes copy-pasteable Go type declarations:
```go
type Status string
const (
    StatusTodo  Status = "TODO"
    StatusDoing Status = "DOING"
    StatusDone  Status = "DONE"
)
func (s *Store) Update(id string, title, description *string, status *Status) (*Task, error)
```

| Run | Acceptance Tests | Best Score | Loops | Cost |
|-----|-----------------|-----------|-------|------|
| 0 | Extracted | 14/15 | 3 (fail) | $0.019 |
| 1 | Extracted | 13/14 | 3 (fail) | $0.018 |
| 2 | Extracted | 10/11 | 3 (fail) | $0.017 |
| 3 | Extracted | 13/14 | 3 (fail) | $0.018 |
| 4 | Extracted | **11/11 PASS** | 1 | $0.010 |

**1/5 (20%) full pass — up from 0% in 17b.** Run 4: all 11 acceptance tests passed on first attempt.

### Progression: 17a → 17b → 17c

| Version | Spec Type | Tests Extracted | Acceptance Score | Pass Rate |
|---------|-----------|----------------|-----------------|-----------|
| 17a | Vague | No (0/3) | N/A (compile only) | 100% (trivial) |
| 17b | Vague | Yes (5/5) | 0/0 (type mismatch) | 0% |
| 17c | Exact types | Yes (5/5) | 10-15 of 11-15 | **20%** |

The exact types fixed the compile mismatch. The remaining 1-test failure is a behavioral edge case (usually List filter or Update) that the cheapest model sometimes gets wrong.

### Next Step: Stronger Blueprint Model

With a stronger model (Sonnet/Opus) writing the spec + acceptance tests, and the cheap model only implementing, the pass rate should increase further. The architecture is proven — it just needs better spec quality.

## Cost

| Item | Cost |
|------|------|
| 17a: 3 runs (compile only) | $0.037 |
| 17b: 5 runs (type mismatch) | $0.078 |
| 17c: 5 runs (exact types) | $0.082 |
| **Total** | **$0.197** |
