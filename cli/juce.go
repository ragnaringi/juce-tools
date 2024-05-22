package main

import (
	"io/fs"
	"log"
	"path/filepath"
	"strings"
)

type JUCE struct {
	path     string
	projucer *Projucer
}

func NewJUCE(rootDirectory string) *JUCE {
	jucePath := findJuceDirectory(rootDirectory)
	if exists, _ := fileExists(jucePath); !exists {
		panic("No JUCE installation found in directory")
	}
	projucer := NewProjucer(jucePath)
	return &JUCE{
		path:     jucePath,
		projucer: projucer,
	}
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
