# dep-doctor — Architecture Specification

## Overview
A Node.js CLI tool that reads `package.json`, analyzes dependency health, and outputs structured reports. Zero npm dependencies — uses only Node.js built-ins.

**All files use `.cjs` extension** (CommonJS) because the parent project has `"type": "module"` in package.json.

## File Structure

```
dep-doctor/
  cli.cjs                          # Entry point (80 lines)
  lib/validator.cjs                # Validation functions (60 lines)
  lib/parser.cjs                   # Parse package.json (60 lines)
  lib/config.cjs                   # Config file I/O (40 lines)
  lib/analyzer.cjs                 # Health analysis (100 lines)
  lib/reporter.cjs                 # Output formatting (80 lines)
  test/dep-doctor.test.cjs         # 18 test cases (120 lines)
  fixtures/valid/package.json      # Valid test data
  fixtures/malformed/package.json  # Invalid JSON test data
  fixtures/empty/package.json      # Empty deps test data
```

## Dependency Order (build in this order)

1. `fixtures/` — test data (no code deps)
2. `lib/validator.cjs` — standalone (no local deps)
3. `lib/parser.cjs` — depends on `fs`, `path` (Node built-ins only)
4. `lib/config.cjs` — depends on `fs`, `path` (Node built-ins only)
5. `lib/analyzer.cjs` — depends on `lib/validator.cjs`
6. `lib/reporter.cjs` — standalone (formats data)
7. `cli.cjs` — depends on all lib modules
8. `test/dep-doctor.test.cjs` — depends on all above

## Module Specifications

### fixtures/valid/package.json
```json
{
  "name": "example-app",
  "version": "1.0.0",
  "description": "A sample application for testing",
  "license": "MIT",
  "dependencies": {
    "express": "^4.18.2",
    "lodash": "^4.17.21",
    "moment": "^2.29.4"
  },
  "devDependencies": {
    "jest": "^29.7.0",
    "eslint": "^8.56.0"
  }
}
```

### fixtures/malformed/package.json
```json
{
  "name": "broken-app"
  "version": "not-semver",
  "dependencies": {
    "fake-pkg": "latest"
  }
```
(Missing comma after "broken-app" — intentionally invalid JSON)

### fixtures/empty/package.json
```json
{
  "name": "empty-app",
  "version": "0.0.1"
}
```

### lib/validator.cjs

**Exports:**
- `isValidSemver(range: string): boolean` — Returns true for valid semver ranges (^1.2.3, ~2.0.0, >=1.0.0, *, latest)
- `isValidSpdx(license: string): boolean` — Returns true for known SPDX identifiers (MIT, ISC, BSD-2-Clause, BSD-3-Clause, Apache-2.0, GPL-2.0-only, GPL-3.0-only, LGPL-2.1-only, MPL-2.0, UNLICENSED)
- `validateDependency(name: string, range: string): Array<{type, message}>` — Returns array of issues

**Validation gate:** `node -e "const v=require('./lib/validator.cjs'); if(typeof v.isValidSemver!=='function') process.exit(1)"`

### lib/parser.cjs

**Exports:**
- `parse(dirPath: string): {name, version, license, deps: [{name, range, type}]}` — Reads package.json from dirPath, returns parsed structure. Throws with `code='ENOENT'` if file missing, `code='INVALID_JSON'` if malformed.

**Validation gate:** `node -e "const p=require('./lib/parser.cjs'); if(typeof p.parse!=='function') process.exit(1)"`

### lib/config.cjs

**Exports:**
- `loadConfig(dirPath: string): {maxAge, denyLicenses, ignorePackages, format}` — Reads .dep-doctor.json or returns defaults
- `saveConfig(dirPath: string, config: object): void` — Writes .dep-doctor.json
- `DEFAULTS` — `{maxAge: 365, denyLicenses: [], ignorePackages: [], format: 'summary'}`
- `CONFIG_FILE` — `'.dep-doctor.json'`

**Validation gate:** `node -e "const c=require('./lib/config.cjs'); if(typeof c.loadConfig!=='function') process.exit(1)"`

