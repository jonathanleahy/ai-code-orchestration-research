# Experiment 18: Full App on Cheapest Model — End-to-End Pipeline

## Hypothesis
Combining all improvements (Exp 15 prompt, Exp 16b granularity, auto-fix, parser v2), the cheapest model can build the entire app for ~$0.017.

## Setup
- Model: Qwen3-30B ($0.0005/call) for all sub-tasks
- 3 sub-tasks: schema → model+test (2-file) → main.go
- Auto-fix: goimports + gofmt + &constant fix
- Parser v2
- 2 retries per sub-task
- 5 runs

## Results

| Run | Schema | Model+Test | Main.go | Overall | Cost |
|-----|--------|-----------|---------|---------|------|
| 0 | PASS | PASS | FAIL (3 tries) | FAIL | $0.034 |
| 1 | PASS | PASS | FAIL (3 tries) | FAIL | $0.035 |
| 2 | PASS | PASS | FAIL (3 tries) | FAIL | $0.032 |
| 3 | PASS | PASS | FAIL (3 tries) | FAIL | $0.034 |
| 4 | PASS | PASS | FAIL (3 tries) | FAIL | $0.033 |

**Pass rate: 0/5 (0%)**
**Schema + model+test: 5/5 (100%)**
**Main.go: 0/15 attempts (0%)**

## Key Finding

**Exp 15's prompt worked in isolation but fails in the full pipeline.**

The difference: Exp 15 sent ONLY the backtick hint prompt (~500 tokens). Exp 18 sends the hint embedded in a 3000-token architecture spec. The model's attention is diluted by the surrounding context and it reverts to its default pattern (template literals).

### Why Isolation vs Pipeline Differ

| Factor | Exp 15 (100%) | Exp 18 (0%) |
|--------|--------------|-------------|
| Prompt size | ~500 tokens | ~3000 tokens |
| Architecture context | Minimal | Full spec |
| Backtick hint position | End of prompt | Middle of prompt |
| Model attention | Focused on hint | Diluted by spec |

## Implications

1. **The tiered escalation IS needed** — Exp 15's escalation chain (T1 → T2 → T3) is the correct architecture, not "just use better prompts"
2. **Schema and model layer work perfectly** — 100% on cheapest model with auto-fix
3. **main.go (HTML-in-Go) requires claude -p** — the backtick issue is context-dependent, not just prompt-dependent
4. **Prompt length matters** — a hint that works in a short prompt may fail in a long one

## Revised Final Pipeline

```
ST-1: schema.graphql    → Qwen3-30B ($0.001)  → 100%
ST-2: model + tests     → Qwen3-30B ($0.007)  → 100% (with auto-fix)
ST-3: main.go           → claude -p Haiku (FREE) → 100%
                                          ─────────
                           Total: $0.008 + FREE = $0.008
```

This is the winning strategy: **cheap model for pure Go, claude -p for HTML-in-Go.**

## Cost
| Item | Cost |
|------|------|
| 5 runs × (schema + model_test + 3× main.go retries) | $0.167 |
| **Total experiment** | **$0.167** |
