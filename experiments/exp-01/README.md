# Experiment 1: Escalation (Cheap Model → Stronger Model on Failure)

## Hypothesis
When Qwen3-30B fails to compile, sending the error + code to Gemini Flash should fix it.

## Setup
- Cheap model: Qwen3-30B ($0.0005/call)
- Strong model: Gemini Flash ($0.008/call)
- Task: Build main.go with embedded HTML for task-board
- 5 runs

## Result: 0/5 — Escalation does NOT fix the backtick issue

Both models make the same mistake: using JavaScript template literals inside Go backtick strings. This is a shared training data blind spot, not a model quality gap.

| Run | Cheap | Escalation | Cost |
|-----|-------|-----------|------|
| 0 | FAIL | FAIL | $0.027 |
| 1 | FAIL | FAIL | $0.026 |
| 2 | FAIL | FAIL | $0.031 |
| 3 | FAIL | FAIL | $0.026 |
| 4 | FAIL | FAIL | $0.024 |

## Key Finding
Escalation works for knowledge gaps (wrong type, missing import). It does NOT work for shared blind spots (language-specific gotchas all models get wrong).

## When Escalation Works
- Unused variable errors
- Missing imports
- Type mismatches

## When Escalation Doesn't Work
- Backtick-in-backtick (Go/JS)
- `&constant` in Go (pointer to constant)
- Language interaction gotchas shared across training data
