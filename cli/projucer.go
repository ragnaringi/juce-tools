package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
)

type Projucer struct {
	path        string
	projectPath string
	binaryPath  string
}

func NewProjucer(jucePath string) *Projucer {
	path := filepath.Join(jucePath, "extras", "Projucer")
	buildsPath := filepath.Join(path, "Builds", "MacOSX")
	projectPath := filepath.Join(buildsPath, "Projucer.xcodeproj")
	binaryPath := filepath.Join(buildsPath, "build", "Release", "Projucer.app", "Contents", "MacOS", "Projucer")
	return &Projucer{
		path,
		projectPath,
		binaryPath,
	}
}

func (p *Projucer) compile() (bool, error) {
	if found, _ := fileExists(p.projectPath); found {
		fmt.Println("Compiling Projucer")
		cmd := exec.Command("xcodebuild", "-project", p.projectPath, "-scheme", "Projucer - App", "-configuration", "Release", "-jobs", "8")
		return run(cmd)
	}
	return false, errors.New("unable to find Projucer Xcode project")
}

func (p *Projucer) open(projectFile string) (bool, error) {
	p.compile()
	cmd := exec.Command(p.binaryPath, projectFile)
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	return true, nil
}

func (p *Projucer) export(projectFile string) (bool, error) {
	p.compile()
	cmd := exec.Command(p.binaryPath, "--resave", projectFile)
	return run(cmd)
}

func (p *Projucer) cleanBuildArtefacts() (bool, error) {
	buildsDir := filepath.Join(p.path, "Builds")
	dir, err := os.ReadDir(buildsDir)
	if err != nil {
		return false, err
	}
	for _, d := range dir {
		os.RemoveAll(path.Join(buildsDir, d.Name(), "build"))
	}
	return true, nil
}
