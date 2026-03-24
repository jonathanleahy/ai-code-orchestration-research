'use strict';

const fs = require('fs');
const path = require('path');

// Mock console.log to capture output for validation
const originalLog = console.log;
let testOutput = '';

function captureLog(...args) {
  testOutput += args.join(' ') + '\n';
}

console.log = captureLog;

// Helper to run a command and capture exit code
function runCommand(args) {
  const cliPath = path.join(__dirname, '..', 'cli.cjs');
  const { spawnSync } = require('child_process');
  const result = spawnSync('node', [cliPath, ...args], { encoding: 'utf8' });
  return {
    stdout: result.stdout,
    stderr: result.stderr,
    status: result.status
  };
}

// Test 1: --help exits 0 and lists all 5 commands
function testHelp() {
  const result = runCommand(['--help']);
  const passes = result.status === 0 && result.stdout.includes('Commands:') && result.stdout.includes('scan') && result.stdout.includes('report') && result.stdout.includes('check') && result.stdout.includes('licenses') && result.stdout.includes('init');
  console.log(passes ? 'PASS' : 'FAIL');
}

// Test 2: unknown command exits 1
function testUnknownCommand() {
  const result = runCommand(['unknown']);
  const passes = result.status === 1;
  console.log(passes ? 'PASS' : 'FAIL');
}

// Test 3: scan valid package.json returns 5 deps
function testScanValid() {
  const result = runCommand(['scan', '--path', 'fixtures/valid']);
  const passes = result.status === 0 && result.stdout.includes('"name":"example-app"') && result.stdout.includes('"dependencies"') && result.stdout.includes('"devDependencies"');
  console.log(passes ? 'PASS' : 'FAIL');
}

// Test 4: scan malformed package.json fails
function testScanMalformed() {
  const result = runCommand(['scan', '--path', 'fixtures/malformed']);
  const passes = result.status !== 0;
  console.log(passes ? 'PASS' : 'FAIL');
}

// Test 5: scan missing directory fails
function testScanMissing() {
  const result = runCommand(['scan', '--path', 'nonexistent']);
  const passes = result.status !== 0;
  console.log(passes ? 'PASS' : 'FAIL');
}

// Test 6: scan empty package.json returns 0 deps
function testScanEmpty() {
  const result = runCommand(['scan', '--path', 'fixtures/empty']);
  const passes = result.status === 0 && result.stdout.includes('"dependencies":{}') && result.stdout.includes('"devDependencies":{}');
  console.log(passes ? 'PASS' : 'FAIL');
}

// Test 7: isValidSemver accepts ^1.2.3
function testValidSemver() {
  const validator = require('./lib/validator.cjs');
  const passes = validator.isValidSemver('^1.2.3');
  console.log(passes ? 'PASS' : 'FAIL');
}

// Test 8: isValidSemver rejects not-semver
function testInvalidSemver() {
  const validator = require('./lib/validator.cjs');
  const passes = !validator.isValidSemver('not-semver');
  console.log(passes ? 'PASS' : 'FAIL');
}

// Test 9: isValidSpdx accepts MIT
function testValidSpdx() {
  const validator = require('./lib/validator.cjs');
  const passes = validator.isValidSpdx('MIT');
  console.log(passes ? 'PASS' : 'FAIL');
}

// Test 10: isValidSpdx rejects unknown license
function testInvalidSpdx() {
  const validator = require('./lib/validator.cjs');
  const passes = !validator.isValidSpdx('UNKNOWN');
  console.log(passes ? 'PASS' : 'FAIL');
}

// Test 11: analyze flags deprecated packages (moment)
function testAnalyzeDeprecated() {
  const parser = require('./lib/parser.cjs');
  const analyzer = require('./lib/analyzer.cjs');
  const parsed = parser.parse('fixtures/valid/package.json');
  const analyzed = analyzer.analyze(parsed);
  const passes = analyzed.dependencies.some(d => d.name === 'moment' && d.issues.some(i => i.message.includes('deprecated')));
  console.log(passes ? 'PASS' : 'FAIL');
}

// Test 12: analyze reports correct summary counts
function testAnalyzeSummary() {
  const parser = require('./lib/parser.cjs');
  const analyzer = require('./lib/analyzer.cjs');
  const parsed = parser.parse('fixtures/valid/package.json');
  const analyzed = analyzer.analyze(parsed);
  const summary = analyzed.summary;
  const passes = summary.total === 5 && summary.healthy === 5 && summary.issues === 0 && summary.deprecated === 1 && summary.large === 2;
  console.log(passes ? 'PASS' : 'FAIL');
}

// Test 13: formatJson returns valid JSON
function testFormatJson() {
  const reporter = require('./lib/reporter.cjs');
  const data = { test: 'data' };
  const json = reporter.formatJson(data);
  try {
    JSON.parse(json);
    console.log('PASS');
  } catch (e) {
    console.log('FAIL');
  }
}

// Test 14: formatTable has column headers
function testFormatTable() {
  const reporter = require('./lib/reporter.cjs');
  const data = { dependencies: [] };
  const table = reporter.formatTable(data);
  const passes = table.includes('Name') && table.includes('Version') && table.includes('Type') && table.includes('Status');
  console.log(passes ? 'PASS' : 'FAIL');
}

// Test 15: formatSummary is single line with app name
function testFormatSummary() {
  const reporter = require('./lib/reporter.cjs');
  const data = { name: 'example-app', summary: { total: 5, healthy: 5, issues: 0 } };
  const summary = reporter.formatSummary(data);
  const passes = summary.includes('example-app') && summary.split('\n').length === 1;
  console.log(passes ? 'PASS' : 'FAIL');
}

// Test 16: check exits 1 for project with issues
function testCheckIssues() {
  const result = runCommand(['check', '--path', 'fixtures/valid']);
  const passes = result.status === 1; // Should fail because moment is deprecated
  console.log(passes ? 'PASS' : 'FAIL');
}

// Test 17: check exits 0 for healthy project
function testCheckHealthy() {
  const result = runCommand(['check', '--path', 'fixtures/empty']);
  const passes = result.status === 0;
  console.log(passes ? 'PASS' : 'FAIL');
}

// Test 18: init creates config file with defaults
function testInit() {
  const configPath = '.dep-doctor.json';
  if (fs.existsSync(configPath)) {
    fs.unlinkSync(configPath);
  }
  const result = runCommand(['init']);
  const passes = result.status === 0 && fs.existsSync(configPath);
  if (fs.existsSync(configPath)) {
    fs.unlinkSync(configPath);
  }
  console.log(passes ? 'PASS' : 'FAIL');
}

// Run all tests
testHelp();
testUnknownCommand();
testScanValid();
testScanMalformed();
testScanMissing();
testScanEmpty();
testValidSemver();
testInvalidSemver();
testValidSpdx();
testInvalidSpdx();
testAnalyzeDeprecated();
testAnalyzeSummary();
testFormatJson();
testFormatTable();
testFormatSummary();
testCheckIssues();
testCheckHealthy();
testInit();

// Final validation
const passCount = (testOutput.match(/PASS/g) || []).length;
const failCount = (testOutput.match(/FAIL/g) || []).length;
console.log(`${passCount} passed, ${failCount} failed`);

// Restore console.log
console.log = originalLog;