#!/bin/bash

set -eux

DST_DIR="/opt/estkme-cloud"

declare -A ESTKME_CLOUD_BINARIES=(
    ["x86_64"]="estkme-cloud-linux-amd64"
    ["aarch64"]="estkme-cloud-linux-arm64"
    ["mips64"]="estkme-cloud-linux-mips64"
    ["riscv64"]="estkme-cloud-linux-riscv64"
)

# Check if the architecture is supported
if [ -z "${ESTKME_CLOUD_BINARIES[$(uname -m)]}" ]; then
    echo "Unsupported architecture"
    exit 1
fi

# Install dependencies.
apt-get update -y && apt-get install -y pkg-config libcurl4-openssl-dev curl

# if estkme-cloud is already running stop it.
if supervisorctl status estkme-cloud | grep -q RUNNING; then
  supervisorctl stop estkme-cloud
fi

# Download and Install estkme-cloud.
ESTKME_CLOUD_VERSION=$(curl -Ls https://api.github.com/repos/damonto/estkme-cloud/releases/latest | grep tag_name | cut -d '"' -f 4)
curl -L -o "$DST_DIR"/estkme-cloud https://github.com/damonto/estkme-cloud/releases/download/"$ESTKME_CLOUD_VERSION"/${ESTKME_CLOUD_BINARIES[$(uname -m)]}
chmod +x "$DST_DIR"/estkme-cloud

START_CMD="/opt/estkme-cloud/estkme-cloud"
if [ $# -ge 1 ] && [ -n "$1" ]; then
    START_CMD="$START_CMD --prompt='$1'"
fi

tee /etc/supervisor/conf.d/estkme-cloud.conf << EOF
[program:estkme-cloud]
command=$START_CMD
autostart=true
autorestart=true
stdout_logfile=/dev/stdout
stdout_logfile_maxbytes=0
stderr_logfile=/dev/stderr
stderr_logfile_maxbytes=0
EOF

supervisorctl update
supervisorctl start estkme-cloud

# Clean up
rm -- "$0"
