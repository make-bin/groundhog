'use strict';

/**
 * OpenClaw Node.js Plugin Bridge Server
 *
 * This script acts as a bridge between the Go gateway and existing JS plugins.
 * It loads plugins from the extensions directory and forwards messages.
 */

const path = require('path');
const fs = require('fs');

const EXTENSIONS_PATH = process.env.EXTENSIONS_PATH || path.join(__dirname, '..');
const PORT = process.env.BRIDGE_PORT || '50052';

/**
 * Load all JS plugins from the extensions directory.
 */
function loadPlugins(extensionsPath) {
  const plugins = [];
  try {
    const entries = fs.readdirSync(extensionsPath, { withFileTypes: true });
    for (const entry of entries) {
      if (entry.isDirectory() && entry.name !== 'bridge') {
        const pluginPath = path.join(extensionsPath, entry.name, 'index.js');
        if (fs.existsSync(pluginPath)) {
          try {
            const plugin = require(pluginPath);
            plugins.push({ name: entry.name, plugin });
            console.log(`Loaded plugin: ${entry.name}`);
          } catch (err) {
            console.error(`Failed to load plugin ${entry.name}:`, err.message);
          }
        }
      }
    }
  } catch (err) {
    console.error('Failed to read extensions directory:', err.message);
  }
  return plugins;
}

/**
 * Main entry point.
 */
function main() {
  console.log(`Node.js Plugin Bridge starting on port ${PORT}`);
  console.log(`Extensions path: ${EXTENSIONS_PATH}`);

  const plugins = loadPlugins(EXTENSIONS_PATH);
  console.log(`Loaded ${plugins.length} plugin(s)`);

  // TODO: Start gRPC server to communicate with Go gateway
  // For now, keep the process alive
  process.on('SIGINT', () => {
    console.log('Node.js Plugin Bridge shutting down...');
    process.exit(0);
  });

  process.on('SIGTERM', () => {
    console.log('Node.js Plugin Bridge shutting down...');
    process.exit(0);
  });

  // Keep process alive
  setInterval(() => {}, 60000);
}

main();
