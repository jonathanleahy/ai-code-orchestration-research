'use strict';

const { validateDependency, isValidSpdx } = require('./validator.cjs');

// Known deprecated packages (static list — no network needed)
const DEPRECATED = new Set(['request', 'moment', 'tslint', 'bower']);

// Known large packages (> 1MB installed)
const LARGE = new Set(['moment', 'lodash', 'rxjs', 'core-js']);

function analyze(parsed) {
  const results = [];

  for (const dep of parsed.deps) {
    const issues = validateDependency(dep.name, dep.range);

    if (DEPRECATED.has(dep.name)) {
      issues.push({ type: 'deprecated', message: `${dep.name} is deprecated` });
    }
    if (LARGE.has(dep.name)) {
      issues.push({ type: 'large', message: `${dep.name} is a large package — consider alternatives` });
    }

    results.push({
      name: dep.name,
      range: dep.range,
      type: dep.type,
      issues,
      healthy: issues.length === 0
    });
  }

  // License check
  const licenseIssues = [];
  if (!isValidSpdx(parsed.license)) {
    licenseIssues.push({ type: 'unknown_license', message: `Unknown license: ${parsed.license}` });
  }

  return {
    name: parsed.name,
    version: parsed.version,
    license: parsed.license,
    licenseIssues,
    dependencies: results,
    summary: {
      total: results.length,
      healthy: results.filter(r => r.healthy).length,
      issues: results.filter(r => !r.healthy).length,
      deprecated: results.filter(r => r.issues.some(i => i.type === 'deprecated')).length,
      large: results.filter(r => r.issues.some(i => i.type === 'large')).length
    }
  };
}

module.exports = { analyze, DEPRECATED, LARGE };
