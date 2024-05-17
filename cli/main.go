package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"

	"github.com/nexidian/gocliselect"
)

func findJucerProjectFile(directory string) fs.FileInfo {
	jucerProjects := findFilesWithExtension(directory, ".jucer")

	if len(jucerProjects) == 0 {
		return nil
	} else if len(jucerProjects) > 1 {
		menu := gocliselect.NewMenu("Multiple JUCE projects in directory. Choose one")

		for _, file := range jucerProjects {
			menu.AddItem(file.Name(), file.Name())
		}

		choice := menu.Display()

		for _, file := range jucerProjects {
			if file.Name() == choice {
				return file
			}
		}
	}

	return jucerProjects[0]
}

// ===============================================================================================
func main() {
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	flag.Parse()

	if flag.NArg() == 0 {
		fmt.Println("No arguments found")
	} else if flag.Arg(0) == "up" {
		fmt.Println("Starting JUCE")

		projectFile := findJucerProjectFile(pwd)

		if projectFile == nil {
			fmt.Println("No JUCE project found in directory", pwd)
			return
		}

		if projucerBinary, _ := findProjucerExecutable(pwd); projucerBinary != "" {
			fmt.Println("Opening", projectFile.Name())

			cmd := exec.Command(projucerBinary, projectFile.Name())
			if err := cmd.Start(); err != nil {
				log.Fatal(err)
			}
		}
	} else if flag.Arg(0) == "export" {
		projectFile := findJucerProjectFile(pwd)

		if projectFile == nil {
			fmt.Println("No JUCE project found in directory", pwd)
			return
		}

		fmt.Println("Exporting project", projectFile.Name())

		exportProject(pwd, projectFile)
	}
}
