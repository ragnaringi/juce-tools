package cli

import (
	"fmt"
	"path/filepath"
)

type JUCECommandContext struct {
	workingDirectory string
	project          *ProjucerProject
	juce             *JUCE
	projucer         *Projucer
}

func loadJUCECommandContext(workingDirectory string) (*JUCECommandContext, error) {
	project, err := NewProject(workingDirectory)
	if err != nil {
		return nil, err
	}

	success("Found %s", filepath.Base(project.jucerFilePath))

	juce, err := NewJUCE(workingDirectory)
	if err != nil {
		return nil, err
	}

	if ok, err := juce.projucer.Build(); !ok {
		return nil, fmt.Errorf("failed to build Projucer: %w", err)
	}

	return &JUCECommandContext{
		workingDirectory: workingDirectory,
		project:          project,
		juce:             juce,
		projucer:         juce.projucer,
	}, nil
}

func runUp(workingDirectory string) error {
	ctx, err := loadJUCECommandContext(workingDirectory)
	if err != nil {
		return err
	}

	success("Opening %s", ctx.project.name)
	if _, err := ctx.projucer.Open(ctx.project.jucerFilePath); err != nil {
		return err
	}

	return nil
}

func runClean(workingDirectory string, args []string) error {
	ctx, err := loadJUCECommandContext(workingDirectory)
	if err != nil {
		return err
	}

	success("Cleaning %s", ctx.project.name)
	if len(args) > 0 && args[0] == "--all" {
		if _, err := ctx.projucer.Clean(); err != nil {
			return err
		}
	}

	if _, err := ctx.project.Clean(); err != nil {
		return err
	}

	return nil
}

func runExport(workingDirectory string) error {
	ctx, err := loadJUCECommandContext(workingDirectory)
	if err != nil {
		return err
	}

	success("Exporting %s", ctx.project.name)
	if _, err := ctx.projucer.Export(ctx.project.jucerFilePath); err != nil {
		return err
	}

	return nil
}

func runCode(workingDirectory string) error {
	ctx, err := loadJUCECommandContext(workingDirectory)
	if err != nil {
		return err
	}

	success("Exporting %s", ctx.project.name)
	if _, err := ctx.projucer.Export(ctx.project.jucerFilePath); err != nil {
		return err
	}

	success("Opening %s", ctx.project.name+ideProjectExtension)
	if _, err := ctx.project.Open(); err != nil {
		return err
	}

	return nil
}

func runBuild(workingDirectory string) error {
	ctx, err := loadJUCECommandContext(workingDirectory)
	if err != nil {
		return err
	}

	success("Exporting %s", ctx.project.name)
	if _, err := ctx.projucer.Export(ctx.project.jucerFilePath); err != nil {
		return err
	}

	success("Building %s", ctx.project.name+ideProjectExtension)
	if _, err := ctx.project.Build(); err != nil {
		return err
	}

	return nil
}
