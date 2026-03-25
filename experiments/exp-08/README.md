# Experiment 8: MiniMax with Explicit Backtick Hint

## Hypothesis
Can MiniMax build main.go if we explicitly explain the Go backtick constraint?

## Setup
- Model: MiniMax M2.7
- Prompt includes: "In JS use string concatenation, NOT template literals"
- 3 runs

## Result: 2/3 — YES, MiniMax can build main.go with the hint!

This is the first time an OpenRouter model compiled main.go with embedded HTML.

## Key Finding
The backtick issue is a **prompt problem**, not a model problem. Explicitly explaining the constraint gives 67% success rate even on cheap models.
