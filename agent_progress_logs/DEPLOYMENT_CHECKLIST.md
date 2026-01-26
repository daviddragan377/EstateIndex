# Estate Index - Deployment Checklist

Use this checklist to ensure your site is production-ready before launching.

## ✅ Pre-Deployment Verification

### Configuration
- [ ] Update `baseURL` in `hugo.toml` to your production domain
- [ ] Update email addresses in `layouts/partials/footer.html`
- [ ] Update CRM endpoint URL in `layouts/listings/single.html` (Request Information link)
- [ ] Review and customize all content pages in `content/pages/`
- [ ] Verify XML feed URL is correct and accessible

### Build & Testing
- [ ] Run `./scripts/setup.sh` to verify all dependencies
- [ ] Run `./scripts/build.sh` successfully (no errors)
- [ ] Verify `public/` directory is created with static files
- [ ] Run `hugo server` and test locally:
  - [ ] Homepage loads correctly
  - [ ] Listings appear (after first sync)
  - [ ] Search/filter functionality works
  - [ ] Comparison feature stores state in localStorage
  - [ ] Navigation links work
  - [ ] Footer links work
  - [ ] Links to external CRM don't error

### Content Validation
- [ ] At least one XML sync has completed successfully
- [ ] Listing pages display with proper metadata
- [ ] Taxonomy pages (countries, locations) are generated
- [ ] No 404 errors on listing pages
- [ ] Meta descriptions are populated

### Performance & SEO
- [ ] Robots.txt is generated (`/robots.txt`)
- [ ] Sitemap is generated (`/sitemap.xml`)
- [ ] Meta tags are present on all pages (title, description, og:*)
- [ ] Canonical URLs are set
- [ ] No console errors in browser DevTools
- [ ] Page load time is acceptable (<3s on 3G)

### Security
- [ ] No sensitive data in templates
- [ ] No hardcoded API keys or credentials
- [ ] External links go to trusted services
- [ ] forms post to external CRM (no internal handling)
- [ ] robots.txt allows/disallows as intended

### Responsive Design
- [ ] Test on mobile (375px width)
- [ ] Test on tablet (768px width)
- [ ] Test on desktop (1920px width)
- [ ] Navigation collapses/expands properly
- [ ] Cards stack vertically on mobile
- [ ] Buttons are tap-friendly (min 48px)

### Browser Compatibility
- [ ] Test in Chrome/Chromium
- [ ] Test in Firefox
- [ ] Test in Safari
- [ ] Test in Edge
- [ ] localStorage works in all browsers

## ✅ Hosting & Deployment

### Netlify

- [ ] Repository connected to Netlify
- [ ] Build command set: `npm install && ./scripts/build.sh`
- [ ] Publish directory set: `public`
- [ ] Environment variables configured (if needed)
- [ ] Deploy webhook configured (optional)
- [ ] Domain configured (custom or netlify.app)
- [ ] HTTPS enabled (auto with Netlify)
- [ ] DNS records updated (if custom domain)

### Vercel

- [ ] Project created in Vercel
- [ ] Repository connected
- [ ] Build command set: `npm install && ./scripts/build.sh`
- [ ] Output directory set: `public`
- [ ] Environment variables configured
- [ ] Domain configured
- [ ] HTTPS enabled
- [ ] DNS records updated

### AWS S3 + CloudFront

- [ ] S3 bucket created and configured
- [ ] CloudFront distribution created
- [ ] Bucket policy allows public read
- [ ] Index document set to `index.html`
- [ ] Error document set to `404.html`
- [ ] SSL certificate configured
- [ ] DNS records updated
- [ ] Cache invalidation tested

### Self-Hosted

- [ ] Server has sufficient storage for `public/` directory
- [ ] Web server (nginx/Apache) configured
- [ ] Static file serving configured
- [ ] 404 error page configured
- [ ] HTTPS certificate installed
- [ ] DNS records updated
- [ ] Firewall allows HTTP/HTTPS

## ✅ Continuous Integration & Deployment

### GitHub Actions

