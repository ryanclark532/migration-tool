package utils

import (
	"os"
	"path/filepath"
	"regexp"
)

func CrawlDir(startingDir string) []string {
	var fileNames []string
	filepath.Walk(startingDir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		pattern := `^\d{8}-[\w\s]+\.sql$`
		re := regexp.MustCompile(pattern)

		if !re.MatchString(info.Name()) {
			return nil
		}
		fileNames = append(fileNames, info.Name())
		return nil
	})
	return fileNames
}
