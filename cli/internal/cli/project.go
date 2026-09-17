package cli

import (
	"errors"
	"fmt"
	"os"
	"path"
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

	return &ProjucerProject{
		directory,
		buildsPath,
		jucerFilePath,
		"",
		name,
		"All",
	}, nil
}

func findExportedProjectFile(buildsPath string, projectName string) (string, error) {
	candidates, err := exportedProjectCandidates(buildsPath, projectName)
	if err != nil {
		return "", err
	}

	for _, candidate := range candidates {
		if found, err := fileExists(candidate); err != nil {
			return "", fmt.Errorf("checking exported project file: %w", err)
		} else if found {
			return candidate, nil
		}
	}

	return "", fmt.Errorf("unable to find exported project file for %q in %s", projectName, buildsPath)
}

func (p *ProjucerProject) resolveBuildFilePath() error {
	if p.buildFilePath != "" {
		return nil
	}

	buildFilePath, err := findExportedProjectFile(p.buildsPath, p.name)
	if err != nil {
		return err
	}

	p.buildFilePath = buildFilePath
	return nil
}

func (p *ProjucerProject) Open() (bool, error) {
	if err := p.resolveBuildFilePath(); err != nil {
		return false, err
	}

	if found, _ := fileExists(p.buildFilePath); found {
		return open(p.buildFilePath)
	}
	return false, errors.New("unable to find build project file")
}

func (p *ProjucerProject) Build() (bool, error) {
	if err := p.resolveBuildFilePath(); err != nil {
		return false, err
	}

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
