package utils

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"
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

func ReadSQLFile(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var content strings.Builder
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		content.WriteString(scanner.Text())
		content.WriteString("\n")
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return content.String(), nil
}
