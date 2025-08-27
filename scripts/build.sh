#!/usr/bin/env bash
set -e

# Move to cli/ where go.mod exists
cd "$(dirname "$0")/../cli"

echo "Fetching Go module dependencies..."
go mod tidy    # ensures go.mod and go.sum are correct
go mod download  # downloads all modules

echo "Building juce-tools..."
go build -o ../juce-tools

echo "Built ./juce-tools"
