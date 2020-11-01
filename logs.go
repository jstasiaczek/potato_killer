package main

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func CheckLogs(directory, prefix string) (int, int64) {
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		log.Panicln(err)
	}
	files = filterDirs(files, prefix)
	lastMod := int64(0)
	for _, file := range files {
		if lastMod < file.ModTime().Unix() {
			lastMod = file.ModTime().Unix()
		}
	}
	return len(files), lastMod
}

func filterDirs(files []os.FileInfo, prefix string) []os.FileInfo {
	if prefix == "" {
		return files
	}
	filtered := make([]os.FileInfo, 0, 1)
	for _, file := range files {
		if strings.Index(file.Name(), prefix) == 0 {
			filtered = append(filtered, file)
		}
	}
	return filtered
}
