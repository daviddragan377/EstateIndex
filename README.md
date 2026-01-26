# Estate Index

A production-ready static property listings website powered by Hugo and pure CSS. This platform ingests XML property feeds, generates fully static pages with excellent SEO, and provides client-side search, filtering, and property comparison.

## Overview

Estate Index is a **read-only, comparison and lead-funnel site** designed for serious investors, advisors, and wealth managers. It features:

- **Static site generation** via Hugo for blazing-fast page loads
- **XML feed integration** with automatic daily sync via cron job
- **Global search** powered by a JSON index generated at build time
- **Property filtering** by location, country, and type
- **Side-by-side comparison** of up to 2 listings
- **Outbound lead routing** to external CRM and investment services
- **Authoritative design** with navy (#25304a), beige (#d0ad72), and minimal aesthetic

### Key Characteristics

- ✅ **Fully static** – No runtime backend, no database required
- ✅ **No authentication** – Read-only public platform
- ✅ **No payments** – Lead funnel only, external CRM integration
- ✅ **Minimal JS** – Client-side only, vanilla JavaScript
- ✅ **SEO-optimized** – Proper meta tags, semantic HTML, fast builds
- ✅ **Mobile-responsive** – Pure CSS responsive design
- ✅ **Automatic updates** – Daily cron job syncs XML feed without UI notifications

## Technology Stack

- **Static Site Generator:** Hugo 0.152+
- **Styling:** Pure CSS (no framework dependencies)
- **JavaScript:** Vanilla ES6 (client-side only)
- **Feed Processing:** Go CLI (`cmd/xmlsync`)
- **Typography:** Georgia serif font (trust, authority)
- **Hosting:** Any static host (Netlify, Vercel, AWS S3, etc.)

## Quick Start

### Prerequisites

- **Hugo** 0.152+ ([install](https://gohugo.io/installation/))
- **Go** 1.21+ ([install](https://golang.org/doc/install))
- **Node.js** 16+ with npm ([install](https://nodejs.org/)) – optional, only for asset copying

### 1. Setup Development Environment

```bash
./scripts/setup.sh
```

This script verifies all dependencies and builds the xmlsync CLI.

### 2. Build the Site

```bash
./scripts/build.sh
```

This will:
1. Build the Go CLI tool
2. Sync the XML feed
3. Run Hugo to generate the static site
4. Output to `public/` directory

### 3. View Locally

```bash
hugo server
```

Visit `http://localhost:1313`

## Daily Automatic Updates (Cron Job)

The `cron_sync_listings.sh` script updates listings daily without UI notifications. All logging goes to `logs/listings_updates.log`.

### Setup

1. Make the cron script executable:
```bash
chmod +x ./scripts/cron_sync_listings.sh
```

2. Add to crontab (runs daily at 2 AM):
```bash
crontab -e
```

Add this line:
```
0 2 * * * /path/to/EstateIndex/scripts/cron_sync_listings.sh
```

### Logging

- **Success**: Logged to `logs/listings_updates.log` with timestamp
- **Errors**: Logged to `logs/listings_updates.log` with full error details
- **Fallback**: If sync fails, previous listings remain in place (no data loss)
- **Never**: Errors/logs never displayed in UI

Example log entry:
```
[2024-01-12 02:00:15] ==========================================
[2024-01-12 02:00:15] Starting daily listings sync
[2024-01-12 02:00:15] Building xmlsync CLI...
[2024-01-12 02:00:18] Fetching and syncing XML feed...
[2024-01-12 02:00:45] SUCCESS: XML feed synced successfully
[2024-01-12 02:00:45] Rebuilding Hugo site...
[2024-01-12 02:01:12] SUCCESS: Site rebuilt successfully
```

## Core Features

### XML Feed Sync

The `xmlsync` CLI ingests properties from an XML feed:

```bash
cd cmd/xmlsync
./xmlsync -feed "https://example.com/feed.xml" -content ./content/listings
```

- Fetches and parses XML properties
- Generates Hugo-compatible Markdown files with YAML frontmatter
- Idempotent: safe to run multiple times
- Adds new listings, updates existing, removes stale entries
- Outputs summary: added/updated/removed counts

### Search & Filtering

- Built-time JSON index generation for fast client-side search
- Filter by country, location, and type
- No server-side queries required
- Instant results

### Property Comparison

- Client-side comparison of up to 2 listings
- State stored in browser localStorage
- Side-by-side attribute comparison
- Works across browser sessions

### Lead Routing

Inquiry forms post directly to external CRM/investment services:

```html
<a href="https://your-crm.io/inquiry?property={{ .Params.id }}">
  Request Information
</a>
```

## Design System

### Color Palette

- **Primary Navy:** `#25304a` (header background)
- **Beige Accent:** `#d0ad72` (buttons, highlights)
- **Off-white:** `#f5f3f0` (footer text, creamy accents)
- **Base Pages:** White and light cream
- **Footer:** Deep black (`#1a1a1a`)
- **Text:** Near-black and muted charcoal

### Typography

- **Font Family:** Georgia serif (trust, authority, editorial seriousness)
- **Responsive hierarchy:** h1-h6, body, captions
- **Weight:** Bold for headings, regular for body

### Animations

- Minimal, smooth fade-in and float-up on listing cards
- Subtle, professional, understated
- No gimmicks or exaggerated transitions

## Deployment

### Netlify

```bash
./scripts/deploy.sh netlify
```

Or via Git integration with build command: `./scripts/build.sh`

### Vercel

```bash
./scripts/deploy.sh vercel
```

### AWS S3

```bash
./scripts/deploy.sh s3 your-bucket-name
```

### Local/Self-Hosted

```bash
./scripts/deploy.sh local /var/www/html
```

## Configuration

### Hugo Config (`hugo.toml`)

Update `baseURL` to your production domain:

```toml
baseURL = "https://estateindex.example.com/"
title = "Estate Index | Global Property Intelligence"
```

### Tailwind Config (`tailwind.config.js`)

Customize colors, fonts, spacing, and animations.

## Performance & SEO

- **All pages fully static HTML** – No runtime rendering
- **Semantic HTML5** with proper meta tags and heading hierarchy
- **Auto-generated sitemap and robots.txt**
- **Minimal JavaScript** – Client-side only, no frameworks
- **Fast builds:** ~5-30 seconds depending on listing count
- **Instant page loads** via CDN

## Project Structure

```
EstateIndex/
├── content/listings/          # Generated listing pages (from XML sync)
├── content/pages/             # Static pages (about, etc.)
├── layouts/
│   ├── listings/              # Listing templates
│   ├── partials/              # Header, footer, cards
│   ├── index.html             # Homepage
│   ├── baseof.html            # Base template
│   └── _default/              # Default templates (comparison, taxonomy)
├── assets/
│   ├── css/main.css           # Tailwind + custom utilities
│   └── js/main.js             # Client-side JS (search, comparison)
├── cmd/xmlsync/               # Go CLI for XML feed sync
├── scripts/
│   ├── build.sh               # Full build orchestration
│   ├── sync.sh                # XML sync only
│   ├── deploy.sh              # Deployment helper
│   └── watch.sh               # Dev mode with auto-rebuild
├── hugo.toml                  # Hugo configuration
├── tailwind.config.js         # Tailwind configuration
├── postcss.config.js          # PostCSS configuration
└── go.mod                     # Go module definition
```

## Troubleshooting

### Build fails: "Hugo not found"

```bash
brew install hugo  # macOS
apt-get install hugo  # Ubuntu/Debian
```

### xmlsync fails: "Failed to fetch feed"

```bash
curl -I https://www.xml2u.com/Xml/International%20Property%20Alerts_3968/7212_Default.xml
```

### Comparison feature not working

1. Check browser console for JS errors
2. Verify `localStorage` is enabled
3. Confirm `assets/js/main.js` is loaded

## FAQ

**Q: Can I use React/Vue/Next.js?**  
A: No. This is intentionally static. No runtime frameworks.

**Q: Where can I host this?**  
A: Any static host: Netlify, Vercel, AWS S3, GitHub Pages, or your own server.

**Q: How often should I sync the XML feed?**  
A: Set up a cron job or GitHub Actions to sync daily (or on schedule).

**Q: Can I add authentication or payments?**  
A: No. This is a read-only lead funnel. Route to external services for those features.

## License

[Add your license here]

## Support

For issues, feature requests, or questions, create a GitHub issue or contact support.

---

**Estate Index** – Global property intelligence for discerning investors.
