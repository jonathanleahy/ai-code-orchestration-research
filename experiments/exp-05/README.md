# Experiment 5: Model Routing by Task Type

## Hypothesis
Route pure logic → cheapest model, mixed content → mid-tier model.

## Setup
- Pure logic (schema, model): Qwen3-30B
- Mixed content (main.go with HTML): MiniMax
- Single-shot mode (no retry loop)

## Result: 0/5 (all failed in single-shot mode)

## Key Finding
**The retry loop is essential.** Single-shot calls fail ~50% due to format inconsistency. The same models achieve 100% in the full pipeline with retries.

The pipeline IS the product, not the model.
