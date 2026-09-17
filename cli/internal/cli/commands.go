package cli

import "fmt"

func runCommand(workingDirectory string, args []string) error {
	switch args[0] {
	case "up":
		return runUp(workingDirectory)
	case "clean":
		return runClean(workingDirectory, args[1:])
	case "export":
		return runExport(workingDirectory)
	case "code":
		return runCode(workingDirectory)
	case "build":
		return runBuild(workingDirectory)
	default:
		return fmt.Errorf("unknown command: %s", args[0])
	}
}
