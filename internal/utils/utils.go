package utils

import (
	"os"
	"path/filepath"
)

func CrawlDir(startingDir string) ([]string, error) {
	var fileNames []string
	err := filepath.Walk(startingDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			relativePath, err := filepath.Rel(startingDir, path)
			if err != nil {
				return err
			}
			fileNames = append(fileNames, relativePath)
		}
		return nil
	})
	return fileNames, err
}
