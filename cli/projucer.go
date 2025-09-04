package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

type Projucer struct {
	path        string
	buildsPath  string
	projectPath string
	binaryPath  string
	rootDir     string
}

// NewProjucer creates a new Projucer instance for a given JUCE path.
func NewProjucer(jucePath string) *Projucer {
	buildsPath := filepath.Join(jucePath, "extras/Projucer/Builds", platformIdentifier)
	projectPath := filepath.Join(buildsPath, "Projucer"+ideProjectExtension)
	binaryPath := initBinaryPath(buildsPath)
	rootDir, _ := os.Getwd()

	return &Projucer{
		path:        filepath.Join(jucePath, "extras/Projucer"),
		buildsPath:  buildsPath,
		projectPath: projectPath,
		binaryPath:  binaryPath,
		rootDir:     rootDir,
	}
}

// needsBuild checks if Projucer binary is missing or outdated.
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
		return true, nil // treat stat errors as needing rebuild
	}
	binaryModTime := binaryInfo.ModTime()

	var newestSource string
	err = filepath.Walk(p.path, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		if info.ModTime().After(binaryModTime) {
			newestSource = path
		}
		return nil
	})
	if err != nil {
		return true, fmt.Errorf("walking Projucer sources: %w", err)
	}

	return newestSource != "", nil
}

// build builds Projucer if needed.
func (p *Projucer) build() (bool, error) {
	needs, err := p.needsBuild()
	if err != nil {
		fail("Could not check Projucer: %v", err)
		return false, err
	}

	if !needs {
		success("Projucer is up to date")
		return true, nil
	}

	warn("Projucer binary not found or outdated, rebuilding...")

	if found, _ := fileExists(p.projectPath); found {
		notice("Building Projucer from: %s", relativePath(p.rootDir, p.projectPath))
		return build(p.projectPath, "Projucer - App")
	}

	fail("Projucer IDE project missing at %s", relativePath(p.path, p.projectPath))
	return false, errors.New("unable to find Projucer IDE project")
}

// open launches the Projucer binary with a given project file.
func (p *Projucer) open(projectFile string) (bool, error) {
	ok, err := p.build()
	if !ok {
		return false, fmt.Errorf("failed to build Projucer: %w", err)
	}

	cmd := exec.Command(p.binaryPath, projectFile)
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	return true, nil
}

// export resaves a project file.
func (p *Projucer) export(projectFile string) (bool, error) {
	if ok, err := p.build(); !ok {
		return false, fmt.Errorf("failed to build Projucer: %w", err)
	}
	notice("Exporting: %s", filepath.Base(projectFile))
	cmd := exec.Command(p.binaryPath, "--resave", projectFile)
	return run(cmd)
}

// cleanBuildArtefacts deletes Projucer build artefacts.
func (p *Projucer) cleanBuildArtefacts() (bool, error) {
	if err := os.RemoveAll(filepath.Join(p.buildsPath, buildArtefactsPath)); err != nil {
		return false, err
	}
	return true, nil
}

// initBinaryPath returns the Projucer binary path for the platform.
func initBinaryPath(buildsPath string) string {
	var binaryPath string

	switch runtime.GOOS {
	case "windows":
		binaryPath = filepath.Join(buildsPath, buildArtefactsPath, "Release", "App", "Projucer.exe")
	case "darwin":
		binaryPath = filepath.Join(buildsPath, buildArtefactsPath, "Release", "Projucer.app", "Contents", "MacOS", "Projucer")
	case "linux":
		binaryPath = filepath.Join(buildsPath, buildArtefactsPath, "Release", "Projucer")
	default:
		fail("Unsupported platform: %s", runtime.GOOS)
		os.Exit(1)
	}

	return binaryPath
}
