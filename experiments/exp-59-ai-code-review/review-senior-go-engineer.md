# Senior Go Engineer Review

## Review Summary

This is a **REQUEST CHANGES** review. The code has several issues across multiple categories that need to be addressed before it can be considered production-ready. The main problems include critical concurrency issues, poor error handling, and anti-patterns in Go idioms.

---

### 1. **Go idioms** - **Major**
- **Issue**: The `genID()` function uses `string(rune(s.nextID))` which is incorrect and will produce invalid IDs.
- **Line**: 32
- **FIX**: Replace with a proper ID generation method like `strconv.Itoa(s.nextID)` or use UUIDs for better uniqueness.
- **VERDICT**: Major

### 2. **Error handling** - **Critical**
- **Issue**: No error handling for critical operations like `json.NewEncoder(w).Encode()` or `http.Redirect()`.
- **Line**: 157 (in `handleHealth`) and 193 (in `handleDashboard`)
- **FIX**: Check and log errors from JSON encoding and redirect operations.
- **VERDICT**: Critical

### 3. **Naming** - **Major**
- **Issue**: Function names like `handleHealth`, `handleDashboard` should be prefixed with `handle` to indicate they're HTTP handlers.
- **Line**: 143, 147, 151, 155, 160, 165, 170, 175, 180, 185
- **FIX**: Ensure all HTTP handler functions are consistently named with `handle` prefix.
- **VERDICT**: Major

### 4. **Function size** - **Major**
- **Issue**: `handleDashboard` function is over 50 lines and contains HTML generation logic.
- **Line**: ~147
- **FIX**: Split into smaller functions for HTML rendering and business logic.
- **VERDICT**: Major

### 5. **DRY** - **Major**
- **Issue**: Duplicated logic in `CreateClient`, `UpdateClient`, `CreateInvoice`, `UpdateInvoice` for setting ID and CreatedAt.
- **Line**: 42-44, 52-54, 67-69, 77-79
- **FIX**: Extract common logic into helper functions.
- **VERDICT**: Major

### 6. **Concurrency** - **Critical**
- **Issue**: The `genID()` method is not thread-safe due to lack of mutex protection.
- **Line**: 32
- **FIX**: Add mutex lock to `genID()` or use atomic operations.
- **VERDICT**: Critical

### 7. **Resource leaks** - **Minor**
- **Issue**: No explicit resource management for HTTP requests or responses.
- **Line**: 157, 193
- **FIX**: Ensure proper handling of request body and response writers.
- **VERDICT**: Minor

### 8. **Testing** - **Major**
- **Issue**: Code is tightly coupled with global state (`s` variable) and lacks testability.
- **Line**: 143, 157, 193
- **FIX**: Use dependency injection for the store and make it testable.
- **VERDICT**: Major

---

## VERDICT: REQUEST CHANGES

The code has several critical and major issues that prevent it from being production-ready. The primary concerns are:
1. **Concurrency bug** in ID generation
2. **Missing error handling**
3. **Poor separation of concerns**
4. **Lack of testability**

These issues must be addressed before the code can be approved.