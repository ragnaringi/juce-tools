package main

import (
	"io/fs"
	"log"
	"os/exec"
)

func exportProject(rootDirectory string, projectFile fs.FileInfo) {
	if projucerBinary, _ := findProjucerExecutable(rootDirectory); projucerBinary != "" {
		cmd := exec.Command(projucerBinary, "--resave", projectFile.Name())
		if err := cmd.Start(); err != nil {
			log.Fatal(err)
		}
	}
}
