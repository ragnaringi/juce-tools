package cli

import (
	"fmt"
	"os"
	"path/filepath"
)

type JUCEInstallation struct {
	path     string
	projucer *Projucer
}

func NewJUCEInstallation(rootDirectory string) (*JUCEInstallation, error) {
	jucePath, err := findJuceDirectory(rootDirectory)
	if err != nil {
		return nil, err
	}

	if exists, err := fileExists(jucePath); err != nil {
		return nil, fmt.Errorf("checking JUCE path: %w", err)
	} else if !exists {
		return nil, fmt.Errorf("no JUCE installation found in directory")
	}

	projucerBuildsPath := filepath.Join(jucePath, "extras", "Projucer", "Builds")
	projucerProjectFile, err := findExportedProjectFile(projucerBuildsPath, "Projucer", "")
	if err != nil {
		return nil, fmt.Errorf("finding Projucer project: %w", err)
	}

	projucer := NewProjucer(projucerProjectFile)

	return &JUCEInstallation{
		path:     jucePath,
		projucer: projucer,
	}, nil
}

// findJuceDirectory attempts to locate a JUCE installation.
// Priority:
// downward scan from startDir
// upward search
// JUCE_PATH environment variable
func findJuceDirectory(startDir string) (string, error) {
	// Downward recursive scan (max depth)
	maxDepth := 5
	found := scanDownwards(startDir, 0, maxDepth)
	if found != "" {
		success("Found JUCE by scanning downward: %s", relativePath(startDir, found))
		return found, nil
	}

	// Walk upwards
	dir := startDir
	for {
		if isJuceDir(dir) {
			success("Found JUCE by walking upwards: %s", relativePath(startDir, dir))
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	// Environment variable fallback
	if envPath := os.Getenv("JUCE_PATH"); envPath != "" {
		if isJuceDir(envPath) {
			success("Found JUCE via JUCE_PATH: %s", envPath)
			return envPath, nil
		}
		warn("JUCE_PATH is set but not valid: %s", envPath)
	}

	return "", fmt.Errorf("no JUCE installation found")
}

// scanDownwards recursively searches subdirectories up to maxDepth
func scanDownwards(dir string, depth, maxDepth int) string {
	if depth > maxDepth {
		return ""
	}

	if isJuceDir(dir) {
		return dir
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return ""
	}

	for _, entry := range entries {
		if entry.IsDir() {
			candidate := filepath.Join(dir, entry.Name())
			found := scanDownwards(candidate, depth+1, maxDepth)
			if found != "" {
				return found
			}
		}
	}

	return ""
}

// isJuceDir checks if a given directory looks like a JUCE installation
func isJuceDir(path string) bool {
	modules := filepath.Join(path, "modules", "juce_core")
	extras := filepath.Join(path, "extras", "Projucer")
	readme := filepath.Join(path, "README.md")

	modulesExists, _ := fileExists(modules)
	extrasExists, _ := fileExists(extras)
	readmeExists, _ := fileExists(readme)

	return modulesExists && extrasExists && readmeExists
}
