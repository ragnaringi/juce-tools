package main

import (
	"fmt"
	"os"
	"runtime"

	"golang.org/x/sys/windows"
)

var (
	colorReset  = ""
	colorRed    = ""
	colorGreen  = ""
	colorYellow = ""
	colorBlue   = ""
)

func init() {
	// Enable ANSI colors on Windows 10+
	if runtime.GOOS == "windows" {
		enableANSI()
	}

	colorGreen = "\033[32m"
	colorYellow = "\033[33m"
	colorRed = "\033[31m"
	colorBlue = "\033[34m"
	colorReset = "\033[0m"
}

// enableANSI enables ANSI escape codes on Windows 10+ consoles.
func enableANSI() {
	stdout := windows.Handle(os.Stdout.Fd())
	var mode uint32
	if err := windows.GetConsoleMode(stdout, &mode); err != nil {
		return
	}
	const ENABLE_VIRTUAL_TERMINAL_PROCESSING = 0x0004
	_ = windows.SetConsoleMode(stdout, mode|ENABLE_VIRTUAL_TERMINAL_PROCESSING)
}

// Log functions without icons
func success(format string, args ...interface{}) {
	fmt.Printf(colorGreen+format+colorReset+"\n", args...)
}

func warn(format string, args ...interface{}) {
	fmt.Printf(colorYellow+format+colorReset+"\n", args...)
}

func fail(format string, args ...interface{}) {
	fmt.Printf(colorRed+format+colorReset+"\n", args...)
}

func notice(format string, args ...interface{}) {
	fmt.Printf(format+"\n", args...)
}
