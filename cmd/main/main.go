package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"ryanclark532/migration-tool/internal/common"
	"ryanclark532/migration-tool/internal/down"
	"ryanclark532/migration-tool/internal/sqlite"
	"ryanclark532/migration-tool/internal/up"
)

type Database interface {
	Connect() (*sql.DB, error)
	GetDatabaseState() (*common.Database, error)
	GetLatestVersion() (int, error)
	Close() error
	Setup(migrationTable string) error
}

func main() {

	c := loadJson()

	dryRun := flag.Bool("dry-run", false, "Run in dry-run mode")

	flag.Parse()

	if c == nil {
		panic("Unable to load config, please provide a json file, config in an .env file, or cli flags, use -h for more information")
	}

	switch c.DbType {
	default:
		s := sqlite.SqLiteServer{
			FilePath: c.FilePath,
		}
		execute(&s, *c, *dryRun)
	}
}

func execute(server Database, config common.Config, dryRun bool) {
	conn, err := server.Connect()
	if err != nil {
		panic(err)
	}

	err = server.Setup(config.MigrationTableName)
	if err != nil {
		panic(err)
	}

	version, err := server.GetLatestVersion()
	if err != nil {
		panic(err)
	}

	if !dryRun {
		original, err := server.GetDatabaseState()
		if err != nil {
			panic(err)
		}

		errs := up.DoMigration(conn, version, config)
		if len(errs) > 0 {
			panic(errs[0])
		}

		post, err := server.GetDatabaseState()
		if err != nil {
			panic(err)
		}

		err = down.GetTableDiff(original.Tables, post.Tables, version, config)
		if err != nil {
			panic(err)
		}

	} else {
		errs := up.DoDryMigration(conn, version, config)
		if len(errs) > 0 {
			panic(errs[0])
		}
	}

	err = server.Close()
	if err != nil {
		panic(err)
	}
}

func loadJson() *common.Config {
	jsonFile, err := os.ReadFile("migration-settings.json")
	if err != nil {
		panic(err)
	}
	var c common.Config
	err = json.Unmarshal(jsonFile, &c)
	if err != nil {
		panic(err)
	}
	fmt.Println("Loaded config from migration-settings.json")
	return &c
}

/*
func loanEnv() *config {
	p, err := strconv.ParseInt(os.Getenv("port"), 0, 8)
	if err != nil {
		panic(err)
	}
	return &config{
		DbType:   os.Getenv("dbType"),
		FilePath: os.Getenv("filePath"),
		Name:     os.Getenv("database"),
		Port:     int(p),
		user:     os.Getenv("user"),
		password: os.Getenv("password"),
		database: os.Getenv("database"),
	}

}

func loadFlags() *config {
	c := config{}

	flag.StringVar(&c.DbType, "type", "Sqlite", "The type of the Database, e.g Sqlite, SqlServer")

	flag.StringVar(&c.FilePath, "path", "./database.db", "If type if Sqlite, path to database file")

	flag.StringVar(&c.Name, "server", "database", "The FQDN of the Server")

	flag.IntVar(&c.Port, "port", 0, "The port the database is listening on")

	flag.StringVar(&c.user, "user", "user", "The username to authenticate against the database")

	flag.StringVar(&c.password, "password", "password", "The password to authenticate against the database")

	flag.StringVar(&c.database, "database", "database", "The name of the database of the server")

	return &c
}
*/
