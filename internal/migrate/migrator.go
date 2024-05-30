package migrate

import (
	"fmt"
	"os"
	"regexp"
	"ryanclark532/migration-tool/internal/sqlserver"
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

var re = regexp.MustCompile(`(?m)^\d{8}-`)

const layout = "02012006"

func DoMigration(server *sqlserver.SqlServer) {
	migrations, _ := GetLastMigrations(server)
	//handle tables folder
	tables := utils.CrawlDir("./testing/tables")
	for _, file := range tables {
		contents, err := os.ReadFile(fmt.Sprintf("./testing/tables/%s", file))
		if err != nil {
			fmt.Printf("Tables: Error processing %s: %s\n", file, err.Error())
			continue
		}
		_, err = server.Conn.Exec(string(contents))
		if err != nil {
			fmt.Printf("Tables: Error processing %s: %s\n", file, err.Error())
			continue
		}

		fmt.Println("Tables: Processed " + file)
	}

	HandleFolder("procs", server, migrations.procs)
	HandleFolder("jobs", server, migrations.jobs)
	HandleFolder("updates", server, migrations.updates)
	HandleFolder("data", server, migrations.data)

	UpdateMigrations("procs", server, migrations.procs.Version)
	UpdateMigrations("jobs", server, migrations.jobs.Version)
	UpdateMigrations("updates", server, migrations.updates.Version)
	UpdateMigrations("data", server, migrations.data.Version)
}

func HandleFolder(folderType string, server *sqlserver.SqlServer, migration Migration) {
	files := utils.CrawlDir(fmt.Sprintf("./testing/%s", folderType))
	for _, file := range files {
		dateS := re.FindString(file)
		if dateS == "" {
			continue
		}
		dateS = strings.TrimSuffix(dateS, "-")

		date, err := time.Parse(layout, dateS)
		if err != nil {
			fmt.Printf("%s: Error processing %s: %s\n", folderType, file, err.Error())
			continue
		}

		if date.Before(migration.LastRun) {
			fmt.Printf("%s: Skipping %s, since its already been run\n", folderType, file)
			continue
		}

		contents, err := os.ReadFile(fmt.Sprintf("./testing/updates/%s", file))
		if err != nil {
			fmt.Printf("%s: Error processing %s: %s\n", folderType, file, err.Error())
			continue
		}

		_, err = server.Conn.Exec(string(contents))
		if err != nil {
			fmt.Printf("%s: Error processing %s: %s\n", folderType, file, err.Error())
			continue
		}

		fmt.Printf("%s: Processed %s", folderType, file)
	}
}

func GetLastMigrations(server *sqlserver.SqlServer) (Migrations, error) {
	sql := `SELECT Type, LastRun, Version FROM Migrations ORDER BY Type`

	rows, err := server.Conn.Query(sql)
	if err != nil {
		return Migrations{}, err
	}

	var result Migrations

	rows.Next()
	err = rows.Scan(&result.data.Type, &result.data.LastRun, &result.data.Version)
	if err != nil {
		return Migrations{}, err
	}
	rows.Next()
	err = rows.Scan(&result.jobs.Type, &result.jobs.LastRun, &result.jobs.Version)
	if err != nil {
		return Migrations{}, err
	}
	rows.Next()
	err = rows.Scan(&result.procs.Type, &result.procs.LastRun, &result.procs.Version)
	if err != nil {
		return Migrations{}, err
	}
	rows.Next()
	err = rows.Scan(&result.updates.Type, &result.updates.LastRun, &result.updates.Version)
	if err != nil {
		return Migrations{}, err
	}
	return result, nil
}

func UpdateMigrations(dataType string, server *sqlserver.SqlServer, version int) {
	sql := fmt.Sprintf("UPDATE Migrations SET LastRun=GETDATE(), Version=%d WHERE Type='%s'", version+1, dataType)
	server.Conn.Exec(sql)
}
