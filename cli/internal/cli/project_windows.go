package cli

import (
	"io/fs"
	"path/filepath"
	"sort"
)

const ideProjectExtension string = ".sln"
const buildArtefactsPath string = "x64"

func exportedProjectCandidates(buildsPath string, projectName string) ([]string, error) {
	matches, err := filepath.Glob(filepath.Join(buildsPath, "VisualStudio*", projectName+ideProjectExtension))
	if err != nil {
		return nil, err
	}

	sort.Sort(sort.Reverse(sort.StringSlice(matches)))

	return matches, nil
}

func findJucerProjectFile(directory string) fs.FileInfo {
	projectFiles := findFilesWithExtension(directory, ".jucer")
	if len(projectFiles) == 0 {
		return nil
	}
	return projectFiles[0]
}
