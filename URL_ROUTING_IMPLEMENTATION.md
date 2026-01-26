# URL Routing Configuration - Implementation Summary

## Overview

Estate Index now supports fully configurable URL routing, allowing deployment to:
- Root domains: `https://estateindex.example.com/`
- Subpaths: `https://company.example.com/properties/`
- Multiple environments with different URLs

All routing is relative and configurable via environment variables.

## What Changed

### New Files Created

1. **`.env.example`** - Environment configuration template
   - Define `BASE_URL` for any deployment scenario
   - Define `HUGO_ENV` for environment-specific settings

2. **`scripts/config.sh`** - Configuration loader script
   - Loads environment variables from `.env`
   - Validates configuration
   - Sets defaults if not specified

3. **`assets/js/path-utils.js`** - JavaScript path utilities
   - Runtime path resolution for dynamic content
   - Helper functions for resource paths
   - Support for subpath deployments

4. **`URL_ROUTING_CONFIG.md`** - Configuration guide
   - Detailed setup instructions
   - Multiple deployment scenarios
   - Troubleshooting guide

5. **`DEPLOYMENT_GUIDE.md`** - Deployment instructions
   - Platform-specific instructions (Netlify, Vercel, VPS, Docker)
   - CI/CD setup with GitHub Actions
   - Post-deployment verification
   - Performance optimization

### Modified Files

#### Hugo Configuration
- **`hugo.toml`**
  - `baseURL` now reads from `$BASE_URL` environment variable
  - Added `canonifyURLs = true` for consistent URL handling
  - Menu URLs changed from absolute (`/listings/`) to relative (`listings/`)

#### Build Scripts
- **`scripts/build.sh`**
  - Now sources `scripts/config.sh` to load environment
  - Passes `BASE_URL` and `HUGO_ENV` to Hugo during build
  - Shows configuration in build output

- **`package.json`**
  - All npm scripts now source `scripts/config.sh`
  - Ensures environment is loaded before builds
  - Includes new path-utils.js in assets build

#### Templates (Updated to use relative paths)
- **`layouts/baseof.html`**
  - CSS: `{{ relref . "/css/main.css" }}`
  - JS: `{{ relref . "/js/main.js" }}`
  - Favicon: `{{ relref . "/favicon.ico" }}`

- **`layouts/partials/header.html`**
  - Logo: `{{ relref . "/" }}`
  - All navigation: `{{ relref . "/listings" }}`, etc.
  - Browse button: `{{ relref . "/listings" }}`

- **`layouts/partials/footer.html`**
  - All footer links use `{{ relref . "..." }}`
  - Logo image: `{{ relref . "/img/estateindex_logo.png" }}`

- **`layouts/index.html`**
  - CTA buttons use `{{ relref . "/listings" }}`
  - View all button uses `{{ relref . "/listings" }}`

- **`layouts/_default/comparison.html`**
  - Browse listings link uses `{{ relref . "/listings" }}`

- **`layouts/listings/list.html`**
  - Dynamic path resolution for `/index.json` fetch
  - Uses `document.currentScript.src` for relative paths

## How It Works

### Configuration Flow

```
.env (user configuration)
  ↓
scripts/config.sh (load & validate)
  ↓
npm scripts / Hugo (build with BASE_URL)
  ↓
hugo.toml (baseURL = $BASE_URL)
  ↓
Templates (relref . "path")
  ↓
Generated HTML (correct relative URLs)
```

### Template URL Generation

When Hugo builds with `BASE_URL=https://example.com/properties/`:

```html
<!-- Template -->
<a href="{{ relref . "/listings" }}">Browse</a>

<!-- Generated (in HTML) -->
<a href="/properties/listings/">Browse</a>

<!-- Actual URL depends on where page is served -->
https://example.com/properties/listings/
```

### Runtime Path Resolution

JavaScript utilities detect current deployment context:

```javascript
getBasePath()      // Returns '/' or '/properties/' based on current URL
resolvePath(path)  // Converts 'listings' to '/listings' or '/properties/listings'
getResourcePath()  // Gets correct path for CSS, JS, images
```

## Usage Examples

### Example 1: Root Domain Deployment

```bash
# .env
BASE_URL=https://mycompany.com/
HUGO_ENV=production

# Build
npm run build

# Result: Homepage at https://mycompany.com/
#         Listings at https://mycompany.com/listings/
#         Detail at https://mycompany.com/listings/[id]/
```

### Example 2: Subpath Deployment

```bash
# .env
BASE_URL=https://mycompany.com/real-estate/
HUGO_ENV=production

# Build
npm run build

# Result: Homepage at https://mycompany.com/real-estate/
#         Listings at https://mycompany.com/real-estate/listings/
#         Detail at https://mycompany.com/real-estate/listings/[id]/
```

### Example 3: Local Development

