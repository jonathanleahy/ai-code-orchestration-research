# Experiment 16: Sub-Task Granularity

## Hypothesis
Combining model/task.go + model/task_test.go into a single sub-task (2 files per call) fixes the model_test 0% problem, because the model writes tests against its own types.

## Setup
- Model: Qwen3-30B ($0.0005/call) via OpenRouter
- 4 configurations with increasing files per task
- 3 runs per config (5 for 16b)
- Auto-fix pipeline: goimports + gofmt + &constant fix

## Results

### Without Auto-Fix (Exp 16)

| Config | Files/Task | Pass Rate | Avg Cost | Failure Reason |
|--------|-----------|-----------|----------|----------------|
| 1file_model_only | 1 | **100%** | $0.0033 | — |
| 1file_test_only | 1 | 0% | $0.0037 | Tests vs golden types |
| 2file_model_and_test | 2 | 0% | $0.0065 | &constant error |
| 3file_model_test_main | 3 | 0% | $0.0143 | Cross-package refs |

### With Auto-Fix (Exp 16b)

| Config | Files/Task | Pass Rate | Avg Cost | Notes |
|--------|-----------|-----------|----------|-------|
| 2file + auto-fix | 2 | **100%** (5/5) | $0.0066 | &constant auto-fixed |

## Key Finding

**2 files per sub-task + auto-fix = 100% pass rate on the cheapest model.**

The original 0% for model_test had two causes:
1. **Architectural**: Tests compiled against golden master types (different from model output)
2. **Go-specific**: `&StatusDoing` fails because Go can't take address of constants

Fix 1: Let the model write BOTH task.go + task_test.go together (tests match its own types).
Fix 2: Auto-fix replaces `&ConstName` with `statusPtr(ConstName)` helper.

## Granularity Threshold

| Files | Quality | Cost | Verdict |
|-------|---------|------|---------|
| 1 | Works for implementation, fails for tests | $0.003 | OK for single files |
| 2 | **100% with auto-fix** | $0.007 | **Sweet spot** |
| 3 | Cross-package reference errors | $0.014 | Too many — model loses context |

**The sweet spot is 2 files per sub-task** when the files are in the same package. 3+ files spanning packages causes the model to lose track of type namespacing.

## Auto-Fix: &Constant

```
BEFORE: store.Update(id, nil, nil, &StatusDoing)
ERROR:  cannot take address of StatusDoing (constant "DOING" of string type Status)

AUTO-FIX adds helper:
  func statusPtr(s Status) *Status { return &s }

AFTER:  store.Update(id, nil, nil, statusPtr(StatusDoing))
```

This is a shared blind spot — all models make this mistake. The auto-fix catches it structurally for free.

## Cost

| Item | Cost |
|------|------|
| Exp 16 (4 configs × 3 runs) | $0.082 |
| Exp 16b (1 config × 5 runs) | $0.033 |
| **Total** | **$0.115** |
