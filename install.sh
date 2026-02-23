#!/bin/sh
set -e

REPO="shanepadgett/agent-extensions"
BINARY="ae"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"

# Detect OS
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
case "$OS" in
    darwin) OS="darwin" ;;
    linux) OS="linux" ;;
    mingw*|msys*|cygwin*) OS="windows" ;;
    *) echo "Unsupported OS: $OS"; exit 1 ;;
esac

# Detect architecture
ARCH=$(uname -m)
case "$ARCH" in
    x86_64|amd64) ARCH="amd64" ;;
    arm64|aarch64) ARCH="arm64" ;;
    *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

# Get latest version
echo "Fetching latest version..."
VERSION=$(curl -sL "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name"' | sed -E 's/.*"([^"]+)".*/\1/')
if [ -z "$VERSION" ]; then
    echo "Failed to get latest version"
    exit 1
fi
echo "Latest version: $VERSION"

# Build download URL
FILENAME="${BINARY}_${VERSION#v}_${OS}_${ARCH}"
if [ "$OS" = "windows" ]; then
    FILENAME="${FILENAME}.zip"
else
    FILENAME="${FILENAME}.tar.gz"
fi
URL="https://github.com/${REPO}/releases/download/${VERSION}/${FILENAME}"

# Download
echo "Downloading ${URL}..."
TMP_DIR=$(mktemp -d)
trap 'rm -rf "$TMP_DIR"' EXIT

curl -sL "$URL" -o "${TMP_DIR}/${FILENAME}"

# Extract
echo "Extracting..."
cd "$TMP_DIR"
if [ "$OS" = "windows" ]; then
    unzip -q "$FILENAME"
else
    tar -xzf "$FILENAME"
fi

# Install
echo "Installing to ${INSTALL_DIR}..."
if [ -w "$INSTALL_DIR" ]; then
    mv "$BINARY" "$INSTALL_DIR/"
else
    echo "Need sudo to install to ${INSTALL_DIR}"
    sudo mv "$BINARY" "$INSTALL_DIR/"
fi

chmod +x "${INSTALL_DIR}/${BINARY}"

echo ""
echo "Successfully installed ${BINARY} ${VERSION} to ${INSTALL_DIR}/${BINARY}"
echo "Run 'ae --help' to get started"
