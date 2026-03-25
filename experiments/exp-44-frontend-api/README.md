# Experiment 44: Frontend + API Security

## 37 checks, 27 pass, 10 fail

### By Category
| Category | Pass | Fail |
|----------|------|------|
| HTML Frontend | 15 | 1 |
| Security Headers | 3 | 6 |
| REST API | 8 | 3 |
| Stored XSS | 1 | 0 |

### Findings
- HTML: Missing lang attribute — accessibility
- Header: Missing — clickjacking risk
- Header: Missing — MIME sniffing risk
- Header: Missing — no HSTS (OK for localhost)
- Header: Missing — no CSP
- Header: Missing — may leak URLs
- Header: Missing — browser features not restricted
- API: No Cache-Control on API responses
- API: No rate limiting headers
- API: No pagination — will be slow with 1000+ clients
