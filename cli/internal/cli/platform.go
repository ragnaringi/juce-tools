package cli

import (
	"fmt"
	"os/exec"
	"runtime"
)

func build(projectFile string, targetName string) (bool, error) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		const buildTool = "MSBuild"
		cmd = exec.Command(buildTool, projectFile, "/property:Configuration=Release")
	case "darwin":
		const buildTool = "xcodebuild"
		cmd = exec.Command(buildTool, "-project", projectFile, "-scheme", targetName, "-configuration", "Release", "-jobs", "8")
	default:
		return false, fmt.Errorf("platform not supported: %s", runtime.GOOS)
	}
	return run(cmd)
}

func open(filePath string) (bool, error) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/C", "start", filePath)
	case "darwin":
		cmd = exec.Command("open", filePath)
	default:
		return false, fmt.Errorf("platform not supported: %s", runtime.GOOS)
	}
	return run(cmd)
}
