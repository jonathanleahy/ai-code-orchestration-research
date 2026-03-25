# Experiment 27: Fully Automated Pipeline

## Brief
"Build a bookmark manager with tags and search"

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
| Store | COMPILES | 4/6 |
| Server | COMPILES (314 lines) | 20/20 |
| **Total** | | **24/26 (92%)** |

## Cost: $0.3229

## Key: Auto-Fix Loop
Store tests go through an auto-fix loop:
- Structural fixes (goimports, gofmt, &constant) — FREE
- Test failure fixes (cheap model reads error + code, outputs fix) — $0.001/attempt
- Max 3 attempts before moving on

Each test function creates its own store (no shared state) to avoid the
Exp 26 ListClients bug.
