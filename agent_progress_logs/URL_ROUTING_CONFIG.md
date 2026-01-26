# URL Routing Configuration Guide

This document explains how to configure Estate Index for different deployment scenarios, including subpath deployments.

## Quick Start

1. **Copy the environment template:**
   ```bash
   cp .env.example .env
   ```

2. **Edit `.env` for your deployment:**
   ```bash
   # For root domain (default)
   BASE_URL=https://estateindex.example.com/

   # For subpath deployment
   BASE_URL=https://example.com/properties/

   # For local development
   BASE_URL=http://localhost:1313/
   ```

3. **Build the site:**
   ```bash
   npm run build
   ```

## Configuration Options

### Environment Variables (`.env`)

- **`BASE_URL`** (required)
  - The full base URL for your site
  - Must end with `/`
  - Examples:
    - `https://estateindex.example.com/` (root domain)
    - `https://example.com/properties/` (subpath)
    - `http://localhost:1313/` (development)

- **`HUGO_ENV`** (optional)
  - Environment mode: `development`, `staging`, or `production`
  - Default: `development`
  - Used for analytics and feature flags

### Hugo Configuration (`.toml`)

The Hugo configuration in `hugo.toml` is automatically updated from the `.env` file during build:

```toml
baseURL = "${BASE_URL:=https://estateindex.example.com/}"
```

This setting:
- Uses the `BASE_URL` environment variable if set
- Falls back to `https://estateindex.example.com/` if not defined
- Is read during the Hugo build process

## How It Works

### Template Links

All internal links in templates use Hugo's `relref` function:

```html
<!-- Before (hardcoded) -->
<a href="/listings/">Browse</a>

<!-- After (relative) -->
<a href="{{ relref . "/listings" }}">Browse</a>
```

This ensures links work correctly regardless of the base URL.

### Asset References

All CSS, JavaScript, and image references use `relref`:

```html
<!-- Stylesheets -->
<link rel="stylesheet" href="{{ relref . "/css/main.css" }}" />

<!-- Scripts -->
<script src="{{ relref . "/js/main.js" }}"></script>

<!-- Images -->
<img src="{{ relref . "/img/logo.png" }}" />
```

### Dynamic Path Resolution

The JavaScript utilities in `assets/js/path-utils.js` provide helper functions for runtime path resolution:

```javascript
// Get the base path for the current deployment
const basePath = getBasePath();  // Returns '/' or '/properties/'

// Resolve a relative path to absolute
const fullPath = resolvePath('listings');  // Returns '/listings' or '/properties/listings'

// Get resource paths
const cssPath = getResourcePath('css/main.css');
```

## Deployment Scenarios

### Scenario 1: Root Domain (Default)

```bash
BASE_URL=https://estateindex.example.com/
```

URLs will be:
- Home: `https://estateindex.example.com/`
- Listings: `https://estateindex.example.com/listings/`
- Detail: `https://estateindex.example.com/listings/property-id/`

### Scenario 2: Subpath on Shared Domain

```bash
BASE_URL=https://company.example.com/properties/
```

URLs will be:
- Home: `https://company.example.com/properties/`
- Listings: `https://company.example.com/properties/listings/`
- Detail: `https://company.example.com/properties/listings/property-id/`

### Scenario 3: Local Development

```bash
BASE_URL=http://localhost:1313/
```

URLs will be:
- Home: `http://localhost:1313/`
- Listings: `http://localhost:1313/listings/`
- Detail: `http://localhost:1313/listings/property-id/`

## Rebuild After Configuration Change

After changing `.env`, you must rebuild the site for changes to take effect:

```bash
# Load new environment and build
source scripts/config.sh
npm run build
```

Or use the build script directly:

```bash
./scripts/build.sh
```

## Nginx Configuration (Subpath Example)

If deploying to a subpath on Nginx, configure the `location` block:

```nginx
location /properties/ {
    try_files $uri $uri/ /properties/index.html;
}
```

This ensures all requests within the subpath are handled by Hugo's routing.

## Apache Configuration (Subpath Example)

If using Apache with mod_rewrite:

```apache
<Directory /var/www/html/properties>
    RewriteEngine On
    RewriteBase /properties/
    RewriteCond %{REQUEST_FILENAME} !-f
    RewriteCond %{REQUEST_FILENAME} !-d
    RewriteRule ^ index.html [L]
</Directory>
```

## Testing Configuration

### Local Testing with Subpath

To test a subpath deployment locally:

1. Update `.env`:
   ```bash
   BASE_URL=http://localhost:1313/test-path/
   ```

2. Rebuild and start server:
   ```bash
   npm run build
   npm run dev
   ```

3. Access at `http://localhost:1313/test-path/`

### Verify Links

Check that all links work correctly:
- Home page loads
- Navigation links work
- Listing cards link properly
- Detail pages load
- All assets (CSS, JS, images) load

## Troubleshooting

### Links Are Broken

**Problem:** Links show 404 errors or don't navigate correctly.

**Solution:** 
1. Verify `BASE_URL` is set correctly in `.env`
2. Ensure the URL ends with `/`
3. Rebuild the site: `npm run build`
4. Clear browser cache

### Styles Not Loading

**Problem:** Page loads but CSS/JS is missing.

**Solution:**
1. Check browser console for 404 errors
2. Verify `BASE_URL` is correct
3. Check `hugo.toml` to ensure `baseURL` is updated
4. Verify assets are in `public/css/` and `public/js/`

### Assets Load Twice or from Wrong Path

**Problem:** Assets load from multiple paths or incorrect location.

**Solution:**
1. Clear `public/` directory: `rm -rf public/`
2. Rebuild: `npm run build`
3. Check HTML source to verify asset paths

## Reference

### Files Modified for URL Configuration

- **`.env.example`** - Environment template
- **`.env`** - Your local configuration (create from template)
- **`hugo.toml`** - Hugo configuration (reads `BASE_URL`)
- **`scripts/config.sh`** - Environment loader script
- **`layouts/baseof.html`** - Base template with relative paths
- **`layouts/partials/header.html`** - Header navigation with `relref`
- **`layouts/partials/footer.html`** - Footer links with `relref`
- **`layouts/index.html`** - Homepage with `relref`
- **`layouts/listings/list.html`** - Listings page with dynamic path resolution
- **`assets/js/path-utils.js`** - JavaScript utilities for runtime path resolution

### Hugo Functions Used

- **`relref`** - Creates relative URLs based on baseURL
- **`absURL`** - Creates absolute URLs (not used, but available)
- **`urls.Parse`** - Parses URLs (used in templates)
- **`strings.TrimSuffix`** - Removes trailing strings

## Questions or Issues?

For configuration issues, verify:

1. ✓ `.env` file exists and is readable
2. ✓ `BASE_URL` is set and ends with `/`
3. ✓ Hugo build completes without errors
4. ✓ `public/` directory contains expected files
5. ✓ Browser console shows no 404 errors
6. ✓ Server is configured correctly for subpath (if applicable)
