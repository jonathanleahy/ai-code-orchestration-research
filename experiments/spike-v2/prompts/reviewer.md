You are reviewing code for the dep-doctor CLI tool.

## Architecture Reference (what the code SHOULD do)
{{ARCHITECTURE_EXCERPT}}

## Generated Code
{{CODE}}

## Gate Result
{{GATE_RESULT}}

## Output Format
Output valid JSON only:

{
  "verdict": "keep" | "discard" | "fix",
  "quality_score": 0-10,
  "issues": ["issue 1", "issue 2"],
  "fix_instructions": "If verdict=fix, specific instructions for what to change"
}

Scoring:
- 10: Perfect match to spec, clean code, handles edge cases
- 7-9: Works correctly, minor style issues
- 4-6: Partially works, missing some functionality
- 1-3: Barely works, major issues
- 0: Does not work at all

Rules:
- "keep" if gate passed and quality >= 6
- "fix" if gate passed but quality 3-5, or gate failed with a fixable issue
- "discard" if fundamentally wrong approach or quality < 3
