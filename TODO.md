# Research TODO

## Completed
- [x] **Autoresearch executor prompts** — 3/4 sub-tasks at 100% on Qwen3-30B
- [x] **Index page with experiment map** — experiments/INDEX.md with Mermaid diagrams, results table
- [x] **Per-experiment directories** — 12 experiment READMEs in exp-01/ through exp-12/
- [x] **Spike V3 Report** — experiments/spike-v3/REPORT.md (architecture, contract-first, backtick case study)
- [x] **Autoresearch Report** — experiments/spike-v3/autoresearch-report.md (methodology, results, winning prompts)
- [x] **Compile-Fix Pipeline** — experiments/spike-v3/compile-fix-pipeline.md (gate pipeline, per-language configs)
- [x] **Escalation Strategy** — experiments/spike-v3/escalation-strategy.md (cost analysis, when it works/doesn't)
- [x] **Escalation on Failure** — exp-01: 0/5, shared blind spots can't be escalated
- [x] **Claude Subscription Models** — exp-06: Haiku 3x faster, both pass
- [x] **Full-Stack App (70s)** — exp-12: 7 files, builds, tests, runs, FREE
- [x] **Auto-Fix Pipeline measurement** — exp-04: 40-60% of errors fixed free
- [x] **MiniMax Backtick Hint** — exp-08: 2/3 with explicit hint
- [x] **V-Model Pattern (conceptual)** — exp-10: pattern works, blocked by parser
- [x] **Model Routing (single-shot)** — exp-05: 0/5, proves retry loop essential

## Experiments to Run

### Parser Hardening (Exp 13)
- [ ] Production-grade file block parser
- [ ] Handle: truncated END markers, code fences, package declarations, filename hints
- [ ] Measure: re-run exp-09 (API-only full app) with improved parser
- [ ] Target: reduce 40% parser failure rate to <5%

### Full Pipeline Retry for Model Routing (Exp 14)
- [ ] Re-run exp-05 (model routing) WITH retry loop instead of single-shot
- [ ] Route: schema/model → Qwen3-30B, main.go → MiniMax or Gemini
- [ ] Measure: does routing save cost vs. one model for everything?

### Re-run V1+V2 with All Improvements (Exp 15)
- [ ] Take winning prompts from autoresearch
- [ ] Add compile gate + auto-fix
- [ ] Re-run V1 (11 models) and V2 (dep-doctor CLI)
- [ ] Compare: improvement delta and cost delta

### Sub-Task Granularity (Exp 16)
- [ ] Test 2-3 files per sub-task (halves API calls)
- [ ] Measure: quality drop threshold
- [ ] Test: task.go + task_test.go as one sub-task (fixes model_test 0%)

### V-Model Full Loop (Exp 17)
- [ ] Fix parser to enable exp-10 completion
- [ ] Blueprint → spec + hidden acceptance tests
- [ ] Executor builds from spec only
- [ ] Run acceptance gate after build
- [ ] Measure: feedback loop iterations to pass

## Future Directions

### Bigger Applications
- [ ] SvelteKit + Go GraphQL full stack (2000+ lines)
- [ ] Test contract-first + mock-driven with parallel frontend/backend
- [ ] Measure: does the approach scale?

### Pipeline Integration
- [ ] Integrate winning strategy into Dark Factory daemon
- [ ] A/B test: new approach vs current claude -p approach
- [ ] Measure: cost per pipeline item, success rate, time
