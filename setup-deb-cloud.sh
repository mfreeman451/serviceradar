#!/bin/bash
# setup-deb-cloud.sh
set -e  # Exit on any error

echo "Setting up package structure..."

# Create package directory structure
PKG_ROOT="homemon-cloud_1.0.0"
mkdir -p "${PKG_ROOT}/DEBIAN"
mkdir -p "${PKG_ROOT}/usr/local/bin"
mkdir -p "${PKG_ROOT}/etc/homemon"
mkdir -p "${PKG_ROOT}/lib/systemd/system"

echo "Building web interface..."

# Build web interface
cd web
npm install
npm run build
cd ..

# Create a directory for the embedded content
mkdir -p pkg/cloud/api/web
cp -r web/dist pkg/cloud/api/web/

echo "Building Go binary..."

# Build Go binary with embedded web content
cd cmd/cloud
GOOS=linux GOARCH=amd64 go build -o "../../${PKG_ROOT}/usr/local/bin/homemon-cloud"
cd ../..

echo "Creating package files..."

# Create control file
cat > "${PKG_ROOT}/DEBIAN/control" << EOF
Package: homemon-cloud
Version: 1.0.0
Section: utils
Priority: optional
Architecture: amd64
Depends: systemd
Maintainer: Your Name <your.email@example.com>
Description: HomeMon cloud service with web interface
 Provides centralized monitoring and web dashboard for HomeMon.
EOF

# Create systemd service file
cat > "${PKG_ROOT}/lib/systemd/system/homemon-cloud.service" << EOF
[Unit]
Description=HomeMon Cloud Service
After=network.target

[Service]
Type=simple
User=homemon
ExecStart=/usr/local/bin/homemon-cloud -config /etc/homemon/cloud.json
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
EOF

# Create default config
cat > "${PKG_ROOT}/etc/homemon/cloud.json" << EOF
{
    "listen_addr": ":8090",
    "alert_threshold": "5m",
    "webhooks": [
        {
            "enabled": false,
            "url": "https://your-webhook-url",
            "cooldown": "15m",
            "headers": [
                {
                    "key": "Authorization",
                    "value": "Bearer your-token"
                }
            ]
        },
        {
            "enabled": false,
            "url": "https://discord.com/api/webhooks/your/webhook",
            "cooldown": "15m",
            "template": "{ \"embeds\": [{ \"title\": \"{{.alert.Title}}\", \"description\": \"{{.alert.Message}}\", \"color\": {{if eq .alert.Level \"error\"}}15158332{{else if eq .alert.Level \"warning\"}}16776960{{else}}3447003{{end}}, \"timestamp\": \"{{.alert.Timestamp}}\", \"fields\": [{ \"name\": \"Hostname\", \"value\": \"{{index .alert.Details \"hostname\"}}\", \"inline\": true }, { \"name\": \"PID\", \"value\": \"{{index .alert.Details \"pid\"}}\", \"inline\": true }, { \"name\": \"Version\", \"value\": \"{{index .alert.Details \"version\"}}\", \"inline\": true }] }] }"
        }
    ]
}
EOF

# Create postinst script
cat > "${PKG_ROOT}/DEBIAN/postinst" << EOF
#!/bin/bash
set -e

# Create homemon user if it doesn't exist
if ! id -u homemon >/dev/null 2>&1; then
    useradd --system --no-create-home --shell /usr/sbin/nologin homemon
fi

# Set permissions
chown -R homemon:homemon /etc/homemon
chmod 755 /usr/local/bin/homemon-cloud

# Enable and start service
systemctl daemon-reload
systemctl enable homemon-cloud
systemctl start homemon-cloud

exit 0
EOF

chmod 755 "${PKG_ROOT}/DEBIAN/postinst"

# Create prerm script
cat > "${PKG_ROOT}/DEBIAN/prerm" << EOF
#!/bin/bash
set -e

# Stop and disable service
systemctl stop homemon-cloud
systemctl disable homemon-cloud

exit 0
EOF

chmod 755 "${PKG_ROOT}/DEBIAN/prerm"

echo "Building Debian package..."

# Build the package
dpkg-deb --build "${PKG_ROOT}"

echo "Package built: ${PKG_ROOT}.deb"