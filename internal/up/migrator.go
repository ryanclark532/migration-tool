package up

import (
	"database/sql"
	"fmt"
	"os"
	"ryanclark532/migration-tool/internal/utils"
	"strings"
	"time"
)

func DoMigration(conn *sql.DB, version int) {
	//handle tables folder
	tables := utils.CrawlDir("./testing/tables")
	for _, file := range tables {
		contents, err := os.ReadFile(fmt.Sprintf("./testing/tables/%s", file))
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

	errors := HandleFolder(conn, &builder, "updates", filesForType, version)
	if len(errors) != 0 {
		fmt.Println("errors: ")
		for _, err := range errors {
			fmt.Println(err.Error())
		}
	}

	if builder.Len() != 0 {
		err := os.WriteFile(fmt.Sprintf("output/up/%d-up.sql", version), []byte(builder.String()), os.ModeAppend)
		if err != nil {
			panic(err)
		}
		_, err = conn.Exec(builder.String())
		if err != nil {
			panic(err)
		}
	}
	fmt.Printf("Migration Successfulll, Version: %d\n", version)

}

func HandleFolder(conn *sql.DB, builder *strings.Builder, folderName string, filesForType map[string]string, version int) []error {
	var errors []error
	for _, file := range utils.CrawlDir(fmt.Sprintf("./testing/%s", folderName)) {
		val, ex := filesForType[file]
		if ex && val == folderName {
			fmt.Printf("%s: Skipping %s, as its already been run\n", folderName, file)
			errors = append(errors, fmt.Errorf("%s: Skipping %s, as its already been run", folderName, file))
			continue
		}

		contents, err := os.ReadFile(fmt.Sprintf("./testing/updates/%s", file))
		if err != nil {
			fmt.Printf("%s: Error processing %s: %s\n", folderName, file, err.Error())
			errors = append(errors, err)
			continue
		}

		tx, _ := conn.Begin()
		_, err = tx.Exec(string(contents))
		if err != nil {
			fmt.Printf("%s: Error processing %s: %s\n", folderName, file, err.Error())
			errors = append(errors, err)
			continue
		}

		err = tx.Rollback()
		if err != nil {
			fmt.Printf("%s: Error processing %s: %s\n", folderName, file, err.Error())
			errors = append(errors, err)
			continue
		}

		sql := fmt.Sprintf("INSERT INTO Migrations(EnterDateTime, Type, Version, FileName) VALUES ('%s', '%s', %d, '%s')", time.Now(),folderName, version, file)
		_, err = conn.Exec(sql)
		if err != nil {
			fmt.Printf("%s: Error processing %s: %s\n", folderName, file, err.Error())
			errors = append(errors, err)
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

