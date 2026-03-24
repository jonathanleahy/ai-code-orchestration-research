'use strict';

function formatJson(analysis) {
  return JSON.stringify(analysis, null, 2);
}

function formatTable(analysis) {
  const { dependencies } = analysis;
  
  // Header
  let output = 'Name         | Range    | Type   | Issues\n';
  output += '-------------|----------|--------|-------\n';
  
  // Rows
  for (const dep of dependencies) {
    const name = dep.name.padEnd(12);
    const range = dep.range.padEnd(8);
    const type = dep.type.padEnd(6);
    const issues = dep.issues.length === 0 ? 'None' : dep.issues.map(i => i.message).join(', ');
    
    output += `${name} | ${range} | ${type} | ${issues}\n`;
  }
  
  return output;
}

function formatSummary(analysis) {
  const { name, version, summary } = analysis;
  const { total, issues, deprecated, large } = summary;
  
  const status = issues > 0 ? 'UNHEALTHY' : 'HEALTHY';
  
  return `${name}@${version}: ${status} (${total} deps, ${issues} issues, ${deprecated} deprecated, ${large} large)`;
}

module.exports = {
  formatJson,
  formatTable,
  formatSummary
};