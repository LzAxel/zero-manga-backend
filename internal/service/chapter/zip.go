package chapter

import (
	"archive/zip"
	"path/filepath"
	"slices"
)

func countValidImagesInZip(files []*zip.File) uint {
	var count uint = 0
	for _, file := range files {
		if isValidPageExtensions(file.Name) {
			count++
		}
	}

	return count
}

func isValidPageExtensions(filename string) bool {
	return slices.Contains(validPageExtensions, filepath.Ext(filename))
}
