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
  if (!fs.existsSync(configPath)) {
    return { ...DEFAULTS };
  }
  try {
    const raw = fs.readFileSync(configPath, 'utf8');
    return { ...DEFAULTS, ...JSON.parse(raw) };
  } catch {
    return { ...DEFAULTS };
  }
}

function saveConfig(dirPath, config) {
  const configPath = path.join(dirPath, CONFIG_FILE);
  fs.writeFileSync(configPath, JSON.stringify({ ...DEFAULTS, ...config }, null, 2) + '\n');
}

module.exports = { loadConfig, saveConfig, DEFAULTS, CONFIG_FILE };
