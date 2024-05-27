package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
)

func run(cmd *exec.Cmd) (bool, error) {
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + string(output))
		return false, err
	}
	return true, nil
}

// ===============================================================================================
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

	project := NewProject(workingDirectory)
	juce := NewJUCE(workingDirectory)

	if flag.Arg(0) == "up" {
		fmt.Println("Opening", project.name)
		juce.projucer.open(project.jucerFilePath)
	} else if flag.Arg(0) == "clean" {
		if flag.Arg(1) == "--all" {
			juce.projucer.cleanBuildArtefacts()
		}
		project.clean()
	} else if flag.Arg(0) == "export" {
		fmt.Println("Exporting", project.name)
		if _, err := juce.projucer.export(project.jucerFilePath); err != nil {
			panic(err)
		}
	} else if flag.Arg(0) == "code" {
		fmt.Println("Exporting", project.name)
		if _, err := juce.projucer.export(project.jucerFilePath); err != nil {
			panic(err)
		}
		fmt.Println("Opening", project.name+ideProjectExtension)
		if _, err := project.open(); err != nil {
			panic(err)
		}
	} else if flag.Arg(0) == "build" {
		fmt.Println("Exporting", project.name)
		if _, err := juce.projucer.export(project.jucerFilePath); err != nil {
			panic(err)
		}
		fmt.Println("Building", project.name+ideProjectExtension)
		if _, err := project.build(); err != nil {
			panic(err)
		}
	}
}
