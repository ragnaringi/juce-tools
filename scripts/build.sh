#!/usr/bin/env bash
set -euo pipefail

TAG="${1:-dev}"  # defaults to "dev" if no argument supplied

cd "$(dirname "$0")/../cli"

echo "Fetching Go module dependencies..."
go mod tidy
go mod download

echo "Building juce-tools for macOS..."

GOOS=darwin GOARCH=amd64 go build -o ../juce-tools_"${TAG}"_darwin_amd64 ./...
GOOS=darwin GOARCH=arm64 go build -o ../juce-tools_"${TAG}"_darwin_arm64 ./...
lipo -create ../juce-tools_"${TAG}"_darwin_amd64 ../juce-tools_"${TAG}"_darwin_arm64 -output ../juce-tools_"${TAG}"_darwin_universal

echo "Built juce-tools_${TAG}_darwin_universal"
