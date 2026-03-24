'use strict';

const fs = require('fs');
const path = require('path');

function parse(dirPath) {
  const packageJsonPath = path.join(dirPath, 'package.json');
  
  let fileContent;
  try {
    fileContent = fs.readFileSync(packageJsonPath, 'utf8');
  } catch (err) {
    if (err.code === 'ENOENT') {
      const e = new Error('package.json not found');
      e.code = 'ENOENT';
      throw e;
    }
    throw err;
  }

  let parsed;
  try {
    parsed = JSON.parse(fileContent);
  } catch (err) {
    const e = new Error('Invalid JSON in package.json');
    e.code = 'INVALID_JSON';
    throw e;
  }

  const result = {
    name: parsed.name,
    version: parsed.version,
    license: parsed.license,
    deps: []
  };

  const depTypes = ['dependencies', 'devDependencies', 'optionalDependencies', 'peerDependencies'];
  for (const type of depTypes) {
    if (parsed[type]) {
      for (const [name, range] of Object.entries(parsed[type])) {
        result.deps.push({ name, range, type });
      }
    }
  }

  return result;
}

module.exports = { parse };