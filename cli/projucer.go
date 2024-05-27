package main

import (
	"errors"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
)

type Projucer struct {
	path        string
	buildsPath  string
	projectPath string
	binaryPath  string
}

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

func (p *Projucer) build() (bool, error) {
	if found, _ := fileExists(p.projectPath); found {
		return build(p.projectPath)
	}
	return false, errors.New("unable to find Projucer IDE project")
}

func (p *Projucer) open(projectFile string) (bool, error) {
	p.build()
	cmd := exec.Command(p.binaryPath, projectFile)
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	return true, nil
}

func (p *Projucer) export(projectFile string) (bool, error) {
	p.build()
	cmd := exec.Command(p.binaryPath, "--resave", projectFile)
	return run(cmd)
}

func (p *Projucer) cleanBuildArtefacts() (bool, error) {
	if err := os.RemoveAll(path.Join(p.buildsPath, buildArtefactsPath)); err != nil {
		return false, err
	}
	return true, nil
}

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
