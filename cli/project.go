package main

import (
	"errors"
	"os"
	"path"
	"path/filepath"
)

type JUCEProject struct {
	directory     string // project root directory
	buildsPath    string // Builds directory
	jucerFilePath string // Projucer project file
	buildFilePath string // IDE project file
	name          string // project name
}

func NewProject(directory string) *JUCEProject {
	projectFile := findJucerProjectFile(directory)
	if projectFile == nil {
		panic("No JUCE projects found in directory")
	}
	name := fileNameWithoutExtension(projectFile.Name())
	buildsPath := path.Join(directory, "Builds")
	jucerFilePath := path.Join(directory, projectFile.Name())
	buildFilePath := filepath.Join(buildsPath, platformIdentifier, name+ideProjectExtension)
	return &JUCEProject{
		directory,
		buildsPath,
		jucerFilePath,
		buildFilePath,
		name,
	}
}

func (p *JUCEProject) open() (bool, error) {
	if found, _ := fileExists(p.buildFilePath); found {
		return open(p.buildFilePath)
	}
	return false, errors.New("unable to find build project file")
}

func (p *JUCEProject) build() (bool, error) {
	if found, _ := fileExists(p.buildFilePath); found {
		return build(p.buildFilePath)
	}
	return false, errors.New("unable to find build project file")
}

func (p *JUCEProject) clean() (bool, error) {
	dir, _ := os.ReadDir(p.buildsPath)
	for _, d := range dir {
		os.RemoveAll(path.Join(p.buildsPath, d.Name()))
	}
	return true, nil
}
