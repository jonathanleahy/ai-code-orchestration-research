'use strict';

const fs = require('fs');
const path = require('path');

function parse(dirPath) {
  const filePath = path.join(dirPath, 'package.json');
  
  let content;
  try {
    content = fs.readFileSync(filePath, 'utf8');
  } catch (err) {
    const error = new Error(`Cannot read file: ${filePath}`);
    error.code = 'ENOENT';
    throw error;
  }

  let pkg;
  try {
    pkg = JSON.parse(content);
  } catch (err) {
    const error = new Error(`Invalid JSON in file: ${filePath}`);
    error.code = 'INVALID_JSON';
    throw error;
  }

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
    name: pkg.name,
    version: pkg.version,
    license: pkg.license,
    deps
  };
}

module.exports = { parse };