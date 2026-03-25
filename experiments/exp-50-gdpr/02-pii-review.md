# PII Code Review

### 1. **PII in logs**: Does `fmt.Println`/`log.Println` output emails, names, phones?

**Answer: LEAK**  
The code uses `log.Println("Server listening on :8080")` but does not log PII directly. However, it's possible that in a full version of the application, logs could contain PII if not carefully managed.

**Line numbers:**
- Line 20: `log.Println("Server listening on :8080")`  
  *Note: This is not PII, but logging server startup info can be risky in some environments.*

✅ **Safe** for this snippet — no PII logged.

**Result: NOT APPLICABLE** (since no PII is logged in the current code)

---

### 2. **PII in error messages**: `"Client john@example.com not found"` leaks email

**Answer: LEAK**

There is no explicit error message like that in the provided code. But if an error handler or similar function were to return such a message, it would be a leak.

✅ **Not present in current code**, but **could be a risk** in extended versions.

**Result: NOT APPLICABLE**

---

### 3. **PII in URLs**: GET `/client?email=john@...` puts PII in server logs

**Answer: LEAK**

This is a potential risk if a route like `/client?email=...` is used, but the current code doesn't do that. It uses `/client/{id}`.

However, if a GET parameter like `?email=...` is passed to any route, it would be visible in logs.

✅ **Not present in current code**, but **could be a risk** in extended versions.

**Result: NOT APPLICABLE**

---

### 4. **PII in HTML comments**: `<!-- client data: ... -->`

**Answer: LEAK**

The HTML is constructed dynamically in Go and does not include comments with PII.

✅ **Not present in current code**

**Result: NOT APPLICABLE**

---

### 5. **Unencrypted PII storage**: Are emails/phones in plain text in memory?

**Answer: LEAK**

The `store.Client` struct holds emails and names in plain text in memory (no encryption). This is a **data at rest** risk if the memory is not protected.

✅ **Yes**, PII is stored in plain text in memory.

**Result: LEAK**

---

### 6. **PII in API responses**: Does list endpoint return full objects or summaries?

**Answer: LEAK**

The `/` endpoint returns a list of clients with full details (name, email). If this were an API, returning full client objects with PII would be a leak unless filtered.

✅ **Yes**, full client data is returned in HTML (not JSON), but still exposes PII.

**Result: LEAK**

---

### 7. **PII in debug output**: Any `fmt.Printf` with client data?

**Answer: LEAK**

There are no `fmt.Printf` calls in the provided code that output PII.

✅ **No PII printed directly**

**Result: NOT APPLICABLE**

---

### 8. **PII retention**: Is old data ever cleaned up?

**Answer: LEAK**

There is no cleanup logic in the provided code. If data is stored long-term, it may retain PII unnecessarily.

✅ **No cleanup logic**

**Result: LEAK**

---

### 9. **PII in exported data**: Are backups/exports encrypted?

**Answer: NOT APPLICABLE**

There is no export or backup logic in the code.

**Result: NOT APPLICABLE**

---

### 10. **PII access control**: Can any request access any client's PII?

**Answer: LEAK**

The `/client/{id}` route is used to access individual clients, but there's no authentication or authorization logic in the provided code. Any user could access any client's data if they know the ID.

✅ **No access control**

**Result: LEAK**

---

## ✅ Summary Table

| Check | Result | Line Numbers |
|-------|--------|--------------|
| 1. PII in logs | NOT APPLICABLE | 20 |
| 2. PII in error messages | NOT APPLICABLE | N/A |
| 3. PII in URLs | NOT APPLICABLE | N/A |
| 4. PII in HTML comments | NOT APPLICABLE | N/A |
| 5. Unencrypted PII storage | LEAK | N/A (in-memory) |
| 6. PII in API responses | LEAK | N/A (HTML rendering) |
| 7. PII in debug output | NOT APPLICABLE | N/A |
| 8. PII retention | LEAK | N/A |
| 9. PII in exported data | NOT APPLICABLE | N/A |
| 10. PII access control | LEAK | N/A |

---

## 🔐 Recommendations

1. **Encrypt sensitive data in memory** (e.g., using `crypto/aes` or similar).
2. **Avoid exposing PII in HTML or API responses** unless strictly necessary.
3. **Implement access control** (e.g., authentication/authorization).
4. **Avoid logging PII** in any form.
5. **Implement data retention policies** to delete old PII.
6. **Sanitize all inputs and outputs** to prevent accidental exposure.
7. **Use secure headers** and HTTPS to protect data in transit.

Let me know if you'd like a revised version of this code with PII protections applied.