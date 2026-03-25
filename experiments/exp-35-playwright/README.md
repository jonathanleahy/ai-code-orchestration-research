# Experiment 35: Playwright Journey Testing

## What We Tested
Ran Playwright against the simplified CRM (Exp 34) to verify user journeys.

## Results

| Journey | Result | Notes |
|---------|--------|-------|
| Add Client | PASS | Name + email only, no phone/address fields |
| View Client Profile | PASS | Shows activities + invoices tabs |
| Invoices Tab | PASS | Shows invoice list with amounts |
| **Edit Client** | **FAIL** | No edit functionality |
| **Add Activity** | **FAIL** | No add activity form |
| **Add Invoice** | **FAIL** | No create invoice form |
| **Delete Client** | **FAIL** | No delete button |
| **Search** | **FAIL** | No search box |

## Key Finding

**The simplicity agent cut features that were in the brief.**

The brief said "client management, activity history, and invoicing" but the simplified app only has:
- Add client (name + email only)
- View client with pre-seeded data
- View activities and invoices (read-only)

Missing: edit, delete, search, add activity, create invoice, phone field, address.

## The Root Cause

The simplicity agent reviewed at the **feature level** ("cut phone field") instead of the **implementation level** ("use one textarea instead of 5 fields for address").

## The Fix

Constrain the simplicity agent:
> "You CANNOT remove features from the brief. You CAN simplify how they're built."

| Feature | Bad Simplification | Good Simplification |
|---------|-------------------|---------------------|
| Address | Cut it entirely | One textarea instead of 5 fields |
| Edit client | Cut it | Reuse add form with pre-filled values |
| Invoice print | Cut it | Browser print button (CSS @media print) |
| Search | Cut it | Simple `strings.Contains` filter |

## Comparison

| Metric | Exp 32 (no simplicity) | Exp 34 (over-simplified) | Target |
|--------|----------------------|-------------------------|--------|
| Store | 380 lines | 173 lines | ~200 lines |
| Server | 1136 lines | 345 lines | ~500 lines |
| Features | All 10 | Missing 5 | All 10, simply |
| Cost | $0.96 | $0.20 | ~$0.40 |
