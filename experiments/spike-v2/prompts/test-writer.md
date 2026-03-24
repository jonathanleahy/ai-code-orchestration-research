You are writing tests for a Node.js CLI tool called dep-doctor.

## Architecture Reference
{{ARCHITECTURE}}

## Sub-Task: Write tests for this module
{{SUB_TASK}}

## Output Format
Output the COMPLETE test file using this exact format:

--- FILE: test/dep-doctor.test.cjs ---
[complete test file content]
--- END FILE ---

Rules:
1. Use .cjs extension (CommonJS)
2. Zero test framework — use console.log("  PASS: ...")/console.log("  FAIL: ...") pattern
3. Use assert from node built-ins or manual checks
4. Tests MUST fail if the implementation doesn't exist yet
5. Include at least 3 test cases per module
6. Output ONLY the file block. No explanation.
