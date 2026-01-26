# Deployment Guide

This guide covers deploying Estate Index to various hosting platforms and environments.

## Pre-Deployment Checklist

- [ ] Environment variables configured in `.env`
- [ ] Site builds without errors: `npm run build`
- [ ] All links work correctly in `public/` directory
- [ ] Assets (CSS, JS, images) load properly
- [ ] Content synced and up-to-date

## Quick Start: Any Platform

### 1. Configure for Your Domain

```bash
# Copy the environment template
cp .env.example .env

# Edit for your deployment
nano .env
```

Set the `BASE_URL` to match your deployment:

```bash
# Root domain
BASE_URL=https://yourdomain.com/

# Subpath
BASE_URL=https://yourdomain.com/properties/

# Staging
BASE_URL=https://staging.yourdomain.com/
```

### 2. Build the Site

```bash
npm run build
```

This will:
- Load configuration from `.env`
- Build the site with Hugo
- Place output in `public/`

### 3. Deploy to Your Server

Copy the `public/` directory to your web server.

## Platform-Specific Deployment

### Netlify

1. **Connect repository:**
   - Go to [netlify.com](https://netlify.com)
   - Click "New site from Git"
   - Select your repository
   - Choose "Main" branch

2. **Configure build settings:**
   - Build command: `npm run build`
   - Publish directory: `public`

3. **Set environment variables:**
   - Go to Site Settings → Build & Deploy → Environment
   - Add `BASE_URL` variable
   - For subpath: `https://yourdomain.com/properties/`

4. **Configure netlify.toml** (if using subpath):

   ```toml
   [build]
   command = "npm run build"
   publish = "public"

   [[redirects]]
   from = "/properties/*"
   to = "/properties/:splat"
   status = 200
   ```

### Vercel

1. **Import project:**
   - Go to [vercel.com](https://vercel.com)
   - Click "New Project"
   - Select your repository

2. **Configure:**
   - Framework: "Other"
   - Build Command: `npm run build`
   - Output Directory: `public`

3. **Set environment variables:**
   - Add `BASE_URL` in Project Settings → Environment Variables

4. **Deploy:**
   - Push to main branch or deploy from Vercel dashboard

### Traditional VPS/Server (Apache)

1. **SSH into server:**
   ```bash
   ssh user@your-server.com
   ```

2. **Clone repository:**
   ```bash
   cd /var/www
   git clone https://github.com/yourusername/EstateIndex.git
   cd EstateIndex
   ```

3. **Install dependencies:**
   ```bash
   # Node.js tools (npm)
   npm install

   # Hugo (if not installed)
   sudo apt-get install hugo
   # or on macOS: brew install hugo
   ```

4. **Configure for your domain:**
   ```bash
   cp .env.example .env
   nano .env
   # Set BASE_URL to your domain with trailing slash
   ```

5. **Build:**
   ```bash
   npm run build
   ```

6. **Configure Apache:**

   ```apache
   <VirtualHost *:80>
       ServerName yourdomain.com
       ServerAlias www.yourdomain.com
       DocumentRoot /var/www/EstateIndex/public
       
       <Directory /var/www/EstateIndex/public>
           Options -MultiViews
           AllowOverride All
           Require all granted
           
           RewriteEngine On
           RewriteCond %{REQUEST_FILENAME} !-f
           RewriteCond %{REQUEST_FILENAME} !-d
           RewriteRule ^ index.html [L]
       </Directory>
   </VirtualHost>
   ```

   For subpath deployments:

   ```apache
   <VirtualHost *:80>
       ServerName yourdomain.com
       DocumentRoot /var/www/html
       
       <Directory /var/www/html/properties>
           Options -MultiViews
           AllowOverride All
           Require all granted
           
           RewriteEngine On
           RewriteBase /properties/
           RewriteCond %{REQUEST_FILENAME} !-f
           RewriteCond %{REQUEST_FILENAME} !-d
           RewriteRule ^ index.html [L]
       </Directory>
   </VirtualHost>
   ```

7. **Enable mod_rewrite:**
   ```bash
   sudo a2enmod rewrite
   sudo systemctl restart apache2
   ```

### Traditional VPS/Server (Nginx)

1. **Setup directory:**
   ```bash
   sudo mkdir -p /var/www/estateindex
   cd /var/www/estateindex
   git clone https://github.com/yourusername/EstateIndex.git .
   ```

2. **Install dependencies and build:**
   ```bash
   npm install
   npm run build
   ```

3. **Configure Nginx:**

   ```nginx
   server {
       listen 80;
       server_name yourdomain.com www.yourdomain.com;
       root /var/www/estateindex/public;
       
       try_files $uri $uri/ /index.html;
       
       location ~* \.(css|js|jpg|jpeg|png|gif|ico|svg)$ {
           expires 1y;
           add_header Cache-Control "public, immutable";
       }
   }
   ```

   For subpath deployments:

   ```nginx
   server {
       listen 80;
       server_name yourdomain.com;
       root /var/www/html;
       
       location /properties/ {
           try_files $uri $uri/ /properties/index.html;
       }
       
       location ~* /properties/.*\.(css|js|jpg|jpeg|png|gif|ico|svg)$ {
           expires 1y;
           add_header Cache-Control "public, immutable";
       }
   }
   ```

4. **Test and reload:**
   ```bash
   sudo nginx -t
   sudo systemctl reload nginx
   ```

### Docker Deployment

1. **Create Dockerfile:**

   ```dockerfile
   FROM node:18-alpine AS builder
   WORKDIR /app
   COPY . .
   RUN npm install
   RUN npm run build
   
   FROM nginx:alpine
   COPY --from=builder /app/public /usr/share/nginx/html
   COPY nginx.conf /etc/nginx/nginx.conf
   EXPOSE 80
   CMD ["nginx", "-g", "daemon off;"]
   ```

2. **Create nginx.conf:**

   ```nginx
   user nginx;
   worker_processes auto;
   error_log /var/log/nginx/error.log warn;
   pid /var/run/nginx.pid;

   events {
       worker_connections 1024;
   }

   http {
       include /etc/nginx/mime.types;
       default_type application/octet-stream;

       log_format main '$remote_addr - $remote_user [$time_local] "$request" '
                       '$status $body_bytes_sent "$http_referer" '
                       '"$http_user_agent" "$http_x_forwarded_for"';

       access_log /var/log/nginx/access.log main;

       sendfile on;
       tcp_nopush on;
       keepalive_timeout 65;
       gzip on;

       server {
           listen 80;
           server_name _;
           root /usr/share/nginx/html;
           
           try_files $uri $uri/ /index.html;
           
           location ~* \.(css|js|jpg|jpeg|png|gif|ico|svg)$ {
               expires 1y;
               add_header Cache-Control "public, immutable";
           }
       }
   }
   ```

3. **Build and run:**

   ```bash
   docker build -t estateindex .
   docker run -p 80:80 estateindex
   ```

## Continuous Deployment (CD)

### GitHub Actions Workflow

Create `.github/workflows/deploy.yml`:

```yaml
name: Deploy to Production

on:
  push:
    branches: [main]

jobs:
  build-and-deploy:
    runs-on: ubuntu-latest
    
    steps:
      - uses: actions/checkout@v3
      
      - name: Setup Node.js
        uses: actions/setup-node@v3
        with:
          node-version: '18'
      
      - name: Install Hugo
        run: |
          wget https://github.com/gohugoio/hugo/releases/download/v0.121.0/hugo_0.121.0_linux-amd64.tar.gz
          tar -xzf hugo_0.121.0_linux-amd64.tar.gz
          sudo mv hugo /usr/local/bin/
      
      - name: Install dependencies
        run: npm install
      
      - name: Build site
        env:
          BASE_URL: ${{ secrets.BASE_URL }}
        run: npm run build
      
      - name: Deploy to server
        uses: appleboy/scp-action@master
        with:
          host: ${{ secrets.HOST }}
          username: ${{ secrets.USERNAME }}
          key: ${{ secrets.DEPLOY_KEY }}
          source: "public/"
          target: "/var/www/estateindex/"
          rm: true
```

## Post-Deployment

### Verify Deployment

1. **Check homepage:** Visit `https://yourdomain.com/`
2. **Test navigation:** Click all header and footer links
3. **Browse listings:** Verify listing pages load
4. **Check assets:** Inspect browser console for 404 errors
5. **Test search:** Use search/filter on listings page
6. **Mobile test:** Test on mobile device or responsive view

### Performance Optimization

1. **Enable GZIP compression:**
   - Apache: Add `mod_deflate` or `mod_gzip`
   - Nginx: Enable `gzip on;`

2. **Add caching headers:**
   ```
   Cache-Control: public, max-age=31536000
   ```
   For images, CSS, JS

3. **Use CDN:** Configure CloudFlare or similar
   - Benefits: Global distribution, caching, DDoS protection
   - Set up automatic cache purge on deployment

### Monitoring

1. **Setup SSL/HTTPS:**
   ```bash
   sudo apt-get install certbot python3-certbot-apache
   sudo certbot --apache -d yourdomain.com
   ```

2. **Monitor uptime:**
   - Use UptimeRobot, StatusPage, or similar
   - Set up alerts for downtime

3. **Monitor errors:**
   - Check server logs regularly
   - Setup error tracking (Sentry, etc.)

## Troubleshooting Deployments

### Links are broken

**Problem:** Navigation doesn't work, shows 404 errors.

**Solution:**
1. Verify `BASE_URL` in `.env` is correct and ends with `/`
2. Rebuild: `npm run build`
3. Check `public/` directory exists and has files
4. Verify web server is configured to serve `public/` as root

### Static assets not loading

**Problem:** CSS, JS, images show 404 in browser console.

**Solution:**
1. Check `BASE_URL` configuration
2. Verify web server can read `public/css/`, `public/js/`, `public/img/`
3. Check file permissions: `chmod 644 public/**/*`
4. Clear browser cache (Ctrl+Shift+Delete)

### Site works locally but not on server

**Problem:** Works with `npm run dev` but broken on production.

**Solution:**
1. Ensure `BASE_URL` is correct for production
2. Check web server document root is `public/`
3. Verify rewrite rules are configured
4. Check file ownership: `chown www-data:www-data /var/www/estateindex/ -R`

## Rollback Procedure

If deployment has issues:

1. **Keep previous build:**
   ```bash
   mv public public.backup
   git checkout HEAD~1
   npm run build
   ```

2. **Or restore from backup:**
   ```bash
   mv public.current public
   mv public.backup public.current
   ```

3. **Reload web server:**
   ```bash
   # Apache
   sudo systemctl restart apache2
   
   # Nginx
   sudo systemctl reload nginx
   ```

## Questions?

For deployment issues, check:
- ✓ `.env` configuration is correct
- ✓ `BASE_URL` matches your domain
- ✓ Web server can read files in `public/`
- ✓ Rewrite rules are configured (if needed)
- ✓ Browser console for 404 errors
- ✓ Server logs for errors
