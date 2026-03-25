# Research TODO

## Active
- [ ] **Autoresearch executor prompts** — running now, finding optimal prompt per sub-task type for Qwen3-30B

## Experiments to Run

### Escalation on Failure (Spike 3b)
- [ ] When cheap model fails and auto-fix doesn't work:
  - Variation 1: Call the SAME cheap model with the error message (retry)
  - Variation 2: Call a STRONGER model (e.g., Gemini Flash) with the error + original code
  - Measure: which is cheaper per successful fix? Does escalation beat retry?
  - The backtick issue is the perfect test case — cheap model fails, does Gemini fix it?

### Re-run Spikes 1+2 with Improved Executor (Spike 1b/2b)
- [ ] Take the winning prompts from autoresearch
- [ ] Re-run the V1 model comparison (11 models) with improved prompts
- [ ] Re-run V2 Node.js CLI with improved prompts + compile gate
- [ ] Compare: how much did prompt engineering + auto-fix improve results?
- [ ] Record improvement delta and cost delta

### More Sub-Task Granularity Testing
- [ ] Can 2-3 files per sub-task work? (halves API calls)
- [ ] What's the maximum files-per-task before quality drops?

### Model Routing by Task Type
- [ ] Route logic tasks → cheapest model, mixed-content tasks → mid-tier
- [ ] Measure if routing saves cost vs using one model for everything

## Documentation

### Index Page with Experiment Map
- [ ] Create an interactive index page (HTML/JS)
- [ ] Node diagram (Mermaid or D3.js) showing:
  - Core idea at centre
  - Connections to each spike (V1, V2, V3)
  - Sub-experiments branching off each spike
  - Results annotated on each node
- [ ] Clickable — each node links to its report

### Spike V3 Report
- [ ] Write dedicated spike-v3/REPORT.md (not just addenda)
- [ ] Include: architecture diagrams, contract-first pattern, mock-driven building
- [ ] claude -p findings with code examples
- [ ] Compile gate + auto-fix pipeline diagram
- [ ] Backtick lesson as a case study

### Autoresearch Report
- [ ] Document the autoresearch methodology
- [ ] Results table: pass rate per prompt variation per sub-task type
- [ ] Winning prompts with explanations of why they work
- [ ] Cost analysis: autoresearch cost vs improvement value

## Architecture Improvements

### Compile-Fix Pipeline
- [ ] Document the full gate pipeline: goimports → gofmt → go build → go vet → go test
- [ ] Add per-language gate configs (Go, Node.js, Python)
- [ ] Measure: what % of errors does each gate level catch?

### Escalation Strategy
- [ ] Document: cheap model → auto-fix → retry → escalate to stronger model
- [ ] Measure cost per successful build with escalation vs without

## New Experiment Ideas

### Spec → Test Data → Use Cases (V-Model Pattern)
- [ ] Blueprint produces: spec + real test data + use cases (acceptance criteria)
- [ ] Executor builds code to pass unit tests (doesn't see use cases)
- [ ] Unit tests run during development
- [ ] Use cases run AFTER complete build (acceptance gate)
- [ ] Prevents "teaching to the test" — executor optimises for spec, not test
- [ ] System loops: fail use cases → back to executor with feedback
- [ ] This is a proper V-model with AI at each layer

### PR Review Gate (Quality + Security)
- [ ] After code passes structural gates, submit as a "PR"
- [ ] AI reviewer checks: code quality, security patterns, error handling
- [ ] Like pr-review.sh but for AI-generated code
- [ ] Could use a different (stronger) model as reviewer
- [ ] Approval required before code enters assembled output
- [ ] Catches issues tests don't: hardcoded secrets, SQL injection, missing auth

### Model Writes Its Own Tests (fix model_test 0%)
- [ ] Instead of testing against golden master types, let the model write task.go + task_test.go together
- [ ] The tests verify the spec, not specific types
- [ ] Then run golden master integration tests as the acceptance gate
- [ ] This matches how real developers work — they write code and tests together

## Future Directions

### Bigger Applications
- [ ] SvelteKit + Go GraphQL full stack (the original V3 plan)
- [ ] Test contract-first + mock-driven with parallel frontend/backend building
- [ ] Measure: does the approach scale to 2000+ line apps?

### Pipeline Integration
- [ ] Integrate winning strategy into Dark Factory daemon
- [ ] A/B test: new approach vs current claude -p approach
- [ ] Measure: cost per pipeline item, success rate, time
