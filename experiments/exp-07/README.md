# Experiment 7: Hybrid Pipeline (Cheapest API + Free Subscription)

## Hypothesis
Use cheapest API model for logic, free subscription model for complex tasks.

## Setup
- Schema + model layer: Qwen3-30B via OpenRouter ($0.009)
- main.go + tests: Haiku via claude -p (FREE)

## Results

| Step | Model | Cost | Result |
|------|-------|------|--------|
| Schema | Qwen3-30B | $0.001 | PASS |
| Model + tests | Qwen3-30B | $0.007 | PASS (10/10) |
| main.go + main_test.go | Haiku | FREE (~$0.10 API) | PASS (62s) |
| **Total** | | **$0.009 + FREE** | **App works!** |

## Key Finding
Hybrid approach gives the best of both: cheapest model for logic ($0.009), subscription model for complex UI (FREE). Total: under 1 cent of API cost.

Haiku even added drag-and-drop to the kanban board (not in the spec!).
