'use strict';

function formatJson(analysis) {
  return JSON.stringify(analysis, null, 2);
}

function formatTable(analysis) {
  const header = 'Name'.padEnd(20) + '|' + 'Range'.padEnd(12) + '|' + 'Type'.padEnd(10) + '|' + 'Issues';
  const separator = '-'.repeat(20) + '+' + '-'.repeat(12) + '+' + '-'.repeat(10) + '+' + '-'.repeat(20);

  const rows = analysis.dependencies.map(dep => {
    const issues = dep.issues.length > 0 ? dep.issues.map(i => i.message).join('; ') : '-';
    return dep.name.padEnd(20) + '|' + dep.range.padEnd(12) + '|' + (dep.type || '-').padEnd(10) + '|' + issues;
  });

  return [header, separator, ...rows].join('\n');
}

function formatSummary(analysis) {
  const { name, version } = analysis;
  const { total, issues, deprecated, large } = analysis.summary;
  const status = issues > 0 ? 'UNHEALTHY' : 'OK';
  return `${name}@${version}: ${status} (${total} deps, ${issues} issues, ${deprecated} deprecated, ${large} large)`;
}

module.exports = {
  formatJson,
  formatTable,
  formatSummary
};