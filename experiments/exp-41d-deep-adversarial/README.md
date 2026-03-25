# Experiment 41d: Deep Adversarial Testing (50 tests)

## 47 tests, 36 pass, 11 fail

### Categories
- Input validation (empty, whitespace, long, null bytes)
- XSS (script tags, img onerror, stored, reflected)
- SQL injection (in name, search, client ID)
- Path traversal (../, encoded, null byte)
- HTTP method abuse (PUT, PATCH, DELETE on wrong endpoints)
- Content type (JSON to form endpoint, XML, empty body)
- Special pages (/admin, /.env, /debug)
- Encoding (unicode, emoji, RTL, HTML entities)
- Concurrency (10 simultaneous creates)
- Header manipulation (Host injection, X-Forwarded-For)

### Failures
- FAIL: SQLi in search → 0
- FAIL: SQLi in client ID → 0
- FAIL: GET /admin → 200
- FAIL: GET /api/internal → 200
- FAIL: GET /.env → 200
- FAIL: GET /debug → 200
- FAIL: Unicode name → 400
- FAIL: Emoji name → 400
- FAIL: RTL text name → 400
- FAIL: HTML entities in name → 400
- FAIL: 404 page → 200
