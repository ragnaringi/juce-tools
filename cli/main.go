package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func run(cmd *exec.Cmd) (bool, error) {
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + string(output))
		return false, err
	}
	return true, nil
}

func main() {
	flag.Parse()

	if flag.NArg() == 0 {
		fmt.Println("No arguments found")
		os.Exit(1)
	}

	workingDirectory, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	// Find the project
	project := NewProject(workingDirectory)
	success("Found %s", filepath.Base(project.jucerFilePath))

	// Locate JUCE
	juce := NewJUCE(workingDirectory)

	// Ensure Projucer is built before using it
	if ok, err := juce.projucer.Build(); !ok {
		fail("Failed to build Projucer: %v", err)
	}

	// Dispatch commands
	switch flag.Arg(0) {
	case "up":
		success("Opening %s", project.name)
		if _, err := juce.projucer.Open(project.jucerFilePath); err != nil {
			panic(err)
		}
	case "clean":
		if flag.Arg(1) == "--all" {
			if _, err := juce.projucer.Clean(); err != nil {
				panic(err)
			}
		}
		if _, err := project.Clean(); err != nil {
			panic(err)
		}
	case "export", "code":
		success("Exporting %s", project.name)
		if _, err := juce.projucer.Export(project.jucerFilePath); err != nil {
			panic(err)
		}

		if flag.Arg(0) == "code" {
			success("Opening %s", project.name+ideProjectExtension)
			if _, err := project.Open(); err != nil {
				panic(err)
			}
		}
	case "build":
		success("Exporting %s", project.name)
		if _, err := juce.projucer.Export(project.jucerFilePath); err != nil {
			panic(err)
		}

		success("Building %s", project.name+ideProjectExtension)
		if _, err := project.Build(); err != nil {
			panic(err)
		}
	default:
		warn("Unknown command: %s", flag.Arg(0))
		os.Exit(1)
	}
}
