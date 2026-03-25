# Experiment 10: V-Model Pattern

## Hypothesis
Blueprint produces spec + hidden acceptance tests. Executor builds from spec only. Acceptance tests run after build as surprise verification.

## Setup
- Blueprint (Gemini Flash): produces spec.md + test-data.json + acceptance_test.go
- Executor (MiniMax): builds from spec only (never sees acceptance tests)
- Unit tests run during build
- Acceptance tests run after complete build

## Result
- Blueprint: $0.014 — created all 3 artifacts ✅
- Executor: $0.018 — built code from spec
- Unit tests: compile issues (parser/type)
- Acceptance tests: could not run (same compile issues)

## Key Finding
The pattern works conceptually — acceptance tests are a valid "surprise" verification gate. With parser improvements, this loop would complete.

```
Blueprint → spec + acceptance criteria (hidden)
    ↓
Executor → code + unit tests (from spec only)
    ↓
Unit gate → pass?
    ↓
Acceptance gate → surprise verification
    ↓
FAIL → back to executor with feedback
```
