URL=https://github.com/ragnaringi/juce-tools/releases/download/0.1.0/juce-tools
echo "Downloading binary from $URL"
curl -L -o ./juce-tools $URL
echo "Setting mode to executable"
chmod +x ./juce-tools
DEST=/usr/local/bin/
echo "Installing binary to $DEST"
mv ./juce-tools $DEST/juce-tools