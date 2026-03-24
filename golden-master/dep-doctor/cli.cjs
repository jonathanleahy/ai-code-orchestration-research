#!/usr/bin/env node
'use strict';

const { parse } = require('./lib/parser.cjs');
const { analyze } = require('./lib/analyzer.cjs');
const { formatJson, formatTable, formatSummary } = require('./lib/reporter.cjs');
const { loadConfig, saveConfig, DEFAULTS } = require('./lib/config.cjs');

const args = process.argv.slice(2);
const command = args[0];

function getFlag(flag, fallback) {
  const idx = args.indexOf(flag);
  return idx !== -1 && args[idx + 1] ? args[idx + 1] : fallback;
}

function printHelp() {
  console.log(`dep-doctor — Dependency Health Checker

Usage: dep-doctor <command> [options]

Commands:
  scan       Scan a directory's package.json and list dependencies
  report     Full health report with analysis
  check      Pass/fail gate (exit 0 = healthy, exit 1 = issues)
  licenses   License inventory for all dependencies
  init       Create a .dep-doctor.json config file

Options:
  --path <dir>       Target directory (default: .)
  --format <fmt>     Output format: json, table, summary (default: summary)
  --max-age <days>   Max acceptable dependency age in days (default: 365)
  --deny-license <l> Deny a specific license (can be repeated)
  --help             Show this help`);
}

if (!command || command === '--help' || command === '-h') {
  printHelp();
  process.exit(0);
}

const dir = getFlag('--path', '.');
const format = getFlag('--format', 'summary');

try {
  switch (command) {
    case 'scan': {
      const parsed = parse(dir);
      console.log(JSON.stringify(parsed, null, 2));
      break;
    }

    case 'report': {
      const parsed = parse(dir);
      const analysis = analyze(parsed);
      if (format === 'json') console.log(formatJson(analysis));
      else if (format === 'table') console.log(formatTable(analysis));
      else console.log(formatSummary(analysis));
      break;
    }

    case 'check': {
      const parsed = parse(dir);
      const analysis = analyze(parsed);
      console.log(formatSummary(analysis));
      process.exit(analysis.summary.issues > 0 ? 1 : 0);
      break;
    }

    case 'licenses': {
      const parsed = parse(dir);
      console.log(`License: ${parsed.license}`);
      for (const dep of parsed.deps) {
        console.log(`  ${dep.name}: ${dep.range} (${dep.type})`);
      }
      break;
    }

    case 'init': {
      saveConfig(dir, DEFAULTS);
      console.log(`Created .dep-doctor.json in ${dir}`);
      break;
    }

    default:
      console.error(`Unknown command: ${command}`);
      console.error('Run dep-doctor --help for usage');
      process.exit(1);
  }
} catch (err) {
  if (err.code === 'ENOENT') {
    console.error(`Error: ${err.message}`);
    process.exit(1);
  }
  if (err.code === 'INVALID_JSON') {
    console.error(`Error: ${err.message}`);
    process.exit(1);
  }
  throw err;
}
