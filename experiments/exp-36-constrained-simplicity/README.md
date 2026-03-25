# Experiment 36: Domain Expert + Constrained Simplicity

## The Balance
Domain expert says WHAT's needed. Simplicity agent says HOW to build it simply.
Simplicity agent CANNOT cut features — only simplify implementation.

## Results

| Metric | Domain Only (32) | Simplicity Only (34) | **Both (36)** |
|--------|-----------------|---------------------|---------------|
| Store | 380 lines | 173 lines | **249 lines** |
| Server | 1136 lines | 345 lines | **942 lines** |
| Total | 1516 lines | 518 lines | **1191 lines** |
| Tests | 22/22 | 13/13 | **29/29** |
| Features | All (complex) | Missing 5 | All (simple) |
| Cost | $0.96 | $0.20 | **$0.3901** |

## Simplification Examples
- Address: single textarea, not 5 fields
- Invoice print: browser print button, not PDF library
- Edit client: reuse add form with pre-filled values
- Search: strings.Contains filter
- Invoice status: simple string, not state machine