```bash
# .env
BASE_URL=http://localhost:1313/
HUGO_ENV=development

# Start server
npm run dev

# Result: Site accessible at http://localhost:1313/
```

## Key Features

✅ **Configurable Base URL**
- Set once in `.env`, used everywhere
- Supports root domains and subpaths
- Environment-specific configuration

✅ **No Hardcoded Paths**
- All internal links use Hugo's `relref` function
- All assets use relative references
- JavaScript has runtime fallbacks

✅ **Easy Deployment**
- Same codebase for all environments
- Change one variable to switch deployment location
- Automatic handling of trailing slashes

✅ **Multiple Deployment Options**
- Netlify, Vercel, traditional VPS, Docker
- CI/CD ready with GitHub Actions
- Automatic HTTPS/SSL configuration

✅ **Environment-Specific Settings**
- Development, staging, production modes
- Different BASE_URLs for each environment
- Analytics and feature flags supported

## Building and Deploying

### Quick Start

```bash
# 1. Configure
cp .env.example .env
nano .env  # Set BASE_URL for your deployment

# 2. Build
npm run build

# 3. Deploy
# Copy public/ directory to your web server
```

### With CI/CD (GitHub Actions)

```bash
# Set repository secrets on GitHub:
# - BASE_URL (e.g., https://mycompany.com/)
# - Other deployment credentials

# Push to main branch - automatically builds and deploys
git push origin main
```

## Backward Compatibility

All changes are backward compatible:
- Existing hardcoded paths still work in templates
- Hugo's `relref` gracefully handles both relative and absolute paths
- No breaking changes to HTML output structure
- CSS and JavaScript remain unchanged

## Migration Guide (If Previously Deployed)

If you had a previous version deployed:

1. **Create `.env` file:**
   ```bash
   cp .env.example .env
   echo "BASE_URL=https://your-current-domain.com/" >> .env
   ```

2. **Rebuild:**
   ```bash
   npm run build
   ```

3. **Deploy:**
   ```bash
   # Copy public/ to your server as before
   ```

No code changes needed - it just works!

## Testing

### Test Locally with Different Base Paths

```bash
# Root domain
echo "BASE_URL=http://localhost:1313/" > .env
npm run dev

# Subpath
echo "BASE_URL=http://localhost:1313/test/" > .env
npm run dev
# Access at http://localhost:1313/test/
```

### Verify in Browser

1. Open browser DevTools (F12)
2. Go to Network tab
3. Navigate site - verify all requests have correct paths
4. Check Console for any 404 errors

## Troubleshooting

**Links don't work after deployment:**
1. Verify `BASE_URL` in `.env` is correct
2. Ensure URL ends with `/`
3. Rebuild: `npm run build`
4. Check `public/` directory has all files
5. Verify web server is serving `public/` correctly

**Assets (CSS/JS/images) return 404:**
1. Check `BASE_URL` configuration
2. Verify asset files exist in `public/css/`, `public/js/`, `public/img/`
3. Check file permissions (should be readable by web server)
4. Clear browser cache

**Subpath doesn't work:**
1. Ensure `BASE_URL` includes the subpath (e.g., `/properties/`)
2. Configure web server to serve from subpath
3. Verify rewrite rules are in place (Apache/Nginx)
4. Check that `public/` is in the correct location

## Files Reference

### Configuration
- `.env.example` - Template for environment variables
- `.env` - Your local configuration (create from template)
- `.gitignore` - Excludes `.env` from version control

### Scripts
- `scripts/config.sh` - Configuration loader (sourced by build scripts)
- `scripts/build.sh` - Main build script (now loads config)
- `package.json` - npm scripts (now source config.sh)

### Templates
- `layouts/baseof.html` - Master template (assets use relref)
- `layouts/index.html` - Homepage (links use relref)
- `layouts/partials/header.html` - Navigation (all links use relref)
- `layouts/partials/footer.html` - Footer (all links use relref)
- `layouts/listings/list.html` - Listing directory (dynamic paths)
- `layouts/_default/comparison.html` - Comparison page (links use relref)

### Documentation
- `URL_ROUTING_CONFIG.md` - Configuration guide
- `DEPLOYMENT_GUIDE.md` - Deployment instructions

## Next Steps

1. **Copy environment template:**
   ```bash
   cp .env.example .env
   ```

2. **Configure for your deployment:**
   ```bash
   # Edit .env and set BASE_URL
   nano .env
   ```

3. **Build and test:**
   ```bash
   npm run build
   npm run dev
   ```

4. **Deploy:**
   - Follow the deployment guide for your platform
   - Set `BASE_URL` environment variable on your server/platform

## Support

For configuration and deployment help, refer to:
- `URL_ROUTING_CONFIG.md` - Complete configuration guide
- `DEPLOYMENT_GUIDE.md` - Platform-specific instructions
- Issue tracker - For bugs or feature requests
