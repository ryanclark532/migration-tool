package up

import (
	"fmt"
	"os"
	"ryanclark532/migration-tool/internal/common"
	"ryanclark532/migration-tool/internal/down"
	"ryanclark532/migration-tool/internal/utils"
	"strings"
	"time"
)

func DoMigration(server common.Server, config common.Config) []error {
	completedFiles, err := common.CompletedFiles(server.GetDB())
	if err != nil {
		errors := []error{err}
		return errors
	}

	var errors []error

	files, err := utils.CrawlDir(config.InputDir)
	if err != nil {
		errors := []error{err}
		return errors
	}

	processedProcs := make(map[string]bool)

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

		original, err := server.GetDatabaseState(config)
		if err != nil {
			errors = append(errors, fmt.Errorf("error processing %s: %s", file, err.Error()))
			continue
		}

		tx, err := server.Begin()
		if err != nil {
			errors = append(errors, fmt.Errorf("error processing %s: %s", file, err.Error()))
			continue
		}

		_, err = tx.Exec(string(contents))
		if err != nil {
			errors = append(errors, fmt.Errorf("error processing %s: %s", file, err.Error()))
			continue
		}

		err = tx.Commit()
		if err != nil {
			errors = append(errors, fmt.Errorf("error processing %s: %s", file, err.Error()))
			continue
		}
		post, err := server.GetDatabaseState(config)
		if err != nil {
			errors = append(errors, fmt.Errorf("error processing %s: %s", file, err.Error()))
			continue
		}

		var builder strings.Builder
		var procBuilder strings.Builder

		down.GetTableDiff(original.Tables, post.Tables, &builder)
		down.GetProcDiff(original.Procs, post.Procs, &procBuilder, processedProcs)

		if builder.Len() != 0 {
			err = os.WriteFile(fmt.Sprintf("%s/%s.down.sql", config.OutputDir, file), []byte(builder.String()), os.ModeAppend)
			if err != nil {
				errors = append(errors, fmt.Errorf("error processing %s: %s", file, err.Error()))
			}

			sqlBatch := fmt.Sprintf("INSERT INTO Migrations(EnterDateTime, Version, FileName) VALUES ('%s', %d, '%s')", time.Now().Format(time.RFC3339), 1, file)
			conn := server.GetDB()
			_, err = conn.Exec(sqlBatch)
			if err != nil {
				errors = append(errors, fmt.Errorf("error processing %s: %s", file, err.Error()))
				continue
			}
		} else if procBuilder.Len() != 0 {
			err = os.WriteFile(fmt.Sprintf("%s/%s.down.sql", config.OutputDir, file), []byte(procBuilder.String()), os.ModeAppend)
			if err != nil {
				errors = append(errors, fmt.Errorf("error processing %s: %s", file, err.Error()))
			}

			sqlBatch := fmt.Sprintf("INSERT INTO Migrations(EnterDateTime, Version, FileName) VALUES ('%s', %d, '%s')", time.Now().Format(time.RFC3339), 1, file)
			conn := server.GetDB()
			_, err = conn.Exec(sqlBatch)
			if err != nil {
				errors = append(errors, fmt.Errorf("error processing %s: %s", file, err.Error()))
				continue
			}
		}
	}

	return errors
}

func DoDryMigration(server common.Server, config common.Config) []error {
	completedfiles, err := common.CompletedFiles(server.GetDB())
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
