#!/usr/bin/env node
/**
 * Vercel Cron Job: Sync XML feed and rebuild site
 * Runs every 6 hours to fetch latest listings and regenerate static site
 * 
 * This is a Vercel serverless function that:
 * 1. Runs the xmlsync CLI to fetch latest listings from XML feed
 * 2. Triggers a static site rebuild
 */

const { execSync } = require('child_process');
const path = require('path');
const fs = require('fs');

module.exports = async (req, res) => {
  // Only allow Vercel cron requests
  if (req.headers['x-vercel-cron'] !== 'true') {
    return res.status(401).json({ error: 'Unauthorized' });
  }

  try {
    console.log('[Cron] Starting XML sync and site rebuild...');
    const startTime = Date.now();

    // Set environment
    const env = {
      ...process.env,
      BASE_URL: process.env.BASE_URL || 'https://estate-index.vercel.app/',
      HUGO_ENV: 'production',
    };

    // Run xmlsync to update listings
    console.log('[Cron] Running xmlsync to fetch latest listings...');
    const cmdDir = path.join(__dirname, '../../cmd/xmlsync');
    const contentDir = path.join(__dirname, '../../content/listings');
    
    try {
      execSync(`cd "${cmdDir}" && go build -o xmlsync . && ./xmlsync -content "${contentDir}"`, {
        env,
        stdio: 'inherit',
      });
      console.log('[Cron] ✓ XML sync completed');
    } catch (error) {
      console.log('[Cron] ✗ XML sync failed, continuing with rebuild anyway...');
      // Don't fail the cron if sync fails - still rebuild the site
    }

    // Rebuild Hugo site
    console.log('[Cron] Rebuilding Hugo site...');
    try {
      execSync('npm run build', {
        cwd: path.join(__dirname, '../..'),
        env,
        stdio: 'inherit',
      });
      console.log('[Cron] ✓ Site rebuild completed');
    } catch (error) {
      console.error('[Cron] ✗ Site rebuild failed:', error.message);
      return res.status(500).json({ 
        error: 'Build failed',
        details: error.message 
      });
    }

    const duration = ((Date.now() - startTime) / 1000).toFixed(2);
    const message = `[Cron] ✓ Sync and rebuild completed in ${duration}s`;
    console.log(message);

    res.status(200).json({
      success: true,
      message,
      duration: `${duration}s`,
      timestamp: new Date().toISOString(),
    });
  } catch (error) {
    console.error('[Cron] Error:', error);
    res.status(500).json({
      error: 'Cron job failed',
      details: error.message,
      timestamp: new Date().toISOString(),
    });
  }
};
