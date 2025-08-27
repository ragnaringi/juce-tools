#!/usr/bin/env bash
set -e

# ----------------------------
# Configuration
# ----------------------------
REPO="ragnaringi/juce-tools"
BINARY="juce-tools"
DEST="/usr/local/bin"

# ----------------------------
# Detect OS
# ----------------------------
OS=$(uname | tr '[:upper:]' '[:lower:]')
EXT=""
ARCH="universal"  # default for macOS universal binary

if [[ "$OS" == "darwin" ]]; then
    :
elif [[ "$OS" == "linux" ]]; then
    ARCH="amd64"
else
    echo "❌ Unsupported OS: $OS"
    exit 1
fi

# ----------------------------
# Fetch latest release from GitHub
# ----------------------------
echo "🔹 Fetching latest release from GitHub..."
RELEASE_JSON=$(curl -s "https://api.github.com/repos/$REPO/releases/latest")

VERSION=$(curl -s "https://api.github.com/repos/$REPO/releases/latest" \
          | grep -m 1 '"tag_name":' | cut -d '"' -f4)
if [ -z "$VERSION" ] || [ "$VERSION" == "null" ]; then
    echo "❌ Could not determine latest release. Make sure a release exists."
    exit 1
fi

echo "🔹 Latest release: $VERSION"

# ----------------------------
# Find correct asset for OS/ARCH
# ----------------------------
URL=$(curl -s "https://api.github.com/repos/$REPO/releases/latest" \
      | grep '"browser_download_url":' \
      | grep "$OS" | cut -d '"' -f4)

if [ -z "$URL" ] || [ "$URL" == "null" ]; then
    echo "❌ Could not find a binary for $OS/$ARCH in the latest release."
    exit 1
fi

echo "⬇️  Downloading $URL"
if ! curl -fL -o "$BINARY" "$URL"; then
    echo "❌ Download failed. The binary may not exist. Check the release assets."
    exit 1
fi

# ----------------------------
# Make executable & install
# ----------------------------
chmod +x "$BINARY"
echo "Installing $BINARY to $DEST/juce-tools"
sudo mv "$BINARY" "$DEST/juce-tools"

echo "Installation complete!"
echo "   Run 'juce-tools --help' to get started."
