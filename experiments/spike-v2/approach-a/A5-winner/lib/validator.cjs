'use strict';

const SPDX_LICENSES = new Set([
  'MIT', 'ISC', 'BSD-2-Clause', 'BSD-3-Clause', 'Apache-2.0',
  'GPL-2.0-only', 'GPL-3.0-only', 'LGPL-2.1-only', 'MPL-2.0', 'UNLICENSED'
]);

const SEMVER_REGEX = /^(\*|latest|>=?<=?|~|(?:(?:\^)?\d+\.\d+\.\d+(?:-\w+(?:\.\w+)*)?(?:\+\w+(?:\.\w+)*)?)|[\d\w]+)$/;
const SEMVER_RANGE_REGEX = /^(\*|latest|>=?<=?|~|(?:(?:\^)?\d+\.\d+\.\d+(?:-\w+(?:\.\w+)*)?(?:\+\w+(?:\.\w+)*)?)|[\d\w]+)$/;

function isValidSemver(range) {
  if (typeof range !== 'string') return false;
  return SEMVER_RANGE_REGEX.test(range);
}

function isValidSpdx(license) {
  if (typeof license !== 'string') return false;
  return SPDX_LICENSES.has(license);
}

function validateDependency(name, range) {
  const issues = [];

  if (typeof name !== 'string' || !name.trim()) {
    issues.push({ type: 'error', message: 'Dependency name must be a non-empty string' });
  }

  if (!isValidSemver(range)) {
    issues.push({ type: 'error', message: `Invalid semver range: ${range}` });
  }

  return issues;
}

module.exports = {
  isValidSemver,
  isValidSpdx,
  validateDependency
};