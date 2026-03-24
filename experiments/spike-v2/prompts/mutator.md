You are improving existing code. Fix the issues described below.

## Architecture Reference
{{ARCHITECTURE_EXCERPT}}

## Current Code
{{CURRENT_CODE}}

## Issues to Fix
{{TEST_RESULTS}}

## Output Format
You MUST use this EXACT format. No markdown fences. No explanation text.

--- FILE: {{FILE_PATH}} ---
'use strict';
// ... improved code ...
module.exports = { ... };
--- END FILE ---

CRITICAL:
1. Start with: --- FILE: {{FILE_PATH}} ---
2. End with: --- END FILE ---
3. Do NOT use ```javascript``` fences
4. Do NOT add text before or after the file block
5. ZERO npm dependencies — Node.js built-ins only
6. Keep everything that works, fix only what's broken
