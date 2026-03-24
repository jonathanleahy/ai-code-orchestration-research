'use strict';

const { spawn } = require('child_process');
const fs = require('fs');
const path = require('path');

const CLI = path.join(__dirname, '..', 'cli.cjs');
const { isValidSemver, isValidSpdx } = require('../lib/validator.cjs');
const { analyze } = require('../lib/analyzer.cjs');
const { formatJson, formatTable, formatSummary } = require('../lib/reporter.cjs');

// Test counters
let passed = 0;
let failed = 0;

function test(name, fn) {
  try {
    fn();
    console.log(`PASS: ${name}`);
    passed++;
  } catch (err) {
    console.log(`FAIL: ${name} — ${err.message}`);
    failed++;
  }
}

function assert(condition, message) {
  if (!condition) throw new Error(message || 'Assertion failed');
}

function assertEqual(actual, expected, message) {
  if (actual !== expected) {
    throw new Error(message || `Expected ${expected}, got ${actual}`);
  }
}

function runCli(args) {
  return new Promise((resolve) => {
    const proc = spawn('node', [CLI, ...args], { cwd: path.join(__dirname, '..') });
    let stdout = '';
    let stderr = '';
    proc.stdout.on('data', d => stdout += d);
    proc.stderr.on('data', d => stderr += d);
    proc.on('close', (code) => resolve({ code, stdout, stderr }));
  });
}

// Test 1: --help exits 0 and lists all 5 commands
test('1. --help exits 0 and lists all 5 commands', async () => {
  const result = await runCli(['--help']);
  assertEqual(result.code, 0, 'Expected exit 0');
  const commands = ['scan', 'report', 'check', 'licenses', 'init'];
  commands.forEach(cmd => {
    assert(result.stdout.includes(cmd), `Expected help to list command: ${cmd}`);
  });
});

// Test 2: unknown command exits 1
test('2. unknown command exits 1', async () => {
  const result = await runCli(['unknown-cmd']);
  assertEqual(result.code, 1, 'Expected exit 1');
});

// Test 3: scan valid package.json returns 5 deps
test('3. scan valid package.json returns 5 deps', async () => {
  const result = await runCli(['scan', '--path', 'fixtures/valid']);
  assertEqual(result.code, 0, 'Expected exit 0');
  const data = JSON.parse(result.stdout);
  assertEqual(data.dependencies.length, 5, 'Expected 5 dependencies');
});

// Test 4: scan malformed package.json fails
test('4. scan malformed package.json fails', async () => {
  const result = await runCli(['scan', '--path', 'fixtures/malformed']);
  assertEqual(result.code, 1, 'Expected exit 1 for malformed JSON');
});

// Test 5: scan missing directory fails
test('5. scan missing directory fails', async () => {
  const result = await runCli(['scan', '--path', 'fixtures/nonexistent']);
  assertEqual(result.code, 1, 'Expected exit 1 for missing directory');
});

// Test 6: scan empty package.json returns 0 deps
test('6. scan empty package.json returns 0 deps', async () => {
  const result = await runCli(['scan', '--path', 'fixtures/empty']);
  assertEqual(result.code, 0, 'Expected exit 0');
  const data = JSON.parse(result.stdout);
  assertEqual(data.dependencies.length, 0, 'Expected 0 dependencies');
});

// Test 7: isValidSemver accepts ^1.2.3
test('7. isValidSemver accepts ^1.2.3', () => {
  assert(isValidSemver('^1.2.3') === true, 'Should accept ^1.2.3');
});

// Test 8: isValidSemver rejects not-semver
test('8. isValidSemver rejects not-semver', () => {
  assert(isValidSemver('not-semver') === false, 'Should reject not-semver');
});

// Test 9: isValidSpdx accepts MIT
test('9. isValidSpdx accepts MIT', () => {
  assert(isValidSpdx('MIT') === true, 'Should accept MIT');
});

// Test 10: isValidSpdx rejects unknown license
test('10. isValidSpdx rejects unknown license', () => {
  assert(isValidSpdx('UnknownLicense') === false, 'Should reject unknown license');
});

