package main

import (
	"errors"
	"io/fs"
	"os"
	"os/exec"
	"path"
	"path/filepath"

	"github.com/nexidian/gocliselect"
)

type JUCEProject struct {
	directory string
	filePath  string
	name      string
}

func NewProject(directory string) *JUCEProject {
	projectFile := findProjectFile(directory)
	if projectFile == nil {
		panic("No JUCE projects found in directory")
	}
	return &JUCEProject{
		directory: directory,
		filePath:  path.Join(directory, projectFile.Name()),
		name:      fileNameWithoutExtension(projectFile.Name()),
	}
}

func (p *JUCEProject) open() (bool, error) {
	xcodeProject := filepath.Join(p.directory, "Builds", "MacOSX", p.name+".xcodeproj")
	if found, _ := fileExists(xcodeProject); found {
		cmd := exec.Command("open", xcodeProject)
		return run(cmd)
	}
	return false, errors.New("unable to find Xcode project")
}

func (p *JUCEProject) build() (bool, error) {
	xcodeProject := filepath.Join(p.directory, "Builds", "MacOSX", p.name+".xcodeproj")
	if found, _ := fileExists(xcodeProject); found {
		cmd := exec.Command("xcodebuild", "-project", xcodeProject, "-scheme", p.name+" - All", "-configuration", "Release", "-jobs", "8")
		return run(cmd)
	}
	return false, errors.New("unable to find Xcode project")
}

func (p *JUCEProject) clean() (bool, error) {
	buildsDir := filepath.Join(p.directory, "Builds")
	dir, _ := os.ReadDir(buildsDir)
	for _, d := range dir {
		os.RemoveAll(path.Join(buildsDir, d.Name()))
	}
	return true, nil
}

func findProjectFile(directory string) fs.FileInfo {
	projectFiles := findFilesWithExtension(directory, ".jucer")

	if len(projectFiles) == 0 {
		return nil
	} else if len(projectFiles) > 1 {
		menu := gocliselect.NewMenu("Multiple JUCE projects in directory. Choose one")

		for _, file := range projectFiles {
			menu.AddItem(file.Name(), file.Name())
		}

		choice := menu.Display()

		for _, file := range projectFiles {
			if file.Name() == choice {
				return file
			}
		}
	}

	return projectFiles[0]
}
