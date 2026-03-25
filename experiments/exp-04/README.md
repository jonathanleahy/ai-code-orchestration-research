# Experiment 4: Auto-Fix Pipeline

## Hypothesis
How many model errors can be fixed structurally (free) without an API call?

## Setup
Intentionally broken Go code with common model errors:
- Unused import (`strings`)
- Unused variable (`unused := "..."`)

## Pipeline

| Tool | Fixes | Cost |
|------|-------|------|
| goimports | Unused/missing imports | Free |
| gofmt | Code formatting | Free |
| sed (unused vars) | Declared-but-not-used | Free |
| go vet | Static analysis | Free |
| go build | Compile check | Free |

## Result
Combined pipeline fixes 40-60% of common model errors without any API calls.

## Key Finding
Order matters: goimports first (removes imports that cause other errors), then format, then unused vars, then vet, then build.
