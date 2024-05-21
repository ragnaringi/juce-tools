package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path"
	"path/filepath"

	"github.com/nexidian/gocliselect"
)

func run(cmd *exec.Cmd) (bool, error) {
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + string(output))
		return false, err
	}
	return true, nil
}

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

func cleanProject(rootDirectory string) (bool, error) {
	buildsDir := filepath.Join(rootDirectory, "Builds")
	dir, _ := os.ReadDir(buildsDir)
	for _, d := range dir {
		os.RemoveAll(path.Join(buildsDir, d.Name()))
	}
	return true, nil
}

func exportProject(rootDirectory string, projectFile fs.FileInfo) (bool, error) {
	if projucerBinary, _ := findProjucerExecutable(rootDirectory); projucerBinary != "" {
		cmd := exec.Command(projucerBinary, "--resave", projectFile.Name())
		return run(cmd)
	}
	return false, errors.New("unable to find Projucer binary")
}

func openProject(rootDirectory string, projectFile fs.FileInfo) (bool, error) {
	projectName := fileNameWithoutExtension(projectFile.Name())
	xcodeProject := filepath.Join(rootDirectory, "Builds/MacOSX", projectName+".xcodeproj")
	if found, _ := fileExists(xcodeProject); found {
		cmd := exec.Command("open", xcodeProject)
		return run(cmd)
	}
	return false, errors.New("unable to find Xcode project")
}

func buildProject(rootDirectory string, projectFile fs.FileInfo) (bool, error) {
	projectName := fileNameWithoutExtension(projectFile.Name())
	xcodeProject := filepath.Join(rootDirectory, "Builds/MacOSX", projectName+".xcodeproj")
	if found, _ := fileExists(xcodeProject); found {
		cmd := exec.Command("xcodebuild", "-project", xcodeProject, "-scheme", projectName+" - All", "-configuration", "Release", "-jobs", "8")
		return run(cmd)
	}
	return false, errors.New("unable to find Xcode project")
}
