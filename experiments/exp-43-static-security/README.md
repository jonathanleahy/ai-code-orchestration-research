# Experiment 43: Static Security Analysis

## Tools Run
| Tool | What It Checks | Issues |
|------|---------------|--------|
| **gosec** | Hardcoded creds, SQL injection patterns, weak crypto, unsafe functions | See gosec.txt |
| **staticcheck** | Code correctness, performance, simplification opportunities | See staticcheck.txt |
| **govulncheck** | Known CVEs in Go dependencies | See govulncheck.txt |
| **go vet** | Suspicious constructs, printf format errors | See govet.txt |

## Key Finding
These tools run on SOURCE CODE — no server needed. They should be
a mandatory gate after every code generation, before the app is started.

## Pipeline Integration
```
Store code generated → gosec + staticcheck + go vet
Server code generated → gosec + staticcheck + go vet
go.mod updated → govulncheck
```

All free, all instant, catches issues before runtime.
