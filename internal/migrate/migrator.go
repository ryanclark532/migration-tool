package migrate

import (
	"database/sql"
	"fmt"
	"os"
	"ryanclark532/migration-tool/internal/utils"
	"strings"
	"time"
)

type Migration struct {
	Type    string
	LastRun time.Time
	Version int
}

type Migrations struct {
	data    Migration
	jobs    Migration
	procs   Migration
	updates Migration
}

func DoMigration(conn *sql.DB) {
	//handle tables folder
	version, err := GetLatestVersion(conn)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

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
	errors := HandleFolder(conn, "updates", filesForType, version)
	if len(errors) != 0 {
		fmt.Println("errors: ")
		for _, err:= range errors{
			fmt.Println(err.Error())
		}
	}
}

func HandleFolder(conn *sql.DB, folderName string, filesForType map[string]string, version int) []error {
	var errors []error
	var builder strings.Builder
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

		sql:= fmt.Sprintf("INSERT INTO Migrations(EnterDateTime, Type, Version, FileName) VALUES (GETDATE(), '%s', %d, '%s')", folderName, version,file)
		_, err = conn.Exec(sql)
		if err != nil {
			fmt.Printf("%s: Error processing %s: %s\n", folderName, file, err.Error())
			errors = append(errors, err)
			continue
		}

		builder.WriteString(string(contents))

		fmt.Printf("%s: Processed %s\n", folderName, file)

	}

	err:= os.WriteFile("test.sql",[]byte(builder.String()), os.ModeAppend)
	if err != nil {
		errors:= []error{err}
		return errors
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

func GetLatestVersion(conn *sql.DB) (int, error) {
	sql := `SELECT MAX(Version) FROM Migrations`

	rows := conn.QueryRow(sql)

	version := 0
	err := rows.Scan(&version)
	if strings.Contains(err.Error(),"converting NULL to int is unsupported"){
		err= nil
	}
	return version, err
}
