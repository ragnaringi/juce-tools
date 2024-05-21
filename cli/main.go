package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
)

// ===============================================================================================
func main() {
	flag.Parse()

	if flag.NArg() == 0 {
		fmt.Println("No arguments found")
		os.Exit(1)
	}

	workingDirectory, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	projectFile := findJucerProjectFile(workingDirectory)
	if projectFile == nil {
		fmt.Println("No JUCE project found in directory", workingDirectory)
		os.Exit(1)
	}

	projectName := fileNameWithoutExtension(projectFile.Name())

	juceDirectory := findJuceDirectory(workingDirectory)
	if juceDirectory == "" {
		fmt.Println("No JUCE installation found")
		os.Exit(1)
	}

	if flag.Arg(0) == "up" {
		if projucerBinary, _ := findProjucerExecutable(workingDirectory); projucerBinary != "" {
			fmt.Println("Opening", projectFile.Name())

			cmd := exec.Command(projucerBinary, projectFile.Name())
			if err := cmd.Start(); err != nil {
				log.Fatal(err)
			}
		}
	} else if flag.Arg(0) == "clean" {
		if flag.Arg(1) == "--all" {
			cleanProjucerBuildArtefacts(workingDirectory)
		}
		cleanProject(workingDirectory)
	} else if flag.Arg(0) == "export" {
		fmt.Println("Exporting", projectName)
		if _, err := exportProject(workingDirectory, projectFile); err != nil {
			log.Fatal(err)
		}
	} else if flag.Arg(0) == "code" {
		fmt.Println("Exporting", projectName)
		if _, err := exportProject(workingDirectory, projectFile); err != nil {
			log.Fatal(err)
		}
		fmt.Println("Opening", projectName+".xcodeproj")
		if _, err := openProject(workingDirectory, projectFile); err != nil {
			log.Fatal(err)
		}
	} else if flag.Arg(0) == "build" {
		fmt.Println("Exporting", projectName)
		if _, err := exportProject(workingDirectory, projectFile); err != nil {
			log.Fatal(err)
		}
		fmt.Println("Building", projectName+".xcodeproj")
		if _, err := buildProject(workingDirectory, projectFile); err != nil {
			log.Fatal(err)
		}
	}
}
