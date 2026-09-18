# juce-tools

A command line utility to simplify working with JUCE projects.

Run it from a JUCE project directory containing one or more .jucer files. The tool uses the local JUCE Projucer to open, export, and build projects.

## Commands

### juce-tools up

Opens a .jucer project using the local Projucer.

### juce-tools export

Exports the selected .jucer project using Projucer.

### juce-tools exporters

Lists exporters configured in the selected .jucer project.

Example output:

    Available exporters for MyProject:
      MacOSX
      iOS
      VisualStudio2026

### juce-tools code

Exports the selected .jucer project and opens the default exported project for the current platform.

To open a specific exporter:

    juce-tools code --exporter iOS

### juce-tools build

Exports the selected .jucer project and builds the default exported project for the current platform.

To build a specific exporter:

    juce-tools build --exporter MacOSX

To specify a build scheme:

    juce-tools build --exporter MacOSX --scheme "MyProject - Standalone Plugin"

### juce-tools clean

Removes the contents of the local Builds directory.

To also remove Projucer build artifacts:

    juce-tools clean --all

## Options

### --verbose

Streams command output directly to the terminal.

Example:

    juce-tools --verbose build --exporter MacOSX --scheme "MyProject - Standalone Plugin"

Without --verbose, commands are printed before they run and output is shown only if a command fails.

## Installation

Clone the repo and build the Go module directly:

    go build -o juce-tools

Or install the pre-built binaries using the provided install scripts.

### macOS

    /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/ragnaringi/juce-tools/main/scripts/install.sh)"

### Windows

Run from an Administrator shell:

    Invoke-WebRequest -Uri https://raw.githubusercontent.com/ragnaringi/juce-tools/main/scripts/install.bat -OutFile .\temp.bat; .\temp.bat; rm .\temp.bat