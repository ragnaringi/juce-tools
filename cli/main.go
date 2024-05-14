package main

import (
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

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

func findJuceDirectory(directory string) (ret string) {
	err := filepath.Walk(directory,
		func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if strings.HasSuffix(path, "/JUCE") {
				if found, _ := fileExists(filepath.Join(path, "extras")); found {
					ret = path
				}
			}

			return nil
		})

	if err != nil {
		log.Println(err)
	}

	return
}

func findProjucerExecutable(directory string) (string, error) {
	if path := findJuceDirectory(directory); path != "" {
		projucerBinary := filepath.Join(path, "extras/Projucer/Builds/MacOSX/build/Release/Projucer.app/Contents/MacOS/Projucer")

		if found, _ := fileExists(projucerBinary); found {
			return projucerBinary, nil
		}

		if success, _ := compileProjucer(directory); success {
			return projucerBinary, nil
		}
	}
	return "", errors.New("unable to find Projucer binary")
}

func compileProjucer(directory string) (bool, error) {
	if path := findJuceDirectory(directory); path != "" {
		projucerProject := filepath.Join(path, "extras/Projucer/Builds/MacOSX/Projucer.xcodeproj")

		if found, _ := fileExists(projucerProject); found {
			fmt.Println("Compiling Projucer")

			cmd := exec.Command("xcodebuild", "-project", projucerProject, "-scheme", "Projucer - App", "-configuration", "Release", "-jobs", "8")
			if err := cmd.Run(); err != nil {
				log.Fatal(err)
			}

			return true, nil
		}
	}

	return false, errors.New("unable to find Projucer Xcode project")
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
	}
}
