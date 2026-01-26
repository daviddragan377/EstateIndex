You are to initialise and scaffold a production-ready static property listings website using Hugo as the static site generator. This project is a read-only comparison and lead-funnel site, not a web app. There is no authentication, no payments, no admin UI, no server-side rendering at runtime.

This repository is the single source of truth. Do not invent features, abstractions, or dependencies not explicitly requested.

CORE OBJECTIVE

Build a sleek, minimal, authoritative property listings site that:

Ingests an external XML feed (≈12,000 listings) via a Go CLI sync script

Generates fully static pages with excellent SEO

Supports global search, filtering, and 2-item comparison client-side

Routes leads outward to an external CRM and future serverless services

Is extremely fast, visually restrained, and trust-inducing

TECHNOLOGY CONSTRAINTS (NON-NEGOTIABLE)

Static site generator: Hugo

Styling: Tailwind CSS only

JavaScript: minimal, client-side only (vanilla JS or Alpine.js)

No React, no Vue, no Next.js

No runtime backend

No admin UI

No authentication

No payment handling

No custom CSS except for defining animation utilities if strictly required

Prefer existing CSS-only animation utilities if suitable

VISUAL / DESIGN SYSTEM

Implement a restrained, authoritative, investment-grade aesthetic:

Colour palette:

Primary accent: deep navy metallic gradient (used sparingly for CTAs, highlights)

Base pages: white and warm beige backgrounds

Footer: pure black

Text: near-black and muted charcoal

Avoid saturated colours or playful tones

Typography:

Use a serif font that conveys trust, authority, and editorial seriousness

Must be highly legible on desktop and mobile

Apply consistent typographic hierarchy across templates

UI behaviour:

Minimal motion

Smooth, subtle float-up / fade-in animations on listing cards

No gimmicks, no exaggerated transitions

All animations must feel professional and understated

ARCHITECTURE OVERVIEW

This is a static corpus + client-side enhancement system.

Content:

Each property listing is a static page (conceptually identical to a blog post)

Listings are categorised via Hugo taxonomies (location, country, type, etc.)

Country / locality landing pages are generated automatically via taxonomies

Interactivity (client-side only):

Global search (JSON index generated at build time)

Filtering (price, location, tags)

Comparison (select up to 2 listings, client-side state only)

Outbound actions:

Inquiry form posts directly to external CRM endpoint

Investment triage CTA is a placeholder link / hook (no implementation yet)

REQUIRED REPOSITORY STRUCTURE

Scaffold the full repo, including:

/content/

/listings/ (generated listing pages)

/layouts/

Base templates

Listing templates

Taxonomy templates

Comparison page template

Reusable partials (header, footer, cards, buttons)

/assets/

Tailwind config

JS for search, filters, comparison

/static/

Generated search index JSON

/cmd/xmlsync/

Go CLI to ingest XML and generate Hugo content

/scripts/

Shell scripts to run sync + build

/README.md

Clear instructions for syncing, building, and deploying

GO XML SYNC TOOL (CRITICAL)

Implement a Go CLI tool that:

Fetches or reads an XML feed

Parses listings into a canonical internal struct

Writes Hugo-compatible Markdown files with frontmatter into /content/listings/

Handles:

Add new listings

Remove listings no longer present

Update changed listings

Outputs a clear summary:

Added / updated / removed counts

Is idempotent and safe to run repeatedly

Can be executed manually or via cron

No database is required for listings. Filesystem is the source of truth.

SEARCH & COMPARISON

Search:

Generate a static index.json at build time

Include only necessary fields (id, title, location, price, tags)

Use a lightweight client-side search library or minimal custom JS

Comparison:

Users can select up to two listings

State is stored client-side (memory or localStorage)

“View comparison” leads to a static comparison page that:

Displays both listings side-by-side

Shows a structured attribute comparison

Includes a short pros/cons summary section (static placeholders)

COMPONENT RULES

Define shared, reusable components:

Buttons

Listing cards

CTAs

Layout primitives

Use Hugo partials consistently

Do not duplicate boilerplate

Templates must be clean, readable, and composable

PERFORMANCE & SEO

All listing and landing pages must be fully static HTML

Proper meta tags and semantic HTML

Fast builds and instant page loads via CDN

No JS dependency for content visibility

WHAT TO GENERATE NOW

Produce:

Full Hugo project scaffold

Tailwind configuration aligned to the design system

Core templates and partials

XML sync Go CLI (functional, not stubbed)

Client-side JS for search, filters, comparison

Clear README with exact commands to:

Sync XML

Build site

Deploy static output

WHAT NOT TO DO

Do not add features not described

Do not introduce auth, admin panels, or dashboards

Do not add unnecessary dependencies

Do not “improve” the product idea

Do not use React or SSR frameworks

This project must be lean, deterministic, and production-ready on first generation.

Proceed.