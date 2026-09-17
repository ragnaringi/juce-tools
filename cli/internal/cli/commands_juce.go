package cli

import (
	"fmt"
	"path/filepath"
)

func loadContext(workingDirectory string) (*JUCEProject, *JUCE, error) {
	project, err := NewProject(workingDirectory)
	if err != nil {
		return nil, nil, err
	}

	success("Found %s", filepath.Base(project.jucerFilePath))

	juce, err := NewJUCE(workingDirectory)
	if err != nil {
		return nil, nil, err
	}

	if ok, err := juce.projucer.Build(); !ok {
		return nil, nil, fmt.Errorf("failed to build Projucer: %w", err)
	}

	return project, juce, nil
}

func runUp(workingDirectory string) error {
	project, juce, err := loadContext(workingDirectory)
	if err != nil {
		return err
	}

	success("Opening %s", project.name)
	if _, err := juce.projucer.Open(project.jucerFilePath); err != nil {
		return err
	}

	return nil
}

func runClean(workingDirectory string, args []string) error {
	project, juce, err := loadContext(workingDirectory)
	if err != nil {
		return err
	}

	if len(args) > 0 && args[0] == "--all" {
		if _, err := juce.projucer.Clean(); err != nil {
			return err
		}
	}

	if _, err := project.Clean(); err != nil {
		return err
	}

	return nil
}

func runExport(workingDirectory string) error {
	project, juce, err := loadContext(workingDirectory)
	if err != nil {
		return err
	}

	success("Exporting %s", project.name)
	if _, err := juce.projucer.Export(project.jucerFilePath); err != nil {
		return err
	}

	return nil
}

func runCode(workingDirectory string) error {
	project, juce, err := loadContext(workingDirectory)
	if err != nil {
		return err
	}

	success("Exporting %s", project.name)
	if _, err := juce.projucer.Export(project.jucerFilePath); err != nil {
		return err
	}

	success("Opening %s", project.name+ideProjectExtension)
	if _, err := project.Open(); err != nil {
		return err
	}

	return nil
}

func runBuild(workingDirectory string) error {
	project, juce, err := loadContext(workingDirectory)
	if err != nil {
		return err
	}

	success("Exporting %s", project.name)
	if _, err := juce.projucer.Export(project.jucerFilePath); err != nil {
		return err
	}

	success("Building %s", project.name+ideProjectExtension)
	if _, err := project.Build(); err != nil {
		return err
	}

	return nil
}
