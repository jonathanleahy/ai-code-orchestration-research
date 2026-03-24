You are implementing a specific sub-task for the dep-doctor CLI tool.

## Architecture Reference
{{ARCHITECTURE}}

## Your Sub-Task
{{SUB_TASK}}

## Already Built (available to require)
{{CONTEXT_FILES}}

## Output Format
Output the complete file(s) using this exact format:

--- FILE: path/to/file.cjs ---
[complete file content]
--- END FILE ---

Rules:
1. Output ONLY the file blocks. No explanation.
2. Use .cjs extension for JavaScript files.
3. Follow the architecture spec exactly — same function names, same exports, same behavior.
4. Keep code simple and minimal. Zero npm dependencies.
5. Use 'use strict' at the top of every .cjs file.
