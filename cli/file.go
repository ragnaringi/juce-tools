package main

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

func fileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func listFiles(directory string) []fs.FileInfo {
	files, _ := os.Open(directory)
	fileInfo, error := files.Readdir(-1)
	if error != nil {
		log.Fatal(error)
	}
	defer files.Close()
	return fileInfo
}

func findFilesWithExtension(directory, ext string) (ret []fs.FileInfo) {
	for _, file := range listFiles(directory) {
		if filepath.Ext(file.Name()) == ext {
			ret = append(ret, file)
		}
	}
	return
}
