package main

import (
	"io/fs"
)

const platformIdentifier string = "VisualStudio2022"
const ideProjectExtension string = ".sln"
const buildArtefactsPath string = "x64"

func findJucerProjectFile(directory string) fs.FileInfo {
	projectFiles := findFilesWithExtension(directory, ".jucer")
	return projectFiles[0]
}
