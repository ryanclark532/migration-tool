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
		if info.IsDir() {
			return nil
		}
		fileNames = append(fileNames, info.Name())
		return nil
	})
	return fileNames, err
}
