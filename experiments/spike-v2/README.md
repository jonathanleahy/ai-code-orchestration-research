# Spike V2: Multi-Model Orchestration Research

## Can AI build a real application from scratch — and what's the cheapest way?

This research spike tests **4 fundamentally different approaches** to AI-driven software development, across **6 models** at price points spanning 200x (from $0.0007 to $0.15 per call). The goal: find the optimal architecture for Dark Factory's autonomous pipeline.

---

## Table of Contents

1. [Background & Motivation](#background--motivation)
2. [The Test Application](#the-test-application)
3. [Four Approaches](#four-approaches)
4. [Models Under Test](#models-under-test)
5. [Experiment Configurations](#experiment-configurations)
6. [Infrastructure](#infrastructure)
7. [Results & Analysis](#results--analysis)
8. [Conclusions & Recommendations](#conclusions--recommendations)

---

## Background & Motivation

### The Problem

Dark Factory's autonomous pipeline uses Claude to write code. But:

- **Sonnet via `claude -p` tool mode:** $1.68 per task, 0% success rate on complex tasks
- **Sonnet via API with imperative prompts:** $0.015 per task, 100% success rate
- **Qwen3 Coder 30B via OpenRouter:** $0.0007 per task, unknown quality on complex tasks

The 2400x cost difference between the worst and best approach means the choice of architecture matters more than the choice of model.

### Spike V1 Findings

We tested 11 models on a simple bash task ("add --model flag to queue-add.sh"):

```mermaid
graph LR
    subgraph "Cost per Task"
        A["Claude Opus<br/>$2.50+"] --> B["Claude Sonnet (claude -p)<br/>$0.34-$1.68"]
        B --> C["Claude Sonnet (API)<br/>$0.015"]
        C --> D["Gemini/MiniMax/DeepSeek<br/>$0.008-$0.015"]
        D --> E["Qwen3 Coder 30B<br/>$0.0007"]
    end

    style A fill:#ff4444
    style B fill:#ff8844
    style C fill:#44aa44
    style D fill:#44aa44
    style E fill:#44ff44
```

**Key insight:** The model isn't the bottleneck — the **prompt pattern** is. All 11 models (including Sonnet) produce correct code when given:
1. The existing file content inline
2. Specific instructions ("add X after Y")
3. The "full file output" format (`--- FILE: path --- ... --- END FILE ---`)

### Research Questions

This spike answers 10 questions:

1. Does careful orchestration (A) beat brute-force generation (B)?
2. Does peer review by multiple models (C) catch bugs that single review misses?
3. Does evolutionary pressure (D) find solutions that single-shot misses?
4. Is a reviewer layer worth its cost?
5. Which model works best at which role (planner/executor/reviewer)?
6. Can the cheapest model ($0.0007/call) produce passing code given enough attempts?
7. Does writing tests before implementation reliably guide the code?
8. What is the minimum cost for a working 500-line application?
9. Which approach produces the most maintainable code?
10. Can approaches be combined for even better results?

---

## The Test Application

### `dep-doctor` — Dependency Health Checker

A Node.js CLI tool that reads `package.json`, analyzes dependency health, and outputs structured reports. Chosen because:

- **Structurally verifiable** — `node cli.js --help` exit code, JSON validity, test pass count
- **Realistic complexity** — 7 files, ~500 lines, multiple modules with dependencies
- **Zero npm dependencies** — the tool itself uses only Node.js built-ins (matches Dark Factory convention)
- **Multiple subcommands** — tests routing, argument parsing, output formatting

```mermaid
graph TD
    CLI["cli.js<br/>Entry point<br/>5 subcommands"]
    CLI --> Parser["lib/parser.js<br/>Read package.json"]
    CLI --> Analyzer["lib/analyzer.js<br/>Health analysis"]
    CLI --> Reporter["lib/reporter.js<br/>JSON/table/summary"]
    CLI --> Config["lib/config.js<br/>Config file I/O"]
    Parser --> Validator["lib/validator.js<br/>Semver + SPDX"]
    Analyzer --> Parser
    Analyzer --> Validator
    Tests["test/dep-doctor.test.js<br/>15+ test cases"]
    Tests --> CLI
    Tests --> Parser
    Tests --> Validator
    Fixtures["fixtures/<br/>valid, malformed,<br/>empty package.json"]
    Tests --> Fixtures
```

### File Structure

| File | Lines | Purpose | Validation Gate |
|------|-------|---------|----------------|
| `fixtures/valid/package.json` | ~20 | Test data | `JSON.parse()` succeeds |
| `fixtures/malformed/package.json` | ~5 | Invalid JSON test data | File exists |
| `fixtures/empty/package.json` | ~3 | Edge case | `JSON.parse()` succeeds |
| `lib/validator.js` | ~60 | Semver + SPDX validation | `typeof isValidSemver === 'function'` |
| `lib/parser.js` | ~60 | Parse package.json | `typeof parse === 'function'` |
| `lib/config.js` | ~40 | Config file read/write | `typeof loadConfig === 'function'` |
| `lib/analyzer.js` | ~100 | Dependency health analysis | `typeof analyze === 'function'` |
| `lib/reporter.js` | ~80 | Output formatting | `typeof formatJson === 'function'` |
| `cli.js` | ~80 | Entry point | `--help` exits 0, `unknown` exits 1 |
| `test/dep-doctor.test.js` | ~120 | 15+ test cases | Exits 0, 15+ "PASS" in stdout |

### Acceptance Criteria

The application is "done" when:
1. `node cli.js --help` exits 0 and lists 5 subcommands
2. `node cli.js scan --path fixtures/valid` outputs valid JSON with dependency list
3. `node cli.js check --path fixtures/valid` exits 0 (healthy)
4. `node cli.js check --path fixtures/malformed` exits 1 (unhealthy)
5. `node test/dep-doctor.test.js` exits 0 with 15+ PASS lines and 0 FAIL lines

---

## Four Approaches

### Approach A: 4-Layer Quality-First

Inspired by Karpathy's **autoresearch** pattern: plan carefully, execute precisely, review critically.

```mermaid
flowchart TD
    Arch["Layer 1: Architecture<br/>(human-written)"]
    Plan["Layer 2: Planner<br/>(1 model call)<br/>→ sub-tasks.json"]
    Exec["Layer 3: Executor<br/>(1 call per sub-task)<br/>→ file blocks"]
    Review["Layer 4: Reviewer<br/>(1 call per sub-task)<br/>→ keep/discard/fix"]

    Arch --> Plan
    Plan --> |"Gate: valid JSON,<br/>8-15 tasks,<br/>no circular deps"| Exec
    Exec --> |"Gate: node --check,<br/>exports exist,<br/>syntax valid"| Review
    Review --> |"keep"| Done["✅ Copy to assembled/"]
    Review --> |"fix"| Exec
    Review --> |"discard"| Skip["⏭️ Next sub-task"]
    Exec --> |"Gate FAIL<br/>(max 3 retries)"| Exec

    style Arch fill:#4488ff
    style Plan fill:#44aa44
    style Exec fill:#ffaa44
    style Review fill:#aa44ff
```

**Cost model:** 1 plan call + N execute calls + N review calls = ~2N+1 calls total
**Estimated cost:** $0.20-$0.60 per build (OpenRouter) or **free** (Sonnet subscription)

### Approach B: Generate-and-Filter

Inspired by DeepMind's **AlphaCode**: generate many candidates cheaply, filter by tests.

```mermaid
flowchart TD
    Arch["Layer 1: Architecture<br/>(human-written)"]
    Plans["Layer 2: Generate 5 plans<br/>(5 model calls)"]
    Best["Select best plan<br/>(structural scoring)"]
    Tests["Layer 3: Write tests first<br/>(1 call per sub-task)"]
    Gen["Layer 4: Generate 10 candidates<br/>(10 cheap calls per sub-task)"]
    Filter["Filter: run tests<br/>(structural gate)"]

    Arch --> Plans
    Plans --> Best
    Best --> Tests
    Tests --> |"Gate: tests exist,<br/>syntax valid"| Gen
    Gen --> Filter
    Filter --> |"First passing<br/>candidate wins"| Done["✅ Copy winner"]
    Filter --> |"All 10 fail"| Retry["Generate 10 more<br/>(or skip)"]

    style Arch fill:#4488ff
    style Plans fill:#44aa44
    style Gen fill:#ffaa44
    style Filter fill:#ff4444
```

**Key insight:** At $0.0007/call (Qwen3 30B), generating 1000 candidates costs $0.70. Quality comes from **filtering**, not from the model.

**Cost model:** 5 plans + N tests + 10N candidates = ~11N+5 calls total
**Estimated cost:** $0.50-$1.50 per build

### Approach C: LLM Council

Inspired by Karpathy's **llm-council**: multiple models generate independently, then peer-review each other.

```mermaid
flowchart TD
    Arch["Layer 1: Architecture<br/>(human-written)"]
    Plan["Layer 2: Planner<br/>(1 call)"]

    subgraph Council["Layer 3: Council (per sub-task)"]
        M1["Model A<br/>generates"]
        M2["Model B<br/>generates"]
        M3["Model C<br/>generates"]
    end

    subgraph Review["Layer 4: Peer Review (anonymous)"]
        R1["Model A reviews<br/>B and C's code"]
        R2["Model B reviews<br/>A and C's code"]
        R3["Model C reviews<br/>A and B's code"]
    end

    Chair["Layer 5: Chairman<br/>picks best or<br/>synthesizes"]

    Arch --> Plan
    Plan --> Council
    Council --> Review
    Review --> Chair
    Chair --> Done["✅ Best code selected"]

    style Arch fill:#4488ff
    style Council fill:#44aa44
    style Review fill:#ffaa44
    style Chair fill:#aa44ff
```

**Key insight:** Models are "surprisingly willing to select another LLM's response as superior to their own." — Karpathy

**Cost model:** 1 plan + 3N generate + 3N review + N chairman = ~7N+1 calls total
**Estimated cost:** $0.50-$0.80 per build

### Approach D: Evolutionary

Inspired by **genetic algorithms**: breed, test, select, mutate, repeat.

```mermaid
flowchart TD
    Arch["Layer 1: Architecture<br/>(human-written)"]
    Plan["Layer 2: Planner<br/>(1 call)"]

    subgraph Evolution["Layer 3: Evolution (per sub-task)"]
        G0["Gen 0: 5 candidates"]
        T0["Test all 5<br/>Score: pass/fail + quality"]
        S0["Select top 2"]
        G1["Gen 1: Mutate top 2<br/>'Improve this code'<br/>'Fix these failures'"]
        T1["Test mutations"]
        S1["Select top 2"]
        G2["Gen 2: Mutate again"]
        T2["Test"]
        GN["... Gen N"]
    end

    Arch --> Plan
    Plan --> Evolution
    G0 --> T0 --> S0 --> G1 --> T1 --> S1 --> G2 --> T2 --> GN
    GN --> Done["✅ Fittest survives"]

    style Arch fill:#4488ff
    style Evolution fill:#44aa44
    style G0 fill:#ffaa44
    style G1 fill:#ffaa44
    style G2 fill:#ffaa44
```

**Key insight:** At $0.0007/call, 5 generations × 5 candidates = 25 calls = $0.0175 per sub-task. Evolution is almost free.

**Cost model:** 1 plan + N × (pop_size × generations) calls = ~125N+1 calls total (but at $0.0007 each!)
**Estimated cost:** $0.20-$0.50 per build

---

## Models Under Test

```mermaid
graph LR
    subgraph "Price Tiers"
        Free["🟢 Free<br/>Sonnet (subscription)"]
        Cheap["🟡 Ultra-cheap<br/>Qwen3 30B: $0.0007"]
        Mid["🟠 Mid-range<br/>Qwen3: $0.008<br/>DeepSeek: $0.011"]
        Quality["🔴 Quality<br/>Gemini Flash: $0.013<br/>MiniMax: $0.013"]
    end
```

| Model | ID | $/call (est) | Strengths | Best For |
|-------|-----|-------------|-----------|----------|
| **Sonnet** | `claude -p` | **free** | Reliable with imperative prompts | Planner, reviewer |
| **Qwen3 30B** | `qwen/qwen3-coder-30b-a3b-instruct` | $0.0007 | Ultra-cheap, fast | Volume executor (B, D) |
| **Qwen3 Coder** | `qwen/qwen3-coder` | $0.008 | Cheapest quality model | Executor (A) |
| **MiniMax M2.7** | `minimax/minimax-m2.7` | $0.013 | Best V1 validation score | Council member (C) |
| **Gemini Flash** | `google/gemini-2.5-flash` | $0.013 | Best V1 preservation | Planner, chairman (C) |
| **DeepSeek V3.2** | `deepseek/deepseek-v3.2` | $0.011 | Strong value | Council member (C) |

---

## Experiment Configurations

### 10 Configs Across 4 Approaches

```mermaid
graph TD
    subgraph "Approach A: Quality-First"
        A1["A1: All Sonnet<br/>FREE baseline"]
        A2["A2: Gemini plan<br/>+ Qwen3 exec"]
        A3["A3: Gemini plan<br/>+ MiniMax exec"]
        A4["A4: Sonnet plan<br/>+ Qwen3-30B exec"]
    end

    subgraph "Approach B: Generate & Filter"
        B1["B1: Gemini plans<br/>+ Qwen3-30B exec×10"]
        B2["B2: Sonnet plans<br/>+ Qwen3-30B exec×10"]
    end

    subgraph "Approach C: Council"
        C1["C1: Gemini chair<br/>+ [Qwen3,MiniMax,DeepSeek]"]
        C2["C2: Sonnet chair<br/>+ [Qwen3,MiniMax,Gemini]"]
    end

    subgraph "Approach D: Evolution"
        D1["D1: Gemini plan<br/>+ Qwen3-30B mutate"]
        D2["D2: Sonnet plan<br/>+ Qwen3-30B mutate"]
    end
```

| Config | Approach | Planner | Executor | Reviewer/Filter | Est Cost | Est Calls |
|--------|----------|---------|----------|----------------|----------|-----------|
| A1 | Quality-first | Sonnet | Sonnet | Sonnet | **free** | ~30 |
| A2 | Quality-first | Gemini | Qwen3 | Gemini | ~$0.40 | ~30 |
| A3 | Quality-first | Gemini | MiniMax | DeepSeek | ~$0.50 | ~30 |
| A4 | Quality-first | Sonnet | Qwen3-30B | Sonnet | ~$0.01 | ~30 |
| B1 | Gen+Filter | Gemini(×5) | Qwen3-30B(×10) | Tests | ~$1.00 | ~150 |
| B2 | Gen+Filter | Sonnet(×5) | Qwen3-30B(×10) | Tests | ~$0.10 | ~150 |
| C1 | Council | Gemini | [Q3,MM,DS] | Peer review | ~$0.70 | ~80 |
| C2 | Council | Sonnet | [Q3,MM,Gem] | Peer review | ~$0.50 | ~80 |
| D1 | Evolution | Gemini | Qwen3-30B | Tests | ~$0.30 | ~250 |
| D2 | Evolution | Sonnet | Qwen3-30B | Tests | ~$0.02 | ~250 |
| **Total** | | | | | **~$3.50** | **~1000** |

---

## Infrastructure

### Runner Architecture

```mermaid
flowchart LR
    Runner["run-experiment.sh<br/>--approach A|B|C|D<br/>--config A1|B1|..."]

    Runner --> CallModel["call-model.py<br/>Unified API caller<br/>OpenRouter or claude -p"]
    Runner --> ParseBlocks["parse-blocks.py<br/>Extract file blocks<br/>Write to disk"]
    Runner --> ValidateGate["validate-gate.sh<br/>Run gate command<br/>Return JSON"]
    Runner --> Results["results.tsv<br/>Append-only<br/>Per sub-task"]

    CallModel --> OpenRouter["OpenRouter API<br/>(Qwen3, Gemini, etc.)"]
    CallModel --> ClaudeP["claude -p<br/>(Sonnet, free)"]
```

### File Layout

```
scripts/dev-spike-v2/
├── README.md                    # This document
├── architecture.md              # Human-written spec (control variable)
├── run-experiment.sh            # Main runner
├── call-model.py                # Unified model caller
├── parse-blocks.py              # File block parser
├── validate-gate.sh             # Structural gate runner
├── results.tsv                  # All experiment data
├── prompts/
│   ├── planner.md               # Sub-task decomposition prompt
│   ├── executor.md              # File generation prompt
│   ├── reviewer.md              # Code review prompt
│   ├── test-writer.md           # Test-first prompt (Approach B)
│   ├── council-review.md        # Peer review prompt (Approach C)
│   ├── chairman.md              # Synthesis prompt (Approach C)
│   └── mutator.md               # "Improve this code" prompt (Approach D)
├── attempts/
│   ├── A1-001/
│   │   ├── plan.json
│   │   ├── ST-01/
│   │   │   ├── prompt.txt
│   │   │   ├── response.json
│   │   │   ├── gate-result.json
│   │   │   └── review.json
│   │   └── assembled/
│   │       └── dep-doctor/      # The final working app
│   └── ...
└── REPORT.md                    # Final analysis

```

### Results Tracking

Every API call produces one row in `results.tsv`:

| Field | Type | Description |
|-------|------|-------------|
| approach | A/B/C/D | Which approach |
| config | A1/B1/... | Which model config |
| attempt | 001 | Attempt number |
| layer | planner/executor/reviewer/... | Which layer |
| sub_task | ST-01 | Which sub-task |
| model | qwen/qwen3-coder | Model used |
| cost_usd | 0.0007 | Cost of this call |
| time_s | 3 | Wall clock seconds |
| tokens_in | 1200 | Input tokens |
| tokens_out | 890 | Output tokens |
| gate_pass | true/false | Structural gate result |
| review_verdict | keep/discard/fix/n/a | Reviewer verdict |
| quality_0_10 | 8 | Quality score |
| retry | 0/1/2 | Retry number |
| files_created | 1 | Files extracted from response |
| error | | Error message if failed |
| timestamp | ISO 8601 | When |

---

## Results & Analysis

*This section will be populated after experiments run.*

### Expected Outputs

#### Build Outcomes Table
| Config | Total Cost | Build Time | Sub-tasks Passed | Tests Passing | Final Working? |
|--------|-----------|-----------|------------------|---------------|----------------|

#### Model Performance by Layer
| Model | Planner Success% | Executor Gate Rate% | Reviewer Agreement% | Avg Cost/Call |
|-------|-----------------|--------------------|--------------------|---------------|

#### Approach Comparison
| Approach | Avg Cost | Success Rate | Best Config | Worst Config |
|----------|---------|-------------|-------------|-------------|

#### Cost Efficiency Frontier

```mermaid
graph LR
    subgraph "Cost vs Quality"
        D2["D2: $0.02<br/>? quality"]
        A4["A4: $0.01<br/>? quality"]
        B2["B2: $0.10<br/>? quality"]
        A2["A2: $0.40<br/>? quality"]
        C1["C1: $0.70<br/>? quality"]
        A1["A1: free<br/>? quality"]
    end
```

---

## Conclusions & Recommendations

*To be written after experiments complete.*

### Template

1. **Winning approach:** A/B/C/D with config X
2. **Best planner model:** ...
3. **Best executor model:** ...
4. **Is reviewer worth it:** yes/no (saves $X per bug caught)
5. **Minimum cost for working app:** $X.XX
6. **Recommended pipeline architecture:** ...
7. **Next steps for Dark Factory integration:** ...

---

## Appendix: Karpathy Research References

| Pattern | Source | How We Use It |
|---------|--------|--------------|
| Autoresearch loop | [karpathy/autoresearch](https://github.com/karpathy/autoresearch) | Time-boxed experiments, keep/discard, results.tsv |
| Verifiability principle | [Blog post](https://karpathy.bearblog.dev/verifiability/) | Structural gates as the quality filter |
| LLM Council | [karpathy/llm-council](https://github.com/karpathy/llm-council) | Approach C: peer review |
| Model tiers | [AnalyticsVidhya](https://www.analyticsvidhya.com/blog/2025/08/llm-workflow-for-developers/) | Opus plan, Sonnet execute, Qwen3 volume |
| AlphaCode | [DeepMind](https://deepmind.google/blog/competitive-programming-with-alphacode/) | Approach B: generate-and-filter |
| TDAD | [arXiv](https://arxiv.org/html/2603.17973) | Specific test targets, not procedural TDD |
| Agentic Engineering | [Karpathy 2026](https://thenewstack.io/vibe-coding-is-passe/) | Structured oversight, quality gates |
