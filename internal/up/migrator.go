package up

import (
	"database/sql"
	"fmt"
	"os"
	"ryanclark532/migration-tool/internal/common"
	"ryanclark532/migration-tool/internal/utils"
	"strings"
	"time"
)

func DoMigration(conn *sql.DB, version int, config common.Config) []error {
	completedfiles, err := CompletedFiles(conn)
	if err != nil {
		errors := []error{err}
		return errors
	}

	var builder strings.Builder
	var errors []error

	files, err := utils.CrawlDir(config.InputDir)
	for _, file := range files {
		_, ex := completedfiles[file]
		if ex == true {
			fmt.Printf("skipping %s, as its already been run\n", file)
			continue
		}

		contents, err := os.ReadFile(fmt.Sprintf("%s/%s", config.InputDir, file))
		if err != nil {
			errors = append(errors, fmt.Errorf("error processing %s: %s", file, err.Error()))
			continue
		}

		tx, _ := conn.Begin()
		_, err = tx.Exec(string(contents))
		if err != nil {
			errors = append(errors, fmt.Errorf("error processing %s: %s", file, err.Error()))
			continue
		}

		err = tx.Rollback()
		if err != nil {
			errors = append(errors, fmt.Errorf("error processing %s: %s", file, err.Error()))
			continue
		}

		sqlBatch := fmt.Sprintf("INSERT INTO Migrations(EnterDateTime, Version, FileName) VALUES ('%s', %d, '%s')", time.Now(), version, file)
		_, err = conn.Exec(sqlBatch)
		if err != nil {
			errors = append(errors, fmt.Errorf("error processing %s: %s", file, err.Error()))
			continue
		}

		builder.WriteString(string(contents))
		fmt.Printf("processed %s\n", file)
	}

	if len(errors) != 0 {
		return errors
	}

	if builder.Len() != 0 {
		err := os.WriteFile(fmt.Sprintf("%s/up/%d-up.sql", config.OutputDir, version), []byte(builder.String()), os.ModeAppend)
		if err != nil {
			panic(err)
		}
		_, err = conn.Exec(builder.String())
		if err != nil {
			panic(err)
		}
	}
	fmt.Printf("Migration Successfull, Version: %d\n", version)
	return nil
}

func DoDryMigration(conn *sql.DB, version int, config common.Config) []error {
	completedfiles, err := CompletedFiles(conn)
	if err != nil {
		errors := []error{err}
		return errors
	}

	var builder strings.Builder
	var errors []error

	files, err := utils.CrawlDir(config.InputDir)
	for _, file := range files {
		_, ex := completedfiles[file]
		if ex == true {
			fmt.Printf("skipping %s, as its already been run\n", file)
			continue
		}

		contents, err := os.ReadFile(fmt.Sprintf("%s/%s", config.InputDir, file))
		if err != nil {
			errors = append(errors, fmt.Errorf("error processing %s: %s", file, err.Error()))
			continue
		}

		tx, _ := conn.Begin()
		_, err = tx.Exec(string(contents))
		if err != nil {
			errors = append(errors, fmt.Errorf("error processing %s: %s", file, err.Error()))
			continue
		}

		err = tx.Rollback()
		if err != nil {
			errors = append(errors, fmt.Errorf("error processing %s: %s", file, err.Error()))
			continue
		}

		sqlBatch := fmt.Sprintf("INSERT INTO Migrations(EnterDateTime, Version, FileName) VALUES ('%s', %d, '%s')", time.Now(), version, file)
		_, err = conn.Exec(sqlBatch)
		if err != nil {
			errors = append(errors, fmt.Errorf("error processing %s: %s", file, err.Error()))
			continue
		}

		builder.WriteString(string(contents))
		fmt.Printf("processed %s\n", file)
	}

	if len(errors) != 0 {
		return errors
	}

	fmt.Printf("Migration Successfull, Version: %d\n", version)
	return nil
}

func CompletedFiles(conn *sql.DB) (map[string]bool, error) {
	sqlBatch := ("SELECT FileName From Migrations")
	fileNames := make(map[string]bool)
	rows, err := conn.Query(sqlBatch)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var fileName string
		err = rows.Scan(&fileName)
		if err != nil {
			return fileNames, err
		}
		fileNames[fileName] = true
	}

	return fileNames, nil
}
