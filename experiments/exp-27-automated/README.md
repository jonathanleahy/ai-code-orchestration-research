# Experiment 27: Fully Automated Pipeline

## Brief
"Build a CRM with invoice generator, client history, comments, address book, and GraphQL gqlgen backend"

## Pipeline (zero intervention)
1. Persona discovery + interviews
2. MVP synthesis + dev review (combined)
3. Go types + screen wireframes
4. Store code (Qwen3-30B)
5. Store tests (Qwen3-30B) + auto-fix loop (3 max attempts)
6. HTTP server (claude -p)
7. HTTP integration tests (claude -p)

## Results

| Component | Status | Tests |
|-----------|--------|-------|
| Store | COMPILES | 16/16 |
| Server | COMPILES (720 lines) | 25/25 |
| **Total** | | **41/41 (100%)** |

## Cost: $0.4858

## Key: Auto-Fix Loop
Store tests go through an auto-fix loop:
- Structural fixes (goimports, gofmt, &constant) — FREE
- Test failure fixes (cheap model reads error + code, outputs fix) — $0.001/attempt
- Max 3 attempts before moving on

Each test function creates its own store (no shared state) to avoid the
Exp 26 ListClients bug.
