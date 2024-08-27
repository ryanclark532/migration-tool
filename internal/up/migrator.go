package up

import (
	"database/sql"
	"fmt"
	"os"
	"ryanclark532/migration-tool/internal/common"
	"ryanclark532/migration-tool/internal/down"
	"ryanclark532/migration-tool/internal/utils"
	"strings"
	"time"
)

func DoMigration(server common.Server, config common.Config) []error {
	completedFiles, err := CompletedFiles(server.GetDB())
	if err != nil {
		errors := []error{err}
		return errors
	}

	verifiedFiles := make(map[string]string)

	var errors []error

	files, err := utils.CrawlDir(config.InputDir)
	//Pass 1, run all migrations in a transaction to validate them
	for _, file := range files {
		_, ex := completedFiles[file]
		if ex == true {
			fmt.Printf("skipping %s, as its already been run\n", file)
			continue
		}

		contents, err := os.ReadFile(fmt.Sprintf("%s/%s", config.InputDir, file))
		if err != nil {
			errors = append(errors, fmt.Errorf("error processing %s: %s", file, err.Error()))
			continue
		}

		tx, err := server.Begin()
		if err != nil {
			errors = append(errors, fmt.Errorf("error processing %s: %s", file, err.Error()))
		}

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
		verifiedFiles[file] = string(contents)
	}

	//Pass 2 execute migration and generate down
	for file, contents := range verifiedFiles {
		original, err := server.GetDatabaseState(config)
		if err != nil {
			errors = append(errors, fmt.Errorf("error processing %s: %s", file, err.Error()))
		}

		conn := server.GetDB()
		_, err = conn.Exec(contents)
		if err != nil {
			errors = append(errors, fmt.Errorf("error processing %s: %s", file, err.Error()))
		}

		post, err := server.GetDatabaseState(config)
		if err != nil {
			errors = append(errors, fmt.Errorf("error processing %s: %s", file, err.Error()))
		}

		var builder strings.Builder

		down.GetTableDiff(original.Tables, post.Tables, 1, &builder)
		down.GetProcDiff(original.Procs, &builder)

		if builder.Len() != 0 {
			err = os.WriteFile(fmt.Sprintf("%s/down/%s.down.sql", config.OutputDir, file), []byte(builder.String()), os.ModeAppend)
			if err != nil {
				errors = append(errors, fmt.Errorf("error processing %s: %s", file, err.Error()))
			}

			sqlBatch := fmt.Sprintf("INSERT INTO Migrations(EnterDateTime, Version, FileName) VALUES ('%s', %d, '%s')", time.Now(), 1, file)
			_, err = conn.Exec(sqlBatch)
			if err != nil {
				errors = append(errors, fmt.Errorf("error processing %s: %s", file, err.Error()))
				continue
			}

		}

	}
	return nil
}

func DoDryMigration(server common.Server, config common.Config) []error {
	completedfiles, err := CompletedFiles(server.GetDB())
	if err != nil {
		errors := []error{err}
		return errors
	}

	var errors []error

	tx, err := server.Begin()

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

		_, err = tx.Exec(string(contents))
		if err != nil {
			errors = append(errors, fmt.Errorf("error processing %s: %s", file, err.Error()))
			continue
		}

		sqlBatch := fmt.Sprintf("INSERT INTO Migrations(EnterDateTime, Version, FileName) VALUES ('%s', %d, '%s')", time.Now(), 0, file)
		_, err = tx.Exec(sqlBatch)
		if err != nil {
			errors = append(errors, fmt.Errorf("error processing %s: %s", file, err.Error()))
			continue
		}
	}
	err = tx.Rollback()
	if err != nil {
		return []error{err}
	}

	return errors
}

type commonDb interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
}

func CompletedFiles(conn commonDb) (map[string]bool, error) {
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
