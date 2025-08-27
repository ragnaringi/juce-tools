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
}

// NewProjucer returns a new Projucer instance for a given JUCE path.
func NewProjucer(jucePath string) *Projucer {
	path := filepath.Join(jucePath, "extras/Projucer")
	buildsPath := filepath.Join(path, "Builds", platformIdentifier)
	projectPath := filepath.Join(buildsPath, "Projucer"+ideProjectExtension)

	return &Projucer{
		path,
		buildsPath,
		projectPath,
		initBinaryPath(buildsPath),
	}
}

// needsBuild returns true if the Projucer binary is missing or outdated.
// It also prints the reason whenever a rebuild is required.
func (p *Projucer) needsBuild() bool {
	found, err := fileExists(p.binaryPath)
	if err != nil {
		fmt.Println("⚠️  Error checking Projucer binary:", err)
		return true
	}
	if !found {
		fmt.Println("🔄 Projucer needs build: binary missing at", p.binaryPath)
		return true
	}

	// Check binary timestamp
	binaryInfo, err := os.Stat(p.binaryPath)
	if err != nil {
		fmt.Println("⚠️  Could not stat Projucer binary:", err)
		return true
	}
	binaryModTime := binaryInfo.ModTime()

	// Walk the Projucer source tree and find any file newer than binary
	var newestSource string
	var newestModTime int64

	err = filepath.Walk(p.path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Ignore unreadable files
		}
		if info.IsDir() {
			return nil
		}
		if info.ModTime().After(binaryModTime) {
			if info.ModTime().Unix() > newestModTime {
				newestSource = path
				newestModTime = info.ModTime().Unix()
			}
		}
		return nil
	})

	if err != nil {
		fmt.Println("⚠️  Failed walking Projucer source files:", err)
		return true
	}

	if newestSource != "" {
		fmt.Printf("🔄 Projucer needs rebuild: newer source file found → %s\n", newestSource)
		return true
	}

	return false
}

// build builds the Projucer IDE project if needed.
func (p *Projucer) build() (bool, error) {
	if !p.needsBuild() {
		fmt.Println("✅ Projucer is up to date:", p.binaryPath)
		return true, nil
	}

	if found, _ := fileExists(p.projectPath); found {
		fmt.Println("🔨 Building Projucer from:", p.projectPath)
		return build(p.projectPath, "Projucer - App")
	}
	return false, errors.New("unable to find Projucer IDE project")
}

// open launches the Projucer binary with a given project file.
func (p *Projucer) open(projectFile string) (bool, error) {
	if ok, err := p.build(); !ok {
		return false, fmt.Errorf("failed to build Projucer: %w", err)
	}

	cmd := exec.Command(p.binaryPath, projectFile)
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	return true, nil
}

// export resaves a project file via Projucer.
func (p *Projucer) export(projectFile string) (bool, error) {
	if ok, err := p.build(); !ok {
		return false, fmt.Errorf("failed to build Projucer: %w", err)
	}

	cmd := exec.Command(p.binaryPath, "--resave", projectFile)
	return run(cmd)
}

// cleanBuildArtefacts deletes build artefacts.
func (p *Projucer) cleanBuildArtefacts() (bool, error) {
	if err := os.RemoveAll(filepath.Join(p.buildsPath, buildArtefactsPath)); err != nil {
		return false, err
	}
	return true, nil
}

// initBinaryPath returns the expected path of the Projucer binary per platform.
func initBinaryPath(buildsPath string) string {
	switch runtime.GOOS {
	case "windows":
		return filepath.Join(buildsPath, buildArtefactsPath, "Release/App/Projucer.exe")
	case "darwin":
		return filepath.Join(buildsPath, buildArtefactsPath, "Release/Projucer.app/Contents/MacOS/Projucer")
	default:
		panic("Platform not supported")
	}
}
