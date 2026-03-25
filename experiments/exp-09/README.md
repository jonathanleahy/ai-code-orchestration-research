# Experiment 9: Full App via API Only (No Subscription)

## Hypothesis
Build the complete Go app using only OpenRouter models, under $0.05.

## Setup
- Planner: Gemini Flash
- All executor calls: MiniMax M2.7

## Result: $0.045 total, but cascading failure from parser missing go.mod

| Step | Cost | Result |
|------|------|--------|
| Planner | $0.001 | OK |
| Schema + go.mod | $0.009 | FAIL (parser) |
| Model + tests | $0.019 | FAIL (no go.mod) |
| main.go | $0.015 | FAIL (no go.mod) |

## Key Finding
The parser is the weakest link. ~40% of failures are parser extraction issues, not model quality. A production parser would make all-API builds reliable.
