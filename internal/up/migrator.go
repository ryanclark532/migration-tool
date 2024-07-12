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

func DoMigration(conn *sql.DB, version int, config common.Config) {
	tables, err := utils.CrawlDir(fmt.Sprintf("%s/tables", config.InputDir))
	if len(tables) == 0 || err != nil {
		fmt.Println("No tables files, or tables folder doesnt exist in input directory")
	}
	for _, file := range tables {
		contents, err := os.ReadFile(fmt.Sprintf("%s/tables/%s", config.InputDir, file))
		if err != nil {
			fmt.Printf("Tables: Error processing %s: %s\n", file, err.Error())
			continue
		}
		_, err = conn.Exec(string(contents))
		if err != nil {
			fmt.Printf("Tables: Error processing %s: %s\n", file, err.Error())
			continue
		}

		fmt.Println("Tables: Processed " + file)
	}

	filesForType, err := GetFilesForType(conn)
	if err != nil {
		fmt.Println(err)
	}

	var builder strings.Builder

	errors := HandleFolder(conn, &builder, "updates", filesForType, version, config)
	if len(errors) != 0 {
		fmt.Println("Skipped Files: ")
		for _, err := range errors {
			fmt.Println(err.Error())
		}
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

}

func HandleFolder(conn *sql.DB, builder *strings.Builder, folderName string, filesForType map[string]string, version int, config common.Config) []error {
	var errors []error
	files, err := utils.CrawlDir(fmt.Sprintf("%s/%s", config.InputDir, folderName))
	if len(files) == 0 || err != nil {
		fmt.Printf("No %s files, or %s folder doesnt exist in input directory\n", folderName, folderName)
	}
	for _, file := range files {
		val, ex := filesForType[file]
		if ex && val == folderName {
			errors = append(errors, fmt.Errorf("%s: Skipping %s, as its already been run", folderName, file))
			continue
		}

		contents, err := os.ReadFile(fmt.Sprintf("%s/updates/%s", config.InputDir, file))
		if err != nil {
			errors = append(errors, fmt.Errorf("Error processing %s: %s", file, err.Error()))
			continue
		}

		tx, _ := conn.Begin()
		_, err = tx.Exec(string(contents))
		if err != nil {
			errors = append(errors, fmt.Errorf("Error processing %s: %s", file, err.Error()))
			continue
		}

		err = tx.Rollback()
		if err != nil {
			errors = append(errors, fmt.Errorf("Error processing %s: %s", file, err.Error()))
			continue
		}

		sql := fmt.Sprintf("INSERT INTO Migrations(EnterDateTime, Type, Version, FileName) VALUES ('%s', '%s', %d, '%s')", time.Now(), folderName, version, file)
		_, err = conn.Exec(sql)
		if err != nil {
			errors = append(errors, fmt.Errorf("Error processing %s: %s", file, err.Error()))
			continue
		}

		builder.WriteString(string(contents))
		fmt.Printf("%s: Processed %s\n", folderName, file)
	}

	return errors
}

func GetFilesForType(conn *sql.DB) (map[string]string, error) {
	sql := ("SELECT FileName, Type From Migrations")
	processedFiles := make(map[string]string)
	rows, err := conn.Query(sql)
	if err != nil {
		return processedFiles, err
	}

	for rows.Next() {
		var x struct {
			FileName string
			TypeName string
		}
		err = rows.Scan(&x.FileName, &x.TypeName)
		if err != nil {
			return processedFiles, err
		}
		processedFiles[x.FileName] = x.TypeName
	}

	return processedFiles, nil
}
