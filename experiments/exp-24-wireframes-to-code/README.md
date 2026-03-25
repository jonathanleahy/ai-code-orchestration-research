# Experiment 24: Wireframes → Code (Close the Loop)

## Hypothesis
Feeding Exp 22's wireframes to claude -p produces a better gateway than Exp 21's guessed UI.

## Setup
- Input: Exp 22's screen map + text wireframes (StatusPulse)
- Builder: claude -p Haiku (FREE on subscription)
- Target: Gateway with 6 routes matching the wireframe screens

## Results

| Metric | Guessed (Exp 21) | Wireframe-Driven (Exp 24) |
|--------|-----------------|--------------------------|
| Lines | 493 | 106 |
| Routes | 3 | **6** |
| Screens | 1 (dashboard only) | Dashboard + Add Check + Incidents + New Incident |
| Build | PASS | **PASS** |
| Cost | FREE | $0.20 (API) / FREE (sub) |

## Key Finding

**Wireframes produce more focused, better-structured code.**

- Exp 21 (guessed): 493 lines for 1 screen — bloated with inline CSS/JS
- Exp 24 (wireframe-driven): 106 lines for 6 routes — clean structure matching the wireframes

The wireframes give claude -p a concrete target instead of "make a status page." It builds exactly the screens specified, with the right navigation and forms.

## Cost
$0.20 (API) — this was a single claude -p call with a large prompt (wireframes + screen map).
On subscription: FREE.
