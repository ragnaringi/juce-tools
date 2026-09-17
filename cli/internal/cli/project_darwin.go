package cli

import (
	"io/fs"
	"path/filepath"

	"github.com/nexidian/gocliselect"
)

const ideProjectExtension string = ".xcodeproj"
const buildArtefactsPath string = "build"

func exportedProjectCandidates(buildsPath string, projectName string) ([]string, error) {
	return []string{
		filepath.Join(buildsPath, "MacOSX", projectName+ideProjectExtension),
	}, nil
}

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
