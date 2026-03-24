'use strict';

const validator = require('./validator.cjs');

const DEPRECATED = new Set(['request', 'moment', 'tslint', 'bower']);
const LARGE = new Set(['moment', 'lodash', 'rxjs', 'core-js']);

function analyze(parsed) {
  const { name, version, license, deps } = parsed;
  
  const dependencies = deps.map(dep => {
    const { name, range, type } = dep;
    const issues = validator.validateDependency(name, range);
    
    const depIssues = [];
    
    // Check for deprecated packages
    if (DEPRECATED.has(name)) {
      depIssues.push({ type: 'warning', message: 'Package is deprecated' });
    }
    
    // Check for large packages
    if (LARGE.has(name)) {
      depIssues.push({ type: 'warning', message: 'Package is large' });
    }
    
    // Check for license issues
    if (license && !validator.isValidSpdx(license)) {
      depIssues.push({ type: 'error', message: `Invalid license: ${license}` });
    }
    
    // Merge issues from validator and custom checks
    const allIssues = [...issues, ...depIssues];
    
    return {
      name,
      range,
      type,
      issues: allIssues,
      healthy: allIssues.length === 0
    };
  });
  
  const total = dependencies.length;
  const healthy = dependencies.filter(d => d.healthy).length;
  const issues = dependencies.filter(d => d.issues.length > 0).length;
  const deprecated = dependencies.filter(d => DEPRECATED.has(d.name)).length;
  const large = dependencies.filter(d => LARGE.has(d.name)).length;
  
  const summary = {
    total,
    healthy,
    issues,
    deprecated,
    large
  };
  
  return {
    name,
    version,
    license,
    licenseIssues: license && !validator.isValidSpdx(license) ? [{ type: 'error', message: `Invalid license: ${license}` }] : [],
    dependencies,
    summary
  };
}

module.exports = {
  analyze,
  DEPRECATED,
  LARGE
};