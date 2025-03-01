#!/bin/bash
# setup-deb-web.sh
set -e  # Exit on any error

echo "Setting up package structure for Next.js web interface..."

VERSION=${VERSION:-1.0.20}

# Create package directory structure
PKG_ROOT="serviceradar-web_${VERSION}"
mkdir -p "${PKG_ROOT}/DEBIAN"
mkdir -p "${PKG_ROOT}/usr/local/share/serviceradar-web"
mkdir -p "${PKG_ROOT}/lib/systemd/system"
mkdir -p "${PKG_ROOT}/etc/serviceradar"
mkdir -p "${PKG_ROOT}/etc/nginx/conf.d"

echo "Building Next.js application..."

# Build Next.js application
cd ./web

# Ensure package.json contains the right scripts and dependencies
if ! grep -q '"next": ' package.json; then
  echo "ERROR: This doesn't appear to be a Next.js app. Check your web directory."
  exit 1
fi

# Install dependencies with pnpm
#npm install -g pnpm
npm install

# Build the Next.js application
echo "Building Next.js application with standalone output..."
npm run build

# Copy the Next.js standalone build
echo "Copying Next.js standalone build to package..."
cp -r .next/standalone/* "../${PKG_ROOT}/usr/local/share/serviceradar-web/"
cp -r .next/standalone/.next "../${PKG_ROOT}/usr/local/share/serviceradar-web/"

# Ensure styled-jsx is properly included
mkdir -p "../${PKG_ROOT}/usr/local/share/serviceradar-web/node_modules/styled-jsx/dist"
cp -r node_modules/styled-jsx/dist/* "../${PKG_ROOT}/usr/local/share/serviceradar-web/node_modules/styled-jsx/dist/"
cp node_modules/styled-jsx/package.json "../${PKG_ROOT}/usr/local/share/serviceradar-web/node_modules/styled-jsx/"

# Copy public files if they exist
if [ -d "public" ]; then
  cp -r public "../${PKG_ROOT}/usr/local/share/serviceradar-web/"
fi

# Cleanup temp directory
rm -rf "$TEMP_DIR"

cd ..

echo "Creating package files..."

# Create default config file
cat > "${PKG_ROOT}/etc/serviceradar/web.json" << EOF
{
  "port": 3000,
  "host": "0.0.0.0",
  "api_url": "http://localhost:8090"
}
EOF

# Create Nginx configuration
cat > "${PKG_ROOT}/etc/nginx/conf.d/serviceradar-web.conf" << EOF
# ServiceRadar Web Interface - Nginx Configuration
server {
    listen 80;
    server_name _; # Catch-all server name (use your domain if you have one)

    access_log /var/log/nginx/serviceradar-web.access.log;
    error_log /var/log/nginx/serviceradar-web.error.log;

    # API proxy (assumes serviceradar-cloud package is installed)
    location /api/ {
        proxy_pass http://localhost:8090;
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
    }

    # Support for Next.js WebSockets (if used)
    location /_next/webpack-hmr {
        proxy_pass http://localhost:3000;
        proxy_http_version 1.1;
        proxy_set_header Upgrade \$http_upgrade;
        proxy_set_header Connection "upgrade";
    }

    # Main app - proxy all requests to Next.js
    location / {
        proxy_pass http://localhost:3000;
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
    }
}
EOF

# Create control file
cat > "${PKG_ROOT}/DEBIAN/control" << EOF
Package: serviceradar-web
Version: ${VERSION}
Section: utils
Priority: optional
Architecture: amd64
Depends: systemd, nodejs (>= 16.0.0), nginx
Recommends: serviceradar-cloud
Maintainer: Michael Freeman <mfreeman451@gmail.com>
Description: ServiceRadar web interface
 Next.js web interface for the ServiceRadar monitoring system.
 Includes Nginx configuration for integrated API and UI access.
Config: /etc/serviceradar/web.json
EOF

# Create conffiles to mark configuration files
cat > "${PKG_ROOT}/DEBIAN/conffiles" << EOF
/etc/serviceradar/web.json
/etc/nginx/conf.d/serviceradar-web.conf
EOF

# Create systemd service file
cat > "${PKG_ROOT}/lib/systemd/system/serviceradar-web.service" << EOF
[Unit]
Description=ServiceRadar Web Interface
After=network.target

[Service]
Type=simple
User=serviceradar
WorkingDirectory=/usr/local/share/serviceradar-web
Environment=NODE_ENV=production
Environment=PORT=3000
ExecStart=/usr/bin/node server.js
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
EOF

# Create postinst script
cat > "${PKG_ROOT}/DEBIAN/postinst" << EOF
#!/bin/bash
set -e

# Check for Nginx
if ! command -v nginx >/dev/null 2>&1; then
    echo "ERROR: Nginx is required but not installed. Please install nginx and try again."
    exit 1
fi

# Create serviceradar user if it doesn't exist
if ! id -u serviceradar >/dev/null 2>&1; then
    useradd --system --no-create-home --shell /usr/sbin/nologin serviceradar
fi

# Install Node.js if not already installed
if ! command -v node >/dev/null 2>&1; then
    echo "Installing Node.js..."
    curl -fsSL https://deb.nodesource.com/setup_18.x | bash -
    apt-get install -y nodejs
fi

# Set permissions
chown -R serviceradar:serviceradar /usr/local/share/serviceradar-web
chown -R serviceradar:serviceradar /etc/serviceradar/web.json
chmod 755 /usr/local/share/serviceradar-web
chmod 644 /etc/serviceradar/web.json

# Configure Nginx
if [ -f /etc/nginx/sites-enabled/default ]; then
    echo "Disabling default Nginx site..."
    rm -f /etc/nginx/sites-enabled/default
fi

# Create symbolic link if Nginx uses sites-enabled pattern
if [ -d /etc/nginx/sites-enabled ]; then
    ln -sf /etc/nginx/conf.d/serviceradar-web.conf /etc/nginx/sites-enabled/
fi

# Test and reload Nginx
echo "Testing Nginx configuration..."
nginx -t || { echo "Warning: Nginx configuration test failed. Please check your configuration."; }
systemctl reload nginx || systemctl restart nginx || echo "Warning: Failed to reload/restart Nginx."

# Enable and start service
systemctl daemon-reload
systemctl enable serviceradar-web
systemctl start serviceradar-web || echo "Failed to start service, please check the logs"

echo "ServiceRadar Web Interface installed successfully!"
echo "Web UI is running on port 3000"
echo "Nginx configured as reverse proxy - you can access the UI at http://localhost/"
echo "Note: For full functionality, install the serviceradar-cloud package"

exit 0
EOF

chmod 755 "${PKG_ROOT}/DEBIAN/postinst"

# Create prerm script
cat > "${PKG_ROOT}/DEBIAN/prerm" << EOF
#!/bin/bash
set -e

# Stop and disable service
systemctl stop serviceradar-web || true
systemctl disable serviceradar-web || true

# Remove Nginx symlink if exists
if [ -f /etc/nginx/sites-enabled/serviceradar-web.conf ]; then
    rm -f /etc/nginx/sites-enabled/serviceradar-web.conf
fi

# Reload Nginx if running
if systemctl is-active --quiet nginx; then
    systemctl reload nginx || true
fi

exit 0
EOF

chmod 755 "${PKG_ROOT}/DEBIAN/prerm"

echo "Building Debian package..."

# Create release-artifacts directory if it doesn't exist
mkdir -p ./release-artifacts

# Build the package with root-owner-group to avoid ownership warnings
dpkg-deb --root-owner-group --build "${PKG_ROOT}"

# Move the deb file to the release-artifacts directory
mv "${PKG_ROOT}.deb" "./release-artifacts/"

echo "Package built: release-artifacts/${PKG_ROOT}.deb"