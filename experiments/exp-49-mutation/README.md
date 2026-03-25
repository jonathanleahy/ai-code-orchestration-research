# Experiment 49: Mutation Testing

## Do the tests catch code changes?

14 mutations applied to Exp 38's CRM (32 tests).

### Results
| Metric | Value |
|--------|-------|
| Mutations | 14 |
| Caught (tests broke) | 8 |
| Survived (tests still pass!) | 1 |
| Skipped (string not found) | 5 |
| **Mutation Score** | **89%** |

### Surviving Mutations (gaps in test coverage)
- Don't set CreatedAt: CreatedAt is zero

### What This Measures
- High score (>80%): tests are strong, catch most changes
- Low score (<50%): tests are weak, many code paths untested
- Each surviving mutation = a bug that could ship undetected
