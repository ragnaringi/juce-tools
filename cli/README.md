A command line utility to simplify working with JUCE projects. Install then run from any directory that contains a JUCE installation and one or more .jucer projects.

## Commands available:

`juce-tools up` - Opens up a .jucer project using the local Projucer
`juce-tools export` - Exports a .jucer project
`juce-tools code` - Exports a .jucer project and opens in the relevant IDE
`juce-tools build` - Compiles an exported .jucer project using the platform build tools

## Steps to install:

### Mac:
`/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/ragnaringi/juce-tools/develop/cli/install.sh)"`