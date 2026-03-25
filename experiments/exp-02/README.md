# Experiment 2: Sub-Task Granularity

## Hypothesis
Can one API call produce multiple files? What's the optimal number?

## Setup
- Model: Qwen3-30B
- Test: 2, 3, and 4 files per call
- 3 runs each

## Results

| Files/Call | Pass Rate | Notes |
|-----------|----------|-------|
| 2 files | 2/3 (67%) | Parser sometimes misses second file |
| 3 files | 0/3 (0%) | Pre-existing go.mod conflicts |
| **4 files** | **3/3 (100%)** | Model owns everything |

## Key Finding
When the model creates ALL files (including go.mod), it works 100%. Pre-existing scaffolding can conflict with the model's output.

**Rule: Give each sub-task full ownership of its files.**
