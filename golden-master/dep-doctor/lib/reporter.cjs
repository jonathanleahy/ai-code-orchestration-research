'use strict';

function formatJson(analysis) {
  return JSON.stringify(analysis, null, 2);
}

function formatTable(analysis) {
  const lines = [];
  lines.push('Name                 Range          Type        Issues');
  lines.push('-------------------  -------------  ----------  ------');
  for (const dep of analysis.dependencies) {
    const name = dep.name.padEnd(20);
    const range = dep.range.padEnd(14);
    const type = dep.type.padEnd(11);
    const issues = dep.issues.length > 0
      ? dep.issues.map(i => i.type).join(', ')
      : 'OK';
    lines.push(`${name} ${range} ${type} ${issues}`);
  }
  lines.push('');
  lines.push(`Total: ${analysis.summary.total} | Healthy: ${analysis.summary.healthy} | Issues: ${analysis.summary.issues}`);
  return lines.join('\n');
}

function formatSummary(analysis) {
  const s = analysis.summary;
  const status = s.issues === 0 ? 'HEALTHY' : 'ISSUES FOUND';
  return `${analysis.name}@${analysis.version}: ${status} (${s.total} deps, ${s.issues} issues, ${s.deprecated} deprecated, ${s.large} large)`;
}

module.exports = { formatJson, formatTable, formatSummary };
