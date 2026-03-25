# Experiment 47: Code Quality Metrics

## Metrics on Exp 38 (Progressive CRM, 32/32 tests)

### Coverage: 57.4%
| Component | Coverage | Gap |
|-----------|---------|-----|
| main.go (HTTP handlers) | 79.5% | handleClientDetail at 51.4% |
| store/store.go | 0% | All 13 functions untested directly |

Store has 0% direct coverage because tests only go through HTTP endpoints.

### Complexity: Max 9 (threshold: 10)
| Function | Complexity | Lines |
|----------|-----------|-------|
| handleClientDetail | 9 | 350 |
| handleCreateInvoice | 8 | 43 |
| TestCreateInvoiceAPI | 9 | (test code) |

All under 10 — acceptable but handleClientDetail is close to the limit.

### File Stats
| File | Lines | Functions | Avg Lines/Function |
|------|-------|-----------|-------------------|
| main.go | 894 | 15 | ~60 |
| store/store.go | 227 | 13 | ~17 |

### Pipeline Integration
Run after every build:
```bash
go test ./... -cover          # Coverage gate: >70%?
gocyclo -over 10 .            # Complexity gate: any >10?
```
