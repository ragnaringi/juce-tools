package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// Projucer wraps a JUCEProject and handles building the Projucer IDE binary.
type Projucer struct {
	project    *JUCEProject // internal JUCEProject for paths
	binaryPath string       // path to Projucer binary
	rootDir    string       // current working directory
}

// NewProjucer creates a new Projucer instance for a given Projucer project path.
func NewProjucer(projectFile string) *Projucer {
	buildsPath := filepath.Dir(projectFile)
	binaryPath := initBinaryPath(buildsPath)
	rootDir, _ := os.Getwd()

	baseProject := &JUCEProject{
		directory:     filepath.Dir(filepath.Dir(buildsPath)),
		buildsPath:    buildsPath,
		jucerFilePath: projectFile,
		buildFilePath: projectFile,
		name:          "Projucer",
		schemeName:    "App",
	}

	return &Projucer{
		project:    baseProject,
		binaryPath: binaryPath,
		rootDir:    rootDir,
	}
}

// Build ensures the Projucer binary exists; rebuilds if missing or outdated.
func (p *Projucer) Build() (bool, error) {
	needs, err := p.needsBuild()
	if err != nil {
		fail("Could not check Projucer: %v", err)
		return false, err
	}

	if !needs {
		success("Projucer is up to date")
		return true, nil
	}

	warn("Projucer binary not found or outdated → rebuilding...")

	notice("Building Projucer from: %s", relativePath(p.rootDir, p.project.buildFilePath))
	return p.project.Build()
}

// Open builds Projucer if necessary, then launches the binary with a project file.
func (p *Projucer) Open(projectFile string) (bool, error) {
	cmd := exec.Command(p.binaryPath, projectFile)
	if err := cmd.Start(); err != nil {
		return false, fmt.Errorf("failed to open project: %w", err)
	}
	return true, nil
}

// Export uses the Projucer binary to resave a given project file.
func (p *Projucer) Export(projectFile string) (bool, error) {
	// Normalize path separators for Windows
	projectFile = filepath.FromSlash(projectFile)

	// Set up the command
	cmd := exec.Command(p.binaryPath, "--resave", projectFile)

	// Set working directory to the directory containing the .jucer file
	cmd.Dir = filepath.Dir(projectFile)

	// Optional: stream stdout/stderr directly to console for debugging
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the command
	err := cmd.Run()
	if err != nil {
		return false, fmt.Errorf("Projucer failed: %w", err)
	}

	return true, nil
}

// Clean deletes only the Projucer build artefacts (same as old cleanBuildArtefacts)
func (p *Projucer) Clean() (bool, error) {
	if err := os.RemoveAll(filepath.Join(p.project.buildsPath, buildArtefactsPath)); err != nil {
		return false, err
	}
	return true, nil
}

// shouldIgnoreFile returns true if the file is a nuisance file (not source).
func shouldIgnoreFile(path string) bool {
	name := filepath.Base(path)
	ext := filepath.Ext(path)

	// macOS artifacts
	if name == ".DS_Store" || ext == ".nib" || ext == ".plist" {
		return true
	}

	// Windows build artifacts (in case anything slips through)
	if ext == ".obj" || ext == ".pdb" || ext == ".ilk" {
		return true
	}

	// JUCE binary itself
	if name == "Projucer.exe" || name == "Projucer" {
		return true
	}

	return false
}

// needsBuild checks if the Projucer binary is missing or older than any source files.
func (p *Projucer) needsBuild() (bool, error) {
	found, err := fileExists(p.binaryPath)
	if err != nil {
		return false, fmt.Errorf("checking Projucer binary: %w", err)
	}
	if !found {
		return true, nil
	}

	binaryInfo, err := os.Stat(p.binaryPath)
	if err != nil {
		return true, nil
	}

	binaryModTime := binaryInfo.ModTime()

	sourceDirs := []string{
		filepath.Join(p.project.directory, "Source"),
		filepath.Join(p.project.directory, "modules"),
	}

	for _, dir := range sourceDirs {
		err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}

			// Skip build folders entirely
			if info.IsDir() && strings.HasPrefix(path, filepath.Join(p.project.directory, "Builds")) {
				return filepath.SkipDir
			}

			if info.IsDir() {
				return nil
			}

			if shouldIgnoreFile(path) {
				return nil
			}

			if info.ModTime().After(binaryModTime) {
				return fmt.Errorf("newer file found") // short-circuit
			}
			return nil
		})
		if err != nil && err.Error() == "newer file found" {
			return true, nil
		}
		if err != nil {
			return true, fmt.Errorf("walking source dir: %w", err)
		}
	}

	return false, nil
}

// initBinaryPath returns the expected Projucer binary path per platform.
func initBinaryPath(buildsPath string) string {
	switch runtime.GOOS {
	case "windows":
		return filepath.Join(buildsPath, buildArtefactsPath, "Release", "App", "Projucer.exe")
	case "darwin":
		return filepath.Join(buildsPath, buildArtefactsPath, "Release", "Projucer.app", "Contents", "MacOS", "Projucer")
	case "linux":
		return filepath.Join(buildsPath, buildArtefactsPath, "Release", "Projucer")
	default:
		fail("Unsupported platform: %s", runtime.GOOS)
		os.Exit(1)
	}
	return ""
}
