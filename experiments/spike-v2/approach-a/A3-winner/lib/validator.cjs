'use strict';

const KNOWN_SPDX = new Set([
  'MIT',
  'ISC',
  'BSD-2-Clause',
  'BSD-3-Clause',
  'Apache-2.0',
  'GPL-2.0-only',
  'GPL-3.0-only',
  'LGPL-2.1-only',
  'MPL-2.0',
  'UNLICENSED'
]);

const SEMVER_REGEX = /^(\^|~|>=|>|<=|<|=|\*|latest|x)?(\d+)\.(\d+)\.(\d+)(-[a-zA-Z0-9.-]+)?$/;

function isValidSemver(range) {
  if (typeof range !== 'string') return false;
  if (range === '*' || range === 'latest') return true;
  return SEMVER_REGEX.test(range);
}

function isValidSpdx(license) {
  if (typeof license !== 'string') return false;
  return KNOWN_SPDX.has(license);
}

function validateDependency(name, range) {
  const issues = [];

  if (typeof name !== 'string' || name.length === 0) {
    issues.push({ type: 'error', message: 'Dependency name must be a non-empty string' });
  }

  if (typeof range !== 'string' || range.length === 0) {
    issues.push({ type: 'error', message: 'Dependency range must be a non-empty string' });
  } else if (!isValidSemver(range)) {
    issues.push({ type: 'error', message: `Invalid semver range: ${range}` });
  }

  return issues;
}

module.exports = {
  isValidSemver,
  isValidSpdx,
  validateDependency
};