'use strict';

// Semver regex: ^, ~, >=, etc. followed by major.minor.patch
const SEMVER_RANGE = /^[\^~>=<]*\d+(\.\d+){0,2}([-.]\w+)*$/;

// Subset of common SPDX license identifiers
const KNOWN_SPDX = new Set([
  'MIT', 'ISC', 'BSD-2-Clause', 'BSD-3-Clause', 'Apache-2.0',
  'GPL-2.0-only', 'GPL-3.0-only', 'LGPL-2.1-only', 'MPL-2.0',
  'UNLICENSED', 'SEE LICENSE IN LICENSE'
]);

function isValidSemver(range) {
  if (typeof range !== 'string') return false;
  if (range === '*' || range === 'latest') return true;
  return SEMVER_RANGE.test(range.trim());
}

function isValidSpdx(license) {
  if (typeof license !== 'string') return false;
  return KNOWN_SPDX.has(license.trim());
}

function validateDependency(name, range) {
  const issues = [];
  if (typeof name !== 'string' || name.length === 0) {
    issues.push({ type: 'invalid_name', message: 'Dependency name is empty or not a string' });
  }
  if (!isValidSemver(range)) {
    issues.push({ type: 'invalid_semver', message: `Invalid semver range: ${range}` });
  }
  if (range === 'latest') {
    issues.push({ type: 'unpinned', message: 'Using "latest" is not recommended' });
  }
  return issues;
}

module.exports = { isValidSemver, isValidSpdx, validateDependency, KNOWN_SPDX };
