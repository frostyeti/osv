#!/usr/bin/env bash
set -e

REPO="frostyeti/osv"
INSTALL_DIR="$HOME/.local/bin"

# Determine OS
OS="$(uname -s)"
case "${OS}" in
    Linux*)     OS_NAME="linux";;
    Darwin*)    OS_NAME="darwin";;
    *)          echo "Unsupported OS: ${OS}"; exit 1;;
esac

# Determine Architecture
ARCH="$(uname -m)"
case "${ARCH}" in
    x86_64)     ARCH_NAME="amd64";;
    arm64)      ARCH_NAME="arm64";;
    aarch64)    ARCH_NAME="arm64";;
    *)          echo "Unsupported architecture: ${ARCH}"; exit 1;;
esac

echo "Detecting latest release for ${OS_NAME}_${ARCH_NAME}..."
LATEST_RELEASE=$(curl -s "https://api.github.com/repos/${REPO}/releases/latest")
ASSET_URL=$(echo "$LATEST_RELEASE" | grep -o "https://github.com/${REPO}/releases/download/[^\"]*${OS_NAME}_${ARCH_NAME}\.tar\.gz" | head -n 1)

if [ -z "$ASSET_URL" ]; then
    echo "Could not find a release for ${OS_NAME} ${ARCH_NAME}."
    exit 1
fi

echo "Downloading ${ASSET_URL}..."
TMP_DIR=$(mktemp -d)
TAR_FILE="${TMP_DIR}/osv.tar.gz"

curl -sL "$ASSET_URL" -o "$TAR_FILE"

echo "Extracting..."
tar -xzf "$TAR_FILE" -C "$TMP_DIR"

echo "Installing to ${INSTALL_DIR}..."
mkdir -p "$INSTALL_DIR"
mv "${TMP_DIR}/osv" "${INSTALL_DIR}/osv"
chmod +x "${INSTALL_DIR}/osv"

rm -rf "$TMP_DIR"

echo "========================================================="
echo "osv was successfully installed to ${INSTALL_DIR}/osv"
echo ""
if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
    echo "WARNING: ${INSTALL_DIR} is not in your PATH."
    echo "Please add the following line to your ~/.bashrc, ~/.zshrc, or profile:"
    echo ""
    echo "    export PATH=\"\$HOME/.local/bin:\$PATH\""
    echo ""
fi
echo "Run 'osv --help' to get started!"
echo "========================================================="