'use strict';

const fs = require('fs');
const path = require('path');

const CONFIG_FILE = '.dep-doctor.json';
const DEFAULTS = {
  maxAge: 365,
  denyLicenses: [],
  ignorePackages: [],
  format: 'summary'
};

function loadConfig(dirPath) {
  const configPath = path.join(dirPath, CONFIG_FILE);
  
  try {
    const fileContent = fs.readFileSync(configPath, 'utf8');
    let parsed;
    try {
      parsed = JSON.parse(fileContent);
    } catch (err) {
      const e = new Error('Invalid JSON in config file');
      e.code = 'INVALID_JSON';
      throw e;
    }

    // Merge with defaults
    return {
      maxAge: parsed.maxAge ?? DEFAULTS.maxAge,
      denyLicenses: parsed.denyLicenses ?? DEFAULTS.denyLicenses,
      ignorePackages: parsed.ignorePackages ?? DEFAULTS.ignorePackages,
      format: parsed.format ?? DEFAULTS.format
    };
  } catch (err) {
    if (err.code === 'ENOENT') {
      return { ...DEFAULTS };
    }
    throw err;
  }
}

function saveConfig(dirPath, config) {
  const configPath = path.join(dirPath, CONFIG_FILE);
  const content = JSON.stringify(config, null, 2);
  fs.writeFileSync(configPath, content, 'utf8');
}

module.exports = {
  loadConfig,
  saveConfig,
  DEFAULTS,
  CONFIG_FILE
};