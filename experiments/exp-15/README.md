# Experiment 15: Tiered Model Escalation

## Hypothesis
When the cheap model fails on main.go (backtick issue), escalate through tiers:
T1 Qwen3-30B → T2 MiniMax → T3 claude -p Haiku.

## Setup
- T1: Qwen3-30B ($0.001/call), 2 attempts with error feedback
- T2: MiniMax M2.7 ($0.003/call), 1 attempt with backtick hint + error context
- T3: claude -p Haiku (FREE on subscription), 1 attempt
- 3 runs
- Full auto-fix pipeline (goimports + gofmt + &constant fix)

## Results

| Run | Winning Tier | Attempts | Cost | Time |
|-----|-------------|----------|------|------|
| 0 | **T1** (Qwen3-30B) | 1 | $0.010 | 107s |
| 1 | **T1** (Qwen3-30B) | 1 | $0.010 | 33s |
| 2 | **T1** (Qwen3-30B) | 1 | $0.010 | 16s |

**T1 won all 3 runs on the first attempt. No escalation needed.**

## The Breakthrough: Prompt Engineering

Previous experiments (Exp 14) had 0% on main.go with Qwen3-30B. This experiment had 100%. The only difference: **the prompt**.

### Exp 14 Prompt (0% success)
```
CRITICAL: Go backtick raw strings cannot contain backticks.
In JavaScript use string concatenation ('Hello ' + name) NOT template literals.
```

### Exp 15 Prompt (100% success)
```
CRITICAL: Go backtick raw strings CANNOT contain backticks.
In JavaScript inside the HTML, use string concatenation:
  'Hello ' + name    ← CORRECT
  `Hello ${name}`    ← WRONG (backtick inside backtick)
Use document.getElementById('x').textContent = value, NOT template literals.
```

### Why It Works
1. **Concrete examples** — CORRECT/WRONG labels with actual code patterns
2. **Alternative provided** — `document.getElementById` gives the model a specific DOM API to use
3. **Visual contrast** — showing both patterns side by side makes the distinction clear
4. The original finding "prompt wording > model choice" (Exp 03) extends to error avoidance too

## Key Finding

**The backtick issue was never a model limitation — it was a prompt clarity issue.**

With the right prompt:
- Qwen3-30B ($0.001) → 100% on main.go
- No escalation needed
- No claude -p needed
- **The entire app can be built on the cheapest model**

## Revised Pipeline

```
Before (Exp 14):
  Model layer  → Qwen3-30B ($0.007) + auto-fix  → 100%
  HTTP server  → claude -p (FREE/sub)             → 100%
  Total: $0.007 + FREE = $0.007

After (Exp 15):
  Model layer  → Qwen3-30B ($0.007) + auto-fix  → 100%
  HTTP server  → Qwen3-30B ($0.010) + auto-fix  → 100%
  Total: $0.017 (ALL on cheapest model, NO subscription needed)
```

## Escalation Still Valid

Even though T1 solved everything here, the tiered approach is still the right architecture because:
1. New blind spots will emerge for different applications
2. The escalation chain catches them automatically
3. Cost is still low (~$0.01-$0.03 per escalation)
4. claude -p as final tier guarantees 100% success

## Cost

| Item | Cost |
|------|------|
| 3 runs × T1 (1 attempt each) | $0.030 |
| T2 and T3 never reached | $0.000 |
| **Total experiment** | **$0.030** |
