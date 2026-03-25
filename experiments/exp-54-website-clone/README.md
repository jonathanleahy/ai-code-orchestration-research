# Experiment 54: Website Clone — eachandother.com

## From live site screenshot to working clone

### Process
1. Playwright navigates to eachandother.com
2. Screenshot captured + DOM snapshot analyzed
3. Exp 53 extracts design system (tokens + components)
4. claude -p builds Go server with embedded HTML matching the layout

### Results
- **622 lines** of Go (server + embedded HTML)
- **13.6KB** HTML output
- **BUILD PASS**
- Matches: nav, hero, logo bar, services, case studies, CTA, blog, footer

### The Clone Includes
- Nav with logo + service links
- Hero: "Competitive Advantage by Design." (pink accent)
- Logo bar: Google, Zurich, Stripe, Coinbase, BNP Paribas
- 3 service items with arrows
- "How we do it" section with 3 columns
- Case studies section
- Blog cards
- Multi-column footer with addresses
- Pink #d4727a accent throughout

### Cost
~$0.20 (claude -p subscription)
