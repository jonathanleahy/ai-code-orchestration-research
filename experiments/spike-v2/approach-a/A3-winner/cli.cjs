'use strict';

const path = require('path');
const fs = require('fs');
const { formatJson, formatTable, formatSummary } = require('./lib/reporter.cjs');
const { analyze } = require('./lib/analyzer.cjs');

// Load parser module
let parser;
try {
  parser = require('./lib/parser.cjs');
} catch (e) {
  // Parser not available, will use manual parsing
  parser = null;
}

// Load config module
let config;
try {
  config = require('./lib/config.cjs');
} catch (e) {
  config = null;
}

const DEFAULT_CONFIG = {
  exclude: [],
  includeDevDependencies: false,
  failOnDeprecated: false,
  failOnLarge: false,
  failOnLicense: false
};

function loadPackageJson(dirPath) {
  const packagePath = path.join(dirPath, 'package.json');
  
  if (!fs.existsSync(dirPath)) {
    console.error(`Error: Directory does not exist: ${dirPath}`);
    process.exit(1);
  }
  
  if (!fs.existsSync(packagePath)) {
    console.error(`Error: package.json not found in ${dirPath}`);
    process.exit(1);
  }
  
  let content;
  try {
    content = fs.readFileSync(packagePath, 'utf8');
    return JSON.parse(content);
  } catch (e) {
    console.error(`Error: Failed to parse package.json: ${e.message}`);
    process.exit(1);
  }
}

function parseDependencies(dirPath) {
  const pkg = loadPackageJson(dirPath);
  
  const deps = [];
  
  if (pkg.dependencies) {
    for (const [name, range] of Object.entries(pkg.dependencies)) {
      deps.push({ name, range, type: 'dependencies' });
    }
  }
  
  if (pkg.devDependencies) {
    for (const [name, range] of Object.entries(pkg.devDependencies)) {
      deps.push({ name, range, type: 'devDependencies' });
    }
  }
  
  return {
    name: pkg.name || 'unknown',
    version: pkg.version || '0.0.0',
    license: pkg.license,
    deps
  };
}

function showHelp() {
  console.log(`Usage: dep-doctor [command] [options]

Commands:
  scan --path <dir>                   Output parsed deps as JSON
  report --path <dir> [--format fmt]  Full analysis report (json|table|summary)
  check --path <dir>                  Exit 0 if healthy, exit 1 if issues found
  licenses --path <dir>               List license + all deps
  init --path <dir>                   Create .dep-doctor.json with defaults
  --help                              Show this help message

Options:
  --path <dir>   Directory containing package.json
  --format       Output format for report (json, table, or summary)`);
}

function cmdScan(args) {
  const dirPath = getPathArg(args);
  const parsed = parseDependencies(dirPath);
  console.log(JSON.stringify(parsed, null, 2));
}

function cmdReport(args) {
  const dirPath = getPathArg(args);
  const format = getFormatArg(args) || 'summary';
  
  const parsed = parseDependencies(dirPath);
  const analysis = analyze(parsed);
  
  let output;
  switch (format) {
    case 'json':
      output = formatJson(analysis);
      break;
    case 'table':
      output = formatTable(analysis);
      break;
    case 'summary':
    default:
      output = formatSummary(analysis);
      break;
  }
  
  console.log(output);
}

function cmdCheck(args) {
  const dirPath = getPathArg(args);
  
  const parsed = parseDependencies(dirPath);
  const analysis = analyze(parsed);
  
  const hasIssues = analysis.licenseIssues.length > 0 || analysis.summary.issues > 0;
  
  if (hasIssues) {
    process.exit(1);
  }
  process.exit(0);
}

function cmdLicenses(args) {
  const dirPath = getPathArg(args);
  const pkg = loadPackageJson(dirPath);
  
  console.log(`License: ${pkg.license || 'UNKNOWN'}`);
  console.log(`Dependencies:`);
  
  const parsed = parseDependencies(dirPath);
  for (const dep of parsed.deps) {
    console.log(`  ${dep.name}@${dep.range} (${dep.type})`);
  }
}

function cmdInit(args) {
  const dirPath = getPathArg(args);
  
  if (!fs.existsSync(dirPath)) {
    console.error(`Error: Directory does not exist: ${dirPath}`);
    process.exit(1);
  }
  
  const configPath = path.join(dirPath, '.dep-doctor.json');
  
  if (fs.existsSync(configPath)) {
    console.error(`Error: .dep-doctor.json already exists in ${dirPath}`);
    process.exit(1);
  }
  
  try {
    fs.writeFileSync(configPath, JSON.stringify(DEFAULT_CONFIG, null, 2) + '\n');
    console.log(`Created .dep-doctor.json in ${dirPath}`);
  } catch (e) {
    console.error(`Error: Failed to write config file: ${e.message}`);
    process.exit(1);
  }
}

function getPathArg(args) {
  const idx = args.indexOf('--path');
  if (idx === -1 || idx + 1 >= args.length) {
    console.error('Error: --path argument is required');
    process.exit(1);
  }
  return args[idx + 1];
}

function getFormatArg(args) {
  const idx = args.indexOf('--format');
  if (idx !== -1 && idx + 1 < args.length) {
    return args[idx + 1];
  }
  return null;
}

function main() {
  const args = process.argv.slice(2);
  
  if (args.length === 0 || args[0] === '--help') {
    showHelp();
    process.exit(0);
  }
  
  const cmd = args[0];
  
  switch (cmd) {
    case 'scan':
      cmdScan(args.slice(1));
      break;
    case 'report':
      cmdReport(args.slice(1));
      break;
    case 'check':
      cmdCheck(args.slice(1));
      break;
    case 'licenses':
      cmdLicenses(args.slice(1));
      break;
    case 'init':
      cmdInit(args.slice(1));
      break;
    default:
      console.error(`Error: Unknown command: ${cmd}`);
      process.exit(1);
  }
}

main();