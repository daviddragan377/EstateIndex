/**
 * Main client-side JavaScript for Estate Index
 * Handles search, filtering, and comparison functionality
 */

document.addEventListener('DOMContentLoaded', function() {
  // Initialize features on page load
  initializeComparison();
  updateComparisonBadge();
});

/**
 * Comparison functionality
 */
class ComparisonManager {
  constructor() {
    this.storageKey = 'estateindex_comparison';
    this.maxItems = 2;
    this.load();
  }

  load() {
    const stored = localStorage.getItem(this.storageKey);
    this.items = stored ? JSON.parse(stored) : [];
  }

  save() {
    localStorage.setItem(this.storageKey, JSON.stringify(this.items.slice(0, this.maxItems)));
    this.notifyListeners();
  }

  add(propertyId) {
    if (!this.items.includes(propertyId) && this.items.length < this.maxItems) {
      this.items.push(propertyId);
      this.save();
      return true;
    }
    return false;
  }

  remove(propertyId) {
    this.items = this.items.filter(id => id !== propertyId);
    this.save();
  }

  has(propertyId) {
    return this.items.includes(propertyId);
  }

  getCount() {
    return this.items.length;
  }

  clear() {
    this.items = [];
    this.save();
  }

  getAll() {
    return [...this.items];
  }

  notifyListeners() {
    window.dispatchEvent(new CustomEvent('comparisonChanged', { detail: { count: this.items.length } }));
  }
}

const comparison = new ComparisonManager();

/**
 * Initialize comparison functionality
 */
function initializeComparison() {
  // Handle URL parameters
  const params = new URLSearchParams(window.location.search);
  
  if (params.has('add')) {
    comparison.add(params.get('add'));
    window.history.replaceState({}, '', window.location.pathname);
  }
  
  if (params.has('remove')) {
    comparison.remove(params.get('remove'));
    window.history.replaceState({}, '', window.location.pathname);
  }

  // Listen for storage changes (for cross-tab updates)
  window.addEventListener('storage', function(e) {
    if (e.key === comparison.storageKey) {
      comparison.load();
      updateComparisonBadge();
    }
  });

  // Listen for comparison changes
  window.addEventListener('comparisonChanged', updateComparisonBadge);
}

/**
 * Update comparison badge in header
 */
function updateComparisonBadge() {
  const badge = document.getElementById('comparison-badge');
  const count = comparison.getCount();

  if (!badge) return;

  if (count > 0) {
    badge.textContent = count;
    badge.classList.remove('hidden');
  } else {
    badge.classList.add('hidden');
  }
}

/**
 * Add property to comparison
 */
function addToComparison(propertyId, event) {
  if (event) {
    event.preventDefault();
    event.stopPropagation();
  }

  const added = comparison.add(propertyId);

  if (added) {
    // Visual feedback
    const btn = event?.target;
    if (btn) {
      const original = btn.textContent;
      btn.textContent = '✓ Added';
      btn.disabled = true;
      btn.classList.add('bg-green-50', 'border-green-300', 'text-green-700');

      setTimeout(() => {
        btn.textContent = original;
        btn.disabled = false;
        btn.classList.remove('bg-green-50', 'border-green-300', 'text-green-700');
      }, 2000);
    }
  } else if (comparison.getCount() >= comparison.maxItems) {
    // Show feedback that max items reached
    if (event?.target) {
      const btn = event.target;
      const original = btn.textContent;
      btn.textContent = '✗ Max items';
      btn.classList.add('bg-red-50', 'border-red-300', 'text-red-700');

      setTimeout(() => {
        btn.textContent = original;
        btn.classList.remove('bg-red-50', 'border-red-300', 'text-red-700');
      }, 2000);
    }
  }
}

/**
 * Generate search index data structure for JSON output
 * This is called at build time via Hugo config
 */
function buildSearchIndex() {
  const listings = document.querySelectorAll('[data-listing]');
  const index = [];

  listings.forEach(item => {
    index.push({
      id: item.getAttribute('data-id'),
      title: item.getAttribute('data-title'),
      location: item.getAttribute('data-location'),
      country: item.getAttribute('data-country'),
      price: item.getAttribute('data-price'),
      tags: (item.getAttribute('data-tags') || '').split(',').filter(t => t),
      url: item.getAttribute('data-url'),
    });
  });

  return index;
}

/**
 * Export for server-side use
 */
if (typeof module !== 'undefined' && module.exports) {
  module.exports = {
    ComparisonManager,
    buildSearchIndex,
  };
}
