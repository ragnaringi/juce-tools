package main

import (
	"os"

	"juce-tools/internal/cli"
)

func main() {
	os.Exit(cli.Run(os.Args[1:]))
}
