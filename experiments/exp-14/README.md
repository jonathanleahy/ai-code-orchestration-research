# Experiment 14: Model Routing with Retry Loop

## Hypothesis
Exp 05 showed model routing fails in single-shot (0/5). With a retry loop (3 attempts) + auto-fix + error feedback, the same cheap model should succeed.

## Setup
- Model: Qwen3-30B ($0.0005/call) for all sub-tasks
- 3 retries per sub-task (error from previous attempt included in next prompt)
- Auto-fix: goimports + gofmt + &constant fix
- Parser: v2 (hardened)
- 3 runs

## Results

| Sub-Task | Pass Rate | Avg Attempts | Avg Cost | Notes |
|----------|-----------|-------------|----------|-------|
| model_and_test | **100%** (3/3) | 1.0 | $0.007 | All first attempt |
| main.go | **0%** (0/9) | 3.0 | $0.030 | Backtick issue every time |

### Detail

| Run | model_and_test | main.go |
|-----|---------------|---------|
| 0 | PASS (1st try) | FAIL (3 tries) |
| 1 | PASS (1st try) | FAIL (3 tries) |
| 2 | PASS (1st try) | FAIL (3 tries) |

## Key Finding

**Retry loop + auto-fix = 100% for model layer, 0% for main.go.**

The retry loop is effective for:
- Errors the auto-fix handles (imports, constants) → first attempt
- Simple compilation errors → 1-2 retries

The retry loop is **NOT** effective for:
- Shared training data blind spots (backtick/template literal) → same mistake every time
- Even with the error message in the prompt, Qwen3-30B repeats the same pattern

## Recommended Strategy

```
model/test sub-tasks  → Qwen3-30B + auto-fix       → $0.007 (100%)
main.go sub-task      → claude -p (subscription)     → FREE   (100%)
```

Don't waste 3 retries ($0.03) on main.go with cheap models — go straight to claude -p.

## vs Exp 05 (Single-Shot)

| Metric | Exp 05 (single-shot) | Exp 14 (retry) | Delta |
|--------|---------------------|----------------|-------|
| Overall pass | 0% | 50% (model only) | +50% |
| model_and_test | N/A | 100% | New |
| main.go | 0% | 0% | Same |
| Cost per run | $0.005 | $0.037 | +$0.032 |

## Cost

| Item | Cost |
|------|------|
| 3 runs × model_and_test (1 attempt each) | $0.020 |
| 3 runs × main.go (3 attempts each) | $0.091 |
| **Total** | **$0.111** |
