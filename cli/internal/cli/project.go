package cli

import (
	"errors"
	"os"
	"path"
	"path/filepath"
)

type ProjucerProject struct {
	directory     string // project root directory
	buildsPath    string // Builds directory
	jucerFilePath string // Projucer project file
	buildFilePath string // IDE project file
	name          string // project name
	schemeName    string // the project scheme name
}

func NewProject(directory string) (*ProjucerProject, error) {
	projectFile := findJucerProjectFile(directory)
	if projectFile == nil {
		return nil, errors.New("no JUCE projects found in directory")
	}
	name := fileNameWithoutExtension(projectFile.Name())
	buildsPath := path.Join(directory, "Builds")
	jucerFilePath := path.Join(directory, projectFile.Name())
	buildFilePath := filepath.Join(buildsPath, platformIdentifier, name+ideProjectExtension)
	return &ProjucerProject{
		directory,
		buildsPath,
		jucerFilePath,
		buildFilePath,
		name,
		"All",
	}, nil
}

func (p *ProjucerProject) Open() (bool, error) {
	if found, _ := fileExists(p.buildFilePath); found {
		return open(p.buildFilePath)
	}
	return false, errors.New("unable to find build project file")
}

func (p *ProjucerProject) Build() (bool, error) {
	if found, _ := fileExists(p.buildFilePath); found {
		return build(p.buildFilePath, p.name+" - "+p.schemeName)
	}
	return false, errors.New("unable to find build project file")
}

func (p *ProjucerProject) Clean() (bool, error) {
	dir, _ := os.ReadDir(p.buildsPath)
	for _, d := range dir {
		os.RemoveAll(path.Join(p.buildsPath, d.Name()))
	}
	return true, nil
}
