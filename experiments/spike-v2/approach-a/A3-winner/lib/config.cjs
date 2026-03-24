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
  const filePath = path.join(dirPath, CONFIG_FILE);

  let content;
  try {
    content = fs.readFileSync(filePath, 'utf8');
  } catch (err) {
    return { ...DEFAULTS };
  }

  let config;
  try {
    config = JSON.parse(content);
  } catch (err) {
    return { ...DEFAULTS };
  }

  return {
    maxAge: config.maxAge !== undefined ? config.maxAge : DEFAULTS.maxAge,
    denyLicenses: Array.isArray(config.denyLicenses) ? config.denyLicenses : DEFAULTS.denyLicenses,
    ignorePackages: Array.isArray(config.ignorePackages) ? config.ignorePackages : DEFAULTS.ignorePackages,
    format: config.format !== undefined ? config.format : DEFAULTS.format
  };
}

function saveConfig(dirPath, config) {
  const filePath = path.join(dirPath, CONFIG_FILE);
  const content = JSON.stringify(config, null, 2);
  fs.writeFileSync(filePath, content, 'utf8');
}

module.exports = {
  loadConfig,
  saveConfig,
  DEFAULTS,
  CONFIG_FILE
};