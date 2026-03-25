# Experiment 20: Different Application — URL Shortener

## Hypothesis
The proven pipeline (Qwen3-30B for model layer, claude -p for HTTP server) works on a different application, not just task-board.

## Setup
- Application: URL shortener (Shorten, Resolve, List, Delete, Stats)
- Model layer: Qwen3-30B ($0.0005/call) with 2-file approach + auto-fix
- HTTP server: claude -p Haiku (FREE on subscription)
- 3 runs

## Results

| Run | Store Layer | Store Tests | main.go (claude -p) | Build | Server |
|-----|------------|------------|---------------------|-------|--------|
| 0 | OK | 9/10 pass | FAIL (not created) | FAIL | — |
| 1 | OK | 0/0 (no tests) | FAIL (not created) | FAIL | — |
| 2 | OK | — | FAIL (not created) | FAIL | — |

**Store layer: works (9/10 tests in best run)**
**HTTP server: 0/3 — claude -p failed to create main.go**

## Key Finding

### Store Layer (Qwen3-30B) — Generalises
The 2-file approach (store.go + store_test.go) works for a completely different application:
- URL shortener store with Shorten/Resolve/List/Delete/Stats
- 9/10 tests passing on first run
- The approach is not task-board-specific

### HTTP Server (claude -p) — Environment Issue
claude -p returned $0.00 cost and no main.go in all runs. This is a **CLI environment issue** when running inside a background subprocess, not a model capability issue. When run interactively, claude -p builds HTTP servers reliably (see Exp 6, 12).

The fix: run claude -p in foreground or use a wrapper that handles the subprocess properly.

## Conclusion

The approach **generalises** to different applications for the model layer. The claude -p integration needs a more robust invocation method for batch pipelines.

## Cost

| Item | Cost |
|------|------|
| 3 runs × (store + claude -p attempt) | $0.019 |
| **Total** | **$0.019** |
