'use strict';

const { validateDependency, isValidSpdx } = require('./validator.cjs');

const DEPRECATED = new Set(['request', 'moment', 'tslint', 'bower']);
const LARGE = new Set(['moment', 'lodash', 'rxjs', 'core-js']);

function analyze(parsed) {
  const { name, version, license, deps = [] } = parsed;

  const licenseIssues = [];
  if (license && !isValidSpdx(license)) {
    licenseIssues.push({
      type: 'warning',
      message: `Potentially invalid license: ${license}`
    });
  }

  const dependencies = deps.map(dep => {
    const issues = validateDependency(dep.name, dep.range);

    if (DEPRECATED.has(dep.name)) {
      issues.push({
        type: 'warning',
        message: `Package ${dep.name} is deprecated`
      });
    }

    if (LARGE.has(dep.name)) {
      issues.push({
        type: 'info',
        message: `Package ${dep.name} is a large package`
      });
    }

    return {
      name: dep.name,
      range: dep.range,
      type: dep.type,
      issues,
      healthy: issues.length === 0
    };
  });

  const summary = {
    total: dependencies.length,
    healthy: dependencies.filter(d => d.healthy).length,
    issues: dependencies.filter(d => !d.healthy).length,
    deprecated: dependencies.filter(d => DEPRECATED.has(d.name)).length,
    large: dependencies.filter(d => LARGE.has(d.name)).length
  };

  return {
    name,
    version,
    license,
    licenseIssues,
    dependencies,
    summary
  };
}

module.exports = {
  analyze,
  DEPRECATED,
  LARGE
};