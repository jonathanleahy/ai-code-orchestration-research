# Experiment 45: Chaos Agent

## Resilience testing — throw random garbage, see what survives

### Tests Run (25)
- Raw HTTP chaos (malformed requests, huge headers, NULL bytes)
- Concurrent chaos (10-50 simultaneous creates, create+delete race)
- Payload chaos (1MB body, 10000 fields, zalgo text, binary data)
- Rapid fire (100 GETs/POSTs in 1 second)
- Recovery check (still alive after all chaos)

### Results
- Survived: **25/25**
- Crashes: 0
- Server alive after chaos: **YES**

### Crashes
None — app is resilient!
