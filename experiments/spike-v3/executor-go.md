You are implementing a Go file for the task-board project.

## Architecture Reference
{{ARCHITECTURE}}

## Your Sub-Task
{{SUB_TASK}}

## Already Built Files
{{CONTEXT_FILES}}

## YOUR OUTPUT MUST START AND END EXACTLY LIKE THIS:

--- FILE: {{PRIMARY_FILE}} ---
package model

import (
	"fmt"
	"sync"
	"time"
)

// ... your implementation ...
--- END FILE ---

CRITICAL RULES:
1. Start with: --- FILE: {{PRIMARY_FILE}} ---
2. End with: --- END FILE ---
3. Do NOT wrap in ```go``` fences
4. Do NOT add text before or after the file block
5. The file must be COMPLETE and compilable — all types, all methods, all imports
6. Use standard Go conventions (gofmt style)
7. Only use Go standard library — no external packages
8. Match the architecture spec EXACTLY — same type names, same method signatures
