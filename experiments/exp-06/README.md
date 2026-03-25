# Experiment 6: Claude Subscription Models (Sonnet vs Haiku)

## Hypothesis
Which Claude model is best for `claude -p` builds?

## Setup
- Sonnet and Haiku build the full task-board (model + main + tests)
- Both via `claude -p --dangerously-skip-permissions`

## Results

| Model | Files | Builds | Tests | Time | Cost (sub) | Cost (API) |
|-------|-------|--------|-------|------|-----------|-----------|
| **Haiku** | 4 | ✅ | ✅ | **80s** | **FREE** | ~$0.10 |
| Sonnet | 4 | ✅ | ✅ | 252s | FREE | ~$0.50 |

## Key Finding
**Haiku is 3x faster with equal quality.** Both solve the backtick issue because `claude -p`'s tool system iterates (write → compile → fix → retry).

For subscription users, Haiku is the optimal choice.
