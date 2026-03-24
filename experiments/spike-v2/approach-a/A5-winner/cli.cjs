'use strict';

const fs = require('fs');
const path = require('path');

const parser = require('./lib/parser.cjs');
const config = require('./lib/config.cjs');
const analyzer = require('./lib/analyzer.cjs');
const reporter = require('./lib/reporter.cjs');

function showHelp() {
  console.log('Usage: dep-doctor <command> [options]');
  console.log('');
  console.log('Commands:');
  console.log('  scan --path <dir>     Output parsed deps as JSON');
  console.log('  report --path <dir> [--format json|table|summary]  Full analysis report');
  console.log('  check --path <dir>    Exit 0 if healthy, exit 1 if issues found');
  console.log('  licenses --path <dir> List license + all deps');
  console.log('  init --path <dir>     Create .dep-doctor.json with defaults');
  console.log('  --help                Show usage');
  console.log('');
  console.log('Options:');
  console.log('  --path <dir>  Path to the project directory (default: current directory)');
  console.log('  --format <fmt> Format for report output (json, table, summary)');
}

function scan(dir) {
  try {
    const pkgPath = path.join(dir, 'package.json');
    const pkg = JSON.parse(fs.readFileSync(pkgPath, 'utf8'));
    const parsed = parser.parse(pkg);
    console.log(JSON.stringify(parsed, null, 2));
  } catch (err) {
    console.error(err.message);
    process.exit(1);
  }
}

function report(dir, format = 'table') {
  try {
    const pkgPath = path.join(dir, 'package.json');
    const pkg = JSON.parse(fs.readFileSync(pkgPath, 'utf8'));
    const parsed = parser.parse(pkg);
    const analyzed = analyzer.analyze(parsed);
    
    if (format === 'json') {
      console.log(reporter.formatJson(analyzed));
    } else if (format === 'summary') {
      console.log(reporter.formatSummary(analyzed));
    } else {
      console.log(reporter.formatTable(analyzed));
    }
  } catch (err) {
    console.error(err.message);
    process.exit(1);
  }
}

function check(dir) {
  try {
    const pkgPath = path.join(dir, 'package.json');
    const pkg = JSON.parse(fs.readFileSync(pkgPath, 'utf8'));
    const parsed = parser.parse(pkg);
    const analyzed = analyzer.analyze(parsed);
    
    if (analyzed.summary.issues > 0) {
      process.exit(1);
    } else {
      process.exit(0);
    }
  } catch (err) {
    console.error(err.message);
    process.exit(1);
  }
}

function licenses(dir) {
  try {
    const pkgPath = path.join(dir, 'package.json');
    const pkg = JSON.parse(fs.readFileSync(pkgPath, 'utf8'));
    const parsed = parser.parse(pkg);
    
    console.log(`License: ${parsed.license || 'UNLICENSED'}`);
    console.log('Dependencies:');
    parsed.deps.forEach(dep => {
      console.log(`  ${dep.name}@${dep.range}`);
    });
  } catch (err) {
    console.error(err.message);
    process.exit(1);
  }
}

function init(dir) {
  try {
    const configPath = path.join(dir, '.dep-doctor.json');
    const defaultConfig = config.getDefaultConfig();
    fs.writeFileSync(configPath, JSON.stringify(defaultConfig, null, 2));
    console.log(`Created ${configPath}`);
  } catch (err) {
    console.error(err.message);
    process.exit(1);
  }
}

function main() {
  const args = process.argv.slice(2);
  
  if (args.length === 0 || args[0] === '--help') {
    showHelp();
    process.exit(0);
  }
  
  const command = args[0];
  const opts = {};
  
  for (let i = 1; i < args.length; i++) {
    if (args[i] === '--path') {
      opts.path = args[++i];
    } else if (args[i] === '--format') {
      opts.format = args[++i];
    }
  }
  
  const dir = opts.path || process.cwd();
  
  switch (command) {
    case 'scan':
      scan(dir);
      break;
    case 'report':
      report(dir, opts.format);
      break;
    case 'check':
      check(dir);
      break;
    case 'licenses':
      licenses(dir);
      break;
    case 'init':
      init(dir);
      break;
    default:
      console.error(`Unknown command: ${command}`);
      process.exit(1);
  }
}

if (require.main === module) {
  main();
}

module.exports = {
  main,
  scan,
  report,
  check,
  licenses,
  init
};