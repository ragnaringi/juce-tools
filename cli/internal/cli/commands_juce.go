package cli

import (
	"flag"
	"fmt"
	"path/filepath"
)

type JUCECommandContext struct {
	workingDirectory string
	project          *ProjucerProject
	juce             *JUCEInstallation
	projucer         *Projucer
}

func loadJUCECommandContext(workingDirectory string) (*JUCECommandContext, error) {
	project, err := NewProject(workingDirectory)
	if err != nil {
		return nil, err
	}

	success("Found %s", filepath.Base(project.jucerFilePath))

	juce, err := NewJUCEInstallation(workingDirectory)
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
	success("Cleaning Builds")

	buildsPath := filepath.Join(workingDirectory, "Builds")
	if err := cleanDirectoryContents(buildsPath); err != nil {
		return err
	}

	if len(args) > 0 && args[0] == "--all" {
		juce, err := NewJUCEInstallation(workingDirectory)
		if err != nil {
			return err
		}

		if _, err := juce.projucer.Clean(); err != nil {
			return err
		}
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

func runExporters(workingDirectory string) error {
	project, err := NewProject(workingDirectory)
	if err != nil {
		return err
	}

	exporters, err := project.AvailableExporters()
	if err != nil {
		return err
	}

	if len(exporters) == 0 {
		warn("No exporters found for %s", project.name)
		return nil
	}

	success("Available exporters for %s:", project.name)
	for _, exporter := range exporters {
		notice("  %s", exporter)
	}

	return nil
}

func runCode(workingDirectory string, args []string) error {
	codeFlags := flag.NewFlagSet("code", flag.ContinueOnError)
	exporter := codeFlags.String("exporter", "", "Projucer exporter to open")
	if err := codeFlags.Parse(args); err != nil {
		return err
	}

	if !canOpenExporter(*exporter) {
		return unsupportedExporterError("open", *exporter)
	}

	ctx, err := loadJUCECommandContext(workingDirectory)
	if err != nil {
		return err
	}

	success("Exporting %s", ctx.project.name)
	if _, err := ctx.projucer.Export(ctx.project.jucerFilePath); err != nil {
		return err
	}

	if *exporter != "" {
		success("Opening %s exporter for %s", *exporter, ctx.project.name)
	} else {
		success("Opening %s", ctx.project.name+ideProjectExtension)
	}

	if _, err := ctx.project.Open(*exporter); err != nil {
		return err
	}

	return nil
}

func runBuild(workingDirectory string, args []string) error {
	buildFlags := flag.NewFlagSet("build", flag.ContinueOnError)
	exporter := buildFlags.String("exporter", "", "Projucer exporter to use")
	if err := buildFlags.Parse(args); err != nil {
		return err
	}

	if !canBuildExporter(*exporter) {
		return unsupportedExporterError("build", *exporter)
	}

	ctx, err := loadJUCECommandContext(workingDirectory)
	if err != nil {
		return err
	}

	success("Exporting %s", ctx.project.name)
	if _, err := ctx.projucer.Export(ctx.project.jucerFilePath); err != nil {
		return err
	}

	success("Building %s", ctx.project.name+ideProjectExtension)
	if _, err := ctx.project.Build(*exporter); err != nil {
		return err
	}

	return nil
}
