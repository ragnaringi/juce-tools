package cli

import "fmt"

const (
	colorReset      = "\033[0m"
	colorRed        = "\033[31m"
	colorGreen      = "\033[32m"
	colorYellow     = "\033[33m"
	colorBrightCyan = "\033[96m"
)

func success(format string, args ...interface{}) {
	fmt.Printf(colorGreen+"✓ "+format+colorReset+"\n", args...)
}

func warn(format string, args ...interface{}) {
	fmt.Printf(colorYellow+"! "+format+colorReset+"\n", args...)
}

func fail(format string, args ...interface{}) {
	fmt.Printf(colorRed+"✗ "+format+colorReset+"\n", args...)
}

func command(format string, args ...interface{}) {
	fmt.Printf(colorBrightCyan+"→ "+format+colorReset+"\n", args...)
}

func notice(format string, args ...interface{}) {
	fmt.Printf(format+"\n", args...)
}
