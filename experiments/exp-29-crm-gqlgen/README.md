# Experiment 29: CRM with gqlgen GraphQL Backend

## Brief (from screenshot)

A Client Management System with:
- **Client Dashboard**: Searchable list of clients (name, company, status)
- **Client Profile**: Central hub per client
- **History & Activity Log**: Chronological timeline of all interactions (emails, calls, meetings). Quick-entry tool to log new history.
- **Billing & Document Module**: Generate/upload invoices per client. Drag-and-drop file attachments (contracts, briefs, photos). Separate from activity log.

## What This Tests

The automated pipeline (Exp 27) assumes Go stdlib + REST. This brief requires:

| Requirement | Current Pipeline | Change Needed |
|-------------|-----------------|---------------|
| External deps (gqlgen) | stdlib only | `go get` step |
| Code generation | Write from scratch | `go generate` + fill stubs |
| GraphQL schema | N/A | Schema-first design |
| Image brief input | Text only | Extract requirements from screenshot |
| File uploads | N/A | Multipart form handling |
| Multi-entity relations | Simple CRUD | Client → Invoices, Client → History |

## Pipeline Changes Discovered

### 1. External Dependencies Need a `go get` Step
The pipeline assumed stdlib only. gqlgen requires `go get github.com/99designs/gqlgen` + transitive deps. **Fix**: Add a dependency resolution step before code generation.

### 2. Code Generation is Multi-Step
gqlgen workflow: schema → `go generate` → fill stubs. The current one-shot pipeline can't handle this. **Fix**: The pipeline needs to support "write schema → run command → read output → fill in stubs" as a pattern.

### 3. AI Must Work WITH Generated Code, Not From Scratch
gqlgen generates `generated.go`, `models_gen.go`, and resolver stubs. The AI should fill stubs, not rewrite everything. **Fix**: Pipeline step that reads generated files and fills `panic("not implemented")` sections.

### 4. Resolver Wiring Needs Context
The AI filled all resolver functions but forgot to add the store to the `Resolver` struct. It wrote `store.clients` instead of `r.clients`. **Fix**: The prompt must include "the Resolver struct is your state — add fields to it, reference them with r."

### 5. Two claude -p Calls Needed
- Call 1: Schema + gqlgen init + generate + fill resolvers
- Call 2: Fix wiring issues (store → Resolver struct)
Total: ~$0.50 for the full build

### 6. Frontend is Separate
The 513-line HTML frontend is separate from the Go code. It uses `fetch('/query')` for GraphQL. This works well — the frontend is a single embedded file.

## Results

| Component | Lines | Status |
|-----------|-------|--------|
| GraphQL schema | 55 | 3 types, 5 queries, 8 mutations |
| Resolvers | 265 | All implemented (0 stubs) |
| Generated code | ~2000 | gqlgen auto-generated |
| main.go | 33 | Server + playground + frontend |
| Frontend HTML | 513 | Client list, profile, history, billing |
| **Total custom code** | **~866** | **BUILD PASS** |

### What Works
- GraphQL playground at `/`
- Full CRUD via GraphQL mutations
- Client search via query parameter
- Activity timeline per client
- Invoice management per client
- 17KB HTML frontend with tabs (Details, History, Billing)

### Cost
- claude -p call 1 (schema + resolvers): ~$0.30
- claude -p call 2 (fix wiring): ~$0.20
- **Total: ~$0.50**

## Key Finding

**gqlgen needs an iterative pipeline, not one-shot.** The current automated pipeline (Exp 27) can't handle:
1. External dependency resolution
2. Code generation (`go generate`)
3. Stub-filling (working WITH generated code)
4. Wiring fixes (store → struct fields)

For Dark Factory integration, gqlgen projects need a "schema-first" pipeline variant:
```
Schema → go generate → fill resolvers → fix wiring → build → test
```
This is 2-3 claude -p calls instead of 1, at ~$0.50 total.
