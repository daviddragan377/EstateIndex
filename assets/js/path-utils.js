/**
 * Path utilities for Estate Index
 * Handles dynamic base path resolution for subpath deployments
 */

// Determine the base path from Hugo's site configuration
// This is injected by Hugo templates
window.siteConfig = window.siteConfig || {
  baseURL: '{{ .Site.BaseURL }}',
  basePath: '{{ strings.TrimSuffix "/" .Site.BaseURL | strings.TrimPrefix (printf "%s://%s" (urls.Parse .Site.BaseURL).Scheme (urls.Parse .Site.BaseURL).Host) }}'
};

/**
 * Get the relative path prefix for the site
 * Handles subpath deployments like /properties/ or /site/
 */
function getBasePath() {
  // Try to extract from current page URL
  const currentPath = window.location.pathname;
  
  // For development/localhost, assume root
  if (window.location.hostname === 'localhost' || window.location.hostname === '127.0.0.1') {
    return '/';
  }
  
  // Use configured base path if available, otherwise try to infer from URL
  if (window.siteConfig && window.siteConfig.basePath) {
    return window.siteConfig.basePath || '/';
  }
  
  return '/';
}

/**
 * Resolve a relative path to an absolute URL
 * @param {string} path - The path to resolve (e.g., 'listings', '/listings/', 'listings/')
 * @returns {string} The resolved absolute path
 */
function resolvePath(path) {
  const basePath = getBasePath();
  let normalizedPath = path.replace(/^\/+/, '').replace(/\/+$/, '');
  
  if (basePath === '/') {
    return '/' + normalizedPath;
  }
  
  return basePath + '/' + normalizedPath;
}

/**
 * Get resource path (CSS, JS, images)
 * @param {string} resource - The resource path (e.g., 'css/main.css', 'js/main.js')
 * @returns {string} The full resource path
 */
function getResourcePath(resource) {
  return resolvePath(resource);
}
