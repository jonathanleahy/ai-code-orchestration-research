# Experiment Map

## Research Overview

```mermaid
graph TD
    Core["🎯 Core Question<br/>Can AI build apps from scratch?<br/>What's the cheapest way?"]

    Core --> V1["📊 Spike V1<br/>Model Comparison<br/>11 models on bash task"]
    Core --> V2["🔧 Spike V2<br/>Node.js CLI (dep-doctor)<br/>4 approaches × 9 configs"]
    Core --> V3["🚀 Spike V3<br/>Go CRUD (task-board)<br/>HTTP server + HTML UI"]

    V1 --> V1R["Results:<br/>All models work at $0.008-$0.015<br/>claude -p tool mode was the bug"]
    V1 --> V1AR["🔬 Autoresearch V1<br/>Prompt optimization<br/>V4: 0% → 100% on cheapest model"]

    V2 --> V2A["Approach A: Quality-First<br/>🏆 A3: 18/18, $0.069<br/>A5: 18/18, $0.10"]
    V2 --> V2B["Approach B: Gen+Filter<br/>B1: 7/8, $0.237"]
    V2 --> V2C["Approach C: LLM Council<br/>C1: 6/8, $0.192"]
    V2 --> V2D["Approach D: Evolutionary<br/>D1: 7/8, $0.202"]

    V3 --> V3S4["🏆 S4: claude -p<br/>22/22 tests, FREE"]
    V3 --> V3S1["S1: Gemini+MiniMax<br/>10/10 model, $0.13"]
    V3 --> V3S2["S2: Gemini+Qwen3-30B<br/>10/10 model, $0.09"]
    V3 --> V3Gate["🔧 Compile Gate<br/>goimports + go build<br/>40% errors fixed free"]
    V3 --> V3AR["🔬 Autoresearch V2<br/>Executor prompt optimization<br/>Running..."]

    V3Gate --> V3Esc["📋 TODO: Escalation<br/>Cheap model fails → stronger model fixes"]
    V3AR --> V3Rerun["📋 TODO: Re-run V1+V2<br/>with improved prompts"]

    style Core fill:#4488ff,color:#fff
    style V2A fill:#22c55e,color:#000
    style V3S4 fill:#22c55e,color:#000
    style V1AR fill:#f59e0b,color:#000
    style V3AR fill:#f59e0b,color:#000
    style V3Gate fill:#8b5cf6,color:#fff
    style V3Esc fill:#fde68a,color:#000
    style V3Rerun fill:#fde68a,color:#000
```

## Experiment Index

| Spike | Application | Files | Tests | Report |
|-------|------------|-------|-------|--------|
| V1 | Bash task (--model flag) | 2 | 5 | [REPORT](spike-v1/REPORT.md) |
| V2 | Node.js CLI (dep-doctor) | 10 | 18 | [REPORT](spike-v2/REPORT.md) |
| V3 | Go CRUD (task-board) | 6 | 22 | [REPORT](spike-v3/REPORT.md) |

## Key Results Timeline

```mermaid
gantt
    title Research Timeline
    dateFormat HH:mm
    axisFormat %H:%M

    section Spike V1
    11 model comparison       :v1, 15:00, 60min
    Autoresearch prompts      :v1ar, 21:00, 15min

    section Spike V2
    Approach A (4 configs)    :v2a, 19:00, 30min
    Approach B (gen+filter)   :v2b, 19:30, 20min
    Approach C (council)      :v2c, 19:30, 20min
    Approach D (evolutionary) :v2d, 19:30, 20min
    A5 Qwen3-30B 18/18       :v2a5, 21:20, 5min

    section Spike V3
    Golden master             :v3gm, 21:30, 15min
    S4 claude -p 22/22       :v3s4, 22:00, 5min
    S1/S2 OpenRouter          :v3or, 22:00, 30min
    Compile gate + auto-fix   :v3fix, 22:30, 30min
    Autoresearch executor     :v3ar, 23:50, 30min
```

## Cost Summary

| Experiment | API Calls | OpenRouter Cost | Subscription Cost |
|-----------|-----------|----------------|-------------------|
| Spike V1 (11 models) | ~30 | $0.35 | — |
| Spike V2 (9 configs) | ~200 | $1.50 | — |
| Spike V2 autoresearch | ~24 | $0.05 | — |
| Spike V3 (3 configs) | ~30 | $0.30 | $0.50 (claude -p) |
| Spike V3 compile fix runs | ~20 | $0.25 | — |
| Spike V3 autoresearch | ~100 | $0.05 | — |
| **Total** | **~400** | **~$2.50** | **~$0.50** |

**Total research cost: ~$3.00**
