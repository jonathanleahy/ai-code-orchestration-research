#!/usr/bin/env node
'use strict';

const { execSync } = require('child_process');
const path = require('path');
const fs = require('fs');

const CLI = path.join(__dirname, '..', 'cli.cjs');
const FIXTURES = path.join(__dirname, '..', 'fixtures');

let pass = 0;
let fail = 0;

function test(name, fn) {
  try {
    fn();
    console.log(`  PASS: ${name}`);
    pass++;
  } catch (e) {
    console.log(`  FAIL: ${name} — ${e.message}`);
    fail++;
  }
}

function assert(condition, message) {
  if (!condition) throw new Error(message || 'Assertion failed');
}

function run(args, expectFail = false) {
  try {
    const out = execSync(`node ${CLI} ${args}`, { encoding: 'utf8', timeout: 5000 });
    if (expectFail) throw new Error('Expected non-zero exit but got 0');
    return out;
  } catch (e) {
    if (!expectFail) throw e;
    return e.stderr || e.stdout || '';
  }
}

console.log('=== dep-doctor tests ===\n');

// --- CLI tests ---
test('--help exits 0 and lists commands', () => {
  const out = run('--help');
  assert(out.includes('scan'), 'missing scan command');
  assert(out.includes('report'), 'missing report command');
  assert(out.includes('check'), 'missing check command');
  assert(out.includes('licenses'), 'missing licenses command');
  assert(out.includes('init'), 'missing init command');
});

test('unknown command exits 1', () => {
  run('nonexistent', true);
});

// --- Parser tests ---
test('scan valid package.json returns deps', () => {
  const out = run(`scan --path ${FIXTURES}/valid`);
  const data = JSON.parse(out);
  assert(data.name === 'example-app', `expected example-app, got ${data.name}`);
  assert(data.deps.length === 5, `expected 5 deps, got ${data.deps.length}`);
});

test('scan malformed package.json fails', () => {
  run(`scan --path ${FIXTURES}/malformed`, true);
});

test('scan missing directory fails', () => {
  run('scan --path /nonexistent/path', true);
});

test('scan empty package.json returns 0 deps', () => {
  const out = run(`scan --path ${FIXTURES}/empty`);
  const data = JSON.parse(out);
  assert(data.deps.length === 0, `expected 0 deps, got ${data.deps.length}`);
});

// --- Validator tests ---
test('isValidSemver accepts ^1.2.3', () => {
  const { isValidSemver } = require('../lib/validator.cjs');
  assert(isValidSemver('^1.2.3') === true);
});

test('isValidSemver rejects not-semver', () => {
  const { isValidSemver } = require('../lib/validator.cjs');
  assert(isValidSemver('not-semver') === false);
});

test('isValidSpdx accepts MIT', () => {
  const { isValidSpdx } = require('../lib/validator.cjs');
  assert(isValidSpdx('MIT') === true);
});

test('isValidSpdx rejects unknown', () => {
  const { isValidSpdx } = require('../lib/validator.cjs');
  assert(isValidSpdx('UNKNOWN-LICENSE') === false);
});

// --- Analyzer tests ---
test('analyze flags deprecated packages', () => {
  const { parse } = require('../lib/parser.cjs');
  const { analyze } = require('../lib/analyzer.cjs');
  const parsed = parse(`${FIXTURES}/valid`);
  const analysis = analyze(parsed);
  const moment = analysis.dependencies.find(d => d.name === 'moment');
  assert(moment, 'moment not found');
  assert(moment.issues.some(i => i.type === 'deprecated'), 'moment not flagged as deprecated');
});

test('analyze reports summary counts', () => {
  const { parse } = require('../lib/parser.cjs');
  const { analyze } = require('../lib/analyzer.cjs');
  const parsed = parse(`${FIXTURES}/valid`);
  const analysis = analyze(parsed);
  assert(analysis.summary.total === 5, `expected 5 total, got ${analysis.summary.total}`);
  assert(analysis.summary.deprecated >= 1, 'expected at least 1 deprecated');
});

// --- Reporter tests ---
test('formatJson returns valid JSON', () => {
  const { parse } = require('../lib/parser.cjs');
  const { analyze } = require('../lib/analyzer.cjs');
  const { formatJson } = require('../lib/reporter.cjs');
  const analysis = analyze(parse(`${FIXTURES}/valid`));
  JSON.parse(formatJson(analysis)); // throws if invalid
});

test('formatTable returns string with headers', () => {
  const { parse } = require('../lib/parser.cjs');
  const { analyze } = require('../lib/analyzer.cjs');
  const { formatTable } = require('../lib/reporter.cjs');
  const table = formatTable(analyze(parse(`${FIXTURES}/valid`)));
  assert(table.includes('Name'), 'missing Name header');
  assert(table.includes('Range'), 'missing Range header');
});

test('formatSummary returns single line', () => {
  const { parse } = require('../lib/parser.cjs');
  const { analyze } = require('../lib/analyzer.cjs');
  const { formatSummary } = require('../lib/reporter.cjs');
  const summary = formatSummary(analyze(parse(`${FIXTURES}/valid`)));
  assert(!summary.includes('\n'), 'summary should be single line');
  assert(summary.includes('example-app'), 'missing app name');
});

// --- Check command ---
test('check exits 1 for project with issues', () => {
  run(`check --path ${FIXTURES}/valid`, true);
  // valid fixture has moment (deprecated) so should fail
});

test('check exits 0 for project with no issues', () => {
  const out = run(`check --path ${FIXTURES}/empty`);
  assert(out.includes('HEALTHY'), 'expected HEALTHY');
});

// --- Config tests ---
test('init creates config file', () => {
  const tmpDir = fs.mkdtempSync('/tmp/dep-doctor-test-');
  run(`init --path ${tmpDir}`);
  assert(fs.existsSync(path.join(tmpDir, '.dep-doctor.json')), 'config not created');
  const config = JSON.parse(fs.readFileSync(path.join(tmpDir, '.dep-doctor.json'), 'utf8'));
  assert(config.maxAge === 365, 'wrong default maxAge');
  fs.rmSync(tmpDir, { recursive: true });
});

console.log(`\n=== Results: ${pass} passed, ${fail} failed ===`);
process.exit(fail > 0 ? 1 : 0);