- [ ] GitHub Actions workflow file exists (`.github/workflows/build-deploy.yml`)
- [ ] Secrets configured (NETLIFY_AUTH_TOKEN, NETLIFY_SITE_ID, etc.)
- [ ] Workflow triggers configured (push, schedule, manual)
- [ ] Cron job for daily syncs configured (if desired)
- [ ] Workflow runs successfully (check Actions tab)
- [ ] Deployments appear in hosting platform

### Manual Deployment Alternative

- [ ] Build command works: `./scripts/build.sh`
- [ ] Deployment command works: `./scripts/deploy.sh netlify`
- [ ] Files deployed successfully
- [ ] Site accessible at production URL

## ✅ Monitoring & Maintenance

### Analytics & Monitoring

- [ ] Google Analytics configured (optional)
- [ ] Error tracking configured (optional - Sentry, etc.)
- [ ] Site uptime monitoring configured (optional)
- [ ] Hosting platform alerts configured

### Regular Maintenance

- [ ] XML feed sync schedule established (daily/weekly)
- [ ] Backup strategy defined
- [ ] Log rotation configured (if self-hosted)
- [ ] SSL certificate renewal tracked
- [ ] Dependency updates scheduled (npm, Hugo, Go)

### Communication

- [ ] Contact information updated
- [ ] Support email verified
- [ ] Privacy policy created (if needed)
- [ ] Terms of service created (if needed)

## ✅ Launch Readiness

### Final Checks

- [ ] All team members aware of launch
- [ ] Emergency contact list defined
- [ ] Rollback plan defined
- [ ] Monitoring dashboards set up
- [ ] Status page configured (optional)

### Launch Day

- [ ] DNS propagation checked
- [ ] HTTPS working
- [ ] All pages accessible
- [ ] Search/filter working
- [ ] Comparison working
- [ ] Lead forms routing correctly
- [ ] Analytics firing (if configured)
- [ ] Team available for support

### Post-Launch

- [ ] Monitor error logs for 24 hours
- [ ] Check user feedback
- [ ] Verify analytics data flowing
- [ ] Confirm XML syncs running automatically
- [ ] Test lead form submissions
- [ ] Document any issues for future reference

---

## 🚨 Common Issues & Solutions

### Build Fails
```bash
# Clear cache and rebuild
rm -rf resources/ public/
./scripts/build.sh
```

### Netlify Build Fails
- Check build logs in Netlify dashboard
- Verify Node.js version (18+ required)
- Verify Hugo version (0.100+ required)
- Test build locally first: `./scripts/build.sh`

### Search/Comparison Not Working
- Check browser console for JS errors
- Verify `localStorage` is enabled
- Confirm CSS is applied (not a display issue)
- Check that `assets/js/main.js` is loaded

### Slow Builds
- This is normal for 12,000+ listings
- Builds typically take 5-30 seconds
- Netlify/Vercel may timeout if build takes >15min
  - Solution: Increase timeout in platform settings
  - Or optimize listing page template

### XML Sync Fails
```bash
# Test feed URL manually
curl -I "https://your-feed-url.xml"

# Run sync with dry-run
./scripts/sync.sh --dry-run

# Check logs for detailed error messages
```

---

## 📋 Pre-Launch Communication

### Internal
- [ ] Notify team of launch date
- [ ] Provide access credentials to all team members
- [ ] Document rollback procedures
- [ ] Establish on-call support schedule

### External (Optional)
- [ ] Announce launch to stakeholders
- [ ] Update marketing materials with new URL
- [ ] Send launch announcement to mailing list
- [ ] Update business cards/collateral

---

## 🎯 Success Criteria

Your deployment is successful when:

✅ Site is accessible at production URL  
✅ All pages load without errors  
✅ Search/filter functionality works  
✅ Comparison feature works  
✅ Lead forms route to CRM  
✅ HTTPS is working  
✅ Performance is acceptable (<3s load time)  
✅ No console errors in browser  
✅ XML syncs are running on schedule  
✅ Monitoring/analytics are active  

---

## 📞 Support

For issues during deployment:

1. Check [README.md](README.md) troubleshooting section
2. Review [BUILD_SUMMARY.md](BUILD_SUMMARY.md) for component details
3. Check hosting platform logs
4. Review GitHub Actions workflow logs
5. Test build locally before deploying

---

**Estate Index is ready for launch. Good luck! 🚀**
