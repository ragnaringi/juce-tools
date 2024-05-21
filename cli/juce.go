package main

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path"
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

		fmt.Println("Projucer binary not found. Compiling")
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
			cmd := exec.Command("xcodebuild", "-project", projucerProject, "-scheme", "Projucer - App", "-configuration", "Release", "-jobs", "8")
			return run(cmd)
		}
	}
	return false, errors.New("unable to find Projucer Xcode project")
}

func cleanProjucerBuildArtefacts(directory string) (bool, error) {
	if jucePath := findJuceDirectory(directory); jucePath != "" {
		buildsDir := filepath.Join(jucePath, "extras", "Projucer", "Builds")
		dir, err := os.ReadDir(buildsDir)
		if err != nil {
			return false, err
		}
		for _, d := range dir {
			os.RemoveAll(path.Join(buildsDir, d.Name(), "build"))
		}
		return true, nil
	}
	return false, errors.New("unable to find JUCE directory")
}
