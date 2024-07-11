package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"os"
	"ryanclark532/migration-tool/internal/common"
	"ryanclark532/migration-tool/internal/down"
	"ryanclark532/migration-tool/internal/sqlite"
	"ryanclark532/migration-tool/internal/sqlserver"
	"ryanclark532/migration-tool/internal/up"
	"strconv"
)

const sqlserverType = "Sqlserver"

type Database interface {
	Connect() (*sql.DB, error)
	GetDatabaseState() (*common.Database, error)
	GetLatestVersion() (int, error)
	Close() error
}

type config struct {
	dbType   string
	filePath string
	name     string
	port     int
	user     string
	password string
	database string
}

func main() {

	c := loadJson()

	if c == nil {
		panic("Unable to load config, please provide a json file, config in an .env file, or cli flags, use -h for more information")
	}

	switch c.dbType {
	case sqlserverType:
		s := sqlserver.SqlServer{
			Server:   c.name,
			Port:     c.port,
			User:     c.user,
			Password: c.password,
			Database: c.database,
		}
		execute(&s)
	default:
		s := sqlite.SqLiteServer{
			FilePath: c.filePath,
		}
		execute(&s)
	}
}

func execute(server Database) {
	conn, err := server.Connect()
	if err != nil {
		panic(err)
	}

	version, err := server.GetLatestVersion()
	if err != nil {
		panic(err)
	}

	original, err := server.GetDatabaseState()
	if err != nil {
		panic(err)
	}

	up.DoMigration(conn, version)

	post, err := server.GetDatabaseState()
	if err != nil {
		panic(err)
	}

	err = down.GetDiff(original.Tables, post.Tables, version)
	if err != nil {
		panic(err)
	}

	err = server.Close()
	if err != nil {
		panic(err)
	}
}

func loadJson() *config {
	jsonFile, err := os.ReadFile("migration-settings.json")
	if err != nil {
		panic(err)
	}
	var c config
	err = json.Unmarshal(jsonFile, &c)
	if err != nil {
		panic(err)
	}
	return &c
}

func loanEnv() *config {
	p, err := strconv.ParseInt(os.Getenv("port"), 0, 8)
	if err != nil {
		panic(err)
	}
	return &config{
		dbType:   os.Getenv("dbType"),
		filePath: os.Getenv("filePath"),
		name:     os.Getenv("database"),
		port:     int(p),
		user:     os.Getenv("user"),
		password: os.Getenv("password"),
		database: os.Getenv("database"),
	}

}

func loadFlags() *config {
	c := config{}

	flag.StringVar(&c.dbType, "type", "Sqlite", "The type of the Database, e.g Sqlite, SqlServer")

	flag.StringVar(&c.filePath, "path", "./database.db", "If type if Sqlite, path to database file")

	flag.StringVar(&c.name, "server", "database", "The FQDN of the Server")

	flag.IntVar(&c.port, "port", 0, "The port the database is listening on")

	flag.StringVar(&c.user, "user", "user", "The username to authenticate against the database")

	flag.StringVar(&c.password, "password", "password", "The password to authenticate against the database")

	flag.StringVar(&c.database, "database", "database", "The name of the database of the server")

	flag.Parse()

	return &c
}
