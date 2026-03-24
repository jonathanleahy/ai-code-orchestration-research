You are implementing a specific sub-task for the dep-doctor CLI tool.

## Architecture Reference
{{ARCHITECTURE}}

## Your Sub-Task
{{SUB_TASK}}

## Already Built (available to require)
{{CONTEXT_FILES}}

## Output Format
You MUST use this EXACT format. Do NOT use markdown code fences. Do NOT add explanation.

--- FILE: path/to/file.cjs ---
'use strict';
// ... your code here ...
module.exports = { ... };
--- END FILE ---

CRITICAL RULES:
1. Start each file with: --- FILE: path/to/file.cjs ---
2. End each file with: --- END FILE ---
3. Do NOT wrap in ```javascript``` fences
4. Do NOT add any text before or after the file blocks
5. Use .cjs extension for all JavaScript files
6. ZERO npm dependencies — use only Node.js built-ins (fs, path, etc.)
7. Use 'use strict' at the top of every .cjs file
8. Follow the architecture spec exactly — same function names, same exports
