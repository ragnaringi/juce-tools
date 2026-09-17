package cli

import (
	"fmt"
	"path/filepath"
)

func runCommand(workingDirectory string, args []string) error {
	switch args[0] {
	case "up":
		return runUp(workingDirectory)
	case "clean":
		return runClean(workingDirectory, args[1:])
	case "export":
		return runExport(workingDirectory)
	case "code":
		return runCode(workingDirectory)
	case "build":
		return runBuild(workingDirectory)
	default:
		return fmt.Errorf("unknown command: %s", args[0])
	}
}

func loadContext(workingDirectory string) (*JUCEProject, *JUCE) {
	project := NewProject(workingDirectory)
	success("Found %s", filepath.Base(project.jucerFilePath))

	juce := NewJUCE(workingDirectory)

	if ok, err := juce.projucer.Build(); !ok {
		fail("Failed to build Projucer: %v", err)
	}

	return project, juce
}

func runUp(workingDirectory string) error {
	project, juce := loadContext(workingDirectory)

	success("Opening %s", project.name)
	if _, err := juce.projucer.Open(project.jucerFilePath); err != nil {
		return err
	}

	return nil
}

func runClean(workingDirectory string, args []string) error {
	project, juce := loadContext(workingDirectory)

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
	project, juce := loadContext(workingDirectory)

	success("Exporting %s", project.name)
	if _, err := juce.projucer.Export(project.jucerFilePath); err != nil {
		return err
	}

	return nil
}

func runCode(workingDirectory string) error {
	project, juce := loadContext(workingDirectory)

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
	project, juce := loadContext(workingDirectory)

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
