'use strict';

const fs = require('fs');
const path = require('path');

function parse(dirPath) {
  const pkgPath = path.join(dirPath, 'package.json');
  if (!fs.existsSync(pkgPath)) {
    const err = new Error(`package.json not found: ${pkgPath}`);
    err.code = 'ENOENT';
    throw err;
  }

  const raw = fs.readFileSync(pkgPath, 'utf8');
  let pkg;
  try {
    pkg = JSON.parse(raw);
  } catch (e) {
    const err = new Error(`Invalid JSON in ${pkgPath}: ${e.message}`);
    err.code = 'INVALID_JSON';
    throw err;
  }

  const deps = [];
  for (const [name, range] of Object.entries(pkg.dependencies || {})) {
    deps.push({ name, range, type: 'production' });
  }
  for (const [name, range] of Object.entries(pkg.devDependencies || {})) {
    deps.push({ name, range, type: 'dev' });
  }

  return {
    name: pkg.name || 'unknown',
    version: pkg.version || '0.0.0',
    license: pkg.license || 'UNLICENSED',
    deps
  };
}

module.exports = { parse };
