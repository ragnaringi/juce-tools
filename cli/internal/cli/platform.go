package cli

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

func canOpenExporter(exporter string) bool {
	if exporter == "" {
		return true
	}

	switch runtime.GOOS {
	case "darwin":
		return exporter == "MacOSX" || exporter == "iOS"
	case "windows":
		return strings.HasPrefix(exporter, "VisualStudio")
	default:
		return false
	}
}

func canBuildExporter(exporter string) bool {
	if exporter == "" {
		return true
	}

	switch runtime.GOOS {
	case "darwin":
		return exporter == "MacOSX" || exporter == "iOS"
	case "windows":
		return strings.HasPrefix(exporter, "VisualStudio")
	default:
		return false
	}
}

func unsupportedExporterError(command string, exporter string) error {
	return fmt.Errorf("cannot %s exporter %q on %s", command, exporter, runtime.GOOS)
}

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
