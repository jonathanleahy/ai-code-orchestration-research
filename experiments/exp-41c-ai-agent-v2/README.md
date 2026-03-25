# Experiment 41c: AI Testing Agent v2

## Fixed: structured JSON extraction with few-shot examples

25 actions executed, 7 unique pages visited.

### Bugs Found: 0
None

### Actions
- POST /api/clients → 201 | invalid email format should be rejected
- POST /api/clients → 201 | test invalid phone format validation
- POST /api/clients → 201 | test for XSS vulnerability in custom field
- POST /api/clients → 201 | invalid date format should be rejected
- POST /api/clients → 201 | invalid phone format should be rejected
- POST /api/clients → 201 | invalid date format in created_at field
- DELETE /api/clients/20260325183834-0010 → 200 | delete newly created client to test if it was actu
- POST /api/clients → 201 | submit duplicate client to test for duplicate dete
- DELETE /api/clients/20260325183834-0010 → 200 | delete a client that was just created
- POST /api/clients → 201 | submit duplicate client data to test for duplicate
- DELETE /api/clients/20260325183834-0010 → 200 | delete existing client record
- POST /api/clients → 201 | submitting duplicate client data
- DELETE /api/clients/20260325183843-0013 → 200 | delete recently created client to test if it's act
- POST /api/clients → 201 | submit duplicate client data to test for duplicate
- DELETE /api/clients/20260325183834-0010 → 200 | delete client that was just created to test idempo
- POST /api/clients → 201 | submit duplicate client data to test for duplicate
- DELETE /api/clients/20260325183849-0015 → 200 | delete recently created client to test if it's act
- POST /api/clients → 201 | submitting duplicate client data
- DELETE /api/clients/20260325183853-0016 → 200 | delete recently created user to test if it's actua
- POST /api/clients → 201 | submit duplicate client data
- DELETE /api/clients/20260325183857-0017 → 200 | delete recently created client to test idempotency
- POST /api/clients → 201 | submit duplicate client data to test for duplicate
- DELETE /api/clients/20260325183843-0013 → 200 | delete a client that was created earlier
- POST /api/clients → 201 | submit duplicate client data to test for duplicate
- DELETE /api/clients/20260325183908-0019 → 200 | delete recently created user to test if it's actua

### Cost: $0.0345
