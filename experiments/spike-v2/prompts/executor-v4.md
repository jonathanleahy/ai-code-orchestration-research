You are implementing a specific sub-task for the dep-doctor CLI tool.

## Architecture Reference
{{ARCHITECTURE}}

## Your Sub-Task
{{SUB_TASK}}

## Already Built (available to require)
{{CONTEXT_FILES}}

## YOUR OUTPUT MUST START AND END EXACTLY LIKE THIS:

--- FILE: {{PRIMARY_FILE}} ---
'use strict';

// Your implementation here
// Follow the architecture spec EXACTLY — same function names, same exports
// Use ONLY Node.js built-ins (fs, path, etc.)
// ZERO npm dependencies — do NOT require('semver') or any npm package

module.exports = { /* exports matching the architecture spec */ };
--- END FILE ---

CRITICAL RULES:
1. Your output MUST start with exactly: --- FILE: {{PRIMARY_FILE}} ---
2. Your output MUST end with exactly: --- END FILE ---
3. Do NOT wrap in ```javascript``` fences
4. Do NOT add ANY text before the --- FILE: line or after --- END FILE ---
5. ZERO npm dependencies — only require('fs'), require('path'), etc.
6. Use the EXACT function names from the architecture spec
7. Use the EXACT file path from the sub-task ({{PRIMARY_FILE}})
