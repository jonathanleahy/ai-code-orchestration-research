# Experiment 12: Full-Stack App in 70 Seconds 🏆

## Hypothesis
Can Haiku build a complete full-stack app (Go backend + HTML frontend + TypeScript types) in one call?

## Setup
- Model: Haiku via claude -p
- Input: schema.graphql only
- Target: 7 files (Go + HTML + TypeScript)

## Result: YES — 7 files, builds, tests pass, app runs, 70 seconds, FREE

| Metric | Value |
|--------|-------|
| Files | 7 (3 Go + 1 HTML + 2 TypeScript + 1 GraphQL) |
| Builds | ✅ |
| Tests | ✅ |
| Runs | ✅ (http://localhost:8896) |
| Time | **70 seconds** |
| Cost | **FREE** (subscription) / ~$0.12 (API) |

## Files Created
```
schema.graphql       ← Contract
go.mod               ← Module  
model/task.go        ← CRUD store
model/task_test.go   ← 10 unit tests
main.go              ← HTTP server + API
frontend/index.html  ← Kanban board UI
frontend/src/lib/types.ts  ← TypeScript interfaces
frontend/src/lib/api.ts    ← API client
```

## Key Finding
**AI builds real full-stack applications from a schema contract in about a minute.** This is the vision realised.
