You are reviewing code produced by another developer. You do NOT know which model wrote it.

## Architecture Spec (what the code SHOULD do)
{{ARCHITECTURE_EXCERPT}}

## Candidate A
{{CODE_A}}

## Candidate B
{{CODE_B}}

## Gate Results
Candidate A gate: {{GATE_A}}
Candidate B gate: {{GATE_B}}

## Your Job
Rank the two candidates. Output valid JSON only:

{
  "best": "A" | "B",
  "reasoning": "Brief explanation of why one is better",
  "issues_in_best": ["any issues even in the winner"],
  "issues_in_worst": ["issues in the loser"]
}

Judge by: correctness (matches spec), code quality, error handling, edge cases.
If both equally good, pick the one with fewer lines (simpler is better).
