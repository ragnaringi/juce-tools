package main

import (
	"io/fs"

	"github.com/nexidian/gocliselect"
)

const platformIdentifier string = "MacOSX"
const ideProjectExtension string = ".xcodeProject"
const buildArtefactsPath string = "build"

func findJucerProjectFile(directory string) fs.FileInfo {
	projectFiles := findFilesWithExtension(directory, ".jucer")

	if len(projectFiles) == 0 {
		return nil
	} else if len(projectFiles) > 1 {
		menu := gocliselect.NewMenu("Multiple JUCE projects in directory. Choose one")

		for _, file := range projectFiles {
			menu.AddItem(file.Name(), file.Name())
		}

		choice := menu.Display()

		for _, file := range projectFiles {
			if file.Name() == choice {
				return file
			}
		}
	}

	return projectFiles[0]
}
