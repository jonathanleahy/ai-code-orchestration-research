## Issue 1
**ISSUE:** No validation or error handling for duplicate client entries - same name/email can be added multiple times
**SEVERITY:** Major
**UNIQUE:** YES

## Issue 2
**ISSUE:** Empty forms don't have clear "save" or "cancel" states - users can't distinguish between new/empty vs saved/loaded forms
**SEVERITY:** Major
**UNIQUE:** YES

## Issue 3
**ISSUE:** No delete cascade logic for client deletion - if a client is deleted, associated invoices and activities should be handled but spec doesn't define this
**SEVERITY:** Critical
**UNIQUE:** YES

## Issue 4
**ISSUE:** Activity timeline shows "chronological activity log" but doesn't specify what happens when multiple activities occur simultaneously or how to handle activity deletion
**SEVERITY:** Minor
**UNIQUE:** NO

## Issue 5
**ISSUE:** Invoice creation form has no validation for required fields or amount formatting - users can submit empty invoices
**SEVERITY:** Major
**UNIQUE:** YES

## Issue 6
**ISSUE:** Search functionality lacks "no results" state handling - users don't know if search failed or just found nothing
**SEVERITY:** Minor
**UNIQUE:** NO

## Issue 7
**ISSUE:** Dashboard quick stats don't specify what happens when there are no active clients or pending invoices (empty states)
**SEVERITY:** Minor
**UNIQUE:** NO

## Issue 8
**ISSUE:** No mention of data persistence for "quick add" buttons on dashboard - what happens if user navigates away?
**SEVERITY:** Major
**UNIQUE:** YES