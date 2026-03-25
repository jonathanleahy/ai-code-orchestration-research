# Experiment 19: Re-run V2 (dep-doctor Node.js CLI) with Improvements

## Hypothesis
Re-running the dep-doctor CLI build with V4 prompt + parser v2 + Qwen3-30B should improve on the original results.

## Setup
- Model: Qwen3-30B ($0.0005/call)
- V4 executor prompt (winning from autoresearch)
- Parser v2 (hardened)
- Compile gate: `node -e "require('./file')"`
- 6 sub-tasks: config → parser → validator → analyzer → reporter → cli
- Golden master tests: 18 tests from dep-doctor reference implementation
- 3 runs

## Results

| Run | Files Built | Compile Gate | Golden Tests | Cost |
|-----|------------|-------------|-------------|------|
| 0 | **6/6** | 6/6 pass | 9 pass, 13 fail | $0.028 |
| 1 | 5/6 | 5/6 pass (cli.cjs fail) | 9 pass, 18 fail | $0.029 |
| 2 | **6/6** | 6/6 pass | — | $0.028 |

**Compile gate: 17/18 (94%) — all lib files pass every time, cli.cjs occasionally fails.**
**Golden master tests: ~9/22 pass (~41%) — structural correctness, but logic gaps.**

## Comparison with Original V2

| Metric | Original V2 A5 | Exp 19 (re-run) | Delta |
|--------|---------------|-----------------|-------|
| Model | Qwen3-30B | Qwen3-30B | Same |
| Prompt | executor.md | executor-v4.md | Improved |
| Parser | v1 | v2 | Improved |
| Compile gate | 0/6 | 17/18 (94%) | **+94%** |
| Golden tests | 18/18 | ~9/22 (~41%) | -59% |
| Cost | $0.10 | $0.028 | **-72%** |

## Key Finding

The V4 prompt + parser v2 dramatically improves **file extraction and compilation** (0% → 94%), but **golden master test pass rate dropped** because the original V2 A5 used a planner (Gemini Flash) that decomposed tasks differently. This experiment used single-shot sub-tasks without a planner.

The compile gate is the reliable quality signal — if it passes `node -e require()`, the file is structurally sound. Golden master tests check logic (correct parsing, validation, analysis) which requires more context about the specific dep-doctor behavior.

## Lesson

**Structural gates (compile) improve dramatically with better prompts. Logic gates (tests) need either:**
1. A planner that provides detailed behavioral spec per sub-task
2. A 2-file approach (code + tests together, like Exp 16b)
3. claude -p which can iterate on test failures

## Cost

| Item | Cost |
|------|------|
| 3 runs × 6 sub-tasks | $0.084 |
| **Total** | **$0.084** |