### lib/analyzer.cjs

**Exports:**
- `analyze(parsed): {name, version, license, licenseIssues, dependencies: [{name, range, type, issues, healthy}], summary: {total, healthy, issues, deprecated, large}}` — Analyzes parsed deps for issues
- `DEPRECATED` — Set of known deprecated packages: `request, moment, tslint, bower`
- `LARGE` — Set of known large packages: `moment, lodash, rxjs, core-js`

**Depends on:** `require('./validator.cjs')`

**Validation gate:** `node -e "const a=require('./lib/analyzer.cjs'); if(typeof a.analyze!=='function') process.exit(1)"`

### lib/reporter.cjs

**Exports:**
- `formatJson(analysis): string` — JSON.stringify with 2-space indent
- `formatTable(analysis): string` — ASCII table with Name/Range/Type/Issues columns
- `formatSummary(analysis): string` — Single line: `name@version: STATUS (N deps, N issues, N deprecated, N large)`

**Validation gate:** `node -e "const r=require('./lib/reporter.cjs'); if(typeof r.formatJson!=='function') process.exit(1)"`

### cli.cjs

**Subcommands:**
1. `scan --path <dir>` — Output parsed deps as JSON
2. `report --path <dir> [--format json|table|summary]` — Full analysis report
3. `check --path <dir>` — Exit 0 if healthy, exit 1 if issues found
4. `licenses --path <dir>` — List license + all deps
5. `init --path <dir>` — Create .dep-doctor.json with defaults
6. `--help` — Show usage

**Depends on:** All lib modules

**Validation gate:** `node cli.cjs --help` exits 0 AND `node cli.cjs nonexistent` exits 1

### test/dep-doctor.test.cjs

**18 test cases** (zero test framework — uses console.log("PASS/FAIL") pattern):

1. --help exits 0 and lists all 5 commands
2. unknown command exits 1
3. scan valid package.json returns 5 deps
4. scan malformed package.json fails
5. scan missing directory fails
6. scan empty package.json returns 0 deps
7. isValidSemver accepts ^1.2.3
8. isValidSemver rejects not-semver
9. isValidSpdx accepts MIT
10. isValidSpdx rejects unknown license
11. analyze flags deprecated packages (moment)
12. analyze reports correct summary counts
13. formatJson returns valid JSON
14. formatTable has column headers
15. formatSummary is single line with app name
16. check exits 1 for project with issues
17. check exits 0 for healthy project
18. init creates config file with defaults

**Validation gate:** `node test/dep-doctor.test.cjs` exits 0, stdout contains "18 passed, 0 failed"

## Golden Master Outputs

The reference implementation produces these exact outputs:

### `node cli.cjs --help` (exit 0)
```
dep-doctor — Dependency Health Checker
...
Commands:
  scan       Scan a directory's package.json and list dependencies
  report     Full health report with analysis
  check      Pass/fail gate (exit 0 = healthy, exit 1 = issues)
  licenses   License inventory for all dependencies
  init       Create a .dep-doctor.json config file
```

### `node cli.cjs scan --path fixtures/valid` (exit 0, valid JSON)
```json
{
  "name": "example-app",
  "version": "1.0.0",
  "license": "MIT",
  "deps": [
    {"name": "express", "range": "^4.18.2", "type": "production"},
    {"name": "lodash", "range": "^4.17.21", "type": "production"},
    {"name": "moment", "range": "^2.29.4", "type": "production"},
    {"name": "jest", "range": "^29.7.0", "type": "dev"},
    {"name": "eslint", "range": "^8.56.0", "type": "dev"}
  ]
}
```

### `node cli.cjs check --path fixtures/empty` (exit 0)
```
empty-app@0.0.1: HEALTHY (0 deps, 0 issues, 0 deprecated, 0 large)
```

### `node test/dep-doctor.test.cjs` (exit 0)
```
=== dep-doctor tests ===
  PASS: --help exits 0 and lists commands
  PASS: unknown command exits 1
  ... (18 PASS lines)
=== Results: 18 passed, 0 failed ===
```
