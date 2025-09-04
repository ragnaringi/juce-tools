package main

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// relativePath converts an absolute path into a relative path
// based on the current working directory. If it fails, it falls
// back to returning the absolute path instead.
func relativePath(baseDir, targetPath string) string {
	rel, err := filepath.Rel(baseDir, targetPath)
	if err != nil {
		return targetPath // fallback to absolute path
	}

	// Normalize to avoid weird "./" prefixes
	if rel == "." {
		return filepath.Base(targetPath)
	}
	return rel
}

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

func fileNameWithoutExtension(fileName string) string {
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}

func findFilesWithExtension(directory, ext string) (ret []fs.FileInfo) {
	for _, file := range listFiles(directory) {
		if filepath.Ext(file.Name()) == ext {
			ret = append(ret, file)
		}
	}
	return
}
