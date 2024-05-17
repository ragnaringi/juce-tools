package main

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os/exec"
	"path/filepath"
	"strings"
)

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
