# Experiment 17: V-Model Full Loop

## Hypothesis
Blueprint produces spec + hidden acceptance tests. Executor builds from spec only. Acceptance tests run after build as surprise verification.

## Setup
- Model: Qwen3-30B for both blueprint and executor
- Blueprint outputs: spec.md + acceptance_test.go (hidden from executor)
- Executor builds: model/task.go from spec only
- Acceptance gate: runs hidden tests after build
- Feedback loop: up to 3 iterations
- 3 runs

## Results

| Run | Blueprint | Spec | Acceptance Tests | Executor | Loops | Cost |
|-----|-----------|------|-----------------|----------|-------|------|
| 0 | OK | OK | MISSING | PASS (compiles) | 1 | $0.011 |
| 1 | OK | OK | MISSING | PASS (compiles) | 1 | $0.013 |
| 2 | OK | OK | MISSING | PASS (compiles) | 1 | $0.012 |

**Pass rate: 3/3 (100%) — but without acceptance test verification.**

## Key Finding

The V-Model pattern works for spec → build, but the **blueprint fails to generate parseable acceptance tests**. The model outputs the test code but the parser doesn't extract `acceptance_test.go` from the response (it gets mixed into explanation text).

### What Works
- Blueprint generates high-quality specs from requirements
- Executor builds compiling code from those specs on first attempt
- The spec → build pipeline is reliable at $0.012/run

### What Needs Work
- Acceptance test extraction (parser needs to handle `.go` files outside standard locations)
- The feedback loop never activates because the compile check passes without acceptance tests
- Need explicit "output acceptance_test.go as a separate file block" in the blueprint prompt

## Cost

| Item | Cost |
|------|------|
| 3 runs × (blueprint + executor) | $0.037 |
| **Total** | **$0.037** |