// Test 11: analyze flags deprecated packages (moment)
test('11. analyze flags deprecated packages (moment)', () => {
  const parsed = {
    name: 'test-app',
    version: '1.0.0',
    license: 'MIT',
    deps: [
      { name: 'moment', range: '^2.29.4', type: 'dependency' }
    ]
  };
  const result = analyze(parsed);
  const moment = result.dependencies.find(d => d.name === 'moment');
  assert(moment, 'Should find moment dependency');
  assert(moment.issues.some(i => i.message.includes('deprecated')), 'Should flag moment as deprecated');
});

// Test 12: analyze reports correct summary counts
test('12. analyze reports correct summary counts', () => {
  const parsed = {
    name: 'test-app',
    version: '1.0.0',
    license: 'MIT',
    deps: [
      { name: 'lodash', range: '^4.17.21', type: 'dependency' },
      { name: 'moment', range: '^2.29.4', type: 'dependency' },
      { name: 'express', range: '~4.18.2', type: 'dependency' }
    ]
  };
  const result = analyze(parsed);
  assertEqual(result.summary.total, 3, 'Should have 3 total deps');
  assertEqual(result.summary.deprecated, 1, 'Should have 1 deprecated (moment)');
  assertEqual(result.summary.large, 2, 'Should have 2 large (lodash, moment)');
});

// Test 13: formatJson returns valid JSON
test('13. formatJson returns valid JSON', () => {
  const analysis = {
    name: 'test',
    version: '1.0.0',
    license: 'MIT',
    licenseIssues: [],
    dependencies: [],
    summary: { total: 0, healthy: 0, issues: 0, deprecated: 0, large: 0 }
  };
  const output = formatJson(analysis);
  const parsed = JSON.parse(output);
  assertEqual(parsed.name, 'test', 'Should parse JSON correctly');
});

// Test 14: formatTable has column headers
test('14. formatTable has column headers', () => {
  const analysis = {
    name: 'test',
    version: '1.0.0',
    license: 'MIT',
    licenseIssues: [],
    dependencies: [],
    summary: { total: 0, healthy: 0, issues: 0, deprecated: 0, large: 0 }
  };
  const output = formatTable(analysis);
  assert(output.includes('Name'), 'Should have Name header');
  assert(output.includes('Range'), 'Should have Range header');
  assert(output.includes('Type'), 'Should have Type header');
  assert(output.includes('Issues'), 'Should have Issues header');
});

// Test 15: formatSummary is single line with app name
test('15. formatSummary is single line with app name', () => {
  const analysis = {
    name: 'my-app',
    version: '2.0.0',
    license: 'MIT',
    licenseIssues: [],
    dependencies: [],
    summary: { total: 5, healthy: 5, issues: 0, deprecated: 0, large: 0 }
  };
  const output = formatSummary(analysis);
  const lines = output.split('\n');
  assertEqual(lines.length, 1, 'Should be single line');
  assert(output.includes('my-app'), 'Should include app name');
});

// Test 16: check exits 1 for project with issues
test('16. check exits 1 for project with issues', async () => {
  const result = await runCli(['check', '--path', 'fixtures/valid']);
  assertEqual(result.code, 1, 'Expected exit 1 for project with issues (moment deprecated)');
});

// Test 17: check exits 0 for healthy project
test('17. check exits 0 for healthy project', async () => {
  const result = await runCli(['check', '--path', 'fixtures/empty']);
  assertEqual(result.code, 0, 'Expected exit 0 for healthy project');
});

// Test 18: init creates config file with defaults
test('18. init creates config file with defaults', async () => {
  const configPath = path.join(__dirname, '..', '.dep-doctor.json');
  if (fs.existsSync(configPath)) {
    fs.unlinkSync(configPath);
  }
  const result = await runCli(['init']);
  assertEqual(result.code, 0, 'Expected exit 0');
  assert(fs.existsSync(configPath), 'Config file should be created');
  const config = JSON.parse(fs.readFileSync(configPath, 'utf8'));
  assert(Array.isArray(config.deprecated), 'Should have deprecated array');
  assert(Array.isArray(config.large), 'Should have large array');
  assertEqual(config.exitCode, 1, 'Should have exitCode');
  assertEqual(config.spdx, true, 'Should have spdx flag');
  if (fs.existsSync(configPath)) {
    fs.unlinkSync(configPath);
  }
});

// Summary
console.log('');
console.log(`${passed} passed, ${failed} failed`);

process.exit(failed > 0 ? 1 : 0);